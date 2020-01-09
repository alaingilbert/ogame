package ogame

// BaseShip base struct for ships
type BaseShip struct {
	BaseDefender
	BaseCargoCapacity int64
	BaseSpeed         int64
	FuelConsumption   int64
}

// GetCargoCapacity returns ship cargo capacity
func (b BaseShip) GetCargoCapacity(techs Researches, probeRaids bool) int64 {
	if b.GetID() == EspionageProbeID && !probeRaids {
		return 0
	}
	return b.BaseCargoCapacity + int64(float64(b.BaseCargoCapacity*techs.HyperspaceTechnology)*0.05)
}

// GetFuelConsumption returns ship fuel consumption
func (b BaseShip) GetFuelConsumption() int64 {
	return b.FuelConsumption
}

// GetSpeed returns speed of the ship
func (b BaseShip) GetSpeed(techs Researches) int64 {
	var techDriveLvl int64 = 0
	if b.ID == SmallCargoID && techs.ImpulseDrive >= 5 {
		baseSpeed := 10000
		return int64(float64(baseSpeed) + (float64(baseSpeed)*0.2)*float64(techs.ImpulseDrive))
	}
	if b.ID == BomberID && techs.HyperspaceDrive >= 8 {
		baseSpeed := 5000
		return int64(float64(baseSpeed) + (float64(baseSpeed)*0.3)*float64(techs.HyperspaceDrive))
	}
	if b.ID == RecyclerID && (techs.ImpulseDrive >= 17 || techs.HyperspaceDrive >= 15) {
		if techs.HyperspaceDrive >= 15 {
			return int64(float64(b.BaseSpeed)+(float64(b.BaseSpeed)*0.3)*float64(techs.HyperspaceDrive)) * 3
		}
		return int64(float64(b.BaseSpeed)+(float64(b.BaseSpeed)*0.2)*float64(techs.ImpulseDrive)) * 2
	}
	if minLvl, ok := b.Requirements[CombustionDrive.ID]; ok {
		techDriveLvl = MaxInt(techs.CombustionDrive, minLvl)
		return int64(float64(b.BaseSpeed) + (float64(b.BaseSpeed)*0.1)*float64(techDriveLvl))
	} else if minLvl, ok := b.Requirements[ImpulseDrive.ID]; ok {
		techDriveLvl = MaxInt(techs.ImpulseDrive, minLvl)
		return int64(float64(b.BaseSpeed) + (float64(b.BaseSpeed)*0.2)*float64(techDriveLvl))
	} else if minLvl, ok := b.Requirements[HyperspaceDrive.ID]; ok {
		techDriveLvl = MaxInt(techs.HyperspaceDrive, minLvl)
		return int64(float64(b.BaseSpeed) + (float64(b.BaseSpeed)*0.3)*float64(techDriveLvl))
	}
	return int64(float64(b.BaseSpeed) + (float64(b.BaseSpeed)*0.2)*float64(techDriveLvl))
}
