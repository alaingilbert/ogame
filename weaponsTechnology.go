package ogame

// WeaponsTechnology ...
type weaponsTechnology struct {
	BaseTechnology
}

// NewWeaponsTechnology ...
func NewWeaponsTechnology() *weaponsTechnology {
	b := new(weaponsTechnology)
	b.ID = WeaponsTechnologyID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Metal: 800, Crystal: 200}
	b.Requirements = map[ID]int{ResearchLabID: 4}
	return b
}
