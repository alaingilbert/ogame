package ogame

import (
	"time"
)

type Celestial interface {
	GetCoordinate() Coordinate
	GetDiameter() int64
	GetFields() Fields
	GetID() CelestialID
	GetImg() string
	GetName() string
	GetType() CelestialType
}

// BaseOgameObj base interface for all ogame objects (buildings, technologies, ships, defenses)
type BaseOgameObj interface {
	ConstructionTime(nbr, universeSpeed int64, acc BuildAccelerators, lfBonuses LfBonuses, class CharacterClass, hasTechnocrat bool) time.Duration
	GetID() ID
	GetName() string
	GetPrice(int64, LfBonuses) Resources
	GetRequirements() map[ID]int64
	IsAvailable(CelestialType, IResourcesBuildings, ILfBuildings, ILfResearches, IFacilities, IResearches, int64, CharacterClass) bool
}

// DefenderObj base interface for all defensive units (ships, defenses)
type DefenderObj interface {
	BaseOgameObj
	DefenderConstructionTime(nbr, universeSpeed int64, acc DefenseAccelerators, lfBonuses LfBonuses) time.Duration
	GetRapidfireAgainst() map[ID]int64
	GetRapidfireFrom() map[ID]int64
	GetShieldPower(IResearches) int64
	GetStructuralIntegrity(IResearches) int64
	GetWeaponPower(IResearches) int64
}

// Ship interface implemented by all ships units
type Ship interface {
	DefenderObj
	GetCargoCapacity(techs IResearches, bonus LfBonuses, characterClass CharacterClass, multiplier float64, probeRaids bool) int64
	GetFuelConsumption(techs IResearches, bonus LfBonuses, characterClass CharacterClass, fleetDeutSaveFactor float64) int64
	GetSpeed(techs IResearches, bonus LfBonuses, characterClass CharacterClass) int64
}

// Defense interface implemented by all defenses units
type Defense interface {
	DefenderObj
}

// Levelable base interface for all levelable ogame objects (buildings, technologies)
type Levelable interface {
	BaseOgameObj
	GetLevel(IResourcesBuildings, IFacilities, IResearches) int64
}

// Technology interface that all technologies implement
type Technology interface {
	Levelable
	TechnologyConstructionTime(nbr, universeSpeed int64, acc TechAccelerators, lfBonuses LfBonuses, class CharacterClass, hasTechnocrat bool) time.Duration
}

// Building interface that all buildings implement
type Building interface {
	Levelable
	BuildingConstructionTime(nbr, universeSpeed int64, acc BuildingAccelerators, lfBonuses LfBonuses) time.Duration
	DeconstructionPrice(lvl int64, techs IResearches) Resources
}

// BuildAccelerators levels of things we need to calculate construction time of anything
type BuildAccelerators interface {
	TechAccelerators
	BuildingAccelerators
	DefenseAccelerators
}

// TechAccelerators to calculate techs construction time, we need research lab level
type TechAccelerators interface {
	GetResearchLab() int64
}

// DefenseAccelerators to calculate defense construction time (ships / defenses), we need nanite and shipyard levels
type DefenseAccelerators interface {
	GetNaniteFactory() int64
	GetShipyard() int64
}

// BuildingAccelerators to calculate building construction time, we need nanite and robotic levels
type BuildingAccelerators interface {
	GetNaniteFactory() int64
	GetRoboticsFactory() int64
}

type IFacilities interface {
	ByID(ID) int64
	GetRoboticsFactory() int64
	GetShipyard() int64
	GetResearchLab() int64
	GetAllianceDepot() int64
	GetMissileSilo() int64
	GetNaniteFactory() int64
	GetTerraformer() int64
	GetSpaceDock() int64
	GetLunarBase() int64
	GetSensorPhalanx() int64
	GetJumpGate() int64
}

type IResearches interface {
	ByID(ID) int64
	GetEnergyTechnology() int64
	GetLaserTechnology() int64
	GetIonTechnology() int64
	GetHyperspaceTechnology() int64
	GetPlasmaTechnology() int64
	GetCombustionDrive() int64
	GetImpulseDrive() int64
	GetHyperspaceDrive() int64
	GetEspionageTechnology() int64
	GetComputerTechnology() int64
	GetAstrophysics() int64
	GetIntergalacticResearchNetwork() int64
	GetGravitonTechnology() int64
	GetWeaponsTechnology() int64
	GetShieldingTechnology() int64
	GetArmourTechnology() int64
}

type IResourcesBuildings interface {
	ByID(ID) int64
	GetMetalMine() int64
	GetCrystalMine() int64
	GetDeuteriumSynthesizer() int64
	GetSolarPlant() int64
	GetFusionReactor() int64
	GetSolarSatellite() int64
	GetMetalStorage() int64
	GetCrystalStorage() int64
	GetDeuteriumTank() int64
}

