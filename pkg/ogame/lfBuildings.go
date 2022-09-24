package ogame

import (
	"math"
)

// LazyLfBuildings ...
type LazyLfBuildings func() LfBuildings

// LfBuildings lifeform buildings
type LfBuildings struct {
	ResidentialSector          int64 // 11101 // Lifeform (humans)
	BiosphereFarm              int64 // 11102
	ResearchCentre             int64 // 11103
	AcademyOfSciences          int64 // 11104
	NeuroCalibrationCentre     int64 // 11105
	HighEnergySmelting         int64 // 11106
	FoodSilo                   int64 // 11107
	FusionPoweredProduction    int64 // 11108
	Skyscraper                 int64 // 11109
	BiotechLab                 int64 // 11110
	Metropolis                 int64 // 11111
	PlanetaryShield            int64 // 11112
	MeditationEnclave          int64 // 12101 // Lifeform (rocktal)
	CrystalFarm                int64 // 12102
	RuneTechnologium           int64 // 12103
	RuneForge                  int64 // 12104
	Oriktorium                 int64 // 12105
	MagmaForge                 int64 // 12106
	DisruptionChamber          int64 // 12107
	Megalith                   int64 // 12108
	CrystalRefinery            int64 // 12109
	DeuteriumSynthesiser       int64 // 12110
	MineralResearchCentre      int64 // 12111
	MetalRecyclingPlant        int64 // 12112
	AssemblyLine               int64 // 13101 // Lifeform (mechas)
	FusionCellFactory          int64 // 13102
	RoboticsResearchCentre     int64 // 13103
	UpdateNetwork              int64 // 12304
	QuantumComputerCentre      int64 // 13105
	AutomatisedAssemblyCentre  int64 // 13106
	HighPerformanceTransformer int64 // 13107
	MicrochipAssemblyLine      int64 // 13108
	ProductionAssemblyHall     int64 // 13109
	HighPerformanceSynthesiser int64 // 13110
	ChipMassProduction         int64 // 13111
	NanoRepairBots             int64 // 13112
	Sanctuary                  int64 // 14101 // Lifeform (kaelesh)
	AntimatterCondenser        int64 // 14102
	VortexChamber              int64 // 14103
	HallsOfRealisation         int64 // 14104
	ForumOfTranscendence       int64 // 14105
	AntimatterConvector        int64 // 14106
	CloningLaboratory          int64 // 14107
	ChrysalisAccelerator       int64 // 14108
	BioModifier                int64 // 14109
	PsionicModulator           int64 // 14110
	ShipManufacturingHall      int64 // 14111
	SupraRefractor             int64 // 14112
}

// Lazy returns a function that return self
func (b LfBuildings) Lazy() LazyLfBuildings {
	return func() LfBuildings { return b }
}

// ByID gets the lfBuilding level by lfBuilding id
func (b LfBuildings) ByID(id ID) int64 {
	switch id {
	case ResidentialSectorID:
		return b.ResidentialSector
	case BiosphereFarmID:
		return b.BiosphereFarm
	case ResearchCentreID:
		return b.ResearchCentre
	case AcademyOfSciencesID:
		return b.AcademyOfSciences
	case NeuroCalibrationCentreID:
		return b.NeuroCalibrationCentre
	case HighEnergySmeltingID:
		return b.HighEnergySmelting
	case FoodSiloID:
		return b.FoodSilo
	case FusionPoweredProductionID:
		return b.FusionPoweredProduction
	case SkyscraperID:
		return b.Skyscraper
	case BiotechLabID:
		return b.BiotechLab
	case MetropolisID:
		return b.Metropolis
	case PlanetaryShieldID:
		return b.PlanetaryShield
	case MeditationEnclaveID:
		return b.MeditationEnclave
	case CrystalFarmID:
		return b.CrystalFarm
	case RuneTechnologiumID:
		return b.RuneTechnologium
	case RuneForgeID:
		return b.RuneForge
	case OriktoriumID:
		return b.Oriktorium
	case MagmaForgeID:
		return b.MagmaForge
	case DisruptionChamberID:
		return b.DisruptionChamber
	case MegalithID:
		return b.Megalith
	case CrystalRefineryID:
		return b.CrystalRefinery
	case DeuteriumSynthesiserID:
		return b.DeuteriumSynthesiser
	case MineralResearchCentreID:
		return b.MineralResearchCentre
	case MetalRecyclingPlantID:
		return b.MetalRecyclingPlant
	case AssemblyLineID:
		return b.AssemblyLine
	case FusionCellFactoryID:
		return b.FusionCellFactory
	case RoboticsResearchCentreID:
		return b.RoboticsResearchCentre
	case UpdateNetworkID:
		return b.UpdateNetwork
	case QuantumComputerCentreID:
		return b.QuantumComputerCentre
	case AutomatisedAssemblyCentreID:
		return b.AutomatisedAssemblyCentre
	case HighPerformanceTransformerID:
		return b.HighPerformanceTransformer
	case MicrochipAssemblyLineID:
		return b.MicrochipAssemblyLine
	case ProductionAssemblyHallID:
		return b.ProductionAssemblyHall
	case HighPerformanceSynthesiserID:
		return b.HighPerformanceSynthesiser
	case ChipMassProductionID:
		return b.ChipMassProduction
	case NanoRepairBotsID:
		return b.NanoRepairBots
	case SanctuaryID:
		return b.Sanctuary
	case AntimatterCondenserID:
		return b.AntimatterCondenser
	case VortexChamberID:
		return b.VortexChamber
	case HallsOfRealisationID:
		return b.HallsOfRealisation
	case ForumOfTranscendenceID:
		return b.ForumOfTranscendence
	case AntimatterConvectorID:
		return b.AntimatterConvector
	case CloningLaboratoryID:
		return b.CloningLaboratory
	case ChrysalisAcceleratorID:
		return b.ChrysalisAccelerator
	case BioModifierID:
		return b.BioModifier
	case PsionicModulatorID:
		return b.PsionicModulator
	case ShipManufacturingHallID:
		return b.ShipManufacturingHall
	case SupraRefractorID:
		return b.SupraRefractor
	}
	return 0
}

