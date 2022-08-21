package wrapper

import (
	"bytes"
	"github.com/alaingilbert/ogame/pkg/extractor/v6"
	"github.com/alaingilbert/ogame/pkg/extractor/v7"
	"github.com/alaingilbert/ogame/pkg/extractor/v71"
	"github.com/alaingilbert/ogame/pkg/extractor/v8"
	"github.com/alaingilbert/ogame/pkg/extractor/v874"
	"github.com/alaingilbert/ogame/pkg/extractor/v9"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/utils"
	"github.com/hashicorp/go-version"
	"io/ioutil"
	"regexp"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/alaingilbert/clockwork"
	"github.com/stretchr/testify/assert"
)

func BenchmarkUserInfoRegex(b *testing.B) {
	extractUserRegex := func(pageHTML []byte) (int, string) {
		playerID := utils.ToInt(regexp.MustCompile(`playerId="(\d+)"`).FindSubmatch(pageHTML)[1])
		playerName := string(regexp.MustCompile(`playerName="([^"]+)"`).FindSubmatch(pageHTML)[1])
		return playerID, playerName
	}
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/overview_inactive.html")
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
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/overview_inactive.html")
	for n := 0; n < b.N; n++ {
		extractUserGoquery(pageHTMLBytes)
	}
}

func TestWrapper(t *testing.T) {
	var bot Wrapper
	bot, _ = NewNoLogin("", "", "", "", "", "", "", 0, nil)
	assert.NotNil(t, bot)
}

func TestExtractCancelFleetTokenFromDocV71(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.5.0/en/cancel_fleet.html")
	token, _ := v71.NewExtractor().ExtractCancelFleetToken(pageHTMLBytes, ogame.FleetID(9078407))
	assert.Equal(t, "db3317fbe004641f7483e8074e34cda1", token)
}

func TestParseInt2(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/deathstar_price.html")
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTMLBytes))
	title := doc.Find("li.metal").AttrOr("title", "")
	metalStr := regexp.MustCompile(`([\d.]+)`).FindStringSubmatch(title)[1]
	metal := utils.ParseInt(metalStr)
	assert.Equal(t, int64(5000000), metal)

	pageHTMLBytes, _ = ioutil.ReadFile("../../samples/unversioned/mrd_price.html")
	doc, _ = goquery.NewDocumentFromReader(bytes.NewReader(pageHTMLBytes))
	title = doc.Find("li.metal").AttrOr("title", "")
	metalStr = regexp.MustCompile(`([\d.]+)`).FindStringSubmatch(title)[1]
	metal = utils.ParseInt(metalStr)
	assert.Equal(t, int64(1555733200), metal)
}

func TestExtractFleetDeutSaveFactor_V6_2_2_1(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/overview_active.html")
	res := v6.NewExtractor().ExtractFleetDeutSaveFactor(pageHTMLBytes)
	assert.Equal(t, 1.0, res)
}

func TestExtractFleetDeutSaveFactor_V6_7_4(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/overview_with_moon.html")
	res := v6.NewExtractor().ExtractFleetDeutSaveFactor(pageHTMLBytes)
	assert.Equal(t, 0.5, res)
}

func TestExtractPlanetCoordinate(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/station.html")
	res, _ := v6.NewExtractor().ExtractPlanetCoordinate(pageHTMLBytes)
	assert.Equal(t, ogame.Coordinate{1, 301, 5, ogame.PlanetType}, res)
}

func TestExtractPlanetCoordinate_moon(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/moon_facilities.html")
	res, _ := v6.NewExtractor().ExtractPlanetCoordinate(pageHTMLBytes)
	assert.Equal(t, ogame.Coordinate{4, 116, 12, ogame.MoonType}, res)
}

func TestExtractPlanetID_planet(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/station.html")
	res, _ := v6.NewExtractor().ExtractPlanetID(pageHTMLBytes)
	assert.Equal(t, ogame.CelestialID(33672410), res)
}

func TestExtractPlanetID_moon(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/moon_facilities.html")
	res, _ := v6.NewExtractor().ExtractPlanetID(pageHTMLBytes)
	assert.Equal(t, ogame.CelestialID(33741598), res)
}

func TestExtractPlanetType_planet(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/station.html")
	res, _ := v6.NewExtractor().ExtractPlanetType(pageHTMLBytes)
	assert.Equal(t, ogame.PlanetType, res)
}

func TestExtractPlanetType_moon(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/moon_facilities.html")
	res, _ := v6.NewExtractor().ExtractPlanetType(pageHTMLBytes)
	assert.Equal(t, ogame.MoonType, res)
}

func TestExtractJumpGate_cooldown(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/jumpgatelayer_charge.html")
	_, _, _, wait := v6.NewExtractor().ExtractJumpGate(pageHTMLBytes)
	assert.Equal(t, int64(1730), wait)
}

func TestExtractJumpGate(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/jumpgatelayer.html")
	ships, token, dests, wait := v6.NewExtractor().ExtractJumpGate(pageHTMLBytes)
	assert.Equal(t, 1, len(dests))
	assert.Equal(t, ogame.MoonID(33743183), dests[0])
	assert.Equal(t, int64(0), wait)
	assert.Equal(t, "7787b530670bc89623b5d65a827e557a", token)
	assert.Equal(t, int64(0), ships.SmallCargo)
	assert.Equal(t, int64(101), ships.LargeCargo)
	assert.Equal(t, int64(1), ships.LightFighter)
}

func TestExtractOgameTimestamp(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/moon_facilities.html")
	res := v6.NewExtractor().ExtractOgameTimestamp(pageHTMLBytes)
	assert.Equal(t, int64(1538912592), res)
}

func TestExtractOgameTimestampFromBytes(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/moon_facilities.html")
	res := v6.NewExtractor().ExtractOGameTimestampFromBytes(pageHTMLBytes)
	assert.Equal(t, int64(1538912592), res)
}

func TestExtractResources(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/moon_facilities.html")
	res := v6.NewExtractor().ExtractResources(pageHTMLBytes)
	assert.Equal(t, int64(280000), res.Metal)
	assert.Equal(t, int64(260000), res.Crystal)
	assert.Equal(t, int64(280000), res.Deuterium)
	assert.Equal(t, int64(0), res.Energy)
	assert.Equal(t, int64(25000), res.Darkmatter)
}

func TestExtractResourcesMobile(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/preferences_mobile.html")
	res := v6.NewExtractor().ExtractResources(pageHTMLBytes)
	assert.Equal(t, int64(7325851), res.Metal)
	assert.Equal(t, int64(1695823), res.Crystal)
	assert.Equal(t, int64(1835627), res.Deuterium)
	assert.Equal(t, int64(-2827), res.Energy)
	assert.Equal(t, int64(19500), res.Darkmatter)
}

func TestExtractResourcesDetailsFromFullPage(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/fleets_1.html")
	res := v6.NewExtractor().ExtractResourcesDetailsFromFullPage(pageHTMLBytes)
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
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7/overview2.html")
	res := v6.NewExtractor().ExtractResourcesDetailsFromFullPage(pageHTMLBytes)
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
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.5.1/en/phalanx_returning.html")
	res, err := v6.NewExtractor().ExtractPhalanx(pageHTMLBytes)
	clock := clockwork.NewFakeClockAt(time.Date(2020, 11, 4, 0, 25, 29, 0, time.UTC))
	assert.Nil(t, err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, ogame.Transport, res[0].Mission)
	assert.Equal(t, true, res[0].ReturnFlight)
	assert.NotNil(t, res[0].ArriveIn)
	assert.Equal(t, clock.Now().Add(10*time.Minute), res[0].ArrivalTime.UTC())
	assert.Equal(t, ogame.Coordinate{4, 116, 9, ogame.PlanetType}, res[0].Origin)
	assert.Equal(t, ogame.Coordinate{4, 116, 10, ogame.PlanetType}, res[0].Destination)
	assert.Equal(t, int64(19), res[0].Ships.SmallCargo)
}

func TestExtractPhalanx(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/phalanx.html")
	res, err := v6.NewExtractor().ExtractPhalanx(pageHTMLBytes)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, ogame.MissionID(3), res[0].Mission)
	assert.Equal(t, true, res[0].ReturnFlight)
	assert.NotNil(t, res[0].ArriveIn)
	assert.Equal(t, ogame.Coordinate{4, 116, 9, ogame.PlanetType}, res[0].Origin)
	assert.Equal(t, ogame.Coordinate{4, 212, 8, ogame.PlanetType}, res[0].Destination)
	assert.Equal(t, int64(100), res[0].Ships.LargeCargo)
}

func TestExtractPhalanx_fromMoon(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/phalanx_from_moon.html")
	res, _ := v6.NewExtractor().ExtractPhalanx(pageHTMLBytes)
	assert.Equal(t, ogame.Coordinate{4, 116, 12, ogame.MoonType}, res[0].Origin)
	assert.Equal(t, ogame.Coordinate{4, 116, 9, ogame.PlanetType}, res[0].Destination)
}

func TestExtractPhalanx_manyFleets(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/phalanx_fleets.html")
	res, err := v6.NewExtractor().ExtractPhalanx(pageHTMLBytes)
	assert.Nil(t, err)
	assert.Equal(t, 12, len(res))
	assert.Equal(t, ogame.Expedition, res[0].Mission)
	assert.False(t, res[0].ReturnFlight)
	assert.Equal(t, ogame.Coordinate{4, 124, 9, ogame.PlanetType}, res[0].Origin)
	assert.Equal(t, ogame.Coordinate{4, 125, 16, ogame.PlanetType}, res[0].Destination)
	assert.Equal(t, int64(250), res[0].Ships.LargeCargo)
	assert.Equal(t, int64(1), res[0].Ships.EspionageProbe)
	assert.Equal(t, int64(1), res[0].Ships.Destroyer)

	assert.Equal(t, ogame.Expedition, res[8].Mission)
	assert.True(t, res[8].ReturnFlight)
	assert.Equal(t, ogame.Coordinate{4, 124, 9, ogame.PlanetType}, res[8].Origin)
	assert.Equal(t, ogame.Coordinate{4, 125, 16, ogame.PlanetType}, res[8].Destination)
	assert.Equal(t, int64(250), res[8].Ships.LargeCargo)
	assert.Equal(t, int64(1), res[8].Ships.EspionageProbe)
	assert.Equal(t, int64(1), res[8].Ships.Destroyer)
}

func TestExtractPhalanx_noFleet(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/phalanx_no_fleet.html")
	res, err := v6.NewExtractor().ExtractPhalanx(pageHTMLBytes)
	assert.Equal(t, 0, len(res))
	assert.Nil(t, err)
}

func TestExtractPhalanx_noDeut(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/phalanx_no_deut.html")
	res, err := v6.NewExtractor().ExtractPhalanx(pageHTMLBytes)
	assert.Equal(t, 0, len(res))
	assert.NotNil(t, err)
}

