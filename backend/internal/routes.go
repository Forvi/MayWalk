package internal

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Forvi/MayWalk/internal/auth/session"
	mw "github.com/Forvi/MayWalk/internal/common/middleware"
)

type RouteConfig struct {
	sessions   *session.SessionCache
	sessionCfg *session.SessConfig
	logger     *slog.Logger
}

func RegisterRoutes(mux *http.ServeMux, handler http.Handler, cfg *RouteConfig) {
	authMiddleware := mw.RequireAuth(cfg.sessions, cfg.sessionCfg.Secret, cfg.logger)

	stack := mw.Chain(
		mw.RequestID,
		mw.RequestLog(cfg.logger),
		mw.Recover(cfg.logger),
		mw.CORS(
			mw.AllowOrigins("http://localhost:8081", "http://backend:8081"),
			mw.AllowMethods("GET", "POST", "PATCH", "DELETE", "OPTIONS"),
			mw.AllowHeaders("Content-Type", "Accept", "Authorization", "X-Request-ID"),
		),
		mw.SecureHeaders(),
		mw.Timeout(10*time.Second),
	)

	// public routes
	mux.Handle("POST /user/register", stack(handler))
	mux.Handle("POST /user/login", stack(handler))
	mux.Handle("POST /user/logout", stack(handler))
	mux.Handle("GET /user/signin/yandex", stack(handler))
	mux.Handle("GET /user/signin/yandex/callback", stack(handler))

	// auth protected routes
	mux.Handle("GET /user/{id}", stack(authMiddleware(handler)))
	mux.Handle("PATCH /user/{id}/inactive", stack(authMiddleware(handler)))
}
