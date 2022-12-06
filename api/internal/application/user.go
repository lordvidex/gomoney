package application

import (
	"context"
	"github.com/google/uuid"
	"github.com/lordvidex/gomoney/api/internal/core"
	"github.com/lordvidex/gomoney/pkg/gomoney"
	"github.com/pkg/errors"
)

type CreateUserParam struct {
	Name     string
	Phone    string
	Password string
}

type CreateUserCommand interface {
	Handle(context.Context, CreateUserParam) (string, error)
}

type createUserCommandImpl struct {
	srv  Service
	repo Repository
	ph   PasswordHasher
}

func NewCreateUserCommand(srv Service, repo Repository, ph PasswordHasher) CreateUserCommand {
	return &createUserCommandImpl{srv, repo, ph}
}

func (c *createUserCommandImpl) Handle(ctx context.Context, p CreateUserParam) (string, error) {
	hashPassword, err := c.ph.CreatePasswordHash(p.Password)
	if err != nil {
		return "", errors.Wrap(err, "error hashing password")
	}

	p.Password = hashPassword
	userID, err := c.srv.CreateUser(ctx, p)
	if err != nil {
		return "", err
	}
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return "", err
	}
	apiUser := core.ApiUser{
		User: gomoney.User{
			Name:  p.Name,
			Phone: p.Phone,
			ID:    userUUID,
		},
		Password: p.Password,
	}
	err = c.repo.SaveUser(ctx, &apiUser)
	if err != nil {
		return "", errors.Wrap(err, "failed to save created user to database")
	}
	return userID, nil
}
