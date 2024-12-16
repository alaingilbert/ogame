package ogame

import (
	"math"
	"time"
)

// LazyLfBuildings ...
type LazyLfBuildings func() LfBuildings

func (b LazyLfBuildings) ByID(id ID) int64                     { return b().ByID(id) }
func (b LazyLfBuildings) GetResidentialSector() int64          { return b().ResidentialSector }
func (b LazyLfBuildings) GetBiosphereFarm() int64              { return b().BiosphereFarm }
func (b LazyLfBuildings) GetResearchCentre() int64             { return b().ResearchCentre }
func (b LazyLfBuildings) GetAcademyOfSciences() int64          { return b().AcademyOfSciences }
func (b LazyLfBuildings) GetNeuroCalibrationCentre() int64     { return b().NeuroCalibrationCentre }
func (b LazyLfBuildings) GetHighEnergySmelting() int64         { return b().HighEnergySmelting }
func (b LazyLfBuildings) GetFoodSilo() int64                   { return b().FoodSilo }
func (b LazyLfBuildings) GetFusionPoweredProduction() int64    { return b().FusionPoweredProduction }
func (b LazyLfBuildings) GetSkyscraper() int64                 { return b().Skyscraper }
func (b LazyLfBuildings) GetBiotechLab() int64                 { return b().BiotechLab }
func (b LazyLfBuildings) GetMetropolis() int64                 { return b().Metropolis }
func (b LazyLfBuildings) GetPlanetaryShield() int64            { return b().PlanetaryShield }
func (b LazyLfBuildings) GetMeditationEnclave() int64          { return b().MeditationEnclave }
func (b LazyLfBuildings) GetCrystalFarm() int64                { return b().CrystalFarm }
func (b LazyLfBuildings) GetRuneTechnologium() int64           { return b().RuneTechnologium }
func (b LazyLfBuildings) GetRuneForge() int64                  { return b().RuneForge }
func (b LazyLfBuildings) GetOriktorium() int64                 { return b().Oriktorium }
func (b LazyLfBuildings) GetMagmaForge() int64                 { return b().MagmaForge }
func (b LazyLfBuildings) GetDisruptionChamber() int64          { return b().DisruptionChamber }
func (b LazyLfBuildings) GetMegalith() int64                   { return b().Megalith }
func (b LazyLfBuildings) GetCrystalRefinery() int64            { return b().CrystalRefinery }
func (b LazyLfBuildings) GetDeuteriumSynthesiser() int64       { return b().DeuteriumSynthesiser }
func (b LazyLfBuildings) GetMineralResearchCentre() int64      { return b().MineralResearchCentre }
func (b LazyLfBuildings) GetAdvancedRecyclingPlant() int64     { return b().AdvancedRecyclingPlant }
func (b LazyLfBuildings) GetAssemblyLine() int64               { return b().AssemblyLine }
func (b LazyLfBuildings) GetFusionCellFactory() int64          { return b().FusionCellFactory }
func (b LazyLfBuildings) GetRoboticsResearchCentre() int64     { return b().RoboticsResearchCentre }
func (b LazyLfBuildings) GetUpdateNetwork() int64              { return b().UpdateNetwork }
func (b LazyLfBuildings) GetQuantumComputerCentre() int64      { return b().QuantumComputerCentre }
func (b LazyLfBuildings) GetAutomatisedAssemblyCentre() int64  { return b().AutomatisedAssemblyCentre }
func (b LazyLfBuildings) GetHighPerformanceTransformer() int64 { return b().HighPerformanceTransformer }
func (b LazyLfBuildings) GetMicrochipAssemblyLine() int64      { return b().MicrochipAssemblyLine }
func (b LazyLfBuildings) GetProductionAssemblyHall() int64     { return b().ProductionAssemblyHall }
func (b LazyLfBuildings) GetHighPerformanceSynthesiser() int64 { return b().HighPerformanceSynthesiser }
func (b LazyLfBuildings) GetChipMassProduction() int64         { return b().ChipMassProduction }
func (b LazyLfBuildings) GetNanoRepairBots() int64             { return b().NanoRepairBots }
func (b LazyLfBuildings) GetSanctuary() int64                  { return b().Sanctuary }
func (b LazyLfBuildings) GetAntimatterCondenser() int64        { return b().AntimatterCondenser }
func (b LazyLfBuildings) GetVortexChamber() int64              { return b().VortexChamber }
func (b LazyLfBuildings) GetHallsOfRealisation() int64         { return b().HallsOfRealisation }
func (b LazyLfBuildings) GetForumOfTranscendence() int64       { return b().ForumOfTranscendence }
func (b LazyLfBuildings) GetAntimatterConvector() int64        { return b().AntimatterConvector }
func (b LazyLfBuildings) GetCloningLaboratory() int64          { return b().CloningLaboratory }
func (b LazyLfBuildings) GetChrysalisAccelerator() int64       { return b().ChrysalisAccelerator }
func (b LazyLfBuildings) GetBioModifier() int64                { return b().BioModifier }
func (b LazyLfBuildings) GetPsionicModulator() int64           { return b().PsionicModulator }
func (b LazyLfBuildings) GetShipManufacturingHall() int64      { return b().ShipManufacturingHall }
func (b LazyLfBuildings) GetSupraRefractor() int64             { return b().SupraRefractor }

