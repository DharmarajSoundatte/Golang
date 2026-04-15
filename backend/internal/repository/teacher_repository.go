package repository

import (
	"context"
	"errors"

	"github.com/DharmarajSoundatte/Golang/backend/internal/models"
	"gorm.io/gorm"
)

// ErrTeacherNotFound is returned when a teacher cannot be found.
var ErrTeacherNotFound = errors.New("teacher not found")

// TeacherRepository defines persistence operations for teachers.
type TeacherRepository interface {
	Create(ctx context.Context, teacher *models.Teacher) error
	FindByID(ctx context.Context, id uint) (*models.Teacher, error)
	FindAll(ctx context.Context, offset, limit int) ([]models.Teacher, int64, error)
	Update(ctx context.Context, teacher *models.Teacher) error
	Deactivate(ctx context.Context, id uint) error
}

type teacherRepository struct {
	db *gorm.DB
}

// NewTeacherRepository returns a new TeacherRepository implementation.
func NewTeacherRepository(db *gorm.DB) TeacherRepository {
	return &teacherRepository{db: db}
}

func (r *teacherRepository) Create(ctx context.Context, teacher *models.Teacher) error {
	return r.db.WithContext(ctx).Create(teacher).Error
}

func (r *teacherRepository) FindByID(ctx context.Context, id uint) (*models.Teacher, error) {
	var teacher models.Teacher
	err := r.db.WithContext(ctx).First(&teacher, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrTeacherNotFound
	}
	return &teacher, err
}

func (r *teacherRepository) FindAll(ctx context.Context, offset, limit int) ([]models.Teacher, int64, error) {
	var teachers []models.Teacher
	var total int64

	if err := r.db.WithContext(ctx).Model(&models.Teacher{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&teachers).Error

	return teachers, total, err
}

func (r *teacherRepository) Update(ctx context.Context, teacher *models.Teacher) error {
	return r.db.WithContext(ctx).Save(teacher).Error
}

func (r *teacherRepository) Deactivate(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&models.Teacher{}).Where("id = ?", id).Update("is_active", false).Error
}
