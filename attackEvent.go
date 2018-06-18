package ogame

import (
	"strconv"
	"time"
)

// AttackEvent ...
type AttackEvent struct {
	MissionType MissionID
	Origin      Coordinate
	Destination Coordinate
	ArrivalTime time.Time
	AttackerID  int
	Missiles    int
}

func (a AttackEvent) String() string {
	return "" +
		"Mission Type: " + strconv.Itoa(int(a.MissionType)) + "\n" +
		"      Origin: " + a.Origin.String() + "\n" +
		" Destination: " + a.Destination.String() + "\n" +
		" ArrivalTime: " + a.ArrivalTime.String() + "\n" +
		"  AttackerID: " + strconv.Itoa(a.AttackerID) + "\n" +
		"    Missiles: " + strconv.Itoa(a.Missiles)
}
