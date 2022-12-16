package application

import (
	"context"
	"github.com/golang/mock/gomock"
	"testing"
)

func Test_createAccountCommandImpl_Handle(t *testing.T) {
	type args struct {
		ctx context.Context
		arg CreateAccountArg
	}
	tests := []struct {
		name    string
		prepare func(repo *MockAccountRepository)
		args    args
		want    int64
		wantErr bool
		verify  func(repo *MockAccountRepository, t *testing.T)
	}{
		{
			name: "repository must be called",
			prepare: func(repo *MockAccountRepository) {
				repo.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).Return(int64(1), nil)
			},
			want:    int64(1),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := NewMockAccountRepository(ctrl)
			c := &createAccountCommandImpl{
				repository: repo,
			}
			if tt.prepare != nil {
				tt.prepare(repo)
			}
			got, err := c.Handle(tt.args.ctx, tt.args.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Handle() got = %v, want %v", got, tt.want)
			}
			if tt.verify != nil {
				tt.verify(repo, t)
			}
		})
	}
}
