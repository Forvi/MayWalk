package app

import (
	"context"
	"fmt"

	"github.com/Forvi/MayWalk/internal/user/domain"
	"github.com/Forvi/MayWalk/internal/user/dto"
	"github.com/jackc/pgx/v5"
)

type UserSaver interface {
	Save(ctx context.Context, user domain.User) (*domain.User, error)
	SaveWithTx(ctx context.Context, user domain.User, tx pgx.Tx) (*domain.User, error)
}

type RegisterUserUsecase struct {
	repo UserSaver
}

func NewRegisterUserUsecase(repo UserSaver) *RegisterUserUsecase {
	return &RegisterUserUsecase{repo: repo}
}

func (uc *RegisterUserUsecase) Execute(ctx context.Context, input domain.UserInput) (*dto.UserResponse, error) {
	user, err := domain.NewUser(&input)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	savedUser, err := uc.repo.Save(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("save user: %w", err)
	}

	var userResponse dto.UserResponse
	userResponse.FromDomain(savedUser)
	return &userResponse, nil
}

func (uc *RegisterUserUsecase) ExecuteWithTx(ctx context.Context, input domain.UserInput, tx pgx.Tx) (*dto.UserResponse, error) {
	user, err := domain.NewUser(&input)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	savedUser, err := uc.repo.SaveWithTx(ctx, user, tx)
	if err != nil {
		return nil, fmt.Errorf("save user: %w", err)
	}

	var userResponse dto.UserResponse
	userResponse.FromDomain(savedUser)
	return &userResponse, nil
}
