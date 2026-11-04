package domain

import (
	"fmt"
	"net/mail"
	"strings"
	"unicode/utf8"

	"github.com/Forvi/MayWalk/internal/common/errs"
)

// Constants
const (
	MinMailLength = 5
	MaxMailLength = 255
)

// Errors
const (
	ErrCodeEmailInvalid = "user.email.invalid"
	EmailMessageInvalid = "email address must be between 5 and 255 characters long and follow the format domain@mail"
)

// UserEmail personal user email.
type UserEmail string

func NewUserEmail(raw string) (UserEmail, error) {
	trimmed := strings.TrimSpace(raw)
	length := utf8.RuneCountInString(trimmed)
	if length < MinMailLength || length > MaxMailLength || !strings.Contains(trimmed, "@") {
		return "", errs.New(errs.TypeValidation, ErrCodeEmailInvalid, EmailMessageInvalid)
	}

	email, err := mail.ParseAddress(trimmed)
	if err != nil {
		return "", fmt.Errorf("parse address: %w", err)
	}

	return UserEmail(email.Address), nil
}

func (e UserEmail) String() string { return string(e) }

func (e UserEmail) Equals(other UserEmail) bool { return e == other }
