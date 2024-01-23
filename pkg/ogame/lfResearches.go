package ogame

import (
	"math"
	"time"
)

// LazyLfResearches ...
type LazyLfResearches func() LfResearches

func (b LazyLfResearches) ByID(id ID) int64                    { return b().ByID(id) }
func (b LazyLfResearches) GetIntergalacticEnvoys() int64       { return b().IntergalacticEnvoys }
func (b LazyLfResearches) GetHighPerformanceExtractors() int64 { return b().HighPerformanceExtractors }
func (b LazyLfResearches) GetFusionDrives() int64              { return b().FusionDrives }
func (b LazyLfResearches) GetStealthFieldGenerator() int64     { return b().StealthFieldGenerator }
func (b LazyLfResearches) GetOrbitalDen() int64                { return b().OrbitalDen }
func (b LazyLfResearches) GetResearchAI() int64                { return b().ResearchAI }
func (b LazyLfResearches) GetHighPerformanceTerraformer() int64 {
	return b().HighPerformanceTerraformer
}
func (b LazyLfResearches) GetEnhancedProductionTechnologies() int64 {
	return b().EnhancedProductionTechnologies
}
func (b LazyLfResearches) GetLightFighterMkII() int64      { return b().LightFighterMkII }
func (b LazyLfResearches) GetCruiserMkII() int64           { return b().CruiserMkII }
func (b LazyLfResearches) GetImprovedLabTechnology() int64 { return b().ImprovedLabTechnology }
func (b LazyLfResearches) GetPlasmaTerraformer() int64     { return b().PlasmaTerraformer }
func (b LazyLfResearches) GetLowTemperatureDrives() int64  { return b().LowTemperatureDrives }
func (b LazyLfResearches) GetBomberMkII() int64            { return b().BomberMkII }
func (b LazyLfResearches) GetDestroyerMkII() int64         { return b().DestroyerMkII }
func (b LazyLfResearches) GetBattlecruiserMkII() int64     { return b().BattlecruiserMkII }
func (b LazyLfResearches) GetRobotAssistants() int64       { return b().RobotAssistants }
func (b LazyLfResearches) GetSupercomputer() int64         { return b().Supercomputer }
func (b LazyLfResearches) GetVolcanicBatteries() int64     { return b().VolcanicBatteries }
func (b LazyLfResearches) GetAcousticScanning() int64      { return b().AcousticScanning }
func (b LazyLfResearches) GetHighEnergyPumpSystems() int64 { return b().HighEnergyPumpSystems }
func (b LazyLfResearches) GetCargoHoldExpansionCivilianShips() int64 {
	return b().CargoHoldExpansionCivilianShips
}
func (b LazyLfResearches) GetMagmaPoweredProduction() int64 { return b().MagmaPoweredProduction }
func (b LazyLfResearches) GetGeothermalPowerPlants() int64  { return b().GeothermalPowerPlants }
func (b LazyLfResearches) GetDepthSounding() int64          { return b().DepthSounding }
func (b LazyLfResearches) GetIonCrystalEnhancementHeavyFighter() int64 {
	return b().IonCrystalEnhancementHeavyFighter
}
func (b LazyLfResearches) GetImprovedStellarator() int64       { return b().ImprovedStellarator }
func (b LazyLfResearches) GetHardenedDiamondDrillHeads() int64 { return b().HardenedDiamondDrillHeads }
func (b LazyLfResearches) GetSeismicMiningTechnology() int64   { return b().SeismicMiningTechnology }
func (b LazyLfResearches) GetMagmaPoweredPumpSystems() int64   { return b().MagmaPoweredPumpSystems }
func (b LazyLfResearches) GetIonCrystalModules() int64         { return b().IonCrystalModules }
func (b LazyLfResearches) GetOptimisedSiloConstructionMethod() int64 {
	return b().OptimisedSiloConstructionMethod
}
func (b LazyLfResearches) GetDiamondEnergyTransmitter() int64 { return b().DiamondEnergyTransmitter }
func (b LazyLfResearches) GetObsidianShieldReinforcement() int64 {
	return b().ObsidianShieldReinforcement
}
func (b LazyLfResearches) GetRuneShields() int64 { return b().RuneShields }
func (b LazyLfResearches) GetRocktalCollectorEnhancement() int64 {
	return b().RocktalCollectorEnhancement
}
func (b LazyLfResearches) GetCatalyserTechnology() int64 { return b().CatalyserTechnology }
func (b LazyLfResearches) GetPlasmaDrive() int64         { return b().PlasmaDrive }
func (b LazyLfResearches) GetEfficiencyModule() int64    { return b().EfficiencyModule }
func (b LazyLfResearches) GetDepotAI() int64             { return b().DepotAI }
func (b LazyLfResearches) GetGeneralOverhaulLightFighter() int64 {
	return b().GeneralOverhaulLightFighter
}
func (b LazyLfResearches) GetAutomatedTransportLines() int64 { return b().AutomatedTransportLines }
func (b LazyLfResearches) GetImprovedDroneAI() int64         { return b().ImprovedDroneAI }
func (b LazyLfResearches) GetExperimentalRecyclingTechnology() int64 {
	return b().ExperimentalRecyclingTechnology
}
func (b LazyLfResearches) GetGeneralOverhaulCruiser() int64 { return b().GeneralOverhaulCruiser }
func (b LazyLfResearches) GetSlingshotAutopilot() int64     { return b().SlingshotAutopilot }
func (b LazyLfResearches) GetHighTemperatureSuperconductors() int64 {
	return b().HighTemperatureSuperconductors
}
func (b LazyLfResearches) GetGeneralOverhaulBattleship() int64 { return b().GeneralOverhaulBattleship }
func (b LazyLfResearches) GetArtificialSwarmIntelligence() int64 {
	return b().ArtificialSwarmIntelligence
}
func (b LazyLfResearches) GetGeneralOverhaulBattlecruiser() int64 {
	return b().GeneralOverhaulBattlecruiser
}
func (b LazyLfResearches) GetGeneralOverhaulBomber() int64    { return b().GeneralOverhaulBomber }
func (b LazyLfResearches) GetGeneralOverhaulDestroyer() int64 { return b().GeneralOverhaulDestroyer }
func (b LazyLfResearches) GetExperimentalWeaponsTechnology() int64 {
	return b().ExperimentalWeaponsTechnology
}
func (b LazyLfResearches) GetMechanGeneralEnhancement() int64 { return b().MechanGeneralEnhancement }
func (b LazyLfResearches) GetHeatRecovery() int64             { return b().HeatRecovery }
func (b LazyLfResearches) GetSulphideProcess() int64          { return b().SulphideProcess }
func (b LazyLfResearches) GetPsionicNetwork() int64           { return b().PsionicNetwork }
func (b LazyLfResearches) GetTelekineticTractorBeam() int64   { return b().TelekineticTractorBeam }
func (b LazyLfResearches) GetEnhancedSensorTechnology() int64 { return b().EnhancedSensorTechnology }
func (b LazyLfResearches) GetNeuromodalCompressor() int64     { return b().NeuromodalCompressor }
func (b LazyLfResearches) GetNeuroInterface() int64           { return b().NeuroInterface }
func (b LazyLfResearches) GetInterplanetaryAnalysisNetwork() int64 {
	return b().InterplanetaryAnalysisNetwork
}
func (b LazyLfResearches) GetOverclockingHeavyFighter() int64 { return b().OverclockingHeavyFighter }
func (b LazyLfResearches) GetTelekineticDrive() int64         { return b().TelekineticDrive }
func (b LazyLfResearches) GetSixthSense() int64               { return b().SixthSense }
func (b LazyLfResearches) GetPsychoharmoniser() int64         { return b().Psychoharmoniser }
func (b LazyLfResearches) GetEfficientSwarmIntelligence() int64 {
	return b().EfficientSwarmIntelligence
}
func (b LazyLfResearches) GetOverclockingLargeCargo() int64 { return b().OverclockingLargeCargo }
func (b LazyLfResearches) GetGravitationSensors() int64     { return b().GravitationSensors }
func (b LazyLfResearches) GetOverclockingBattleship() int64 { return b().OverclockingBattleship }
func (b LazyLfResearches) GetPsionicShieldMatrix() int64    { return b().PsionicShieldMatrix }
func (b LazyLfResearches) GetKaeleshDiscovererEnhancement() int64 {
	return b().KaeleshDiscovererEnhancement
}

