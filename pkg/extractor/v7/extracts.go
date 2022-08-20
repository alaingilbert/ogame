package v7

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/alaingilbert/ogame/pkg/extractor/v6"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/utils"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/alaingilbert/clockwork"
)

func GetNbrV7(doc *goquery.Document, name string) int64 {
	val := utils.DoParseI64(doc.Find("span."+name+" span.level").First().AttrOr("data-value", "0"))
	return val
}

func getNbrV7Ships(doc *goquery.Document, name string) int64 {
	val := utils.DoParseI64(doc.Find("span."+name+" span.amount").First().AttrOr("data-value", "0"))
	return val
}

func extractPremiumTokenV7(pageHTML []byte, days int64) (token string, err error) {
	rgx := regexp.MustCompile(`\?page=premium&buynow=1&type=\d&days=` + utils.FI64(days) + `&token=(\w+)`)
	m := rgx.FindSubmatch(pageHTML)
	if len(m) < 2 {
		return "", errors.New("unable to find token")
	}
	token = string(m[1])
	return
}

func extractResourcesDetailsFromFullPageFromDocV7(doc *goquery.Document) ogame.ResourcesDetails {
	out := ogame.ResourcesDetails{}
	out.Metal.Available = utils.ParseInt(strings.Split(doc.Find("span#resources_metal").AttrOr("data-raw", "0"), ".")[0])
	out.Crystal.Available = utils.ParseInt(strings.Split(doc.Find("span#resources_crystal").AttrOr("data-raw", "0"), ".")[0])
	out.Deuterium.Available = utils.ParseInt(strings.Split(doc.Find("span#resources_deuterium").AttrOr("data-raw", "0"), ".")[0])
	out.Energy.Available = utils.ParseInt(doc.Find("span#resources_energy").Text())
	out.Darkmatter.Available = utils.ParseInt(doc.Find("span#resources_darkmatter").Text())
	metalDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(doc.Find("li#metal_box").AttrOr("title", "")))
	crystalDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(doc.Find("li#crystal_box").AttrOr("title", "")))
	deuteriumDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(doc.Find("li#deuterium_box").AttrOr("title", "")))
	energyDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(doc.Find("li#energy_box").AttrOr("title", "")))
	darkmatterDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(doc.Find("li#darkmatter_box").AttrOr("title", "")))
	out.Metal.StorageCapacity = utils.ParseInt(metalDoc.Find("table tr").Eq(1).Find("td").Eq(0).Text())
	out.Metal.CurrentProduction = utils.ParseInt(metalDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	out.Crystal.StorageCapacity = utils.ParseInt(crystalDoc.Find("table tr").Eq(1).Find("td").Eq(0).Text())
	out.Crystal.CurrentProduction = utils.ParseInt(crystalDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	out.Deuterium.StorageCapacity = utils.ParseInt(deuteriumDoc.Find("table tr").Eq(1).Find("td").Eq(0).Text())
	out.Deuterium.CurrentProduction = utils.ParseInt(deuteriumDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	out.Energy.CurrentProduction = utils.ParseInt(energyDoc.Find("table tr").Eq(1).Find("td").Eq(0).Text())
	out.Energy.Consumption = utils.ParseInt(energyDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	out.Darkmatter.Purchased = utils.ParseInt(darkmatterDoc.Find("table tr").Eq(1).Find("td").Eq(0).Text())
	out.Darkmatter.Found = utils.ParseInt(darkmatterDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	return out
}

func ExtractFacilitiesFromDocV7(doc *goquery.Document) (ogame.Facilities, error) {
	res := ogame.Facilities{}
	res.RoboticsFactory = GetNbrV7(doc, "roboticsFactory")
	res.Shipyard = GetNbrV7(doc, "shipyard")
	res.ResearchLab = GetNbrV7(doc, "researchLaboratory")
	res.AllianceDepot = GetNbrV7(doc, "allianceDepot")
	res.MissileSilo = GetNbrV7(doc, "missileSilo")
	res.NaniteFactory = GetNbrV7(doc, "naniteFactory")
	res.Terraformer = GetNbrV7(doc, "terraformer")
	res.SpaceDock = GetNbrV7(doc, "repairDock")
	res.LunarBase = GetNbrV7(doc, "lunarBase")         // TODO: ensure name is correct
	res.SensorPhalanx = GetNbrV7(doc, "sensorPhalanx") // TODO: ensure name is correct
	res.JumpGate = GetNbrV7(doc, "jumpGate")           // TODO: ensure name is correct
	return res, nil
}

func extractDefenseFromDocV7(doc *goquery.Document) (ogame.DefensesInfos, error) {
	res := ogame.DefensesInfos{}
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

func extractResearchFromDocV7(doc *goquery.Document) ogame.Researches {
	doc.Find("span.textlabel").Remove()
	res := ogame.Researches{}
	res.EnergyTechnology = GetNbrV7(doc, "energyTechnology")
	res.LaserTechnology = GetNbrV7(doc, "laserTechnology")
	res.IonTechnology = GetNbrV7(doc, "ionTechnology")
	res.HyperspaceTechnology = GetNbrV7(doc, "hyperspaceTechnology")
	res.PlasmaTechnology = GetNbrV7(doc, "plasmaTechnology")
	res.CombustionDrive = GetNbrV7(doc, "combustionDriveTechnology")
	res.ImpulseDrive = GetNbrV7(doc, "impulseDriveTechnology")
	res.HyperspaceDrive = GetNbrV7(doc, "hyperspaceDriveTechnology")
	res.EspionageTechnology = GetNbrV7(doc, "espionageTechnology")
	res.ComputerTechnology = GetNbrV7(doc, "computerTechnology")
	res.Astrophysics = GetNbrV7(doc, "astrophysicsTechnology")
	res.IntergalacticResearchNetwork = GetNbrV7(doc, "researchNetworkTechnology")
	res.GravitonTechnology = GetNbrV7(doc, "gravitonTechnology")
	res.WeaponsTechnology = GetNbrV7(doc, "weaponsTechnology")
	res.ShieldingTechnology = GetNbrV7(doc, "shieldingTechnology")
	res.ArmourTechnology = GetNbrV7(doc, "armorTechnology")
	return res
}

func extractShipsFromDocV7(doc *goquery.Document) (ogame.ShipsInfos, error) {
	res := ogame.ShipsInfos{}
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

func extractResourcesBuildingsFromDocV7(doc *goquery.Document) (ogame.ResourcesBuildings, error) {
	res := ogame.ResourcesBuildings{}
	res.MetalMine = GetNbrV7(doc, "metalMine")
	res.CrystalMine = GetNbrV7(doc, "crystalMine")
	res.DeuteriumSynthesizer = GetNbrV7(doc, "deuteriumSynthesizer")
	res.SolarPlant = GetNbrV7(doc, "solarPlant")
	res.FusionReactor = GetNbrV7(doc, "fusionPlant")
	res.SolarSatellite = getNbrV7Ships(doc, "solarSatellite")
	res.MetalStorage = GetNbrV7(doc, "metalStorage")
	res.CrystalStorage = GetNbrV7(doc, "crystalStorage")
	res.DeuteriumTank = GetNbrV7(doc, "deuteriumStorage")
	return res, nil
}

type resourcesRespV7 struct {
	Metal struct {
		ActualFormat string
		Actual       int64
		Max          int64
		Production   float64
		Tooltip      string
		Class        string
	}
	Crystal struct {
		ActualFormat string
		Actual       int64
		Max          int64
		Production   float64
		Tooltip      string
		Class        string
	}
	Deuterium struct {
		ActualFormat string
		Actual       int64
		Max          int64
		Production   float64
		Tooltip      string
		Class        string
	}
	Energy struct {
		ActualFormat string
		Actual       int64
		Tooltip      string
		Class        string
	}
	Darkmatter struct {
		ActualFormat string
		Actual       int64
		String       string
		Tooltip      string
	}
	HonorScore int64
}

func extractResourcesDetailsV7(pageHTML []byte) (out ogame.ResourcesDetails, err error) {
	var res resourcesRespV7
	if err = json.Unmarshal(pageHTML, &res); err != nil {
		if v6.IsLogged(pageHTML) {
			return out, ogame.ErrInvalidPlanetID
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
	out.Metal.CurrentProduction = utils.ParseInt(metalDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	out.Crystal.CurrentProduction = utils.ParseInt(crystalDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	out.Deuterium.CurrentProduction = utils.ParseInt(deuteriumDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	out.Energy.CurrentProduction = utils.ParseInt(energyDoc.Find("table tr").Eq(1).Find("td").Eq(0).Text())
	out.Energy.Consumption = utils.ParseInt(energyDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	out.Darkmatter.Purchased = utils.ParseInt(darkmatterDoc.Find("table tr").Eq(1).Find("td").Eq(0).Text())
	out.Darkmatter.Found = utils.ParseInt(darkmatterDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	return
}

func ExtractConstructionsV7(pageHTML []byte, clock clockwork.Clock) (buildingID ogame.ID, buildingCountdown int64, researchID ogame.ID, researchCountdown int64) {
	buildingCountdownMatch := regexp.MustCompile(`var restTimebuilding = (\d+) -`).FindSubmatch(pageHTML)
	if len(buildingCountdownMatch) > 0 {
		buildingCountdown = int64(utils.ToInt(buildingCountdownMatch[1])) - clock.Now().Unix()
		buildingIDInt := utils.ToInt(regexp.MustCompile(`onclick="cancelbuilding\((\d+),`).FindSubmatch(pageHTML)[1])
		buildingID = ogame.ID(buildingIDInt)
	}
	researchCountdownMatch := regexp.MustCompile(`var restTimeresearch = (\d+) -`).FindSubmatch(pageHTML)
	if len(researchCountdownMatch) > 0 {
		researchCountdown = int64(utils.ToInt(researchCountdownMatch[1])) - clock.Now().Unix()
		researchIDInt := utils.ToInt(regexp.MustCompile(`onclick="cancelresearch\((\d+),`).FindSubmatch(pageHTML)[1])
		researchID = ogame.ID(researchIDInt)
	}
	return
}

func extractIPMFromDocV7(doc *goquery.Document) (duration, max int64, token string) {
	duration = utils.DoParseI64(doc.Find("span#timer").AttrOr("data-duration", "0"))
	max = utils.DoParseI64(doc.Find("input[name=missileCount]").AttrOr("data-max", "0"))
	token = doc.Find("input[name=token]").AttrOr("value", "")
	return
}

func extractFleet1ShipsFromDocV7(doc *goquery.Document) (s ogame.ShipsInfos) {
	onclick := doc.Find("div#fleetdispatchcomponent")
	h, _ := onclick.Html()
	matches := regexp.MustCompile(`var shipsOnPlanet = ([^;]+);`).FindStringSubmatch(h)
	if len(matches) == 0 {
		return
	}
	m := matches[1]
	var res []struct {
		ID     int64 `json:"id"`
		Number int64 `json:"number"`
	}
	if err := json.Unmarshal([]byte(m), &res); err != nil {
		return
	}
	for _, obj := range res {
		s.Set(ogame.ID(obj.ID), obj.Number)
	}
	return
}

func extractCombatReportMessagesFromDocV7(doc *goquery.Document) ([]ogame.CombatReportSummary, int64) {
	msgs := make([]ogame.CombatReportSummary, 0)
	nbPage := utils.DoParseI64(doc.Find("ul.pagination li").Last().AttrOr("data-page", "1"))
	doc.Find("li.msg").Each(func(i int, s *goquery.Selection) {
		if idStr, exists := s.Attr("data-msg-id"); exists {
			if id, err := utils.ParseI64(idStr); err == nil {
				report := ogame.CombatReportSummary{ID: id}
				report.Destination = v6.ExtractCoordV6(s.Find("div.msg_head a").Text())
				if s.Find("div.msg_head figure").HasClass("planet") {
					report.Destination.Type = ogame.PlanetType
				} else if s.Find("div.msg_head figure").HasClass("moon") {
					report.Destination.Type = ogame.MoonType
				} else {
					report.Destination.Type = ogame.PlanetType
				}
				apiKeyTitle := s.Find("span.icon_apikey").AttrOr("title", "")
				m := regexp.MustCompile(`'(cr-[^']+)'`).FindStringSubmatch(apiKeyTitle)
				if len(m) == 2 {
					report.APIKey = m[1]
				}
				resTitle := s.Find("span.msg_content div.combatLeftSide span").Eq(1).AttrOr("title", "")
				m = regexp.MustCompile(`([\d.,]+)<br/>[^\d]*([\d.,]+)<br/>[^\d]*([\d.,]+)`).FindStringSubmatch(resTitle)
				if len(m) == 4 {
					report.Metal = utils.ParseInt(m[1])
					report.Crystal = utils.ParseInt(m[2])
					report.Deuterium = utils.ParseInt(m[3])
				}
				debrisFieldTitle := s.Find("span.msg_content div.combatLeftSide span").Eq(2).AttrOr("title", "0")
				report.DebrisField = utils.ParseInt(debrisFieldTitle)
				resText := s.Find("span.msg_content div.combatLeftSide span").Eq(1).Text()
				m = regexp.MustCompile(`[\d.,]+[^\d]*([\d.,]+)`).FindStringSubmatch(resText)
				if len(m) == 2 {
					report.Loot = utils.ParseInt(m[1])
				}
				msgDate, _ := time.Parse("02.01.2006 15:04:05", s.Find("span.msg_date").Text())
				report.CreatedAt = msgDate

				link := s.Find("div.msg_actions a span.icon_attack").Parent().AttrOr("href", "")
				m = regexp.MustCompile(`page=ingame&component=fleetdispatch&galaxy=(\d+)&system=(\d+)&position=(\d+)&type=(\d+)&`).FindStringSubmatch(link)
				if len(m) != 5 {
					return
				}
				galaxy := utils.DoParseI64(m[1])
				system := utils.DoParseI64(m[2])
				position := utils.DoParseI64(m[3])
				planetType := utils.DoParseI64(m[4])
				report.Origin = &ogame.Coordinate{galaxy, system, position, ogame.CelestialType(planetType)}
				if report.Origin.Equal(report.Destination) {
					report.Origin = nil
				}

				msgs = append(msgs, report)
			}
		}
	})
	return msgs, nbPage
}

func extractEspionageReportFromDocV7(doc *goquery.Document, location *time.Location) (ogame.EspionageReport, error) {
	report := ogame.EspionageReport{}
	report.ID = utils.DoParseI64(doc.Find("div.detail_msg").AttrOr("data-msg-id", "0"))
	spanLink := doc.Find("span.msg_title a").First()
	txt := spanLink.Text()
	figure := spanLink.Find("figure").First()
	r := regexp.MustCompile(`([^\[]+) \[(\d+):(\d+):(\d+)]`)
	m := r.FindStringSubmatch(txt)
	if len(m) == 5 {
		report.Coordinate.Galaxy = utils.DoParseI64(m[2])
		report.Coordinate.System = utils.DoParseI64(m[3])
		report.Coordinate.Position = utils.DoParseI64(m[4])
	} else {
		return report, errors.New("failed to extract coordinate")
	}
	if figure.HasClass("planet") {
		report.Coordinate.Type = ogame.PlanetType
	} else if figure.HasClass("moon") {
		report.Coordinate.Type = ogame.MoonType
	}
	messageType := ogame.Report
	if doc.Find("span.espionageDefText").Size() > 0 {
		messageType = ogame.Action
	}
	report.Type = messageType
	msgDateRaw := doc.Find("span.msg_date").Text()
	msgDate, _ := time.ParseInLocation("02.01.2006 15:04:05", msgDateRaw, location)
	report.Date = msgDate.In(time.Local)

	username := doc.Find("div.detail_txt").First().Find("span span").First().Text()
	username = strings.TrimSpace(username)
	split := strings.Split(username, "(i")
	if len(split) > 0 {
		report.Username = strings.TrimSpace(split[0])
	}

	// Bandit, Starlord
	banditstarlord := doc.Find("div.detail_txt").First().Find("span")
	if banditstarlord.HasClass("honorRank") {
		report.IsBandit = banditstarlord.HasClass("rank_bandit1") || banditstarlord.HasClass("rank_bandit2") || banditstarlord.HasClass("rank_bandit3")
		report.IsStarlord = banditstarlord.HasClass("rank_starlord1") || banditstarlord.HasClass("rank_starlord2") || banditstarlord.HasClass("rank_starlord3")
	}

	// IsInactive, IsLongInactive
	inactive := doc.Find("div.detail_txt").First().Find("span")
	if inactive.HasClass("status_abbr_longinactive") {
		report.IsInactive = true
		report.IsLongInactive = true
	} else if inactive.HasClass("status_abbr_inactive") {
		report.IsInactive = true
	}

	// APIKey
	apikey, _ := doc.Find("span.icon_apikey").Attr("title")
	apiDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(apikey))
	report.APIKey = apiDoc.Find("input").First().AttrOr("value", "")

	// Inactivity timer
	activity := doc.Find("div.detail_txt").Eq(2).Find("font")
	if len(activity.Text()) == 2 {
		report.LastActivity = utils.ParseInt(activity.Text())
	}

	// CounterEspionage
	ceTxt := doc.Find("div.detail_txt").Eq(2).Text()
	m1 := regexp.MustCompile(`(\d+)%`).FindStringSubmatch(ceTxt)
	if len(m1) == 2 {
		report.CounterEspionage = utils.DoParseI64(m1[1])
	}

	hasError := false
	doc.Find("ul.detail_list").Each(func(i int, s *goquery.Selection) {
		dataType := s.AttrOr("data-type", "")
		if dataType == "resources" {
			report.Metal = utils.ParseInt(s.Find("li").Eq(0).AttrOr("title", "0"))
			report.Crystal = utils.ParseInt(s.Find("li").Eq(1).AttrOr("title", "0"))
			report.Deuterium = utils.ParseInt(s.Find("li").Eq(2).AttrOr("title", "0"))
			report.Energy = utils.ParseInt(s.Find("li").Eq(3).AttrOr("title", "0"))
		} else if dataType == "buildings" {
			report.HasBuildingsInformation = s.Find("li.detail_list_fail").Size() == 0
			s.Find("li.detail_list_el").EachWithBreak(func(i int, s2 *goquery.Selection) bool {
				img := s2.Find("img")
				if img.Size() == 0 {
					hasError = true
					return false
				}
				imgClass := img.AttrOr("class", "")
				r := regexp.MustCompile(`building(\d+)`)
				buildingID := utils.DoParseI64(r.FindStringSubmatch(imgClass)[1])
				l := utils.ParseInt(s2.Find("span.fright").Text())
				level := &l
				switch ogame.ID(buildingID) {
				case ogame.MetalMine.ID:
					report.MetalMine = level
				case ogame.CrystalMine.ID:
					report.CrystalMine = level
				case ogame.DeuteriumSynthesizer.ID:
					report.DeuteriumSynthesizer = level
				case ogame.SolarPlant.ID:
					report.SolarPlant = level
				case ogame.FusionReactor.ID:
					report.FusionReactor = level
				case ogame.MetalStorage.ID:
					report.MetalStorage = level
				case ogame.CrystalStorage.ID:
					report.CrystalStorage = level
				case ogame.DeuteriumTank.ID:
					report.DeuteriumTank = level
				case ogame.AllianceDepot.ID:
					report.AllianceDepot = level
				case ogame.RoboticsFactory.ID:
					report.RoboticsFactory = level
				case ogame.Shipyard.ID:
					report.Shipyard = level
				case ogame.ResearchLab.ID:
					report.ResearchLab = level
				case ogame.MissileSilo.ID:
					report.MissileSilo = level
				case ogame.NaniteFactory.ID:
					report.NaniteFactory = level
				case ogame.Terraformer.ID:
					report.Terraformer = level
				case ogame.SpaceDock.ID:
					report.SpaceDock = level
				case ogame.LunarBase.ID:
					report.LunarBase = level
				case ogame.SensorPhalanx.ID:
					report.SensorPhalanx = level
				case ogame.JumpGate.ID:
					report.JumpGate = level
				}
				return true
			})
		} else if dataType == "research" {
			report.HasResearchesInformation = s.Find("li.detail_list_fail").Size() == 0
			s.Find("li.detail_list_el").EachWithBreak(func(i int, s2 *goquery.Selection) bool {
				img := s2.Find("img")
				if img.Size() == 0 {
					hasError = true
					return false
				}
				imgClass := img.AttrOr("class", "")
				r := regexp.MustCompile(`research(\d+)`)
				researchID := utils.DoParseI64(r.FindStringSubmatch(imgClass)[1])
				l := utils.ParseInt(s2.Find("span.fright").Text())
				level := &l
				switch ogame.ID(researchID) {
				case ogame.EspionageTechnology.ID:
					report.EspionageTechnology = level
				case ogame.ComputerTechnology.ID:
					report.ComputerTechnology = level
				case ogame.WeaponsTechnology.ID:
					report.WeaponsTechnology = level
				case ogame.ShieldingTechnology.ID:
					report.ShieldingTechnology = level
				case ogame.ArmourTechnology.ID:
					report.ArmourTechnology = level
				case ogame.EnergyTechnology.ID:
					report.EnergyTechnology = level
				case ogame.HyperspaceTechnology.ID:
					report.HyperspaceTechnology = level
				case ogame.CombustionDrive.ID:
					report.CombustionDrive = level
				case ogame.ImpulseDrive.ID:
					report.ImpulseDrive = level
				case ogame.HyperspaceDrive.ID:
					report.HyperspaceDrive = level
				case ogame.LaserTechnology.ID:
					report.LaserTechnology = level
				case ogame.IonTechnology.ID:
					report.IonTechnology = level
				case ogame.PlasmaTechnology.ID:
					report.PlasmaTechnology = level
				case ogame.IntergalacticResearchNetwork.ID:
					report.IntergalacticResearchNetwork = level
				case ogame.Astrophysics.ID:
					report.Astrophysics = level
				case ogame.GravitonTechnology.ID:
					report.GravitonTechnology = level
				}
				return true
			})
		} else if dataType == "ships" {
			report.HasFleetInformation = s.Find("li.detail_list_fail").Size() == 0
			s.Find("li.detail_list_el").EachWithBreak(func(i int, s2 *goquery.Selection) bool {
				img := s2.Find("img")
				if img.Size() == 0 {
					hasError = true
					return false
				}
				imgClass := img.AttrOr("class", "")
				r := regexp.MustCompile(`tech(\d+)`)
				shipID := utils.DoParseI64(r.FindStringSubmatch(imgClass)[1])
				l := utils.ParseInt(s2.Find("span.fright").Text())
				level := &l
				switch ogame.ID(shipID) {
				case ogame.SmallCargo.ID:
					report.SmallCargo = level
				case ogame.LargeCargo.ID:
					report.LargeCargo = level
				case ogame.LightFighter.ID:
					report.LightFighter = level
				case ogame.HeavyFighter.ID:
					report.HeavyFighter = level
				case ogame.Cruiser.ID:
					report.Cruiser = level
				case ogame.Battleship.ID:
					report.Battleship = level
				case ogame.ColonyShip.ID:
					report.ColonyShip = level
				case ogame.Recycler.ID:
					report.Recycler = level
				case ogame.EspionageProbe.ID:
					report.EspionageProbe = level
				case ogame.Bomber.ID:
					report.Bomber = level
				case ogame.SolarSatellite.ID:
					report.SolarSatellite = level
				case ogame.Destroyer.ID:
					report.Destroyer = level
				case ogame.Deathstar.ID:
					report.Deathstar = level
				case ogame.Battlecruiser.ID:
					report.Battlecruiser = level
				case ogame.Crawler.ID:
					report.Crawler = level
				case ogame.Reaper.ID:
					report.Reaper = level
				case ogame.Pathfinder.ID:
					report.Pathfinder = level
				}
				return true
			})
		} else if dataType == "defense" {
			report.HasDefensesInformation = s.Find("li.detail_list_fail").Size() == 0
			s.Find("li.detail_list_el").EachWithBreak(func(i int, s2 *goquery.Selection) bool {
				img := s2.Find("img")
				if img.Size() == 0 {
					hasError = true
					return false
				}
				imgClass := img.AttrOr("class", "")
				r := regexp.MustCompile(`defense(\d+)`)
				defenceID := utils.DoParseI64(r.FindStringSubmatch(imgClass)[1])
				l := utils.ParseInt(s2.Find("span.fright").Text())
				level := &l
				switch ogame.ID(defenceID) {
				case ogame.RocketLauncher.ID:
					report.RocketLauncher = level
				case ogame.LightLaser.ID:
					report.LightLaser = level
				case ogame.HeavyLaser.ID:
					report.HeavyLaser = level
				case ogame.GaussCannon.ID:
					report.GaussCannon = level
				case ogame.IonCannon.ID:
					report.IonCannon = level
				case ogame.PlasmaTurret.ID:
					report.PlasmaTurret = level
				case ogame.SmallShieldDome.ID:
					report.SmallShieldDome = level
				case ogame.LargeShieldDome.ID:
					report.LargeShieldDome = level
				case ogame.AntiBallisticMissiles.ID:
					report.AntiBallisticMissiles = level
				case ogame.InterplanetaryMissiles.ID:
					report.InterplanetaryMissiles = level
				}
				return true
			})
		}
	})
	if hasError {
		return report, ogame.ErrDeactivateHidePictures
	}
	return report, nil
}

func ExtractCancelInfos(pageHTML []byte, linkVarName, fnName string, tableIdx int) (token string, id, listID int64, err error) {
	r1 := regexp.MustCompile(linkVarName + `[^?]+\?page=ingame&component=overview&modus=2&token=(\w+)&action=cancel`)
	m1 := r1.FindSubmatch(pageHTML)
	if len(m1) < 2 {
		return "", 0, 0, errors.New("unable to find token")
	}
	token = string(m1[1])
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	t := doc.Find("table.construction").Eq(tableIdx)
	a, _ := t.Find("a.abortNow").Attr("onclick")
	r := regexp.MustCompile(fnName + `\((\d+),\s?(\d+),`)
	m := r.FindStringSubmatch(a)
	if len(m) < 3 {
		return "", 0, 0, errors.New("unable to find id/listid")
	}
	id = utils.DoParseI64(m[1])
	listID = utils.DoParseI64(m[2])
	return
}

func extractCancelBuildingInfosV7(pageHTML []byte) (token string, techID, listID int64, err error) {
	return ExtractCancelInfos(pageHTML, "cancelLinkbuilding", "cancelbuilding", 0)
}

func extractCancelResearchInfosV7(pageHTML []byte) (token string, techID, listID int64, err error) {
	return ExtractCancelInfos(pageHTML, "cancelLinkresearch", "cancelresearch", 1)
}

func extractResourceSettingsFromDocV7(doc *goquery.Document) (ogame.ResourceSettings, error) {
	bodyID := v6.ExtractBodyIDFromDocV6(doc)
	if bodyID == "overview" {
		return ogame.ResourceSettings{}, ogame.ErrInvalidPlanetID
	}
	vals := make([]int64, 0)
	doc.Find("option").Each(func(i int, s *goquery.Selection) {
		_, selectedExists := s.Attr("selected")
		if selectedExists {
			a, _ := s.Attr("value")
			val := utils.DoParseI64(a)
			vals = append(vals, val)
		}
	})
	if len(vals) != 7 {
		return ogame.ResourceSettings{}, errors.New("failed to find all resource settings")
	}

	res := ogame.ResourceSettings{}
	res.MetalMine = vals[0]
	res.CrystalMine = vals[1]
	res.DeuteriumSynthesizer = vals[2]
	res.SolarPlant = vals[3]
	res.FusionReactor = vals[4]
	res.SolarSatellite = vals[5]
	res.Crawler = vals[6]

	return res, nil
}

func extractOverviewProductionFromDocV7(doc *goquery.Document) ([]ogame.Quantifiable, error) {
	res := make([]ogame.Quantifiable, 0)
	active := doc.Find("table.construction").Eq(2)
	href, _ := active.Find("td a").Attr("href")
	m := regexp.MustCompile(`openTech=(\d+)`).FindStringSubmatch(href)
	if len(m) == 0 {
		return []ogame.Quantifiable{}, nil
	}
	idInt := utils.DoParseI64(m[1])
	activeID := ogame.ID(idInt)
	activeNbr := utils.DoParseI64(active.Find("div.shipSumCount").Text())
	res = append(res, ogame.Quantifiable{ID: activeID, Nbr: activeNbr})
	active.Parent().Find("table.queue td").Each(func(i int, s *goquery.Selection) {
		img := s.Find("img")
		alt := img.AttrOr("alt", "")
		m := regexp.MustCompile(`techId_(\d+)`).FindStringSubmatch(alt)
		if len(m) == 0 {
			return
		}
		idInt := utils.DoParseI64(m[1])
		activeID := ogame.ID(idInt)
		activeNbr := utils.ParseInt(s.Text())
		res = append(res, ogame.Quantifiable{ID: activeID, Nbr: activeNbr})
	})
	return res, nil
}

func extractOverviewShipSumCountdownFromBytesV7(pageHTML []byte) int64 {
	var shipSumCountdown int64
	shipSumCountdownMatch := regexp.MustCompile(`var restTimeship2 = (\d+);`).FindSubmatch(pageHTML)
	if len(shipSumCountdownMatch) > 0 {
		shipSumCountdown = int64(utils.ToInt(shipSumCountdownMatch[1]))
	}
	return shipSumCountdown
}

func extractCharacterClassFromDocV7(doc *goquery.Document) (ogame.CharacterClass, error) {
	characterClassDiv := doc.Find("div#characterclass a div")
	if characterClassDiv.HasClass("miner") {
		return ogame.Collector, nil
	} else if characterClassDiv.HasClass("warrior") {
		return ogame.General, nil
	} else if characterClassDiv.HasClass("explorer") {
		return ogame.Discoverer, nil
	}
	return 0, errors.New("character class not found")
}

func extractExpeditionMessagesFromDocV7(doc *goquery.Document, location *time.Location) ([]ogame.ExpeditionMessage, int64, error) {
	msgs := make([]ogame.ExpeditionMessage, 0)
	nbPage := utils.DoParseI64(doc.Find("ul.pagination li").Last().AttrOr("data-page", "1"))
	doc.Find("li.msg").Each(func(i int, s *goquery.Selection) {
		if idStr, exists := s.Attr("data-msg-id"); exists {
			if id, err := utils.ParseI64(idStr); err == nil {
				msg := ogame.ExpeditionMessage{ID: id}
				msg.CreatedAt, _ = time.ParseInLocation("02.01.2006 15:04:05", s.Find(".msg_date").Text(), location)
				msg.Coordinate = v6.ExtractCoordV6(s.Find(".msg_title a").Text())
				msg.Coordinate.Type = ogame.PlanetType
				msg.Content, _ = s.Find("span.msg_content").Html()
				msg.Content = strings.TrimSpace(msg.Content)
				msgs = append(msgs, msg)
			}
		}
	})
	return msgs, nbPage, nil
}

func extractMarketplaceMessagesFromDocV7(doc *goquery.Document, location *time.Location) ([]ogame.MarketplaceMessage, int64, error) {
	msgs := make([]ogame.MarketplaceMessage, 0)
	tab := utils.DoParseI64(doc.Find("ul.pagination li").Last().AttrOr("data-tab", ""))
	nbPage := utils.DoParseI64(doc.Find("ul.pagination li").Last().AttrOr("data-page", "1"))
	doc.Find("li.msg").Each(func(i int, s *goquery.Selection) {
		if idStr, exists := s.Attr("data-msg-id"); exists {
			if id, err := utils.ParseI64(idStr); err == nil {
				href := s.Find("a.js_actionCollect").AttrOr("href", "")
				m := regexp.MustCompile(`token=([^&]+)`).FindStringSubmatch(href)
				var token string
				var marketTransactionID int64
				if len(m) == 2 {
					token = m[1]
				}
				m = regexp.MustCompile(`marketTransactionId=([^&]+)`).FindStringSubmatch(href)
				if len(m) == 2 {
					marketTransactionIDStr := m[1]
					marketTransactionID = utils.DoParseI64(marketTransactionIDStr)
				}
				msg := ogame.MarketplaceMessage{ID: id}
				msg.Type = tab
				msg.CreatedAt, _ = time.ParseInLocation("02.01.2006 15:04:05", s.Find(".msg_date").Text(), location)
				msg.Token = token
				msg.MarketTransactionID = marketTransactionID
				msgs = append(msgs, msg)
			}
		}
	})
	return msgs, nbPage, nil
}
