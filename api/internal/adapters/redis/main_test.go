package redis

import (
	"github.com/go-redis/redis/v9"
	"github.com/lordvidex/gomoney/pkg/config"
	"testing"
)

var client *redis.Client

func TestMain(t *testing.M) {
	c := config.New()
	client = NewConnection(c, TestCache)
}
