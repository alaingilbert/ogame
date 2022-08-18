package v6

import (
	"github.com/alaingilbert/clockwork"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
	"time"
)

func TestExtractAttacks(t *testing.T) {
	clock := clockwork.NewFakeClockAt(time.Date(2016, 8, 23, 17, 48, 13, 0, time.UTC))
	pageHTMLBytes, _ := ioutil.ReadFile("../../../samples/unversioned/event_list_attack.html")
	attacks, _ := NewExtractor().extractAttacks(pageHTMLBytes, clock, nil)
	assert.Equal(t, 1, len(attacks))
	assert.Equal(t, "Homeworld", attacks[0].DestinationName)
	assert.Equal(t, clock.Now().Add(14*time.Minute), attacks[0].ArrivalTime.UTC())
	assert.Equal(t, int64(14*60), attacks[0].ArriveIn)
}

func TestExtractAttacksFromFullPage(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../../samples/unversioned/overview_always_events.html")
	attacks, err := NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(attacks))
	assert.Equal(t, int64(1), attacks[0].Ships.SmallCargo)

	pageHTMLBytes, _ = ioutil.ReadFile("../../../samples/unversioned/overview_active.html")
	_, err = NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), nil)
	assert.EqualError(t, err, ogame.ErrEventsBoxNotDisplayed.Error())

	pageHTMLBytes, _ = ioutil.ReadFile("../../../samples/unversioned/eventlist_loggedout.html")
	_, err = NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), nil)
	assert.EqualError(t, err, ogame.ErrNotLogged.Error())
}

func TestExtractAttacksPhoneDisplay(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../../samples/unversioned/event_list_attack_phone.html")
	attacks, _ := NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), nil)
	assert.Equal(t, 1, len(attacks))
	assert.Equal(t, int64(106734), attacks[0].AttackerID)
	assert.Equal(t, "", attacks[0].AttackerName, "should not be able to get the name")
}

func TestExtractAttacksMeAttacking(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../../samples/unversioned/eventlist_me_attacking.html")
	attacks, _ := NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), nil)
	assert.Equal(t, 0, len(attacks))
}

func TestExtractAttacksWithoutShips(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../../samples/unversioned/event_list_attack.html")
	attacks, _ := NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), nil)
	assert.Equal(t, 1, len(attacks))
	assert.Equal(t, int64(100771), attacks[0].AttackerID)
	assert.Equal(t, int64(0), attacks[0].Missiles)
	assert.Nil(t, attacks[0].Ships)
}

func TestExtractAttacksWithShips(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../../samples/unversioned/eventList_attack_ships.html")
	attacks, _ := NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), nil)
	assert.Equal(t, 1, len(attacks))
	assert.Equal(t, "hammad", attacks[0].AttackerName)
	assert.Equal(t, int64(107088), attacks[0].AttackerID)
	assert.NotNil(t, attacks[0].Ships)
	assert.Equal(t, ogame.PlanetType, attacks[0].Destination.Type)
	assert.Equal(t, int64(197), attacks[0].Ships.LargeCargo)
	assert.Equal(t, int64(3), attacks[0].Ships.LightFighter)
	assert.Equal(t, int64(8), attacks[0].Ships.HeavyFighter)
	assert.Equal(t, int64(92), attacks[0].Ships.Cruiser)
	assert.Equal(t, int64(571), attacks[0].Ships.EspionageProbe)
	assert.Equal(t, int64(27), attacks[0].Ships.Bomber)
	assert.Equal(t, int64(4), attacks[0].Ships.Destroyer)
	assert.Equal(t, int64(11), attacks[0].Ships.Battlecruiser)
}

func TestExtractAttacksMoon(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../../samples/unversioned/eventlist_moon_attacked.html")
	attacks, _ := NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), nil)
	assert.Equal(t, 1, len(attacks))
	assert.NotNil(t, attacks[0].Ships)
	assert.Equal(t, int64(107009), attacks[0].AttackerID)
	assert.Equal(t, ogame.Coordinate{4, 212, 8, ogame.PlanetType}, attacks[0].Origin)
	assert.Equal(t, ogame.Coordinate{4, 116, 12, ogame.MoonType}, attacks[0].Destination)
	assert.Equal(t, ogame.MoonType, attacks[0].Destination.Type)
	assert.Equal(t, int64(1), attacks[0].Ships.SmallCargo)
	assert.Equal(t, "Moon", attacks[0].DestinationName)
}

