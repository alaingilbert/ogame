package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCombustionDriveCost(t *testing.T) {
	cd := newCombustionDrive()
	assert.Equal(t, Resources{Metal: 12800, Deuterium: 19200}, cd.GetPrice(6))
}

func TestCombustionDrive_IsAvailable(t *testing.T) {
	cd := newCombustionDrive()
	assert.False(t, cd.IsAvailable(PlanetType, lazyResourcesBuildings, lazyFacilities, newLazyResearches(Researches{EnergyTechnology: 1}), 0))
	assert.True(t, cd.IsAvailable(PlanetType, lazyResourcesBuildings, newLazyFacilities(Facilities{ResearchLab: 1}), newLazyResearches(Researches{EnergyTechnology: 1}), 0))
	assert.False(t, cd.IsAvailable(MoonType, lazyResourcesBuildings, lazyFacilities, newLazyResearches(Researches{EnergyTechnology: 1}), 0))
}
