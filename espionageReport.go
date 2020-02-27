package ogame

import "time"

// EspionageReport detailed espionage report
type EspionageReport struct {
	Resources
	ID                           int64
	Username                     string
	LastActivity                 int64
	CounterEspionage             int64
	APIKey                       string
	HasFleet                     bool
	HasDefenses                  bool
	HasBuildings                 bool
	HasResearches                bool
	HonorableTarget              bool
	IsBandit                     bool
	IsStarlord                   bool
	IsInactive                   bool
	IsLongInactive               bool
	MetalMine                    *int64 // ResourcesBuildings
	CrystalMine                  *int64
	DeuteriumSynthesizer         *int64
	SolarPlant                   *int64
	FusionReactor                *int64
	SolarSatellite               *int64
	MetalStorage                 *int64
	CrystalStorage               *int64
	DeuteriumTank                *int64
	RoboticsFactory              *int64 // Facilities
	Shipyard                     *int64
	ResearchLab                  *int64
	AllianceDepot                *int64
	MissileSilo                  *int64
	NaniteFactory                *int64
	Terraformer                  *int64
	SpaceDock                    *int64
	LunarBase                    *int64
	SensorPhalanx                *int64
	JumpGate                     *int64
	EnergyTechnology             *int64 // Researches
	LaserTechnology              *int64
	IonTechnology                *int64
	HyperspaceTechnology         *int64
	PlasmaTechnology             *int64
	CombustionDrive              *int64
	ImpulseDrive                 *int64
	HyperspaceDrive              *int64
	EspionageTechnology          *int64
	ComputerTechnology           *int64
	Astrophysics                 *int64
	IntergalacticResearchNetwork *int64
	GravitonTechnology           *int64
	WeaponsTechnology            *int64
	ShieldingTechnology          *int64
	ArmourTechnology             *int64
	RocketLauncher               *int64 // Defenses
	LightLaser                   *int64
	HeavyLaser                   *int64
	GaussCannon                  *int64
	IonCannon                    *int64
	PlasmaTurret                 *int64
	SmallShieldDome              *int64
	LargeShieldDome              *int64
	AntiBallisticMissiles        *int64
	InterplanetaryMissiles       *int64
	LightFighter                 *int64 // Fleets
	HeavyFighter                 *int64
	Cruiser                      *int64
	Battleship                   *int64
	Battlecruiser                *int64
	Bomber                       *int64
	Destroyer                    *int64
	Deathstar                    *int64
	SmallCargo                   *int64
	LargeCargo                   *int64
	ColonyShip                   *int64
	Recycler                     *int64
	EspionageProbe               *int64
	Crawler                      *int64
	Reaper                       *int64
	Pathfinder                   *int64
	Coordinate                   Coordinate
	Type                         EspionageReportType
	Date                         time.Time
}

func i64(v *int64) int64 {
	if v == nil {
		return 0
	}
	return *v
}

// ShipsInfos returns a ShipsInfos struct from the espionage report
func (r EspionageReport) ShipsInfos() *ShipsInfos {
	if !r.HasFleet {
		return nil
	}
	return &ShipsInfos{
		LightFighter:   i64(r.LightFighter),
		HeavyFighter:   i64(r.HeavyFighter),
		Cruiser:        i64(r.Cruiser),
		Battleship:     i64(r.Battleship),
		Battlecruiser:  i64(r.Battlecruiser),
		Bomber:         i64(r.Bomber),
		Destroyer:      i64(r.Destroyer),
		Deathstar:      i64(r.Deathstar),
		SmallCargo:     i64(r.SmallCargo),
		LargeCargo:     i64(r.LargeCargo),
		ColonyShip:     i64(r.ColonyShip),
		Recycler:       i64(r.Recycler),
		EspionageProbe: i64(r.EspionageProbe),
		SolarSatellite: i64(r.SolarSatellite),
		Crawler:        i64(r.Crawler),
		Reaper:         i64(r.Reaper),
		Pathfinder:     i64(r.Pathfinder),
	}
}

// DefensesInfos returns a DefensesInfos struct from the espionage report
func (r EspionageReport) DefensesInfos() *DefensesInfos {
	if !r.HasDefenses {
		return nil
	}
	return &DefensesInfos{
		RocketLauncher:         i64(r.RocketLauncher),
		LightLaser:             i64(r.LightLaser),
		HeavyLaser:             i64(r.HeavyLaser),
		GaussCannon:            i64(r.GaussCannon),
		IonCannon:              i64(r.IonCannon),
		PlasmaTurret:           i64(r.PlasmaTurret),
		SmallShieldDome:        i64(r.SmallShieldDome),
		LargeShieldDome:        i64(r.LargeShieldDome),
		AntiBallisticMissiles:  i64(r.AntiBallisticMissiles),
		InterplanetaryMissiles: i64(r.InterplanetaryMissiles),
	}
}

// PlunderRatio returns the plunder ratio
func (r EspionageReport) PlunderRatio(characterClass CharacterClass) float64 {
	plunderRatio := 0.5
	if r.IsInactive && characterClass == Discoverer {
		plunderRatio = 0.75
	}
	if r.IsBandit {
		plunderRatio = 1
	} else if !r.IsInactive && r.IsStarlord {
		plunderRatio = 0.75
	}
	return plunderRatio
}

// Loot returns the possible loot
func (r EspionageReport) Loot(characterClass CharacterClass) Resources {
	plunderRatio := r.PlunderRatio(characterClass)
	return Resources{
		Metal:     int64(float64(r.Metal) * plunderRatio),
		Crystal:   int64(float64(r.Crystal) * plunderRatio),
		Deuterium: int64(float64(r.Deuterium) * plunderRatio),
	}
}
