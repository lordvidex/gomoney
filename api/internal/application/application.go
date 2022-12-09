package application

func New(repo UserRepository, maker TokenHelper, service Service, hasher PasswordHasher) *Usecases {
	return &Usecases{
		Query: Query{
			GetAPIUser:             NewAPIUserQuery(repo, service, maker),
			GetAccounts:            NewGetAccountsQuery(service),
			GetTransactionsSummary: NewGetTransactionsSummaryQuery(service),
			GetAccountTransactions: NewGetTransactionsQuery(service),
		},
		Command: Command{
			Login:         NewLoginCommand(repo, maker, service, hasher),
			CreateUser:    NewCreateUserCommand(repo, service, hasher),
			CreateAccount: NewCreateAccountCommand(service),
			Transfer:      NewTransferCommand(service),
			Deposit:       NewDepositCommand(service),
			Withdraw:      NewWithdrawCommand(service),
		},
	}
}

type Usecases struct {
	Query
	Command
}

type Query struct {
	GetAPIUser             APIUserQuery
	GetAccounts            GetAccountsQuery
	GetTransactionsSummary GetTransactionSummaryQuery
	GetAccountTransactions GetTransactionQuery
}

type Command struct {
	Login         LoginCommand
	CreateUser    CreateUserCommand
	CreateAccount CreateAccountCommand
	Transfer      TransferCommand
	Deposit       DepositCommand
	Withdraw      WithdrawCommand
}
