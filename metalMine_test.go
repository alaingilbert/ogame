package ogame

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMetalMineProduction(t *testing.T) {
	mm := newMetalMine()
	assert.Equal(t, int64(30), mm.Production(1, 1, 1, 0, 0))
	assert.Equal(t, int64(63), mm.Production(1, 1, 1, 0, 1))
	assert.Equal(t, int64(120), mm.Production(4, 1, 1, 0, 0))
	assert.Equal(t, int64(252), mm.Production(4, 1, 1, 0, 1))
	assert.Equal(t, int64(96606+6762+210), mm.Production(7, 1, 1, 7, 29))
}

func TestMetalMineConstructionTime(t *testing.T) {
	mm := newMetalMine()
	assert.Equal(t, 8550*time.Second, mm.ConstructionTime(20, 7, Facilities{RoboticsFactory: 3}, false, false))
	assert.Equal(t, 30*time.Second, mm.ConstructionTime(4, 6, Facilities{}, false, false))
}

func TestMetalMine_GetLevel(t *testing.T) {
	mm := newMetalMine()
	assert.Equal(t, int64(0), mm.GetLevel(ResourcesBuildings{}.Lazy(), Facilities{}.Lazy(), Researches{}.Lazy()))
	assert.Equal(t, int64(3), mm.GetLevel(ResourcesBuildings{MetalMine: 3}.Lazy(), Facilities{}.Lazy(), Researches{}.Lazy()))
}

func TestMetalMine_EnergyConsumption(t *testing.T) {
	mm := newMetalMine()
	assert.Equal(t, int64(4601), mm.EnergyConsumption(29))
}

func TestMetalMine_IsAvailable(t *testing.T) {
	mm := newMetalMine()
	assert.True(t, mm.IsAvailable(PlanetType, ResourcesBuildings{}.Lazy(), Facilities{}.Lazy(), Researches{}.Lazy(), 0))
	assert.False(t, mm.IsAvailable(DebrisType, ResourcesBuildings{}.Lazy(), Facilities{}.Lazy(), Researches{}.Lazy(), 0))
	assert.False(t, mm.IsAvailable(MoonType, ResourcesBuildings{}.Lazy(), Facilities{}.Lazy(), Researches{}.Lazy(), 0))
}
