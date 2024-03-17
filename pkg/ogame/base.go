package ogame

// Base struct for all ogame objects
type Base struct {
	ID           ID
	Name         string
	Requirements map[ID]int64
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
func (b Base) GetRequirements() map[ID]int64 {
	return b.Requirements
}

// IsAvailable returns either or not the object is available to us
func (b Base) IsAvailable(t CelestialType, resourcesBuildings IResourcesBuildings, lfBuildings ILfBuildings,
	lfResearches ILfResearches, facilities IFacilities, researches IResearches, energy int64,
	characterClass CharacterClass) bool {
	if t != PlanetType && t != MoonType {
		return false
	}
	if t == PlanetType {
		if b.ID == LunarBaseID ||
			b.ID == SensorPhalanxID ||
			b.ID == JumpGateID {
			return false
		}
	} else if t == MoonType {
		if b.ID == MetalMineID ||
			b.ID == CrystalMineID ||
			b.ID == DeuteriumSynthesizerID ||
			b.ID == SolarPlantID ||
			b.ID == FusionReactorID ||
			b.ID == ResearchLabID ||
			b.ID == AllianceDepotID ||
			b.ID == MissileSiloID ||
			b.ID == NaniteFactoryID ||
			b.ID == TerraformerID ||
			b.ID == SpaceDockID {
			return false
		}
		if b.ID.IsTech() {
			return false
		}
	}
	if b.ID == GravitonTechnologyID && energy < 300000 {
		return false
	}
	if (b.ID == ReaperID && characterClass != General) ||
		(b.ID == PathfinderID && characterClass != Discoverer) ||
		(b.ID == CrawlerID && characterClass != Collector) {
		return false
	}
	type requirement struct {
		ID  ID
		Lvl int64
	}
	q := make([]requirement, 0)
	for id, levelNeeded := range b.Requirements {
		q = append(q, requirement{id, levelNeeded})
	}
	for len(q) > 0 {
		var req requirement
		req, q = q[0], q[1:]
		if t == PlanetType && b.ID.IsTech() {
			reqs := Objs.ByID(req.ID).GetRequirements()
			for k, v := range reqs {
				q = append(q, requirement{k, v})
			}
		}
		id := req.ID
		levelNeeded := req.Lvl
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
		} else if id.IsLfBuilding() {
			if lfBuildings.ByID(id) < levelNeeded {
				return false
			}
		} else if id.IsLfTech() {
			lfResearch := lfResearches.ByID(id)
			if lfResearch == nil || *lfResearch < levelNeeded {
				return false
			}
		}
	}
	return true
}
