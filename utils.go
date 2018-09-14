package ogame

import (
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
