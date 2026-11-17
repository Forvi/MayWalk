package errors

import "fmt"

// DomainError represents a business logic error with HTTP status
type DomainError struct {
	Code    string
	Message string
	Status  int
}

func (e *DomainError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewDomainError creates a new domain error
func NewDomainError(code, message string, status int) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}
