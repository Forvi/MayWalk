package domain

import "github.com/Forvi/MayWalk/internal/common/errors"

import "net/http"

var (
	ErrAlreadyExists      = errors.NewDomainError("user.conflict", "user with this name or email already exists", http.StatusConflict)
	ErrNotFound           = errors.NewDomainError("user.not_found", "user not found", http.StatusNotFound)
	ErrInvalidCredentials = errors.NewDomainError("auth.invalid_credentials", "invalid credentials", http.StatusBadRequest)
	ErrAccountInactive    = errors.NewDomainError("auth.account_inactive", "account is inactive", http.StatusForbidden)
)
