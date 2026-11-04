package log

import (
	"log/slog"
	"os"
	"strings"

	"github.com/lmittmann/tint"
)

func NewLogger(cfg LoggerConfig) *slog.Logger {
	level := slog.LevelInfo

	switch strings.ToUpper(cfg.Level) {
	case "DEBUG":
		level = slog.LevelDebug
	case "INFO":
		level = slog.LevelInfo
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	}

	var handler slog.Handler
	switch cfg.Format {
	case Json:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     level,
			AddSource: false,
		})
	default:
		handler = tint.NewTextHandler(os.Stdout, &tint.Options{
			Level:     level,
			AddSource: true,
		})
	}

	logger := slog.New(handler)
	return logger
}
