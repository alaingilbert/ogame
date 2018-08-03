package ogame

import "math"

// StorageBuilding ...
type StorageBuilding struct {
	BaseBuilding
}

// Capacity ...
func (s StorageBuilding) Capacity(lvl int) int {
	return 5000 * int(2.5*math.Pow(math.E, (20*float64(lvl))/33))
}
