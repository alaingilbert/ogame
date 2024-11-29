package ogame

// Slots ...
type Slots struct {
	InUse    int64
	Total    int64
	ExpInUse int64
	ExpTotal int64
}

func (s Slots) IsAllSlotsInUse(mission MissionID) bool {
	return (s.InUse == s.Total) ||
		(mission == Expedition && s.ExpInUse == s.ExpTotal)
}
