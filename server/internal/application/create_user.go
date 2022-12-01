package application

import "github.com/google/uuid"

type CreateUserArg struct {
	Name  string
	Phone string
}

type CreateUserCommand interface {
	Handle(CreateUserArg) (uuid.UUID, error)
}

type createUserCommandImpl struct {
	repository UserRepository
}

func (c *createUserCommandImpl) Handle(arg CreateUserArg) (uuid.UUID, error) {
	return c.repository.CreateUser(arg)
}

func NewCreateUserCommand(repository UserRepository) CreateUserCommand {
	return &createUserCommandImpl{repository: repository}
}
