package wrapper

import (
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/go-errors/errors"
)

// Moon ogame moon object
type Moon struct {
	ogame.Moon
	ogame *OGame
}

// GetProduction get what is in the production queue.
// (ships & defense being built)
func (m Moon) GetProduction() ([]ogame.Quantifiable, int64, error) {
	return m.ogame.GetProduction(m.ID.Celestial())
}

// ConstructionsBeingBuilt returns the building & research being built, and the time remaining (secs)
func (m Moon) ConstructionsBeingBuilt() (ogame.ID, int64, ogame.ID, int64, ogame.ID, int64, ogame.ID, int64) {
	return m.ogame.ConstructionsBeingBuilt(ogame.CelestialID(m.ID))
}

// GetResourcesBuildings gets the resources buildings levels
func (m Moon) GetResourcesBuildings(options ...Option) (ogame.ResourcesBuildings, error) {
	return m.ogame.GetResourcesBuildings(m.ID.Celestial(), options...)
}

// GetDefense gets all the defenses units information
func (m Moon) GetDefense(options ...Option) (ogame.DefensesInfos, error) {
	return m.ogame.GetDefense(m.ID.Celestial(), options...)
}

// GetShips gets all ships units information
func (m Moon) GetShips(options ...Option) (ogame.ShipsInfos, error) {
	return m.ogame.GetShips(m.ID.Celestial(), options...)
}

// Build builds any ogame objects (building, technology, ship, defence)
func (m Moon) Build(id ogame.ID, nbr int64) error {
	return m.ogame.Build(ogame.CelestialID(m.ID), id, nbr)
}

// TearDown tears down any ogame building
func (m Moon) TearDown(buildingID ogame.ID) error {
	return m.ogame.TearDown(ogame.CelestialID(m.ID), buildingID)
}

// BuildTechnology ensure that we're trying to build a technology
func (m Moon) BuildTechnology(technologyID ogame.ID) error {
	return errors.New("cannot build technology on a moon")
}

// BuildDefense builds a defense unit
func (m Moon) BuildDefense(defenseID ogame.ID, nbr int64) error {
	return m.ogame.BuildDefense(ogame.CelestialID(m.ID), defenseID, nbr)
}

// BuildBuilding ensure what is being built is a building
func (m Moon) BuildBuilding(buildingID ogame.ID) error {
	return m.ogame.BuildBuilding(ogame.CelestialID(m.ID), buildingID)
}

// CancelBuilding cancel the construction of a building
func (m Moon) CancelBuilding() error {
	return m.ogame.CancelBuilding(ogame.CelestialID(m.ID))
}

// CancelLfBuilding cancel the construction of a lifeform building
func (p Moon) CancelLfBuilding() error {
	return p.ogame.CancelLfBuilding(ogame.CelestialID(p.ID))
}

// CancelResearch cancel the research
func (m Moon) CancelResearch() error {
	return m.ogame.CancelResearch(m.ID.Celestial())
}

// SendFleet sends a fleet
func (m Moon) SendFleet(ships ogame.ShipsInfos, speed ogame.Speed, where ogame.Coordinate,
	mission ogame.MissionID, resources ogame.Resources, holdingTime, unionID int64) (ogame.Fleet, error) {
	return m.ogame.SendFleet(ogame.CelestialID(m.ID), ships, speed, where, mission, resources, holdingTime, unionID)
}

// EnsureFleet either sends all the requested ships or fail
func (m Moon) EnsureFleet(ships ogame.ShipsInfos, speed ogame.Speed, where ogame.Coordinate,
	mission ogame.MissionID, resources ogame.Resources, holdingTime, unionID int64) (ogame.Fleet, error) {
	return m.ogame.EnsureFleet(ogame.CelestialID(m.ID), ships, speed, where, mission, resources, holdingTime, unionID)
}

// GetResources gets moon resources
func (m Moon) GetResources() (ogame.Resources, error) {
	return m.ogame.GetResources(ogame.CelestialID(m.ID))
}

// GetResourcesDetails gets resources details
func (m Moon) GetResourcesDetails() (ogame.ResourcesDetails, error) {
	return m.ogame.GetResourcesDetails(ogame.CelestialID(m.ID))
}

// GetFacilities gets the moon facilities
func (m Moon) GetFacilities(options ...Option) (ogame.Facilities, error) {
	return m.ogame.GetFacilities(m.ID.Celestial(), options...)
}

// GetItems get all items information
func (m Moon) GetItems() ([]ogame.Item, error) {
	return m.ogame.GetItems(m.ID.Celestial())
}

// ActivateItem activate an item
func (m Moon) ActivateItem(ref string) error {
	return m.ogame.ActivateItem(ref, m.ID.Celestial())
}

//// BuildFacility build a facility
//func (m *Moon) BuildFacility(ID) error {
//	return nil
//}
//
//func (m *Moon) Downgrade(ID) error {
//	return nil
//}
//

// Phalanx uses 5000 deuterium to scan a coordinate
func (m Moon) Phalanx(coord ogame.Coordinate) ([]ogame.PhalanxFleet, error) {
	return m.ogame.Phalanx(m.ID, coord)
}

//
//func (m *Moon) IsJumpGateReady() bool {
//	return false
//}
//
//func (m *Moon) UseJumpGate() error {
//	return nil
//}

// GetLfBuildings gets the lifeform buildings levels
func (m Moon) GetLfBuildings(options ...Option) (ogame.LfBuildings, error) {
	return m.ogame.GetLfBuildings(m.ID.Celestial(), options...)
}

// GetLfResearch gets the lifeform techs levels
func (m Moon) GetLfResearch(options ...Option) (ogame.LfResearches, error) {
	return m.ogame.GetLfResearch(m.ID.Celestial(), options...)
}

// GetTechs gets (ogame.ResourcesBuildings, ogame.Facilities, ogame.ShipsInfos, ogame.DefensesInfos, ogame.Researches)
func (m Moon) GetTechs() (ogame.ResourcesBuildings, ogame.Facilities, ogame.ShipsInfos, ogame.DefensesInfos, ogame.Researches, ogame.LfBuildings, ogame.LfResearches, error) {
	return m.ogame.GetTechs(m.ID.Celestial())
}
