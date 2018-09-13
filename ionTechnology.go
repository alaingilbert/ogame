package ogame

type ionTechnology struct {
	BaseTechnology
}

func newIonTechnology() *ionTechnology {
	b := new(ionTechnology)
	b.Name = "ion technology"
	b.ID = IonTechnologyID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Metal: 1000, Crystal: 300, Deuterium: 100}
	b.Requirements = map[ID]int{ResearchLabID: 4, EnergyTechnologyID: 4, LaserTechnologyID: 5}
	return b
}
