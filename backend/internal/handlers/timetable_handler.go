package handlers

import (
	"net/http"

	"github.com/DharmarajSoundatte/Golang/backend/internal/services"
	"github.com/DharmarajSoundatte/Golang/backend/pkg/response"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// TimetableHandler handles HTTP requests for timetable management.
type TimetableHandler struct {
	svc      services.TimetableService
	validate *validator.Validate
	log      *zap.Logger
}

// NewTimetableHandler creates a new TimetableHandler.
func NewTimetableHandler(svc services.TimetableService, log *zap.Logger) *TimetableHandler {
	return &TimetableHandler{svc: svc, validate: validator.New(), log: log}
}

// GetByClass returns the timetable for a specific class.
func (h *TimetableHandler) GetByClass(w http.ResponseWriter, r *http.Request) {
	classID, err := parseUintParam(r, "class_id")
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid class_id")
		return
	}

	entries, err := h.svc.GetByClass(r.Context(), classID)
	if err != nil {
		h.log.Error("get timetable", zap.Error(err))
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"data": entries,
	})
}

// Create adds a new timetable entry (admin only).
func (h *TimetableHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req services.CreateTimetableRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.WriteValidationError(w, err)
		return
	}

	entry, err := h.svc.Create(r.Context(), req)
	if err != nil {
		h.log.Error("create timetable entry", zap.Error(err))
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.WriteJSON(w, http.StatusCreated, entry)
}
