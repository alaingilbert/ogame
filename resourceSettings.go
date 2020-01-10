package ogame

import "strconv"

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
		"           Metal Mine: " + strconv.FormatInt(r.MetalMine, 10) + "\n" +
		"         Crystal Mine: " + strconv.FormatInt(r.CrystalMine, 10) + "\n" +
		"Deuterium Synthesizer: " + strconv.FormatInt(r.DeuteriumSynthesizer, 10) + "\n" +
		"          Solar Plant: " + strconv.FormatInt(r.SolarPlant, 10) + "\n" +
		"       Fusion Reactor: " + strconv.FormatInt(r.FusionReactor, 10) + "\n" +
		"      Solar Satellite: " + strconv.FormatInt(r.SolarSatellite, 10) + "\n" +
		"              Crawler: " + strconv.FormatInt(r.Crawler, 10)
}