// BaseLfBuilding base struct for Lifeform buildings
type BaseLfBuilding struct {
	BaseBuilding
	energyIncreaseFactor float64
}

// GetPrice returns the price to build the given level
func (b BaseLfBuilding) GetPrice(level int64) Resources {
	tmp := func(baseCost int64, increaseFactor float64, level int64) int64 {
		return int64(float64(baseCost) * math.Pow(increaseFactor, float64(level-1)) * float64(level))
	}
	return Resources{
		Metal:     tmp(b.BaseCost.Metal, b.IncreaseFactor, level),
		Crystal:   tmp(b.BaseCost.Crystal, b.IncreaseFactor, level),
		Deuterium: tmp(b.BaseCost.Deuterium, b.IncreaseFactor, level),
		Energy:    tmp(b.BaseCost.Energy, b.energyIncreaseFactor, level),
	}
}

type residentialSector struct {
	BaseLfBuilding
}

func newResidentialSector() *residentialSector {
	b := new(residentialSector)
	b.Name = "residential sector"
	b.ID = ResidentialSectorID
	b.IncreaseFactor = 1.2
	b.BaseCost = Resources{Metal: 7, Crystal: 2}
	b.Requirements = map[ID]int64{}
	return b
}

type biosphereFarm struct {
	BaseLfBuilding
}

func newBiosphereFarm() *biosphereFarm {
	b := new(biosphereFarm)
	b.Name = "biosphere farm"
	b.ID = BiosphereFarmID
	b.IncreaseFactor = 1.23
	b.energyIncreaseFactor = 1.021
	b.BaseCost = Resources{Metal: 5, Crystal: 2, Energy: 8}
	b.Requirements = map[ID]int64{}
	return b
}

type researchCentre struct {
	BaseLfBuilding
}

func newResearchCentre() *researchCentre {
	b := new(researchCentre)
	b.Name = "research centre"
	b.ID = ResearchCentreID
	b.IncreaseFactor = 1.3
	b.BaseCost = Resources{Metal: 20000, Crystal: 25000, Deuterium: 10000}
	b.Requirements = map[ID]int64{ResidentialSectorID: 12, BiosphereFarmID: 13}
	return b
}

type academyOfSciences struct {
	BaseLfBuilding
}

func newAcademyOfSciences() *academyOfSciences {
	b := new(academyOfSciences)
	b.Name = "academy of sciences"
	b.ID = AcademyOfSciencesID
	b.IncreaseFactor = 1 // TODO
	b.BaseCost = Resources{Metal: 5000, Crystal: 3200, Deuterium: 1500, Lifeform: 20000000}
	b.Requirements = map[ID]int64{ResidentialSectorID: 40}
	return b
}

type neuroCalibrationCentre struct {
	BaseLfBuilding
}

