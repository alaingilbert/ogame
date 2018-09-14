package ogame

import "strconv"

// FleetID represent a fleet id
type FleetID int

func (f FleetID) String() string {
	return strconv.Itoa(int(f))
}
