package ogame

type plasmaTechnology struct {
	BaseTechnology
}

func newPlasmaTechnology() *plasmaTechnology {
	b := new(plasmaTechnology)
	b.Name = "plasma technology"
	b.ID = PlasmaTechnologyID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Metal: 2000, Crystal: 4000, Deuterium: 1000}
	b.Requirements = map[ID]int{IonTechnologyID: 5, EnergyTechnologyID: 8, LaserTechnologyID: 10}
	return b
}
