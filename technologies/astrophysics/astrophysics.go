package astrophysics

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/technologies/baseTechnology"
)

// Astrophysics ...
type Astrophysics struct {
	baseTechnology.BaseTechnology
}

// New ...
func New() *Astrophysics {
	b := new(Astrophysics)
	b.OGameID = 124
	b.IncreaseFactor = 1.75
	b.BaseCost = ogame.Resources{Metal: 4000, Crystal: 8000, Deuterium: 4000}
	b.Requirements = map[ogame.ID]int{ogame.EspionageTechnology: 4, ogame.ImpulseDrive: 3}
	return b
}
