package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Forvi/MayWalk/internal/auth/session"
	commonErrors "github.com/Forvi/MayWalk/internal/common/errors"
)

type sessionHandler interface {
	Save(ctx context.Context, session *session.Session) error
	Get(ctx context.Context, sessionID session.SessionID) (*session.Session, error)
	Delete(ctx context.Context, sessionID session.SessionID) error
	Refresh(ctx context.Context, session *session.Session) error
}

type contextKey struct{}

var userIDKey = contextKey{}

func RequireAuth(sessHandler sessionHandler, secret string, logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_id")
			if err != nil {
				commonErrors.WriteError(w, commonErrors.NewDomainError("auth.unauthorized", "unauthorized", http.StatusUnauthorized), logger)
				return
			}

			id, ok := session.VerifyID(cookie.Value, secret)
			if !ok {
				commonErrors.WriteError(w, commonErrors.NewDomainError("auth.unauthorized", "unauthorized", http.StatusUnauthorized), logger)
				return
			}

			sessionID := session.SessionID(id)

			sess, err := sessHandler.Get(r.Context(), sessionID)
			if err != nil {
				if errors.Is(err, session.ErrCacheNotFound) {
					commonErrors.WriteError(w, commonErrors.NewDomainError("auth.unauthorized", "unauthorized", http.StatusUnauthorized), logger)
				} else {
					commonErrors.WriteError(w, commonErrors.NewDomainError("auth.internal_error", "internal server error", http.StatusInternalServerError), logger)
				}
				return
			}

			if sess.IsExpired() {
				_ = sessHandler.Delete(r.Context(), sessionID)
				commonErrors.WriteError(w, commonErrors.NewDomainError("auth.unauthorized", "unauthorized", http.StatusUnauthorized), logger)
				return
			}

			if err := sessHandler.Refresh(r.Context(), sess); err != nil {
				commonErrors.WriteError(w, commonErrors.NewDomainError("auth.unauthorized", "unauthorized", http.StatusUnauthorized), logger)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, sess.UserID.String())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
