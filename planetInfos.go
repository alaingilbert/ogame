package ogame

// PlanetInfos ...
type PlanetInfos struct {
	Activity        int
	Name            string
	Img             string
	Coordinate      Coordinate
	Administrator   bool
	Inactive        bool
	Vacation        bool
	StrongPlayer    bool
	HonorableTarget bool
	Debris          struct {
		Metal           int
		Crystal         int
		RecyclersNeeded int
	}
	Player struct {
		ID   int
		Name string
		Rank int
	}
	Alliance struct {
		ID     int
		Name   string
		Rank   int
		Member int
	}
}
