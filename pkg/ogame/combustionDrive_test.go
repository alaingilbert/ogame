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
	assert.False(t, cd.IsAvailable(PlanetType, ResourcesBuildings{}, LfBuildings{}, LfResearches{}, Facilities{}, Researches{EnergyTechnology: 1}, 0, NoClass))
	assert.True(t, cd.IsAvailable(PlanetType, ResourcesBuildings{}, LfBuildings{}, LfResearches{}, Facilities{ResearchLab: 1}, Researches{EnergyTechnology: 1}, 0, NoClass))
	assert.False(t, cd.IsAvailable(MoonType, ResourcesBuildings{}, LfBuildings{}, LfResearches{}, Facilities{}, Researches{EnergyTechnology: 1}, 0, NoClass))
}
