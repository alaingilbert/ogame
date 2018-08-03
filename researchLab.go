package ogame

// ResearchLab ...
type researchLab struct {
	BaseBuilding
}

// NewResearchLab ...
func NewResearchLab() *researchLab {
	b := new(researchLab)
	b.ID = ResearchLabID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Metal: 200, Crystal: 400, Deuterium: 200}
	return b
}
