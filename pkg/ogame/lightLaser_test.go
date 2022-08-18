package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLightLaser_IsAvailable(t *testing.T) {
	ll := newLightLaser()
	assert.True(t, ll.IsAvailable(PlanetType, ResourcesBuildings{}.Lazy(), Facilities{ResearchLab: 0, Shipyard: 2, RoboticsFactory: 2}.Lazy(), Researches{EnergyTechnology: 2, LaserTechnology: 3}.Lazy(), 0, NoClass))
}
