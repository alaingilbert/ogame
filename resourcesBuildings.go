package ogame

import "strconv"

// ResourcesBuildings ...
type ResourcesBuildings struct {
	MetalMine            int
	CrystalMine          int
	DeuteriumSynthesizer int
	SolarPlant           int
	FusionReactor        int
	SolarSatellite       int
	MetalStorage         int
	CrystalStorage       int
	DeuteriumTank        int
}

// ByOGameID ...
func (r ResourcesBuildings) ByOGameID(ogameID ID) int {
	if ogameID == MetalMine {
		return r.MetalMine
	} else if ogameID == CrystalMine {
		return r.CrystalMine
	} else if ogameID == DeuteriumSynthesizer {
		return r.DeuteriumSynthesizer
	} else if ogameID == SolarPlant {
		return r.SolarPlant
	} else if ogameID == FusionReactor {
		return r.FusionReactor
	} else if ogameID == SolarSatellite {
		return r.SolarSatellite
	} else if ogameID == MetalStorage {
		return r.MetalStorage
	} else if ogameID == CrystalStorage {
		return r.CrystalStorage
	} else if ogameID == DeuteriumTank {
		return r.DeuteriumTank
	}
	return 0
}

func (r ResourcesBuildings) String() string {
	return "\n" +
		"           Metal Mine: " + strconv.Itoa(r.MetalMine) + "\n" +
		"         Crystal Mine: " + strconv.Itoa(r.CrystalMine) + "\n" +
		"Deuterium Synthesizer: " + strconv.Itoa(r.DeuteriumSynthesizer) + "\n" +
		"          Solar Plant: " + strconv.Itoa(r.SolarPlant) + "\n" +
		"       Fusion Reactor: " + strconv.Itoa(r.FusionReactor) + "\n" +
		"      Solar Satellite: " + strconv.Itoa(r.SolarSatellite) + "\n" +
		"        Metal Storage: " + strconv.Itoa(r.MetalStorage) + "\n" +
		"      Crystal Storage: " + strconv.Itoa(r.CrystalStorage) + "\n" +
		"       Deuterium Tank: " + strconv.Itoa(r.DeuteriumTank)
}
