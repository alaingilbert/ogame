package wrapper

import (
	"github.com/alaingilbert/ogame/pkg/ogame"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFields_HasFieldAvailable(t *testing.T) {
	assert.True(t, ogame.Fields{Built: 10, Total: 11}.HasFieldAvailable())
	assert.False(t, ogame.Fields{Built: 11, Total: 11}.HasFieldAvailable())
}

func TestTemperature_Mean(t *testing.T) {
	assert.Equal(t, int64(5), ogame.Temperature{Min: 0, Max: 10}.Mean())
	assert.Equal(t, int64(0), ogame.Temperature{Min: -10, Max: 10}.Mean())
}

func TestPlanet_String(t *testing.T) {
	assert.Equal(t, "Earth [P:1:1:3]", Planet{Planet: ogame.Planet{Name: "Earth", Coordinate: ogame.Coordinate{Galaxy: 1, System: 1, Position: 3, Type: ogame.PlanetType}}}.String())
}
