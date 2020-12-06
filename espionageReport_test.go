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

func TestEspionageReport_IsDefenceless(t *testing.T) {
	two := int64(2)
	assert.True(t, EspionageReport{Resources: Resources{Metal: 100}, HasFleetInformation: true, HasDefensesInformation: true}.IsDefenceless())
	assert.False(t, EspionageReport{Resources: Resources{Metal: 100}, HasFleetInformation: true, HasDefensesInformation: true, LightFighter: &two}.IsDefenceless())
	assert.False(t, EspionageReport{Resources: Resources{Metal: 100}, HasFleetInformation: true, HasDefensesInformation: true, RocketLauncher: &two}.IsDefenceless())
	assert.False(t, EspionageReport{Resources: Resources{Metal: 100}, HasFleetInformation: true, HasDefensesInformation: false}.IsDefenceless())
	assert.False(t, EspionageReport{Resources: Resources{Metal: 100}, HasFleetInformation: false, HasDefensesInformation: true}.IsDefenceless())
	assert.False(t, EspionageReport{Resources: Resources{Metal: 100}, HasFleetInformation: false, HasDefensesInformation: false}.IsDefenceless())
}

func TestShipsInfos(t *testing.T) {
	er := EspionageReport{HasFleetInformation: true, SmallCargo: I64Ptr(3), LightFighter: I64Ptr(5)}
	assert.Equal(t, int64(8), er.ShipsInfos().CountShips())

	er = EspionageReport{HasFleetInformation: false}
	var nilShipsInfos *ShipsInfos = nil
	assert.Equal(t, nilShipsInfos, er.ShipsInfos())
}
