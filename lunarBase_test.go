package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLunarBasePrice(t *testing.T) {
	lb := newLunarBase()
	assert.Equal(t, Resources{Metal: 20000, Crystal: 40000, Deuterium: 20000}, lb.GetPrice(1))
	assert.Equal(t, Resources{Metal: 40000, Crystal: 80000, Deuterium: 40000}, lb.GetPrice(2))
	assert.Equal(t, Resources{Metal: 80000, Crystal: 160000, Deuterium: 80000}, lb.GetPrice(3))
}

func TestLunarBase_IsAvailable(t *testing.T) {
	lb := newLunarBase()
	assert.False(t, lb.IsAvailable(PlanetDest, ResourcesBuildings{}, Facilities{}, Researches{}, 0))
	assert.False(t, lb.IsAvailable(DebrisDest, ResourcesBuildings{}, Facilities{}, Researches{}, 0))
	assert.True(t, lb.IsAvailable(MoonDest, ResourcesBuildings{}, Facilities{}, Researches{}, 0))
}
