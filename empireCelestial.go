package ogame

// EmpireCelestial celestial information extracted from empire page (commander only)
type EmpireCelestial struct {
	Name        string
	Diameter    int64
	Img         string
	ID          CelestialID
	Type        CelestialType
	Fields      Fields
	Temperature Temperature
	Coordinate  Coordinate
	Resources   Resources
	Supplies    ResourcesBuildings
	Facilities  Facilities
	Defenses    DefensesInfos
	Researches  Researches
	Ships       ShipsInfos
}
