package ogame

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFusionReactorCapacity(t *testing.T) {
	fr := newFusionReactor()
	assert.Equal(t, 38*time.Second, fr.ConstructionTime(2, 7, Facilities{RoboticsFactory: 3}))
}

func TestFusionReactor_IsAvailable(t *testing.T) {
	fr := newFusionReactor()
	assert.False(t, fr.IsAvailable(PlanetType, ResourcesBuildings{}.Lazy(), Facilities{}.Lazy(), Researches{}.Lazy(), 0))
	assert.False(t, fr.IsAvailable(PlanetType, ResourcesBuildings{DeuteriumSynthesizer: 4}.Lazy(), Facilities{ResearchLab: 1}.Lazy(), Researches{EnergyTechnology: 3}.Lazy(), 0))
	assert.False(t, fr.IsAvailable(PlanetType, ResourcesBuildings{DeuteriumSynthesizer: 5}.Lazy(), Facilities{ResearchLab: 1}.Lazy(), Researches{EnergyTechnology: 2}.Lazy(), 0))
	assert.True(t, fr.IsAvailable(PlanetType, ResourcesBuildings{DeuteriumSynthesizer: 5}.Lazy(), Facilities{ResearchLab: 1}.Lazy(), Researches{EnergyTechnology: 3}.Lazy(), 0))
	assert.True(t, fr.IsAvailable(PlanetType, ResourcesBuildings{DeuteriumSynthesizer: 6}.Lazy(), Facilities{ResearchLab: 1}.Lazy(), Researches{EnergyTechnology: 4}.Lazy(), 0))
}

func TestFusionReactor_Production(t *testing.T) {
	fr := newFusionReactor()
	assert.Equal(t, 3002, fr.Production(12, 13))
}

func TestFusionReactor_GetFuelConsumption(t *testing.T) {
	fr := newFusionReactor()
	assert.Equal(t, 1486, fr.GetFuelConsumption(7, 1.0, 9))
	assert.Equal(t, 1040, fr.GetFuelConsumption(7, 0.7, 9))
}
