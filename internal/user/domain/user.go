package domain

import (
	"fmt"
	"time"
)

var ()

type UserInput struct {
	Name     string
	Email    string
	Password string
}

type RawUser struct {
	ID        string
	Name      string
	Email     string
	Password  string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	ID        UserID
	Name      UserName
	Email     UserEmail
	Password  UserPassword
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(input *UserInput) (User, error) {
	id, err := NewUserID()
	if err != nil {
		return User{}, fmt.Errorf("create id: %w", err)
	}

	name, err := NewUserName(input.Name)
	if err != nil {
		return User{}, fmt.Errorf("parse username: %w", err)
	}

	email, err := NewUserEmail(input.Email)
	if err != nil {
		return User{}, fmt.Errorf("parse user email: %w", err)
	}

	passwordHash, err := NewUserPassword(input.Password)
	if err != nil {
		return User{}, fmt.Errorf("parse user password: %w", err)
	}

	now := time.Now().UTC()
	return User{
		ID:        id,
		Name:      name,
		Email:     email,
		Password:  passwordHash,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func UserFromRaw(raw RawUser) (User, error) {
	var id UserID
	err := id.Scan(raw.ID)
	if err != nil {
		return User{}, fmt.Errorf("create id: %w", err)
	}

	name, err := NewUserName(raw.Name)
	if err != nil {
		return User{}, fmt.Errorf("parse username: %w", err)
	}

	email, err := NewUserEmail(raw.Email)
	if err != nil {
		return User{}, fmt.Errorf("parse user email: %w", err)
	}

	passwordHash, err := NewUserPassword(raw.Password)
	if err != nil {
		return User{}, fmt.Errorf("parse user password: %w", err)
	}

	now := time.Now().UTC()
	return User{
		ID:        id,
		Name:      name,
		Email:     email,
		Password:  passwordHash,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (u *User) Rename(raw string) error {
	name, err := NewUserName(raw)
	if err != nil {
		return fmt.Errorf("parse username: %w", err)
	}
	u.Name = name
	return nil
}

func (u *User) UpdateEmail(raw string) error {
	email, err := NewUserEmail(raw)
	if err != nil {
		return fmt.Errorf("parse user email: %w", err)
	}
	u.Email = email
	return nil
}

func (u *User) UpdatePassword(newPass, oldPass string) error {
	pass, err := NewUserPassword(newPass)
	if err != nil {
		return fmt.Errorf("parse user email: %w", err)
	}
	if ok := pass.Verify(oldPass); !ok {
		return ErrInvalidCredentials
	}

	u.Password = pass
	return nil
}

func (u *User) Inactive() {
	u.IsActive = false
}

func (u *User) Active() {
	u.IsActive = true
}
