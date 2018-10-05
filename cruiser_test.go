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

func TestCruiser_GetPrice(t *testing.T) {
	c := newCruiser()
	assert.Equal(t, Resources{Metal: 20000, Crystal: 7000, Deuterium: 2000}, c.GetPrice(1))
	assert.Equal(t, Resources{Metal: 60000, Crystal: 21000, Deuterium: 6000}, c.GetPrice(3))
}