type ILfBuildings interface {
	ByID(ID) int64
	GetResidentialSector() int64
	GetBiosphereFarm() int64
	GetResearchCentre() int64
	GetAcademyOfSciences() int64
	GetNeuroCalibrationCentre() int64
	GetHighEnergySmelting() int64
	GetFoodSilo() int64
	GetFusionPoweredProduction() int64
	GetSkyscraper() int64
	GetBiotechLab() int64
	GetMetropolis() int64
	GetPlanetaryShield() int64
	GetMeditationEnclave() int64
	GetCrystalFarm() int64
	GetRuneTechnologium() int64
	GetRuneForge() int64
	GetOriktorium() int64
	GetMagmaForge() int64
	GetDisruptionChamber() int64
	GetMegalith() int64
	GetCrystalRefinery() int64
	GetDeuteriumSynthesiser() int64
	GetMineralResearchCentre() int64
	GetAdvancedRecyclingPlant() int64
	GetAssemblyLine() int64
	GetFusionCellFactory() int64
	GetRoboticsResearchCentre() int64
	GetUpdateNetwork() int64
	GetQuantumComputerCentre() int64
	GetAutomatisedAssemblyCentre() int64
	GetHighPerformanceTransformer() int64
	GetMicrochipAssemblyLine() int64
	GetProductionAssemblyHall() int64
	GetHighPerformanceSynthesiser() int64
	GetChipMassProduction() int64
	GetNanoRepairBots() int64
	GetSanctuary() int64
	GetAntimatterCondenser() int64
	GetVortexChamber() int64
	GetHallsOfRealisation() int64
	GetForumOfTranscendence() int64
	GetAntimatterConvector() int64
	GetCloningLaboratory() int64
	GetChrysalisAccelerator() int64
	GetBioModifier() int64
	GetPsionicModulator() int64
	GetShipManufacturingHall() int64
	GetSupraRefractor() int64
}

type ILfResearches interface {
	ByID(ID) *int64
	GetIntergalacticEnvoys() *int64
	GetHighPerformanceExtractors() *int64
	GetFusionDrives() *int64
	GetStealthFieldGenerator() *int64
	GetOrbitalDen() *int64
	GetResearchAI() *int64
	GetHighPerformanceTerraformer() *int64
	GetEnhancedProductionTechnologies() *int64
	GetLightFighterMkII() *int64
	GetCruiserMkII() *int64
	GetImprovedLabTechnology() *int64
	GetPlasmaTerraformer() *int64
	GetLowTemperatureDrives() *int64
	GetBomberMkII() *int64
	GetDestroyerMkII() *int64
	GetBattlecruiserMkII() *int64
	GetRobotAssistants() *int64
	GetSupercomputer() *int64
	GetVolcanicBatteries() *int64
	GetAcousticScanning() *int64
	GetHighEnergyPumpSystems() *int64
	GetCargoHoldExpansionCivilianShips() *int64
	GetMagmaPoweredProduction() *int64
	GetGeothermalPowerPlants() *int64
	GetDepthSounding() *int64
	GetIonCrystalEnhancementHeavyFighter() *int64
	GetImprovedStellarator() *int64
	GetHardenedDiamondDrillHeads() *int64
	GetSeismicMiningTechnology() *int64
	GetMagmaPoweredPumpSystems() *int64
	GetIonCrystalModules() *int64
	GetOptimisedSiloConstructionMethod() *int64
	GetDiamondEnergyTransmitter() *int64
	GetObsidianShieldReinforcement() *int64
	GetRuneShields() *int64
	GetRocktalCollectorEnhancement() *int64
	GetCatalyserTechnology() *int64
	GetPlasmaDrive() *int64
	GetEfficiencyModule() *int64
	GetDepotAI() *int64
	GetGeneralOverhaulLightFighter() *int64
	GetAutomatedTransportLines() *int64
	GetImprovedDroneAI() *int64
	GetExperimentalRecyclingTechnology() *int64
	GetGeneralOverhaulCruiser() *int64
	GetSlingshotAutopilot() *int64
	GetHighTemperatureSuperconductors() *int64
	GetGeneralOverhaulBattleship() *int64
	GetArtificialSwarmIntelligence() *int64
	GetGeneralOverhaulBattlecruiser() *int64
	GetGeneralOverhaulBomber() *int64
	GetGeneralOverhaulDestroyer() *int64
	GetExperimentalWeaponsTechnology() *int64
	GetMechanGeneralEnhancement() *int64
	GetHeatRecovery() *int64
	GetSulphideProcess() *int64
	GetPsionicNetwork() *int64
	GetTelekineticTractorBeam() *int64
	GetEnhancedSensorTechnology() *int64
	GetNeuromodalCompressor() *int64
	GetNeuroInterface() *int64
	GetInterplanetaryAnalysisNetwork() *int64
	GetOverclockingHeavyFighter() *int64
	GetTelekineticDrive() *int64
	GetSixthSense() *int64
	GetPsychoharmoniser() *int64
	GetEfficientSwarmIntelligence() *int64
	GetOverclockingLargeCargo() *int64
	GetGravitationSensors() *int64
	GetOverclockingBattleship() *int64
	GetPsionicShieldMatrix() *int64
	GetKaeleshDiscovererEnhancement() *int64
}
