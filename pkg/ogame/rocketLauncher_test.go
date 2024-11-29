package ogame

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRocketLauncherConstructionTime(t *testing.T) {
	rl := newRocketLauncher()
	assert.Equal(t, 82*time.Second, rl.ConstructionTime(1, 7, Facilities{Shipyard: 4}, LfBonuses{}, NoClass, false))
	assert.Equal(t, 164*time.Second, rl.ConstructionTime(2, 7, Facilities{Shipyard: 4}, LfBonuses{}, NoClass, false))
}

func TestRocketLauncher_GetName(t *testing.T) {
	rl := newRocketLauncher()
	assert.Equal(t, "rocket launcher", rl.GetName())
}

func TestRocketLauncher_GetRequirements(t *testing.T) {
	rl := newRocketLauncher()
	assert.Equal(t, map[ID]int64{ShipyardID: 1}, rl.GetRequirements())
}

func TestRocketLauncher_GetPrice(t *testing.T) {
	rl := newRocketLauncher()
	assert.Equal(t, Resources{Metal: 2000}, rl.GetPrice(1, LfBonuses{}))
	assert.Equal(t, Resources{Metal: 6000}, rl.GetPrice(3, LfBonuses{}))
}

func TestRocketLauncher_GetRapidfireFrom(t *testing.T) {
	rl := newRocketLauncher()
	assert.Equal(t, map[ID]int64{CruiserID: 10, BomberID: 20, DeathstarID: 200}, rl.GetRapidfireFrom())
}

func TestRocketLauncher_GetStructuralIntegrity(t *testing.T) {
	rl := newRocketLauncher()
	assert.Equal(t, rl.StructuralIntegrity, rl.GetStructuralIntegrity(Researches{ArmourTechnology: 0}))
}

func TestRocketLauncher_GetShieldPower(t *testing.T) {
	rl := newRocketLauncher()
	assert.Equal(t, rl.ShieldPower, rl.GetShieldPower(Researches{ShieldingTechnology: 0}))
}

func TestRocketLauncher_GetWeaponPower(t *testing.T) {
	rl := newRocketLauncher()
	assert.Equal(t, rl.WeaponPower, rl.GetWeaponPower(Researches{WeaponsTechnology: 0}))
}
