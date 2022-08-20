package ogame

import (
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
)

func TestAttackEvent_String(t *testing.T) {
	a := AttackEvent{
		ID:              123,
		MissionType:     3,
		Origin:          Coordinate{1, 2, 3, PlanetType},
		Destination:     Coordinate{4, 5, 6, PlanetType},
		DestinationName: "Homeworld",
		ArrivalTime:     time.Date(2018, 9, 11, 1, 2, 3, 4, time.UTC),
		AttackerID:      456,
		Missiles:        0,
		Ships:           &ShipsInfos{LargeCargo: 10},
	}
	expected := "" +
		"               ID: 123\n" +
		"     Mission Type: 3\n" +
		"           Origin: [P:1:2:3]\n" +
		"      Destination: [P:4:5:6]\n" +
		" Destination Name: Homeworld\n" +
		"      ArrivalTime: 2018-09-11 01:02:03.000000004 +0000 UTC\n" +
		"       AttackerID: 456\n" +
		"          UnionID: 0\n" +
		"         Missiles: 0"
	assert.Equal(t, expected, a.String())
}
