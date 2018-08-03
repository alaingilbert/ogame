package ogame

// Fleet ...
type Fleet struct {
	Mission      MissionID
	ReturnFlight bool
	ID           FleetID
	Resources    Resources
	Origin       Coordinate
	Destination  Coordinate
	Ships        ShipsInfos
	ArriveIn     int
}
