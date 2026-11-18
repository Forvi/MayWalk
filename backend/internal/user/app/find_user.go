package app

import (
	"context"
	"fmt"

	"github.com/Forvi/MayWalk/internal/user/domain"
	"github.com/Forvi/MayWalk/internal/user/dto"
)

type UserFinder interface {
	FindByID(ctx context.Context, id domain.UserID) (*domain.User, error)
}

type FindUserUsecase struct {
	repo UserFinder
}

func NewFindUsecase(repo UserFinder) *FindUserUsecase {
	return &FindUserUsecase{repo: repo}
}

func (uc *FindUserUsecase) Execute(ctx context.Context, id string) (*dto.UserResponse, error) {
	uid, err := domain.ParseUserID(id)
	if err != nil {
		return nil, fmt.Errorf("parse user id: %w", err)
	}

	user, err := uc.repo.FindByID(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}

	var userResponse dto.UserResponse
	userResponse.FromDomain(user)
	return &userResponse, nil
}
