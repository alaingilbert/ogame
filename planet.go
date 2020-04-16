package ogame

import "math"

// Fields planet fields stats
type Fields struct {
	Built int64
	Total int64
}

// HasFieldAvailable returns either or not we can still build on this planet
func (f Fields) HasFieldAvailable() bool {
	return f.Built < f.Total
}

// Temperature planet temperature values
type Temperature struct {
	Min int64
	Max int64
}

// Mean returns the planet mean temperature
func (t Temperature) Mean() int64 {
	return int64(math.Round(float64(t.Min+t.Max) / 2))
}

// Planet ogame planet object
type Planet struct {
	ogame       *OGame
	Img         string
	ID          PlanetID
	Name        string
	Diameter    int64
	Coordinate  Coordinate
	Fields      Fields
	Temperature Temperature
	Moon        *Moon
}

// GetName ...
func (p Planet) GetName() string {
	return p.Name
}

// GetDiameter ...
func (p Planet) GetDiameter() int64 {
	return p.Diameter
}

// GetID ...
func (p Planet) GetID() CelestialID {
	return p.ID.Celestial()
}

// GetType ...
func (p Planet) GetType() CelestialType {
	return PlanetType
}

// GetCoordinate ...
func (p Planet) GetCoordinate() Coordinate {
	return p.Coordinate
}

// GetFields ...
func (p Planet) GetFields() Fields {
	return p.Fields
}

// GetProduction get what is in the production queue.
// (ships & defense being built)
func (p Planet) GetProduction() ([]Quantifiable, int64, error) {
	return p.ogame.GetProduction(p.ID.Celestial())
}

// GetResourceSettings gets the resources settings for specified planetID
func (p *Planet) GetResourceSettings() (ResourceSettings, error) {
	return p.ogame.GetResourceSettings(p.ID)
}

// GetResourcesBuildings gets the resources buildings levels
func (p Planet) GetResourcesBuildings() (ResourcesBuildings, error) {
	return p.ogame.GetResourcesBuildings(p.ID.Celestial())
}

// GetDefense gets all the defenses units information
func (p Planet) GetDefense() (DefensesInfos, error) {
	return p.ogame.GetDefense(p.ID.Celestial())
}

// GetShips gets all ships units information
func (p Planet) GetShips() (ShipsInfos, error) {
	return p.ogame.GetShips(p.ID.Celestial())
}

// GetFacilities  gets all facilities information
func (p Planet) GetFacilities() (Facilities, error) {
	return p.ogame.GetFacilities(p.ID.Celestial())
}

// Build builds any ogame objects (building, technology, ship, defence)
func (p Planet) Build(id ID, nbr int64) error {
	return p.ogame.Build(CelestialID(p.ID), id, nbr)
}

// TearDown tears down any ogame building
func (p Planet) TearDown(buildingID ID) error {
	return p.ogame.TearDown(CelestialID(p.ID), buildingID)
}

// BuildCancelable builds any cancelable ogame objects (building, technology)
func (p *Planet) BuildCancelable(id ID) error {
	return p.ogame.BuildCancelable(CelestialID(p.ID), id)
}

// BuildBuilding ensure what is being built is a building
func (p Planet) BuildBuilding(buildingID ID) error {
	return p.ogame.BuildBuilding(CelestialID(p.ID), buildingID)
}

// BuildDefense builds a defense unit
func (p Planet) BuildDefense(defenseID ID, nbr int64) error {
	return p.ogame.BuildDefense(CelestialID(p.ID), defenseID, nbr)
}

// BuildShips builds a ship unit
func (p *Planet) BuildShips(shipID ID, nbr int64) error {
	return p.ogame.BuildShips(CelestialID(p.ID), shipID, nbr)
}

// BuildTechnology ensure that we're trying to build a technology
func (p Planet) BuildTechnology(technologyID ID) error {
	return p.ogame.BuildTechnology(p.ID.Celestial(), technologyID)
}

// GetResources gets user resources
func (p Planet) GetResources() (Resources, error) {
	return p.ogame.GetResources(CelestialID(p.ID))
}

// GetResourcesDetails gets resources details
func (p Planet) GetResourcesDetails() (ResourcesDetails, error) {
	return p.ogame.GetResourcesDetails(CelestialID(p.ID))
}

// SendFleet sends a fleet
func (p Planet) SendFleet(ships []Quantifiable, speed Speed, where Coordinate,
	mission MissionID, resources Resources, expeditiontime, unionID int64) (Fleet, error) {
	return p.ogame.SendFleet(CelestialID(p.ID), ships, speed, where, mission, resources, expeditiontime, unionID)
}

// EnsureFleet either sends all the requested ships or fail
func (p Planet) EnsureFleet(ships []Quantifiable, speed Speed, where Coordinate,
	mission MissionID, resources Resources, expeditiontime, unionID int64) (Fleet, error) {
	return p.ogame.EnsureFleet(CelestialID(p.ID), ships, speed, where, mission, resources, expeditiontime, unionID)
}

// ConstructionsBeingBuilt returns the building & research being built, and the time remaining (secs)
func (p Planet) ConstructionsBeingBuilt() (ID, int64, ID, int64) {
	return p.ogame.ConstructionsBeingBuilt(CelestialID(p.ID))
}

// CancelBuilding cancel the construction of a building
func (p Planet) CancelBuilding() error {
	return p.ogame.CancelBuilding(CelestialID(p.ID))
}

// CancelResearch cancel the research
func (p Planet) CancelResearch() error {
	return p.ogame.CancelResearch(p.ID.Celestial())
}

// GetItems get all items information
func (p Planet) GetItems() ([]Item, error) {
	return p.ogame.GetItems(p.ID.Celestial())
}

// ActivateItem activate an item
func (p Planet) ActivateItem(ref string) error {
	return p.ogame.ActivateItem(ref, p.ID.Celestial())
}

// GetResourcesProductions gets the resources production
func (p *Planet) GetResourcesProductions() (Resources, error) {
	return p.ogame.GetResourcesProductions(p.ID)
}

// FlightTime calculate flight time and fuel needed
func (p *Planet) FlightTime(destination Coordinate, speed Speed, ships ShipsInfos) (secs, fuel int64) {
	return p.ogame.FlightTime(p.Coordinate, destination, speed, ships)
}

// SendIPM send interplanetary missiles
func (p *Planet) SendIPM(planetID PlanetID, coord Coordinate, nbr int64, priority ID) (int64, error) {
	return p.ogame.SendIPM(planetID, coord, nbr, priority)
}
