package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColonyShip_GetSpeed(t *testing.T) {
	cs := newColonyShip()
	speed := cs.GetSpeed(Researches{ImpulseDrive: 6}, false, false)
	assert.Equal(t, int64(5500), speed)

}

func TestColony_GetCargoCapacity(t *testing.T) {
	cs := newColonyShip()
	assert.Equal(t, int64(10500), cs.GetCargoCapacity(Researches{HyperspaceTechnology: 8}, false, false, false))

}
