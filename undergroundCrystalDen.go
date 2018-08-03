package ogame

// UndergroundCrystalDen ...
type undergroundCrystalDen struct {
	BaseBuilding
}

// NewUndergroundCrystalDen ...
func NewUndergroundCrystalDen() *undergroundCrystalDen {
	b := new(undergroundCrystalDen)
	b.ID = UndergroundCrystalDenID
	b.IncreaseFactor = 2.3
	b.BaseCost = Resources{Metal: 2645, Crystal: 1322}
	return b
}
