package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSystemInfos_Position(t *testing.T) {
	si := SystemInfos{}
	var nilPlanetInfo *PlanetInfos
	assert.Equal(t, nilPlanetInfo, si.Position(0))
	assert.Equal(t, nilPlanetInfo, si.Position(16))
}
