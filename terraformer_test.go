package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTerraformerIsAvailable(t *testing.T) {
	resourcesBuildings := ResourcesBuildings{}
	facilities := Facilities{NaniteFactory: 1}
	researches := Researches{EnergyTechnology: 12}
	b := newTerraformer()
	avail := b.IsAvailable(PlanetType, resourcesBuildings, facilities, researches, 0)
	assert.True(t, avail)
}

func TestTerraformerIsAvailable_NoTech(t *testing.T) {
	resourcesBuildings := ResourcesBuildings{}
	facilities := Facilities{NaniteFactory: 1}
	researches := Researches{EnergyTechnology: 11}
	b := newTerraformer()
	avail := b.IsAvailable(PlanetType, resourcesBuildings, facilities, researches, 0)
	assert.False(t, avail)
}

func TestTerraformerGetPrice(t *testing.T) {
	b := newTerraformer()
	assert.Equal(t, Resources{Crystal: 50000, Deuterium: 100000, Energy: 1000}, b.GetPrice(1))
	assert.Equal(t, Resources{Crystal: 100000, Deuterium: 200000, Energy: 2000}, b.GetPrice(2))
	assert.Equal(t, Resources{Crystal: 200000, Deuterium: 400000, Energy: 4000}, b.GetPrice(3))
}
