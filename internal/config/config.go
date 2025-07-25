package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type AppConfig struct {
	DB     DBConfig
	Server ServerConfig
}

type DBConfig struct {
	User   string `env:"DB_USER"     envDefault:"admin"`
	Passwd string `env:"DB_PASSWORD" envDefault:"admin"`
	DBName string `env:"DB_NAME"     envDefault:"postgres"`
	Host   string `env:"HOST"        envDefault:"localhost"`
	Port   string `env:"PORT"        envDefault:"5432"`
}

type ServerConfig struct {
	Port            int           `env:"SERVER_PORT"      envDefault:"8080"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"5s"`
	EnvType         string        `env:"ENV_TYPE"         envDefault:"local"`
	MigrationPath   string        `env:"MIGRATION_PATH"   envDefault:"./internal/migrations"`
}

func New() (cfg *AppConfig, err error) {
	cfgEnv := AppConfig{}
	if err := env.Parse(&cfgEnv); err != nil {
		return nil, fmt.Errorf("failed to parse env: %w", err)
	}

	return &cfgEnv, nil
}
