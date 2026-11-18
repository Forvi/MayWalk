package session

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type SessConfig struct {
	IdleTimeout     time.Duration `env:"IDLE_TIMEOUT" envDefault:"8h"`
	RememberTimeout time.Duration `env:"REMEMBER_TIMEOUT" envDefault:"720h"`
	Secret          string        `env:"SECRET,required"`
}

func LoadConfig() (*SessConfig, error) {
	cfg := &SessConfig{}
	opts := env.Options{
		Prefix: "SESS_",
	}
	if err := env.ParseWithOptions(cfg, opts); err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	return cfg, nil
}
