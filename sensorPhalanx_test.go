package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSensorPhalanxPrice(t *testing.T) {
	sp := newSensorPhalanx()
	assert.Equal(t, Resources{Metal: 20000, Crystal: 40000, Deuterium: 20000}, sp.GetPrice(1))
	assert.Equal(t, Resources{Metal: 40000, Crystal: 80000, Deuterium: 40000}, sp.GetPrice(2))
	assert.Equal(t, Resources{Metal: 80000, Crystal: 160000, Deuterium: 80000}, sp.GetPrice(3))
}

func TestSensorPhalanx_IsAvailable(t *testing.T) {
	sp := newSensorPhalanx()
	assert.False(t, sp.IsAvailable(MoonType, newLazyResourcesBuildings(ResourcesBuildings{}), newLazyFacilities(Facilities{}), newLazyResearches(Researches{}), 0))
	assert.True(t, sp.IsAvailable(MoonType, newLazyResourcesBuildings(ResourcesBuildings{}), newLazyFacilities(Facilities{LunarBase: 1}), newLazyResearches(Researches{}), 0))
}

func TestSensorPhalanx_GetRange(t *testing.T) {
	sp := newSensorPhalanx()
	assert.Equal(t, 0, sp.GetRange(0))
	assert.Equal(t, 1, sp.GetRange(1))
	assert.Equal(t, 3, sp.GetRange(2))
	assert.Equal(t, 8, sp.GetRange(3))
	assert.Equal(t, 15, sp.GetRange(4))
}

func TestSensorPhalanx_ScanConsumption(t *testing.T) {
	sp := newSensorPhalanx()
	assert.Equal(t, 5000, sp.ScanConsumption())
}
