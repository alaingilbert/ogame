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
func (b BaseDefender) GetStructuralIntegrity(researches IResearches) int64 {
	return int64(float64(b.StructuralIntegrity) * (1 + float64(researches.GetArmourTechnology())*0.1))
}

// GetShieldPower returns shield power of a defender unit
func (b BaseDefender) GetShieldPower(researches IResearches) int64 {
	return int64(float64(b.ShieldPower) * (1 + float64(researches.GetShieldingTechnology())*0.1))
}

// GetWeaponPower returns weapon power of a defender unit
func (b BaseDefender) GetWeaponPower(researches IResearches) int64 {
	return int64(float64(b.WeaponPower) * (1 + float64(researches.GetWeaponsTechnology())*0.1))
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
func (b BaseDefender) DefenderConstructionTime(nbr, universeSpeed int64, acc DefenseAccelerators, lfBonuses LfBonuses) time.Duration {
	shipyardLvl := float64(acc.GetShipyard())
	naniteLvl := float64(acc.GetNaniteFactory())
	hours := float64(b.StructuralIntegrity) / (2500 * (1 + shipyardLvl) * float64(universeSpeed) * math.Pow(2, naniteLvl))
	secs := math.Max(1, hours*3600)
	dur := time.Duration(int64(math.Floor(secs))*nbr) * time.Second
	bonus := lfBonuses.CostTimeBonuses[b.ID].Duration
	return time.Duration(float64(dur) - float64(dur)*bonus)
}

// ConstructionTime same as DefenderConstructionTime, needed for BaseOgameObj implementation
// func (b BaseDefender) ConstructionTime(nbr, universeSpeed int64, acc BuildAccelerators, _, _ bool) time.Duration {
func (b BaseDefender) ConstructionTime(nbr, universeSpeed int64, acc BuildAccelerators, lfBonuses LfBonuses, _ CharacterClass, _ bool) time.Duration {
	return b.DefenderConstructionTime(nbr, universeSpeed, acc, lfBonuses)
}

// GetPrice returns the price of nbr defender units
func (b BaseDefender) GetPrice(nbr int64, lfBonuses LfBonuses) Resources {
	price := b.Price.Mul(nbr)
	bonus := lfBonuses.CostTimeBonuses[b.ID].Cost
	return price.SubPercent(bonus)
}
