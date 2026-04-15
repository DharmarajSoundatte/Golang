package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/DharmarajSoundatte/Golang/backend/internal/repository"
	"github.com/DharmarajSoundatte/Golang/backend/internal/services"
	"github.com/DharmarajSoundatte/Golang/backend/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// UserHandler handles HTTP requests for user management.
type UserHandler struct {
	userSvc  services.UserService
	validate *validator.Validate
	log      *zap.Logger
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(userSvc services.UserService, log *zap.Logger) *UserHandler {
	return &UserHandler{
		userSvc:  userSvc,
		validate: validator.New(),
		log:      log,
	}
}

// GetAll godoc
// @Summary  List all users (paginated)
// @Tags     users
// @Produce  json
// @Param    page      query int false "Page number"      default(1)
// @Param    page_size query int false "Items per page"  default(20)
// @Success  200 {object} map[string]interface{}
// @Security BearerAuth
// @Router   /users [get]
func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	users, total, err := h.userSvc.GetAll(r.Context(), page, pageSize)
	if err != nil {
		h.log.Error("get all users", zap.Error(err))
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"data":  users,
		"total": total,
		"page":  page,
	})
}

// GetByID godoc
// @Summary  Get a user by ID
// @Tags     users
// @Produce  json
// @Param    id  path int true "User ID"
// @Success  200 {object} models.UserResponse
// @Security BearerAuth
// @Router   /users/{id} [get]
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseUintParam(r, "id")
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.userSvc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			response.WriteError(w, http.StatusNotFound, "user not found")
			return
		}
		h.log.Error("get user by id", zap.Error(err))
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.WriteJSON(w, http.StatusOK, user)
}

// UpdateUser godoc
// @Summary  Update a user's name or active status
// @Tags     users
// @Accept   json
// @Produce  json
// @Param    id   path int                true "User ID"
// @Param    body body updateUserRequest  true "Update payload"
// @Success  200  {object} models.UserResponse
// @Security BearerAuth
// @Router   /users/{id} [put]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseUintParam(r, "id")
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var req updateUserRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.userSvc.UpdateUser(r.Context(), id, req.Name, req.IsActive)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			response.WriteError(w, http.StatusNotFound, "user not found")
			return
		}
		h.log.Error("update user", zap.Error(err))
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.WriteJSON(w, http.StatusOK, user)
}

// DeleteUser godoc
// @Summary  Soft-delete a user
// @Tags     users
// @Param    id path int true "User ID"
// @Success  204
// @Security BearerAuth
// @Router   /users/{id} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseUintParam(r, "id")
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	if err := h.userSvc.DeleteUser(r.Context(), id); err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			response.WriteError(w, http.StatusNotFound, "user not found")
			return
		}
		h.log.Error("delete user", zap.Error(err))
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ── helpers ──────────────────────────────────────────────────────────────────

type updateUserRequest struct {
	Name     string `json:"name"`
	IsActive *bool  `json:"is_active"`
}

func parseUintParam(r *http.Request, key string) (uint, error) {
	val, err := strconv.ParseUint(chi.URLParam(r, key), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}