type LfResearches struct {
	IntergalacticEnvoys               int64 // 11201 // Humans techs
	HighPerformanceExtractors         int64 // 11202
	FusionDrives                      int64 // 11203
	StealthFieldGenerator             int64 // 11204
	OrbitalDen                        int64 // 11205
	ResearchAI                        int64 // 11206
	HighPerformanceTerraformer        int64 // 11207
	EnhancedProductionTechnologies    int64 // 11208
	LightFighterMkII                  int64 // 11209
	CruiserMkII                       int64 // 11210
	ImprovedLabTechnology             int64 // 11211
	PlasmaTerraformer                 int64 // 11212
	LowTemperatureDrives              int64 // 11213
	BomberMkII                        int64 // 11214
	DestroyerMkII                     int64 // 11215
	BattlecruiserMkII                 int64 // 11216
	RobotAssistants                   int64 // 11217
	Supercomputer                     int64 // 11218
	VolcanicBatteries                 int64 // 12201 // Rocktal techs
	AcousticScanning                  int64 // 12202
	HighEnergyPumpSystems             int64 // 12203
	CargoHoldExpansionCivilianShips   int64 // 12204
	MagmaPoweredProduction            int64 // 12205
	GeothermalPowerPlants             int64 // 12206
	DepthSounding                     int64 // 12207
	IonCrystalEnhancementHeavyFighter int64 // 12208
	ImprovedStellarator               int64 // 12209
	HardenedDiamondDrillHeads         int64 // 12210
	SeismicMiningTechnology           int64 // 12211
	MagmaPoweredPumpSystems           int64 // 12212
	IonCrystalModules                 int64 // 12213
	OptimisedSiloConstructionMethod   int64 // 12214
	DiamondEnergyTransmitter          int64 // 12215
	ObsidianShieldReinforcement       int64 // 12216
	RuneShields                       int64 // 12217
	RocktalCollectorEnhancement       int64 // 12218
	CatalyserTechnology               int64 // 13201 // Mechas techs
	PlasmaDrive                       int64 // 13202
	EfficiencyModule                  int64 // 13203
	DepotAI                           int64 // 13204
	GeneralOverhaulLightFighter       int64 // 13205
	AutomatedTransportLines           int64 // 13206
	ImprovedDroneAI                   int64 // 13207
	ExperimentalRecyclingTechnology   int64 // 13208
	GeneralOverhaulCruiser            int64 // 13209
	SlingshotAutopilot                int64 // 13210
	HighTemperatureSuperconductors    int64 // 13211
	GeneralOverhaulBattleship         int64 // 13212
	ArtificialSwarmIntelligence       int64 // 13213
	GeneralOverhaulBattlecruiser      int64 // 13214
	GeneralOverhaulBomber             int64 // 13215
	GeneralOverhaulDestroyer          int64 // 13216
	ExperimentalWeaponsTechnology     int64 // 13217
	MechanGeneralEnhancement          int64 // 13218
	HeatRecovery                      int64 // 14201 // Kaelesh techs
	SulphideProcess                   int64 // 14202
	PsionicNetwork                    int64 // 14203
	TelekineticTractorBeam            int64 // 14204
	EnhancedSensorTechnology          int64 // 14205
	NeuromodalCompressor              int64 // 14206
	NeuroInterface                    int64 // 14207
	InterplanetaryAnalysisNetwork     int64 // 14208
	OverclockingHeavyFighter          int64 // 14209
	TelekineticDrive                  int64 // 14210
	SixthSense                        int64 // 14211
	Psychoharmoniser                  int64 // 14212
	EfficientSwarmIntelligence        int64 // 14213
	OverclockingLargeCargo            int64 // 14214
	GravitationSensors                int64 // 14215
	OverclockingBattleship            int64 // 14216
	PsionicShieldMatrix               int64 // 14217
	KaeleshDiscovererEnhancement      int64 // 14218
}

