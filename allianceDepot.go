package ogame

// AllianceDepot ...
type allianceDepot struct {
	BaseBuilding
}

// New ...
func NewAllianceDepot() *allianceDepot {
	b := new(allianceDepot)
	b.ID = AllianceDepotID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Metal: 20000, Crystal: 40000}
	return b
}
