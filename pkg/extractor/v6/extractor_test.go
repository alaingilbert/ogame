package v6

import (
	"github.com/alaingilbert/clockwork"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestExtractAttacks(t *testing.T) {
	clock := clockwork.NewFakeClockAt(time.Date(2016, 8, 23, 17, 48, 13, 0, time.UTC))
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/event_list_attack.html")
	attacks, _ := NewExtractor().extractAttacks(pageHTMLBytes, clock, nil)
	assert.Equal(t, 1, len(attacks))
	assert.Equal(t, "Homeworld", attacks[0].DestinationName)
	assert.Equal(t, clock.Now().Add(14*time.Minute), attacks[0].ArrivalTime.UTC())
	assert.Equal(t, int64(14*60), attacks[0].ArriveIn)
}

func TestExtractAttacksFromFullPage(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/overview_always_events.html")
	attacks, err := NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(attacks))
	assert.Equal(t, int64(1), attacks[0].Ships.SmallCargo)

	pageHTMLBytes, _ = os.ReadFile("../../../samples/unversioned/overview_active.html")
	_, err = NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), nil)
	assert.EqualError(t, err, ogame.ErrEventsBoxNotDisplayed.Error())

	pageHTMLBytes, _ = os.ReadFile("../../../samples/unversioned/eventlist_loggedout.html")
	_, err = NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), nil)
	assert.EqualError(t, err, ogame.ErrNotLogged.Error())
}

func TestExtractAttacksPhoneDisplay(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/event_list_attack_phone.html")
	attacks, _ := NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), nil)
	assert.Equal(t, 1, len(attacks))
	assert.Equal(t, int64(106734), attacks[0].AttackerID)
	assert.Equal(t, "", attacks[0].AttackerName, "should not be able to get the name")
}

func TestExtractAttacksMeAttacking(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/eventlist_me_attacking.html")
	attacks, _ := NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), nil)
	assert.Equal(t, 0, len(attacks))
}

func TestExtractAttacksWithoutShips(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/event_list_attack.html")
	attacks, _ := NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), nil)
	assert.Equal(t, 1, len(attacks))
	assert.Equal(t, int64(100771), attacks[0].AttackerID)
	assert.Equal(t, int64(0), attacks[0].Missiles)
	assert.Nil(t, attacks[0].Ships)
}

func TestExtractAttacksWithShips(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/eventList_attack_ships.html")
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
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/eventlist_moon_attacked.html")
	attacks, _ := NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), nil)
	assert.Equal(t, 1, len(attacks))
	assert.NotNil(t, attacks[0].Ships)
	assert.Equal(t, int64(107009), attacks[0].AttackerID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 212, Position: 8, Type: ogame.PlanetType}, attacks[0].Origin)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 116, Position: 12, Type: ogame.MoonType}, attacks[0].Destination)
	assert.Equal(t, ogame.MoonType, attacks[0].Destination.Type)
	assert.Equal(t, int64(1), attacks[0].Ships.SmallCargo)
	assert.Equal(t, "Moon", attacks[0].DestinationName)
}

func TestExtractAttacksMoonDestruction(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/eventlist_moon_destruction.html")
	attacks, _ := NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), nil)
	assert.Equal(t, 1, len(attacks))
	assert.NotNil(t, attacks[0].Ships)
	assert.Equal(t, int64(106734), attacks[0].AttackerID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 116, Position: 12, Type: ogame.PlanetType}, attacks[0].Origin)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 116, Position: 9, Type: ogame.MoonType}, attacks[0].Destination)
	assert.Equal(t, ogame.MoonType, attacks[0].Destination.Type)
	assert.Equal(t, int64(1), attacks[0].Ships.Deathstar)
}

func TestExtractAttacksWithThousandsShips(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/eventlist_attack_thousands.html")
	ownCoords := make([]ogame.Coordinate, 0)
	attacks, _ := NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), ownCoords)
	assert.Equal(t, 2, len(attacks))
	assert.Equal(t, int64(1012), attacks[1].Ships.Cruiser)
	assert.Equal(t, int64(1000), attacks[1].Ships.LargeCargo)
}

func TestExtractAttacksUnknownShips(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/eventlist_unknown_ships_nbr.html")
	attacks, _ := NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), nil)
	assert.Equal(t, 1, len(attacks))
	assert.Equal(t, int64(-1), attacks[0].Ships.Cruiser)
	assert.Equal(t, int64(0), attacks[0].Ships.Destroyer)
}

func TestExtractAttacksACS(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/eventlist_acs.html")
	attacks, _ := NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), nil)
	assert.Equal(t, 1, len(attacks))
	assert.Equal(t, ogame.GroupedAttack, attacks[0].MissionType)
	assert.Equal(t, int64(10), attacks[0].Ships.LightFighter)
	assert.Equal(t, int64(2176), attacks[0].Ships.Battlecruiser)
}

func TestExtractAttacksACSMany(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/eventlist_acs_multiple.html")
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
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/eventlist_acs2.html")
	attacks, _ := NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), nil)
	assert.Equal(t, 2, len(attacks))
	assert.Equal(t, ogame.GroupedAttack, attacks[0].MissionType)
	assert.Equal(t, int64(106734), attacks[0].AttackerID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 116, Position: 12, Type: ogame.PlanetType}, attacks[0].Origin)
	assert.Equal(t, int64(4), attacks[0].Ships.SmallCargo)
	assert.Equal(t, int64(3), attacks[0].Ships.Battlecruiser)
	assert.Equal(t, ogame.GroupedAttack, attacks[1].MissionType)
	assert.Equal(t, int64(106734), attacks[1].AttackerID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 116, Position: 12, Type: ogame.PlanetType}, attacks[1].Origin)
	assert.Equal(t, int64(7), attacks[1].Ships.SmallCargo)
	assert.Equal(t, int64(11), attacks[1].Ships.Battlecruiser)
	assert.Equal(t, int64(2), attacks[1].Ships.EspionageProbe)
}

func TestExtractAttacks_spy(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/event_list_spy.html")
	attacks, _ := NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), nil)
	assert.Equal(t, 1, len(attacks))
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 212, Position: 8, Type: ogame.PlanetType}, attacks[0].Origin)
	assert.Equal(t, int64(107009), attacks[0].AttackerID)
}

func TestExtractAttacks1(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/event_list_missile.html")
	attacks, _ := NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), nil)
	assert.Equal(t, 1, len(attacks))
	assert.Equal(t, int64(1), attacks[0].Missiles)
	assert.Nil(t, attacks[0].Ships)
}

func TestExtractLifeformEnabled(t *testing.T) {
	pageHTML, _ := os.ReadFile("../../../samples/unversioned/overview_active.html")
	assert.False(t, NewExtractor().ExtractLifeformEnabled(pageHTML))

	pageHTML, _ = os.ReadFile("../../../samples/v9.0.2/en/overview_all_queues.html")
	assert.False(t, NewExtractor().ExtractLifeformEnabled(pageHTML))

	pageHTML, _ = os.ReadFile("../../../samples/v9.0.2/en/lifeform/overview_all_queues.html")
	assert.True(t, NewExtractor().ExtractLifeformEnabled(pageHTML))
}

func TestExtractFleetDeutSaveFactor_V6_2_2_1(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/overview_active.html")
	res := NewExtractor().ExtractFleetDeutSaveFactor(pageHTMLBytes)
	assert.Equal(t, 1.0, res)
}

func TestExtractFleetDeutSaveFactor_V6_7_4(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/overview_with_moon.html")
	res := NewExtractor().ExtractFleetDeutSaveFactor(pageHTMLBytes)
	assert.Equal(t, 0.5, res)
}

func TestExtractPlanetCoordinate(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/station.html")
	res, _ := NewExtractor().ExtractPlanetCoordinate(pageHTMLBytes)
	assert.Equal(t, ogame.Coordinate{Galaxy: 1, System: 301, Position: 5, Type: ogame.PlanetType}, res)
}

func TestExtractPlanetCoordinate_moon(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/moon_facilities.html")
	res, _ := NewExtractor().ExtractPlanetCoordinate(pageHTMLBytes)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 116, Position: 12, Type: ogame.MoonType}, res)
}

func TestExtractPlanetID_planet(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/station.html")
	res, _ := NewExtractor().ExtractPlanetID(pageHTMLBytes)
	assert.Equal(t, ogame.CelestialID(33672410), res)
}

func TestExtractPlanetID_moon(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/moon_facilities.html")
	res, _ := NewExtractor().ExtractPlanetID(pageHTMLBytes)
	assert.Equal(t, ogame.CelestialID(33741598), res)
}

func TestExtractPlanetType_planet(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/station.html")
	res, _ := NewExtractor().ExtractPlanetType(pageHTMLBytes)
	assert.Equal(t, ogame.PlanetType, res)
}

func TestExtractPlanetType_moon(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/moon_facilities.html")
	res, _ := NewExtractor().ExtractPlanetType(pageHTMLBytes)
	assert.Equal(t, ogame.MoonType, res)
}

func TestExtractJumpGate_cooldown(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/jumpgatelayer_charge.html")
	_, _, _, wait := NewExtractor().ExtractJumpGate(pageHTMLBytes)
	assert.Equal(t, int64(1730), wait)
}

func TestExtractJumpGate(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/jumpgatelayer.html")
	ships, token, dests, wait := NewExtractor().ExtractJumpGate(pageHTMLBytes)
	assert.Equal(t, 1, len(dests))
	assert.Equal(t, ogame.MoonID(33743183), dests[0])
	assert.Equal(t, int64(0), wait)
	assert.Equal(t, "7787b530670bc89623b5d65a827e557a", token)
	assert.Equal(t, int64(0), ships.SmallCargo)
	assert.Equal(t, int64(101), ships.LargeCargo)
	assert.Equal(t, int64(1), ships.LightFighter)
}

func TestExtractOgameTimestamp(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/moon_facilities.html")
	res := NewExtractor().ExtractOgameTimestamp(pageHTMLBytes)
	assert.Equal(t, int64(1538912592), res)
}

func TestExtractOgameTimestampFromBytes(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/moon_facilities.html")
	res := NewExtractor().ExtractOGameTimestampFromBytes(pageHTMLBytes)
	assert.Equal(t, int64(1538912592), res)
}

func TestExtractResources(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/moon_facilities.html")
	res := NewExtractor().ExtractResources(pageHTMLBytes)
	assert.Equal(t, int64(280000), res.Metal)
	assert.Equal(t, int64(260000), res.Crystal)
	assert.Equal(t, int64(280000), res.Deuterium)
	assert.Equal(t, int64(0), res.Energy)
	assert.Equal(t, int64(25000), res.Darkmatter)
}

func TestExtractResourcesMobile(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/preferences_mobile.html")
	res := NewExtractor().ExtractResources(pageHTMLBytes)
	assert.Equal(t, int64(7325851), res.Metal)
	assert.Equal(t, int64(1695823), res.Crystal)
	assert.Equal(t, int64(1835627), res.Deuterium)
	assert.Equal(t, int64(-2827), res.Energy)
	assert.Equal(t, int64(19500), res.Darkmatter)
}

func TestExtractResourcesDetailsFromFullPage(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/fleets_1.html")
	res := NewExtractor().ExtractResourcesDetailsFromFullPage(pageHTMLBytes)
	assert.Equal(t, int64(1959227), res.Metal.Available)
	assert.Equal(t, int64(37818), res.Metal.CurrentProduction)
	assert.Equal(t, int64(5355000), res.Metal.StorageCapacity)
	assert.Equal(t, int64(327916), res.Crystal.Available)
	assert.Equal(t, int64(21862), res.Crystal.CurrentProduction)
	assert.Equal(t, int64(865000), res.Crystal.StorageCapacity)
	assert.Equal(t, int64(618155), res.Deuterium.Available)
	assert.Equal(t, int64(7508), res.Deuterium.CurrentProduction)
	assert.Equal(t, int64(865000), res.Deuterium.StorageCapacity)
	assert.Equal(t, int64(220), res.Energy.Available)
	assert.Equal(t, int64(17597), res.Energy.CurrentProduction)
	assert.Equal(t, int64(-17377), res.Energy.Consumption)
	assert.Equal(t, int64(25000), res.Darkmatter.Available)
	assert.Equal(t, int64(0), res.Darkmatter.Purchased)
	assert.Equal(t, int64(25000), res.Darkmatter.Found)
}

