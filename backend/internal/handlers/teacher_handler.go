package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/DharmarajSoundatte/Golang/backend/internal/repository"
	"github.com/DharmarajSoundatte/Golang/backend/internal/services"
	"github.com/DharmarajSoundatte/Golang/backend/pkg/response"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// TeacherHandler handles HTTP requests for teacher management.
type TeacherHandler struct {
	svc      services.TeacherService
	validate *validator.Validate
	log      *zap.Logger
}

// NewTeacherHandler creates a new TeacherHandler.
func NewTeacherHandler(svc services.TeacherService, log *zap.Logger) *TeacherHandler {
	return &TeacherHandler{svc: svc, validate: validator.New(), log: log}
}

// List returns all teachers (paginated).
func (h *TeacherHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	teachers, total, err := h.svc.GetAll(r.Context(), page, pageSize)
	if err != nil {
		h.log.Error("list teachers", zap.Error(err))
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"data":  teachers,
		"total": total,
		"page":  page,
	})
}

// Add creates a new teacher (admin only).
func (h *TeacherHandler) Add(w http.ResponseWriter, r *http.Request) {
	var req services.CreateTeacherRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.WriteValidationError(w, err)
		return
	}

	teacher, err := h.svc.Add(r.Context(), req)
	if err != nil {
		if errors.Is(err, services.ErrEmailAlreadyExists) {
			response.WriteError(w, http.StatusConflict, "email already in use")
			return
		}
		h.log.Error("add teacher", zap.Error(err))
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.WriteJSON(w, http.StatusCreated, teacher)
}

// GetByID returns a single teacher by ID.
func (h *TeacherHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseUintParam(r, "id")
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid teacher id")
		return
	}

	teacher, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrTeacherNotFound) {
			response.WriteError(w, http.StatusNotFound, "teacher not found")
			return
		}
		h.log.Error("get teacher", zap.Error(err))
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.WriteJSON(w, http.StatusOK, teacher)
}

// Update modifies a teacher record (admin only).
func (h *TeacherHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseUintParam(r, "id")
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid teacher id")
		return
	}

	var req services.UpdateTeacherRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	teacher, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, repository.ErrTeacherNotFound) {
			response.WriteError(w, http.StatusNotFound, "teacher not found")
			return
		}
		h.log.Error("update teacher", zap.Error(err))
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.WriteJSON(w, http.StatusOK, teacher)
}
