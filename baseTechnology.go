package ogame

import (
	"math"
	"time"
)

// BaseTechnology ...
type BaseTechnology struct {
	Base
	BaseCost       Resources
	IncreaseFactor float64
}

// GetPrice ...
func (b BaseTechnology) GetPrice(level int) Resources {
	return Resources{
		Metal:     researchCost(b.BaseCost.Metal, b.IncreaseFactor, level),
		Crystal:   researchCost(b.BaseCost.Crystal, b.IncreaseFactor, level),
		Deuterium: researchCost(b.BaseCost.Deuterium, b.IncreaseFactor, level),
		Energy:    researchCost(b.BaseCost.Energy, b.IncreaseFactor, level),
	}
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

func researchCost(baseCost int, increaseFactor float64, level int) int {
	return int(float64(baseCost) * math.Pow(increaseFactor, float64(level-1)))
}
