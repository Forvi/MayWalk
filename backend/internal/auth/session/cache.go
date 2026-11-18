package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrCacheNotFound   = errors.New("cached value not found")
	ErrCacheMarshaling = errors.New("cache marshaling")
	ErrDeleteCache     = errors.New("delete cache")
	ErrExpiredSession  = errors.New("session is expired")
)

const prefix = "session"

type SessionCache struct {
	cache   *redis.Client
	sessCfg *SessConfig
}

func NewSessionCache(redis *redis.Client, sessCfg *SessConfig) *SessionCache {
	return &SessionCache{
		cache:   redis,
		sessCfg: sessCfg,
	}
}

func (r *SessionCache) Save(ctx context.Context, session *Session) error {
	data, err := json.Marshal(toSessionDTO(session))
	if err != nil {
		return fmt.Errorf("cache marshaling: %w", err)
	}

	exp := r.sessCfg.IdleTimeout
	if session.RememberMe {
		exp = r.sessCfg.RememberTimeout
	}

	key := createKey(string(session.ID))
	err = r.cache.Set(ctx, key, data, exp).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *SessionCache) Refresh(ctx context.Context, session *Session) error {
	session.LastUsedAt = time.Now()
	data, err := json.Marshal(toSessionDTO(session))
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	key := createKey(string(session.ID))
	exp := r.sessCfg.IdleTimeout
	if session.RememberMe {
		exp = r.sessCfg.RememberTimeout
	}

	return r.cache.Set(ctx, key, data, exp).Err()
}

func (r *SessionCache) Get(ctx context.Context, sessionID SessionID) (*Session, error) {
	sess := &sessionDTO{}
	key := createKey(string(sessionID))
	data, err := r.cache.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrCacheNotFound
		}
		return nil, err
	}

	err = json.Unmarshal(data, sess)
	if err != nil {
		return nil, ErrCacheMarshaling
	}

	return sess.toSession(r.sessCfg)
}

func (r *SessionCache) Delete(ctx context.Context, sessionID SessionID) error {
	key := createKey(string(sessionID))
	err := r.cache.Del(ctx, key).Err()
	if err != nil {
		return ErrDeleteCache
	}
	return nil
}

func createKey(sessionID string) string {
	key := []string{prefix, sessionID}
	return strings.Join(key, ":")
}
