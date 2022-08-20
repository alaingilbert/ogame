package ogame

import (
	"testing"

	"github.com/google/gxui/math"
	"github.com/stretchr/testify/assert"
)

func TestValue(t *testing.T) {
	assert.Equal(t, int64(3), Resources{Deuterium: 1}.Value())
	assert.Equal(t, int64(2), Resources{Crystal: 1}.Value())
	assert.Equal(t, int64(1), Resources{Metal: 1}.Value())
	assert.Equal(t, int64(6), Resources{Deuterium: 1, Crystal: 1, Metal: 1}.Value())
}

func TestTotal(t *testing.T) {
	assert.Equal(t, int64(4), Resources{Deuterium: 4}.Total())
	assert.Equal(t, int64(2), Resources{Crystal: 2}.Total())
	assert.Equal(t, int64(1), Resources{Metal: 1}.Total())
	assert.Equal(t, int64(7), Resources{Deuterium: 1, Crystal: 2, Metal: 4}.Total())
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

func TestDiv(t *testing.T) {
	first := Resources{Metal: 100, Crystal: 200, Deuterium: 400}
	assert.Equal(t, int64(1), first.Div(Resources{Metal: 100, Crystal: 2, Deuterium: 4}))
	assert.Equal(t, int64(1), first.Div(Resources{Metal: 1, Crystal: 200, Deuterium: 4}))
	assert.Equal(t, int64(1), first.Div(Resources{Metal: 1, Crystal: 2, Deuterium: 400}))
	assert.Equal(t, int64(100), first.Div(Resources{Metal: 1, Crystal: 2, Deuterium: 4}))
	assert.Equal(t, int64(100), first.Div(Resources{Metal: 0, Crystal: 2, Deuterium: 4}))
	assert.Equal(t, int64(100), first.Div(Resources{Metal: 1, Crystal: 0, Deuterium: 4}))
	assert.Equal(t, int64(100), first.Div(Resources{Metal: 1, Crystal: 2, Deuterium: 0}))
	assert.Equal(t, int64(math.MaxInt), first.Div(Resources{Metal: 0, Crystal: 0, Deuterium: 0}))
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

func TestResourcesDetails_Available(t *testing.T) {
	d := ResourcesDetails{}
	d.Metal.Available = 1
	d.Crystal.Available = 2
	d.Deuterium.Available = 3
	d.Energy.Available = 4
	d.Darkmatter.Available = 5
	assert.Equal(t, Resources{1, 2, 3, 4, 5}, d.Available())
}

func TestResources_FitsIn(t *testing.T) {
	assert.Equal(t, int64(1), Resources{Metal: 100, Crystal: 200, Deuterium: 300}.FitsIn(SmallCargo, Researches{}, false, false, false))
	assert.Equal(t, int64(2), Resources{Metal: 1001, Crystal: 2000, Deuterium: 2000}.FitsIn(SmallCargo, Researches{}, false, false, false))
	assert.Equal(t, int64(2), Resources{Metal: 1000, Crystal: 2000, Deuterium: 3000}.FitsIn(SmallCargo, Researches{}, false, false, false))
	assert.Equal(t, int64(2), Resources{Metal: 999, Crystal: 4000, Deuterium: 5000}.FitsIn(SmallCargo, Researches{}, false, false, false))
	assert.Equal(t, int64(0), Resources{Metal: 0, Crystal: 0, Deuterium: 0}.FitsIn(SmallCargo, Researches{}, false, false, false))
	assert.Equal(t, int64(0), Resources{Metal: 100, Crystal: 200, Deuterium: 300}.FitsIn(EspionageProbe, Researches{}, false, false, false))
	assert.Equal(t, int64(120), Resources{Metal: 100, Crystal: 200, Deuterium: 300}.FitsIn(EspionageProbe, Researches{}, true, false, false))
}
