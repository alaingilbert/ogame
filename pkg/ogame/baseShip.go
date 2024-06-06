package ogame

import (
	"math"
)

// BaseShip base struct for ships
type BaseShip struct {
	BaseDefender
	BaseCargoCapacity int64
	BaseSpeed         int64
	FuelConsumption   int64
}

// GetCargoCapacity returns ship cargo capacity
func (b BaseShip) GetCargoCapacity(techs IResearches, lfBonuses LfBonuses, characterClass CharacterClass, hyperspaceBonusMultiplier float64, probeRaids bool) int64 {
	id := b.GetID()
	if id == EspionageProbeID && !probeRaids {
		return 0
	}
	baseCargoCapacity := b.BaseCargoCapacity
	hyperspaceTech := techs.GetHyperspaceTechnology()
	lfBonus := int64(float64(baseCargoCapacity) * lfBonuses.LfShipBonuses[b.ID].CargoCapacity)
	hyperspaceBonus := int64(float64(baseCargoCapacity*hyperspaceTech) * hyperspaceBonusMultiplier)
	cargo := baseCargoCapacity + lfBonus + hyperspaceBonus
	if characterClass.IsCollector() && (id == SmallCargoID || id == LargeCargoID) {
		cargo += int64(float64(baseCargoCapacity) * 0.25)
	}
	if characterClass.IsGeneral() && id == RecyclerID {
		cargo += int64(float64(baseCargoCapacity) * 0.2)
	}
	return cargo
}

// GetFuelConsumption returns ship fuel consumption
func (b BaseShip) GetFuelConsumption(techs IResearches, lfBonuses LfBonuses, characterClass CharacterClass, fleetDeutSaveFactor float64) int64 {
	fuelConsumption := b.FuelConsumption
	if b.ID == SmallCargoID && techs.GetImpulseDrive() >= 5 {
		fuelConsumption *= 2
	} else if b.ID == RecyclerID && techs.GetHyperspaceDrive() >= 15 {
		fuelConsumption *= 3
	} else if b.ID == RecyclerID && techs.GetImpulseDrive() >= 17 {
		fuelConsumption *= 2
	}
	fuelConsumption = int64(fleetDeutSaveFactor * float64(fuelConsumption))
	lfBonus := float64(fuelConsumption) * lfBonuses.LfShipBonuses[b.ID].FuelConsumption
	if characterClass.IsGeneral() {
		fuelConsumption = int64(float64(fuelConsumption) / 2)
	}
	return fuelConsumption + int64(lfBonus)
}

// GetSpeed returns speed of the ship
func (b BaseShip) GetSpeed(techs IResearches, lfBonuses LfBonuses, characterClass CharacterClass) int64 {
	var techDriveLvl int64
	driveFactor := 0.2
	baseSpeed := float64(b.BaseSpeed)
	multiplier := int64(1)
	if b.ID == SmallCargoID && techs.GetImpulseDrive() >= 5 {
		baseSpeed = 10000
		techDriveLvl = techs.GetImpulseDrive()
	} else if b.ID == BomberID && techs.GetHyperspaceDrive() >= 8 {
		baseSpeed = 5000
		techDriveLvl = techs.GetHyperspaceDrive()
		driveFactor = 0.3
	} else if b.ID == RecyclerID && techs.GetHyperspaceDrive() >= 15 {
		techDriveLvl = techs.GetHyperspaceDrive()
		multiplier = 3
		driveFactor = 0.3
	} else if b.ID == RecyclerID && techs.GetImpulseDrive() >= 17 {
		techDriveLvl = techs.GetImpulseDrive()
		multiplier = 2
	} else if _, ok := b.Requirements[CombustionDriveID]; ok {
		techDriveLvl = techs.GetCombustionDrive()
		driveFactor = 0.1
	} else if _, ok := b.Requirements[ImpulseDriveID]; ok {
		techDriveLvl = techs.GetImpulseDrive()
	} else if _, ok := b.Requirements[HyperspaceDriveID]; ok {
		techDriveLvl = techs.GetHyperspaceDrive()
		driveFactor = 0.3
	}
	techDriveLvlF := float64(techDriveLvl)
	lfBonus := baseSpeed * lfBonuses.LfShipBonuses[b.ID].Speed
	driveBonus := baseSpeed * driveFactor * techDriveLvlF
	speed := baseSpeed + driveBonus + lfBonus
	if characterClass.IsCollector() && (b.ID == SmallCargoID || b.ID == LargeCargoID) {
		speed += baseSpeed
	} else if characterClass.IsGeneral() && (b.ID == RecyclerID || b.ID.IsCombatShip()) && b.ID != DeathstarID {
		speed += baseSpeed
	}
	return int64(math.Round(speed * float64(multiplier)))
}
