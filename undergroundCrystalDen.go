package ogame

// UndergroundCrystalDen ...
type undergroundCrystalDen struct {
	BaseBuilding
}

func newUndergroundCrystalDen() *undergroundCrystalDen {
	b := new(undergroundCrystalDen)
	b.Name = "underground crystal den"
	b.ID = UndergroundCrystalDenID
	b.IncreaseFactor = 2.3
	b.BaseCost = Resources{Metal: 2645, Crystal: 1322}
	return b
}
