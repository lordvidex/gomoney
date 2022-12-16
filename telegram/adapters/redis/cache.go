package redis

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v9"
	"github.com/lordvidex/gomoney/pkg/gomoney"
	app "github.com/lordvidex/gomoney/telegram/application"
)



type redisCache struct {
	*redis.Client
}

func NewCache(rc *redis.Client) app.Cache {
	return &redisCache{rc}
}

func (r *redisCache) GetUserFromChatID(ctx context.Context, id string) (*gomoney.User, bool) {
	uid, err := r.Get(ctx, id).Result()
	if err != nil {
		return nil, false
	}
	u := &gomoney.User{}
	s := []byte(uid)
	err = json.Unmarshal(s, u)
	return u, true
}
func (r *redisCache) SetUserWithChatID(ctx context.Context, id string, user gomoney.User) error {
	return r.Set(ctx, id, user, 0).Err()
}
