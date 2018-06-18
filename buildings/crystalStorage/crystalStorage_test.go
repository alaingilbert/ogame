package crystalStorage

import (
	"testing"

	"github.com/alaingilbert/ogame"
	"github.com/stretchr/testify/assert"
)

func TestCost(t *testing.T) {
	cs := New()
	assert.Equal(t, ogame.Resources{Metal: 1000, Crystal: 500}, cs.GetPrice(1))
	assert.Equal(t, ogame.Resources{Metal: 2000, Crystal: 1000}, cs.GetPrice(2))
	assert.Equal(t, ogame.Resources{Metal: 4000, Crystal: 2000}, cs.GetPrice(3))
	assert.Equal(t, ogame.Resources{Metal: 8000, Crystal: 4000}, cs.GetPrice(4))
	assert.Equal(t, ogame.Resources{Metal: 16000, Crystal: 8000}, cs.GetPrice(5))
	assert.Equal(t, ogame.Resources{Metal: 32000, Crystal: 16000}, cs.GetPrice(6))
}

func TestCapacity(t *testing.T) {
	cs := New()
	assert.Equal(t, 10000, cs.Capacity(0))
	assert.Equal(t, 20000, cs.Capacity(1))
	assert.Equal(t, 40000, cs.Capacity(2))
	assert.Equal(t, 75000, cs.Capacity(3))
	assert.Equal(t, 140000, cs.Capacity(4))
	assert.Equal(t, 255000, cs.Capacity(5))
}
