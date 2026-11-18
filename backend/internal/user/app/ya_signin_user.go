package app

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"unicode"

	"github.com/Forvi/MayWalk/internal/auth/oauth"
	"github.com/Forvi/MayWalk/internal/auth/session"
	commonErrors "github.com/Forvi/MayWalk/internal/common/errors"
	"github.com/Forvi/MayWalk/internal/user/adapters/repository/models"
	"github.com/Forvi/MayWalk/internal/user/domain"
	"github.com/Forvi/MayWalk/internal/user/dto"
	"github.com/jackc/pgx/v5"
)

type UserNameFindSaver interface {
	BeginTx(ctx context.Context) (pgx.Tx, error)
	Save(ctx context.Context, user domain.User) (*domain.User, error)
	FindByName(ctx context.Context, name domain.UserName) (*domain.User, error)
	FindByEmail(ctx context.Context, email domain.UserEmail) (*domain.User, error)
	FindByID(ctx context.Context, id domain.UserID) (*domain.User, error)
	IsUsernameAvailable(ctx context.Context, name domain.UserName) (bool, error)
}

type UserRegister interface {
	ExecuteWithTx(ctx context.Context, input domain.UserInput, tx pgx.Tx) (*dto.UserResponse, error)
}

type UserLoginer interface {
	Execute(ctx context.Context, data dto.UserLogin, secret string) (string, error)
}

type OAuthAccountRepository interface {
	FindByProvider(ctx context.Context, provider, providerUserID string) (*models.OAuthAccount, error)
	CreateWithTx(ctx context.Context, tx pgx.Tx, account models.OAuthAccountCreate) error
}

type YaSignInUserUsecase struct {
	oauth        *oauth.YaOAuth
	sessCfg      *session.SessConfig
	sessHandler  SessionHandler
	userRepo     UserNameFindSaver
	oauthRepo    OAuthAccountRepository
	userRegister UserRegister
	userLoginer  UserLoginer
}

func NewYaSignInUsecase(
	oauth *oauth.YaOAuth,
	sessCfg *session.SessConfig,
	sessHandler SessionHandler,
	userRepo UserNameFindSaver,
	oauthRepo OAuthAccountRepository,
	userRegister UserRegister,
	userLoginer UserLoginer,
) *YaSignInUserUsecase {
	return &YaSignInUserUsecase{
		oauth:        oauth,
		sessCfg:      sessCfg,
		sessHandler:  sessHandler,
		userRepo:     userRepo,
		oauthRepo:    oauthRepo,
		userRegister: userRegister,
		userLoginer:  userLoginer,
	}
}

func (uc *YaSignInUserUsecase) SignIn(ctx context.Context) (redirectURL, verifier string) {
	verifier = uc.oauth.GeneratePKCE()
	redirectURL = uc.oauth.SignIn(ctx, verifier)
	return redirectURL, verifier
}

func (uc *YaSignInUserUsecase) Auth(ctx context.Context, code, verifier string, rememberMe bool) (string, error) {
	authorized, err := uc.oauth.Authorize(ctx, code, verifier)
	if err != nil {
		return "", fmt.Errorf("authorize: %w", err)
	}

	oauthAccount, err := uc.oauthRepo.FindByProvider(ctx, models.YandexProvider, authorized.ID)
	var userID string
	switch {
	case err == nil:
		userID, err = uc.handleExistingOAuthAccount(ctx, oauthAccount)
		if err != nil {
			return "", err
		}

	case errors.Is(err, domain.ErrNotFound):
		userID, err = uc.findOrCreateUserByEmail(ctx, authorized)
		if err != nil {
			return "", err
		}

	default:
		return "", fmt.Errorf("find oauth account: %w", err)
	}

	return uc.loginUser(ctx, userID, rememberMe)
}

func (uc *YaSignInUserUsecase) handleExistingOAuthAccount(ctx context.Context, oauthAccount *models.OAuthAccount) (string, error) {
	userID, err := domain.ParseUserID(oauthAccount.UserID)
	if err != nil {
		return "", fmt.Errorf("parse user id: %w", err)
	}

	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("find user: %w", err)
	}

	if !user.IsActive {
		return "", domain.ErrAccountInactive
	}

	return userID.String(), nil
}