type LifeformType int64

const (
	NoneLfType LifeformType = iota
	Humans
	Rocktal
	Mechas
	Kaelesh
)

func (l *LifeformType) String() string {
	switch *l {
	case Humans:
		return "humans"
	case Rocktal:
		return "rocktal"
	case Mechas:
		return "mechas"
	case Kaelesh:
		return "kaelesh"
	default:
		return "none"
	}
}

// LfBuildings lifeform buildings
type LfBuildings struct {
	LifeformType               LifeformType
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
	AdvancedRecyclingPlant     int64 // 12112
	AssemblyLine               int64 // 13101 // Lifeform (mechas)
	FusionCellFactory          int64 // 13102
	RoboticsResearchCentre     int64 // 13103
	UpdateNetwork              int64 // 13104
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

func (b LfBuildings) GetResidentialSector() int64          { return b.ResidentialSector }
func (b LfBuildings) GetBiosphereFarm() int64              { return b.BiosphereFarm }
func (b LfBuildings) GetResearchCentre() int64             { return b.ResearchCentre }
func (b LfBuildings) GetAcademyOfSciences() int64          { return b.AcademyOfSciences }
func (b LfBuildings) GetNeuroCalibrationCentre() int64     { return b.NeuroCalibrationCentre }
func (b LfBuildings) GetHighEnergySmelting() int64         { return b.HighEnergySmelting }
func (b LfBuildings) GetFoodSilo() int64                   { return b.FoodSilo }
func (b LfBuildings) GetFusionPoweredProduction() int64    { return b.FusionPoweredProduction }
func (b LfBuildings) GetSkyscraper() int64                 { return b.Skyscraper }
func (b LfBuildings) GetBiotechLab() int64                 { return b.BiotechLab }
func (b LfBuildings) GetMetropolis() int64                 { return b.Metropolis }
func (b LfBuildings) GetPlanetaryShield() int64            { return b.PlanetaryShield }
func (b LfBuildings) GetMeditationEnclave() int64          { return b.MeditationEnclave }
func (b LfBuildings) GetCrystalFarm() int64                { return b.CrystalFarm }
func (b LfBuildings) GetRuneTechnologium() int64           { return b.RuneTechnologium }
func (b LfBuildings) GetRuneForge() int64                  { return b.RuneForge }
func (b LfBuildings) GetOriktorium() int64                 { return b.Oriktorium }
func (b LfBuildings) GetMagmaForge() int64                 { return b.MagmaForge }
func (b LfBuildings) GetDisruptionChamber() int64          { return b.DisruptionChamber }
func (b LfBuildings) GetMegalith() int64                   { return b.Megalith }
func (b LfBuildings) GetCrystalRefinery() int64            { return b.CrystalRefinery }
func (b LfBuildings) GetDeuteriumSynthesiser() int64       { return b.DeuteriumSynthesiser }
func (b LfBuildings) GetMineralResearchCentre() int64      { return b.MineralResearchCentre }
func (b LfBuildings) GetAdvancedRecyclingPlant() int64     { return b.AdvancedRecyclingPlant }
func (b LfBuildings) GetAssemblyLine() int64               { return b.AssemblyLine }
func (b LfBuildings) GetFusionCellFactory() int64          { return b.FusionCellFactory }
func (b LfBuildings) GetRoboticsResearchCentre() int64     { return b.RoboticsResearchCentre }
func (b LfBuildings) GetUpdateNetwork() int64              { return b.UpdateNetwork }
func (b LfBuildings) GetQuantumComputerCentre() int64      { return b.QuantumComputerCentre }
func (b LfBuildings) GetAutomatisedAssemblyCentre() int64  { return b.AutomatisedAssemblyCentre }
func (b LfBuildings) GetHighPerformanceTransformer() int64 { return b.HighPerformanceTransformer }
func (b LfBuildings) GetMicrochipAssemblyLine() int64      { return b.MicrochipAssemblyLine }
func (b LfBuildings) GetProductionAssemblyHall() int64     { return b.ProductionAssemblyHall }
func (b LfBuildings) GetHighPerformanceSynthesiser() int64 { return b.HighPerformanceSynthesiser }
func (b LfBuildings) GetChipMassProduction() int64         { return b.ChipMassProduction }
func (b LfBuildings) GetNanoRepairBots() int64             { return b.NanoRepairBots }
func (b LfBuildings) GetSanctuary() int64                  { return b.Sanctuary }
func (b LfBuildings) GetAntimatterCondenser() int64        { return b.AntimatterCondenser }
func (b LfBuildings) GetVortexChamber() int64              { return b.VortexChamber }
func (b LfBuildings) GetHallsOfRealisation() int64         { return b.HallsOfRealisation }
func (b LfBuildings) GetForumOfTranscendence() int64       { return b.ForumOfTranscendence }
func (b LfBuildings) GetAntimatterConvector() int64        { return b.AntimatterConvector }
func (b LfBuildings) GetCloningLaboratory() int64          { return b.CloningLaboratory }
func (b LfBuildings) GetChrysalisAccelerator() int64       { return b.ChrysalisAccelerator }
func (b LfBuildings) GetBioModifier() int64                { return b.BioModifier }
func (b LfBuildings) GetPsionicModulator() int64           { return b.PsionicModulator }
func (b LfBuildings) GetShipManufacturingHall() int64      { return b.ShipManufacturingHall }
func (b LfBuildings) GetSupraRefractor() int64             { return b.SupraRefractor }

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
	case AdvancedRecyclingPlantID:
		return b.AdvancedRecyclingPlant
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
	energyIncreaseFactor     float64
	populationIncreaseFactor float64
	durationBase             float64
	durationFactor           float64
}

