package ogame

import "strconv"

// PlanetID represent a planet id
type PlanetID int

func (p PlanetID) String() string {
	return strconv.Itoa(int(p))
}
