package application

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/lordvidex/gomoney/pkg/gomoney"
	"testing"
)

func Test_deleteAccountImpl_Handle(t *testing.T) {
	uid := uuid.New()
	type args struct {
		ctx context.Context
		arg DeleteAccountArg
	}
	tests := []struct {
		name    string
		prepare func(repo *MockAccountRepository)
		args    args
		wantErr bool
	}{
		{
			name: "only owners can delete accounts",
			prepare: func(repo *MockAccountRepository) {
				repo.
					EXPECT().
					GetAccountByID(gomock.Any(), gomock.Any()).
					Return(&gomoney.Account{}, nil).
					Times(1)
			},
			args: args{
				ctx: context.Background(),
				arg: DeleteAccountArg{
					ActorID:   uid,
					AccountID: 1,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := NewMockAccountRepository(ctrl)
			tt.prepare(repo)
			d := &deleteAccountImpl{
				repo: repo,
			}
			if err := d.Handle(tt.args.ctx, tt.args.arg); (err != nil) != tt.wantErr {
				t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
