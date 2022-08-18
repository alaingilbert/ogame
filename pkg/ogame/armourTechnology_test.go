package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArmourTechnologyIsAvailable(t *testing.T) {
	b := newArmourTechnology()
	avail := b.IsAvailable(PlanetType, ResourcesBuildings{}.Lazy(), Facilities{ResearchLab: 2}.Lazy(), Researches{}.Lazy(), 0, NoClass)
	assert.True(t, avail)
}

func TestArmourTechnologyIsAvailable_NoBuilding(t *testing.T) {
	b := newArmourTechnology()
	avail := b.IsAvailable(PlanetType, ResourcesBuildings{}.Lazy(), Facilities{ResearchLab: 1}.Lazy(), Researches{}.Lazy(), 0, NoClass)
	assert.False(t, avail)
}