func TestExtractResourcesDetailsFromFullPageV7(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7/overview2.html")
	res := NewExtractor().ExtractResourcesDetailsFromFullPage(pageHTMLBytes)
	assert.Equal(t, int64(36800), res.Metal.Available)
	assert.Equal(t, int64(396), res.Metal.CurrentProduction)
	assert.Equal(t, int64(40000), res.Metal.StorageCapacity)
	assert.Equal(t, int64(56524), res.Crystal.Available)
	assert.Equal(t, int64(143), res.Crystal.CurrentProduction)
	assert.Equal(t, int64(75000), res.Crystal.StorageCapacity)
	assert.Equal(t, int64(18401), res.Deuterium.Available)
	assert.Equal(t, int64(128), res.Deuterium.CurrentProduction)
	assert.Equal(t, int64(20000), res.Deuterium.StorageCapacity)
	assert.Equal(t, int64(-4), res.Energy.Available)
	assert.Equal(t, int64(79), res.Energy.CurrentProduction)
	assert.Equal(t, int64(-83), res.Energy.Consumption)
	assert.Equal(t, int64(19348523), res.Darkmatter.Available)
	assert.Equal(t, int64(0), res.Darkmatter.Purchased)
	assert.Equal(t, int64(19348523), res.Darkmatter.Found)
}

func TestExtractPhalanx_75(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.5.1/en/phalanx_returning.html")
	res, err := NewExtractor().ExtractPhalanx(pageHTMLBytes)
	clock := clockwork.NewFakeClockAt(time.Date(2020, 11, 4, 0, 25, 29, 0, time.UTC))
	assert.Nil(t, err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, ogame.Transport, res[0].Mission)
	assert.Equal(t, true, res[0].ReturnFlight)
	assert.NotNil(t, res[0].ArriveIn)
	assert.Equal(t, clock.Now().Add(10*time.Minute), res[0].ArrivalTime.UTC())
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 116, Position: 9, Type: ogame.PlanetType}, res[0].Origin)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 116, Position: 10, Type: ogame.PlanetType}, res[0].Destination)
	assert.Equal(t, int64(19), res[0].Ships.SmallCargo)
}

func TestExtractPhalanx(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/phalanx.html")
	res, err := NewExtractor().ExtractPhalanx(pageHTMLBytes)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, ogame.MissionID(3), res[0].Mission)
	assert.Equal(t, true, res[0].ReturnFlight)
	assert.NotNil(t, res[0].ArriveIn)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 116, Position: 9, Type: ogame.PlanetType}, res[0].Origin)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 212, Position: 8, Type: ogame.PlanetType}, res[0].Destination)
	assert.Equal(t, int64(100), res[0].Ships.LargeCargo)
}

func TestExtractPhalanx_fromMoon(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/phalanx_from_moon.html")
	res, _ := NewExtractor().ExtractPhalanx(pageHTMLBytes)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 116, Position: 12, Type: ogame.MoonType}, res[0].Origin)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 116, Position: 9, Type: ogame.PlanetType}, res[0].Destination)
}

func TestExtractPhalanx_manyFleets(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/phalanx_fleets.html")
	res, err := NewExtractor().ExtractPhalanx(pageHTMLBytes)
	assert.Nil(t, err)
	assert.Equal(t, 12, len(res))
	assert.Equal(t, ogame.Expedition, res[0].Mission)
	assert.False(t, res[0].ReturnFlight)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 124, Position: 9, Type: ogame.PlanetType}, res[0].Origin)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 125, Position: 16, Type: ogame.PlanetType}, res[0].Destination)
	assert.Equal(t, int64(250), res[0].Ships.LargeCargo)
	assert.Equal(t, int64(1), res[0].Ships.EspionageProbe)
	assert.Equal(t, int64(1), res[0].Ships.Destroyer)

	assert.Equal(t, ogame.Expedition, res[8].Mission)
	assert.True(t, res[8].ReturnFlight)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 124, Position: 9, Type: ogame.PlanetType}, res[8].Origin)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 125, Position: 16, Type: ogame.PlanetType}, res[8].Destination)
	assert.Equal(t, int64(250), res[8].Ships.LargeCargo)
	assert.Equal(t, int64(1), res[8].Ships.EspionageProbe)
	assert.Equal(t, int64(1), res[8].Ships.Destroyer)
}

func TestExtractPhalanx_noFleet(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/phalanx_no_fleet.html")
	res, err := NewExtractor().ExtractPhalanx(pageHTMLBytes)
	assert.Equal(t, 0, len(res))
	assert.Nil(t, err)
}

func TestExtractPhalanx_noDeut(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/phalanx_no_deut.html")
	res, err := NewExtractor().ExtractPhalanx(pageHTMLBytes)
	assert.Equal(t, 0, len(res))
	assert.NotNil(t, err)
}

func TestExtractResearch(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/research_bonus.html")
	res := NewExtractor().ExtractResearch(pageHTMLBytes)
	assert.Equal(t, int64(12), res.EnergyTechnology)
	assert.Equal(t, int64(12), res.LaserTechnology)
	assert.Equal(t, int64(7), res.IonTechnology)
	assert.Equal(t, int64(6), res.HyperspaceTechnology)
	assert.Equal(t, int64(7), res.PlasmaTechnology)
	assert.Equal(t, int64(15), res.CombustionDrive)
	assert.Equal(t, int64(7), res.ImpulseDrive)
	assert.Equal(t, int64(8), res.HyperspaceDrive)
	assert.Equal(t, int64(10), res.EspionageTechnology)
	assert.Equal(t, int64(14), res.ComputerTechnology)
	assert.Equal(t, int64(13), res.Astrophysics)
	assert.Equal(t, int64(0), res.IntergalacticResearchNetwork)
	assert.Equal(t, int64(0), res.GravitonTechnology)
	assert.Equal(t, int64(13), res.WeaponsTechnology)
	assert.Equal(t, int64(12), res.ShieldingTechnology)
	assert.Equal(t, int64(12), res.ArmourTechnology)
}

func TestExtractResourcesBuildings(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/resource_inconstruction.html")
	res, _ := NewExtractor().ExtractResourcesBuildings(pageHTMLBytes)
	assert.Equal(t, int64(19), res.MetalMine)
	assert.Equal(t, int64(17), res.CrystalMine)
	assert.Equal(t, int64(13), res.DeuteriumSynthesizer)
	assert.Equal(t, int64(20), res.SolarPlant)
	assert.Equal(t, int64(3), res.FusionReactor)
	assert.Equal(t, int64(0), res.SolarSatellite)
	assert.Equal(t, int64(5), res.MetalStorage)
	assert.Equal(t, int64(4), res.CrystalStorage)
	assert.Equal(t, int64(3), res.DeuteriumTank)
}

func TestExtractFacilities(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/facility_inconstruction.html")
	res, _ := NewExtractor().ExtractFacilities(pageHTMLBytes)
	assert.Equal(t, int64(7), res.RoboticsFactory)
	assert.Equal(t, int64(7), res.Shipyard)
	assert.Equal(t, int64(7), res.ResearchLab)
	assert.Equal(t, int64(0), res.AllianceDepot)
	assert.Equal(t, int64(0), res.MissileSilo)
	assert.Equal(t, int64(0), res.NaniteFactory)
	assert.Equal(t, int64(0), res.Terraformer)
	assert.Equal(t, int64(3), res.SpaceDock)
}

func TestExtractMoonFacilities(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/moon_facilities.html")
	res, _ := NewExtractor().ExtractFacilities(pageHTMLBytes)
	assert.Equal(t, int64(1), res.RoboticsFactory)
	assert.Equal(t, int64(2), res.Shipyard)
	assert.Equal(t, int64(3), res.LunarBase)
	assert.Equal(t, int64(4), res.SensorPhalanx)
	assert.Equal(t, int64(5), res.JumpGate)
}

func TestExtractDefense(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/defence.html")
	defense, _ := NewExtractor().ExtractDefense(pageHTMLBytes)
	assert.Equal(t, int64(1), defense.RocketLauncher)
	assert.Equal(t, int64(2), defense.LightLaser)
	assert.Equal(t, int64(3), defense.HeavyLaser)
	assert.Equal(t, int64(4), defense.GaussCannon)
	assert.Equal(t, int64(5), defense.IonCannon)
	assert.Equal(t, int64(6), defense.PlasmaTurret)
	assert.Equal(t, int64(0), defense.SmallShieldDome)
	assert.Equal(t, int64(0), defense.LargeShieldDome)
	assert.Equal(t, int64(7), defense.AntiBallisticMissiles)
	assert.Equal(t, int64(8), defense.InterplanetaryMissiles)
}

func TestExtractFleet1Ships(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/fleet1.html")
	s := NewExtractor().ExtractFleet1Ships(pageHTMLBytes)
	assert.Equal(t, int64(3), s.LightFighter)
	assert.Equal(t, int64(0), s.HeavyFighter)
	assert.Equal(t, int64(1012), s.Cruiser)
	assert.Equal(t, int64(0), s.Battleship)
	assert.Equal(t, int64(0), s.SmallCargo)
	assert.Equal(t, int64(1003), s.LargeCargo)
	assert.Equal(t, int64(1), s.ColonyShip)
	assert.Equal(t, int64(200), s.Battlecruiser)
	assert.Equal(t, int64(100), s.Bomber)
	assert.Equal(t, int64(200), s.Destroyer)
	assert.Equal(t, int64(0), s.Deathstar)
	assert.Equal(t, int64(30), s.Recycler)
	assert.Equal(t, int64(1001), s.EspionageProbe)
}

func TestExtractFleet1Ships_NoShips(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/fleet1_no_ships.html")
	s := NewExtractor().ExtractFleet1Ships(pageHTMLBytes)
	assert.Equal(t, int64(0), s.LightFighter)
	assert.Equal(t, int64(0), s.HeavyFighter)
	assert.Equal(t, int64(0), s.Cruiser)
	assert.Equal(t, int64(0), s.Battleship)
	assert.Equal(t, int64(0), s.SmallCargo)
	assert.Equal(t, int64(0), s.LargeCargo)
	assert.Equal(t, int64(0), s.ColonyShip)
	assert.Equal(t, int64(0), s.Battlecruiser)
	assert.Equal(t, int64(0), s.Bomber)
	assert.Equal(t, int64(0), s.Destroyer)
	assert.Equal(t, int64(0), s.Deathstar)
	assert.Equal(t, int64(0), s.Recycler)
	assert.Equal(t, int64(0), s.EspionageProbe)
}

