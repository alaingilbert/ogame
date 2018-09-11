package ogame

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSmallCargoConstructionTime(t *testing.T) {
	sc := newSmallCargo()
	assert.Equal(t, 164*time.Second, sc.ConstructionTime(1, 7, Facilities{Shipyard: 4}))
	assert.Equal(t, 328*time.Second, sc.ConstructionTime(2, 7, Facilities{Shipyard: 4}))
}

func TestSmallCargoSpeed(t *testing.T) {
	sc := newSmallCargo()
	assert.Equal(t, 12000, sc.GetSpeed(Researches{CombustionDrive: 2}))
	assert.Equal(t, 16000, sc.GetSpeed(Researches{CombustionDrive: 6}))
	assert.Equal(t, 16000, sc.GetSpeed(Researches{CombustionDrive: 6, ImpulseDrive: 4}))
	assert.Equal(t, 20000, sc.GetSpeed(Researches{CombustionDrive: 6, ImpulseDrive: 5}))
	assert.Equal(t, 22000, sc.GetSpeed(Researches{CombustionDrive: 10, ImpulseDrive: 6}))
}
