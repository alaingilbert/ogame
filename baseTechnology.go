package ogame

import (
	"math"
	"time"
)

// BaseTechnology base struct for technologies
type BaseTechnology struct {
	BaseLevelable
}

// ConstructionTime returns the duration it takes to build given technology
func (b BaseTechnology) ConstructionTime(level, universeSpeed int64, facilities Facilities, hasTechnocrat, isDiscoverer bool) time.Duration {
	price := b.GetPrice(int64(level))
	metalCost := float64(price.Metal)
	crystalCost := float64(price.Crystal)
	researchLabLvl := float64(facilities.ResearchLab)
	hours := (metalCost + crystalCost) / (1000 * (1 + researchLabLvl) * float64(universeSpeed))
	twentyFivePct := 0.25 * hours
	if hasTechnocrat {
		hours -= twentyFivePct
	}
	if isDiscoverer {
		hours -= twentyFivePct
	}
	secs := math.Max(1, hours*3600)
	return time.Duration(int64(math.Floor(secs))) * time.Second
}

// GetLevel returns current level of a technology
func (b BaseTechnology) GetLevel(_ LazyResourcesBuildings, _ LazyFacilities, lazyResearches LazyResearches) int64 {
	return lazyResearches().ByID(b.ID)
}
