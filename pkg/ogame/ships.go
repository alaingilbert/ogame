package ogame

import (
	"github.com/alaingilbert/ogame/pkg/utils"
	"iter"
	"math"
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
	return utils.Count2(s.Iter()) > 0
}

// IsEmpty returns true if no ships are set
func (s ShipsInfos) IsEmpty() bool {
	return !s.HasShips()
}

// HasFlyableShips returns either or not at least one flyable ship is present
func (s ShipsInfos) HasFlyableShips() bool {
	return utils.Any2(s.Iter(), func(shipID ID, i int64) bool { return shipID.IsFlyableShip() })
}

// HasCombatShips returns either or not at least one combat ship is present
func (s ShipsInfos) HasCombatShips() bool {
	return utils.Any2(s.Iter(), func(shipID ID, i int64) bool { return shipID.IsCombatShip() })
}

// HasCivilShips returns either or not at least one civil ship is present
func (s ShipsInfos) HasCivilShips() bool {
	return utils.Any2(s.Iter(), func(shipID ID, i int64) bool { return shipID.IsCivilShip() })
}

// Speed returns the speed of the slowest ship
func (s ShipsInfos) Speed(techs IResearches, lfBonuses LfBonuses, characterClass CharacterClass, allianceClass AllianceClass) int64 {
	var minSpeed int64 = math.MaxInt64
	for shipID := range s.IterFlyable() {
		shipSpeed := Objs.GetShip(shipID).GetSpeed(techs, lfBonuses, characterClass, allianceClass)
		minSpeed = min(minSpeed, shipSpeed)
	}
	return minSpeed
}

// ToQuantifiables convert a ShipsInfos to an array of Quantifiable
func (s ShipsInfos) ToQuantifiables() []Quantifiable {
	out := make([]Quantifiable, 0)
	for shipID, nb := range s.IterFlyable() {
		out = append(out, Quantifiable{ID: shipID, Nbr: nb})
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
func (s ShipsInfos) Cargo(techs IResearches, lfBonuses LfBonuses, characterClass CharacterClass, multiplier float64, probeRaids bool) (out int64) {
	for shipID, nb := range s.Iter() {
		out += Objs.GetShip(shipID).GetCargoCapacity(techs, lfBonuses, characterClass, multiplier, probeRaids) * nb
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
func (s ShipsInfos) FleetValue(lfBonuses LfBonuses) (out int64) {
	for shipID, nb := range s.Iter() {
		out += Objs.GetShip(shipID).GetPrice(nb, lfBonuses).Total()
	}
	return
}

// FleetCost returns the cost of the fleet
func (s ShipsInfos) FleetCost(lfBonuses LfBonuses) (out Resources) {
	for shipID, nb := range s.Iter() {
		out = out.Add(Objs.GetShip(shipID).GetPrice(nb, lfBonuses))
	}
	return
}

// CountShips returns the count of ships
func (s ShipsInfos) CountShips() (out int64) {
	for _, nb := range s.Iter() {
		out += nb
	}
	return
}

// Add adds two ShipsInfos together
func (s *ShipsInfos) Add(v ShipsInfos) {
	for _, ship := range Ships {
		shipID := ship.GetID()
		ownNb := s.ByID(shipID)
		otherNb := v.ByID(shipID)
		toSet := utils.Ternary(otherNb != -1, max(ownNb+otherNb, 0), -1)
		s.Set(shipID, toSet)
	}
}

// Sub subtracts v from s
func (s *ShipsInfos) Sub(v ShipsInfos) {
	for _, ship := range Ships {
		shipID := ship.GetID()
		s.Set(shipID, max(s.ByID(shipID)-v.ByID(shipID), 0))
	}
}

// AddShips adds some ships
func (s *ShipsInfos) AddShips(shipID ID, nb int64) {
	s.Set(shipID, max(s.ByID(shipID)+nb, 0))
}

// SubShips subtracts some ships
func (s *ShipsInfos) SubShips(shipID ID, nb int64) {
	s.AddShips(shipID, -1*nb)
}

func (s ShipsInfos) each(clb func(shipID ID, nb int64) bool) {
	for _, ship := range Ships {
		shipID := ship.GetID()
		nb := s.ByID(shipID)
		if nb > 0 {
			if !clb(shipID, nb) {
				return
			}
		}
	}
}

func (s ShipsInfos) eachFlyable(clb func(shipID ID, nb int64) bool) {
	for shipID, nb := range s.Iter() {
		if shipID.IsFlyableShip() {
			if !clb(shipID, nb) {
				return
			}
		}
	}
}

// Each calls clb callback for every ships that has a value higher than zero
func (s ShipsInfos) Each(clb func(shipID ID, nb int64)) {
	s.each(func(shipID ID, nb int64) bool {
		clb(shipID, nb)
		return true
	})
}

// EachFlyable calls clb callback for every ships that has a value higher than zero and is flyable
func (s ShipsInfos) EachFlyable(clb func(shipID ID, nb int64)) {
	s.eachFlyable(func(shipID ID, nb int64) bool {
		clb(shipID, nb)
		return true
	})
}

// Iter implements iterator so that we can use in a for loop and avoid having to deal with closure
func (s ShipsInfos) Iter() iter.Seq2[ID, int64] {
	return func(yield func(shipID ID, nb int64) bool) {
		s.each(yield)
	}
}

// IterFlyable implements iterator so that we can use in a for loop and avoid having to deal with closure
func (s ShipsInfos) IterFlyable() iter.Seq2[ID, int64] {
	return func(yield func(shipID ID, nb int64) bool) {
		s.eachFlyable(yield)
	}
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

// ByShip get number of ships given a "Ship"
func (s ShipsInfos) ByShip(ship Ship) int64 {
	return s.ByID(ship.GetID())
}

// Get gets number of ships
func (s ShipsInfos) Get(v any) int64 {
	switch vv := v.(type) {
	case ID:
		return s.ByID(vv)
	case Ship:
		return s.ByShip(vv)
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

// SetShip sets the ships value using the "Ship" id
func (s *ShipsInfos) SetShip(ship Ship, val int64) {
	s.Set(ship.GetID(), val)
}

func (s ShipsInfos) String() string {
	return "\n" +
		"  Light Fighter: " + utils.FI64(s.LightFighter) + "\n" +
		"  Heavy Fighter: " + utils.FI64(s.HeavyFighter) + "\n" +
		"        Cruiser: " + utils.FI64(s.Cruiser) + "\n" +
		"     Battleship: " + utils.FI64(s.Battleship) + "\n" +
		"  Battlecruiser: " + utils.FI64(s.Battlecruiser) + "\n" +
		"         Bomber: " + utils.FI64(s.Bomber) + "\n" +
		"      Destroyer: " + utils.FI64(s.Destroyer) + "\n" +
		"      Deathstar: " + utils.FI64(s.Deathstar) + "\n" +
		"    Small Cargo: " + utils.FI64(s.SmallCargo) + "\n" +
		"    Large Cargo: " + utils.FI64(s.LargeCargo) + "\n" +
		"    Colony Ship: " + utils.FI64(s.ColonyShip) + "\n" +
		"       Recycler: " + utils.FI64(s.Recycler) + "\n" +
		"Espionage Probe: " + utils.FI64(s.EspionageProbe) + "\n" +
		"Solar Satellite: " + utils.FI64(s.SolarSatellite) + "\n" +
		"        Crawler: " + utils.FI64(s.Crawler) + "\n" +
		"         Reaper: " + utils.FI64(s.Reaper) + "\n" +
		"     Pathfinder: " + utils.FI64(s.Pathfinder)
}
