package hyperspaceDrive

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/technologies/baseTechnology"
)

// HyperspaceDrive ...
type HyperspaceDrive struct {
	baseTechnology.BaseTechnology
}

// New ...
func New() *HyperspaceDrive {
	b := new(HyperspaceDrive)
	b.OGameID = 118
	b.IncreaseFactor = 2.0
	b.BaseCost = ogame.Resources{Metal: 10000, Crystal: 20000, Deuterium: 6000}
	b.Requirements = map[ogame.ID]int{ogame.HyperspaceTechnology: 3}
	return b
}
