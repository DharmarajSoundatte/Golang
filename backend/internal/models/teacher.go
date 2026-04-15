package models

import (
	"time"

	"gorm.io/gorm"
)

// Teacher represents a teacher in the school system.
type Teacher struct {
	ID        uint           `gorm:"primarykey"                    json:"id"`
	CreatedAt time.Time      `                                     json:"created_at"`
	UpdatedAt time.Time      `                                     json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"                         json:"-"`

	Name     string `gorm:"size:255;not null"             json:"name"     validate:"required,min=2,max=255"`
	Email    string `gorm:"size:255;not null;uniqueIndex"  json:"email"    validate:"required,email"`
	Phone    string `gorm:"size:20"                       json:"phone"`
	Subject  string `gorm:"size:255"                      json:"subject"`
	IsActive bool   `gorm:"default:true"                  json:"is_active"`
}

// TableName overrides the default table name.
func (Teacher) TableName() string {
	return "teachers"
}
