package log

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type LogFormat string

const (
	Json  LogFormat = "JSON"
	Plain LogFormat = "PLAIN"
)

type LoggerConfig struct {
	/** log output level; can be: INFO, DEBUG, WARN */
	Level string `env:"LEVEL" envDefault:"INFO"`

	/** log output format; can be: JSON, PLAIN */
	Format LogFormat `env:"FORMAT" envDefault:"JSON"`
}

func LoadConfig() (*LoggerConfig, error) {
	cfg := &LoggerConfig{}
	opts := env.Options{
		Prefix: "LOG_",
	}
	if err := env.ParseWithOptions(cfg, opts); err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	return cfg, nil
}
