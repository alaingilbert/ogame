package ogame

import (
	"fmt"
	"strings"
)

// Coordinate represent an ogame coordinate
type Coordinate struct {
	Galaxy   int64
	System   int64
	Position int64
	Type     CelestialType `gorm:"column:planet_type"`
}

func (c Coordinate) String() string {
	return fmt.Sprintf("[%c:%d:%d:%d]", strings.ToUpper(c.Type.String())[0], c.Galaxy, c.System, c.Position)
}

// Equal returns either two coordinates are equal or not
func (c Coordinate) Equal(v Coordinate) bool {
	return c.Galaxy == v.Galaxy &&
		c.System == v.System &&
		c.Position == v.Position &&
		c.Type == v.Type
}

// IsPlanet return true if coordinate is a planet
func (c Coordinate) IsPlanet() bool {
	return c.Type == PlanetType
}

// IsMoon return true if coordinate is a moon
func (c Coordinate) IsMoon() bool {
	return c.Type == MoonType
}

// IsDebris return true if coordinate is a debris field
func (c Coordinate) IsDebris() bool {
	return c.Type == DebrisType
}

// Planet return a new coordinate with planet type
func (c Coordinate) Planet() Coordinate {
	return Coordinate{Galaxy: c.Galaxy, System: c.System, Position: c.Position, Type: PlanetType}
}

// Moon return a new coordinate with moon type
func (c Coordinate) Moon() Coordinate {
	return Coordinate{Galaxy: c.Galaxy, System: c.System, Position: c.Position, Type: MoonType}
}

// Debris return a new coordinate with debris type
func (c Coordinate) Debris() Coordinate {
	return Coordinate{Galaxy: c.Galaxy, System: c.System, Position: c.Position, Type: DebrisType}
}
