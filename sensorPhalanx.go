package ogame

import "math"

type sensorPhalanx struct {
	BaseBuilding
	scanConsumption int64
}

func newSensorPhalanx() *sensorPhalanx {
	b := new(sensorPhalanx)
	b.Name = "sensor phalanx"
	b.ID = SensorPhalanxID
	b.IncreaseFactor = 2
	b.BaseCost = Resources{Metal: 20000, Crystal: 40000, Deuterium: 20000}
	b.Requirements = map[ID]int64{LunarBaseID: 1}
	b.scanConsumption = 5000
	return b
}

// ScanConsumption deuterium consumption to scan
func (p sensorPhalanx) ScanConsumption() int64 {
	return p.scanConsumption
}

// GetRange gets sensor range
func (p sensorPhalanx) GetRange(lvl int64) int64 {
	var phalanxRange int64
	if lvl == 0 {
		phalanxRange = 0
	} else if lvl == 1 {
		phalanxRange = 1
	} else {
		phalanxRange = int64(math.Pow(float64(lvl), 2)) - 1
	}
	return phalanxRange
}
