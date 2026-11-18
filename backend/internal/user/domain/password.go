package domain

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strings"
	"unicode/utf8"

	commonErrors "github.com/Forvi/MayWalk/internal/common/errors"
	"golang.org/x/crypto/bcrypt"
)

const (
	MinPasswordLength      = 5
	ErrCodePasswordTooWeak = "user.password.too_weak"
)

type UserPassword []byte

// NewUserPassword validates and creates a Bcrypt hash of the password.
func NewUserPassword(raw *string) (UserPassword, error) {
	if raw == nil {
		return UserPassword(nil), nil
	}
	trimmed := strings.TrimSpace(*raw)
	if length := utf8.RuneCountInString(trimmed); length < MinPasswordLength {
		return []byte{}, commonErrors.NewDomainError(ErrCodePasswordTooWeak, "password too weak", http.StatusBadRequest)
	}

	shaSum := sha256.Sum256([]byte(trimmed))
	bytes, err := bcrypt.GenerateFromPassword(shaSum[:], bcrypt.DefaultCost)
	if err != nil {
		return []byte{}, fmt.Errorf("generate hash from password: %w", err)
	}

	return UserPassword(bytes), nil
}

func PasswordFromRaw(password string) UserPassword {
	return UserPassword(password)
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
