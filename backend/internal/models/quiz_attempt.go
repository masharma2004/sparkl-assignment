package models

import "time"

type QuizAttempt struct {
	QuizAttemptID uint      `gorm:"primaryKey"`
	QuizID        uint      `gorm:"not null;index"`
	StudentID     uint      `gorm:"not null;index"`
	Status        string    `gorm:"size:20;not null"`
	StartedAt     time.Time `gorm:"not null"`
	SubmittedAt   *time.Time
	Score         int             `gorm:"default:0"`
	Answers       []AttemptAnswer `gorm:"foreignKey:AttemptID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
