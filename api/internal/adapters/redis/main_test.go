package redis

import (
	"os"
	"testing"

	"github.com/go-redis/redis/v9"
	"github.com/lordvidex/gomoney/pkg/config"
)

var testClient *redis.Client

func TestMain(m *testing.M) {
	c := config.New()
	testClient = NewConnection(c)

	code := m.Run()
	os.Exit(code)
}
