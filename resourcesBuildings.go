package ogame

import "strconv"

// LazyResourcesBuildings ...
type LazyResourcesBuildings func() ResourcesBuildings

// ResourcesBuildings represent a planet resource buildings
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

// Lazy returns a function that return self
func (r ResourcesBuildings) Lazy() LazyResourcesBuildings {
	return func() ResourcesBuildings { return r }
}

// ByID gets the resource building level from a building id
func (r ResourcesBuildings) ByID(id ID) int {
	if id == MetalMine.ID {
		return r.MetalMine
	} else if id == CrystalMine.ID {
		return r.CrystalMine
	} else if id == DeuteriumSynthesizer.ID {
		return r.DeuteriumSynthesizer
	} else if id == SolarPlant.ID {
		return r.SolarPlant
	} else if id == FusionReactor.ID {
		return r.FusionReactor
	} else if id == SolarSatellite.ID {
		return r.SolarSatellite
	} else if id == MetalStorage.ID {
		return r.MetalStorage
	} else if id == CrystalStorage.ID {
		return r.CrystalStorage
	} else if id == DeuteriumTank.ID {
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
