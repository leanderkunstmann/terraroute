package config

import (
	"cmp"
	"errors"
	"fmt"
	"reflect"

	"github.com/leanderkunstmann/terraroute/backend/database"
	"github.com/lvlcn-t/go-kit/config"
)

const (
	appName           = "terraroute"
	defaultListenAddr = ":8080"
	defaultCORSOrigin = "http://localhost:5173"
)

type Config struct {
	Database       database.Config `json:"database" mapstructure:"database"`
	AllowedOrigins []string        `json:"allowed_origins" mapstructure:"allowed_origins"`
	ListenAddr     string          `json:"listen_addr" mapstructure:"listen_addr"`
}

func (c Config) IsEmpty() bool {
	return reflect.DeepEqual(c, Config{})
}

func (c *Config) withDefaults() *Config {
	if c == nil {
		return (&Config{}).withDefaults()
	}
	c.ListenAddr = cmp.Or(c.ListenAddr, defaultListenAddr)
	c.Database = cmp.Or(c.Database, database.Config{LocalDB: true})

	if len(c.AllowedOrigins) == 0 {
		c.AllowedOrigins = []string{defaultCORSOrigin}
	}

	return c
}

// redacted returns a copy of the [Config] with sensitive information redacted.
// This is useful for logging or displaying the config without exposing sensitive data.
func (c *Config) redacted() *Config {
	redacted := *c
	if redacted.Database.Password != "" {
		redacted.Database.Password = "[REDACTED]"
	}
	return &redacted
}

// Load loads the config from the given path.
// The config can be in every format supported
// by [github.com/spf13/viper].
func Load(path string) (*Config, error) {
	if path == "" {
		return nil, errors.New("path to config file is empty")
	}
	fmt.Printf("Loading config file: %s\n", path)

	config.SetName(appName)
	cfg, err := config.Load[Config](path)
	if err != nil && !errors.Is(err, &config.ErrConfigEmpty{}) {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	c := cfg.withDefaults()
	fmt.Printf("Used Config: %+v\n", c.redacted())
	return c, nil
}
