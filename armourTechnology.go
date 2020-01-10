package ogame

type armourTechnology struct {
	BaseTechnology
}

func newArmourTechnology() *armourTechnology {
	b := new(armourTechnology)
	b.Name = "armour technology"
	b.ID = ArmourTechnologyID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Metal: 1000}
	b.Requirements = map[ID]int64{ResearchLabID: 2}
	return b
}
