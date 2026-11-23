package middleware

import (
	"net/http"

	ctxUtil "github.com/Forvi/MayWalk/internal/common/ctx"
	"github.com/google/uuid"
)

const requestIDHeaderKey = "X-Request-ID"

// RequestID generates or propagates an X-Request-ID header.
// If the incoming request has X-Request-ID, it's preserved. Otherwise a new one is generated.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestID := r.Header.Get(requestIDHeaderKey)
		if requestID == "" {
			requestID = uuid.NewString()
		}

		w.Header().Set(requestIDHeaderKey, requestID)
		ctx = ctxUtil.SetRequestID(ctx, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
