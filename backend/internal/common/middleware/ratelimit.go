package middleware

import (
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	l *redis_rate.Limiter
}

func NewLimiter(r *redis.Client) *RateLimiter {
	return &RateLimiter{
		l: redis_rate.NewLimiter(r),
	}
}

func RateLimit(rdb *redis.Client, logger *slog.Logger, limit int) Middleware {
	limiter := redis_rate.NewLimiter(rdb)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := "limit:" + clientIP(r)

			res, err := limiter.Allow(r.Context(), key, redis_rate.PerMinute(limit))
			if err != nil {
				logger.ErrorContext(r.Context(), "rate limit error", "error", err)
				next.ServeHTTP(w, r)
				return
			}

			h := w.Header()
			h.Set("RateLimit-Limit", strconv.Itoa(limit))
			h.Set("RateLimit-Remaining", strconv.Itoa(res.Remaining))

			if res.Allowed == 0 {
				seconds := int(res.RetryAfter / time.Second)
				h.Set("RateLimit-RetryAfter", strconv.Itoa(seconds))
				w.WriteHeader(http.StatusTooManyRequests)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"code":    "rate.limit_exceeded",
					"message": "Too many requests",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
