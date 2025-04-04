package ogame

import (
	"github.com/alaingilbert/ogame/pkg/utils"
	"time"
)

// EspionageReport detailed espionage report
type EspionageReport struct {
	Resources
	ID                           int64
	PlayerID                     int64
	Username                     string
	CharacterClass               CharacterClass
	AllianceClass                AllianceClass
	LastActivity                 int64
	CounterEspionage             int64
	APIKey                       string
	HasFleetInformation          bool // Either or not we sent enough probes to get the fleet information
	HasDefensesInformation       bool // Either or not we sent enough probes to get the defenses information
	HasBuildingsInformation      bool // Either or not we sent enough probes to get the buildings information
	HasResearchesInformation     bool // Either or not we sent enough probes to get the researches information
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

// ResourcesBuildings returns a ResourcesBuildings struct from the espionage report
func (r EspionageReport) ResourcesBuildings() *ResourcesBuildings {
	if !r.HasBuildingsInformation {
		return nil
	}
	return &ResourcesBuildings{
		MetalMine:            utils.Deref(r.MetalMine),
		CrystalMine:          utils.Deref(r.CrystalMine),
		DeuteriumSynthesizer: utils.Deref(r.DeuteriumSynthesizer),
		SolarPlant:           utils.Deref(r.SolarPlant),
		FusionReactor:        utils.Deref(r.FusionReactor),
		SolarSatellite:       utils.Deref(r.SolarSatellite),
		MetalStorage:         utils.Deref(r.MetalStorage),
		CrystalStorage:       utils.Deref(r.CrystalStorage),
		DeuteriumTank:        utils.Deref(r.DeuteriumTank),
	}
}

// Facilities returns a Facilities struct from the espionage report
func (r EspionageReport) Facilities() *Facilities {
	if !r.HasBuildingsInformation {
		return nil
	}
	return &Facilities{
		RoboticsFactory: utils.Deref(r.RoboticsFactory),
		Shipyard:        utils.Deref(r.Shipyard),
		ResearchLab:     utils.Deref(r.ResearchLab),
		AllianceDepot:   utils.Deref(r.AllianceDepot),
		MissileSilo:     utils.Deref(r.MissileSilo),
		NaniteFactory:   utils.Deref(r.NaniteFactory),
		Terraformer:     utils.Deref(r.Terraformer),
		SpaceDock:       utils.Deref(r.SpaceDock),
		LunarBase:       utils.Deref(r.LunarBase),
		SensorPhalanx:   utils.Deref(r.SensorPhalanx),
		JumpGate:        utils.Deref(r.JumpGate),
	}
}

// Researches returns a Researches struct from the espionage report
func (r EspionageReport) Researches() *Researches {
	if !r.HasResearchesInformation {
		return nil
	}
	return &Researches{
		EnergyTechnology:             utils.Deref(r.EnergyTechnology),
		LaserTechnology:              utils.Deref(r.LaserTechnology),
		IonTechnology:                utils.Deref(r.IonTechnology),
		HyperspaceTechnology:         utils.Deref(r.HyperspaceTechnology),
		PlasmaTechnology:             utils.Deref(r.PlasmaTechnology),
		CombustionDrive:              utils.Deref(r.CombustionDrive),
		ImpulseDrive:                 utils.Deref(r.ImpulseDrive),
		HyperspaceDrive:              utils.Deref(r.HyperspaceDrive),
		EspionageTechnology:          utils.Deref(r.EspionageTechnology),
		ComputerTechnology:           utils.Deref(r.ComputerTechnology),
		Astrophysics:                 utils.Deref(r.Astrophysics),
		IntergalacticResearchNetwork: utils.Deref(r.IntergalacticResearchNetwork),
		GravitonTechnology:           utils.Deref(r.GravitonTechnology),
		WeaponsTechnology:            utils.Deref(r.WeaponsTechnology),
		ShieldingTechnology:          utils.Deref(r.ShieldingTechnology),
		ArmourTechnology:             utils.Deref(r.ArmourTechnology),
	}
}

// ShipsInfos returns a ShipsInfos struct from the espionage report
func (r EspionageReport) ShipsInfos() *ShipsInfos {
	if !r.HasFleetInformation {
		return nil
	}
	return &ShipsInfos{
		LightFighter:   utils.Deref(r.LightFighter),
		HeavyFighter:   utils.Deref(r.HeavyFighter),
		Cruiser:        utils.Deref(r.Cruiser),
		Battleship:     utils.Deref(r.Battleship),
		Battlecruiser:  utils.Deref(r.Battlecruiser),
		Bomber:         utils.Deref(r.Bomber),
		Destroyer:      utils.Deref(r.Destroyer),
		Deathstar:      utils.Deref(r.Deathstar),
		SmallCargo:     utils.Deref(r.SmallCargo),
		LargeCargo:     utils.Deref(r.LargeCargo),
		ColonyShip:     utils.Deref(r.ColonyShip),
		Recycler:       utils.Deref(r.Recycler),
		EspionageProbe: utils.Deref(r.EspionageProbe),
		SolarSatellite: utils.Deref(r.SolarSatellite),
		Crawler:        utils.Deref(r.Crawler),
		Reaper:         utils.Deref(r.Reaper),
		Pathfinder:     utils.Deref(r.Pathfinder),
	}
}

// DefensesInfos returns a DefensesInfos struct from the espionage report
func (r EspionageReport) DefensesInfos() *DefensesInfos {
	if !r.HasDefensesInformation {
		return nil
	}
	return &DefensesInfos{
		RocketLauncher:         utils.Deref(r.RocketLauncher),
		LightLaser:             utils.Deref(r.LightLaser),
		HeavyLaser:             utils.Deref(r.HeavyLaser),
		GaussCannon:            utils.Deref(r.GaussCannon),
		IonCannon:              utils.Deref(r.IonCannon),
		PlasmaTurret:           utils.Deref(r.PlasmaTurret),
		SmallShieldDome:        utils.Deref(r.SmallShieldDome),
		LargeShieldDome:        utils.Deref(r.LargeShieldDome),
		AntiBallisticMissiles:  utils.Deref(r.AntiBallisticMissiles),
		InterplanetaryMissiles: utils.Deref(r.InterplanetaryMissiles),
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

// IsDefenceless returns either or not the scanned planet has any defense (either ships or defense) against an attack
// with ships. If no ShipsInfos or DefensesInfos is including in the espionage report due to the lack of enough probes,
// the planet is assumed to be not defenceless.
func (r EspionageReport) IsDefenceless() bool {
	return r.HasFleetInformation &&
		r.HasDefensesInformation &&
		!r.ShipsInfos().HasShips() &&
		!r.DefensesInfos().HasShipDefense()
}
