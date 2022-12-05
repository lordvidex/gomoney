package postgres

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/lordvidex/gomoney/pkg/gomoney"
	"github.com/lordvidex/gomoney/server/internal/adapters/postgres/sqlgen"
	app "github.com/lordvidex/gomoney/server/internal/application"
	"github.com/pkg/errors"
)

type accountRepo struct {
	*sqlgen.Queries
	c *pgx.Conn
}

func NewAccount(conn *pgx.Conn) app.AccountRepository {
	return &accountRepo{
		Queries: sqlgen.New(conn),
		c:       conn,
	}
}
func (r *accountRepo) GetAccountsForUser(ctx context.Context, userID uuid.UUID) ([]gomoney.Account, error) {
	accs, err := r.Queries.GetUserAccounts(ctx, uuid.NullUUID{UUID: userID, Valid: true})
	if err != nil {
		// TODO: handle errors
		return nil, err
	}
	return convertAccounts(accs), nil
}

func convertAccounts(sl []*sqlgen.Account) []gomoney.Account {
	x := make([]gomoney.Account, len(sl))
	for i := 0; i < len(x); i++ {
		curr := sl[i]
		x[i] = convertAccount(curr)
	}
	return x
}

func convertAccount(curr *sqlgen.Account) gomoney.Account {
	x := gomoney.Account{
		Id:       curr.ID,
		Title:    curr.Title,
		Currency: gomoney.Currency(curr.Currency),
	}
	if curr.Description.Valid {
		x.Description = curr.Description.String
	}
	if curr.IsBlocked.Valid {
		x.IsBlocked = curr.IsBlocked.Bool
	}
	_ = curr.Balance.AssignTo(&(x.Balance))
	return x
}

func (r *accountRepo) CreateAccount(ctx context.Context, arg app.CreateAccountArg) (int64, error) {
	id, err := r.Queries.CreateAccount(ctx, sqlgen.CreateAccountParams{
		Title:       arg.Title,
		Description: sql.NullString{String: arg.Description, Valid: arg.Description != ""},
		Currency:    sqlgen.Currency(arg.Currency),
		UserID:      uuid.NullUUID{UUID: arg.UserID, Valid: arg.UserID != uuid.Nil},
	})
	if err != nil {
		// TODO: handle errors and map them to domain errors
		return 0, err
	}
	return id, nil
}

func (r *accountRepo) Transfer(ctx context.Context, transaction *gomoney.Transaction) error {
	tx, err := r.c.Begin(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create database transaction while making transfer")
	}
	defer tx.Rollback(ctx)
	q := r.WithTx(tx)
	// withdraw from source account
	n := pgtype.Numeric{}
	_ = n.Scan(transaction.Amount)
	err = q.Withdraw(ctx, sqlgen.WithdrawParams{
		ID:     transaction.From.Id,
		Amount: n,
	})
	if err != nil {
		return errors.Wrap(err, "failed to withdraw from source account")
	}

	// deposit to destination account
	err = q.Deposit(ctx, sqlgen.DepositParams{
		ID:     transaction.To.Id,
		Amount: n,
	})
	if err != nil {
		return errors.Wrap(err, "failed to deposit to destination account")
	}

	// create transaction record
	return q.SaveTransaction(ctx, sqlgen.SaveTransactionParams{
		ID:            transaction.ID,
		CreatedAt:     transaction.Created,
		FromAccountID: mustInt64(transaction.From.Id),
		ToAccountID:   mustInt64(transaction.To.Id),
		Amount:        mustNumeric(transaction.Amount),
		Type:          mapTxType(transaction.Type),
	})
}

func (r *accountRepo) GetAccountByID(ctx context.Context, accountID int64) (*gomoney.Account, error) {
	acc, err := r.Queries.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, app.ErrAccountNotFound
		}
		return nil, err
	}
	conv := convertAccount(acc)
	return &conv, nil
}

func mustInt64(s int64) sql.NullInt64 {
	t := sql.NullInt64{}
	_ = t.Scan(s)
	return t
}

func mustNumeric(s any) pgtype.Numeric {
	t := pgtype.Numeric{}
	_ = t.Scan(s)
	return t
}

func mapTxType(t gomoney.TransactionType) sqlgen.TransactionType {
	switch t {
	case gomoney.Deposit:
		return sqlgen.TransactionTypeDeposit
	case gomoney.Withdrawal:
		return sqlgen.TransactionTypeWithdrawal
	case gomoney.Transfer:
		return sqlgen.TransactionTypeTransfer
	}
	return ""
}
