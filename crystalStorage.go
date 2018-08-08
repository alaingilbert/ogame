package ogame

// CrystalStorage ...
type crystalStorage struct {
	StorageBuilding
}

func newCrystalStorage() *crystalStorage {
	b := new(crystalStorage)
	b.ID = CrystalStorageID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Metal: 1000, Crystal: 500}
	return b
}