func (b LfResearches) GetIntergalacticEnvoys() int64        { return b.IntergalacticEnvoys }
func (b LfResearches) GetHighPerformanceExtractors() int64  { return b.HighPerformanceExtractors }
func (b LfResearches) GetFusionDrives() int64               { return b.FusionDrives }
func (b LfResearches) GetStealthFieldGenerator() int64      { return b.StealthFieldGenerator }
func (b LfResearches) GetOrbitalDen() int64                 { return b.OrbitalDen }
func (b LfResearches) GetResearchAI() int64                 { return b.ResearchAI }
func (b LfResearches) GetHighPerformanceTerraformer() int64 { return b.HighPerformanceTerraformer }
func (b LfResearches) GetEnhancedProductionTechnologies() int64 {
	return b.EnhancedProductionTechnologies
}
func (b LfResearches) GetLightFighterMkII() int64      { return b.LightFighterMkII }
func (b LfResearches) GetCruiserMkII() int64           { return b.CruiserMkII }
func (b LfResearches) GetImprovedLabTechnology() int64 { return b.ImprovedLabTechnology }
func (b LfResearches) GetPlasmaTerraformer() int64     { return b.PlasmaTerraformer }
func (b LfResearches) GetLowTemperatureDrives() int64  { return b.LowTemperatureDrives }
func (b LfResearches) GetBomberMkII() int64            { return b.BomberMkII }
func (b LfResearches) GetDestroyerMkII() int64         { return b.DestroyerMkII }
func (b LfResearches) GetBattlecruiserMkII() int64     { return b.BattlecruiserMkII }
func (b LfResearches) GetRobotAssistants() int64       { return b.RobotAssistants }
func (b LfResearches) GetSupercomputer() int64         { return b.Supercomputer }
func (b LfResearches) GetVolcanicBatteries() int64     { return b.VolcanicBatteries }
func (b LfResearches) GetAcousticScanning() int64      { return b.AcousticScanning }
func (b LfResearches) GetHighEnergyPumpSystems() int64 { return b.HighEnergyPumpSystems }
func (b LfResearches) GetCargoHoldExpansionCivilianShips() int64 {
	return b.CargoHoldExpansionCivilianShips
}
func (b LfResearches) GetMagmaPoweredProduction() int64 { return b.MagmaPoweredProduction }
func (b LfResearches) GetGeothermalPowerPlants() int64  { return b.GeothermalPowerPlants }
func (b LfResearches) GetDepthSounding() int64          { return b.DepthSounding }
func (b LfResearches) GetIonCrystalEnhancementHeavyFighter() int64 {
	return b.IonCrystalEnhancementHeavyFighter
}
func (b LfResearches) GetImprovedStellarator() int64       { return b.ImprovedStellarator }
func (b LfResearches) GetHardenedDiamondDrillHeads() int64 { return b.HardenedDiamondDrillHeads }
func (b LfResearches) GetSeismicMiningTechnology() int64   { return b.SeismicMiningTechnology }
func (b LfResearches) GetMagmaPoweredPumpSystems() int64   { return b.MagmaPoweredPumpSystems }
func (b LfResearches) GetIonCrystalModules() int64         { return b.IonCrystalModules }
func (b LfResearches) GetOptimisedSiloConstructionMethod() int64 {
	return b.OptimisedSiloConstructionMethod
}
func (b LfResearches) GetDiamondEnergyTransmitter() int64    { return b.DiamondEnergyTransmitter }
func (b LfResearches) GetObsidianShieldReinforcement() int64 { return b.ObsidianShieldReinforcement }
func (b LfResearches) GetRuneShields() int64                 { return b.RuneShields }
func (b LfResearches) GetRocktalCollectorEnhancement() int64 { return b.RocktalCollectorEnhancement }
func (b LfResearches) GetCatalyserTechnology() int64         { return b.CatalyserTechnology }
func (b LfResearches) GetPlasmaDrive() int64                 { return b.PlasmaDrive }
func (b LfResearches) GetEfficiencyModule() int64            { return b.EfficiencyModule }
func (b LfResearches) GetDepotAI() int64                     { return b.DepotAI }
func (b LfResearches) GetGeneralOverhaulLightFighter() int64 { return b.GeneralOverhaulLightFighter }
func (b LfResearches) GetAutomatedTransportLines() int64     { return b.AutomatedTransportLines }
func (b LfResearches) GetImprovedDroneAI() int64             { return b.ImprovedDroneAI }
func (b LfResearches) GetExperimentalRecyclingTechnology() int64 {
	return b.ExperimentalRecyclingTechnology
}
func (b LfResearches) GetGeneralOverhaulCruiser() int64 { return b.GeneralOverhaulCruiser }
func (b LfResearches) GetSlingshotAutopilot() int64     { return b.SlingshotAutopilot }
func (b LfResearches) GetHighTemperatureSuperconductors() int64 {
	return b.HighTemperatureSuperconductors
}
func (b LfResearches) GetGeneralOverhaulBattleship() int64    { return b.GeneralOverhaulBattleship }
func (b LfResearches) GetArtificialSwarmIntelligence() int64  { return b.ArtificialSwarmIntelligence }
func (b LfResearches) GetGeneralOverhaulBattlecruiser() int64 { return b.GeneralOverhaulBattlecruiser }
func (b LfResearches) GetGeneralOverhaulBomber() int64        { return b.GeneralOverhaulBomber }
func (b LfResearches) GetGeneralOverhaulDestroyer() int64     { return b.GeneralOverhaulDestroyer }
func (b LfResearches) GetExperimentalWeaponsTechnology() int64 {
	return b.ExperimentalWeaponsTechnology
}
func (b LfResearches) GetMechanGeneralEnhancement() int64 { return b.MechanGeneralEnhancement }
func (b LfResearches) GetHeatRecovery() int64             { return b.HeatRecovery }
func (b LfResearches) GetSulphideProcess() int64          { return b.SulphideProcess }
func (b LfResearches) GetPsionicNetwork() int64           { return b.PsionicNetwork }
func (b LfResearches) GetTelekineticTractorBeam() int64   { return b.TelekineticTractorBeam }
func (b LfResearches) GetEnhancedSensorTechnology() int64 { return b.EnhancedSensorTechnology }
func (b LfResearches) GetNeuromodalCompressor() int64     { return b.NeuromodalCompressor }
func (b LfResearches) GetNeuroInterface() int64           { return b.NeuroInterface }
func (b LfResearches) GetInterplanetaryAnalysisNetwork() int64 {
	return b.InterplanetaryAnalysisNetwork
}
func (b LfResearches) GetOverclockingHeavyFighter() int64     { return b.OverclockingHeavyFighter }
func (b LfResearches) GetTelekineticDrive() int64             { return b.TelekineticDrive }
func (b LfResearches) GetSixthSense() int64                   { return b.SixthSense }
func (b LfResearches) GetPsychoharmoniser() int64             { return b.Psychoharmoniser }
func (b LfResearches) GetEfficientSwarmIntelligence() int64   { return b.EfficientSwarmIntelligence }
func (b LfResearches) GetOverclockingLargeCargo() int64       { return b.OverclockingLargeCargo }
func (b LfResearches) GetGravitationSensors() int64           { return b.GravitationSensors }
func (b LfResearches) GetOverclockingBattleship() int64       { return b.OverclockingBattleship }
func (b LfResearches) GetPsionicShieldMatrix() int64          { return b.PsionicShieldMatrix }
func (b LfResearches) GetKaeleshDiscovererEnhancement() int64 { return b.KaeleshDiscovererEnhancement }

func (b LfResearches) Lazy() LazyLfResearches {
	return func() LfResearches { return b }
}

