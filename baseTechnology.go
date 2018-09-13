package ogame

import (
	"math"
	"time"
)

// BaseTechnology ...
type BaseTechnology struct {
	BaseLevelable
}

// ConstructionTime ...
func (b BaseTechnology) ConstructionTime(level, universeSpeed int, facilities Facilities) time.Duration {
	price := b.GetPrice(level)
	metalCost := float64(price.Metal)
	crystalCost := float64(price.Crystal)
	researchLabLvl := float64(facilities.ResearchLab)
	hours := (metalCost + crystalCost) / (1000 * (1 + researchLabLvl) * float64(universeSpeed))
	secs := hours * 3600
	return time.Duration(int(math.Floor(secs))) * time.Second
}

// GetLevel ...
func (b BaseTechnology) GetLevel(resourcesBuildings ResourcesBuildings, facilities Facilities, researches Researches) int {
	return researches.ByID(b.ID)
}
