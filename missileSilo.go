package ogame

type missileSilo struct {
	BaseBuilding
}

func newMissileSilo() *missileSilo {
	b := new(missileSilo)
	b.Name = "missile silo"
	b.ID = MissileSiloID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Metal: 20000, Crystal: 20000, Deuterium: 1000}
	b.Requirements = map[ID]int64{ShipyardID: 1}
	return b
}
