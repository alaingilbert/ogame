package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCruiser_RapidfireAgainst(t *testing.T) {
	c := newCruiser()
	assert.Equal(t, map[ID]int{EspionageProbeID: 5, SolarSatelliteID: 5, LightFighterID: 6, RocketLauncherID: 10}, c.GetRapidfireAgainst())
}

func TestCruiser_GetCargoCapacity(t *testing.T) {
	c := newCruiser()
	assert.Equal(t, 800, c.GetCargoCapacity())
}

func TestCruiser_GetFuelConsumption(t *testing.T) {
	c := newCruiser()
	assert.Equal(t, 300, c.GetFuelConsumption())
}
