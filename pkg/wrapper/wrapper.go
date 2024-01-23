package wrapper

import (
	"github.com/alaingilbert/ogame/pkg/ogame"
	"time"
)

// IntoCoordinate any types that can be turned into a Coordinate
type IntoCoordinate any

// IntoCelestial any types that can be turned into a Celestial
type IntoCelestial any

// IntoPlanet any types that can be turned into a Planet
type IntoPlanet any

// IntoMoon any types that can be turned into a Moon
type IntoMoon any

// Options ...
type Options struct {
	DebugGalaxy     bool
	SkipInterceptor bool
	SkipRetry       bool
	ChangePlanet    ogame.CelestialID // cp parameter
	Delay           time.Duration
}

// Option functions to be passed to public interface to change behaviors
type Option func(*Options)

// DebugGalaxy option to debug galaxy
func DebugGalaxy(opt *Options) {
	opt.DebugGalaxy = true
}

// SkipInterceptor option to skip html interceptors
func SkipInterceptor(opt *Options) {
	opt.SkipInterceptor = true
}

// SkipRetry option to skip retry
func SkipRetry(opt *Options) {
	opt.SkipRetry = true
}

// ChangePlanet set the cp parameter
func ChangePlanet(celestialID ogame.CelestialID) Option {
	return func(opt *Options) {
		opt.ChangePlanet = celestialID
	}
}

// Delay delays a page load; simulating slow request/response from ogame server
func Delay(dur time.Duration) Option {
	return func(opt *Options) {
		opt.Delay = dur
	}
}
