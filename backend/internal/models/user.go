package models

import "time"

type User struct {
	UserID         uint   `gorm:"primary_key;auto_increment"`
	Username       string `gorm:"size:50;uniqueIndex;not null"`
	HashedPassword string `gorm:"not null"`
	Role           string `gorm:"size:10;not null"`
	FullName       string `gorm:"size:100;not null;"`
	Email          string `gorm:"size:50;uniqueIndex;not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
