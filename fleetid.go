package ogame

import "strconv"

// FleetID ...
type FleetID int

func (f FleetID) String() string {
	return strconv.Itoa(int(f))
}
