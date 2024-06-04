package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTerraformerIsAvailable(t *testing.T) {
	resourcesBuildings := ResourcesBuildings{}
	facilities := Facilities{NaniteFactory: 1, RoboticsFactory: 10, ResearchLab: 1}
	researches := Researches{EnergyTechnology: 12, ComputerTechnology: 10}
	b := newTerraformer()
	avail := b.IsAvailable(PlanetType, resourcesBuildings, LfBuildings{}, LfResearches{}, facilities, researches, 0, NoClass)
	assert.True(t, avail)
}

func TestTerraformerIsAvailable_NoTech(t *testing.T) {
	resourcesBuildings := ResourcesBuildings{}
	facilities := Facilities{NaniteFactory: 1}
	researches := Researches{EnergyTechnology: 11}
	b := newTerraformer()
	avail := b.IsAvailable(PlanetType, resourcesBuildings, LfBuildings{}, LfResearches{}, facilities, researches, 0, NoClass)
	assert.False(t, avail)
}

func TestTerraformerGetPrice(t *testing.T) {
	b := newTerraformer()
	assert.Equal(t, Resources{Crystal: 50000, Deuterium: 100000, Energy: 1000}, b.GetPrice(1, LfBonuses{}))
	assert.Equal(t, Resources{Crystal: 100000, Deuterium: 200000, Energy: 2000}, b.GetPrice(2, LfBonuses{}))
	assert.Equal(t, Resources{Crystal: 200000, Deuterium: 400000, Energy: 4000}, b.GetPrice(3, LfBonuses{}))
}
