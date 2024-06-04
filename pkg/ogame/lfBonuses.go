package ogame

type LfBonuses struct {
	LfResourceBonuses
	LfShipBonuses   LfShipBonuses
	CostTimeBonuses CostTimeBonuses
}

func NewLfBonuses() *LfBonuses {
	return &LfBonuses{
		LfShipBonuses:   make(LfShipBonuses),
		CostTimeBonuses: make(CostTimeBonuses),
	}
}

type CostTimeBonuses map[ID]CostTimeBonus

type CostTimeBonus struct {
	Cost     float64
	Duration float64
}

type LfShipBonuses map[ID]LfShipBonus

type LfShipBonus struct {
	ID                  ID
	StructuralIntegrity float64
	ShieldPower         float64
	WeaponPower         float64
	Speed               float64
	CargoCapacity       float64
	FuelConsumption     float64
}

type LfResourceBonuses struct {
}