func TestExtractPlanet_en(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/overview_queues.html")
	planet, _ := NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33677371))
	assert.Equal(t, "C1", planet.Name)
	assert.Equal(t, int64(14615), planet.Diameter)
	assert.Equal(t, int64(-2), planet.Temperature.Min)
	assert.Equal(t, int64(38), planet.Temperature.Max)
	assert.Equal(t, int64(35), planet.Fields.Built)
	assert.Equal(t, int64(238), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33677371), planet.ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 1, System: 301, Position: 8, Type: ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_fr(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/fr_overview.html")
	planet, _ := NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33629512))
	assert.Equal(t, "planète mère", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(48), planet.Temperature.Min)
	assert.Equal(t, int64(88), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(188), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33629512), planet.ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 2, System: 180, Position: 4, Type: ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_de(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/de_overview.html")
	planet, _ := NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33630447))
	assert.Equal(t, "Heimatplanet", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(21), planet.Temperature.Min)
	assert.Equal(t, int64(61), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(188), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33630447), planet.ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 2, System: 175, Position: 8, Type: ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_dk(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/dk_overview.html")
	planet, _ := NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33627426))
	assert.Equal(t, "Hjemme verden", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(-23), planet.Temperature.Min)
	assert.Equal(t, int64(17), planet.Temperature.Max)
	assert.Equal(t, int64(5), planet.Fields.Built)
	assert.Equal(t, int64(193), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33627426), planet.ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 1, System: 148, Position: 12, Type: ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_es(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v6/es/shipyard.html")
	planet, _ := NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33644981))
	assert.Equal(t, "Planeta Principal", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(-18), planet.Temperature.Min)
	assert.Equal(t, int64(22), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(193), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33644981), planet.ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 2, System: 493, Position: 10, Type: ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_br(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v6/br/overview.html")
	planet, _ := NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33633767))
	assert.Equal(t, "Planeta Principal", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(-13), planet.Temperature.Min)
	assert.Equal(t, int64(27), planet.Temperature.Max)
	assert.Equal(t, int64(5), planet.Fields.Built)
	assert.Equal(t, int64(193), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33633767), planet.ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 1, System: 449, Position: 12, Type: ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_it(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v6/it/overview.html")
	planet, _ := NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33634944))
	assert.Equal(t, "Pianeta Madre", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(28), planet.Temperature.Min)
	assert.Equal(t, int64(68), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(188), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33634944), planet.ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 2, System: 58, Position: 8, Type: ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_jp(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/jp_overview.html")
	planet, _ := NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33620484))
	assert.Equal(t, "ホームワールド", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(40), planet.Temperature.Min)
	assert.Equal(t, int64(80), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(188), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33620484), planet.ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 1, System: 18, Position: 4, Type: ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_tw(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v6/tw/overview.html")
	planet, _ := NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33626432))
	assert.Equal(t, "母星", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(29), planet.Temperature.Min)
	assert.Equal(t, int64(69), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(188), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33626432), planet.ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 1, System: 206, Position: 8, Type: ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_hr(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v6/hr/overview.html")
	planet, _ := NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33627961))
	assert.Equal(t, "Glavni Planet", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(-33), planet.Temperature.Min)
	assert.Equal(t, int64(7), planet.Temperature.Max)
	assert.Equal(t, int64(4), planet.Fields.Built)
	assert.Equal(t, int64(193), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33627961), planet.ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 1, System: 236, Position: 12, Type: ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_no(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v6/no/overview.html")
	planet, _ := NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33624646))
	assert.Equal(t, "Hjemmeverden", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(-24), planet.Temperature.Min)
	assert.Equal(t, int64(16), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(188), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33624646), planet.ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 1, System: 99, Position: 10, Type: ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_sk(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v6/sk/overview.html")
	planet, _ := NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33625241))
	assert.Equal(t, "Domovská planéta", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(-12), planet.Temperature.Min)
	assert.Equal(t, int64(28), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(163), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33625241), planet.ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 1, System: 94, Position: 10, Type: ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_si(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.2/si/overview.html")
	planet, _ := NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33625245))
	assert.Equal(t, "Glavni Planet", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(41), planet.Temperature.Min)
	assert.Equal(t, int64(81), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(188), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33625245), planet.ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 1, System: 70, Position: 6, Type: ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_hu(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.2/hu/overview.html")
	planet, _ := NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33621505))
	assert.Equal(t, "Otthon", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(-18), planet.Temperature.Min)
	assert.Equal(t, int64(22), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(188), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33621505), planet.ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 1, System: 162, Position: 12, Type: ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_fi(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.2/fi/overview.html")
	planet, _ := NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33625483))
	assert.Equal(t, "Kotimaailma", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(15), planet.Temperature.Min)
	assert.Equal(t, int64(55), planet.Temperature.Max)
	assert.Equal(t, int64(5), planet.Fields.Built)
	assert.Equal(t, int64(188), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33625483), planet.ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 1, System: 94, Position: 6, Type: ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_ba(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.5.1/ba/overview.html")
	planet, _ := NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33621433))
	assert.Equal(t, "Glavni Planet", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(70), planet.Temperature.Min)
	assert.Equal(t, int64(110), planet.Temperature.Max)
	assert.Equal(t, int64(5), planet.Fields.Built)
	assert.Equal(t, int64(193), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33621433), planet.ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 1, System: 55, Position: 4, Type: ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_gr(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v6/gr/overview.html")
	planet, _ := NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33629206))
	assert.Equal(t, "Κύριος Πλανήτης", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(37), planet.Temperature.Min)
	assert.Equal(t, int64(77), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(193), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33629206), planet.ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 1, System: 182, Position: 6, Type: ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_mx(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v6/mx/overview.html")
	planet, _ := NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33624669))
	assert.Equal(t, "Planeta Principal", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(33), planet.Temperature.Min)
	assert.Equal(t, int64(73), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(188), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33624669), planet.ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 1, System: 390, Position: 6, Type: ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_cz(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v6/cz/overview.html")
	planet, _ := NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33622822))
	assert.Equal(t, "Domovska planeta", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(-13), planet.Temperature.Min)
	assert.Equal(t, int64(27), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(188), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33622822), planet.ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 1, System: 221, Position: 12, Type: ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_jp1(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v6/jp/overview.html")
	planet, _ := NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33623513))
	assert.Equal(t, "ホームワールド", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(58), planet.Temperature.Min)
	assert.Equal(t, int64(98), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(188), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33623513), planet.ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 1, System: 85, Position: 4, Type: ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_pl(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/pl_overview.html")
	planet, _ := NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33669699))
	assert.Equal(t, "Planeta matka", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(-21), planet.Temperature.Min)
	assert.Equal(t, int64(19), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(188), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33669699), planet.ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 248, Position: 12, Type: ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_tr(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/tr_overview.html")
	planet, _ := NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33650421))
	assert.Equal(t, "Ana Gezegen", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(9), planet.Temperature.Min)
	assert.Equal(t, int64(49), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(188), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33650421), planet.ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 3, System: 143, Position: 10, Type: ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_pt(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/pt_overview.html")
	planet, _ := NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33635398))
	assert.Equal(t, "Planeta Principal", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(33), planet.Temperature.Min)
	assert.Equal(t, int64(73), planet.Temperature.Max)
	assert.Equal(t, int64(4), planet.Fields.Built)
	assert.Equal(t, int64(193), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33635398), planet.ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 2, System: 241, Position: 6, Type: ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_nl(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/nl_overview.html")
	planet, _ := NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33624684))
	assert.Equal(t, "Hoofdplaneet", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(-47), planet.Temperature.Min)
	assert.Equal(t, int64(-7), planet.Temperature.Max)
	assert.Equal(t, int64(5), planet.Fields.Built)
	assert.Equal(t, int64(188), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33624684), planet.ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 1, System: 178, Position: 12, Type: ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_ar(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v6/ar/overview.html")
	planet, _ := NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33629527))
	assert.Equal(t, "Planeta Principal", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(37), planet.Temperature.Min)
	assert.Equal(t, int64(77), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(193), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33629527), planet.ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 1, System: 367, Position: 4, Type: ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_ru(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v6/ru/overview.html")
	planet, _ := NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33629521))
	assert.Equal(t, "Главная планета", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(23), planet.Temperature.Min)
	assert.Equal(t, int64(63), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(163), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33629521), planet.ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 1, System: 374, Position: 6, Type: ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_notExists(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/overview_queues.html")
	_, err := NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(12345))
	assert.NotNil(t, err)
}

func TestExtractPlanetByCoord(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/overview_queues.html")
	planet, _ := NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.Coordinate{Galaxy: 1, System: 301, Position: 8, Type: ogame.PlanetType})
	assert.Equal(t, "C1", planet.Name)
	assert.Equal(t, int64(14615), planet.Diameter)
}

func TestExtractPlanetByCoord_notExists(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/overview_queues.html")
	_, err := NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.Coordinate{Galaxy: 1, System: 2, Position: 3, Type: ogame.PlanetType})
	assert.NotNil(t, err)
}

func TestExtractShips(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/shipyard_thousands_ships.html")
	ships, _ := NewExtractor().ExtractShips(pageHTMLBytes)
	assert.Equal(t, int64(1000), ships.LargeCargo)
	assert.Equal(t, int64(1000), ships.EspionageProbe)
	assert.Equal(t, int64(700), ships.Cruiser)
}

func TestExtractShipsMillions(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/shipyard_millions_ships.html")
	ships, _ := NewExtractor().ExtractShips(pageHTMLBytes)
	assert.Equal(t, int64(15000001), ships.LightFighter)
}

func TestExtractShipsWhileBeingBuilt(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/shipyard_ship_being_built.html")
	ships, _ := NewExtractor().ExtractShips(pageHTMLBytes)
	assert.Equal(t, int64(213), ships.EspionageProbe)
}

func TestExtractEspionageReportMessageIDs(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/messages.html")
	msgs, _, _ := NewExtractor().ExtractEspionageReportMessageIDs(pageHTMLBytes)
	assert.Equal(t, 2, len(msgs))
	assert.Equal(t, ogame.Report, msgs[0].Type)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 117, Position: 6, Type: ogame.PlanetType}, msgs[0].Target)
	assert.Equal(t, 0.5, msgs[0].LootPercentage)
	assert.Equal(t, "Fleet Command", msgs[0].From)
	assert.Equal(t, ogame.Action, msgs[1].Type)
	assert.Equal(t, "Space Monitoring", msgs[1].From)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 117, Position: 9, Type: ogame.PlanetType}, msgs[1].Target)
}

func TestExtractEspionageReportMessageIDsLootPercentage(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/messages_loot_percentage.html")
	msgs, _, _ := NewExtractor().ExtractEspionageReportMessageIDs(pageHTMLBytes)
	assert.Equal(t, 1.0, msgs[0].LootPercentage)
	assert.Equal(t, 0.5, msgs[1].LootPercentage)
	assert.Equal(t, 0.5, msgs[2].LootPercentage)
}

func TestExtractCombatReportMessages(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/combat_reports_msgs.html")
	msgs, _, _ := NewExtractor().ExtractCombatReportMessagesSummary(pageHTMLBytes)
	assert.Equal(t, 9, len(msgs))
}

func TestExtractCombatReportAttackingMessages(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/combat_reports_msgs_attacking.html")
	msgs, _, _ := NewExtractor().ExtractCombatReportMessagesSummary(pageHTMLBytes)
	assert.Equal(t, int64(7945368), msgs[0].ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 233, Position: 11, Type: ogame.PlanetType}, msgs[0].Destination)
	assert.Equal(t, int64(50), msgs[0].Loot)
	assert.Equal(t, int64(74495), msgs[0].Metal)
	assert.Equal(t, int64(88280), msgs[0].Crystal)
	assert.Equal(t, int64(21572), msgs[0].Deuterium)
	assert.Equal(t, int64(3500), msgs[0].DebrisField)
	assert.Equal(t, int64(25200), msgs[1].DebrisField)
	assert.Equal(t, int64(0), msgs[2].DebrisField)
	assert.Equal(t, "08.09.2018 09:33:18", msgs[0].CreatedAt.Format("02.01.2006 15:04:05"))
}

func TestExtractCombatReportMessagesSummary(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/combat_reports_msgs_2.html")
	msgs, nbPages, _ := NewExtractor().ExtractCombatReportMessagesSummary(pageHTMLBytes)
	assert.Equal(t, 10, len(msgs))
	assert.Equal(t, int64(44), nbPages)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 116, Position: 12, Type: ogame.MoonType}, msgs[1].Destination)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 127, Position: 9, Type: ogame.MoonType}, *msgs[1].Origin)
}

func TestExtractResourcesProductions(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/resource_settings.html")
	prods, _ := NewExtractor().ExtractResourcesProductions(pageHTMLBytes)
	assert.Equal(t, ogame.Resources{Metal: 10352, Crystal: 5104, Deuterium: 1282, Energy: -52}, prods)
}

func TestExtractResourceSettings(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/resource_settings.html")
	settings, _, _ := NewExtractor().ExtractResourceSettings(pageHTMLBytes)
	assert.Equal(t, ogame.ResourceSettings{MetalMine: 100, CrystalMine: 100, DeuteriumSynthesizer: 100, SolarPlant: 100, FusionReactor: 0, SolarSatellite: 100}, settings)
}

func TestExtractNbProbes(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/preferences.html")
	probes := NewExtractor().ExtractSpioAnz(pageHTMLBytes)
	assert.Equal(t, int64(10), probes)

	pageHTMLBytes, _ = os.ReadFile("../../../samples/unversioned/preferences_mobile.html")
	probes = NewExtractor().ExtractSpioAnz(pageHTMLBytes)
	assert.Equal(t, int64(3), probes)
}

func TestExtractPreferencesShowActivityMinutes(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/preferences.html")
	checked := NewExtractor().ExtractPreferencesShowActivityMinutes(pageHTMLBytes)
	assert.True(t, checked)

	pageHTMLBytes, _ = os.ReadFile("../../../samples/unversioned/preferences_mobile.html")
	checked = NewExtractor().ExtractPreferencesShowActivityMinutes(pageHTMLBytes)
	assert.True(t, checked)

	pageHTMLBytes, _ = os.ReadFile("../../../samples/unversioned/preferences_without_detailed_activities.html")
	checked = NewExtractor().ExtractPreferencesShowActivityMinutes(pageHTMLBytes)
	assert.False(t, checked)
}

