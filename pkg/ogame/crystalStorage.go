package ogame

type crystalStorage struct {
	storageBuilding
}

func newCrystalStorage() *crystalStorage {
	b := new(crystalStorage)
	b.Name = "crystal storage"
	b.ID = CrystalStorageID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Metal: 1000, Crystal: 500}
	return b
}
