package metadata

import (
	"context"
	"reflect"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestContextWith(t *testing.T) {
	authKey := "a123456"
	reqID := "xxx123"
	userID := "tester"
	factoryIDs := []string{"KY", "KU", "KS"}
	type args struct {
		ctx MetaContext
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "get outgoing context success with a factory id",
			args: args{
				ctx: NewOutgoingContext(authKey, reqID, userID, factoryIDs[:1]),
			},
			want: metadata.NewOutgoingContext(context.Background(), metadata.Pairs(
				AuthenticationKey, authKey,
				RequestID, reqID,
				UserID, userID,
				FactoryIDs, factoryIDs[0],
			)),
		}, {
			name: "get outgoing context success with many factory ids",
			args: args{
				ctx: NewOutgoingContext(authKey, reqID, userID, factoryIDs),
			},
			want: metadata.NewOutgoingContext(context.Background(), metadata.Pairs(
				AuthenticationKey, authKey,
				RequestID, reqID,
				UserID, userID,
				FactoryIDs, factoryIDs[0],
				FactoryIDs, factoryIDs[1],
				FactoryIDs, factoryIDs[2],
			)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContextWith(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ContextWith() = %v, want %v", got, tt.want)
			}
		})
	}
}
