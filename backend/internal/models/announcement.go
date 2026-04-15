package models

import (
	"time"

	"gorm.io/gorm"
)

// TargetAudience defines who an announcement is addressed to.
type TargetAudience string

const (
	AudienceAll   TargetAudience = "all"
	AudienceClass TargetAudience = "class"
	AudienceRole  TargetAudience = "role"
)

// Announcement is a message posted by admin or teachers to a target audience.
type Announcement struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `                  json:"created_at"`
	UpdatedAt time.Time      `                  json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"      json:"-"`

	Title          string         `gorm:"size:255;not null"                        json:"title"           validate:"required"`
	Content        string         `gorm:"type:text;not null"                       json:"content"         validate:"required"`
	AuthorID       uint           `gorm:"not null"                                 json:"author_id"`
	TargetAudience TargetAudience `gorm:"size:20;not null;default:'all'"           json:"target_audience" validate:"required,oneof=all class role"`
	TargetClassID  *uint          `                                                json:"target_class_id"`
	TargetRole     *string        `gorm:"size:50"                                  json:"target_role"`
}

// TableName overrides the default table name.
func (Announcement) TableName() string {
	return "announcements"
}
