package ogame

type combustionDrive struct {
	BaseTechnology
}

func newCombustionDrive() *combustionDrive {
	b := new(combustionDrive)
	b.Name = "combustion drive"
	b.ID = CombustionDriveID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Metal: 400, Deuterium: 600}
	b.Requirements = map[ID]int64{EnergyTechnologyID: 1}
	return b
}
