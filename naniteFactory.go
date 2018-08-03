package ogame

// NaniteFactory ...
type naniteFactory struct {
	BaseBuilding
}

// NewNaniteFactory ...
func NewNaniteFactory() *naniteFactory {
	b := new(naniteFactory)
	b.ID = NaniteFactoryID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Metal: 1000000, Crystal: 500000, Deuterium: 100000}
	b.Requirements = map[ID]int{RoboticsFactoryID: 10, ComputerTechnologyID: 10}
	return b
}
