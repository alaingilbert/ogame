package technologies

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/technologies/armourTechnology"
	"github.com/alaingilbert/ogame/technologies/astrophysics"
	"github.com/alaingilbert/ogame/technologies/combustionDrive"
	"github.com/alaingilbert/ogame/technologies/computerTechnology"
	"github.com/alaingilbert/ogame/technologies/energyTechnology"
	"github.com/alaingilbert/ogame/technologies/espionageTechnology"
	"github.com/alaingilbert/ogame/technologies/gravitonTechnology"
	"github.com/alaingilbert/ogame/technologies/hyperspaceDrive"
	"github.com/alaingilbert/ogame/technologies/hyperspaceTechnology"
	"github.com/alaingilbert/ogame/technologies/impulseDrive"
	"github.com/alaingilbert/ogame/technologies/intergalacticResearchNetwork"
	"github.com/alaingilbert/ogame/technologies/ionTechnology"
	"github.com/alaingilbert/ogame/technologies/laserTechnology"
	"github.com/alaingilbert/ogame/technologies/plasmaTechnology"
	"github.com/alaingilbert/ogame/technologies/shieldingTechnology"
	"github.com/alaingilbert/ogame/technologies/weaponsTechnology"
)

// Technology ...
type Technology interface {
	GetOGameID() ogame.ID
	GetBaseCost() ogame.Resources
	GetIncreaseFactor() float64
	GetRequirements() map[ogame.ID]int
	IsAvailable(ogame.ResourcesBuildings, ogame.Facilities, ogame.Researches, int) bool
	GetLevel(ogame.ResourcesBuildings, ogame.Facilities, ogame.Researches) int
	GetPrice(level int) ogame.Resources
	ConstructionTime(level, universeSpeed int, facilities ogame.Facilities) int
}

// Technologies
var (
	ArmourTechnology             = armourTechnology.New()
	Astrophysics                 = astrophysics.New()
	CombustionDrive              = combustionDrive.New()
	ComputerTechnology           = computerTechnology.New()
	EnergyTechnology             = energyTechnology.New()
	EspionageTechnology          = espionageTechnology.New()
	GravitonTechnology           = gravitonTechnology.New()
	HyperspaceDrive              = hyperspaceDrive.New()
	HyperspaceTechnology         = hyperspaceTechnology.New()
	ImpulseDrive                 = impulseDrive.New()
	IntergalacticResearchNetwork = intergalacticResearchNetwork.New()
	IonTechnology                = ionTechnology.New()
	LaserTechnology              = laserTechnology.New()
	PlasmaTechnology             = plasmaTechnology.New()
	ShieldingTechnology          = shieldingTechnology.New()
	WeaponsTechnology            = weaponsTechnology.New()

	All = []Technology{ArmourTechnology, Astrophysics, CombustionDrive, ComputerTechnology, EnergyTechnology,
		EspionageTechnology, GravitonTechnology, HyperspaceDrive, HyperspaceTechnology, ImpulseDrive,
		IntergalacticResearchNetwork, IonTechnology, LaserTechnology, PlasmaTechnology, ShieldingTechnology,
		WeaponsTechnology}
)

// GetByID ...
func GetByID(ogameID ogame.ID) Technology {
	for _, b := range All {
		if b.GetOGameID() == ogameID {
			return b
		}
	}
	return nil
}
