package models

import (
	"time"

	"gorm.io/gorm"
)

// Class represents a class/section in the school (e.g., Grade 1 - Section A).
type Class struct {
	ID        uint           `gorm:"primarykey"   json:"id"`
	CreatedAt time.Time      `                    json:"created_at"`
	UpdatedAt time.Time      `                    json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"        json:"-"`

	Name           string   `gorm:"size:100;not null"  json:"name"             validate:"required"`
	Section        string   `gorm:"size:10;not null"   json:"section"          validate:"required"`
	ClassTeacherID *uint    `                          json:"class_teacher_id"`
	ClassTeacher   *Teacher `gorm:"foreignKey:ClassTeacherID" json:"class_teacher,omitempty"`
}

// TableName overrides the default table name.
func (Class) TableName() string {
	return "classes"
}

// ClassStudent links a student (by ID) to a class.
type ClassStudent struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `                  json:"created_at"`

	ClassID   uint `gorm:"not null;index"                                  json:"class_id"`
	StudentID uint `gorm:"not null;index;uniqueIndex:idx_class_student"    json:"student_id"`
}

// TableName overrides the default table name.
func (ClassStudent) TableName() string {
	return "class_students"
}
