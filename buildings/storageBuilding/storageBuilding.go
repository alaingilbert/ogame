package storageBuilding

import (
	"math"

	"github.com/alaingilbert/ogame/buildings/baseBuilding"
)

// StorageBuilding ...
type StorageBuilding struct {
	baseBuilding.BaseBuilding
}

// Capacity ...
func (s StorageBuilding) Capacity(lvl int) int {
	return 5000 * int(2.5*math.Pow(math.E, (20*float64(lvl))/33))
}
