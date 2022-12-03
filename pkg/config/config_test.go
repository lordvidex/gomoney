package config

import (
	"os"
	"testing"
)

func createTestEnvFile(t *testing.T) {
	f, err := os.Create("test.txt")
	defer f.Close()
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.WriteString("value_file")
	if err != nil {
		t.Fatal(err)
	}
}
func TestConfigGet(t *testing.T) {
	createTestEnvFile(t)
	defer os.Remove("test.txt")
	tests := []struct {
		name    string
		prepare func(c *Config)
		clean   func() error
		key     string
		want    string
	}{
		{
			name: "get existing key", prepare: func(c *Config) {
				c.m = map[string]string{"key": "value"}
			},
			key:  "key",
			want: "value",
		},
		{
			name: "using secret file",
			prepare: func(c *Config) {
				_ = os.Setenv("key_FILE", "test.txt")
			},
			clean: func() error {
				return os.Unsetenv("key_FILE")
			},
			key:  "key",
			want: "value_file",
		},
		{
			name: "using env variable",
			prepare: func(c *Config) {
				_ = os.Setenv("key", "value_env")
			},
			clean: func() error {
				return os.Unsetenv("key")
			},
			key:  "key",
			want: "value_env",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			tt.prepare(c)
			defer func() {
				if tt.clean != nil {
					tt.clean()
				}
			}()
			if got := c.Get(tt.key); got != tt.want {
				t.Errorf("Config.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