// ByID gets the research level by lfResearch id
func (b LfResearches) ByID(id ID) int64 {
	switch id {
	case IntergalacticEnvoysID:
		return b.IntergalacticEnvoys
	case HighPerformanceExtractorsID:
		return b.HighPerformanceExtractors
	case FusionDrivesID:
		return b.FusionDrives
	case StealthFieldGeneratorID:
		return b.StealthFieldGenerator
	case OrbitalDenID:
		return b.OrbitalDen
	case ResearchAIID:
		return b.ResearchAI
	case HighPerformanceTerraformerID:
		return b.HighPerformanceTerraformer
	case EnhancedProductionTechnologiesID:
		return b.EnhancedProductionTechnologies
	case LightFighterMkIIID:
		return b.LightFighterMkII
	case CruiserMkIIID:
		return b.CruiserMkII
	case ImprovedLabTechnologyID:
		return b.ImprovedLabTechnology
	case PlasmaTerraformerID:
		return b.PlasmaTerraformer
	case LowTemperatureDrivesID:
		return b.LowTemperatureDrives
	case BomberMkIIID:
		return b.BomberMkII
	case DestroyerMkIIID:
		return b.DestroyerMkII
	case BattlecruiserMkIIID:
		return b.BattlecruiserMkII
	case RobotAssistantsID:
		return b.RobotAssistants
	case SupercomputerID:
		return b.Supercomputer
	case VolcanicBatteriesID:
		return b.VolcanicBatteries
	case AcousticScanningID:
		return b.AcousticScanning
	case HighEnergyPumpSystemsID:
		return b.HighEnergyPumpSystems
	case CargoHoldExpansionCivilianShipsID:
		return b.CargoHoldExpansionCivilianShips
	case MagmaPoweredProductionID:
		return b.MagmaPoweredProduction
	case GeothermalPowerPlantsID:
		return b.GeothermalPowerPlants
	case DepthSoundingID:
		return b.DepthSounding
	case IonCrystalEnhancementHeavyFighterID:
		return b.IonCrystalEnhancementHeavyFighter
	case ImprovedStellaratorID:
		return b.ImprovedStellarator
	case HardenedDiamondDrillHeadsID:
		return b.HardenedDiamondDrillHeads
	case SeismicMiningTechnologyID:
		return b.SeismicMiningTechnology
	case MagmaPoweredPumpSystemsID:
		return b.MagmaPoweredPumpSystems
	case IonCrystalModulesID:
		return b.IonCrystalModules
	case OptimisedSiloConstructionMethodID:
		return b.OptimisedSiloConstructionMethod
	case DiamondEnergyTransmitterID:
		return b.DiamondEnergyTransmitter
	case ObsidianShieldReinforcementID:
		return b.ObsidianShieldReinforcement
	case RuneShieldsID:
		return b.RuneShields
	case RocktalCollectorEnhancementID:
		return b.RocktalCollectorEnhancement
	case CatalyserTechnologyID:
		return b.CatalyserTechnology
	case PlasmaDriveID:
		return b.PlasmaDrive
	case EfficiencyModuleID:
		return b.EfficiencyModule
	case DepotAIID:
		return b.DepotAI
	case GeneralOverhaulLightFighterID:
		return b.GeneralOverhaulLightFighter
	case AutomatedTransportLinesID:
		return b.AutomatedTransportLines
	case ImprovedDroneAIID:
		return b.ImprovedDroneAI
	case ExperimentalRecyclingTechnologyID:
		return b.ExperimentalRecyclingTechnology
	case GeneralOverhaulCruiserID:
		return b.GeneralOverhaulCruiser
	case SlingshotAutopilotID:
		return b.SlingshotAutopilot
	case HighTemperatureSuperconductorsID:
		return b.HighTemperatureSuperconductors
	case GeneralOverhaulBattleshipID:
		return b.GeneralOverhaulBattleship
	case ArtificialSwarmIntelligenceID:
		return b.ArtificialSwarmIntelligence
	case GeneralOverhaulBattlecruiserID:
		return b.GeneralOverhaulBattlecruiser
	case GeneralOverhaulBomberID:
		return b.GeneralOverhaulBomber
	case GeneralOverhaulDestroyerID:
		return b.GeneralOverhaulDestroyer
	case ExperimentalWeaponsTechnologyID:
		return b.ExperimentalWeaponsTechnology
	case MechanGeneralEnhancementID:
		return b.MechanGeneralEnhancement
	case HeatRecoveryID:
		return b.HeatRecovery
	case SulphideProcessID:
		return b.SulphideProcess
	case PsionicNetworkID:
		return b.PsionicNetwork
	case TelekineticTractorBeamID:
		return b.TelekineticTractorBeam
	case EnhancedSensorTechnologyID:
		return b.EnhancedSensorTechnology
	case NeuromodalCompressorID:
		return b.NeuromodalCompressor
	case NeuroInterfaceID:
		return b.NeuroInterface
	case InterplanetaryAnalysisNetworkID:
		return b.InterplanetaryAnalysisNetwork
	case OverclockingHeavyFighterID:
		return b.OverclockingHeavyFighter
	case TelekineticDriveID:
		return b.TelekineticDrive
	case SixthSenseID:
		return b.SixthSense
	case PsychoharmoniserID:
		return b.Psychoharmoniser
	case EfficientSwarmIntelligenceID:
		return b.EfficientSwarmIntelligence
	case OverclockingLargeCargoID:
		return b.OverclockingLargeCargo
	case GravitationSensorsID:
		return b.GravitationSensors
	case OverclockingBattleshipID:
		return b.OverclockingBattleship
	case PsionicShieldMatrixID:
		return b.PsionicShieldMatrix
	case KaeleshDiscovererEnhancementID:
		return b.KaeleshDiscovererEnhancement
	}
	return 0
}

// BaseLfResearch base struct for Lifeform techs
type BaseLfResearch struct {
	BaseTechnology
	durationBase   float64
	durationFactor float64
}

// TechnologyConstructionTime returns the duration it takes to build given technology
func (b BaseLfResearch) TechnologyConstructionTime(level, universeSpeed int64, acc TechAccelerators, hasTechnocrat, isDiscoverer bool) time.Duration {
	levelF := float64(level)
	secs := levelF * b.durationBase * math.Pow(b.durationFactor, levelF)
	secs /= float64(universeSpeed)
	secs = math.Max(1, secs)
	return time.Duration(int64(math.Floor(secs))) * time.Second
}

// ConstructionTime same as TechnologyConstructionTime, needed for BaseOgameObj implementation
func (b BaseLfResearch) ConstructionTime(level, universeSpeed int64, facilities BuildAccelerators, hasTechnocrat, isDiscoverer bool) time.Duration {
	return b.TechnologyConstructionTime(level, universeSpeed, facilities, hasTechnocrat, isDiscoverer)
}

// GetPrice returns the price to build the given level
func (b BaseLfResearch) GetPrice(level int64) Resources {
	tmp := func(baseCost int64, increaseFactor float64, level int64) int64 {
		return int64(float64(baseCost) * math.Pow(increaseFactor, float64(level-1)) * float64(level))
	}
	return Resources{
		Metal:     tmp(b.BaseCost.Metal, b.IncreaseFactor, level),
		Crystal:   tmp(b.BaseCost.Crystal, b.IncreaseFactor, level),
		Deuterium: tmp(b.BaseCost.Deuterium, b.IncreaseFactor, level),
	}
}

// Humans
type intergalacticEnvoys struct {
	BaseLfResearch
}

func newIntergalacticEnvoys() *intergalacticEnvoys {
	b := new(intergalacticEnvoys)
	b.Name = "IntergalacticEnvoys"
	b.ID = IntergalacticEnvoysID
	b.durationBase = 1000
	b.durationFactor = 1.2
	b.IncreaseFactor = 1.30
	b.BaseCost = Resources{Metal: 5000, Crystal: 2500, Deuterium: 500}
	b.Requirements = map[ID]int64{}
	return b
}

type highPerformanceExtractors struct {
	BaseLfResearch
}

func newHighPerformanceExtractors() *highPerformanceExtractors {
	b := new(highPerformanceExtractors)
	b.Name = "HighPerformanceExtractors"
	b.ID = HighPerformanceExtractorsID
	b.durationBase = 2000
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 7000, Crystal: 10000, Deuterium: 5000}
	b.Requirements = map[ID]int64{}
	return b
}

type fusionDrives struct {
	BaseLfResearch
}

func newFusionDrives() *fusionDrives {
	b := new(fusionDrives)
	b.Name = "FusionDrives"
	b.ID = FusionDrivesID
	b.durationBase = 2500
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.30
	b.BaseCost = Resources{Metal: 15000, Crystal: 10000, Deuterium: 5000}
	b.Requirements = map[ID]int64{}
	return b
}

type stealthFieldGenerator struct {
	BaseLfResearch
}

func newStealthFieldGenerator() *stealthFieldGenerator {
	b := new(stealthFieldGenerator)
	b.Name = "StealthFieldGenerator"
	b.ID = StealthFieldGeneratorID
	b.durationBase = 3500
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.30
	b.BaseCost = Resources{Metal: 20000, Crystal: 15000, Deuterium: 7500}
	b.Requirements = map[ID]int64{}
	return b
}

type orbitalDen struct {
	BaseLfResearch
}

func newOrbitalDen() *orbitalDen {
	b := new(orbitalDen)
	b.Name = "OrbitalDen"
	b.ID = OrbitalDenID
	b.durationBase = 4500
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.40
	b.BaseCost = Resources{Metal: 25000, Crystal: 20000, Deuterium: 10000}
	b.Requirements = map[ID]int64{}
	return b
}

type researchAI struct {
	BaseLfResearch
}

func newResearchAI() *researchAI {
	b := new(researchAI)
	b.Name = "ResearchAI"
	b.ID = ResearchAIID
	b.durationBase = 5000
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 35000, Crystal: 25000, Deuterium: 15000}
	b.Requirements = map[ID]int64{}
	return b
}

type highPerformanceTerraformer struct {
	BaseLfResearch
}

