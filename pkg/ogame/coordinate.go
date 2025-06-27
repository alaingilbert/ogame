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
	Type     CelestialType
}

// NewCoordinate creates a new Coordinate
func NewCoordinate(galaxy, system, position int64, typ CelestialType) Coordinate {
	return Coordinate{Galaxy: galaxy, System: system, Position: position, Type: typ}
}

// NewPlanetCoordinate creates a new planet Coordinate
func NewPlanetCoordinate(galaxy, system, position int64) Coordinate {
	return Coordinate{Galaxy: galaxy, System: system, Position: position, Type: PlanetType}
}

// NewDebrisCoordinate creates a new debris Coordinate
func NewDebrisCoordinate(galaxy, system, position int64) Coordinate {
	return Coordinate{Galaxy: galaxy, System: system, Position: position, Type: DebrisType}
}

// NewMoonCoordinate creates a new moon Coordinate
func NewMoonCoordinate(galaxy, system, position int64) Coordinate {
	return Coordinate{Galaxy: galaxy, System: system, Position: position, Type: MoonType}
}

// String ...
func (c Coordinate) String() string {
	return fmt.Sprintf("[%c:%d:%d:%d]", strings.ToUpper(c.Type.String())[0], c.Galaxy, c.System, c.Position)
}

// IsZero reports whether c represents the zero coordinate,
func (c Coordinate) IsZero() bool {
	return c.Galaxy == 0 && c.System == 0 && c.Position == 0 && c.Type == 0
}

// Equal returns either two coordinates are equal or not
func (c Coordinate) Equal(v Coordinate) bool {
	return c.Cmp(v) == 0
}

// Cmp returns -1 if c < v, 0 if c == v, 1 if c > v
// Planet type is considered < Moon type
func (c Coordinate) Cmp(v Coordinate) int {
	switch {
	case c.Galaxy < v.Galaxy:
		return -1
	case c.Galaxy > v.Galaxy:
		return 1
	case c.System < v.System:
		return -1
	case c.System > v.System:
		return 1
	case c.Position < v.Position:
		return -1
	case c.Position > v.Position:
		return 1
	case c.Type < v.Type:
		return -1
	case c.Type > v.Type:
		return 1
	default:
		return 0
	}
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
