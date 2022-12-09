package application

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/lordvidex/gomoney/pkg/gomoney"
)

type CreateUserParam struct {
	Phone  string
	Name   string
	ChatID string
}

type UseCases struct {
	srv Service
	c   Cache
}

func New(srv Service, c Cache) *UseCases {
	return &UseCases{srv, c}
}

func (u *UseCases) CreateUser(ctx context.Context, param CreateUserParam) (id string, err error) {
	userID, err := u.srv.CreateUser(ctx, param)
	if err != nil {
		return "", err
	}
	u.c.SetUserWithChatID(ctx, param.ChatID, gomoney.User{
		ID:    uuid.MustParse(userID),
		Name:  param.Name,
		Phone: param.Phone,
	})
	return userID, nil
}

func (u *UseCases) GetUserByPhone(ctx context.Context, phone string, chatID string) (*gomoney.User, error) {
	user, err := u.srv.GetUserByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}
	err = u.c.SetUserWithChatID(ctx, chatID, *user)
	if err != nil {
		log.Println(err)
	}
	return user, nil
}

func (u *UseCases) GetUserByChatID(ctx context.Context, chatID string) (*gomoney.User, error) {
	user, ok := u.c.GetUserFromChatID(ctx, chatID)
	if !ok {
		return nil, errors.Wrap(gomoney.ErrNotFound, "telegram user")
	}

	return u.srv.GetUserByPhone(ctx, user.Phone)
}

func (u *UseCases) GetAccounts(ctx context.Context, chatID string) ([]gomoney.Account, error) {
	user, ok := u.c.GetUserFromChatID(ctx, chatID)
	if !ok {
		return nil, errors.Wrap(gomoney.ErrNotFound, "telegram user")
	}
	return u.srv.GetAccounts(ctx, user.ID.String())
}

func (u *UseCases) CreateAccount(ctx context.Context, chatID string, account *gomoney.Account) (int64, error) {
	user, ok := u.c.GetUserFromChatID(ctx, chatID)
	if !ok {
		return 0, errors.Wrap(gomoney.ErrNotFound, "telegram user")
	}
	return u.srv.CreateAccount(ctx, user.ID.String(), account)
}

func (u *UseCases) DeleteAccount(ctx context.Context, accountID int64, chatID string) error {
	user, ok := u.c.GetUserFromChatID(ctx, chatID)
	if !ok {
		return errors.Wrap(gomoney.ErrNotFound, "telegram user")
	}
	return u.srv.DeleteAccount(ctx, user.ID.String(), accountID)
}

func (u *UseCases) GetAccountTransfers(ctx context.Context, accountID int64, chatID string) ([]gomoney.Transaction, error) {
	user, ok := u.c.GetUserFromChatID(ctx, chatID)
	if !ok {
		return nil, errors.Wrap(gomoney.ErrNotFound, "telegram user")
	}
	return u.srv.GetAccountTransfers(ctx, accountID, user.ID)
}
