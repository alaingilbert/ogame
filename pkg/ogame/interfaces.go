package ogame

import (
	"time"
)

type Celestial interface {
	GetID() CelestialID
	GetName() string
	GetDiameter() int64
	GetFields() Fields
	GetCoordinate() Coordinate
	GetImg() string
}

// BaseOgameObj base interface for all ogame objects (buildings, technologies, ships, defenses)
type BaseOgameObj interface {
	ConstructionTime(nbr, universeSpeed int64, facilities Facilities, hasTechnocrat, isDiscoverer bool) time.Duration
	GetID() ID
	GetName() string
	GetPrice(int64) Resources
	GetRequirements() map[ID]int64
	IsAvailable(CelestialType, LazyResourcesBuildings, LazyFacilities, LazyResearches, int64, CharacterClass) bool
}

// DefenderObj base interface for all defensive units (ships, defenses)
type DefenderObj interface {
	BaseOgameObj
	GetRapidfireAgainst() map[ID]int64
	GetRapidfireFrom() map[ID]int64
	GetShieldPower(Researches) int64
	GetStructuralIntegrity(Researches) int64
	GetWeaponPower(Researches) int64
}

// Ship interface implemented by all ships units
type Ship interface {
	DefenderObj
	GetCargoCapacity(techs Researches, probeRaids, isCollector, isPioneers bool) int64
	GetFuelConsumption(techs Researches, fleetDeutSaveFactor float64, isGeneral bool) int64
	GetSpeed(techs Researches, isCollector, isGeneral bool) int64
}

// Defense interface implemented by all defenses units
type Defense interface {
	DefenderObj
}

// Levelable base interface for all levelable ogame objects (buildings, technologies)
type Levelable interface {
	BaseOgameObj
	GetLevel(LazyResourcesBuildings, LazyFacilities, LazyResearches) int64
}

// Technology interface that all technologies implement
type Technology interface {
	Levelable
}

// Building interface that all buildings implement
type Building interface {
	Levelable
	DeconstructionPrice(lvl int64, techs Researches) Resources
}
