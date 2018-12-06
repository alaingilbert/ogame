package ogame

import (
	"math"
	"time"
)

// BaseDefender base for defender units (ships, defenses)
type BaseDefender struct {
	Base
	StructuralIntegrity int
	ShieldPower         int
	WeaponPower         int
	RapidfireFrom       map[ID]int
	Price               Resources
}

// GetStructuralIntegrity returns structural integrity of a defender unit
func (b BaseDefender) GetStructuralIntegrity(researches Researches) int {
	return int(float64(b.StructuralIntegrity) * (1 + float64(researches.ArmourTechnology)*0.1))
}

// GetShieldPower returns shield power of a defender unit
func (b BaseDefender) GetShieldPower(researches Researches) int {
	return int(float64(b.ShieldPower) * (1 + float64(researches.ShieldingTechnology)*0.1))
}

// GetWeaponPower returns weapon power of a defender unit
func (b BaseDefender) GetWeaponPower(researches Researches) int {
	return int(float64(b.WeaponPower) * (1 + float64(researches.WeaponsTechnology)*0.1))
}

// GetRapidfireFrom returns which ships have rapid fire against the defender unit
func (b BaseDefender) GetRapidfireFrom() map[ID]int {
	return b.RapidfireFrom
}

// ConstructionTime returns the duration it takes to build nbr defender units
func (b BaseDefender) ConstructionTime(nbr, universeSpeed int, facilities Facilities) time.Duration {
	shipyardLvl := float64(facilities.Shipyard)
	naniteLvl := float64(facilities.NaniteFactory)
	hours := float64(b.StructuralIntegrity) / (2500 * (1 + shipyardLvl) * float64(universeSpeed) * math.Pow(2, naniteLvl))
	secs := math.Max(1, hours*3600)
	return time.Duration(int(math.Floor(secs))*nbr) * time.Second
}

// GetPrice returns the price of nbr defender units
func (b BaseDefender) GetPrice(nbr int) Resources {
	return b.Price.Mul(nbr)
}
