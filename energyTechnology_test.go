package ogame

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEnergyTechnologyConstructionTime(t *testing.T) {
	mm := newEnergyTechnology()
	universeSpeed := 7
	researchSpeed := universeSpeed
	ct := mm.ConstructionTime(5, researchSpeed, Facilities{ResearchLab: 3})
	assert.Equal(t, 1645*time.Second, ct)

	researchSpeed = universeSpeed * 2
	ct = mm.ConstructionTime(5, researchSpeed, Facilities{ResearchLab: 3})
	assert.Equal(t, 822*time.Second, ct)
}

func TestEnergyTechnology_GetLevel(t *testing.T) {
	et := newEnergyTechnology()
	l := et.GetLevel(lazyResourcesBuildings, lazyFacilities, newLazyResearches(Researches{EnergyTechnology: 4}))
	assert.Equal(t, 4, l)
}
