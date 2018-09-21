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
