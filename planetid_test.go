package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCelestial(t *testing.T) {
	assert.Equal(t, CelestialID(123), PlanetID(123).Celestial())
}

func TestPlantID_String(t *testing.T) {
	assert.Equal(t, "123", PlanetID(123).String())
}