func (b BaseLfBuilding) BuildingConstructionTime(level, universeSpeed int64, acc BuildingAccelerators, lfBonuses LfBonuses) time.Duration {
	roboticLvl := float64(acc.GetRoboticsFactory())
	naniteLvl := float64(acc.GetNaniteFactory())
	levelF := float64(level)
	secs := levelF * b.durationBase * math.Pow(b.durationFactor, levelF) / ((1 + roboticLvl) * math.Pow(2, naniteLvl))
	secs /= float64(universeSpeed)
	secs = math.Max(1, secs)
	dur := time.Duration(int64(math.Floor(secs))) * time.Second
	bonus := lfBonuses.CostTimeBonuses[b.ID].Duration
	return time.Duration(float64(dur) - float64(dur)*bonus)
}

// ConstructionTime returns the duration it takes to build given level. Deconstruction time is the same function.
func (b BaseLfBuilding) ConstructionTime(level, universeSpeed int64, facilities BuildAccelerators, lfBonuses LfBonuses, _ CharacterClass, _ bool) time.Duration {
	return b.BuildingConstructionTime(level, universeSpeed, facilities, lfBonuses)
}

// GetPrice returns the price to build the given level
func (b BaseLfBuilding) GetPrice(level int64, _ LfBonuses) Resources {
	resourcesFormula := func(baseCost int64, increaseFactor float64, level int64) int64 {
		return int64(float64(baseCost) * math.Pow(increaseFactor, float64(level-1)) * float64(level))
	}
	populationFormula := func(baseCost int64, increaseFactor float64, level int64) int64 {
		return int64(float64(baseCost) * math.Pow(increaseFactor, float64(level-1)))
	}
	return Resources{
		Metal:      resourcesFormula(b.BaseCost.Metal, b.IncreaseFactor, level),
		Crystal:    resourcesFormula(b.BaseCost.Crystal, b.IncreaseFactor, level),
		Deuterium:  resourcesFormula(b.BaseCost.Deuterium, b.IncreaseFactor, level),
		Energy:     resourcesFormula(b.BaseCost.Energy, b.energyIncreaseFactor, level),
		Population: populationFormula(b.BaseCost.Population, b.populationIncreaseFactor, level),
	}
}

// Humans
type residentialSector struct {
	BaseLfBuilding
}

func newResidentialSector() *residentialSector {
	b := new(residentialSector)
	b.Name = "residential sector"
	b.ID = ResidentialSectorID
	b.IncreaseFactor = 1.20
	b.durationBase = 40
	b.durationFactor = 1.21
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
	b.durationBase = 40
	b.durationFactor = 1.25
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
	b.durationBase = 16000
	b.durationFactor = 1.25
	b.BaseCost = Resources{Metal: 20000, Crystal: 25000, Deuterium: 10000}
	b.Requirements = map[ID]int64{ResidentialSectorID: 21, BiosphereFarmID: 22}
	return b
}

type academyOfSciences struct {
	BaseLfBuilding
}

func newAcademyOfSciences() *academyOfSciences {
	b := new(academyOfSciences)
	b.Name = "academy of sciences"
	b.ID = AcademyOfSciencesID
	b.IncreaseFactor = 1.70
	b.populationIncreaseFactor = 1.10
	b.durationBase = 16000
	b.durationFactor = 1.6
	b.BaseCost = Resources{Metal: 5000, Crystal: 3200, Deuterium: 1500, Population: 20000000}
	b.Requirements = map[ID]int64{ResidentialSectorID: 41}
	return b
}

type neuroCalibrationCentre struct {
	BaseLfBuilding
}

func newNeuroCalibrationCentre() *neuroCalibrationCentre {
	b := new(neuroCalibrationCentre)
	b.Name = "neuro calibration centre"
	b.ID = NeuroCalibrationCentreID
	b.IncreaseFactor = 1.70
	b.populationIncreaseFactor = 1.10
	b.durationBase = 64000
	b.durationFactor = 1.7
	b.BaseCost = Resources{Metal: 50000, Crystal: 40000, Deuterium: 50000, Population: 100000000}
	b.Requirements = map[ID]int64{ResidentialSectorID: 41, AcademyOfSciencesID: 1, FusionPoweredProductionID: 1, SkyscraperID: 5}
	return b
}

