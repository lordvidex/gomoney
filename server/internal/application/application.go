// application.go acts as a binder to the application layer exposing all the services
package application

func New(ur UserRepository, ar AccountRepository) *UseCases {
	return &UseCases{
		Query: Query{
			GetUser: NewGetUserQuery(ur),
		},
		Command: Command{
			CreateUser: NewCreateUserCommand(ur),
		},
	}
}

type UseCases struct {
	Query
	Command
}

type Query struct {
	GetUser GetUserQuery
}

type Command struct {
	CreateUser CreateUserCommand
}
