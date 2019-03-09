package ogame

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// ParseInt ...
func ParseInt(val string) int {
	val = strings.Replace(val, ".", "", -1)
	val = strings.Replace(val, ",", "", -1)
	val = strings.TrimSpace(val)
	res, _ := strconv.Atoi(val)
	return res
}

func toInt(buf []byte) (n int) {
	for _, v := range buf {
		n = n*10 + int(v-'0')
	}
	return
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

// ParseCoord parse a coordinate from a string
func ParseCoord(str string) (coord Coordinate, err error) {
	m := regexp.MustCompile(`^\[?(([PMD]):)?(\d{1,3}):(\d{1,3}):(\d{1,3})]?$`).FindStringSubmatch(str)
	if len(m) == 5 {
		galaxy, _ := strconv.Atoi(m[2])
		system, _ := strconv.Atoi(m[3])
		position, _ := strconv.Atoi(m[4])
		planetType := PlanetType
		return Coordinate{galaxy, system, position, planetType}, nil
	} else if len(m) == 6 {
		planetTypeStr := m[2]
		galaxy, _ := strconv.Atoi(m[3])
		system, _ := strconv.Atoi(m[4])
		position, _ := strconv.Atoi(m[5])
		planetType := PlanetType
		if planetTypeStr == "M" {
			planetType = MoonType
		} else if planetTypeStr == "D" {
			planetType = DebrisType
		}
		return Coordinate{galaxy, system, position, planetType}, nil
	}
	return coord, errors.New("unable to parse coordinate")
}
