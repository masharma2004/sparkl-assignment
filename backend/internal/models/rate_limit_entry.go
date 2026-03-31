package models

import "time"

type RateLimitEntry struct {
	RateLimitEntryID uint      `gorm:"primaryKey"`
	ScopeKey         string    `gorm:"size:191;uniqueIndex;not null"`
	HitCount         int       `gorm:"not null"`
	ResetAt          time.Time `gorm:"not null;index"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