func TestExtractResearch(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/research_bonus.html")
	res := v6.NewExtractor().ExtractResearch(pageHTMLBytes)
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

func TestExtractResearchV7(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7/researches.html")
	res := v7.NewExtractor().ExtractResearch(pageHTMLBytes)
	assert.Equal(t, int64(2), res.EnergyTechnology)
	assert.Equal(t, int64(4), res.LaserTechnology)
	assert.Equal(t, int64(0), res.IonTechnology)
	assert.Equal(t, int64(0), res.HyperspaceTechnology)
	assert.Equal(t, int64(0), res.PlasmaTechnology)
	assert.Equal(t, int64(5), res.CombustionDrive)
	assert.Equal(t, int64(4), res.ImpulseDrive)
	assert.Equal(t, int64(0), res.HyperspaceDrive)
	assert.Equal(t, int64(4), res.EspionageTechnology)
	assert.Equal(t, int64(1), res.ComputerTechnology)
	assert.Equal(t, int64(3), res.Astrophysics)
	assert.Equal(t, int64(0), res.IntergalacticResearchNetwork)
	assert.Equal(t, int64(0), res.GravitonTechnology)
	assert.Equal(t, int64(0), res.WeaponsTechnology)
	assert.Equal(t, int64(0), res.ShieldingTechnology)
	assert.Equal(t, int64(4), res.ArmourTechnology)
}

func TestExtractResearchV7_2(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7/researches2.html")
	res := v7.NewExtractor().ExtractResearch(pageHTMLBytes)
	assert.Equal(t, int64(1), res.EnergyTechnology)
	assert.Equal(t, int64(0), res.LaserTechnology)
	assert.Equal(t, int64(0), res.IonTechnology)
	assert.Equal(t, int64(0), res.HyperspaceTechnology)
	assert.Equal(t, int64(0), res.PlasmaTechnology)
	assert.Equal(t, int64(3), res.CombustionDrive)
	assert.Equal(t, int64(1), res.ImpulseDrive)
	assert.Equal(t, int64(0), res.HyperspaceDrive)
	assert.Equal(t, int64(0), res.EspionageTechnology)
	assert.Equal(t, int64(1), res.ComputerTechnology)
	assert.Equal(t, int64(0), res.Astrophysics)
	assert.Equal(t, int64(0), res.IntergalacticResearchNetwork)
	assert.Equal(t, int64(0), res.GravitonTechnology)
	assert.Equal(t, int64(0), res.WeaponsTechnology)
	assert.Equal(t, int64(0), res.ShieldingTechnology)
	assert.Equal(t, int64(0), res.ArmourTechnology)
}

func TestExtractResourcesBuildings(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/resource_inconstruction.html")
	res, _ := v6.NewExtractor().ExtractResourcesBuildings(pageHTMLBytes)
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

func TestExtractResourcesBuildingsV7(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7/supplies.html")
	res, _ := v7.NewExtractor().ExtractResourcesBuildings(pageHTMLBytes)
	assert.Equal(t, int64(2), res.MetalMine)
	assert.Equal(t, int64(1), res.CrystalMine)
	assert.Equal(t, int64(2), res.DeuteriumSynthesizer)
	assert.Equal(t, int64(3), res.SolarPlant)
	assert.Equal(t, int64(0), res.FusionReactor)
	assert.Equal(t, int64(0), res.SolarSatellite)
	assert.Equal(t, int64(2), res.MetalStorage)
	assert.Equal(t, int64(3), res.CrystalStorage)
	assert.Equal(t, int64(1), res.DeuteriumTank)
}

func TestExtractFacilities(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/facility_inconstruction.html")
	res, _ := v6.NewExtractor().ExtractFacilities(pageHTMLBytes)
	assert.Equal(t, int64(7), res.RoboticsFactory)
	assert.Equal(t, int64(7), res.Shipyard)
	assert.Equal(t, int64(7), res.ResearchLab)
	assert.Equal(t, int64(0), res.AllianceDepot)
	assert.Equal(t, int64(0), res.MissileSilo)
	assert.Equal(t, int64(0), res.NaniteFactory)
	assert.Equal(t, int64(0), res.Terraformer)
	assert.Equal(t, int64(3), res.SpaceDock)
}

func TestExtractFacilitiesV7(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7/facilities.html")
	res, _ := v7.NewExtractor().ExtractFacilities(pageHTMLBytes)
	assert.Equal(t, int64(3), res.RoboticsFactory)
	assert.Equal(t, int64(7), res.Shipyard)
	assert.Equal(t, int64(6), res.ResearchLab)
	assert.Equal(t, int64(0), res.AllianceDepot)
	assert.Equal(t, int64(0), res.MissileSilo)
	assert.Equal(t, int64(0), res.NaniteFactory)
	assert.Equal(t, int64(0), res.Terraformer)
	assert.Equal(t, int64(0), res.SpaceDock)
}

func TestExtractMoonFacilities(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/moon_facilities.html")
	res, _ := v6.NewExtractor().ExtractFacilities(pageHTMLBytes)
	assert.Equal(t, int64(1), res.RoboticsFactory)
	assert.Equal(t, int64(2), res.Shipyard)
	assert.Equal(t, int64(3), res.LunarBase)
	assert.Equal(t, int64(4), res.SensorPhalanx)
	assert.Equal(t, int64(5), res.JumpGate)
}

func TestExtractMoonFacilitiesV71(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.1/en/moon_facilities.html")
	res, _ := v71.NewExtractor().ExtractFacilities(pageHTMLBytes)
	assert.Equal(t, int64(10), res.RoboticsFactory)
	assert.Equal(t, int64(1), res.Shipyard)
	assert.Equal(t, int64(10), res.LunarBase)
	assert.Equal(t, int64(6), res.SensorPhalanx)
	assert.Equal(t, int64(1), res.JumpGate)
}

func TestExtractDefense(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/defence.html")
	defense, _ := v6.NewExtractor().ExtractDefense(pageHTMLBytes)
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

func TestExtractDefenseV7(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7/defenses.html")
	defense, _ := v7.NewExtractor().ExtractDefense(pageHTMLBytes)
	assert.Equal(t, int64(0), defense.RocketLauncher)
	assert.Equal(t, int64(2), defense.LightLaser)
	assert.Equal(t, int64(0), defense.HeavyLaser)
	assert.Equal(t, int64(0), defense.GaussCannon)
	assert.Equal(t, int64(0), defense.IonCannon)
	assert.Equal(t, int64(0), defense.PlasmaTurret)
	assert.Equal(t, int64(0), defense.SmallShieldDome)
	assert.Equal(t, int64(0), defense.LargeShieldDome)
	assert.Equal(t, int64(0), defense.AntiBallisticMissiles)
	assert.Equal(t, int64(0), defense.InterplanetaryMissiles)
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

func TestExtractFleet1Ships(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/fleet1.html")
	s := v6.NewExtractor().ExtractFleet1Ships(pageHTMLBytes)
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
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/fleet1_no_ships.html")
	s := v6.NewExtractor().ExtractFleet1Ships(pageHTMLBytes)
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
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/overview_queues.html")
	planet, _ := v6.NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33677371))
	assert.Equal(t, "C1", planet.Name)
	assert.Equal(t, int64(14615), planet.Diameter)
	assert.Equal(t, int64(-2), planet.Temperature.Min)
	assert.Equal(t, int64(38), planet.Temperature.Max)
	assert.Equal(t, int64(35), planet.Fields.Built)
	assert.Equal(t, int64(238), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33677371), planet.ID)
	assert.Equal(t, ogame.Coordinate{1, 301, 8, ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_fr(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/fr_overview.html")
	planet, _ := v6.NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33629512))
	assert.Equal(t, "planète mère", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(48), planet.Temperature.Min)
	assert.Equal(t, int64(88), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(188), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33629512), planet.ID)
	assert.Equal(t, ogame.Coordinate{2, 180, 4, ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_de(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/de_overview.html")
	planet, _ := v6.NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33630447))
	assert.Equal(t, "Heimatplanet", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(21), planet.Temperature.Min)
	assert.Equal(t, int64(61), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(188), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33630447), planet.ID)
	assert.Equal(t, ogame.Coordinate{2, 175, 8, ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_dk(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/dk_overview.html")
	planet, _ := v6.NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33627426))
	assert.Equal(t, "Hjemme verden", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(-23), planet.Temperature.Min)
	assert.Equal(t, int64(17), planet.Temperature.Max)
	assert.Equal(t, int64(5), planet.Fields.Built)
	assert.Equal(t, int64(193), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33627426), planet.ID)
	assert.Equal(t, ogame.Coordinate{1, 148, 12, ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_es(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v6/es/shipyard.html")
	planet, _ := v6.NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33644981))
	assert.Equal(t, "Planeta Principal", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(-18), planet.Temperature.Min)
	assert.Equal(t, int64(22), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(193), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33644981), planet.ID)
	assert.Equal(t, ogame.Coordinate{2, 493, 10, ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_br(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v6/br/overview.html")
	planet, _ := v6.NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33633767))
	assert.Equal(t, "Planeta Principal", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(-13), planet.Temperature.Min)
	assert.Equal(t, int64(27), planet.Temperature.Max)
	assert.Equal(t, int64(5), planet.Fields.Built)
	assert.Equal(t, int64(193), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33633767), planet.ID)
	assert.Equal(t, ogame.Coordinate{1, 449, 12, ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_it(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v6/it/overview.html")
	planet, _ := v6.NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33634944))
	assert.Equal(t, "Pianeta Madre", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(28), planet.Temperature.Min)
	assert.Equal(t, int64(68), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(188), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33634944), planet.ID)
	assert.Equal(t, ogame.Coordinate{2, 58, 8, ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_jp(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/jp_overview.html")
	planet, _ := v6.NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33620484))
	assert.Equal(t, "ホームワールド", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(40), planet.Temperature.Min)
	assert.Equal(t, int64(80), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(188), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33620484), planet.ID)
	assert.Equal(t, ogame.Coordinate{1, 18, 4, ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_tw(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v6/tw/overview.html")
	planet, _ := v6.NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33626432))
	assert.Equal(t, "母星", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(29), planet.Temperature.Min)
	assert.Equal(t, int64(69), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(188), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33626432), planet.ID)
	assert.Equal(t, ogame.Coordinate{1, 206, 8, ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_hr(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v6/hr/overview.html")
	planet, _ := v6.NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33627961))
	assert.Equal(t, "Glavni Planet", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(-33), planet.Temperature.Min)
	assert.Equal(t, int64(7), planet.Temperature.Max)
	assert.Equal(t, int64(4), planet.Fields.Built)
	assert.Equal(t, int64(193), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33627961), planet.ID)
	assert.Equal(t, ogame.Coordinate{1, 236, 12, ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_no(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v6/no/overview.html")
	planet, _ := v6.NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33624646))
	assert.Equal(t, "Hjemmeverden", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(-24), planet.Temperature.Min)
	assert.Equal(t, int64(16), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(188), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33624646), planet.ID)
	assert.Equal(t, ogame.Coordinate{1, 99, 10, ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_ro(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.1/ro/overview.html")
	planet, _ := v71.NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33629199))
	assert.Equal(t, "Planeta Principala", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(31), planet.Temperature.Min)
	assert.Equal(t, int64(71), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(193), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33629199), planet.ID)
	assert.Equal(t, ogame.Coordinate{1, 185, 4, ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_sk(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v6/sk/overview.html")
	planet, _ := v6.NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33625241))
	assert.Equal(t, "Domovská planéta", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(-12), planet.Temperature.Min)
	assert.Equal(t, int64(28), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(163), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33625241), planet.ID)
	assert.Equal(t, ogame.Coordinate{1, 94, 10, ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_si(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.2/si/overview.html")
	planet, _ := v6.NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33625245))
	assert.Equal(t, "Glavni Planet", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(41), planet.Temperature.Min)
	assert.Equal(t, int64(81), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(188), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33625245), planet.ID)
	assert.Equal(t, ogame.Coordinate{1, 70, 6, ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_hu(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.2/hu/overview.html")
	planet, _ := v6.NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33621505))
	assert.Equal(t, "Otthon", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(-18), planet.Temperature.Min)
	assert.Equal(t, int64(22), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(188), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33621505), planet.ID)
	assert.Equal(t, ogame.Coordinate{1, 162, 12, ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_fi(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.2/fi/overview.html")
	planet, _ := v6.NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33625483))
	assert.Equal(t, "Kotimaailma", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(15), planet.Temperature.Min)
	assert.Equal(t, int64(55), planet.Temperature.Max)
	assert.Equal(t, int64(5), planet.Fields.Built)
	assert.Equal(t, int64(188), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33625483), planet.ID)
	assert.Equal(t, ogame.Coordinate{1, 94, 6, ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_ba(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.5.1/ba/overview.html")
	planet, _ := v6.NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33621433))
	assert.Equal(t, "Glavni Planet", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(70), planet.Temperature.Min)
	assert.Equal(t, int64(110), planet.Temperature.Max)
	assert.Equal(t, int64(5), planet.Fields.Built)
	assert.Equal(t, int64(193), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33621433), planet.ID)
	assert.Equal(t, ogame.Coordinate{1, 55, 4, ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_gr(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v6/gr/overview.html")
	planet, _ := v6.NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33629206))
	assert.Equal(t, "Κύριος Πλανήτης", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(37), planet.Temperature.Min)
	assert.Equal(t, int64(77), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(193), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33629206), planet.ID)
	assert.Equal(t, ogame.Coordinate{1, 182, 6, ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_mx(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v6/mx/overview.html")
	planet, _ := v6.NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33624669))
	assert.Equal(t, "Planeta Principal", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(33), planet.Temperature.Min)
	assert.Equal(t, int64(73), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(188), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33624669), planet.ID)
	assert.Equal(t, ogame.Coordinate{1, 390, 6, ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_cz(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v6/cz/overview.html")
	planet, _ := v6.NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33622822))
	assert.Equal(t, "Domovska planeta", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(-13), planet.Temperature.Min)
	assert.Equal(t, int64(27), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(188), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33622822), planet.ID)
	assert.Equal(t, ogame.Coordinate{1, 221, 12, ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_jp1(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v6/jp/overview.html")
	planet, _ := v6.NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33623513))
	assert.Equal(t, "ホームワールド", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(58), planet.Temperature.Min)
	assert.Equal(t, int64(98), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(188), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33623513), planet.ID)
	assert.Equal(t, ogame.Coordinate{1, 85, 4, ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_pl(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/pl_overview.html")
	planet, _ := v6.NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33669699))
	assert.Equal(t, "Planeta matka", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(-21), planet.Temperature.Min)
	assert.Equal(t, int64(19), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(188), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33669699), planet.ID)
	assert.Equal(t, ogame.Coordinate{4, 248, 12, ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_tr(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/tr_overview.html")
	planet, _ := v6.NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33650421))
	assert.Equal(t, "Ana Gezegen", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(9), planet.Temperature.Min)
	assert.Equal(t, int64(49), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(188), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33650421), planet.ID)
	assert.Equal(t, ogame.Coordinate{3, 143, 10, ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_pt(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/pt_overview.html")
	planet, _ := v6.NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33635398))
	assert.Equal(t, "Planeta Principal", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(33), planet.Temperature.Min)
	assert.Equal(t, int64(73), planet.Temperature.Max)
	assert.Equal(t, int64(4), planet.Fields.Built)
	assert.Equal(t, int64(193), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33635398), planet.ID)
	assert.Equal(t, ogame.Coordinate{2, 241, 6, ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_nl(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/nl_overview.html")
	planet, _ := v6.NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33624684))
	assert.Equal(t, "Hoofdplaneet", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(-47), planet.Temperature.Min)
	assert.Equal(t, int64(-7), planet.Temperature.Max)
	assert.Equal(t, int64(5), planet.Fields.Built)
	assert.Equal(t, int64(188), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33624684), planet.ID)
	assert.Equal(t, ogame.Coordinate{1, 178, 12, ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_ar(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v6/ar/overview.html")
	planet, _ := v6.NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33629527))
	assert.Equal(t, "Planeta Principal", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(37), planet.Temperature.Min)
	assert.Equal(t, int64(77), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(193), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33629527), planet.ID)
	assert.Equal(t, ogame.Coordinate{1, 367, 4, ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_ru(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v6/ru/overview.html")
	planet, _ := v6.NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33629521))
	assert.Equal(t, "Главная планета", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(23), planet.Temperature.Min)
	assert.Equal(t, int64(63), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(163), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33629521), planet.ID)
	assert.Equal(t, ogame.Coordinate{1, 374, 6, ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_notExists(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/overview_queues.html")
	_, err := v6.NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(12345))
	assert.NotNil(t, err)
}

func TestExtractPlanetByCoord(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/overview_queues.html")
	planet, _ := v6.NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.Coordinate{1, 301, 8, ogame.PlanetType})
	assert.Equal(t, "C1", planet.Name)
	assert.Equal(t, int64(14615), planet.Diameter)
}

func TestExtractPlanetByCoord_notExists(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/overview_queues.html")
	_, err := v6.NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.Coordinate{1, 2, 3, ogame.PlanetType})
	assert.NotNil(t, err)
}

func TestFindSlowestSpeed(t *testing.T) {
	assert.Equal(t, int64(8000), findSlowestSpeed(ogame.ShipsInfos{SmallCargo: 1, LargeCargo: 1}, ogame.Researches{CombustionDrive: 6}, false, false))
}

func TestExtractShips(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/shipyard_thousands_ships.html")
	ships, _ := v6.NewExtractor().ExtractShips(pageHTMLBytes)
	assert.Equal(t, int64(1000), ships.LargeCargo)
	assert.Equal(t, int64(1000), ships.EspionageProbe)
	assert.Equal(t, int64(700), ships.Cruiser)
}

func TestExtractShipsV7(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7/shipyard.html")
	ships, _ := v7.NewExtractor().ExtractShips(pageHTMLBytes)
	assert.Equal(t, int64(6), ships.SmallCargo)
	assert.Equal(t, int64(1), ships.ColonyShip)
	assert.Equal(t, int64(9), ships.Crawler)
}

func TestExtractShipsV7_build(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7/shipyard_build.html")
	ships, _ := v7.NewExtractor().ExtractShips(pageHTMLBytes)
	assert.Equal(t, int64(33), ships.Cruiser)
}

func TestExtractShipsV7_fleetdispatch(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7/fleetdispatch.html")
	ships := v7.NewExtractor().ExtractFleet1Ships(pageHTMLBytes)
	assert.Equal(t, int64(6), ships.SmallCargo)
	assert.Equal(t, int64(1), ships.ColonyShip)
	assert.Equal(t, int64(0), ships.Crawler)
}

func TestExtractShipsMillions(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/shipyard_millions_ships.html")
	ships, _ := v6.NewExtractor().ExtractShips(pageHTMLBytes)
	assert.Equal(t, int64(15000001), ships.LightFighter)
}

func TestExtractShipsWhileBeingBuilt(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/shipyard_ship_being_built.html")
	ships, _ := v6.NewExtractor().ExtractShips(pageHTMLBytes)
	assert.Equal(t, int64(213), ships.EspionageProbe)
}

func TestExtractExpeditionMessages(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.2/en/expedition_messages.html")
	e := v7.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	msgs, nbPages, _ := e.ExtractExpeditionMessages(pageHTMLBytes)
	assert.Equal(t, int64(10), nbPages)
	assert.Equal(t, 10, len(msgs))
	assert.Equal(t, time.Date(2020, 04, 21, 23, 12, 6, 0, time.UTC), msgs[0].CreatedAt.UTC())
	assert.Equal(t, int64(11199359), msgs[0].ID)
	assert.Equal(t, ogame.Coordinate{1, 8, 16, ogame.PlanetType}, msgs[0].Coordinate)
	assert.Equal(t, `We came across the remains of a previous expedition! Our technicians will try to get some of the ships to work again.<br/><br/>The following ships are now part of the fleet:<br/>Espionage Probe: 1880<br/>Light Fighter: 161<br/>Small Cargo: 156`,
		msgs[0].Content)
}

func TestExtractMarketplaceMessages(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.2/en/sales_messages.html")
	e := v7.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	msgs, _, _ := e.ExtractMarketplaceMessages(pageHTMLBytes)
	assert.Equal(t, 9, len(msgs))
	assert.Equal(t, int64(12912161), msgs[3].ID)
	assert.Equal(t, int64(27), msgs[3].Type)
	assert.Equal(t, int64(1379), msgs[3].MarketTransactionID)
	assert.Equal(t, "164ba9f6e5cbfdaa03c061730767d779", msgs[3].Token)
}

func TestExtractEspionageReportMessageIDs(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/messages.html")
	msgs, _ := v6.NewExtractor().ExtractEspionageReportMessageIDs(pageHTMLBytes)
	assert.Equal(t, 2, len(msgs))
	assert.Equal(t, ogame.Report, msgs[0].Type)
	assert.Equal(t, ogame.Coordinate{4, 117, 6, ogame.PlanetType}, msgs[0].Target)
	assert.Equal(t, 0.5, msgs[0].LootPercentage)
	assert.Equal(t, "Fleet Command", msgs[0].From)
	assert.Equal(t, ogame.Action, msgs[1].Type)
	assert.Equal(t, "Space Monitoring", msgs[1].From)
	assert.Equal(t, ogame.Coordinate{4, 117, 9, ogame.PlanetType}, msgs[1].Target)
}

func TestExtractEspionageReportMessageIDsLootPercentage(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/messages_loot_percentage.html")
	msgs, _ := v6.NewExtractor().ExtractEspionageReportMessageIDs(pageHTMLBytes)
	assert.Equal(t, 1.0, msgs[0].LootPercentage)
	assert.Equal(t, 0.5, msgs[1].LootPercentage)
	assert.Equal(t, 0.5, msgs[2].LootPercentage)
}

func TestV71ExtractEspionageReportMessages(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/messages_loot_percentage.html")
	msgs, _ := v71.NewExtractor().ExtractEspionageReportMessageIDs(pageHTMLBytes)
	assert.Equal(t, 1.0, msgs[0].LootPercentage)
	assert.Equal(t, 0.5, msgs[1].LootPercentage)
	assert.Equal(t, 0.5, msgs[2].LootPercentage)
}

func TestExtractCombatReportMessagesV7(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7/combat_reports_msgs.html")
	msgs, _ := v7.NewExtractor().ExtractCombatReportMessagesSummary(pageHTMLBytes)
	assert.Equal(t, 10, len(msgs))
}

func TestExtractCombatReportMessagesV7_Debris(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7/combat_reports_debris.html")
	msgs, _ := v7.NewExtractor().ExtractCombatReportMessagesSummary(pageHTMLBytes)
	assert.Equal(t, int64(2400), msgs[0].DebrisField)
}

func TestExtractCombatReportMessagesV71(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.1/en/combat_reports.html")
	msgs, _ := v71.NewExtractor().ExtractCombatReportMessagesSummary(pageHTMLBytes)
	assert.Equal(t, "cr-us-149-fe449460902860455db7ef57a522ae341f931a59", msgs[0].APIKey)
}

func TestExtractCombatReportMessagesV71_lossContact(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.1/en/combat_reports_loss_contact.html")
	msgs, _ := v71.NewExtractor().ExtractCombatReportMessagesSummary(pageHTMLBytes)
	assert.Equal(t, 8, len(msgs))
}

func TestExtractCombatReportMessages(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/combat_reports_msgs.html")
	msgs, _ := v6.NewExtractor().ExtractCombatReportMessagesSummary(pageHTMLBytes)
	assert.Equal(t, 9, len(msgs))
}

func TestExtractCombatReportAttackingMessages(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/combat_reports_msgs_attacking.html")
	msgs, _ := v6.NewExtractor().ExtractCombatReportMessagesSummary(pageHTMLBytes)
	assert.Equal(t, int64(7945368), msgs[0].ID)
	assert.Equal(t, ogame.Coordinate{4, 233, 11, ogame.PlanetType}, msgs[0].Destination)
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
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/combat_reports_msgs_2.html")
	msgs, nbPages := v6.NewExtractor().ExtractCombatReportMessagesSummary(pageHTMLBytes)
	assert.Equal(t, 10, len(msgs))
	assert.Equal(t, int64(44), nbPages)
	assert.Equal(t, ogame.Coordinate{4, 116, 12, ogame.MoonType}, msgs[1].Destination)
	assert.Equal(t, ogame.Coordinate{4, 127, 9, ogame.MoonType}, *msgs[1].Origin)
}

func TestExtractResourcesProductions(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/resource_settings.html")
	prods, _ := v6.NewExtractor().ExtractResourcesProductions(pageHTMLBytes)
	assert.Equal(t, ogame.Resources{Metal: 10352, Crystal: 5104, Deuterium: 1282, Energy: -52}, prods)
}

func TestExtractResourceSettings(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/resource_settings.html")
	settings, _ := v6.NewExtractor().ExtractResourceSettings(pageHTMLBytes)
	assert.Equal(t, ogame.ResourceSettings{MetalMine: 100, CrystalMine: 100, DeuteriumSynthesizer: 100, SolarPlant: 100, FusionReactor: 0, SolarSatellite: 100}, settings)
}

func TestExtractResourceSettingsV7(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7/resource_settings.html")
	settings, _ := v7.NewExtractor().ExtractResourceSettings(pageHTMLBytes)
	assert.Equal(t, ogame.ResourceSettings{MetalMine: 100, CrystalMine: 100, DeuteriumSynthesizer: 100, SolarPlant: 100, FusionReactor: 0, SolarSatellite: 0, Crawler: 0}, settings)
}

func TestExtractNbProbes(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/preferences.html")
	probes := v6.NewExtractor().ExtractSpioAnz(pageHTMLBytes)
	assert.Equal(t, int64(10), probes)

	pageHTMLBytes, _ = ioutil.ReadFile("../../samples/unversioned/preferences_mobile.html")
	probes = v6.NewExtractor().ExtractSpioAnz(pageHTMLBytes)
	assert.Equal(t, int64(3), probes)
}

func TestExtractPreferencesShowActivityMinutes(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/preferences.html")
	checked := v6.NewExtractor().ExtractPreferencesShowActivityMinutes(pageHTMLBytes)
	assert.True(t, checked)

	pageHTMLBytes, _ = ioutil.ReadFile("../../samples/unversioned/preferences_mobile.html")
	checked = v6.NewExtractor().ExtractPreferencesShowActivityMinutes(pageHTMLBytes)
	assert.True(t, checked)

	pageHTMLBytes, _ = ioutil.ReadFile("../../samples/unversioned/preferences_without_detailed_activities.html")
	checked = v6.NewExtractor().ExtractPreferencesShowActivityMinutes(pageHTMLBytes)
	assert.False(t, checked)
}

func TestExtractPreferences(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/preferences.html")
	prefs := v6.NewExtractor().ExtractPreferences(pageHTMLBytes)
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

	pageHTMLBytes, _ = ioutil.ReadFile("../../samples/unversioned/preferences_reverse.html")
	prefs = v6.NewExtractor().ExtractPreferences(pageHTMLBytes)
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

	pageHTMLBytes, _ = ioutil.ReadFile("../../samples/unversioned/preferences_mobile.html")
	prefs = v6.NewExtractor().ExtractPreferences(pageHTMLBytes)
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

	pageHTMLBytes, _ = ioutil.ReadFile("../../samples/unversioned/preferences_reverse_mobile.html")
	prefs = v6.NewExtractor().ExtractPreferences(pageHTMLBytes)
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
//	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/traderOverview.html")
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

func TestExtractOfferOfTheDayPriceV874(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v8.7.4/en/traderImportExport.html")
	price, token, _, _, _ := v874.NewExtractor().ExtractOfferOfTheDay(pageHTMLBytes)
	assert.Equal(t, int64(178224), price)
	assert.Equal(t, "2a38193e2fa6047e1d92d2f2c71c00fd", token)
}

func TestExtractOfferOfTheDayPrice(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/traderOverview.html")
	price, token, _, _, _ := v6.NewExtractor().ExtractOfferOfTheDay(pageHTMLBytes)
	assert.Equal(t, int64(54243), price)
	assert.Equal(t, "8128c0ba0c9981599a87d818003f95e1", token)
}

func TestExtractOfferOfTheDayPrice1(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.4/en/traderOverview.html")
	price, token, _, _, _ := v6.NewExtractor().ExtractOfferOfTheDay(pageHTMLBytes)
	assert.Equal(t, int64(822159), price)
	assert.Equal(t, "2c829372796443bf6994cbfa051e4cd2", token)
}

func TestExtractCargoCapacity(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/sendfleet3.htm")
	fleet3Doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTMLBytes))
	cargo := utils.ParseInt(fleet3Doc.Find("#maxresources").Text())
	assert.Equal(t, int64(442500), cargo)
}

