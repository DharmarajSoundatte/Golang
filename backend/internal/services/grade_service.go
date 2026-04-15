package services

import (
	"context"

	"github.com/DharmarajSoundatte/Golang/backend/internal/models"
	"github.com/DharmarajSoundatte/Golang/backend/internal/repository"
)

// EnterGradeRequest is the DTO for entering marks.
type EnterGradeRequest struct {
	StudentID uint    `json:"student_id" validate:"required"`
	ClassID   uint    `json:"class_id"   validate:"required"`
	Subject   string  `json:"subject"    validate:"required"`
	ExamName  string  `json:"exam_name"  validate:"required"`
	Marks     float64 `json:"marks"      validate:"required,min=0"`
	MaxMarks  float64 `json:"max_marks"  validate:"required,min=1"`
}

// GradeService defines business logic for grade management.
type GradeService interface {
	Enter(ctx context.Context, req EnterGradeRequest, teacherID uint) (*models.Grade, error)
	GetAll(ctx context.Context, studentID, classID uint, page, pageSize int) ([]models.Grade, int64, error)
}

type gradeService struct {
	repo repository.GradeRepository
}

// NewGradeService returns a new GradeService implementation.
func NewGradeService(repo repository.GradeRepository) GradeService {
	return &gradeService{repo: repo}
}

func (s *gradeService) Enter(ctx context.Context, req EnterGradeRequest, teacherID uint) (*models.Grade, error) {
	grade := &models.Grade{
		StudentID: req.StudentID,
		ClassID:   req.ClassID,
		Subject:   req.Subject,
		ExamName:  req.ExamName,
		Marks:     req.Marks,
		MaxMarks:  req.MaxMarks,
		TeacherID: teacherID,
	}
	if err := s.repo.Create(ctx, grade); err != nil {
		return nil, err
	}
	return grade, nil
}

func (s *gradeService) GetAll(ctx context.Context, studentID, classID uint, page, pageSize int) ([]models.Grade, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.repo.FindAll(ctx, studentID, classID, (page-1)*pageSize, pageSize)
}
