package ogame

// Base ...
type Base struct {
	ID           ID
	Name         string
	Requirements map[ID]int
}

// GetID ...
func (b Base) GetID() ID {
	return b.ID
}

// GetName ...
func (b Base) GetName() string {
	return b.Name
}

// GetRequirements ...
func (b Base) GetRequirements() map[ID]int {
	return b.Requirements
}

// IsAvailable ...
func (b Base) IsAvailable(resourcesBuildings ResourcesBuildings, facilities Facilities, researches Researches, _ int) bool {
	for id, levelNeeded := range b.Requirements {
		if id.IsResourceBuilding() {
			if resourcesBuildings.ByID(id) < levelNeeded {
				return false
			}
		} else if id.IsFacility() {
			if facilities.ByID(id) < levelNeeded {
				return false
			}
		} else if id.IsTech() {
			if researches.ByID(id) < levelNeeded {
				return false
			}
		}
	}
	return true
}
