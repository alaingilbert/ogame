package wrapper

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/utils"
	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/assert"
	"os"
	"regexp"
	"testing"
)

func BenchmarkUserInfoRegex(b *testing.B) {
	extractUserRegex := func(pageHTML []byte) (int, string) {
		playerID := utils.ToInt(regexp.MustCompile(`playerId="(\d+)"`).FindSubmatch(pageHTML)[1])
		playerName := string(regexp.MustCompile(`playerName="([^"]+)"`).FindSubmatch(pageHTML)[1])
		return playerID, playerName
	}
	pageHTMLBytes, _ := os.ReadFile("../../samples/overview_inactive.html")
	for n := 0; n < b.N; n++ {
		extractUserRegex(pageHTMLBytes)
	}
}

func BenchmarkUserInfoGoquery(b *testing.B) {
	extractUserGoquery := func(pageHTML []byte) (int64, string) {
		doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
		playerID := utils.ParseInt(doc.Find("meta[name=ogame-player-id]").AttrOr("content", "0"))
		playerName := doc.Find("meta[name=ogame-player-name]").AttrOr("content", "")
		return playerID, playerName
	}
	pageHTMLBytes, _ := os.ReadFile("../../samples/overview_inactive.html")
	for n := 0; n < b.N; n++ {
		extractUserGoquery(pageHTMLBytes)
	}
}

func TestWrapper(t *testing.T) {
	var bot Wrapper
	bot, _ = NewNoLogin("", "", "", "", "", "", 0, nil)
	assert.NotNil(t, bot)
}

//func TestGetResourcesProductionsLight(t *testing.T) {
//	supplies := ResourcesBuildings{
//		MetalMine:            32,
//		CrystalMine:          28,
//		DeuteriumSynthesizer: 28,
//		SolarPlant:           30,
//		FusionReactor:        9,
//		SolarSatellite:       0,
//	}
//	researches := Researches{
//		EnergyTechnology: 18,
//		PlasmaTechnology: 15,
//	}
//	resSettings := ResourceSettings{
//		MetalMine:            100,
//		CrystalMine:          100,
//		DeuteriumSynthesizer: 60,
//		SolarPlant:           100,
//		FusionReactor:        0,
//		SolarSatellite:       100,
//		Crawler:              0,
//	}
//	temp := Temperature{Min: -23, Max: 17}
//	prod := getResourcesProductionsLight(supplies, researches, resSettings, temp, 7)
//	assert.Equal(t, Resources{Metal: 109444, Crystal: 41697, Deuterium: 16347, Energy: -5169}, prod)
//}

func TestProductionRatio(t *testing.T) {
	ratio := productionRatio(
		ogame.Temperature{-23, 17},
		ogame.ResourcesBuildings{MetalMine: 29, CrystalMine: 16, DeuteriumSynthesizer: 26, SolarPlant: 29, FusionReactor: 13, SolarSatellite: 51},
		ogame.ResourceSettings{MetalMine: 100, CrystalMine: 100, DeuteriumSynthesizer: 100, SolarPlant: 100, FusionReactor: 100, SolarSatellite: 100},
		12,
	)
	assert.Equal(t, 1.0, ratio)
}

func TestEnergyNeeded(t *testing.T) {
	needed := energyNeeded(
		ogame.ResourcesBuildings{MetalMine: 29, CrystalMine: 16, DeuteriumSynthesizer: 26},
		ogame.ResourceSettings{MetalMine: 100, CrystalMine: 100, DeuteriumSynthesizer: 100},
	)
	assert.Equal(t, int64(4601+736+6198), needed)
}

func TestEnergyProduced(t *testing.T) {
	produced := energyProduced(
		ogame.Temperature{-23, 17},
		ogame.ResourcesBuildings{SolarPlant: 29, FusionReactor: 13, SolarSatellite: 51},
		ogame.ResourceSettings{SolarPlant: 100, FusionReactor: 100, SolarSatellite: 100},
		12,
	)
	assert.Equal(t, int64(9200+3002+1326), produced)
}

func TestExtractCargoCapacity(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../samples/unversioned/sendfleet3.htm")
	fleet3Doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTMLBytes))
	cargo := utils.ParseInt(fleet3Doc.Find("#maxresources").Text())
	assert.Equal(t, int64(442500), cargo)
}

//func TestExtractGalaxyInfosPlanetActivityWithoutDetailedActivity(t *testing.T) {
//	pageHTMLBytes, _ := os.ReadFile("../../samples/galaxy_planet_activity_without_detailed_activity.html")
//	infos, _ := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
//	assert.Equal(t, 49, infos.Position(5).Activity)
//}

