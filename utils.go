package ogame

import (
	"errors"
	"strconv"
	"strings"
)

func parseInt(val string) int {
	val = strings.Replace(val, ".", "", -1)
	val = strings.Replace(val, ",", "", -1)
	val = strings.Trim(val, " \t\r\n")
	res, _ := strconv.Atoi(val)
	return res
}

func parseShip(name string) (ID, error) {
	name = strings.ToLower(name)
	switch name {
	case "small cargo":
		return SmallCargo, nil
	case "large cargo":
		return LargeCargo, nil
	case "light fighter":
		return LightFighter, nil
	case "heavy fighter":
		return HeavyFighter, nil
	case "cruiser":
		return Cruiser, nil
	case "battleship":
		return Battleship, nil
	case "colony ship":
		return ColonyShip, nil
	case "recycler":
		return Recycler, nil
	case "espionage probe":
		return EspionageProbe, nil
	case "bomber":
		return Bomber, nil
	case "solar satellite":
		return SolarSatellite, nil
	case "destroyer":
		return Destroyer, nil
	case "deathstar":
		return Deathstar, nil
	case "battlecruiser":
		return Battlecruiser, nil

	case "petit transporteur":
		return SmallCargo, nil
	case "grand transporteur":
		return LargeCargo, nil
	case "chasseur léger":
		return LightFighter, nil
	case "chasseur lourd":
		return HeavyFighter, nil
	case "croiseur":
		return Cruiser, nil
	case "vaisseau de bataille":
		return Battleship, nil
	case "vaisseau de colonisation":
		return ColonyShip, nil
	case "recycleur":
		return Recycler, nil
	case "sonde d`espionnage":
		return EspionageProbe, nil
	case "bombardier":
		return Bomber, nil
	case "satellite solaire":
		return SolarSatellite, nil
	case "destructeur":
		return Destroyer, nil
	case "étoile de la mort":
		return Deathstar, nil
	case "traqueur":
		return Battlecruiser, nil
	}
	return 0, errors.New("unable to parse ship " + name)
}

// IsDefenseID ...
func IsDefenseID(id int) bool {
	ogameID := ID(id)
	return ogameID == RocketLauncher ||
		ogameID == LightLaser ||
		ogameID == HeavyLaser ||
		ogameID == GaussCannon ||
		ogameID == IonCannon ||
		ogameID == PlasmaTurret ||
		ogameID == SmallShieldDome ||
		ogameID == LargeShieldDome ||
		ogameID == AntiBallisticMissiles ||
		ogameID == InterplanetaryMissiles
}

// IsShipID ...
func IsShipID(id int) bool {
	ogameID := ID(id)
	return ogameID == SmallCargo ||
		ogameID == LargeCargo ||
		ogameID == LightFighter ||
		ogameID == HeavyFighter ||
		ogameID == Cruiser ||
		ogameID == Battleship ||
		ogameID == ColonyShip ||
		ogameID == Recycler ||
		ogameID == EspionageProbe ||
		ogameID == Bomber ||
		ogameID == SolarSatellite ||
		ogameID == Destroyer ||
		ogameID == Deathstar ||
		ogameID == Battlecruiser
}

// IsTechID ...
func IsTechID(id int) bool {
	ogameID := ID(id)
	return ogameID == EspionageTechnology ||
		ogameID == ComputerTechnology ||
		ogameID == WeaponsTechnology ||
		ogameID == ShieldingTechnology ||
		ogameID == ArmourTechnology ||
		ogameID == EnergyTechnology ||
		ogameID == HyperspaceTechnology ||
		ogameID == CombustionDrive ||
		ogameID == ImpulseDrive ||
		ogameID == HyperspaceDrive ||
		ogameID == LaserTechnology ||
		ogameID == IonTechnology ||
		ogameID == PlasmaTechnology ||
		ogameID == IntergalacticResearchNetwork ||
		ogameID == Astrophysics ||
		ogameID == GravitonTechnology
}

// IsBuildingID ...
func IsBuildingID(id int) bool {
	return IsResourceBuildingID(id) || IsFacilityID(id)
}

// IsResourceBuildingID ...
func IsResourceBuildingID(id int) bool {
	ogameID := ID(id)
	return ogameID == MetalMine ||
		ogameID == CrystalMine ||
		ogameID == DeuteriumSynthesizer ||
		ogameID == SolarPlant ||
		ogameID == FusionReactor ||
		ogameID == MetalStorage ||
		ogameID == CrystalStorage ||
		ogameID == DeuteriumTank ||
		ogameID == ShieldedMetalDen ||
		ogameID == UndergroundCrystalDen ||
		ogameID == SeabedDeuteriumDen
}

// IsFacilityID ...
func IsFacilityID(id int) bool {
	ogameID := ID(id)
	return ogameID == AllianceDepot ||
		ogameID == RoboticsFactory ||
		ogameID == Shipyard ||
		ogameID == ResearchLab ||
		ogameID == MissileSilo ||
		ogameID == NaniteFactory ||
		ogameID == Terraformer ||
		ogameID == SpaceDock
}
