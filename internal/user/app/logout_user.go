package app

import (
	"context"
	"fmt"

	"github.com/Forvi/MayWalk/internal/auth/session"
)

type SessionDeleter interface {
	Delete(ctx context.Context, sessionID session.SessionID) error
}

type LogoutUserUsecase struct {
	sessHandler SessionDeleter
}

func NewLogoutUsecase(sessHandler SessionDeleter) *LogoutUserUsecase {
	return &LogoutUserUsecase{sessHandler: sessHandler}
}

func (uc *LogoutUserUsecase) Execute(ctx context.Context, sessionID session.SessionID) error {
	err := uc.sessHandler.Delete(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}
