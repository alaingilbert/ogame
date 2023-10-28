package ogame

// All ogame objects
var (
	AllianceDepot                     = register[*allianceDepot](newAllianceDepot) // Buildings
	CrystalMine                       = register[*crystalMine](newCrystalMine)
	CrystalStorage                    = register[*crystalStorage](newCrystalStorage)
	DeuteriumSynthesizer              = register[*deuteriumSynthesizer](newDeuteriumSynthesizer)
	DeuteriumTank                     = register[*deuteriumTank](newDeuteriumTank)
	FusionReactor                     = register[*fusionReactor](newFusionReactor)
	MetalMine                         = register[*metalMine](newMetalMine)
	MetalStorage                      = register[*metalStorage](newMetalStorage)
	MissileSilo                       = register[*missileSilo](newMissileSilo)
	NaniteFactory                     = register[*naniteFactory](newNaniteFactory)
	ResearchLab                       = register[*researchLab](newResearchLab)
	RoboticsFactory                   = register[*roboticsFactory](newRoboticsFactory)
	SeabedDeuteriumDen                = register[*seabedDeuteriumDen](newSeabedDeuteriumDen)
	ShieldedMetalDen                  = register[*shieldedMetalDen](newShieldedMetalDen)
	Shipyard                          = register[*shipyard](newShipyard)
	SolarPlant                        = register[*solarPlant](newSolarPlant)
	SpaceDock                         = register[*spaceDock](newSpaceDock)
	LunarBase                         = register[*lunarBase](newLunarBase)
	SensorPhalanx                     = register[*sensorPhalanx](newSensorPhalanx)
	JumpGate                          = register[*jumpGate](newJumpGate)
	Terraformer                       = register[*terraformer](newTerraformer)
	UndergroundCrystalDen             = register[*undergroundCrystalDen](newUndergroundCrystalDen)
	SolarSatellite                    = register[*solarSatellite](newSolarSatellite)
	AntiBallisticMissiles             = register[*antiBallisticMissiles](newAntiBallisticMissiles) // Defense
	GaussCannon                       = register[*gaussCannon](newGaussCannon)
	HeavyLaser                        = register[*heavyLaser](newHeavyLaser)
	InterplanetaryMissiles            = register[*interplanetaryMissiles](newInterplanetaryMissiles)
	IonCannon                         = register[*ionCannon](newIonCannon)
	LargeShieldDome                   = register[*largeShieldDome](newLargeShieldDome)
	LightLaser                        = register[*lightLaser](newLightLaser)
	PlasmaTurret                      = register[*plasmaTurret](newPlasmaTurret)
	RocketLauncher                    = register[*rocketLauncher](newRocketLauncher)
	SmallShieldDome                   = register[*smallShieldDome](newSmallShieldDome)
	Battlecruiser                     = register[*battlecruiser](newBattlecruiser) // Ships
	Battleship                        = register[*battleship](newBattleship)
	Bomber                            = register[*bomber](newBomber)
	ColonyShip                        = register[*colonyShip](newColonyShip)
	Cruiser                           = register[*cruiser](newCruiser)
	Deathstar                         = register[*deathstar](newDeathstar)
	Destroyer                         = register[*destroyer](newDestroyer)
	EspionageProbe                    = register[*espionageProbe](newEspionageProbe)
	HeavyFighter                      = register[*heavyFighter](newHeavyFighter)
	LargeCargo                        = register[*largeCargo](newLargeCargo)
	LightFighter                      = register[*lightFighter](newLightFighter)
	Recycler                          = register[*recycler](newRecycler)
	SmallCargo                        = register[*smallCargo](newSmallCargo)
	Crawler                           = register[*crawler](newCrawler)
	Reaper                            = register[*reaper](newReaper)
	Pathfinder                        = register[*pathfinder](newPathfinder)
	ArmourTechnology                  = register[*armourTechnology](newArmourTechnology) // Technologies
	Astrophysics                      = register[*astrophysics](newAstrophysics)
	CombustionDrive                   = register[*combustionDrive](newCombustionDrive)
	ComputerTechnology                = register[*computerTechnology](newComputerTechnology)
	EnergyTechnology                  = register[*energyTechnology](newEnergyTechnology)
	EspionageTechnology               = register[*espionageTechnology](newEspionageTechnology)
	GravitonTechnology                = register[*gravitonTechnology](newGravitonTechnology)
	HyperspaceDrive                   = register[*hyperspaceDrive](newHyperspaceDrive)
	HyperspaceTechnology              = register[*hyperspaceTechnology](newHyperspaceTechnology)
	ImpulseDrive                      = register[*impulseDrive](newImpulseDrive)
	IntergalacticResearchNetwork      = register[*intergalacticResearchNetwork](newIntergalacticResearchNetwork)
	IonTechnology                     = register[*ionTechnology](newIonTechnology)
	LaserTechnology                   = register[*laserTechnology](newLaserTechnology)
	PlasmaTechnology                  = register[*plasmaTechnology](newPlasmaTechnology)
	ShieldingTechnology               = register[*shieldingTechnology](newShieldingTechnology)
	WeaponsTechnology                 = register[*weaponsTechnology](newWeaponsTechnology)
	ResidentialSector                 = register[*residentialSector](newResidentialSector) // Humans
	BiosphereFarm                     = register[*biosphereFarm](newBiosphereFarm)
	ResearchCentre                    = register[*researchCentre](newResearchCentre)
	AcademyOfSciences                 = register[*academyOfSciences](newAcademyOfSciences)
	NeuroCalibrationCentre            = register[*neuroCalibrationCentre](newNeuroCalibrationCentre)
	HighEnergySmelting                = register[*highEnergySmelting](newHighEnergySmelting)
	FoodSilo                          = register[*foodSilo](newFoodSilo)
	FusionPoweredProduction           = register[*fusionPoweredProduction](newFusionPoweredProduction)
	Skyscraper                        = register[*skyscraper](newSkyscraper)
	BiotechLab                        = register[*biotechLab](newBiotechLab)
	Metropolis                        = register[*metropolis](newMetropolis)
	PlanetaryShield                   = register[*planetaryShield](newPlanetaryShield)
	MeditationEnclave                 = register[*meditationEnclave](newMeditationEnclave) //Rocktal
	CrystalFarm                       = register[*crystalFarm](newCrystalFarm)
	RuneTechnologium                  = register[*runeTechnologium](newRuneTechnologium)
	RuneForge                         = register[*runeForge](newRuneForge)
	Oriktorium                        = register[*oriktorium](newOriktorium)
	MagmaForge                        = register[*magmaForge](newMagmaForge)
	DisruptionChamber                 = register[*disruptionChamber](newDisruptionChamber)
	Megalith                          = register[*megalith](newMegalith)
	CrystalRefinery                   = register[*crystalRefinery](newCrystalRefinery)
	DeuteriumSynthesiser              = register[*deuteriumSynthesiser](newDeuteriumSynthesiser)
	MineralResearchCentre             = register[*mineralResearchCentre](newMineralResearchCentre)
	AdvancedRecyclingPlant            = register[*advancedRecyclingPlant](newAdvancedRecyclingPlant)
	AssemblyLine                      = register[*assemblyLine](newAssemblyLine) //Mechas
	FusionCellFactory                 = register[*fusionCellFactory](newFusionCellFactory)
	RoboticsResearchCentre            = register[*roboticsResearchCentre](newRoboticsResearchCentre)
	UpdateNetwork                     = register[*updateNetwork](newUpdateNetwork)
	QuantumComputerCentre             = register[*quantumComputerCentre](newQuantumComputerCentre)
	AutomatisedAssemblyCentre         = register[*automatisedAssemblyCentre](newAutomatisedAssemblyCentre)
	HighPerformanceTransformer        = register[*highPerformanceTransformer](newHighPerformanceTransformer)
	MicrochipAssemblyLine             = register[*microchipAssemblyLine](newMicrochipAssemblyLine)
	ProductionAssemblyHall            = register[*productionAssemblyHall](newProductionAssemblyHall)
	HighPerformanceSynthesiser        = register[*highPerformanceSynthesiser](newHighPerformanceSynthesiser)
	ChipMassProduction                = register[*chipMassProduction](newChipMassProduction)
	NanoRepairBots                    = register[*nanoRepairBots](newNanoRepairBots)
	Sanctuary                         = register[*sanctuary](newSanctuary) //Kaelesh
	AntimatterCondenser               = register[*antimatterCondenser](newAntimatterCondenser)
	VortexChamber                     = register[*vortexChamber](newVortexChamber)
	HallsOfRealisation                = register[*hallsOfRealisation](newHallsOfRealisation)
	ForumOfTranscendence              = register[*forumOfTranscendence](newForumOfTranscendence)
	AntimatterConvector               = register[*antimatterConvector](newAntimatterConvector)
	CloningLaboratory                 = register[*cloningLaboratory](newCloningLaboratory)
	ChrysalisAccelerator              = register[*chrysalisAccelerator](newChrysalisAccelerator)
	BioModifier                       = register[*bioModifier](newBioModifier)
	PsionicModulator                  = register[*psionicModulator](newPsionicModulator)
	ShipManufacturingHall             = register[*shipManufacturingHall](newShipManufacturingHall)
	SupraRefractor                    = register[*supraRefractor](newSupraRefractor)
	IntergalacticEnvoys               = register[*intergalacticEnvoys](newIntergalacticEnvoys) // Humans tech
	HighPerformanceExtractors         = register[*highPerformanceExtractors](newHighPerformanceExtractors)
	FusionDrives                      = register[*fusionDrives](newFusionDrives)
	StealthFieldGenerator             = register[*stealthFieldGenerator](newStealthFieldGenerator)
	OrbitalDen                        = register[*orbitalDen](newOrbitalDen)
	ResearchAI                        = register[*researchAI](newResearchAI)
	HighPerformanceTerraformer        = register[*highPerformanceTerraformer](newHighPerformanceTerraformer)
	EnhancedProductionTechnologies    = register[*enhancedProductionTechnologies](newEnhancedProductionTechnologies)
	LightFighterMkII                  = register[*lightFighterMkII](newLightFighterMkII)
	CruiserMkII                       = register[*cruiserMkII](newCruiserMkII)
	ImprovedLabTechnology             = register[*improvedLabTechnology](newImprovedLabTechnology)
	PlasmaTerraformer                 = register[*plasmaTerraformer](newPlasmaTerraformer)
	LowTemperatureDrives              = register[*lowTemperatureDrives](newLowTemperatureDrives)
	BomberMkII                        = register[*bomberMkII](newBomberMkII)
	DestroyerMkII                     = register[*destroyerMkII](newDestroyerMkII)
	BattlecruiserMkII                 = register[*battlecruiserMkII](newBattlecruiserMkII)
	RobotAssistants                   = register[*robotAssistants](newRobotAssistants)
	Supercomputer                     = register[*supercomputer](newSupercomputer)
	VolcanicBatteries                 = register[*volcanicBatteries](newVolcanicBatteries) //Rocktal techs
	AcousticScanning                  = register[*acousticScanning](newAcousticScanning)
	HighEnergyPumpSystems             = register[*highEnergyPumpSystems](newHighEnergyPumpSystems)
	CargoHoldExpansionCivilianShips   = register[*cargoHoldExpansionCivilianShips](newCargoHoldExpansionCivilianShips)
	MagmaPoweredProduction            = register[*magmaPoweredProduction](newMagmaPoweredProduction)
	GeothermalPowerPlants             = register[*geothermalPowerPlants](newGeothermalPowerPlants)
	DepthSounding                     = register[*depthSounding](newDepthSounding)
	IonCrystalEnhancementHeavyFighter = register[*ionCrystalEnhancementHeavyFighter](newIonCrystalEnhancementHeavyFighter)
	ImprovedStellarator               = register[*improvedStellarator](newImprovedStellarator)
	HardenedDiamondDrillHeads         = register[*hardenedDiamondDrillHeads](newHardenedDiamondDrillHeads)
	SeismicMiningTechnology           = register[*seismicMiningTechnology](newSeismicMiningTechnology)
	MagmaPoweredPumpSystems           = register[*magmaPoweredPumpSystems](newMagmaPoweredPumpSystems)
	IonCrystalModules                 = register[*ionCrystalModules](newIonCrystalModules)
	OptimisedSiloConstructionMethod   = register[*optimisedSiloConstructionMethod](newOptimisedSiloConstructionMethod)
	DiamondEnergyTransmitter          = register[*diamondEnergyTransmitter](newDiamondEnergyTransmitter)
	ObsidianShieldReinforcement       = register[*obsidianShieldReinforcement](newObsidianShieldReinforcement)
	RuneShields                       = register[*runeShields](newRuneShields)
	RocktalCollectorEnhancement       = register[*rocktalCollectorEnhancement](newRocktalCollectorEnhancement)
	CatalyserTechnology               = register[*catalyserTechnology](newCatalyserTechnology) //Mechas techs
	PlasmaDrive                       = register[*plasmaDrive](newPlasmaDrive)
	EfficiencyModule                  = register[*efficiencyModule](newEfficiencyModule)
	DepotAI                           = register[*depotAI](newDepotAI)
	GeneralOverhaulLightFighter       = register[*generalOverhaulLightFighter](newGeneralOverhaulLightFighter)
	AutomatedTransportLines           = register[*automatedTransportLines](newAutomatedTransportLines)
	ImprovedDroneAI                   = register[*improvedDroneAI](newImprovedDroneAI)
	ExperimentalRecyclingTechnology   = register[*experimentalRecyclingTechnology](newExperimentalRecyclingTechnology)
	GeneralOverhaulCruiser            = register[*generalOverhaulCruiser](newGeneralOverhaulCruiser)
	SlingshotAutopilot                = register[*slingshotAutopilot](newSlingshotAutopilot)
	HighTemperatureSuperconductors    = register[*highTemperatureSuperconductors](newHighTemperatureSuperconductors)
	GeneralOverhaulBattleship         = register[*generalOverhaulBattleship](newGeneralOverhaulBattleship)
	ArtificialSwarmIntelligence       = register[*artificialSwarmIntelligence](newArtificialSwarmIntelligence)
	GeneralOverhaulBattlecruiser      = register[*generalOverhaulBattlecruiser](newGeneralOverhaulBattlecruiser)
	GeneralOverhaulBomber             = register[*generalOverhaulBomber](newGeneralOverhaulBomber)
	GeneralOverhaulDestroyer          = register[*generalOverhaulDestroyer](newGeneralOverhaulDestroyer)
	ExperimentalWeaponsTechnology     = register[*experimentalWeaponsTechnology](newExperimentalWeaponsTechnology)
	MechanGeneralEnhancement          = register[*mechanGeneralEnhancement](newMechanGeneralEnhancement)
	HeatRecovery                      = register[*heatRecovery](newHeatRecovery) //Kaelesh techs
	SulphideProcess                   = register[*sulphideProcess](newSulphideProcess)
	PsionicNetwork                    = register[*psionicNetwork](newPsionicNetwork)
	TelekineticTractorBeam            = register[*telekineticTractorBeam](newTelekineticTractorBeam)
	EnhancedSensorTechnology          = register[*enhancedSensorTechnology](newEnhancedSensorTechnology)
	NeuromodalCompressor              = register[*neuromodalCompressor](newNeuromodalCompressor)
	NeuroInterface                    = register[*neuroInterface](newNeuroInterface)
	InterplanetaryAnalysisNetwork     = register[*interplanetaryAnalysisNetwork](newInterplanetaryAnalysisNetwork)
	OverclockingHeavyFighter          = register[*overclockingHeavyFighter](newOverclockingHeavyFighter)
	TelekineticDrive                  = register[*telekineticDrive](newTelekineticDrive)
	SixthSense                        = register[*sixthSense](newSixthSense)
	Psychoharmoniser                  = register[*psychoharmoniser](newPsychoharmoniser)
	EfficientSwarmIntelligence        = register[*efficientSwarmIntelligence](newEfficientSwarmIntelligence)
	OverclockingLargeCargo            = register[*overclockingLargeCargo](newOverclockingLargeCargo)
	GravitationSensors                = register[*gravitationSensors](newGravitationSensors)
	OverclockingBattleship            = register[*overclockingBattleship](newOverclockingBattleship)
	PsionicShieldMatrix               = register[*psionicShieldMatrix](newPsionicShieldMatrix)
	KaeleshDiscovererEnhancement      = register[*kaeleshDiscovererEnhancement](newKaeleshDiscovererEnhancement)
)

