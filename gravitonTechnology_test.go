package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGravitonTechnology_IsAvailable(t *testing.T) {
	gt := newGravitonTechnology()
	assert.False(t, gt.IsAvailable(PlanetType, newLazyResourcesBuildings(ResourcesBuildings{}), newLazyFacilities(Facilities{}), newLazyResearches(Researches{}), 0))
	assert.False(t, gt.IsAvailable(PlanetType, newLazyResourcesBuildings(ResourcesBuildings{}), newLazyFacilities(Facilities{ResearchLab: 12}), newLazyResearches(Researches{}), 299999))
	assert.True(t, gt.IsAvailable(PlanetType, newLazyResourcesBuildings(ResourcesBuildings{}), newLazyFacilities(Facilities{ResearchLab: 12}), newLazyResearches(Researches{}), 300000))
}
