package middleware

import (
	"net/http"
	"time"
)

// Timeout enforces a per-request timeout. If the downstream handler
// does not complete within d, the client receives 503 Service Unavailable.
func Timeout(d time.Duration) Middleware {
	return func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, d, http.StatusText(http.StatusServiceUnavailable))
	}
}
