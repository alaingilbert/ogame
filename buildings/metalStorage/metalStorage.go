package metalStorage

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/buildings/storageBuilding"
)

// MetalStorage ...
type MetalStorage struct {
	storageBuilding.StorageBuilding
}

// New ...
func New() *MetalStorage {
	b := new(MetalStorage)
	b.OGameID = 22
	b.IncreaseFactor = 2.0
	b.BaseCost = ogame.Resources{Metal: 1000}
	return b
}
