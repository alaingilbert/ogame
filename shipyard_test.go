package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShipyardCost(t *testing.T) {
	sy := newShipyard()
	assert.Equal(t, Resources{Metal: 3200, Crystal: 1600, Deuterium: 800}, sy.GetPrice(4))
}

func TestShipyard_GetLevel(t *testing.T) {
	s := newShipyard()
	assert.Equal(t, 0, s.GetLevel(newLazyResourcesBuildings(ResourcesBuildings{}), newLazyFacilities(Facilities{}), newLazyResearches(Researches{})))
	assert.Equal(t, 3, s.GetLevel(newLazyResourcesBuildings(ResourcesBuildings{}), newLazyFacilities(Facilities{Shipyard: 3}), newLazyResearches(Researches{})))
}
