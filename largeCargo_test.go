package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLargeCargo_GetSpeed(t *testing.T) {
	lc := newLargeCargo()
	assert.Equal(t, 12000, lc.GetSpeed(Researches{}))
}

func TestLargeCargo_GetCargoCapacity(t *testing.T) {
	lc := newLargeCargo()
	assert.Equal(t, 35000, lc.GetCargoCapacity(Researches{HyperspaceTechnology: 8}, false))
	assert.Equal(t, 37500, lc.GetCargoCapacity(Researches{HyperspaceTechnology: 10}, false))
}
