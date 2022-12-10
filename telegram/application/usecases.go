package application

import (
	"context"

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
	u.c.SetUserIDWithChatID(param.ChatID, userID)
	return userID, nil
}

func (u *UseCases) GetUserByPhone(ctx context.Context, phone string) (*gomoney.User, error) {
	return u.srv.GetUserByPhone(ctx, phone)
}

func (u *UseCases) GetUserByChatID(ctx context.Context, chatID string) (*gomoney.User, error) {
	userID, ok := u.c.GetUserIDFromChatID(chatID)
	if !ok {
		return nil, errors.Wrap(gomoney.ErrNotFound, "telegram user")
	}
	return u.srv.GetUserByPhone(ctx, userID)
}

func (u *UseCases) GetAccounts(ctx context.Context, chatID string) ([]gomoney.Account, error) {
	userID, ok := u.c.GetUserIDFromChatID(chatID)
	if !ok {
		return nil, errors.Wrap(gomoney.ErrNotFound, "telegram user")
	}
	return u.srv.GetAccounts(ctx, userID)
}

func (u *UseCases) CreateAccount(ctx context.Context, chatID string, account *gomoney.Account) (int64, error) {
	userID, ok := u.c.GetUserIDFromChatID(chatID)
	if !ok {
		return 0, errors.Wrap(gomoney.ErrNotFound, "telegram user")
	}
	return u.srv.CreateAccount(ctx, userID, account)
}

func (u *UseCases) GetAccountTransfers(ctx context.Context, accountID int64, chatID string) ([]gomoney.Transaction, error) {
	userID, ok := u.c.GetUserIDFromChatID(chatID)
	if !ok {
		return nil, errors.Wrap(gomoney.ErrNotFound, "telegram user")
	}
	return u.srv.GetAccountTransfers(ctx, accountID, uuid.MustParse(userID))
}