package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMoonID_Celestial(t *testing.T) {
	assert.Equal(t, CelestialID(123), MoonID(123).Celestial())
}
