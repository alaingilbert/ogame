package ogame

// BaseDefender ...
type BaseDefender struct {
	Base
	StructuralIntegrity int
	ShieldPower         int
	WeaponPower         int
	RapidfireFrom       map[ID]int
}

// GetStructuralIntegrity ...
func (b BaseDefender) GetStructuralIntegrity(researches Researches) int {
	return int(float64(b.StructuralIntegrity) * (1 + float64(researches.ArmourTechnology)*0.1))
}

// GetShieldPower ...
func (b BaseDefender) GetShieldPower(researches Researches) int {
	return int(float64(b.ShieldPower) * (1 + float64(researches.ShieldingTechnology)*0.1))
}

// GetWeaponPower ...
func (b BaseDefender) GetWeaponPower(researches Researches) int {
	return int(float64(b.WeaponPower) * (1 + float64(researches.WeaponsTechnology)*0.1))
}

// GetRapidfireFrom ...
func (b BaseDefender) GetRapidfireFrom() map[ID]int {
	return b.RapidfireFrom
}
