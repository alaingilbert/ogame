package ogame

import (
	"strconv"
)

// FleetID represent a fleet id
type FleetID int64

// String ...
func (f FleetID) String() string {
	return strconv.FormatInt(int64(f), 10)
}

// IsSet ...
func (f FleetID) IsSet() bool {
	return f > 0
}
