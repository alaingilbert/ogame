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
		return SmallCargo.ID, nil
	case "large cargo":
		return LargeCargo.ID, nil
	case "light fighter":
		return LightFighter.ID, nil
	case "heavy fighter":
		return HeavyFighter.ID, nil
	case "cruiser":
		return Cruiser.ID, nil
	case "battleship":
		return Battleship.ID, nil
	case "colony ship":
		return ColonyShip.ID, nil
	case "recycler":
		return Recycler.ID, nil
	case "espionage probe":
		return EspionageProbe.ID, nil
	case "bomber":
		return Bomber.ID, nil
	case "solar satellite":
		return SolarSatellite.ID, nil
	case "destroyer":
		return Destroyer.ID, nil
	case "deathstar":
		return Deathstar.ID, nil
	case "battlecruiser":
		return Battlecruiser.ID, nil

	case "petit transporteur":
		return SmallCargo.ID, nil
	case "grand transporteur":
		return LargeCargo.ID, nil
	case "chasseur léger":
		return LightFighter.ID, nil
	case "chasseur lourd":
		return HeavyFighter.ID, nil
	case "croiseur":
		return Cruiser.ID, nil
	case "vaisseau de bataille":
		return Battleship.ID, nil
	case "vaisseau de colonisation":
		return ColonyShip.ID, nil
	case "recycleur":
		return Recycler.ID, nil
	case "sonde d`espionnage":
		return EspionageProbe.ID, nil
	case "bombardier":
		return Bomber.ID, nil
	case "satellite solaire":
		return SolarSatellite.ID, nil
	case "destructeur":
		return Destroyer.ID, nil
	case "étoile de la mort":
		return Deathstar.ID, nil
	case "traqueur":
		return Battlecruiser.ID, nil
	}
	return 0, errors.New("unable to parse ship " + name)
}

// IsDefenseID ...
func IsDefenseID(id int) bool {
	ogameID := ID(id)
	return ogameID == RocketLauncher.ID ||
		ogameID == LightLaser.ID ||
		ogameID == HeavyLaser.ID ||
		ogameID == GaussCannon.ID ||
		ogameID == IonCannon.ID ||
		ogameID == PlasmaTurret.ID ||
		ogameID == SmallShieldDome.ID ||
		ogameID == LargeShieldDome.ID ||
		ogameID == AntiBallisticMissiles.ID ||
		ogameID == InterplanetaryMissiles.ID
}

// IsShipID ...
func IsShipID(id int) bool {
	ogameID := ID(id)
	return ogameID == SmallCargo.ID ||
		ogameID == LargeCargo.ID ||
		ogameID == LightFighter.ID ||
		ogameID == HeavyFighter.ID ||
		ogameID == Cruiser.ID ||
		ogameID == Battleship.ID ||
		ogameID == ColonyShip.ID ||
		ogameID == Recycler.ID ||
		ogameID == EspionageProbe.ID ||
		ogameID == Bomber.ID ||
		ogameID == SolarSatellite.ID ||
		ogameID == Destroyer.ID ||
		ogameID == Deathstar.ID ||
		ogameID == Battlecruiser.ID
}

// IsTechID ...
func IsTechID(id int) bool {
	ogameID := ID(id)
	return ogameID == EspionageTechnology.ID ||
		ogameID == ComputerTechnology.ID ||
		ogameID == WeaponsTechnology.ID ||
		ogameID == ShieldingTechnology.ID ||
		ogameID == ArmourTechnology.ID ||
		ogameID == EnergyTechnology.ID ||
		ogameID == HyperspaceTechnology.ID ||
		ogameID == CombustionDrive.ID ||
		ogameID == ImpulseDrive.ID ||
		ogameID == HyperspaceDrive.ID ||
		ogameID == LaserTechnology.ID ||
		ogameID == IonTechnology.ID ||
		ogameID == PlasmaTechnology.ID ||
		ogameID == IntergalacticResearchNetwork.ID ||
		ogameID == Astrophysics.ID ||
		ogameID == GravitonTechnology.ID
}

// IsBuildingID ...
func IsBuildingID(id int) bool {
	return IsResourceBuildingID(id) || IsFacilityID(id)
}

// IsResourceBuildingID ...
func IsResourceBuildingID(id int) bool {
	ogameID := ID(id)
	return ogameID == MetalMine.ID ||
		ogameID == CrystalMine.ID ||
		ogameID == DeuteriumSynthesizer.ID ||
		ogameID == SolarPlant.ID ||
		ogameID == FusionReactor.ID ||
		ogameID == MetalStorage.ID ||
		ogameID == CrystalStorage.ID ||
		ogameID == DeuteriumTank.ID ||
		ogameID == ShieldedMetalDen.ID ||
		ogameID == UndergroundCrystalDen.ID ||
		ogameID == SeabedDeuteriumDen.ID
}

// IsFacilityID ...
func IsFacilityID(id int) bool {
	ogameID := ID(id)
	return ogameID == AllianceDepot.ID ||
		ogameID == RoboticsFactory.ID ||
		ogameID == Shipyard.ID ||
		ogameID == ResearchLab.ID ||
		ogameID == MissileSilo.ID ||
		ogameID == NaniteFactory.ID ||
		ogameID == Terraformer.ID ||
		ogameID == SpaceDock.ID
}
