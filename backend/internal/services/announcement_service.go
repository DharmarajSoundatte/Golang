package services

import (
	"context"

	"github.com/DharmarajSoundatte/Golang/backend/internal/models"
	"github.com/DharmarajSoundatte/Golang/backend/internal/repository"
)

// PostAnnouncementRequest is the DTO for posting an announcement.
type PostAnnouncementRequest struct {
	Title          string                `json:"title"           validate:"required"`
	Content        string                `json:"content"         validate:"required"`
	TargetAudience models.TargetAudience `json:"target_audience" validate:"required,oneof=all class role"`
	TargetClassID  *uint                 `json:"target_class_id"`
	TargetRole     *string               `json:"target_role"`
}

// AnnouncementService defines business logic for announcements.
type AnnouncementService interface {
	Post(ctx context.Context, req PostAnnouncementRequest, authorID uint) (*models.Announcement, error)
	GetAll(ctx context.Context, page, pageSize int) ([]models.Announcement, int64, error)
}

type announcementService struct {
	repo repository.AnnouncementRepository
}

// NewAnnouncementService returns a new AnnouncementService implementation.
func NewAnnouncementService(repo repository.AnnouncementRepository) AnnouncementService {
	return &announcementService{repo: repo}
}

func (s *announcementService) Post(ctx context.Context, req PostAnnouncementRequest, authorID uint) (*models.Announcement, error) {
	a := &models.Announcement{
		Title:          req.Title,
		Content:        req.Content,
		AuthorID:       authorID,
		TargetAudience: req.TargetAudience,
		TargetClassID:  req.TargetClassID,
		TargetRole:     req.TargetRole,
	}
	if err := s.repo.Create(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *announcementService) GetAll(ctx context.Context, page, pageSize int) ([]models.Announcement, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.repo.FindAll(ctx, (page-1)*pageSize, pageSize)
}
