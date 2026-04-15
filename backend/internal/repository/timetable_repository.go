package repository

import (
	"context"
	"errors"

	"github.com/DharmarajSoundatte/Golang/backend/internal/models"
	"gorm.io/gorm"
)

// ErrTimetableNotFound is returned when a timetable entry cannot be found.
var ErrTimetableNotFound = errors.New("timetable entry not found")

// TimetableRepository defines persistence operations for timetable entries.
type TimetableRepository interface {
	Create(ctx context.Context, entry *models.Timetable) error
	FindByClass(ctx context.Context, classID uint) ([]models.Timetable, error)
}

type timetableRepository struct {
	db *gorm.DB
}

// NewTimetableRepository returns a new TimetableRepository implementation.
func NewTimetableRepository(db *gorm.DB) TimetableRepository {
	return &timetableRepository{db: db}
}

func (r *timetableRepository) Create(ctx context.Context, entry *models.Timetable) error {
	return r.db.WithContext(ctx).Create(entry).Error
}

func (r *timetableRepository) FindByClass(ctx context.Context, classID uint) ([]models.Timetable, error) {
	var entries []models.Timetable
	err := r.db.WithContext(ctx).
		Preload("Teacher").
		Where("class_id = ?", classID).
		Order("day_of_week ASC, period_number ASC").
		Find(&entries).Error
	return entries, err
}
