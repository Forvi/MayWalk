package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Forvi/MayWalk/internal/auth/session"
)

type sessionHandler interface {
	Save(ctx context.Context, session *session.Session) error
	Get(ctx context.Context, sessionID session.SessionID) (*session.Session, error)
	Delete(ctx context.Context, sessionID session.SessionID) error
	Refresh(ctx context.Context, session *session.Session) error
}

type contextKey struct{}

var userIDKey = contextKey{}

func RequireAuth(sessHandler sessionHandler, secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_id")
			if err != nil {
				sendUnauthorizedErr(w)
				return
			}

			id, ok := session.VerifyID(cookie.Value, secret)
			if !ok {
				sendUnauthorizedErr(w)
				return
			}

			sessionID := session.SessionID(id)

			sess, err := sessHandler.Get(r.Context(), sessionID)
			if err != nil {
				if errors.Is(err, session.ErrCacheNotFound) {
					sendUnauthorizedErr(w)
				} else {
					sendUnexpectedErr(w)
				}
				return
			}

			if sess.IsExpired() {
				_ = sessHandler.Delete(r.Context(), sessionID)
				sendUnauthorizedErr(w)
				return
			}

			if err := sessHandler.Refresh(r.Context(), sess); err != nil {
				sendUnauthorizedErr(w)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, sess.UserID.String())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func sendUnauthorizedErr(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "unauthorized",
	})
}

func sendUnexpectedErr(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "internal_server_error",
	})
}
