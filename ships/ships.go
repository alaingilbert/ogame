package ships

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/ships/battlecruiser"
	"github.com/alaingilbert/ogame/ships/battleship"
	"github.com/alaingilbert/ogame/ships/bomber"
	"github.com/alaingilbert/ogame/ships/colonyShip"
	"github.com/alaingilbert/ogame/ships/cruiser"
	"github.com/alaingilbert/ogame/ships/deathstar"
	"github.com/alaingilbert/ogame/ships/destroyer"
	"github.com/alaingilbert/ogame/ships/espionageProbe"
	"github.com/alaingilbert/ogame/ships/heavyFighter"
	"github.com/alaingilbert/ogame/ships/largeCargo"
	"github.com/alaingilbert/ogame/ships/lightFighter"
	"github.com/alaingilbert/ogame/ships/recycler"
	"github.com/alaingilbert/ogame/ships/smallCargo"
	"github.com/alaingilbert/ogame/ships/solarSatellite"
)

// Ship ...
type Ship interface {
	GetOGameID() ogame.ID
	GetRequirements() map[ogame.ID]int
	IsAvailable(ogame.ResourcesBuildings, ogame.Facilities, ogame.Researches, int) bool
	GetPrice(int) ogame.Resources
	GetStructuralIntegrity() int
	GetShieldPower() int
	GetWeaponPower() int
	GetCargoCapacity() int
	GetBaseSpeed() int
	GetSpeed(researches ogame.Researches) int
	GetFuelConsumption() int
	GetRapidfireFrom() map[ogame.ID]int
	GetRapidfireAgainst() map[ogame.ID]int
}

// Ships ...
var (
	Battlecruiser  = battlecruiser.New()
	Battelship     = battleship.New()
	Bomber         = bomber.New()
	ColonyShip     = colonyShip.New()
	Cruiser        = cruiser.New()
	Deathstar      = deathstar.New()
	Destroyer      = destroyer.New()
	EspionageProbe = espionageProbe.New()
	HeavyFighter   = heavyFighter.New()
	LargeCargo     = largeCargo.New()
	LightFighter   = lightFighter.New()
	Recycler       = recycler.New()
	SmallCargo     = smallCargo.New()
	SolarSatellite = solarSatellite.New()

	All = []Ship{Battlecruiser, Battelship, Bomber, ColonyShip, Cruiser, Deathstar, Destroyer, EspionageProbe,
		HeavyFighter, LargeCargo, LightFighter, Recycler, SmallCargo, SolarSatellite}
)

// GetByID ...
func GetByID(ogameID ogame.ID) Ship {
	for _, b := range All {
		if b.GetOGameID() == ogameID {
			return b
		}
	}
	return nil
}
