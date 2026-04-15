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

// AnnouncementHandler handles HTTP requests for announcement management.
type AnnouncementHandler struct {
	svc      services.AnnouncementService
	validate *validator.Validate
	log      *zap.Logger
}

// NewAnnouncementHandler creates a new AnnouncementHandler.
func NewAnnouncementHandler(svc services.AnnouncementService, log *zap.Logger) *AnnouncementHandler {
	return &AnnouncementHandler{svc: svc, validate: validator.New(), log: log}
}

// List returns all announcements (paginated).
func (h *AnnouncementHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))

	announcements, total, err := h.svc.GetAll(r.Context(), page, pageSize)
	if err != nil {
		h.log.Error("list announcements", zap.Error(err))
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"data":  announcements,
		"total": total,
		"page":  page,
	})
}

// Post creates a new announcement (admin or teacher).
func (h *AnnouncementHandler) Post(w http.ResponseWriter, r *http.Request) {
	authorID, ok := middleware.GetUserIDFromCtx(r.Context())
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req services.PostAnnouncementRequest
	if err := response.DecodeJSON(r, &req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.WriteValidationError(w, err)
		return
	}

	a, err := h.svc.Post(r.Context(), req, authorID)
	if err != nil {
		h.log.Error("post announcement", zap.Error(err))
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.WriteJSON(w, http.StatusCreated, a)
}
