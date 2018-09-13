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

// IsDefenseID helper returns if an integer is a defense id
func IsDefenseID(id int) bool {
	return ID(id).IsDefense()
}

// IsShipID helper returns if an integer is a ship id
func IsShipID(id int) bool {
	return ID(id).IsShip()
}

// IsTechID helper returns if an integer is a tech id
func IsTechID(id int) bool {
	return ID(id).IsTech()
}

// IsBuildingID helper returns if an integer is a building id
func IsBuildingID(id int) bool {
	return ID(id).IsBuilding()
}

// IsResourceBuildingID helper returns if an integer is a resource defense id
func IsResourceBuildingID(id int) bool {
	return ID(id).IsResourceBuilding()
}

// IsFacilityID helper returns if an integer is a facility id
func IsFacilityID(id int) bool {
	return ID(id).IsFacility()
}
