package application

import (
	"context"
	"github.com/google/uuid"
	g "github.com/lordvidex/gomoney/pkg/gomoney"
	"github.com/pkg/errors"
)

const (
	SummariesPerAccount = 10
)

var (
	ErrSameAccount     = g.Err().WithCode(g.ErrInvalidInput).WithMessage("cannot transfer to the same account")
	ErrAccountNotFound = g.Err().WithCode(g.ErrNotFound).WithMessage("account not found")
	ErrOwnerAction     = g.Err().WithCode(g.ErrInvalidInput).WithMessage("only owner can perform this action")
	ErrAmountTooSmall  = g.Err().WithCode(g.ErrInvalidInput).WithMessage("amount must be greater than zero")
)

type TransferArg struct {
	FromAccountID int64
	ToAccountID   int64
	Amount        float64
	ActorID       uuid.UUID
}
type TransferCommand interface {
	Handle(ctx context.Context, arg TransferArg) (*g.Transaction, error)
}

type transferImpl struct {
	repo AccountRepository
	l    TxLocker
}

func (t *transferImpl) Handle(ctx context.Context, arg TransferArg) (*g.Transaction, error) {
	if arg.Amount <= 0 {
		return nil, ErrAmountTooSmall
	}
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

	// check that the actor owns the account from
	if accFrom.OwnerID != arg.ActorID {
		return nil, ErrOwnerAction
	}

	accTo, err := t.repo.GetAccountByID(ctx, arg.ToAccountID)
	if err != nil {
		return nil, err
	}

	// create a new transaction
	tx := g.NewTransaction(accFrom, accTo, arg.Amount, g.Transfer)
	if err = tx.Validate(); err != nil {
		return nil, err
	}

	err = t.repo.Transfer(ctx, tx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save transfer")
	}
	tx.From.Balance -= tx.Amount
	tx.To.Balance += tx.Amount
	return tx, nil
}

func NewTransferCommand(repo AccountRepository, l TxLocker) TransferCommand {
	return &transferImpl{repo, l}
}

type DepositArg struct {
	AccountID int64
	Amount    float64
	ActorID   uuid.UUID
}

// DepositCommand is a command that the account owner uses to add money to their accounts
type DepositCommand interface {
	Handle(ctx context.Context, arg DepositArg) (*g.Transaction, error)
}

type depositImpl struct {
	repo AccountRepository
	l    TxLocker
}

func (d *depositImpl) Handle(ctx context.Context, arg DepositArg) (*g.Transaction, error) {
	// validate the amount
	if arg.Amount <= 0 {
		return nil, ErrAmountTooSmall
	}

	// lock the account
	unlock := d.l.Lock(arg.AccountID)
	defer unlock()

	accTo, err := d.repo.GetAccountByID(ctx, arg.AccountID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account")
	}

	// check that the actor owns the account
	if accTo.OwnerID != arg.ActorID {
		return nil, ErrOwnerAction
	}

	// create a new transaction
	tx := g.NewTransaction(nil, accTo, arg.Amount, g.Deposit)
	if err = tx.Validate(); err != nil {
		return nil, err
	}
	err = d.repo.Transfer(ctx, tx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save transaction")
	}
	tx.To.Balance += tx.Amount
	return tx, nil
}

func NewDepositCommand(repo AccountRepository, l TxLocker) DepositCommand {
	return &depositImpl{repo, l}
}

type WithdrawArg struct {
	AccountID int64
	Amount    float64
	ActorID   uuid.UUID
}

type WithdrawCommand interface {
	Handle(ctx context.Context, arg WithdrawArg) (*g.Transaction, error)
}

type withdrawImpl struct {
	repo AccountRepository
	l    TxLocker
}

func (w *withdrawImpl) Handle(ctx context.Context, arg WithdrawArg) (*g.Transaction, error) {
	if arg.Amount <= 0 {
		return nil, ErrAmountTooSmall
	}
	// lock the account
	unlock := w.l.Lock(arg.AccountID)
	defer unlock()

	accFrom, err := w.repo.GetAccountByID(ctx, arg.AccountID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get account %d", arg.AccountID)
	}

	// check that the actor owns the account
	if accFrom.OwnerID != arg.ActorID {
		return nil, ErrOwnerAction
	}

	// create the new transaction
	tx := g.NewTransaction(accFrom, nil, arg.Amount, g.Withdrawal)
	if err = tx.Validate(); err != nil {
		return nil, err
	}
	err = w.repo.Transfer(ctx, tx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to save transaction")
	}
	tx.From.Balance -= tx.Amount
	return tx, nil
}

func NewWithdrawCommand(repo AccountRepository, l TxLocker) WithdrawCommand {
	return &withdrawImpl{repo, l}
}

type TransactionSummaryArg struct {
	ActorID uuid.UUID
}

type TransactionSummaryQuery interface {
	Handle(ctx context.Context, arg TransactionSummaryArg) ([]g.TransactionSummary, error)
}

type transactionSummaryImpl struct {
	repo AccountRepository
}

func (t *transactionSummaryImpl) Handle(ctx context.Context, arg TransactionSummaryArg) ([]g.TransactionSummary, error) {
	accs, err := t.repo.GetAccountsForUser(ctx, arg.ActorID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get accounts")
	}
	txs := make([]g.TransactionSummary, len(accs))

	// an alternative approach (faster but messier) is to get all transactions for the user sorted by date
	// and then filter them by account INSTEAD of repetitive repo calls.

	for i, acc := range accs {
		// get the transactions for the account
		tmp := acc
		accTxs, err := t.repo.GetLastNTransactions(ctx, acc.Id, SummariesPerAccount)
		if err != nil {
			return nil, err
		}
		txs[i] = g.TransactionSummary{
			Account:      &tmp,
			Transactions: accTxs,
		}
	}
	return txs, nil
}

func NewTransactionSummaryQuery(repo AccountRepository) TransactionSummaryQuery {
	return &transactionSummaryImpl{repo}
}

type TransactionsArg struct {
	AccountID int64
	ActorID   uuid.UUID
}

type TransactionsQuery interface {
	Handle(ctx context.Context, arg TransactionsArg) ([]g.Transaction, error)
}

type transactionsImpl struct {
	repo AccountRepository
}

func (t *transactionsImpl) Handle(ctx context.Context, arg TransactionsArg) ([]g.Transaction, error) {
	acc, err := t.repo.GetAccountByID(ctx, arg.AccountID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account")
	}
	// check that the actor owns the account
	if acc.OwnerID != arg.ActorID {
		return nil, ErrOwnerAction
	}
	return t.repo.GetAllTransactions(ctx, acc.Id)
}

func NewTransactionsQuery(repo AccountRepository) TransactionsQuery {
	return &transactionsImpl{repo}
}
