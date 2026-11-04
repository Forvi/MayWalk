package domain

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/Forvi/MayWalk/internal/common/errs"
	"golang.org/x/crypto/bcrypt"
)

const (
	MinPasswordLength      = 5
	ErrCodePasswordTooWeak = "user.password.too_weak"
)

type UserPassword []byte

// NewUserPassword validates and creates a Bcrypt hash of the password.
func NewUserPassword(raw string) (UserPassword, error) {
	trimmed := strings.TrimSpace(raw)
	if length := utf8.RuneCountInString(trimmed); length < MinPasswordLength {
		return []byte{}, errs.New(errs.TypeValidation, ErrCodePasswordTooWeak, "password too weak")
	}

	shaSum := sha256.Sum256([]byte(trimmed))
	bytes, err := bcrypt.GenerateFromPassword(shaSum[:], bcrypt.DefaultCost)
	if err != nil {
		return []byte{}, fmt.Errorf("generate hash from password: %w", err)
	}

	return UserPassword(bytes), nil
}

func (p UserPassword) Verify(raw string) bool {
	trimmed := strings.TrimSpace(raw)
	shaSum := sha256.Sum256([]byte(trimmed))
	err := bcrypt.CompareHashAndPassword([]byte(p), shaSum[:])
	return err == nil
}

func (p UserPassword) String() string {
	return string(p)
}
