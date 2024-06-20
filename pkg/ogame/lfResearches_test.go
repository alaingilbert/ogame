package ogame

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestIntergalacticEnvoysConstructionTime(t *testing.T) {
	ie := newIntergalacticEnvoys()
	assert.Equal(t, (6*60+0)*time.Second, ie.ConstructionTime(2, 8, Facilities{}, LfBonuses{}, NoClass, false))
}

func TestCatalyserTechnologyPrice(t *testing.T) {
	ie := newCatalyserTechnology()
	assert.Equal(t, Resources{Metal: 177_347_025, Crystal: 106_408_215, Deuterium: 17_734_702}, ie.GetPrice(18, LfBonuses{}))
	assert.Equal(t, Resources{Metal: 165_819_468, Crystal: 99_491_681, Deuterium: 16_581_946}, ie.GetPrice(18, LfBonuses{PlanetLfResearchCostTimeBonus: CostTimeBonus{Cost: 0.065}}))
}