func TestExtractFleetsFromEventList(t *testing.T) {
	//pageHTMLBytes, _ := os.ReadFile("../../samples/eventlist_test.html")
	//fleets := NewExtractor().ExtractFleetsFromEventList(pageHTMLBytes)
	//assert.Equal(t, 4, len(fleets))
}

func TestGalaxyDistance(t *testing.T) {
	assert.Equal(t, int64(60000), galaxyDistance(6, 3, 6, false))
	assert.Equal(t, int64(20000), galaxyDistance(1, 2, 6, false))
	assert.Equal(t, int64(40000), galaxyDistance(1, 3, 6, false))
	assert.Equal(t, int64(60000), galaxyDistance(1, 4, 6, false))
	assert.Equal(t, int64(80000), galaxyDistance(1, 5, 6, false))
	assert.Equal(t, int64(100000), galaxyDistance(1, 6, 6, false))

	assert.Equal(t, int64(20000), galaxyDistance(1, 2, 6, true))
	assert.Equal(t, int64(40000), galaxyDistance(1, 3, 6, true))
	assert.Equal(t, int64(60000), galaxyDistance(1, 4, 6, true))
	assert.Equal(t, int64(40000), galaxyDistance(1, 5, 6, true))
	assert.Equal(t, int64(20000), galaxyDistance(1, 6, 6, true))
	assert.Equal(t, int64(20000), galaxyDistance(6, 1, 6, true))
}

func TestSystemDistance(t *testing.T) {
	assert.Equal(t, int64(3175), flightSystemDistance(499, 35, 30, false))

	assert.Equal(t, int64(2795), flightSystemDistance(499, 1, 2, true))
	assert.Equal(t, int64(2795), flightSystemDistance(499, 1, 499, true))
	assert.Equal(t, int64(2890), flightSystemDistance(499, 1, 3, true))
	assert.Equal(t, int64(2890), flightSystemDistance(499, 1, 498, true))
	assert.Equal(t, int64(2890), flightSystemDistance(499, 498, 1, true))
}

func TestPlanetDistance(t *testing.T) {
	assert.Equal(t, int64(1015), planetDistance(6, 3))
}

func TestDistance(t *testing.T) {
	assert.Equal(t, int64(1015), Distance(ogame.Coordinate{1, 1, 3, ogame.PlanetType}, ogame.Coordinate{1, 1, 6, ogame.PlanetType}, 6, 499, true, true))
	assert.Equal(t, int64(2890), Distance(ogame.Coordinate{1, 1, 3, ogame.PlanetType}, ogame.Coordinate{1, 498, 6, ogame.PlanetType}, 6, 499, true, true))
	assert.Equal(t, int64(20000), Distance(ogame.Coordinate{6, 1, 3, ogame.PlanetType}, ogame.Coordinate{1, 498, 6, ogame.PlanetType}, 6, 499, true, true))
	assert.Equal(t, int64(5), Distance(ogame.Coordinate{6, 1, 3, ogame.PlanetType}, ogame.Coordinate{6, 1, 3, ogame.MoonType}, 6, 499, true, true))
}

