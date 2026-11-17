package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Forvi/MayWalk/internal/common/errs"
	"github.com/Forvi/MayWalk/internal/user/domain"
)

func mapError(err error, logger *slog.Logger) (int, errs.ErrorHttpResponse) {
	if errors.Is(err, domain.ErrAccountInactive) {
		logger.Warn("account inactive", "error", err)
		return http.StatusForbidden, errs.ErrorHttpResponse{
			Code:    "auth.account_inactive",
			Message: "Your account has been deactivated. Please contact support.",
		}
	}

	if appErr, ok := errors.AsType[*errs.Err](err); ok {
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
