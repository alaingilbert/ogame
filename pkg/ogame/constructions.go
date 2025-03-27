package ogame

type Constructions struct {
	Building   Construction
	Research   Construction
	LfBuilding Construction
	LfResearch Construction
}

type Construction struct {
	ID        ID
	Countdown int64
}
