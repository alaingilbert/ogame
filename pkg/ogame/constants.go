package ogame

import (
	"strconv"
)

// MissionID represent a mission id
type MissionID int

func (m MissionID) String() string {
	switch m {
	case Relocate:
		return "Relocate"
	case Attack:
		return "Attack"
	case GroupedAttack:
		return "GroupedAttack"
	case Transport:
		return "Transport"
	case Park:
		return "Park"
	case ParkInThatAlly:
		return "ParkInThatAlly"
	case Spy:
		return "Spy"
	case Colonize:
		return "Colonize"
	case RecycleDebrisField:
		return "RecycleDebrisField"
	case Destroy:
		return "Destroy"
	case MissileAttack:
		return "MissileAttack"
	case Expedition:
		return "Expedition"
	case SearchForLifeforms:
		return "SearchForLifeform"
	default:
		return strconv.FormatInt(int64(m), 10)
	}
}

// CelestialType destination type might be planet/moon/debris
type CelestialType int64

func (d CelestialType) String() string {
	switch d {
	case PlanetType:
		return "planet"
	case MoonType:
		return "moon"
	case DebrisType:
		return "debris"
	default:
		return strconv.FormatInt(int64(d), 10)
	}
}

// Int64 returns an integer value of the CelestialType
func (d CelestialType) Int64() int64 {
	return int64(d)
}

// Int returns an integer value of the CelestialType
// Deprecated: backward compatibility
func (d CelestialType) Int() int64 {
	return int64(d)
}

// AllianceClass ...
type AllianceClass int64

// IsWarrior ...
func (c AllianceClass) IsWarrior() bool {
	return c == Warrior
}

// IsTrader ...
func (c AllianceClass) IsTrader() bool {
	return c == Trader
}

// IsResearcher ...
func (c AllianceClass) IsResearcher() bool {
	return c == Researcher
}

// String ...
func (c AllianceClass) String() string {
	switch c {
	case NoAllianceClass:
		return "NoAllianceClass"
	case Warrior:
		return "Warrior"
	case Trader:
		return "Trader"
	case Researcher:
		return "Researcher"
	default:
		return strconv.FormatInt(int64(c), 10)
	}
}

// DMType ...
type DMType string

// IsValid ...
func (t DMType) IsValid() bool {
	return t == BuildingsDmType || t == ResearchDmType || t == ShipyardDmType
}

// CharacterClass ...
type CharacterClass int64

func (c CharacterClass) IsCollector() bool {
	return c == Collector
}

func (c CharacterClass) IsGeneral() bool {
	return c == General
}

func (c CharacterClass) IsDiscoverer() bool {
	return c == Discoverer
}

