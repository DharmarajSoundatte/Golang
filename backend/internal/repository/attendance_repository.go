package repository

import (
	"context"
	"errors"
	"time"

	"github.com/DharmarajSoundatte/Golang/backend/internal/models"
	"gorm.io/gorm"
)

// ErrAttendanceNotFound is returned when an attendance record cannot be found.
var ErrAttendanceNotFound = errors.New("attendance record not found")

// AttendanceRepository defines persistence operations for attendance.
type AttendanceRepository interface {
	Create(ctx context.Context, a *models.Attendance) error
	FindByID(ctx context.Context, id uint) (*models.Attendance, error)
	FindByClassAndDate(ctx context.Context, classID uint, date time.Time) ([]models.Attendance, error)
	Update(ctx context.Context, a *models.Attendance) error
}

type attendanceRepository struct {
	db *gorm.DB
}

// NewAttendanceRepository returns a new AttendanceRepository implementation.
func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &attendanceRepository{db: db}
}

func (r *attendanceRepository) Create(ctx context.Context, a *models.Attendance) error {
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *attendanceRepository) FindByID(ctx context.Context, id uint) (*models.Attendance, error) {
	var a models.Attendance
	err := r.db.WithContext(ctx).First(&a, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrAttendanceNotFound
	}
	return &a, err
}

func (r *attendanceRepository) FindByClassAndDate(ctx context.Context, classID uint, date time.Time) ([]models.Attendance, error) {
	var records []models.Attendance
	// Match the entire day: from 00:00:00 to 23:59:59 UTC
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)

	err := r.db.WithContext(ctx).
		Where("class_id = ? AND date >= ? AND date < ?", classID, start, end).
		Order("student_id ASC").
		Find(&records).Error

	return records, err
}

func (r *attendanceRepository) Update(ctx context.Context, a *models.Attendance) error {
	return r.db.WithContext(ctx).Save(a).Error
}
