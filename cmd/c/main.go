package main

import (
	"C"
	"github.com/alaingilbert/ogame"
)

var bot *ogame.OGame

//export OGame
func OGame(universe, username, password *C.char) (errorMsg *C.char) {
	var err error
	bot, err = ogame.New(C.GoString(universe), C.GoString(username), C.GoString(password))
	if err != nil {
		errorMsg = C.CString(err.Error())
	}
	return
}

//export SetUserAgent
func SetUserAgent(newUserAgent string) {
	bot.SetUserAgent(newUserAgent)
}

//export ServerURL
func ServerURL() string {
	return bot.ServerURL()
}

//export Login
func Login() error {
	return bot.Login()
}

//export Logout
func Logout() {
	bot.Logout()
}

//export GetUniverseSpeed
func GetUniverseSpeed() C.int {
	return C.int(bot.GetUniverseSpeed())
}

//export ServerVersion
func ServerVersion() *C.char {
	return C.CString(bot.ServerVersion())
}

//export IsUnderAttack
func IsUnderAttack() C.int {
	if bot.IsUnderAttack() {
		return C.int(1)
	}
	return C.int(0)
}

//export GetUserInfos
func GetUserInfos() (playerID C.int, playerName *C.char, points, rank, total, honourPoints C.int) {
	i := bot.GetUserInfos()
	return C.int(i.PlayerID), C.CString(i.PlayerName), C.int(i.Points), C.int(i.Rank), C.int(i.Total), C.int(i.HonourPoints)
}

//export SendMessage
func SendMessage(playerID C.int, msg *C.char) (errorMsg *C.char) {
	err := bot.SendMessage(int(playerID), C.GoString(msg))
	if err != nil {
		errorMsg = C.CString(err.Error())
	}
	return errorMsg
}

// TODO: GetFleets
func GetFleets() {
}

//export CancelFleet
func CancelFleet(fleetID C.int) (errorMsg *C.char) {
	err := bot.CancelFleet(ogame.FleetID(int(fleetID)))
	if err != nil {
		errorMsg = C.CString(err.Error())
	}
	return
}

// TODO: GetAttacks
func GetAttacks() {
}

// TODO: GalaxyInfos
func GalaxyInfos() {
}

//export GetResearch
func GetResearch() (energyTechnology, laserTechnology, ionTechnology, hyperspaceTechnology, plasmaTechnology,
	combustionDrive, impulseDrive, hyperspaceDrive, espionageTechnology, computerTechnology, astrophysics,
	intergalacticResearchNetwork, gravitonTechnology, weaponsTechnology, shieldingTechnology, armourTechnology C.int) {
	r := bot.GetResearch()
	return C.int(r.EnergyTechnology), C.int(r.LaserTechnology), C.int(r.IonTechnology), C.int(r.HyperspaceTechnology),
		C.int(r.PlasmaTechnology), C.int(r.CombustionDrive), C.int(r.ImpulseDrive), C.int(r.HyperspaceDrive),
		C.int(r.EspionageTechnology), C.int(r.ComputerTechnology), C.int(r.Astrophysics), C.int(r.IntergalacticResearchNetwork),
		C.int(r.GravitonTechnology), C.int(r.WeaponsTechnology), C.int(r.ShieldingTechnology), C.int(r.ArmourTechnology)
}

// TODO: GetPlanets

//export GetPlanetByCoord
func GetPlanetByCoord(galaxyIn, systemIn, positionIn C.int) (id C.int, name *C.char, diameter, galaxy, system, position, fieldsBuilt, fieldsTotal,
	temperatureMin, temperatureMax C.int, img *C.char, errorMsg *C.char) {
	p, err := bot.GetPlanetByCoord(ogame.Coordinate{Galaxy: int(galaxyIn), System: int(systemIn), Position: int(positionIn)})
	if err != nil {
		errorMsg = C.CString(err.Error())
	}
	return C.int(int(p.ID)), C.CString(p.Name), C.int(p.Diameter), C.int(p.Coordinate.Galaxy), C.int(p.Coordinate.System),
		C.int(p.Coordinate.Position), C.int(p.Fields.Built), C.int(p.Fields.Total), C.int(p.Temperature.Min),
		C.int(p.Temperature.Max), C.CString(p.Img), errorMsg
}

