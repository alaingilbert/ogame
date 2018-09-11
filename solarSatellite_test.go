package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSolarSatelliteSpeed(t *testing.T) {
	ss := newSolarSatellite()
	assert.Equal(t, 0, ss.GetSpeed(Researches{CombustionDrive: 10, ImpulseDrive: 6}))
}
