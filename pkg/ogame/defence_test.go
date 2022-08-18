package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefence_String(t *testing.T) {
	r := DefensesInfos{
		RocketLauncher:         1,
		LightLaser:             2,
		HeavyLaser:             3,
		GaussCannon:            4,
		IonCannon:              5,
		PlasmaTurret:           6,
		SmallShieldDome:        7,
		LargeShieldDome:        8,
		AntiBallisticMissiles:  9,
		InterplanetaryMissiles: 10,
	}
	expected := "\n" +
		"        Rocket Launcher: 1\n" +
		"            Light Laser: 2\n" +
		"            Heavy Laser: 3\n" +
		"           Gauss Cannon: 4\n" +
		"             Ion Cannon: 5\n" +
		"          Plasma Turret: 6\n" +
		"      Small Shield Dome: 7\n" +
		"      Large Shield Dome: 8\n" +
		"Anti Ballistic Missiles: 9\n" +
		"Interplanetary Missiles: 10"
	assert.Equal(t, expected, r.String())
}

func TestDefence_AttackableValue(t *testing.T) {
	r := DefensesInfos{RocketLauncher: 2}
	assert.Equal(t, int64(4000), r.AttackableValue())
	r = DefensesInfos{RocketLauncher: 2, LightLaser: 4}
	assert.Equal(t, int64(12000), r.AttackableValue())
}

func TestDefenceByID(t *testing.T) {
	assert.Equal(t, int64(0), DefensesInfos{}.ByID(123456))
	assert.Equal(t, int64(2), DefensesInfos{RocketLauncher: 2}.ByID(RocketLauncherID))
	assert.Equal(t, int64(2), DefensesInfos{LightLaser: 2}.ByID(LightLaserID))
	assert.Equal(t, int64(2), DefensesInfos{HeavyLaser: 2}.ByID(HeavyLaserID))
	assert.Equal(t, int64(2), DefensesInfos{GaussCannon: 2}.ByID(GaussCannonID))
	assert.Equal(t, int64(2), DefensesInfos{IonCannon: 2}.ByID(IonCannonID))
	assert.Equal(t, int64(2), DefensesInfos{PlasmaTurret: 2}.ByID(PlasmaTurretID))
	assert.Equal(t, int64(2), DefensesInfos{SmallShieldDome: 2}.ByID(SmallShieldDomeID))
	assert.Equal(t, int64(2), DefensesInfos{LargeShieldDome: 2}.ByID(LargeShieldDomeID))
	assert.Equal(t, int64(2), DefensesInfos{AntiBallisticMissiles: 2}.ByID(AntiBallisticMissilesID))
	assert.Equal(t, int64(2), DefensesInfos{InterplanetaryMissiles: 2}.ByID(InterplanetaryMissilesID))
}

func TestDefenceSet(t *testing.T) {
	s := DefensesInfos{}
	s.Set(RocketLauncherID, 1)
	s.Set(LightLaserID, 2)
	s.Set(HeavyLaserID, 3)
	s.Set(GaussCannonID, 4)
	s.Set(IonCannonID, 5)
	s.Set(PlasmaTurretID, 6)
	s.Set(SmallShieldDomeID, 7)
	s.Set(LargeShieldDomeID, 8)
	s.Set(AntiBallisticMissilesID, 9)
	s.Set(InterplanetaryMissilesID, 10)
	assert.Equal(t, int64(1), s.ByID(RocketLauncherID))
	assert.Equal(t, int64(2), s.ByID(LightLaserID))
	assert.Equal(t, int64(3), s.ByID(HeavyLaserID))
	assert.Equal(t, int64(4), s.ByID(GaussCannonID))
	assert.Equal(t, int64(5), s.ByID(IonCannonID))
	assert.Equal(t, int64(6), s.ByID(PlasmaTurretID))
	assert.Equal(t, int64(7), s.ByID(SmallShieldDomeID))
	assert.Equal(t, int64(8), s.ByID(LargeShieldDomeID))
	assert.Equal(t, int64(9), s.ByID(AntiBallisticMissilesID))
	assert.Equal(t, int64(10), s.ByID(InterplanetaryMissilesID))
}

func TestDefence_HasMissilesDefenses(t *testing.T) {
	assert.True(t, DefensesInfos{RocketLauncher: 2, PlasmaTurret: 3, InterplanetaryMissiles: 5, AntiBallisticMissiles: 2}.HasMissilesDefense())
	assert.False(t, DefensesInfos{}.HasMissilesDefense())
	assert.False(t, DefensesInfos{InterplanetaryMissiles: 2}.HasMissilesDefense())
	assert.True(t, DefensesInfos{InterplanetaryMissiles: 2, AntiBallisticMissiles: 3}.HasMissilesDefense())
}

func TestDefence_HasShipDefenses(t *testing.T) {
	assert.True(t, DefensesInfos{RocketLauncher: 2, PlasmaTurret: 3, InterplanetaryMissiles: 5, AntiBallisticMissiles: 2}.HasShipDefense())
	assert.False(t, DefensesInfos{}.HasShipDefense())
	assert.False(t, DefensesInfos{InterplanetaryMissiles: 2}.HasShipDefense())
	assert.False(t, DefensesInfos{InterplanetaryMissiles: 2, AntiBallisticMissiles: 3}.HasShipDefense())
}

func TestDefence_CountShipDefenses(t *testing.T) {
	assert.Equal(t, int64(5), DefensesInfos{RocketLauncher: 2, PlasmaTurret: 3, AntiBallisticMissiles: 4, InterplanetaryMissiles: 5}.CountShipDefenses())
}
