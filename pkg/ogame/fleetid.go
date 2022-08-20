package ogame

import (
	"strconv"
)

// FleetID represent a fleet id
type FleetID int64

func (f FleetID) String() string {
	return strconv.FormatInt(int64(f), 10)
}