//export GetPlanet
func GetPlanet(planetID C.int) (id C.int, name *C.char, diameter, galaxy, system, position, fieldsBuilt, fieldsTotal,
	temperatureMin, temperatureMax C.int, img *C.char, errorMsg *C.char) {
	p, err := bot.GetPlanet(ogame.PlanetID(planetID))
	if err != nil {
		errorMsg = C.CString(err.Error())
	}
	return C.int(int(p.ID)), C.CString(p.Name), C.int(p.Diameter), C.int(p.Coordinate.Galaxy), C.int(p.Coordinate.System),
		C.int(p.Coordinate.Position), C.int(p.Fields.Built), C.int(p.Fields.Total), C.int(p.Temperature.Min),
		C.int(p.Temperature.Max), C.CString(p.Img), errorMsg
}

// TODO: GetEspionageReportMessageIDs

//export GetEspionageReport
func GetEspionageReport(msgID C.int) (metal, crystal, deuterium, metalMine, crystalMine, deuteriumSynthesizer,
	solarPlant, fusionReactor, metalStorage, crystalStorage, deuteriumTank, roboticsFactory, shipyard,
	researchLab, allianceDepot, missileSilo, naniteFactory, terraformer, spaceDock, energyTechnology, laserTechnology,
	ionTechnology, hyperspaceTechnology, plasmaTechnology, combustionDrive, impulseDrive, hyperspaceDrive,
	espionageTechnology, computerTechnology, astrophysics, intergalacticResearchNetwork, gravitonTechnology,
	weaponsTechnology, shieldingTechnology, armourTechnology, lightFighter, heavyFighter, cruiser, battleship,
	battlecruiser, bomber, destroyer, deathstar, smallCargo, largeCargo, colonyShip, recycler, espionageProbe,
	solarSatellite, rocketLauncher, lightLaser, heavyLaser, gaussCannon, ionCannon, plasmaTurret,
	smallShieldDome, largeShieldDome, antiBallisticMissiles, interplanetaryMissiles,
	galaxy, system, position C.int, errorMsg *C.char) {
	r, err := bot.GetEspionageReport(int(msgID))
	if err != nil {
		errorMsg = C.CString(err.Error())
		return
	}
	return C.int(r.Resources.Metal), C.int(r.Resources.Crystal), C.int(r.Resources.Deuterium),
		C.int(r.ResourcesBuildings.MetalMine), C.int(r.ResourcesBuildings.CrystalMine),
		C.int(r.ResourcesBuildings.DeuteriumSynthesizer), C.int(r.ResourcesBuildings.SolarPlant),
		C.int(r.ResourcesBuildings.FusionReactor), C.int(r.ResourcesBuildings.MetalStorage),
		C.int(r.ResourcesBuildings.CrystalStorage), C.int(r.ResourcesBuildings.DeuteriumTank),
		C.int(r.Facilities.RoboticsFactory), C.int(r.Facilities.Shipyard), C.int(r.Facilities.ResearchLab),
		C.int(r.Facilities.AllianceDepot), C.int(r.Facilities.MissileSilo), C.int(r.Facilities.NaniteFactory),
		C.int(r.Facilities.Terraformer), C.int(r.Facilities.SpaceDock), C.int(r.Researches.EnergyTechnology),
		C.int(r.Researches.LaserTechnology), C.int(r.Researches.IonTechnology),
		C.int(r.Researches.HyperspaceTechnology), C.int(r.Researches.PlasmaTechnology),
		C.int(r.Researches.CombustionDrive), C.int(r.Researches.ImpulseDrive), C.int(r.Researches.HyperspaceDrive),
		C.int(r.Researches.EspionageTechnology), C.int(r.Researches.ComputerTechnology),
		C.int(r.Researches.Astrophysics), C.int(r.Researches.IntergalacticResearchNetwork),
		C.int(r.Researches.GravitonTechnology), C.int(r.Researches.WeaponsTechnology),
		C.int(r.Researches.ShieldingTechnology), C.int(r.Researches.ArmourTechnology), C.int(r.Ships.LightFighter),
		C.int(r.Ships.HeavyFighter), C.int(r.Ships.Cruiser), C.int(r.Ships.Battleship), C.int(r.Ships.Battlecruiser),
		C.int(r.Ships.Bomber), C.int(r.Ships.Destroyer), C.int(r.Ships.Deathstar), C.int(r.Ships.SmallCargo),
		C.int(r.Ships.LargeCargo), C.int(r.Ships.ColonyShip), C.int(r.Ships.Recycler), C.int(r.Ships.EspionageProbe),
		C.int(r.Ships.SolarSatellite), C.int(r.Defense.RocketLauncher), C.int(r.Defense.LightLaser),
		C.int(r.Defense.HeavyLaser), C.int(r.Defense.GaussCannon), C.int(r.Defense.IonCannon),
		C.int(r.Defense.PlasmaTurret), C.int(r.Defense.SmallShieldDome), C.int(r.Defense.LargeShieldDome),
		C.int(r.Defense.AntiBallisticMissiles), C.int(r.Defense.InterplanetaryMissiles), C.int(r.Coordinate.Galaxy),
		C.int(r.Coordinate.System), C.int(r.Coordinate.Position), errorMsg
}

