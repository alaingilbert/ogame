package ogame

type allianceDepot struct {
	BaseBuilding
}

func newAllianceDepot() *allianceDepot {
	b := new(allianceDepot)
	b.Name = "alliance depot"
	b.ID = AllianceDepotID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Metal: 20000, Crystal: 40000}

	return b
}
