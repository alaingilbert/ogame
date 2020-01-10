package ogame

// Fleet represent a player fleet information
type Fleet struct {
	Mission        MissionID
	ReturnFlight   bool
	ID             FleetID
	Resources      Resources
	Origin         Coordinate
	Destination    Coordinate
	Ships          ShipsInfos
	ArriveIn       int64
	BackIn         int64
	UnionID        int64
	TargetPlanetID int64
}
