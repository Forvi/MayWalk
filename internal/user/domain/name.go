package domain

import (
	"strings"

	"github.com/Forvi/MayWalk/internal/common/errs"
)

var (
	ErrCodeNameInvalid = "user.name.invalid"
	NameMessageInvalid = "name is empty"
)

type UserName string

func NewUserName(raw string) (UserName, error) {
	name := strings.TrimSpace(raw)
	if name == "" {
		return "", errs.New(errs.TypeValidation, ErrCodeNameInvalid, NameMessageInvalid)
	}

	return UserName(name), nil
}

func (n UserName) String() string { return string(n) }

func (n UserName) Equals(other UserName) bool { return n == other }