type highEnergySmelting struct {
	BaseLfBuilding
}

func newHighEnergySmelting() *highEnergySmelting {
	b := new(highEnergySmelting)
	b.Name = "high energy smelting"
	b.ID = HighEnergySmeltingID
	b.durationBase = 2000
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 9000, Crystal: 6000, Deuterium: 3000}
	b.Requirements = map[ID]int64{ResidentialSectorID: 21, BiosphereFarmID: 22, ResearchCentreID: 5}
	return b
}

type foodSilo struct {
	BaseLfBuilding
}

func newFoodSilo() *foodSilo {
	b := new(foodSilo)
	b.Name = "food silo"
	b.ID = FoodSiloID
	b.durationBase = 12000
	b.durationFactor = 1.17
	b.IncreaseFactor = 1.09
	b.BaseCost = Resources{Metal: 25000, Crystal: 13000, Deuterium: 7000}
	b.Requirements = map[ID]int64{ResidentialSectorID: 21, BiosphereFarmID: 22, ResearchCentreID: 5, HighEnergySmeltingID: 3}
	return b
}

type fusionPoweredProduction struct {
	BaseLfBuilding
}

func newFusionPoweredProduction() *fusionPoweredProduction {
	b := new(fusionPoweredProduction)
	b.Name = "fusion powered production"
	b.ID = FusionPoweredProductionID
	b.durationBase = 28000
	b.durationFactor = 1.2
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 50000, Crystal: 25000, Deuterium: 15000}
	b.Requirements = map[ID]int64{ResidentialSectorID: 41, AcademyOfSciencesID: 1}
	return b
}

type skyscraper struct {
	BaseLfBuilding
}

func newSkyscraper() *skyscraper {
	b := new(skyscraper)
	b.Name = "skyscraper"
	b.ID = SkyscraperID
	b.durationBase = 40000
	b.durationFactor = 1.2
	b.IncreaseFactor = 1.09
	b.BaseCost = Resources{Metal: 75000, Crystal: 20000, Deuterium: 25000}
	b.Requirements = map[ID]int64{ResidentialSectorID: 41, AcademyOfSciencesID: 1, FusionPoweredProductionID: 1}
	return b
}

type biotechLab struct {
	BaseLfBuilding
}

func newBiotechLab() *biotechLab {
	b := new(biotechLab)
	b.Name = "biotech lab"
	b.ID = BiotechLabID
	b.durationBase = 52000
	b.durationFactor = 1.2
	b.IncreaseFactor = 1.12
	b.BaseCost = Resources{Metal: 150000, Crystal: 30000, Deuterium: 15000}
	b.Requirements = map[ID]int64{ResidentialSectorID: 41, AcademyOfSciencesID: 1, FusionPoweredProductionID: 2}
	return b
}

type metropolis struct {
	BaseLfBuilding
}

func newMetropolis() *metropolis {
	b := new(metropolis)
	b.Name = "metropolis"
	b.ID = MetropolisID
	b.durationBase = 90000
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.5
	b.BaseCost = Resources{Metal: 80000, Crystal: 35000, Deuterium: 60000}
	b.Requirements = map[ID]int64{ResidentialSectorID: 41, AcademyOfSciencesID: 1, FusionPoweredProductionID: 1, SkyscraperID: 6, NeuroCalibrationCentreID: 1}
	return b
}

type planetaryShield struct {
	BaseLfBuilding
}

func newPlanetaryShield() *planetaryShield {
	b := new(planetaryShield)
	b.Name = "planetary shield"
	b.ID = PlanetaryShieldID
	b.durationBase = 95000
	b.durationFactor = 1.2
	b.IncreaseFactor = 1.20
	b.BaseCost = Resources{Metal: 250000, Crystal: 125000, Deuterium: 125000}
	b.Requirements = map[ID]int64{
		ResidentialSectorID:       41,
		BiosphereFarmID:           22,
		ResearchCentreID:          5,
		AcademyOfSciencesID:       1,
		FusionPoweredProductionID: 5,
		SkyscraperID:              6,
		HighEnergySmeltingID:      3,
		MetropolisID:              5,
		FoodSiloID:                4,
		NeuroCalibrationCentreID:  5}
	return b
}

// Rocktal
type meditationEnclave struct {
	BaseLfBuilding
}

func newMeditationEnclave() *meditationEnclave {
	b := new(meditationEnclave)
	b.Name = "meditation enclave"
	b.ID = MeditationEnclaveID
	b.durationBase = 40
	b.durationFactor = 1.21
	b.IncreaseFactor = 1.20
	b.BaseCost = Resources{Metal: 9, Crystal: 3}
	b.Requirements = map[ID]int64{}
	return b
}

type crystalFarm struct {
	BaseLfBuilding
}

func newCrystalFarm() *crystalFarm {
	b := new(crystalFarm)
	b.Name = "crystal farm"
	b.ID = CrystalFarmID
	b.durationBase = 40
	b.durationFactor = 1.21
	b.IncreaseFactor = 1.20
	b.energyIncreaseFactor = 1.03
	b.BaseCost = Resources{Metal: 7, Crystal: 2, Energy: 10}
	b.Requirements = map[ID]int64{}
	return b
}

