package app

import (
	"context"
	"fmt"

	"github.com/Forvi/MayWalk/internal/user/domain"
)

type UserFindSaver interface {
	Save(ctx context.Context, user domain.User) error
	FindByID(ctx context.Context, id domain.UserID) (*domain.User, error)
}

type InactiveUserUsecase struct {
	repo UserFindSaver
}

func NewInactiveUsecase(repo UserFindSaver) *InactiveUserUsecase {
	return &InactiveUserUsecase{repo: repo}
}

func (uc *InactiveUserUsecase) Execute(ctx context.Context, id string) error {
	uid, err := domain.ParseUserID(id)
	if err != nil {
		return fmt.Errorf("parse uuid: %w", err)
	}

	user, err := uc.repo.FindByID(ctx, uid)
	if err != nil {
		return fmt.Errorf("find user: %w", err)
	}
	user.Inactive()
	if err := uc.repo.Save(ctx, *user); err != nil {
		return fmt.Errorf("inactive user by id: %w", err)
	}
	return nil
}
