package ogame

import "encoding/json"

// SystemInfos planets information for a specific system
type SystemInfos struct {
	galaxy           int64
	system           int64
	planets          [15]*PlanetInfos
	ExpeditionDebris struct {
		Metal             int64
		Crystal           int64
		PathfindersNeeded int64
	}
	Events struct {
		Darkmatter  int64
		HasAsteroid bool
	}
}

// Galaxy returns galaxy info
func (s SystemInfos) Galaxy() int64 {
	return s.galaxy
}

// System returns system info
func (s SystemInfos) System() int64 {
	return s.system
}

// Position returns planet at position idx in the SystemInfos
func (s SystemInfos) Position(idx int64) *PlanetInfos {
	if idx < 1 || idx > 15 {
		return nil
	}
	return s.planets[idx-1]
}

// Each will execute provided callback for every positions in the system
func (s SystemInfos) Each(clb func(planetInfo *PlanetInfos)) {
	var i int64
	for i = 1; i <= 15; i++ {
		clb(s.Position(i))
	}
}

// MarshalJSON export private fields to json for ogamed
func (s SystemInfos) MarshalJSON() ([]byte, error) {
	var tmp struct {
		Galaxy           int64
		System           int64
		Planets          [15]*PlanetInfos
		ExpeditionDebris struct {
			Metal             int64
			Crystal           int64
			PathfindersNeeded int64
		}
	}
	tmp.Galaxy = s.galaxy
	tmp.System = s.system
	tmp.Planets = s.planets
	tmp.ExpeditionDebris.Metal = s.ExpeditionDebris.Metal
	tmp.ExpeditionDebris.Crystal = s.ExpeditionDebris.Crystal
	tmp.ExpeditionDebris.PathfindersNeeded = s.ExpeditionDebris.PathfindersNeeded
	return json.Marshal(tmp)
}

// MoonInfos public information of a moon in the galaxy page
type MoonInfos struct {
	ID       int64
	Diameter int64
	Activity int64
}

// AllianceInfos public information of an alliance in the galaxy page
type AllianceInfos struct {
	ID     int64
	Name   string
	Rank   int64
	Member int64
}

// PlanetInfos public information of a planet in the galaxy page
type PlanetInfos struct {
	ID              int64
	Activity        int64 // no activity: 0, active: 15, inactive: [16, 59]
	Name            string
	Img             string
	Coordinate      Coordinate
	Administrator   bool
	Inactive        bool
	Vacation        bool
	StrongPlayer    bool
	Newbie          bool
	HonorableTarget bool
	Banned          bool
	Debris          struct {
		Metal           int64
		Crystal         int64
		RecyclersNeeded int64
	}
	Moon   *MoonInfos
	Player struct {
		ID         int64
		Name       string
		Rank       int64
		IsBandit   bool
		IsStarlord bool
	}
	Alliance *AllianceInfos
}
