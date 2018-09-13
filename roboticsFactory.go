package ogame

type roboticsFactory struct {
	BaseBuilding
}

func newRoboticsFactory() *roboticsFactory {
	b := new(roboticsFactory)
	b.Name = "robotics factory"
	b.ID = RoboticsFactoryID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Metal: 400, Crystal: 120, Deuterium: 200}
	return b
}
