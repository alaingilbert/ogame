package armourTechnology

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/technologies/baseTechnology"
)

// ArmourTechnology ...
type ArmourTechnology struct {
	baseTechnology.BaseTechnology
}

// New ...
func New() *ArmourTechnology {
	b := new(ArmourTechnology)
	b.OGameID = 111
	b.IncreaseFactor = 2.0
	b.BaseCost = ogame.Resources{Metal: 1000}
	b.Requirements = map[ogame.ID]int{ogame.ResearchLab: 2}
	return b
}
