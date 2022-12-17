package redis

import (
	"context"
	"log"

	"github.com/go-redis/redis/v9"
	"github.com/lordvidex/gomoney/pkg/config"
)

func NewConnection(c *config.Config) *redis.Client {
	opt, err := redis.ParseURL(c.Get("REDIS_URL"))
	if err != nil {
		log.Fatal(err, "failed to parse redis url")
	}
	client := redis.NewClient(opt)
	err = client.Ping(context.Background()).Err()
	if err != nil {
		log.Fatal(err, "failed to ping cache")
	}
	return client
}
