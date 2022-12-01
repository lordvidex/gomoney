// application.go acts as a binder to the application layer exposing all the services
package application

func New(r Repository) *UseCases {
	return &UseCases{
		Query: Query{
			GetUser: NewGetUserQuery(r),
		},
		Command: Command{
			CreateUser: NewCreateUserCommand(r),
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
