package ogame

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExtractEspionageReportMessageIDs(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/messages.html")
	msgs, _ := extractEspionageReportMessageIDs(string(pageHTMLBytes))
	assert.Equal(t, 2, len(msgs))
	assert.Equal(t, Report, msgs[0].Type)
	assert.Equal(t, Coordinate{4, 117, 6}, msgs[0].Target)
	assert.Equal(t, "Fleet Command", msgs[0].From)
	assert.Equal(t, Action, msgs[1].Type)
	assert.Equal(t, "Space Monitoring", msgs[1].From)
	assert.Equal(t, Coordinate{4, 117, 9}, msgs[1].Target)
}

func TestExtractResourcesProductions(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/resource_settings.html")
	prods, _ := extractResourcesProductions(string(pageHTMLBytes))
	assert.Equal(t, Resources{Metal: 10352, Crystal: 5104, Deuterium: 1282, Energy: -52}, prods)
}

func TestExtractAttacks(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/event_list_attack.html")
	attacks := extractAttacks(string(pageHTMLBytes))
	assert.Equal(t, 1, len(attacks))
}

func TestExtractAttacks1(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/event_list_missile.html")
	attacks := extractAttacks(string(pageHTMLBytes))
	assert.Equal(t, 1, len(attacks))
	assert.Equal(t, 1, attacks[0].Missiles)
}

func TestExtractGalaxyInfos(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/galaxy_ajax.html")
	infos, _ := extractGalaxyInfos(string(pageHTMLBytes))
	assert.Equal(t, 5, len(infos))
}

func TestExtractUserInfos(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/overview_inactive.html")
	infos, _ := extractUserInfos(string(pageHTMLBytes), "en")
	assert.Equal(t, 1295, infos.Points)
}

func TestExtractUserInfos_de(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/de_overview.html")
	infos, _ := extractUserInfos(string(pageHTMLBytes), "de")
	assert.Equal(t, 0, infos.Points)
	assert.Equal(t, 2980, infos.Rank)
	assert.Equal(t, 2980, infos.Total)
}

func TestExtractUserInfos_jp(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/jp_overview.html")
	infos, _ := extractUserInfos(string(pageHTMLBytes), "jp")
	assert.Equal(t, 0, infos.Points)
	assert.Equal(t, 73, infos.Rank)
	assert.Equal(t, 73, infos.Total)
}

func TestExtractUserInfos_fr(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/fr_overview.html")
	infos, _ := extractUserInfos(string(pageHTMLBytes), "fr")
	assert.Equal(t, 0, infos.Points)
	assert.Equal(t, 3197, infos.Rank)
	assert.Equal(t, 3348, infos.Total)
}

func TestExtractPlanets(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/overview_inactive.html")
	planets := extractPlanets(string(pageHTMLBytes), nil)
	assert.Equal(t, 1, len(planets))
	assert.Equal(t, PlanetID(33672410), planets[0].ID)
	assert.Equal(t, Coordinate{Galaxy: 1, System: 301, Position: 5}, planets[0].Coordinate)
	assert.Equal(t, "Homeworld", planets[0].Name)
	assert.Equal(t, "https://gf2.geo.gfsrv.net/cdn46/9f84a481c0c9a83d3b000d801d9d9d.png", planets[0].Img)
	assert.Equal(t, 31, planets[0].Temperature.Min)
	assert.Equal(t, 71, planets[0].Temperature.Max)
	assert.Equal(t, 89, planets[0].Fields.Built)
	assert.Equal(t, 188, planets[0].Fields.Total)
	assert.Equal(t, 12800, planets[0].Diameter)
}

func TestExtractPlanets_fr(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/fr_overview.html")
	planets := extractPlanets(string(pageHTMLBytes), &OGame{language: "fr"})
	assert.Equal(t, 1, len(planets))
	assert.Equal(t, PlanetID(33629512), planets[0].ID)
	assert.Equal(t, Coordinate{Galaxy: 2, System: 180, Position: 4}, planets[0].Coordinate)
	assert.Equal(t, "planète mère", planets[0].Name)
	assert.Equal(t, "https://gf1.geo.gfsrv.net/cdn35/9545f984bcd53c816a1a8452356d00.png", planets[0].Img)
	assert.Equal(t, 48, planets[0].Temperature.Min)
	assert.Equal(t, 88, planets[0].Temperature.Max)
	assert.Equal(t, 0, planets[0].Fields.Built)
	assert.Equal(t, 188, planets[0].Fields.Total)
	assert.Equal(t, 12800, planets[0].Diameter)
}

