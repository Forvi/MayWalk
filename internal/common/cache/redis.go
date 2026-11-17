package cache

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	Redis *redis.Client
}

func NewCache(cfg *CacheConfig, logger *slog.Logger) (*Cache, error) {
	redis, err := newRedis(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("redis connection: %w", err)
	}

	logger.Info(
		"Redis successfully connected",
		"host", cfg.Host,
		"port", cfg.Port,
		"poolsize", cfg.PoolSize,
	)

	return &Cache{
		Redis: redis,
	}, nil
}

func (c *Cache) Close() error {
	if c.Redis != nil {
		return c.Redis.Close()
	}

	return nil
}

func newRedis(cfg *CacheConfig, logger *slog.Logger) (*redis.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cfg.DialTimeout)
	defer cancel()

	client := redis.NewClient(configureRedis(cfg))
	if err := client.Ping(ctx).Err(); err != nil {
		if closeErr := client.Close(); closeErr != nil {
			logger.Error("close redis after ping error", "error", closeErr)
		}

		return nil, fmt.Errorf("initialize redis: %w", err)
	}

	return client, nil
}

func configureRedis(cfg *CacheConfig) *redis.Options {
	return &redis.Options{
		Addr:         fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		MaxRetries:   cfg.MaxRetries,
	}
}
