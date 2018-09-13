package ogame

import (
	"math"
	"time"
)

// BaseDefender ...
type BaseDefender struct {
	Base
	StructuralIntegrity int
	ShieldPower         int
	WeaponPower         int
	RapidfireFrom       map[ID]int
	Price               Resources
}

// GetStructuralIntegrity ...
func (b BaseDefender) GetStructuralIntegrity(researches Researches) int {
	return int(float64(b.StructuralIntegrity) * (1 + float64(researches.ArmourTechnology)*0.1))
}

// GetShieldPower ...
func (b BaseDefender) GetShieldPower(researches Researches) int {
	return int(float64(b.ShieldPower) * (1 + float64(researches.ShieldingTechnology)*0.1))
}

// GetWeaponPower ...
func (b BaseDefender) GetWeaponPower(researches Researches) int {
	return int(float64(b.WeaponPower) * (1 + float64(researches.WeaponsTechnology)*0.1))
}

// GetRapidfireFrom ...
func (b BaseDefender) GetRapidfireFrom() map[ID]int {
	return b.RapidfireFrom
}

// ConstructionTime ...
func (b BaseDefender) ConstructionTime(nbr, universeSpeed int, facilities Facilities) time.Duration {
	shipyardLvl := float64(facilities.Shipyard)
	naniteLvl := float64(facilities.NaniteFactory)
	hours := float64(b.StructuralIntegrity) / (2500 * (1 + shipyardLvl) * float64(universeSpeed) * math.Pow(2, naniteLvl))
	secs := hours * 3600
	return time.Duration(int(math.Floor(secs))*nbr) * time.Second
}

// GetPrice ...
func (b BaseDefender) GetPrice(nbr int) Resources {
	return b.Price.Mul(nbr)
}
