package ogame

// BaseShip base struct for ships
type BaseShip struct {
	BaseDefender
	BaseCargoCapacity int64
	BaseSpeed         int64
	FuelConsumption   int64
}

// GetCargoCapacity returns ship cargo capacity
func (b BaseShip) GetCargoCapacity(techs Researches, probeRaids, isCollector bool) int64 {
	if b.GetID() == EspionageProbeID && !probeRaids {
		return 0
	}
	cargo := b.BaseCargoCapacity + int64(float64(b.BaseCargoCapacity*techs.HyperspaceTechnology)*0.05)
	if isCollector && (b.ID == SmallCargoID || b.ID == LargeCargoID) {
		cargo += int64(float64(b.BaseCargoCapacity) * 0.25)
	}
	return cargo
}

// GetFuelConsumption returns ship fuel consumption
func (b BaseShip) GetFuelConsumption(techs Researches) int64 {
	fuelConsumption := b.FuelConsumption
	if b.ID == SmallCargoID && techs.ImpulseDrive >= 5 {
		fuelConsumption *= 2
	} else if b.ID == RecyclerID && techs.HyperspaceDrive >= 15 {
		fuelConsumption *= 3
	} else if b.ID == RecyclerID && techs.ImpulseDrive >= 17 {
		fuelConsumption *= 2
	}
	return fuelConsumption
}

// GetSpeed returns speed of the ship
func (b BaseShip) GetSpeed(techs Researches, isCollector, isGeneral bool) int64 {
	techDriveLvl := 0.0
	driveFactor := 0.2
	baseSpeed := float64(b.BaseSpeed)
	multiplier := int64(1)
	if b.ID == SmallCargoID && techs.ImpulseDrive >= 5 {
		baseSpeed = 10000
		techDriveLvl = float64(techs.ImpulseDrive)
	} else if b.ID == BomberID && techs.HyperspaceDrive >= 8 {
		baseSpeed = 5000
		techDriveLvl = float64(techs.HyperspaceDrive)
		driveFactor = 0.3
	} else if b.ID == RecyclerID && techs.HyperspaceDrive >= 15 {
		techDriveLvl = float64(techs.HyperspaceDrive)
		multiplier = 3
		driveFactor = 0.3
	} else if b.ID == RecyclerID && techs.ImpulseDrive >= 17 {
		techDriveLvl = float64(techs.ImpulseDrive)
		multiplier = 2
	} else if minLvl, ok := b.Requirements[CombustionDrive.ID]; ok {
		techDriveLvl = float64(MaxInt(techs.CombustionDrive, minLvl))
		driveFactor = 0.1
	} else if minLvl, ok := b.Requirements[ImpulseDrive.ID]; ok {
		techDriveLvl = float64(MaxInt(techs.ImpulseDrive, minLvl))
	} else if minLvl, ok := b.Requirements[HyperspaceDrive.ID]; ok {
		techDriveLvl = float64(MaxInt(techs.HyperspaceDrive, minLvl))
		driveFactor = 0.3
	}
	speed := baseSpeed + (baseSpeed*driveFactor)*techDriveLvl
	if isCollector && (b.ID == SmallCargoID || b.ID == LargeCargoID) {
		speed += baseSpeed
	} else if isGeneral && (b.ID == RecyclerID || b.ID.IsCombatShip()) {
		speed += baseSpeed
	}
	return int64(speed) * multiplier
}