func newHighPerformanceTerraformer() *highPerformanceTerraformer {
	b := new(highPerformanceTerraformer)
	b.Name = "HighPerformanceTerraformer"
	b.ID = HighPerformanceTerraformerID
	b.durationBase = 8000
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.30
	b.BaseCost = Resources{Metal: 70000, Crystal: 40000, Deuterium: 20000}
	b.Requirements = map[ID]int64{}
	return b
}

type enhancedProductionTechnologies struct {
	BaseLfResearch
}

func newEnhancedProductionTechnologies() *enhancedProductionTechnologies {
	b := new(enhancedProductionTechnologies)
	b.Name = "EnhancedProductionTechnologies"
	b.ID = EnhancedProductionTechnologiesID
	b.durationBase = 6000
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 80000, Crystal: 50000, Deuterium: 20000}
	b.Requirements = map[ID]int64{}
	return b
}

type lightFighterMkII struct {
	BaseLfResearch
}

func newLightFighterMkII() *lightFighterMkII {
	b := new(lightFighterMkII)
	b.Name = "LightFighterMkII"
	b.ID = LightFighterMkIIID
	b.durationBase = 6500
	b.durationFactor = 1.4
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 320000, Crystal: 240000, Deuterium: 100000}
	b.Requirements = map[ID]int64{}
	return b
}

type cruiserMkII struct {
	BaseLfResearch
}

func newCruiserMkII() *cruiserMkII {
	b := new(cruiserMkII)
	b.Name = "CruiserMkII"
	b.ID = CruiserMkIIID
	b.durationBase = 7000
	b.durationFactor = 1.4
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 320000, Crystal: 240000, Deuterium: 100000}
	b.Requirements = map[ID]int64{}
	return b
}

type improvedLabTechnology struct {
	BaseLfResearch
}

func newImprovedLabTechnology() *improvedLabTechnology {
	b := new(improvedLabTechnology)
	b.Name = "ImprovedLabTechnology"
	b.ID = ImprovedLabTechnologyID
	b.durationBase = 7500
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 120000, Crystal: 30000, Deuterium: 25000}
	b.Requirements = map[ID]int64{}
	return b
}

type plasmaTerraformer struct {
	BaseLfResearch
}

func newPlasmaTerraformer() *plasmaTerraformer {
	b := new(plasmaTerraformer)
	b.Name = "PlasmaTerraformer"
	b.ID = PlasmaTerraformerID
	b.durationBase = 10000
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.30
	b.BaseCost = Resources{Metal: 100000, Crystal: 40000, Deuterium: 30000}
	b.Requirements = map[ID]int64{}
	return b
}

type lowTemperatureDrives struct {
	BaseLfResearch
}

func newLowTemperatureDrives() *lowTemperatureDrives {
	b := new(lowTemperatureDrives)
	b.Name = "LowTemperatureDrives"
	b.ID = LowTemperatureDrivesID
	b.durationBase = 8500
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.30
	b.BaseCost = Resources{Metal: 200000, Crystal: 100000, Deuterium: 100000}
	b.Requirements = map[ID]int64{}
	return b
}

type bomberMkII struct {
	BaseLfResearch
}

func newBomberMkII() *bomberMkII {
	b := new(bomberMkII)
	b.Name = "BomberMkII"
	b.ID = BomberMkIIID
	b.durationBase = 9000
	b.durationFactor = 1.4
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 160000, Crystal: 120000, Deuterium: 50000}
	b.Requirements = map[ID]int64{}
	return b
}

type destroyerMkII struct {
	BaseLfResearch
}

func newDestroyerMkII() *destroyerMkII {
	b := new(destroyerMkII)
	b.Name = "DestroyerMkII"
	b.ID = DestroyerMkIIID
	b.durationBase = 9500
	b.durationFactor = 1.4
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 160000, Crystal: 120000, Deuterium: 50000}
	b.Requirements = map[ID]int64{}
	return b
}

type battlecruiserMkII struct {
	BaseLfResearch
}

func newBattlecruiserMkII() *battlecruiserMkII {
	b := new(battlecruiserMkII)
	b.Name = "BattlecruiserMkII"
	b.ID = BattlecruiserMkIIID
	b.durationBase = 10000
	b.durationFactor = 1.4
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 320000, Crystal: 240000, Deuterium: 100000}
	b.Requirements = map[ID]int64{}
	return b
}

type robotAssistants struct {
	BaseLfResearch
}

func newRobotAssistants() *robotAssistants {
	b := new(robotAssistants)
	b.Name = "robotAssistants"
	b.ID = RobotAssistantsID
	b.durationBase = 11000
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 300000, Crystal: 180000, Deuterium: 120000}
	b.Requirements = map[ID]int64{}
	return b
}

type supercomputer struct {
	BaseLfResearch
}

func newSupercomputer() *supercomputer {
	b := new(supercomputer)
	b.Name = "Supercomputer"
	b.ID = SupercomputerID
	b.durationBase = 13000
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.30
	b.BaseCost = Resources{Metal: 500000, Crystal: 300000, Deuterium: 200000}
	b.Requirements = map[ID]int64{}
	return b
}

// Rocktal
type volcanicBatteries struct {
	BaseLfResearch
}

func newVolcanicBatteries() *volcanicBatteries {
	b := new(volcanicBatteries)
	b.Name = "VolcanicBatteries"
	b.ID = VolcanicBatteriesID
	b.durationBase = 1000
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 10000, Crystal: 6000, Deuterium: 1000}
	b.Requirements = map[ID]int64{}
	return b
}

type acousticScanning struct {
	BaseLfResearch
}

func newAcousticScanning() *acousticScanning {
	b := new(acousticScanning)
	b.Name = "AcousticScanning"
	b.ID = AcousticScanningID
	b.durationBase = 2000
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 7500, Crystal: 12500, Deuterium: 5000}
	b.Requirements = map[ID]int64{}
	return b
}

type highEnergyPumpSystems struct {
	BaseLfResearch
}

func newHighEnergyPumpSystems() *highEnergyPumpSystems {
	b := new(highEnergyPumpSystems)
	b.Name = "HighEnergyPumpSystems"
	b.ID = HighEnergyPumpSystemsID
	b.durationBase = 2500
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 15000, Crystal: 10000, Deuterium: 5000}
	b.Requirements = map[ID]int64{}
	return b
}

type cargoHoldExpansionCivilianShips struct {
	BaseLfResearch
}

func newCargoHoldExpansionCivilianShips() *cargoHoldExpansionCivilianShips {
	b := new(cargoHoldExpansionCivilianShips)
	b.Name = "CargoHoldExpansionCivilianShips"
	b.ID = CargoHoldExpansionCivilianShipsID
	b.durationBase = 3500
	b.durationFactor = 1.4
	b.IncreaseFactor = 1.30
	b.BaseCost = Resources{Metal: 20000, Crystal: 15000, Deuterium: 7500}
	b.Requirements = map[ID]int64{}
	return b
}

type magmaPoweredProduction struct {
	BaseLfResearch
}

func newMagmaPoweredProduction() *magmaPoweredProduction {
	b := new(magmaPoweredProduction)
	b.Name = "MagmaPoweredProduction"
	b.ID = MagmaPoweredProductionID
	b.durationBase = 4500
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 25000, Crystal: 20000, Deuterium: 10000}
	b.Requirements = map[ID]int64{}
	return b
}

type geothermalPowerPlants struct {
	BaseLfResearch
}

func newGeothermalPowerPlants() *geothermalPowerPlants {
	b := new(geothermalPowerPlants)
	b.Name = "GeothermalPowerPlants"
	b.ID = GeothermalPowerPlantsID
	b.durationBase = 5000
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 50000, Crystal: 50000, Deuterium: 20000}
	b.Requirements = map[ID]int64{}
	return b
}

