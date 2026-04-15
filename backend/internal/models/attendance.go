package models

import "time"

// AttendanceStatus represents attendance state for a student on a given day.
type AttendanceStatus string

const (
	AttendancePresent AttendanceStatus = "present"
	AttendanceAbsent  AttendanceStatus = "absent"
	AttendanceLate    AttendanceStatus = "late"
)

// Attendance records a student's attendance for a class on a specific date.
type Attendance struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `                  json:"created_at"`
	UpdatedAt time.Time `                  json:"updated_at"`

	ClassID    uint             `gorm:"not null;index"                                          json:"class_id"    validate:"required"`
	StudentID  uint             `gorm:"not null;index"                                          json:"student_id"  validate:"required"`
	Date       time.Time        `gorm:"not null;index"                                          json:"date"        validate:"required"`
	Status     AttendanceStatus `gorm:"size:20;not null;default:'present'"                      json:"status"      validate:"required,oneof=present absent late"`
	MarkedByID uint             `gorm:"not null"                                                json:"marked_by_id"`
}

// TableName overrides the default table name.
func (Attendance) TableName() string {
	return "attendances"
}
