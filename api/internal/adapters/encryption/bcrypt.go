package encryption

import (
	"github.com/lordvidex/gomoney/api/internal/application"
	"golang.org/x/crypto/bcrypt"
)

type bcryptPasswordHasherImpl struct{}

func NewBcryptPasswordHasher() application.PasswordHasher {
	return &bcryptPasswordHasherImpl{}
}

func (h *bcryptPasswordHasherImpl) CreatePasswordHash(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashPassword), nil
}

func (h *bcryptPasswordHasherImpl) CheckPasswordHash(hashPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
}
