package redis

import (
	"context"
	"log"

	"github.com/go-redis/redis/v9"
	"github.com/lordvidex/gomoney/pkg/config"
)

func NewConn(ctx context.Context, c *config.Config) *redis.Client {
	opt, err := redis.ParseURL(c.Get("REDIS_URL"))
	if err != nil {
		log.Fatal(err, "failed to parse redis url")
	}
	cl := redis.NewClient(opt)
	err = cl.Ping(ctx).Err()
	if err != nil {
		log.Fatal(err, "failed to connect to cache")
	}
	return cl
}
