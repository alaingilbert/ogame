package crystalStorage

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/buildings/storageBuilding"
)

// CrystalStorage ...
type CrystalStorage struct {
	storageBuilding.StorageBuilding
}

// New ...
func New() *CrystalStorage {
	b := new(CrystalStorage)
	b.OGameID = 23
	b.IncreaseFactor = 2.0
	b.BaseCost = ogame.Resources{Metal: 1000, Crystal: 500}
	return b
}
