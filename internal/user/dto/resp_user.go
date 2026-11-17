package dto

import (
	"time"

	"github.com/Forvi/MayWalk/internal/user/domain"
)

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *UserResponse) FromDomain(user *domain.User) {
	u.ID = user.ID.String()
	u.Name = user.Name.String()
	u.Email = user.Email.String()
	u.IsActive = user.IsActive
	u.CreatedAt = user.CreatedAt
	u.UpdatedAt = user.UpdatedAt
}

func (u *UserResponse) ToDomain() (domain.User, error) {
	return domain.UserFromRaw(domain.RawUser{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Password:  nil,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	})
}
