package models

import "time"

// Timetable defines a period slot for a class on a specific day.
type Timetable struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `                  json:"created_at"`
	UpdatedAt time.Time `                  json:"updated_at"`

	ClassID      uint     `gorm:"not null;index"    json:"class_id"      validate:"required"`
	DayOfWeek    string   `gorm:"size:20;not null"  json:"day_of_week"   validate:"required,oneof=Monday Tuesday Wednesday Thursday Friday Saturday"`
	PeriodNumber int      `gorm:"not null"          json:"period_number" validate:"required,min=1"`
	Subject      string   `gorm:"size:100;not null" json:"subject"       validate:"required"`
	TeacherID    uint     `gorm:"not null"          json:"teacher_id"    validate:"required"`
	Teacher      *Teacher `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
	StartTime    string   `gorm:"size:10;not null"  json:"start_time"    validate:"required"`
	EndTime      string   `gorm:"size:10;not null"  json:"end_time"      validate:"required"`
}

// TableName overrides the default table name.
func (Timetable) TableName() string {
	return "timetables"
}
