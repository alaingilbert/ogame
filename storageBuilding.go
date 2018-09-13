package ogame

import "math"

type storageBuilding struct {
	BaseBuilding
}

// Capacity returns the capacity of a storage building
func (s storageBuilding) Capacity(lvl int) int {
	return 5000 * int(2.5*math.Pow(math.E, (20*float64(lvl))/33))
}
