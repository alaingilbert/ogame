package ogame

type gravitonTechnology struct {
	BaseTechnology
}

func newGravitonTechnology() *gravitonTechnology {
	b := new(gravitonTechnology)
	b.Name = "graviton technology"
	b.ID = GravitonTechnologyID
	b.IncreaseFactor = 3.0
	b.BaseCost = Resources{Energy: 300000}
	b.Requirements = map[ID]int{ResearchLabID: 12}
	return b
}
