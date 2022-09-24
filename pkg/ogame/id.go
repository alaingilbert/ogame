package ogame

import (
	"strconv"
)

// ID represent an ogame id
type ID int64

// IsSet returns either or not the id is set to a value different than 0
func (o ID) IsSet() bool {
	return o.Int64() != 0
}

// Int64 returns an integer value of the id
func (o ID) Int64() int64 {
	return int64(o)
}

// Int returns an integer value of the id
// Deprecated: backward compatibility
func (o ID) Int() int64 {
	return int64(o)
}

func (o ID) String() string {
	res := ""
	switch o {
	case AllianceDepotID:
		res += "AllianceDepot"
	case RoboticsFactoryID:
		res += "RoboticsFactory"
	case ShipyardID:
		res += "Shipyard"
	case ResearchLabID:
		res += "ResearchLab"
	case MissileSiloID:
		res += "MissileSilo"
	case NaniteFactoryID:
		res += "NaniteFactory"
	case TerraformerID:
		res += "Terraformer"
	case SpaceDockID:
		res += "SpaceDock"
	case LunarBaseID:
		res += "LunarBase"
	case SensorPhalanxID:
		res += "SensorPhalanx"
	case JumpGateID:
		res += "JumpGate"
	case MetalMineID:
		res += "MetalMine"
	case CrystalMineID:
		res += "CrystalMine"
	case DeuteriumSynthesizerID:
		res += "DeuteriumSynthesizer"
	case SolarPlantID:
		res += "SolarPlant"
	case FusionReactorID:
		res += "FusionReactor"
	case MetalStorageID:
		res += "MetalStorage"
	case CrystalStorageID:
		res += "CrystalStorage"
	case DeuteriumTankID:
		res += "DeuteriumTank"
	case ShieldedMetalDenID:
		res += "ShieldedMetalDen"
	case UndergroundCrystalDenID:
		res += "UndergroundCrystalDen"
	case SeabedDeuteriumDenID:
		res += "SeabedDeuteriumDen"
	case RocketLauncherID:
		res += "RocketLauncher"
	case LightLaserID:
		res += "LightLaser"
	case HeavyLaserID:
		res += "HeavyLaser"
	case GaussCannonID:
		res += "GaussCannon"
	case IonCannonID:
		res += "IonCannon"
	case PlasmaTurretID:
		res += "PlasmaTurret"
	case SmallShieldDomeID:
		res += "SmallShieldDome"
	case LargeShieldDomeID:
		res += "LargeShieldDome"
	case AntiBallisticMissilesID:
		res += "AntiBallisticMissiles"
	case InterplanetaryMissilesID:
		res += "InterplanetaryMissiles"
	case SmallCargoID:
		res += "SmallCargo"
	case LargeCargoID:
		res += "LargeCargo"
	case LightFighterID:
		res += "LightFighter"
	case HeavyFighterID:
		res += "HeavyFighter"
	case CruiserID:
		res += "Cruiser"
	case BattleshipID:
		res += "Battleship"
	case ColonyShipID:
		res += "ColonyShip"
	case RecyclerID:
		res += "Recycler"
	case EspionageProbeID:
		res += "EspionageProbe"
	case BomberID:
		res += "Bomber"
	case SolarSatelliteID:
		res += "SolarSatellite"
	case DestroyerID:
		res += "Destroyer"
	case DeathstarID:
		res += "Deathstar"
	case BattlecruiserID:
		res += "Battlecruiser"
	case CrawlerID:
		res += "Crawler"
	case ReaperID:
		res += "Reaper"
	case PathfinderID:
		res += "Pathfinder"
	case EspionageTechnologyID:
		res += "EspionageTechnology"
	case ComputerTechnologyID:
		res += "ComputerTechnology"
	case WeaponsTechnologyID:
		res += "WeaponsTechnology"
	case ShieldingTechnologyID:
		res += "ShieldingTechnology"
	case ArmourTechnologyID:
		res += "ArmourTechnology"
	case EnergyTechnologyID:
		res += "EnergyTechnology"
	case HyperspaceTechnologyID:
		res += "HyperspaceTechnology"
	case CombustionDriveID:
		res += "CombustionDrive"
	case ImpulseDriveID:
		res += "ImpulseDrive"
	case HyperspaceDriveID:
		res += "HyperspaceDrive"
	case LaserTechnologyID:
		res += "LaserTechnology"
	case IonTechnologyID:
		res += "IonTechnology"
	case PlasmaTechnologyID:
		res += "PlasmaTechnology"
	case IntergalacticResearchNetworkID:
		res += "IntergalacticResearchNetwork"
	case AstrophysicsID:
		res += "Astrophysics"
	case GravitonTechnologyID:
		res += "GravitonTechnology"
	default:
		res += "Invalid(" + strconv.FormatInt(int64(o), 10) + ")"
	}
	return res
}

