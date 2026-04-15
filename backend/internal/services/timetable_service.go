package services

import (
	"context"

	"github.com/DharmarajSoundatte/Golang/backend/internal/models"
	"github.com/DharmarajSoundatte/Golang/backend/internal/repository"
)

// CreateTimetableRequest is the DTO for adding a timetable entry.
type CreateTimetableRequest struct {
	ClassID      uint   `json:"class_id"      validate:"required"`
	DayOfWeek    string `json:"day_of_week"   validate:"required,oneof=Monday Tuesday Wednesday Thursday Friday Saturday"`
	PeriodNumber int    `json:"period_number" validate:"required,min=1"`
	Subject      string `json:"subject"       validate:"required"`
	TeacherID    uint   `json:"teacher_id"    validate:"required"`
	StartTime    string `json:"start_time"    validate:"required"`
	EndTime      string `json:"end_time"      validate:"required"`
}

// TimetableService defines business logic for timetable management.
type TimetableService interface {
	Create(ctx context.Context, req CreateTimetableRequest) (*models.Timetable, error)
	GetByClass(ctx context.Context, classID uint) ([]models.Timetable, error)
}

type timetableService struct {
	repo repository.TimetableRepository
}

// NewTimetableService returns a new TimetableService implementation.
func NewTimetableService(repo repository.TimetableRepository) TimetableService {
	return &timetableService{repo: repo}
}

func (s *timetableService) Create(ctx context.Context, req CreateTimetableRequest) (*models.Timetable, error) {
	entry := &models.Timetable{
		ClassID:      req.ClassID,
		DayOfWeek:    req.DayOfWeek,
		PeriodNumber: req.PeriodNumber,
		Subject:      req.Subject,
		TeacherID:    req.TeacherID,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
	}
	if err := s.repo.Create(ctx, entry); err != nil {
		return nil, err
	}
	return entry, nil
}

func (s *timetableService) GetByClass(ctx context.Context, classID uint) ([]models.Timetable, error) {
	return s.repo.FindByClass(ctx, classID)
}
