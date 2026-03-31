package models

import "time"

type QuizQuestion struct {
	QuizQuestionID uint `gorm:"primaryKey"`
	QuizID         uint `gorm:"not null;index"`
	QuestionID     uint `gorm:"not null"`
	SequenceNumber int  `gorm:"not null"`
	Marks          int  `gorm:"not null"`
	CreatedAt      time.Time
}
