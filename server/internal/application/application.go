// application.go acts as a binder to the application layer exposing all the services
package application

func New(ur UserRepository, ar AccountRepository, l TxLocker) *UseCases {
	return &UseCases{
		Query: Query{
			GetUser:                      NewGetUserQuery(ur),
			GetAccountsForUser:           NewGetUserAccountsQuery(ar),
			GetTransactionSummaryForUser: NewTransactionSummaryQuery(ar),
			GetTransactionsInAccount:     NewTransactionsQuery(ar),
		},
		Command: Command{
			CreateUser:    NewCreateUserCommand(ur),
			CreateAccount: NewCreateAccountCommand(ar),
			Transfer:      NewTransferCommand(ar, l),
			Deposit:       NewDepositCommand(ar, l),
			Withdrawal:    NewWithdrawCommand(ar, l),
		},
	}
}

type UseCases struct {
	Query
	Command
}

type Query struct {
	GetUser                      GetUserQuery
	GetAccountsForUser           GetUserAccountsQuery
	GetTransactionSummaryForUser TransactionSummaryQuery
	GetTransactionsInAccount     TransactionsQuery
}

type Command struct {
	CreateUser    CreateUserCommand
	CreateAccount CreateAccountCommand
	Transfer      TransferCommand
	Deposit       DepositCommand
	Withdrawal    WithdrawCommand
}
