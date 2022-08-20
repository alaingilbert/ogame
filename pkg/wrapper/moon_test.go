package wrapper

import (
	"github.com/alaingilbert/ogame/pkg/ogame"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMoonID_Celestial(t *testing.T) {
	assert.Equal(t, ogame.CelestialID(123), ogame.MoonID(123).Celestial())
}
