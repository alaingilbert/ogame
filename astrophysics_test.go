package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAstrophysicsCost(t *testing.T) {
	a := NewAstrophysics()
	assert.Equal(t, Resources{Metal: 7000, Crystal: 14000, Deuterium: 7000}, a.GetPrice(2))
}
