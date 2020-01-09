package ogame

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSmallCargoConstructionTime(t *testing.T) {
	sc := newSmallCargo()
	assert.Equal(t, 164*time.Second, sc.ConstructionTime(1, 7, Facilities{Shipyard: 4}, false, false))
	assert.Equal(t, 328*time.Second, sc.ConstructionTime(2, 7, Facilities{Shipyard: 4}, false, false))
}

func TestSmallCargoSpeed(t *testing.T) {
	sc := newSmallCargo()
	assert.Equal(t, int64(6000), sc.GetSpeed(Researches{CombustionDrive: 2}, false, false))
	assert.Equal(t, int64(8000), sc.GetSpeed(Researches{CombustionDrive: 6}, false, false))
	assert.Equal(t, int64(8000), sc.GetSpeed(Researches{CombustionDrive: 6, ImpulseDrive: 4}, false, false))
	assert.Equal(t, int64(20000), sc.GetSpeed(Researches{CombustionDrive: 6, ImpulseDrive: 5}, false, false))
	assert.Equal(t, int64(22000), sc.GetSpeed(Researches{CombustionDrive: 10, ImpulseDrive: 6}, false, false))
}
