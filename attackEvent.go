package ogame

import (
	"strconv"
	"time"
)

// AttackEvent all information available about an enemy attack
type AttackEvent struct {
	MissionType  MissionID
	Origin       Coordinate
	Destination  Coordinate
	ArrivalTime  time.Time
	AttackerName string
	AttackerID   int
	UnionID      int
	Missiles     int
	Ships        *ShipsInfos
}

func (a AttackEvent) String() string {
	return "" +
		"Mission Type: " + strconv.Itoa(int(a.MissionType)) + "\n" +
		"      Origin: " + a.Origin.String() + "\n" +
		" Destination: " + a.Destination.String() + "\n" +
		" ArrivalTime: " + a.ArrivalTime.String() + "\n" +
		"  AttackerID: " + strconv.Itoa(a.AttackerID) + "\n" +
		"     UnionID: " + strconv.Itoa(a.UnionID) + "\n" +
		"    Missiles: " + strconv.Itoa(a.Missiles)
}
