package ogame

type hyperspaceTechnology struct {
	BaseTechnology
}

func newHyperspaceTechnology() *hyperspaceTechnology {
	b := new(hyperspaceTechnology)
	b.Name = "hyperspace technology"
	b.ID = HyperspaceTechnologyID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Crystal: 4000, Deuterium: 2000}
	b.Requirements = map[ID]int64{ResearchLabID: 7, ShieldingTechnologyID: 5, EnergyTechnologyID: 5}
	return b
}
