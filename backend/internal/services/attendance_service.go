package services

import (
	"context"
	"time"

	"github.com/DharmarajSoundatte/Golang/backend/internal/models"
	"github.com/DharmarajSoundatte/Golang/backend/internal/repository"
)

// MarkAttendanceRequest is the DTO for marking a student's attendance.
type MarkAttendanceRequest struct {
	ClassID   uint                    `json:"class_id"   validate:"required"`
	StudentID uint                    `json:"student_id" validate:"required"`
	Date      time.Time               `json:"date"       validate:"required"`
	Status    models.AttendanceStatus `json:"status"     validate:"required,oneof=present absent late"`
}

// EditAttendanceRequest is the DTO for editing an attendance record.
type EditAttendanceRequest struct {
	Status models.AttendanceStatus `json:"status" validate:"required,oneof=present absent late"`
}

// AttendanceService defines business logic for attendance management.
type AttendanceService interface {
	Mark(ctx context.Context, req MarkAttendanceRequest, markedByID uint) (*models.Attendance, error)
	GetByClassAndDate(ctx context.Context, classID uint, date time.Time) ([]models.Attendance, error)
	Edit(ctx context.Context, id uint, req EditAttendanceRequest) (*models.Attendance, error)
}

type attendanceService struct {
	repo repository.AttendanceRepository
}

// NewAttendanceService returns a new AttendanceService implementation.
func NewAttendanceService(repo repository.AttendanceRepository) AttendanceService {
	return &attendanceService{repo: repo}
}

func (s *attendanceService) Mark(ctx context.Context, req MarkAttendanceRequest, markedByID uint) (*models.Attendance, error) {
	a := &models.Attendance{
		ClassID:    req.ClassID,
		StudentID:  req.StudentID,
		Date:       req.Date.UTC(),
		Status:     req.Status,
		MarkedByID: markedByID,
	}
	if err := s.repo.Create(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *attendanceService) GetByClassAndDate(ctx context.Context, classID uint, date time.Time) ([]models.Attendance, error) {
	return s.repo.FindByClassAndDate(ctx, classID, date)
}

func (s *attendanceService) Edit(ctx context.Context, id uint, req EditAttendanceRequest) (*models.Attendance, error) {
	a, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	a.Status = req.Status
	if err := s.repo.Update(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}
