package internal

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Forvi/MayWalk/internal/auth/oauth"
	"github.com/Forvi/MayWalk/internal/auth/session"
	"github.com/Forvi/MayWalk/internal/common/cache"
	"github.com/Forvi/MayWalk/internal/common/config"
	"github.com/Forvi/MayWalk/internal/common/database"
	"github.com/Forvi/MayWalk/internal/common/log"
	"github.com/Forvi/MayWalk/internal/common/middleware"
	"github.com/Forvi/MayWalk/internal/health"
	userRepoPsql "github.com/Forvi/MayWalk/internal/user/adapters/repository/postgres"
	userHttp "github.com/Forvi/MayWalk/internal/user/adapters/transport/http"
	userUsecases "github.com/Forvi/MayWalk/internal/user/app"

	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
)

type App struct {
	listener net.Listener
	server   *http.Server
	config   *config.AppConfig
	logger   *slog.Logger
}

func InitApp() (*App, error) {
	config, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	logger := log.NewLogger(config.LoggerConfig)

	mux := http.NewServeMux()
	if err := buildApp(mux, logger); err != nil {
		logger.Error("failed to build app", "err", err)
		os.Exit(1)
	}

	server := &http.Server{
		Addr:         config.Port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &App{
		server: server,
		config: config,
		logger: logger,
	}, nil
}

func (app *App) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.start()
	}()

	select {
	case <-ctx.Done():
		return app.shutdown(context.Background())
	case err := <-errCh:
		return err
	}
}

func buildApp(mux *http.ServeMux, logger *slog.Logger) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	if cfg.DBConfig.Automigrate {
		if err := database.RunMigrate(cfg.DBConfig.MigrateURL, logger); err != nil {
			return fmt.Errorf("make migration: %w", err)
		}
	}
	db, err := database.NewDatabase(cfg.DBConfig)
	if err != nil {
		return fmt.Errorf("init database: %w", err)
	}

	cache, err := cache.NewCache(&cfg.CacheConfig, logger)
	if err != nil {
		return fmt.Errorf("init database: %w", err)
	}

	sessions := session.NewSessionCache(cache.Redis, &cfg.SessionConfig)
	auth := middleware.RequireAuth(sessions, cfg.SessionConfig.Secret)

	healthHandler := health.NewHealthHandler()
	mux.HandleFunc("/health", healthHandler.Health)

	yaOAuth := oauth.NewYaOAuth(&cfg.OAuthConfig)
	OAuthRepo := userRepoPsql.NewOAuthRepo(db)
	userRepo := userRepoPsql.NewUserRepo(db)
	userRegisterUsecase := userUsecases.NewRegisterUserUsecase(userRepo)
	userInactiveUsecase := userUsecases.NewInactiveUsecase(userRepo)
	userFindUsecase := userUsecases.NewFindUsecase(userRepo)
	userLoginUsecase := userUsecases.NewLoginUsecase(&cfg.SessionConfig, sessions, userRepo)
	userLogoutUsecase := userUsecases.NewLogoutUsecase(sessions)
	userRefreshUsecase := userUsecases.NewRefreshUsecase(&cfg.SessionConfig, sessions)
	yaSignInUsecase := userUsecases.NewYaSignInUsecase(
		yaOAuth,
		&cfg.SessionConfig,
		sessions,
		userRepo,
		OAuthRepo,
		userRegisterUsecase,
		userLoginUsecase,
	)
	userHandler := userHttp.NewUserHandler(
		logger,
		&cfg.SessionConfig,
		userRegisterUsecase,
		userInactiveUsecase,
		userFindUsecase,
		userLoginUsecase,
		userLogoutUsecase,
		userRefreshUsecase,
		yaSignInUsecase,
	)
	userHttp.Routes(mux, userHandler, auth)

	return nil
}

func (app *App) shutdown(ctx context.Context) error {
	app.logger.Info("Graceful Shutdown")
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var err error
	if app.listener != nil && app.server != nil {
		if serr := app.server.Shutdown(shutdownCtx); serr != nil {
			err = errors.Join(err, serr)
		}
	}

	return err
}

func (app *App) start() error {
	ln, err := net.Listen("tcp", app.server.Addr)
	if err != nil {
		return err
	}
	app.listener = ln
	app.logger.Info("server started", "addr", app.server.Addr)
	return app.server.Serve(ln)
}
