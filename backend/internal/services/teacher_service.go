package services

import (
	"context"

	"github.com/DharmarajSoundatte/Golang/backend/internal/models"
	"github.com/DharmarajSoundatte/Golang/backend/internal/repository"
)

// CreateTeacherRequest is the DTO for adding a new teacher.
type CreateTeacherRequest struct {
	Name    string `json:"name"    validate:"required,min=2,max=255"`
	Email   string `json:"email"   validate:"required,email"`
	Phone   string `json:"phone"`
	Subject string `json:"subject"`
}

// UpdateTeacherRequest is the DTO for updating a teacher.
type UpdateTeacherRequest struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Subject  string `json:"subject"`
	IsActive *bool  `json:"is_active"`
}

// TeacherService defines the business logic for teacher management.
type TeacherService interface {
	Add(ctx context.Context, req CreateTeacherRequest) (*models.Teacher, error)
	GetAll(ctx context.Context, page, pageSize int) ([]models.Teacher, int64, error)
	GetByID(ctx context.Context, id uint) (*models.Teacher, error)
	Update(ctx context.Context, id uint, req UpdateTeacherRequest) (*models.Teacher, error)
	Deactivate(ctx context.Context, id uint) error
}

type teacherService struct {
	repo repository.TeacherRepository
}

// NewTeacherService returns a new TeacherService implementation.
func NewTeacherService(repo repository.TeacherRepository) TeacherService {
	return &teacherService{repo: repo}
}

func (s *teacherService) Add(ctx context.Context, req CreateTeacherRequest) (*models.Teacher, error) {
	teacher := &models.Teacher{
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Subject:  req.Subject,
		IsActive: true,
	}
	if err := s.repo.Create(ctx, teacher); err != nil {
		if isUniqueConstraintError(err) {
			return nil, ErrEmailAlreadyExists
		}
		return nil, err
	}
	return teacher, nil
}

// isUniqueConstraintError detects Postgres duplicate-key violations by message text.
func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return len(msg) > 0 && (contains(msg, "duplicate key") || contains(msg, "unique constraint") || contains(msg, "23505"))
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && indexStr(s, sub) >= 0)
}

func indexStr(s, sub string) int {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

func (s *teacherService) GetAll(ctx context.Context, page, pageSize int) ([]models.Teacher, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.repo.FindAll(ctx, offset, pageSize)
}

func (s *teacherService) GetByID(ctx context.Context, id uint) (*models.Teacher, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *teacherService) Update(ctx context.Context, id uint, req UpdateTeacherRequest) (*models.Teacher, error) {
	teacher, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		teacher.Name = req.Name
	}
	if req.Phone != "" {
		teacher.Phone = req.Phone
	}
	if req.Subject != "" {
		teacher.Subject = req.Subject
	}
	if req.IsActive != nil {
		teacher.IsActive = *req.IsActive
	}

	if err := s.repo.Update(ctx, teacher); err != nil {
		return nil, err
	}
	return teacher, nil
}

func (s *teacherService) Deactivate(ctx context.Context, id uint) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.Deactivate(ctx, id)
}
