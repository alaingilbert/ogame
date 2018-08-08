package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColonyShipSpeed(t *testing.T) {
	cs := newColonyShip()
	speed := cs.GetSpeed(Researches{ImpulseDrive: 6})
	assert.Equal(t, 5500, speed)
}
