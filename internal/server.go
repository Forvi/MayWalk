package internal

import (
	"context"

	"github.com/Forvi/MayWalk/internal/api"
	userhttp "github.com/Forvi/MayWalk/internal/user/adapters/transport/http"
)

type Server struct {
	user *userhttp.UserHandler
}

func NewServer(user *userhttp.UserHandler) *Server {
	return &Server{user: user}
}

// ================================================================================================
// 													USER
// ================================================================================================

func (s *Server) PostUserRegister(ctx context.Context, req api.PostUserRegisterRequestObject) (api.PostUserRegisterResponseObject, error) {
	return s.user.PostUserRegister(ctx, req)
}

func (s *Server) PostUserLogin(ctx context.Context, req api.PostUserLoginRequestObject) (api.PostUserLoginResponseObject, error) {
	return s.user.PostUserLogin(ctx, req)
}

func (s *Server) PostUserLogout(ctx context.Context, req api.PostUserLogoutRequestObject) (api.PostUserLogoutResponseObject, error) {
	return s.user.PostUserLogout(ctx, req)
}

func (s *Server) GetUserSigninYandex(ctx context.Context, req api.GetUserSigninYandexRequestObject) (api.GetUserSigninYandexResponseObject, error) {
	return s.user.GetUserSigninYandex(ctx, req)
}

func (s *Server) GetUserSigninYandexCallback(ctx context.Context, req api.GetUserSigninYandexCallbackRequestObject) (api.GetUserSigninYandexCallbackResponseObject, error) {
	return s.user.GetUserSigninYandexCallback(ctx, req)
}

func (s *Server) GetUserId(ctx context.Context, req api.GetUserIdRequestObject) (api.GetUserIdResponseObject, error) {
	return s.user.GetUserId(ctx, req)
}

func (s *Server) PatchUserIdInactive(ctx context.Context, req api.PatchUserIdInactiveRequestObject) (api.PatchUserIdInactiveResponseObject, error) {
	return s.user.PatchUserIdInactive(ctx, req)
}
