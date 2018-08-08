package ogame

// GravitonTechnology ...
type gravitonTechnology struct {
	BaseTechnology
}

func newGravitonTechnology() *gravitonTechnology {
	b := new(gravitonTechnology)
	b.ID = GravitonTechnologyID
	b.IncreaseFactor = 3.0
	b.BaseCost = Resources{Energy: 300000}
	b.Requirements = map[ID]int{ResearchLabID: 12}
	return b
}

// IsAvailable ...
func (b gravitonTechnology) IsAvailable(resourcesBuildings ResourcesBuildings, facilities Facilities,
	researches Researches, energy int) bool {
	if energy < 300000 {
		return false
	}
	for ogameID, levelNeeded := range b.Requirements {
		if ogameID.IsResourceBuilding() {
			if resourcesBuildings.ByOGameID(ogameID) < levelNeeded {
				return false
			}
		} else if ogameID.IsFacility() {
			if facilities.ByOGameID(ogameID) < levelNeeded {
				return false
			}
		} else if ogameID.IsTech() {
			if researches.ByOGameID(ogameID) < levelNeeded {
				return false
			}
		}
	}
	return true
}
