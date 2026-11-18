package main

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/Forvi/MayWalk/docs"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Port string `env:"SWAGGER_PORT,required"`
}

func main() {
	_ = godotenv.Load("../.env")
	_ = godotenv.Load(".env")

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		slog.Error("parse env", "err", err)
		return
	}

	h, err := docs.Handler()
	if err != nil {
		slog.Error("docs handler", "err", err)
		return
	}

	mux := http.NewServeMux()
	mux.Handle("/swagger/", http.StripPrefix("/swagger", h))
	mux.Handle("/openapi.yml", docs.OpenAPIHandler())

	srv := &http.Server{
		Addr:              cfg.Port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	slog.Info("docs listening", "port", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("serve", "err", err)
		return
	}
}
