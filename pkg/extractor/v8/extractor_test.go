package v8

import (
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestExtractEspionageReport(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v8.5/en/spy_report.html")
	e := NewExtractor()
	e.SetLocation(time.FixedZone("OGT", 3600))
	infos, _ := e.ExtractEspionageReport(pageHTMLBytes)
	assert.Equal(t, int64(15), infos.LastActivity)
}

func TestExtractPlanets(t *testing.T) {
	pageHTMLBytes, _ := os.ReadFile("../../../samples/v9.0.0/en/overview.html")
	planets := NewExtractor().ExtractPlanets(pageHTMLBytes)
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