// OGame constants
const (
	BuildingsDmType DMType = "buildings"
	ResearchDmType  DMType = "research"
	ShipyardDmType  DMType = "shipyard"

	NoClass    CharacterClass = 0
	Collector  CharacterClass = 1
	General    CharacterClass = 2
	Discoverer CharacterClass = 3

	NoAllianceClass AllianceClass = 0
	Warrior         AllianceClass = 1
	Trader          AllianceClass = 2
	Researcher      AllianceClass = 3

	PlanetType CelestialType = 1
	DebrisType CelestialType = 2
	MoonType   CelestialType = 3

	//Buildings
	MetalMineID                         ID = 1
	CrystalMineID                       ID = 2
	DeuteriumSynthesizerID              ID = 3
	SolarPlantID                        ID = 4
	FusionReactorID                     ID = 12
	MetalStorageID                      ID = 22
	CrystalStorageID                    ID = 23
	DeuteriumTankID                     ID = 24
	ShieldedMetalDenID                  ID = 25
	UndergroundCrystalDenID             ID = 26
	SeabedDeuteriumDenID                ID = 27
	AllianceDepotID                     ID = 34 // Facilities
	RoboticsFactoryID                   ID = 14
	ShipyardID                          ID = 21
	ResearchLabID                       ID = 31
	MissileSiloID                       ID = 44
	NaniteFactoryID                     ID = 15
	TerraformerID                       ID = 33
	SpaceDockID                         ID = 36
	LunarBaseID                         ID = 41 // Moon facilities
	SensorPhalanxID                     ID = 42
	JumpGateID                          ID = 43
	RocketLauncherID                    ID = 401 // Defense
	LightLaserID                        ID = 402
	HeavyLaserID                        ID = 403
	GaussCannonID                       ID = 404
	IonCannonID                         ID = 405
	PlasmaTurretID                      ID = 406
	SmallShieldDomeID                   ID = 407
	LargeShieldDomeID                   ID = 408
	AntiBallisticMissilesID             ID = 502
	InterplanetaryMissilesID            ID = 503
	SmallCargoID                        ID = 202 // Ships
	LargeCargoID                        ID = 203
	LightFighterID                      ID = 204
	HeavyFighterID                      ID = 205
	CruiserID                           ID = 206
	BattleshipID                        ID = 207
	ColonyShipID                        ID = 208
	RecyclerID                          ID = 209
	EspionageProbeID                    ID = 210
	BomberID                            ID = 211
	SolarSatelliteID                    ID = 212
	DestroyerID                         ID = 213
	DeathstarID                         ID = 214
	BattlecruiserID                     ID = 215
	CrawlerID                           ID = 217
	ReaperID                            ID = 218
	PathfinderID                        ID = 219
	EspionageTechnologyID               ID = 106 // Research
	ComputerTechnologyID                ID = 108
	WeaponsTechnologyID                 ID = 109
	ShieldingTechnologyID               ID = 110
	ArmourTechnologyID                  ID = 111
	EnergyTechnologyID                  ID = 113
	HyperspaceTechnologyID              ID = 114
	CombustionDriveID                   ID = 115
	ImpulseDriveID                      ID = 117
	HyperspaceDriveID                   ID = 118
	LaserTechnologyID                   ID = 120
	IonTechnologyID                     ID = 121
	PlasmaTechnologyID                  ID = 122
	IntergalacticResearchNetworkID      ID = 123
	AstrophysicsID                      ID = 124
	GravitonTechnologyID                ID = 199
	ResidentialSectorID                 ID = 11101 // Lifeform (humans)
	BiosphereFarmID                     ID = 11102
	ResearchCentreID                    ID = 11103
	AcademyOfSciencesID                 ID = 11104
	NeuroCalibrationCentreID            ID = 11105
	HighEnergySmeltingID                ID = 11106
	FoodSiloID                          ID = 11107
	FusionPoweredProductionID           ID = 11108
	SkyscraperID                        ID = 11109
	BiotechLabID                        ID = 11110
	MetropolisID                        ID = 11111
	PlanetaryShieldID                   ID = 11112
	MeditationEnclaveID                 ID = 12101 // Lifeform (rocktal)
	CrystalFarmID                       ID = 12102
	RuneTechnologiumID                  ID = 12103
	RuneForgeID                         ID = 12104
	OriktoriumID                        ID = 12105
	MagmaForgeID                        ID = 12106
	DisruptionChamberID                 ID = 12107
	MegalithID                          ID = 12108
	CrystalRefineryID                   ID = 12109
	DeuteriumSynthesiserID              ID = 12110
	MineralResearchCentreID             ID = 12111
	AdvancedRecyclingPlantID            ID = 12112
	AssemblyLineID                      ID = 13101 // Lifeform (mechas)
	FusionCellFactoryID                 ID = 13102
	RoboticsResearchCentreID            ID = 13103
	UpdateNetworkID                     ID = 13104
	QuantumComputerCentreID             ID = 13105
	AutomatisedAssemblyCentreID         ID = 13106
	HighPerformanceTransformerID        ID = 13107
	MicrochipAssemblyLineID             ID = 13108
	ProductionAssemblyHallID            ID = 13109
	HighPerformanceSynthesiserID        ID = 13110
	ChipMassProductionID                ID = 13111
	NanoRepairBotsID                    ID = 13112
	SanctuaryID                         ID = 14101 // Lifeform (kaelesh)
	AntimatterCondenserID               ID = 14102
	VortexChamberID                     ID = 14103
	HallsOfRealisationID                ID = 14104
	ForumOfTranscendenceID              ID = 14105
	AntimatterConvectorID               ID = 14106
	CloningLaboratoryID                 ID = 14107
	ChrysalisAcceleratorID              ID = 14108
	BioModifierID                       ID = 14109
	PsionicModulatorID                  ID = 14110
	ShipManufacturingHallID             ID = 14111
	SupraRefractorID                    ID = 14112
	IntergalacticEnvoysID               ID = 11201 // Human techs
	HighPerformanceExtractorsID         ID = 11202
	FusionDrivesID                      ID = 11203
	StealthFieldGeneratorID             ID = 11204
	OrbitalDenID                        ID = 11205
	ResearchAIID                        ID = 11206
	HighPerformanceTerraformerID        ID = 11207
	EnhancedProductionTechnologiesID    ID = 11208
	LightFighterMkIIID                  ID = 11209
	CruiserMkIIID                       ID = 11210
	ImprovedLabTechnologyID             ID = 11211
	PlasmaTerraformerID                 ID = 11212
	LowTemperatureDrivesID              ID = 11213
	BomberMkIIID                        ID = 11214
	DestroyerMkIIID                     ID = 11215
	BattlecruiserMkIIID                 ID = 11216
	RobotAssistantsID                   ID = 11217
	SupercomputerID                     ID = 11218
	VolcanicBatteriesID                 ID = 12201 //Rocktal techs
	AcousticScanningID                  ID = 12202
	HighEnergyPumpSystemsID             ID = 12203
	CargoHoldExpansionCivilianShipsID   ID = 12204
	MagmaPoweredProductionID            ID = 12205
	GeothermalPowerPlantsID             ID = 12206
	DepthSoundingID                     ID = 12207
	IonCrystalEnhancementHeavyFighterID ID = 12208
	ImprovedStellaratorID               ID = 12209
	HardenedDiamondDrillHeadsID         ID = 12210
	SeismicMiningTechnologyID           ID = 12211
	MagmaPoweredPumpSystemsID           ID = 12212
	IonCrystalModulesID                 ID = 12213
	OptimisedSiloConstructionMethodID   ID = 12214
	DiamondEnergyTransmitterID          ID = 12215
	ObsidianShieldReinforcementID       ID = 12216
	RuneShieldsID                       ID = 12217
	RocktalCollectorEnhancementID       ID = 12218
	CatalyserTechnologyID               ID = 13201 //Mechas techs
	PlasmaDriveID                       ID = 13202
	EfficiencyModuleID                  ID = 13203
	DepotAIID                           ID = 13204
	GeneralOverhaulLightFighterID       ID = 13205
	AutomatedTransportLinesID           ID = 13206
	ImprovedDroneAIID                   ID = 13207
	ExperimentalRecyclingTechnologyID   ID = 13208
	GeneralOverhaulCruiserID            ID = 13209
	SlingshotAutopilotID                ID = 13210
	HighTemperatureSuperconductorsID    ID = 13211
	GeneralOverhaulBattleshipID         ID = 13212
	ArtificialSwarmIntelligenceID       ID = 13213
	GeneralOverhaulBattlecruiserID      ID = 13214
	GeneralOverhaulBomberID             ID = 13215
	GeneralOverhaulDestroyerID          ID = 13216
	ExperimentalWeaponsTechnologyID     ID = 13217
	MechanGeneralEnhancementID          ID = 13218
	HeatRecoveryID                      ID = 14201 //Kaelesh techs
	SulphideProcessID                   ID = 14202
	PsionicNetworkID                    ID = 14203
	TelekineticTractorBeamID            ID = 14204
	EnhancedSensorTechnologyID          ID = 14205
	NeuromodalCompressorID              ID = 14206
	NeuroInterfaceID                    ID = 14207
	InterplanetaryAnalysisNetworkID     ID = 14208
	OverclockingHeavyFighterID          ID = 14209
	TelekineticDriveID                  ID = 14210
	SixthSenseID                        ID = 14211
	PsychoharmoniserID                  ID = 14212
	EfficientSwarmIntelligenceID        ID = 14213
	OverclockingLargeCargoID            ID = 14214
	GravitationSensorsID                ID = 14215
	OverclockingBattleshipID            ID = 14216
	PsionicShieldMatrixID               ID = 14217
	KaeleshDiscovererEnhancementID      ID = 14218

	// Missions
	Relocate           MissionID = 0
	Attack             MissionID = 1
	GroupedAttack      MissionID = 2
	Transport          MissionID = 3
	Park               MissionID = 4
	ParkInThatAlly     MissionID = 5
	Spy                MissionID = 6
	Colonize           MissionID = 7
	RecycleDebrisField MissionID = 8
	Destroy            MissionID = 9
	MissileAttack      MissionID = 10
	Expedition         MissionID = 15
	SearchForLifeforms MissionID = 18

	// Speeds
	TenPercent         Speed = 1
	TwentyPercent      Speed = 2
	ThirtyPercent      Speed = 3
	FourtyPercent      Speed = 4
	FiftyPercent       Speed = 5
	SixtyPercent       Speed = 6
	SeventyPercent     Speed = 7
	EightyPercent      Speed = 8
	NinetyPercent      Speed = 9
	HundredPercent     Speed = 10
	FivePercent        Speed = 0.5 // General class only detailed speeds
	FifteenPercent     Speed = 1.5
	TwentyFivePercent  Speed = 2.5
	ThirtyFivePercent  Speed = 3.5
	FourtyFivePercent  Speed = 4.5
	FiftyFivePercent   Speed = 5.5
	SixtyFivePercent   Speed = 6.5
	SeventyFivePercent Speed = 7.5
	EightyFivePercent  Speed = 8.5
	NinetyFivePercent  Speed = 9.5
)
