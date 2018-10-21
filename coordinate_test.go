package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoordinate_String(t *testing.T) {
	assert.Equal(t, "[1:2:3]", Coordinate{1, 2, 3, PlanetDest}.String())
}

func TestCoordinate_Equal(t *testing.T) {
	assert.True(t, Coordinate{1, 2, 3, PlanetDest}.Equal(Coordinate{1, 2, 3, PlanetDest}))
	assert.False(t, Coordinate{1, 2, 3, PlanetDest}.Equal(Coordinate{2, 2, 3, PlanetDest}))
	assert.False(t, Coordinate{1, 2, 3, PlanetDest}.Equal(Coordinate{1, 3, 3, PlanetDest}))
	assert.False(t, Coordinate{1, 2, 3, PlanetDest}.Equal(Coordinate{1, 2, 4, PlanetDest}))
}
