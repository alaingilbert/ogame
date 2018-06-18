package intergalacticResearchNetwork

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/technologies/baseTechnology"
)

// IntergalacticResearchNetwork ...
type IntergalacticResearchNetwork struct {
	baseTechnology.BaseTechnology
}

// New ...
func New() *IntergalacticResearchNetwork {
	b := new(IntergalacticResearchNetwork)
	b.OGameID = 123
	b.IncreaseFactor = 2.0
	b.BaseCost = ogame.Resources{Metal: 240000, Crystal: 400000, Deuterium: 160000}
	b.Requirements = map[ogame.ID]int{ogame.ResearchLab: 10, ogame.ComputerTechnology: 8, ogame.HyperspaceTechnology: 8}
	return b
}
