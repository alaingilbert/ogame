package ogame

import "encoding/json"

// SystemInfos planets information for a specific system
type SystemInfos struct {
	galaxy  int
	system  int
	planets [15]*PlanetInfos
}

// Galaxy returns galaxy info
func (s SystemInfos) Galaxy() int {
	return s.galaxy
}

// System returns system info
func (s SystemInfos) System() int {
	return s.system
}

// Position returns planet at position idx in the SystemInfos
func (s SystemInfos) Position(idx int) *PlanetInfos {
	if idx < 1 || idx > 15 {
		return nil
	}
	return s.planets[idx-1]
}

// Each will execute provided callback for every positions in the system
func (s SystemInfos) Each(clb func(planetInfo *PlanetInfos)) {
	for i := 1; i <= 15; i++ {
		clb(s.Position(i))
	}
}

// MarshalJSON export private fields to json for ogamed
func (s SystemInfos) MarshalJSON() ([]byte, error) {
	var tmp struct {
		Galaxy  int
		System  int
		Planets [15]*PlanetInfos
	}
	tmp.Galaxy = s.galaxy
	tmp.System = s.system
	tmp.Planets = s.planets
	return json.Marshal(tmp)
}

// MoonInfos public information of a moon in the galaxy page
type MoonInfos struct {
	ID       int
	Diameter int
	Activity int
}

// AllianceInfos public information of an alliance in the galaxy page
type AllianceInfos struct {
	ID     int
	Name   string
	Rank   int
	Member int
}

// PlanetInfos public information of a planet in the galaxy page
type PlanetInfos struct {
	ID              int
	Activity        int
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
		Metal           int
		Crystal         int
		RecyclersNeeded int
	}
	Moon   *MoonInfos
	Player struct {
		ID         int
		Name       string
		Rank       int
		IsBandit   bool
		IsStarlord bool
	}
	Alliance *AllianceInfos
}
