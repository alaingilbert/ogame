package ogame

import (
	"errors"
	"strconv"
	"strings"
)

func parseInt(val string) int {
	val = strings.Replace(val, ".", "", -1)
	val = strings.Replace(val, ",", "", -1)
	val = strings.TrimSpace(val)
	res, _ := strconv.Atoi(val)
	return res
}

func parseShip(name string) (ID, error) {
	name = strings.ToLower(name)
	switch name {
	case "small cargo":
		return SmallCargoID, nil
	case "large cargo":
		return LargeCargoID, nil
	case "light fighter":
		return LightFighterID, nil
	case "heavy fighter":
		return HeavyFighterID, nil
	case "cruiser":
		return CruiserID, nil
	case "battleship":
		return BattleshipID, nil
	case "colony ship":
		return ColonyShipID, nil
	case "recycler":
		return RecyclerID, nil
	case "espionage probe":
		return EspionageProbeID, nil
	case "bomber":
		return BomberID, nil
	case "solar satellite":
		return SolarSatelliteID, nil
	case "destroyer":
		return DestroyerID, nil
	case "deathstar":
		return DeathstarID, nil
	case "battlecruiser":
		return BattlecruiserID, nil

	case "petit transporteur":
		return SmallCargoID, nil
	case "grand transporteur":
		return LargeCargoID, nil
	case "chasseur léger":
		return LightFighterID, nil
	case "chasseur lourd":
		return HeavyFighterID, nil
	case "croiseur":
		return CruiserID, nil
	case "vaisseau de bataille":
		return BattleshipID, nil
	case "vaisseau de colonisation":
		return ColonyShipID, nil
	case "recycleur":
		return RecyclerID, nil
	case "sonde d`espionnage":
		return EspionageProbeID, nil
	case "bombardier":
		return BomberID, nil
	case "satellite solaire":
		return SolarSatelliteID, nil
	case "destructeur":
		return DestroyerID, nil
	case "étoile de la mort":
		return DeathstarID, nil
	case "traqueur":
		return BattlecruiserID, nil
	}
	return 0, errors.New("unable to parse ship " + name)
}

// IsDefenseID ...
func IsDefenseID(id int) bool {
	ogameID := ID(id)
	return ogameID == RocketLauncherID ||
		ogameID == LightLaserID ||
		ogameID == HeavyLaserID ||
		ogameID == GaussCannonID ||
		ogameID == IonCannonID ||
		ogameID == PlasmaTurretID ||
		ogameID == SmallShieldDomeID ||
		ogameID == LargeShieldDomeID ||
		ogameID == AntiBallisticMissilesID ||
		ogameID == InterplanetaryMissilesID
}

// IsShipID ...
func IsShipID(id int) bool {
	ogameID := ID(id)
	return ogameID == SmallCargoID ||
		ogameID == LargeCargoID ||
		ogameID == LightFighterID ||
		ogameID == HeavyFighterID ||
		ogameID == CruiserID ||
		ogameID == BattleshipID ||
		ogameID == ColonyShipID ||
		ogameID == RecyclerID ||
		ogameID == EspionageProbeID ||
		ogameID == BomberID ||
		ogameID == SolarSatelliteID ||
		ogameID == DestroyerID ||
		ogameID == DeathstarID ||
		ogameID == BattlecruiserID
}

// IsTechID ...
func IsTechID(id int) bool {
	ogameID := ID(id)
	return ogameID == EspionageTechnologyID ||
		ogameID == ComputerTechnologyID ||
		ogameID == WeaponsTechnologyID ||
		ogameID == ShieldingTechnologyID ||
		ogameID == ArmourTechnologyID ||
		ogameID == EnergyTechnologyID ||
		ogameID == HyperspaceTechnologyID ||
		ogameID == CombustionDriveID ||
		ogameID == ImpulseDriveID ||
		ogameID == HyperspaceDriveID ||
		ogameID == LaserTechnologyID ||
		ogameID == IonTechnologyID ||
		ogameID == PlasmaTechnologyID ||
		ogameID == IntergalacticResearchNetworkID ||
		ogameID == AstrophysicsID ||
		ogameID == GravitonTechnologyID
}

// IsBuildingID ...
func IsBuildingID(id int) bool {
	return IsResourceBuildingID(id) || IsFacilityID(id)
}

// IsResourceBuildingID ...
func IsResourceBuildingID(id int) bool {
	ogameID := ID(id)
	return ogameID == MetalMineID ||
		ogameID == CrystalMineID ||
		ogameID == DeuteriumSynthesizerID ||
		ogameID == SolarPlantID ||
		ogameID == FusionReactorID ||
		ogameID == MetalStorageID ||
		ogameID == CrystalStorageID ||
		ogameID == DeuteriumTankID ||
		ogameID == ShieldedMetalDenID ||
		ogameID == UndergroundCrystalDenID ||
		ogameID == SeabedDeuteriumDenID
}

// IsFacilityID ...
func IsFacilityID(id int) bool {
	ogameID := ID(id)
	return ogameID == AllianceDepotID ||
		ogameID == RoboticsFactoryID ||
		ogameID == ShipyardID ||
		ogameID == ResearchLabID ||
		ogameID == MissileSiloID ||
		ogameID == NaniteFactoryID ||
		ogameID == TerraformerID ||
		ogameID == SpaceDockID
}
