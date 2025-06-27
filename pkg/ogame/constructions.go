package ogame

import "time"

type Constructions struct {
	Building   Construction
	Research   Construction
	LfBuilding Construction
	LfResearch Construction
}

type Construction struct {
	ID        ID
	Countdown time.Duration
	Level     int64
}
