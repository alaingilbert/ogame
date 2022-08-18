package ogame

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsShipID(t *testing.T) {
	assert.True(t, IsShipID(int64(SmallCargoID)))
	assert.False(t, IsShipID(int64(RocketLauncherID)))
}

func TestIsDefenseID(t *testing.T) {
	assert.True(t, IsDefenseID(int64(RocketLauncherID)))
	assert.False(t, IsDefenseID(int64(SmallCargoID)))
}

func TestIsTechID(t *testing.T) {
	assert.True(t, IsTechID(int64(CombustionDriveID)))
	assert.False(t, IsTechID(int64(SmallCargoID)))
}

func TestIsBuildingID(t *testing.T) {
	assert.True(t, IsBuildingID(int64(MetalMineID)))
	assert.True(t, IsBuildingID(int64(RoboticsFactoryID)))
	assert.False(t, IsBuildingID(int64(SmallCargoID)))
}

func TestIsResourceBuildingID(t *testing.T) {
	assert.True(t, IsResourceBuildingID(int64(MetalMineID)))
	assert.False(t, IsResourceBuildingID(int64(RoboticsFactoryID)))
	assert.False(t, IsResourceBuildingID(int64(SmallCargoID)))
}

func TestIsFacilityID(t *testing.T) {
	assert.False(t, IsFacilityID(int64(MetalMineID)))
	assert.True(t, IsFacilityID(int64(RoboticsFactoryID)))
	assert.False(t, IsFacilityID(int64(SmallCargoID)))
}

func TestIsValidIPMTarget(t *testing.T) {
	assert.True(t, ID(0).IsValidIPMTarget())
	assert.True(t, RocketLauncherID.IsValidIPMTarget())
	assert.False(t, AntiBallisticMissilesID.IsValidIPMTarget())
	assert.False(t, InterplanetaryMissilesID.IsValidIPMTarget())
	assert.False(t, MetalMineID.IsValidIPMTarget())
}

func TestID_IsSet(t *testing.T) {
	assert.True(t, AllianceDepotID.IsSet())
	assert.False(t, ID(0).IsSet())
}

func TestID_Int(t *testing.T) {
	assert.Equal(t, int64(34), AllianceDepotID.Int64())
}

func TestID_String(t *testing.T) {
	assert.Equal(t, "Invalid(123456)", ID(123456).String())
	assert.Equal(t, "AllianceDepot", AllianceDepotID.String())
	assert.Equal(t, "RoboticsFactory", RoboticsFactoryID.String())
	assert.Equal(t, "Shipyard", ShipyardID.String())
	assert.Equal(t, "ResearchLab", ResearchLabID.String())
	assert.Equal(t, "MissileSilo", MissileSiloID.String())
	assert.Equal(t, "NaniteFactory", NaniteFactoryID.String())
	assert.Equal(t, "Terraformer", TerraformerID.String())
	assert.Equal(t, "SpaceDock", SpaceDockID.String())
	assert.Equal(t, "LunarBase", LunarBaseID.String())
	assert.Equal(t, "SensorPhalanx", SensorPhalanxID.String())
	assert.Equal(t, "JumpGate", JumpGateID.String())
	assert.Equal(t, "MetalMine", MetalMineID.String())
	assert.Equal(t, "CrystalMine", CrystalMineID.String())
	assert.Equal(t, "DeuteriumSynthesizer", DeuteriumSynthesizerID.String())
	assert.Equal(t, "SolarPlant", SolarPlantID.String())
	assert.Equal(t, "FusionReactor", FusionReactorID.String())
	assert.Equal(t, "MetalStorage", MetalStorageID.String())
	assert.Equal(t, "CrystalStorage", CrystalStorageID.String())
	assert.Equal(t, "DeuteriumTank", DeuteriumTankID.String())
	assert.Equal(t, "ShieldedMetalDen", ShieldedMetalDenID.String())
	assert.Equal(t, "UndergroundCrystalDen", UndergroundCrystalDenID.String())
	assert.Equal(t, "SeabedDeuteriumDen", SeabedDeuteriumDenID.String())
	assert.Equal(t, "RocketLauncher", RocketLauncherID.String())
	assert.Equal(t, "LightLaser", LightLaserID.String())
	assert.Equal(t, "HeavyLaser", HeavyLaserID.String())
	assert.Equal(t, "GaussCannon", GaussCannonID.String())
	assert.Equal(t, "IonCannon", IonCannonID.String())
	assert.Equal(t, "PlasmaTurret", PlasmaTurretID.String())
	assert.Equal(t, "SmallShieldDome", SmallShieldDomeID.String())
	assert.Equal(t, "LargeShieldDome", LargeShieldDomeID.String())
	assert.Equal(t, "AntiBallisticMissiles", AntiBallisticMissilesID.String())
	assert.Equal(t, "InterplanetaryMissiles", InterplanetaryMissilesID.String())
	assert.Equal(t, "SmallCargo", SmallCargoID.String())
	assert.Equal(t, "LargeCargo", LargeCargoID.String())
	assert.Equal(t, "LightFighter", LightFighterID.String())
	assert.Equal(t, "HeavyFighter", HeavyFighterID.String())
	assert.Equal(t, "Cruiser", CruiserID.String())
	assert.Equal(t, "Battleship", BattleshipID.String())
	assert.Equal(t, "ColonyShip", ColonyShipID.String())
	assert.Equal(t, "Recycler", RecyclerID.String())
	assert.Equal(t, "EspionageProbe", EspionageProbeID.String())
	assert.Equal(t, "Bomber", BomberID.String())
	assert.Equal(t, "SolarSatellite", SolarSatelliteID.String())
	assert.Equal(t, "Destroyer", DestroyerID.String())
	assert.Equal(t, "Deathstar", DeathstarID.String())
	assert.Equal(t, "Battlecruiser", BattlecruiserID.String())
	assert.Equal(t, "Crawler", CrawlerID.String())
	assert.Equal(t, "Reaper", ReaperID.String())
	assert.Equal(t, "Pathfinder", PathfinderID.String())
	assert.Equal(t, "EspionageTechnology", EspionageTechnologyID.String())
	assert.Equal(t, "ComputerTechnology", ComputerTechnologyID.String())
	assert.Equal(t, "WeaponsTechnology", WeaponsTechnologyID.String())
	assert.Equal(t, "ShieldingTechnology", ShieldingTechnologyID.String())
	assert.Equal(t, "ArmourTechnology", ArmourTechnologyID.String())
	assert.Equal(t, "EnergyTechnology", EnergyTechnologyID.String())
	assert.Equal(t, "HyperspaceTechnology", HyperspaceTechnologyID.String())
	assert.Equal(t, "CombustionDrive", CombustionDriveID.String())
	assert.Equal(t, "ImpulseDrive", ImpulseDriveID.String())
	assert.Equal(t, "HyperspaceDrive", HyperspaceDriveID.String())
	assert.Equal(t, "LaserTechnology", LaserTechnologyID.String())
	assert.Equal(t, "IonTechnology", IonTechnologyID.String())
	assert.Equal(t, "PlasmaTechnology", PlasmaTechnologyID.String())
	assert.Equal(t, "IntergalacticResearchNetwork", IntergalacticResearchNetworkID.String())
	assert.Equal(t, "Astrophysics", AstrophysicsID.String())
	assert.Equal(t, "GravitonTechnology", GravitonTechnologyID.String())
}
