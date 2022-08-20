package ogame

type computerTechnology struct {
	BaseTechnology
}

func newComputerTechnology() *computerTechnology {
	b := new(computerTechnology)
	b.Name = "computer technology"
	b.ID = ComputerTechnologyID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Crystal: 400, Deuterium: 600}
	b.Requirements = map[ID]int64{ResearchLabID: 1}
	return b
}
