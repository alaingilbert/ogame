package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArmourTechnologyIsAvailable(t *testing.T) {
	b := newArmourTechnology()
	avail := b.IsAvailable(PlanetType, ResourcesBuildings{}, LfBuildings{}, LfResearches{}, Facilities{ResearchLab: 2}, Researches{}, 0, NoClass)
	assert.True(t, avail)
}

func TestArmourTechnologyIsAvailable_NoBuilding(t *testing.T) {
	b := newArmourTechnology()
	avail := b.IsAvailable(PlanetType, ResourcesBuildings{}, LfBuildings{}, LfResearches{}, Facilities{ResearchLab: 1}, Researches{}, 0, NoClass)
	assert.False(t, avail)
}
