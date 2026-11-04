package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Forvi/MayWalk/internal/common/errs"
)

func mapError(err error, logger *slog.Logger) (int, errs.ErrorHttpResponse) {
	var appErr *errs.Err

	if errors.As(err, &appErr) {
		logger.Warn(appErr.Message, "type", appErr.Type, "code", appErr.Code)
		switch appErr.Type {
		case errs.TypeValidation:
			return http.StatusBadRequest, errs.ErrorHttpResponse{Code: appErr.Code, Message: appErr.Message}
		case errs.TypeNotFound:
			return http.StatusNotFound, errs.ErrorHttpResponse{Code: appErr.Code, Message: appErr.Message}
		case errs.TypeConflict:
			return http.StatusConflict, errs.ErrorHttpResponse{Code: appErr.Code, Message: appErr.Message}
		default:
			return http.StatusInternalServerError, errs.ErrorHttpResponse{Code: "internal.error", Message: "Internal server error"}
		}
	}

	logger.Error("unhandled system error", "error", err)
	return http.StatusInternalServerError, errs.ErrorHttpResponse{
		Code:    "internal.server_error",
		Message: "An unexpected error occurred. Please try again later.",
	}
}