type depthSounding struct {
	BaseLfResearch
}

func newDepthSounding() *depthSounding {
	b := new(depthSounding)
	b.Name = "DepthSounding"
	b.ID = DepthSoundingID
	b.durationBase = 5000
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 70000, Crystal: 40000, Deuterium: 20000}
	b.Requirements = map[ID]int64{}
	return b
}

type ionCrystalEnhancementHeavyFighter struct {
	BaseLfResearch
}

func newIonCrystalEnhancementHeavyFighter() *ionCrystalEnhancementHeavyFighter {
	b := new(ionCrystalEnhancementHeavyFighter)
	b.Name = "IonCrystalEnhancementHeavyFighter"
	b.ID = IonCrystalEnhancementHeavyFighterID
	b.durationBase = 6000
	b.durationFactor = 1.4
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 160000, Crystal: 120000, Deuterium: 50000}
	b.Requirements = map[ID]int64{}
	return b
}

type improvedStellarator struct {
	BaseLfResearch
}

func newImprovedStellarator() *improvedStellarator {
	b := new(improvedStellarator)
	b.Name = "ImprovedStellarator"
	b.ID = ImprovedStellaratorID
	b.durationBase = 6500
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 75000, Crystal: 55000, Deuterium: 25000}
	b.Requirements = map[ID]int64{}
	return b
}

type hardenedDiamondDrillHeads struct {
	BaseLfResearch
}

func newHardenedDiamondDrillHeads() *hardenedDiamondDrillHeads {
	b := new(hardenedDiamondDrillHeads)
	b.Name = "HardenedDiamondDrillHeads"
	b.ID = HardenedDiamondDrillHeadsID
	b.durationBase = 7000
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 85000, Crystal: 40000, Deuterium: 35000}
	b.Requirements = map[ID]int64{}
	return b
}

type seismicMiningTechnology struct {
	BaseLfResearch
}

func newSeismicMiningTechnology() *seismicMiningTechnology {
	b := new(seismicMiningTechnology)
	b.Name = "SeismicMiningTechnology"
	b.ID = SeismicMiningTechnologyID
	b.durationBase = 7500
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 120000, Crystal: 30000, Deuterium: 25000}
	b.Requirements = map[ID]int64{}
	return b
}

type magmaPoweredPumpSystems struct {
	BaseLfResearch
}

func newMagmaPoweredPumpSystems() *magmaPoweredPumpSystems {
	b := new(magmaPoweredPumpSystems)
	b.Name = "MagmaPoweredPumpSystems"
	b.ID = MagmaPoweredPumpSystemsID
	b.durationBase = 8000
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 100000, Crystal: 40000, Deuterium: 30000}
	b.Requirements = map[ID]int64{}
	return b
}

type ionCrystalModules struct {
	BaseLfResearch
}

func newIonCrystalModules() *ionCrystalModules {
	b := new(ionCrystalModules)
	b.Name = "IonCrystalModules"
	b.ID = IonCrystalModulesID
	b.durationBase = 8500
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.20
	b.BaseCost = Resources{Metal: 200000, Crystal: 100000, Deuterium: 100000}
	b.Requirements = map[ID]int64{}
	return b
}

type optimisedSiloConstructionMethod struct {
	BaseLfResearch
}

func newOptimisedSiloConstructionMethod() *optimisedSiloConstructionMethod {
	b := new(optimisedSiloConstructionMethod)
	b.Name = "OptimisedSiloConstructionMethod"
	b.ID = OptimisedSiloConstructionMethodID
	b.durationBase = 9000
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.30
	b.BaseCost = Resources{Metal: 220000, Crystal: 110000, Deuterium: 110000}
	b.Requirements = map[ID]int64{}
	return b
}

type diamondEnergyTransmitter struct {
	BaseLfResearch
}

func newDiamondEnergyTransmitter() *diamondEnergyTransmitter {
	b := new(diamondEnergyTransmitter)
	b.Name = "DiamondEnergyTransmitter"
	b.ID = DiamondEnergyTransmitterID
	b.durationBase = 9500
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.30
	b.BaseCost = Resources{Metal: 240000, Crystal: 120000, Deuterium: 120000}
	b.Requirements = map[ID]int64{}
	return b
}

type obsidianShieldReinforcement struct {
	BaseLfResearch
}

func newObsidianShieldReinforcement() *obsidianShieldReinforcement {
	b := new(obsidianShieldReinforcement)
	b.Name = "ObsidianShieldReinforcement"
	b.ID = ObsidianShieldReinforcementID
	b.durationBase = 10000
	b.durationFactor = 1.4
	b.IncreaseFactor = 1.40
	b.BaseCost = Resources{Metal: 250000, Crystal: 250000, Deuterium: 250000}
	b.Requirements = map[ID]int64{}
	return b
}

type runeShields struct {
	BaseLfResearch
}

func newRuneShields() *runeShields {
	b := new(runeShields)
	b.Name = "RuneShields"
	b.ID = RuneShieldsID
	b.durationBase = 13000
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 500000, Crystal: 300000, Deuterium: 200000}
	b.Requirements = map[ID]int64{}
	return b
}

type rocktalCollectorEnhancement struct {
	BaseLfResearch
}

func newRocktalCollectorEnhancement() *rocktalCollectorEnhancement {
	b := new(rocktalCollectorEnhancement)
	b.Name = "RocktalCollectorEnhancement"
	b.ID = RocktalCollectorEnhancementID
	b.durationBase = 11000
	b.durationFactor = 1.4
	b.IncreaseFactor = 1.70
	b.BaseCost = Resources{Metal: 300000, Crystal: 180000, Deuterium: 120000}
	b.Requirements = map[ID]int64{}
	return b
}

//Mechas

type catalyserTechnology struct {
	BaseLfResearch
}

func newCatalyserTechnology() *catalyserTechnology {
	b := new(catalyserTechnology)
	b.Name = "CatalyserTechnology"
	b.ID = CatalyserTechnologyID
	b.durationBase = 1000
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 10000, Crystal: 6000, Deuterium: 1000}
	b.Requirements = map[ID]int64{}
	return b
}

type plasmaDrive struct {
	BaseLfResearch
}

func newPlasmaDrive() *plasmaDrive {
	b := new(plasmaDrive)
	b.Name = "PlasmaDrive"
	b.ID = PlasmaDriveID
	b.durationBase = 2000
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.30
	b.BaseCost = Resources{Metal: 7500, Crystal: 12500, Deuterium: 5000}
	b.Requirements = map[ID]int64{}
	return b
}

type efficiencyModule struct {
	BaseLfResearch
}

func newEfficiencyModule() *efficiencyModule {
	b := new(efficiencyModule)
	b.Name = "EfficiencyModule"
	b.ID = EfficiencyModuleID
	b.durationBase = 2500
	b.durationFactor = 1.4
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 15000, Crystal: 10000, Deuterium: 5000}
	b.Requirements = map[ID]int64{}
	return b
}

type depotAI struct {
	BaseLfResearch
}

func newDepotAI() *depotAI {
	b := new(depotAI)
	b.Name = "DepotAI"
	b.ID = DepotAIID
	b.durationBase = 3500
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.30
	b.BaseCost = Resources{Metal: 20000, Crystal: 15000, Deuterium: 7500}
	b.Requirements = map[ID]int64{}
	return b
}

type generalOverhaulLightFighter struct {
	BaseLfResearch
}

