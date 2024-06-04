package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpaceDockGetPrice(t *testing.T) {
	sd := newSpaceDock()
	assert.Equal(t, Resources{Metal: 200, Deuterium: 50, Energy: 50}, sd.GetPrice(1, LfBonuses{}))
	assert.Equal(t, Resources{Metal: 1000, Deuterium: 250, Energy: 125}, sd.GetPrice(2, LfBonuses{}))
	assert.Equal(t, Resources{Metal: 5000, Deuterium: 1250, Energy: 312}, sd.GetPrice(3, LfBonuses{}))
	assert.Equal(t, Resources{Metal: 25000, Deuterium: 6250, Energy: 781}, sd.GetPrice(4, LfBonuses{}))
	assert.Equal(t, Resources{Metal: 125000, Deuterium: 31250, Energy: 1953}, sd.GetPrice(5, LfBonuses{}))
	assert.Equal(t, Resources{Metal: 15625000, Deuterium: 3906250, Energy: 30517}, sd.GetPrice(8, LfBonuses{}))
}
