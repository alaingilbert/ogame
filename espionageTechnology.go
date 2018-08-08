package ogame

// EspionageTechnology ...
type espionageTechnology struct {
	BaseTechnology
}

func newEspionageTechnology() *espionageTechnology {
	b := new(espionageTechnology)
	b.ID = EspionageTechnologyID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Metal: 200, Crystal: 1000, Deuterium: 200}
	b.Requirements = map[ID]int{ResearchLabID: 3}
	return b
}
