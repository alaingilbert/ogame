package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGravitonTechnology_IsAvailable(t *testing.T) {
	gt := newGravitonTechnology()
	assert.False(t, gt.IsAvailable(PlanetDest, ResourcesBuildings{}, Facilities{}, Researches{}, 0))
	assert.False(t, gt.IsAvailable(PlanetDest, ResourcesBuildings{}, Facilities{ResearchLab: 12}, Researches{}, 299999))
	assert.True(t, gt.IsAvailable(PlanetDest, ResourcesBuildings{}, Facilities{ResearchLab: 12}, Researches{}, 300000))
}
