package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/Forvi/MayWalk/internal/auth/session"
	"github.com/Forvi/MayWalk/internal/user/domain"
	"github.com/Forvi/MayWalk/internal/user/dto"
)

type SessionSaver interface {
	Save(ctx context.Context, session *session.Session) error
}

type UserLoginFinder interface {
	FindByName(ctx context.Context, name domain.UserName) (*domain.User, error)
	FindByEmail(ctx context.Context, email domain.UserEmail) (*domain.User, error)
}

type LoginUserUsecase struct {
	sessCfg     *session.SessConfig
	sessHandler SessionSaver
	repo        UserLoginFinder
}

func NewLoginUsecase(sessCfg *session.SessConfig, sessHandler SessionSaver, repo UserLoginFinder) *LoginUserUsecase {
	return &LoginUserUsecase{
		sessCfg:     sessCfg,
		sessHandler: sessHandler,
		repo:        repo,
	}
}

func (uc *LoginUserUsecase) Execute(ctx context.Context, data dto.UserLogin, secret string) (string, error) {
	var exists *domain.User
	var err error

	if name, nameErr := domain.NewUserName(data.Login); nameErr == nil {
		exists, err = uc.repo.FindByName(ctx, name)
		if err == nil {
			goto checkPassword
		}
		if !errors.Is(err, domain.ErrNotFound) {
			return "", domain.ErrInvalidCredentials
		}
	}

	if email, emailErr := domain.NewUserEmail(data.Login); emailErr == nil {
		exists, err = uc.repo.FindByEmail(ctx, email)
		if err != nil {
			return "", domain.ErrInvalidCredentials
		}
	}

	if exists == nil {
		return "", domain.ErrInvalidCredentials
	}

checkPassword:
	if data.Password != nil {
		if exists.Password == nil {
			return "", domain.ErrInvalidCredentials
		}
		if ok := exists.Password.Verify(*data.Password); !ok {
			return "", domain.ErrInvalidCredentials
		}
	} else {
		return "", domain.ErrInvalidCredentials
	}

	sess, err := session.New(exists.ID, data.RememberMe, uc.sessCfg)
	if err != nil {
		return "", fmt.Errorf("create session: %w", err)
	}
	if err := uc.sessHandler.Save(ctx, sess); err != nil {
		return "", fmt.Errorf("save session: %w", err)
	}

	return session.MakeCookieValue(sess.ID.String(), secret), nil
}
