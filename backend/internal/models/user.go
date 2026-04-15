package models

import (
	"time"

	"gorm.io/gorm"
)

// Role represents user roles in the system.
type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

// User is the primary user model stored in PostgreSQL.
type User struct {
	ID        uint           `gorm:"primarykey"            json:"id"`
	CreatedAt time.Time      `                             json:"created_at"`
	UpdatedAt time.Time      `                             json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"                 json:"-"`

	Name     string `gorm:"size:255;not null"              json:"name"     validate:"required,min=2,max=255"`
	Email    string `gorm:"size:255;not null;uniqueIndex"   json:"email"    validate:"required,email"`
	Password string `gorm:"size:255;not null"               json:"-"`
	Role     Role   `gorm:"size:50;not null;default:'user'" json:"role"`
	IsActive bool   `gorm:"default:true"                    json:"is_active"`
}

// TableName overrides the default table name.
func (User) TableName() string {
	return "users"
}

// UserResponse is the safe, outward-facing representation of a user.
type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      Role      `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse converts a User model to a UserResponse (strips sensitive fields).
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
	}
}
