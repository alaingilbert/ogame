package ogame

// CombustionDrive ...
type combustionDrive struct {
	BaseTechnology
}

// NewCombustionDrive ...
func NewCombustionDrive() *combustionDrive {
	b := new(combustionDrive)
	b.ID = CombustionDriveID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Metal: 400, Deuterium: 600}
	b.Requirements = map[ID]int{EnergyTechnologyID: 1}
	return b
}
