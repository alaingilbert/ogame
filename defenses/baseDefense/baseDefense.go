package baseDefense

import (
	"math"

	"github.com/alaingilbert/ogame"
)

// BaseDefense ...
type BaseDefense struct {
	OGameID             ogame.ID
	Price               ogame.Resources
	StructuralIntegrity int
	ShieldPower         int
	WeaponPower         int
	RapidfireFrom       map[ogame.ID]int
	Requirements        map[ogame.ID]int
}

// GetOGameID ...
func (b BaseDefense) GetOGameID() ogame.ID {
	return b.OGameID
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
func (b BaseDefense) GetRapidfireFrom() map[ogame.ID]int {
	return b.RapidfireFrom
}

// GetPrice ...
func (b BaseDefense) GetPrice(int) ogame.Resources {
	return b.Price
}

// ConstructionTime ...
func (b BaseDefense) ConstructionTime(nbr, universeSpeed int, facilities ogame.Facilities) int {
	shipyardLvl := float64(facilities.Shipyard)
	naniteLvl := float64(facilities.NaniteFactory)
	hours := float64(b.StructuralIntegrity) / (2500 * (1 + shipyardLvl) * float64(universeSpeed) * math.Pow(2, naniteLvl))
	secs := hours * 3600
	return int(math.Floor(secs)) * nbr
}

// GetRequirements ...
func (b BaseDefense) GetRequirements() map[ogame.ID]int {
	return b.Requirements
}

// IsAvailable ...
func (b BaseDefense) IsAvailable(resourcesBuildings ogame.ResourcesBuildings, facilities ogame.Facilities, researches ogame.Researches, _ int) bool {
	for ogameID, levelNeeded := range b.Requirements {
		if ogameID.IsResourceBuilding() {
			if resourcesBuildings.ByOGameID(ogameID) < levelNeeded {
				return false
			}
		} else if ogameID.IsFacility() {
			if facilities.ByOGameID(ogameID) < levelNeeded {
				return false
			}
		} else if ogameID.IsTech() {
			if researches.ByOGameID(ogameID) < levelNeeded {
				return false
			}
		}
	}
	return true
}
