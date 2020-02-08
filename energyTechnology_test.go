package ogame

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestEnergyTechnologyConstructionTime(t *testing.T) {
	mm := newEnergyTechnology()
	universeSpeed := int64(7)
	researchSpeed := int64(1)
	ct := mm.ConstructionTime(5, universeSpeed*researchSpeed, Facilities{ResearchLab: 3}, false, false)
	assert.Equal(t, 1645*time.Second, ct)

	researchSpeed = 2
	ct = mm.ConstructionTime(5, universeSpeed*researchSpeed, Facilities{ResearchLab: 3}, false, false)
	assert.Equal(t, 822*time.Second, ct)

	universeSpeed = 6
	researchSpeed = 1
	ct = mm.ConstructionTime(1, universeSpeed*researchSpeed, Facilities{ResearchLab: 0}, false, false)
	assert.Equal(t, 8*time.Minute, ct)

	researchSpeed = 2
	ct = mm.ConstructionTime(1, universeSpeed*researchSpeed, Facilities{ResearchLab: 0}, false, false)
	assert.Equal(t, 4*time.Minute, ct)

	universeSpeed = 1
	researchSpeed = 1
	ct = mm.ConstructionTime(1, universeSpeed*researchSpeed, Facilities{ResearchLab: 10}, false, false)
	assert.Equal(t, 261*time.Second, ct)
	ct = mm.ConstructionTime(1, universeSpeed*researchSpeed, Facilities{ResearchLab: 10}, true, false)
	assert.Equal(t, 196*time.Second, ct)
	ct = mm.ConstructionTime(1, universeSpeed*researchSpeed, Facilities{ResearchLab: 10}, true, true)
	assert.Equal(t, 147*time.Second, ct)
}

func TestEnergyTechnology_GetLevel(t *testing.T) {
	et := newEnergyTechnology()
	l := et.GetLevel(ResourcesBuildings{}.Lazy(), Facilities{}.Lazy(), Researches{EnergyTechnology: 4}.Lazy())
	assert.Equal(t, int64(4), l)
}
