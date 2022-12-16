package redis

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v9"
	"github.com/lordvidex/gomoney/api/internal/application"
	"github.com/lordvidex/gomoney/api/internal/core"
	"github.com/pkg/errors"
)

type userRepo struct {
	*redis.Client
}

func NewUserRepo(rc *redis.Client) application.UserRepository {
	return &userRepo{rc}
}

func (r *userRepo) SaveUser(ctx context.Context, user *core.ApiUser) error {
	err := r.Set(ctx, user.Phone, user, 0).Err()
	if err != nil {
		return errors.Wrap(err, "failed to save user")
	}
	return nil
}

func (r *userRepo) GetUserFromPhone(ctx context.Context, phone string) (*core.ApiUser, error) {
	// get user from db
	user, err := r.Get(ctx, phone).Result()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user")
	}

	// unmarshal user string to `core.ApiUser`
	var u *core.ApiUser
	err = json.Unmarshal([]byte(user), &u)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal api user")
	}

	return u, nil
}
