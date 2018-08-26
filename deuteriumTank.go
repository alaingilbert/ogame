package ogame

// DeuteriumTank ...
type deuteriumTank struct {
	StorageBuilding
}

func newDeuteriumTank() *deuteriumTank {
	b := new(deuteriumTank)
	b.Name = "deuterium tank"
	b.ID = DeuteriumTankID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Metal: 1000, Crystal: 1000}
	return b
}