//export DeleteMessage
func DeleteMessage(msgID C.int) (errorMsg *C.char) {
	err := bot.DeleteMessage(int(msgID))
	if err != nil {
		errorMsg = C.CString(err.Error())
	}
	return
}

//export GetResourceSettings
func GetResourceSettings(planetID C.int) (metalMine, crystalMine, deuteriumSynthesizer, solarPlant, fusionReactor,
	solarSatellite C.int, errorMsg *C.char) {
	r, err := bot.GetResourceSettings(ogame.PlanetID(planetID))
	if err != nil {
		errorMsg = C.CString(err.Error())
		return
	}
	return C.int(r.MetalMine), C.int(r.CrystalMine), C.int(r.DeuteriumSynthesizer), C.int(r.SolarPlant),
		C.int(r.FusionReactor), C.int(r.SolarSatellite), errorMsg
}

//export SetResourceSettings
func SetResourceSettings(planetID, metalMine, crystalMine, deuteriumSynthesizer, solarPlant, fusionReactor,
	solarSatellite C.int) (errorMsg *C.char) {
	settings := ogame.ResourceSettings{
		MetalMine:            int(metalMine),
		CrystalMine:          int(crystalMine),
		DeuteriumSynthesizer: int(deuteriumSynthesizer),
		SolarPlant:           int(solarPlant),
		FusionReactor:        int(fusionReactor),
		SolarSatellite:       int(solarSatellite),
	}
	err := bot.SetResourceSettings(ogame.PlanetID(planetID), settings)
	if err != nil {
		errorMsg = C.CString(err.Error())
	}
	return
}

//export GetResourcesBuildings
func GetResourcesBuildings(planetID C.int) (metalMine, crystalMine, deuteriumSynthesizer, solarPlant, fusionReactor,
	solarSatellite, metalStorage, crystalStorage, deuteriumTank C.int, errorMsg *C.char) {
	r, err := bot.GetResourcesBuildings(ogame.PlanetID(planetID))
	if err != nil {
		errorMsg = C.CString(err.Error())
		return
	}
	return C.int(r.MetalMine), C.int(r.CrystalMine), C.int(r.DeuteriumSynthesizer), C.int(r.SolarPlant),
		C.int(r.FusionReactor), C.int(r.SolarSatellite), C.int(r.MetalStorage), C.int(r.CrystalStorage),
		C.int(r.DeuteriumTank), errorMsg
}

//export GetDefense
func GetDefense(planetID C.int) (rocketLauncher, lightLaser, heavyLaser, gaussCannon, ionCannon, plasmaTurret,
	smallShieldDome, largeShieldDome, antiBallisticMissiles, interplanetaryMissiles C.int, errorMsg *C.char) {
	d, err := bot.GetDefense(ogame.PlanetID(planetID))
	if err != nil {
		errorMsg = C.CString(err.Error())
		return
	}
	return C.int(d.RocketLauncher), C.int(d.LightLaser), C.int(d.HeavyLaser), C.int(d.GaussCannon), C.int(d.IonCannon),
		C.int(d.PlasmaTurret), C.int(d.SmallShieldDome), C.int(d.LargeShieldDome), C.int(d.AntiBallisticMissiles),
		C.int(d.InterplanetaryMissiles), errorMsg
}

