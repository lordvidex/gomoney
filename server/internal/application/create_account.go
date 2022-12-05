package application

import (
	"context"
	"github.com/google/uuid"
	"github.com/lordvidex/gomoney/pkg/gomoney"
)

type CreateAccountArg struct {
	UserID      uuid.UUID
	Title       string
	Description string
	Currency    gomoney.Currency
}

type CreateAccountCommand interface {
	Handle(context.Context, CreateAccountArg) (int64, error)
}

type createAccountCommandImpl struct {
	repository AccountRepository
}

func (c *createAccountCommandImpl) Handle(ctx context.Context, arg CreateAccountArg) (int64, error) {
	return c.repository.CreateAccount(ctx, arg)
}

func NewCreateAccountCommand(repository AccountRepository) CreateAccountCommand {
	return &createAccountCommandImpl{repository: repository}
}
