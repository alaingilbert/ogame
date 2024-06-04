package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCruiser_RapidfireAgainst(t *testing.T) {
	c := newCruiser()
	assert.Equal(t, map[ID]int64{EspionageProbeID: 5, SolarSatelliteID: 5, LightFighterID: 6, RocketLauncherID: 10, CrawlerID: 5}, c.GetRapidfireAgainst())
}

func TestCruiser_GetCargoCapacity(t *testing.T) {
	c := newCruiser()
	assert.Equal(t, int64(800), c.GetCargoCapacity(Researches{HyperspaceTechnology: 0}, LfBonuses{}, NoClass, 0.05, false))
	assert.Equal(t, int64(1120), c.GetCargoCapacity(Researches{HyperspaceTechnology: 8}, LfBonuses{}, NoClass, 0.05, false))
}

func TestCruiser_GetFuelConsumption(t *testing.T) {
	c := newCruiser()
	assert.Equal(t, int64(300), c.GetFuelConsumption(Researches{}, LfBonuses{}, NoClass, 1))
}

func TestCruiser_GetPrice(t *testing.T) {
	c := newCruiser()
	assert.Equal(t, Resources{Metal: 20000, Crystal: 7000, Deuterium: 2000}, c.GetPrice(1, LfBonuses{}))
	assert.Equal(t, Resources{Metal: 60000, Crystal: 21000, Deuterium: 6000}, c.GetPrice(3, LfBonuses{}))
}
