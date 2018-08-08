package ogame

// MetalStorage ...
type metalStorage struct {
	StorageBuilding
}

func newMetalStorage() *metalStorage {
	b := new(metalStorage)
	b.ID = MetalStorageID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Metal: 1000}
	return b
}