func TestExtractGalaxyInfos_vacationMode(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7/galaxy_vacation_mode.html")
	_, err := v6.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.EqualError(t, err, "account in vacation mode")
}

func TestExtractGalaxyInfos_bandit(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/galaxy_inactive_bandit_lord.html")
	infos, _ := v6.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.True(t, infos.Position(6).Player.IsBandit)
	assert.False(t, infos.Position(6).Player.IsStarlord)
}

func TestExtractGalaxyInfos_starlord(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/galaxy_inactive_emperor.html")
	infos, _ := v6.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.True(t, infos.Position(7).Player.IsStarlord)
	assert.False(t, infos.Position(7).Player.IsBandit)
}

func TestExtractGalaxyInfos_destroyedPlanet(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/galaxy_destroyed_planet.html")
	infos, _ := v6.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.NotNil(t, infos.Position(8))
	assert.True(t, infos.Position(8).Destroyed)
}

func TestExtractGalaxyInfos_destroyedPlanetAndMoon(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/galaxy_destroyed_planet_and_moon2.html")
	infos, _ := v6.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.NotNil(t, infos.Position(15))
	assert.True(t, infos.Position(15).Destroyed)
}

func TestExtractGalaxyInfos_banned(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/galaxy_banned.html")
	infos, _ := v6.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, true, infos.Position(1).Banned)
	assert.Equal(t, false, infos.Position(9).Banned)
}

