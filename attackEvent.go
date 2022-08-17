package ogame

import (
	"time"
)

// AttackEvent all information available about an enemy attack
type AttackEvent struct {
	ID              int64
	MissionType     MissionID
	Origin          Coordinate
	Destination     Coordinate
	DestinationName string
	ArrivalTime     time.Time
	ArriveIn        int64
	AttackerName    string
	AttackerID      int64
	UnionID         int64
	Missiles        int64
	Ships           *ShipsInfos
}

func (a AttackEvent) String() string {
	return "" +
		"               ID: " + FI64(a.ID) + "\n" +
		"     Mission Type: " + FI64(int64(a.MissionType)) + "\n" +
		"           Origin: " + a.Origin.String() + "\n" +
		"      Destination: " + a.Destination.String() + "\n" +
		" Destination Name: " + a.DestinationName + "\n" +
		"      ArrivalTime: " + a.ArrivalTime.String() + "\n" +
		"       AttackerID: " + FI64(a.AttackerID) + "\n" +
		"          UnionID: " + FI64(a.UnionID) + "\n" +
		"         Missiles: " + FI64(a.Missiles)
}
