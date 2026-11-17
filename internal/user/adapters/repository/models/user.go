package models

import (
	"fmt"
	"time"

	"github.com/Forvi/MayWalk/internal/user/domain"
)

type DBUser struct {
	ID        string    `db:"id"`
	Name      string    `db:"username"`
	Email     string    `db:"email"`
	Password  *string   `db:"password_hash"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (u *DBUser) ToDomain() (*domain.User, error) {
	user, err := domain.UserFromRaw(domain.RawUser{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.CreatedAt,
	})
	if err != nil {
		return nil, fmt.Errorf("user mapping to domain: %w", err)
	}
	return &user, nil
}
