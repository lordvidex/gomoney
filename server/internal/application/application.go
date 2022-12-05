// application.go acts as a binder to the application layer exposing all the services
package application

func New(ur UserRepository, ar AccountRepository) *UseCases {
	return &UseCases{
		Query: Query{
			GetUser:            NewGetUserQuery(ur),
			GetAccountsForUser: NewGetUserAccountsQuery(ar),
		},
		Command: Command{
			CreateUser:    NewCreateUserCommand(ur),
			CreateAccount: NewCreateAccountCommand(ar),
		},
	}
}

type UseCases struct {
	Query
	Command
}

type Query struct {
	GetUser            GetUserQuery
	GetAccountsForUser GetUserAccountsQuery
}

type Command struct {
	CreateUser    CreateUserCommand
	CreateAccount CreateAccountCommand
}
