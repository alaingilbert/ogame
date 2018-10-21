package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoordinate_String(t *testing.T) {
	assert.Equal(t, "[1:2:3]", Coordinate{1, 2, 3, PlanetType}.String())
}

func TestCoordinate_Equal(t *testing.T) {
	assert.True(t, Coordinate{1, 2, 3, PlanetType}.Equal(Coordinate{1, 2, 3, PlanetType}))
	assert.False(t, Coordinate{1, 2, 3, PlanetType}.Equal(Coordinate{2, 2, 3, PlanetType}))
	assert.False(t, Coordinate{1, 2, 3, PlanetType}.Equal(Coordinate{1, 3, 3, PlanetType}))
	assert.False(t, Coordinate{1, 2, 3, PlanetType}.Equal(Coordinate{1, 2, 4, PlanetType}))
}
