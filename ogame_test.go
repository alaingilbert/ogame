package ogame

import (
	"bytes"
	"io/ioutil"
	"regexp"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func BenchmarkUserInfoRegex(b *testing.B) {
	extractUserRegex := func(pageHTML []byte) (int, string) {
		playerID := toInt(regexp.MustCompile(`playerId="(\d+)"`).FindSubmatch(pageHTML)[1])
		playerName := string(regexp.MustCompile(`playerName="([^"]+)"`).FindSubmatch(pageHTML)[1])
		return playerID, playerName
	}
	pageHTMLBytes, _ := ioutil.ReadFile("samples/overview_inactive.html")
	for n := 0; n < b.N; n++ {
		extractUserRegex(pageHTMLBytes)
	}
}

func BenchmarkUserInfoGoquery(b *testing.B) {
	extractUserGoquery := func(pageHTML []byte) (int, string) {
		doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
		playerID := ParseInt(doc.Find("meta[name=ogame-player-id]").AttrOr("content", "0"))
		playerName := doc.Find("meta[name=ogame-player-name]").AttrOr("content", "")
		return playerID, playerName
	}
	pageHTMLBytes, _ := ioutil.ReadFile("samples/overview_inactive.html")
	for n := 0; n < b.N; n++ {
		extractUserGoquery(pageHTMLBytes)
	}
}

func TestWrapper(t *testing.T) {
	var bot Wrapper
	bot = NewNoLogin("", "", "", "")
	assert.NotNil(t, bot)
}

func TestParseInt2(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/deathstar_price.html")
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTMLBytes))
	title := doc.Find("li.metal").AttrOr("title", "")
	metalStr := regexp.MustCompile(`([\d.]+)`).FindStringSubmatch(title)[1]
	metal := ParseInt(metalStr)
	assert.Equal(t, 5000000, metal)

	pageHTMLBytes, _ = ioutil.ReadFile("samples/mrd_price.html")
	doc, _ = goquery.NewDocumentFromReader(bytes.NewReader(pageHTMLBytes))
	title = doc.Find("li.metal").AttrOr("title", "")
	metalStr = regexp.MustCompile(`([\d.]+)`).FindStringSubmatch(title)[1]
	metal = ParseInt(metalStr)
	assert.Equal(t, 1555733200, metal)
}

func TestExtractFleetDeutSaveFactor_V6_2_2_1(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/overview_active.html")
	res := ExtractFleetDeutSaveFactor(pageHTMLBytes)
	assert.Equal(t, 1.0, res)
}

func TestExtractFleetDeutSaveFactor_V6_7_4(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/overview_with_moon.html")
	res := ExtractFleetDeutSaveFactor(pageHTMLBytes)
	assert.Equal(t, 0.5, res)
}

func TestExtractPlanetCoordinate(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/station.html")
	res, _ := ExtractPlanetCoordinate(pageHTMLBytes)
	assert.Equal(t, Coordinate{1, 301, 5, PlanetType}, res)
}

func TestExtractPlanetCoordinate_moon(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/moon_facilities.html")
	res, _ := ExtractPlanetCoordinate(pageHTMLBytes)
	assert.Equal(t, Coordinate{4, 116, 12, MoonType}, res)
}

func TestExtractPlanetID_planet(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/station.html")
	res, _ := ExtractPlanetID(pageHTMLBytes)
	assert.Equal(t, CelestialID(33672410), res)
}

func TestExtractPlanetID_moon(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/moon_facilities.html")
	res, _ := ExtractPlanetID(pageHTMLBytes)
	assert.Equal(t, CelestialID(33741598), res)
}

func TestExtractPlanetType_planet(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/station.html")
	res, _ := ExtractPlanetType(pageHTMLBytes)
	assert.Equal(t, PlanetType, res)
}

func TestExtractPlanetType_moon(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/moon_facilities.html")
	res, _ := ExtractPlanetType(pageHTMLBytes)
	assert.Equal(t, MoonType, res)
}

func TestExtractJumpGate_cooldown(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/jumpgatelayer_charge.html")
	_, _, _, wait := extractJumpGate(pageHTMLBytes)
	assert.Equal(t, 1730, wait)
}

func TestExtractJumpGate(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/jumpgatelayer.html")
	ships, token, dests, wait := extractJumpGate(pageHTMLBytes)
	assert.Equal(t, 1, len(dests))
	assert.Equal(t, MoonID(33743183), dests[0])
	assert.Equal(t, 0, wait)
	assert.Equal(t, "7787b530670bc89623b5d65a827e557a", token)
	assert.Equal(t, 0, ships.SmallCargo)
	assert.Equal(t, 101, ships.LargeCargo)
	assert.Equal(t, 1, ships.LightFighter)
}

func TestExtractOgameTimestamp(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/moon_facilities.html")
	res := ExtractOgameTimestamp(pageHTMLBytes)
	assert.Equal(t, int64(1538912592), res)
}

func TestExtractOgameTimestampFromBytes(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/moon_facilities.html")
	res := ExtractOGameTimestampFromBytes(pageHTMLBytes)
	assert.Equal(t, int64(1538912592), res)
}

func TestExtractResources(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/moon_facilities.html")
	res := ExtractResources(pageHTMLBytes)
	assert.Equal(t, 280000, res.Metal)
	assert.Equal(t, 260000, res.Crystal)
	assert.Equal(t, 280000, res.Deuterium)
	assert.Equal(t, 0, res.Energy)
	assert.Equal(t, 25000, res.Darkmatter)
}

func TestExtractResourcesDetailsFromFullPage(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/fleets_1.html")
	res := ExtractResourcesDetailsFromFullPage(pageHTMLBytes)
	assert.Equal(t, 1959227, res.Metal.Available)
	assert.Equal(t, 37818, res.Metal.CurrentProduction)
	assert.Equal(t, 5355000, res.Metal.StorageCapacity)
	assert.Equal(t, 327916, res.Crystal.Available)
	assert.Equal(t, 21862, res.Crystal.CurrentProduction)
	assert.Equal(t, 865000, res.Crystal.StorageCapacity)
	assert.Equal(t, 618155, res.Deuterium.Available)
	assert.Equal(t, 7508, res.Deuterium.CurrentProduction)
	assert.Equal(t, 865000, res.Deuterium.StorageCapacity)
	assert.Equal(t, 220, res.Energy.Available)
	assert.Equal(t, 17597, res.Energy.CurrentProduction)
	assert.Equal(t, -17377, res.Energy.Consumption)
	assert.Equal(t, 25000, res.Darkmatter.Available)
	assert.Equal(t, 0, res.Darkmatter.Purchased)
	assert.Equal(t, 25000, res.Darkmatter.Found)
}

func TestExtractPhalanx(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/phalanx.html")
	res, err := extractPhalanx(pageHTMLBytes)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(res))
	assert.Equal(t, MissionID(3), res[0].Mission)
	assert.Equal(t, true, res[0].ReturnFlight)
	assert.NotNil(t, res[0].ArriveIn)
	assert.Equal(t, Coordinate{4, 116, 9, PlanetType}, res[0].Origin)
	assert.Equal(t, Coordinate{4, 212, 8, PlanetType}, res[0].Destination)
	assert.Equal(t, 100, res[0].Ships.LargeCargo)
}

func TestExtractPhalanx_fromMoon(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/phalanx_from_moon.html")
	res, _ := extractPhalanx(pageHTMLBytes)
	assert.Equal(t, Coordinate{4, 116, 12, MoonType}, res[0].Origin)
	assert.Equal(t, Coordinate{4, 116, 9, PlanetType}, res[0].Destination)
}

func TestExtractPhalanx_manyFleets(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/phalanx_fleets.html")
	res, err := extractPhalanx(pageHTMLBytes)
	assert.Nil(t, err)
	assert.Equal(t, 12, len(res))
	assert.Equal(t, Expedition, res[0].Mission)
	assert.False(t, res[0].ReturnFlight)
	assert.Equal(t, Coordinate{4, 124, 9, PlanetType}, res[0].Origin)
	assert.Equal(t, Coordinate{4, 125, 16, PlanetType}, res[0].Destination)
	assert.Equal(t, 250, res[0].Ships.LargeCargo)
	assert.Equal(t, 1, res[0].Ships.EspionageProbe)
	assert.Equal(t, 1, res[0].Ships.Destroyer)

	assert.Equal(t, Expedition, res[8].Mission)
	assert.True(t, res[8].ReturnFlight)
	assert.Equal(t, Coordinate{4, 124, 9, PlanetType}, res[8].Origin)
	assert.Equal(t, Coordinate{4, 125, 16, PlanetType}, res[8].Destination)
	assert.Equal(t, 250, res[8].Ships.LargeCargo)
	assert.Equal(t, 1, res[8].Ships.EspionageProbe)
	assert.Equal(t, 1, res[8].Ships.Destroyer)
}

func TestExtractPhalanx_noFleet(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/phalanx_no_fleet.html")
	res, err := extractPhalanx(pageHTMLBytes)
	assert.Equal(t, 0, len(res))
	assert.Nil(t, err)
}

func TestExtractPhalanx_noDeut(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/phalanx_no_deut.html")
	res, err := extractPhalanx(pageHTMLBytes)
	assert.Equal(t, 0, len(res))
	assert.NotNil(t, err)
}

func TestExtractResearch(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/research_bonus.html")
	res := ExtractResearch(pageHTMLBytes)
	assert.Equal(t, 12, res.EnergyTechnology)
	assert.Equal(t, 12, res.LaserTechnology)
	assert.Equal(t, 7, res.IonTechnology)
	assert.Equal(t, 6, res.HyperspaceTechnology)
	assert.Equal(t, 7, res.PlasmaTechnology)
	assert.Equal(t, 15, res.CombustionDrive)
	assert.Equal(t, 7, res.ImpulseDrive)
	assert.Equal(t, 8, res.HyperspaceDrive)
	assert.Equal(t, 10, res.EspionageTechnology)
	assert.Equal(t, 14, res.ComputerTechnology)
	assert.Equal(t, 13, res.Astrophysics)
	assert.Equal(t, 0, res.IntergalacticResearchNetwork)
	assert.Equal(t, 0, res.GravitonTechnology)
	assert.Equal(t, 13, res.WeaponsTechnology)
	assert.Equal(t, 12, res.ShieldingTechnology)
	assert.Equal(t, 12, res.ArmourTechnology)
}

func TestExtractResourcesBuildings(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/resource_inconstruction.html")
	res, _ := ExtractResourcesBuildings(pageHTMLBytes)
	assert.Equal(t, 19, res.MetalMine)
	assert.Equal(t, 17, res.CrystalMine)
	assert.Equal(t, 13, res.DeuteriumSynthesizer)
	assert.Equal(t, 20, res.SolarPlant)
	assert.Equal(t, 3, res.FusionReactor)
	assert.Equal(t, 0, res.SolarSatellite)
	assert.Equal(t, 5, res.MetalStorage)
	assert.Equal(t, 4, res.CrystalStorage)
	assert.Equal(t, 3, res.DeuteriumTank)
}