func newGeneralOverhaulLightFighter() *generalOverhaulLightFighter {
	b := new(generalOverhaulLightFighter)
	b.Name = "GeneralOverhaulLightFighter"
	b.ID = GeneralOverhaulLightFighterID
	b.durationBase = 4500
	b.durationFactor = 1.4
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 160000, Crystal: 120000, Deuterium: 50000}
	b.Requirements = map[ID]int64{}
	return b
}

type automatedTransportLines struct {
	BaseLfResearch
}

func newAutomatedTransportLines() *automatedTransportLines {
	b := new(automatedTransportLines)
	b.Name = "AutomatedTransportLines"
	b.ID = AutomatedTransportLinesID
	b.durationBase = 5000
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 50000, Crystal: 50000, Deuterium: 20000}
	b.Requirements = map[ID]int64{}
	return b
}

type improvedDroneAI struct {
	BaseLfResearch
}

func newImprovedDroneAI() *improvedDroneAI {
	b := new(improvedDroneAI)
	b.Name = "ImprovedDroneAI"
	b.ID = ImprovedDroneAIID
	b.durationBase = 5500
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.30
	b.BaseCost = Resources{Metal: 70000, Crystal: 40000, Deuterium: 20000}
	b.Requirements = map[ID]int64{}
	return b
}

type experimentalRecyclingTechnology struct {
	BaseLfResearch
}

func newExperimentalRecyclingTechnology() *experimentalRecyclingTechnology {
	b := new(experimentalRecyclingTechnology)
	b.Name = "ExperimentalRecyclingTechnology"
	b.ID = ExperimentalRecyclingTechnologyID
	b.durationBase = 6000
	b.durationFactor = 1.4
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 160000, Crystal: 120000, Deuterium: 50000}
	b.Requirements = map[ID]int64{}
	return b
}

type generalOverhaulCruiser struct {
	BaseLfResearch
}

func newGeneralOverhaulCruiser() *generalOverhaulCruiser {
	b := new(generalOverhaulCruiser)
	b.Name = "GeneralOverhaulCruiser"
	b.ID = GeneralOverhaulCruiserID
	b.durationBase = 6500
	b.durationFactor = 1.4
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 160000, Crystal: 120000, Deuterium: 50000}
	b.Requirements = map[ID]int64{}
	return b
}

type slingshotAutopilot struct {
	BaseLfResearch
}

func newSlingshotAutopilot() *slingshotAutopilot {
	b := new(slingshotAutopilot)
	b.Name = "SlingshotAutopilot"
	b.ID = SlingshotAutopilotID
	b.durationBase = 7000
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.20
	b.BaseCost = Resources{Metal: 85000, Crystal: 40000, Deuterium: 35000}
	b.Requirements = map[ID]int64{}
	return b
}

type highTemperatureSuperconductors struct {
	BaseLfResearch
}

func newHighTemperatureSuperconductors() *highTemperatureSuperconductors {
	b := new(highTemperatureSuperconductors)
	b.Name = "HighTemperatureSuperconductors"
	b.ID = HighTemperatureSuperconductorsID
	b.durationBase = 7500
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.30
	b.BaseCost = Resources{Metal: 120000, Crystal: 30000, Deuterium: 25000}
	b.Requirements = map[ID]int64{}
	return b
}

type generalOverhaulBattleship struct {
	BaseLfResearch
}

func newGeneralOverhaulBattleship() *generalOverhaulBattleship {
	b := new(generalOverhaulBattleship)
	b.Name = "GeneralOverhaulBattleship"
	b.ID = GeneralOverhaulBattleshipID
	b.durationBase = 8000
	b.durationFactor = 1.4
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 160000, Crystal: 120000, Deuterium: 50000}
	b.Requirements = map[ID]int64{}
	return b
}

type artificialSwarmIntelligence struct {
	BaseLfResearch
}

func newArtificialSwarmIntelligence() *artificialSwarmIntelligence {
	b := new(artificialSwarmIntelligence)
	b.Name = "ArtificialSwarmIntelligence"
	b.ID = ArtificialSwarmIntelligenceID
	b.durationBase = 8500
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 200000, Crystal: 100000, Deuterium: 100000}
	b.Requirements = map[ID]int64{}
	return b
}

type generalOverhaulBattlecruiser struct {
	BaseLfResearch
}

func newGeneralOverhaulBattlecruiser() *generalOverhaulBattlecruiser {
	b := new(generalOverhaulBattlecruiser)
	b.Name = "GeneralOverhaulBattlecruiser"
	b.ID = GeneralOverhaulBattlecruiserID
	b.durationBase = 9000
	b.durationFactor = 1.4
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 160000, Crystal: 120000, Deuterium: 50000}
	b.Requirements = map[ID]int64{}
	return b
}

type generalOverhaulBomber struct {
	BaseLfResearch
}

func newGeneralOverhaulBomber() *generalOverhaulBomber {
	b := new(generalOverhaulBomber)
	b.Name = "GeneralOverhaulBomber"
	b.ID = GeneralOverhaulBomberID
	b.durationBase = 9500
	b.durationFactor = 1.4
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 320000, Crystal: 240000, Deuterium: 100000}
	b.Requirements = map[ID]int64{}
	return b
}

type generalOverhaulDestroyer struct {
	BaseLfResearch
}

func newGeneralOverhaulDestroyer() *generalOverhaulDestroyer {
	b := new(generalOverhaulDestroyer)
	b.Name = "GeneralOverhaulDestroyer"
	b.ID = GeneralOverhaulDestroyerID
	b.durationBase = 10000
	b.durationFactor = 1.4
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 320000, Crystal: 240000, Deuterium: 100000}
	b.Requirements = map[ID]int64{}
	return b
}

type experimentalWeaponsTechnology struct {
	BaseLfResearch
}

func newExperimentalWeaponsTechnology() *experimentalWeaponsTechnology {
	b := new(experimentalWeaponsTechnology)
	b.Name = "ExperimentalWeaponsTechnology"
	b.ID = ExperimentalWeaponsTechnologyID
	b.durationBase = 13000
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 500000, Crystal: 300000, Deuterium: 200000}
	b.Requirements = map[ID]int64{}
	return b
}

type mechanGeneralEnhancement struct {
	BaseLfResearch
}

func newMechanGeneralEnhancement() *mechanGeneralEnhancement {
	b := new(mechanGeneralEnhancement)
	b.Name = "MechanGeneralEnhancement"
	b.ID = MechanGeneralEnhancementID
	b.durationBase = 11000
	b.durationFactor = 1.4
	b.IncreaseFactor = 1.70
	b.BaseCost = Resources{Metal: 300000, Crystal: 180000, Deuterium: 120000}
	b.Requirements = map[ID]int64{}
	return b
}

// Kaelesh
type heatRecovery struct {
	BaseLfResearch
}

func newHeatRecovery() *heatRecovery {
	b := new(heatRecovery)
	b.Name = "HeatRecovery"
	b.ID = HeatRecoveryID
	b.durationBase = 1000
	b.durationFactor = 1.4
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 10000, Crystal: 6000, Deuterium: 1000}
	b.Requirements = map[ID]int64{}
	return b
}

type sulphideProcess struct {
	BaseLfResearch
}

func newSulphideProcess() *sulphideProcess {
	b := new(sulphideProcess)
	b.Name = "SulphideProcess"
	b.ID = SulphideProcessID
	b.durationBase = 2000
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 7500, Crystal: 12500, Deuterium: 5000}
	b.Requirements = map[ID]int64{}
	return b
}

