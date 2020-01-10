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
	assert.Equal(t, int64(0), s.GetLevel(ResourcesBuildings{}.Lazy(), Facilities{}.Lazy(), Researches{}.Lazy()))
	assert.Equal(t, int64(3), s.GetLevel(ResourcesBuildings{}.Lazy(), Facilities{Shipyard: 3}.Lazy(), Researches{}.Lazy()))
}
