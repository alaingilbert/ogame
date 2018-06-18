package researchLab

import (
	"testing"

	"github.com/alaingilbert/ogame"
	"github.com/stretchr/testify/assert"
)

func TestCost(t *testing.T) {
	rl := New()
	assert.Equal(t, ogame.Resources{Metal: 1600, Crystal: 3200, Deuterium: 1600}, rl.GetPrice(4))
	assert.Equal(t, ogame.Resources{Metal: 6400, Crystal: 12800, Deuterium: 6400}, rl.GetPrice(6))
}
