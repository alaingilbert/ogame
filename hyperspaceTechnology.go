package ogame

// HyperspaceTechnology ...
type hyperspaceTechnology struct {
	BaseTechnology
}

func newHyperspaceTechnology() *hyperspaceTechnology {
	b := new(hyperspaceTechnology)
	b.ID = HyperspaceTechnologyID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Crystal: 4000, Deuterium: 2000}
	b.Requirements = map[ID]int{ResearchLabID: 7, ShieldingTechnologyID: 5, EnergyTechnologyID: 5}
	return b
}
