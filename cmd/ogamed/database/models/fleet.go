package models

import (
	"github.com/faunX/ogame"
	"gorm.io/gorm"
)

type Fleet struct {
	gorm.Model
	ogame.Fleet `gorm:"embedded"`
}
