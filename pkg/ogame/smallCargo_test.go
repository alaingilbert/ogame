package ogame

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSmallCargoConstructionTime(t *testing.T) {
	sc := newSmallCargo()
	assert.Equal(t, 164*time.Second, sc.ConstructionTime(1, 7, Facilities{Shipyard: 4}, LfBonuses{}, NoClass, false))
	assert.Equal(t, 328*time.Second, sc.ConstructionTime(2, 7, Facilities{Shipyard: 4}, LfBonuses{}, NoClass, false))
}

func TestSmallCargoSpeed(t *testing.T) {
	sc := newSmallCargo()
	assert.Equal(t, int64(6000), sc.GetSpeed(Researches{CombustionDrive: 2}, LfBonuses{}, NoClass))
	assert.Equal(t, int64(8000), sc.GetSpeed(Researches{CombustionDrive: 6}, LfBonuses{}, NoClass))
	assert.Equal(t, int64(8000), sc.GetSpeed(Researches{CombustionDrive: 6, ImpulseDrive: 4}, LfBonuses{}, NoClass))
	assert.Equal(t, int64(20000), sc.GetSpeed(Researches{CombustionDrive: 6, ImpulseDrive: 5}, LfBonuses{}, NoClass))
	assert.Equal(t, int64(22000), sc.GetSpeed(Researches{CombustionDrive: 10, ImpulseDrive: 6}, LfBonuses{}, NoClass))
	lfBonuses := LfBonuses{LfShipBonuses: make(LfShipBonuses)}
	lfBonuses.LfShipBonuses[SmallCargoID] = LfShipBonus{Speed: 0.836775}
	assert.Equal(t, int64(48368), sc.GetSpeed(Researches{CombustionDrive: 15, ImpulseDrive: 15}, lfBonuses, Discoverer))
}

func TestSmallCargoFuelConsumption(t *testing.T) {
	sc := newSmallCargo()
	assert.Equal(t, int64(10), sc.GetFuelConsumption(Researches{}, LfBonuses{}, NoClass, 1))
	assert.Equal(t, int64(20), sc.GetFuelConsumption(Researches{ImpulseDrive: 5}, LfBonuses{}, NoClass, 1))
}
