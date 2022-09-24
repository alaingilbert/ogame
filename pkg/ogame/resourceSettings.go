package ogame

import "github.com/alaingilbert/ogame/pkg/utils"

// ResourceSettings represent a planet resource settings
type ResourceSettings struct {
	MetalMine            int64
	CrystalMine          int64
	DeuteriumSynthesizer int64
	SolarPlant           int64
	FusionReactor        int64
	SolarSatellite       int64
	Crawler              int64
}

func (r ResourceSettings) String() string {
	return "\n" +
		"           Metal Mine: " + utils.FI64(r.MetalMine) + "\n" +
		"         Crystal Mine: " + utils.FI64(r.CrystalMine) + "\n" +
		"Deuterium Synthesizer: " + utils.FI64(r.DeuteriumSynthesizer) + "\n" +
		"          Solar Plant: " + utils.FI64(r.SolarPlant) + "\n" +
		"       Fusion Reactor: " + utils.FI64(r.FusionReactor) + "\n" +
		"      Solar Satellite: " + utils.FI64(r.SolarSatellite) + "\n" +
		"              Crawler: " + utils.FI64(r.Crawler) + "\n"
}
