package ogame

type crawler struct {
	BaseShip
}

func newCrawler() *crawler {
	c := new(crawler)
	c.Name = "crawler"
	c.ID = CrawlerID
	c.StructuralIntegrity = 4000
	c.ShieldPower = 1
	c.WeaponPower = 1
	c.BaseCargoCapacity = 0
	c.BaseSpeed = 0
	c.FuelConsumption = 0
	c.RapidfireFrom = map[ID]int64{LightFighterID: 5, HeavyFighterID: 5, CruiserID: 5, BattleshipID: 5,
		BattlecruiserID: 5, BomberID: 5, DestroyerID: 5, DeathstarID: 1250, ReaperID: 5, PathfinderID: 5,
		SmallCargoID: 5, LargeCargoID: 5, ColonyShipID: 5, RecyclerID: 5}
	c.RapidfireAgainst = map[ID]int64{}
	c.Price = Resources{Metal: 2000, Crystal: 2000, Deuterium: 1000}
	c.Requirements = map[ID]int64{ShipyardID: 5, CombustionDriveID: 4, ArmourTechnologyID: 4, LaserTechnologyID: 4}
	return c
}