func TestExtractFacilities(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/facility_inconstruction.html")
	res, _ := ExtractFacilities(pageHTMLBytes)
	assert.Equal(t, 7, res.RoboticsFactory)
	assert.Equal(t, 7, res.Shipyard)
	assert.Equal(t, 7, res.ResearchLab)
	assert.Equal(t, 0, res.AllianceDepot)
	assert.Equal(t, 0, res.MissileSilo)
	assert.Equal(t, 0, res.NaniteFactory)
	assert.Equal(t, 0, res.Terraformer)
	assert.Equal(t, 3, res.SpaceDock)
}

func TestExtractMoonFacilities(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/moon_facilities.html")
	res, _ := ExtractFacilities(pageHTMLBytes)
	assert.Equal(t, 1, res.RoboticsFactory)
	assert.Equal(t, 2, res.Shipyard)
	assert.Equal(t, 3, res.LunarBase)
	assert.Equal(t, 4, res.SensorPhalanx)
	assert.Equal(t, 5, res.JumpGate)
}

func TestExtractDefense(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/defence.html")
	defense, _ := ExtractDefense(pageHTMLBytes)
	assert.Equal(t, 1, defense.RocketLauncher)
	assert.Equal(t, 2, defense.LightLaser)
	assert.Equal(t, 3, defense.HeavyLaser)
	assert.Equal(t, 4, defense.GaussCannon)
	assert.Equal(t, 5, defense.IonCannon)
	assert.Equal(t, 6, defense.PlasmaTurret)
	assert.Equal(t, 0, defense.SmallShieldDome)
	assert.Equal(t, 0, defense.LargeShieldDome)
	assert.Equal(t, 7, defense.AntiBallisticMissiles)
	assert.Equal(t, 8, defense.InterplanetaryMissiles)
}

func TestProductionRatio(t *testing.T) {
	ratio := productionRatio(
		Temperature{-23, 17},
		ResourcesBuildings{MetalMine: 29, CrystalMine: 16, DeuteriumSynthesizer: 26, SolarPlant: 29, FusionReactor: 13, SolarSatellite: 51},
		ResourceSettings{MetalMine: 100, CrystalMine: 100, DeuteriumSynthesizer: 100, SolarPlant: 100, FusionReactor: 100, SolarSatellite: 100},
		12,
	)
	assert.Equal(t, 1.0, ratio)
}

func TestEnergyNeeded(t *testing.T) {
	needed := energyNeeded(
		ResourcesBuildings{MetalMine: 29, CrystalMine: 16, DeuteriumSynthesizer: 26},
		ResourceSettings{MetalMine: 100, CrystalMine: 100, DeuteriumSynthesizer: 100},
	)
	assert.Equal(t, 4601+736+6198, needed)
}

func TestEnergyProduced(t *testing.T) {
	produced := energyProduced(
		Temperature{-23, 17},
		ResourcesBuildings{SolarPlant: 29, FusionReactor: 13, SolarSatellite: 51},
		ResourceSettings{SolarPlant: 100, FusionReactor: 100, SolarSatellite: 100},
		12,
	)
	assert.Equal(t, 9200+3002+1326, produced)
}

func TestExtractFleet1Ships(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/fleet1.html")
	s := ExtractFleet1Ships(pageHTMLBytes)
	assert.Equal(t, 3, s.LightFighter)
	assert.Equal(t, 0, s.HeavyFighter)
	assert.Equal(t, 1012, s.Cruiser)
	assert.Equal(t, 0, s.Battleship)
	assert.Equal(t, 0, s.SmallCargo)
	assert.Equal(t, 1003, s.LargeCargo)
	assert.Equal(t, 1, s.ColonyShip)
	assert.Equal(t, 200, s.Battlecruiser)
	assert.Equal(t, 100, s.Bomber)
	assert.Equal(t, 200, s.Destroyer)
	assert.Equal(t, 0, s.Deathstar)
	assert.Equal(t, 30, s.Recycler)
	assert.Equal(t, 1001, s.EspionageProbe)
}

func TestExtractFleet1Ships_NoShips(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/fleet1_no_ships.html")
	s := ExtractFleet1Ships(pageHTMLBytes)
	assert.Equal(t, 0, s.LightFighter)
	assert.Equal(t, 0, s.HeavyFighter)
	assert.Equal(t, 0, s.Cruiser)
	assert.Equal(t, 0, s.Battleship)
	assert.Equal(t, 0, s.SmallCargo)
	assert.Equal(t, 0, s.LargeCargo)
	assert.Equal(t, 0, s.ColonyShip)
	assert.Equal(t, 0, s.Battlecruiser)
	assert.Equal(t, 0, s.Bomber)
	assert.Equal(t, 0, s.Destroyer)
	assert.Equal(t, 0, s.Deathstar)
	assert.Equal(t, 0, s.Recycler)
	assert.Equal(t, 0, s.EspionageProbe)
}

