package application

func New(repo Repository, maker TokenHelper, service Service) *Usecases {
	return &Usecases{
		Query: Query{
			ViewAccounts: NewViewAccountsQuery(service),
			GetAPIUser:   NewAPIUserQuery(repo, maker),
			//ViewTransactions: nil,
		},
		Command: Command{
			CreateUser:    NewCreateUserCommand(service, repo),
			Login:         NewLoginCommand(repo, maker, service),
			CreateAccount: NewCreateAccountCommand(service),
		},
	}
}

type Usecases struct {
	Query
	Command
}

type Query struct {
	ViewAccounts ViewAccountsQuery
	GetAPIUser   APIUserQuery
	//ViewTransactions ViewTransactionsQuery
}

type Command struct {
	CreateUser    CreateUserCommand
	Login         LoginCommand
	CreateAccount CreateAccountCommand
}
