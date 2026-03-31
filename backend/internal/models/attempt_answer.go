package models

import (
	"time"

	"gorm.io/datatypes"
)

type AttemptAnswer struct {
	AttemptAnswerID uint           `gorm:"primaryKey"`
	AttemptID       uint           `gorm:"not null;index"`
	QuizQuestionID  uint           `gorm:"not null"`
	ChosenOptions   datatypes.JSON `gorm:"type:jsonb;not null"`
	AwardedMarks    int            `gorm:"default:0"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
