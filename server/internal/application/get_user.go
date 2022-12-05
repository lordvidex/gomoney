package application

import (
	"context"
	"github.com/lordvidex/gomoney/pkg/gomoney"
)

type GetUserQueryArg struct {
	Phone string
}

type GetUserQuery interface {
	Handle(context.Context, GetUserQueryArg) (gomoney.User, error)
}

type getUserQueryImpl struct {
	repository UserRepository
}

func (q *getUserQueryImpl) Handle(ctx context.Context, arg GetUserQueryArg) (gomoney.User, error) {
	return q.repository.GetUserByPhone(ctx, arg.Phone)
}

func NewGetUserQuery(repository UserRepository) GetUserQuery {
	return &getUserQueryImpl{repository: repository}
}
