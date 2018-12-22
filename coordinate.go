package ogame

import (
	"fmt"
	"strings"
)

// Coordinate represent an ogame coordinate
type Coordinate struct {
	Galaxy   int
	System   int
	Position int
	Type     CelestialType
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
