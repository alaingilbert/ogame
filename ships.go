package ogame

import (
	"math"
	"strconv"
)

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

// HasShips returns either or not at least one ship is present
func (s ShipsInfos) HasShips() bool {
	for _, ship := range Ships {
		if s.ByID(ship.GetID()) > 0 {
			return true
		}
	}
	return false
}

// Speed returns the speed of the slowest ship
func (s ShipsInfos) Speed(techs Researches) int {
	minSpeed := math.MaxInt32
	for _, ship := range Ships {
		nbr := s.ByID(ship.GetID())
		if nbr > 0 {
			shipSpeed := ship.GetSpeed(techs)
			if shipSpeed < minSpeed {
				minSpeed = shipSpeed
			}
		}
	}
	return minSpeed
}

// ToQuantifiables convert a ShipsInfos to an array of Quantifiable
func (s ShipsInfos) ToQuantifiables() []Quantifiable {
	out := make([]Quantifiable, 0)
	for _, ship := range Ships {
		if ship.GetID() == SolarSatelliteID {
			continue
		}
		shipID := ship.GetID()
		nbr := s.ByID(shipID)
		if nbr > 0 {
			out = append(out, Quantifiable{ID: shipID, Nbr: nbr})
		}
	}
	return out
}

// Cargo returns the total cargo of the ships
func (s ShipsInfos) Cargo(techs Researches) (out int) {
	for _, ship := range Ships {
		out += ship.GetCargoCapacity(techs) * s.ByID(ship.GetID())
	}
	return
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
	val += s.SolarSatellite * SolarSatellite.Price.Total()
	return val
}

// FleetCost returns the cost of the fleet
func (s ShipsInfos) FleetCost() Resources {
	val := LightFighter.Price.Mul(s.LightFighter)
	val = val.Add(HeavyFighter.Price.Mul(s.HeavyFighter))
	val = val.Add(Cruiser.Price.Mul(s.Cruiser))
	val = val.Add(Battleship.Price.Mul(s.Battleship))
	val = val.Add(Battlecruiser.Price.Mul(s.Battlecruiser))
	val = val.Add(Bomber.Price.Mul(s.Bomber))
	val = val.Add(Destroyer.Price.Mul(s.Destroyer))
	val = val.Add(Deathstar.Price.Mul(s.Deathstar))
	val = val.Add(SmallCargo.Price.Mul(s.SmallCargo))
	val = val.Add(LargeCargo.Price.Mul(s.LargeCargo))
	val = val.Add(ColonyShip.Price.Mul(s.ColonyShip))
	val = val.Add(Recycler.Price.Mul(s.Recycler))
	val = val.Add(EspionageProbe.Price.Mul(s.EspionageProbe))
	val = val.Add(SolarSatellite.Price.Mul(s.SolarSatellite))
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
