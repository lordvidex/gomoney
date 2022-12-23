package grpc

import (
	"reflect"
	"testing"

	"github.com/lordvidex/gomoney/pkg/gomoney"
	"github.com/lordvidex/gomoney/pkg/grpc"
)

func Test_mapPbAccToAcc(t *testing.T) {
	type args struct {
		pbAcc *grpc.Account
	}

	tests := []struct {
		name    string
		args    args
		want    *gomoney.Account
		wantErr bool
	}{
		{
			name: "OK",
			args: args{&grpc.Account{
				Id:          1,
				Title:       "test",
				Description: "test",
				Balance:     100,
				Currency:    grpc.Currency_USD,
				IsBlocked:   false,
			},
			},
			want: &gomoney.Account{
				Id:          1,
				Title:       "test",
				Description: "test",
				Balance:     100,
				Currency:    "USD",
				IsBlocked:   false,
			},
			wantErr: false,
		},
		{
			name: "NOT_EQUAL",
			args: args{&grpc.Account{
				Id:          1,
				Title:       "test",
				Description: "test",
				Balance:     100,
				Currency:    grpc.Currency_USD,
				IsBlocked:   false,
			},
			},
			want: &gomoney.Account{
				Id:          2,
				Title:       "test",
				Description: "test",
				Balance:     100,
				Currency:    "USD",
				IsBlocked:   false,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mapPbAccToAcc(tt.args.pbAcc); reflect.DeepEqual(got, tt.want) == tt.wantErr {
				t.Errorf("MapPbAccToAcc() = %v, want %v", got, tt.want)
			}
		})
	}
}
