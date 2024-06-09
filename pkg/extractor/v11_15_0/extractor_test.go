package v11_15_0

import (
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/utils"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestExtractEspionageReport(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v11.15.0/en/spy_report.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, err := e.ExtractEspionageReport(pageHTMLBytes)
	assert.NoError(t, err)
	assert.Equal(t, time.Date(2024, 05, 23, 14, 20, 42, 0, time.Local), infos.Date) // "23.05.2024 22:20:42"
	assert.Equal(t, "Hmoa", infos.Username)
	assert.Equal(t, ogame.Collector, infos.CharacterClass)
	assert.Equal(t, int64(2_920_000), infos.Resources.Metal)
	assert.Equal(t, int64(5_743), infos.Resources.Energy)
	assert.Equal(t, utils.I64Ptr(24), infos.MetalMine)
	assert.Equal(t, utils.I64Ptr(9), infos.MetalStorage)
	assert.Equal(t, utils.I64Ptr(8), infos.EnergyTechnology)
	assert.Equal(t, utils.I64Ptr(10), infos.LaserTechnology)
	assert.Equal(t, int64(0), infos.LastActivity)
}

func TestExtractLfBonuses(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v11.15.4/en/lfbonuses.html")
	e := NewExtractor()
	bonuses, _ := e.ExtractLfBonuses(pageHTMLBytes)
	assert.Equal(t, 0.012, bonuses.LfShipBonuses[ogame.LightFighterID].CargoCapacity)
	assert.Equal(t, 0.006, bonuses.CostTimeBonuses[ogame.AllianceDepotID].Cost)
	assert.Equal(t, 0.012, bonuses.CostTimeBonuses[ogame.AllianceDepotID].Duration)
}

func TestExtractAllianceClass(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v11.15.5/en/allianceOverviewTab.html")
	e := NewExtractor()
	c, _ := e.ExtractAllianceClass(pageHTMLBytes)
	assert.Equal(t, ogame.Researcher, c)
}

func TestExtractPhalanx(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v11.15.5/pl/phalanx_acs.html")
	res, err := NewExtractor().ExtractPhalanx(pageHTMLBytes)
	assert.Nil(t, err)
	assert.Equal(t, 20, len(res))
	assert.Equal(t, ogame.DoParseCoord("M:6:228:7"), res[0].Origin)
	assert.Equal(t, int64(1), res[0].Ships.Deathstar)
	assert.Equal(t, int64(730), res[0].BaseSpeed)
	assert.Equal(t, ogame.DoParseCoord("6:228:9"), res[13].Origin)
	assert.Equal(t, ogame.DoParseCoord("6:229:9"), res[13].Destination)
	assert.Equal(t, int64(1_111_111), res[13].Ships.Bomber)
	assert.Equal(t, ogame.DoParseCoord("M:6:228:7"), res[18].Origin)
	assert.Equal(t, int64(997), res[18].Ships.Battleship)
}
