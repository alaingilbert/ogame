package ogame

// BaseShip base struct for ships
type BaseShip struct {
	BaseDefender
	BaseCargoCapacity int64
	BaseSpeed         int64
	FuelConsumption   int64
}

// GetCargoCapacity returns ship cargo capacity
func (b BaseShip) GetCargoCapacity(techs IResearches, probeRaids, isCollector, isPioneers bool) int64 {
	if b.GetID() == EspionageProbeID && !probeRaids {
		return 0
	}
	hyperspaceBonus := 0.05
	if isPioneers {
		hyperspaceBonus = 0.02
	}
	cargo := b.BaseCargoCapacity + int64(float64(b.BaseCargoCapacity*techs.GetHyperspaceTechnology())*hyperspaceBonus)
	if isCollector && (b.ID == SmallCargoID || b.ID == LargeCargoID) {
		cargo += int64(float64(b.BaseCargoCapacity) * 0.25)
	}
	return cargo
}

// GetFuelConsumption returns ship fuel consumption
func (b BaseShip) GetFuelConsumption(techs IResearches, fleetDeutSaveFactor float64, isGeneral bool) int64 {
	fuelConsumption := b.FuelConsumption
	if b.ID == SmallCargoID && techs.GetImpulseDrive() >= 5 {
		fuelConsumption *= 2
	} else if b.ID == RecyclerID && techs.GetHyperspaceDrive() >= 15 {
		fuelConsumption *= 3
	} else if b.ID == RecyclerID && techs.GetImpulseDrive() >= 17 {
		fuelConsumption *= 2
	}
	fuelConsumption = int64(fleetDeutSaveFactor * float64(fuelConsumption))
	if isGeneral {
		fuelConsumption = int64(float64(fuelConsumption) / 2)
	}
	return fuelConsumption
}

// GetSpeed returns speed of the ship
func (b BaseShip) GetSpeed(techs IResearches, isCollector, isGeneral bool) int64 {
	techDriveLvl := 0.0
	driveFactor := 0.2
	baseSpeed := float64(b.BaseSpeed)
	multiplier := int64(1)
	if b.ID == SmallCargoID && techs.GetImpulseDrive() >= 5 {
		baseSpeed = 10000
		techDriveLvl = float64(techs.GetImpulseDrive())
	} else if b.ID == BomberID && techs.GetHyperspaceDrive() >= 8 {
		baseSpeed = 5000
		techDriveLvl = float64(techs.GetHyperspaceDrive())
		driveFactor = 0.3
	} else if b.ID == RecyclerID && techs.GetHyperspaceDrive() >= 15 {
		techDriveLvl = float64(techs.GetHyperspaceDrive())
		multiplier = 3
		driveFactor = 0.3
	} else if b.ID == RecyclerID && techs.GetImpulseDrive() >= 17 {
		techDriveLvl = float64(techs.GetImpulseDrive())
		multiplier = 2
	} else if _, ok := b.Requirements[CombustionDriveID]; ok {
		techDriveLvl = float64(techs.GetCombustionDrive())
		driveFactor = 0.1
	} else if _, ok := b.Requirements[ImpulseDriveID]; ok {
		techDriveLvl = float64(techs.GetImpulseDrive())
	} else if _, ok := b.Requirements[HyperspaceDriveID]; ok {
		techDriveLvl = float64(techs.GetHyperspaceDrive())
		driveFactor = 0.3
	}
	speed := baseSpeed + (baseSpeed*driveFactor)*techDriveLvl
	if isCollector && (b.ID == SmallCargoID || b.ID == LargeCargoID) {
		speed += baseSpeed
	} else if isGeneral && (b.ID == RecyclerID || b.ID.IsCombatShip()) && b.ID != DeathstarID {
		speed += baseSpeed
	}
	return int64(speed) * multiplier
}
