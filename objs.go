package ogame

// All ogame objects
var (
	AllianceDepot                = newAllianceDepot() // Buildings
	CrystalMine                  = newCrystalMine()
	CrystalStorage               = newCrystalStorage()
	DeuteriumSynthesizer         = newDeuteriumSynthesizer()
	DeuteriumTank                = newDeuteriumTank()
	FusionReactor                = newFusionReactor()
	MetalMine                    = newMetalMine()
	MetalStorage                 = newMetalStorage()
	MissileSilo                  = newMissileSilo()
	NaniteFactory                = newNaniteFactory()
	ResearchLab                  = newResearchLab()
	RoboticsFactory              = newRoboticsFactory()
	SeabedDeuteriumDen           = newSeabedDeuteriumDen()
	ShieldedMetalDen             = newShieldedMetalDen()
	Shipyard                     = newShipyard()
	SolarPlant                   = newSolarPlant()
	SpaceDock                    = newSpaceDock()
	LunarBase                    = newLunarBase()
	SensorPhalanx                = newSensorPhalanx()
	JumpGate                     = newJumpGate()
	Terraformer                  = newTerraformer()
	UndergroundCrystalDen        = newUndergroundCrystalDen()
	SolarSatellite               = newSolarSatellite()
	AntiBallisticMissiles        = newAntiBallisticMissiles() // Defense
	GaussCannon                  = newGaussCannon()
	HeavyLaser                   = newHeavyLaser()
	InterplanetaryMissiles       = newInterplanetaryMissiles()
	IonCannon                    = newIonCannon()
	LargeShieldDome              = newLargeShieldDome()
	LightLaser                   = newLightLaser()
	PlasmaTurret                 = newPlasmaTurret()
	RocketLauncher               = newRocketLauncher()
	SmallShieldDome              = newSmallShieldDome()
	Battlecruiser                = newBattlecruiser() // Ships
	Battleship                   = newBattleship()
	Bomber                       = newBomber()
	ColonyShip                   = newColonyShip()
	Cruiser                      = newCruiser()
	Deathstar                    = newDeathstar()
	Destroyer                    = newDestroyer()
	EspionageProbe               = newEspionageProbe()
	HeavyFighter                 = newHeavyFighter()
	LargeCargo                   = newLargeCargo()
	LightFighter                 = newLightFighter()
	Recycler                     = newRecycler()
	SmallCargo                   = newSmallCargo()
	ArmourTechnology             = newArmourTechnology() // Technologies
	Astrophysics                 = newAstrophysics()
	CombustionDrive              = newCombustionDrive()
	ComputerTechnology           = newComputerTechnology()
	EnergyTechnology             = newEnergyTechnology()
	EspionageTechnology          = newEspionageTechnology()
	GravitonTechnology           = newGravitonTechnology()
	HyperspaceDrive              = newHyperspaceDrive()
	HyperspaceTechnology         = newHyperspaceTechnology()
	ImpulseDrive                 = newImpulseDrive()
	IntergalacticResearchNetwork = newIntergalacticResearchNetwork()
	IonTechnology                = newIonTechnology()
	LaserTechnology              = newLaserTechnology()
	PlasmaTechnology             = newPlasmaTechnology()
	ShieldingTechnology          = newShieldingTechnology()
	WeaponsTechnology            = newWeaponsTechnology()
)

