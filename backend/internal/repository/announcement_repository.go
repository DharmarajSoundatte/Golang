package repository

import (
	"context"
	"errors"

	"github.com/DharmarajSoundatte/Golang/backend/internal/models"
	"gorm.io/gorm"
)

// ErrAnnouncementNotFound is returned when an announcement cannot be found.
var ErrAnnouncementNotFound = errors.New("announcement not found")

// AnnouncementRepository defines persistence operations for announcements.
type AnnouncementRepository interface {
	Create(ctx context.Context, a *models.Announcement) error
	FindAll(ctx context.Context, offset, limit int) ([]models.Announcement, int64, error)
	FindByID(ctx context.Context, id uint) (*models.Announcement, error)
}

type announcementRepository struct {
	db *gorm.DB
}

// NewAnnouncementRepository returns a new AnnouncementRepository implementation.
func NewAnnouncementRepository(db *gorm.DB) AnnouncementRepository {
	return &announcementRepository{db: db}
}

func (r *announcementRepository) Create(ctx context.Context, a *models.Announcement) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *announcementRepository) FindAll(ctx context.Context, offset, limit int) ([]models.Announcement, int64, error) {
	var announcements []models.Announcement
	var total int64

	if err := r.db.WithContext(ctx).Model(&models.Announcement{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&announcements).Error

	return announcements, total, err
}

func (r *announcementRepository) FindByID(ctx context.Context, id uint) (*models.Announcement, error) {
	var a models.Announcement
	err := r.db.WithContext(ctx).First(&a, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrAnnouncementNotFound
	}
	return &a, err
}
