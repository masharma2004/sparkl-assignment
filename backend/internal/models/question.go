package models

import (
	"time"

	"gorm.io/datatypes"
)

type Question struct {
	QuestionID     uint           `gorm:"primaryKey"`
	Category       string         `gorm:"size:80;not null;default:General"`
	Prompt         string         `gorm:"type:text;not null"`
	Options        datatypes.JSON `gorm:"type:jsonb;not null"`
	CorrectOptions datatypes.JSON `gorm:"type:jsonb;not null"`
	Solution       string         `gorm:"type:text"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
