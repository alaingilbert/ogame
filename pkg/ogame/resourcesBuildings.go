package ogame

import (
	"github.com/alaingilbert/ogame/pkg/utils"
)

// LazyResourcesBuildings ...
type LazyResourcesBuildings func() ResourcesBuildings

func (r LazyResourcesBuildings) GetMetalMine() int64            { return r().MetalMine }
func (r LazyResourcesBuildings) GetCrystalMine() int64          { return r().CrystalMine }
func (r LazyResourcesBuildings) GetDeuteriumSynthesizer() int64 { return r().DeuteriumSynthesizer }
func (r LazyResourcesBuildings) GetSolarPlant() int64           { return r().SolarPlant }
func (r LazyResourcesBuildings) GetFusionReactor() int64        { return r().FusionReactor }
func (r LazyResourcesBuildings) GetSolarSatellite() int64       { return r().SolarSatellite }
func (r LazyResourcesBuildings) GetMetalStorage() int64         { return r().MetalStorage }
func (r LazyResourcesBuildings) GetCrystalStorage() int64       { return r().CrystalStorage }
func (r LazyResourcesBuildings) GetDeuteriumTank() int64        { return r().DeuteriumTank }

// ResourcesBuildings represent a planet resource buildings
type ResourcesBuildings struct {
	MetalMine            int64 // 1
	CrystalMine          int64 // 2
	DeuteriumSynthesizer int64 // 3
	SolarPlant           int64 // 4
	FusionReactor        int64 // 12
	SolarSatellite       int64 // 212
	MetalStorage         int64 // 22
	CrystalStorage       int64 // 23
	DeuteriumTank        int64 // 24
}

func (r ResourcesBuildings) GetMetalMine() int64            { return r.MetalMine }
func (r ResourcesBuildings) GetCrystalMine() int64          { return r.CrystalMine }
func (r ResourcesBuildings) GetDeuteriumSynthesizer() int64 { return r.DeuteriumSynthesizer }
func (r ResourcesBuildings) GetSolarPlant() int64           { return r.SolarPlant }
func (r ResourcesBuildings) GetFusionReactor() int64        { return r.FusionReactor }
func (r ResourcesBuildings) GetSolarSatellite() int64       { return r.SolarSatellite }
func (r ResourcesBuildings) GetMetalStorage() int64         { return r.MetalStorage }
func (r ResourcesBuildings) GetCrystalStorage() int64       { return r.CrystalStorage }
func (r ResourcesBuildings) GetDeuteriumTank() int64        { return r.DeuteriumTank }

// Lazy returns a function that return self
func (r ResourcesBuildings) Lazy() LazyResourcesBuildings {
	return func() ResourcesBuildings { return r }
}

// ByID gets the resource building level from a building id
func (r ResourcesBuildings) ByID(id ID) int64 {
	return resourcesBuildingByID(id, r)
}

func resourcesBuildingByID(id ID, r IResourcesBuildings) int64 {
	switch id {
	case MetalMineID:
		return r.GetMetalMine()
	case CrystalMineID:
		return r.GetCrystalMine()
	case DeuteriumSynthesizerID:
		return r.GetDeuteriumSynthesizer()
	case SolarPlantID:
		return r.GetSolarPlant()
	case FusionReactorID:
		return r.GetFusionReactor()
	case SolarSatelliteID:
		return r.GetSolarSatellite()
	case MetalStorageID:
		return r.GetMetalStorage()
	case CrystalStorageID:
		return r.GetCrystalStorage()
	case DeuteriumTankID:
		return r.GetDeuteriumTank()
	}
	return 0
}

func (r ResourcesBuildings) String() string {
	return "\n" +
		"           Metal Mine: " + utils.FI64(r.MetalMine) + "\n" +
		"         Crystal Mine: " + utils.FI64(r.CrystalMine) + "\n" +
		"Deuterium Synthesizer: " + utils.FI64(r.DeuteriumSynthesizer) + "\n" +
		"          Solar Plant: " + utils.FI64(r.SolarPlant) + "\n" +
		"       Fusion Reactor: " + utils.FI64(r.FusionReactor) + "\n" +
		"      Solar Satellite: " + utils.FI64(r.SolarSatellite) + "\n" +
		"        Metal Storage: " + utils.FI64(r.MetalStorage) + "\n" +
		"      Crystal Storage: " + utils.FI64(r.CrystalStorage) + "\n" +
		"       Deuterium Tank: " + utils.FI64(r.DeuteriumTank)
}
