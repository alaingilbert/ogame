package ogame

import "math"

// Fields planet fields stats
type Fields struct {
	Built int
	Total int
}

// HasFieldAvailable returns either or not we can still build on this planet
func (f Fields) HasFieldAvailable() bool {
	return f.Built < f.Total
}

// Temperature planet temperature values
type Temperature struct {
	Min int
	Max int
}

// Mean returns the planet mean temperature
func (t Temperature) Mean() int {
	return int(math.Round(float64(t.Min+t.Max) / 2))
}

// Planet ogame planet object
type Planet struct {
	ogame       *OGame
	Img         string
	ID          PlanetID
	Name        string
	Diameter    int
	Coordinate  Coordinate
	Fields      Fields
	Temperature Temperature
	Moon        *Moon
}

// GetResourceSettings gets the resources settings for specified planetID
func (p *Planet) GetResourceSettings() (ResourceSettings, error) {
	return p.ogame.GetResourceSettings(p.ID)
}

// GetResourcesBuildings gets the resources buildings levels
func (p *Planet) GetResourcesBuildings() (ResourcesBuildings, error) {
	return p.ogame.GetResourcesBuildings(p.ID)
}

// GetDefense gets all the defenses units information
func (p *Planet) GetDefense() (DefensesInfos, error) {
	return p.ogame.GetDefense(p.ID)
}

// GetShips gets all ships units information
func (p *Planet) GetShips() (ShipsInfos, error) {
	return p.ogame.GetShips(p.ID)
}

// GetFacilities  gets all facilities information
func (p *Planet) GetFacilities() (Facilities, error) {
	return p.ogame.GetFacilities(p.ID)
}

// Build builds any ogame objects (building, technology, ship, defence)
func (p *Planet) Build(id ID, nbr int) error {
	return p.ogame.Build(p.ID, id, nbr)
}

// BuildCancelable builds any cancelable ogame objects (building, technology)
func (p *Planet) BuildCancelable(id ID) error {
	return p.ogame.BuildCancelable(p.ID, id)
}

// BuildBuilding ensure what is being built is a building
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

// BuildTechnology ensure that we're trying to build a technology
func (p *Planet) BuildTechnology(technologyID ID) error {
	return p.ogame.BuildTechnology(p.ID, technologyID)
}

// GetResources gets user resources
func (p *Planet) GetResources() (Resources, error) {
	return p.ogame.GetResources(CelestialID(p.ID))
}

// SendFleet sends a fleet
func (p *Planet) SendFleet(ships []Quantifiable, speed Speed, where Coordinate, destType DestinationType,
	mission MissionID, resources Resources) (FleetID, error) {
	return p.ogame.SendFleet(CelestialID(p.ID), ships, speed, where, destType, mission, resources)
}

// ConstructionsBeingBuilt returns the building & research being built, and the time remaining (secs)
func (p *Planet) ConstructionsBeingBuilt() (ID, int, ID, int) {
	return p.ogame.ConstructionsBeingBuilt(p.ID)
}

// CancelBuilding cancel the construction of a building
func (p *Planet) CancelBuilding() error {
	return p.ogame.CancelBuilding(p.ID)
}

// CancelResearch cancel the research
func (p *Planet) CancelResearch() error {
	return p.ogame.CancelResearch(p.ID)
}

// GetResourcesProductions gets the resources production
func (p *Planet) GetResourcesProductions() (Resources, error) {
	return p.ogame.GetResourcesProductions(p.ID)
}

// FlightTime calculate flight time and fuel needed
func (p *Planet) FlightTime(destination Coordinate, speed Speed, ships ShipsInfos) (secs, fuel int) {
	return p.ogame.FlightTime(p.Coordinate, destination, speed, ships)
}
