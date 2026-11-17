package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Forvi/MayWalk/internal/auth/session"
	"github.com/Forvi/MayWalk/internal/common/errs"
	"github.com/Forvi/MayWalk/internal/user/app"
	"github.com/Forvi/MayWalk/internal/user/domain"
	"github.com/Forvi/MayWalk/internal/user/dto"
)

type Middleware func(http.Handler) http.Handler

func Routes(mux *http.ServeMux, h *UserHandler, auth Middleware) {
	// auth not required
	mux.HandleFunc("POST /user/register", h.RegisterUser)
	mux.HandleFunc("POST /user/login", h.LoginUser)
	mux.HandleFunc("POST /user/logout", h.LogoutUser)
	mux.HandleFunc("GET /user/signin/yandex", h.UserYaSignIn)
	mux.HandleFunc("GET /user/signin/yandex/callback", h.UserYaSignInCallback)

	// auth required
	mux.Handle("PATCH /user/{id}/inactive", auth(http.HandlerFunc(h.InactiveUser)))
	mux.Handle("GET /user/{id}", auth(http.HandlerFunc(h.FindUser)))
}

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
	inacticeUsecase *app.InactiveUserUsecase,
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
		inactiveUsecase: inacticeUsecase,
		findUsecase:     findUsecase,
		loginUsecase:    loginUsecase,
		logoutUsecase:   logoutUsecase,
		refreshUsecase:  refreshUsecase,
		yaSignInUsecase: yaSignInUsecase,
	}
}

type UserReq struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserReq struct {
	Login      string `json:"login"`
	Password   string `json:"password"`
	RememberMe bool   `json:"remember_me"`
}

func (h *UserHandler) UserYaSignInCallback(w http.ResponseWriter, r *http.Request) {
	rememberMe := true
	cookie, err := r.Cookie("pkce_verifier")
	if err != nil {
		h.sendError(w, err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "pkce_verifier",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	params := r.URL.Query()
	var code string
	if codes, exists := params["code"]; !exists || strings.TrimSpace(codes[0]) == "" {
		h.sendError(w, domain.ErrInvalidCredentials)
		return
	} else {
		code = codes[0]
	}

	cookieVal, err := h.yaSignInUsecase.Auth(r.Context(), code, cookie.Value, rememberMe)
	if err != nil {
		h.sendError(w, err)
		return
	}

	maxAge := 0
	if rememberMe {
		maxAge = int(h.sessCfg.RememberTimeout.Seconds())
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    cookieVal,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAge,
	})

	h.respond(w, http.StatusOK, nil)
}

func (h *UserHandler) UserYaSignIn(w http.ResponseWriter, r *http.Request) {
	url, verifier := h.yaSignInUsecase.SignIn(r.Context())
	http.SetCookie(w, &http.Cookie{
		Name:     "pkce_verifier",
		Value:    verifier,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   600,
	})
	http.Redirect(w, r, url, http.StatusFound)
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		sess, err := h.refreshUsecase.Execute(r.Context(), cookie.Value)
		if err == nil {
			maxAge := 0
			if sess.RememberMe {
				maxAge = int(h.sessCfg.RememberTimeout.Seconds())
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    cookie.Value,
				Path:     "/",
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
				MaxAge:   maxAge,
			})

			h.respond(w, http.StatusNoContent, nil)
			return
		}
	}

	var req LoginUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, errs.New(errs.TypeValidation, "http.bad_request", "Invalid JSON format"))
		return
	}

	login := dto.UserLogin{
		Login:      req.Login,
		Password:   &req.Password,
		RememberMe: req.RememberMe,
	}

	id, err := h.loginUsecase.Execute(r.Context(), login, h.sessCfg.Secret)
	if err != nil {
		h.sendError(w, err)
		return
	}

	maxAge := 0
	if req.RememberMe {
		maxAge = int(h.sessCfg.RememberTimeout.Seconds())
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    id,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   maxAge,
	})

	h.respond(w, http.StatusNoContent, nil)
}

func (h *UserHandler) LogoutUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		h.respond(w, http.StatusNoContent, nil)
		return
	}

	id, ok := session.VerifyID(cookie.Value, h.sessCfg.Secret)
	if !ok {
		h.respond(w, http.StatusNoContent, nil)
		return
	}

	sessionID := session.SessionID(id)
	if err := h.logoutUsecase.Execute(r.Context(), sessionID); err != nil {
		h.sendError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1, // remove cookie
	})

	h.respond(w, http.StatusNoContent, nil)
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req UserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, errs.New(errs.TypeValidation, "http.bad_request", "Invalid JSON format"))
		return
	}

	input := domain.UserInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: &req.Password,
	}

	data, err := h.registerUsecase.Execute(r.Context(), input)
	if err != nil {
		h.sendError(w, err)
		return
	}

	h.respond(w, http.StatusCreated, data)
}

func (h *UserHandler) InactiveUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.inactiveUsecase.Execute(r.Context(), id); err != nil {
		h.sendError(w, err)
		return
	}
	h.respond(w, http.StatusNoContent, nil)
}

func (h *UserHandler) FindUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	data, err := h.findUsecase.Execute(r.Context(), id)
	if err != nil {
		h.sendError(w, err)
		return
	}
	h.respond(w, http.StatusOK, data)
}

func (h *UserHandler) respond(w http.ResponseWriter, status int, data any) {
	if data == nil {
		w.WriteHeader(status)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func (h *UserHandler) sendError(w http.ResponseWriter, err error) {
	statusCode, errorJSON := mapError(err, h.logger)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(errorJSON)
}
