package errors

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	ctxUtil "github.com/Forvi/MayWalk/internal/common/ctx"
)

// ErrorMiddleware handles errors from transport-handlers
func ErrorMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}

// WriteError writes error response
func WriteError(w http.ResponseWriter, r *http.Request, err error, logger *slog.Logger) {
	if domainErr, ok := errors.AsType[*DomainError](err); ok {
		reqID := ctxUtil.RequestID(r.Context())
		logger.WarnContext(r.Context(), "domain error", "code", domainErr.Code, "status", domainErr.Status)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(domainErr.Status)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"code":       domainErr.Code,
			"message":    domainErr.Message,
			"request_id": reqID,
		})
		return
	}

	logger.ErrorContext(r.Context(), "unhandled error", "error", err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"code":    "internal.error",
		"message": "Internal server error",
	})
}