func TestExtractPreferences(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/preferences.html")
	prefs := NewExtractor().ExtractPreferences(pageHTMLBytes)
	assert.Equal(t, int64(10), prefs.SpioAnz)
	assert.False(t, prefs.UrlaubsModus)
	assert.False(t, prefs.DisableChatBar)
	assert.False(t, prefs.DisableOutlawWarning)
	assert.False(t, prefs.MobileVersion)
	assert.False(t, prefs.ShowOldDropDowns)
	assert.False(t, prefs.ActivateAutofocus)
	assert.Equal(t, int64(1), prefs.EventsShow)
	assert.Equal(t, int64(0), prefs.SortSetting)
	assert.Equal(t, int64(0), prefs.SortOrder)
	assert.True(t, prefs.ShowDetailOverlay)
	assert.True(t, prefs.AnimatedSliders)
	assert.True(t, prefs.AnimatedOverview)
	assert.False(t, prefs.PopupsNotices)
	assert.False(t, prefs.PopopsCombatreport)
	assert.False(t, prefs.SpioReportPictures)
	assert.Equal(t, int64(10), prefs.MsgResultsPerPage)
	assert.True(t, prefs.AuctioneerNotifications)
	assert.False(t, prefs.EconomyNotifications)
	assert.True(t, prefs.ShowActivityMinutes)
	assert.False(t, prefs.PreserveSystemOnPlanetChange)

	pageHTMLBytes, _ = os.ReadFile("../../../samples/unversioned/preferences_reverse.html")
	prefs = NewExtractor().ExtractPreferences(pageHTMLBytes)
	assert.Equal(t, int64(2), prefs.SpioAnz)
	assert.False(t, prefs.UrlaubsModus)
	assert.True(t, prefs.DisableChatBar)
	assert.True(t, prefs.DisableOutlawWarning)
	assert.False(t, prefs.MobileVersion)
	assert.True(t, prefs.ShowOldDropDowns)
	assert.True(t, prefs.ActivateAutofocus)
	assert.Equal(t, int64(3), prefs.EventsShow)
	assert.Equal(t, int64(3), prefs.SortSetting)
	assert.Equal(t, int64(1), prefs.SortOrder)
	assert.False(t, prefs.ShowDetailOverlay)
	assert.False(t, prefs.AnimatedSliders)
	assert.False(t, prefs.AnimatedOverview)
	assert.True(t, prefs.PopupsNotices)
	assert.True(t, prefs.PopopsCombatreport)
	assert.True(t, prefs.SpioReportPictures)
	assert.Equal(t, int64(50), prefs.MsgResultsPerPage)
	assert.False(t, prefs.AuctioneerNotifications)
	assert.True(t, prefs.EconomyNotifications)
	assert.False(t, prefs.ShowActivityMinutes)
	assert.True(t, prefs.PreserveSystemOnPlanetChange)

	pageHTMLBytes, _ = os.ReadFile("../../../samples/unversioned/preferences_mobile.html")
	prefs = NewExtractor().ExtractPreferences(pageHTMLBytes)
	assert.Equal(t, int64(3), prefs.SpioAnz)
	assert.False(t, prefs.UrlaubsModus)
	assert.False(t, prefs.DisableChatBar) // no mobile
	assert.False(t, prefs.DisableOutlawWarning)
	assert.True(t, prefs.MobileVersion)
	assert.False(t, prefs.ShowOldDropDowns)
	assert.False(t, prefs.ActivateAutofocus)
	assert.Equal(t, int64(2), prefs.EventsShow)
	assert.Equal(t, int64(0), prefs.SortSetting)
	assert.Equal(t, int64(0), prefs.SortOrder)
	assert.True(t, prefs.ShowDetailOverlay)
	assert.False(t, prefs.AnimatedSliders)    // no mobile
	assert.False(t, prefs.AnimatedOverview)   // no mobile
	assert.False(t, prefs.PopupsNotices)      // no mobile
	assert.False(t, prefs.PopopsCombatreport) // no mobile
	assert.False(t, prefs.SpioReportPictures)
	assert.Equal(t, int64(10), prefs.MsgResultsPerPage)
	assert.True(t, prefs.AuctioneerNotifications)
	assert.False(t, prefs.EconomyNotifications)
	assert.True(t, prefs.ShowActivityMinutes)
	assert.False(t, prefs.PreserveSystemOnPlanetChange)

	//assert.True(t, prefs.Notifications.BuildList)
	//assert.True(t, prefs.Notifications.FriendlyFleetActivities)
	//assert.True(t, prefs.Notifications.HostileFleetActivities)
	//assert.True(t, prefs.Notifications.ForeignEspionage)
	//assert.True(t, prefs.Notifications.AllianceBroadcasts)
	//assert.True(t, prefs.Notifications.AllianceMessages)
	//assert.True(t, prefs.Notifications.Auctions)
	//assert.True(t, prefs.Notifications.Account)

	pageHTMLBytes, _ = os.ReadFile("../../../samples/unversioned/preferences_reverse_mobile.html")
	prefs = NewExtractor().ExtractPreferences(pageHTMLBytes)
	assert.Equal(t, int64(2), prefs.SpioAnz)
	assert.False(t, prefs.UrlaubsModus)
	assert.False(t, prefs.DisableChatBar) // no mobile
	assert.True(t, prefs.DisableOutlawWarning)
	assert.True(t, prefs.MobileVersion)
	assert.True(t, prefs.ShowOldDropDowns)
	assert.True(t, prefs.ActivateAutofocus)
	assert.Equal(t, int64(3), prefs.EventsShow)
	assert.Equal(t, int64(3), prefs.SortSetting)
	assert.Equal(t, int64(1), prefs.SortOrder)
	assert.False(t, prefs.ShowDetailOverlay)
	assert.False(t, prefs.AnimatedSliders)    // no mobile
	assert.False(t, prefs.AnimatedOverview)   // no mobile
	assert.False(t, prefs.PopupsNotices)      // no mobile
	assert.False(t, prefs.PopopsCombatreport) // no mobile
	assert.True(t, prefs.SpioReportPictures)
	assert.Equal(t, int64(50), prefs.MsgResultsPerPage)
	assert.False(t, prefs.AuctioneerNotifications)
	assert.True(t, prefs.EconomyNotifications)
	assert.False(t, prefs.ShowActivityMinutes)
	assert.True(t, prefs.PreserveSystemOnPlanetChange)

	//assert.False(t, prefs.Notifications.BuildList)
	//assert.False(t, prefs.Notifications.FriendlyFleetActivities)
	//assert.False(t, prefs.Notifications.HostileFleetActivities)
	//assert.False(t, prefs.Notifications.ForeignEspionage)
	//assert.False(t, prefs.Notifications.AllianceBroadcasts)
	//assert.False(t, prefs.Notifications.AllianceMessages)
	//assert.False(t, prefs.Notifications.Auctions)
	//assert.False(t, prefs.Notifications.Account)
}

//func TestCalcResources(t *testing.T) {
//	pageHTMLBytes, _ := os.ReadFile("../../../samples/traderOverview.html")
//	price, _, planetResources, multiplier, _ := NewExtractor().ExtractOfferOfTheDay(pageHTMLBytes)
//	actual := calcResources(price, planetResources, multiplier)
//	expected := url.Values{
//		"bid[planets][33711028][crystal]":   []string{"0"},
//		"bid[planets][33711028][deuterium]": []string{"0"},
//		"bid[planets][33711028][metal]":     []string{"54243"},
//		"bid[planets][33738397][crystal]":   []string{"0"},
//		"bid[planets][33738397][deuterium]": []string{"0"},
//		"bid[planets][33738397][metal]":     []string{"0"},
//		"bid[planets][33738457][crystal]":   []string{"0"},
//		"bid[planets][33738457][deuterium]": []string{"0"},
//		"bid[planets][33738457][metal]":     []string{"0"},
//		"bid[planets][33739506][crystal]":   []string{"0"},
//		"bid[planets][33739506][deuterium]": []string{"0"},
//		"bid[planets][33739506][metal]":     []string{"0"},
//		"bid[planets][33760932][crystal]":   []string{"0"},
//		"bid[planets][33760932][deuterium]": []string{"0"},
//		"bid[planets][33760932][metal]":     []string{"0"},
//		"bid[planets][33760935][crystal]":   []string{"0"},
//		"bid[planets][33760935][deuterium]": []string{"0"},
//		"bid[planets][33760935][metal]":     []string{"0"},
//		"bid[planets][33760958][crystal]":   []string{"0"},
//		"bid[planets][33760958][deuterium]": []string{"0"},
//		"bid[planets][33760958][metal]":     []string{"0"},
//		"bid[planets][33762073][crystal]":   []string{"0"},
//		"bid[planets][33762073][deuterium]": []string{"0"},
//		"bid[planets][33762073][metal]":     []string{"0"},
//		"bid[planets][33765791][crystal]":   []string{"0"},
//		"bid[planets][33765791][deuterium]": []string{"0"},
//		"bid[planets][33765791][metal]":     []string{"0"},
//		"bid[planets][33792134][crystal]":   []string{"0"},
//		"bid[planets][33792134][deuterium]": []string{"0"},
//		"bid[planets][33792134][metal]":     []string{"0"}}
//	assert.Equal(t, expected, actual)
//}

func TestExtractOfferOfTheDayPrice(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/traderOverview.html")
	price, token, _, _, _ := NewExtractor().ExtractOfferOfTheDay(pageHTMLBytes)
	assert.Equal(t, int64(54243), price)
	assert.Equal(t, "8128c0ba0c9981599a87d818003f95e1", token)
}

func TestExtractOfferOfTheDayPrice1(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.4/en/traderOverview.html")
	price, token, _, _, _ := NewExtractor().ExtractOfferOfTheDay(pageHTMLBytes)
	assert.Equal(t, int64(822159), price)
	assert.Equal(t, "2c829372796443bf6994cbfa051e4cd2", token)
}

func TestExtractGalaxyInfos_vacationMode(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7/galaxy_vacation_mode.html")
	_, err := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.EqualError(t, err, "account in vacation mode")
}

func TestExtractGalaxyInfos_bandit(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/galaxy_inactive_bandit_lord.html")
	infos, _ := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.True(t, infos.Position(6).Player.IsBandit)
	assert.False(t, infos.Position(6).Player.IsStarlord)
}

func TestExtractGalaxyInfos_starlord(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/galaxy_inactive_emperor.html")
	infos, _ := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.True(t, infos.Position(7).Player.IsStarlord)
	assert.False(t, infos.Position(7).Player.IsBandit)
}

func TestExtractGalaxyInfos_destroyedPlanet(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/galaxy_destroyed_planet.html")
	infos, _ := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.NotNil(t, infos.Position(8))
	assert.True(t, infos.Position(8).Destroyed)
}

func TestExtractGalaxyInfos_destroyedPlanetAndMoon(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/galaxy_destroyed_planet_and_moon2.html")
	infos, _ := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.NotNil(t, infos.Position(15))
	assert.True(t, infos.Position(15).Destroyed)
}

func TestExtractGalaxyInfos_banned(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/galaxy_banned.html")
	infos, _ := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, true, infos.Position(1).Banned)
	assert.Equal(t, false, infos.Position(9).Banned)
}

func TestExtractGalaxyInfos(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/galaxy_ajax.html")
	infos, _ := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, int64(4), infos.Galaxy())
	assert.Equal(t, int64(116), infos.System())
	assert.Equal(t, int64(33698600), infos.Position(4).ID)
	assert.Equal(t, int64(33698645), infos.Position(6).ID)
	assert.Equal(t, int64(106733), infos.Position(6).Player.ID)
	assert.Equal(t, "Origin", infos.Position(6).Player.Name)
	assert.Equal(t, int64(1671), infos.Position(6).Player.Rank)
	assert.Equal(t, "Ra", infos.Position(6).Name)
}

func TestExtractGalaxyInfosOwnPlanet(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/galaxy_ajax.html")
	infos, _ := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, int64(33698658), infos.Position(12).ID)
	assert.Equal(t, "Commodore Nomade", infos.Position(12).Player.Name)
	assert.Equal(t, int64(123), infos.Position(12).Player.ID)
	assert.Equal(t, int64(456), infos.Position(12).Player.Rank)
	assert.Equal(t, "Homeworld", infos.Position(12).Name)
}

func TestExtractGalaxyInfosPlanetNoActivity(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/galaxy_planet_activity.html")
	infos, _ := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, int64(0), infos.Position(15).Activity)
}

func TestExtractGalaxyInfosPlanetActivity15(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/galaxy_planet_activity.html")
	infos, _ := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, int64(15), infos.Position(8).Activity)
}

func TestExtractGalaxyInfosPlanetActivity23(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/galaxy_planet_activity.html")
	infos, _ := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, int64(23), infos.Position(9).Activity)
}

func TestExtractGalaxyInfosMoonActivity(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/galaxy_moon_activity.html")
	infos, _ := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, int64(33732827), infos.Position(3).Moon.ID)
	assert.Equal(t, int64(5830), infos.Position(3).Moon.Diameter)
	assert.Equal(t, int64(18), infos.Position(3).Moon.Activity)
}

func TestExtractGalaxyInfosMoonNoActivity(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/galaxy_moon_no_activity.html")
	infos, _ := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, int64(33650476), infos.Position(2).Moon.ID)
	assert.Equal(t, int64(7897), infos.Position(2).Moon.Diameter)
	assert.Equal(t, int64(0), infos.Position(2).Moon.Activity)
}