func TestExtractGalaxyV7ExpeditionDebrisDM(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.1/fr/galaxy_darkmatter_df.html")
	infos, err := v7.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.NoError(t, err)
	assert.Equal(t, int64(3137), infos.Events.Darkmatter)
	assert.False(t, infos.Events.HasAsteroid)
}

func TestExtractGalaxyAsteroid(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.2/en/galaxyContent_asteroid.html")
	infos, err := v7.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.NoError(t, err)
	assert.True(t, infos.Events.HasAsteroid)
}

func TestExtractGalaxyV7ExpeditionDebris(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7/galaxy_debris16.html")
	infos, err := v7.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), infos.ExpeditionDebris.Metal)
	assert.Equal(t, int64(2300), infos.ExpeditionDebris.Crystal)
	assert.Equal(t, int64(1), infos.ExpeditionDebris.PathfindersNeeded)
}

func TestExtractGalaxyV752TWExpeditionDebris(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.5.2/tw/galaxy_debris16.html")
	infos, err := v7.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.NoError(t, err)
	assert.Equal(t, int64(4275000), infos.ExpeditionDebris.Metal)
	assert.Equal(t, int64(2953000), infos.ExpeditionDebris.Crystal)
	assert.Equal(t, int64(467), infos.ExpeditionDebris.PathfindersNeeded)
}

func TestExtractGalaxyV7ExpeditionDebris2(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7/galaxy_debris16_2.html")
	infos, err := v7.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.NoError(t, err)
	assert.Equal(t, int64(7200), infos.ExpeditionDebris.Metal)
	assert.Equal(t, int64(7200), infos.ExpeditionDebris.Crystal)
	assert.Equal(t, int64(1), infos.ExpeditionDebris.PathfindersNeeded)
}

func TestExtractGalaxyV7ExpeditionDebrisMobile(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7/galaxy_debris16_mobile.html")
	_, err := v7.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.EqualError(t, err, "mobile view not supported")
}

func TestExtractGalaxyV7IgnoredPlayer(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7/galaxy_ignored_player.html")
	infos, err := v7.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.NoError(t, err)
	assert.Equal(t, int64(33637068), infos.Position(7).ID)
	assert.Equal(t, int64(102418), infos.Position(7).Player.ID)
	assert.Equal(t, "Procurator Serpentis", infos.Position(7).Player.Name)
}

func TestExtractGalaxyV7NoExpeditionDebris(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7/galaxy_no_debris16.html")
	infos, err := v7.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), infos.ExpeditionDebris.Metal)
	assert.Equal(t, int64(0), infos.ExpeditionDebris.Crystal)
	assert.Equal(t, int64(0), infos.ExpeditionDebris.PathfindersNeeded)
}

func TestExtractGalaxyInfos(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/galaxy_ajax.html")
	infos, _ := v6.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
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
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/galaxy_ajax.html")
	infos, _ := v6.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, int64(33698658), infos.Position(12).ID)
	assert.Equal(t, "Commodore Nomade", infos.Position(12).Player.Name)
	assert.Equal(t, int64(123), infos.Position(12).Player.ID)
	assert.Equal(t, int64(456), infos.Position(12).Player.Rank)
	assert.Equal(t, "Homeworld", infos.Position(12).Name)
}

func TestExtractGalaxyInfosPlanetNoActivity(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/galaxy_planet_activity.html")
	infos, _ := v6.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, int64(0), infos.Position(15).Activity)
}

func TestExtractGalaxyInfosPlanetActivity15(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/galaxy_planet_activity.html")
	infos, _ := v6.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, int64(15), infos.Position(8).Activity)
}

func TestExtractGalaxyInfosPlanetActivity23(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/galaxy_planet_activity.html")
	infos, _ := v6.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, int64(23), infos.Position(9).Activity)
}

//func TestExtractGalaxyInfosPlanetActivityWithoutDetailedActivity(t *testing.T) {
//	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/galaxy_planet_activity_without_detailed_activity.html")
//	infos, _ := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
//	assert.Equal(t, 49, infos.Position(5).Activity)
//}

func TestExtractGalaxyInfosMoonActivity(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/galaxy_moon_activity.html")
	infos, _ := v6.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, int64(33732827), infos.Position(3).Moon.ID)
	assert.Equal(t, int64(5830), infos.Position(3).Moon.Diameter)
	assert.Equal(t, int64(18), infos.Position(3).Moon.Activity)
}

func TestExtractGalaxyInfosMoonNoActivity(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/galaxy_moon_no_activity.html")
	infos, _ := v6.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, int64(33650476), infos.Position(2).Moon.ID)
	assert.Equal(t, int64(7897), infos.Position(2).Moon.Diameter)
	assert.Equal(t, int64(0), infos.Position(2).Moon.Activity)
}

func TestExtractGalaxyInfosMoonActivity15(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/galaxy_moon_activity_unprecise.html")
	infos, _ := v6.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, int64(0), infos.Position(11).Activity)
	assert.Equal(t, int64(33730993), infos.Position(11).Moon.ID)
	assert.Equal(t, int64(8944), infos.Position(11).Moon.Diameter)
	assert.Equal(t, int64(15), infos.Position(11).Moon.Activity)
}

func TestExtractUserInfosV7(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7/overview.html")
	e := v6.NewExtractor()
	e.SetLanguage("en")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(538), infos.Points)
	assert.Equal(t, int64(1402), infos.Rank)
	assert.Equal(t, int64(3179), infos.Total)
	assert.Equal(t, "Governor Meridian", infos.PlayerName)
}

func TestExtractUserInfos(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/overview_inactive.html")
	e := v6.NewExtractor()
	e.SetLanguage("en")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(1295), infos.Points)
}

func TestExtractUserInfos_hr(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v6/hr/overview.html")
	e := v6.NewExtractor()
	e.SetLanguage("hr")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(214), infos.Rank)
	assert.Equal(t, int64(252), infos.Total)
}

func TestExtractUserInfos_tw(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v6/tw/overview.html")
	e := v6.NewExtractor()
	e.SetLanguage("tw")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(212), infos.Rank)
	assert.Equal(t, int64(212), infos.Total)
}

func TestExtractUserInfos_no(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v6/no/overview.html")
	e := v6.NewExtractor()
	e.SetLanguage("no")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(49), infos.Rank)
	assert.Equal(t, int64(50), infos.Total)
}

func TestExtractUserInfos_sk(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v6/sk/overview.html")
	e := v6.NewExtractor()
	e.SetLanguage("sk")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(89), infos.Rank)
	assert.Equal(t, int64(90), infos.Total)
}

func TestExtractUserInfos_fi(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.2/fi/overview.html")
	e := v6.NewExtractor()
	e.SetLanguage("fi")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(46), infos.Rank)
	assert.Equal(t, int64(51), infos.Total)
}

func TestExtractUserInfos_si(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.2/si/overview.html")
	e := v6.NewExtractor()
	e.SetLanguage("si")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(59), infos.Rank)
	assert.Equal(t, int64(60), infos.Total)
}

func TestExtractUserInfos_hu(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.2/hu/overview.html")
	e := v6.NewExtractor()
	e.SetLanguage("hu")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(635), infos.Rank)
	assert.Equal(t, int64(636), infos.Total)
}

func TestExtractUserInfos_gr(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v6/gr/overview.html")
	e := v6.NewExtractor()
	e.SetLanguage("gr")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(162), infos.Rank)
	assert.Equal(t, int64(163), infos.Total)
}

func TestExtractUserInfos_ro(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.1/ro/overview.html")
	e := v6.NewExtractor()
	e.SetLanguage("ro")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(108), infos.Rank)
	assert.Equal(t, int64(109), infos.Total)
}

func TestExtractUserInfos_mx(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v6/mx/overview.html")
	e := v6.NewExtractor()
	e.SetLanguage("mx")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(916), infos.Rank)
	assert.Equal(t, int64(917), infos.Total)
}

func TestExtractUserInfos_de(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/de_overview.html")
	e := v6.NewExtractor()
	e.SetLanguage("de")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(2980), infos.Rank)
	assert.Equal(t, int64(2980), infos.Total)
}

func TestExtractUserInfos_dk(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/dk_overview.html")
	e := v6.NewExtractor()
	e.SetLanguage("dk")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(253), infos.Rank)
	assert.Equal(t, int64(254), infos.Total)
	assert.Equal(t, "Procurator Zibal", infos.PlayerName)
}

func TestExtractUserInfos_jp(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/jp_overview.html")
	e := v6.NewExtractor()
	e.SetLanguage("jp")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(73), infos.Rank)
	assert.Equal(t, int64(73), infos.Total)
}

func TestExtractUserInfos_jp1(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v6/jp/overview.html")
	e := v6.NewExtractor()
	e.SetLanguage("jp")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(85), infos.Rank)
	assert.Equal(t, int64(86), infos.Total)
}

func TestExtractUserInfos_cz(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v6/cz/overview.html")
	e := v6.NewExtractor()
	e.SetLanguage("cz")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(1008), infos.Rank)
	assert.Equal(t, int64(1009), infos.Total)
}

func TestExtractUserInfos_fr(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/fr_overview.html")
	e := v6.NewExtractor()
	e.SetLanguage("fr")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(3197), infos.Rank)
	assert.Equal(t, int64(3348), infos.Total)
	assert.Equal(t, "Bandit Pégasus", infos.PlayerName)
}

func TestExtractUserInfos_nl(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/nl_overview.html")
	e := v6.NewExtractor()
	e.SetLanguage("nl")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(482), infos.Rank)
	assert.Equal(t, int64(542), infos.Total)
	assert.Equal(t, "Bandit Japetus", infos.PlayerName)
}

func TestExtractUserInfos_pl(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/pl_overview.html")
	e := v6.NewExtractor()
	e.SetLanguage("pl")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(5873), infos.Rank)
	assert.Equal(t, int64(5876), infos.Total)
	assert.Equal(t, "Constable Leonis", infos.PlayerName)
}

func TestExtractUserInfos_br(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v6/br/overview.html")
	e := v6.NewExtractor()
	e.SetLanguage("br")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(1026), infos.Rank)
	assert.Equal(t, int64(1268), infos.Total)
}

