package ogame

type lunarBase struct {
	BaseBuilding
}

func newLunarBase() *lunarBase {
	b := new(lunarBase)
	b.Name = "lunar base"
	b.ID = LunarBaseID
	b.IncreaseFactor = 2
	b.BaseCost = Resources{Metal: 20000, Crystal: 40000, Deuterium: 20000}
	b.Requirements = map[ID]int64{}
	return b
}

// DeconstructionPrice lunar base cannot be deconstructed
func (s *lunarBase) DeconstructionPrice(level int64, techs IResearches) Resources {
	return Resources{}
}
