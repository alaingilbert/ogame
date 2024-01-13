package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGravitonTechnology_IsAvailable(t *testing.T) {
	gt := newGravitonTechnology()
	assert.False(t, gt.IsAvailable(PlanetType, ResourcesBuildings{}, LfBuildings{}, LfResearches{}, Facilities{}, Researches{}, 0, NoClass))
	assert.False(t, gt.IsAvailable(PlanetType, ResourcesBuildings{}, LfBuildings{}, LfResearches{}, Facilities{ResearchLab: 12}, Researches{}, 299999, NoClass))
	assert.True(t, gt.IsAvailable(PlanetType, ResourcesBuildings{}, LfBuildings{}, LfResearches{}, Facilities{ResearchLab: 12}, Researches{}, 300000, NoClass))
}
