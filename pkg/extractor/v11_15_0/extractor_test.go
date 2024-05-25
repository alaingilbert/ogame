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