func (uc *YaSignInUserUsecase) findOrCreateUserByEmail(ctx context.Context, authorized *oauth.AuthorizedUser) (string, error) {
	email := getEmailFromAuthorized(authorized)

	user, err := uc.userRepo.FindByEmail(ctx, domain.UserEmail(email))
	switch {
	case err == nil:
		return uc.handleExistingUserByEmail(user)

	case errors.Is(err, domain.ErrNotFound):
		return uc.createNewUserFromOAuth(ctx, authorized, email)

	default:
		return "", fmt.Errorf("find user by email: %w", err)
	}
}

func getEmailFromAuthorized(authorized *oauth.AuthorizedUser) string {
	email := authorized.DefaultEmail
	if email == "" && len(authorized.Emails) > 0 {
		email = authorized.Emails[0]
	}
	return email
}

func (uc *YaSignInUserUsecase) handleExistingUserByEmail(user *domain.User) (string, error) {
	if !user.IsActive {
		return "", domain.ErrAccountInactive
	}

	if user.Password != nil {
		return "", commonErrors.NewDomainError("auth.email_already_registered",
			"email is already registered with password.", http.StatusConflict)
	}

	return user.ID.String(), nil
}

func (uc *YaSignInUserUsecase) createNewUserFromOAuth(ctx context.Context, authorized *oauth.AuthorizedUser, email string) (string, error) {
	tx, txErr := uc.userRepo.BeginTx(ctx)
	if txErr != nil {
		return "", fmt.Errorf("begin tx: %w", txErr)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	username, err := uc.generateUniqueUsername(ctx, authorized.Login)
	if err != nil {
		return "", fmt.Errorf("generate username: %w", err)
	}

	input := domain.UserInput{
		Name:     username,
		Email:    email,
		Password: nil,
	}
	newUser, err := uc.userRegister.ExecuteWithTx(ctx, input, tx)
	if err != nil {
		return "", fmt.Errorf("register user: %w", err)
	}

	err = uc.oauthRepo.CreateWithTx(ctx, tx, models.OAuthAccountCreate{
		UserID:         newUser.ID,
		Provider:       models.YandexProvider,
		ProviderUserID: authorized.ID,
	})
	if err != nil {
		return "", fmt.Errorf("create oauth account: %w", err)
	}

	txErr = tx.Commit(ctx)
	if txErr != nil {
		return "", fmt.Errorf("commit tx: %w", txErr)
	}

	return newUser.ID, nil
}

func (uc *YaSignInUserUsecase) generateUniqueUsername(ctx context.Context, baseUsername string) (string, error) {
	username := normalizeUsername(baseUsername)
	if username == "" {
		return "", fmt.Errorf("invalid username")
	}

	available, err := uc.userRepo.IsUsernameAvailable(ctx, domain.UserName(username))
	if err != nil {
		return "", fmt.Errorf("check username availability: %w", err)
	}

	if available {
		return username, nil
	}

	for range 10 {
		suffix, err := randomDigit(int64(1000))
		if err != nil {
			return "", err
		}
		newUsername := fmt.Sprintf("%s_%d", username, suffix)

		available, err := uc.userRepo.IsUsernameAvailable(ctx, domain.UserName(newUsername))
		if err != nil {
			return "", fmt.Errorf("check username availability: %w", err)
		}

		if available {
			return newUsername, nil
		}
	}

	return "", fmt.Errorf("could not generate unique username")
}

func normalizeUsername(username string) string {
	var result strings.Builder
	for _, r := range username {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func (uc *YaSignInUserUsecase) loginUser(ctx context.Context, userID string, rememberMe bool) (string, error) {
	uid, err := domain.ParseUserID(userID)
	if err != nil {
		return "", fmt.Errorf("parse user id: %w", err)
	}
	sess, err := session.New(uid, rememberMe, uc.sessCfg)
	if err != nil {
		return "", fmt.Errorf("create session: %w", err)
	}

	if err := uc.sessHandler.Save(ctx, sess); err != nil {
		return "", fmt.Errorf("save session: %w", err)
	}

	return session.MakeCookieValue(sess.ID.String(), uc.sessCfg.Secret), nil
}

func randomDigit(edge int64) (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(edge))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()), nil
}
