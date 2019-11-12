package ogame

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getNbrV7(doc *goquery.Document, name string) int {
	val, _ := strconv.Atoi(doc.Find("span."+name+" span").First().AttrOr("data-value", "0"))
	return val
}

func extractFacilitiesFromDocV7(doc *goquery.Document) (Facilities, error) {
	res := Facilities{}
	res.RoboticsFactory = getNbrV7(doc, "roboticsFactory")
	res.Shipyard = getNbrV7(doc, "shipyard")
	res.ResearchLab = getNbrV7(doc, "researchLaboratory")
	res.AllianceDepot = getNbrV7(doc, "allianceDepot")
	res.MissileSilo = getNbrV7(doc, "missileSilo")
	res.NaniteFactory = getNbrV7(doc, "naniteFactory")
	res.Terraformer = getNbrV7(doc, "terraformer")
	res.SpaceDock = getNbrV7(doc, "repairDock")
	res.LunarBase = getNbrV7(doc, "lunarBase")         // TODO: ensure name is correct
	res.SensorPhalanx = getNbrV7(doc, "sensorPhalanx") // TODO: ensure name is correct
	res.JumpGate = getNbrV7(doc, "jumpGate")           // TODO: ensure name is correct
	return res, nil
}

func extractDefenseFromDocV7(doc *goquery.Document) (DefensesInfos, error) {
	res := DefensesInfos{}
	res.RocketLauncher = getNbrV7(doc, "rocketLauncher")
	res.LightLaser = getNbrV7(doc, "laserCannonLight")
	res.HeavyLaser = getNbrV7(doc, "laserCannonHeavy")
	res.GaussCannon = getNbrV7(doc, "gaussCannon")
	res.IonCannon = getNbrV7(doc, "ionCannon")
	res.PlasmaTurret = getNbrV7(doc, "plasmaCannon")
	res.SmallShieldDome = getNbrV7(doc, "shieldDomeSmall")
	res.LargeShieldDome = getNbrV7(doc, "shieldDomeLarge")
	res.AntiBallisticMissiles = getNbrV7(doc, "missileInterceptor")
	res.InterplanetaryMissiles = getNbrV7(doc, "missileInterplanetary")
	return res, nil
}

func extractResearchFromDocV7(doc *goquery.Document) Researches {
	doc.Find("span.textlabel").Remove()
	res := Researches{}
	res.EnergyTechnology = getNbrV7(doc, "energyTechnology")
	res.LaserTechnology = getNbrV7(doc, "laserTechnology")
	res.IonTechnology = getNbrV7(doc, "ionTechnology")
	res.HyperspaceTechnology = getNbrV7(doc, "hyperspaceTechnology")
	res.PlasmaTechnology = getNbrV7(doc, "plasmaTechnology")
	res.CombustionDrive = getNbrV7(doc, "combustionDriveTechnology")
	res.ImpulseDrive = getNbrV7(doc, "impulseDriveTechnology")
	res.HyperspaceDrive = getNbrV7(doc, "hyperspaceDriveTechnology")
	res.EspionageTechnology = getNbrV7(doc, "espionageTechnology")
	res.ComputerTechnology = getNbrV7(doc, "computerTechnology")
	res.Astrophysics = getNbrV7(doc, "astrophysicsTechnology")
	res.IntergalacticResearchNetwork = getNbrV7(doc, "researchNetworkTechnology")
	res.GravitonTechnology = getNbrV7(doc, "gravitonTechnology")
	res.WeaponsTechnology = getNbrV7(doc, "weaponsTechnology")
	res.ShieldingTechnology = getNbrV7(doc, "shieldingTechnology")
	res.ArmourTechnology = getNbrV7(doc, "armorTechnology")
	return res
}

