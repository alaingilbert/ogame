package deuteriumTank

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/buildings/storageBuilding"
)

// DeuteriumTank ...
type DeuteriumTank struct {
	storageBuilding.StorageBuilding
}

// New ...
func New() *DeuteriumTank {
	b := new(DeuteriumTank)
	b.OGameID = 24
	b.IncreaseFactor = 2.0
	b.BaseCost = ogame.Resources{Metal: 1000, Crystal: 1000}
	return b
}
