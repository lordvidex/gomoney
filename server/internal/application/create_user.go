package application

type CreateUserArg struct {
	Name  string
	Phone string
}

type CreateUserCommand interface {
	Handle(CreateUserArg) error
}

type createUserCommandImpl struct {
	repository Repository
}

func (c *createUserCommandImpl) Handle(arg CreateUserArg) error {
	return c.repository.CreateUser(arg)
}

func NewCreateUserCommand(repository Repository) CreateUserCommand {
	return &createUserCommandImpl{repository: repository}
}
