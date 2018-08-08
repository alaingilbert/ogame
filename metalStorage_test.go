package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetalStorageCost(t *testing.T) {
	ms := newMetalStorage()
	assert.Equal(t, Resources{Metal: 1000}, ms.GetPrice(1))
	assert.Equal(t, Resources{Metal: 2000}, ms.GetPrice(2))
	assert.Equal(t, Resources{Metal: 4000}, ms.GetPrice(3))
	assert.Equal(t, Resources{Metal: 8000}, ms.GetPrice(4))
	assert.Equal(t, Resources{Metal: 16000}, ms.GetPrice(5))
	assert.Equal(t, Resources{Metal: 32000}, ms.GetPrice(6))
}

func TestMetalStorageCapacity(t *testing.T) {
	ms := newMetalStorage()
	assert.Equal(t, 10000, ms.Capacity(0))
	assert.Equal(t, 20000, ms.Capacity(1))
	assert.Equal(t, 40000, ms.Capacity(2))
	assert.Equal(t, 75000, ms.Capacity(3))
	assert.Equal(t, 140000, ms.Capacity(4))
	assert.Equal(t, 255000, ms.Capacity(5))
}
