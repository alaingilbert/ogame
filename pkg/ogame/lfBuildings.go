package ogame

import (
	"math"
)

// LazyLfBuildings ...
type LazyLfBuildings func() LfBuildings

// LfBuildings lifeform buildings
type LfBuildings struct {
	ResidentialSector       int64 // 11101
	BiosphereFarm           int64 // 11102
	ResearchCentre          int64 // 11103
	AcademyOfSciences       int64 // 11104
	NeuroCalibrationCentre  int64 // 11105
	HighEnergySmelting      int64 // 11106
	FoodSilo                int64 // 11107
	FusionPoweredProduction int64 // 11108
	Skyscraper              int64 // 11109
	BiotechLab              int64 // 11110
	Metropolis              int64 // 11111
	PlanetaryShield         int64 // 11112
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
