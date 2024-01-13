package wrapper

import (
	"fmt"

	"github.com/alaingilbert/ogame/pkg/ogame"
)

// Planet ogame planet object
type Planet struct {
	ogame.Planet
	ogame *OGame
	Moon  *Moon
}

// String ..
func (p Planet) String() string {
	return fmt.Sprintf("%s %s", p.Name, p.Coordinate)
}

func (p Planet) GetMoon() *Moon { return p.Moon }

// GetProduction get what is in the production queue.
// (ships & defense being built)
func (p Planet) GetProduction() ([]ogame.Quantifiable, int64, error) {
	return p.ogame.GetProduction(p.ID.Celestial())
}

// GetResourceSettings gets the resources settings for specified planetID
func (p Planet) GetResourceSettings(options ...Option) (ogame.ResourceSettings, error) {
	return p.ogame.GetResourceSettings(p.ID, options...)
}

// GetResourcesBuildings gets the resources buildings levels
func (p Planet) GetResourcesBuildings(options ...Option) (ogame.ResourcesBuildings, error) {
	return p.ogame.GetResourcesBuildings(p.ID.Celestial(), options...)
}

// GetDefense gets all the defenses units information
func (p Planet) GetDefense(options ...Option) (ogame.DefensesInfos, error) {
	return p.ogame.GetDefense(p.ID.Celestial(), options...)
}

// GetShips gets all ships units information
func (p Planet) GetShips(options ...Option) (ogame.ShipsInfos, error) {
	return p.ogame.GetShips(p.ID.Celestial(), options...)
}

// GetFacilities  gets all facilities information
func (p Planet) GetFacilities(options ...Option) (ogame.Facilities, error) {
	return p.ogame.GetFacilities(p.ID.Celestial(), options...)
}

// Build builds any ogame objects (building, technology, ship, defence)
func (p Planet) Build(id ogame.ID, nbr int64) error {
	return p.ogame.Build(ogame.CelestialID(p.ID), id, nbr)
}

// TearDown tears down any ogame building
func (p Planet) TearDown(buildingID ogame.ID) error {
	return p.ogame.TearDown(ogame.CelestialID(p.ID), buildingID)
}

// BuildCancelable builds any cancelable ogame objects (building, technology)
func (p Planet) BuildCancelable(id ogame.ID) error {
	return p.ogame.BuildCancelable(ogame.CelestialID(p.ID), id)
}

// BuildBuilding ensure what is being built is a building
func (p Planet) BuildBuilding(buildingID ogame.ID) error {
	return p.ogame.BuildBuilding(ogame.CelestialID(p.ID), buildingID)
}

// BuildDefense builds a defense unit
func (p Planet) BuildDefense(defenseID ogame.ID, nbr int64) error {
	return p.ogame.BuildDefense(ogame.CelestialID(p.ID), defenseID, nbr)
}

// BuildShips builds a ship unit
func (p Planet) BuildShips(shipID ogame.ID, nbr int64) error {
	return p.ogame.BuildShips(ogame.CelestialID(p.ID), shipID, nbr)
}

// BuildTechnology ensure that we're trying to build a technology
func (p Planet) BuildTechnology(technologyID ogame.ID) error {
	return p.ogame.BuildTechnology(p.ID.Celestial(), technologyID)
}

// GetResources gets user resources
func (p Planet) GetResources() (ogame.Resources, error) {
	return p.ogame.GetResources(ogame.CelestialID(p.ID))
}

// GetResourcesDetails gets resources details
func (p Planet) GetResourcesDetails() (ogame.ResourcesDetails, error) {
	return p.ogame.GetResourcesDetails(ogame.CelestialID(p.ID))
}

// SendFleet sends a fleet
func (p Planet) SendFleet(ships []ogame.Quantifiable, speed ogame.Speed, where ogame.Coordinate,
	mission ogame.MissionID, resources ogame.Resources, holdingTime, unionID int64) (ogame.Fleet, error) {
	return p.ogame.SendFleet(ogame.CelestialID(p.ID), ships, speed, where, mission, resources, holdingTime, unionID)
}

// EnsureFleet either sends all the requested ships or fail
func (p Planet) EnsureFleet(ships []ogame.Quantifiable, speed ogame.Speed, where ogame.Coordinate,
	mission ogame.MissionID, resources ogame.Resources, holdingTime, unionID int64) (ogame.Fleet, error) {
	return p.ogame.EnsureFleet(ogame.CelestialID(p.ID), ships, speed, where, mission, resources, holdingTime, unionID)
}

// ConstructionsBeingBuilt returns the building & research being built, and the time remaining (secs)
func (p Planet) ConstructionsBeingBuilt() (ogame.ID, int64, ogame.ID, int64, ogame.ID, int64, ogame.ID, int64) {
	return p.ogame.ConstructionsBeingBuilt(ogame.CelestialID(p.ID))
}

// CancelBuilding cancel the construction of a building
func (p Planet) CancelBuilding() error {
	return p.ogame.CancelBuilding(ogame.CelestialID(p.ID))
}

// CancelLfBuilding cancel the construction of a lifeform building
func (p Planet) CancelLfBuilding() error {
	return p.ogame.CancelLfBuilding(ogame.CelestialID(p.ID))
}

// CancelResearch cancel the research
func (p Planet) CancelResearch() error {
	return p.ogame.CancelResearch(p.ID.Celestial())
}

// GetItems get all items information
func (p Planet) GetItems() ([]ogame.Item, error) {
	return p.ogame.GetItems(p.ID.Celestial())
}

// ActivateItem activate an item
func (p Planet) ActivateItem(ref string) error {
	return p.ogame.ActivateItem(ref, p.ID.Celestial())
}

// GetResourcesProductions gets the resources production
func (p Planet) GetResourcesProductions() (ogame.Resources, error) {
	return p.ogame.GetResourcesProductions(p.ID)
}

// FlightTime calculate flight time and fuel needed
func (p Planet) FlightTime(destination ogame.Coordinate, speed ogame.Speed, ships ogame.ShipsInfos, missionID ogame.MissionID) (secs, fuel int64) {
	return p.ogame.FlightTime(p.Coordinate, destination, speed, ships, missionID)
}

// SendIPM send interplanetary missiles
func (p Planet) SendIPM(planetID ogame.PlanetID, coord ogame.Coordinate, nbr int64, priority ogame.ID) (int64, error) {
	return p.ogame.SendIPM(planetID, coord, nbr, priority)
}

// GetLfBuildings gets the lifeform buildings levels
func (p Planet) GetLfBuildings(options ...Option) (ogame.LfBuildings, error) {
	return p.ogame.getLfBuildings(p.ID.Celestial(), options...)
}

// GetLfResearch gets the lifeform techs levels
func (p Planet) GetLfResearch(options ...Option) (ogame.LfResearches, error) {
	return p.ogame.getLfResearch(p.ID.Celestial(), options...)
}

// GetTechs gets (ogame.ResourcesBuildings, ogame.Facilities, ogame.ShipsInfos, ogame.DefensesInfos, ogame.Researches)
func (p Planet) GetTechs() (ogame.ResourcesBuildings, ogame.Facilities, ogame.ShipsInfos, ogame.DefensesInfos, ogame.Researches, ogame.LfBuildings, ogame.LfResearches, error) {
	return p.ogame.GetTechs(p.ID.Celestial())
}