// IsValid returns either or not the id is valid
func (o ID) IsValid() bool {
	return o.IsDefense() || o.IsShip() || o.IsTech() || o.IsBuilding() || o.IsLfBuilding() || o.IsLfTech()
}

// IsFacility returns either or not the id is a facility
func (o ID) IsFacility() bool {
	return o == AllianceDepotID ||
		o == RoboticsFactoryID ||
		o == ShipyardID ||
		o == ResearchLabID ||
		o == MissileSiloID ||
		o == NaniteFactoryID ||
		o == TerraformerID ||
		o == SpaceDockID ||
		o == LunarBaseID ||
		o == SensorPhalanxID ||
		o == JumpGateID
}

// IsResourceBuilding returns either or not the id is a resource building
func (o ID) IsResourceBuilding() bool {
	return o == MetalMineID ||
		o == CrystalMineID ||
		o == DeuteriumSynthesizerID ||
		o == SolarPlantID ||
		o == FusionReactorID ||
		o == MetalStorageID ||
		o == CrystalStorageID ||
		o == DeuteriumTankID ||
		o == ShieldedMetalDenID ||
		o == UndergroundCrystalDenID ||
		o == SeabedDeuteriumDenID
}

func (o ID) IsLfBuilding() bool {
	return o == ResidentialSectorID || // humans
		o == BiosphereFarmID ||
		o == ResearchCentreID ||
		o == AcademyOfSciencesID ||
		o == NeuroCalibrationCentreID ||
		o == HighEnergySmeltingID ||
		o == FoodSiloID ||
		o == FusionPoweredProductionID ||
		o == SkyscraperID ||
		o == BiotechLabID ||
		o == MetropolisID ||
		o == PlanetaryShieldID || // rocktal
		o == MeditationEnclaveID ||
		o == CrystalFarmID ||
		o == RuneTechnologiumID ||
		o == RuneForgeID ||
		o == OriktoriumID ||
		o == MagmaForgeID ||
		o == DisruptionChamberID ||
		o == MegalithID ||
		o == CrystalRefineryID ||
		o == DeuteriumSynthesiserID ||
		o == MineralResearchCentreID ||
		o == MetalRecyclingPlantID ||
		o == AssemblyLineID || // mechas
		o == FusionCellFactoryID ||
		o == RoboticsResearchCentreID ||
		o == UpdateNetworkID ||
		o == QuantumComputerCentreID ||
		o == AutomatisedAssemblyCentreID ||
		o == HighPerformanceTransformerID ||
		o == MicrochipAssemblyLineID ||
		o == ProductionAssemblyHallID ||
		o == HighPerformanceSynthesiserID ||
		o == ChipMassProductionID ||
		o == NanoRepairBotsID ||
		o == SanctuaryID || // kaelesh
		o == AntimatterCondenserID ||
		o == VortexChamberID ||
		o == HallsOfRealisationID ||
		o == ForumOfTranscendenceID ||
		o == AntimatterConvectorID ||
		o == CloningLaboratoryID ||
		o == ChrysalisAcceleratorID ||
		o == BioModifierID ||
		o == PsionicModulatorID ||
		o == ShipManufacturingHallID ||
		o == SupraRefractorID
}

// IsBuilding returns either or not the id is a building (facility, resource building)
func (o ID) IsBuilding() bool {
	return o.IsResourceBuilding() || o.IsLfBuilding() || o.IsFacility()
}

// IsTech returns either or not the id is a technology
func (o ID) IsTech() bool {
	return o == EspionageTechnologyID ||
		o == ComputerTechnologyID ||
		o == WeaponsTechnologyID ||
		o == ShieldingTechnologyID ||
		o == ArmourTechnologyID ||
		o == EnergyTechnologyID ||
		o == HyperspaceTechnologyID ||
		o == CombustionDriveID ||
		o == ImpulseDriveID ||
		o == HyperspaceDriveID ||
		o == LaserTechnologyID ||
		o == IonTechnologyID ||
		o == PlasmaTechnologyID ||
		o == IntergalacticResearchNetworkID ||
		o == AstrophysicsID ||
		o == GravitonTechnologyID
}

