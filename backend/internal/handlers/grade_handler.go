package handlers

import (
	"net/http"
	"strconv"

	"github.com/DharmarajSoundatte/Golang/backend/internal/middleware"
	"github.com/DharmarajSoundatte/Golang/backend/internal/services"
	"github.com/DharmarajSoundatte/Golang/backend/pkg/response"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// GradeHandler handles HTTP requests for grade management.
type GradeHandler struct {
	svc      services.GradeService
	validate *validator.Validate
	log      *zap.Logger
}

// NewGradeHandler creates a new GradeHandler.
func NewGradeHandler(svc services.GradeService, log *zap.Logger) *GradeHandler {
	return &GradeHandler{svc: svc, validate: validator.New(), log: log}
}

// List returns grades filtered by optional student_id and class_id query params.
func (h *GradeHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	studentIDVal, _ := strconv.ParseUint(r.URL.Query().Get("student_id"), 10, 64)
	classIDVal, _ := strconv.ParseUint(r.URL.Query().Get("class_id"), 10, 64)

	grades, total, err := h.svc.GetAll(r.Context(), uint(studentIDVal), uint(classIDVal), page, pageSize)
	if err != nil {
		h.log.Error("list grades", zap.Error(err))
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"data":  grades,
		"total": total,
		"page":  page,
	})
}

// Enter records marks for a student (teacher only).
func (h *GradeHandler) Enter(w http.ResponseWriter, r *http.Request) {
	teacherID, ok := middleware.GetUserIDFromCtx(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req services.EnterGradeRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.WriteValidationError(w, err)
		return
	}

	grade, err := h.svc.Enter(r.Context(), req, teacherID)
	if err != nil {
		h.log.Error("enter grade", zap.Error(err))
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.WriteJSON(w, http.StatusCreated, grade)
}
