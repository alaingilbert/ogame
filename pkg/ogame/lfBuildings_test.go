package ogame

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestResidentialSectorCost(t *testing.T) {
	a := newResidentialSector()
	assert.Equal(t, Resources{Metal: 7, Crystal: 2}, a.GetPrice(1, LfBonuses{}))
	assert.Equal(t, Resources{Metal: 16, Crystal: 4}, a.GetPrice(2, LfBonuses{}))
	assert.Equal(t, Resources{Metal: 30, Crystal: 8}, a.GetPrice(3, LfBonuses{}))
	assert.Equal(t, Resources{Metal: 120594, Crystal: 34455}, a.GetPrice(35, LfBonuses{}))
}

func TestBiosphereFarmCost(t *testing.T) {
	a := newBiosphereFarm()
	assert.Equal(t, Resources{Metal: 5, Crystal: 2, Energy: 8}, a.GetPrice(1, LfBonuses{}))
	assert.Equal(t, Resources{Metal: 8499, Crystal: 3399, Energy: 272}, a.GetPrice(22, LfBonuses{}))
}

func TestResearchCenterCost(t *testing.T) {
	a := newResearchCentre()
	assert.Equal(t, Resources{Metal: 20000, Crystal: 25000, Deuterium: 10000}, a.GetPrice(1, LfBonuses{}))
	assert.Equal(t, Resources{Metal: 52000, Crystal: 65000, Deuterium: 26000}, a.GetPrice(2, LfBonuses{}))
}

func TestResidentialSectorConstructionTime(t *testing.T) {
	// https://proxyforgame.com/en/ogame/calc/lfcosts.php
	rs := newResidentialSector()
	assert.Equal(t, (25*60+36)*time.Second, rs.BuildingConstructionTime(23, 8, Facilities{RoboticsFactory: 5}, LfBonuses{}))
	assert.Equal(t, (25*60+36)*time.Second, rs.ConstructionTime(23, 8, Facilities{RoboticsFactory: 5}, LfBonuses{}, NoClass, false))
}

func TestResearchCentreConstructionTime(t *testing.T) {
	rc := newResearchCentre()
	assert.Equal(t, (17*60+21)*time.Second, rc.BuildingConstructionTime(2, 8, Facilities{RoboticsFactory: 5}, LfBonuses{}))
}

func TestAssemblyLineConstructionTime(t *testing.T) {
	al := newAssemblyLine()
	assert.Equal(t, (10*60+32)*time.Second, al.BuildingConstructionTime(42, 8, Facilities{RoboticsFactory: 10, NaniteFactory: 7}, LfBonuses{}))
}

func TestAntimatterCondenserCost(t *testing.T) {
	a := newAntimatterCondenser()
	testCases := []struct {
		level    int64
		expected Resources
	}{
		{1, Resources{Metal: 6, Crystal: 3, Energy: 9}},
		{2, Resources{Metal: 14, Crystal: 7, Energy: 18}},
		{3, Resources{Metal: 25, Crystal: 12, Energy: 28}},
		{4, Resources{Metal: 41, Crystal: 20, Energy: 38}},
		//{5, Resources{Metal: 62, Crystal: 31, Energy: 49}},
		//{6, Resources{Metal: 89, Crystal: 44, Energy: 60}},
		//{7, Resources{Metal: 125, Crystal: 62, Energy: 72}},
		//{8, Resources{Metal: 171, Crystal: 85, Energy: 84}},
		//{28, Resources{Metal: 23_078, Crystal: 11_539, Energy: 438}},
		//{70, Resources{Metal: 122_111_134, Crystal: 61_055_567, Energy: 2519}},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Level %d", tc.level), func(t *testing.T) {
			assert.Equal(t, tc.expected, a.GetPrice(tc.level, LfBonuses{}))
		})
	}
}
