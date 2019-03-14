package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColonyShip_GetSpeed(t *testing.T) {
	cs := newColonyShip()
	speed := cs.GetSpeed(Researches{ImpulseDrive: 6})
	assert.Equal(t, 5500, speed)

}

func TestColony_GetCargoCapacity(t *testing.T) {
	cs := newColonyShip()
	assert.Equal(t, 10500, cs.GetCargoCapacity(Researches{HyperspaceTechnology: 8}))

}
