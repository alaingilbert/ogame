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
	assert.Equal(t, 8550*time.Second, mm.ConstructionTime(20, 7, Facilities{RoboticsFactory: 3}, LfBonuses{}, NoClass, false))
	assert.Equal(t, 30*time.Second, mm.ConstructionTime(4, 6, Facilities{}, LfBonuses{}, NoClass, false))
}

func TestMetalMine_GetLevel(t *testing.T) {
	mm := newMetalMine()
	assert.Equal(t, int64(0), mm.GetLevel(ResourcesBuildings{}, Facilities{}, Researches{}))
	assert.Equal(t, int64(3), mm.GetLevel(ResourcesBuildings{MetalMine: 3}, Facilities{}, Researches{}))
}

func TestMetalMine_EnergyConsumption(t *testing.T) {
	mm := newMetalMine()
	assert.Equal(t, int64(4601), mm.EnergyConsumption(29))
}

func TestMetalMine_IsAvailable(t *testing.T) {
	mm := newMetalMine()
	assert.True(t, mm.IsAvailable(PlanetType, ResourcesBuildings{}, LfBuildings{}, LfResearches{}, Facilities{}, Researches{}, 0, NoClass))
	assert.False(t, mm.IsAvailable(DebrisType, ResourcesBuildings{}, LfBuildings{}, LfResearches{}, Facilities{}, Researches{}, 0, NoClass))
	assert.False(t, mm.IsAvailable(MoonType, ResourcesBuildings{}, LfBuildings{}, LfResearches{}, Facilities{}, Researches{}, 0, NoClass))
}

func TestDeconstructionPrice(t *testing.T) {
	mm := newMetalMine()
	assert.Equal(t, Resources{Metal: 3681620, Crystal: 920404}, mm.DeconstructionPrice(31, Researches{IonTechnology: 17}))
}
