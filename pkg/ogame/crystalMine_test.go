package ogame

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCrystalMineConstructionTime(t *testing.T) {
	cm := newCrystalMine()
	assert.Equal(t, 75*time.Second, cm.ConstructionTime(5, 6, Facilities{}, LfBonuses{}, NoClass, false))
}

func TestCrystalMine_EnergyConsumption(t *testing.T) {
	cm := newCrystalMine()
	assert.Equal(t, int64(736), cm.EnergyConsumption(16))
}

func TestCrystalMine_Production(t *testing.T) {
	cm := newCrystalMine()
	assert.Equal(t, int64(37921+1752+105), cm.Production(7, 1, 1, 7, 25))
}
