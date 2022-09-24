package v9

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/alaingilbert/clockwork"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/stretchr/testify/assert"
)

func TestExtractResourcesDetailsFromFullPage(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../../samples/v9.0.0/en/overview2.html")
	res := NewExtractor().ExtractResourcesDetailsFromFullPage(pageHTMLBytes)
	assert.Equal(t, int64(6182), res.Metal.Available)
	assert.Equal(t, int64(10060), res.Metal.CurrentProduction)
	assert.Equal(t, int64(1590000), res.Metal.StorageCapacity)
	assert.Equal(t, int64(84388), res.Crystal.Available)
	assert.Equal(t, int64(4989), res.Crystal.CurrentProduction)
	assert.Equal(t, int64(1590000), res.Crystal.StorageCapacity)
	assert.Equal(t, int64(100188), res.Deuterium.Available)
	assert.Equal(t, int64(3499), res.Deuterium.CurrentProduction)
	assert.Equal(t, int64(865000), res.Deuterium.StorageCapacity)
	assert.Equal(t, int64(-1679), res.Energy.Available)
	assert.Equal(t, int64(2690), res.Energy.CurrentProduction)
	assert.Equal(t, int64(-4369), res.Energy.Consumption)
	assert.Equal(t, int64(8000), res.Darkmatter.Available)
	assert.Equal(t, int64(0), res.Darkmatter.Purchased)
	assert.Equal(t, int64(8000), res.Darkmatter.Found)
}

func TestExtractResources(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../../samples/v9.0.0/en/overview.html")
	res := NewExtractor().ExtractResources(pageHTMLBytes)
	assert.Equal(t, int64(10000), res.Metal)
	assert.Equal(t, int64(10000), res.Crystal)
	assert.Equal(t, int64(7829), res.Deuterium)
	assert.Equal(t, int64(26), res.Energy)
	assert.Equal(t, int64(10000000), res.Darkmatter)
}

func TestExtractEspionageReport(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../../samples/v9.0.0/en/spy_report.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, err := e.ExtractEspionageReport(pageHTMLBytes)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), infos.LastActivity)
}

func TestExtractOverviewProduction(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../../samples/v9.0.2/en/lifeform/overview_all_queues.html")
	prods, countdown, _ := NewExtractor().ExtractOverviewProduction(pageHTMLBytes)
	assert.Equal(t, 4, len(prods))
	assert.Equal(t, int64(1660), countdown)
	assert.Equal(t, ogame.SmallCargoID, prods[0].ID)
	assert.Equal(t, int64(1), prods[0].Nbr)
	assert.Equal(t, ogame.SmallCargoID, prods[1].ID)
	assert.Equal(t, int64(2), prods[1].Nbr)
	assert.Equal(t, ogame.LightFighterID, prods[2].ID)
	assert.Equal(t, int64(3), prods[2].Nbr)
	assert.Equal(t, ogame.RocketLauncherID, prods[3].ID)
	assert.Equal(t, int64(2), prods[3].Nbr)
}

func TestCancelResearch(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../../samples/v9.0.2/en/overview_all_queues.html")
	token, id, listID, _ := NewExtractor().ExtractCancelResearchInfos(pageHTMLBytes)
	assert.Equal(t, "66f639922a3c76fe6074d12ae36e573e", token)
	assert.Equal(t, int64(108), id)
	assert.Equal(t, int64(3469490), listID)
}

func TestCancelResearchLF(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../../samples/v9.0.2/en/lifeform/overview_all_queues.html")
	e := NewExtractor()
	e.SetLifeformEnabled(true)
	token, id, listID, _ := e.ExtractCancelResearchInfos(pageHTMLBytes)
	assert.Equal(t, "07287218c9661bcc67b05ec1b6171fe8", token)
	assert.Equal(t, int64(113), id)
	assert.Equal(t, int64(3998106), listID)
}

func TestCancelLfBuildingLF(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../../samples/v9.0.2/en/lifeform/overview_all_queues.html")
	token, id, listID, _ := NewExtractor().ExtractCancelLfBuildingInfos(pageHTMLBytes)
	assert.Equal(t, "07287218c9661bcc67b05ec1b6171fe8", token)
	assert.Equal(t, int64(11101), id)
	assert.Equal(t, int64(3998104), listID)
}

func TestCancelBuilding(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../../samples/v9.0.2/en/overview_all_queues.html")
	token, id, listID, _ := NewExtractor().ExtractCancelBuildingInfos(pageHTMLBytes)
	assert.Equal(t, "66f639922a3c76fe6074d12ae36e573e", token)
	assert.Equal(t, int64(1), id)
	assert.Equal(t, int64(3469488), listID)
}

