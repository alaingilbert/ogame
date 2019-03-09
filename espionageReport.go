package ogame

import "time"

// EspionageReport detailed espionage report
type EspionageReport struct {
	Resources
	ID                           int
	Username                     string
	LastActivity                 int
	CounterEspionage             int
	APIKey                       string
	HasFleet                     bool
	HasDefenses                  bool
	HasBuildings                 bool
	HasResearches                bool
	IsBandit                     bool
	IsStarlord                   bool
	IsInactive                   bool
	IsLongInactive               bool
	MetalMine                    *int // ResourcesBuildings
	CrystalMine                  *int
	DeuteriumSynthesizer         *int
	SolarPlant                   *int
	FusionReactor                *int
	SolarSatellite               *int
	MetalStorage                 *int
	CrystalStorage               *int
	DeuteriumTank                *int
	RoboticsFactory              *int // Facilities
	Shipyard                     *int
	ResearchLab                  *int
	AllianceDepot                *int
	MissileSilo                  *int
	NaniteFactory                *int
	Terraformer                  *int
	SpaceDock                    *int
	LunarBase                    *int
	SensorPhalanx                *int
	JumpGate                     *int
	EnergyTechnology             *int // Researches
	LaserTechnology              *int
	IonTechnology                *int
	HyperspaceTechnology         *int
	PlasmaTechnology             *int
	CombustionDrive              *int
	ImpulseDrive                 *int
	HyperspaceDrive              *int
	EspionageTechnology          *int
	ComputerTechnology           *int
	Astrophysics                 *int
	IntergalacticResearchNetwork *int
	GravitonTechnology           *int
	WeaponsTechnology            *int
	ShieldingTechnology          *int
	ArmourTechnology             *int
	RocketLauncher               *int // Defenses
	LightLaser                   *int
	HeavyLaser                   *int
	GaussCannon                  *int
	IonCannon                    *int
	PlasmaTurret                 *int
	SmallShieldDome              *int
	LargeShieldDome              *int
	AntiBallisticMissiles        *int
	InterplanetaryMissiles       *int
	LightFighter                 *int // Fleets
	HeavyFighter                 *int
	Cruiser                      *int
	Battleship                   *int
	Battlecruiser                *int
	Bomber                       *int
	Destroyer                    *int
	Deathstar                    *int
	SmallCargo                   *int
	LargeCargo                   *int
	ColonyShip                   *int
	Recycler                     *int
	EspionageProbe               *int
	Coordinate                   Coordinate
	Type                         EspionageReportType
	Date                         time.Time
}
