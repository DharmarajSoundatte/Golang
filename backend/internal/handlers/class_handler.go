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

// ClassHandler handles HTTP requests for class management.
type ClassHandler struct {
	svc      services.ClassService
	validate *validator.Validate
	log      *zap.Logger
}

// NewClassHandler creates a new ClassHandler.
func NewClassHandler(svc services.ClassService, log *zap.Logger) *ClassHandler {
	return &ClassHandler{svc: svc, validate: validator.New(), log: log}
}

// List returns all classes (paginated).
func (h *ClassHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	classes, total, err := h.svc.GetAll(r.Context(), page, pageSize)
	if err != nil {
		h.log.Error("list classes", zap.Error(err))
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"data":  classes,
		"total": total,
		"page":  page,
	})
}

// Create adds a new class (admin only).
func (h *ClassHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req services.CreateClassRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.WriteValidationError(w, err)
		return
	}

	class, err := h.svc.Create(r.Context(), req)
	if err != nil {
		h.log.Error("create class", zap.Error(err))
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.WriteJSON(w, http.StatusCreated, class)
}

// AssignStudent links a student to a class (admin only).
func (h *ClassHandler) AssignStudent(w http.ResponseWriter, r *http.Request) {
	classID, err := parseUintParam(r, "id")
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid class id")
		return
	}

	var req services.AssignStudentRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.WriteValidationError(w, err)
		return
	}

	cs, err := h.svc.AssignStudent(r.Context(), classID, req)
	if err != nil {
		if errors.Is(err, repository.ErrClassNotFound) {
			response.WriteError(w, http.StatusNotFound, "class not found")
			return
		}
		h.log.Error("assign student to class", zap.Error(err))
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.WriteJSON(w, http.StatusCreated, cs)
}
