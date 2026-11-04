package config

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/Forvi/MayWalk/internal/common/database"
	"github.com/Forvi/MayWalk/internal/common/log"
	"github.com/joho/godotenv"
)

const portRegExp = `^:\d{1,5}$`

var ErrInvalidPort = errors.New("port should be in string format ':<value>'; example: :8080")

type AppConfig struct {
	Port         string
	DBConfig     database.DBConfig
	LoggerConfig log.LoggerConfig
}

func LoadConfig() (*AppConfig, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println(".env not found")
	}

	port, err := getPort()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	databaseCfg, err := database.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("load db config: %w", err)
	}

	loggerCfg, err := log.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("load logger config: %w", err)
	}

	return &AppConfig{
		Port:         port,
		DBConfig:     *databaseCfg,
		LoggerConfig: *loggerCfg,
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