func TestExtractGalaxyInfos_honorableTarget(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/galaxy_debris.html")
	infos, _ := extractGalaxyInfos(string(pageHTMLBytes))
	assert.False(t, infos[0].HonorableTarget)
	assert.True(t, infos[1].HonorableTarget)
	assert.False(t, infos[2].HonorableTarget)
}

func TestExtractGalaxyInfos_inactive(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/galaxy_debris.html")
	infos, _ := extractGalaxyInfos(string(pageHTMLBytes))
	assert.True(t, infos[0].Inactive)
	assert.False(t, infos[1].Inactive)
	assert.False(t, infos[2].Inactive)
}

func TestExtractGalaxyInfos_strongPlayer(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/galaxy_debris.html")
	infos, _ := extractGalaxyInfos(string(pageHTMLBytes))
	assert.False(t, infos[0].StrongPlayer)
	assert.True(t, infos[1].StrongPlayer)
	assert.False(t, infos[2].StrongPlayer)
}

func TestExtractGalaxyInfos_debris(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/galaxy_debris.html")
	infos, _ := extractGalaxyInfos(string(pageHTMLBytes))
	assert.Equal(t, 0, infos[0].Debris.Metal)
	assert.Equal(t, 700, infos[0].Debris.Crystal)
	assert.Equal(t, 1, infos[0].Debris.RecyclersNeeded)
}

func TestExtractGalaxyInfos_vacation(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/galaxy_ajax.html")
	infos, _ := extractGalaxyInfos(string(pageHTMLBytes))
	assert.Equal(t, 5, len(infos))
	assert.False(t, infos[0].Vacation)
	assert.True(t, infos[1].Vacation)
	assert.True(t, infos[2].Vacation)
	assert.False(t, infos[3].Vacation)
	assert.False(t, infos[4].Vacation)
}

func TestExtractGalaxyInfos_alliance(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/galaxy_ajax.html")
	infos, _ := extractGalaxyInfos(string(pageHTMLBytes))
	assert.Equal(t, 303, infos[3].Alliance.ID)
	assert.Equal(t, "Qrvix", infos[3].Alliance.Name)
	assert.Equal(t, 27, infos[3].Alliance.Rank)
	assert.Equal(t, 16, infos[3].Alliance.Member)
}

func TestUniverseSpeed(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/techtree_universe_speed.html")
	universeSpeed := extractUniverseSpeed(string(pageHTMLBytes))
	assert.Equal(t, 7, universeSpeed)
}

func TestCancel(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/overview_active_queue2.html")
	token, techID, listID, _ := extractCancelBuildingInfos(string(pageHTMLBytes))
	assert.Equal(t, "fef7488e4809150cd16e3fa8fa14db37", token)
	assert.Equal(t, 4, techID)
	assert.Equal(t, 2099434, listID)
}

func TestCancelResearch(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/overview_active_queue2.html")
	token, techID, listID, _ := extractCancelResearchInfos(string(pageHTMLBytes))
	assert.Equal(t, "fff7488e4809150cd16e3fa8fa14db37", token)
	assert.Equal(t, 120, techID)
	assert.Equal(t, 1769925, listID)
}

func TestGetConstructions(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/overview_active.html")
	buildingID, buildingCountdown, researchID, researchCountdown := extractConstructions(string(pageHTMLBytes))
	assert.Equal(t, ID(2), buildingID)
	assert.Equal(t, 731, buildingCountdown)
	assert.Equal(t, ID(115), researchID)
	assert.Equal(t, 927, researchCountdown)
}

func TestExtractFleet(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/fleets_1.html")
	fleets := extractFleets(string(pageHTMLBytes))
	assert.Equal(t, 1, len(fleets))
	assert.Equal(t, 4134, fleets[0].ArriveIn)
	assert.Equal(t, Coordinate{4, 116, 12}, fleets[0].Origin)
	assert.Equal(t, Coordinate{4, 117, 9}, fleets[0].Destination)
	assert.Equal(t, MissionID(3), fleets[0].Mission)
	assert.Equal(t, false, fleets[0].ReturnFlight)
	assert.Equal(t, FleetID(4494950), fleets[0].ID)
	assert.Equal(t, 1, fleets[0].Ships.SmallCargo)
	assert.Equal(t, 8, fleets[0].Ships.LargeCargo)
	assert.Equal(t, 1, fleets[0].Ships.LightFighter)
	assert.Equal(t, 1, fleets[0].Ships.ColonyShip)
	assert.Equal(t, 1, fleets[0].Ships.EspionageProbe)
	assert.Equal(t, Resources{Metal: 123, Crystal: 456, Deuterium: 789}, fleets[0].Resources)
}

