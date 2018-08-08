package ogame

// ArmourTechnology ...
type armourTechnology struct {
	BaseTechnology
}

func newArmourTechnology() *armourTechnology {
	b := new(armourTechnology)
	b.ID = ArmourTechnologyID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Metal: 1000}
	b.Requirements = map[ID]int{ResearchLabID: 2}
	return b
}