func TestExtractUserInfos_tr(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/tr_overview.html")
	e := v6.NewExtractor()
	e.SetLanguage("tr")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(3655), infos.Rank)
	assert.Equal(t, int64(3656), infos.Total)
	assert.Equal(t, "Chief Apus", infos.PlayerName)
}

func TestExtractUserInfos_ar(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v6/ar/overview.html")
	e := v6.NewExtractor()
	e.SetLanguage("ar")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(1158), infos.Rank)
	assert.Equal(t, int64(1159), infos.Total)
	assert.Equal(t, "Chief Lambda", infos.PlayerName)
}

func TestExtractUserInfos_it(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v6/it/overview.html")
	e := v6.NewExtractor()
	e.SetLanguage("it")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(1776), infos.Rank)
	assert.Equal(t, int64(1777), infos.Total)
	assert.Equal(t, "President Fidis", infos.PlayerName)
}

func TestExtractUserInfos_pt(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/pt_overview.html")
	e := v6.NewExtractor()
	e.SetLanguage("pt")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(1762), infos.Rank)
	assert.Equal(t, int64(1862), infos.Total)
	assert.Equal(t, "Director Europa", infos.PlayerName)
}

func TestExtractUserInfos_ru(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v6/ru/overview.html")
	e := v6.NewExtractor()
	e.SetLanguage("ru")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(1067), infos.Rank)
	assert.Equal(t, int64(1068), infos.Total)
	assert.Equal(t, "Viceregent Horizon", infos.PlayerName)
}

func TestExtractUserInfos_ba(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.5.1/ba/overview.html")
	e := v6.NewExtractor()
	e.SetLanguage("ba")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(138), infos.Rank)
	assert.Equal(t, int64(139), infos.Total)
	assert.Equal(t, "Governor Hunter", infos.PlayerName)
}

func TestExtractUserInfos_es(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.6.5/es/overview.html")
	e := v7.NewExtractor()
	e.SetLanguage("es")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(2976), infos.Rank)
	assert.Equal(t, int64(2977), infos.Total)
	assert.Equal(t, "Commodore Navi", infos.PlayerName)
}

func TestExtractMoons(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/overview_with_moon.html")
	moons := v6.NewExtractor().ExtractMoons(pageHTMLBytes)
	assert.Equal(t, 1, len(moons))
}

func TestExtractMoons2(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/overview_with_many_moon.html")
	moons := v6.NewExtractor().ExtractMoons(pageHTMLBytes)
	assert.Equal(t, 2, len(moons))
}

func TestExtractMoon_exists(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/overview_with_moon.html")
	_, err := v6.NewExtractor().ExtractMoon(pageHTMLBytes, ogame.MoonID(33741598))
	assert.Nil(t, err)
}

func TestExtractMoon_notExists(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/overview_with_moon.html")
	_, err := v6.NewExtractor().ExtractMoon(pageHTMLBytes, ogame.MoonID(12345))
	assert.NotNil(t, err)
}

func TestExtractMoonByCoord_exists(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/overview_with_moon.html")
	_, err := v6.NewExtractor().ExtractMoon(pageHTMLBytes, ogame.Coordinate{4, 116, 12, ogame.MoonType})
	assert.Nil(t, err)
}

func TestExtractMoonByCoord_notExists(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/overview_with_moon.html")
	_, err := v6.NewExtractor().ExtractMoon(pageHTMLBytes, ogame.Coordinate{1, 2, 3, ogame.PlanetType})
	assert.NotNil(t, err)
}

func TestExtractIsInVacationFromDoc(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v6/es/overview_vacation.html")
	assert.True(t, v6.NewExtractor().ExtractIsInVacation(pageHTMLBytes))
	pageHTMLBytes, _ = ioutil.ReadFile("../../samples/v6/es/fleet1_vacation.html")
	assert.True(t, v6.NewExtractor().ExtractIsInVacation(pageHTMLBytes))
	pageHTMLBytes, _ = ioutil.ReadFile("../../samples/v6/es/shipyard.html")
	assert.False(t, v6.NewExtractor().ExtractIsInVacation(pageHTMLBytes))
}

func TestExtractPlanetsMoon(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/overview_with_moon.html")
	planets := v6.NewExtractor().ExtractPlanets(pageHTMLBytes)
	assert.Equal(t, ogame.MoonID(33741598), planets[0].Moon.ID)
	assert.Equal(t, "Moon", planets[0].Moon.Name)
	assert.Equal(t, "https://gf1.geo.gfsrv.net/cdn9d/8e0e6034049bd64e18a1804b42f179.gif", planets[0].Moon.Img)
	assert.Equal(t, int64(8774), planets[0].Moon.Diameter)
	assert.Equal(t, ogame.Coordinate{4, 116, 12, ogame.MoonType}, planets[0].Moon.Coordinate)
	assert.Equal(t, int64(0), planets[0].Moon.Fields.Built)
	assert.Equal(t, int64(1), planets[0].Moon.Fields.Total)
	assert.Nil(t, planets[1].Moon)
}

func TestExtractPlanets_fieldsFilled(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/overview_fields_filled.html")
	planets := v6.NewExtractor().ExtractPlanets(pageHTMLBytes)
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

func TestExtractPlanetsV9(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v9.0.0/en/overview.html")
	planets := v8.NewExtractor().ExtractPlanets(pageHTMLBytes)
	assert.Equal(t, 1, len(planets))
	assert.Equal(t, ogame.PlanetID(34071290), planets[0].ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 292, Position: 4, Type: ogame.PlanetType}, planets[0].Coordinate)
	assert.Equal(t, "Homeworld", planets[0].Name)
	assert.Equal(t, "https://gf2.geo.gfsrv.net/cdn7a/ca5a968aa62c0441a62334221eaa74.png", planets[0].Img)
	assert.Equal(t, int64(70), planets[0].Temperature.Min)
	assert.Equal(t, int64(110), planets[0].Temperature.Max)
	assert.Equal(t, int64(3), planets[0].Fields.Built)
	assert.Equal(t, int64(163), planets[0].Fields.Total)
	assert.Equal(t, int64(12800), planets[0].Diameter)
}

func TestExtractPlanetsEsV902(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v9.0.2/es/overview.html")
	planets := v6.NewExtractor().ExtractPlanets(pageHTMLBytes)
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
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v9.0.2/tw/overview.html")
	planets := v6.NewExtractor().ExtractPlanets(pageHTMLBytes)
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
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/overview_inactive.html")
	planets := v6.NewExtractor().ExtractPlanets(pageHTMLBytes)
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
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/overview_es.html")
	planets := v6.NewExtractor().ExtractPlanets(pageHTMLBytes)
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
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/fr_overview.html")
	planets := v6.NewExtractor().ExtractPlanets(pageHTMLBytes)
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
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v6/fr/overview.html")
	planets := v6.NewExtractor().ExtractPlanets(pageHTMLBytes)
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
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v6/br/overview.html")
	planets := v6.NewExtractor().ExtractPlanets(pageHTMLBytes)
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
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/galaxy_debris.html")
	infos, _ := v6.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.False(t, infos.Position(6).HonorableTarget)
	assert.True(t, infos.Position(8).HonorableTarget)
	assert.False(t, infos.Position(9).HonorableTarget)
}

func TestExtractGalaxyInfos_inactive(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/galaxy_debris.html")
	infos, _ := v6.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.True(t, infos.Position(6).Inactive)
	assert.False(t, infos.Position(8).Inactive)
	assert.False(t, infos.Position(9).Inactive)
}

func TestExtractGalaxyInfos_strongPlayer(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/galaxy_debris.html")
	infos, _ := v6.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.False(t, infos.Position(6).StrongPlayer)
	assert.True(t, infos.Position(8).StrongPlayer)
	assert.False(t, infos.Position(9).StrongPlayer)
}

func TestExtractGalaxyInfos_newbie(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/galaxy_newbie.html")
	infos, _ := v6.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.True(t, infos.Position(4).Newbie)
}

func TestExtractGalaxyInfos_moon(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/galaxy_debris.html")
	infos, _ := v6.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.NotNil(t, infos.Position(6).Moon)
	assert.Equal(t, int64(33701543), infos.Position(6).Moon.ID)
	assert.Equal(t, int64(8366), infos.Position(6).Moon.Diameter)
}

func TestExtractGalaxyInfos_debris(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/galaxy_debris.html")
	infos, _ := v6.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, int64(0), infos.Position(6).Debris.Metal)
	assert.Equal(t, int64(700), infos.Position(6).Debris.Crystal)
	assert.Equal(t, int64(1), infos.Position(6).Debris.RecyclersNeeded)
}

func TestExtractGalaxyInfos_debris_es(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/galaxy_debris_es.html")
	infos, _ := v6.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, int64(0), infos.Position(12).Debris.Metal)
	assert.Equal(t, int64(128000), infos.Position(12).Debris.Crystal)
	assert.Equal(t, int64(7), infos.Position(12).Debris.RecyclersNeeded)
}

func TestExtractGalaxyInfos_debris_fr(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/galaxy_debris_fr.html")
	infos, _ := v6.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, int64(100), infos.Position(7).Debris.Metal)
	assert.Equal(t, int64(600), infos.Position(7).Debris.Crystal)
	assert.Equal(t, int64(1), infos.Position(7).Debris.RecyclersNeeded)
}

func TestExtractGalaxyInfos_debris_de(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/galaxy_debris_de.html")
	infos, _ := v6.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, int64(100), infos.Position(9).Debris.Metal)
	assert.Equal(t, int64(2500), infos.Position(9).Debris.Crystal)
	assert.Equal(t, int64(1), infos.Position(9).Debris.RecyclersNeeded)
}

func TestExtractGalaxyInfos_vacation(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/galaxy_ajax.html")
	infos, _ := v6.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.False(t, infos.Position(4).Vacation)
	assert.True(t, infos.Position(6).Vacation)
	assert.True(t, infos.Position(8).Vacation)
	assert.False(t, infos.Position(10).Vacation)
	assert.False(t, infos.Position(12).Vacation)
}

func TestExtractGalaxyInfos_alliance(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/galaxy_ajax.html")
	infos, _ := v6.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, int64(303), infos.Position(10).Alliance.ID)
	assert.Equal(t, "Qrvix", infos.Position(10).Alliance.Name)
	assert.Equal(t, int64(27), infos.Position(10).Alliance.Rank)
	assert.Equal(t, int64(16), infos.Position(10).Alliance.Member)
}

func TestExtractGalaxyInfos_alliance_fr(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v6/fr/galaxy.html")
	infos, _ := v6.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, int64(635), infos.Position(5).Alliance.ID)
	assert.Equal(t, "leretour", infos.Position(5).Alliance.Name)
	assert.Equal(t, int64(24), infos.Position(5).Alliance.Rank)
	assert.Equal(t, int64(11), infos.Position(5).Alliance.Member)
}

func TestExtractGalaxyInfos_alliance_es(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v6/es/galaxy.html")
	infos, _ := v6.NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, int64(500053), infos.Position(4).Alliance.ID)
	assert.Equal(t, "Los Aliens Grises", infos.Position(4).Alliance.Name)
	assert.Equal(t, int64(8), infos.Position(4).Alliance.Rank)
	assert.Equal(t, int64(70), infos.Position(4).Alliance.Member)
}

func TestUniverseSpeed(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/techtree_universe_speed.html")
	universeSpeed := v6.ExtractUniverseSpeed(pageHTMLBytes)
	assert.Equal(t, int64(7), universeSpeed)
}

func TestCancel(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/overview_active_queue2.html")
	token, techID, listID, _ := v6.NewExtractor().ExtractCancelBuildingInfos(pageHTMLBytes)
	assert.Equal(t, "fef7488e4809150cd16e3fa8fa14db37", token)
	assert.Equal(t, int64(4), techID)
	assert.Equal(t, int64(2099434), listID)
}

func TestCancelV7(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7/overview_cancels.html")
	token, techID, listID, _ := v7.NewExtractor().ExtractCancelBuildingInfos(pageHTMLBytes)
	assert.Equal(t, "cf00a76b307f5cabf867af0d61ad1991", token)
	assert.Equal(t, int64(23), techID)
	assert.Equal(t, int64(1336041), listID)
}

func TestCancelBuildingV902(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v9.0.2/en/overview_all_queues.html")
	token, id, listID, _ := v9.NewExtractor().ExtractCancelBuildingInfos(pageHTMLBytes)
	assert.Equal(t, "66f639922a3c76fe6074d12ae36e573e", token)
	assert.Equal(t, int64(1), id)
	assert.Equal(t, int64(3469488), listID)
}

func TestCancelResearch(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/overview_active_queue2.html")
	token, techID, listID, _ := v6.NewExtractor().ExtractCancelResearchInfos(pageHTMLBytes)
	assert.Equal(t, "fff7488e4809150cd16e3fa8fa14db37", token)
	assert.Equal(t, int64(120), techID)
	assert.Equal(t, int64(1769925), listID)
}
func TestCancelResearchV7(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7/overview_cancels.html")
	token, techID, listID, _ := v7.NewExtractor().ExtractCancelResearchInfos(pageHTMLBytes)
	assert.Equal(t, "9d44b41d8136dffadab759749508105e", token)
	assert.Equal(t, int64(124), techID)
	assert.Equal(t, int64(1324883), listID)
}

