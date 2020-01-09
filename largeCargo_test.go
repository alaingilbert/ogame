package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLargeCargo_GetSpeed(t *testing.T) {
	lc := newLargeCargo()
	assert.Equal(t, int64(12000), lc.GetSpeed(Researches{}, false))
}

func TestLargeCargo_GetCargoCapacity(t *testing.T) {
	lc := newLargeCargo()
	assert.Equal(t, int64(35000), lc.GetCargoCapacity(Researches{HyperspaceTechnology: 8}, false, false))
	assert.Equal(t, int64(37500), lc.GetCargoCapacity(Researches{HyperspaceTechnology: 10}, false, false))
}
