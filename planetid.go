package ogame

import "strconv"

// PlanetID ...
type PlanetID int

func (p PlanetID) String() string {
	return strconv.Itoa(int(p))
}
