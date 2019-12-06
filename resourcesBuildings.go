package ogame

import "strconv"

// LazyResourcesBuildings ...
type LazyResourcesBuildings func() ResourcesBuildings

// ResourcesBuildings represent a planet resource buildings
type ResourcesBuildings struct {
	MetalMine            int64
	CrystalMine          int64
	DeuteriumSynthesizer int64
	SolarPlant           int64
	FusionReactor        int64
	SolarSatellite       int64
	MetalStorage         int64
	CrystalStorage       int64
	DeuteriumTank        int64
}

// Lazy returns a function that return self
func (r ResourcesBuildings) Lazy() LazyResourcesBuildings {
	return func() ResourcesBuildings { return r }
}

// ByID gets the resource building level from a building id
func (r ResourcesBuildings) ByID(id ID) int64 {
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
		"           Metal Mine: " + strconv.FormatInt(r.MetalMine, 10) + "\n" +
		"         Crystal Mine: " + strconv.FormatInt(r.CrystalMine, 10) + "\n" +
		"Deuterium Synthesizer: " + strconv.FormatInt(r.DeuteriumSynthesizer, 10) + "\n" +
		"          Solar Plant: " + strconv.FormatInt(r.SolarPlant, 10) + "\n" +
		"       Fusion Reactor: " + strconv.FormatInt(r.FusionReactor, 10) + "\n" +
		"      Solar Satellite: " + strconv.FormatInt(r.SolarSatellite, 10) + "\n" +
		"        Metal Storage: " + strconv.FormatInt(r.MetalStorage, 10) + "\n" +
		"      Crystal Storage: " + strconv.FormatInt(r.CrystalStorage, 10) + "\n" +
		"       Deuterium Tank: " + strconv.FormatInt(r.DeuteriumTank, 10)
}
