package ogame

import "strconv"

// ResourceSettings represent a planet resource settings
type ResourceSettings struct {
	MetalMine            int
	CrystalMine          int
	DeuteriumSynthesizer int
	SolarPlant           int
	FusionReactor        int
	SolarSatellite       int
	Crawler              int
}

func (r ResourceSettings) String() string {
	return "\n" +
		"           Metal Mine: " + strconv.Itoa(r.MetalMine) + "\n" +
		"         Crystal Mine: " + strconv.Itoa(r.CrystalMine) + "\n" +
		"Deuterium Synthesizer: " + strconv.Itoa(r.DeuteriumSynthesizer) + "\n" +
		"          Solar Plant: " + strconv.Itoa(r.SolarPlant) + "\n" +
		"       Fusion Reactor: " + strconv.Itoa(r.FusionReactor) + "\n" +
		"      Solar Satellite: " + strconv.Itoa(r.SolarSatellite) + "\n" +
		"              Crawler: " + strconv.Itoa(r.Crawler)
}