func TestCalcFlightTime(t *testing.T) {
	// Test from https://ogame.fandom.com/wiki/Talk:Fuel_Consumption
	secs, fuel := CalcFlightTime(ogame.Coordinate{1, 1, 1, ogame.PlanetType}, ogame.Coordinate{1, 5, 3, ogame.PlanetType},
		1, 499, false, false, 1, 0.8, 1, ogame.ShipsInfos{LightFighter: 16, HeavyFighter: 8, Cruiser: 4}, ogame.Researches{CombustionDrive: 10, ImpulseDrive: 7}, ogame.NoClass)
	assert.Equal(t, int64(4966), secs)
	assert.Equal(t, int64(550), fuel)

	// Different fleetDeutSaveFactor
	secs, fuel = CalcFlightTime(ogame.Coordinate{4, 116, 12, ogame.PlanetType}, ogame.Coordinate{3, 116, 12, ogame.PlanetType},
		6, 499, true, true, 0.5, 1, 2, ogame.ShipsInfos{LargeCargo: 1931}, ogame.Researches{CombustionDrive: 18, ImpulseDrive: 15, HyperspaceDrive: 13}, ogame.Discoverer)
	assert.Equal(t, int64(5406), secs)
	assert.Equal(t, int64(110336), fuel)

	// Test with solar satellite
	secs, fuel = CalcFlightTime(ogame.Coordinate{1, 1, 1, ogame.PlanetType}, ogame.Coordinate{1, 1, 15, ogame.PlanetType},
		6, 499, false, false, 1, 1, 4, ogame.ShipsInfos{LargeCargo: 100, SolarSatellite: 50}, ogame.Researches{CombustionDrive: 16, ImpulseDrive: 13, HyperspaceDrive: 15}, ogame.NoClass)
	assert.Equal(t, int64(651), secs)
	assert.Equal(t, int64(612), fuel)

	// General tests
	secs, fuel = CalcFlightTime(
		ogame.Coordinate{2, 68, 4, ogame.MoonType},
		ogame.Coordinate{1, 313, 9, ogame.PlanetType},
		5, 499, true, true, 1, 1, 2,
		ogame.ShipsInfos{LightFighter: 1, HeavyFighter: 1, Cruiser: 1, Battleship: 1, SmallCargo: 1, LargeCargo: 1, Recycler: 1, ColonyShip: 1, EspionageProbe: 1},
		ogame.Researches{CombustionDrive: 7, ImpulseDrive: 5, HyperspaceDrive: 0}, ogame.Discoverer)
	assert.Equal(t, int64(13427), secs)
	assert.Equal(t, int64(3808), fuel)

	secs, fuel = CalcFlightTime(
		ogame.Coordinate{1, 230, 7, ogame.MoonType},
		ogame.Coordinate{1, 318, 4, ogame.MoonType},
		5, 499, true, true, 0.5, 1, 6,
		ogame.ShipsInfos{LightFighter: 1, HeavyFighter: 1, Cruiser: 1, Battleship: 1, SmallCargo: 1, LargeCargo: 1, Recycler: 1, EspionageProbe: 1, Pathfinder: 1},
		ogame.Researches{CombustionDrive: 10, ImpulseDrive: 6, HyperspaceDrive: 4}, ogame.Discoverer)
	assert.Equal(t, int64(3069), secs)
	assert.Equal(t, int64(584), fuel)

	secs, fuel = CalcFlightTime(
		ogame.Coordinate{1, 230, 7, ogame.MoonType},
		ogame.Coordinate{1, 318, 4, ogame.MoonType},
		5, 499, true, true, 0.5, 1, 6,
		ogame.ShipsInfos{EspionageProbe: 9000},
		ogame.Researches{CombustionDrive: 10, ImpulseDrive: 6, HyperspaceDrive: 4}, ogame.Discoverer)
	assert.Equal(t, int64(15), secs)
	assert.Equal(t, int64(1), fuel)

	secs, fuel = CalcFlightTime(
		ogame.Coordinate{1, 230, 7, ogame.MoonType},
		ogame.Coordinate{1, 318, 4, ogame.MoonType},
		5, 499, true, true, 1, 1, 6,
		ogame.ShipsInfos{EspionageProbe: 9000},
		ogame.Researches{CombustionDrive: 10, ImpulseDrive: 6, HyperspaceDrive: 4}, ogame.General)
	assert.Equal(t, int64(15), secs)
	assert.Equal(t, int64(1), fuel)
}

func TestFixAttackEvents(t *testing.T) {
	// Test when moon name matches
	p1 := Planet{}
	p1.Name = "My Planet"
	p1.Coordinate = ogame.Coordinate{1, 2, 3, ogame.PlanetType}
	p1.Moon = &Moon{Moon: ogame.Moon{Name: "VeryLongName Moon"}}
	planets := []Planet{p1}

	attacks := []ogame.AttackEvent{
		{DestinationName: "VeryLongName Moon", Destination: ogame.Coordinate{1, 2, 3, ogame.PlanetType}},
	}
	fixAttackEvents(attacks, planets)
	assert.Equal(t, ogame.MoonType, attacks[0].Destination.Type) // Fixed to moon type

	// Test when the moon name doesn't match
	attacks = []ogame.AttackEvent{
		{DestinationName: "My Planet", Destination: ogame.Coordinate{1, 2, 3, ogame.PlanetType}},
	}
	fixAttackEvents(attacks, planets)
	assert.Equal(t, ogame.PlanetType, attacks[0].Destination.Type) // Did not change
}

func TestVersion(t *testing.T) {
	assert.False(t, version.Must(version.NewVersion("8.7.4-pl3")).GreaterThanOrEqual(version.Must(version.NewVersion("8.7.4"))))
	assert.True(t, version.Must(version.NewVersion("8.7.4-pl3")).GreaterThanOrEqual(version.Must(version.NewVersion("8.7.4-pl3"))))
	assert.True(t, version.Must(version.NewVersion("8.7.4-pl4")).GreaterThanOrEqual(version.Must(version.NewVersion("8.7.4-pl3"))))
	assert.True(t, version.Must(version.NewVersion("8.7.5-pl3")).GreaterThanOrEqual(version.Must(version.NewVersion("8.7.5-pl3"))))
}

func TestFindSlowestSpeed(t *testing.T) {
	assert.Equal(t, int64(8000), findSlowestSpeed(ogame.ShipsInfos{SmallCargo: 1, LargeCargo: 1}, ogame.Researches{CombustionDrive: 6}, false, false))
}
