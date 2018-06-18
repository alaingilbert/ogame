package ogame

// Planet ...
type Planet struct {
	ogame      *OGame
	Img        string
	ID         PlanetID
	Name       string
	Diameter   int
	Coordinate Coordinate
	Fields     struct {
		Built int
		Total int
	}
	Temperature struct {
		Min int
		Max int
	}
}

// GetInfos ...
func (p *Planet) GetInfos() error {
	planet, err := p.ogame.GetPlanet(p.ID)
	if err != nil {
		return err
	}
	p.Img = planet.Img
	p.Name = planet.Name
	p.Coordinate = planet.Coordinate
	p.Diameter = planet.Diameter
	p.Fields = planet.Fields
	p.Temperature = planet.Temperature
	return nil
}

// GetResourceSettings gets the resources settings for specified planetID
func (p *Planet) GetResourceSettings() (ResourceSettings, error) {
	return p.ogame.GetResourceSettings(p.ID)
}

// GetResourcesBuildings gets the resources buildings levels
func (p *Planet) GetResourcesBuildings() (ResourcesBuildings, error) {
	return p.ogame.GetResourcesBuildings(p.ID)
}

// GetDefense ...
func (p *Planet) GetDefense() (Defense, error) {
	return p.ogame.GetDefense(p.ID)
}

// GetShips ...
func (p *Planet) GetShips() (Ships, error) {
	return p.ogame.GetShips(p.ID)
}

// GetFacilities ...
func (p *Planet) GetFacilities() (Facilities, error) {
	return p.ogame.GetFacilities(p.ID)
}

// BuildBuilding ...
func (p *Planet) BuildBuilding(buildingID ID) error {
	return p.ogame.BuildBuilding(p.ID, buildingID)
}

// BuildDefense builds a defense unit
func (p *Planet) BuildDefense(defenseID ID, nbr int) error {
	return p.ogame.BuildDefense(p.ID, defenseID, nbr)
}

// BuildShips builds a ship unit
func (p *Planet) BuildShips(shipID ID, nbr int) error {
	return p.ogame.BuildShips(p.ID, shipID, nbr)
}

// BuildTechnology ...
func (p *Planet) BuildTechnology(technologyID ID) error {
	return p.ogame.BuildTechnology(p.ID, technologyID)
}

// GetResources gets user resources
func (p *Planet) GetResources() (Resources, error) {
	return p.ogame.GetResources(p.ID)
}

// SendFleet ...
func (p *Planet) SendFleet(ships []Quantifiable, speed Speed, where Coordinate, mission MissionID,
	resources Resources) (FleetID, error) {
	return p.ogame.SendFleet(p.ID, ships, speed, where, mission, resources)
}
