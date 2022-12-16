package application

import (
	"context"
	"log"

	"github.com/google/uuid"

	g "github.com/lordvidex/gomoney/pkg/gomoney"
)

var (
	ErrTUserNotFound = g.Err().WithCode(g.ErrNotFound).WithMessage("telegram user not found")
)

type CreateUserParam struct {
	Phone  string
	Name   string
	ChatID string
}

type TransferParam struct {
	From    int64
	To      int64
	Amount  float64
	ActorID uuid.UUID
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
	u.c.SetUserWithChatID(ctx, param.ChatID, g.User{
		ID:    uuid.MustParse(userID),
		Name:  param.Name,
		Phone: param.Phone,
	})
	return userID, nil
}

func (u *UseCases) GetUserByPhone(ctx context.Context, phone string, chatID string) (*g.User, error) {
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

func (u *UseCases) GetUserByChatID(ctx context.Context, chatID string) (*g.User, error) {
	user, ok := u.c.GetUserFromChatID(ctx, chatID)
	if !ok {
		return nil,ErrTUserNotFound
	}

	return u.srv.GetUserByPhone(ctx, user.Phone)
}

func (u *UseCases) GetAccounts(ctx context.Context, chatID string) ([]g.Account, error) {
	user, ok := u.c.GetUserFromChatID(ctx, chatID)
	if !ok {
		return nil, ErrTUserNotFound
	}
	return u.srv.GetAccounts(ctx, user.ID.String())
}

func (u *UseCases) CreateAccount(ctx context.Context, chatID string, account *g.Account) (int64, error) {
	user, ok := u.c.GetUserFromChatID(ctx, chatID)
	if !ok {
		return 0, ErrTUserNotFound
	}
	return u.srv.CreateAccount(ctx, user.ID.String(), account)
}

func (u *UseCases) DeleteAccount(ctx context.Context, accountID int64, chatID string) error {
	user, ok := u.c.GetUserFromChatID(ctx, chatID)
	if !ok {
		return ErrTUserNotFound
	}
	return u.srv.DeleteAccount(ctx, user.ID.String(), accountID)
}

func (u *UseCases) GetAccountTransfers(ctx context.Context, accountID int64, chatID string) ([]g.Transaction, error) {
	user, ok := u.c.GetUserFromChatID(ctx, chatID)
	if !ok {
		return nil, ErrTUserNotFound
	}
	return u.srv.GetAccountTransfers(ctx, accountID, user.ID)
}

func (u *UseCases) GetTransferSummary(ctx context.Context, chatID string) ([]g.TransactionSummary, error) {
	user, ok := u.c.GetUserFromChatID(ctx, chatID)
	if !ok {
		return nil, ErrTUserNotFound
	}
	return u.srv.GetTransferSummary(ctx, user.ID)
}

func (u *UseCases) Deposit(ctx context.Context, amount float64, accountID int64, chatID string) error {
	user, ok := u.c.GetUserFromChatID(ctx, chatID)
	if !ok {
		return ErrTUserNotFound
	}
	param := TransferParam{
		From:    0,
		To:      accountID,
		Amount:  amount,
		ActorID: user.ID,
	}
	return u.srv.Deposit(ctx, param)
}

func (u *UseCases) Withdraw(ctx context.Context, amount float64, accountID int64, chatID string) error {
	user, ok := u.c.GetUserFromChatID(ctx, chatID)
	if !ok {
		return ErrTUserNotFound
	}
	param := TransferParam{
		From:    accountID,
		To:      0,
		Amount:  amount,
		ActorID: user.ID,
	}
	return u.srv.Withdraw(ctx, param)
}

func (u *UseCases) Transfer(ctx context.Context, param TransferParam, chatID string) error {
	user, ok := u.c.GetUserFromChatID(ctx, chatID)
	if !ok {
		return ErrTUserNotFound
	}
	param.ActorID = user.ID
	return u.srv.Transfer(ctx, param)
}
