package models

import "gorm.io/gorm"

type Bot struct {
	gorm.Model
	ServerID uint
	Server   Server `gorm:"foreignKey:ServerID"`
	UserID   uint
	User     User `gorm:"foreignKey:UserID"`
}
