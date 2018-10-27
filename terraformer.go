package ogame

type terraformer struct {
	BaseBuilding
}

func newTerraformer() *terraformer {
	b := new(terraformer)
	b.Name = "terraformer"
	b.ID = TerraformerID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Metal: 0, Crystal: 50000, Deuterium: 100000, Energy: 1000}
	b.Requirements = map[ID]int{NaniteFactoryID: 1, EnergyTechnologyID: 12}
	return b
}
