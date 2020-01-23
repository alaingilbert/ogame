package ogame

import "strconv"

// MissionID represent a mission id
type MissionID int

func (m MissionID) String() string {
	switch m {
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
	default:
		return strconv.FormatInt(int64(m), 10)
	}
}

// Speed represent a fleet speed
type Speed int

// Int64 returns an integer value of the speed
func (s Speed) Int64() int64 {
	return int64(s)
}

// Int returns an integer value of the speed
// Deprecated: backward compatibility
func (s Speed) Int() int64 {
	return int64(s)
}

func (s Speed) String() string {
	switch s {
	case TenPercent:
		return "10%"
	case TwentyPercent:
		return "20%"
	case ThirtyPercent:
		return "30%"
	case FourtyPercent:
		return "40%"
	case FiftyPercent:
		return "50%"
	case SixtyPercent:
		return "60%"
	case SeventyPercent:
		return "70%"
	case EightyPercent:
		return "80%"
	case NinetyPercent:
		return "90%"
	case HundredPercent:
		return "100%"
	default:
		return strconv.FormatInt(int64(s), 10)
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

// CharacterClass ...
type CharacterClass int64

// OGame constants
const (
	NoClass    CharacterClass = 0
	Collector  CharacterClass = 1
	General    CharacterClass = 2
	Discoverer CharacterClass = 3

	PlanetType CelestialType = 1
	DebrisType CelestialType = 2
	MoonType   CelestialType = 3

	//Buildings
	MetalMineID                    ID = 1
	CrystalMineID                  ID = 2
	DeuteriumSynthesizerID         ID = 3
	SolarPlantID                   ID = 4
	FusionReactorID                ID = 12
	MetalStorageID                 ID = 22
	CrystalStorageID               ID = 23
	DeuteriumTankID                ID = 24
	ShieldedMetalDenID             ID = 25
	UndergroundCrystalDenID        ID = 26
	SeabedDeuteriumDenID           ID = 27
	AllianceDepotID                ID = 34 // Facilities
	RoboticsFactoryID              ID = 14
	ShipyardID                     ID = 21
	ResearchLabID                  ID = 31
	MissileSiloID                  ID = 44
	NaniteFactoryID                ID = 15
	TerraformerID                  ID = 33
	SpaceDockID                    ID = 36
	LunarBaseID                    ID = 41 // Moon facilities
	SensorPhalanxID                ID = 42
	JumpGateID                     ID = 43
	RocketLauncherID               ID = 401 // Defense
	LightLaserID                   ID = 402
	HeavyLaserID                   ID = 403
	GaussCannonID                  ID = 404
	IonCannonID                    ID = 405
	PlasmaTurretID                 ID = 406
	SmallShieldDomeID              ID = 407
	LargeShieldDomeID              ID = 408
	AntiBallisticMissilesID        ID = 502
	InterplanetaryMissilesID       ID = 503
	SmallCargoID                   ID = 202 // Ships
	LargeCargoID                   ID = 203
	LightFighterID                 ID = 204
	HeavyFighterID                 ID = 205
	CruiserID                      ID = 206
	BattleshipID                   ID = 207
	ColonyShipID                   ID = 208
	RecyclerID                     ID = 209
	EspionageProbeID               ID = 210
	BomberID                       ID = 211
	SolarSatelliteID               ID = 212
	DestroyerID                    ID = 213
	DeathstarID                    ID = 214
	BattlecruiserID                ID = 215
	CrawlerID                      ID = 217
	ReaperID                       ID = 218
	PathfinderID                   ID = 219
	EspionageTechnologyID          ID = 106 // Research
	ComputerTechnologyID           ID = 108
	WeaponsTechnologyID            ID = 109
	ShieldingTechnologyID          ID = 110
	ArmourTechnologyID             ID = 111
	EnergyTechnologyID             ID = 113
	HyperspaceTechnologyID         ID = 114
	CombustionDriveID              ID = 115
	ImpulseDriveID                 ID = 117
	HyperspaceDriveID              ID = 118
	LaserTechnologyID              ID = 120
	IonTechnologyID                ID = 121
	PlasmaTechnologyID             ID = 122
	IntergalacticResearchNetworkID ID = 123
	AstrophysicsID                 ID = 124
	GravitonTechnologyID           ID = 199

	// Missions
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

	// Speeds
	TenPercent     Speed = 1
	TwentyPercent  Speed = 2
	ThirtyPercent  Speed = 3
	FourtyPercent  Speed = 4
	FiftyPercent   Speed = 5
	SixtyPercent   Speed = 6
	SeventyPercent Speed = 7
	EightyPercent  Speed = 8
	NinetyPercent  Speed = 9
	HundredPercent Speed = 10
)
