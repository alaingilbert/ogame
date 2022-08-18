package ogame

type jumpGate struct {
	BaseBuilding
}

func newJumpGate() *jumpGate {
	b := new(jumpGate)
	b.Name = "jump gate"
	b.ID = JumpGateID
	b.IncreaseFactor = 2
	b.BaseCost = Resources{Metal: 2000000, Crystal: 4000000, Deuterium: 2000000}
	b.Requirements = map[ID]int64{LunarBaseID: 1, HyperspaceTechnologyID: 7}
	return b
}
