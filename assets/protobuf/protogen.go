package protobuf

// NOTE: $GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.11.3 should
//       match the version in go.mod and Dockerfile. Once everyone is on non-Windows
//       systems we can use go list to get the path

//go:generate protoc -I=. --go_out=paths=source_relative:../../pkg/protobuf kenda/types/types.proto
//go:generate protoc -I=. --go_out=paths=source_relative:../../pkg/protobuf kenda/mesync/user.proto
//go:generate protoc -I=. --go_out=paths=source_relative:../../pkg/protobuf kenda/mesync/department.proto
//go:generate protoc -I=. --go_out=paths=source_relative:../../pkg/protobuf kenda/mesync/station.proto
//go:generate protoc -I=. --go_out=paths=source_relative:../../pkg/protobuf kenda/mesync/plan.proto
//go:generate protoc -I=. --go_out=paths=source_relative:../../pkg/protobuf kenda/mesync/record.proto
//go:generate protoc -I=. --go_out=paths=source_relative:../../pkg/protobuf kenda/mesync/recipe.proto
//go:generate protoc -I=. --go_out=paths=source_relative:../../pkg/protobuf kenda/mesync/blob.proto
//go:generate protoc -I=. --go_out=paths=source_relative:../../pkg/protobuf kenda/mesync/limitary_hour.proto

//go:generate protoc -I=. -I$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.11.3/third_party/googleapis --go_out=paths=source_relative,plugins=grpc:../../pkg/protobuf kenda/mesync/services.proto
//go:generate protoc -I=. -I$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.11.3/third_party/googleapis --grpc-gateway_out=paths=source_relative,logtostderr=true,allow_delete_body=true:../../pkg/protobuf kenda/mesync/services.proto
//go:generate protoc -I=. -I$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.11.3/third_party/googleapis --swagger_out=allow_delete_body=true,logtostderr=true:../../pkg/protobuf kenda/mesync/services.proto

//go:generate protoc -I=. -I$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.11.3/third_party/googleapis --go_out=paths=source_relative,plugins=grpc:../../pkg/protobuf kenda/mesync/cloud.proto
//go:generate protoc -I=. -I$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.11.3/third_party/googleapis --grpc-gateway_out=paths=source_relative,logtostderr=true,allow_delete_body=true:../../pkg/protobuf kenda/mesync/cloud.proto
//go:generate protoc -I=. -I$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.11.3/third_party/googleapis --swagger_out=allow_delete_body=true,logtostderr=true:../../pkg/protobuf kenda/mesync/cloud.proto
