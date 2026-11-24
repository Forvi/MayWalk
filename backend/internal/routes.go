package internal

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Forvi/MayWalk/internal/auth/session"
	mw "github.com/Forvi/MayWalk/internal/common/middleware"
	"github.com/redis/go-redis/v9"
)

const (
	authLimit   = 5
	commonLimit = 30
)

type RouteConfig struct {
	sessions   *session.SessionCache
	sessionCfg *session.SessConfig
	logger     *slog.Logger
	cache      *redis.Client
}

func RegisterRoutes(mux *http.ServeMux, handler http.Handler, cfg *RouteConfig) {
	authMiddleware := mw.RequireAuth(cfg.sessions, cfg.sessionCfg.Secret, cfg.logger)
	authRateLimit := mw.RateLimit(cfg.cache, cfg.logger, authLimit)
	commonRateLimit := mw.RateLimit(cfg.cache, cfg.logger, commonLimit)

	stack := mw.Chain(
		mw.RequestID,
		mw.RequestLog(cfg.logger),
		mw.Recover(cfg.logger),
		mw.CORS(
			mw.AllowOrigins("http://localhost:8081", "http://swagger:8081"),
			mw.AllowMethods("GET", "POST", "PATCH", "DELETE", "OPTIONS"),
			mw.AllowHeaders("Content-Type", "Accept", "Authorization", "X-Request-ID"),
		),
		mw.SecureHeaders(),
		mw.Timeout(10*time.Second),
	)

	// public routes
	mux.Handle("POST /user/register", stack(handler))
	mux.Handle("POST /user/login", authRateLimit(stack(handler)))
	mux.Handle("POST /user/logout", commonRateLimit(stack(handler)))
	mux.Handle("GET /user/signin/yandex", authRateLimit(stack(handler)))
	mux.Handle("GET /user/signin/yandex/callback", authRateLimit(stack(handler)))

	// auth protected routes
	mux.Handle("GET /user/{id}", commonRateLimit(stack(authMiddleware(handler))))
	mux.Handle("PATCH /user/{id}/inactive", commonRateLimit(stack(authMiddleware(handler))))
}
