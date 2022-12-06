package application

import (
	"context"
	"github.com/lordvidex/gomoney/api/internal/core"
	"github.com/lordvidex/gomoney/pkg/gomoney"
	"time"
)

type Repository interface {
	SaveUser(context.Context, *core.ApiUser) error
	GetUserFromPhone(context.Context, string) (*core.ApiUser, error)
}

type TokenHelper interface {
	CreateToken(payload core.Payload) (string, error)
	VerifyToken(token string) (core.Payload, error)
	TokenDuration() time.Duration
}

type Service interface {
	CreateUser(ctx context.Context, param CreateUserParam) (string, error)
	GetUserByPhone(ctx context.Context, phone string) (*core.ApiUser, error)

	GetAccount(ctx context.Context, accountID int64) (gomoney.Account, error)
	GetAccounts(ctx context.Context, ID string) ([]gomoney.Account, error)
	CreateAccount(ctx context.Context, userID string, account *gomoney.Account) (int64, error)

	GetAccountTransfers(ctx context.Context, accountID int64) ([]gomoney.Transaction, error)
	GetTransfers(ctx context.Context, userID string) ([]gomoney.Transaction, error)
	CreateTransfer(ctx context.Context, param CreateTransferParam) (int64, error)
}
