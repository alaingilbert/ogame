package ogame

type naniteFactory struct {
	BaseBuilding
}

func newNaniteFactory() *naniteFactory {
	b := new(naniteFactory)
	b.Name = "nanite factory"
	b.ID = NaniteFactoryID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Metal: 1000000, Crystal: 500000, Deuterium: 100000}
	b.Requirements = map[ID]int64{RoboticsFactoryID: 10, ComputerTechnologyID: 10}
	return b
}
