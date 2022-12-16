package redis

import (
	"context"
	"log"

	"github.com/go-redis/redis/v9"
	"github.com/lordvidex/gomoney/pkg/config"
)

func NewConn(ctx context.Context, c *config.Config) *redis.Client {
	cl := redis.NewClient(&redis.Options{
		Addr:     c.Get("REDIS_URL"),
		Password: c.Get("REDIS_PASSWORD"),
	})
	err := cl.Ping(ctx).Err()
	if err != nil {
		log.Fatal(err, "failed to connect to cache")
	}
	return cl
}
