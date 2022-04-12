package models

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	AuthToken string
	ExpiresAt time.Time
	UserID    uint
	User      User `gorm:"foreignKey:UserID"`
	ClientIP  string
	UserAgent string
}
