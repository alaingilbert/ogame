package ogame

import (
	"math"
)

// MetalMineProduction ...
func MetalMineProduction(level int) float64 {
	return 30 * float64(level) * math.Pow(1.1, float64(level))
}

// CrystalMineProduction ...
func CrystalMineProduction(level int) float64 {
	return 20 * float64(level) * math.Pow(1.1, float64(level))
}

// DeuteriumSynthesizerProduction ...
func DeuteriumSynthesizerProduction(level int, avgTemp int) float64 {
	return 10 * float64(level) * math.Pow(1.1, float64(level)) * (-0.002*float64(avgTemp) + 1.28)
}

// SolarPlantProduction ...
func SolarPlantProduction(level int) float64 {
	return 20 * float64(level) * math.Pow(1.1, float64(level))
}

// FusionReactorProduction ...
func FusionReactorProduction(level int, energyTechno int) float64 {
	return 30 * float64(level) * math.Pow(1.05+float64(energyTechno)*0.01, float64(level))
}

// SolarSatellitesProduction ...
func SolarSatellitesProduction(avgTemp int) float64 {
	return float64(avgTemp)/4 + 20
}

func buildingCost(baseCost int, increaseFactor float64, level int) int {
	return int(math.Floor(float64(baseCost) * math.Pow(increaseFactor, float64(level)-1)))
}

// MetalMineCost ...
func MetalMineCost(level int) Resources {
	return Resources{
		Metal:   buildingCost(60, 1.5, level),
		Crystal: buildingCost(15, 1.5, level),
	}
}
