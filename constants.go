package ogame

import "strconv"

// MissionID ...
type MissionID int

func (m MissionID) String() string {
	return strconv.Itoa(int(m))
}

// Speed ...
type Speed int

func (s Speed) String() string {
	return strconv.Itoa(int(s))
}

// OGame constants
const (
	// Buildings
	MetalMine                    ID = 1
	CrystalMine                  ID = 2
	DeuteriumSynthesizer         ID = 3
	SolarPlant                   ID = 4
	FusionReactor                ID = 12
	MetalStorage                 ID = 22
	CrystalStorage               ID = 23
	DeuteriumTank                ID = 24
	ShieldedMetalDen             ID = 25
	UndergroundCrystalDen        ID = 26
	SeabedDeuteriumDen           ID = 27
	AllianceDepot                ID = 34 // Facilities
	RoboticsFactory              ID = 14
	Shipyard                     ID = 21
	ResearchLab                  ID = 31
	MissileSilo                  ID = 44
	NaniteFactory                ID = 15
	Terraformer                  ID = 33
	SpaceDock                    ID = 36
	RocketLauncher               ID = 401 // Defense
	LightLaser                   ID = 402
	HeavyLaser                   ID = 403
	GaussCannon                  ID = 404
	IonCannon                    ID = 405
	PlasmaTurret                 ID = 406
	SmallShieldDome              ID = 407
	LargeShieldDome              ID = 408
	AntiBallisticMissiles        ID = 502
	InterplanetaryMissiles       ID = 503
	SmallCargo                   ID = 202 // Ships
	LargeCargo                   ID = 203
	LightFighter                 ID = 204
	HeavyFighter                 ID = 205
	Cruiser                      ID = 206
	Battleship                   ID = 207
	ColonyShip                   ID = 208
	Recycler                     ID = 209
	EspionageProbe               ID = 210
	Bomber                       ID = 211
	SolarSatellite               ID = 212
	Destroyer                    ID = 213
	Deathstar                    ID = 214
	Battlecruiser                ID = 215
	EspionageTechnology          ID = 106 // Research
	ComputerTechnology           ID = 108
	WeaponsTechnology            ID = 109
	ShieldingTechnology          ID = 110
	ArmourTechnology             ID = 111
	EnergyTechnology             ID = 113
	HyperspaceTechnology         ID = 114
	CombustionDrive              ID = 115
	ImpulseDrive                 ID = 117
	HyperspaceDrive              ID = 118
	LaserTechnology              ID = 120
	IonTechnology                ID = 121
	PlasmaTechnology             ID = 122
	IntergalacticResearchNetwork ID = 123
	Astrophysics                 ID = 124
	GravitonTechnology           ID = 199

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
