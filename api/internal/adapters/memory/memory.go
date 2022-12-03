package memory

import (
	"context"
	"github.com/lordvidex/gomoney/api/internal/application"
	"github.com/lordvidex/gomoney/api/internal/core"
	"github.com/pkg/errors"
)

type repo struct {
	m map[string]*core.ApiUser
}

func (r *repo) SaveUser(ctx context.Context, user *core.ApiUser) error {
	r.m[user.Phone] = user
	return nil
}

func (r *repo) GetUserFromPhone(ctx context.Context, s string) (*core.ApiUser, error) {
	user, exists := r.m[s]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func New() application.Repository {
	return &repo{
		m: make(map[string]*core.ApiUser),
	}
}
