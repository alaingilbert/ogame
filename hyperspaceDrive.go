package ogame

type hyperspaceDrive struct {
	BaseTechnology
}

func newHyperspaceDrive() *hyperspaceDrive {
	b := new(hyperspaceDrive)
	b.Name = "graviton technology"
	b.ID = HyperspaceDriveID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Metal: 10000, Crystal: 20000, Deuterium: 6000}
	b.Requirements = map[ID]int{HyperspaceTechnologyID: 3}
	return b
}
