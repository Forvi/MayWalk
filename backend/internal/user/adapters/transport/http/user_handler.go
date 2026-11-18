package http

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Forvi/MayWalk/internal/api"
	"github.com/Forvi/MayWalk/internal/auth/session"
	"github.com/Forvi/MayWalk/internal/user/app"
	"github.com/Forvi/MayWalk/internal/user/domain"
	"github.com/Forvi/MayWalk/internal/user/dto"
)

type UserHandler struct {
	logger          *slog.Logger
	sessCfg         *session.SessConfig
	registerUsecase *app.RegisterUserUsecase
	inactiveUsecase *app.InactiveUserUsecase
	findUsecase     *app.FindUserUsecase
	loginUsecase    *app.LoginUserUsecase
	logoutUsecase   *app.LogoutUserUsecase
	refreshUsecase  *app.RefreshUserUsecase
	yaSignInUsecase *app.YaSignInUserUsecase
}

func NewUserHandler(
	logger *slog.Logger,
	sessCfg *session.SessConfig,
	registerUsecase *app.RegisterUserUsecase,
	inactiveUsecase *app.InactiveUserUsecase,
	findUsecase *app.FindUserUsecase,
	loginUsecase *app.LoginUserUsecase,
	logoutUsecase *app.LogoutUserUsecase,
	refreshUsecase *app.RefreshUserUsecase,
	yaSignInUsecase *app.YaSignInUserUsecase,
) *UserHandler {
	return &UserHandler{
		logger:          logger,
		sessCfg:         sessCfg,
		registerUsecase: registerUsecase,
		inactiveUsecase: inactiveUsecase,
		findUsecase:     findUsecase,
		loginUsecase:    loginUsecase,
		logoutUsecase:   logoutUsecase,
		refreshUsecase:  refreshUsecase,
		yaSignInUsecase: yaSignInUsecase,
	}
}

// PostUserRegister implements api.StrictServerInterface
func (h *UserHandler) PostUserRegister(
	ctx context.Context,
	req api.PostUserRegisterRequestObject,
) (api.PostUserRegisterResponseObject, error) {
	if req.Body == nil {
		return api.PostUserRegister400JSONResponse{
			Code:    "http.bad_request",
			Message: "Request body is required",
		}, nil
	}

	input := domain.UserInput{
		Name:     req.Body.Name,
		Email:    string(req.Body.Email),
		Password: &req.Body.Password,
	}

	data, err := h.registerUsecase.Execute(ctx, input)
	if err != nil {
		return nil, err
	}

	return api.PostUserRegister201JSONResponse{
		Id:        data.ID,
		Name:      data.Name,
		Email:     data.Email,
		IsActive:  data.IsActive,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}, nil
}

// PostUserLogin implements api.StrictServerInterface
func (h *UserHandler) PostUserLogin(
	ctx context.Context,
	req api.PostUserLoginRequestObject,
) (api.PostUserLoginResponseObject, error) {
	// refresh existing session
	if req.Params.SessionId != nil && *req.Params.SessionId != "" {
		sess, err := h.refreshUsecase.Execute(ctx, *req.Params.SessionId)
		if err == nil {
			maxAge := 0
			if sess.RememberMe {
				maxAge = int(h.sessCfg.RememberTimeout.Seconds())
			}

			cookie := &http.Cookie{
				Name:     "session_id",
				Value:    *req.Params.SessionId,
				Path:     "/",
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
				MaxAge:   maxAge,
			}
			cookieStr := cookie.String()

			return api.PostUserLogin204Response{
				Headers: api.PostUserLogin204ResponseHeaders{
					SetCookie: &cookieStr,
				},
			}, nil
		}
	}

	// login
	if req.Body == nil {
		return api.PostUserLogin400JSONResponse{
			Code:    "http.bad_request",
			Message: "Request body is required",
		}, nil
	}

	rememberMe := false
	if req.Body.RememberMe != nil {
		rememberMe = *req.Body.RememberMe
	}

	login := dto.UserLogin{
		Login:      req.Body.Login,
		Password:   &req.Body.Password,
		RememberMe: rememberMe,
	}

	id, err := h.loginUsecase.Execute(ctx, login, h.sessCfg.Secret)
	if err != nil {
		return nil, err
	}

	maxAge := 0
	if rememberMe {
		maxAge = int(h.sessCfg.RememberTimeout.Seconds())
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    id,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAge,
	}
	cookieStr := cookie.String()

	return api.PostUserLogin204Response{
		Headers: api.PostUserLogin204ResponseHeaders{
			SetCookie: &cookieStr,
		},
	}, nil
}

