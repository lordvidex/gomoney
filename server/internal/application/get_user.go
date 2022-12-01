package application

import "github.com/lordvidex/gomoney/pkg/gomoney"

type GetUserQueryArg struct {
	Phone string
}

type GetUserQuery interface {
	Handle(GetUserQueryArg) (gomoney.User, error)
}

type getUserQueryImpl struct {
	repository Repository
}

func (q *getUserQueryImpl) Handle(arg GetUserQueryArg) (gomoney.User, error) {
	return q.repository.GetUserByPhone(arg.Phone)
}

func NewGetUserQuery(repository Repository) GetUserQuery {
	return &getUserQueryImpl{repository: repository}
}