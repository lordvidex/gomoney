package application

import (
	"context"
	"github.com/lordvidex/gomoney/pkg/gomoney"
)

type CreateUserParam struct {
	Phone string
	Name  string
}
type UseCases struct {
	srv Service
	c   Cache
}

func New(srv Service) *UseCases {
	return &UseCases{
		srv: srv,
	}
}

func (u *UseCases) CreateUser(ctx context.Context, param CreateUserParam) (string, error) {
	return u.srv.CreateUser(ctx, param)
}

func (u *UseCases) GetUserByPhone(ctx context.Context, phone string) (*gomoney.User, error) {
	return u.srv.GetUserByPhone(ctx, phone)
}

func (u *UseCases) GetAccounts(ctx context.Context, userID string) ([]gomoney.Account, error) {
	return u.srv.GetAccounts(ctx, userID)
}

func (u *UseCases) CreateAccount(ctx context.Context, userID string, account *gomoney.Account) (int64, error) {
	return u.srv.CreateAccount(ctx, userID, account)
}

func (u *UseCases) GetAccountTransfers(ctx context.Context, accountID int64, chatID string) ([]gomoney.Transaction, error) {
	// if not found, get from the service

	// defer add to cache

	//return u.srv.GetAccountTransfers(ctx, accountID, chatID)
	return nil, nil
}