func TestGetConstructions(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/overview_active.html")
	buildingID, buildingCountdown, researchID, researchCountdown := v6.NewExtractor().ExtractConstructions(pageHTMLBytes)
	assert.Equal(t, ogame.CrystalMineID, buildingID)
	assert.Equal(t, int64(731), buildingCountdown)
	assert.Equal(t, ogame.CombustionDriveID, researchID)
	assert.Equal(t, int64(927), researchCountdown)
}

func TestGetConstructionsV7(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7/overview_supplies_in_construction.html")
	clock := clockwork.NewFakeClockAt(time.Date(2019, 11, 12, 9, 6, 43, 0, time.UTC))
	buildingID, buildingCountdown, researchID, researchCountdown := v7.ExtractConstructions(pageHTMLBytes, clock)
	assert.Equal(t, ogame.MetalMineID, buildingID)
	assert.Equal(t, int64(62), buildingCountdown)
	assert.Equal(t, ogame.EnergyTechnologyID, researchID)
	assert.Equal(t, int64(271), researchCountdown)
}

func TestExtractFleetsFromEventList(t *testing.T) {
	//pageHTMLBytes, _ := ioutil.ReadFile("../../samples/eventlist_test.html")
	//fleets := NewExtractor().ExtractFleetsFromEventList(pageHTMLBytes)
	//assert.Equal(t, 4, len(fleets))
}

func TestExtractIPM(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/missileattacklayer.html")
	duration, max, token := v6.NewExtractor().ExtractIPM(pageHTMLBytes)
	assert.Equal(t, "26a08f4cc0c0b513e1e8c10d49c14a27", token)
	assert.Equal(t, int64(17), max)
	assert.Equal(t, int64(15), duration)
}

func TestExtractFleetV71(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.1/en/movement.html")
	e := v6.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	fleets := e.ExtractFleets(pageHTMLBytes)
	assert.Equal(t, 1, len(fleets))
	assert.Equal(t, int64(8271), fleets[0].ArriveIn)
	assert.Equal(t, int64(16545), fleets[0].BackIn)
	assert.Equal(t, ogame.Coordinate{1, 432, 6, ogame.PlanetType}, fleets[0].Origin)
	assert.Equal(t, ogame.Coordinate{1, 432, 5, ogame.PlanetType}, fleets[0].Destination)
	assert.Equal(t, ogame.Transport, fleets[0].Mission)
	assert.Equal(t, false, fleets[0].ReturnFlight)
	assert.Equal(t, ogame.FleetID(1674510), fleets[0].ID)
	assert.Equal(t, int64(250), fleets[0].Ships.SmallCargo)
	assert.Equal(t, int64(2), fleets[0].Ships.Pathfinder)
	assert.Equal(t, ogame.Resources{}, fleets[0].Resources)
}

func TestExtractFleetV72(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.2/de/movement.html")
	clock := clockwork.NewFakeClockAt(time.Date(2020, 3, 6, 11, 43, 15, 0, time.UTC))
	e := v6.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	fleets := e.ExtractFleets(pageHTMLBytes)
	assert.Equal(t, clock.Now().Add(-5031*time.Second), fleets[0].StartTime.UTC())
	assert.Equal(t, clock.Now().Add(-5041*time.Second), fleets[1].StartTime.UTC())
}

func TestExtractFleetV71_2(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.1/en/movement2.html")
	clock := clockwork.NewFakeClockAt(time.Date(2020, 1, 12, 1, 45, 34, 0, time.UTC))
	e := v6.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 0))
	fleets := e.ExtractFleets(pageHTMLBytes)
	assert.Equal(t, 2, len(fleets))
	assert.Equal(t, int64(621), fleets[0].ArriveIn)
	assert.Equal(t, int64(1245), fleets[0].BackIn)
	assert.Equal(t, ogame.Coordinate{4, 116, 12, ogame.MoonType}, fleets[0].Origin)
	assert.Equal(t, ogame.Coordinate{4, 116, 12, ogame.PlanetType}, fleets[0].Destination)
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
	assert.Equal(t, ogame.Coordinate{4, 208, 10, ogame.PlanetType}, fleets[1].Origin)
	assert.Equal(t, ogame.Coordinate{4, 116, 12, ogame.PlanetType}, fleets[1].Destination)
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
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.6.7/en/movement.html")
	e := v6.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	fleets := e.ExtractFleets(pageHTMLBytes)
	assert.Equal(t, 1, len(fleets))
	assert.Equal(t, time.Date(2021, 6, 1, 9, 28, 2, 0, time.UTC), fleets[0].StartTime.UTC())
	assert.Equal(t, time.Date(2021, 6, 1, 9, 51, 10, 0, time.UTC), fleets[0].ArrivalTime.UTC())
	assert.Equal(t, time.Date(2021, 6, 1, 10, 14, 18, 0, time.UTC), fleets[0].BackTime.UTC())
}

func TestExtractFleetV7(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7/movement.html")
	e := v6.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	fleets := e.ExtractFleets(pageHTMLBytes)
	assert.Equal(t, 1, len(fleets))
	assert.Equal(t, int64(1010), fleets[0].ArriveIn)
	assert.Equal(t, int64(2030), fleets[0].BackIn)
	assert.Equal(t, ogame.Coordinate{9, 297, 12, ogame.PlanetType}, fleets[0].Origin)
	assert.Equal(t, ogame.Coordinate{9, 297, 9, ogame.PlanetType}, fleets[0].Destination)
	assert.Equal(t, ogame.Transport, fleets[0].Mission)
	assert.Equal(t, false, fleets[0].ReturnFlight)
	assert.Equal(t, ogame.FleetID(4218727), fleets[0].ID)
	assert.Equal(t, int64(2), fleets[0].Ships.SmallCargo)
	assert.Equal(t, ogame.Resources{}, fleets[0].Resources)
}

func TestExtractFleet(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/fleets_1.html")
	e := v6.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	fleets := e.ExtractFleets(pageHTMLBytes)
	assert.Equal(t, 1, len(fleets))
	assert.Equal(t, int64(4134), fleets[0].ArriveIn)
	assert.Equal(t, int64(8277), fleets[0].BackIn)
	assert.Equal(t, ogame.Coordinate{4, 116, 12, ogame.PlanetType}, fleets[0].Origin)
	assert.Equal(t, ogame.Coordinate{4, 117, 9, ogame.PlanetType}, fleets[0].Destination)
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
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/fleets_expedition.html")
	e := v6.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	fleets := e.ExtractFleets(pageHTMLBytes)
	assert.Equal(t, 2, len(fleets))
	assert.Equal(t, int64(2), fleets[1].Ships.LargeCargo)
	assert.Equal(t, ogame.Expedition, fleets[1].Mission)
	assert.False(t, fleets[1].ReturnFlight)
	assert.Equal(t, ogame.Coordinate{4, 116, 12, ogame.MoonType}, fleets[1].Origin)
	assert.Equal(t, ogame.Coordinate{4, 116, 16, ogame.PlanetType}, fleets[1].Destination)
}

func TestExtractFleet_harvest(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/fleets_harvest.html")
	e := v6.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	fleets := e.ExtractFleets(pageHTMLBytes)
	assert.Equal(t, ogame.Coordinate{4, 116, 12, ogame.PlanetType}, fleets[5].Origin)
	assert.Equal(t, ogame.Coordinate{4, 116, 9, ogame.DebrisType}, fleets[5].Destination)
}

func TestExtractFleet_returningTransport(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/fleets_2.html")
	e := v6.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	fleets := e.ExtractFleets(pageHTMLBytes)
	assert.Equal(t, 1, len(fleets))
	assert.Equal(t, int64(-1), fleets[0].ArriveIn)
	assert.Equal(t, int64(36), fleets[0].BackIn)
}

func TestExtractFleet_deployment(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/fleets_moon_to_moon.html")
	e := v6.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	fleets := e.ExtractFleets(pageHTMLBytes)
	assert.Equal(t, int64(210), fleets[0].ArriveIn)
	assert.Equal(t, int64(426), fleets[0].BackIn)
}

func TestExtractFleetThousands(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/fleets_thousands.html")
	e := v6.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	fleets := e.ExtractFleets(pageHTMLBytes)
	assert.Equal(t, ogame.Transport, fleets[0].Mission)
	assert.Equal(t, int64(210), fleets[0].Ships.LargeCargo)
	assert.Equal(t, ogame.Resources{Metal: 207862, Crystal: 78903, Deuterium: 42956}, fleets[0].Resources)
}

func TestExtractFleet_returning(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/fleets_2.html")
	e := v6.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	fleets := e.ExtractFleets(pageHTMLBytes)
	assert.Equal(t, 1, len(fleets))
	assert.Equal(t, ogame.Coordinate{4, 116, 12, ogame.PlanetType}, fleets[0].Origin)
	assert.Equal(t, ogame.Coordinate{4, 117, 9, ogame.PlanetType}, fleets[0].Destination)
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
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.2/en/fleets_expeditions.html")
	e := v6.NewExtractor()
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
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/fleets_moon_to_moon.html")
	e := v6.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	fleets := e.ExtractFleets(pageHTMLBytes)
	assert.Equal(t, int64(0), fleets[0].TargetPlanetID)
	assert.Equal(t, int64(0), fleets[1].TargetPlanetID)
	assert.Equal(t, int64(33702114), fleets[2].TargetPlanetID)
	assert.Equal(t, int64(33699325), fleets[3].TargetPlanetID)
}

func TestExtractFleet_unionID(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/fleets_no_union.html")
	e := v6.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	fleets := e.ExtractFleets(pageHTMLBytes)
	assert.Equal(t, int64(0), fleets[0].UnionID)

	pageHTMLBytes, _ = ioutil.ReadFile("../../samples/unversioned/fleets_union_alone.html")
	e = v6.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	fleets = e.ExtractFleets(pageHTMLBytes)
	assert.Equal(t, int64(13558), fleets[0].UnionID)
}

func TestExtractOverviewProduction(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/overview_shipyard_queue_full.html")
	prods, countdown, _ := v6.NewExtractor().ExtractOverviewProduction(pageHTMLBytes)
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

func TestV71ExtractOverviewProduction(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.1/en/overview_shipyard_queue.html")
	prods, countdown, _ := v7.NewExtractor().ExtractOverviewProduction(pageHTMLBytes)
	assert.Equal(t, 4, len(prods))
	assert.Equal(t, int64(542), countdown)
	assert.Equal(t, ogame.GaussCannonID, prods[0].ID)
	assert.Equal(t, int64(2), prods[0].Nbr)
	assert.Equal(t, ogame.RocketLauncherID, prods[1].ID)
	assert.Equal(t, int64(2), prods[1].Nbr)
	assert.Equal(t, ogame.SmallShieldDomeID, prods[2].ID)
	assert.Equal(t, int64(1), prods[2].Nbr)
	assert.Equal(t, ogame.LightLaserID, prods[3].ID)
	assert.Equal(t, int64(3), prods[3].Nbr)
}

func TestExtractV71Production(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.1/en/shipyard_queue.html")
	prods, secs, _ := v71.NewExtractor().ExtractProduction(pageHTMLBytes)
	assert.Equal(t, 4, len(prods))
	assert.Equal(t, int64(977), secs)
	assert.Equal(t, ogame.SmallCargoID, prods[0].ID)
	assert.Equal(t, int64(12), prods[0].Nbr)
	assert.Equal(t, ogame.SmallCargoID, prods[1].ID)
	assert.Equal(t, int64(5), prods[1].Nbr)
	assert.Equal(t, ogame.LargeCargoID, prods[2].ID)
	assert.Equal(t, int64(3), prods[2].Nbr)
	assert.Equal(t, ogame.SmallCargoID, prods[3].ID)
	assert.Equal(t, int64(5), prods[3].Nbr)
}

func TestExtractProduction(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/shipyard_queue.html")
	prods, secs, _ := v6.NewExtractor().ExtractProduction(pageHTMLBytes)
	assert.Equal(t, 20, len(prods))
	assert.Equal(t, int64(16254), secs)
	assert.Equal(t, ogame.LargeCargoID, prods[0].ID)
	assert.Equal(t, int64(4), prods[0].Nbr)
}

func TestExtractProduction2(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/shipyard_queue2.html")
	prods, secs, _ := v6.NewExtractor().ExtractProduction(pageHTMLBytes)
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
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/production_with_abm.html")
	prods, secs, _ := v6.NewExtractor().ExtractProduction(pageHTMLBytes)
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
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v6/dk/production_with_abm.html")
	prods, secs, _ := v6.NewExtractor().ExtractProduction(pageHTMLBytes)
	assert.Equal(t, 2, len(prods))
	assert.Equal(t, int64(641), secs)
	assert.Equal(t, ogame.AntiBallisticMissilesID, prods[0].ID)
	assert.Equal(t, int64(1), prods[0].Nbr)
	assert.Equal(t, ogame.AntiBallisticMissilesID, prods[1].ID)
	assert.Equal(t, int64(1), prods[1].Nbr)
}

func TestI64Ptr(t *testing.T) {
	v := int64(6)
	assert.Equal(t, &v, utils.I64Ptr(6))
}

func TestExtractEspionageReport_tz(t *testing.T) {
	clock := clockwork.NewFakeClockAt(time.Date(2019, 10, 27, 0, 26, 4, 0, time.UTC))
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/spy_report_17h26-7Z.html")
	e := v6.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, clock.Now(), infos.Date.UTC())
}