type ObjsStruct struct{ m map[ID]BaseOgameObj }

func (o ObjsStruct) ByID(id ID) BaseOgameObj {
	return o.m[id]
}

var Objs = ObjsStruct{m: make(map[ID]BaseOgameObj)}

func register[T BaseOgameObj](constructorFn func() T) T {
	inst := constructorFn()
	Objs.m[inst.GetID()] = inst
	return inst
}

// Defenses array of all defenses objects
var Defenses = []Defense{
	RocketLauncher,
	LightLaser,
	HeavyLaser,
	GaussCannon,
	IonCannon,
	PlasmaTurret,
	SmallShieldDome,
	LargeShieldDome,
	AntiBallisticMissiles,
	InterplanetaryMissiles,
}

// Ships array of all ships objects
var Ships = []Ship{
	LightFighter,
	HeavyFighter,
	Cruiser,
	Battleship,
	Battlecruiser,
	Bomber,
	Destroyer,
	Deathstar,
	SmallCargo,
	LargeCargo,
	ColonyShip,
	Recycler,
	EspionageProbe,
	SolarSatellite,
	Crawler,
	Reaper,
	Pathfinder,
}

// Buildings array of all buildings/facilities objects
var Buildings = []Building{
	MetalMine,
	CrystalMine,
	DeuteriumSynthesizer,
	SolarPlant,
	FusionReactor,
	SolarSatellite,
	MetalStorage,
	CrystalStorage,
	DeuteriumTank,
	ShieldedMetalDen,
	UndergroundCrystalDen,
	SeabedDeuteriumDen,
	RoboticsFactory,
	Shipyard,
	ResearchLab,
	AllianceDepot,
	MissileSilo,
	NaniteFactory,
	Terraformer,
	SpaceDock,
	LunarBase,
	SensorPhalanx,
	JumpGate,
}

