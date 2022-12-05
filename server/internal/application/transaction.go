package application

import (
	"context"
	"github.com/lordvidex/gomoney/pkg/gomoney"
	"github.com/pkg/errors"
)

var (
	ErrSameAccount     = errors.New("cannot transfer to the same account")
	ErrAccountNotFound = errors.New("account not found")
)

type TransferArg struct {
	FromAccountID int64
	ToAccountID   int64
	Amount        float64
}
type TransferCommand interface {
	Handle(ctx context.Context, arg TransferArg) (*gomoney.Transaction, error)
}

type transferImpl struct {
	repo AccountRepository
	l    TxLocker
}

func (t *transferImpl) Handle(ctx context.Context, arg TransferArg) (*gomoney.Transaction, error) {
	if arg.FromAccountID == arg.ToAccountID {
		return nil, ErrSameAccount
	}

	// lock the accounts to prevent race condition
	unlock := t.l.Lock(arg.FromAccountID, arg.ToAccountID)
	defer unlock()

	// get the accounts from the database
	accFrom, err := t.repo.GetAccountByID(ctx, arg.FromAccountID)
	if err != nil {
		return nil, err
	}
	accTo, err := t.repo.GetAccountByID(ctx, arg.ToAccountID)
	if err != nil {
		return nil, err
	}

	// create a new transaction
	tx := gomoney.NewTransaction(accFrom, accTo, arg.Amount, gomoney.Transfer)
	if err = tx.Validate(); err != nil {
		return nil, err
	}

	err = t.repo.Transfer(ctx, tx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save transfer")
	}
	return tx, nil
}

func NewTransferCommand(repo AccountRepository, l TxLocker) TransferCommand {
	return &transferImpl{repo, l}
}
