package ogame

import "strconv"

// PlanetID represent a planet id
type PlanetID CelestialID

func (p PlanetID) String() string {
	return strconv.Itoa(int(p))
}

// Celestial convert a PlanetID to a CelestialID
func (p PlanetID) Celestial() CelestialID {
	return CelestialID(p)
}
