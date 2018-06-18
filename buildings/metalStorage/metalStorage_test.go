package metalStorage

import (
	"testing"

	"github.com/alaingilbert/ogame"
	"github.com/stretchr/testify/assert"
)

func TestCost(t *testing.T) {
	ms := New()
	assert.Equal(t, ogame.Resources{Metal: 1000}, ms.GetPrice(1))
	assert.Equal(t, ogame.Resources{Metal: 2000}, ms.GetPrice(2))
	assert.Equal(t, ogame.Resources{Metal: 4000}, ms.GetPrice(3))
	assert.Equal(t, ogame.Resources{Metal: 8000}, ms.GetPrice(4))
	assert.Equal(t, ogame.Resources{Metal: 16000}, ms.GetPrice(5))
	assert.Equal(t, ogame.Resources{Metal: 32000}, ms.GetPrice(6))
}

func TestCapacity(t *testing.T) {
	ms := New()
	assert.Equal(t, 10000, ms.Capacity(0))
	assert.Equal(t, 20000, ms.Capacity(1))
	assert.Equal(t, 40000, ms.Capacity(2))
	assert.Equal(t, 75000, ms.Capacity(3))
	assert.Equal(t, 140000, ms.Capacity(4))
	assert.Equal(t, 255000, ms.Capacity(5))
}
