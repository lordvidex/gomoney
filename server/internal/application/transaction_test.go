package application

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	g "github.com/lordvidex/gomoney/pkg/gomoney"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransfer(t *testing.T) {
	tests := []struct {
		name    string
		arg     TransferArg
		prepare func(r *MockAccountRepository, l *MockTxLocker)
		want    *g.Transaction
		wantErr error
	}{
		{
			name: "zero amount",
			arg: TransferArg{
				FromAccountID: 1,
				ToAccountID:   2,
				Amount:        0,
			},
			wantErr: ErrAmountTooSmall,
			want:    nil,
		},
		{
			name: "negative amount",
			arg: TransferArg{
				FromAccountID: 1,
				ToAccountID:   2,
				Amount:        -1,
			},
			wantErr: ErrAmountTooSmall,
			want:    nil,
		},
		{
			name: "same amount",
			arg: TransferArg{
				FromAccountID: 1,
				ToAccountID:   1,
				Amount:        10,
			},
			wantErr: ErrSameAccount,
			want:    nil,
		},
		{
			name: "accounts not found",
			arg: TransferArg{
				FromAccountID: 10,
				ToAccountID:   20,
				Amount:        100,
			},
			prepare: func(r *MockAccountRepository, l *MockTxLocker) {
				r.EXPECT().GetAccountByID(gomock.Any(), gomock.Any()).Return(nil, ErrAccountNotFound)
				l.EXPECT().Lock(gomock.Any(), gomock.Any()).Return(func() {})
			},
			wantErr: ErrAccountNotFound,
			want:    nil,
		},
		{
			name: "account to not found",
			arg: TransferArg{
				FromAccountID: 1,
				ToAccountID:   3,
				Amount:        100,
			},
			prepare: func(r *MockAccountRepository, l *MockTxLocker) {
				r.EXPECT().GetAccountByID(gomock.Any(), int64(1)).Return(&g.Account{Id: 1}, nil).Times(1)
				r.EXPECT().GetAccountByID(gomock.Any(), gomock.Any()).Return(nil, ErrAccountNotFound).Times(1)
				l.EXPECT().Lock(gomock.Any(), gomock.Any()).Return(func() {})
			},
			wantErr: ErrAccountNotFound,
			want:    nil,
		},
		{
			name: "owner does not own account from",
			arg: TransferArg{
				FromAccountID: 1,
				ToAccountID:   2,
				Amount:        100,
				ActorID:       uuid.Nil,
			},
			prepare: func(r *MockAccountRepository, l *MockTxLocker) {
				uid := uuid.New()
				r.EXPECT().GetAccountByID(gomock.Any(), int64(1)).Return(&g.Account{Id: 1, OwnerID: uid}, nil).Times(1)
				r.EXPECT().GetAccountByID(gomock.Any(), int64(2)).Return(&g.Account{Id: 2, OwnerID: uid}, nil).MaxTimes(1)
				l.EXPECT().Lock(gomock.Any(), gomock.Any()).Return(func() {})
			},
			wantErr: ErrOwnerAction,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := NewMockAccountRepository(ctrl)
			locker := NewMockTxLocker(ctrl)
			if tt.prepare != nil {
				tt.prepare(repo, locker)
			}
			transfer := NewTransferCommand(repo, locker)
			tx, err := transfer.Handle(context.TODO(), tt.arg)
			require.Equal(t, tx, tt.want)
			require.Equal(t, err, tt.wantErr)
		})
	}
}
