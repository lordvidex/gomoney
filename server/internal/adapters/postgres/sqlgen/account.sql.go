// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: account.sql

package sqlgen

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createAccount = `-- name: CreateAccount :one
INSERT INTO "accounts" (title, description, currency, user_id) VALUES ($1, $2, $3, $4) RETURNING id
`

type CreateAccountParams struct {
	Title       string
	Description sql.NullString
	Currency    Currency
	UserID      uuid.NullUUID
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (int64, error) {
	row := q.db.QueryRow(ctx, createAccount,
		arg.Title,
		arg.Description,
		arg.Currency,
		arg.UserID,
	)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const getAccount = `-- name: GetAccount :one
SELECT id, title, description, balance, currency, is_blocked, user_id FROM "accounts" WHERE id = $1
`

func (q *Queries) GetAccount(ctx context.Context, id int64) (*Account, error) {
	row := q.db.QueryRow(ctx, getAccount, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.Balance,
		&i.Currency,
		&i.IsBlocked,
		&i.UserID,
	)
	return &i, err
}

const getUserAccounts = `-- name: GetUserAccounts :many
SELECT id, title, description, balance, currency, is_blocked, user_id FROM "accounts" WHERE user_id = $1
`

func (q *Queries) GetUserAccounts(ctx context.Context, userID uuid.NullUUID) ([]*Account, error) {
	rows, err := q.db.Query(ctx, getUserAccounts, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*Account{}
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.Balance,
			&i.Currency,
			&i.IsBlocked,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