func TestExtractGalaxyInfosMoonActivity15(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/galaxy_moon_activity_unprecise.html")
	infos, _ := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, int64(0), infos.Position(11).Activity)
	assert.Equal(t, int64(33730993), infos.Position(11).Moon.ID)
	assert.Equal(t, int64(8944), infos.Position(11).Moon.Diameter)
	assert.Equal(t, int64(15), infos.Position(11).Moon.Activity)
}

func TestExtractUserInfosV7(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7/overview.html")
	e := NewExtractor()
	e.SetLanguage("en")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(538), infos.Points)
	assert.Equal(t, int64(1402), infos.Rank)
	assert.Equal(t, int64(3179), infos.Total)
	assert.Equal(t, "Governor Meridian", infos.PlayerName)
}

func TestExtractUserInfos(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/overview_inactive.html")
	e := NewExtractor()
	e.SetLanguage("en")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(1295), infos.Points)
}

func TestExtractUserInfos_hr(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v6/hr/overview.html")
	e := NewExtractor()
	e.SetLanguage("hr")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(214), infos.Rank)
	assert.Equal(t, int64(252), infos.Total)
}

func TestExtractUserInfos_tw(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v6/tw/overview.html")
	e := NewExtractor()
	e.SetLanguage("tw")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(212), infos.Rank)
	assert.Equal(t, int64(212), infos.Total)
}

func TestExtractUserInfos_no(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v6/no/overview.html")
	e := NewExtractor()
	e.SetLanguage("no")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(49), infos.Rank)
	assert.Equal(t, int64(50), infos.Total)
}

func TestExtractUserInfos_sk(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v6/sk/overview.html")
	e := NewExtractor()
	e.SetLanguage("sk")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(89), infos.Rank)
	assert.Equal(t, int64(90), infos.Total)
}

func TestExtractUserInfos_fi(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.2/fi/overview.html")
	e := NewExtractor()
	e.SetLanguage("fi")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(46), infos.Rank)
	assert.Equal(t, int64(51), infos.Total)
}

func TestExtractUserInfos_si(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.2/si/overview.html")
	e := NewExtractor()
	e.SetLanguage("si")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(59), infos.Rank)
	assert.Equal(t, int64(60), infos.Total)
}

func TestExtractUserInfos_hu(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.2/hu/overview.html")
	e := NewExtractor()
	e.SetLanguage("hu")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(635), infos.Rank)
	assert.Equal(t, int64(636), infos.Total)
}

func TestExtractUserInfos_gr(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v6/gr/overview.html")
	e := NewExtractor()
	e.SetLanguage("gr")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(162), infos.Rank)
	assert.Equal(t, int64(163), infos.Total)
}

func TestExtractUserInfos_ro(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.1/ro/overview.html")
	e := NewExtractor()
	e.SetLanguage("ro")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(108), infos.Rank)
	assert.Equal(t, int64(109), infos.Total)
}

func TestExtractUserInfos_mx(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v6/mx/overview.html")
	e := NewExtractor()
	e.SetLanguage("mx")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(916), infos.Rank)
	assert.Equal(t, int64(917), infos.Total)
}

func TestExtractUserInfos_de(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/de_overview.html")
	e := NewExtractor()
	e.SetLanguage("de")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(2980), infos.Rank)
	assert.Equal(t, int64(2980), infos.Total)
}

func TestExtractUserInfos_dk(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/dk_overview.html")
	e := NewExtractor()
	e.SetLanguage("dk")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(253), infos.Rank)
	assert.Equal(t, int64(254), infos.Total)
	assert.Equal(t, "Procurator Zibal", infos.PlayerName)
}

func TestExtractUserInfos_jp(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/jp_overview.html")
	e := NewExtractor()
	e.SetLanguage("jp")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(73), infos.Rank)
	assert.Equal(t, int64(73), infos.Total)
}

func TestExtractUserInfos_jp1(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v6/jp/overview.html")
	e := NewExtractor()
	e.SetLanguage("jp")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(85), infos.Rank)
	assert.Equal(t, int64(86), infos.Total)
}

func TestExtractUserInfos_cz(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v6/cz/overview.html")
	e := NewExtractor()
	e.SetLanguage("cz")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(1008), infos.Rank)
	assert.Equal(t, int64(1009), infos.Total)
}

func TestExtractUserInfos_fr(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/fr_overview.html")
	e := NewExtractor()
	e.SetLanguage("fr")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(3197), infos.Rank)
	assert.Equal(t, int64(3348), infos.Total)
	assert.Equal(t, "Bandit Pégasus", infos.PlayerName)
}

func TestExtractUserInfos_nl(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/nl_overview.html")
	e := NewExtractor()
	e.SetLanguage("nl")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(482), infos.Rank)
	assert.Equal(t, int64(542), infos.Total)
	assert.Equal(t, "Bandit Japetus", infos.PlayerName)
}

func TestExtractUserInfos_pl(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/pl_overview.html")
	e := NewExtractor()
	e.SetLanguage("pl")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(5873), infos.Rank)
	assert.Equal(t, int64(5876), infos.Total)
	assert.Equal(t, "Constable Leonis", infos.PlayerName)
}

func TestExtractUserInfos_br(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v6/br/overview.html")
	e := NewExtractor()
	e.SetLanguage("br")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(1026), infos.Rank)
	assert.Equal(t, int64(1268), infos.Total)
}

func TestExtractUserInfos_tr(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/tr_overview.html")
	e := NewExtractor()
	e.SetLanguage("tr")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(3655), infos.Rank)
	assert.Equal(t, int64(3656), infos.Total)
	assert.Equal(t, "Chief Apus", infos.PlayerName)
}

func TestExtractUserInfos_ar(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v6/ar/overview.html")
	e := NewExtractor()
	e.SetLanguage("ar")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(1158), infos.Rank)
	assert.Equal(t, int64(1159), infos.Total)
	assert.Equal(t, "Chief Lambda", infos.PlayerName)
}

func TestExtractUserInfos_it(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v6/it/overview.html")
	e := NewExtractor()
	e.SetLanguage("it")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(1776), infos.Rank)
	assert.Equal(t, int64(1777), infos.Total)
	assert.Equal(t, "President Fidis", infos.PlayerName)
}

func TestExtractUserInfos_pt(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/pt_overview.html")
	e := NewExtractor()
	e.SetLanguage("pt")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(1762), infos.Rank)
	assert.Equal(t, int64(1862), infos.Total)
	assert.Equal(t, "Director Europa", infos.PlayerName)
}

func TestExtractUserInfos_ru(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v6/ru/overview.html")
	e := NewExtractor()
	e.SetLanguage("ru")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(1067), infos.Rank)
	assert.Equal(t, int64(1068), infos.Total)
	assert.Equal(t, "Viceregent Horizon", infos.PlayerName)
}

func TestExtractUserInfos_ba(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.5.1/ba/overview.html")
	e := NewExtractor()
	e.SetLanguage("ba")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(138), infos.Rank)
	assert.Equal(t, int64(139), infos.Total)
	assert.Equal(t, "Governor Hunter", infos.PlayerName)
}

func TestExtractMoons(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/overview_with_moon.html")
	moons := NewExtractor().ExtractMoons(pageHTMLBytes)
	assert.Equal(t, 1, len(moons))
}

func TestExtractMoons2(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/overview_with_many_moon.html")
	moons := NewExtractor().ExtractMoons(pageHTMLBytes)
	assert.Equal(t, 2, len(moons))
}

func TestExtractMoon_exists(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/overview_with_moon.html")
	_, err := NewExtractor().ExtractMoon(pageHTMLBytes, ogame.MoonID(33741598))
	assert.Nil(t, err)
}

func TestExtractMoon_notExists(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/overview_with_moon.html")
	_, err := NewExtractor().ExtractMoon(pageHTMLBytes, ogame.MoonID(12345))
	assert.NotNil(t, err)
}

func TestExtractMoonByCoord_exists(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/overview_with_moon.html")
	_, err := NewExtractor().ExtractMoon(pageHTMLBytes, ogame.Coordinate{Galaxy: 4, System: 116, Position: 12, Type: ogame.MoonType})
	assert.Nil(t, err)
}

func TestExtractMoonByCoord_notExists(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/overview_with_moon.html")
	_, err := NewExtractor().ExtractMoon(pageHTMLBytes, ogame.Coordinate{Galaxy: 1, System: 2, Position: 3, Type: ogame.PlanetType})
	assert.NotNil(t, err)
}

func TestExtractIsInVacationFromDoc(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v6/es/overview_vacation.html")
	assert.True(t, NewExtractor().ExtractIsInVacation(pageHTMLBytes))
	pageHTMLBytes, _ = os.ReadFile("../../../samples/v6/es/fleet1_vacation.html")
	assert.True(t, NewExtractor().ExtractIsInVacation(pageHTMLBytes))
	pageHTMLBytes, _ = os.ReadFile("../../../samples/v6/es/shipyard.html")
	assert.False(t, NewExtractor().ExtractIsInVacation(pageHTMLBytes))
}

func TestExtractPlanetsMoon(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/overview_with_moon.html")
	planets := NewExtractor().ExtractPlanets(pageHTMLBytes)
	assert.Equal(t, ogame.MoonID(33741598), planets[0].Moon.ID)
	assert.Equal(t, "Moon", planets[0].Moon.Name)
	assert.Equal(t, "https://gf1.geo.gfsrv.net/cdn9d/8e0e6034049bd64e18a1804b42f179.gif", planets[0].Moon.Img)
	assert.Equal(t, int64(8774), planets[0].Moon.Diameter)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 116, Position: 12, Type: ogame.MoonType}, planets[0].Moon.Coordinate)
	assert.Equal(t, int64(0), planets[0].Moon.Fields.Built)
	assert.Equal(t, int64(1), planets[0].Moon.Fields.Total)
	assert.Nil(t, planets[1].Moon)
}

func TestExtractPlanets_fieldsFilled(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/overview_fields_filled.html")
	planets := NewExtractor().ExtractPlanets(pageHTMLBytes)
	assert.Equal(t, 5, len(planets))
	assert.Equal(t, ogame.PlanetID(33698658), planets[0].ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 116, Position: 12, Type: ogame.PlanetType}, planets[0].Coordinate)
	assert.Equal(t, "Homeworld", planets[0].Name)
	assert.Equal(t, "https://gf3.geo.gfsrv.net/cdnea/7d7ba402d90247ef7d89aa1035e525.png", planets[0].Img)
	assert.Equal(t, int64(-23), planets[0].Temperature.Min)
	assert.Equal(t, int64(17), planets[0].Temperature.Max)
	assert.Equal(t, int64(188), planets[0].Fields.Built)
	assert.Equal(t, int64(188), planets[0].Fields.Total)
	assert.Equal(t, int64(12800), planets[0].Diameter)
}

func TestExtractPlanetsEsV902(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v9.0.2/es/overview.html")
	planets := NewExtractor().ExtractPlanets(pageHTMLBytes)
	assert.Equal(t, 1, len(planets))
	assert.Equal(t, ogame.PlanetID(33620383), planets[0].ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 1, System: 41, Position: 4, Type: ogame.PlanetType}, planets[0].Coordinate)
	assert.Equal(t, "Planeta principal", planets[0].Name)
	assert.Equal(t, "https://gf1.geo.gfsrv.net/cdnf7/6177c65f05d9039be190b926a43f91.png", planets[0].Img)
	assert.Equal(t, int64(69), planets[0].Temperature.Min)
	assert.Equal(t, int64(109), planets[0].Temperature.Max)
	assert.Equal(t, int64(1), planets[0].Fields.Built)
	assert.Equal(t, int64(188), planets[0].Fields.Total)
	assert.Equal(t, int64(12800), planets[0].Diameter)
}

func TestExtractPlanetsTwV902(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v9.0.2/tw/overview.html")
	planets := NewExtractor().ExtractPlanets(pageHTMLBytes)
	assert.Equal(t, 1, len(planets))
	assert.Equal(t, ogame.PlanetID(33620229), planets[0].ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 1, System: 10, Position: 6, Type: ogame.PlanetType}, planets[0].Coordinate)
	assert.Equal(t, "母星", planets[0].Name)
	assert.Equal(t, "https://gf2.geo.gfsrv.net/cdn46/9f84a481c0c9a83d3b000d801d9d9d.png", planets[0].Img)
	assert.Equal(t, int64(13), planets[0].Temperature.Min)
	assert.Equal(t, int64(53), planets[0].Temperature.Max)
	assert.Equal(t, int64(0), planets[0].Fields.Built)
	assert.Equal(t, int64(193), planets[0].Fields.Total)
	assert.Equal(t, int64(12800), planets[0].Diameter)
}

