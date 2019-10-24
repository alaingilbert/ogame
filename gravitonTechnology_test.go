package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGravitonTechnology_IsAvailable(t *testing.T) {
	gt := newGravitonTechnology()
	assert.False(t, gt.IsAvailable(PlanetType, ResourcesBuildings{}.Lazy(), Facilities{}.Lazy(), Researches{}.Lazy(), 0))
	assert.False(t, gt.IsAvailable(PlanetType, ResourcesBuildings{}.Lazy(), Facilities{ResearchLab: 12}.Lazy(), Researches{}.Lazy(), 299999))
	assert.True(t, gt.IsAvailable(PlanetType, ResourcesBuildings{}.Lazy(), Facilities{ResearchLab: 12}.Lazy(), Researches{}.Lazy(), 300000))
}
