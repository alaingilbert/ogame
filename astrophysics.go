package ogame

// Astrophysics ...
type astrophysics struct {
	BaseTechnology
}

// NewAstrophysics ...
func NewAstrophysics() *astrophysics {
	b := new(astrophysics)
	b.ID = AstrophysicsID
	b.IncreaseFactor = 1.75
	b.BaseCost = Resources{Metal: 4000, Crystal: 8000, Deuterium: 4000}
	b.Requirements = map[ID]int{EspionageTechnologyID: 4, ImpulseDriveID: 3}
	return b
}
