package ogame

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMetalMineProduction(t *testing.T) {
	mm := newMetalMine()
	assert.Equal(t, 30, mm.Production(1, 1, 1, 0, 0))
	assert.Equal(t, 63, mm.Production(1, 1, 1, 0, 1))
	assert.Equal(t, 120, mm.Production(4, 1, 1, 0, 0))
	assert.Equal(t, 252, mm.Production(4, 1, 1, 0, 1))
	assert.Equal(t, 96606+6762+210, mm.Production(7, 1, 1, 7, 29))
}

func TestMetalMineConstructionTime(t *testing.T) {
	mm := newMetalMine()
	assert.Equal(t, 8550*time.Second, mm.ConstructionTime(20, 7, Facilities{RoboticsFactory: 3}))
	assert.Equal(t, 30*time.Second, mm.ConstructionTime(4, 6, Facilities{}))
}

func TestMetalMine_GetLevel(t *testing.T) {
	mm := newMetalMine()
	assert.Equal(t, 0, mm.GetLevel(ResourcesBuildings{}.Lazy(), Facilities{}.Lazy(), Researches{}.Lazy()))
	assert.Equal(t, 3, mm.GetLevel(ResourcesBuildings{MetalMine: 3}.Lazy(), Facilities{}.Lazy(), Researches{}.Lazy()))
}

func TestMetalMine_EnergyConsumption(t *testing.T) {
	mm := newMetalMine()
	assert.Equal(t, 4601, mm.EnergyConsumption(29))
}

func TestMetalMine_IsAvailable(t *testing.T) {
	mm := newMetalMine()
	assert.True(t, mm.IsAvailable(PlanetType, ResourcesBuildings{}.Lazy(), Facilities{}.Lazy(), Researches{}.Lazy(), 0))
	assert.False(t, mm.IsAvailable(DebrisType, ResourcesBuildings{}.Lazy(), Facilities{}.Lazy(), Researches{}.Lazy(), 0))
	assert.False(t, mm.IsAvailable(MoonType, ResourcesBuildings{}.Lazy(), Facilities{}.Lazy(), Researches{}.Lazy(), 0))
}
