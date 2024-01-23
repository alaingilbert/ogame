package v71

import (
	"github.com/alaingilbert/clockwork"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestExtractAttacksACSAttackSelf(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v8.6/en/eventlist_acs_attack_self.html")
	ownCoords := []ogame.Coordinate{{Galaxy: 4, System: 116, Position: 9, Type: ogame.PlanetType}}
	attacks, _ := NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), ownCoords)
	assert.Equal(t, 1, len(attacks))
	assert.Equal(t, ogame.GroupedAttack, attacks[0].MissionType)
	assert.Equal(t, int64(1), attacks[0].Ships.LightFighter)
}

func TestExtractAttacksACS_v71(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.1/en/eventlist_acs.html")
	ownCoords := make([]ogame.Coordinate, 0)
	attacks, _ := NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), ownCoords)
	assert.Equal(t, 1, len(attacks))
	assert.Equal(t, ogame.GroupedAttack, attacks[0].MissionType)
	assert.Equal(t, int64(200), attacks[0].Ships.LightFighter)
	assert.Equal(t, int64(9), attacks[0].Ships.SmallCargo)
}

func TestExtractAttacksACS_v72(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.2/en/eventlist_multipleACS.html")
	ownCoords := make([]ogame.Coordinate, 0)
	attacks, _ := NewExtractor().extractAttacks(pageHTMLBytes, clockwork.NewFakeClock(), ownCoords)
	assert.Equal(t, 3, len(attacks))
	assert.Equal(t, ogame.GroupedAttack, attacks[0].MissionType)
	assert.Equal(t, int64(14028), attacks[0].ID)
	assert.Equal(t, int64(14029), attacks[1].ID)
	assert.Equal(t, int64(673019), attacks[2].ID)
}

func TestExtractIsMobile(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.1/en/movement.html")
	isMobile := NewExtractor().ExtractIsMobile(pageHTMLBytes)
	assert.False(t, isMobile)

	pageHTMLBytes, _ = os.ReadFile("../../../samples/v7.2/en/movement_mobile.html")
	isMobile = NewExtractor().ExtractIsMobile(pageHTMLBytes)
	assert.True(t, isMobile)
}

