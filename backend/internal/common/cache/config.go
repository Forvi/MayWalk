package cache

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type CacheConfig struct {
	Host     string `env:"HOST,required"`
	Port     string `env:"PORT,required"`
	Password string `env:"PASSWORD,required"`

	// Кол-во TCP-соединений
	PoolSize int `env:"POOLSIZE" envDefault:"10"`

	// Кол-во неиспользуемых соединений
	MinIdleConns int `env:"MIN_IDLE_CONNS" envDefault:"5"`

	// Время на установку TCP-соединения
	DialTimeout time.Duration `env:"DIAL_TIMEOUT" envDefault:"3s"`

	// Время на ожидание ответа после отправки команды
	ReadTimeout time.Duration `env:"READ_TIMEOUT" envDefault:"3s"`

	// Кол-во неиспользуемых соединений
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" envDefault:"3s"`

	// Кол-во повторений попыток
	MaxRetries int `env:"MAX_RETRIES" envDefault:"5"`
}

func LoadConfig() (*CacheConfig, error) {
	cfg := &CacheConfig{}
	opts := env.Options{
		Prefix: "REDIS_",
	}
	if err := env.ParseWithOptions(cfg, opts); err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}
	return cfg, nil
}
