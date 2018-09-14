package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestByID(t *testing.T) {
	assert.Equal(t, 0, ShipsInfos{}.ByID(123456))
}

func TestSet(t *testing.T) {
	s := ShipsInfos{}
	s.Set(BattleshipID, 1)
	s.Set(DeathstarID, 2)
	s.Set(SolarSatelliteID, 4)
	assert.Equal(t, 1, s.ByID(BattleshipID))
	assert.Equal(t, 2, s.ByID(DeathstarID))
	assert.Equal(t, 4, s.ByID(SolarSatelliteID))
}
