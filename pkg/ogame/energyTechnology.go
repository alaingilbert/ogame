package ogame

type energyTechnology struct {
	BaseTechnology
}

func newEnergyTechnology() *energyTechnology {
	b := new(energyTechnology)
	b.Name = "energy technology"
	b.ID = EnergyTechnologyID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Crystal: 800, Deuterium: 400}
	b.Requirements = map[ID]int64{ResearchLabID: 1}
	return b
}
