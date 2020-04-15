package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAstrophysicsCost(t *testing.T) {
	a := newAstrophysics()
	assert.Equal(t, Resources{Metal: 7000, Crystal: 14000, Deuterium: 7000}, a.GetPrice(2))
	assert.Equal(t, Resources{Metal: 351900, Crystal: 703700, Deuterium: 351900}, a.GetPrice(9))
}
