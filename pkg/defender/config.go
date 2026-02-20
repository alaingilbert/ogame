package defender

import (
	"github.com/alaingilbert/ogame/pkg/ogame"
)

// CelestialOrigin defines which celestial to check for attacks.
type CelestialOrigin struct {
	// Random when true, a random celestial is chosen for each check cycle.
	Random bool
	// CelestialID when Random is false, use this specific celestial.
	CelestialID ogame.CelestialID
}

// Config holds all configurable parameters for the defender bot.
type Config struct {
	// Enabled activates or deactivates the defender.
	Enabled bool

	// CheckIntervalSec is the interval between attack checks in seconds (e.g. 300–900).
	CheckIntervalSec int

	// EvacSpeed is the fleet speed to use when evacuating on destroy attacks.
	// Valid values: 1–10 (representing 10%–100%).
	EvacSpeed ogame.Speed

	// EvacSecondsBefore is how many seconds before attack arrival to send the fleet (e.g. 20–30).
	EvacSecondsBefore int

	// RecallSecondsAfter is how many seconds after the attack was supposed to arrive
	// to recall the evacuated fleet. 0 means no recall.
	// The fleet is only recalled if the attacking fleet is no longer present.
	RecallSecondsAfter int

	// DumpResourcesToDefense is reserved for future implementation.
	// When true, remaining resources that do not fit in fleet cargo will be invested in defenses.
	DumpResourcesToDefense bool

	// DumpResourcesToDefenseMoons also applies the dump logic to moons.
	DumpResourcesToDefenseMoons bool

	// Origin determines which celestial to use for attack checks.
	Origin CelestialOrigin
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Enabled:            false,
		CheckIntervalSec:   300,
		EvacSpeed:          ogame.HundredPercent,
		EvacSecondsBefore:  30,
		RecallSecondsAfter: 0,
		Origin:             CelestialOrigin{Random: true},
	}
}

// Validate checks the config values are within acceptable ranges and adjusts them if needed.
func (c *Config) Validate() {
	if c.CheckIntervalSec < 60 {
		c.CheckIntervalSec = 60
	}
	if c.EvacSecondsBefore < 5 {
		c.EvacSecondsBefore = 5
	}
	if c.RecallSecondsAfter < 0 {
		c.RecallSecondsAfter = 0
	}
}
