package ogame

import "strconv"

// DefensesInfos represent a planet defenses information
type DefensesInfos struct {
	RocketLauncher         int64
	LightLaser             int64
	HeavyLaser             int64
	GaussCannon            int64
	IonCannon              int64
	PlasmaTurret           int64
	SmallShieldDome        int64
	LargeShieldDome        int64
	AntiBallisticMissiles  int64
	InterplanetaryMissiles int64
}

// AttackableValue returns the value of the defenses that can be attacked
func (d DefensesInfos) AttackableValue() int64 {
	val := d.RocketLauncher * RocketLauncher.Price.Total()
	val += d.LightLaser * LightLaser.Price.Total()
	val += d.HeavyLaser * HeavyLaser.Price.Total()
	val += d.GaussCannon * GaussCannon.Price.Total()
	val += d.IonCannon * IonCannon.Price.Total()
	val += d.PlasmaTurret * PlasmaTurret.Price.Total()
	val += d.SmallShieldDome * SmallShieldDome.Price.Total()
	val += d.LargeShieldDome * LargeShieldDome.Price.Total()
	return val
}

func (d DefensesInfos) String() string {
	return "\n" +
		"        Rocket Launcher: " + strconv.FormatInt(d.RocketLauncher, 10) + "\n" +
		"            Light Laser: " + strconv.FormatInt(d.LightLaser, 10) + "\n" +
		"            Heavy Laser: " + strconv.FormatInt(d.HeavyLaser, 10) + "\n" +
		"           Gauss Cannon: " + strconv.FormatInt(d.GaussCannon, 10) + "\n" +
		"             Ion Cannon: " + strconv.FormatInt(d.IonCannon, 10) + "\n" +
		"          Plasma Turret: " + strconv.FormatInt(d.PlasmaTurret, 10) + "\n" +
		"      Small Shield Dome: " + strconv.FormatInt(d.SmallShieldDome, 10) + "\n" +
		"      Large Shield Dome: " + strconv.FormatInt(d.LargeShieldDome, 10) + "\n" +
		"Anti Ballistic Missiles: " + strconv.FormatInt(d.AntiBallisticMissiles, 10) + "\n" +
		"Interplanetary Missiles: " + strconv.FormatInt(d.InterplanetaryMissiles, 10)
}

// ByID get number of defenses by defense id
func (d DefensesInfos) ByID(id ID) int64 {
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
func (d *DefensesInfos) Set(id ID, val int64) {
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