func TestExtractPlanets(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/overview_inactive.html")
	planets := NewExtractor().ExtractPlanets(pageHTMLBytes)
	assert.Equal(t, 1, len(planets))
	assert.Equal(t, ogame.PlanetID(33672410), planets[0].ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 1, System: 301, Position: 5, Type: ogame.PlanetType}, planets[0].Coordinate)
	assert.Equal(t, "Homeworld", planets[0].Name)
	assert.Equal(t, "https://gf2.geo.gfsrv.net/cdn46/9f84a481c0c9a83d3b000d801d9d9d.png", planets[0].Img)
	assert.Equal(t, int64(31), planets[0].Temperature.Min)
	assert.Equal(t, int64(71), planets[0].Temperature.Max)
	assert.Equal(t, int64(89), planets[0].Fields.Built)
	assert.Equal(t, int64(188), planets[0].Fields.Total)
	assert.Equal(t, int64(12800), planets[0].Diameter)
}

func TestExtractPlanets_es(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/overview_es.html")
	planets := NewExtractor().ExtractPlanets(pageHTMLBytes)
	assert.Equal(t, 1, len(planets))
	assert.Equal(t, ogame.PlanetID(33630486), planets[0].ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 2, System: 147, Position: 8, Type: ogame.PlanetType}, planets[0].Coordinate)
	assert.Equal(t, "Planeta Principal", planets[0].Name)
	assert.Equal(t, "https://gf2.geo.gfsrv.net/cdnd1/83579badf7c16d217b06afda455cfe.png", planets[0].Img)
	assert.Equal(t, int64(18), planets[0].Temperature.Min)
	assert.Equal(t, int64(58), planets[0].Temperature.Max)
	assert.Equal(t, int64(0), planets[0].Fields.Built)
	assert.Equal(t, int64(193), planets[0].Fields.Total)
	assert.Equal(t, int64(12800), planets[0].Diameter)
}

func TestExtractPlanets_fr(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/fr_overview.html")
	planets := NewExtractor().ExtractPlanets(pageHTMLBytes)
	assert.Equal(t, 1, len(planets))
	assert.Equal(t, ogame.PlanetID(33629512), planets[0].ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 2, System: 180, Position: 4, Type: ogame.PlanetType}, planets[0].Coordinate)
	assert.Equal(t, "planète mère", planets[0].Name)
	assert.Equal(t, "https://gf1.geo.gfsrv.net/cdn35/9545f984bcd53c816a1a8452356d00.png", planets[0].Img)
	assert.Equal(t, int64(48), planets[0].Temperature.Min)
	assert.Equal(t, int64(88), planets[0].Temperature.Max)
	assert.Equal(t, int64(0), planets[0].Fields.Built)
	assert.Equal(t, int64(188), planets[0].Fields.Total)
	assert.Equal(t, int64(12800), planets[0].Diameter)
}

func TestExtractPlanets_fr1(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v6/fr/overview.html")
	planets := NewExtractor().ExtractPlanets(pageHTMLBytes)
	assert.Equal(t, 1, len(planets))
	assert.Equal(t, ogame.PlanetID(33693887), planets[0].ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 5, System: 201, Position: 4, Type: ogame.PlanetType}, planets[0].Coordinate)
	assert.Equal(t, "planète mère", planets[0].Name)
	assert.Equal(t, "https://gf1.geo.gfsrv.net/cdnf7/6177c65f05d9039be190b926a43f91.png", planets[0].Img)
	assert.Equal(t, int64(70), planets[0].Temperature.Min)
	assert.Equal(t, int64(110), planets[0].Temperature.Max)
	assert.Equal(t, int64(3), planets[0].Fields.Built)
	assert.Equal(t, int64(188), planets[0].Fields.Total)
	assert.Equal(t, int64(12800), planets[0].Diameter)
}

func TestExtractPlanets_br(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v6/br/overview.html")
	planets := NewExtractor().ExtractPlanets(pageHTMLBytes)
	assert.Equal(t, 1, len(planets))
	assert.Equal(t, ogame.PlanetID(33633767), planets[0].ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 1, System: 449, Position: 12, Type: ogame.PlanetType}, planets[0].Coordinate)
	assert.Equal(t, "Planeta Principal", planets[0].Name)
	assert.Equal(t, "https://gf3.geo.gfsrv.net/cdne8/41d05740ce1a534f5ec77feb11f100.png", planets[0].Img)
	assert.Equal(t, int64(-13), planets[0].Temperature.Min)
	assert.Equal(t, int64(27), planets[0].Temperature.Max)
	assert.Equal(t, int64(5), planets[0].Fields.Built)
	assert.Equal(t, int64(193), planets[0].Fields.Total)
	assert.Equal(t, int64(12800), planets[0].Diameter)
}

func TestExtractGalaxyInfos_honorableTarget(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/galaxy_debris.html")
	infos, _ := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.False(t, infos.Position(6).HonorableTarget)
	assert.True(t, infos.Position(8).HonorableTarget)
	assert.False(t, infos.Position(9).HonorableTarget)
}

func TestExtractGalaxyInfos_inactive(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/galaxy_debris.html")
	infos, _ := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.True(t, infos.Position(6).Inactive)
	assert.False(t, infos.Position(8).Inactive)
	assert.False(t, infos.Position(9).Inactive)
}

func TestExtractGalaxyInfos_strongPlayer(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/galaxy_debris.html")
	infos, _ := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.False(t, infos.Position(6).StrongPlayer)
	assert.True(t, infos.Position(8).StrongPlayer)
	assert.False(t, infos.Position(9).StrongPlayer)
}

func TestExtractGalaxyInfos_newbie(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/galaxy_newbie.html")
	infos, _ := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.True(t, infos.Position(4).Newbie)
}

func TestExtractGalaxyInfos_moon(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/galaxy_debris.html")
	infos, _ := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.NotNil(t, infos.Position(6).Moon)
	assert.Equal(t, int64(33701543), infos.Position(6).Moon.ID)
	assert.Equal(t, int64(8366), infos.Position(6).Moon.Diameter)
}

func TestExtractGalaxyInfos_debris(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/galaxy_debris.html")
	infos, _ := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, int64(0), infos.Position(6).Debris.Metal)
	assert.Equal(t, int64(700), infos.Position(6).Debris.Crystal)
	assert.Equal(t, int64(1), infos.Position(6).Debris.RecyclersNeeded)
}

func TestExtractGalaxyInfos_debris_es(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/galaxy_debris_es.html")
	infos, _ := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, int64(0), infos.Position(12).Debris.Metal)
	assert.Equal(t, int64(128000), infos.Position(12).Debris.Crystal)
	assert.Equal(t, int64(7), infos.Position(12).Debris.RecyclersNeeded)
}

func TestExtractGalaxyInfos_debris_fr(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/galaxy_debris_fr.html")
	infos, _ := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, int64(100), infos.Position(7).Debris.Metal)
	assert.Equal(t, int64(600), infos.Position(7).Debris.Crystal)
	assert.Equal(t, int64(1), infos.Position(7).Debris.RecyclersNeeded)
}

func TestExtractGalaxyInfos_debris_de(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/galaxy_debris_de.html")
	infos, _ := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, int64(100), infos.Position(9).Debris.Metal)
	assert.Equal(t, int64(2500), infos.Position(9).Debris.Crystal)
	assert.Equal(t, int64(1), infos.Position(9).Debris.RecyclersNeeded)
}

func TestExtractGalaxyInfos_vacation(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/galaxy_ajax.html")
	infos, _ := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.False(t, infos.Position(4).Vacation)
	assert.True(t, infos.Position(6).Vacation)
	assert.True(t, infos.Position(8).Vacation)
	assert.False(t, infos.Position(10).Vacation)
	assert.False(t, infos.Position(12).Vacation)
}

func TestExtractGalaxyInfos_alliance(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/galaxy_ajax.html")
	infos, _ := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, int64(303), infos.Position(10).Alliance.ID)
	assert.Equal(t, "Qrvix", infos.Position(10).Alliance.Name)
	assert.Equal(t, int64(27), infos.Position(10).Alliance.Rank)
	assert.Equal(t, int64(16), infos.Position(10).Alliance.Member)
}

func TestExtractGalaxyInfos_alliance_fr(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v6/fr/galaxy.html")
	infos, _ := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, int64(635), infos.Position(5).Alliance.ID)
	assert.Equal(t, "leretour", infos.Position(5).Alliance.Name)
	assert.Equal(t, int64(24), infos.Position(5).Alliance.Rank)
	assert.Equal(t, int64(11), infos.Position(5).Alliance.Member)
}

func TestExtractGalaxyInfos_alliance_es(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v6/es/galaxy.html")
	infos, _ := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, int64(500053), infos.Position(4).Alliance.ID)
	assert.Equal(t, "Los Aliens Grises", infos.Position(4).Alliance.Name)
	assert.Equal(t, int64(8), infos.Position(4).Alliance.Rank)
	assert.Equal(t, int64(70), infos.Position(4).Alliance.Member)
}

func TestUniverseSpeed(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/techtree_universe_speed.html")
	universeSpeed := ExtractUniverseSpeed(pageHTMLBytes)
	assert.Equal(t, int64(7), universeSpeed)
}

func TestCancel(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/overview_active_queue2.html")
	token, techID, listID, _ := NewExtractor().ExtractCancelBuildingInfos(pageHTMLBytes)
	assert.Equal(t, "fef7488e4809150cd16e3fa8fa14db37", token)
	assert.Equal(t, int64(4), techID)
	assert.Equal(t, int64(2099434), listID)
}

func TestCancelResearch(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/overview_active_queue2.html")
	token, techID, listID, _ := NewExtractor().ExtractCancelResearchInfos(pageHTMLBytes)
	assert.Equal(t, "fff7488e4809150cd16e3fa8fa14db37", token)
	assert.Equal(t, int64(120), techID)
	assert.Equal(t, int64(1769925), listID)
}

func TestGetConstructions(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/overview_active.html")
	buildingID, buildingCountdown, researchID, researchCountdown, _, _, _, _ := NewExtractor().ExtractConstructions(pageHTMLBytes)
	assert.Equal(t, ogame.CrystalMineID, buildingID)
	assert.Equal(t, int64(731), buildingCountdown)
	assert.Equal(t, ogame.CombustionDriveID, researchID)
	assert.Equal(t, int64(927), researchCountdown)
}

func TestExtractIPM(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/missileattacklayer.html")
	duration, max, token := NewExtractor().ExtractIPM(pageHTMLBytes)
	assert.Equal(t, "26a08f4cc0c0b513e1e8c10d49c14a27", token)
	assert.Equal(t, int64(17), max)
	assert.Equal(t, int64(15), duration)
}

func TestExtractFleetV71(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.1/en/movement.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	fleets := e.ExtractFleets(pageHTMLBytes)
	assert.Equal(t, 1, len(fleets))
	assert.Equal(t, int64(8271), fleets[0].ArriveIn)
	assert.Equal(t, int64(16545), fleets[0].BackIn)
	assert.Equal(t, ogame.Coordinate{Galaxy: 1, System: 432, Position: 6, Type: ogame.PlanetType}, fleets[0].Origin)
	assert.Equal(t, ogame.Coordinate{Galaxy: 1, System: 432, Position: 5, Type: ogame.PlanetType}, fleets[0].Destination)
	assert.Equal(t, ogame.Transport, fleets[0].Mission)
	assert.Equal(t, false, fleets[0].ReturnFlight)
	assert.Equal(t, ogame.FleetID(1674510), fleets[0].ID)
	assert.Equal(t, int64(250), fleets[0].Ships.SmallCargo)
	assert.Equal(t, int64(2), fleets[0].Ships.Pathfinder)
	assert.Equal(t, ogame.Resources{}, fleets[0].Resources)
}

func TestExtractFleetV72(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.2/de/movement.html")
	clock := clockwork.NewFakeClockAt(time.Date(2020, 3, 6, 11, 43, 15, 0, time.UTC))
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	fleets := e.ExtractFleets(pageHTMLBytes)
	assert.Equal(t, clock.Now().Add(-5031*time.Second), fleets[0].StartTime.UTC())
	assert.Equal(t, clock.Now().Add(-5041*time.Second), fleets[1].StartTime.UTC())
}

