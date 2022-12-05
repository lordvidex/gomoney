package application

import (
	"context"
	"github.com/google/uuid"
	"github.com/lordvidex/gomoney/pkg/gomoney"
)

type GetUserAccountsParam struct {
	UserID uuid.UUID
}

type GetUserAccountsQuery interface {
	Handle(ctx context.Context, param GetUserAccountsParam) ([]*gomoney.Account, error)
}

type getUserAccountsImpl struct {
	repo AccountRepository
}

func NewGetUserAccountsQuery(repo AccountRepository) GetUserAccountsQuery {
	return &getUserAccountsImpl{repo: repo}
}

func (g *getUserAccountsImpl) Handle(ctx context.Context, param GetUserAccountsParam) ([]*gomoney.Account, error) {
	return g.repo.GetAccountsForUser(ctx, param.UserID)
}
