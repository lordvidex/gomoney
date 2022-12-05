package postgres

import (
	"context"
	"database/sql"
	"github.com/lordvidex/gomoney/pkg/gomoney"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/lordvidex/gomoney/server/internal/adapters/postgres/sqlgen"
	app "github.com/lordvidex/gomoney/server/internal/application"
)

type accountRepo struct {
	*sqlgen.Queries
}

func NewAccount(conn *pgx.Conn) app.AccountRepository {
	return &accountRepo{
		Queries: sqlgen.New(conn),
	}
}
func (r *accountRepo) GetAccountsForUser(ctx context.Context, userID uuid.UUID) ([]*gomoney.Account, error) {
	accs, err := r.Queries.GetUserAccounts(ctx, uuid.NullUUID{UUID: userID, Valid: true})
	if err != nil {
		// TODO: handle errors
		return nil, err
	}
	return r.convertAccount(accs), nil
}

func (r *accountRepo) convertAccount(sl []*sqlgen.Account) []*gomoney.Account {
	x := make([]*gomoney.Account, len(sl))
	for i := 0; i < len(x); i++ {
		curr := sl[i]
		x[i] = &gomoney.Account{
			Id:       curr.ID,
			Title:    curr.Title,
			Currency: gomoney.Currency(curr.Currency),
		}
		if curr.Description.Valid {
			x[i].Description = curr.Description.String
		}
		if curr.IsBlocked.Valid {
			x[i].IsBlocked = curr.IsBlocked.Bool
		}
		_ = curr.Balance.AssignTo(&(x[i].Balance))
	}
	return x
}

func (r *accountRepo) CreateAccount(ctx context.Context, arg app.CreateAccountArg) (int64, error) {
	id, err := r.Queries.CreateAccount(ctx, sqlgen.CreateAccountParams{
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
