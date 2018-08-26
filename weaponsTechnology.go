package ogame

// WeaponsTechnology ...
type weaponsTechnology struct {
	BaseTechnology
}

func newWeaponsTechnology() *weaponsTechnology {
	b := new(weaponsTechnology)
	b.Name = "weapons technology"
	b.ID = WeaponsTechnologyID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Metal: 800, Crystal: 200}
	b.Requirements = map[ID]int{ResearchLabID: 4}
	return b
}