func TestExtractEspionageReport_action(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/message_foreign_fleet_sighted.html")
	e := v6.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, ogame.Action, infos.Type)
	assert.Equal(t, int64(6970988), infos.ID)
}

func TestExtractEspionageReport(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/spy_report_res_buildings_researches.html")
	e := v6.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, ogame.Coordinate{4, 212, 6, ogame.PlanetType}, infos.Coordinate)
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
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/spy_report_no_pics.html")
	e := v6.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, err := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, ogame.ErrDeactivateHidePictures, err)
	assert.Equal(t, ogame.Coordinate{4, 203, 6, ogame.PlanetType}, infos.Coordinate)
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
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/spy_report_moon.html")
	e := v6.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, ogame.Coordinate{4, 116, 12, ogame.MoonType}, infos.Coordinate)
	assert.Equal(t, int64(6), *infos.LunarBase)
	assert.Equal(t, int64(4), *infos.SensorPhalanx)
	assert.Nil(t, infos.JumpGate)
}

func TestExtractEspionageReport1(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/spy_report_res_buildings_researches_fleet.html")
	e := v6.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, int64(2), *infos.Battleship)
	assert.Equal(t, int64(1), *infos.Bomber)
}

func TestExtractEspionageReportV7(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7/spy_report.html")
	e := v7.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, int64(15), infos.LastActivity)
	assert.Equal(t, int64(7), *infos.SmallCargo)
}

func TestExtractEspionageReportV71(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.1/en/spy_report.html")
	e := v71.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.False(t, infos.HonorableTarget)
	assert.Equal(t, int64(66331), infos.Metal)
	assert.Equal(t, int64(58452), infos.Crystal)
	assert.Equal(t, int64(0), infos.Deuterium)
	assert.Equal(t, ogame.Collector, infos.CharacterClass)
}

func TestExtractEspionageReportV8(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v8.5/en/spy_report.html")
	e := v8.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, int64(15), infos.LastActivity)
}

func TestExtractEspionageReportAllianceClass(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v8.1/en/spy_report_alliance_class_trader.html")
	e := v71.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, ogame.Trader, infos.AllianceClass)

	pageHTMLBytes, _ = ioutil.ReadFile("../../samples/v8.1/en/spy_report_alliance_class_warrior.html")
	e = v71.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ = e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, ogame.Warrior, infos.AllianceClass)

	//pageHTMLBytes, _ = ioutil.ReadFile("../../samples/v8.1/en/spy_report_alliance_class_researcher.html")
	//infos, _ = NewExtractor().ExtractEspionageReport(pageHTMLBytes, time.FixedZone("OGT", 3600))
	//assert.Equal(t, Researcher, infos.AllianceClass)

	pageHTMLBytes, _ = ioutil.ReadFile("../../samples/v8.1/en/spy_report_alliance_no_class.html")
	e = v71.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ = e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, ogame.NoAllianceClass, infos.AllianceClass)
}

func TestExtractEspionageReportHonorableV71(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.1/en/spy_report_honorable.html")
	e := v71.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.True(t, infos.HonorableTarget)
}

func TestExtractEspionageReportHonorableStrongV71(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.1/en/spy_report_honorable_strong.html")
	e := v71.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.True(t, infos.HonorableTarget)
}

func TestExtractEspionageReportThousands(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/spy_report_thousand_units.html")
	e := v6.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, int64(4000), *infos.RocketLauncher)
	assert.Equal(t, int64(3882), *infos.LargeCargo)
	assert.Equal(t, int64(374), *infos.SolarSatellite)
}

func TestExtractEspionageReport_defence(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/spy_report_res_fleet_defences.html")
	e := v6.NewExtractor()
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
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/spy_report_inactive_bandit_lord.html")
	e := v6.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, true, infos.IsBandit)
	assert.Equal(t, false, infos.IsStarlord)
}

func TestExtractEspionageReport_starlord(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/spy_report_active_star_lord.html")
	e := v6.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, false, infos.IsBandit)
	assert.Equal(t, true, infos.IsStarlord)
}

func TestExtractEspionageReport_norank(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/spy_report_res_buildings_researches_fleet.html")
	e := v6.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, false, infos.IsBandit)
	assert.Equal(t, false, infos.IsStarlord)

}

func TestExtractEspionageReport_username1(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/spy_report_inactive_bandit_lord.html")
	e := v6.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, "Cid Granjeador", infos.Username)
}

func TestExtractEspionageReport_username2(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/spy_report_active_star_lord.html")
	e := v6.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, "Commodore Nomad", infos.Username)
}

func TestExtractEspionageReport_username_outlaw(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/spy_report_outlaw.html")
	e := v6.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, "Constable Telesto", infos.Username)
}

func TestExtractEspionageReport_apiKey(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/spy_report_active_star_lord.html")
	e := v6.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, "sr-en-152-ea0b59302bfad7e3ab0f2d15f7ef2c6a4633b4ba", infos.APIKey)
}

func TestExtractEspionageReport_inactivetimer_within15(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/spy_report_res_buildings.html")
	e := v6.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, int64(15), infos.LastActivity)
}

func TestExtractEspionageReport_inactivetimer_29(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/spy_report_res_buildings_researches.html")
	e := v6.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, int64(29), infos.LastActivity)
}

func TestExtractEspionageReport_inactivetimer_over1h(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/spy_report_inactive_bandit_lord.html")
	e := v6.NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, int64(0), infos.LastActivity)
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

func TestExtractFleetSlot_FleetDispatch_V7(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7/fleetdispatch.html")
	s := v7.NewExtractor().ExtractSlots(pageHTMLBytes)
	assert.Equal(t, int64(0), s.InUse)
	assert.Equal(t, int64(4), s.Total)
	assert.Equal(t, int64(0), s.ExpInUse)
	assert.Equal(t, int64(1), s.ExpTotal)
}

func TestExtractFleetSlotV7_movement(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7/movement.html")
	s := v6.NewExtractor().ExtractSlots(pageHTMLBytes)
	assert.Equal(t, int64(1), s.InUse)
	assert.Equal(t, int64(2), s.Total)
	assert.Equal(t, int64(0), s.ExpInUse)
	assert.Equal(t, int64(1), s.ExpTotal)
}

func TestExtractFleetSlot_fleet1(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/fleet1.html")
	s := v6.NewExtractor().ExtractSlots(pageHTMLBytes)
	assert.Equal(t, int64(2), s.InUse)
	assert.Equal(t, int64(14), s.Total)
	assert.Equal(t, int64(0), s.ExpInUse)
	assert.Equal(t, int64(3), s.ExpTotal)
}

func TestExtractFleetSlot_movement(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/fleets_1.html")
	s := v6.NewExtractor().ExtractSlots(pageHTMLBytes)
	assert.Equal(t, int64(1), s.InUse)
	assert.Equal(t, int64(11), s.Total)
	assert.Equal(t, int64(0), s.ExpInUse)
	assert.Equal(t, int64(2), s.ExpTotal)
}

func TestExtractFleetSlot_commanders(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/fleet1_extract_slots_with_commanders.html")
	s := v6.NewExtractor().ExtractSlots(pageHTMLBytes)
	assert.Equal(t, int64(13), s.InUse)
	assert.Equal(t, int64(14), s.Total)
	assert.Equal(t, int64(2), s.ExpInUse)
	assert.Equal(t, int64(3), s.ExpTotal)
}

func TestGetResourcesDetails(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/unversioned/fetch_resources.html")
	res, _ := v6.NewExtractor().ExtractResourcesDetails(pageHTMLBytes)
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

func TestGetResourcesDetailsV7(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7/fetchResources.html")
	res, _ := v7.NewExtractor().ExtractResourcesDetails(pageHTMLBytes)
	assert.Equal(t, int64(415), res.Metal.Available)
	assert.Equal(t, int64(10000), res.Metal.StorageCapacity)
	assert.Equal(t, int64(150), res.Metal.CurrentProduction)

	assert.Equal(t, int64(501), res.Crystal.Available)
	assert.Equal(t, int64(10000), res.Crystal.StorageCapacity)
	assert.Equal(t, int64(75), res.Crystal.CurrentProduction)

	assert.Equal(t, int64(73), res.Deuterium.Available)
	assert.Equal(t, int64(10000), res.Deuterium.StorageCapacity)
	assert.Equal(t, int64(66), res.Deuterium.CurrentProduction)

	assert.Equal(t, int64(0), res.Energy.Available)
	assert.Equal(t, int64(22), res.Energy.CurrentProduction)
	assert.Equal(t, int64(-22), res.Energy.Consumption)

	assert.Equal(t, int64(8000), res.Darkmatter.Available)
	assert.Equal(t, int64(0), res.Darkmatter.Purchased)
	assert.Equal(t, int64(8000), res.Darkmatter.Found)
}

func TestGetResourcesDetailsV71(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.1/en/fetchResources.html")
	res, _ := v71.NewExtractor().ExtractResourcesDetails(pageHTMLBytes)
	assert.Equal(t, int64(260120), res.Metal.Available)
	assert.Equal(t, int64(470000), res.Metal.StorageCapacity)
	assert.Equal(t, int64(13915), res.Metal.CurrentProduction)

	assert.Equal(t, int64(95684), res.Crystal.Available)
	assert.Equal(t, int64(255000), res.Crystal.StorageCapacity)
	assert.Equal(t, int64(5984), res.Crystal.CurrentProduction)

	assert.Equal(t, int64(140000), res.Deuterium.Available)
	assert.Equal(t, int64(140000), res.Deuterium.StorageCapacity)
	assert.Equal(t, int64(0), res.Deuterium.CurrentProduction)

	assert.Equal(t, int64(-1865), res.Energy.Available)
	assert.Equal(t, int64(2690), res.Energy.CurrentProduction)
	assert.Equal(t, int64(-4555), res.Energy.Consumption)

	assert.Equal(t, int64(8000), res.Darkmatter.Available)
	assert.Equal(t, int64(0), res.Darkmatter.Purchased)
	assert.Equal(t, int64(8000), res.Darkmatter.Found)
}

func TestExtractDestroyRockets(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.6.2/en/destroy_rockets.html")
	abm, ipm, token, _ := v71.NewExtractor().ExtractDestroyRockets(pageHTMLBytes)
	assert.Equal(t, "3a1148bb0d2c6a18f323cf7f0ce09d2b", token)
	assert.Equal(t, int64(24), abm)
	assert.Equal(t, int64(6), ipm)
}

func TestExtractIPMV71(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.1/nl/ipm_missile_launch.html")
	duration, max, token := v71.NewExtractor().ExtractIPM(pageHTMLBytes)
	assert.Equal(t, "95b68270230217f7e9a813e4a4beb20e", token)
	assert.Equal(t, int64(25), max)
	assert.Equal(t, int64(248), duration)
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

func TestExtractEmpirePlanets(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v8.1/en/empire_planets.html")
	res, _ := v6.NewExtractor().ExtractEmpire(pageHTMLBytes)
	assert.Equal(t, 8, len(res))
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 208, Position: 8, Type: ogame.PlanetType}, res[0].Coordinate)
	assert.Equal(t, int64(-3199), res[0].Resources.Energy)
	assert.Equal(t, int64(13904), res[0].Diameter)
}

func TestExtractEmpireMoons(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v8.1/en/empire_moons.html")
	res, _ := v6.NewExtractor().ExtractEmpire(pageHTMLBytes)
	assert.Equal(t, 3, len(res))
	assert.Equal(t, ogame.Coordinate{Galaxy: 4, System: 116, Position: 9, Type: ogame.MoonType}, res[0].Coordinate)
	assert.Equal(t, int64(0), res[0].Resources.Energy)
	assert.Equal(t, int64(-19), res[0].Temperature.Min)
	assert.Equal(t, int64(21), res[0].Temperature.Max)
	assert.Equal(t, int64(5783), res[0].Diameter)
}

func TestExtractAuctionV874(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v8.7.4/en/traderAuctioneer.html")
	res, _ := v874.NewExtractor().ExtractAuction(pageHTMLBytes)
	assert.Equal(t, "43576386810cdf91a833a6239f323f66", res.Token)
}

