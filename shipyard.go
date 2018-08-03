package ogame

// Shipyard ...
type shipyard struct {
	BaseBuilding
}

// NewShipyard ...
func NewShipyard() *shipyard {
	b := new(shipyard)
	b.ID = ShipyardID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Metal: 400, Crystal: 200, Deuterium: 100}
	b.Requirements = map[ID]int{RoboticsFactoryID: 2}
	return b
}
