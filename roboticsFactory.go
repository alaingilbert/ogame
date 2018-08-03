package ogame

// RoboticsFactory ...
type roboticsFactory struct {
	BaseBuilding
}

// NewRoboticsFactory ...
func NewRoboticsFactory() *roboticsFactory {
	b := new(roboticsFactory)
	b.ID = RoboticsFactoryID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Metal: 400, Crystal: 120, Deuterium: 200}
	return b
}
