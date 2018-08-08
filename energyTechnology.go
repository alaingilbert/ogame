package ogame

// EnergyTechnology ...
type energyTechnology struct {
	BaseTechnology
}

func newEnergyTechnology() *energyTechnology {
	b := new(energyTechnology)
	b.ID = EnergyTechnologyID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Crystal: 800, Deuterium: 400}
	b.Requirements = map[ID]int{ResearchLabID: 1}
	return b
}

// IsAvailable ...
func (t *energyTechnology) IsAvailable(_ ResourcesBuildings, facilities Facilities, _ Researches, _ int) bool {
	return facilities.ResearchLab >= 1
}
