package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Forvi/MayWalk/internal/common/errs"
	"github.com/Forvi/MayWalk/internal/user/app"
	"github.com/Forvi/MayWalk/internal/user/domain"
)

func Routes(mux *http.ServeMux, h *UserHandler) {
	mux.HandleFunc("POST /user", h.CreateUser)
	mux.HandleFunc("PATCH /user/{id}/inactive", h.InactiveUser)
}

type UserHandler struct {
	logger          *slog.Logger
	createUsecase   *app.CreateUserUsecase
	inactiveUsecase *app.InactiveUserUsecase
}

func NewUserHandler(
	logger *slog.Logger,
	createUsecase *app.CreateUserUsecase,
	inacticeUsecase *app.InactiveUserUsecase,
) *UserHandler {
	return &UserHandler{
		logger:          logger,
		createUsecase:   createUsecase,
		inactiveUsecase: inacticeUsecase,
	}
}

type UserReq struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req UserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, errs.New(errs.TypeValidation, "http.bad_request", "Invalid JSON format"))
		return
	}

	input := domain.UserInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	err := h.createUsecase.Execute(r.Context(), input)
	if err != nil {
		h.sendError(w, err)
		return
	}

	h.respond(w, http.StatusCreated, nil)
}

func (h *UserHandler) InactiveUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.inactiveUsecase.Execute(r.Context(), id); err != nil {
		h.sendError(w, err)
		return
	}
	h.respond(w, http.StatusOK, nil)
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
