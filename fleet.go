package ogame

// Fleet represent a player fleet information
type Fleet struct {
	Mission      MissionID
	ReturnFlight bool
	ID           FleetID
	Resources    Resources
	Origin       Coordinate
	Destination  Coordinate
	Ships        ShipsInfos
	ArriveIn     int
	BackIn       int
}

// Fleet is not returning and we can calculate the in air time in sec. The function returns 0 if the fleet is returning!
func (f Fleet) InAir() int64 {
	if !f.ReturnFlight {
		TotalFleetFlightTime := (f.BackIn - f.ArriveIn)
		ArriveAt := time.Now().Unix() + f.ArriveIn
		FleetStartAt := (ArriveAt - TotalFleetFlightTime)
		InAirForSec := time.Now().Unix() - FleetStartAt
		return InAirForSec
	}
	return 0
}

// Fleet is not returning and we can calculate the total time the Fleet needs to arrive. The function returns 0 if the fleet is returning!
func (f Fleet) Flighttime() int64 {
	if !f.ReturnFlight {
		return (f.BackIn - f.ArriveIn)
	} else {
		return 0
	}
}

// Fleet calculate the ArriveAt Unix.Timestamp for returning and not returning fleets.
func (f Fleet) ArriveAt() int64 {
	if f.ReturnFlight {
		return time.Now().Unix() + f.BackIn
	} else {
		return time.Now().Unix() + (f.BackIn - f.ArriveIn)
	}
}