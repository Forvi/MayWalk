package errors

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

// WriteError writes error response
func WriteError(w http.ResponseWriter, r *http.Request, err error, logger *slog.Logger) {
	if domainErr, ok := errors.AsType[*DomainError](err); ok {
		logger.WarnContext(r.Context(), "domain error", "code", domainErr.Code, "status", domainErr.Status)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(domainErr.Status)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"code":    domainErr.Code,
			"message": domainErr.Message,
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
