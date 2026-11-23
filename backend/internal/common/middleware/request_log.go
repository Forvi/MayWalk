package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type loggingRec struct {
	http.ResponseWriter
	status int
}

func newLoggingRec(w http.ResponseWriter) *loggingRec {
	return &loggingRec{ResponseWriter: w, status: http.StatusOK}
}

func RequestLog(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := newLoggingRec(w)
			next.ServeHTTP(rec, r)
			logger.InfoContext(r.Context(), "request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rec.status,
				"duration_ms", time.Since(start).Milliseconds(),
			)
		})
	}
}
