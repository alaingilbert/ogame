package wrapper

import (
	"github.com/alaingilbert/ogame/pkg/ogame"
)

type Options struct {
	DebugGalaxy     bool
	SkipInterceptor bool
	SkipRetry       bool
	ChangePlanet    ogame.CelestialID // cp parameter
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
