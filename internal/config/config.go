package config

import (
	"github.com/kelseyhightower/envconfig"
)

const prefix = "APP"

type Config struct {
	Environment string         `envconfig:"ENV" default:"development"`
	HTTP        HTTPConfig     `envconfig:"HTTP"`
	Database    DatabaseConfig `envconfig:"DB"`
	Auth        AuthConfig     `envconfig:"AUTH"`
}

type HTTPConfig struct {
	Host string `envconfig:"HOST" default:"0.0.0.0"`
	Port int    `envconfig:"PORT" default:"8080"`
}

type DatabaseConfig struct {
	URL            string `envconfig:"URL" default:"postgres://postgres:postgres@localhost:5432/smartscheduling?sslmode=disable"`
	MaxOpenConns   int    `envconfig:"MAX_OPEN_CONNS" default:"10"`
	MaxIdleConns   int    `envconfig:"MAX_IDLE_CONNS" default:"5"`
	ConnMaxIdleSec int    `envconfig:"MAX_IDLE_SECONDS" default:"60"`
}

type AuthConfig struct {
	JWTSecret string `envconfig:"JWT_SECRET" default:"dev-secret"`
	JWTIssuer string `envconfig:"JWT_ISSUER" default:"smartscheduling"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process(prefix, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
