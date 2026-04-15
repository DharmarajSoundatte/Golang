package services

import (
	"context"
	"errors"

	"github.com/DharmarajSoundatte/Golang/backend/internal/models"
	"github.com/DharmarajSoundatte/Golang/backend/internal/repository"
)

// ErrEmailAlreadyExists is returned on duplicate email registration attempts.
var ErrEmailAlreadyExists = errors.New("email already in use")

// UserService defines the business-logic contract for user operations.
type UserService interface {
	GetByID(ctx context.Context, id uint) (*models.UserResponse, error)
	GetAll(ctx context.Context, page, pageSize int) ([]models.UserResponse, int64, error)
	UpdateUser(ctx context.Context, id uint, name string, isActive *bool) (*models.UserResponse, error)
	DeleteUser(ctx context.Context, id uint) error
}

type userService struct {
	repo repository.UserRepository
}

// NewUserService returns a new UserService implementation.
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetByID(ctx context.Context, id uint) (*models.UserResponse, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := user.ToResponse()
	return &resp, nil
}

func (s *userService) GetAll(ctx context.Context, page, pageSize int) ([]models.UserResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	users, total, err := s.repo.FindAll(ctx, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]models.UserResponse, len(users))
	for i, u := range users {
		responses[i] = u.ToResponse()
	}
	return responses, total, nil
}

func (s *userService) UpdateUser(ctx context.Context, id uint, name string, isActive *bool) (*models.UserResponse, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if name != "" {
		user.Name = name
	}
	if isActive != nil {
		user.IsActive = *isActive
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	resp := user.ToResponse()
	return &resp, nil
}

func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
