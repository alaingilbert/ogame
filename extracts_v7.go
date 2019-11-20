package ogame

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/alaingilbert/clockwork"
)

func getNbrV7(doc *goquery.Document, name string) int {
	val, _ := strconv.Atoi(doc.Find("span."+name+" span.level").First().AttrOr("data-value", "0"))
	return val
}

func getNbrV7Ships(doc *goquery.Document, name string) int {
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
	res.RocketLauncher = getNbrV7Ships(doc, "rocketLauncher")
	res.LightLaser = getNbrV7Ships(doc, "laserCannonLight")
	res.HeavyLaser = getNbrV7Ships(doc, "laserCannonHeavy")
	res.GaussCannon = getNbrV7Ships(doc, "gaussCannon")
	res.IonCannon = getNbrV7Ships(doc, "ionCannon")
	res.PlasmaTurret = getNbrV7Ships(doc, "plasmaCannon")
	res.SmallShieldDome = getNbrV7Ships(doc, "shieldDomeSmall")
	res.LargeShieldDome = getNbrV7Ships(doc, "shieldDomeLarge")
	res.AntiBallisticMissiles = getNbrV7Ships(doc, "missileInterceptor")
	res.InterplanetaryMissiles = getNbrV7Ships(doc, "missileInterplanetary")
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
	res.LightFighter = getNbrV7Ships(doc, "fighterLight")
	res.HeavyFighter = getNbrV7Ships(doc, "fighterHeavy")
	res.Cruiser = getNbrV7Ships(doc, "cruiser")
	res.Battleship = getNbrV7Ships(doc, "battleship")
	res.Battlecruiser = getNbrV7Ships(doc, "interceptor")
	res.Bomber = getNbrV7Ships(doc, "bomber")
	res.Destroyer = getNbrV7Ships(doc, "destroyer")
	res.Deathstar = getNbrV7Ships(doc, "deathstar")
	res.Reaper = getNbrV7Ships(doc, "reaper")
	res.Pathfinder = getNbrV7Ships(doc, "explorer")
	res.SmallCargo = getNbrV7Ships(doc, "transporterSmall")
	res.LargeCargo = getNbrV7Ships(doc, "transporterLarge")
	res.ColonyShip = getNbrV7Ships(doc, "colonyShip")
	res.Recycler = getNbrV7Ships(doc, "recycler")
	res.EspionageProbe = getNbrV7Ships(doc, "espionageProbe")
	res.SolarSatellite = getNbrV7Ships(doc, "solarSatellite")
	res.Crawler = getNbrV7Ships(doc, "resbuggy")
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

func extractConstructionsV7(pageHTML []byte, clock clockwork.Clock) (buildingID ID, buildingCountdown int, researchID ID, researchCountdown int) {
	buildingCountdownMatch := regexp.MustCompile(`var restTimebuilding = (\d+) -`).FindSubmatch(pageHTML)
	if len(buildingCountdownMatch) > 0 {
		buildingCountdown = int(int64(toInt(buildingCountdownMatch[1])) - clock.Now().Unix())
		buildingIDInt := toInt(regexp.MustCompile(`onclick="cancelbuilding\((\d+),`).FindSubmatch(pageHTML)[1])
		buildingID = ID(buildingIDInt)
	}
	researchCountdownMatch := regexp.MustCompile(`var restTimeresearch = (\d+) -`).FindSubmatch(pageHTML)
	if len(researchCountdownMatch) > 0 {
		researchCountdown = int(int64(toInt(researchCountdownMatch[1])) - clock.Now().Unix())
		researchIDInt := toInt(regexp.MustCompile(`onclick="cancelresearch\((\d+),`).FindSubmatch(pageHTML)[1])
		researchID = ID(researchIDInt)
	}
	return
}

func extractFleet1ShipsFromDocV7(doc *goquery.Document) (s ShipsInfos) {
	onclick := doc.Find("div#fleetdispatchcomponent")
	h, _ := onclick.Html()
	matches := regexp.MustCompile(`var shipsOnPlanet = ([^;]+);`).FindStringSubmatch(h)
	if len(matches) == 0 {
		return
	}
	m := matches[1]
	var res []struct {
		ID     int `json:"id"`
		Number int `json:"number"`
	}
	if err := json.Unmarshal([]byte(m), &res); err != nil {
		return
	}
	for _, obj := range res {
		s.Set(ID(obj.ID), obj.Number)
	}
	return
}

func extractCombatReportMessagesFromDocV7(doc *goquery.Document) ([]CombatReportSummary, int) {
	msgs := make([]CombatReportSummary, 0)
	nbPage, _ := strconv.Atoi(doc.Find("ul.pagination li").Last().AttrOr("data-page", "1"))
	doc.Find("li.msg").Each(func(i int, s *goquery.Selection) {
		if idStr, exists := s.Attr("data-msg-id"); exists {
			if id, err := strconv.Atoi(idStr); err == nil {
				report := CombatReportSummary{ID: id}
				report.Destination = extractCoordV6(s.Find("div.msg_head a").Text())
				if s.Find("div.msg_head figure").HasClass("planet") {
					report.Destination.Type = PlanetType
				} else if s.Find("div.msg_head figure").HasClass("moon") {
					report.Destination.Type = MoonType
				} else {
					report.Destination.Type = PlanetType
				}
				resTitle := s.Find("span.msg_content div.combatLeftSide span").Eq(1).AttrOr("title", "")
				m := regexp.MustCompile(`([\d.]+)<br/>[^\d]*([\d.]+)<br/>[^\d]*([\d.]+)`).FindStringSubmatch(resTitle)
				if len(m) == 4 {
					report.Metal = ParseInt(m[1])
					report.Crystal = ParseInt(m[2])
					report.Deuterium = ParseInt(m[3])
				}
				resText := s.Find("span.msg_content div.combatLeftSide span").Eq(1).Text()
				m = regexp.MustCompile(`[\d.]+[^\d]*([\d.]+)`).FindStringSubmatch(resText)
				if len(m) == 2 {
					report.Loot = ParseInt(m[1])
				}
				msgDate, _ := time.Parse("02.01.2006 15:04:05", s.Find("span.msg_date").Text())
				report.CreatedAt = msgDate

				link := s.Find("div.msg_actions a span.icon_attack").Parent().AttrOr("href", "")
				m = regexp.MustCompile(`page=ingame&component=fleetdispatch&galaxy=(\d+)&system=(\d+)&position=(\d+)&type=(\d+)&`).FindStringSubmatch(link)
				if len(m) != 5 {
					return
				}
				galaxy, _ := strconv.Atoi(m[1])
				system, _ := strconv.Atoi(m[2])
				position, _ := strconv.Atoi(m[3])
				planetType, _ := strconv.Atoi(m[4])
				report.Origin = &Coordinate{galaxy, system, position, CelestialType(planetType)}
				if report.Origin.Equal(report.Destination) {
					report.Origin = nil
				}

				msgs = append(msgs, report)
			}
		}
	})
	return msgs, nbPage
}