// IsLsTech returns either or not the id is a lifeform technology
func (o ID) IsLfTech() bool {
		return o ==  IntergalacticEnvoysID || //Humans
		o ==  HighPerformanceExtractorsID ||
		o ==  FusionDrivesID ||
		o ==  StealthFieldGeneratorID ||
		o ==  OrbitalDenID ||
		o ==  ResearchAIID ||
		o ==  HighPerformanceTerraformerID ||
		o ==  EnhancedProductionTechnologiesID ||
		o ==  LightFighterMkIIID ||
		o ==  CruiserMkIIID ||
		o ==  ImprovedLabTechnologyID ||
		o ==  PlasmaTerraformerID ||
		o ==  LowTemperatureDrivesID ||
		o ==  BomberMkIIID ||
		o ==  DestroyerMkIIID ||
		o ==  BattlecruiserMkIIID ||
		o ==  RobotAssistantsID ||
		o ==  SupercomputerID || 
		o ==  VolcanicBatteriesID || //Rocktal
		o ==  AcousticScanningID ||
		o ==  HighEnergyPumpSystemsID ||
		o ==  CargoHoldExpansionCivilianShipsID ||
		o ==  MagmaPoweredProductionID ||
		o ==  GeothermalPowerPlantsID ||
		o ==  DepthSoundingID ||
		o ==  IonCrystalEnhancementHeavyFighterID ||
		o ==  ImprovedStellaratorID ||
		o ==  HardenedDiamondDrillHeadsID ||
		o ==  SeismicMiningTechnologyID ||
		o ==  MagmaPoweredPumpSystemsID ||
		o ==  IonCrystalModulesID ||
		o ==  OptimisedSiloConstructionMethodID ||
		o ==  DiamondEnergyTransmitterID ||
		o ==  ObsidianShieldReinforcementID ||
		o ==  RuneShieldsID ||
		o ==  RocktalCollectorEnhancementID ||
		o ==  CatalyserTechnologyID || //Mechas
		o ==  PlasmaDriveID ||
		o ==  EfficiencyModuleID ||
		o ==  DepotAIID ||
		o ==  GeneralOverhaulLightFighterID ||
		o ==  AutomatedTransportLinesID ||
		o ==  ImprovedDroneAIID ||
		o ==  ExperimentalRecyclingTechnologyID ||
		o ==  GeneralOverhaulCruiserID ||
		o ==  SlingshotAutopilotID ||
		o ==  HighTemperatureSuperconductorsID ||
		o ==  GeneralOverhaulBattleshipID ||
		o ==  ArtificialSwarmIntelligenceID ||
		o ==  GeneralOverhaulBattlecruiserID ||
		o ==  GeneralOverhaulBomberID ||
		o ==  GeneralOverhaulDestroyerID ||
		o ==  ExperimentalWeaponsTechnologyID ||
		o ==  MechanGeneralEnhancementID ||
		o ==  HeatRecoveryID || //Kaelesh
		o ==  SulphideProcessID ||
		o ==  PsionicNetworkID ||
		o ==  TelekineticTractorBeamID ||
		o ==  EnhancedSensorTechnologyID ||
		o ==  NeuromodalCompressorID ||
		o ==  NeuroInterfaceID ||
		o ==  InterplanetaryAnalysisNetworkID ||
		o ==  OverclockingHeavyFighterID ||
		o ==  TelekineticDriveID ||
		o ==  SixthSenseID ||
		o ==  PsychoharmoniserID ||
		o ==  EfficientSwarmIntelligenceID ||
		o ==  OverclockingLargeCargoID ||
		o ==  GravitationSensorsID ||
		o ==  OverclockingBattleshipID ||
		o ==  PsionicShieldMatrixID ||
		o ==  KaeleshDiscovererEnhancementID
}

// IsDefense returns either or not the id is a defense
func (o ID) IsDefense() bool {
	return o == RocketLauncherID ||
		o == LightLaserID ||
		o == HeavyLaserID ||
		o == GaussCannonID ||
		o == IonCannonID ||
		o == PlasmaTurretID ||
		o == SmallShieldDomeID ||
		o == LargeShieldDomeID ||
		o == AntiBallisticMissilesID ||
		o == InterplanetaryMissilesID
}

// IsShip returns either or not the id is a ship
func (o ID) IsShip() bool {
	return o == SmallCargoID ||
		o == LargeCargoID ||
		o == LightFighterID ||
		o == HeavyFighterID ||
		o == CruiserID ||
		o == BattleshipID ||
		o == ColonyShipID ||
		o == RecyclerID ||
		o == EspionageProbeID ||
		o == BomberID ||
		o == SolarSatelliteID ||
		o == DestroyerID ||
		o == DeathstarID ||
		o == BattlecruiserID ||
		o == CrawlerID ||
		o == ReaperID ||
		o == PathfinderID
}

// IsFlyableShip returns either or not the id is a ship that can fly
func (o ID) IsFlyableShip() bool {
	if o == SolarSatelliteID || o == CrawlerID {
		return false
	}
	return o.IsShip()
}

// IsCombatShip ...
func (o ID) IsCombatShip() bool {
	return o == LightFighterID ||
		o == HeavyFighterID ||
		o == CruiserID ||
		o == BattleshipID ||
		o == BomberID ||
		o == DestroyerID ||
		o == DeathstarID ||
		o == BattlecruiserID ||
		o == ReaperID
}

func (o ID) IsValidIPMTarget() bool {
	return !o.IsSet() || (o.IsDefense() && o != AntiBallisticMissilesID && o != InterplanetaryMissilesID)
}
