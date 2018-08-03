package ogame

// IsAvailable ...
func IsAvailable() bool {
	return true
}

// SeabedDeuteriumDen ...
type seabedDeuteriumDen struct {
	BaseBuilding
}

// NewSeabedDeuteriumDen ...
func NewSeabedDeuteriumDen() *seabedDeuteriumDen {
	b := new(seabedDeuteriumDen)
	b.ID = SeabedDeuteriumDenID
	b.IncreaseFactor = 2.3
	b.BaseCost = Resources{Metal: 2645, Crystal: 2645}
	return b
}
