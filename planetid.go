package ogame

// PlanetID represent a planet id
type PlanetID CelestialID

func (p PlanetID) String() string {
	return FI64(int64(p))
}

// Celestial convert a PlanetID to a CelestialID
func (p PlanetID) Celestial() CelestialID {
	return CelestialID(p)
}
