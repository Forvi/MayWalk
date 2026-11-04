package app

import (
	"context"
	"fmt"

	"github.com/Forvi/MayWalk/internal/user/domain"
)

type UserSaver interface {
	Save(ctx context.Context, user domain.User) error
}

type CreateUserUsecase struct {
	repo UserSaver
}

func NewCreateUserUsecase(repo UserSaver) *CreateUserUsecase {
	return &CreateUserUsecase{repo: repo}
}

func (uc *CreateUserUsecase) Execute(ctx context.Context, input domain.UserInput) error {
	user, err := domain.NewUser(&input)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return uc.repo.Save(ctx, user)
}
