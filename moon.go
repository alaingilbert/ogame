package ogame

// MoonID ...
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
	Diameter   int
	Coordinate Coordinate
	Fields     Fields
}

// SendFleet sends a fleet
func (m *Moon) SendFleet(ships []Quantifiable, speed Speed, where Coordinate, destType DestinationType,
	mission MissionID, resources Resources) (FleetID, error) {
	return m.ogame.SendFleet(CelestialID(m.ID), ships, speed, where, destType, mission, resources)
}

// GetResources gets moon resources
func (m *Moon) GetResources() (Resources, error) {
	return m.ogame.GetResources(CelestialID(m.ID))
}

// GetFacilities gets the moon facilities
func (m *Moon) GetFacilities() (MoonFacilities, error) {
	return m.ogame.GetMoonFacilities(m.ID)
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
