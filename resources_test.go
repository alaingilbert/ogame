package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValue(t *testing.T) {
	assert.Equal(t, 3, Resources{Deuterium: 1}.Value())
	assert.Equal(t, 2, Resources{Crystal: 1}.Value())
	assert.Equal(t, 1, Resources{Metal: 1}.Value())
	assert.Equal(t, 6, Resources{Deuterium: 1, Crystal: 1, Metal: 1}.Value())
}

func TestTotal(t *testing.T) {
	assert.Equal(t, 4, Resources{Deuterium: 4}.Total())
	assert.Equal(t, 2, Resources{Crystal: 2}.Total())
	assert.Equal(t, 1, Resources{Metal: 1}.Total())
	assert.Equal(t, 7, Resources{Deuterium: 1, Crystal: 2, Metal: 4}.Total())
}

func TestSub(t *testing.T) {
	first := Resources{Metal: 2, Crystal: 3, Deuterium: 4}
	second := Resources{Metal: 1, Crystal: 1, Deuterium: 1}
	assert.Equal(t, Resources{Metal: 1, Crystal: 2, Deuterium: 3}, first.Sub(second))

	assert.Equal(t, Resources{Metal: 75, Crystal: 0, Deuterium: 0}, Resources{Metal: 100, Crystal: 10, Deuterium: 0}.Sub(Resources{Metal: 25, Crystal: 40, Deuterium: 30}))
}

func TestAdd(t *testing.T) {
	first := Resources{Metal: 1, Crystal: 2, Deuterium: 4}
	second := Resources{Metal: 8, Crystal: 16, Deuterium: 32}
	assert.Equal(t, Resources{Metal: 9, Crystal: 18, Deuterium: 36}, first.Add(second))
}

func TestMul(t *testing.T) {
	first := Resources{Metal: 1, Crystal: 2, Deuterium: 4}
	assert.Equal(t, Resources{Metal: 2, Crystal: 4, Deuterium: 8}, first.Mul(2))
}

func TestGte(t *testing.T) {
	assert.True(t, Resources{Metal: 1, Crystal: 2, Deuterium: 4}.Gte(Resources{Metal: 1, Crystal: 2, Deuterium: 4}))
	assert.True(t, Resources{Metal: 1, Crystal: 2, Deuterium: 4}.Gte(Resources{Metal: 1, Crystal: 1, Deuterium: 1}))
	assert.False(t, Resources{Metal: 1, Crystal: 2, Deuterium: 4}.Gte(Resources{Metal: 1, Crystal: 2, Deuterium: 5}))
	assert.False(t, Resources{Metal: 1, Crystal: 2, Deuterium: 4}.Gte(Resources{Metal: 1, Crystal: 3, Deuterium: 4}))
	assert.False(t, Resources{Metal: 1, Crystal: 2, Deuterium: 4}.Gte(Resources{Metal: 2, Crystal: 2, Deuterium: 4}))
}

func TestLte(t *testing.T) {
	assert.True(t, Resources{Metal: 1, Crystal: 2, Deuterium: 4}.Lte(Resources{Metal: 1, Crystal: 2, Deuterium: 4}))
	assert.False(t, Resources{Metal: 1, Crystal: 2, Deuterium: 4}.Lte(Resources{Metal: 1, Crystal: 1, Deuterium: 1}))
	assert.True(t, Resources{Metal: 1, Crystal: 2, Deuterium: 4}.Lte(Resources{Metal: 1, Crystal: 2, Deuterium: 5}))
	assert.True(t, Resources{Metal: 1, Crystal: 2, Deuterium: 4}.Lte(Resources{Metal: 1, Crystal: 3, Deuterium: 4}))
	assert.True(t, Resources{Metal: 1, Crystal: 2, Deuterium: 4}.Lte(Resources{Metal: 2, Crystal: 2, Deuterium: 4}))
}

func TestCanAford(t *testing.T) {
	assert.True(t, Resources{Metal: 1, Crystal: 2, Deuterium: 4}.CanAfford(Resources{Metal: 1, Crystal: 2, Deuterium: 4}))
	assert.True(t, Resources{Metal: 1, Crystal: 2, Deuterium: 4}.CanAfford(Resources{Metal: 1, Crystal: 1, Deuterium: 1}))
	assert.False(t, Resources{Metal: 1, Crystal: 2, Deuterium: 4}.CanAfford(Resources{Metal: 1, Crystal: 2, Deuterium: 5}))
	assert.False(t, Resources{Metal: 1, Crystal: 2, Deuterium: 4}.CanAfford(Resources{Metal: 1, Crystal: 3, Deuterium: 4}))
	assert.False(t, Resources{Metal: 1, Crystal: 2, Deuterium: 4}.CanAfford(Resources{Metal: 2, Crystal: 2, Deuterium: 4}))
}

func TestString(t *testing.T) {
	assert.Equal(t, "[1|2|3]", Resources{Metal: 1, Crystal: 2, Deuterium: 3}.String())
	assert.Equal(t, "[1,000,000|2,000,000|3,000,000]", Resources{Metal: 1000000, Crystal: 2000000, Deuterium: 3000000}.String())
}
