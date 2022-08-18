package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourcesBuildings_ByID(t *testing.T) {
	r := ResourcesBuildings{
		MetalMine:            1,
		CrystalMine:          2,
		DeuteriumSynthesizer: 3,
		SolarPlant:           4,
		FusionReactor:        5,
		SolarSatellite:       6,
		MetalStorage:         7,
		CrystalStorage:       8,
		DeuteriumTank:        9,
	}
	assert.Equal(t, int64(1), r.ByID(MetalMineID))
	assert.Equal(t, int64(2), r.ByID(CrystalMineID))
	assert.Equal(t, int64(3), r.ByID(DeuteriumSynthesizerID))
	assert.Equal(t, int64(4), r.ByID(SolarPlantID))
	assert.Equal(t, int64(5), r.ByID(FusionReactorID))
	assert.Equal(t, int64(6), r.ByID(SolarSatelliteID))
	assert.Equal(t, int64(7), r.ByID(MetalStorageID))
	assert.Equal(t, int64(8), r.ByID(CrystalStorageID))
	assert.Equal(t, int64(9), r.ByID(DeuteriumTankID))
	assert.Equal(t, int64(0), r.ByID(ID(12345)))
}

func TestResourcesBuildings_String(t *testing.T) {
	r := ResourcesBuildings{
		MetalMine:            1,
		CrystalMine:          2,
		DeuteriumSynthesizer: 3,
		SolarPlant:           4,
		FusionReactor:        5,
		SolarSatellite:       6,
		MetalStorage:         7,
		CrystalStorage:       8,
		DeuteriumTank:        9,
	}
	expected := "\n" +
		"           Metal Mine: 1\n" +
		"         Crystal Mine: 2\n" +
		"Deuterium Synthesizer: 3\n" +
		"          Solar Plant: 4\n" +
		"       Fusion Reactor: 5\n" +
		"      Solar Satellite: 6\n" +
		"        Metal Storage: 7\n" +
		"      Crystal Storage: 8\n" +
		"       Deuterium Tank: 9"
	assert.Equal(t, expected, r.String())
}