func TestGetConstructions(t *testing.T) {
	// Without lifeform
	pageHTMLBytes, _ := ioutil.ReadFile("../../../samples/v9.0.2/en/overview_all_queues.html")
	clock := clockwork.NewFakeClockAt(time.Date(2022, 8, 20, 12, 43, 11, 0, time.UTC))
	buildingID, buildingCountdown, researchID, researchCountdown := ExtractConstructions(pageHTMLBytes, clock)
	assert.Equal(t, ogame.MetalMineID, buildingID)
	assert.Equal(t, int64(5413), buildingCountdown)
	assert.Equal(t, ogame.ComputerTechnologyID, researchID)
	assert.Equal(t, int64(7), researchCountdown)

	// With lifeform
	pageHTMLBytes, _ = ioutil.ReadFile("../../../samples/v9.0.2/en/lifeform/overview_all_queues2.html")
	clock = clockwork.NewFakeClockAt(time.Date(2022, 8, 28, 17, 22, 26, 0, time.UTC))
	buildingID, buildingCountdown, researchID, researchCountdown = ExtractConstructions(pageHTMLBytes, clock)
	assert.Equal(t, ogame.MetalStorageID, buildingID)
	assert.Equal(t, int64(33483), buildingCountdown)
	assert.Equal(t, ogame.ComputerTechnologyID, researchID)
	assert.Equal(t, int64(18355), researchCountdown)
}

func TestExtractUserInfos(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../../samples/v9.0.4/en/overview.html")
	info, err := NewExtractor().ExtractUserInfos(pageHTMLBytes)
	assert.NoError(t, err)
	assert.Equal(t, int64(30478), info.Points)
	assert.Equal(t, int64(1102), info.Rank)
	assert.Equal(t, int64(2931), info.Total)
}

func TestExtractFleetResources(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../../samples/v9.0.4/en/lifeform/movement.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	e.SetLifeformEnabled(true)
	fleets := e.ExtractFleets(pageHTMLBytes)
	assert.Equal(t, int64(1), fleets[0].Resources.Metal)
	assert.Equal(t, int64(2), fleets[0].Resources.Crystal)
	assert.Equal(t, int64(3), fleets[0].Resources.Deuterium)
}

func TestExtractLfBuildings(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../../samples/v9.0.4/en/lfbuildings.html")
	res, _ := NewExtractor().ExtractLfBuildings(pageHTMLBytes)
	assert.Equal(t, int64(2), res.ResidentialSector)
	assert.Equal(t, int64(1), res.BiosphereFarm)
	assert.Equal(t, int64(0), res.ResearchCentre)
	assert.Equal(t, int64(0), res.AcademyOfSciences)
	assert.Equal(t, int64(0), res.NeuroCalibrationCentre)
	assert.Equal(t, int64(0), res.HighEnergySmelting)
	assert.Equal(t, int64(0), res.FoodSilo)
	assert.Equal(t, int64(0), res.FusionPoweredProduction)
	assert.Equal(t, int64(0), res.Skyscraper)
	assert.Equal(t, int64(0), res.BiotechLab)
	assert.Equal(t, int64(0), res.Metropolis)
	assert.Equal(t, int64(0), res.PlanetaryShield)
}

func TestExtractLfBuildingsRocktal(t *testing.T) {
	pageHTMLBytes, _ := ioutil.ReadFile("../../../samples/v9.0.4/en/lifeform/lfbuildings_rocktal.html")
	res, _ := NewExtractor().ExtractLfBuildings(pageHTMLBytes)
	assert.Equal(t, int64(0), res.ResidentialSector)
	assert.Equal(t, int64(0), res.BiosphereFarm)
	assert.Equal(t, int64(0), res.ResearchCentre)
	assert.Equal(t, int64(0), res.AcademyOfSciences)
	assert.Equal(t, int64(0), res.NeuroCalibrationCentre)
	assert.Equal(t, int64(0), res.HighEnergySmelting)
	assert.Equal(t, int64(0), res.FoodSilo)
	assert.Equal(t, int64(0), res.FusionPoweredProduction)
	assert.Equal(t, int64(0), res.Skyscraper)
	assert.Equal(t, int64(0), res.BiotechLab)
	assert.Equal(t, int64(0), res.Metropolis)
	assert.Equal(t, int64(0), res.PlanetaryShield)
	assert.Equal(t, int64(2), res.MeditationEnclave)
	assert.Equal(t, int64(1), res.CrystalFarm)
}
