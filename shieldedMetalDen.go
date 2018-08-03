package ogame

// ShieldedMetalDen ...
type shieldedMetalDen struct {
	BaseBuilding
}

// NewShieldedMetalDen ...
func NewShieldedMetalDen() *shieldedMetalDen {
	b := new(shieldedMetalDen)
	b.ID = ShieldedMetalDenID
	b.IncreaseFactor = 2.3
	b.BaseCost = Resources{Metal: 2645}
	return b
}
