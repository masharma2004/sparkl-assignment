package models

import "time"

type Quiz struct {
	QuizID          uint           `gorm:"primary_key;AUTO_INCREMENT"`
	Category        string         `gorm:"size:80;not null;default:General"`
	Title           string         `gorm:"size:250;not null"`
	QuestionCount   int            `gorm:"not null"`
	TotalMarks      int            `gorm:"not null"`
	DurationMinutes int            `gorm:"not null"`
	CreatedBy       uint           `gorm:"not null"`
	QuizQuestions   []QuizQuestion `gorm:"foreignkey:QuizID"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
