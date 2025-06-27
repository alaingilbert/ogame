package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoordinate_String(t *testing.T) {
	assert.Equal(t, "[P:1:2:3]", Coordinate{1, 2, 3, PlanetType}.String())
	assert.Equal(t, "[M:1:2:3]", Coordinate{1, 2, 3, MoonType}.String())
	assert.Equal(t, "[D:1:2:3]", Coordinate{1, 2, 3, DebrisType}.String())
}

func TestCoordinate_Equal(t *testing.T) {
	assert.True(t, Coordinate{1, 2, 3, PlanetType}.Equal(Coordinate{1, 2, 3, PlanetType}))
	assert.False(t, Coordinate{1, 2, 3, PlanetType}.Equal(Coordinate{2, 2, 3, PlanetType}))
	assert.False(t, Coordinate{1, 2, 3, PlanetType}.Equal(Coordinate{1, 3, 3, PlanetType}))
	assert.False(t, Coordinate{1, 2, 3, PlanetType}.Equal(Coordinate{1, 2, 4, PlanetType}))
	assert.False(t, Coordinate{1, 2, 3, PlanetType}.Equal(Coordinate{1, 2, 3, MoonType}))
	assert.False(t, Coordinate{1, 2, 3, PlanetType} == Coordinate{1, 2, 3, MoonType})
	assert.True(t, Coordinate{1, 2, 3, PlanetType} == Coordinate{1, 2, 3, PlanetType})
	assert.True(t, Coordinate{1, 2, 3, MoonType} == Coordinate{1, 2, 3, MoonType})
	assert.False(t, Coordinate{1, 2, 3, MoonType} == Coordinate{1, 2, 4, MoonType})
}

func TestCoordinate_IsPlanet(t *testing.T) {
	assert.True(t, Coordinate{1, 2, 3, PlanetType}.IsPlanet())
	assert.False(t, Coordinate{1, 2, 3, MoonType}.IsPlanet())
	assert.False(t, Coordinate{1, 2, 3, DebrisType}.IsPlanet())
}

func TestCoordinate_IsMoon(t *testing.T) {
	assert.False(t, Coordinate{1, 2, 3, PlanetType}.IsMoon())
	assert.True(t, Coordinate{1, 2, 3, MoonType}.IsMoon())
	assert.False(t, Coordinate{1, 2, 3, DebrisType}.IsMoon())
}

func TestCoordinate_IsDebris(t *testing.T) {
	assert.False(t, Coordinate{1, 2, 3, PlanetType}.IsDebris())
	assert.False(t, Coordinate{1, 2, 3, MoonType}.IsDebris())
	assert.True(t, Coordinate{1, 2, 3, DebrisType}.IsDebris())
}

func TestCoordinate_Planet(t *testing.T) {
	assert.Equal(t, Coordinate{1, 2, 3, PlanetType}.Planet(), Coordinate{1, 2, 3, PlanetType})
	assert.Equal(t, Coordinate{1, 2, 3, MoonType}.Planet(), Coordinate{1, 2, 3, PlanetType})
	assert.Equal(t, Coordinate{1, 2, 3, DebrisType}.Planet(), Coordinate{1, 2, 3, PlanetType})
}

func TestCoordinate_Moon(t *testing.T) {
	assert.Equal(t, Coordinate{1, 2, 3, PlanetType}.Moon(), Coordinate{1, 2, 3, MoonType})
	assert.Equal(t, Coordinate{1, 2, 3, MoonType}.Moon(), Coordinate{1, 2, 3, MoonType})
	assert.Equal(t, Coordinate{1, 2, 3, DebrisType}.Moon(), Coordinate{1, 2, 3, MoonType})
}

func TestCoordinate_Debris(t *testing.T) {
	assert.Equal(t, Coordinate{1, 2, 3, PlanetType}.Debris(), Coordinate{1, 2, 3, DebrisType})
	assert.Equal(t, Coordinate{1, 2, 3, MoonType}.Debris(), Coordinate{1, 2, 3, DebrisType})
	assert.Equal(t, Coordinate{1, 2, 3, DebrisType}.Debris(), Coordinate{1, 2, 3, DebrisType})
}

func TestCoordinate_Cmp(t *testing.T) {
	mustParseCoord := MustParseCoord
	assert.Equal(t, 0, mustParseCoord("2:2:2").Cmp(mustParseCoord("2:2:2")))
	assert.Equal(t, -1, mustParseCoord("2:2:1").Cmp(mustParseCoord("2:2:2")))
	assert.Equal(t, 1, mustParseCoord("2:2:2").Cmp(mustParseCoord("2:2:1")))
	assert.Equal(t, -1, mustParseCoord("2:1:2").Cmp(mustParseCoord("2:2:2")))
	assert.Equal(t, 1, mustParseCoord("2:2:2").Cmp(mustParseCoord("2:1:2")))
	assert.Equal(t, -1, mustParseCoord("1:2:2").Cmp(mustParseCoord("2:2:2")))
	assert.Equal(t, 1, mustParseCoord("2:2:2").Cmp(mustParseCoord("1:2:2")))
	assert.Equal(t, -1, mustParseCoord("2:2:2").Cmp(mustParseCoord("M:2:2:2")))
	assert.Equal(t, 1, mustParseCoord("M:2:2:2").Cmp(mustParseCoord("2:2:2")))
	assert.Equal(t, 0, mustParseCoord("M:2:2:2").Cmp(mustParseCoord("M:2:2:2")))
}