func TestExtractPlanet_en(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/overview_queues.html")
	planet, _ := ExtractPlanet(pageHTMLBytes, PlanetID(33677371), &OGame{language: "en"})
	assert.Equal(t, "C1", planet.Name)
	assert.Equal(t, 14615, planet.Diameter)
	assert.Equal(t, -2, planet.Temperature.Min)
	assert.Equal(t, 38, planet.Temperature.Max)
	assert.Equal(t, 35, planet.Fields.Built)
	assert.Equal(t, 238, planet.Fields.Total)
	assert.Equal(t, PlanetID(33677371), planet.ID)
	assert.Equal(t, Coordinate{1, 301, 8, PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_fr(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/fr_overview.html")
	planet, _ := ExtractPlanet(pageHTMLBytes, PlanetID(33629512), &OGame{language: "fr"})
	assert.Equal(t, "planète mère", planet.Name)
	assert.Equal(t, 12800, planet.Diameter)
	assert.Equal(t, 48, planet.Temperature.Min)
	assert.Equal(t, 88, planet.Temperature.Max)
	assert.Equal(t, 0, planet.Fields.Built)
	assert.Equal(t, 188, planet.Fields.Total)
	assert.Equal(t, PlanetID(33629512), planet.ID)
	assert.Equal(t, Coordinate{2, 180, 4, PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_de(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/de_overview.html")
	planet, _ := ExtractPlanet(pageHTMLBytes, PlanetID(33630447), &OGame{language: "de"})
	assert.Equal(t, "Heimatplanet", planet.Name)
	assert.Equal(t, 12800, planet.Diameter)
	assert.Equal(t, 21, planet.Temperature.Min)
	assert.Equal(t, 61, planet.Temperature.Max)
	assert.Equal(t, 0, planet.Fields.Built)
	assert.Equal(t, 188, planet.Fields.Total)
	assert.Equal(t, PlanetID(33630447), planet.ID)
	assert.Equal(t, Coordinate{2, 175, 8, PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_dk(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/dk_overview.html")
	planet, _ := ExtractPlanet(pageHTMLBytes, PlanetID(33627426), &OGame{language: "dk"})
	assert.Equal(t, "Hjemme verden", planet.Name)
	assert.Equal(t, 12800, planet.Diameter)
	assert.Equal(t, -23, planet.Temperature.Min)
	assert.Equal(t, 17, planet.Temperature.Max)
	assert.Equal(t, 5, planet.Fields.Built)
	assert.Equal(t, 193, planet.Fields.Total)
	assert.Equal(t, PlanetID(33627426), planet.ID)
	assert.Equal(t, Coordinate{1, 148, 12, PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_es(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/es/shipyard.html")
	planet, _ := ExtractPlanet(pageHTMLBytes, PlanetID(33644981), &OGame{language: "es"})
	assert.Equal(t, "Planeta Principal", planet.Name)
	assert.Equal(t, 12800, planet.Diameter)
	assert.Equal(t, -18, planet.Temperature.Min)
	assert.Equal(t, 22, planet.Temperature.Max)
	assert.Equal(t, 0, planet.Fields.Built)
	assert.Equal(t, 193, planet.Fields.Total)
	assert.Equal(t, PlanetID(33644981), planet.ID)
	assert.Equal(t, Coordinate{2, 493, 10, PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_br(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/br/overview.html")
	planet, _ := ExtractPlanet(pageHTMLBytes, PlanetID(33633767), &OGame{language: "br"})
	assert.Equal(t, "Planeta Principal", planet.Name)
	assert.Equal(t, 12800, planet.Diameter)
	assert.Equal(t, -13, planet.Temperature.Min)
	assert.Equal(t, 27, planet.Temperature.Max)
	assert.Equal(t, 5, planet.Fields.Built)
	assert.Equal(t, 193, planet.Fields.Total)
	assert.Equal(t, PlanetID(33633767), planet.ID)
	assert.Equal(t, Coordinate{1, 449, 12, PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_it(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/it/overview.html")
	planet, _ := ExtractPlanet(pageHTMLBytes, PlanetID(33634944), &OGame{language: "it"})
	assert.Equal(t, "Pianeta Madre", planet.Name)
	assert.Equal(t, 12800, planet.Diameter)
	assert.Equal(t, 28, planet.Temperature.Min)
	assert.Equal(t, 68, planet.Temperature.Max)
	assert.Equal(t, 0, planet.Fields.Built)
	assert.Equal(t, 188, planet.Fields.Total)
	assert.Equal(t, PlanetID(33634944), planet.ID)
	assert.Equal(t, Coordinate{2, 58, 8, PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_jp(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/jp_overview.html")
	planet, _ := ExtractPlanet(pageHTMLBytes, PlanetID(33620484), &OGame{language: "jp"})
	assert.Equal(t, "ホームワールド", planet.Name)
	assert.Equal(t, 12800, planet.Diameter)
	assert.Equal(t, 40, planet.Temperature.Min)
	assert.Equal(t, 80, planet.Temperature.Max)
	assert.Equal(t, 0, planet.Fields.Built)
	assert.Equal(t, 188, planet.Fields.Total)
	assert.Equal(t, PlanetID(33620484), planet.ID)
	assert.Equal(t, Coordinate{1, 18, 4, PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_tw(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/tw/overview.html")
	planet, _ := ExtractPlanet(pageHTMLBytes, PlanetID(33626432), &OGame{language: "tw"})
	assert.Equal(t, "母星", planet.Name)
	assert.Equal(t, 12800, planet.Diameter)
	assert.Equal(t, 29, planet.Temperature.Min)
	assert.Equal(t, 69, planet.Temperature.Max)
	assert.Equal(t, 0, planet.Fields.Built)
	assert.Equal(t, 188, planet.Fields.Total)
	assert.Equal(t, PlanetID(33626432), planet.ID)
	assert.Equal(t, Coordinate{1, 206, 8, PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_hr(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/hr/overview.html")
	planet, _ := ExtractPlanet(pageHTMLBytes, PlanetID(33627961), &OGame{language: "hr"})
	assert.Equal(t, "Glavni Planet", planet.Name)
	assert.Equal(t, 12800, planet.Diameter)
	assert.Equal(t, -33, planet.Temperature.Min)
	assert.Equal(t, 7, planet.Temperature.Max)
	assert.Equal(t, 4, planet.Fields.Built)
	assert.Equal(t, 193, planet.Fields.Total)
	assert.Equal(t, PlanetID(33627961), planet.ID)
	assert.Equal(t, Coordinate{1, 236, 12, PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_mx(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/mx/overview.html")
	planet, _ := ExtractPlanet(pageHTMLBytes, PlanetID(33624669), &OGame{language: "mx"})
	assert.Equal(t, "Planeta Principal", planet.Name)
	assert.Equal(t, 12800, planet.Diameter)
	assert.Equal(t, 33, planet.Temperature.Min)
	assert.Equal(t, 73, planet.Temperature.Max)
	assert.Equal(t, 0, planet.Fields.Built)
	assert.Equal(t, 188, planet.Fields.Total)
	assert.Equal(t, PlanetID(33624669), planet.ID)
	assert.Equal(t, Coordinate{1, 390, 6, PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_cz(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/cz/overview.html")
	planet, _ := ExtractPlanet(pageHTMLBytes, PlanetID(33622822), &OGame{language: "cz"})
	assert.Equal(t, "Domovska planeta", planet.Name)
	assert.Equal(t, 12800, planet.Diameter)
	assert.Equal(t, -13, planet.Temperature.Min)
	assert.Equal(t, 27, planet.Temperature.Max)
	assert.Equal(t, 0, planet.Fields.Built)
	assert.Equal(t, 188, planet.Fields.Total)
	assert.Equal(t, PlanetID(33622822), planet.ID)
	assert.Equal(t, Coordinate{1, 221, 12, PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_jp1(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/jp/overview.html")
	planet, _ := ExtractPlanet(pageHTMLBytes, PlanetID(33623513), &OGame{language: "jp"})
	assert.Equal(t, "ホームワールド", planet.Name)
	assert.Equal(t, 12800, planet.Diameter)
	assert.Equal(t, 58, planet.Temperature.Min)
	assert.Equal(t, 98, planet.Temperature.Max)
	assert.Equal(t, 0, planet.Fields.Built)
	assert.Equal(t, 188, planet.Fields.Total)
	assert.Equal(t, PlanetID(33623513), planet.ID)
	assert.Equal(t, Coordinate{1, 85, 4, PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_pl(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/pl_overview.html")
	planet, _ := ExtractPlanet(pageHTMLBytes, PlanetID(33669699), &OGame{language: "pl"})
	assert.Equal(t, "Planeta matka", planet.Name)
	assert.Equal(t, 12800, planet.Diameter)
	assert.Equal(t, -21, planet.Temperature.Min)
	assert.Equal(t, 19, planet.Temperature.Max)
	assert.Equal(t, 0, planet.Fields.Built)
	assert.Equal(t, 188, planet.Fields.Total)
	assert.Equal(t, PlanetID(33669699), planet.ID)
	assert.Equal(t, Coordinate{4, 248, 12, PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_tr(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/tr_overview.html")
	planet, _ := ExtractPlanet(pageHTMLBytes, PlanetID(33650421), &OGame{language: "tr"})
	assert.Equal(t, "Ana Gezegen", planet.Name)
	assert.Equal(t, 12800, planet.Diameter)
	assert.Equal(t, 9, planet.Temperature.Min)
	assert.Equal(t, 49, planet.Temperature.Max)
	assert.Equal(t, 0, planet.Fields.Built)
	assert.Equal(t, 188, planet.Fields.Total)
	assert.Equal(t, PlanetID(33650421), planet.ID)
	assert.Equal(t, Coordinate{3, 143, 10, PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_pt(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/pt_overview.html")
	planet, _ := ExtractPlanet(pageHTMLBytes, PlanetID(33635398), &OGame{language: "pt"})
	assert.Equal(t, "Planeta Principal", planet.Name)
	assert.Equal(t, 12800, planet.Diameter)
	assert.Equal(t, 33, planet.Temperature.Min)
	assert.Equal(t, 73, planet.Temperature.Max)
	assert.Equal(t, 4, planet.Fields.Built)
	assert.Equal(t, 193, planet.Fields.Total)
	assert.Equal(t, PlanetID(33635398), planet.ID)
	assert.Equal(t, Coordinate{2, 241, 6, PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_nl(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/nl_overview.html")
	planet, _ := ExtractPlanet(pageHTMLBytes, PlanetID(33624684), &OGame{language: "nl"})
	assert.Equal(t, "Hoofdplaneet", planet.Name)
	assert.Equal(t, 12800, planet.Diameter)
	assert.Equal(t, -47, planet.Temperature.Min)
	assert.Equal(t, -7, planet.Temperature.Max)
	assert.Equal(t, 5, planet.Fields.Built)
	assert.Equal(t, 188, planet.Fields.Total)
	assert.Equal(t, PlanetID(33624684), planet.ID)
	assert.Equal(t, Coordinate{1, 178, 12, PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_ar(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/ar/overview.html")
	planet, _ := ExtractPlanet(pageHTMLBytes, PlanetID(33629527), &OGame{language: "ar"})
	assert.Equal(t, "Planeta Principal", planet.Name)
	assert.Equal(t, 12800, planet.Diameter)
	assert.Equal(t, 37, planet.Temperature.Min)
	assert.Equal(t, 77, planet.Temperature.Max)
	assert.Equal(t, 0, planet.Fields.Built)
	assert.Equal(t, 193, planet.Fields.Total)
	assert.Equal(t, PlanetID(33629527), planet.ID)
	assert.Equal(t, Coordinate{1, 367, 4, PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_ru(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/ru/overview.html")
	planet, _ := ExtractPlanet(pageHTMLBytes, PlanetID(33629521), &OGame{language: "ru"})
	assert.Equal(t, "Главная планета", planet.Name)
	assert.Equal(t, 12800, planet.Diameter)
	assert.Equal(t, 23, planet.Temperature.Min)
	assert.Equal(t, 63, planet.Temperature.Max)
	assert.Equal(t, 0, planet.Fields.Built)
	assert.Equal(t, 163, planet.Fields.Total)
	assert.Equal(t, PlanetID(33629521), planet.ID)
	assert.Equal(t, Coordinate{1, 374, 6, PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestExtractPlanet_notExists(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/overview_queues.html")
	_, err := ExtractPlanet(pageHTMLBytes, PlanetID(12345), &OGame{language: "en"})
	assert.NotNil(t, err)
}

func TestExtractPlanetByCoord(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/overview_queues.html")
	planet, _ := ExtractPlanetByCoord(pageHTMLBytes, &OGame{language: "en"}, Coordinate{1, 301, 8, PlanetType})
	assert.Equal(t, "C1", planet.Name)
	assert.Equal(t, 14615, planet.Diameter)
}

func TestExtractPlanetByCoord_notExists(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/overview_queues.html")
	_, err := ExtractPlanetByCoord(pageHTMLBytes, &OGame{language: "en"}, Coordinate{1, 2, 3, PlanetType})
	assert.NotNil(t, err)
}

func TestFindSlowestSpeed(t *testing.T) {
	assert.Equal(t, 8000, findSlowestSpeed(ShipsInfos{SmallCargo: 1, LargeCargo: 1}, Researches{CombustionDrive: 6}))
}

func TestExtractShips(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/shipyard_thousands_ships.html")
	ships, _ := ExtractShips(pageHTMLBytes)
	assert.Equal(t, 1000, ships.LargeCargo)
	assert.Equal(t, 1000, ships.EspionageProbe)
	assert.Equal(t, 700, ships.Cruiser)
}

func TestExtractShipsMillions(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/shipyard_millions_ships.html")
	ships, _ := ExtractShips(pageHTMLBytes)
	assert.Equal(t, 15000001, ships.LightFighter)
}

func TestExtractShipsWhileBeingBuilt(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/shipyard_ship_being_built.html")
	ships, _ := ExtractShips(pageHTMLBytes)
	assert.Equal(t, 213, ships.EspionageProbe)
}

func TestExtractEspionageReportMessageIDs(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/messages.html")
	msgs, _ := extractEspionageReportMessageIDs(pageHTMLBytes)
	assert.Equal(t, 2, len(msgs))
	assert.Equal(t, Report, msgs[0].Type)
	assert.Equal(t, Coordinate{4, 117, 6, PlanetType}, msgs[0].Target)
	assert.Equal(t, 0.5, msgs[0].LootPercentage)
	assert.Equal(t, "Fleet Command", msgs[0].From)
	assert.Equal(t, Action, msgs[1].Type)
	assert.Equal(t, "Space Monitoring", msgs[1].From)
	assert.Equal(t, Coordinate{4, 117, 9, PlanetType}, msgs[1].Target)
}

func TestExtractEspionageReportMessageIDsLootPercentage(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/messages_loot_percentage.html")
	msgs, _ := extractEspionageReportMessageIDs(pageHTMLBytes)
	assert.Equal(t, 1.0, msgs[0].LootPercentage)
	assert.Equal(t, 0.5, msgs[1].LootPercentage)
	assert.Equal(t, 0.5, msgs[2].LootPercentage)
}

func TestExtractCombatReportMessages(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/combat_reports_msgs.html")
	msgs, _ := extractCombatReportMessagesSummary(pageHTMLBytes)
	assert.Equal(t, 9, len(msgs))
}

func TestExtractCombatReportAttackingMessages(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/combat_reports_msgs_attacking.html")
	msgs, _ := extractCombatReportMessagesSummary(pageHTMLBytes)
	assert.Equal(t, 7945368, msgs[0].ID)
	assert.Equal(t, Coordinate{4, 233, 11, PlanetType}, msgs[0].Destination)
	assert.Equal(t, 50, msgs[0].Loot)
	assert.Equal(t, 74495, msgs[0].Metal)
	assert.Equal(t, 88280, msgs[0].Crystal)
	assert.Equal(t, 21572, msgs[0].Deuterium)
	assert.Equal(t, "08.09.2018 09:33:18", msgs[0].CreatedAt.Format("02.01.2006 15:04:05"))
}

func TestExtractCombatReportMessagesSummary(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/combat_reports_msgs_2.html")
	msgs, nbPages := extractCombatReportMessagesSummary(pageHTMLBytes)
	assert.Equal(t, 10, len(msgs))
	assert.Equal(t, 44, nbPages)
	assert.Equal(t, Coordinate{4, 116, 12, MoonType}, msgs[1].Destination)
	assert.Equal(t, Coordinate{4, 127, 9, MoonType}, *msgs[1].Origin)
}

func TestExtractResourcesProductions(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/resource_settings.html")
	prods, _ := ExtractResourcesProductions(pageHTMLBytes)
	assert.Equal(t, Resources{Metal: 10352, Crystal: 5104, Deuterium: 1282, Energy: -52}, prods)
}

func TestExtractNbProbes(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/preferences.html")
	probes := ExtractSpioAnz(pageHTMLBytes)
	assert.Equal(t, 10, probes)

	pageHTMLBytes, _ = ioutil.ReadFile("samples/preferences_mobile.html")
	probes = ExtractSpioAnz(pageHTMLBytes)
	assert.Equal(t, 10, probes)
}

func TestExtractPreferencesShowActivityMinutes(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/preferences.html")
	checked := ExtractPreferencesShowActivityMinutes(pageHTMLBytes)
	assert.True(t, checked)

	pageHTMLBytes, _ = ioutil.ReadFile("samples/preferences_mobile.html")
	checked = ExtractPreferencesShowActivityMinutes(pageHTMLBytes)
	assert.True(t, checked)

	pageHTMLBytes, _ = ioutil.ReadFile("samples/preferences_without_detailed_activities.html")
	checked = ExtractPreferencesShowActivityMinutes(pageHTMLBytes)
	assert.False(t, checked)
}

func TestExtractPreferences(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/preferences.html")
	prefs := ExtractPreferences(pageHTMLBytes)
	assert.Equal(t, 10, prefs.SpioAnz)
	assert.False(t, prefs.DisableChatBar)
	assert.False(t, prefs.DisableOutlawWarning)
	assert.False(t, prefs.MobileVersion)
	assert.False(t, prefs.ShowOldDropDowns)
	assert.False(t, prefs.ActivateAutofocus)
	assert.Equal(t, 1, prefs.EventsShow)
	assert.Equal(t, 0, prefs.SortSetting)
	assert.Equal(t, 0, prefs.SortOrder)
	assert.True(t, prefs.ShowDetailOverlay)
	assert.True(t, prefs.AnimatedSliders)
	assert.True(t, prefs.AnimatedOverview)
	assert.False(t, prefs.PopupsNotices)
	assert.False(t, prefs.PopopsCombatreport)
	assert.False(t, prefs.SpioReportPictures)
	assert.Equal(t, 10, prefs.MsgResultsPerPage)
	assert.True(t, prefs.AuctioneerNotifications)
	assert.False(t, prefs.EconomyNotifications)
	assert.True(t, prefs.ShowActivityMinutes)
	assert.False(t, prefs.PreserveSystemOnPlanetChange)

	pageHTMLBytes, _ = ioutil.ReadFile("samples/preferences_reverse.html")
	prefs = ExtractPreferences(pageHTMLBytes)
	assert.Equal(t, 2, prefs.SpioAnz)
	assert.True(t, prefs.DisableChatBar)
	assert.True(t, prefs.DisableOutlawWarning)
	assert.False(t, prefs.MobileVersion)
	assert.True(t, prefs.ShowOldDropDowns)
	assert.True(t, prefs.ActivateAutofocus)
	assert.Equal(t, 3, prefs.EventsShow)
	assert.Equal(t, 3, prefs.SortSetting)
	assert.Equal(t, 1, prefs.SortOrder)
	assert.False(t, prefs.ShowDetailOverlay)
	assert.False(t, prefs.AnimatedSliders)
	assert.False(t, prefs.AnimatedOverview)
	assert.True(t, prefs.PopupsNotices)
	assert.True(t, prefs.PopopsCombatreport)
	assert.True(t, prefs.SpioReportPictures)
	assert.Equal(t, 50, prefs.MsgResultsPerPage)
	assert.False(t, prefs.AuctioneerNotifications)
	assert.True(t, prefs.EconomyNotifications)
	assert.False(t, prefs.ShowActivityMinutes)
	assert.True(t, prefs.PreserveSystemOnPlanetChange)

	pageHTMLBytes, _ = ioutil.ReadFile("samples/preferences_mobile.html")
	prefs = ExtractPreferences(pageHTMLBytes)
	assert.Equal(t, 3, prefs.SpioAnz)
	assert.False(t, prefs.DisableChatBar) // no mobile
	assert.False(t, prefs.DisableOutlawWarning)
	assert.True(t, prefs.MobileVersion)
	assert.False(t, prefs.ShowOldDropDowns)
	assert.False(t, prefs.ActivateAutofocus)
	assert.Equal(t, 2, prefs.EventsShow)
	assert.Equal(t, 0, prefs.SortSetting)
	assert.Equal(t, 0, prefs.SortOrder)
	assert.True(t, prefs.ShowDetailOverlay)
	assert.False(t, prefs.AnimatedSliders)    // no mobile
	assert.False(t, prefs.AnimatedOverview)   // no mobile
	assert.False(t, prefs.PopupsNotices)      // no mobile
	assert.False(t, prefs.PopopsCombatreport) // no mobile
	assert.False(t, prefs.SpioReportPictures)
	assert.Equal(t, 10, prefs.MsgResultsPerPage)
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

	pageHTMLBytes, _ = ioutil.ReadFile("samples/preferences_reverse_mobile.html")
	prefs = ExtractPreferences(pageHTMLBytes)
	assert.Equal(t, 2, prefs.SpioAnz)
	assert.False(t, prefs.DisableChatBar) // no mobile
	assert.True(t, prefs.DisableOutlawWarning)
	assert.True(t, prefs.MobileVersion)
	assert.True(t, prefs.ShowOldDropDowns)
	assert.True(t, prefs.ActivateAutofocus)
	assert.Equal(t, 3, prefs.EventsShow)
	assert.Equal(t, 3, prefs.SortSetting)
	assert.Equal(t, 1, prefs.SortOrder)
	assert.False(t, prefs.ShowDetailOverlay)
	assert.False(t, prefs.AnimatedSliders)    // no mobile
	assert.False(t, prefs.AnimatedOverview)   // no mobile
	assert.False(t, prefs.PopupsNotices)      // no mobile
	assert.False(t, prefs.PopopsCombatreport) // no mobile
	assert.True(t, prefs.SpioReportPictures)
	assert.Equal(t, 50, prefs.MsgResultsPerPage)
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
//	pageHTMLBytes, _ := ioutil.ReadFile("samples/traderOverview.html")
//	price, _, planetResources, multiplier, _ := ExtractOfferOfTheDay(pageHTMLBytes)
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
	pageHTMLBytes, _ := ioutil.ReadFile("samples/traderOverview.html")
	price, token, _, _, _ := ExtractOfferOfTheDay(pageHTMLBytes)
	assert.Equal(t, 54243, price)
	assert.Equal(t, "8128c0ba0c9981599a87d818003f95e1", token)
}

func TestExtractAttacks(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/event_list_attack.html")
	attacks, _ := ExtractAttacks(pageHTMLBytes)
	assert.Equal(t, 1, len(attacks))
}

func TestExtractAttacksPhoneDisplay(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/event_list_attack_phone.html")
	attacks, _ := ExtractAttacks(pageHTMLBytes)
	assert.Equal(t, 1, len(attacks))
}

func TestExtractAttacksMeAttacking(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/eventlist_me_attacking.html")
	attacks, _ := ExtractAttacks(pageHTMLBytes)
	assert.Equal(t, 0, len(attacks))
}

func TestExtractAttacksWithoutShips(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/event_list_attack.html")
	attacks, _ := ExtractAttacks(pageHTMLBytes)
	assert.Equal(t, 1, len(attacks))
	assert.Equal(t, 100771, attacks[0].AttackerID)
	assert.Equal(t, 0, attacks[0].Missiles)
	assert.Nil(t, attacks[0].Ships)
}

func TestExtractAttacksWithShips(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/eventList_attack_ships.html")
	attacks, _ := ExtractAttacks(pageHTMLBytes)
	assert.Equal(t, 1, len(attacks))
	assert.NotNil(t, attacks[0].Ships)
	assert.Equal(t, PlanetType, attacks[0].Destination.Type)
	assert.Equal(t, 197, attacks[0].Ships.LargeCargo)
	assert.Equal(t, 3, attacks[0].Ships.LightFighter)
	assert.Equal(t, 8, attacks[0].Ships.HeavyFighter)
	assert.Equal(t, 92, attacks[0].Ships.Cruiser)
	assert.Equal(t, 571, attacks[0].Ships.EspionageProbe)
	assert.Equal(t, 27, attacks[0].Ships.Bomber)
	assert.Equal(t, 4, attacks[0].Ships.Destroyer)
	assert.Equal(t, 11, attacks[0].Ships.Battlecruiser)
}

func TestExtractAttacksMoon(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/eventlist_moon_attacked.html")
	attacks, _ := ExtractAttacks(pageHTMLBytes)
	assert.Equal(t, 1, len(attacks))
	assert.NotNil(t, attacks[0].Ships)
	assert.Equal(t, 107009, attacks[0].AttackerID)
	assert.Equal(t, Coordinate{4, 212, 8, PlanetType}, attacks[0].Origin)
	assert.Equal(t, Coordinate{4, 116, 12, MoonType}, attacks[0].Destination)
	assert.Equal(t, MoonType, attacks[0].Destination.Type)
	assert.Equal(t, 1, attacks[0].Ships.SmallCargo)
}

func TestExtractAttacksMoonDestruction(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/eventlist_moon_destruction.html")
	attacks, _ := ExtractAttacks(pageHTMLBytes)
	assert.Equal(t, 1, len(attacks))
	assert.NotNil(t, attacks[0].Ships)
	assert.Equal(t, 106734, attacks[0].AttackerID)
	assert.Equal(t, Coordinate{4, 116, 12, PlanetType}, attacks[0].Origin)
	assert.Equal(t, Coordinate{4, 116, 9, MoonType}, attacks[0].Destination)
	assert.Equal(t, MoonType, attacks[0].Destination.Type)
	assert.Equal(t, 1, attacks[0].Ships.Deathstar)
}

func TestExtractAttacksWithThousandsShips(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/eventlist_attack_thousands.html")
	attacks, _ := ExtractAttacks(pageHTMLBytes)
	assert.Equal(t, 2, len(attacks))
	assert.Equal(t, 1012, attacks[1].Ships.Cruiser)
	assert.Equal(t, 1000, attacks[1].Ships.LargeCargo)
}

func TestExtractAttacksUnknownShips(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/eventlist_unknown_ships_nbr.html")
	attacks, _ := ExtractAttacks(pageHTMLBytes)
	assert.Equal(t, 1, len(attacks))
	assert.Equal(t, -1, attacks[0].Ships.Cruiser)
	assert.Equal(t, 0, attacks[0].Ships.Destroyer)
}

func TestExtractAttacks_spy(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/event_list_spy.html")
	attacks, _ := ExtractAttacks(pageHTMLBytes)
	assert.Equal(t, 1, len(attacks))
	assert.Equal(t, Coordinate{4, 212, 8, PlanetType}, attacks[0].Origin)
	assert.Equal(t, 107009, attacks[0].AttackerID)
}

func TestExtractAttacks1(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/event_list_missile.html")
	attacks, _ := ExtractAttacks(pageHTMLBytes)
	assert.Equal(t, 1, len(attacks))
	assert.Equal(t, 1, attacks[0].Missiles)
	assert.Nil(t, attacks[0].Ships)
}

func TestExtractCargoCapacity(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/sendfleet3.htm")
	fleet3Doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTMLBytes))
	cargo := ParseInt(fleet3Doc.Find("#maxresources").Text())
	assert.Equal(t, 442500, cargo)
}

func TestExtractGalaxyInfos_bandit(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/galaxy_inactive_bandit_lord.html")
	infos, _ := ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.True(t, infos.Position(6).Player.IsBandit)
	assert.False(t, infos.Position(6).Player.IsStarlord)
}

func TestExtractGalaxyInfos_starlord(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/galaxy_inactive_emperor.html")
	infos, _ := ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.True(t, infos.Position(7).Player.IsStarlord)
	assert.False(t, infos.Position(7).Player.IsBandit)
}

func TestExtractGalaxyInfos_destroyedPlanet(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/galaxy_destroyed_planet.html")
	infos, _ := ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Nil(t, infos.Position(8))
}

func TestExtractGalaxyInfos_destroyedPlanetAndMoon(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/galaxy_destroyed_planet_and_moon2.html")
	infos, _ := ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Nil(t, infos.Position(15))
}

func TestExtractGalaxyInfos_banned(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/galaxy_banned.html")
	infos, _ := ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, true, infos.Position(1).Banned)
	assert.Equal(t, false, infos.Position(9).Banned)
}

func TestExtractGalaxyInfos(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/galaxy_ajax.html")
	infos, _ := ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, 4, infos.Galaxy())
	assert.Equal(t, 116, infos.System())
	assert.Equal(t, 33698600, infos.Position(4).ID)
	assert.Equal(t, 33698645, infos.Position(6).ID)
	assert.Equal(t, 106733, infos.Position(6).Player.ID)
	assert.Equal(t, "Origin", infos.Position(6).Player.Name)
	assert.Equal(t, 1671, infos.Position(6).Player.Rank)
	assert.Equal(t, "Ra", infos.Position(6).Name)
}

func TestExtractGalaxyInfosOwnPlanet(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/galaxy_ajax.html")
	infos, _ := ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, 33698658, infos.Position(12).ID)
	assert.Equal(t, "Commodore Nomade", infos.Position(12).Player.Name)
	assert.Equal(t, 123, infos.Position(12).Player.ID)
	assert.Equal(t, 456, infos.Position(12).Player.Rank)
	assert.Equal(t, "Homeworld", infos.Position(12).Name)
}

func TestExtractGalaxyInfosPlanetNoActivity(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/galaxy_planet_activity.html")
	infos, _ := ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, 0, infos.Position(15).Activity)
}

func TestExtractGalaxyInfosPlanetActivity15(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/galaxy_planet_activity.html")
	infos, _ := ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, 15, infos.Position(8).Activity)
}

func TestExtractGalaxyInfosPlanetActivity23(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/galaxy_planet_activity.html")
	infos, _ := ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, 23, infos.Position(9).Activity)
}

//func TestExtractGalaxyInfosPlanetActivityWithoutDetailedActivity(t *testing.T) {
//	pageHTMLBytes, _ := ioutil.ReadFile("samples/galaxy_planet_activity_without_detailed_activity.html")
//	infos, _ := ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
//	assert.Equal(t, 49, infos.Position(5).Activity)
//}

func TestExtractGalaxyInfosMoonActivity(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/galaxy_moon_activity.html")
	infos, _ := ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, 33732827, infos.Position(3).Moon.ID)
	assert.Equal(t, 5830, infos.Position(3).Moon.Diameter)
	assert.Equal(t, 18, infos.Position(3).Moon.Activity)
}

func TestExtractGalaxyInfosMoonNoActivity(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/galaxy_moon_no_activity.html")
	infos, _ := ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, 33650476, infos.Position(2).Moon.ID)
	assert.Equal(t, 7897, infos.Position(2).Moon.Diameter)
	assert.Equal(t, 0, infos.Position(2).Moon.Activity)
}

func TestExtractGalaxyInfosMoonActivity15(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/galaxy_moon_activity_unprecise.html")
	infos, _ := ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, 0, infos.Position(11).Activity)
	assert.Equal(t, 33730993, infos.Position(11).Moon.ID)
	assert.Equal(t, 8944, infos.Position(11).Moon.Diameter)
	assert.Equal(t, 15, infos.Position(11).Moon.Activity)
}

func TestExtractUserInfos(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/overview_inactive.html")
	infos, _ := ExtractUserInfos(pageHTMLBytes, "en")
	assert.Equal(t, 1295, infos.Points)
}

func TestExtractUserInfos_hr(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/hr/overview.html")
	infos, _ := ExtractUserInfos(pageHTMLBytes, "hr")
	assert.Equal(t, 0, infos.Points)
	assert.Equal(t, 214, infos.Rank)
	assert.Equal(t, 252, infos.Total)
}

func TestExtractUserInfos_mx(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/mx/overview.html")
	infos, _ := ExtractUserInfos(pageHTMLBytes, "mx")
	assert.Equal(t, 0, infos.Points)
	assert.Equal(t, 916, infos.Rank)
	assert.Equal(t, 917, infos.Total)
}

func TestExtractUserInfos_de(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/de_overview.html")
	infos, _ := ExtractUserInfos(pageHTMLBytes, "de")
	assert.Equal(t, 0, infos.Points)
	assert.Equal(t, 2980, infos.Rank)
	assert.Equal(t, 2980, infos.Total)
}

func TestExtractUserInfos_dk(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/dk_overview.html")
	infos, _ := ExtractUserInfos(pageHTMLBytes, "dk")
	assert.Equal(t, 0, infos.Points)
	assert.Equal(t, 253, infos.Rank)
	assert.Equal(t, 254, infos.Total)
	assert.Equal(t, "Procurator Zibal", infos.PlayerName)
}

func TestExtractUserInfos_jp(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/jp_overview.html")
	infos, _ := ExtractUserInfos(pageHTMLBytes, "jp")
	assert.Equal(t, 0, infos.Points)
	assert.Equal(t, 73, infos.Rank)
	assert.Equal(t, 73, infos.Total)
}

func TestExtractUserInfos_jp1(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/jp/overview.html")
	infos, _ := ExtractUserInfos(pageHTMLBytes, "jp")
	assert.Equal(t, 0, infos.Points)
	assert.Equal(t, 85, infos.Rank)
	assert.Equal(t, 86, infos.Total)
}

func TestExtractUserInfos_cz(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/cz/overview.html")
	infos, _ := ExtractUserInfos(pageHTMLBytes, "cz")
	assert.Equal(t, 0, infos.Points)
	assert.Equal(t, 1008, infos.Rank)
	assert.Equal(t, 1009, infos.Total)
}

func TestExtractUserInfos_fr(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/fr_overview.html")
	infos, _ := ExtractUserInfos(pageHTMLBytes, "fr")
	assert.Equal(t, 0, infos.Points)
	assert.Equal(t, 3197, infos.Rank)
	assert.Equal(t, 3348, infos.Total)
	assert.Equal(t, "Bandit Pégasus", infos.PlayerName)
}

func TestExtractUserInfos_nl(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/nl_overview.html")
	infos, _ := ExtractUserInfos(pageHTMLBytes, "nl")
	assert.Equal(t, 0, infos.Points)
	assert.Equal(t, 482, infos.Rank)
	assert.Equal(t, 542, infos.Total)
	assert.Equal(t, "Bandit Japetus", infos.PlayerName)
}

func TestExtractUserInfos_pl(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/pl_overview.html")
	infos, _ := ExtractUserInfos(pageHTMLBytes, "pl")
	assert.Equal(t, 0, infos.Points)
	assert.Equal(t, 5873, infos.Rank)
	assert.Equal(t, 5876, infos.Total)
	assert.Equal(t, "Constable Leonis", infos.PlayerName)
}

func TestExtractUserInfos_br(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/br/overview.html")
	infos, _ := ExtractUserInfos(pageHTMLBytes, "br")
	assert.Equal(t, 0, infos.Points)
	assert.Equal(t, 1026, infos.Rank)
	assert.Equal(t, 1268, infos.Total)
}

func TestExtractUserInfos_tr(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/tr_overview.html")
	infos, _ := ExtractUserInfos(pageHTMLBytes, "tr")
	assert.Equal(t, 0, infos.Points)
	assert.Equal(t, 3655, infos.Rank)
	assert.Equal(t, 3656, infos.Total)
	assert.Equal(t, "Chief Apus", infos.PlayerName)
}

func TestExtractUserInfos_ar(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/ar/overview.html")
	infos, _ := ExtractUserInfos(pageHTMLBytes, "ar")
	assert.Equal(t, 0, infos.Points)
	assert.Equal(t, 1158, infos.Rank)
	assert.Equal(t, 1159, infos.Total)
	assert.Equal(t, "Chief Lambda", infos.PlayerName)
}

func TestExtractUserInfos_it(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/it/overview.html")
	infos, _ := ExtractUserInfos(pageHTMLBytes, "it")
	assert.Equal(t, 0, infos.Points)
	assert.Equal(t, 1776, infos.Rank)
	assert.Equal(t, 1777, infos.Total)
	assert.Equal(t, "President Fidis", infos.PlayerName)
}

func TestExtractUserInfos_pt(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/pt_overview.html")
	infos, _ := ExtractUserInfos(pageHTMLBytes, "pt")
	assert.Equal(t, 0, infos.Points)
	assert.Equal(t, 1762, infos.Rank)
	assert.Equal(t, 1862, infos.Total)
	assert.Equal(t, "Director Europa", infos.PlayerName)
}

func TestExtractUserInfos_ru(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/ru/overview.html")
	infos, _ := ExtractUserInfos(pageHTMLBytes, "ru")
	assert.Equal(t, 0, infos.Points)
	assert.Equal(t, 1067, infos.Rank)
	assert.Equal(t, 1068, infos.Total)
	assert.Equal(t, "Viceregent Horizon", infos.PlayerName)
}

func TestExtractMoons(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/overview_with_moon.html")
	moons := ExtractMoons(pageHTMLBytes, &OGame{language: "en"})
	assert.Equal(t, 1, len(moons))
}

func TestExtractMoons2(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/overview_with_many_moon.html")
	moons := ExtractMoons(pageHTMLBytes, &OGame{language: "en"})
	assert.Equal(t, 2, len(moons))
}

func TestExtractMoon_exists(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/overview_with_moon.html")
	_, err := ExtractMoon(pageHTMLBytes, &OGame{language: "en"}, MoonID(33741598))
	assert.Nil(t, err)
}

func TestExtractMoon_notExists(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/overview_with_moon.html")
	_, err := ExtractMoon(pageHTMLBytes, &OGame{language: "en"}, MoonID(12345))
	assert.NotNil(t, err)
}

func TestExtractMoonByCoord_exists(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/overview_with_moon.html")
	_, err := ExtractMoonByCoord(pageHTMLBytes, &OGame{language: "en"}, Coordinate{4, 116, 12, MoonType})
	assert.Nil(t, err)
}

func TestExtractMoonByCoord_notExists(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/overview_with_moon.html")
	_, err := ExtractMoonByCoord(pageHTMLBytes, &OGame{language: "en"}, Coordinate{1, 2, 3, PlanetType})
	assert.NotNil(t, err)
}

func TestExtractIsInVacationFromDoc(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/es/overview_vacation.html")
	assert.True(t, ExtractIsInVacation(pageHTMLBytes))
	pageHTMLBytes, _ = ioutil.ReadFile("samples/es/fleet1_vacation.html")
	assert.True(t, ExtractIsInVacation(pageHTMLBytes))
	pageHTMLBytes, _ = ioutil.ReadFile("samples/es/shipyard.html")
	assert.False(t, ExtractIsInVacation(pageHTMLBytes))
}

func TestExtractPlanetsMoon(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/overview_with_moon.html")
	planets := ExtractPlanets(pageHTMLBytes, nil)
	assert.Equal(t, MoonID(33741598), planets[0].Moon.ID)
	assert.Equal(t, "Moon", planets[0].Moon.Name)
	assert.Equal(t, "https://gf1.geo.gfsrv.net/cdn9d/8e0e6034049bd64e18a1804b42f179.gif", planets[0].Moon.Img)
	assert.Equal(t, 8774, planets[0].Moon.Diameter)
	assert.Equal(t, Coordinate{4, 116, 12, MoonType}, planets[0].Moon.Coordinate)
	assert.Equal(t, 0, planets[0].Moon.Fields.Built)
	assert.Equal(t, 1, planets[0].Moon.Fields.Total)
	assert.Nil(t, planets[1].Moon)
}

func TestExtractPlanets_fieldsFilled(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/overview_fields_filled.html")
	planets := ExtractPlanets(pageHTMLBytes, nil)
	assert.Equal(t, 5, len(planets))
	assert.Equal(t, PlanetID(33698658), planets[0].ID)
	assert.Equal(t, Coordinate{Galaxy: 4, System: 116, Position: 12, Type: PlanetType}, planets[0].Coordinate)
	assert.Equal(t, "Homeworld", planets[0].Name)
	assert.Equal(t, "https://gf3.geo.gfsrv.net/cdnea/7d7ba402d90247ef7d89aa1035e525.png", planets[0].Img)
	assert.Equal(t, -23, planets[0].Temperature.Min)
	assert.Equal(t, 17, planets[0].Temperature.Max)
	assert.Equal(t, 188, planets[0].Fields.Built)
	assert.Equal(t, 188, planets[0].Fields.Total)
	assert.Equal(t, 12800, planets[0].Diameter)
}

func TestExtractPlanets(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/overview_inactive.html")
	planets := ExtractPlanets(pageHTMLBytes, nil)
	assert.Equal(t, 1, len(planets))
	assert.Equal(t, PlanetID(33672410), planets[0].ID)
	assert.Equal(t, Coordinate{Galaxy: 1, System: 301, Position: 5, Type: PlanetType}, planets[0].Coordinate)
	assert.Equal(t, "Homeworld", planets[0].Name)
	assert.Equal(t, "https://gf2.geo.gfsrv.net/cdn46/9f84a481c0c9a83d3b000d801d9d9d.png", planets[0].Img)
	assert.Equal(t, 31, planets[0].Temperature.Min)
	assert.Equal(t, 71, planets[0].Temperature.Max)
	assert.Equal(t, 89, planets[0].Fields.Built)
	assert.Equal(t, 188, planets[0].Fields.Total)
	assert.Equal(t, 12800, planets[0].Diameter)
}

func TestExtractPlanets_es(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/overview_es.html")
	planets := ExtractPlanets(pageHTMLBytes, &OGame{language: "es"})
	assert.Equal(t, 1, len(planets))
	assert.Equal(t, PlanetID(33630486), planets[0].ID)
	assert.Equal(t, Coordinate{Galaxy: 2, System: 147, Position: 8, Type: PlanetType}, planets[0].Coordinate)
	assert.Equal(t, "Planeta Principal", planets[0].Name)
	assert.Equal(t, "https://gf2.geo.gfsrv.net/cdnd1/83579badf7c16d217b06afda455cfe.png", planets[0].Img)
	assert.Equal(t, 18, planets[0].Temperature.Min)
	assert.Equal(t, 58, planets[0].Temperature.Max)
	assert.Equal(t, 0, planets[0].Fields.Built)
	assert.Equal(t, 193, planets[0].Fields.Total)
	assert.Equal(t, 12800, planets[0].Diameter)
}

func TestExtractPlanets_fr(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/fr_overview.html")
	planets := ExtractPlanets(pageHTMLBytes, &OGame{language: "fr"})
	assert.Equal(t, 1, len(planets))
	assert.Equal(t, PlanetID(33629512), planets[0].ID)
	assert.Equal(t, Coordinate{Galaxy: 2, System: 180, Position: 4, Type: PlanetType}, planets[0].Coordinate)
	assert.Equal(t, "planète mère", planets[0].Name)
	assert.Equal(t, "https://gf1.geo.gfsrv.net/cdn35/9545f984bcd53c816a1a8452356d00.png", planets[0].Img)
	assert.Equal(t, 48, planets[0].Temperature.Min)
	assert.Equal(t, 88, planets[0].Temperature.Max)
	assert.Equal(t, 0, planets[0].Fields.Built)
	assert.Equal(t, 188, planets[0].Fields.Total)
	assert.Equal(t, 12800, planets[0].Diameter)
}

func TestExtractPlanets_br(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/br/overview.html")
	planets := ExtractPlanets(pageHTMLBytes, &OGame{language: "br"})
	assert.Equal(t, 1, len(planets))
	assert.Equal(t, PlanetID(33633767), planets[0].ID)
	assert.Equal(t, Coordinate{Galaxy: 1, System: 449, Position: 12, Type: PlanetType}, planets[0].Coordinate)
	assert.Equal(t, "Planeta Principal", planets[0].Name)
	assert.Equal(t, "https://gf3.geo.gfsrv.net/cdne8/41d05740ce1a534f5ec77feb11f100.png", planets[0].Img)
	assert.Equal(t, -13, planets[0].Temperature.Min)
	assert.Equal(t, 27, planets[0].Temperature.Max)
	assert.Equal(t, 5, planets[0].Fields.Built)
	assert.Equal(t, 193, planets[0].Fields.Total)
	assert.Equal(t, 12800, planets[0].Diameter)
}

func TestExtractGalaxyInfos_honorableTarget(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/galaxy_debris.html")
	infos, _ := ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.False(t, infos.Position(6).HonorableTarget)
	assert.True(t, infos.Position(8).HonorableTarget)
	assert.False(t, infos.Position(9).HonorableTarget)
}

func TestExtractGalaxyInfos_inactive(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/galaxy_debris.html")
	infos, _ := ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.True(t, infos.Position(6).Inactive)
	assert.False(t, infos.Position(8).Inactive)
	assert.False(t, infos.Position(9).Inactive)
}

func TestExtractGalaxyInfos_strongPlayer(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/galaxy_debris.html")
	infos, _ := ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.False(t, infos.Position(6).StrongPlayer)
	assert.True(t, infos.Position(8).StrongPlayer)
	assert.False(t, infos.Position(9).StrongPlayer)
}

func TestExtractGalaxyInfos_newbie(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/galaxy_newbie.html")
	infos, _ := ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.True(t, infos.Position(4).Newbie)
}

func TestExtractGalaxyInfos_moon(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/galaxy_debris.html")
	infos, _ := ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.NotNil(t, infos.Position(6).Moon)
	assert.Equal(t, 33701543, infos.Position(6).Moon.ID)
	assert.Equal(t, 8366, infos.Position(6).Moon.Diameter)
}

func TestExtractGalaxyInfos_debris(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/galaxy_debris.html")
	infos, _ := ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, 0, infos.Position(6).Debris.Metal)
	assert.Equal(t, 700, infos.Position(6).Debris.Crystal)
	assert.Equal(t, 1, infos.Position(6).Debris.RecyclersNeeded)
}

func TestExtractGalaxyInfos_debris_es(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/galaxy_debris_es.html")
	infos, _ := ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, 0, infos.Position(12).Debris.Metal)
	assert.Equal(t, 128000, infos.Position(12).Debris.Crystal)
	assert.Equal(t, 7, infos.Position(12).Debris.RecyclersNeeded)
}

func TestExtractGalaxyInfos_debris_fr(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/galaxy_debris_fr.html")
	infos, _ := ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, 100, infos.Position(7).Debris.Metal)
	assert.Equal(t, 600, infos.Position(7).Debris.Crystal)
	assert.Equal(t, 1, infos.Position(7).Debris.RecyclersNeeded)
}

func TestExtractGalaxyInfos_debris_de(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/galaxy_debris_de.html")
	infos, _ := ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, 100, infos.Position(9).Debris.Metal)
	assert.Equal(t, 2500, infos.Position(9).Debris.Crystal)
	assert.Equal(t, 1, infos.Position(9).Debris.RecyclersNeeded)
}

func TestExtractGalaxyInfos_vacation(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/galaxy_ajax.html")
	infos, _ := ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.False(t, infos.Position(4).Vacation)
	assert.True(t, infos.Position(6).Vacation)
	assert.True(t, infos.Position(8).Vacation)
	assert.False(t, infos.Position(10).Vacation)
	assert.False(t, infos.Position(12).Vacation)
}

func TestExtractGalaxyInfos_alliance(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/galaxy_ajax.html")
	infos, _ := ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, 303, infos.Position(10).Alliance.ID)
	assert.Equal(t, "Qrvix", infos.Position(10).Alliance.Name)
	assert.Equal(t, 27, infos.Position(10).Alliance.Rank)
	assert.Equal(t, 16, infos.Position(10).Alliance.Member)
}

func TestExtractGalaxyInfos_alliance_fr(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/fr/galaxy.html")
	infos, _ := ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, 635, infos.Position(5).Alliance.ID)
	assert.Equal(t, "leretour", infos.Position(5).Alliance.Name)
	assert.Equal(t, 24, infos.Position(5).Alliance.Rank)
	assert.Equal(t, 11, infos.Position(5).Alliance.Member)
}

func TestExtractGalaxyInfos_alliance_es(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/es/galaxy.html")
	infos, _ := ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.Equal(t, 500053, infos.Position(4).Alliance.ID)
	assert.Equal(t, "Los Aliens Grises", infos.Position(4).Alliance.Name)
	assert.Equal(t, 8, infos.Position(4).Alliance.Rank)
	assert.Equal(t, 70, infos.Position(4).Alliance.Member)
}

func TestUniverseSpeed(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/techtree_universe_speed.html")
	universeSpeed := extractUniverseSpeed(pageHTMLBytes)
	assert.Equal(t, 7, universeSpeed)
}

func TestCancel(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/overview_active_queue2.html")
	token, techID, listID, _ := extractCancelBuildingInfos(pageHTMLBytes)
	assert.Equal(t, "fef7488e4809150cd16e3fa8fa14db37", token)
	assert.Equal(t, 4, techID)
	assert.Equal(t, 2099434, listID)
}

func TestCancelResearch(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/overview_active_queue2.html")
	token, techID, listID, _ := extractCancelResearchInfos(pageHTMLBytes)
	assert.Equal(t, "fff7488e4809150cd16e3fa8fa14db37", token)
	assert.Equal(t, 120, techID)
	assert.Equal(t, 1769925, listID)
}

func TestGetConstructions(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/overview_active.html")
	buildingID, buildingCountdown, researchID, researchCountdown := ExtractConstructions(pageHTMLBytes)
	assert.Equal(t, CrystalMineID, buildingID)
	assert.Equal(t, 731, buildingCountdown)
	assert.Equal(t, CombustionDriveID, researchID)
	assert.Equal(t, 927, researchCountdown)
}

func TestExtractFleetsFromEventList(t *testing.T) {
	//pageHTMLBytes, _ := ioutil.ReadFile("samples/eventlist_test.html")
	//fleets := ExtractFleetsFromEventList(pageHTMLBytes)
	//assert.Equal(t, 4, len(fleets))
}

func TestExtractIPM(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/missileattacklayer.html")
	duration, max, token := ExtractIPM(pageHTMLBytes)
	assert.Equal(t, "26a08f4cc0c0b513e1e8c10d49c14a27", token)
	assert.Equal(t, 17, max)
	assert.Equal(t, 15, duration)
}

func TestExtractFleet(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/fleets_1.html")
	fleets := ExtractFleets(pageHTMLBytes)
	assert.Equal(t, 1, len(fleets))
	assert.Equal(t, 4134, fleets[0].ArriveIn)
	assert.Equal(t, 8277, fleets[0].BackIn)
	assert.Equal(t, Coordinate{4, 116, 12, PlanetType}, fleets[0].Origin)
	assert.Equal(t, Coordinate{4, 117, 9, PlanetType}, fleets[0].Destination)
	assert.Equal(t, Transport, fleets[0].Mission)
	assert.Equal(t, false, fleets[0].ReturnFlight)
	assert.Equal(t, FleetID(4494950), fleets[0].ID)
	assert.Equal(t, 1, fleets[0].Ships.SmallCargo)
	assert.Equal(t, 8, fleets[0].Ships.LargeCargo)
	assert.Equal(t, 1, fleets[0].Ships.LightFighter)
	assert.Equal(t, 1, fleets[0].Ships.ColonyShip)
	assert.Equal(t, 1, fleets[0].Ships.EspionageProbe)
	assert.Equal(t, Resources{Metal: 123, Crystal: 456, Deuterium: 789}, fleets[0].Resources)
}

func TestExtractFleet_expedition(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/fleets_expedition.html")
	fleets := ExtractFleets(pageHTMLBytes)
	assert.Equal(t, 2, len(fleets))
	assert.Equal(t, 2, fleets[1].Ships.LargeCargo)
	assert.Equal(t, Expedition, fleets[1].Mission)
	assert.False(t, fleets[1].ReturnFlight)
	assert.Equal(t, Coordinate{4, 116, 12, MoonType}, fleets[1].Origin)
	assert.Equal(t, Coordinate{4, 116, 16, PlanetType}, fleets[1].Destination)
}

func TestExtractFleet_harvest(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/fleets_harvest.html")
	fleets := ExtractFleets(pageHTMLBytes)
	assert.Equal(t, Coordinate{4, 116, 12, PlanetType}, fleets[5].Origin)
	assert.Equal(t, Coordinate{4, 116, 9, DebrisType}, fleets[5].Destination)
}

func TestExtractFleet_returningTransport(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/fleets_2.html")
	fleets := ExtractFleets(pageHTMLBytes)
	assert.Equal(t, 1, len(fleets))
	assert.Equal(t, -1, fleets[0].ArriveIn)
	assert.Equal(t, 36, fleets[0].BackIn)
}

func TestExtractFleet_deployment(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/fleets_moon_to_moon.html")
	fleets := ExtractFleets(pageHTMLBytes)
	assert.Equal(t, 210, fleets[0].ArriveIn)
	assert.Equal(t, 426, fleets[0].BackIn)
}

func TestExtractFleetThousands(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/fleets_thousands.html")
	fleets := ExtractFleets(pageHTMLBytes)
	assert.Equal(t, Transport, fleets[0].Mission)
	assert.Equal(t, 210, fleets[0].Ships.LargeCargo)
	assert.Equal(t, Resources{Metal: 207862, Crystal: 78903, Deuterium: 42956}, fleets[0].Resources)
}

func TestExtractFleet_returning(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/fleets_2.html")
	fleets := ExtractFleets(pageHTMLBytes)
	assert.Equal(t, 1, len(fleets))
	assert.Equal(t, Coordinate{4, 116, 12, PlanetType}, fleets[0].Origin)
	assert.Equal(t, Coordinate{4, 117, 9, PlanetType}, fleets[0].Destination)
	assert.Equal(t, Transport, fleets[0].Mission)
	assert.Equal(t, true, fleets[0].ReturnFlight)
	assert.Equal(t, FleetID(0), fleets[0].ID)
	assert.Equal(t, 1, fleets[0].Ships.SmallCargo)
	assert.Equal(t, 8, fleets[0].Ships.LargeCargo)
	assert.Equal(t, 1, fleets[0].Ships.LightFighter)
	assert.Equal(t, 1, fleets[0].Ships.ColonyShip)
	assert.Equal(t, 1, fleets[0].Ships.EspionageProbe)
	assert.Equal(t, Resources{Metal: 123, Crystal: 456, Deuterium: 789}, fleets[0].Resources)
}

func TestExtractOverviewProduction(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/overview_shipyard_queue_full.html")
	prods, _ := ExtractOverviewProduction(pageHTMLBytes)
	assert.Equal(t, 6, len(prods))
	assert.Equal(t, HeavyFighterID, prods[0].ID)
	assert.Equal(t, 1, prods[0].Nbr)
	assert.Equal(t, HeavyFighterID, prods[1].ID)
	assert.Equal(t, 2, prods[1].Nbr)
	assert.Equal(t, HeavyFighterID, prods[2].ID)
	assert.Equal(t, 3, prods[2].Nbr)
	assert.Equal(t, HeavyFighterID, prods[3].ID)
	assert.Equal(t, 4, prods[3].Nbr)
	assert.Equal(t, HeavyFighterID, prods[4].ID)
	assert.Equal(t, 5, prods[4].Nbr)
	assert.Equal(t, HeavyFighterID, prods[5].ID)
	assert.Equal(t, 6, prods[5].Nbr)
}

func TestExtractProduction(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/shipyard_queue.html")
	prods, _ := ExtractProduction(pageHTMLBytes)
	assert.Equal(t, 20, len(prods))
	assert.Equal(t, LargeCargoID, prods[0].ID)
	assert.Equal(t, 4, prods[0].Nbr)
}

func TestExtractProduction2(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/shipyard_queue2.html")
	prods, _ := ExtractProduction(pageHTMLBytes)
	assert.Equal(t, BattlecruiserID, prods[0].ID)
	assert.Equal(t, 18, prods[0].Nbr)
	assert.Equal(t, PlasmaTurretID, prods[1].ID)
	assert.Equal(t, 8, prods[1].Nbr)
	assert.Equal(t, RocketLauncherID, prods[2].ID)
	assert.Equal(t, 1000, prods[2].Nbr)
	assert.Equal(t, LightFighterID, prods[10].ID)
	assert.Equal(t, 1, prods[10].Nbr)
}

func TestExtractProductionWithABM(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/production_with_abm.html")
	prods, _ := ExtractProduction(pageHTMLBytes)
	assert.Equal(t, 4, len(prods))
	assert.Equal(t, DeathstarID, prods[0].ID)
	assert.Equal(t, 1, prods[0].Nbr)
	assert.Equal(t, AntiBallisticMissilesID, prods[1].ID)
	assert.Equal(t, 1, prods[1].Nbr)
	assert.Equal(t, InterplanetaryMissilesID, prods[2].ID)
	assert.Equal(t, 1, prods[2].Nbr)
}

func TestExtractDKProductionWithABM(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/dk/production_with_abm.html")
	prods, _ := ExtractProduction(pageHTMLBytes)
	assert.Equal(t, 2, len(prods))
	assert.Equal(t, AntiBallisticMissilesID, prods[0].ID)
	assert.Equal(t, 1, prods[0].Nbr)
	assert.Equal(t, AntiBallisticMissilesID, prods[1].ID)
	assert.Equal(t, 1, prods[1].Nbr)
}

func TestIsShipID(t *testing.T) {
	assert.True(t, IsShipID(int(SmallCargoID)))
	assert.False(t, IsShipID(int(RocketLauncherID)))
}

func TestIsDefenseID(t *testing.T) {
	assert.True(t, IsDefenseID(int(RocketLauncherID)))
	assert.False(t, IsDefenseID(int(SmallCargoID)))
}

func TestIsTechID(t *testing.T) {
	assert.True(t, IsTechID(int(CombustionDriveID)))
	assert.False(t, IsTechID(int(SmallCargoID)))
}

func TestIsBuildingID(t *testing.T) {
	assert.True(t, IsBuildingID(int(MetalMineID)))
	assert.True(t, IsBuildingID(int(RoboticsFactoryID)))
	assert.False(t, IsBuildingID(int(SmallCargoID)))
}

func TestIsResourceBuildingID(t *testing.T) {
	assert.True(t, IsResourceBuildingID(int(MetalMineID)))
	assert.False(t, IsResourceBuildingID(int(RoboticsFactoryID)))
	assert.False(t, IsResourceBuildingID(int(SmallCargoID)))
}

func TestIsFacilityID(t *testing.T) {
	assert.False(t, IsFacilityID(int(MetalMineID)))
	assert.True(t, IsFacilityID(int(RoboticsFactoryID)))
	assert.False(t, IsFacilityID(int(SmallCargoID)))
}

func TestExtractEspionageReport_action(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/message_foreign_fleet_sighted.html")
	infos, _ := extractEspionageReport(pageHTMLBytes, time.FixedZone("OGT", 3600))
	assert.Equal(t, Action, infos.Type)
	assert.Equal(t, 6970988, infos.ID)
}

func TestExtractEspionageReport(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/spy_report_res_buildings_researches.html")
	infos, _ := extractEspionageReport(pageHTMLBytes, time.FixedZone("OGT", 3600))
	assert.Equal(t, Coordinate{4, 212, 6, PlanetType}, infos.Coordinate)
	assert.Equal(t, Report, infos.Type)
	assert.True(t, infos.HasFleet)
	assert.True(t, infos.HasDefenses)
	assert.True(t, infos.HasBuildings)
	assert.True(t, infos.HasResearches)
	assert.Equal(t, 6862893, infos.ID)
	assert.Equal(t, 0, infos.CounterEspionage)
	assert.Equal(t, 227034, infos.Metal)
	assert.Equal(t, 146970, infos.Crystal)
	assert.Equal(t, 24751, infos.Deuterium)
	assert.Equal(t, 2324, infos.Energy)
	assert.Equal(t, 20, *infos.MetalMine)
	assert.Equal(t, 14, *infos.CrystalMine)
	assert.Equal(t, 8, *infos.DeuteriumSynthesizer)
	assert.Equal(t, 19, *infos.SolarPlant)
	assert.Equal(t, 5, *infos.RoboticsFactory)
	assert.Equal(t, 2, *infos.Shipyard)
	assert.Equal(t, 5, *infos.MetalStorage)
	assert.Equal(t, 5, *infos.CrystalStorage)
	assert.Equal(t, 2, *infos.DeuteriumTank)
	assert.Equal(t, 3, *infos.ResearchLab)
	assert.Equal(t, 2, *infos.EspionageTechnology)
	assert.Equal(t, 1, *infos.ComputerTechnology)
	assert.Equal(t, 1, *infos.ArmourTechnology)
	assert.Equal(t, 1, *infos.EnergyTechnology)
	assert.Equal(t, 7, *infos.CombustionDrive)
	assert.Equal(t, 2, *infos.ImpulseDrive)
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
	pageHTMLBytes, _ := ioutil.ReadFile("samples/spy_report_no_pics.html")
	infos, err := extractEspionageReport(pageHTMLBytes, time.FixedZone("OGT", 3600))
	assert.Equal(t, ErrDeactivateHidePictures, err)
	assert.Equal(t, Coordinate{4, 203, 6, PlanetType}, infos.Coordinate)
	assert.Equal(t, Report, infos.Type)
	assert.True(t, infos.HasFleet)
	assert.True(t, infos.HasDefenses)
	assert.True(t, infos.HasBuildings)
	assert.True(t, infos.HasResearches)
	assert.Equal(t, 9142399, infos.ID)
	assert.Equal(t, 0, infos.CounterEspionage)
	assert.Equal(t, 1131895, infos.Metal)
	assert.Equal(t, 432515, infos.Crystal)
	assert.Equal(t, 114957, infos.Deuterium)
	assert.Equal(t, 4727, infos.Energy)
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
	pageHTMLBytes, _ := ioutil.ReadFile("samples/spy_report_moon.html")
	infos, _ := extractEspionageReport(pageHTMLBytes, time.FixedZone("OGT", 3600))
	assert.Equal(t, Coordinate{4, 116, 12, MoonType}, infos.Coordinate)
	assert.Equal(t, 6, *infos.LunarBase)
	assert.Equal(t, 4, *infos.SensorPhalanx)
	assert.Nil(t, infos.JumpGate)
}

func TestExtractEspionageReport1(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/spy_report_res_buildings_researches_fleet.html")
	infos, _ := extractEspionageReport(pageHTMLBytes, time.FixedZone("OGT", 3600))
	assert.Equal(t, 2, *infos.Battleship)
	assert.Equal(t, 1, *infos.Bomber)
}

func TestExtractEspionageReportThousands(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/spy_report_thousand_units.html")
	infos, _ := extractEspionageReport(pageHTMLBytes, time.FixedZone("OGT", 3600))
	assert.Equal(t, 4000, *infos.RocketLauncher)
	assert.Equal(t, 3882, *infos.LargeCargo)
	assert.Equal(t, 374, *infos.SolarSatellite)
}

func TestExtractEspionageReport_defence(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/spy_report_res_fleet_defences.html")
	infos, _ := extractEspionageReport(pageHTMLBytes, time.FixedZone("OGT", 3600))
	assert.True(t, infos.HasFleet)
	assert.True(t, infos.HasDefenses)
	assert.False(t, infos.HasBuildings)
	assert.False(t, infos.HasResearches)
	assert.Equal(t, 13, infos.CounterEspionage)
	assert.Equal(t, 57, *infos.RocketLauncher)
	assert.Equal(t, 57, *infos.LightLaser)
	assert.Equal(t, 61, *infos.HeavyLaser)
	assert.Nil(t, infos.GaussCannon)
}

func TestExtractEspionageReport_bandit(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/spy_report_inactive_bandit_lord.html")
	infos, _ := extractEspionageReport(pageHTMLBytes, time.FixedZone("OGT", 3600))
	assert.Equal(t, true, infos.IsBandit)
	assert.Equal(t, false, infos.IsStarlord)
}

func TestExtractEspionageReport_starlord(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/spy_report_active_star_lord.html")
	infos, _ := extractEspionageReport(pageHTMLBytes, time.FixedZone("OGT", 3600))
	assert.Equal(t, false, infos.IsBandit)
	assert.Equal(t, true, infos.IsStarlord)
}

func TestExtractEspionageReport_norank(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/spy_report_res_buildings_researches_fleet.html")
	infos, _ := extractEspionageReport(pageHTMLBytes, time.FixedZone("OGT", 3600))
	assert.Equal(t, false, infos.IsBandit)
	assert.Equal(t, false, infos.IsStarlord)

}

func TestExtractEspionageReport_username1(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/spy_report_inactive_bandit_lord.html")
	infos, _ := extractEspionageReport(pageHTMLBytes, time.FixedZone("OGT", 3600))
	assert.Equal(t, "Cid Granjeador", infos.Username)
}

func TestExtractEspionageReport_username2(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/spy_report_active_star_lord.html")
	infos, _ := extractEspionageReport(pageHTMLBytes, time.FixedZone("OGT", 3600))
	assert.Equal(t, "Commodore Nomad", infos.Username)
}

func TestExtractEspionageReport_apiKey(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/spy_report_active_star_lord.html")
	infos, _ := extractEspionageReport(pageHTMLBytes, time.FixedZone("OGT", 3600))
	assert.Equal(t, "sr-en-152-ea0b59302bfad7e3ab0f2d15f7ef2c6a4633b4ba", infos.APIKey)
}

func TestExtractEspionageReport_inactivetimer_within15(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/spy_report_res_buildings.html")
	infos, _ := extractEspionageReport(pageHTMLBytes, time.FixedZone("OGT", 3600))
	assert.Equal(t, 15, infos.LastActivity)
}

func TestExtractEspionageReport_inactivetimer_29(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/spy_report_res_buildings_researches.html")
	infos, _ := extractEspionageReport(pageHTMLBytes, time.FixedZone("OGT", 3600))
	assert.Equal(t, 29, infos.LastActivity)
}

func TestExtractEspionageReport_inactivetimer_over1h(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/spy_report_inactive_bandit_lord.html")
	infos, _ := extractEspionageReport(pageHTMLBytes, time.FixedZone("OGT", 3600))
	assert.Equal(t, 0, infos.LastActivity)
}

func TestGalaxyDistance(t *testing.T) {
	assert.Equal(t, 60000, galaxyDistance(6, 3, 6, false))
	assert.Equal(t, 20000, galaxyDistance(1, 2, 6, false))
	assert.Equal(t, 40000, galaxyDistance(1, 3, 6, false))
	assert.Equal(t, 60000, galaxyDistance(1, 4, 6, false))
	assert.Equal(t, 80000, galaxyDistance(1, 5, 6, false))
	assert.Equal(t, 100000, galaxyDistance(1, 6, 6, false))

	assert.Equal(t, 20000, galaxyDistance(1, 2, 6, true))
	assert.Equal(t, 40000, galaxyDistance(1, 3, 6, true))
	assert.Equal(t, 60000, galaxyDistance(1, 4, 6, true))
	assert.Equal(t, 40000, galaxyDistance(1, 5, 6, true))
	assert.Equal(t, 20000, galaxyDistance(1, 6, 6, true))
	assert.Equal(t, 20000, galaxyDistance(6, 1, 6, true))
}

func TestSystemDistance(t *testing.T) {
	assert.Equal(t, 3175, flightSystemDistance(35, 30, false))

	assert.Equal(t, 2795, flightSystemDistance(1, 2, true))
	assert.Equal(t, 2795, flightSystemDistance(1, 499, true))
	assert.Equal(t, 2890, flightSystemDistance(1, 3, true))
	assert.Equal(t, 2890, flightSystemDistance(1, 498, true))
	assert.Equal(t, 2890, flightSystemDistance(498, 1, true))
}

func TestPlanetDistance(t *testing.T) {
	assert.Equal(t, 1015, planetDistance(6, 3))
}

func TestDistance(t *testing.T) {
	assert.Equal(t, 1015, Distance(Coordinate{1, 1, 3, PlanetType}, Coordinate{1, 1, 6, PlanetType}, 6, true, true))
	assert.Equal(t, 2890, Distance(Coordinate{1, 1, 3, PlanetType}, Coordinate{1, 498, 6, PlanetType}, 6, true, true))
	assert.Equal(t, 20000, Distance(Coordinate{6, 1, 3, PlanetType}, Coordinate{1, 498, 6, PlanetType}, 6, true, true))
	assert.Equal(t, 5, Distance(Coordinate{6, 1, 3, PlanetType}, Coordinate{6, 1, 3, MoonType}, 6, true, true))
}

func TestCalcFlightTime(t *testing.T) {
	//secs, fuel := calcFlightTime(Coordinate{1, 1, 1}, Coordinate{1, 1, 2},
	//	1, false, false, 1, 1,
	//	ShipsInfos{LightFighter: 1}, Researches{})
	//assert.Equal(t, 2121, secs)
	//assert.Equal(t, 3, fuel)
}

func TestExtractFleetSlot_fleet1(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/fleet1.html")
	s := ExtractSlots(pageHTMLBytes)
	assert.Equal(t, 2, s.InUse)
	assert.Equal(t, 14, s.Total)
	assert.Equal(t, 0, s.ExpInUse)
	assert.Equal(t, 3, s.ExpTotal)
}

func TestExtractFleetSlot_movement(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/fleets_1.html")
	s := ExtractSlots(pageHTMLBytes)
	assert.Equal(t, 1, s.InUse)
	assert.Equal(t, 11, s.Total)
	assert.Equal(t, 0, s.ExpInUse)
	assert.Equal(t, 2, s.ExpTotal)
}

func TestExtractFleetSlot_commanders(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/fleet1_extract_slots_with_commanders.html")
	s := ExtractSlots(pageHTMLBytes)
	assert.Equal(t, 13, s.InUse)
	assert.Equal(t, 14, s.Total)
	assert.Equal(t, 2, s.ExpInUse)
	assert.Equal(t, 3, s.ExpTotal)
}

func TestGetResourcesDetails(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/fetch_resources.html")
	res, _ := ExtractResourcesDetails(pageHTMLBytes)
	assert.Equal(t, 380030343, res.Metal.Available)
	assert.Equal(t, 60510000, res.Metal.StorageCapacity)
	assert.Equal(t, 0, res.Metal.CurrentProduction)

	assert.Equal(t, 19320, res.Crystal.Available)
	assert.Equal(t, 9820000, res.Crystal.StorageCapacity)
	assert.Equal(t, 40636, res.Crystal.CurrentProduction)

	assert.Equal(t, 24902, res.Deuterium.Available)
	assert.Equal(t, 18005000, res.Deuterium.StorageCapacity)
	assert.Equal(t, 22508, res.Deuterium.CurrentProduction)

	assert.Equal(t, -8402, res.Energy.Available)
	assert.Equal(t, 10469, res.Energy.CurrentProduction)
	assert.Equal(t, -18871, res.Energy.Consumption)

	assert.Equal(t, 28500, res.Darkmatter.Available)
	assert.Equal(t, 0, res.Darkmatter.Purchased)
	assert.Equal(t, 28500, res.Darkmatter.Found)
}
