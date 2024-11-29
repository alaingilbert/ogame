package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeuteriumTankCost(t *testing.T) {
	dt := newDeuteriumTank()
	assert.Equal(t, Resources{Metal: 1000, Crystal: 1000}, dt.GetPrice(1, LfBonuses{}))
	assert.Equal(t, Resources{Metal: 2000, Crystal: 2000}, dt.GetPrice(2, LfBonuses{}))
	assert.Equal(t, Resources{Metal: 4000, Crystal: 4000}, dt.GetPrice(3, LfBonuses{}))
	assert.Equal(t, Resources{Metal: 8000, Crystal: 8000}, dt.GetPrice(4, LfBonuses{}))
	assert.Equal(t, Resources{Metal: 16000, Crystal: 16000}, dt.GetPrice(5, LfBonuses{}))
	assert.Equal(t, Resources{Metal: 32000, Crystal: 32000}, dt.GetPrice(6, LfBonuses{}))
}

func TestDeuteriumTankCapacity(t *testing.T) {
	dt := newDeuteriumTank()
	assert.Equal(t, int64(10000), dt.Capacity(0))
	assert.Equal(t, int64(20000), dt.Capacity(1))
	assert.Equal(t, int64(40000), dt.Capacity(2))
	assert.Equal(t, int64(75000), dt.Capacity(3))
	assert.Equal(t, int64(140000), dt.Capacity(4))
	assert.Equal(t, int64(255000), dt.Capacity(5))
}
