package redis

import (
	"context"
	"github.com/go-redis/redis/v9"
	"github.com/lordvidex/gomoney/pkg/config"
	"log"
)

func NewConnection(c *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     c.Get("REDIS_URL"),
		Password: c.Get("REDIS_PASSWORD"),
		DB:       0,
	})
	err := client.Ping(context.Background()).Err()
	if err != nil {
		log.Fatal(err, "failed to connect to cache")
	}
	return client
}
