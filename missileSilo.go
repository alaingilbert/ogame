package ogame

// MissileSilo ...
type missileSilo struct {
	BaseBuilding
}

// NewMissileSilo ...
func NewMissileSilo() *missileSilo {
	b := new(missileSilo)
	b.ID = MissileSiloID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Metal: 20000, Crystal: 20000, Deuterium: 1000}
	b.Requirements = map[ID]int{ShipyardID: 1}
	return b
}
