// ports.go contains the application ports needed by the application layer to perform the business functions
package application

import (
	"github.com/google/uuid"
	"github.com/lordvidex/gomoney/pkg/gomoney"
)

type UserRepository interface {
	CreateUser(CreateUserArg) (uuid.UUID, error)
	GetUserByPhone(phone string) (gomoney.User, error)
}

type AccountRepository interface {
	CreateAccount(arg CreateAccountArg) (int64, error)
	// TODO: Add more repo methods based on the application needs
}
