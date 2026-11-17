package domain

import "github.com/Forvi/MayWalk/internal/common/errs"

var (
	ErrAlreadyExists      = errs.New(errs.TypeConflict, "user.conflict", "user with this name or email already exists")
	ErrNotFound           = errs.New(errs.TypeNotFound, "user.not_found", "user not found")
	ErrInvalidCredentials = errs.New(errs.TypeValidation, "auth.invalid_credentials", "invalid credentials")
	ErrAccountInactive    = errs.New(errs.TypeValidation, "auth.account_inactive", "account is inactive")
)