//export GetShips
func GetShips(planetID C.int) (lightFighter, heavyFighter, cruiser, battleship, battlecruiser, bomber, destroyer,
	deathstar, smallCargo, largeCargo, colonyShip, recycler, espionageProbe, solarSatellite C.int, errorMsg *C.char) {
	s, err := bot.GetShips(ogame.PlanetID(planetID))
	if err != nil {
		errorMsg = C.CString(err.Error())
		return
	}
	return C.int(s.LightFighter), C.int(s.HeavyFighter), C.int(s.Cruiser), C.int(s.Battleship), C.int(s.Battlecruiser),
		C.int(s.Bomber), C.int(s.Destroyer), C.int(s.Deathstar), C.int(s.SmallCargo), C.int(s.LargeCargo),
		C.int(s.ColonyShip), C.int(s.Recycler), C.int(s.EspionageProbe), C.int(s.SolarSatellite), errorMsg
}

//export GetFacilities
func GetFacilities(planetID C.int) (roboticsFactory, shipyard, researchLab, allianceDepot, missileSilo, naniteFactory,
	terraformer, spaceDock C.int, errorMsg *C.char) {
	f, err := bot.GetFacilities(ogame.PlanetID(planetID))
	if err != nil {
		errorMsg = C.CString(err.Error())
		return
	}
	return C.int(f.RoboticsFactory), C.int(f.Shipyard), C.int(f.ResearchLab), C.int(f.AllianceDepot),
		C.int(f.MissileSilo), C.int(f.NaniteFactory), C.int(f.Terraformer), C.int(f.SpaceDock), errorMsg
}

//export Build
func Build(planetID, ogameID, nbr C.int) (errorMsg *C.char) {
	err := bot.Build(ogame.PlanetID(planetID), ogame.ID(ogameID), int(nbr))
	if err != nil {
		errorMsg = C.CString(err.Error())
	}
	return
}

//export BuildCancelable
func BuildCancelable(planetID, ogameID C.int) (errorMsg *C.char) {
	err := bot.BuildCancelable(ogame.PlanetID(planetID), ogame.ID(ogameID))
	if err != nil {
		errorMsg = C.CString(err.Error())
	}
	return
}

//export BuildProduction
func BuildProduction(planetID, ogameID, nbr C.int) (errorMsg *C.char) {
	err := bot.BuildProduction(ogame.PlanetID(planetID), ogame.ID(ogameID), int(nbr))
	if err != nil {
		errorMsg = C.CString(err.Error())
	}
	return
}

//export BuildBuilding
func BuildBuilding(planetID, buildingID C.int) (errorMsg *C.char) {
	err := bot.BuildBuilding(ogame.PlanetID(planetID), ogame.ID(buildingID))
	if err != nil {
		errorMsg = C.CString(err.Error())
	}
	return
}

//export BuildTechnology
func BuildTechnology(planetID, technologyID C.int) (errorMsg *C.char) {
	err := bot.BuildTechnology(ogame.PlanetID(planetID), ogame.ID(technologyID))
	if err != nil {
		errorMsg = C.CString(err.Error())
	}
	return
}

//export BuildDefense
func BuildDefense(planetID, defenseID, nbr C.int) (errorMsg *C.char) {
	err := bot.BuildDefense(ogame.PlanetID(planetID), ogame.ID(defenseID), int(nbr))
	if err != nil {
		errorMsg = C.CString(err.Error())
	}
	return
}

//export BuildShips
func BuildShips(planetID, shipID, nbr C.int) (errorMsg *C.char) {
	err := bot.BuildShips(ogame.PlanetID(planetID), ogame.ID(shipID), int(nbr))
	if err != nil {
		errorMsg = C.CString(err.Error())
	}
	return
}

// TODO: GetProduction(PlanetID) ([]Quantifiable, error)

//export ConstructionsBeingBuilt
func ConstructionsBeingBuilt(planetID C.int) (buildingID, buildingCountdown, researchID, researchCountdown C.int) {
	a, b, c, d := bot.ConstructionsBeingBuilt(ogame.PlanetID(planetID))
	return C.int(a), C.int(b), C.int(c), C.int(d)
}

