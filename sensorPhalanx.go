package ogame

type sensorPhalanx struct {
	BaseBuilding
}

func newSensorPhalanx() *sensorPhalanx {
	b := new(sensorPhalanx)
	b.Name = "sensor phalanx"
	b.ID = SensorPhalanxID
	b.IncreaseFactor = 2
	b.BaseCost = Resources{Metal: 20000, Crystal: 40000, Deuterium: 20000}
	b.Requirements = map[ID]int{LunarBaseID: 1}
	return b
}
