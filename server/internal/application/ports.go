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
	GetAccountsForUser(ctx context.Context, userID uuid.UUID) ([]gomoney.Account, error)
}
