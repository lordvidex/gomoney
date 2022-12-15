package redis

import (
	"context"
	"log"

	"github.com/go-redis/redis/v9"
	"github.com/lordvidex/gomoney/pkg/config"
)

const (
	MainCache = iota
	TestCache
)

func NewConnection(c *config.Config, db int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     c.Get("REDIS_URL"),
		Password: c.Get("REDIS_PASSWORD"),
		DB:       db,
	})
	err := client.Ping(context.Background()).Err()
	if err != nil {
		log.Fatal(err, "failed to connect to cache")
	}
	return client
}
