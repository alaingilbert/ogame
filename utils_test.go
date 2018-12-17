package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseInt(t *testing.T) {
	assert.Equal(t, 1234567890, ParseInt("1.234.567.890"))
}

func TestToInt(t *testing.T) {
	assert.Equal(t, 1234567890, toInt([]byte("1234567890")))
}

func TestParseCoord(t *testing.T) {
	coord, err := ParseCoord("[P:1:2:3]")
	assert.Equal(t, Coordinate{1, 2, 3, PlanetType}, coord)
	coord, err = ParseCoord("[M:1:2:3]")
	assert.Equal(t, Coordinate{1, 2, 3, MoonType}, coord)
	coord, err = ParseCoord("M:1:2:3")
	assert.Equal(t, Coordinate{1, 2, 3, MoonType}, coord)
	coord, err = ParseCoord("1:2:3")
	assert.Equal(t, Coordinate{1, 2, 3, PlanetType}, coord)
	coord, err = ParseCoord("1:2:3")
	assert.Equal(t, Coordinate{1, 2, 3, PlanetType}, coord)
	coord, err = ParseCoord("D:1:2:3")
	assert.Equal(t, Coordinate{1, 2, 3, DebrisType}, coord)
	coord, err = ParseCoord("[D:1:2:3]")
	assert.Equal(t, Coordinate{1, 2, 3, DebrisType}, coord)
	coord, err = ParseCoord("[A:1:2:3]")
	assert.NotNil(t, err)
	coord, err = ParseCoord("aP:1:2:3")
	assert.NotNil(t, err)
	coord, err = ParseCoord("P:1234:2:3")
	assert.NotNil(t, err)
	coord, err = ParseCoord("P:1:2345:3")
	assert.NotNil(t, err)
	coord, err = ParseCoord("P:1:2:3456")
	assert.NotNil(t, err)
}
