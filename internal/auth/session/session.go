package session

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"strings"
	"time"

	user "github.com/Forvi/MayWalk/internal/user/domain"
)

const tokenBytes = 16

type SessionID string

func (s SessionID) String() string { return string(s) }

func newSessionID() (SessionID, error) {
	id, err := generateToken()
	if err != nil {
		return "", err
	}
	return SessionID(id), nil
}

type Session struct {
	ID         SessionID   `json:"id"`
	UserID     user.UserID `json:"user_id"`
	RememberMe bool        `json:"remember_me"`
	Config     *SessConfig `json:"-"`
	CreatedAt  time.Time   `json:"created_at"`
	LastUsedAt time.Time   `json:"last_used_at"`
}

func New(userID user.UserID, rememberMe bool, cfg *SessConfig) (*Session, error) {
	id, err := newSessionID()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	return &Session{
		ID:         id,
		UserID:     userID,
		RememberMe: rememberMe,
		Config:     cfg,
		CreatedAt:  now,
		LastUsedAt: now,
	}, nil
}

func (s *Session) IsExpired() bool {
	exp := s.Config.IdleTimeout
	if s.RememberMe {
		exp = s.Config.RememberTimeout
	}

	return time.Since(s.LastUsedAt) >= exp
}

func generateToken() (string, error) {
	bytes := make([]byte, tokenBytes)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func SignID(id, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(id))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func MakeCookieValue(id, secret string) string {
	return id + "." + SignID(id, secret)
}

func VerifyID(cookieVal string, secret string) (string, bool) {
	parts := strings.SplitN(cookieVal, ".", 2)
	if len(parts) != 2 {
		return "", false
	}
	id, sig := parts[0], parts[1]
	expected := SignID(id, secret)
	if !hmac.Equal([]byte(expected), []byte(sig)) {
		return "", false
	}
	return id, true
}
