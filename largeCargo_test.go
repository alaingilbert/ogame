package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLargeCargo_GetSpeed(t *testing.T) {
	lc := newLargeCargo()
	assert.Equal(t, 12000, lc.GetSpeed(Researches{}))
	assert.Equal(t, 29000, lc.GetCargoCapacity(Researches{HyperspaceTechnology: 8}))
	assert.Equal(t, 30000, lc.GetCargoCapacity(Researches{HyperspaceTechnology: 10}))
}
