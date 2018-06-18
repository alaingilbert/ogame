package buildings

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/buildings/allianceDepot"
	"github.com/alaingilbert/ogame/buildings/crystalMine"
	"github.com/alaingilbert/ogame/buildings/crystalStorage"
	"github.com/alaingilbert/ogame/buildings/deuteriumSynthesizer"
	"github.com/alaingilbert/ogame/buildings/deuteriumTank"
	"github.com/alaingilbert/ogame/buildings/fusionReactor"
	"github.com/alaingilbert/ogame/buildings/metalMine"
	"github.com/alaingilbert/ogame/buildings/metalStorage"
	"github.com/alaingilbert/ogame/buildings/missileSilo"
	"github.com/alaingilbert/ogame/buildings/naniteFactory"
	"github.com/alaingilbert/ogame/buildings/researchLab"
	"github.com/alaingilbert/ogame/buildings/roboticsFactory"
	"github.com/alaingilbert/ogame/buildings/seabedDeuteriumDen"
	"github.com/alaingilbert/ogame/buildings/shieldedMetalDen"
	"github.com/alaingilbert/ogame/buildings/shipyard"
	"github.com/alaingilbert/ogame/buildings/solarPlant"
	"github.com/alaingilbert/ogame/buildings/spaceDock"
	"github.com/alaingilbert/ogame/buildings/terraformer"
	"github.com/alaingilbert/ogame/buildings/undergroundCrystalDen"
)

// Building ...
type Building interface {
	GetOGameID() ogame.ID
	GetBaseCost() ogame.Resources
	GetIncreaseFactor() float64
	GetRequirements() map[ogame.ID]int
	IsAvailable(ogame.ResourcesBuildings, ogame.Facilities, ogame.Researches, int) bool
	GetLevel(ogame.ResourcesBuildings, ogame.Facilities, ogame.Researches) int
	GetPrice(level int) ogame.Resources
	ConstructionTime(level, universeSpeed int, facilities ogame.Facilities) int
}

// Buildings ...
var (
	AllianceDepot         = allianceDepot.New()
	CrystalMine           = crystalMine.New()
	CrystalStorage        = crystalStorage.New()
	DeuteriumSynthesizer  = deuteriumSynthesizer.New()
	DeuteriumTank         = deuteriumTank.New()
	FusionReactor         = fusionReactor.New()
	MetalMine             = metalMine.New()
	MetalStorage          = metalStorage.New()
	MissileSilo           = missileSilo.New()
	NaniteFactory         = naniteFactory.New()
	ResearchLab           = researchLab.New()
	RoboticsFactory       = roboticsFactory.New()
	SeabedDeuteriumDen    = seabedDeuteriumDen.New()
	ShieldedMetalDen      = shieldedMetalDen.New()
	Shipyard              = shipyard.New()
	SolarPlant            = solarPlant.New()
	SpaceDock             = spaceDock.New()
	Terraformer           = terraformer.New()
	UndergroundCrystalDen = undergroundCrystalDen.New()

	All = []Building{AllianceDepot, CrystalMine, CrystalStorage, DeuteriumSynthesizer, DeuteriumTank, FusionReactor,
		MetalMine, MetalStorage, MissileSilo, NaniteFactory, ResearchLab, RoboticsFactory, SeabedDeuteriumDen,
		ShieldedMetalDen, Shipyard, SolarPlant, SpaceDock, Terraformer, UndergroundCrystalDen}
)

// GetByID ...
func GetByID(ogameID ogame.ID) Building {
	for _, b := range All {
		if b.GetOGameID() == ogameID {
			return b
		}
	}
	return nil
}
