package defenses

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/defenses/antiBallisticMissiles"
	"github.com/alaingilbert/ogame/defenses/gaussCannon"
	"github.com/alaingilbert/ogame/defenses/heavyLaser"
	"github.com/alaingilbert/ogame/defenses/interplanetaryMissiles"
	"github.com/alaingilbert/ogame/defenses/ionCannon"
	"github.com/alaingilbert/ogame/defenses/largeShieldDome"
	"github.com/alaingilbert/ogame/defenses/lightLaser"
	"github.com/alaingilbert/ogame/defenses/plasmaTurret"
	"github.com/alaingilbert/ogame/defenses/rocketLauncher"
	"github.com/alaingilbert/ogame/defenses/smallShieldDome"
)

// Defense ...
type Defense interface {
	GetOGameID() ogame.ID
	GetPrice(int) ogame.Resources
	GetRequirements() map[ogame.ID]int
	IsAvailable(ogame.ResourcesBuildings, ogame.Facilities, ogame.Researches, int) bool
	GetStructuralIntegrity() int
	GetShieldPower() int
	GetWeaponPower() int
	GetRapidfireFrom() map[ogame.ID]int
}

// Defenses ...
var (
	AntiBallisticMissiles  = antiBallisticMissiles.New()
	GaussCannon            = gaussCannon.New()
	HeavyLaser             = heavyLaser.New()
	InterplanetaryMissiles = interplanetaryMissiles.New()
	IonCannon              = ionCannon.New()
	LargeShieldDome        = largeShieldDome.New()
	LightLaser             = lightLaser.New()
	PlasmaTurret           = plasmaTurret.New()
	RocketLauncher         = rocketLauncher.New()
	SmallShieldDome        = smallShieldDome.New()

	All = []Defense{AntiBallisticMissiles, GaussCannon, HeavyLaser, InterplanetaryMissiles, IonCannon, LargeShieldDome,
		LightLaser, PlasmaTurret, RocketLauncher, SmallShieldDome}
)

// GetByID ...
func GetByID(ogameID ogame.ID) Defense {
	for _, b := range All {
		if b.GetOGameID() == ogameID {
			return b
		}
	}
	return nil
}
