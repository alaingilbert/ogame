package models

import "github.com/faunX/ogame"

type Test struct {
	ID             uint
	MyType         `gorm:"-"`
	ogame.PlanetID `gorm:"type:uint"`
}

type MyType int64