type runeTechnologium struct {
	BaseLfBuilding
}

func newRuneTechnologium() *runeTechnologium {
	b := new(runeTechnologium)
	b.Name = "rune technologium"
	b.ID = RuneTechnologiumID
	b.durationBase = 16000
	b.durationFactor = 1.25
	b.IncreaseFactor = 1.30
	b.BaseCost = Resources{Metal: 40000, Crystal: 10000, Deuterium: 15000}
	b.Requirements = map[ID]int64{MeditationEnclaveID: 21, CrystalFarmID: 22}
	return b
}

type runeForge struct {
	BaseLfBuilding
}

func newRuneForge() *runeForge {
	b := new(runeForge)
	b.Name = "rune forge"
	b.ID = RuneForgeID
	b.durationBase = 16000
	b.durationFactor = 1.6
	b.IncreaseFactor = 1.70
	b.populationIncreaseFactor = 1.14
	b.BaseCost = Resources{Metal: 5000, Crystal: 3800, Deuterium: 1000, Population: 16000000}
	b.Requirements = map[ID]int64{MeditationEnclaveID: 41}
	return b
}

type oriktorium struct {
	BaseLfBuilding
}

func newOriktorium() *oriktorium {
	b := new(oriktorium)
	b.Name = "oriktorium"
	b.ID = OriktoriumID
	b.durationBase = 64000
	b.durationFactor = 1.7
	b.IncreaseFactor = 1.65
	b.populationIncreaseFactor = 1.1
	b.BaseCost = Resources{Metal: 50000, Crystal: 40000, Deuterium: 50000, Population: 90000000}
	b.Requirements = map[ID]int64{MeditationEnclaveID: 41, RuneForgeID: 1, MegalithID: 1, CrystalRefineryID: 5}
	return b
}

type magmaForge struct {
	BaseLfBuilding
}

func newMagmaForge() *magmaForge {
	b := new(magmaForge)
	b.Name = "magma forge"
	b.ID = MagmaForgeID
	b.durationBase = 2000
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.40
	b.BaseCost = Resources{Metal: 10000, Crystal: 8000, Deuterium: 1000}
	b.Requirements = map[ID]int64{MeditationEnclaveID: 21, CrystalFarmID: 22, RuneTechnologiumID: 5}
	return b
}

type disruptionChamber struct {
	BaseLfBuilding
}

func newDisruptionChamber() *disruptionChamber {
	b := new(disruptionChamber)
	b.Name = "disruption chamber"
	b.ID = DisruptionChamberID
	b.durationBase = 16000
	b.durationFactor = 1.25
	b.IncreaseFactor = 1.20
	b.BaseCost = Resources{Metal: 20000, Crystal: 15000, Deuterium: 10000}
	b.Requirements = map[ID]int64{MeditationEnclaveID: 21, CrystalFarmID: 22, RuneTechnologiumID: 5, MagmaForgeID: 3}
	return b
}

type megalith struct {
	BaseLfBuilding
}

func newMegalith() *megalith {
	b := new(megalith)
	b.Name = "megalith"
	b.ID = MegalithID
	b.durationBase = 40000
	b.durationFactor = 1.4
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 50000, Crystal: 35000, Deuterium: 15000}
	b.Requirements = map[ID]int64{MeditationEnclaveID: 41, RuneForgeID: 1}
	return b
}

type crystalRefinery struct {
	BaseLfBuilding
}

func newCrystalRefinery() *crystalRefinery {
	b := new(crystalRefinery)
	b.Name = "crystal refinery"
	b.ID = CrystalRefineryID
	b.durationBase = 40000
	b.durationFactor = 1.2
	b.IncreaseFactor = 1.40
	b.BaseCost = Resources{Metal: 85000, Crystal: 44000, Deuterium: 25000}
	b.Requirements = map[ID]int64{MeditationEnclaveID: 41, RuneForgeID: 1, MegalithID: 1}
	return b
}

type deuteriumSynthesiser struct {
	BaseLfBuilding
}

func newDeuteriumSynthesiser() *deuteriumSynthesiser {
	b := new(deuteriumSynthesiser)
	b.Name = "deuterium synthesiser"
	b.ID = DeuteriumSynthesiserID
	b.durationBase = 52000
	b.durationFactor = 1.2
	b.IncreaseFactor = 1.40
	b.BaseCost = Resources{Metal: 120000, Crystal: 50000, Deuterium: 20000}
	b.Requirements = map[ID]int64{MeditationEnclaveID: 41, RuneForgeID: 1, MegalithID: 2}
	return b
}

type mineralResearchCentre struct {
	BaseLfBuilding
}

func newMineralResearchCentre() *mineralResearchCentre {
	b := new(mineralResearchCentre)
	b.Name = "mineral research centre"
	b.ID = MineralResearchCentreID
	b.durationBase = 90000
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.80
	b.BaseCost = Resources{Metal: 250000, Crystal: 150000, Deuterium: 100000}
	b.Requirements = map[ID]int64{MeditationEnclaveID: 41, RuneForgeID: 1, MegalithID: 1, CrystalRefineryID: 6, OriktoriumID: 1}
	return b
}