func TestExtractAuction_playerBid(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.5.0/en/auction_player_bid.html")
	res, _ := v6.NewExtractor().ExtractAuction(pageHTMLBytes)
	assert.Equal(t, int64(1603000), res.AlreadyBid)
}

func TestExtractAuction_noPlayerBid(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.5.0/en/auction_no_player_bid.html")
	res, _ := v6.NewExtractor().ExtractAuction(pageHTMLBytes)
	assert.Equal(t, int64(0), res.AlreadyBid)
}

func TestExtractAuction_ongoing2(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.4/en/traderAuctioneer_ongoing.html")
	res, _ := v6.NewExtractor().ExtractAuction(pageHTMLBytes)
	assert.Equal(t, int64(1800), res.Endtime)
}

func TestExtractAuction_ongoing(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.1/en/traderOverview_ongoing.html")
	res, _ := v6.NewExtractor().ExtractAuction(pageHTMLBytes)
	assert.Equal(t, int64(1200), res.Endtime)
}

func TestExtractAuction_waiting(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.1/en/traderOverview_waiting.html")
	res, _ := v6.NewExtractor().ExtractAuction(pageHTMLBytes)
	assert.Equal(t, int64(6202), res.Endtime)
}

func TestExtractHighscore(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.1/en/highscore.html")
	highscore, _ := v71.NewExtractor().ExtractHighscore(pageHTMLBytes)
	assert.Equal(t, int64(7), highscore.NbPage)
	assert.Equal(t, int64(1), highscore.CurrPage)
	assert.Equal(t, int64(1), highscore.Category)
	assert.Equal(t, int64(0), highscore.Type)
	assert.Equal(t, 100, len(highscore.Players))
	assert.Equal(t, int64(103636), highscore.Players[0].ID)
	assert.Equal(t, int64(525), highscore.Players[0].AllianceID)
	assert.Equal(t, int64(3299957), highscore.Players[0].HonourPoints)
	assert.Equal(t, int64(320933389), highscore.Players[0].Score)
	assert.Equal(t, "blondie", highscore.Players[0].Name)
	assert.Equal(t, ogame.Coordinate{2, 356, 15, ogame.PlanetType}, highscore.Players[0].Homeworld)

	pageHTMLBytes, _ = ioutil.ReadFile("../../samples/v7.1/en/highscore_withSelf.html")
	highscore, _ = v71.NewExtractor().ExtractHighscore(pageHTMLBytes)
	assert.Equal(t, int64(7), highscore.NbPage)
	assert.Equal(t, int64(2), highscore.CurrPage)
	assert.Equal(t, int64(1), highscore.Category)
	assert.Equal(t, int64(0), highscore.Type)
	assert.Equal(t, 100, len(highscore.Players))
	assert.Equal(t, "Bob", highscore.Players[7].Name)
	assert.Equal(t, int64(0), highscore.Players[7].ID)         // Player ID is broken for self
	assert.Equal(t, int64(0), highscore.Players[7].AllianceID) // Alliance ID is broken for self

	pageHTMLBytes, _ = ioutil.ReadFile("../../samples/v7.1/en/highscore_fullPage.html")
	highscore, _ = v71.NewExtractor().ExtractHighscore(pageHTMLBytes)
	assert.Equal(t, int64(7), highscore.NbPage)
	assert.Equal(t, int64(2), highscore.CurrPage)
	assert.Equal(t, int64(1), highscore.Category)
	assert.Equal(t, int64(0), highscore.Type)
	assert.Equal(t, 100, len(highscore.Players))
	assert.Equal(t, "Bob", highscore.Players[7].Name)
	assert.Equal(t, int64(0), highscore.Players[7].ID)         // Player ID is broken for self
	assert.Equal(t, int64(0), highscore.Players[7].AllianceID) // Alliance ID is broken for self

	pageHTMLBytes, _ = ioutil.ReadFile("../../samples/v7.1/en/highscore_withShips.html")
	highscore, _ = v71.NewExtractor().ExtractHighscore(pageHTMLBytes)
	assert.Equal(t, "malakopipis", highscore.Players[0].Name)
	assert.Equal(t, int64(125758), highscore.Players[0].Ships)
}

func TestExtractAllResources(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.1/en/traderOverview_waiting.html")
	resources, _ := v71.NewExtractor().ExtractAllResources(pageHTMLBytes)
	assert.Equal(t, 12, len(resources))
	assert.Equal(t, ogame.Resources{Metal: 97696396, Crystal: 30582087, Deuterium: 32752170}, resources[33698658])
	assert.Equal(t, ogame.Resources{Metal: 133578, Crystal: 74977, Deuterium: 66899}, resources[33702461])
	assert.Equal(t, ogame.Resources{Metal: 0, Crystal: 0, Deuterium: 2676231}, resources[33741598])
}

func TestExtractAllResourcesTwV902(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v9.0.2/tw/traderauctioneer.html")
	resources, _ := v71.NewExtractor().ExtractAllResources(pageHTMLBytes)
	assert.Equal(t, 1, len(resources))
	assert.Equal(t, ogame.Resources{Metal: 1005, Crystal: 1002, Deuterium: 0}, resources[33620229])
}

func TestExtractDMCosts(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.1/en/overview_allDM.html")
	dmCosts, _ := v71.NewExtractor().ExtractDMCosts(pageHTMLBytes)
	assert.Equal(t, ogame.SolarPlantID, dmCosts.Buildings.OGameID)
	assert.Equal(t, int64(30), dmCosts.Buildings.Nbr)
	assert.Equal(t, false, dmCosts.Buildings.Complete)
	assert.Equal(t, true, dmCosts.Buildings.CanBuy)
	assert.Equal(t, int64(5250), dmCosts.Buildings.Cost)
	assert.Equal(t, "cb4fd53e61feced0d52cfc4c1ce383bad9c05f67", dmCosts.Buildings.BuyAndActivateToken)
	assert.Equal(t, "335606c6b7a472c4685ecf74667a9da4", dmCosts.Buildings.Token)
	assert.Equal(t, ogame.EnergyTechnologyID, dmCosts.Research.OGameID)
	assert.Equal(t, int64(13), dmCosts.Research.Nbr)
	assert.Equal(t, false, dmCosts.Research.Complete)
	assert.Equal(t, false, dmCosts.Research.CanBuy)
	assert.Equal(t, int64(70500), dmCosts.Research.Cost)
	assert.Equal(t, "14c17d49462963f5e5b67efa1257622ce1b866ac", dmCosts.Research.BuyAndActivateToken)
	assert.Equal(t, "03983fb66c7cae9caf2c2e8af91c7285", dmCosts.Research.Token)
	assert.Equal(t, ogame.BomberID, dmCosts.Shipyard.OGameID)
	assert.Equal(t, int64(2), dmCosts.Shipyard.Nbr)
	assert.Equal(t, false, dmCosts.Shipyard.Complete)
	assert.Equal(t, true, dmCosts.Shipyard.CanBuy)
	assert.Equal(t, int64(750), dmCosts.Shipyard.Cost)
	assert.Equal(t, "75accaa0d1bc22b78d83b89cd437bdccd6a58887", dmCosts.Shipyard.BuyAndActivateToken)
	assert.Equal(t, "d5303c8e555e21b406976731b3283a26", dmCosts.Shipyard.Token)

	pageHTMLBytes, _ = ioutil.ReadFile("../../samples/v7.1/en/overview_shipyard_queue.html")
	dmCosts, _ = v71.NewExtractor().ExtractDMCosts(pageHTMLBytes)
	assert.Equal(t, ogame.ID(0), dmCosts.Buildings.OGameID)
	assert.Equal(t, int64(0), dmCosts.Buildings.Nbr)
	assert.Equal(t, false, dmCosts.Buildings.Complete)
	assert.Equal(t, false, dmCosts.Buildings.CanBuy)
	assert.Equal(t, int64(0), dmCosts.Buildings.Cost)
	assert.Equal(t, "", dmCosts.Buildings.BuyAndActivateToken)
	assert.Equal(t, "", dmCosts.Buildings.Token)
	assert.Equal(t, ogame.LaserTechnologyID, dmCosts.Research.OGameID)
	assert.Equal(t, int64(12), dmCosts.Research.Nbr)
	assert.Equal(t, false, dmCosts.Research.Complete)
	assert.Equal(t, false, dmCosts.Research.CanBuy)
	assert.Equal(t, int64(9000), dmCosts.Research.Cost)
	assert.Equal(t, "14c17d49462963f5e5b67efa1257622ce1b866ac", dmCosts.Research.BuyAndActivateToken)
	assert.Equal(t, "aa979762eed1b27c05db6fa7d5eb20a6", dmCosts.Research.Token)
	assert.Equal(t, ogame.GaussCannonID, dmCosts.Shipyard.OGameID)
	assert.Equal(t, int64(2), dmCosts.Shipyard.Nbr)
	assert.Equal(t, false, dmCosts.Shipyard.Complete)
	assert.Equal(t, true, dmCosts.Shipyard.CanBuy)
	assert.Equal(t, int64(1500), dmCosts.Shipyard.Cost)
	assert.Equal(t, "75accaa0d1bc22b78d83b89cd437bdccd6a58887", dmCosts.Shipyard.BuyAndActivateToken)
	assert.Equal(t, "f34eddc43aaeb43f9e7c6971c87eea2f", dmCosts.Shipyard.Token)
}

func TestExtractBuffActivation(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.1/en/buffActivation.html")
	token, items, _ := v71.NewExtractor().ExtractBuffActivation(pageHTMLBytes)
	assert.Equal(t, "081876002bf5791011097597836d3f5c", token)
	assert.Equal(t, 31, len(items))
}

func TestExtractOGameSession(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7/overview.html")
	session := v6.NewExtractor().ExtractOGameSession(pageHTMLBytes)
	assert.Equal(t, "0a724276a3ddbe9949f62bdae48d71c1a16adf20", session)

	pageHTMLBytes, _ = ioutil.ReadFile("../../samples/v7/overview_mobile.html")
	session = v6.NewExtractor().ExtractOGameSession(pageHTMLBytes)
	assert.Equal(t, "c1626ce8228ac5986e3808a7d42d4afc764c1b68", session)
}

func TestExtractIsMobile(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.1/en/movement.html")
	isMobile := v71.NewExtractor().ExtractIsMobile(pageHTMLBytes)
	assert.False(t, isMobile)

	pageHTMLBytes, _ = ioutil.ReadFile("../../samples/v7.2/en/movement_mobile.html")
	isMobile = v71.NewExtractor().ExtractIsMobile(pageHTMLBytes)
	assert.True(t, isMobile)
}

func TestExtractActiveItems(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../samples/v7.6.6/en/overview_with_active_items.html")
	items, _ := v71.NewExtractor().ExtractActiveItems(pageHTMLBytes)
	assert.Equal(t, 2, len(items))
	assert.Equal(t, int64(69994), items[0].ID)
	assert.Equal(t, "ba85cc2b8a5d986bbfba6954e2164ef71af95d4a", items[0].Ref)
	assert.Equal(t, "Silver Metal Booster", items[0].Name)
	assert.Equal(t, int64(604800), items[0].TotalDuration)
	assert.Equal(t, int64(579307), items[0].TimeRemaining)
	assert.Equal(t, "https://s152-en.ogame.gameforge.com/cdn/img/item-images/1ab70d0954b4ebbb91e020c60afbaacb28707e5d-small.png", items[0].ImgSmall)

	assert.Equal(t, int64(69995), items[1].ID)
	assert.Equal(t, "5560a1580a0330e8aadf05cb5bfe6bc3200406e2", items[1].Ref)
	assert.Equal(t, "Gold Deuterium Booster", items[1].Name)
	assert.Equal(t, int64(604800), items[1].TotalDuration)
	assert.Equal(t, int64(579827), items[1].TimeRemaining)
	assert.Equal(t, "https://s152-en.ogame.gameforge.com/cdn/img/item-images/db408084e3b2b7b0e1fe13d9f234d2ebd76f11c5-small.png", items[1].ImgSmall)
}

func TestVersion(t *testing.T) {
	assert.False(t, version.Must(version.NewVersion("8.7.4-pl3")).GreaterThanOrEqual(version.Must(version.NewVersion("8.7.4"))))
	assert.True(t, version.Must(version.NewVersion("8.7.4-pl3")).GreaterThanOrEqual(version.Must(version.NewVersion("8.7.4-pl3"))))
	assert.True(t, version.Must(version.NewVersion("8.7.4-pl4")).GreaterThanOrEqual(version.Must(version.NewVersion("8.7.4-pl3"))))
	assert.True(t, version.Must(version.NewVersion("8.7.5-pl3")).GreaterThanOrEqual(version.Must(version.NewVersion("8.7.5-pl3"))))
}
