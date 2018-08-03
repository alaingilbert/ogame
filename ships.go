package ogame

import "strconv"

// ShipsInfos ...
type ShipsInfos struct {
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
func (s ShipsInfos) ByOGameID(ogameID ID) int {
	if ogameID == LightFighter.ID {
		return s.LightFighter
	} else if ogameID == HeavyFighter.ID {
		return s.HeavyFighter
	} else if ogameID == Cruiser.ID {
		return s.Cruiser
	} else if ogameID == Battleship.ID {
		return s.Battleship
	} else if ogameID == Battlecruiser.ID {
		return s.Battlecruiser
	} else if ogameID == Bomber.ID {
		return s.Bomber
	} else if ogameID == Destroyer.ID {
		return s.Destroyer
	} else if ogameID == Deathstar.ID {
		return s.Deathstar
	} else if ogameID == SmallCargo.ID {
		return s.SmallCargo
	} else if ogameID == LargeCargo.ID {
		return s.LargeCargo
	} else if ogameID == ColonyShip.ID {
		return s.ColonyShip
	} else if ogameID == Recycler.ID {
		return s.Recycler
	} else if ogameID == EspionageProbe.ID {
		return s.EspionageProbe
	} else if ogameID == SolarSatellite.ID {
		return s.SolarSatellite
	}
	return 0
}

// Set ...
func (s *ShipsInfos) Set(ogameID ID, val int) {
	if ogameID == LightFighter.ID {
		s.LightFighter = val
	} else if ogameID == HeavyFighter.ID {
		s.HeavyFighter = val
	} else if ogameID == Cruiser.ID {
		s.Cruiser = val
	} else if ogameID == Battleship.ID {
		s.Battleship = val
	} else if ogameID == Battlecruiser.ID {
		s.Battlecruiser = val
	} else if ogameID == Bomber.ID {
		s.Bomber = val
	} else if ogameID == Destroyer.ID {
		s.Destroyer = val
	} else if ogameID == Deathstar.ID {
		s.Deathstar = val
	} else if ogameID == SmallCargo.ID {
		s.SmallCargo = val
	} else if ogameID == LargeCargo.ID {
		s.LargeCargo = val
	} else if ogameID == ColonyShip.ID {
		s.ColonyShip = val
	} else if ogameID == Recycler.ID {
		s.Recycler = val
	} else if ogameID == EspionageProbe.ID {
		s.EspionageProbe = val
	} else if ogameID == SolarSatellite.ID {
		s.SolarSatellite = val
	}
}

func (s ShipsInfos) String() string {
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
