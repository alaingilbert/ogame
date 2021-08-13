package ogame

import (
	"math"
	"strconv"
)

// ShipsInfos represent a planet ships information
type ShipsInfos struct {
	LightFighter   int64 // 204
	HeavyFighter   int64 // 205
	Cruiser        int64 // 206
	Battleship     int64 // 207
	Battlecruiser  int64 // 215
	Bomber         int64 // 211
	Destroyer      int64 // 213
	Deathstar      int64 // 214
	SmallCargo     int64 // 202
	LargeCargo     int64 // 203
	ColonyShip     int64 // 208
	Recycler       int64 // 209
	EspionageProbe int64 // 210
	SolarSatellite int64 // 212
	Crawler        int64 // 217
	Reaper         int64 // 218
	Pathfinder     int64 // 219
}

// ToPtr returns a pointer to self
func (s ShipsInfos) ToPtr() *ShipsInfos {
	return &s
}

// Equal either or not two ShipsInfos are equal
func (s ShipsInfos) Equal(other ShipsInfos) bool {
	for _, ship := range Ships {
		shipID := ship.GetID()
		if s.ByID(shipID) != other.ByID(shipID) {
			return false
		}
	}
	return true
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

// HasMovableShips returns either or not at least one ship that can be moved is present
func (s ShipsInfos) HasMovableShips() bool {
	for _, ship := range Ships {
		if ship.GetID() != SolarSatelliteID || ship.GetID() != CrawlerID {

		}
		switch ship.GetID() {
		case SolarSatelliteID:
		case CrawlerID:
		default:
			if s.ByID(ship.GetID()) > 0 {
				return true
			}
		}
	}
	return false
}

// Speed returns the speed of the slowest ship
func (s ShipsInfos) Speed(techs Researches, isCollector, isGeneral bool) int64 {
	var minSpeed int64 = math.MaxInt64
	for _, ship := range Ships {
		if ship.GetID() == SolarSatelliteID {
			continue
		}
		nbr := s.ByID(ship.GetID())
		if nbr > 0 {
			shipSpeed := ship.GetSpeed(techs, isCollector, isGeneral)
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
		if ship.GetID() == SolarSatelliteID || ship.GetID() == CrawlerID {
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

// FromQuantifiables convert an array of Quantifiable to a ShipsInfos
func (s ShipsInfos) FromQuantifiables(in []Quantifiable) (out ShipsInfos) {
	for _, item := range in {
		out.Set(item.ID, item.Nbr)
	}
	return
}

// Cargo returns the total cargo of the ships
func (s ShipsInfos) Cargo(techs Researches, probeRaids, isCollector, isPioneers bool) (out int64) {
	for _, ship := range Ships {
		out += ship.GetCargoCapacity(techs, probeRaids, isCollector, isPioneers) * s.ByID(ship.GetID())
	}
	return
}

// FuelCapacity returns the total Capacity for fuel
func (s ShipsInfos) FuelCapacity() (out int64) {
	for _, ship := range Ships {
		out += ship.GetFuelCapacity() * s.ByID(ship.GetID())
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
func (s ShipsInfos) FleetValue() (out int64) {
	for _, ship := range Ships {
		out += ship.GetPrice(s.ByID(ship.GetID())).Total()
	}
	return
}

// FleetCost returns the cost of the fleet
func (s ShipsInfos) FleetCost() (out Resources) {
	for _, ship := range Ships {
		out = out.Add(ship.GetPrice(s.ByID(ship.GetID())))
	}
	return
}

// CountShips returns the count of ships
func (s ShipsInfos) CountShips() (out int64) {
	for _, ship := range Ships {
		out += s.ByID(ship.GetID())
	}
	return
}

// Add adds two ShipsInfos together
func (s *ShipsInfos) Add(v ShipsInfos) {
	for _, ship := range Ships {
		shipID := ship.GetID()
		s.Set(shipID, s.ByID(shipID)+v.ByID(shipID))
	}
}

// AddShips adds some ships
func (s *ShipsInfos) AddShips(shipID ID, nb int64) {
	s.Set(shipID, s.ByID(shipID)+nb)
}

// SubShips subtracts some ships
func (s *ShipsInfos) SubShips(shipID ID, nb int64) {
	s.AddShips(shipID, -1*nb)
}

// ByID get number of ships by ship id
func (s ShipsInfos) ByID(id ID) int64 {
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
	case CrawlerID:
		return s.Crawler
	case ReaperID:
		return s.Reaper
	case PathfinderID:
		return s.Pathfinder
	default:
		return 0
	}
}

// Set sets the ships value using the ship id
func (s *ShipsInfos) Set(id ID, val int64) {
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
	case CrawlerID:
		s.Crawler = val
	case ReaperID:
		s.Reaper = val
	case PathfinderID:
		s.Pathfinder = val
	}
}

func (s ShipsInfos) String() string {
	return "\n" +
		"  Light Fighter: " + strconv.FormatInt(s.LightFighter, 10) + "\n" +
		"  Heavy Fighter: " + strconv.FormatInt(s.HeavyFighter, 10) + "\n" +
		"        Cruiser: " + strconv.FormatInt(s.Cruiser, 10) + "\n" +
		"     Battleship: " + strconv.FormatInt(s.Battleship, 10) + "\n" +
		"  Battlecruiser: " + strconv.FormatInt(s.Battlecruiser, 10) + "\n" +
		"         Bomber: " + strconv.FormatInt(s.Bomber, 10) + "\n" +
		"      Destroyer: " + strconv.FormatInt(s.Destroyer, 10) + "\n" +
		"      Deathstar: " + strconv.FormatInt(s.Deathstar, 10) + "\n" +
		"    Small Cargo: " + strconv.FormatInt(s.SmallCargo, 10) + "\n" +
		"    Large Cargo: " + strconv.FormatInt(s.LargeCargo, 10) + "\n" +
		"    Colony Ship: " + strconv.FormatInt(s.ColonyShip, 10) + "\n" +
		"       Recycler: " + strconv.FormatInt(s.Recycler, 10) + "\n" +
		"Espionage Probe: " + strconv.FormatInt(s.EspionageProbe, 10) + "\n" +
		"Solar Satellite: " + strconv.FormatInt(s.SolarSatellite, 10) + "\n" +
		"        Crawler: " + strconv.FormatInt(s.Crawler, 10) + "\n" +
		"         Reaper: " + strconv.FormatInt(s.Reaper, 10) + "\n" +
		"     Pathfinder: " + strconv.FormatInt(s.Pathfinder, 10)
}