type advancedRecyclingPlant struct {
	BaseLfBuilding
}

func newAdvancedRecyclingPlant() *advancedRecyclingPlant {
	b := new(advancedRecyclingPlant)
	b.Name = "metal recycling plant"
	b.ID = AdvancedRecyclingPlantID
	b.durationBase = 95000
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 250000, Crystal: 125000, Deuterium: 125000}
	b.Requirements = map[ID]int64{MeditationEnclaveID: 41, CrystalFarmID: 22, RuneForgeID: 1, MegalithID: 5, CrystalRefineryID: 6, OriktoriumID: 5, RuneTechnologiumID: 5, MagmaForgeID: 3, DisruptionChamberID: 4, MineralResearchCentreID: 5}
	return b
}

// Mechas
type assemblyLine struct {
	BaseLfBuilding
}

func newAssemblyLine() *assemblyLine {
	b := new(assemblyLine)
	b.Name = "assembly line"
	b.ID = AssemblyLineID
	b.durationBase = 40
	b.durationFactor = 1.22
	b.IncreaseFactor = 1.21
	b.BaseCost = Resources{Metal: 6, Crystal: 2}
	b.Requirements = map[ID]int64{}
	return b
}

type fusionCellFactory struct {
	BaseLfBuilding
}

func newFusionCellFactory() *fusionCellFactory {
	b := new(fusionCellFactory)
	b.Name = "fusion cell factory"
	b.ID = FusionCellFactoryID
	b.durationBase = 48
	b.durationFactor = 1.2
	b.IncreaseFactor = 1.18
	b.energyIncreaseFactor = 1.02
	b.BaseCost = Resources{Metal: 5, Crystal: 2, Energy: 8}
	b.Requirements = map[ID]int64{}
	return b
}

type roboticsResearchCentre struct {
	BaseLfBuilding
}

func newRoboticsResearchCentre() *roboticsResearchCentre {
	b := new(roboticsResearchCentre)
	b.Name = "robotics research centre"
	b.ID = RoboticsResearchCentreID
	b.durationBase = 16000
	b.durationFactor = 1.25
	b.IncreaseFactor = 1.30
	b.BaseCost = Resources{Metal: 30000, Crystal: 20000, Deuterium: 10000}
	b.Requirements = map[ID]int64{AssemblyLineID: 17, FusionCellFactoryID: 20}
	return b
}

type updateNetwork struct {
	BaseLfBuilding
}

func newUpdateNetwork() *updateNetwork {
	b := new(updateNetwork)
	b.Name = "update network"
	b.ID = UpdateNetworkID
	b.durationBase = 16000
	b.durationFactor = 1.6
	b.IncreaseFactor = 1.80
	b.populationIncreaseFactor = 1.10
	b.BaseCost = Resources{Metal: 5000, Crystal: 3800, Deuterium: 1000, Population: 40000000}
	b.Requirements = map[ID]int64{AssemblyLineID: 41}
	return b
}

type quantumComputerCentre struct {
	BaseLfBuilding
}

func newQuantumComputerCentre() *quantumComputerCentre {
	b := new(quantumComputerCentre)
	b.Name = "quantum computer centre"
	b.ID = QuantumComputerCentreID
	b.durationBase = 64000
	b.durationFactor = 1.7
	b.IncreaseFactor = 1.80
	b.populationIncreaseFactor = 1.10
	b.BaseCost = Resources{Metal: 50000, Crystal: 40000, Deuterium: 50000, Population: 130000000}
	b.Requirements = map[ID]int64{AssemblyLineID: 41, UpdateNetworkID: 1, MicrochipAssemblyLineID: 1, ProductionAssemblyHallID: 5}
	return b
}

type automatisedAssemblyCentre struct {
	BaseLfBuilding
}

func newAutomatisedAssemblyCentre() *automatisedAssemblyCentre {
	b := new(automatisedAssemblyCentre)
	b.Name = "automatised assembly centre"
	b.ID = AutomatisedAssemblyCentreID
	b.durationBase = 2000
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.30
	b.BaseCost = Resources{Metal: 7500, Crystal: 7000, Deuterium: 1000}
	b.Requirements = map[ID]int64{AssemblyLineID: 17, FusionCellFactoryID: 20, RoboticsResearchCentreID: 5}
	return b
}

type highPerformanceTransformer struct {
	BaseLfBuilding
}

func newHighPerformanceTransformer() *highPerformanceTransformer {
	b := new(highPerformanceTransformer)
	b.Name = "high performance transformer"
	b.ID = HighPerformanceTransformerID
	b.durationBase = 16000
	b.durationFactor = 1.4
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 35000, Crystal: 15000, Deuterium: 10000}
	b.Requirements = map[ID]int64{AssemblyLineID: 17, FusionCellFactoryID: 20, RoboticsResearchCentreID: 5, AutomatisedAssemblyCentreID: 3}
	return b
}

type microchipAssemblyLine struct {
	BaseLfBuilding
}