func TestExtractFleet_returning(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/fleets_2.html")
	fleets := extractFleets(string(pageHTMLBytes))
	assert.Equal(t, 1, len(fleets))
	assert.Equal(t, Coordinate{4, 116, 12}, fleets[0].Origin)
	assert.Equal(t, Coordinate{4, 117, 9}, fleets[0].Destination)
	assert.Equal(t, MissionID(3), fleets[0].Mission)
	assert.Equal(t, true, fleets[0].ReturnFlight)
	assert.Equal(t, FleetID(0), fleets[0].ID)
	assert.Equal(t, 1, fleets[0].Ships.SmallCargo)
	assert.Equal(t, 8, fleets[0].Ships.LargeCargo)
	assert.Equal(t, 1, fleets[0].Ships.LightFighter)
	assert.Equal(t, 1, fleets[0].Ships.ColonyShip)
	assert.Equal(t, 1, fleets[0].Ships.EspionageProbe)
	assert.Equal(t, Resources{Metal: 123, Crystal: 456, Deuterium: 789}, fleets[0].Resources)
}

func TestExtractProduction(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/shipyard_queue.html")
	prods, _ := extractProduction(string(pageHTMLBytes))
	assert.Equal(t, 20, len(prods))
	assert.Equal(t, ID(203), prods[0].ID)
	assert.Equal(t, 4, prods[0].Nbr)
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
	infos, _ := extractEspionageReport(string(pageHTMLBytes), time.FixedZone("OGT", 3600))
	assert.Equal(t, Action, infos.Type)
}

func TestExtractEspionageReport(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/spy_report_res_buildings_researches.html")
	infos, _ := extractEspionageReport(string(pageHTMLBytes), time.FixedZone("OGT", 3600))
	assert.Equal(t, Coordinate{4, 212, 6}, infos.Coordinate)
	assert.Equal(t, Report, infos.Type)
	assert.Equal(t, 227034, infos.Metal)
	assert.Equal(t, 146970, infos.Crystal)
	assert.Equal(t, 24751, infos.Deuterium)
	assert.Equal(t, 2324, infos.Energy)
	assert.Equal(t, 20, infos.MetalMine)
	assert.Equal(t, 14, infos.CrystalMine)
	assert.Equal(t, 8, infos.DeuteriumSynthesizer)
	assert.Equal(t, 19, infos.SolarPlant)
	assert.Equal(t, 5, infos.RoboticsFactory)
	assert.Equal(t, 2, infos.Shipyard)
	assert.Equal(t, 5, infos.MetalStorage)
	assert.Equal(t, 5, infos.CrystalStorage)
	assert.Equal(t, 2, infos.DeuteriumTank)
	assert.Equal(t, 3, infos.ResearchLab)
	assert.Equal(t, 2, infos.EspionageTechnology)
	assert.Equal(t, 1, infos.ComputerTechnology)
	assert.Equal(t, 1, infos.ArmourTechnology)
	assert.Equal(t, 1, infos.EnergyTechnology)
	assert.Equal(t, 7, infos.CombustionDrive)
	assert.Equal(t, 2, infos.ImpulseDrive)
}

func TestExtractEspionageReport1(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/spy_report_res_buildings_researches_fleet.html")
	infos, _ := extractEspionageReport(string(pageHTMLBytes), time.FixedZone("OGT", 3600))
	assert.Equal(t, 2, infos.Battleship)
	assert.Equal(t, 1, infos.Bomber)
}

func TestExtractEspionageReport_defence(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("samples/spy_report_res_fleet_defences.html")
	infos, _ := extractEspionageReport(string(pageHTMLBytes), time.FixedZone("OGT", 3600))
	assert.Equal(t, 57, infos.RocketLauncher)
	assert.Equal(t, 57, infos.LightLaser)
	assert.Equal(t, 61, infos.HeavyLaser)
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
}

func TestSystemDistance(t *testing.T) {
	assert.Equal(t, 3175, systemDistance(35, 30, false))

	assert.Equal(t, 2795, systemDistance(1, 2, true))
	assert.Equal(t, 2795, systemDistance(1, 499, true))
	assert.Equal(t, 2890, systemDistance(1, 3, true))
	assert.Equal(t, 2890, systemDistance(1, 498, true))
}

func TestPlanetDistance(t *testing.T) {
	assert.Equal(t, 1015, planetDistance(6, 3))
}
