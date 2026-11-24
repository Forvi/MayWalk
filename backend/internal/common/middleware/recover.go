package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	commonErrors "github.com/Forvi/MayWalk/internal/common/errors"
)

// Recover returns middleware that recovers from panics in downstream handlers.
// When a panic occurs, it logs the panic value and stack trace to logs and
// responds with a 500 Internal Server Error.
func Recover(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.ErrorContext(r.Context(), "panic",
						"err", err,
						"stack", string(debug.Stack()),
					)
					commonErrors.WriteError(
						w, r,
						commonErrors.NewDomainError("http.internal_error", "internal server error", http.StatusInternalServerError),
						logger,
					)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
