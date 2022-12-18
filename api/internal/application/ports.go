package application

import (
	"context"
	"github.com/lordvidex/gomoney/api/internal/core"
	"github.com/lordvidex/gomoney/pkg/gomoney"
	"time"
)

type UserRepository interface {
	SaveUser(context.Context, *core.ApiUser) error
	GetUserFromPhone(context.Context, string) (*core.ApiUser, error)
}

type TokenHelper interface {
	CreateToken(payload core.Payload) (string, error)
	VerifyToken(token string) (core.Payload, error)
	TokenDuration() time.Duration
}

type PasswordHasher interface {
	CreatePasswordHash(password string) (string, error)
	CheckPasswordHash(hashPassword, password string) error
}

type Service interface {
	CreateUser(ctx context.Context, param CreateUserParam) (string, error)
	GetUserByPhone(ctx context.Context, phone string) (*core.ApiUser, error)

	GetAccounts(ctx context.Context, ID string) ([]gomoney.Account, error)
	CreateAccount(ctx context.Context, param CreateAccountParam) (int64, error)

	Transfer(ctx context.Context, param CreateTransferParam) (*gomoney.Transaction, error)
	Deposit(ctx context.Context, param DepositParam) (*gomoney.Transaction, error)
	Withdraw(ctx context.Context, param WithdrawParam) (*gomoney.Transaction, error)

	GetTransactionSummary(ctx context.Context, userID string) ([]gomoney.TransactionSummary, error)
	GetTransactions(ctx context.Context, param UserWithAccount) (gomoney.TransactionSummary, error)
}
