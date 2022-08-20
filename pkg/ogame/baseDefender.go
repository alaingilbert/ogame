package ogame

import (
	"math"
	"time"
)

// BaseDefender base for defender units (ships, defenses)
type BaseDefender struct {
	Base
	StructuralIntegrity int64
	ShieldPower         int64
	WeaponPower         int64
	RapidfireFrom       map[ID]int64
	RapidfireAgainst    map[ID]int64
	Price               Resources
}

// GetStructuralIntegrity returns structural integrity of a defender unit
func (b BaseDefender) GetStructuralIntegrity(researches Researches) int64 {
	return int64(float64(b.StructuralIntegrity) * (1 + float64(researches.ArmourTechnology)*0.1))
}

// GetShieldPower returns shield power of a defender unit
func (b BaseDefender) GetShieldPower(researches Researches) int64 {
	return int64(float64(b.ShieldPower) * (1 + float64(researches.ShieldingTechnology)*0.1))
}

// GetWeaponPower returns weapon power of a defender unit
func (b BaseDefender) GetWeaponPower(researches Researches) int64 {
	return int64(float64(b.WeaponPower) * (1 + float64(researches.WeaponsTechnology)*0.1))
}

// GetRapidfireFrom returns which ships have rapid fire against the defender unit
func (b BaseDefender) GetRapidfireFrom() map[ID]int64 {
	return b.RapidfireFrom
}

// GetRapidfireAgainst returns which ships/defenses we have rapid fire against
func (b BaseDefender) GetRapidfireAgainst() map[ID]int64 {
	return b.RapidfireAgainst
}

// DefenderConstructionTime returns the duration it takes to build nbr defender units
func (b BaseDefender) DefenderConstructionTime(nbr, universeSpeed int64, acc DefenseAccelerators) time.Duration {
	shipyardLvl := float64(acc.GetShipyard())
	naniteLvl := float64(acc.GetNaniteFactory())
	hours := float64(b.StructuralIntegrity) / (2500 * (1 + shipyardLvl) * float64(universeSpeed) * math.Pow(2, naniteLvl))
	secs := math.Max(1, hours*3600)
	return time.Duration(int64(math.Floor(secs))*nbr) * time.Second
}

// ConstructionTime same as DefenderConstructionTime, needed for BaseOgameObj implementation
func (b BaseDefender) ConstructionTime(nbr, universeSpeed int64, acc BuildAccelerators, _, _ bool) time.Duration {
	return b.DefenderConstructionTime(nbr, universeSpeed, acc)
}

// GetPrice returns the price of nbr defender units
func (b BaseDefender) GetPrice(nbr int64) Resources {
	return b.Price.Mul(nbr)
}
