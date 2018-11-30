package ogame

import "strconv"

// ShipsInfos represent a planet ships information
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

// Has returns true if v is contained by s
func (s ShipsInfos) Has(v ShipsInfos) bool {
	for _, ship := range Ships {
		needed := v.ByID(ship.GetID())
		current := s.ByID(ship.GetID())
		if needed > 0 && needed > current {
			return false
		}
	}
	return true
}

// FleetValue returns the value of the fleet
func (s ShipsInfos) FleetValue() int {
	val := s.LightFighter * LightFighter.Price.Total()
	val += s.HeavyFighter * HeavyFighter.Price.Total()
	val += s.Cruiser * Cruiser.Price.Total()
	val += s.Battleship * Battleship.Price.Total()
	val += s.Battlecruiser * Battlecruiser.Price.Total()
	val += s.Bomber * Bomber.Price.Total()
	val += s.Destroyer * Destroyer.Price.Total()
	val += s.Deathstar * Deathstar.Price.Total()
	val += s.SmallCargo * SmallCargo.Price.Total()
	val += s.LargeCargo * LargeCargo.Price.Total()
	val += s.ColonyShip * ColonyShip.Price.Total()
	val += s.Recycler * Recycler.Price.Total()
	val += s.EspionageProbe * EspionageProbe.Price.Total()
	return val
}

// Add adds two ShipsInfos together
func (s *ShipsInfos) Add(v ShipsInfos) {
	s.LightFighter += v.LightFighter
	s.HeavyFighter += v.HeavyFighter
	s.Cruiser += v.Cruiser
	s.Battleship += v.Battleship
	s.Battlecruiser += v.Battlecruiser
	s.Bomber += v.Bomber
	s.Destroyer += v.Destroyer
	s.Deathstar += v.Deathstar
	s.SmallCargo += v.SmallCargo
	s.LargeCargo += v.LargeCargo
	s.ColonyShip += v.ColonyShip
	s.Recycler += v.Recycler
	s.EspionageProbe += v.EspionageProbe
	s.SolarSatellite += v.SolarSatellite
}

// ByID get number of ships by ship id
func (s ShipsInfos) ByID(id ID) int {
	switch id {
	case LightFighterID:
		return s.LightFighter
	case HeavyFighterID:
		return s.HeavyFighter
	case CruiserID:
		return s.Cruiser
	case BattleshipID:
		return s.Battleship
	case BattlecruiserID:
		return s.Battlecruiser
	case BomberID:
		return s.Bomber
	case DestroyerID:
		return s.Destroyer
	case DeathstarID:
		return s.Deathstar
	case SmallCargoID:
		return s.SmallCargo
	case LargeCargoID:
		return s.LargeCargo
	case ColonyShipID:
		return s.ColonyShip
	case RecyclerID:
		return s.Recycler
	case EspionageProbeID:
		return s.EspionageProbe
	case SolarSatelliteID:
		return s.SolarSatellite
	default:
		return 0
	}
}

// Set sets the ships value using the ship id
func (s *ShipsInfos) Set(id ID, val int) {
	switch id {
	case LightFighterID:
		s.LightFighter = val
	case HeavyFighterID:
		s.HeavyFighter = val
	case CruiserID:
		s.Cruiser = val
	case BattleshipID:
		s.Battleship = val
	case BattlecruiserID:
		s.Battlecruiser = val
	case BomberID:
		s.Bomber = val
	case DestroyerID:
		s.Destroyer = val
	case DeathstarID:
		s.Deathstar = val
	case SmallCargoID:
		s.SmallCargo = val
	case LargeCargoID:
		s.LargeCargo = val
	case ColonyShipID:
		s.ColonyShip = val
	case RecyclerID:
		s.Recycler = val
	case EspionageProbeID:
		s.EspionageProbe = val
	case SolarSatelliteID:
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
