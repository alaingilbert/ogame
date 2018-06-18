package computerTechnology

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/technologies/baseTechnology"
)

// ComputerTechnology ...
type ComputerTechnology struct {
	baseTechnology.BaseTechnology
}

// New ...
func New() *ComputerTechnology {
	b := new(ComputerTechnology)
	b.OGameID = 108
	b.IncreaseFactor = 2.0
	b.BaseCost = ogame.Resources{Crystal: 400, Deuterium: 600}
	b.Requirements = map[ogame.ID]int{ogame.ResearchLab: 1}
	return b
}
