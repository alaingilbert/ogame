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

func TestSub(t *testing.T) {
	first := Resources{Metal: 2, Crystal: 3, Deuterium: 4}
	second := Resources{Metal: 1, Crystal: 1, Deuterium: 1}
	assert.Equal(t, Resources{Metal: 1, Crystal: 2, Deuterium: 3}, first.Sub(second))
}

func TestString(t *testing.T) {
	assert.Equal(t, "[1|2|3]", Resources{Metal: 1, Crystal: 2, Deuterium: 3}.String())
	assert.Equal(t, "[1,000,000|2,000,000|3,000,000]", Resources{Metal: 1000000, Crystal: 2000000, Deuterium: 3000000}.String())
}
