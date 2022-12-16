package postgres

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	g "github.com/lordvidex/gomoney/pkg/gomoney"
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
		// check violation of unique constraint
		if strings.Contains(err.Error(), "unique_phone") {
			return uuid.Nil,
				g.Err().
					WithCode(g.ErrAlreadyExists).
					WithMessage("user with phone already exists")
		}
		return uuid.Nil, err
	}
	return id, nil
}
func (r *userRepo) GetUserByPhone(ctx context.Context, phone string) (g.User, error) {
	user, err := r.Queries.GetUserByPhone(ctx, phone)
	if err != nil {
		if err == pgx.ErrNoRows {
			return g.User{}, g.Err().WithCode(g.ErrNotFound)
		}
		return g.User{}, err
	}
	return g.User{
		ID:    user.ID,
		Name:  user.Name,
		Phone: user.Phone,
	}, nil
}