// PlanetBuildings arrays of planet specific buildings
var PlanetBuildings = []Building{
	MetalMine,
	CrystalMine,
	DeuteriumSynthesizer,
	SolarPlant,
	FusionReactor,
	SolarSatellite,
	MetalStorage,
	CrystalStorage,
	DeuteriumTank,
	ShieldedMetalDen,
	UndergroundCrystalDen,
	SeabedDeuteriumDen,
	RoboticsFactory,
	Shipyard,
	ResearchLab,
	AllianceDepot,
	MissileSilo,
	NaniteFactory,
	Terraformer,
	SpaceDock,
}

// MoonBuildings arrays of moon specific buildings
var MoonBuildings = []Building{
	SolarSatellite,
	MetalStorage,
	CrystalStorage,
	DeuteriumTank,
	RoboticsFactory,
	Shipyard,
	LunarBase,
	SensorPhalanx,
	JumpGate,
}

// Technologies array of all technologies objects
var Technologies = []Technology{
	EnergyTechnology,
	LaserTechnology,
	IonTechnology,
	HyperspaceTechnology,
	PlasmaTechnology,
	CombustionDrive,
	ImpulseDrive,
	HyperspaceDrive,
	EspionageTechnology,
	ComputerTechnology,
	Astrophysics,
	IntergalacticResearchNetwork,
	GravitonTechnology,
	WeaponsTechnology,
	ShieldingTechnology,
	ArmourTechnology,
}