// ObjsStruct structure containing all possible ogame objects
type ObjsStruct struct {
	AllianceDepot                *allianceDepot
	CrystalMine                  *crystalMine
	CrystalStorage               *crystalStorage
	DeuteriumSynthesizer         *deuteriumSynthesizer
	DeuteriumTank                *deuteriumTank
	FusionReactor                *fusionReactor
	MetalMine                    *metalMine
	MetalStorage                 *metalStorage
	MissileSilo                  *missileSilo
	NaniteFactory                *naniteFactory
	ResearchLab                  *researchLab
	RoboticsFactory              *roboticsFactory
	SeabedDeuteriumDen           *seabedDeuteriumDen
	ShieldedMetalDen             *shieldedMetalDen
	Shipyard                     *shipyard
	SolarPlant                   *solarPlant
	SpaceDock                    *spaceDock
	LunarBase                    *lunarBase
	SensorPhalanx                *sensorPhalanx
	JumpGate                     *jumpGate
	Terraformer                  *terraformer
	UndergroundCrystalDen        *undergroundCrystalDen
	SolarSatellite               *solarSatellite
	AntiBallisticMissiles        *antiBallisticMissiles
	GaussCannon                  *gaussCannon
	HeavyLaser                   *heavyLaser
	InterplanetaryMissiles       *interplanetaryMissiles
	IonCannon                    *ionCannon
	LargeShieldDome              *largeShieldDome
	LightLaser                   *lightLaser
	PlasmaTurret                 *plasmaTurret
	RocketLauncher               *rocketLauncher
	SmallShieldDome              *smallShieldDome
	Battlecruiser                *battlecruiser
	Battleship                   *battleship
	Bomber                       *bomber
	ColonyShip                   *colonyShip
	Cruiser                      *cruiser
	Deathstar                    *deathstar
	Destroyer                    *destroyer
	EspionageProbe               *espionageProbe
	HeavyFighter                 *heavyFighter
	LargeCargo                   *largeCargo
	LightFighter                 *lightFighter
	Recycler                     *recycler
	SmallCargo                   *smallCargo
	ArmourTechnology             *armourTechnology
	Astrophysics                 *astrophysics
	CombustionDrive              *combustionDrive
	ComputerTechnology           *computerTechnology
	EnergyTechnology             *energyTechnology
	EspionageTechnology          *espionageTechnology
	GravitonTechnology           *gravitonTechnology
	HyperspaceDrive              *hyperspaceDrive
	HyperspaceTechnology         *hyperspaceTechnology
	ImpulseDrive                 *impulseDrive
	IntergalacticResearchNetwork *intergalacticResearchNetwork
	IonTechnology                *ionTechnology
	LaserTechnology              *laserTechnology
	PlasmaTechnology             *plasmaTechnology
	ShieldingTechnology          *shieldingTechnology
	WeaponsTechnology            *weaponsTechnology
}

// ByID gets an object by id
func (o *ObjsStruct) ByID(id ID) BaseOgameObj {
	switch id {
	case AllianceDepotID:
		return o.AllianceDepot
	case CrystalMineID:
		return o.CrystalMine
	case CrystalStorageID:
		return o.CrystalStorage
	case DeuteriumSynthesizerID:
		return o.DeuteriumSynthesizer
	case DeuteriumTankID:
		return o.DeuteriumTank
	case FusionReactorID:
		return o.FusionReactor
	case MetalMineID:
		return o.MetalMine
	case MetalStorageID:
		return o.MetalStorage
	case MissileSiloID:
		return o.MissileSilo
	case NaniteFactoryID:
		return o.NaniteFactory
	case ResearchLabID:
		return o.ResearchLab
	case RoboticsFactoryID:
		return o.RoboticsFactory
	case SeabedDeuteriumDenID:
		return o.SeabedDeuteriumDen
	case ShieldedMetalDenID:
		return o.ShieldedMetalDen
	case ShipyardID:
		return o.Shipyard
	case SolarPlantID:
		return o.SolarPlant
	case SpaceDockID:
		return o.SpaceDock
	case LunarBaseID:
		return o.LunarBase
	case SensorPhalanxID:
		return o.SensorPhalanx
	case JumpGateID:
		return o.JumpGate
	case TerraformerID:
		return o.Terraformer
	case UndergroundCrystalDenID:
		return o.UndergroundCrystalDen
	case SolarSatelliteID:
		return o.SolarSatellite
	case AntiBallisticMissilesID:
		return o.AntiBallisticMissiles
	case GaussCannonID:
		return o.GaussCannon
	case HeavyLaserID:
		return o.HeavyLaser
	case InterplanetaryMissilesID:
		return o.InterplanetaryMissiles
	case IonCannonID:
		return o.IonCannon
	case LargeShieldDomeID:
		return o.LargeShieldDome
	case LightLaserID:
		return o.LightLaser
	case PlasmaTurretID:
		return o.PlasmaTurret
	case RocketLauncherID:
		return o.RocketLauncher
	case SmallShieldDomeID:
		return o.SmallShieldDome
	case BattlecruiserID:
		return o.Battlecruiser
	case BattleshipID:
		return o.Battleship
	case BomberID:
		return o.Bomber
	case ColonyShipID:
		return o.ColonyShip
	case CruiserID:
		return o.Cruiser
	case DeathstarID:
		return o.Deathstar
	case DestroyerID:
		return o.Destroyer
	case EspionageProbeID:
		return o.EspionageProbe
	case HeavyFighterID:
		return o.HeavyFighter
	case LargeCargoID:
		return o.LargeCargo
	case LightFighterID:
		return o.LightFighter
	case RecyclerID:
		return o.Recycler
	case SmallCargoID:
		return o.SmallCargo
	case ArmourTechnologyID:
		return o.ArmourTechnology
	case AstrophysicsID:
		return o.Astrophysics
	case CombustionDriveID:
		return o.CombustionDrive
	case ComputerTechnologyID:
		return o.ComputerTechnology
	case EnergyTechnologyID:
		return o.EnergyTechnology
	case EspionageTechnologyID:
		return o.EspionageTechnology
	case GravitonTechnologyID:
		return o.GravitonTechnology
	case HyperspaceDriveID:
		return o.HyperspaceDrive
	case HyperspaceTechnologyID:
		return o.HyperspaceTechnology
	case ImpulseDriveID:
		return o.ImpulseDrive
	case IntergalacticResearchNetworkID:
		return o.IntergalacticResearchNetwork
	case IonTechnologyID:
		return o.IonTechnology
	case LaserTechnologyID:
		return o.LaserTechnology
	case PlasmaTechnologyID:
		return o.PlasmaTechnology
	case ShieldingTechnologyID:
		return o.ShieldingTechnology
	case WeaponsTechnologyID:
		return o.WeaponsTechnology
	}
	return nil
}