func extractShipsFromDocV7(doc *goquery.Document) (ShipsInfos, error) {
	res := ShipsInfos{}
	res.LightFighter = getNbrV7(doc, "fighterLight")
	res.HeavyFighter = getNbrV7(doc, "fighterHeavy")
	res.Cruiser = getNbrV7(doc, "cruiser")
	res.Battleship = getNbrV7(doc, "battleship")
	res.Battlecruiser = getNbrV7(doc, "interceptor")
	res.Bomber = getNbrV7(doc, "bomber")
	res.Destroyer = getNbrV7(doc, "destroyer")
	res.Deathstar = getNbrV7(doc, "deathstar")
	res.Reaper = getNbrV7(doc, "reaper")
	res.Pathfinder = getNbrV7(doc, "explorer")
	res.SmallCargo = getNbrV7(doc, "transporterSmall")
	res.LargeCargo = getNbrV7(doc, "transporterLarge")
	res.ColonyShip = getNbrV7(doc, "colonyShip")
	res.Recycler = getNbrV7(doc, "recycler")
	res.EspionageProbe = getNbrV7(doc, "espionageProbe")
	res.SolarSatellite = getNbrV7(doc, "solarSatellite")
	res.Crawler = getNbrV7(doc, "resbuggy")
	return res, nil
}

func extractResourcesBuildingsFromDocV7(doc *goquery.Document) (ResourcesBuildings, error) {
	res := ResourcesBuildings{}
	res.MetalMine = getNbrV7(doc, "metalMine")
	res.CrystalMine = getNbrV7(doc, "crystalMine")
	res.DeuteriumSynthesizer = getNbrV7(doc, "deuteriumSynthesizer")
	res.SolarPlant = getNbrV7(doc, "solarPlant")
	res.FusionReactor = getNbrV7(doc, "fusionPlant")
	res.SolarSatellite = getNbrV7(doc, "solarSatellite")
	res.MetalStorage = getNbrV7(doc, "metalStorage")
	res.CrystalStorage = getNbrV7(doc, "crystalStorage")
	res.DeuteriumTank = getNbrV7(doc, "deuteriumStorage")
	return res, nil
}

type resourcesRespV7 struct {
	Metal struct {
		ActualFormat string
		Actual       int
		Max          int
		Production   float64
		Tooltip      string
		Class        string
	}
	Crystal struct {
		ActualFormat string
		Actual       int
		Max          int
		Production   float64
		Tooltip      string
		Class        string
	}
	Deuterium struct {
		ActualFormat string
		Actual       int
		Max          int
		Production   float64
		Tooltip      string
		Class        string
	}
	Energy struct {
		ActualFormat string
		Actual       int
		Tooltip      string
		Class        string
	}
	Darkmatter struct {
		ActualFormat string
		Actual       int
		String       string
		Tooltip      string
	}
	HonorScore int
}

func extractResourcesDetailsV7(pageHTML []byte) (out ResourcesDetails, err error) {
	var res resourcesRespV7
	if err = json.Unmarshal(pageHTML, &res); err != nil {
		if isLogged(pageHTML) {
			return out, ErrInvalidPlanetID
		}
		return
	}
	out.Metal.Available = res.Metal.Actual
	out.Metal.StorageCapacity = res.Metal.Max
	out.Crystal.Available = res.Crystal.Actual
	out.Crystal.StorageCapacity = res.Crystal.Max
	out.Deuterium.Available = res.Deuterium.Actual
	out.Deuterium.StorageCapacity = res.Deuterium.Max
	out.Energy.Available = res.Energy.Actual
	out.Darkmatter.Available = res.Darkmatter.Actual
	metalDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(res.Metal.Tooltip))
	crystalDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(res.Crystal.Tooltip))
	deuteriumDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(res.Deuterium.Tooltip))
	darkmatterDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(res.Darkmatter.Tooltip))
	energyDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(res.Energy.Tooltip))
	out.Metal.CurrentProduction = ParseInt(metalDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	out.Crystal.CurrentProduction = ParseInt(crystalDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	out.Deuterium.CurrentProduction = ParseInt(deuteriumDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	out.Energy.CurrentProduction = ParseInt(energyDoc.Find("table tr").Eq(1).Find("td").Eq(0).Text())
	out.Energy.Consumption = ParseInt(energyDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	out.Darkmatter.Purchased = ParseInt(darkmatterDoc.Find("table tr").Eq(1).Find("td").Eq(0).Text())
	out.Darkmatter.Found = ParseInt(darkmatterDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	return
}
