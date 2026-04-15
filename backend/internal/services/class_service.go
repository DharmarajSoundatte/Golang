package services

import (
	"context"

	"github.com/DharmarajSoundatte/Golang/backend/internal/models"
	"github.com/DharmarajSoundatte/Golang/backend/internal/repository"
)

// CreateClassRequest is the DTO for creating a class.
type CreateClassRequest struct {
	Name           string `json:"name"             validate:"required"`
	Section        string `json:"section"          validate:"required"`
	ClassTeacherID *uint  `json:"class_teacher_id"`
}

// AssignStudentRequest is the DTO for assigning a student to a class.
type AssignStudentRequest struct {
	StudentID uint `json:"student_id" validate:"required"`
}

// ClassService defines business logic for class management.
type ClassService interface {
	Create(ctx context.Context, req CreateClassRequest) (*models.Class, error)
	GetAll(ctx context.Context, page, pageSize int) ([]models.Class, int64, error)
	GetByID(ctx context.Context, id uint) (*models.Class, error)
	AssignStudent(ctx context.Context, classID uint, req AssignStudentRequest) (*models.ClassStudent, error)
}

type classService struct {
	repo repository.ClassRepository
}

// NewClassService returns a new ClassService implementation.
func NewClassService(repo repository.ClassRepository) ClassService {
	return &classService{repo: repo}
}

func (s *classService) Create(ctx context.Context, req CreateClassRequest) (*models.Class, error) {
	class := &models.Class{
		Name:           req.Name,
		Section:        req.Section,
		ClassTeacherID: req.ClassTeacherID,
	}
	if err := s.repo.Create(ctx, class); err != nil {
		return nil, err
	}
	return class, nil
}

func (s *classService) GetAll(ctx context.Context, page, pageSize int) ([]models.Class, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.repo.FindAll(ctx, (page-1)*pageSize, pageSize)
}

func (s *classService) GetByID(ctx context.Context, id uint) (*models.Class, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *classService) AssignStudent(ctx context.Context, classID uint, req AssignStudentRequest) (*models.ClassStudent, error) {
	cs := &models.ClassStudent{
		ClassID:   classID,
		StudentID: req.StudentID,
	}
	if err := s.repo.AssignStudent(ctx, cs); err != nil {
		return nil, err
	}
	return cs, nil
}
