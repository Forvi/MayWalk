package session

import (
	"time"

	"github.com/Forvi/MayWalk/internal/user/domain"
)

type sessionDTO struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	RememberMe bool      `json:"remember_me"`
	CreatedAt  time.Time `json:"created_at"`
	LastUsedAt time.Time `json:"last_used_at"`
}

func toSessionDTO(session *Session) sessionDTO {
	return sessionDTO{
		ID:         string(session.ID),
		UserID:     session.UserID.String(),
		RememberMe: session.RememberMe,
		CreatedAt:  session.CreatedAt,
		LastUsedAt: session.LastUsedAt,
	}
}

func (s *sessionDTO) toSession(cfg *SessConfig) (*Session, error) {
	uid, err := domain.ParseUserID(s.UserID)
	if err != nil {
		return nil, err
	}

	return &Session{
		ID:         SessionID(s.ID),
		UserID:     uid,
		RememberMe: s.RememberMe,
		Config:     cfg,
		CreatedAt:  s.CreatedAt,
		LastUsedAt: s.LastUsedAt,
	}, nil
}
