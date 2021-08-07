package ogame

type TechInfos struct {
	//ResourcesBuildings
	MetalMine            int64
	CrystalMine          int64
	DeuteriumSynthesizer int64
	SolarPlant           int64
	FusionReactor        int64
	SolarSatellite       int64
	MetalStorage         int64
	CrystalStorage       int64
	DeuteriumTank        int64

	//Facilities
	RoboticsFactory int64
	Shipyard        int64
	ResearchLab     int64
	AllianceDepot   int64
	MissileSilo     int64
	NaniteFactory   int64
	Terraformer     int64
	SpaceDock       int64
	LunarBase       int64
	SensorPhalanx   int64
	JumpGate        int64

	//ShipsInfos
	LightFighter   int64
	HeavyFighter   int64
	Cruiser        int64
	Battleship     int64
	Battlecruiser  int64
	Bomber         int64
	Destroyer      int64
	Deathstar      int64
	SmallCargo     int64
	LargeCargo     int64
	ColonyShip     int64
	Recycler       int64
	EspionageProbe int64
	/*SolarSatellite int64*/
	Crawler        int64
	Reaper         int64
	Pathfinder     int64

	//DefensesInfos
	RocketLauncher         int64
	LightLaser             int64
	HeavyLaser             int64
	GaussCannon            int64
	IonCannon              int64
	PlasmaTurret           int64
	SmallShieldDome        int64
	LargeShieldDome        int64
	AntiBallisticMissiles  int64
	InterplanetaryMissiles int64

	// Researches
	EnergyTechnology             int64
	LaserTechnology              int64
	IonTechnology                int64
	HyperspaceTechnology         int64
	PlasmaTechnology             int64
	CombustionDrive              int64
	ImpulseDrive                 int64
	HyperspaceDrive              int64
	EspionageTechnology          int64
	ComputerTechnology           int64
	Astrophysics                 int64
	IntergalacticResearchNetwork int64
	GravitonTechnology           int64
	WeaponsTechnology            int64
	ShieldingTechnology          int64
	ArmourTechnology             int64
}

// ToPtr returns a pointer to self
func (s TechInfos) ToPtr() *TechInfos {
	return &s
}

// ByID get number/level of tech by tech id
func (s TechInfos) ByID(id ID) int64 {
	switch id {
	// ResourcesBuildings
	case MetalMineID:
		return s.MetalMine
	case CrystalMineID:
		return s.CrystalMine
	case DeuteriumSynthesizerID:
		return s.DeuteriumSynthesizer
	case SolarPlantID:
		return s.SolarPlant
	case FusionReactorID:
		return s.FusionReactor
	case SolarSatelliteID:
		return s.SolarSatellite
	case MetalStorageID:
		return s.MetalStorage
	case CrystalStorageID:
		return s.CrystalStorage
	case DeuteriumTankID:
		return s.LunarBase

	// Facilities
	case RoboticsFactoryID:
		return s.RoboticsFactory
	case ShipyardID:
		return s.Shipyard
	case ResearchLabID:
		return s.ResearchLab
	case AllianceDepotID:
		return s.AllianceDepot
	case MissileSiloID:
		return s.MissileSilo
	case NaniteFactoryID:
		return s.NaniteFactory
	case TerraformerID:
		return s.Terraformer
	case SpaceDockID:
		return s.SpaceDock
	case LunarBaseID:
		return s.LunarBase
	case SensorPhalanxID:
		return s.SensorPhalanx
	case JumpGateID:
		return s.JumpGate

	//ShipsInfos
	case LightFighterID:
		return s.LightFighter
	case HeavyFighterID:
		return s.HeavyFighter
	case CruiserID:
		return s.Cruiser
	case BattleshipID:
		return s.Battleship
	case BattlecruiserID:
		return s.Battlecruiser
	case BomberID:
		return s.Bomber
	case DestroyerID:
		return s.Destroyer
	case DeathstarID:
		return s.Deathstar
	case SmallCargoID:
		return s.SmallCargo
	case LargeCargoID:
		return s.LargeCargo
	case ColonyShipID:
		return s.ColonyShip
	case RecyclerID:
		return s.Recycler
	case EspionageProbeID:
		return s.EspionageProbe
//	case SolarSatelliteID:
//		return s.SolarSatellite
	case CrawlerID:
		return s.Crawler
	case ReaperID:
		return s.Reaper
	case PathfinderID:
		return s.Pathfinder

	// DefensesInfos
	case RocketLauncherID:
		return s.RocketLauncher
	case LightLaserID:
		return s.LightLaser
	case HeavyLaserID:
		return s.HeavyLaser
	case GaussCannonID:
		return s.GaussCannon
	case IonCannonID:
		return s.IonCannon
	case PlasmaTurretID:
		return s.PlasmaTurret
	case SmallShieldDomeID:
		return s.SmallShieldDome
	case LargeShieldDomeID:
		return s.LargeShieldDome
	case AntiBallisticMissilesID:
		return s.AntiBallisticMissiles
	case InterplanetaryMissilesID:
		return s.InterplanetaryMissiles

// Researches
	case EnergyTechnologyID:
		return s.EnergyTechnology
	case LaserTechnologyID:
		return s.LaserTechnology
	case IonTechnologyID:
		return s.IonTechnology
	case HyperspaceTechnologyID:
		return s.HyperspaceTechnology
	case PlasmaTechnologyID:
		return s.PlasmaTechnology
	case CombustionDriveID:
		return s.CombustionDrive
	case ImpulseDriveID:
		return s.ImpulseDrive
	case HyperspaceDriveID:
		return s.HyperspaceDrive
	case EspionageTechnologyID:
		return s.EspionageTechnology
	case ComputerTechnologyID:
		return s.ComputerTechnology
	case AstrophysicsID:
		return s.Astrophysics
	case IntergalacticResearchNetworkID:
		return s.IntergalacticResearchNetwork
	case GravitonTechnologyID:
		return s.GravitonTechnology
	case WeaponsTechnologyID:
		return s.WeaponsTechnology
	case ShieldingTechnologyID:
		return s.ShieldingTechnology
	case ArmourTechnologyID:
		return s.ArmourTechnology

	default:
		return 0
	}
}

