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
	Crawler        int
	Reaper         int
	Pathfinder     int
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

// Speed returns the speed of the slowest ship
func (s ShipsInfos) Speed(techs Researches) int {
	minSpeed := math.MaxInt32
	for _, ship := range Ships {
		if ship.GetID() == SolarSatelliteID {
			continue
		}
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
func (s ShipsInfos) Cargo(techs Researches, probeRaids bool) (out int) {
	for _, ship := range Ships {
		out += ship.GetCargoCapacity(techs, probeRaids) * s.ByID(ship.GetID())
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
func (s ShipsInfos) FleetValue() (out int) {
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

// Add adds two ShipsInfos together
func (s *ShipsInfos) Add(v ShipsInfos) {
	for _, ship := range Ships {
		shipID := ship.GetID()
		s.Set(shipID, s.ByID(shipID)+v.ByID(shipID))
	}
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
		"Solar Satellite: " + strconv.Itoa(s.SolarSatellite) + "\n" +
		"        Crawler: " + strconv.Itoa(s.Crawler) + "\n" +
		"         Reaper: " + strconv.Itoa(s.Reaper) + "\n" +
		"     Pathfinder: " + strconv.Itoa(s.Pathfinder)
}
