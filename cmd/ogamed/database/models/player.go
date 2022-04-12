package models

import (
	"github.com/faunX/ogame"
	"gorm.io/gorm"
)

type Player struct {
	gorm.Model
	ogame.UserInfos
}