func TestExtractFleetV71_2(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.1/en/movement2.html")
	clock := clockwork.NewFakeClockAt(time.Date(2020, 1, 12, 1, 45, 34, 0, time.UTC))
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 0))
	fleets := e.ExtractFleets(pageHTMLBytes)
	assert.Equal(t, 2, len(fleets))
	assert.Equal(t, int64(621), fleets[0].ArriveIn)
	assert.Equal(t, int64(1245), fleets[0].BackIn)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 116, Position: 12, Type: ogame.MoonType}, fleets[0].Origin)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 116, Position: 12, Type: ogame.PlanetType}, fleets[0].Destination)
	assert.Equal(t, ogame.Transport, fleets[0].Mission)
	assert.Equal(t, false, fleets[0].ReturnFlight)
	assert.Equal(t, ogame.FleetID(8441918), fleets[0].ID)
	assert.Equal(t, int64(12), fleets[0].Ships.SmallCargo)
	assert.Equal(t, ogame.Resources{}, fleets[0].Resources)
	assert.Equal(t, clock.Now().Add(-3*time.Second), fleets[0].StartTime.UTC())
	assert.Equal(t, clock.Now().Add(621*time.Second), fleets[0].ArrivalTime.UTC())
	assert.Equal(t, clock.Now().Add(1245*time.Second), fleets[0].BackTime.UTC())

	assert.Equal(t, int64(-1), fleets[1].ArriveIn)
	assert.Equal(t, int64(2815), fleets[1].BackIn)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 208, Position: 10, Type: ogame.PlanetType}, fleets[1].Origin)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 116, Position: 12, Type: ogame.PlanetType}, fleets[1].Destination)
	assert.Equal(t, ogame.Transport, fleets[1].Mission)
	assert.Equal(t, true, fleets[1].ReturnFlight)
	assert.Equal(t, ogame.FleetID(8441803), fleets[1].ID)
	assert.Equal(t, int64(11), fleets[1].Ships.LargeCargo)
	assert.Equal(t, ogame.Resources{}, fleets[1].Resources)
	assert.Equal(t, clock.Now().Add(-1275*time.Second), fleets[1].StartTime.UTC())
	assert.Equal(t, clock.Now().Add(2815*time.Second), fleets[1].ArrivalTime.UTC())
	assert.Equal(t, clock.Now().Add(2815*time.Second), fleets[1].BackTime.UTC())
}

func TestExtractFleetV767(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.6.7/en/movement.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	fleets := e.ExtractFleets(pageHTMLBytes)
	assert.Equal(t, 1, len(fleets))
	assert.Equal(t, time.Date(2021, 6, 1, 9, 28, 2, 0, time.UTC), fleets[0].StartTime.UTC())
	assert.Equal(t, time.Date(2021, 6, 1, 9, 51, 10, 0, time.UTC), fleets[0].ArrivalTime.UTC())
	assert.Equal(t, time.Date(2021, 6, 1, 10, 14, 18, 0, time.UTC), fleets[0].BackTime.UTC())
}

func TestExtractFleetV7(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7/movement.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	fleets := e.ExtractFleets(pageHTMLBytes)
	assert.Equal(t, 1, len(fleets))
	assert.Equal(t, int64(1010), fleets[0].ArriveIn)
	assert.Equal(t, int64(2030), fleets[0].BackIn)
	assert.Equal(t, ogame.Coordinate{Galaxy: 9, System: 297, Position: 12, Type: ogame.PlanetType}, fleets[0].Origin)
	assert.Equal(t, ogame.Coordinate{Galaxy: 9, System: 297, Position: 9, Type: ogame.PlanetType}, fleets[0].Destination)
	assert.Equal(t, ogame.Transport, fleets[0].Mission)
	assert.Equal(t, false, fleets[0].ReturnFlight)
	assert.Equal(t, ogame.FleetID(4218727), fleets[0].ID)
	assert.Equal(t, int64(2), fleets[0].Ships.SmallCargo)
	assert.Equal(t, ogame.Resources{}, fleets[0].Resources)
}

func TestExtractFleet(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/fleets_1.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	fleets := e.ExtractFleets(pageHTMLBytes)
	assert.Equal(t, 1, len(fleets))
	assert.Equal(t, int64(4134), fleets[0].ArriveIn)
	assert.Equal(t, int64(8277), fleets[0].BackIn)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 116, Position: 12, Type: ogame.PlanetType}, fleets[0].Origin)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 117, Position: 9, Type: ogame.PlanetType}, fleets[0].Destination)
	assert.Equal(t, ogame.Transport, fleets[0].Mission)
	assert.Equal(t, false, fleets[0].ReturnFlight)
	assert.Equal(t, ogame.FleetID(4494950), fleets[0].ID)
	assert.Equal(t, int64(1), fleets[0].Ships.SmallCargo)
	assert.Equal(t, int64(8), fleets[0].Ships.LargeCargo)
	assert.Equal(t, int64(1), fleets[0].Ships.LightFighter)
	assert.Equal(t, int64(1), fleets[0].Ships.ColonyShip)
	assert.Equal(t, int64(1), fleets[0].Ships.EspionageProbe)
	assert.Equal(t, ogame.Resources{Metal: 123, Crystal: 456, Deuterium: 789}, fleets[0].Resources)
}

func TestExtractFleet_expedition(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/fleets_expedition.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	fleets := e.ExtractFleets(pageHTMLBytes)
	assert.Equal(t, 2, len(fleets))
	assert.Equal(t, int64(2), fleets[1].Ships.LargeCargo)
	assert.Equal(t, ogame.Expedition, fleets[1].Mission)
	assert.False(t, fleets[1].ReturnFlight)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 116, Position: 12, Type: ogame.MoonType}, fleets[1].Origin)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 116, Position: 16, Type: ogame.PlanetType}, fleets[1].Destination)
}

func TestExtractFleet_harvest(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/fleets_harvest.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	fleets := e.ExtractFleets(pageHTMLBytes)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 116, Position: 12, Type: ogame.PlanetType}, fleets[5].Origin)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 116, Position: 9, Type: ogame.DebrisType}, fleets[5].Destination)
}

func TestExtractFleet_returningTransport(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/fleets_2.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	fleets := e.ExtractFleets(pageHTMLBytes)
	assert.Equal(t, 1, len(fleets))
	assert.Equal(t, int64(-1), fleets[0].ArriveIn)
	assert.Equal(t, int64(36), fleets[0].BackIn)
}

func TestExtractFleet_deployment(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/fleets_moon_to_moon.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	fleets := e.ExtractFleets(pageHTMLBytes)
	assert.Equal(t, int64(210), fleets[0].ArriveIn)
	assert.Equal(t, int64(426), fleets[0].BackIn)
}

func TestExtractFleetThousands(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/fleets_thousands.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	fleets := e.ExtractFleets(pageHTMLBytes)
	assert.Equal(t, ogame.Transport, fleets[0].Mission)
	assert.Equal(t, int64(210), fleets[0].Ships.LargeCargo)
	assert.Equal(t, ogame.Resources{Metal: 207862, Crystal: 78903, Deuterium: 42956}, fleets[0].Resources)
}

func TestExtractFleet_returning(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/fleets_2.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	fleets := e.ExtractFleets(pageHTMLBytes)
	assert.Equal(t, 1, len(fleets))
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 116, Position: 12, Type: ogame.PlanetType}, fleets[0].Origin)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 117, Position: 9, Type: ogame.PlanetType}, fleets[0].Destination)
	assert.Equal(t, ogame.Transport, fleets[0].Mission)
	assert.Equal(t, true, fleets[0].ReturnFlight)
	assert.Equal(t, ogame.FleetID(4494950), fleets[0].ID)
	assert.Equal(t, int64(1), fleets[0].Ships.SmallCargo)
	assert.Equal(t, int64(8), fleets[0].Ships.LargeCargo)
	assert.Equal(t, int64(1), fleets[0].Ships.LightFighter)
	assert.Equal(t, int64(1), fleets[0].Ships.ColonyShip)
	assert.Equal(t, int64(1), fleets[0].Ships.EspionageProbe)
	assert.Equal(t, ogame.Resources{Metal: 123, Crystal: 456, Deuterium: 789}, fleets[0].Resources)
}

func TestExtractFleet_deepspace(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.2/en/fleets_expeditions.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	fleets := e.ExtractFleets(pageHTMLBytes)
	assert.Equal(t, 5, len(fleets))
	assert.False(t, fleets[0].InDeepSpace)
	assert.False(t, fleets[1].InDeepSpace)
	assert.False(t, fleets[2].InDeepSpace)
	assert.False(t, fleets[3].InDeepSpace)
	assert.True(t, fleets[4].InDeepSpace)
}

func TestExtractFleet_targetPlanetID(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/fleets_moon_to_moon.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	fleets := e.ExtractFleets(pageHTMLBytes)
	assert.Equal(t, int64(0), fleets[0].TargetPlanetID)
	assert.Equal(t, int64(0), fleets[1].TargetPlanetID)
	assert.Equal(t, int64(33702114), fleets[2].TargetPlanetID)
	assert.Equal(t, int64(33699325), fleets[3].TargetPlanetID)
}

func TestExtractFleet_unionID(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/fleets_no_union.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	fleets := e.ExtractFleets(pageHTMLBytes)
	assert.Equal(t, int64(0), fleets[0].UnionID)

	pageHTMLBytes, _ = os.ReadFile("../../../samples/unversioned/fleets_union_alone.html")
	e = NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	fleets = e.ExtractFleets(pageHTMLBytes)
	assert.Equal(t, int64(13558), fleets[0].UnionID)
}

func TestExtractOverviewProduction(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/overview_shipyard_queue_full.html")
	prods, countdown, _ := NewExtractor().ExtractOverviewProduction(pageHTMLBytes)
	assert.Equal(t, 6, len(prods))
	assert.Equal(t, int64(3399), countdown)
	assert.Equal(t, ogame.HeavyFighterID, prods[0].ID)
	assert.Equal(t, int64(1), prods[0].Nbr)
	assert.Equal(t, ogame.HeavyFighterID, prods[1].ID)
	assert.Equal(t, int64(2), prods[1].Nbr)
	assert.Equal(t, ogame.HeavyFighterID, prods[2].ID)
	assert.Equal(t, int64(3), prods[2].Nbr)
	assert.Equal(t, ogame.HeavyFighterID, prods[3].ID)
	assert.Equal(t, int64(4), prods[3].Nbr)
	assert.Equal(t, ogame.HeavyFighterID, prods[4].ID)
	assert.Equal(t, int64(5), prods[4].Nbr)
	assert.Equal(t, ogame.HeavyFighterID, prods[5].ID)
	assert.Equal(t, int64(6), prods[5].Nbr)
}

func TestExtractProduction(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/shipyard_queue.html")
	prods, secs, _ := NewExtractor().ExtractProduction(pageHTMLBytes)
	assert.Equal(t, 20, len(prods))
	assert.Equal(t, int64(16254), secs)
	assert.Equal(t, ogame.LargeCargoID, prods[0].ID)
	assert.Equal(t, int64(4), prods[0].Nbr)
}

func TestExtractProduction2(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/shipyard_queue2.html")
	prods, secs, _ := NewExtractor().ExtractProduction(pageHTMLBytes)
	assert.Equal(t, int64(7082), secs)
	assert.Equal(t, ogame.BattlecruiserID, prods[0].ID)
	assert.Equal(t, int64(18), prods[0].Nbr)
	assert.Equal(t, ogame.PlasmaTurretID, prods[1].ID)
	assert.Equal(t, int64(8), prods[1].Nbr)
	assert.Equal(t, ogame.RocketLauncherID, prods[2].ID)
	assert.Equal(t, int64(1000), prods[2].Nbr)
	assert.Equal(t, ogame.LightFighterID, prods[10].ID)
	assert.Equal(t, int64(1), prods[10].Nbr)
}

func TestExtractProductionWithABM(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/production_with_abm.html")
	prods, secs, _ := NewExtractor().ExtractProduction(pageHTMLBytes)
	assert.Equal(t, 4, len(prods))
	assert.Equal(t, int64(220), secs)
	assert.Equal(t, ogame.DeathstarID, prods[0].ID)
	assert.Equal(t, int64(1), prods[0].Nbr)
	assert.Equal(t, ogame.AntiBallisticMissilesID, prods[1].ID)
	assert.Equal(t, int64(1), prods[1].Nbr)
	assert.Equal(t, ogame.InterplanetaryMissilesID, prods[2].ID)
	assert.Equal(t, int64(1), prods[2].Nbr)
}

func TestExtractDKProductionWithABM(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v6/dk/production_with_abm.html")
	prods, secs, _ := NewExtractor().ExtractProduction(pageHTMLBytes)
	assert.Equal(t, 2, len(prods))
	assert.Equal(t, int64(641), secs)
	assert.Equal(t, ogame.AntiBallisticMissilesID, prods[0].ID)
	assert.Equal(t, int64(1), prods[0].Nbr)
	assert.Equal(t, ogame.AntiBallisticMissilesID, prods[1].ID)
	assert.Equal(t, int64(1), prods[1].Nbr)
}

func TestExtractEspionageReport_tz(t *testing.T) {
	clock := clockwork.NewFakeClockAt(time.Date(2019, 10, 27, 0, 26, 4, 0, time.UTC))
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/spy_report_17h26-7Z.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, clock.Now(), infos.Date.UTC())
}

