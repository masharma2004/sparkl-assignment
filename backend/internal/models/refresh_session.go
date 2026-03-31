package models

import "time"

type RefreshSession struct {
	RefreshSessionID uint      `gorm:"primaryKey"`
	UserID           uint      `gorm:"not null;index"`
	TokenHash        string    `gorm:"size:64;uniqueIndex;not null"`
	ExpiresAt        time.Time `gorm:"not null;index"`
	RevokedAt        *time.Time
	ReplacedByHash   string `gorm:"size:64"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
