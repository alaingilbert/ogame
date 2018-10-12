package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCelestial(t *testing.T) {
	assert.Equal(t, CelestialID(123), PlanetID(123).Celestial())
}