func TestExtractActiveItems(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.6.6/en/overview_with_active_items.html")
	items, _ := NewExtractor().ExtractActiveItems(pageHTMLBytes)
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

func TestExtractBuffActivation(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.1/en/buffActivation.html")
	token, items, _ := NewExtractor().ExtractBuffActivation(pageHTMLBytes)
	assert.Equal(t, "081876002bf5791011097597836d3f5c", token)
	assert.Equal(t, 31, len(items))
}

func TestExtractDMCosts(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.1/en/overview_allDM.html")
	dmCosts, _ := NewExtractor().ExtractDMCosts(pageHTMLBytes)
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

	pageHTMLBytes, _ = os.ReadFile("../../../samples/v7.1/en/overview_shipyard_queue.html")
	dmCosts, _ = NewExtractor().ExtractDMCosts(pageHTMLBytes)
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

func TestExtractAllResources(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.1/en/traderOverview_waiting.html")
	resources, _ := NewExtractor().ExtractAllResources(pageHTMLBytes)
	assert.Equal(t, 12, len(resources))
	assert.Equal(t, ogame.Resources{Metal: 97696396, Crystal: 30582087, Deuterium: 32752170}, resources[33698658])
	assert.Equal(t, ogame.Resources{Metal: 133578, Crystal: 74977, Deuterium: 66899}, resources[33702461])
	assert.Equal(t, ogame.Resources{Metal: 0, Crystal: 0, Deuterium: 2676231}, resources[33741598])
}

func TestExtractAllResourcesTwV902(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v9.0.2/tw/traderauctioneer.html")
	resources, _ := NewExtractor().ExtractAllResources(pageHTMLBytes)
	assert.Equal(t, 1, len(resources))
	assert.Equal(t, ogame.Resources{Metal: 1005, Crystal: 1002, Deuterium: 0}, resources[33620229])
}

func TestExtractHighscore(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.1/en/highscore.html")
	highscore, _ := NewExtractor().ExtractHighscore(pageHTMLBytes)
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
	assert.Equal(t, ogame.Coordinate{Galaxy: 2, System: 356, Position: 15, Type: ogame.PlanetType}, highscore.Players[0].Homeworld)

	pageHTMLBytes, _ = os.ReadFile("../../../samples/v7.1/en/highscore_withSelf.html")
	highscore, _ = NewExtractor().ExtractHighscore(pageHTMLBytes)
	assert.Equal(t, int64(7), highscore.NbPage)
	assert.Equal(t, int64(2), highscore.CurrPage)
	assert.Equal(t, int64(1), highscore.Category)
	assert.Equal(t, int64(0), highscore.Type)
	assert.Equal(t, 100, len(highscore.Players))
	assert.Equal(t, "Bob", highscore.Players[7].Name)
	assert.Equal(t, int64(0), highscore.Players[7].ID)         // Player ID is broken for self
	assert.Equal(t, int64(0), highscore.Players[7].AllianceID) // Alliance ID is broken for self

	pageHTMLBytes, _ = os.ReadFile("../../../samples/v7.1/en/highscore_fullPage.html")
	highscore, _ = NewExtractor().ExtractHighscore(pageHTMLBytes)
	assert.Equal(t, int64(7), highscore.NbPage)
	assert.Equal(t, int64(2), highscore.CurrPage)
	assert.Equal(t, int64(1), highscore.Category)
	assert.Equal(t, int64(0), highscore.Type)
	assert.Equal(t, 100, len(highscore.Players))
	assert.Equal(t, "Bob", highscore.Players[7].Name)
	assert.Equal(t, int64(0), highscore.Players[7].ID)         // Player ID is broken for self
	assert.Equal(t, int64(0), highscore.Players[7].AllianceID) // Alliance ID is broken for self

	pageHTMLBytes, _ = os.ReadFile("../../../samples/v7.1/en/highscore_withShips.html")
	highscore, _ = NewExtractor().ExtractHighscore(pageHTMLBytes)
	assert.Equal(t, "malakopipis", highscore.Players[0].Name)
	assert.Equal(t, int64(125758), highscore.Players[0].Ships)
}

func TestExtractProduction(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.1/en/shipyard_queue.html")
	prods, secs, _ := NewExtractor().ExtractProduction(pageHTMLBytes)
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

func TestExtractIPM(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.1/nl/ipm_missile_launch.html")
	duration, max, token := NewExtractor().ExtractIPM(pageHTMLBytes)
	assert.Equal(t, "95b68270230217f7e9a813e4a4beb20e", token)
	assert.Equal(t, int64(25), max)
	assert.Equal(t, int64(248), duration)
}

func TestExtractDestroyRockets(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.6.2/en/destroy_rockets.html")
	abm, ipm, token, _ := NewExtractor().ExtractDestroyRockets(pageHTMLBytes)
	assert.Equal(t, "3a1148bb0d2c6a18f323cf7f0ce09d2b", token)
	assert.Equal(t, int64(24), abm)
	assert.Equal(t, int64(6), ipm)
}

func TestExtractEspionageReport(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.1/en/spy_report.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.False(t, infos.HonorableTarget)
	assert.Equal(t, int64(66331), infos.Metal)
	assert.Equal(t, int64(58452), infos.Crystal)
	assert.Equal(t, int64(0), infos.Deuterium)
	assert.Equal(t, ogame.Collector, infos.CharacterClass)
}

func TestExtractEspionageReportAllianceClass(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v8.1/en/spy_report_alliance_class_trader.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, ogame.Trader, infos.AllianceClass)

	pageHTMLBytes, _ = os.ReadFile("../../../samples/v8.1/en/spy_report_alliance_class_warrior.html")
	e = NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ = e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, ogame.Warrior, infos.AllianceClass)

	//pageHTMLBytes, _ = os.ReadFile("../../samples/v8.1/en/spy_report_alliance_class_researcher.html")
	//infos, _ = NewExtractor().ExtractEspionageReport(pageHTMLBytes, time.FixedZone("OGT", 3600))
	//assert.Equal(t, Researcher, infos.AllianceClass)

	pageHTMLBytes, _ = os.ReadFile("../../../samples/v8.1/en/spy_report_alliance_no_class.html")
	e = NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ = e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, ogame.NoAllianceClass, infos.AllianceClass)
}

