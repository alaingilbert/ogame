package ogame

// DeuteriumTank ...
type deuteriumTank struct {
	StorageBuilding
}

// NewDeuteriumTank ...
func NewDeuteriumTank() *deuteriumTank {
	b := new(deuteriumTank)
	b.ID = DeuteriumTankID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Metal: 1000, Crystal: 1000}
	return b
}