func newMicrochipAssemblyLine() *microchipAssemblyLine {
	b := new(microchipAssemblyLine)
	b.Name = "microchip assembly line"
	b.ID = MicrochipAssemblyLineID
	b.durationBase = 12000
	b.durationFactor = 1.17
	b.IncreaseFactor = 1.07
	b.BaseCost = Resources{Metal: 50000, Crystal: 20000, Deuterium: 30000}
	b.Requirements = map[ID]int64{AssemblyLineID: 41, UpdateNetworkID: 1}
	return b
}

type productionAssemblyHall struct {
	BaseLfBuilding
}

func newProductionAssemblyHall() *productionAssemblyHall {
	b := new(productionAssemblyHall)
	b.Name = "production assembly hall"
	b.ID = ProductionAssemblyHallID
	b.durationBase = 40000
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.14
	b.BaseCost = Resources{Metal: 100000, Crystal: 10000, Deuterium: 3000}
	b.Requirements = map[ID]int64{AssemblyLineID: 41, UpdateNetworkID: 1, MicrochipAssemblyLineID: 1}
	return b
}

type highPerformanceSynthesiser struct {
	BaseLfBuilding
}

func newHighPerformanceSynthesiser() *highPerformanceSynthesiser {
	b := new(highPerformanceSynthesiser)
	b.Name = "high performance synthesiser"
	b.ID = HighPerformanceSynthesiserID
	b.durationBase = 52000
	b.durationFactor = 1.2
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 100000, Crystal: 40000, Deuterium: 20000}
	b.Requirements = map[ID]int64{AssemblyLineID: 41, UpdateNetworkID: 1, MicrochipAssemblyLineID: 2}
	return b
}

type chipMassProduction struct {
	BaseLfBuilding
}

func newChipMassProduction() *chipMassProduction {
	b := new(chipMassProduction)
	b.Name = "chip mass production"
	b.ID = ChipMassProductionID
	b.durationBase = 50000
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 55000, Crystal: 50000, Deuterium: 30000}
	b.Requirements = map[ID]int64{AssemblyLineID: 41, UpdateNetworkID: 1, MicrochipAssemblyLineID: 1, ProductionAssemblyHallID: 6, QuantumComputerCentreID: 1}
	return b
}

type nanoRepairBots struct {
	BaseLfBuilding
}

func newNanoRepairBots() *nanoRepairBots {
	b := new(nanoRepairBots)
	b.Name = "nano repair bots"
	b.ID = NanoRepairBotsID
	b.durationBase = 95000
	b.durationFactor = 1.4
	b.IncreaseFactor = 1.40
	b.BaseCost = Resources{Metal: 250000, Crystal: 125000, Deuterium: 125000}
	b.Requirements = map[ID]int64{AssemblyLineID: 41, FusionCellFactoryID: 20, MicrochipAssemblyLineID: 5, RoboticsResearchCentreID: 5, HighPerformanceTransformerID: 4, ProductionAssemblyHallID: 6, QuantumComputerCentreID: 5, ChipMassProductionID: 11}
	return b
}

// Kaelesh
type sanctuary struct {
	BaseLfBuilding
}

func newSanctuary() *sanctuary {
	b := new(sanctuary)
	b.Name = "sanctuary"
	b.ID = SanctuaryID
	b.durationBase = 40
	b.durationFactor = 1.22
	b.IncreaseFactor = 1.21
	b.BaseCost = Resources{Metal: 4, Crystal: 3}
	b.Requirements = map[ID]int64{}
	return b
}

type antimatterCondenser struct {
	BaseLfBuilding
}

func newAntimatterCondenser() *antimatterCondenser {
	b := new(antimatterCondenser)
	b.Name = "antimatter condenser"
	b.ID = AntimatterCondenserID
	b.durationBase = 40
	b.durationFactor = 1.22
	b.IncreaseFactor = 1.20
	b.energyIncreaseFactor = 1.02
	b.BaseCost = Resources{Metal: 6, Crystal: 3, Energy: 9}
	b.Requirements = map[ID]int64{}
	return b
}

type vortexChamber struct {
	BaseLfBuilding
}

func newVortexChamber() *vortexChamber {
	b := new(vortexChamber)
	b.Name = "vortex chamber"
	b.ID = VortexChamberID
	b.durationBase = 16000
	b.durationFactor = 1.25
	b.IncreaseFactor = 1.30
	b.BaseCost = Resources{Metal: 20000, Crystal: 20000, Deuterium: 30000}
	b.Requirements = map[ID]int64{SanctuaryID: 20, AntimatterCondenserID: 21}
	return b
}

type hallsOfRealisation struct {
	BaseLfBuilding
}

func newHallsOfRealisation() *hallsOfRealisation {
	b := new(hallsOfRealisation)
	b.Name = "halls of realisation"
	b.ID = HallsOfRealisationID
	b.durationBase = 16000
	b.durationFactor = 1.7
	b.IncreaseFactor = 1.80
	b.populationIncreaseFactor = 1.10
	b.BaseCost = Resources{Metal: 7500, Crystal: 5000, Deuterium: 800, Population: 30000000}
	b.Requirements = map[ID]int64{SanctuaryID: 42}
	return b
}

