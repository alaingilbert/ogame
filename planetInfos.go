package ogame

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
		ID   int
		Name string
		Rank int
	}
	Alliance *AllianceInfos
}
