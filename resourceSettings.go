package ogame

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
		"           Metal Mine: " + FI64(r.MetalMine) + "\n" +
		"         Crystal Mine: " + FI64(r.CrystalMine) + "\n" +
		"Deuterium Synthesizer: " + FI64(r.DeuteriumSynthesizer) + "\n" +
		"          Solar Plant: " + FI64(r.SolarPlant) + "\n" +
		"       Fusion Reactor: " + FI64(r.FusionReactor) + "\n" +
		"      Solar Satellite: " + FI64(r.SolarSatellite) + "\n" +
		"              Crawler: " + FI64(r.Crawler)
}
