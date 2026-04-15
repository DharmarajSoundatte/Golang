package services

import (
	"context"
	"errors"

	"github.com/DharmarajSoundatte/Golang/backend/internal/config"
	"github.com/DharmarajSoundatte/Golang/backend/internal/models"
	"github.com/DharmarajSoundatte/Golang/backend/internal/repository"
	"github.com/DharmarajSoundatte/Golang/backend/pkg/hash"
	jwtpkg "github.com/DharmarajSoundatte/Golang/backend/pkg/jwt"
)

// Auth-specific errors
var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserInactive       = errors.New("account is deactivated")
)

// RegisterRequest is the DTO for user registration.
type RegisterRequest struct {
	Name     string `json:"name"     validate:"required,min=2,max=255"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginRequest is the DTO for user login.
type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse is returned on successful register/login.
type AuthResponse struct {
	Token string              `json:"token"`
	User  models.UserResponse `json:"user"`
}

// AuthService defines the contract for authentication operations.
type AuthService interface {
	Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
}

type authService struct {
	repo repository.UserRepository
	cfg  *config.Config
}

// NewAuthService returns a new AuthService implementation.
func NewAuthService(repo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{repo: repo, cfg: cfg}
}

func (s *authService) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	// Check for duplicate email
	_, err := s.repo.FindByEmail(ctx, req.Email)
	if err == nil {
		return nil, ErrEmailAlreadyExists
	}
	if !errors.Is(err, repository.ErrUserNotFound) {
		return nil, err
	}

	// Hash password
	hashed, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashed,
		Role:     models.RoleUser,
		IsActive: true,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	token, err := jwtpkg.GenerateToken(user.ID, string(user.Role), s.cfg.JWTSecret, s.cfg.JWTExpireHours)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token, User: user.ToResponse()}, nil
}

func (s *authService) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, ErrUserInactive
	}

	if !hash.CheckPassword(req.Password, user.Password) {
		return nil, ErrInvalidCredentials
	}

	token, err := jwtpkg.GenerateToken(user.ID, string(user.Role), s.cfg.JWTSecret, s.cfg.JWTExpireHours)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{Token: token, User: user.ToResponse()}, nil
}