type psionicNetwork struct {
	BaseLfResearch
}

func newPsionicNetwork() *psionicNetwork {
	b := new(psionicNetwork)
	b.Name = "PsionicNetwork"
	b.ID = PsionicNetworkID
	b.durationBase = 2500
	b.durationFactor = 1.4
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 15000, Crystal: 10000, Deuterium: 5000}
	b.Requirements = map[ID]int64{}
	return b
}

type telekineticTractorBeam struct {
	BaseLfResearch
}

func newTelekineticTractorBeam() *telekineticTractorBeam {
	b := new(telekineticTractorBeam)
	b.Name = "TelekineticTractorBeam"
	b.ID = TelekineticTractorBeamID
	b.durationBase = 3500
	b.durationFactor = 1.4
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 20000, Crystal: 15000, Deuterium: 7500}
	b.Requirements = map[ID]int64{}
	return b
}

type enhancedSensorTechnology struct {
	BaseLfResearch
}

func newEnhancedSensorTechnology() *enhancedSensorTechnology {
	b := new(enhancedSensorTechnology)
	b.Name = "EnhancedSensorTechnology"
	b.ID = EnhancedSensorTechnologyID
	b.durationBase = 4500
	b.durationFactor = 1.4
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 25000, Crystal: 20000, Deuterium: 10000}
	b.Requirements = map[ID]int64{}
	return b
}

type neuromodalCompressor struct {
	BaseLfResearch
}

func newNeuromodalCompressor() *neuromodalCompressor {
	b := new(neuromodalCompressor)
	b.Name = "NeuromodalCompressor"
	b.ID = NeuromodalCompressorID
	b.durationBase = 5000
	b.durationFactor = 1.4
	b.IncreaseFactor = 1.30
	b.BaseCost = Resources{Metal: 50000, Crystal: 50000, Deuterium: 20000}
	b.Requirements = map[ID]int64{}
	return b
}

type neuroInterface struct {
	BaseLfResearch
}

func newNeuroInterface() *neuroInterface {
	b := new(neuroInterface)
	b.Name = "NeuroInterface"
	b.ID = NeuroInterfaceID
	b.durationBase = 5500
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 70000, Crystal: 40000, Deuterium: 20000}
	b.Requirements = map[ID]int64{}
	return b
}

type interplanetaryAnalysisNetwork struct {
	BaseLfResearch
}

func newInterplanetaryAnalysisNetwork() *interplanetaryAnalysisNetwork {
	b := new(interplanetaryAnalysisNetwork)
	b.Name = "InterplanetaryAnalysisNetwork"
	b.ID = InterplanetaryAnalysisNetworkID
	b.durationBase = 6000
	b.durationFactor = 1.2
	b.IncreaseFactor = 1.20
	b.BaseCost = Resources{Metal: 80000, Crystal: 50000, Deuterium: 20000}
	b.Requirements = map[ID]int64{}
	return b
}

type overclockingHeavyFighter struct {
	BaseLfResearch
}

func newOverclockingHeavyFighter() *overclockingHeavyFighter {
	b := new(overclockingHeavyFighter)
	b.Name = "OverclockingHeavyFighter"
	b.ID = OverclockingHeavyFighterID
	b.durationBase = 6500
	b.durationFactor = 1.4
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 320000, Crystal: 240000, Deuterium: 100000}
	b.Requirements = map[ID]int64{}
	return b
}

type telekineticDrive struct {
	BaseLfResearch
}

func newTelekineticDrive() *telekineticDrive {
	b := new(telekineticDrive)
	b.Name = "TelekineticDrive"
	b.ID = TelekineticDriveID
	b.durationBase = 7000
	b.durationFactor = 1.2
	b.IncreaseFactor = 1.20
	b.BaseCost = Resources{Metal: 85000, Crystal: 40000, Deuterium: 35000}
	b.Requirements = map[ID]int64{}
	return b
}

type sixthSense struct {
	BaseLfResearch
}

func newSixthSense() *sixthSense {
	b := new(sixthSense)
	b.Name = "SixthSense"
	b.ID = SixthSenseID
	b.durationBase = 7500
	b.durationFactor = 1.4
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 120000, Crystal: 30000, Deuterium: 25000}
	b.Requirements = map[ID]int64{}
	return b
}

type psychoharmoniser struct {
	BaseLfResearch
}

func newPsychoharmoniser() *psychoharmoniser {
	b := new(psychoharmoniser)
	b.Name = "Psychoharmoniser"
	b.ID = PsychoharmoniserID
	b.durationBase = 8000
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 100000, Crystal: 40000, Deuterium: 30000}
	b.Requirements = map[ID]int64{}
	return b
}

type efficientSwarmIntelligence struct {
	BaseLfResearch
}

func newEfficientSwarmIntelligence() *efficientSwarmIntelligence {
	b := new(efficientSwarmIntelligence)
	b.Name = "EfficientSwarmIntelligence"
	b.ID = EfficientSwarmIntelligenceID
	b.durationBase = 8500
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 200000, Crystal: 100000, Deuterium: 100000}
	b.Requirements = map[ID]int64{}
	return b
}

type overclockingLargeCargo struct {
	BaseLfResearch
}

func newOverclockingLargeCargo() *overclockingLargeCargo {
	b := new(overclockingLargeCargo)
	b.Name = "OverclockingLargeCargo"
	b.ID = OverclockingLargeCargoID
	b.durationBase = 9000
	b.durationFactor = 1.4
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 160000, Crystal: 120000, Deuterium: 50000}
	b.Requirements = map[ID]int64{}
	return b
}

type gravitationSensors struct {
	BaseLfResearch
}

func newGravitationSensors() *gravitationSensors {
	b := new(gravitationSensors)
	b.Name = "GravitationSensors"
	b.ID = GravitationSensorsID
	b.durationBase = 9500
	b.durationFactor = 1.4
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 240000, Crystal: 120000, Deuterium: 120000}
	b.Requirements = map[ID]int64{}
	return b
}

type overclockingBattleship struct {
	BaseLfResearch
}

func newOverclockingBattleship() *overclockingBattleship {
	b := new(overclockingBattleship)
	b.Name = "OverclockingBattleship"
	b.ID = OverclockingBattleshipID
	b.durationBase = 10000
	b.durationFactor = 1.4
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 320000, Crystal: 240000, Deuterium: 100000}
	b.Requirements = map[ID]int64{}
	return b
}

type psionicShieldMatrix struct {
	BaseLfResearch
}

func newPsionicShieldMatrix() *psionicShieldMatrix {
	b := new(psionicShieldMatrix)
	b.Name = "PsionicShieldMatrix"
	b.ID = PsionicShieldMatrixID
	b.durationBase = 13000
	b.durationFactor = 1.3
	b.IncreaseFactor = 1.50
	b.BaseCost = Resources{Metal: 500000, Crystal: 300000, Deuterium: 200000}
	b.Requirements = map[ID]int64{}
	return b
}

type kaeleshDiscovererEnhancement struct {
	BaseLfResearch
}

func newKaeleshDiscovererEnhancement() *kaeleshDiscovererEnhancement {
	b := new(kaeleshDiscovererEnhancement)
	b.Name = "KaeleshDiscovererEnhancement"
	b.ID = KaeleshDiscovererEnhancementID
	b.durationBase = 11000
	b.durationFactor = 1.4
	b.IncreaseFactor = 1.70
	b.BaseCost = Resources{Metal: 300000, Crystal: 180000, Deuterium: 120000}
	b.Requirements = map[ID]int64{}
	return b
}
