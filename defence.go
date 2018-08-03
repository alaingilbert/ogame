package ogame

import "strconv"

// Defense ...
type Defenses struct {
	RocketLauncher         int
	LightLaser             int
	HeavyLaser             int
	GaussCannon            int
	IonCannon              int
	PlasmaTurret           int
	SmallShieldDome        int
	LargeShieldDome        int
	AntiBallisticMissiles  int
	InterplanetaryMissiles int
}

func (d Defenses) String() string {
	return "\n" +
		"        Rocket Launcher: " + strconv.Itoa(d.RocketLauncher) + "\n" +
		"            Light Laser: " + strconv.Itoa(d.LightLaser) + "\n" +
		"            Heavy Laser: " + strconv.Itoa(d.HeavyLaser) + "\n" +
		"           Gauss Cannon: " + strconv.Itoa(d.GaussCannon) + "\n" +
		"             Ion Cannon: " + strconv.Itoa(d.IonCannon) + "\n" +
		"          Plasma Turret: " + strconv.Itoa(d.PlasmaTurret) + "\n" +
		"      Small Shield Dome: " + strconv.Itoa(d.SmallShieldDome) + "\n" +
		"      Large Shield Dome: " + strconv.Itoa(d.LargeShieldDome) + "\n" +
		"Anti Ballistic Missiles: " + strconv.Itoa(d.AntiBallisticMissiles) + "\n" +
		"Interplanetary Missiles: " + strconv.Itoa(d.InterplanetaryMissiles)
}
