package ogame

// FleetID represent a fleet id
type FleetID int64

func (f FleetID) String() string {
	return FI64(f)
}
