package ogame

type researchLab struct {
	BaseBuilding
}

func newResearchLab() *researchLab {
	b := new(researchLab)
	b.Name = "research lab"
	b.ID = ResearchLabID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Metal: 200, Crystal: 400, Deuterium: 200}
	return b
}
