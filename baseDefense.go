package ogame

import (
	"math"
	"time"
)

// BaseDefense ...
type BaseDefense struct {
	ID                  ID
	Name                string
	Price               Resources
	StructuralIntegrity int
	ShieldPower         int
	WeaponPower         int
	RapidfireFrom       map[ID]int
	Requirements        map[ID]int
}

// GetID ...
func (b BaseDefense) GetID() ID {
	return b.ID
}

// GetName ...
func (b BaseDefense) GetName() string {
	return b.Name
}

// GetStructuralIntegrity ...
func (b BaseDefense) GetStructuralIntegrity() int {
	return b.StructuralIntegrity
}

// GetShieldPower ...
func (b BaseDefense) GetShieldPower() int {
	return b.ShieldPower
}

// GetWeaponPower ...
func (b BaseDefense) GetWeaponPower() int {
	return b.WeaponPower
}

// GetRapidfireFrom ...
func (b BaseDefense) GetRapidfireFrom() map[ID]int {
	return b.RapidfireFrom
}

// GetPrice ...
func (b BaseDefense) GetPrice(int) Resources {
	return b.Price
}

// ConstructionTime ...
func (b BaseDefense) ConstructionTime(nbr, universeSpeed int, facilities Facilities) time.Duration {
	shipyardLvl := float64(facilities.Shipyard)
	naniteLvl := float64(facilities.NaniteFactory)
	hours := float64(b.StructuralIntegrity) / (2500 * (1 + shipyardLvl) * float64(universeSpeed) * math.Pow(2, naniteLvl))
	secs := hours * 3600
	return time.Duration(int(math.Floor(secs))*nbr) * time.Second
}

// GetRequirements ...
func (b BaseDefense) GetRequirements() map[ID]int {
	return b.Requirements
}

// IsAvailable ...
func (b BaseDefense) IsAvailable(resourcesBuildings ResourcesBuildings, facilities Facilities, researches Researches, _ int) bool {
	for id, levelNeeded := range b.Requirements {
		if id.IsResourceBuilding() {
			if resourcesBuildings.ByID(id) < levelNeeded {
				return false
			}
		} else if id.IsFacility() {
			if facilities.ByID(id) < levelNeeded {
				return false
			}
		} else if id.IsTech() {
			if researches.ByID(id) < levelNeeded {
				return false
			}
		}
	}
	return true
}