func newNeuroCalibrationCentre() *neuroCalibrationCentre {
	b := new(neuroCalibrationCentre)
	b.Name = "neuro calibration centre"
	b.ID = NeuroCalibrationCentreID
	b.IncreaseFactor = 1 // TODO
	b.BaseCost = Resources{Metal: 50000, Crystal: 40000, Deuterium: 50000, Lifeform: 100000000}
	b.Requirements = map[ID]int64{ResidentialSectorID: 40, AcademyOfSciencesID: 1, FusionPoweredProductionID: 1, SkyscraperID: 5}
	return b
}

type highEnergySmelting struct {
	BaseLfBuilding
}

func newHighEnergySmelting() *highEnergySmelting {
	b := new(highEnergySmelting)
	b.Name = "high energy smelting"
	b.ID = HighEnergySmeltingID
	b.IncreaseFactor = 1 // TODO
	b.BaseCost = Resources{Metal: 7500, Crystal: 5000, Deuterium: 3000}
	b.Requirements = map[ID]int64{ResidentialSectorID: 12, BiosphereFarmID: 13, ResearchCentreID: 5}
	return b
}

type foodSilo struct {
	BaseLfBuilding
}

func newFoodSilo() *foodSilo {
	b := new(foodSilo)
	b.Name = "food silo"
	b.ID = FoodSiloID
	b.IncreaseFactor = 1 // TODO
	b.BaseCost = Resources{Metal: 25000, Crystal: 13000, Deuterium: 7000}
	b.Requirements = map[ID]int64{ResidentialSectorID: 12, BiosphereFarmID: 13, ResearchCentreID: 5, HighEnergySmeltingID: 3}
	return b
}

type fusionPoweredProduction struct {
	BaseLfBuilding
}

func newFusionPoweredProduction() *fusionPoweredProduction {
	b := new(fusionPoweredProduction)
	b.Name = "fusion powered production"
	b.ID = FusionPoweredProductionID
	b.IncreaseFactor = 1 // TODO
	b.BaseCost = Resources{Metal: 50000, Crystal: 25000, Deuterium: 25000}
	b.Requirements = map[ID]int64{ResidentialSectorID: 40, AcademyOfSciencesID: 1}
	return b
}

type skyscraper struct {
	BaseLfBuilding
}

func newSkyscraper() *skyscraper {
	b := new(skyscraper)
	b.Name = "skyscraper"
	b.ID = SkyscraperID
	b.IncreaseFactor = 1 // TODO
	b.BaseCost = Resources{Metal: 75000, Crystal: 20000, Deuterium: 25000}
	b.Requirements = map[ID]int64{ResidentialSectorID: 40, AcademyOfSciencesID: 1, FusionPoweredProductionID: 1}
	return b
}

type biotechLab struct {
	BaseLfBuilding
}

func newBiotechLab() *biotechLab {
	b := new(biotechLab)
	b.Name = "biotech lab"
	b.ID = BiotechLabID
	b.IncreaseFactor = 1 // TODO
	b.BaseCost = Resources{Metal: 150000, Crystal: 30000, Deuterium: 15000}
	b.Requirements = map[ID]int64{ResidentialSectorID: 40, AcademyOfSciencesID: 1, FusionPoweredProductionID: 2}
	return b
}

type metropolis struct {
	BaseLfBuilding
}

func newMetropolis() *metropolis {
	b := new(metropolis)
	b.Name = "metropolis"
	b.ID = MetropolisID
	b.IncreaseFactor = 1 // TODO
	b.BaseCost = Resources{Metal: 80000, Crystal: 35000, Deuterium: 60000}
	b.Requirements = map[ID]int64{ResidentialSectorID: 40, AcademyOfSciencesID: 1, FusionPoweredProductionID: 1, SkyscraperID: 5, NeuroCalibrationCentreID: 1}
	return b
}

type planetaryShield struct {
	BaseLfBuilding
}

func newPlanetaryShield() *planetaryShield {
	b := new(planetaryShield)
	b.Name = "planetary shield"
	b.ID = PlanetaryShieldID
	b.IncreaseFactor = 1 // TODO
	b.BaseCost = Resources{Metal: 250000, Crystal: 125000, Deuterium: 125000}
	b.Requirements = map[ID]int64{
		ResidentialSectorID:       40,
		BiosphereFarmID:           13,
		ResearchCentreID:          5,
		AcademyOfSciencesID:       1,
		FusionPoweredProductionID: 5,
		SkyscraperID:              5,
		HighEnergySmeltingID:      3,
		MetropolisID:              5,
		FoodSiloID:                4,
		NeuroCalibrationCentreID:  5}
	return b
}

// BaseLfTechnology base struct for lifeform technologies
type BaseLfTechnology struct {
	BaseLevelable
}
