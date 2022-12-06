package application

func New(repo Repository, maker TokenHelper, service Service) *Usecases {
	return &Usecases{
		Query: Query{
			GetAPIUser:           NewAPIUserQuery(repo, maker),
			ViewAccounts:         NewViewAccountsQuery(service),
			ViewAccount:          NewViewAccountQuery(service),
			ViewAccountTransfers: NewViewAccountTransfersQuery(service, repo),
			ViewTransfers:        NewViewTransfersQuery(service, repo),
		},
		Command: Command{
			CreateUser:     NewCreateUserCommand(service, repo),
			Login:          NewLoginCommand(repo, maker, service),
			CreateAccount:  NewCreateAccountCommand(service),
			CreateTransfer: NewCreateTransferCommand(service, repo),
		},
	}
}

type Usecases struct {
	Query
	Command
}

type Query struct {
	GetAPIUser           APIUserQuery
	ViewAccount          ViewAccountQuery
	ViewAccounts         ViewAccountsQuery
	ViewTransfers        ViewTransfersQuery
	ViewAccountTransfers ViewAccountTransfersQuery
}

type Command struct {
	CreateUser     CreateUserCommand
	Login          LoginCommand
	CreateAccount  CreateAccountCommand
	CreateTransfer CreateTransferCommand
}
