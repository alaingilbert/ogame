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
	coord, _ := ParseCoord("[P:1:2:3]")
	assert.Equal(t, Coordinate{1, 2, 3, PlanetType}, coord)
	coord, _ = ParseCoord("[M:1:2:3]")
	assert.Equal(t, Coordinate{1, 2, 3, MoonType}, coord)
	coord, _ = ParseCoord("M:1:2:3")
	assert.Equal(t, Coordinate{1, 2, 3, MoonType}, coord)
	coord, _ = ParseCoord("1:2:3")
	assert.Equal(t, Coordinate{1, 2, 3, PlanetType}, coord)
	coord, _ = ParseCoord("1:2:3")
	assert.Equal(t, Coordinate{1, 2, 3, PlanetType}, coord)
	coord, _ = ParseCoord("D:1:2:3")
	assert.Equal(t, Coordinate{1, 2, 3, DebrisType}, coord)
	coord, _ = ParseCoord("[D:1:2:3]")
	assert.Equal(t, Coordinate{1, 2, 3, DebrisType}, coord)
	_, err := ParseCoord("[A:1:2:3]")
	assert.NotNil(t, err)
	_, err = ParseCoord("aP:1:2:3")
	assert.NotNil(t, err)
	_, err = ParseCoord("P:1234:2:3")
	assert.NotNil(t, err)
	_, err = ParseCoord("P:1:2345:3")
	assert.NotNil(t, err)
	_, err = ParseCoord("P:1:2:3456")
	assert.NotNil(t, err)
}
