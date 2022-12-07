package grpc

import (
	"github.com/google/uuid"
	"github.com/lordvidex/gomoney/pkg/gomoney"
	pb "github.com/lordvidex/gomoney/pkg/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"reflect"
	"testing"
	"time"
)

func Test_mapTxToPb(t *testing.T) {
	testUUID := uuid.New()
	testTime := time.Now()
	tests := []struct {
		name string
		t    *gomoney.Transaction
		want *pb.Transaction
	}{
		{name: "empty", t: nil, want: nil},
		{name: "withdraw", t: &gomoney.Transaction{
			ID:      testUUID,
			Amount:  100,
			From:    nil,
			To:      nil,
			Created: testTime,
			Type:    gomoney.Withdrawal,
		},
			want: &pb.Transaction{
				Id:        &pb.StringID{Id: testUUID.String()},
				Amount:    100,
				From:      nil,
				To:        nil,
				CreatedAt: timestamppb.New(testTime),
				Type:      pb.TransactionType_WITHDRAWAL,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mapTxToPb(tt.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mapTxToPb() = %v, want %v", got, tt.want)
			}
		})
	}
}
