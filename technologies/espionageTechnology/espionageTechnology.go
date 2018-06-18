package espionageTechnology

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/technologies/baseTechnology"
)

// EspionageTechnology ...
type EspionageTechnology struct {
	baseTechnology.BaseTechnology
}

// New ...
func New() *EspionageTechnology {
	b := new(EspionageTechnology)
	b.OGameID = 106
	b.IncreaseFactor = 2.0
	b.BaseCost = ogame.Resources{Metal: 200, Crystal: 1000, Deuterium: 200}
	b.Requirements = map[ogame.ID]int{ogame.ResearchLab: 3}
	return b
}
