package models

import "gorm.io/gorm"

type Script struct {
	gorm.Model
	Name   string
	Script string
}