// PostUserLogout implements api.StrictServerInterface
func (h *UserHandler) PostUserLogout(
	ctx context.Context,
	req api.PostUserLogoutRequestObject,
) (api.PostUserLogoutResponseObject, error) {
	if req.Params.SessionId == nil || *req.Params.SessionId == "" {
		cookie := &http.Cookie{
			Name:     "session_id",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   -1,
		}
		cookieStr := cookie.String()
		return api.PostUserLogout204Response{
			Headers: api.PostUserLogout204ResponseHeaders{
				SetCookie: &cookieStr,
			},
		}, nil
	}

	id, ok := session.VerifyID(*req.Params.SessionId, h.sessCfg.Secret)
	if !ok {
		cookie := &http.Cookie{
			Name:     "session_id",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   -1,
		}
		cookieStr := cookie.String()
		return api.PostUserLogout204Response{
			Headers: api.PostUserLogout204ResponseHeaders{
				SetCookie: &cookieStr,
			},
		}, nil
	}

	sessionID := session.SessionID(id)
	if err := h.logoutUsecase.Execute(ctx, sessionID); err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	}
	cookieStr := cookie.String()

	return api.PostUserLogout204Response{
		Headers: api.PostUserLogout204ResponseHeaders{
			SetCookie: &cookieStr,
		},
	}, nil
}

// GetUserSigninYandex implements api.StrictServerInterface
func (h *UserHandler) GetUserSigninYandex(
	ctx context.Context,
	req api.GetUserSigninYandexRequestObject,
) (api.GetUserSigninYandexResponseObject, error) {
	url, verifier := h.yaSignInUsecase.SignIn(ctx)

	cookie := &http.Cookie{
		Name:     "pkce_verifier",
		Value:    verifier,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   600,
	}
	cookieStr := cookie.String()

	return api.GetUserSigninYandex302Response{
		Headers: api.GetUserSigninYandex302ResponseHeaders{
			Location:  &url,
			SetCookie: &cookieStr,
		},
	}, nil
}

// YandexCallbackResponse type for update cookies when oauth callback
type YandexCallbackResponse struct {
	sessionCookie *http.Cookie
	clearCookie   *http.Cookie
}

// VisitGetUserSigninYandexCallbackResponse implement GetUserSigninYandexCallbackResponseObject for SetCookie
func (r YandexCallbackResponse) VisitGetUserSigninYandexCallbackResponse(w http.ResponseWriter) error {
	http.SetCookie(w, r.clearCookie)
	http.SetCookie(w, r.sessionCookie)
	w.WriteHeader(200)
	return nil
}

// GetUserSigninYandexCallback implements api.StrictServerInterface
func (h *UserHandler) GetUserSigninYandexCallback(
	ctx context.Context,
	req api.GetUserSigninYandexCallbackRequestObject,
) (api.GetUserSigninYandexCallbackResponseObject, error) {
	rememberMe := true

	clearCookie := &http.Cookie{
		Name:     "pkce_verifier",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1, // clear cookie
	}

	code := req.Params.Code
	pkceVerifier := req.Params.PkceVerifier

	cookieVal, err := h.yaSignInUsecase.Auth(ctx, code, pkceVerifier, rememberMe)
	if err != nil {
		return nil, err
	}

	maxAge := 0
	if rememberMe {
		maxAge = int(h.sessCfg.RememberTimeout.Seconds())
	}

	sessionCookie := &http.Cookie{
		Name:     "session_id",
		Value:    cookieVal,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAge,
	}

	return &YandexCallbackResponse{
		sessionCookie: sessionCookie,
		clearCookie:   clearCookie,
	}, nil
}

// GetUserId implements api.StrictServerInterface
func (h *UserHandler) GetUserId(
	ctx context.Context,
	req api.GetUserIdRequestObject,
) (api.GetUserIdResponseObject, error) {
	data, err := h.findUsecase.Execute(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return api.GetUserId200JSONResponse{
		Id:        data.ID,
		Name:      data.Name,
		Email:     data.Email,
		IsActive:  data.IsActive,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}, nil
}

// PatchUserIdInactive implements api.StrictServerInterface
func (h *UserHandler) PatchUserIdInactive(
	ctx context.Context,
	req api.PatchUserIdInactiveRequestObject,
) (api.PatchUserIdInactiveResponseObject, error) {
	if err := h.inactiveUsecase.Execute(ctx, req.Id); err != nil {
		return nil, err
	}

	return api.PatchUserIdInactive204Response{}, nil
}
