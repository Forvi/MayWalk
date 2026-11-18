package config

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/Forvi/MayWalk/internal/auth/oauth"
	"github.com/Forvi/MayWalk/internal/auth/session"
	"github.com/Forvi/MayWalk/internal/common/cache"
	"github.com/Forvi/MayWalk/internal/common/database"
	"github.com/Forvi/MayWalk/internal/common/log"
	"github.com/joho/godotenv"
)

const portRegExp = `^:\d{1,5}$`

var ErrInvalidPort = errors.New("port should be in string format ':<value>'; example: :8080")

type AppConfig struct {
	Port          string
	DBConfig      database.DBConfig
	CacheConfig   cache.CacheConfig
	LoggerConfig  log.LoggerConfig
	SessionConfig session.SessConfig
	OAuthConfig   oauth.OAuthConfig
}

func LoadConfig() (*AppConfig, error) {
	_ = godotenv.Load("../.env")
	_ = godotenv.Load(".env")

	port, err := getPort()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	databaseCfg, err := database.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("load db config: %w", err)
	}

	cacheCfg, err := cache.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("load cache config: %w", err)
	}

	loggerCfg, err := log.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("load logger config: %w", err)
	}

	sessCfg, err := session.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("load session config: %w", err)
	}

	oauthCfg, err := oauth.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("load oauth config: %w", err)
	}

	return &AppConfig{
		Port:          port,
		DBConfig:      *databaseCfg,
		CacheConfig:   *cacheCfg,
		LoggerConfig:  *loggerCfg,
		SessionConfig: *sessCfg,
		OAuthConfig:   *oauthCfg,
	}, nil
}

func getPort() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return ":8080", nil
	}

	strictPortRegex := regexp.MustCompile(portRegExp)
	if !strictPortRegex.MatchString(port) {
		return "", ErrInvalidPort
	}

	return port, nil
}