// Set sets the techs value using the id
func (s *TechInfos) Set(id ID, val int64) {
	switch id {
	// ResourcesBuildings
	case MetalMineID:
		s.MetalMine = val
	case CrystalMineID:
		s.CrystalMine = val
	case DeuteriumSynthesizerID:
		s.DeuteriumSynthesizer = val
	case SolarPlantID:
		s.SolarPlant = val
	case FusionReactorID:
		s.FusionReactor = val
	case SolarSatelliteID:
		s.SolarSatellite = val
	case MetalStorageID:
		s.MetalStorage = val
	case CrystalStorageID:
		s.CrystalStorage = val
	case DeuteriumTankID:
		s.DeuteriumTank	= val

	// Facilities
	case RoboticsFactoryID:
		s.RoboticsFactory = val
	case ShipyardID:
		s.Shipyard = val
	case ResearchLabID:
		s.ResearchLab = val
	case AllianceDepotID:
		s.AllianceDepot = val
	case MissileSiloID:
		s.MissileSilo = val
	case NaniteFactoryID:
		s.NaniteFactory = val
	case TerraformerID:
		s.Terraformer = val
	case SpaceDockID:
		s.SpaceDock = val
	case LunarBaseID:
		s.LunarBase = val
	case SensorPhalanxID:
		s.SensorPhalanx = val
	case JumpGateID:
		s.JumpGate = val


	// ShipsInfos
	case LightFighterID:
		s.LightFighter = val
	case HeavyFighterID:
		s.HeavyFighter = val
	case CruiserID:
		s.Cruiser = val
	case BattleshipID:
		s.Battleship = val
	case BattlecruiserID:
		s.Battlecruiser = val
	case BomberID:
		s.Bomber = val
	case DestroyerID:
		s.Destroyer = val
	case DeathstarID:
		s.Deathstar = val
	case SmallCargoID:
		s.SmallCargo = val
	case LargeCargoID:
		s.LargeCargo = val
	case ColonyShipID:
		s.ColonyShip = val
	case RecyclerID:
		s.Recycler = val
	case EspionageProbeID:
		s.EspionageProbe = val
//	case SolarSatelliteID:
//		s.SolarSatellite = val
	case CrawlerID:
		s.Crawler = val
	case ReaperID:
		s.Reaper = val
	case PathfinderID:
		s.Pathfinder = val

	// DefensesInfos
	case RocketLauncherID:
		s.RocketLauncher = val
	case LightLaserID:
		s.LightLaser = val
	case HeavyLaserID:
		s.HeavyLaser = val
	case GaussCannonID:
		s.GaussCannon = val
	case IonCannonID:
		s.IonCannon = val
	case PlasmaTurretID:
		s.PlasmaTurret = val
	case SmallShieldDomeID:
		s.SmallShieldDome = val
	case LargeShieldDomeID:
		s.LargeShieldDome = val
	case AntiBallisticMissilesID:
		s.AntiBallisticMissiles = val
	case InterplanetaryMissilesID:
		s.InterplanetaryMissiles = val

		// Researches
	case EnergyTechnologyID:
		s.EnergyTechnology = val
	case LaserTechnologyID:
		s.LaserTechnology = val
	case IonTechnologyID:
		s.IonTechnology = val
	case HyperspaceTechnologyID:
		s.HyperspaceTechnology = val
	case PlasmaTechnologyID:
		s.PlasmaTechnology = val
	case CombustionDriveID:
		s.CombustionDrive = val
	case ImpulseDriveID:
		s.ImpulseDrive = val
	case HyperspaceDriveID:
		s.HyperspaceDrive = val
	case EspionageTechnologyID:
		s.EspionageTechnology = val
	case ComputerTechnologyID:
		s.ComputerTechnology = val
	case AstrophysicsID:
		s.Astrophysics = val
	case IntergalacticResearchNetworkID:
		s.IntergalacticResearchNetwork = val
	case GravitonTechnologyID:
		s.GravitonTechnology = val
	case WeaponsTechnologyID:
		s.WeaponsTechnology = val
	case ShieldingTechnologyID:
		s.ShieldingTechnology = val
	case ArmourTechnologyID:
		s.ArmourTechnology	 = val
	}
}

