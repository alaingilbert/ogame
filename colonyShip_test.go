package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColonyShipSpeed(t *testing.T) {
	cs := newColonyShip()
	speed := cs.GetSpeed(Researches{ImpulseDrive: 6})
	assert.Equal(t, 5500, speed)
	assert.Equal(t, 8700, cs.GetCargoCapacity(Researches{HyperspaceTechnology: 8}))

}
