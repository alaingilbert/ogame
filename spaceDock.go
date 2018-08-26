package ogame

// SpaceDock ...
type spaceDock struct {
	BaseBuilding
}

func newSpaceDock() *spaceDock {
	b := new(spaceDock)
	b.Name = "space dock"
	b.ID = SpaceDockID
	b.IncreaseFactor = 5
	b.BaseCost = Resources{Metal: 200, Crystal: 50, Energy: 50}
	b.Requirements = map[ID]int{ShipyardID: 2}
	return b
}