//export CancelBuilding
func CancelBuilding(planetID C.int) (errorMsg *C.char) {
	err := bot.CancelBuilding(ogame.PlanetID(planetID))
	if err != nil {
		errorMsg = C.CString(err.Error())
	}
	return
}

//export CancelResearch
func CancelResearch(planetID C.int) (errorMsg *C.char) {
	err := bot.CancelResearch(ogame.PlanetID(planetID))
	if err != nil {
		errorMsg = C.CString(err.Error())
	}
	return
}

//export GetResources
func GetResources(planetID C.int) (metal, crystal, deuterium, energy, darkmatter C.int, errorMsg *C.char) {
	r, err := bot.GetResources(ogame.PlanetID(planetID))
	if err != nil {
		errorMsg = C.CString(err.Error())
	}
	return C.int(r.Metal), C.int(r.Crystal), C.int(r.Deuterium), C.int(r.Energy), C.int(r.Darkmatter), errorMsg
}

//export SendFleet
func SendFleet(planetID, lightFighter, heavyFighter, cruiser, battleship, battlecruiser, bomber, destroyer, deathstar,
	smallCargo, largeCargo, colonyShip, recycler, espionageProbe, speed, galaxy, system, position, mission,
	metal, crystal, deuterium C.int) (fleetID C.int, errorMsg *C.char) {
	ships := make([]ogame.Quantifiable, 0)
	if int(lightFighter) > 0 {
		ships = append(ships, ogame.Quantifiable{ID: ogame.LightFighter, Nbr: int(lightFighter)})
	}
	if int(heavyFighter) > 0 {
		ships = append(ships, ogame.Quantifiable{ID: ogame.LightFighter, Nbr: int(heavyFighter)})
	}
	if int(cruiser) > 0 {
		ships = append(ships, ogame.Quantifiable{ID: ogame.LightFighter, Nbr: int(cruiser)})
	}
	if int(battleship) > 0 {
		ships = append(ships, ogame.Quantifiable{ID: ogame.LightFighter, Nbr: int(battleship)})
	}
	if int(battlecruiser) > 0 {
		ships = append(ships, ogame.Quantifiable{ID: ogame.LightFighter, Nbr: int(battlecruiser)})
	}
	if int(bomber) > 0 {
		ships = append(ships, ogame.Quantifiable{ID: ogame.LightFighter, Nbr: int(bomber)})
	}
	if int(destroyer) > 0 {
		ships = append(ships, ogame.Quantifiable{ID: ogame.LightFighter, Nbr: int(destroyer)})
	}
	if int(deathstar) > 0 {
		ships = append(ships, ogame.Quantifiable{ID: ogame.LightFighter, Nbr: int(deathstar)})
	}
	if int(smallCargo) > 0 {
		ships = append(ships, ogame.Quantifiable{ID: ogame.LightFighter, Nbr: int(smallCargo)})
	}
	if int(largeCargo) > 0 {
		ships = append(ships, ogame.Quantifiable{ID: ogame.LightFighter, Nbr: int(largeCargo)})
	}
	if int(colonyShip) > 0 {
		ships = append(ships, ogame.Quantifiable{ID: ogame.LightFighter, Nbr: int(colonyShip)})
	}
	if int(recycler) > 0 {
		ships = append(ships, ogame.Quantifiable{ID: ogame.LightFighter, Nbr: int(recycler)})
	}
	if int(espionageProbe) > 0 {
		ships = append(ships, ogame.Quantifiable{ID: ogame.LightFighter, Nbr: int(espionageProbe)})
	}
	fleetIDInt, err := bot.SendFleet(
		ogame.PlanetID(planetID),
		ships, ogame.Speed(speed),
		ogame.Coordinate{int(galaxy), int(system), int(position)},
		ogame.MissionID(mission),
		ogame.Resources{Metal: int(metal), Crystal: int(crystal), Deuterium: int(deuterium)},
	)
	if err != nil {
		errorMsg = C.CString(err.Error())
	}
	return C.int(fleetIDInt), errorMsg
}

func main() {}
