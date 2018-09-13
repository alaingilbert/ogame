package ogame

type seabedDeuteriumDen struct {
	BaseBuilding
}

func newSeabedDeuteriumDen() *seabedDeuteriumDen {
	b := new(seabedDeuteriumDen)
	b.Name = "seabed deuterium den"
	b.ID = SeabedDeuteriumDenID
	b.IncreaseFactor = 2.3
	b.BaseCost = Resources{Metal: 2645, Crystal: 2645}
	return b
}
