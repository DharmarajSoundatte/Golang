package repository

import (
	"context"
	"errors"

	"github.com/DharmarajSoundatte/Golang/backend/internal/models"
	"gorm.io/gorm"
)

// ErrClassNotFound is returned when a class cannot be found.
var ErrClassNotFound = errors.New("class not found")

// ClassRepository defines persistence operations for classes.
type ClassRepository interface {
	Create(ctx context.Context, class *models.Class) error
	FindByID(ctx context.Context, id uint) (*models.Class, error)
	FindAll(ctx context.Context, offset, limit int) ([]models.Class, int64, error)
	AssignStudent(ctx context.Context, cs *models.ClassStudent) error
}

type classRepository struct {
	db *gorm.DB
}

// NewClassRepository returns a new ClassRepository implementation.
func NewClassRepository(db *gorm.DB) ClassRepository {
	return &classRepository{db: db}
}

func (r *classRepository) Create(ctx context.Context, class *models.Class) error {
	return r.db.WithContext(ctx).Create(class).Error
}

func (r *classRepository) FindByID(ctx context.Context, id uint) (*models.Class, error) {
	var class models.Class
	err := r.db.WithContext(ctx).Preload("ClassTeacher").First(&class, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrClassNotFound
	}
	return &class, err
}

func (r *classRepository) FindAll(ctx context.Context, offset, limit int) ([]models.Class, int64, error) {
	var classes []models.Class
	var total int64

	if err := r.db.WithContext(ctx).Model(&models.Class{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.WithContext(ctx).
		Preload("ClassTeacher").
		Offset(offset).
		Limit(limit).
		Order("name ASC, section ASC").
		Find(&classes).Error

	return classes, total, err
}

func (r *classRepository) AssignStudent(ctx context.Context, cs *models.ClassStudent) error {
	return r.db.WithContext(ctx).Create(cs).Error
}
