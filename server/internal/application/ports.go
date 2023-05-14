// ports.go contains the application ports needed by the application layer to perform the business functions
package application

import (
	"context"

	"github.com/google/uuid"
	"github.com/lordvidex/gomoney/pkg/gomoney"
)

type UserRepository interface {
	CreateUser(ctx context.Context, arg CreateUserArg) (uuid.UUID, error)
	GetUserByPhone(ctx context.Context, phone string) (gomoney.User, error)
}

type AccountRepository interface {
	CreateAccount(ctx context.Context, arg CreateAccountArg) (int64, error)
	DeleteAccount(ctx context.Context, accountID int64) error
	// Transfer must atomically update the balance of both accounts in a transaction
	// and save the transaction itself in the storage layer
	Transfer(ctx context.Context, tx *gomoney.Transaction) error
	GetAccountsForUser(ctx context.Context, userID uuid.UUID) ([]gomoney.Account, error)
	GetAccountByID(ctx context.Context, accountID int64) (*gomoney.Account, error)

	GetLastNTransactions(ctx context.Context, accountID int64, n int) ([]gomoney.Transaction, error)
	GetAllTransactions(ctx context.Context, accountID int64) ([]gomoney.Transaction, error)
}

// TxLocker is a transaction lock interface that locks `keys`
// and returns a function that unlocks the keys
type TxLocker interface {
	Lock(x any, y ...any) func()
}
