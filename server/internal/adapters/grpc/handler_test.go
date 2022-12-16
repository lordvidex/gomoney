package grpc

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lordvidex/gomoney/pkg/gomoney"
	pb "github.com/lordvidex/gomoney/pkg/grpc"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_mapAcctToPb(t *testing.T) {
	tests := []struct {
		name string
		a    *gomoney.Account
		want *pb.Account
	}{
		{name: "empty", a: nil, want: nil},
		{name: "full",
			a: &gomoney.Account{
				Id:          2,
				Title:       "test",
				Description: "test description",
				Balance:     100,
				Currency:    gomoney.USD,
				IsBlocked:   false,
			},
			want: &pb.Account{
				Id:          2,
				Title:       "test",
				Description: "test description",
				Balance:     100,
				Currency:    pb.Currency_USD,
				IsBlocked:   false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapAcctToPb(tt.a)
			assert.Equal(t, tt.want, got)
		})
	}
}

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
			got := mapTxToPb(tt.t)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMapTxTypeToPb(t *testing.T) {
	tests := []struct {
		from gomoney.TransactionType
		want pb.TransactionType
	}{
		{gomoney.Deposit, pb.TransactionType_DEPOSIT},
		{gomoney.Transfer, pb.TransactionType_TRANSFER},
		{gomoney.Withdrawal, pb.TransactionType_WITHDRAWAL},
		{gomoney.TransactionType(100), -1},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("from %v to %v", tt.from, tt.want), func(t *testing.T) {
			if got := mapTxTypeToPb(tt.from); got != tt.want {
				t.Errorf("mapTxTypeToPb() = %v, want %v", got, tt.want)
			}
		})
	}
}
