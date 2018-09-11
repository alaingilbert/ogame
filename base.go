package ogame

import (
	"math"
	"time"
)

// BaseShip ...
type BaseShip struct {
	ID                  ID
	Name                string
	StructuralIntegrity int
	ShieldPower         int
	WeaponPower         int
	CargoCapacity       int
	BaseSpeed           int
	FuelConsumption     int
	RapidfireFrom       map[ID]int
	RapidfireAgainst    map[ID]int
	Requirements        map[ID]int
	Price               Resources
}

// GetID ...
func (b BaseShip) GetID() ID {
	return b.ID
}

// GetName ...
func (b BaseShip) GetName() string {
	return b.Name
}

// GetStructuralIntegrity ...
func (b BaseShip) GetStructuralIntegrity(researches Researches) int {
	return int(float64(b.StructuralIntegrity) * (1 + float64(researches.ArmourTechnology)*0.1))
}

// GetShieldPower ...
func (b BaseShip) GetShieldPower(researches Researches) int {
	return int(float64(b.ShieldPower) * (1 + float64(researches.ShieldingTechnology)*0.1))
}

// GetWeaponPower ...
func (b BaseShip) GetWeaponPower(researches Researches) int {
	return int(float64(b.WeaponPower) * (1 + float64(researches.WeaponsTechnology)*0.1))
}

// GetCargoCapacity ...
func (b BaseShip) GetCargoCapacity() int {
	return b.CargoCapacity
}

// GetBaseSpeed ...
func (b BaseShip) GetBaseSpeed() int {
	return b.BaseSpeed
}

// GetSpeed ...
func (b BaseShip) GetSpeed(techs Researches) int {
	techDriveLvl := 0
	if b.ID == SmallCargoID && techs.ImpulseDrive >= 5 {
		return int(float64(b.BaseSpeed) + (float64(b.BaseSpeed)*0.2)*float64(techs.ImpulseDrive))
	}
	if minLvl, ok := b.Requirements[CombustionDrive.ID]; ok {
		techDriveLvl = techs.CombustionDrive
		if techDriveLvl < minLvl {
			techDriveLvl = minLvl
		}
		return int(float64(b.BaseSpeed) + (float64(b.BaseSpeed)*0.1)*float64(techDriveLvl))
	} else if minLvl, ok := b.Requirements[ImpulseDrive.ID]; ok {
		techDriveLvl = techs.ImpulseDrive
		if techDriveLvl < minLvl {
			techDriveLvl = minLvl
		}
		return int(float64(b.BaseSpeed) + (float64(b.BaseSpeed)*0.2)*float64(techDriveLvl))
	} else if minLvl, ok := b.Requirements[HyperspaceDrive.ID]; ok {
		techDriveLvl = techs.HyperspaceDrive
		if techDriveLvl < minLvl {
			techDriveLvl = minLvl
		}
		return int(float64(b.BaseSpeed) + (float64(b.BaseSpeed)*0.3)*float64(techDriveLvl))
	}
	return int(float64(b.BaseSpeed) + (float64(b.BaseSpeed)*0.2)*float64(techDriveLvl))
}

// GetFuelConsumption ...
func (b BaseShip) GetFuelConsumption() int {
	return b.FuelConsumption
}

// GetRapidfireFrom ...
func (b BaseShip) GetRapidfireFrom() map[ID]int {
	return b.RapidfireFrom
}

// GetRapidfireAgainst ...
func (b BaseShip) GetRapidfireAgainst() map[ID]int {
	return b.RapidfireAgainst
}

// GetPrice ...
func (b BaseShip) GetPrice(nbr int) Resources {
	return b.Price.Mul(nbr)
}

// ConstructionTime ...
func (b BaseShip) ConstructionTime(nbr, universeSpeed int, facilities Facilities) time.Duration {
	shipyardLvl := float64(facilities.Shipyard)
	naniteLvl := float64(facilities.NaniteFactory)
	hours := float64(b.StructuralIntegrity) / (2500 * (1 + shipyardLvl) * float64(universeSpeed) * math.Pow(2, naniteLvl))
	secs := hours * 3600
	return time.Duration(int(math.Floor(secs))*nbr) * time.Second
}

// GetRequirements ...
func (b BaseShip) GetRequirements() map[ID]int {
	return b.Requirements
}

// IsAvailable ...
func (b BaseShip) IsAvailable(_ ResourcesBuildings, facilities Facilities, researches Researches, _ int) bool {
	for id, levelNeeded := range b.Requirements {
		if id.IsFacility() {
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
