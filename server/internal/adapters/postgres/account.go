package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/lordvidex/gomoney/server/internal/adapters/postgres/sqlgen"
	app "github.com/lordvidex/gomoney/server/internal/application"
)

type accountRepo struct {
	*sqlgen.Queries
	c context.Context	
}

func NewAccount(conn *pgx.Conn) *accountRepo {
	return &accountRepo{
		c:       context.TODO(),
		Queries: sqlgen.New(conn),
	}
}

func (r *accountRepo) CreateAccount(arg app.CreateAccountArg) (int64, error) {
	id, err := r.Queries.CreateAccount(r.c, sqlgen.CreateAccountParams{
		Title:       arg.Title,
		Description: sql.NullString{String: arg.Description, Valid: arg.Description == ""},
		Currency:    sqlgen.Currency(arg.Currency),
		UserID:      uuid.NullUUID{UUID: arg.UserID, Valid: arg.UserID != uuid.Nil},
	})
	if err != nil {
		// TODO: handle errors and map them to domain errors
		return 0, err
	}
	return id, nil
}
