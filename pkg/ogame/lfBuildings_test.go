package ogame

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestResidentialSectorCost(t *testing.T) {
	a := newResidentialSector()
	assert.Equal(t, Resources{Metal: 7, Crystal: 2}, a.GetPrice(1))
	assert.Equal(t, Resources{Metal: 16, Crystal: 4}, a.GetPrice(2))
	assert.Equal(t, Resources{Metal: 30, Crystal: 8}, a.GetPrice(3))
	assert.Equal(t, Resources{Metal: 120594, Crystal: 34455}, a.GetPrice(35))
}

func TestBiosphereFarmCost(t *testing.T) {
	a := newBiosphereFarm()
	assert.Equal(t, Resources{Metal: 5, Crystal: 2, Energy: 8}, a.GetPrice(1))
	assert.Equal(t, Resources{Metal: 8499, Crystal: 3399, Energy: 272}, a.GetPrice(22))
}

func TestResearchCenterCost(t *testing.T) {
	a := newResearchCentre()
	assert.Equal(t, Resources{Metal: 20000, Crystal: 25000, Deuterium: 10000}, a.GetPrice(1))
	assert.Equal(t, Resources{Metal: 52000, Crystal: 65000, Deuterium: 26000}, a.GetPrice(2))
}

func TestResidentialSectorConstructionTime(t *testing.T) {
	// https://proxyforgame.com/en/ogame/calc/lfcosts.php
	rs := newResidentialSector()
	assert.Equal(t, (25*60+36)*time.Second, rs.BuildingConstructionTime(23, 8, Facilities{RoboticsFactory: 5}))
	assert.Equal(t, (25*60+36)*time.Second, rs.ConstructionTime(23, 8, Facilities{RoboticsFactory: 5}, false, false))
}

func TestResearchCentreConstructionTime(t *testing.T) {
	rc := newResearchCentre()
	assert.Equal(t, (17*60+21)*time.Second, rc.BuildingConstructionTime(2, 8, Facilities{RoboticsFactory: 5}))
}
