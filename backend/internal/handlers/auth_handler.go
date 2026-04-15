package handlers

import (
	"errors"
	"net/http"

	"github.com/DharmarajSoundatte/Golang/backend/internal/services"
	"github.com/DharmarajSoundatte/Golang/backend/pkg/response"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// AuthHandler handles HTTP requests for authentication.
type AuthHandler struct {
	authSvc  services.AuthService
	validate *validator.Validate
	log      *zap.Logger
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(authSvc services.AuthService, log *zap.Logger) *AuthHandler {
	return &AuthHandler{
		authSvc:  authSvc,
		validate: validator.New(),
		log:      log,
	}
}

// Register godoc
// @Summary      Register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body services.RegisterRequest true "Register payload"
// @Success      201  {object} services.AuthResponse
// @Failure      400  {object} response.ErrorBody
// @Router       /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req services.RegisterRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.WriteValidationError(w, err)
		return
	}

	resp, err := h.authSvc.Register(r.Context(), req)
	if err != nil {
		if errors.Is(err, services.ErrEmailAlreadyExists) {
			response.WriteError(w, http.StatusConflict, err.Error())
			return
		}
		h.log.Error("register error", zap.Error(err))
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.WriteJSON(w, http.StatusCreated, resp)
}

// Login godoc
// @Summary      Login with email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body services.LoginRequest true "Login payload"
// @Success      200  {object} services.AuthResponse
// @Failure      401  {object} response.ErrorBody
// @Router       /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req services.LoginRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.WriteValidationError(w, err)
		return
	}

	resp, err := h.authSvc.Login(r.Context(), req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			response.WriteError(w, http.StatusUnauthorized, err.Error())
			return
		}
		if errors.Is(err, services.ErrUserInactive) {
			response.WriteError(w, http.StatusForbidden, err.Error())
			return
		}
		h.log.Error("login error", zap.Error(err))
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.WriteJSON(w, http.StatusOK, resp)
}
