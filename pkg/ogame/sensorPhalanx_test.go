package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSensorPhalanxPrice(t *testing.T) {
	sp := newSensorPhalanx()
	assert.Equal(t, Resources{Metal: 20000, Crystal: 40000, Deuterium: 20000}, sp.GetPrice(1, LfBonuses{}))
	assert.Equal(t, Resources{Metal: 40000, Crystal: 80000, Deuterium: 40000}, sp.GetPrice(2, LfBonuses{}))
	assert.Equal(t, Resources{Metal: 80000, Crystal: 160000, Deuterium: 80000}, sp.GetPrice(3, LfBonuses{}))
}

func TestSensorPhalanx_IsAvailable(t *testing.T) {
	sp := newSensorPhalanx()
	assert.False(t, sp.IsAvailable(MoonType, ResourcesBuildings{}, LfBuildings{}, LfResearches{}, Facilities{}, Researches{}, 0, NoClass))
	assert.True(t, sp.IsAvailable(MoonType, ResourcesBuildings{}, LfBuildings{}, LfResearches{}, Facilities{LunarBase: 1}, Researches{}, 0, NoClass))
}

func TestSensorPhalanx_GetRange(t *testing.T) {
	sp := newSensorPhalanx()
	assert.Equal(t, int64(0), sp.GetRange(0, false))
	assert.Equal(t, int64(1), sp.GetRange(1, false))
	assert.Equal(t, int64(3), sp.GetRange(2, false))
	assert.Equal(t, int64(8), sp.GetRange(3, false))
	assert.Equal(t, int64(15), sp.GetRange(4, false))

	assert.Equal(t, int64(1), sp.GetRange(1, true))
	assert.Equal(t, int64(18), sp.GetRange(4, true))
}

func TestSensorPhalanx_ScanConsumption(t *testing.T) {
	sp := newSensorPhalanx()
	assert.Equal(t, int64(5000), sp.ScanConsumption())
}
