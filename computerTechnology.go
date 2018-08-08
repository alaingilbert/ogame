package ogame

// ComputerTechnology ...
type computerTechnology struct {
	BaseTechnology
}

func newComputerTechnology() *computerTechnology {
	b := new(computerTechnology)
	b.ID = ComputerTechnologyID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Crystal: 400, Deuterium: 600}
	b.Requirements = map[ID]int{ResearchLabID: 1}
	return b
}
