package database

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type DBConfig struct {
	MaxConns          int32         `env:"MAX_CONNS" envDefault:"25"`
	MinConns          int32         `env:"MIN_CONNS" envDefault:"2"`
	MaxConnIdleTime   time.Duration `env:"MAX_CONN_IDLE_TIME" envDefault:"5m"`
	MaxConnLifetime   time.Duration `env:"MAX_CONNLIFETIME" envDefault:"60m"`
	HealthCheckPeriod time.Duration `env:"HEALTH_CHECK_PERIOD" envDefault:"1m"`
	MigratePath       string        `env:"MIGRATION_PATH" envDefault:"./migrations"`
	Automigrate       bool          `env:"AUTOMIGRATRE" envDefault:"true"`
	MigrateURL        string        `env:"MIGRATE_URL,required"`
	URL               string        `env:"URL,required"`
}

func LoadConfig() (*DBConfig, error) {
	cfg := &DBConfig{}
	opts := env.Options{
		Prefix: "DB_",
	}
	if err := env.ParseWithOptions(cfg, opts); err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	return cfg, nil
}
