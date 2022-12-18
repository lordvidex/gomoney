package application

import (
	"context"
	"github.com/google/uuid"
	"github.com/lordvidex/gomoney/pkg/gomoney"
)

// ------------------ Create new user account ------------------

type CreateAccountParam struct {
	UserID      uuid.UUID
	Title       string
	Description string
	Currency    gomoney.Currency
}

type CreateAccountCommand interface {
	Handle(context.Context, CreateAccountParam) (int64, error)
}

type createAccountImpl struct {
	srv Service
}

func NewCreateAccountCommand(srv Service) CreateAccountCommand {
	return &createAccountImpl{srv}
}

func (c *createAccountImpl) Handle(ctx context.Context, param CreateAccountParam) (int64, error) {
	return c.srv.CreateAccount(ctx, param)
}

// ------------------ View transfers for an account ------------------

type GetAccountsParam struct {
	UserID uuid.UUID
}

type GetAccountsQuery interface {
	Handle(ctx context.Context, param GetAccountsParam) ([]gomoney.Account, error)
}

type getAccountsImpl struct {
	srv Service
}

func NewGetAccountsQuery(srv Service) GetAccountsQuery {
	return &getAccountsImpl{srv}
}

func (v *getAccountsImpl) Handle(ctx context.Context, param GetAccountsParam) ([]gomoney.Account, error) {
	return v.srv.GetAccounts(ctx, param.UserID.String())
}