func TestExtractEspionageReportHonorable(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.1/en/spy_report_honorable.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.True(t, infos.HonorableTarget)
}

func TestExtractEspionageReportHonorableStrong(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.1/en/spy_report_honorable_strong.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.True(t, infos.HonorableTarget)
}

func TestGetResourcesDetails(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.1/en/fetchResources.html")
	res, _ := NewExtractor().ExtractResourcesDetails(pageHTMLBytes)
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

func TestExtractMoonFacilities(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.1/en/moon_facilities.html")
	res, _ := NewExtractor().ExtractFacilities(pageHTMLBytes)
	assert.Equal(t, int64(10), res.RoboticsFactory)
	assert.Equal(t, int64(1), res.Shipyard)
	assert.Equal(t, int64(10), res.LunarBase)
	assert.Equal(t, int64(6), res.SensorPhalanx)
	assert.Equal(t, int64(1), res.JumpGate)
}

func TestExtractCancelFleetTokenFromDoc(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.5.0/en/cancel_fleet.html")
	token, _ := NewExtractor().ExtractCancelFleetToken(pageHTMLBytes, ogame.FleetID(9078407))
	assert.Equal(t, "db3317fbe004641f7483e8074e34cda1", token)
}

func TestExtractCombatReportMessages(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.1/en/combat_reports.html")
	msgs, _, _ := NewExtractor().ExtractCombatReportMessagesSummary(pageHTMLBytes)
	assert.Equal(t, "cr-us-149-fe449460902860455db7ef57a522ae341f931a59", msgs[0].APIKey)
}

func TestExtractCombatReportMessages_lossContact(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.1/en/combat_reports_loss_contact.html")
	msgs, _, _ := NewExtractor().ExtractCombatReportMessagesSummary(pageHTMLBytes)
	assert.Equal(t, 8, len(msgs))
}

func TestExtractPlanet_ro(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v7.1/ro/overview.html")
	planet, _ := NewExtractor().ExtractPlanet(pageHTMLBytes, ogame.PlanetID(33629199))
	assert.Equal(t, "Planeta Principala", planet.Name)
	assert.Equal(t, int64(12800), planet.Diameter)
	assert.Equal(t, int64(31), planet.Temperature.Min)
	assert.Equal(t, int64(71), planet.Temperature.Max)
	assert.Equal(t, int64(0), planet.Fields.Built)
	assert.Equal(t, int64(193), planet.Fields.Total)
	assert.Equal(t, ogame.PlanetID(33629199), planet.ID)
	assert.Equal(t, ogame.Coordinate{Galaxy: 1, System: 185, Position: 4, Type: ogame.PlanetType}, planet.Coordinate)
	assert.Nil(t, planet.Moon)
}

func TestV71ExtractEspionageReportMessages(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/unversioned/messages_loot_percentage.html")
	msgs, _, _ := NewExtractor().ExtractEspionageReportMessageIDs(pageHTMLBytes)
	assert.Equal(t, 1.0, msgs[0].LootPercentage)
	assert.Equal(t, 0.5, msgs[1].LootPercentage)
	assert.Equal(t, 0.5, msgs[2].LootPercentage)
}
