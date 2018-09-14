package ogame

type shieldedMetalDen struct {
	BaseBuilding
}

func newShieldedMetalDen() *shieldedMetalDen {
	b := new(shieldedMetalDen)
	b.Name = "shielded metal den"
	b.ID = ShieldedMetalDenID
	b.IncreaseFactor = 2.3
	b.BaseCost = Resources{Metal: 2645}
	return b
}
