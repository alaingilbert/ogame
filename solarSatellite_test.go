package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSolarSatelliteSpeed(t *testing.T) {
	ss := newSolarSatellite()
	assert.Equal(t, 0, ss.GetSpeed(Researches{CombustionDrive: 10, ImpulseDrive: 6}))
}

func TestSolarSatellite_GetLevel(t *testing.T) {
	ss := newSolarSatellite()
	assert.Equal(t, 0, ss.GetLevel(ResourcesBuildings{SolarSatellite: 10}, Facilities{}, Researches{}))
}

func TestSolarSatellite_Production(t *testing.T) {
	ss := newSolarSatellite()
	assert.Equal(t, 1326, ss.Production(Temperature{-23, 17}, 51))
}
