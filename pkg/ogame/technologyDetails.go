package ogame

import "time"

type TechnologyDetails struct {
	TechnologyID       ID
	ProductionDuration time.Duration
	Price              Resources
	Level              int64
	TearDownEnabled    bool
}
