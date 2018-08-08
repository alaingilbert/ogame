package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShipyardCost(t *testing.T) {
	sy := newShipyard()
	assert.Equal(t, Resources{Metal: 3200, Crystal: 1600, Deuterium: 800}, sy.GetPrice(4))
}
