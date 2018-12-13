package ogame

import "strconv"

// DefensesInfos represent a planet defenses information
type DefensesInfos struct {
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

func (d DefensesInfos) String() string {
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

// ByID get number of defenses by defense id
func (d DefensesInfos) ByID(id ID) int {
	switch id {
	case RocketLauncherID:
		return d.RocketLauncher
	case LightLaserID:
		return d.LightLaser
	case HeavyLaserID:
		return d.HeavyLaser
	case GaussCannonID:
		return d.GaussCannon
	case IonCannonID:
		return d.IonCannon
	case PlasmaTurretID:
		return d.PlasmaTurret
	case SmallShieldDomeID:
		return d.SmallShieldDome
	case LargeShieldDomeID:
		return d.LargeShieldDome
	case AntiBallisticMissilesID:
		return d.AntiBallisticMissiles
	case InterplanetaryMissilesID:
		return d.InterplanetaryMissiles
	default:
		return 0
	}
}

// Set sets the defenses value using the defense id
func (d *DefensesInfos) Set(id ID, val int) {
	switch id {
	case RocketLauncherID:
		d.RocketLauncher = val
	case LightLaserID:
		d.LightLaser = val
	case HeavyLaserID:
		d.HeavyLaser = val
	case GaussCannonID:
		d.GaussCannon = val
	case IonCannonID:
		d.IonCannon = val
	case PlasmaTurretID:
		d.PlasmaTurret = val
	case SmallShieldDomeID:
		d.SmallShieldDome = val
	case LargeShieldDomeID:
		d.LargeShieldDome = val
	case AntiBallisticMissilesID:
		d.AntiBallisticMissiles = val
	case InterplanetaryMissilesID:
		d.InterplanetaryMissiles = val
	}
}
