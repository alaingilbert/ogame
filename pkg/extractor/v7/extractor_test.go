package v7

import (
	"github.com/alaingilbert/clockwork"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestCancelResearch(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7/overview_cancels.html")
	token, techID, listID, _ := NewExtractor().ExtractCancelResearchInfos(pageHTMLBytes)
	assert.Equal(t, "9d44b41d8136dffadab759749508105e", token)
	assert.Equal(t, int64(124), techID)
	assert.Equal(t, int64(1324883), listID)
}

func TestCancel(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7/overview_cancels.html")
	token, techID, listID, _ := NewExtractor().ExtractCancelBuildingInfos(pageHTMLBytes)
	assert.Equal(t, "cf00a76b307f5cabf867af0d61ad1991", token)
	assert.Equal(t, int64(23), techID)
	assert.Equal(t, int64(1336041), listID)
}

func TestExtractOverviewProduction(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.1/en/overview_shipyard_queue.html")
	prods, countdown, _ := NewExtractor().ExtractOverviewProduction(pageHTMLBytes)
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

func TestExtractEspionageReport(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7/spy_report.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, int64(15), infos.LastActivity)
	assert.Equal(t, int64(7), *infos.SmallCargo)
}

func TestExtractCombatReportMessages(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7/combat_reports_msgs.html")
	msgs, _, _ := NewExtractor().ExtractCombatReportMessagesSummary(pageHTMLBytes)
	assert.Equal(t, 10, len(msgs))
}

func TestExtractCombatReportMessages_Debris(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7/combat_reports_debris.html")
	msgs, _, _ := NewExtractor().ExtractCombatReportMessagesSummary(pageHTMLBytes)
	assert.Equal(t, int64(2400), msgs[0].DebrisField)
}

func TestExtractShips_fleetdispatch(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7/fleetdispatch.html")
	ships := NewExtractor().ExtractFleet1Ships(pageHTMLBytes)
	assert.Equal(t, int64(6), ships.SmallCargo)
	assert.Equal(t, int64(1), ships.ColonyShip)
	assert.Equal(t, int64(0), ships.Crawler)
}

func TestGetResourcesDetails(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7/fetchResources.html")
	res, _ := NewExtractor().ExtractResourcesDetails(pageHTMLBytes)
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

func TestExtractResourcesBuildings(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7/supplies.html")
	res, _ := NewExtractor().ExtractResourcesBuildings(pageHTMLBytes)
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

func TestExtractResourceSettings(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7/resource_settings.html")
	settings, _, _ := NewExtractor().ExtractResourceSettings(pageHTMLBytes)
	assert.Equal(t, ogame.ResourceSettings{MetalMine: 100, CrystalMine: 100, DeuteriumSynthesizer: 100, SolarPlant: 100, FusionReactor: 0, SolarSatellite: 0, Crawler: 0}, settings)
}

func TestExtractShips(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7/shipyard.html")
	ships, _ := NewExtractor().ExtractShips(pageHTMLBytes)
	assert.Equal(t, int64(6), ships.SmallCargo)
	assert.Equal(t, int64(1), ships.ColonyShip)
	assert.Equal(t, int64(9), ships.Crawler)
}

func TestExtractShips_build(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7/shipyard_build.html")
	ships, _ := NewExtractor().ExtractShips(pageHTMLBytes)
	assert.Equal(t, int64(33), ships.Cruiser)
}

func TestExtractResearch(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7/researches.html")
	res := NewExtractor().ExtractResearch(pageHTMLBytes)
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

func TestExtractResearch_2(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7/researches2.html")
	res := NewExtractor().ExtractResearch(pageHTMLBytes)
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

func TestExtractFacilities(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7/facilities.html")
	res, _ := NewExtractor().ExtractFacilities(pageHTMLBytes)
	assert.Equal(t, int64(3), res.RoboticsFactory)
	assert.Equal(t, int64(7), res.Shipyard)
	assert.Equal(t, int64(6), res.ResearchLab)
	assert.Equal(t, int64(0), res.AllianceDepot)
	assert.Equal(t, int64(0), res.MissileSilo)
	assert.Equal(t, int64(0), res.NaniteFactory)
	assert.Equal(t, int64(0), res.Terraformer)
	assert.Equal(t, int64(0), res.SpaceDock)
}

func TestExtractDefense(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7/defenses.html")
	defense, _ := NewExtractor().ExtractDefense(pageHTMLBytes)
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

func TestExtractMarketplaceMessages(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.2/en/sales_messages.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	msgs, _, _ := e.ExtractMarketplaceMessages(pageHTMLBytes)
	assert.Equal(t, 9, len(msgs))
	assert.Equal(t, int64(12912161), msgs[3].ID)
	assert.Equal(t, int64(27), msgs[3].Type)
	assert.Equal(t, int64(1379), msgs[3].MarketTransactionID)
	assert.Equal(t, "164ba9f6e5cbfdaa03c061730767d779", msgs[3].Token)
}

func TestExtractExpeditionMessages(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.2/en/expedition_messages.html")
	e := NewExtractor()
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

func TestExtractGalaxyExpeditionDebrisDM(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.1/fr/galaxy_darkmatter_df.html")
	infos, err := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.NoError(t, err)
	assert.Equal(t, int64(3137), infos.Events.Darkmatter)
	assert.False(t, infos.Events.HasAsteroid)
}

func TestExtractGalaxyAsteroid(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.2/en/galaxyContent_asteroid.html")
	infos, err := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.NoError(t, err)
	assert.True(t, infos.Events.HasAsteroid)
}

func TestExtractGalaxyExpeditionDebris(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7/galaxy_debris16.html")
	infos, err := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), infos.ExpeditionDebris.Metal)
	assert.Equal(t, int64(2300), infos.ExpeditionDebris.Crystal)
	assert.Equal(t, int64(1), infos.ExpeditionDebris.PathfindersNeeded)
}

func TestExtractGalaxyTWExpeditionDebris(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.5.2/tw/galaxy_debris16.html")
	infos, err := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.NoError(t, err)
	assert.Equal(t, int64(4275000), infos.ExpeditionDebris.Metal)
	assert.Equal(t, int64(2953000), infos.ExpeditionDebris.Crystal)
	assert.Equal(t, int64(467), infos.ExpeditionDebris.PathfindersNeeded)
}

func TestExtractGalaxyExpeditionDebris2(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7/galaxy_debris16_2.html")
	infos, err := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.NoError(t, err)
	assert.Equal(t, int64(7200), infos.ExpeditionDebris.Metal)
	assert.Equal(t, int64(7200), infos.ExpeditionDebris.Crystal)
	assert.Equal(t, int64(1), infos.ExpeditionDebris.PathfindersNeeded)
}

func TestExtractGalaxyExpeditionDebrisMobile(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7/galaxy_debris16_mobile.html")
	_, err := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.EqualError(t, err, "mobile view not supported")
}

func TestExtractGalaxyIgnoredPlayer(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7/galaxy_ignored_player.html")
	infos, err := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.NoError(t, err)
	assert.Equal(t, int64(33637068), infos.Position(7).ID)
	assert.Equal(t, int64(102418), infos.Position(7).Player.ID)
	assert.Equal(t, "Procurator Serpentis", infos.Position(7).Player.Name)
}

func TestExtractGalaxyV7NoExpeditionDebris(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7/galaxy_no_debris16.html")
	infos, err := NewExtractor().ExtractGalaxyInfos(pageHTMLBytes, "Commodore Nomade", 123, 456)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), infos.ExpeditionDebris.Metal)
	assert.Equal(t, int64(0), infos.ExpeditionDebris.Crystal)
	assert.Equal(t, int64(0), infos.ExpeditionDebris.PathfindersNeeded)
}

func TestExtractUserInfos_es(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.6.5/es/overview.html")
	e := NewExtractor()
	e.SetLanguage("es")
	infos, _ := e.ExtractUserInfos(pageHTMLBytes)

	assert.Equal(t, int64(0), infos.Points)
	assert.Equal(t, int64(2976), infos.Rank)
	assert.Equal(t, int64(2977), infos.Total)
	assert.Equal(t, "Commodore Navi", infos.PlayerName)
}

func TestExtractFleetSlot_FleetDispatch(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7/fleetdispatch.html")
	s, _ := NewExtractor().ExtractSlots(pageHTMLBytes)
	assert.Equal(t, int64(0), s.InUse)
	assert.Equal(t, int64(4), s.Total)
	assert.Equal(t, int64(0), s.ExpInUse)
	assert.Equal(t, int64(1), s.ExpTotal)
}

func TestGetConstructionsV7(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7/overview_supplies_in_construction.html")
	clock := clockwork.NewFakeClockAt(time.Date(2019, 11, 12, 9, 6, 43, 0, time.UTC))
	buildingID, buildingCountdown, researchID, researchCountdown, _, _, _, _ := ExtractConstructions(pageHTMLBytes, clock)
	assert.Equal(t, ogame.MetalMineID, buildingID)
	assert.Equal(t, int64(62), buildingCountdown)
	assert.Equal(t, ogame.EnergyTechnologyID, researchID)
	assert.Equal(t, int64(271), researchCountdown)
}
