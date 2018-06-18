package shipyard

import (
	"testing"

	"github.com/alaingilbert/ogame"
	"github.com/stretchr/testify/assert"
)

func TestCost(t *testing.T) {
	sy := New()
	assert.Equal(t, ogame.Resources{Metal: 3200, Crystal: 1600, Deuterium: 800}, sy.GetPrice(4))
}
