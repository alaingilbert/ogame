package ogame

// Base struct for all ogame objects
type Base struct {
	ID           ID
	Name         string
	Requirements map[ID]int
}

// GetID returns the ogame id of the object
func (b Base) GetID() ID {
	return b.ID
}

// GetName returns the printable name of the object
func (b Base) GetName() string {
	return b.Name
}

// GetRequirements returns the requirements to have this object available
func (b Base) GetRequirements() map[ID]int {
	return b.Requirements
}

// IsAvailable returns either or not the object is available to us
func (b Base) IsAvailable(resourcesBuildings ResourcesBuildings, facilities Facilities, researches Researches, energy int) bool {
	if b.ID == GravitonTechnologyID && energy < 300000 {
		return false
	}
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
