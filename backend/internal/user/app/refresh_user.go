package app

import (
	"context"
	"fmt"

	"github.com/Forvi/MayWalk/internal/auth/session"
)

type SessionHandler interface {
	Save(ctx context.Context, session *session.Session) error
	Get(ctx context.Context, sessionID session.SessionID) (*session.Session, error)
	Delete(ctx context.Context, sessionID session.SessionID) error
	Refresh(ctx context.Context, session *session.Session) error
}

type RefreshUserUsecase struct {
	sessCfg     *session.SessConfig
	sessHandler SessionHandler
}

func NewRefreshUsecase(sessCfg *session.SessConfig, sessHandler SessionHandler) *RefreshUserUsecase {
	return &RefreshUserUsecase{sessCfg: sessCfg, sessHandler: sessHandler}
}

func (uc *RefreshUserUsecase) Execute(ctx context.Context, cookieValue string) (*session.Session, error) {
	id, ok := session.VerifyID(cookieValue, uc.sessCfg.Secret)
	if !ok {
		return nil, fmt.Errorf("invalid session signature")
	}

	sessionID := session.SessionID(id)

	sess, err := uc.sessHandler.Get(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("get session: %w", err)
	}

	if sess.IsExpired() {
		_ = uc.sessHandler.Delete(ctx, sessionID)
		return nil, fmt.Errorf("session expired")
	}

	if err := uc.sessHandler.Refresh(ctx, sess); err != nil {
		return nil, fmt.Errorf("refresh session: %w", err)
	}
	return sess, nil
}
