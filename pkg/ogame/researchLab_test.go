package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResearchLabCost(t *testing.T) {
	rl := newResearchLab()
	assert.Equal(t, Resources{Metal: 1600, Crystal: 3200, Deuterium: 1600}, rl.GetPrice(4, LfBonuses{}))
	assert.Equal(t, Resources{Metal: 6400, Crystal: 12800, Deuterium: 6400}, rl.GetPrice(6, LfBonuses{}))
}
