package domain

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
)

// UserID personal user identifier in the format UUIDv7.
type UserID uuid.UUID

// NewUserID generates a time-ordered UUIDv7.
// It may return an error if the system time is malfunctioning or cannot be retrieved.
func NewUserID() (UserID, error) {
	id, err := uuid.NewV7()
	return UserID(id), err
}

// ParseUserID created user identifier in the format UUIDv7 from raw string.
func ParseUserID(s string) (UserID, error) {
	uid, err := uuid.Parse(s)
	if err != nil {
		return UserID{}, fmt.Errorf("parse user id: %w", err)
	}
	return UserID(uid), nil
}

// String return UserID as string.
func (id UserID) String() string {
	return id.AsUUID().String()
}

// AsUUID return UserID as uuid.UUID type (google/uuid).
func (id UserID) AsUUID() uuid.UUID {
	return uuid.UUID(id)
}

// Value implement driver.Valuer interface.
func (id UserID) Value() (driver.Value, error) {
	return id.AsUUID().Value()
}

// Scan implement sql.Scanner interface.
func (id *UserID) Scan(src any) error {
	if src == nil {
		return nil
	}

	var uid uuid.UUID
	if err := uid.Scan(src); err != nil {
		return fmt.Errorf("scan UserID: %w", err)
	}

	*id = UserID(uid)
	return nil
}