type forumOfTranscendence struct {
	BaseLfBuilding
}

func newForumOfTranscendence() *forumOfTranscendence {
	b := new(forumOfTranscendence)
	b.Name = "forum of transcendence"
	b.ID = ForumOfTranscendenceID
	b.durationBase = 64000
	b.durationFactor = 1.8
	b.IncreaseFactor = 1.80
	b.populationIncreaseFactor = 1.10
	b.BaseCost = Resources{Metal: 60000, Crystal: 30000, Deuterium: 50000, Population: 100000000}
	b.Requirements = map[ID]int64{SanctuaryID: 42, HallsOfRealisationID: 1, ChrysalisAcceleratorID: 1, BioModifierID: 5}
	return b
}

type antimatterConvector struct {
	BaseLfBuilding
}

func newAntimatterConvector() *antimatterConvector {
	b := new(antimatterConvector)
	b.Name = "antimatter convector"
	b.ID = AntimatterConvectorID
	b.durationBase = 2000
	b.durationFactor = 1.35
	b.IncreaseFactor = 1.25
	b.BaseCost = Resources{Metal: 8500, Crystal: 5000, Deuterium: 3000}
	b.Requirements = map[ID]int64{SanctuaryID: 20, AntimatterCondenserID: 21, VortexChamberID: 5}
	return b
}

type cloningLaboratory struct {
	BaseLfBuilding
}

func newCloningLaboratory() *cloningLaboratory {
	b := new(cloningLaboratory)
	b.Name = "cloning laboratory"
	b.ID = CloningLaboratoryID
	b.durationBase = 12000
	b.durationFactor = 1.20
	b.IncreaseFactor = 1.20
	b.BaseCost = Resources{Metal: 15000, Crystal: 15000, Deuterium: 5000}
	b.Requirements = map[ID]int64{SanctuaryID: 20, AntimatterCondenserID: 21, VortexChamberID: 5, AntimatterConvectorID: 3}
	return b
}

type chrysalisAccelerator struct {
	BaseLfBuilding
}

func newChrysalisAccelerator() *chrysalisAccelerator {
	b := new(chrysalisAccelerator)
	b.Name = "chrysalis accelerator"
	b.ID = ChrysalisAcceleratorID
	b.durationBase = 16000
	b.durationFactor = 1.18
	b.IncreaseFactor = 1.05
	b.BaseCost = Resources{Metal: 75000, Crystal: 25000, Deuterium: 30000}
	b.Requirements = map[ID]int64{SanctuaryID: 42, HallsOfRealisationID: 1}
	return b
}

type bioModifier struct {
	BaseLfBuilding
}

func newBioModifier() *bioModifier {
	b := new(bioModifier)
	b.Name = "bio modifier"
	b.ID = BioModifierID
	b.durationBase = 40000
	b.durationFactor = 1.2
	b.IncreaseFactor = 1.20
	b.BaseCost = Resources{Metal: 87500, Crystal: 25000, Deuterium: 30000}
	b.Requirements = map[ID]int64{SanctuaryID: 42, HallsOfRealisationID: 1, ChrysalisAcceleratorID: 1}
	return b
}

type psionicModulator struct {
	BaseLfBuilding
}

func newPsionicModulator() *psionicModulator {
	b := new(psionicModulator)
	b.Name = "psionic modulator"
	b.ID = PsionicModulatorID
	b.durationBase = 52000
	b.durationFactor = 1.8
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 150000, Crystal: 30000, Deuterium: 30000}
	b.Requirements = map[ID]int64{SanctuaryID: 42, HallsOfRealisationID: 1, ChrysalisAcceleratorID: 2}
	return b
}

type shipManufacturingHall struct {
	BaseLfBuilding
}

func newShipManufacturingHall() *shipManufacturingHall {
	b := new(shipManufacturingHall)
	b.Name = "ship manufacturing hall"
	b.ID = ShipManufacturingHallID
	b.durationBase = 90000
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.20
	b.BaseCost = Resources{Metal: 75000, Crystal: 50000, Deuterium: 55000}
	b.Requirements = map[ID]int64{SanctuaryID: 42, HallsOfRealisationID: 1, ChrysalisAcceleratorID: 1, BioModifierID: 6, ForumOfTranscendenceID: 1}
	return b
}

type supraRefractor struct {
	BaseLfBuilding
}

func newSupraRefractor() *supraRefractor {
	b := new(supraRefractor)
	b.Name = "suprarefractor"
	b.ID = SupraRefractorID
	b.durationBase = 95000
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.40
	b.BaseCost = Resources{Metal: 500000, Crystal: 250000, Deuterium: 250000}
	b.Requirements = map[ID]int64{SanctuaryID: 42, AntimatterCondenserID: 21, VortexChamberID: 5, AntimatterConvectorID: 3, CloningLaboratoryID: 4, HallsOfRealisationID: 1, ChrysalisAcceleratorID: 5, BioModifierID: 6, ForumOfTranscendenceID: 5, ShipManufacturingHallID: 5}
	return b
}
