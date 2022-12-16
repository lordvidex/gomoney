package application

import (
	"context"
	"github.com/google/uuid"
)

type DeleteAccountArg struct {
	ActorID   uuid.UUID
	AccountID int64
}

type DeleteAccountCommand interface {
	Handle(context.Context, DeleteAccountArg) error
}

type deleteAccountImpl struct {
	repo AccountRepository
}

func (d *deleteAccountImpl) Handle(ctx context.Context, arg DeleteAccountArg) error {
	acc, err := d.repo.GetAccountByID(ctx, arg.AccountID)
	if err != nil {
		return err
	}
	if acc.OwnerID != arg.ActorID {
		return ErrOwnerAction
	}
	return d.repo.DeleteAccount(ctx, arg.AccountID)
}

func NewDeleteAccountCommand(repo AccountRepository) DeleteAccountCommand {
	return &deleteAccountImpl{repo: repo}
}
