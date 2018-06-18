package colonyShip

import (
	"testing"

	"github.com/alaingilbert/ogame"
	"github.com/stretchr/testify/assert"
)

func TestSpeed(t *testing.T) {
	cs := New()
	speed := cs.GetSpeed(ogame.Researches{ImpulseDrive: 6})
	assert.Equal(t, 5500, speed)
}
