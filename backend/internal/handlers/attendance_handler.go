package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/DharmarajSoundatte/Golang/backend/internal/middleware"
	"github.com/DharmarajSoundatte/Golang/backend/internal/repository"
	"github.com/DharmarajSoundatte/Golang/backend/internal/services"
	"github.com/DharmarajSoundatte/Golang/backend/pkg/response"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// AttendanceHandler handles HTTP requests for attendance management.
type AttendanceHandler struct {
	svc      services.AttendanceService
	validate *validator.Validate
	log      *zap.Logger
}

// NewAttendanceHandler creates a new AttendanceHandler.
func NewAttendanceHandler(svc services.AttendanceService, log *zap.Logger) *AttendanceHandler {
	return &AttendanceHandler{svc: svc, validate: validator.New(), log: log}
}

// Get returns attendance records filtered by class_id and date query params.
func (h *AttendanceHandler) Get(w http.ResponseWriter, r *http.Request) {
	classIDStr := r.URL.Query().Get("class_id")
	dateStr := r.URL.Query().Get("date")

	classIDVal, err := strconv.ParseUint(classIDStr, 10, 64)
	if err != nil || classIDVal == 0 {
		response.WriteError(w, http.StatusBadRequest, "class_id query param is required")
		return
	}

	date := time.Now().UTC()
	if dateStr != "" {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			response.WriteError(w, http.StatusBadRequest, "date must be in YYYY-MM-DD format")
			return
		}
	}

	records, err := h.svc.GetByClassAndDate(r.Context(), uint(classIDVal), date)
	if err != nil {
		h.log.Error("get attendance", zap.Error(err))
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"data": records,
		"date": date.Format("2006-01-02"),
	})
}

// Mark records attendance for a student (teacher only).
func (h *AttendanceHandler) Mark(w http.ResponseWriter, r *http.Request) {
	markedByID, ok := middleware.GetUserIDFromCtx(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req services.MarkAttendanceRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.WriteValidationError(w, err)
		return
	}

	a, err := h.svc.Mark(r.Context(), req, markedByID)
	if err != nil {
		h.log.Error("mark attendance", zap.Error(err))
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.WriteJSON(w, http.StatusCreated, a)
}

// Edit updates an existing attendance record (teacher only).
func (h *AttendanceHandler) Edit(w http.ResponseWriter, r *http.Request) {
	id, err := parseUintParam(r, "id")
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid attendance id")
		return
	}

	var req services.EditAttendanceRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.WriteValidationError(w, err)
		return
	}

	a, err := h.svc.Edit(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, repository.ErrAttendanceNotFound) {
			response.WriteError(w, http.StatusNotFound, "attendance record not found")
			return
		}
		h.log.Error("edit attendance", zap.Error(err))
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.WriteJSON(w, http.StatusOK, a)
}