func TestExtractEspionageReport_action(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/message_foreign_fleet_sighted.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, ogame.Action, infos.Type)
	assert.Equal(t, int64(6970988), infos.ID)
}

func TestExtractEspionageReport(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/spy_report_res_buildings_researches.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 212, Position: 6, Type: ogame.PlanetType}, infos.Coordinate)
	assert.Equal(t, ogame.Report, infos.Type)
	assert.True(t, infos.HasFleetInformation)
	assert.True(t, infos.HasDefensesInformation)
	assert.True(t, infos.HasBuildingsInformation)
	assert.True(t, infos.HasResearchesInformation)
	assert.Equal(t, int64(6862893), infos.ID)
	assert.Equal(t, int64(0), infos.CounterEspionage)
	assert.Equal(t, int64(227034), infos.Metal)
	assert.Equal(t, int64(146970), infos.Crystal)
	assert.Equal(t, int64(24751), infos.Deuterium)
	assert.Equal(t, int64(2324), infos.Energy)
	assert.Equal(t, int64(20), *infos.MetalMine)
	assert.Equal(t, int64(14), *infos.CrystalMine)
	assert.Equal(t, int64(8), *infos.DeuteriumSynthesizer)
	assert.Equal(t, int64(19), *infos.SolarPlant)
	assert.Equal(t, int64(5), *infos.RoboticsFactory)
	assert.Equal(t, int64(2), *infos.Shipyard)
	assert.Equal(t, int64(5), *infos.MetalStorage)
	assert.Equal(t, int64(5), *infos.CrystalStorage)
	assert.Equal(t, int64(2), *infos.DeuteriumTank)
	assert.Equal(t, int64(3), *infos.ResearchLab)
	assert.Equal(t, int64(2), *infos.EspionageTechnology)
	assert.Equal(t, int64(1), *infos.ComputerTechnology)
	assert.Equal(t, int64(1), *infos.ArmourTechnology)
	assert.Equal(t, int64(1), *infos.EnergyTechnology)
	assert.Equal(t, int64(7), *infos.CombustionDrive)
	assert.Equal(t, int64(2), *infos.ImpulseDrive)
	assert.Nil(t, infos.LightFighter)
	assert.Nil(t, infos.HeavyFighter)
	assert.Nil(t, infos.Cruiser)
	assert.Nil(t, infos.Battleship)
	assert.Nil(t, infos.Battlecruiser)
	assert.Nil(t, infos.Bomber)
	assert.Nil(t, infos.Destroyer)
	assert.Nil(t, infos.Deathstar)
	assert.Nil(t, infos.SmallCargo)
	assert.Nil(t, infos.LargeCargo)
	assert.Nil(t, infos.ColonyShip)
	assert.Nil(t, infos.Recycler)
	assert.Nil(t, infos.EspionageProbe)
	assert.Nil(t, infos.SolarSatellite)
}

func TestExtractEspionageReport_noPictures(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/spy_report_no_pics.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, err := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, ogame.ErrDeactivateHidePictures, err)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 203, Position: 6, Type: ogame.PlanetType}, infos.Coordinate)
	assert.Equal(t, ogame.Report, infos.Type)
	assert.True(t, infos.HasFleetInformation)
	assert.True(t, infos.HasDefensesInformation)
	assert.True(t, infos.HasBuildingsInformation)
	assert.True(t, infos.HasResearchesInformation)
	assert.Equal(t, int64(9142399), infos.ID)
	assert.Equal(t, int64(0), infos.CounterEspionage)
	assert.Equal(t, int64(1131895), infos.Metal)
	assert.Equal(t, int64(432515), infos.Crystal)
	assert.Equal(t, int64(114957), infos.Deuterium)
	assert.Equal(t, int64(4727), infos.Energy)
	assert.Nil(t, infos.MetalMine)
	assert.Nil(t, infos.CrystalMine)
	assert.Nil(t, infos.DeuteriumSynthesizer)
	assert.Nil(t, infos.SolarPlant)
	assert.Nil(t, infos.RoboticsFactory)
	assert.Nil(t, infos.Shipyard)
	assert.Nil(t, infos.MetalStorage)
	assert.Nil(t, infos.CrystalStorage)
}

func TestExtractEspionageReportMoon(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/spy_report_moon.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 116, Position: 12, Type: ogame.MoonType}, infos.Coordinate)
	assert.Equal(t, int64(6), *infos.LunarBase)
	assert.Equal(t, int64(4), *infos.SensorPhalanx)
	assert.Nil(t, infos.JumpGate)
}

func TestExtractEspionageReport1(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/spy_report_res_buildings_researches_fleet.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, int64(2), *infos.Battleship)
	assert.Equal(t, int64(1), *infos.Bomber)
}

func TestExtractEspionageReportThousands(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/spy_report_thousand_units.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, int64(4000), *infos.RocketLauncher)
	assert.Equal(t, int64(3882), *infos.LargeCargo)
	assert.Equal(t, int64(374), *infos.SolarSatellite)
}

func TestExtractEspionageReport_defence(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/spy_report_res_fleet_defences.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.True(t, infos.HasFleetInformation)
	assert.True(t, infos.HasDefensesInformation)
	assert.False(t, infos.HasBuildingsInformation)
	assert.False(t, infos.HasResearchesInformation)
	assert.Equal(t, int64(13), infos.CounterEspionage)
	assert.Equal(t, int64(57), *infos.RocketLauncher)
	assert.Equal(t, int64(57), *infos.LightLaser)
	assert.Equal(t, int64(61), *infos.HeavyLaser)
	assert.Nil(t, infos.GaussCannon)
}

func TestExtractEspionageReport_bandit(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/spy_report_inactive_bandit_lord.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, true, infos.IsBandit)
	assert.Equal(t, false, infos.IsStarlord)
}

func TestExtractEspionageReport_starlord(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/spy_report_active_star_lord.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, false, infos.IsBandit)
	assert.Equal(t, true, infos.IsStarlord)
}

func TestExtractEspionageReport_norank(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/spy_report_res_buildings_researches_fleet.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, false, infos.IsBandit)
	assert.Equal(t, false, infos.IsStarlord)

}

func TestExtractEspionageReport_username1(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/spy_report_inactive_bandit_lord.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, "Cid Granjeador", infos.Username)
}

func TestExtractEspionageReport_username2(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/spy_report_active_star_lord.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, "Commodore Nomad", infos.Username)
}

func TestExtractEspionageReport_username_outlaw(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/spy_report_outlaw.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, "Constable Telesto", infos.Username)
}

func TestExtractEspionageReport_apiKey(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/spy_report_active_star_lord.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, "sr-en-152-ea0b59302bfad7e3ab0f2d15f7ef2c6a4633b4ba", infos.APIKey)
}

func TestExtractEspionageReport_inactivetimer_within15(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/spy_report_res_buildings.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, int64(15), infos.LastActivity)
}

func TestExtractEspionageReport_inactivetimer_29(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/spy_report_res_buildings_researches.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, int64(29), infos.LastActivity)
}

func TestExtractEspionageReport_inactivetimer_over1h(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/spy_report_inactive_bandit_lord.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, int64(0), infos.LastActivity)
}

func TestExtractFleetSlotV7_movement(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7/movement.html")
	s, _ := NewExtractor().ExtractSlots(pageHTMLBytes)
	assert.Equal(t, int64(1), s.InUse)
	assert.Equal(t, int64(2), s.Total)
	assert.Equal(t, int64(0), s.ExpInUse)
	assert.Equal(t, int64(1), s.ExpTotal)
}

func TestExtractFleetSlot_fleet1(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/fleet1.html")
	s, _ := NewExtractor().ExtractSlots(pageHTMLBytes)
	assert.Equal(t, int64(2), s.InUse)
	assert.Equal(t, int64(14), s.Total)
	assert.Equal(t, int64(0), s.ExpInUse)
	assert.Equal(t, int64(3), s.ExpTotal)
}

func TestExtractFleetSlot_movement(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/fleets_1.html")
	s, _ := NewExtractor().ExtractSlots(pageHTMLBytes)
	assert.Equal(t, int64(1), s.InUse)
	assert.Equal(t, int64(11), s.Total)
	assert.Equal(t, int64(0), s.ExpInUse)
	assert.Equal(t, int64(2), s.ExpTotal)
}

func TestExtractFleetSlot_commanders(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/fleet1_extract_slots_with_commanders.html")
	s, _ := NewExtractor().ExtractSlots(pageHTMLBytes)
	assert.Equal(t, int64(13), s.InUse)
	assert.Equal(t, int64(14), s.Total)
	assert.Equal(t, int64(2), s.ExpInUse)
	assert.Equal(t, int64(3), s.ExpTotal)
}

func TestGetResourcesDetails(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/fetch_resources.html")
	res, _ := NewExtractor().ExtractResourcesDetails(pageHTMLBytes)
	assert.Equal(t, int64(380030343), res.Metal.Available)
	assert.Equal(t, int64(60510000), res.Metal.StorageCapacity)
	assert.Equal(t, int64(0), res.Metal.CurrentProduction)

	assert.Equal(t, int64(19320), res.Crystal.Available)
	assert.Equal(t, int64(9820000), res.Crystal.StorageCapacity)
	assert.Equal(t, int64(40636), res.Crystal.CurrentProduction)

	assert.Equal(t, int64(24902), res.Deuterium.Available)
	assert.Equal(t, int64(18005000), res.Deuterium.StorageCapacity)
	assert.Equal(t, int64(22508), res.Deuterium.CurrentProduction)

	assert.Equal(t, int64(-8402), res.Energy.Available)
	assert.Equal(t, int64(10469), res.Energy.CurrentProduction)
	assert.Equal(t, int64(-18871), res.Energy.Consumption)

	assert.Equal(t, int64(28500), res.Darkmatter.Available)
	assert.Equal(t, int64(0), res.Darkmatter.Purchased)
	assert.Equal(t, int64(28500), res.Darkmatter.Found)
}

func TestExtractEmpirePlanets(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v8.1/en/empire_planets.html")
	res, _ := NewExtractor().ExtractEmpire(pageHTMLBytes)
	assert.Equal(t, 8, len(res))
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 208, Position: 8, Type: ogame.PlanetType}, res[0].Coordinate)
	assert.Equal(t, int64(-3199), res[0].Resources.Energy)
	assert.Equal(t, int64(13904), res[0].Diameter)
}

func TestExtractEmpireMoons(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v8.1/en/empire_moons.html")
	res, _ := NewExtractor().ExtractEmpire(pageHTMLBytes)
	assert.Equal(t, 3, len(res))
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 116, Position: 9, Type: ogame.MoonType}, res[0].Coordinate)
	assert.Equal(t, int64(0), res[0].Resources.Energy)
	assert.Equal(t, int64(-19), res[0].Temperature.Min)
	assert.Equal(t, int64(21), res[0].Temperature.Max)
	assert.Equal(t, int64(5783), res[0].Diameter)
}

func TestExtractAuction_playerBid(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.5.0/en/auction_player_bid.html")
	res, _ := NewExtractor().ExtractAuction(pageHTMLBytes)
	assert.Equal(t, int64(1603000), res.AlreadyBid)
}

func TestExtractAuction_noPlayerBid(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.5.0/en/auction_no_player_bid.html")
	res, _ := NewExtractor().ExtractAuction(pageHTMLBytes)
	assert.Equal(t, int64(0), res.AlreadyBid)
}

func TestExtractAuction_ongoing2(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.4/en/traderAuctioneer_ongoing.html")
	res, _ := NewExtractor().ExtractAuction(pageHTMLBytes)
	assert.Equal(t, int64(1800), res.Endtime)
}

func TestExtractAuction_ongoing(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.1/en/traderOverview_ongoing.html")
	res, _ := NewExtractor().ExtractAuction(pageHTMLBytes)
	assert.Equal(t, int64(1200), res.Endtime)
}

func TestExtractAuction_waiting(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.1/en/traderOverview_waiting.html")
	res, _ := NewExtractor().ExtractAuction(pageHTMLBytes)
	assert.Equal(t, int64(6202), res.Endtime)
}

func TestExtractOGameSession(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7/overview.html")
	session := NewExtractor().ExtractOGameSession(pageHTMLBytes)
	assert.Equal(t, "0a724276a3ddbe9949f62bdae48d71c1a16adf20", session)

	pageHTMLBytes, _ = os.ReadFile("../../../samples/v7/overview_mobile.html")
	session = NewExtractor().ExtractOGameSession(pageHTMLBytes)
	assert.Equal(t, "c1626ce8228ac5986e3808a7d42d4afc764c1b68", session)
}
