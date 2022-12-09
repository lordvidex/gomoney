package application

import (
	"context"
	"github.com/google/uuid"
	"github.com/lordvidex/gomoney/pkg/gomoney"
)

type Service interface {
	CreateUser(ctx context.Context, param CreateUserParam) (string, error)
	GetUserByPhone(ctx context.Context, phone string) (*gomoney.User, error)

	GetAccounts(ctx context.Context, userID string) ([]gomoney.Account, error)
	CreateAccount(ctx context.Context, userID string, account *gomoney.Account) (int64, error)
	DeleteAccount(ctx context.Context, userID string, accountID int64) error
	GetAccountTransfers(ctx context.Context, accountID int64, userID uuid.UUID) ([]gomoney.Transaction, error)
	GetTransfers(ctx context.Context, userID string) ([]gomoney.Transaction, error)
	//CreateTransfer(ctx context.Context, param CreateTransferParam) (int64, error)
}

type Cache interface {
	GetUserFromChatID(ctx context.Context, id string) (*gomoney.User, bool)
	SetUserWithChatID(ctx context.Context, id string, user gomoney.User) error
}
