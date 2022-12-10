package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/lordvidex/gomoney/pkg/config"
)

func NewConn(c *config.Config) {
	redis.NewClient(&redis.Options{
		
})
}
