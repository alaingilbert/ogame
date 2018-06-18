package bomber

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/ships/baseShip"
)

// Bomber ...
type Bomber struct {
	baseShip.BaseShip
}

// New ...
func New() *Bomber {
	b := new(Bomber)
	b.OGameID = 211
	b.StructuralIntegrity = 75000
	b.ShieldPower = 500
	b.WeaponPower = 1000
	b.CargoCapacity = 500
	b.BaseSpeed = 4000
	b.FuelConsumption = 1000
	b.RapidfireFrom = map[ogame.ID]int{ogame.Deathstar: 25}
	b.RapidfireAgainst = map[ogame.ID]int{ogame.EspionageProbe: 5, ogame.SolarSatellite: 5,
		ogame.RocketLauncher: 20, ogame.LightLaser: 20, ogame.HeavyLaser: 10, ogame.IonCannon: 10}
	b.Price = ogame.Resources{Metal: 50000, Crystal: 25000, Deuterium: 15000}
	b.Requirements = map[ogame.ID]int{ogame.ImpulseDrive: 6, ogame.Shipyard: 8, ogame.PlasmaTechnology: 5}
	return b
}
