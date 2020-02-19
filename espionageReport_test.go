package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEspionageReport_PlunderRatio(t *testing.T) {
	er := EspionageReport{}
	assert.Equal(t, 0.5, er.PlunderRatio(NoClass))

	er = EspionageReport{IsInactive: true}
	assert.Equal(t, 0.5, er.PlunderRatio(NoClass))
	assert.Equal(t, 0.75, er.PlunderRatio(Discoverer))

	er = EspionageReport{IsBandit: true}
	assert.Equal(t, 1.0, er.PlunderRatio(NoClass))
	assert.Equal(t, 1.0, er.PlunderRatio(Discoverer))

	er = EspionageReport{IsStarlord: true}
	assert.Equal(t, 0.75, er.PlunderRatio(NoClass))

	er = EspionageReport{IsStarlord: true, IsInactive: true}
	assert.Equal(t, 0.5, er.PlunderRatio(NoClass))
}

func TestEspionageReport_Loot(t *testing.T) {
	er := EspionageReport{Resources: Resources{Metal: 100}}
	assert.Equal(t, Resources{Metal: 50}, er.Loot(NoClass))
}

func TestShipsInfos(t *testing.T) {
	er := EspionageReport{HasFleet: true, SmallCargo: I64Ptr(3), LightFighter: I64Ptr(5)}
	assert.Equal(t, int64(8), er.ShipsInfos().CountShips())

	er = EspionageReport{HasFleet: false}
	var nilShipsInfos *ShipsInfos = nil
	assert.Equal(t, nilShipsInfos, er.ShipsInfos())
}
