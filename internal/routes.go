package internal

import (
	"log/slog"
	"net/http"

	"github.com/Forvi/MayWalk/internal/auth/session"
	"github.com/Forvi/MayWalk/internal/common/middleware"
)

type RouteConfig struct {
	sessions   *session.SessionCache
	sessionCfg *session.SessConfig
	logger     *slog.Logger
}

func RegisterRoutes(mux *http.ServeMux, handler http.Handler, cfg *RouteConfig) {
	authMiddleware := middleware.RequireAuth(cfg.sessions, cfg.sessionCfg.Secret, cfg.logger)

	// public routes
	mux.Handle("POST /user/register", handler)
	mux.Handle("POST /user/login", handler)
	mux.Handle("POST /user/logout", handler)
	mux.Handle("GET /user/signin/yandex", handler)
	mux.Handle("GET /user/signin/yandex/callback", handler)

	// auth protected routes
	mux.Handle("GET /user/{id}", authMiddleware(handler))
	mux.Handle("PATCH /user/{id}/inactive", authMiddleware(handler))
}
