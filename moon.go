package ogame

import "github.com/go-errors/errors"

// MoonID represent a moon id
type MoonID CelestialID

// Celestial convert a MoonID to a CelestialID
func (m MoonID) Celestial() CelestialID {
	return CelestialID(m)
}

// Moon ogame moon object
type Moon struct {
	ogame      *OGame
	ID         MoonID
	Img        string
	Name       string
	Diameter   int64
	Coordinate Coordinate
	Fields     Fields
}

// GetName ...
func (m Moon) GetName() string {
	return m.Name
}

// GetDiameter ...
func (m Moon) GetDiameter() int64 {
	return m.Diameter
}

// GetID ...
func (m Moon) GetID() CelestialID {
	return m.ID.Celestial()
}

// GetType ...
func (m Moon) GetType() CelestialType {
	return MoonType
}

// GetCoordinate ...
func (m Moon) GetCoordinate() Coordinate {
	return m.Coordinate
}

// GetFields ...
func (m Moon) GetFields() Fields {
	return m.Fields
}

// GetProduction get what is in the production queue.
// (ships & defense being built)
func (m Moon) GetProduction() ([]Quantifiable, int64, error) {
	return m.ogame.GetProduction(m.ID.Celestial())
}

// ConstructionsBeingBuilt returns the building & research being built, and the time remaining (secs)
func (m Moon) ConstructionsBeingBuilt() (ID, int64, ID, int64) {
	return m.ogame.ConstructionsBeingBuilt(CelestialID(m.ID))
}

// GetResourcesBuildings gets the resources buildings levels
func (m Moon) GetResourcesBuildings() (ResourcesBuildings, error) {
	return m.ogame.GetResourcesBuildings(m.ID.Celestial())
}

// GetDefense gets all the defenses units information
func (m Moon) GetDefense() (DefensesInfos, error) {
	return m.ogame.GetDefense(m.ID.Celestial())
}

// GetShips gets all ships units information
func (m Moon) GetShips() (ShipsInfos, error) {
	return m.ogame.GetShips(m.ID.Celestial())
}

// Build builds any ogame objects (building, technology, ship, defence)
func (m Moon) Build(id ID, nbr int64) error {
	return m.ogame.Build(CelestialID(m.ID), id, nbr)
}

// TearDown tears down any ogame building
func (m Moon) TearDown(buildingID ID) error {
	return m.ogame.TearDown(CelestialID(m.ID), buildingID)
}

// BuildTechnology ensure that we're trying to build a technology
func (m Moon) BuildTechnology(technologyID ID) error {
	return errors.New("cannot build technology on a moon")
}

// BuildDefense builds a defense unit
func (m Moon) BuildDefense(defenseID ID, nbr int64) error {
	return m.ogame.BuildDefense(CelestialID(m.ID), defenseID, nbr)
}

// BuildBuilding ensure what is being built is a building
func (m Moon) BuildBuilding(buildingID ID) error {
	return m.ogame.BuildBuilding(CelestialID(m.ID), buildingID)
}

// CancelBuilding cancel the construction of a building
func (m Moon) CancelBuilding() error {
	return m.ogame.CancelBuilding(CelestialID(m.ID))
}

// CancelResearch cancel the research
func (m Moon) CancelResearch() error {
	return m.ogame.CancelResearch(m.ID.Celestial())
}

// SendFleet sends a fleet
func (m Moon) SendFleet(ships []Quantifiable, speed Speed, where Coordinate,
	mission MissionID, resources Resources, expeditiontime, unionID int64) (Fleet, error) {
	return m.ogame.SendFleet(CelestialID(m.ID), ships, speed, where, mission, resources, expeditiontime, unionID)
}

// EnsureFleet either sends all the requested ships or fail
func (m Moon) EnsureFleet(ships []Quantifiable, speed Speed, where Coordinate,
	mission MissionID, resources Resources, expeditiontime, unionID int64) (Fleet, error) {
	return m.ogame.EnsureFleet(CelestialID(m.ID), ships, speed, where, mission, resources, expeditiontime, unionID)
}

// GetResources gets moon resources
func (m Moon) GetResources() (Resources, error) {
	return m.ogame.GetResources(CelestialID(m.ID))
}

// GetResourcesDetails gets resources details
func (m Moon) GetResourcesDetails() (ResourcesDetails, error) {
	return m.ogame.GetResourcesDetails(CelestialID(m.ID))
}

// GetFacilities gets the moon facilities
func (m Moon) GetFacilities() (Facilities, error) {
	return m.ogame.GetFacilities(m.ID.Celestial())
}

// GetItems get all items information
func (m Moon) GetItems() ([]Item, error) {
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
func (m *Moon) Phalanx(coord Coordinate) ([]Fleet, error) {
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
