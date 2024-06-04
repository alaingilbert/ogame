package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetalStorageCost(t *testing.T) {
	ms := newMetalStorage()
	assert.Equal(t, Resources{Metal: 1000}, ms.GetPrice(1, LfBonuses{}))
	assert.Equal(t, Resources{Metal: 2000}, ms.GetPrice(2, LfBonuses{}))
	assert.Equal(t, Resources{Metal: 4000}, ms.GetPrice(3, LfBonuses{}))
	assert.Equal(t, Resources{Metal: 8000}, ms.GetPrice(4, LfBonuses{}))
	assert.Equal(t, Resources{Metal: 16000}, ms.GetPrice(5, LfBonuses{}))
	assert.Equal(t, Resources{Metal: 32000}, ms.GetPrice(6, LfBonuses{}))
}

func TestMetalStorageCapacity(t *testing.T) {
	ms := newMetalStorage()
	assert.Equal(t, int64(10000), ms.Capacity(0))
	assert.Equal(t, int64(20000), ms.Capacity(1))
	assert.Equal(t, int64(40000), ms.Capacity(2))
	assert.Equal(t, int64(75000), ms.Capacity(3))
	assert.Equal(t, int64(140000), ms.Capacity(4))
	assert.Equal(t, int64(255000), ms.Capacity(5))
}
