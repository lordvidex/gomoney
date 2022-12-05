package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/lordvidex/gomoney/pkg/gomoney"
	"github.com/lordvidex/gomoney/server/internal/adapters/postgres/sqlgen"
	app "github.com/lordvidex/gomoney/server/internal/application"
)

type userRepo struct {
	*sqlgen.Queries
}

func NewUser(conn *pgx.Conn) app.UserRepository {
	return &userRepo{
		Queries: sqlgen.New(conn),
	}
}
func (r *userRepo) CreateUser(ctx context.Context, arg app.CreateUserArg) (uuid.UUID, error) {
	id, err := r.Queries.CreateUser(ctx, sqlgen.CreateUserParams{Name: arg.Name, Phone: arg.Phone})
	if err != nil {
		// TODO: handle errors and map them to domain errors
		return uuid.Nil, err
	}
	return id, nil
}
func (r *userRepo) GetUserByPhone(ctx context.Context, phone string) (gomoney.User, error) {
	user, err := r.Queries.GetUserByPhone(ctx, phone)
	if err != nil {
		// TODO: handle errors and map them to domain errors
		return gomoney.User{}, err
	}
	return gomoney.User{
		ID:    user.ID,
		Name:  user.Name,
		Phone: user.Phone,
	}, nil
}
