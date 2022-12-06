package config

import (
	"os"
)

type Config struct {
	m map[string]string
}

func New() *Config {
	return &Config{
		m: make(map[string]string),
	}
}

// Get returns the value of the key
// if the key is not found in the map, it checks if the key_FILE exists
// then it checks if the key itself is set from env
// otherwise it returns an empty string
func (c *Config) Get(key string) (result string) {
	// check if key exists in map
	if val, ok := c.m[key]; ok {
		return val
	}

	defer func() {
		if result != "" {
			c.m[key] = result
		}
	}()

	// check if file flag exists
	file := key + "_FILE"
	secret, set := os.LookupEnv(file)
	if set {
		b, err := os.ReadFile(secret)
		if err != nil {
			return ""
		}
		return string(b)
	} else {
		return os.Getenv(key)
	}
}

func (c *Config) Set(key, value string) {
	c.m[key] = value
}
