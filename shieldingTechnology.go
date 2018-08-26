package ogame

// ShieldingTechnology ...
type shieldingTechnology struct {
	BaseTechnology
}

func newShieldingTechnology() *shieldingTechnology {
	b := new(shieldingTechnology)
	b.Name = "shielding technology"
	b.ID = ShieldingTechnologyID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Metal: 200, Crystal: 600}
	b.Requirements = map[ID]int{ResearchLabID: 6, EnergyTechnologyID: 3}
	return b
}