// FromQuantifiables convert an array of Quantifiable to a TechInfos
func (s TechInfos) FromQuantifiables(in []Quantifiable) (out TechInfos) {
	for _, item := range in {
		out.Set(item.ID, item.Nbr)
	}
	return
}

// ToQuantifiables convert a TechInfos to an array of Quantifiable
func (s TechInfos) ToQuantifiables() []Quantifiable {
	out := make([]Quantifiable, 0)

	out = append(out, Quantifiable{ID: MetalMineID, Nbr: s.MetalMine})
	out = append(out, Quantifiable{ID: CrystalMineID, Nbr: s.CrystalMine})
	out = append(out, Quantifiable{ID: DeuteriumSynthesizerID, Nbr: s.DeuteriumSynthesizer})
	out = append(out, Quantifiable{ID: SolarPlantID, Nbr: s.SolarPlant})
	out = append(out, Quantifiable{ID: FusionReactorID, Nbr: s.FusionReactor})
	out = append(out, Quantifiable{ID: SolarSatelliteID, Nbr: s.SolarSatellite})
	out = append(out, Quantifiable{ID: MetalStorageID, Nbr: s.MetalStorage})
	out = append(out, Quantifiable{ID: CrystalStorageID, Nbr: s.CrystalStorage})
	out = append(out, Quantifiable{ID: DeuteriumTankID, Nbr: s.DeuteriumTank})

	out = append(out, Quantifiable{ID: RoboticsFactoryID, Nbr: s.RoboticsFactory})
	out = append(out, Quantifiable{ID: ShipyardID, Nbr: s.Shipyard})
	out = append(out, Quantifiable{ID: ResearchLabID, Nbr: s.ResearchLab})
	out = append(out, Quantifiable{ID: AllianceDepotID, Nbr: s.AllianceDepot})
	out = append(out, Quantifiable{ID: MissileSiloID, Nbr: s.MissileSilo})
	out = append(out, Quantifiable{ID: NaniteFactoryID, Nbr: s.NaniteFactory})
	out = append(out, Quantifiable{ID: TerraformerID, Nbr: s.Terraformer})
	out = append(out, Quantifiable{ID: SpaceDockID, Nbr: s.SpaceDock})
	out = append(out, Quantifiable{ID: LunarBaseID, Nbr: s.LunarBase})
	out = append(out, Quantifiable{ID: SensorPhalanxID, Nbr: s.SensorPhalanx})
	out = append(out, Quantifiable{ID: JumpGateID, Nbr: s.JumpGate})

	out = append(out, Quantifiable{ID: LightFighterID, Nbr: s.LightFighter})
	out = append(out, Quantifiable{ID: HeavyFighterID, Nbr: s.HeavyFighter})
	out = append(out, Quantifiable{ID: CruiserID, Nbr: s.Cruiser})
	out = append(out, Quantifiable{ID: BattleshipID, Nbr: s.Battleship})
	out = append(out, Quantifiable{ID: BattlecruiserID, Nbr: s.Battlecruiser})
	out = append(out, Quantifiable{ID: BomberID, Nbr: s.Bomber})
	out = append(out, Quantifiable{ID: DestroyerID, Nbr: s.Destroyer})
	out = append(out, Quantifiable{ID: DeathstarID, Nbr: s.Deathstar})
	out = append(out, Quantifiable{ID: SmallCargoID, Nbr: s.SmallCargo})
	out = append(out, Quantifiable{ID: LargeCargoID, Nbr: s.LargeCargo})
	out = append(out, Quantifiable{ID: ColonyShipID, Nbr: s.ColonyShip})
	out = append(out, Quantifiable{ID: RecyclerID, Nbr: s.Recycler})
	out = append(out, Quantifiable{ID: EspionageProbeID, Nbr: s.EspionageProbe})
	//out = append(out, Quantifiable{ID: //SolarSatelliteID, Nbr: s.//SolarSatellite})
	out = append(out, Quantifiable{ID: CrawlerID, Nbr: s.Crawler})
	out = append(out, Quantifiable{ID: ReaperID, Nbr: s.Reaper})
	out = append(out, Quantifiable{ID: PathfinderID, Nbr: s.Pathfinder})

	out = append(out, Quantifiable{ID: RocketLauncherID, Nbr: s.RocketLauncher})
	out = append(out, Quantifiable{ID: LightLaserID, Nbr: s.LightLaser})
	out = append(out, Quantifiable{ID: HeavyLaserID, Nbr: s.HeavyLaser})
	out = append(out, Quantifiable{ID: GaussCannonID, Nbr: s.GaussCannon})
	out = append(out, Quantifiable{ID: IonCannonID, Nbr: s.IonCannon})
	out = append(out, Quantifiable{ID: PlasmaTurretID, Nbr: s.PlasmaTurret})
	out = append(out, Quantifiable{ID: SmallShieldDomeID, Nbr: s.SmallShieldDome})
	out = append(out, Quantifiable{ID: LargeShieldDomeID, Nbr: s.LargeShieldDome})
	out = append(out, Quantifiable{ID: AntiBallisticMissilesID, Nbr: s.AntiBallisticMissiles})
	out = append(out, Quantifiable{ID: InterplanetaryMissilesID, Nbr: s.InterplanetaryMissiles})

	out = append(out, Quantifiable{ID: EnergyTechnologyID, Nbr: s.EnergyTechnology})
	out = append(out, Quantifiable{ID: LaserTechnologyID, Nbr: s.LaserTechnology})
	out = append(out, Quantifiable{ID: IonTechnologyID, Nbr: s.IonTechnology})
	out = append(out, Quantifiable{ID: HyperspaceTechnologyID, Nbr: s.HyperspaceTechnology})
	out = append(out, Quantifiable{ID: PlasmaTechnologyID, Nbr: s.PlasmaTechnology})
	out = append(out, Quantifiable{ID: CombustionDriveID, Nbr: s.CombustionDrive})
	out = append(out, Quantifiable{ID: ImpulseDriveID, Nbr: s.ImpulseDrive})
	out = append(out, Quantifiable{ID: HyperspaceDriveID, Nbr: s.HyperspaceDrive})
	out = append(out, Quantifiable{ID: EspionageTechnologyID, Nbr: s.EspionageTechnology})
	out = append(out, Quantifiable{ID: ComputerTechnologyID, Nbr: s.ComputerTechnology})
	out = append(out, Quantifiable{ID: AstrophysicsID, Nbr: s.Astrophysics})
	out = append(out, Quantifiable{ID: IntergalacticResearchNetworkID, Nbr: s.IntergalacticResearchNetwork})
	out = append(out, Quantifiable{ID: GravitonTechnologyID, Nbr: s.GravitonTechnology})
	out = append(out, Quantifiable{ID: WeaponsTechnologyID, Nbr: s.WeaponsTechnology})
	out = append(out, Quantifiable{ID: ShieldingTechnologyID, Nbr: s.ShieldingTechnology})
	out = append(out, Quantifiable{ID: ArmourTechnologyID, Nbr: s.ArmourTechnology})

	return out
}