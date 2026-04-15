package models

import "time"

// Grade records marks for a student in a subject/exam.
type Grade struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `                  json:"created_at"`
	UpdatedAt time.Time `                  json:"updated_at"`

	StudentID uint    `gorm:"not null;index" json:"student_id" validate:"required"`
	ClassID   uint    `gorm:"not null;index" json:"class_id"   validate:"required"`
	Subject   string  `gorm:"size:100;not null" json:"subject" validate:"required"`
	ExamName  string  `gorm:"size:100;not null" json:"exam_name" validate:"required"`
	Marks     float64 `gorm:"not null"       json:"marks"     validate:"required,min=0"`
	MaxMarks  float64 `gorm:"not null;default:100" json:"max_marks" validate:"required,min=1"`
	TeacherID uint    `gorm:"not null"       json:"teacher_id" validate:"required"`
}

// TableName overrides the default table name.
func (Grade) TableName() string {
	return "grades"
}
