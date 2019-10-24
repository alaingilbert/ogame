package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArmourTechnologyIsAvailable(t *testing.T) {
	b := newArmourTechnology()
	avail := b.IsAvailable(PlanetType, newLazyResourcesBuildings(ResourcesBuildings{}), newLazyFacilities(Facilities{ResearchLab: 2}), newLazyResearches(Researches{}), 0)
	assert.True(t, avail)
}

func TestArmourTechnologyIsAvailable_NoBuilding(t *testing.T) {
	b := newArmourTechnology()
	avail := b.IsAvailable(PlanetType, newLazyResourcesBuildings(ResourcesBuildings{}), newLazyFacilities(Facilities{ResearchLab: 1}), newLazyResearches(Researches{}), 0)
	assert.False(t, avail)
}
