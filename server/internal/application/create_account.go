package application

import (
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
	Handle(CreateAccountArg) error
}

type createAccountCommandImpl struct {
	repository Repository
}

func (c *createAccountCommandImpl) Handle(arg CreateAccountArg) error {
	return c.repository.CreateAccount(arg)
}

func NewCreateAccountCommand(repository Repository) CreateAccountCommand {
	return &createAccountCommandImpl{repository: repository}
}