// Objs all ogame objects
var Objs = ObjsStruct{
	AllianceDepot:                AllianceDepot,
	CrystalMine:                  CrystalMine,
	CrystalStorage:               CrystalStorage,
	DeuteriumSynthesizer:         DeuteriumSynthesizer,
	DeuteriumTank:                DeuteriumTank,
	FusionReactor:                FusionReactor,
	MetalMine:                    MetalMine,
	MetalStorage:                 MetalStorage,
	MissileSilo:                  MissileSilo,
	NaniteFactory:                NaniteFactory,
	ResearchLab:                  ResearchLab,
	RoboticsFactory:              RoboticsFactory,
	SeabedDeuteriumDen:           SeabedDeuteriumDen,
	ShieldedMetalDen:             ShieldedMetalDen,
	Shipyard:                     Shipyard,
	SolarPlant:                   SolarPlant,
	SpaceDock:                    SpaceDock,
	LunarBase:                    LunarBase,
	SensorPhalanx:                SensorPhalanx,
	JumpGate:                     JumpGate,
	Terraformer:                  Terraformer,
	UndergroundCrystalDen:        UndergroundCrystalDen,
	SolarSatellite:               SolarSatellite,
	AntiBallisticMissiles:        AntiBallisticMissiles,
	GaussCannon:                  GaussCannon,
	HeavyLaser:                   HeavyLaser,
	InterplanetaryMissiles:       InterplanetaryMissiles,
	IonCannon:                    IonCannon,
	LargeShieldDome:              LargeShieldDome,
	LightLaser:                   LightLaser,
	PlasmaTurret:                 PlasmaTurret,
	RocketLauncher:               RocketLauncher,
	SmallShieldDome:              SmallShieldDome,
	Battlecruiser:                Battlecruiser,
	Battleship:                   Battleship,
	Bomber:                       Bomber,
	ColonyShip:                   ColonyShip,
	Cruiser:                      Cruiser,
	Deathstar:                    Deathstar,
	Destroyer:                    Destroyer,
	EspionageProbe:               EspionageProbe,
	HeavyFighter:                 HeavyFighter,
	LargeCargo:                   LargeCargo,
	LightFighter:                 LightFighter,
	Recycler:                     Recycler,
	SmallCargo:                   SmallCargo,
	ArmourTechnology:             ArmourTechnology,
	Astrophysics:                 Astrophysics,
	CombustionDrive:              CombustionDrive,
	ComputerTechnology:           ComputerTechnology,
	EnergyTechnology:             EnergyTechnology,
	EspionageTechnology:          EspionageTechnology,
	GravitonTechnology:           GravitonTechnology,
	HyperspaceDrive:              HyperspaceDrive,
	HyperspaceTechnology:         HyperspaceTechnology,
	ImpulseDrive:                 ImpulseDrive,
	IntergalacticResearchNetwork: IntergalacticResearchNetwork,
	IonTechnology:                IonTechnology,
	LaserTechnology:              LaserTechnology,
	PlasmaTechnology:             PlasmaTechnology,
	ShieldingTechnology:          ShieldingTechnology,
	WeaponsTechnology:            WeaponsTechnology,
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
