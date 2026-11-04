package errs

import "fmt"

type ErrorType string

const (
	TypeValidation ErrorType = "VALIDATION"
	TypeNotFound   ErrorType = "NOT_FOUND"
	TypeConflict   ErrorType = "CONFLICT"
	TypeInternal   ErrorType = "INTERNAL"
)

type Err struct {
	Type    ErrorType
	Code    string
	Message string
}

func New(errType ErrorType, code string, msg string) *Err {
	return &Err{
		Type:    errType,
		Code:    code,
		Message: msg,
	}
}

func (e *Err) Error() string {
	return fmt.Sprintf("[%s] %s: %s", e.Type, e.Code, e.Message)
}
