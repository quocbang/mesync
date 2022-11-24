package server

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rs/xid"
	"go.uber.org/zap"

	commonsCtx "gitlab.kenda.com.tw/kenda/commons/v2/utils/context"
	"gitlab.kenda.com.tw/kenda/mcom"
	"gitlab.kenda.com.tw/kenda/mcom/utils/types"

	"gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/mesync"
	grpc_context "gitlab.kenda.com.tw/kenda/mesync/pkg/protocol/grpc/context"
)

// UploadBlob implements gitlab.kenda.com.tw/kenda/mesync/pkg/protobuf/kenda/mesync MesyncServer interface
func (s CloudServer) UploadBlob(stream mesync.Cloud_UploadBlobServer) error {
	ch, done := s.supplier(stream)

	for i := 0; i < s.cloud.configs.UploadRetryParameters.MaxConcurrentJobs; i++ {
		go consumer(ch, func(data *mesync.UploadBlobRequest) {
			if data == nil {
				return
			}

			mcomReq, err := s.uploadBlob(stream.Context(), data)
			if err != nil {
				s.getLogger(stream.Context()).Error("fail to upload station parameters to azure blob", zap.Error(err))
				return
			}

			if err := s.cloud.dm.CreateBlobResourceRecord(commonsCtx.WithLogger(context.Background(), s.getLogger(stream.Context())), mcomReq); err != nil {
				s.getLogger(stream.Context()).Error("fail to save record to cloud database", zap.Error(err))
			}
		})
	}

	if err := <-done; err != nil {
		return err
	}
	return nil
}

func consumer(ch <-chan *mesync.UploadBlobRequest, delegate func(data *mesync.UploadBlobRequest)) {
	for data := range ch {
		delegate(data)
	}
}

func (s CloudServer) supplier(stream mesync.Cloud_UploadBlobServer) (ch chan *mesync.UploadBlobRequest, done chan error) {
	ch, done = make(chan *mesync.UploadBlobRequest, 5), make(chan error)
	go func() {
		defer close(ch)
		defer close(done)

		for {
			data, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) {
					stream.SendAndClose(&empty.Empty{})
					break
				}
				s.getLogger(stream.Context()).Error("stream err", zap.Error(err))
				done <- err
				break
			}
			ch <- data
		}
	}()
	return ch, done
}

func (s CloudServer) uploadBlob(ctx context.Context, req *mesync.UploadBlobRequest) (mcomReq mcom.CreateBlobResourceRecordRequest, err error) {
	files := parseUploadBlobReq(req)
	rep, err := proto.Marshal(&files)
	if err != nil {
		return mcom.CreateBlobResourceRecordRequest{}, fmt.Errorf("cannot marshal proto message to binary: %w", err)
	}

	blobName, t, err := s.getBlobNameAndTime(files, req.StationId)
	if err != nil {
		return mcom.CreateBlobResourceRecordRequest{}, err
	}

	containerName, err := s.getContainerName(ctx)
	if err != nil {
		return mcom.CreateBlobResourceRecordRequest{}, err
	}

	token := s.cloud.configs.ContainerToken[containerName]

	rub := reUploadBlob{
		ContainerName:   containerName,
		BlobAccountName: token.AccountName,
		BlobSasToken:    token.SasToken,
		BlobName:        blobName,
		Contents:        rep,
		Resources:       req.Resources,
		StationId:       req.StationId,
		DateTime:        types.ToTimeNano(t),
	}

	s.doUpload(ctx, rub)

	return rub.ToMcomRequest(), nil
}

func (s CloudServer) getBlobNameAndTime(files mesync.BlobStoreFile, stationID string) (string, time.Time, error) {
	if len(files.Parameters.Details) == 0 {
		return "", time.Time{}, fmt.Errorf("empty details detected")
	}

	lastFile := files.Parameters.Details[len(files.Parameters.Details)-1]
	t, err := ptypes.Timestamp(lastFile.GetDateTime())
	if err != nil {
		return "", time.Time{}, err
	}

	stationID = sha256String(stationID)

	blobName := t.Format("2006/01/02/") + stationID + t.Format("/150405")
	return blobName, t, nil
}

