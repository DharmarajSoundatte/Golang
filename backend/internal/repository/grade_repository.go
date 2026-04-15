package repository

import (
	"context"
	"errors"

	"github.com/DharmarajSoundatte/Golang/backend/internal/models"
	"gorm.io/gorm"
)

// ErrGradeNotFound is returned when a grade record cannot be found.
var ErrGradeNotFound = errors.New("grade not found")

// GradeRepository defines persistence operations for grades.
type GradeRepository interface {
	Create(ctx context.Context, grade *models.Grade) error
	FindAll(ctx context.Context, studentID, classID uint, offset, limit int) ([]models.Grade, int64, error)
	FindByID(ctx context.Context, id uint) (*models.Grade, error)
}

type gradeRepository struct {
	db *gorm.DB
}

// NewGradeRepository returns a new GradeRepository implementation.
func NewGradeRepository(db *gorm.DB) GradeRepository {
	return &gradeRepository{db: db}
}

func (r *gradeRepository) Create(ctx context.Context, grade *models.Grade) error {
	return r.db.WithContext(ctx).Create(grade).Error
}

func (r *gradeRepository) FindAll(ctx context.Context, studentID, classID uint, offset, limit int) ([]models.Grade, int64, error) {
	var grades []models.Grade
	var total int64

	q := r.db.WithContext(ctx).Model(&models.Grade{})
	if studentID != 0 {
		q = q.Where("student_id = ?", studentID)
	}
	if classID != 0 {
		q = q.Where("class_id = ?", classID)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := q.Offset(offset).Limit(limit).Order("created_at DESC").Find(&grades).Error
	return grades, total, err
}

func (r *gradeRepository) FindByID(ctx context.Context, id uint) (*models.Grade, error) {
	var grade models.Grade
	err := r.db.WithContext(ctx).First(&grade, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrGradeNotFound
	}
	return &grade, err
}
