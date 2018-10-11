package ogame

import "math"

type sensorPhalanx struct {
	BaseBuilding
	scanConsumption int
}

func newSensorPhalanx() *sensorPhalanx {
	b := new(sensorPhalanx)
	b.Name = "sensor phalanx"
	b.ID = SensorPhalanxID
	b.IncreaseFactor = 2
	b.BaseCost = Resources{Metal: 20000, Crystal: 40000, Deuterium: 20000}
	b.Requirements = map[ID]int{LunarBaseID: 1}
	b.scanConsumption = 5000
	return b
}

// ScanConsumption deuterium consumption to scan
func (p sensorPhalanx) ScanConsumption() int {
	return p.scanConsumption
}

// GetRange gets sensor range
func (p sensorPhalanx) GetRange(lvl int) int {
	phalanxRange := 0
	if lvl == 0 {
		phalanxRange = 0
	} else if lvl == 1 {
		phalanxRange = 1
	} else {
		phalanxRange = int(math.Pow(float64(lvl), 2)) - 1
	}
	return phalanxRange
}
