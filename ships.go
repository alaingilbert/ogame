package ogame

import "strconv"

// Ships ...
type Ships struct {
	LightFighter   int
	HeavyFighter   int
	Cruiser        int
	Battleship     int
	Battlecruiser  int
	Bomber         int
	Destroyer      int
	Deathstar      int
	SmallCargo     int
	LargeCargo     int
	ColonyShip     int
	Recycler       int
	EspionageProbe int
	SolarSatellite int
}

// ByOGameID ...
func (s Ships) ByOGameID(ogameID ID) int {
	if ogameID == LightFighter {
		return s.LightFighter
	} else if ogameID == HeavyFighter {
		return s.HeavyFighter
	} else if ogameID == Cruiser {
		return s.Cruiser
	} else if ogameID == Battleship {
		return s.Battleship
	} else if ogameID == Battlecruiser {
		return s.Battlecruiser
	} else if ogameID == Bomber {
		return s.Bomber
	} else if ogameID == Destroyer {
		return s.Destroyer
	} else if ogameID == Deathstar {
		return s.Deathstar
	} else if ogameID == SmallCargo {
		return s.SmallCargo
	} else if ogameID == LargeCargo {
		return s.LargeCargo
	} else if ogameID == ColonyShip {
		return s.ColonyShip
	} else if ogameID == Recycler {
		return s.Recycler
	} else if ogameID == EspionageProbe {
		return s.EspionageProbe
	} else if ogameID == SolarSatellite {
		return s.SolarSatellite
	}
	return 0
}

// Set ...
func (s *Ships) Set(ogameID ID, val int) {
	if ogameID == LightFighter {
		s.LightFighter = val
	} else if ogameID == HeavyFighter {
		s.HeavyFighter = val
	} else if ogameID == Cruiser {
		s.Cruiser = val
	} else if ogameID == Battleship {
		s.Battleship = val
	} else if ogameID == Battlecruiser {
		s.Battlecruiser = val
	} else if ogameID == Bomber {
		s.Bomber = val
	} else if ogameID == Destroyer {
		s.Destroyer = val
	} else if ogameID == Deathstar {
		s.Deathstar = val
	} else if ogameID == SmallCargo {
		s.SmallCargo = val
	} else if ogameID == LargeCargo {
		s.LargeCargo = val
	} else if ogameID == ColonyShip {
		s.ColonyShip = val
	} else if ogameID == Recycler {
		s.Recycler = val
	} else if ogameID == EspionageProbe {
		s.EspionageProbe = val
	} else if ogameID == SolarSatellite {
		s.SolarSatellite = val
	}
}

func (s Ships) String() string {
	return "\n" +
		"  Light Fighter: " + strconv.Itoa(s.LightFighter) + "\n" +
		"  Heavy Fighter: " + strconv.Itoa(s.HeavyFighter) + "\n" +
		"        Cruiser: " + strconv.Itoa(s.Cruiser) + "\n" +
		"     Battleship: " + strconv.Itoa(s.Battleship) + "\n" +
		"  Battlecruiser: " + strconv.Itoa(s.Battlecruiser) + "\n" +
		"         Bomber: " + strconv.Itoa(s.Bomber) + "\n" +
		"      Destroyer: " + strconv.Itoa(s.Destroyer) + "\n" +
		"      Deathstar: " + strconv.Itoa(s.Deathstar) + "\n" +
		"    Small Cargo: " + strconv.Itoa(s.SmallCargo) + "\n" +
		"    Large Cargo: " + strconv.Itoa(s.LargeCargo) + "\n" +
		"    Colony Ship: " + strconv.Itoa(s.ColonyShip) + "\n" +
		"       Recycler: " + strconv.Itoa(s.Recycler) + "\n" +
		"Espionage Probe: " + strconv.Itoa(s.EspionageProbe) + "\n" +
		"Solar Satellite: " + strconv.Itoa(s.SolarSatellite)
}