func sha256String(stationID string) string {
	h := sha256.New()
	h.Write([]byte(stationID))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (s CloudServer) getContainerName(ctx context.Context) (string, error) {
	factoriesID := grpc_context.GetFactoryIDs(ctx)
	if len(factoriesID) == 0 {
		return "", fmt.Errorf("missing factory id")
	} else if len(factoriesID) > 1 {
		logger := s.getLogger(ctx)
		logger.Warn("find multiple factories in the context: ", zap.Strings("factory ids", factoriesID))
	}

	containerName, ok := s.cloud.configs.FactoryContainerMap[factoriesID[0]]
	if !ok {
		return "", fmt.Errorf("undefined factoryID-container mapping")
	}
	return containerName, nil
}

func (s CloudServer) StartReUploadTicker() {
	ticker := time.NewTicker(time.Duration(s.cloud.configs.UploadRetryParameters.Interval) * time.Minute)
	go func() {
		for range ticker.C {
			logger := zap.L().With(zap.String("request_id", xid.New().String()))
			s.reUploadBlob(commonsCtx.WithLogger(context.Background(), logger))
		}
	}()
}

func (s CloudServer) reUploadBlob(ctx context.Context) {
	dirPath := s.cloud.configs.UploadRetryParameters.UnsuccessStoragePath
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return
	}

	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		commonsCtx.Logger(ctx).Error("fail to read cloud unsuccess storage path.", zap.String("path", dirPath), zap.Error(err))
	}
	if len(files) == 0 {
		return
	}

	commonsCtx.Logger(ctx).Info("retry uploading station parameters to azure blob")
	if err != nil {
		commonsCtx.Logger(ctx).Error("retry uploading station parameters to azure blob failed", zap.Error(err))
		return
	}

	limitChannel := make(chan bool, s.cloud.configs.UploadRetryParameters.MaxConcurrentJobs)
	for _, file := range files {
		limitChannel <- true
		go func(file fs.FileInfo) {
			defer func() { <-limitChannel }()
			fileName := dirPath + "/" + file.Name()
			f, err := os.Open(fileName)
			if err != nil {
				commonsCtx.Logger(ctx).Error("retry uploading station parameters to azure blob failed", zap.Error(err))
				return
			}
			defer f.Close()

			var rub reUploadBlob
			decoder := gob.NewDecoder(f)
			if err := decoder.Decode(&rub); err != nil {
				commonsCtx.Logger(ctx).Error("retry uploading station parameters to azure blob failed", zap.Error(err))
				return
			}

			s.doUpload(ctx, rub)

			if err := s.cloud.dm.CreateBlobResourceRecord(ctx, rub.ToMcomRequest()); err != nil {
				commonsCtx.Logger(ctx).Error("retry uploading station parameters to azure blob failed", zap.Error(err))
				return
			}

			if err := os.Remove(fileName); err != nil {
				commonsCtx.Logger(ctx).Error("reUpload failed: ", zap.String("file name", fileName), zap.Error(err))
				return
			}

			commonsCtx.Logger(ctx).Info("retry uploading station parameters to azure blob completed", zap.String("file name", fileName))
		}(file)
	}
	close(limitChannel)
}

func (s CloudServer) doUpload(ctx context.Context, rub reUploadBlob) {
	var err error
	retryTimes := 1
	loopCondition := func() bool { return err != nil && retryTimes <= s.cloud.configs.UploadRetryParameters.MaxTimes }
	for err = sendHttpRequest(rub); loopCondition(); retryTimes++ {
		commonsCtx.Logger(ctx).Warn("failed to upload to cloud", zap.Int("retry_time", retryTimes), zap.String("retry_progress", fmt.Sprintf("%02f%%", float32(retryTimes)/float32(s.cloud.configs.UploadRetryParameters.MaxTimes)*100)), zap.String("blob", rub.BlobName), zap.Error(err))
		err = sendHttpRequest(rub)
	}

	if err != nil {
		s.uploadFailHandle(ctx, rub)
	}
}

func (s CloudServer) uploadFailHandle(ctx context.Context, rub reUploadBlob) {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(rub)
	if err != nil {
		commonsCtx.Logger(ctx).Error("gob encoding fail", zap.Error(err))
		return
	}

	dirPath := s.cloud.configs.UploadRetryParameters.UnsuccessStoragePath
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		commonsCtx.Logger(ctx).Error("fail to create directory", zap.String("directory path", dirPath), zap.Error(err))
		return
	}

	fPath := filepath.Join(dirPath, xid.New().String())
	err = ioutil.WriteFile(fPath, buf.Bytes(), fs.ModePerm)
	if err != nil {
		commonsCtx.Logger(ctx).Error("fail to write file", zap.String("file path", fPath), zap.Error(err))
		return
	}
}

type reUploadBlob struct {
	ContainerName   string
	BlobAccountName string
	BlobSasToken    string
	BlobName        string
	Contents        []byte
	Resources       []string
	StationId       string
	DateTime        types.TimeNano
}

func (rub reUploadBlob) ToMcomRequest() mcom.CreateBlobResourceRecordRequest {
	return mcom.CreateBlobResourceRecordRequest{
		Details: []mcom.CreateBlobResourceRecordDetail{{
			BlobURI:       rub.BlobName,
			Resources:     rub.Resources,
			Station:       rub.StationId,
			DateTime:      rub.DateTime,
			ContainerName: rub.ContainerName,
		}},
	}
}

func parseUploadBlobReq(req *mesync.UploadBlobRequest) mesync.BlobStoreFile {
	res := make([]*mesync.BlobUnitFile, len(req.Detail))
	for i, d := range req.Detail {
		res[i] = &mesync.BlobUnitFile{
			DateTime:              d.GetDateTime(),
			ManufactureParameters: d.GetManufactureParameters(),
		}
	}
	return mesync.BlobStoreFile{
		Batch:     req.GetBatch(),
		Resources: req.GetResources(),
		StationId: req.GetStationId(),
		Parameters: &mesync.StationParameters{
			ParametersHeader: req.GetParametersHeader(),
			Details:          res,
		},
	}
}

func sendHttpRequest(rub reUploadBlob) error {
	url := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s?%s", rub.BlobAccountName, rub.ContainerName, url.QueryEscape(rub.BlobName), rub.BlobSasToken)
	request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(rub.Contents))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "text/plain")
	request.Header.Set("x-ms-blob-type", "BlockBlob")
	request.Header.Set("Accept", "*/*")
	request.Header.Set("Accept-Encoding", "gzip, deflate, br")
	request.Header.Set("Connection", "keep-alive")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.Status[0] != '2' {
		var resBody []byte
		response.Body.Read(resBody)
		return fmt.Errorf("unsuccess response: %s, body: %s", response.Status, string(resBody))
	}
	return nil
}
