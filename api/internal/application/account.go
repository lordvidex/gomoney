package application

import (
	"context"
	"github.com/google/uuid"
	"github.com/lordvidex/gomoney/pkg/gomoney"
)

type CreateAccountParam struct {
	UserID  uuid.UUID
	Account gomoney.Account
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
	return c.srv.CreateAccount(ctx, param.UserID.String(), &param.Account)
}

type ViewAccountParam struct {
	AccountID int64
}

type ViewAccountQuery interface {
	Handle(ctx context.Context, param ViewAccountParam) (gomoney.Account, error)
}

type viewAccountImpl struct {
	srv Service
}

func NewViewAccountQuery(srv Service) ViewAccountQuery {
	return &viewAccountImpl{srv}
}

func (v *viewAccountImpl) Handle(ctx context.Context, param ViewAccountParam) (gomoney.Account, error) {
	return v.srv.GetAccount(ctx, param.AccountID)
}

type ViewAccountsParam struct {
	UserID uuid.UUID
}

type ViewAccountsQuery interface {
	Handle(ctx context.Context, param ViewAccountsParam) ([]gomoney.Account, error)
}

type viewAccountsImpl struct {
	srv Service
}

func NewViewAccountsQuery(srv Service) ViewAccountsQuery {
	return &viewAccountsImpl{srv}
}

func (v *viewAccountsImpl) Handle(ctx context.Context, param ViewAccountsParam) ([]gomoney.Account, error) {
	return v.srv.GetAccounts(ctx, param.UserID.String())
}
