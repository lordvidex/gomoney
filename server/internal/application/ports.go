// ports.go contains the application ports needed by the application layer to perform the business functions
package application

import (
	"github.com/lordvidex/gomoney/pkg/gomoney"
)

type Repository interface {
	CreateUser(CreateUserArg) error
	GetUserByPhone(phone string) (gomoney.User, error)
	CreateAccount(arg CreateAccountArg) error
}
