package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLunarBasePrice(t *testing.T) {
	lb := newLunarBase()
	assert.Equal(t, Resources{Metal: 20000, Crystal: 40000, Deuterium: 20000}, lb.GetPrice(1, LfBonuses{}))
	assert.Equal(t, Resources{Metal: 40000, Crystal: 80000, Deuterium: 40000}, lb.GetPrice(2, LfBonuses{}))
	assert.Equal(t, Resources{Metal: 80000, Crystal: 160000, Deuterium: 80000}, lb.GetPrice(3, LfBonuses{}))
}

func TestLunarBase_IsAvailable(t *testing.T) {
	lb := newLunarBase()
	assert.False(t, lb.IsAvailable(PlanetType, ResourcesBuildings{}, LfBuildings{}, LfResearches{}, Facilities{}, Researches{}, 0, NoClass))
	assert.False(t, lb.IsAvailable(DebrisType, ResourcesBuildings{}, LfBuildings{}, LfResearches{}, Facilities{}, Researches{}, 0, NoClass))
	assert.True(t, lb.IsAvailable(MoonType, ResourcesBuildings{}, LfBuildings{}, LfResearches{}, Facilities{}, Researches{}, 0, NoClass))
}
