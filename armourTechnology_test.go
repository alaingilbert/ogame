package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArmourTechnologyIsAvailable(t *testing.T) {
	resourcesBuildings := ResourcesBuildings{}
	facilities := Facilities{ResearchLab: 2}
	researches := Researches{}
	b := newArmourTechnology()
	avail := b.IsAvailable(PlanetDest, resourcesBuildings, facilities, researches, 0)
	assert.True(t, avail)
}

func TestArmourTechnologyIsAvailable_NoBuilding(t *testing.T) {
	resourcesBuildings := ResourcesBuildings{}
	facilities := Facilities{ResearchLab: 1}
	researches := Researches{}
	b := newArmourTechnology()
	avail := b.IsAvailable(PlanetDest, resourcesBuildings, facilities, researches, 0)
	assert.False(t, avail)
}
