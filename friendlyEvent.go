package ogame

import (
	"strconv"
	"time"
)

// FriendlyEvent all information available about an friendly fleet
type FriendlyEvent struct {
	ID              int64
	MissionType     MissionID
	Origin          Coordinate
	Destination     Coordinate
	DestinationName string
	ArrivalTime     time.Time
	ArriveIn        int64
	PlayerName      string
	PlayerID        int64
	Ships           *ShipsInfos
	Resources       Resources
}

func (a FriendlyEvent) String() string {
	return "" +
		"               ID: " + strconv.FormatInt(a.ID, 10) + "\n" +
		"     Mission Type: " + strconv.FormatInt(int64(a.MissionType), 10) + "\n" +
		"           Origin: " + a.Origin.String() + "\n" +
		"      Destination: " + a.Destination.String() + "\n" +
		" Destination Name: " + a.DestinationName + "\n" +
		"      ArrivalTime: " + a.ArrivalTime.String() + "\n" +
		"         PlayerID: " + strconv.FormatInt(a.PlayerID, 10) + "\n" +
		"       PlayerName: " + a.PlayerName + "\n" +
		"        Resources: " + a.Resources.String()
}
