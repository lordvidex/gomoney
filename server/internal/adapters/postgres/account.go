package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/lordvidex/gomoney/pkg/gomoney"
	"github.com/lordvidex/gomoney/server/internal/adapters/postgres/sqlgen"
	app "github.com/lordvidex/gomoney/server/internal/application"
)

type repository struct {
	*sqlgen.Queries
	c context.Context
}

func NewRepository(conn *pgx.Conn) *repository {
	return &repository{
		c: context.TODO(),
	}
}
func (r *repository) CreateUser(arg app.CreateUserArg) error {
	_, err := r.Queries.CreateUser(r.c, sqlgen.CreateUserParams{Name: arg.Name, Phone: arg.Phone})
	if err != nil {
		return err
	}
	return nil
}
func (r *repository) GetUserByPhone(phone string) (gomoney.User, error) {
	user, err := r.Queries.GetUserByPhone(r.c, phone)
	if err != nil {
		return gomoney.User{}, err
	}
	return gomoney.User{
		ID:    user.ID,
		Name:  user.Name,
		Phone: user.Phone,
	}, nil
}

func (r *repository) CreateAccount(arg app.CreateAccountArg) error {
	_, err := r.Queries.CreateAccount(r.c, sqlgen.CreateAccountParams{
		Title:       arg.Title,
		Description: sql.NullString{String: arg.Description, Valid: arg.Description == ""},
		Currency:    sqlgen.Currency(arg.Currency),
		UserID:      uuid.NullUUID{UUID: arg.UserID, Valid: arg.UserID != uuid.Nil},
	})
	if err != nil {
		return err
	}
	return nil
}
