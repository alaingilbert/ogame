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

// ByID ...
func (s ShipsInfos) ByID(id ID) int {
	if id == LightFighter.ID {
		return s.LightFighter
	} else if id == HeavyFighter.ID {
		return s.HeavyFighter
	} else if id == Cruiser.ID {
		return s.Cruiser
	} else if id == Battleship.ID {
		return s.Battleship
	} else if id == Battlecruiser.ID {
		return s.Battlecruiser
	} else if id == Bomber.ID {
		return s.Bomber
	} else if id == Destroyer.ID {
		return s.Destroyer
	} else if id == Deathstar.ID {
		return s.Deathstar
	} else if id == SmallCargo.ID {
		return s.SmallCargo
	} else if id == LargeCargo.ID {
		return s.LargeCargo
	} else if id == ColonyShip.ID {
		return s.ColonyShip
	} else if id == Recycler.ID {
		return s.Recycler
	} else if id == EspionageProbe.ID {
		return s.EspionageProbe
	} else if id == SolarSatellite.ID {
		return s.SolarSatellite
	}
	return 0
}

// Set ...
func (s *ShipsInfos) Set(id ID, val int) {
	if id == LightFighter.ID {
		s.LightFighter = val
	} else if id == HeavyFighter.ID {
		s.HeavyFighter = val
	} else if id == Cruiser.ID {
		s.Cruiser = val
	} else if id == Battleship.ID {
		s.Battleship = val
	} else if id == Battlecruiser.ID {
		s.Battlecruiser = val
	} else if id == Bomber.ID {
		s.Bomber = val
	} else if id == Destroyer.ID {
		s.Destroyer = val
	} else if id == Deathstar.ID {
		s.Deathstar = val
	} else if id == SmallCargo.ID {
		s.SmallCargo = val
	} else if id == LargeCargo.ID {
		s.LargeCargo = val
	} else if id == ColonyShip.ID {
		s.ColonyShip = val
	} else if id == Recycler.ID {
		s.Recycler = val
	} else if id == EspionageProbe.ID {
		s.EspionageProbe = val
	} else if id == SolarSatellite.ID {
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
