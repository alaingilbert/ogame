package ogame

type LfBonuses struct {
	LfResourceBonuses       LfResourceBonuses
	CharacterClassesBonuses CharacterClassesBonuses
	LfShipBonuses           LfShipBonuses
	CostTimeBonuses         CostTimeBonuses
	MiscBonuses             MiscBonuses

	// Following lifeform buildings decreases the costs and duration for researching new technologies.
	// Humans ResearchCentre / Rocktal RuneTechnologium / Mechas RoboticsResearchCentre / Kaelesh VortexChamber
	PlanetLfResearchCostTimeBonus CostTimeBonus

	// LifeformLevels is the experience level of each discovered lifeform (lfbonuses page).
	// It adds level x 0.1 % to every lifeform technology on planets where it is active.
	LifeformLevels map[LifeformType]int64
}

func NewLfBonuses() *LfBonuses {
	return &LfBonuses{
		LfShipBonuses:   make(LfShipBonuses),
		CostTimeBonuses: make(CostTimeBonuses),
		LifeformLevels:  make(map[LifeformType]int64),
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
	ResourcesExpedition float64
}

type CharacterClassesBonuses struct {
	Characterclasses3 float64
}

type MiscBonuses struct {
	PhalanxRange float64
}
