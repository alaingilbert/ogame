package ogame

type intergalacticResearchNetwork struct {
	BaseTechnology
}

func newIntergalacticResearchNetwork() *intergalacticResearchNetwork {
	b := new(intergalacticResearchNetwork)
	b.Name = "intergalactic research network"
	b.ID = IntergalacticResearchNetworkID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Metal: 240000, Crystal: 400000, Deuterium: 160000}
	b.Requirements = map[ID]int64{ResearchLabID: 10, ComputerTechnologyID: 8, HyperspaceTechnologyID: 8}
	return b
}
