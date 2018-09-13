package ogame

import (
	"math"
	"time"
)

// BaseDefense ...
type BaseDefense struct {
	BaseDefender
	Price Resources
}

// GetPrice ...
func (b BaseDefense) GetPrice(nbr int) Resources {
	return b.Price.Mul(nbr)
}

// ConstructionTime ...
func (b BaseDefense) ConstructionTime(nbr, universeSpeed int, facilities Facilities) time.Duration {
	shipyardLvl := float64(facilities.Shipyard)
	naniteLvl := float64(facilities.NaniteFactory)
	hours := float64(b.StructuralIntegrity) / (2500 * (1 + shipyardLvl) * float64(universeSpeed) * math.Pow(2, naniteLvl))
	secs := hours * 3600
	return time.Duration(int(math.Floor(secs))*nbr) * time.Second
}
