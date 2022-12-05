package application

import (
	"context"
	"github.com/google/uuid"
)

type CreateUserArg struct {
	Name  string
	Phone string
}

type CreateUserCommand interface {
	Handle(context.Context, CreateUserArg) (uuid.UUID, error)
}

type createUserCommandImpl struct {
	repository UserRepository
}

func (c *createUserCommandImpl) Handle(ctx context.Context, arg CreateUserArg) (uuid.UUID, error) {
	return c.repository.CreateUser(ctx, arg)
}

func NewCreateUserCommand(repository UserRepository) CreateUserCommand {
	return &createUserCommandImpl{repository: repository}
}