func TestExtractAttacksMoonDestruction(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../../samples/unversioned/eventlist_moon_destruction.html")
	attacks, _ := NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), nil)
	assert.Equal(t, 1, len(attacks))
	assert.NotNil(t, attacks[0].Ships)
	assert.Equal(t, int64(106734), attacks[0].AttackerID)
	assert.Equal(t, ogame.Coordinate{4, 116, 12, ogame.PlanetType}, attacks[0].Origin)
	assert.Equal(t, ogame.Coordinate{4, 116, 9, ogame.MoonType}, attacks[0].Destination)
	assert.Equal(t, ogame.MoonType, attacks[0].Destination.Type)
	assert.Equal(t, int64(1), attacks[0].Ships.Deathstar)
}

func TestExtractAttacksWithThousandsShips(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../../samples/unversioned/eventlist_attack_thousands.html")
	ownCoords := make([]ogame.Coordinate, 0)
	attacks, _ := NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), ownCoords)
	assert.Equal(t, 2, len(attacks))
	assert.Equal(t, int64(1012), attacks[1].Ships.Cruiser)
	assert.Equal(t, int64(1000), attacks[1].Ships.LargeCargo)
}

func TestExtractAttacksUnknownShips(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../../samples/unversioned/eventlist_unknown_ships_nbr.html")
	attacks, _ := NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), nil)
	assert.Equal(t, 1, len(attacks))
	assert.Equal(t, int64(-1), attacks[0].Ships.Cruiser)
	assert.Equal(t, int64(0), attacks[0].Ships.Destroyer)
}

func TestExtractAttacksACS(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../../samples/unversioned/eventlist_acs.html")
	attacks, _ := NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), nil)
	assert.Equal(t, 1, len(attacks))
	assert.Equal(t, ogame.GroupedAttack, attacks[0].MissionType)
	assert.Equal(t, int64(10), attacks[0].Ships.LightFighter)
	assert.Equal(t, int64(2176), attacks[0].Ships.Battlecruiser)
}

func TestExtractAttacksACSMany(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../../samples/unversioned/eventlist_acs_multiple.html")
	attacks, _ := NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), nil)
	assert.Equal(t, 3, len(attacks))
	assert.Equal(t, ogame.GroupedAttack, attacks[0].MissionType)
	assert.Equal(t, int64(2), attacks[0].Ships.LightFighter)
	assert.Equal(t, int64(3), attacks[0].Ships.Battlecruiser)
	assert.Equal(t, ogame.GroupedAttack, attacks[1].MissionType)
	assert.Equal(t, int64(4), attacks[1].Ships.LightFighter)
	assert.Equal(t, int64(5), attacks[1].Ships.Battlecruiser)
	assert.Equal(t, ogame.Attack, attacks[2].MissionType)
	assert.Equal(t, int64(1), attacks[2].Ships.LightFighter)
	assert.Equal(t, int64(7), attacks[2].Ships.Battlecruiser)
}

func TestExtractAttacksACS2(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../../samples/unversioned/eventlist_acs2.html")
	attacks, _ := NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), nil)
	assert.Equal(t, 2, len(attacks))
	assert.Equal(t, ogame.GroupedAttack, attacks[0].MissionType)
	assert.Equal(t, int64(106734), attacks[0].AttackerID)
	assert.Equal(t, ogame.Coordinate{4, 116, 12, ogame.PlanetType}, attacks[0].Origin)
	assert.Equal(t, int64(4), attacks[0].Ships.SmallCargo)
	assert.Equal(t, int64(3), attacks[0].Ships.Battlecruiser)
	assert.Equal(t, ogame.GroupedAttack, attacks[1].MissionType)
	assert.Equal(t, int64(106734), attacks[1].AttackerID)
	assert.Equal(t, ogame.Coordinate{4, 116, 12, ogame.PlanetType}, attacks[1].Origin)
	assert.Equal(t, int64(7), attacks[1].Ships.SmallCargo)
	assert.Equal(t, int64(11), attacks[1].Ships.Battlecruiser)
	assert.Equal(t, int64(2), attacks[1].Ships.EspionageProbe)
}

func TestExtractAttacks_spy(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../../samples/unversioned/event_list_spy.html")
	attacks, _ := NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), nil)
	assert.Equal(t, 1, len(attacks))
	assert.Equal(t, ogame.Coordinate{4, 212, 8, ogame.PlanetType}, attacks[0].Origin)
	assert.Equal(t, int64(107009), attacks[0].AttackerID)
}

func TestExtractAttacks1(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../../samples/unversioned/event_list_missile.html")
	attacks, _ := NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), nil)
	assert.Equal(t, 1, len(attacks))
	assert.Equal(t, int64(1), attacks[0].Missiles)
	assert.Nil(t, attacks[0].Ships)
}
