package ogame

import "strconv"

// PlanetID represent a planet id
type PlanetID CelestialID

func (p PlanetID) String() string {
	return strconv.FormatInt(int64(p), 10)
}

// Celestial convert a PlanetID to a CelestialID
func (p PlanetID) Celestial() CelestialID {
	return CelestialID(p)
}
