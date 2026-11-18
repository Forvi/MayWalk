package domain

import (
	"net/http"
	"strings"
	"unicode"

	commonErrors "github.com/Forvi/MayWalk/internal/common/errors"
)

var (
	ErrCodeNameInvalid = "user.name.invalid"
	NameMessageInvalid = "name is empty or contains invalid characters"
)

type UserName string

func NewUserName(raw string) (UserName, error) {
	name := strings.TrimSpace(raw)
	if name == "" {
		return "", commonErrors.NewDomainError(ErrCodeNameInvalid, "name is empty", http.StatusBadRequest)
	}

	if strings.Contains(name, ":") {
		return "", commonErrors.NewDomainError(ErrCodeNameInvalid, NameMessageInvalid, http.StatusBadRequest)
	}

	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '-' {
			return "", commonErrors.NewDomainError(ErrCodeNameInvalid, NameMessageInvalid, http.StatusBadRequest)
		}
	}

	return UserName(name), nil
}

func (n UserName) String() string { return string(n) }

func (n UserName) Equals(other UserName) bool { return n == other }
