package ogame

import (
	"bytes"
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/alaingilbert/clockwork"
)

func getNameV7(doc *goquery.Document, name string) string {
	val := doc.Find("span."+name).First().AttrOr("aria-label", "0")
	return val
}

func getNbrV7(doc *goquery.Document, name string) int64 {
	val, _ := strconv.ParseInt(doc.Find("span."+name+" span.level").First().AttrOr("data-value", "0"), 10, 64)
	return val
}

func getNbrV7Ships(doc *goquery.Document, name string) int64 {
	val, _ := strconv.ParseInt(doc.Find("span."+name+" span.amount").First().AttrOr("data-value", "0"), 10, 64)
	return val
}

func extractPremiumTokenV7(pageHTML []byte, days int64) (token string, err error) {
	rgx := regexp.MustCompile(`\?page=premium&buynow=1&type=\d&days=` + strconv.FormatInt(days, 10) + `&token=(\w+)`)
	m := rgx.FindSubmatch(pageHTML)
	if len(m) < 2 {
		return "", errors.New("unable to find token")
	}
	token = string(m[1])
	return
}

func extractResourcesDetailsFromFullPageFromDocV7(doc *goquery.Document) ResourcesDetails {
	out := ResourcesDetails{}
	out.Metal.Available = ParseInt(strings.Split(doc.Find("span#resources_metal").AttrOr("data-raw", "0"), ".")[0])
	out.Crystal.Available = ParseInt(strings.Split(doc.Find("span#resources_crystal").AttrOr("data-raw", "0"), ".")[0])
	out.Deuterium.Available = ParseInt(strings.Split(doc.Find("span#resources_deuterium").AttrOr("data-raw", "0"), ".")[0])
	out.Energy.Available = ParseInt(doc.Find("span#resources_energy").Text())
	out.Darkmatter.Available = ParseInt(doc.Find("span#resources_darkmatter").Text())
	metalDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(doc.Find("li#metal_box").AttrOr("title", "")))
	crystalDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(doc.Find("li#crystal_box").AttrOr("title", "")))
	deuteriumDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(doc.Find("li#deuterium_box").AttrOr("title", "")))
	energyDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(doc.Find("li#energy_box").AttrOr("title", "")))
	darkmatterDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(doc.Find("li#darkmatter_box").AttrOr("title", "")))
	out.Metal.StorageCapacity = ParseInt(metalDoc.Find("table tr").Eq(1).Find("td").Eq(0).Text())
	out.Metal.CurrentProduction = ParseInt(metalDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	out.Crystal.StorageCapacity = ParseInt(crystalDoc.Find("table tr").Eq(1).Find("td").Eq(0).Text())
	out.Crystal.CurrentProduction = ParseInt(crystalDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	out.Deuterium.StorageCapacity = ParseInt(deuteriumDoc.Find("table tr").Eq(1).Find("td").Eq(0).Text())
	out.Deuterium.CurrentProduction = ParseInt(deuteriumDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	out.Energy.CurrentProduction = ParseInt(energyDoc.Find("table tr").Eq(1).Find("td").Eq(0).Text())
	out.Energy.Consumption = ParseInt(energyDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	out.Darkmatter.Purchased = ParseInt(darkmatterDoc.Find("table tr").Eq(1).Find("td").Eq(0).Text())
	out.Darkmatter.Found = ParseInt(darkmatterDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	return out
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
	res.SolarSatellite = getNbrV7Ships(doc, "solarSatellite")
	res.MetalStorage = getNbrV7(doc, "metalStorage")
	res.CrystalStorage = getNbrV7(doc, "crystalStorage")
	res.DeuteriumTank = getNbrV7(doc, "deuteriumStorage")
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

func extractConstructionsV7(pageHTML []byte, clock clockwork.Clock) (buildingID ID, buildingCountdown int64, researchID ID, researchCountdown int64) {
	buildingCountdownMatch := regexp.MustCompile(`var restTimebuilding = (\d+) -`).FindSubmatch(pageHTML)
	if len(buildingCountdownMatch) > 0 {
		buildingCountdown = int64(toInt(buildingCountdownMatch[1])) - clock.Now().Unix()
		buildingIDInt := toInt(regexp.MustCompile(`onclick="cancelbuilding\((\d+),`).FindSubmatch(pageHTML)[1])
		buildingID = ID(buildingIDInt)
	}
	researchCountdownMatch := regexp.MustCompile(`var restTimeresearch = (\d+) -`).FindSubmatch(pageHTML)
	if len(researchCountdownMatch) > 0 {
		researchCountdown = int64(toInt(researchCountdownMatch[1])) - clock.Now().Unix()
		researchIDInt := toInt(regexp.MustCompile(`onclick="cancelresearch\((\d+),`).FindSubmatch(pageHTML)[1])
		researchID = ID(researchIDInt)
	}
	return
}

func extractIPMFromDocV7(doc *goquery.Document) (duration, max int64, token string) {
	duration, _ = strconv.ParseInt(doc.Find("span#timer").AttrOr("data-duration", "0"), 10, 64)
	max, _ = strconv.ParseInt(doc.Find("input[name=missileCount]").AttrOr("data-max", "0"), 10, 64)
	token = doc.Find("input[name=token]").AttrOr("value", "")
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
		ID     int64 `json:"id"`
		Number int64 `json:"number"`
	}
	if err := json.Unmarshal([]byte(m), &res); err != nil {
		return
	}
	for _, obj := range res {
		s.Set(ID(obj.ID), obj.Number)
	}
	return
}

func extractCombatReportMessagesFromDocV7(doc *goquery.Document) ([]CombatReportSummary, int64) {
	msgs := make([]CombatReportSummary, 0)
	nbPage, _ := strconv.ParseInt(doc.Find("ul.pagination li").Last().AttrOr("data-page", "1"), 10, 64)
	doc.Find("li.msg").Each(func(i int, s *goquery.Selection) {
		if idStr, exists := s.Attr("data-msg-id"); exists {
			if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
				report := CombatReportSummary{ID: id}
				report.Destination = extractCoordV6(s.Find("div.msg_head a").Text())
				if s.Find("div.msg_head figure").HasClass("planet") {
					report.Destination.Type = PlanetType
				} else if s.Find("div.msg_head figure").HasClass("moon") {
					report.Destination.Type = MoonType
				} else {
					report.Destination.Type = PlanetType
				}
				apiKeyTitle := s.Find("span.icon_apikey").AttrOr("title", "")
				m := regexp.MustCompile(`'(cr-[^']+)'`).FindStringSubmatch(apiKeyTitle)
				if len(m) == 2 {
					report.APIKey = m[1]
				}
				resTitle := s.Find("span.msg_content div.combatLeftSide span").Eq(1).AttrOr("title", "")
				m = regexp.MustCompile(`([\d.,]+)<br/>[^\d]*([\d.,]+)<br/>[^\d]*([\d.,]+)`).FindStringSubmatch(resTitle)
				if len(m) == 4 {
					report.Metal = ParseInt(m[1])
					report.Crystal = ParseInt(m[2])
					report.Deuterium = ParseInt(m[3])
				}
				debrisFieldTitle := s.Find("span.msg_content div.combatLeftSide span").Eq(2).AttrOr("title", "0")
				report.DebrisField = ParseInt(debrisFieldTitle)
				resText := s.Find("span.msg_content div.combatLeftSide span").Eq(1).Text()
				m = regexp.MustCompile(`[\d.,]+[^\d]*([\d.,]+)`).FindStringSubmatch(resText)
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
				galaxy, _ := strconv.ParseInt(m[1], 10, 64)
				system, _ := strconv.ParseInt(m[2], 10, 64)
				position, _ := strconv.ParseInt(m[3], 10, 64)
				planetType, _ := strconv.ParseInt(m[4], 10, 64)
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

func extractEspionageReportFromDocV7(doc *goquery.Document, location *time.Location) (EspionageReport, error) {
	report := EspionageReport{}
	report.ID, _ = strconv.ParseInt(doc.Find("div.detail_msg").AttrOr("data-msg-id", "0"), 10, 64)
	spanLink := doc.Find("span.msg_title a").First()
	txt := spanLink.Text()
	figure := spanLink.Find("figure").First()
	r := regexp.MustCompile(`([^\[]+) \[(\d+):(\d+):(\d+)]`)
	m := r.FindStringSubmatch(txt)
	if len(m) == 5 {
		report.Coordinate.Galaxy, _ = strconv.ParseInt(m[2], 10, 64)
		report.Coordinate.System, _ = strconv.ParseInt(m[3], 10, 64)
		report.Coordinate.Position, _ = strconv.ParseInt(m[4], 10, 64)
	} else {
		return report, errors.New("failed to extract coordinate")
	}
	if figure.HasClass("planet") {
		report.Coordinate.Type = PlanetType
	} else if figure.HasClass("moon") {
		report.Coordinate.Type = MoonType
	}
	messageType := Report
	if doc.Find("span.espionageDefText").Size() > 0 {
		messageType = Action
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
		report.LastActivity = ParseInt(activity.Text())
	}

	// CounterEspionage
	ceTxt := doc.Find("div.detail_txt").Eq(2).Text()
	m1 := regexp.MustCompile(`(\d+)%`).FindStringSubmatch(ceTxt)
	if len(m1) == 2 {
		report.CounterEspionage, _ = strconv.ParseInt(m1[1], 10, 64)
	}

	hasError := false
	doc.Find("ul.detail_list").Each(func(i int, s *goquery.Selection) {
		dataType := s.AttrOr("data-type", "")
		if dataType == "resources" {
			report.Metal = ParseInt(s.Find("li").Eq(0).AttrOr("title", "0"))
			report.Crystal = ParseInt(s.Find("li").Eq(1).AttrOr("title", "0"))
			report.Deuterium = ParseInt(s.Find("li").Eq(2).AttrOr("title", "0"))
			report.Energy = ParseInt(s.Find("li").Eq(3).AttrOr("title", "0"))
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
				buildingID, _ := strconv.ParseInt(r.FindStringSubmatch(imgClass)[1], 10, 64)
				l := ParseInt(s2.Find("span.fright").Text())
				level := &l
				switch ID(buildingID) {
				case MetalMine.ID:
					report.MetalMine = level
				case CrystalMine.ID:
					report.CrystalMine = level
				case DeuteriumSynthesizer.ID:
					report.DeuteriumSynthesizer = level
				case SolarPlant.ID:
					report.SolarPlant = level
				case FusionReactor.ID:
					report.FusionReactor = level
				case MetalStorage.ID:
					report.MetalStorage = level
				case CrystalStorage.ID:
					report.CrystalStorage = level
				case DeuteriumTank.ID:
					report.DeuteriumTank = level
				case AllianceDepot.ID:
					report.AllianceDepot = level
				case RoboticsFactory.ID:
					report.RoboticsFactory = level
				case Shipyard.ID:
					report.Shipyard = level
				case ResearchLab.ID:
					report.ResearchLab = level
				case MissileSilo.ID:
					report.MissileSilo = level
				case NaniteFactory.ID:
					report.NaniteFactory = level
				case Terraformer.ID:
					report.Terraformer = level
				case SpaceDock.ID:
					report.SpaceDock = level
				case LunarBase.ID:
					report.LunarBase = level
				case SensorPhalanx.ID:
					report.SensorPhalanx = level
				case JumpGate.ID:
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
				researchID, _ := strconv.ParseInt(r.FindStringSubmatch(imgClass)[1], 10, 64)
				l := ParseInt(s2.Find("span.fright").Text())
				level := &l
				switch ID(researchID) {
				case EspionageTechnology.ID:
					report.EspionageTechnology = level
				case ComputerTechnology.ID:
					report.ComputerTechnology = level
				case WeaponsTechnology.ID:
					report.WeaponsTechnology = level
				case ShieldingTechnology.ID:
					report.ShieldingTechnology = level
				case ArmourTechnology.ID:
					report.ArmourTechnology = level
				case EnergyTechnology.ID:
					report.EnergyTechnology = level
				case HyperspaceTechnology.ID:
					report.HyperspaceTechnology = level
				case CombustionDrive.ID:
					report.CombustionDrive = level
				case ImpulseDrive.ID:
					report.ImpulseDrive = level
				case HyperspaceDrive.ID:
					report.HyperspaceDrive = level
				case LaserTechnology.ID:
					report.LaserTechnology = level
				case IonTechnology.ID:
					report.IonTechnology = level
				case PlasmaTechnology.ID:
					report.PlasmaTechnology = level
				case IntergalacticResearchNetwork.ID:
					report.IntergalacticResearchNetwork = level
				case Astrophysics.ID:
					report.Astrophysics = level
				case GravitonTechnology.ID:
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
				shipID, _ := strconv.ParseInt(r.FindStringSubmatch(imgClass)[1], 10, 64)
				l := ParseInt(s2.Find("span.fright").Text())
				level := &l
				switch ID(shipID) {
				case SmallCargo.ID:
					report.SmallCargo = level
				case LargeCargo.ID:
					report.LargeCargo = level
				case LightFighter.ID:
					report.LightFighter = level
				case HeavyFighter.ID:
					report.HeavyFighter = level
				case Cruiser.ID:
					report.Cruiser = level
				case Battleship.ID:
					report.Battleship = level
				case ColonyShip.ID:
					report.ColonyShip = level
				case Recycler.ID:
					report.Recycler = level
				case EspionageProbe.ID:
					report.EspionageProbe = level
				case Bomber.ID:
					report.Bomber = level
				case SolarSatellite.ID:
					report.SolarSatellite = level
				case Destroyer.ID:
					report.Destroyer = level
				case Deathstar.ID:
					report.Deathstar = level
				case Battlecruiser.ID:
					report.Battlecruiser = level
				case Crawler.ID:
					report.Crawler = level
				case Reaper.ID:
					report.Reaper = level
				case Pathfinder.ID:
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
				defenceID, _ := strconv.ParseInt(r.FindStringSubmatch(imgClass)[1], 10, 64)
				l := ParseInt(s2.Find("span.fright").Text())
				level := &l
				switch ID(defenceID) {
				case RocketLauncher.ID:
					report.RocketLauncher = level
				case LightLaser.ID:
					report.LightLaser = level
				case HeavyLaser.ID:
					report.HeavyLaser = level
				case GaussCannon.ID:
					report.GaussCannon = level
				case IonCannon.ID:
					report.IonCannon = level
				case PlasmaTurret.ID:
					report.PlasmaTurret = level
				case SmallShieldDome.ID:
					report.SmallShieldDome = level
				case LargeShieldDome.ID:
					report.LargeShieldDome = level
				case AntiBallisticMissiles.ID:
					report.AntiBallisticMissiles = level
				case InterplanetaryMissiles.ID:
					report.InterplanetaryMissiles = level
				}
				return true
			})
		}
	})
	if hasError {
		return report, ErrDeactivateHidePictures
	}
	return report, nil
}

func extractCancelBuildingInfosV7(pageHTML []byte) (token string, techID, listID int64, err error) {
	r1 := regexp.MustCompile(`cancelLinkbuilding[^?]+\?page=ingame&component=overview&modus=2&token=(\w+)&action=cancel`)
	m1 := r1.FindSubmatch(pageHTML)
	if len(m1) < 2 {
		return "", 0, 0, errors.New("unable to find token")
	}
	token = string(m1[1])
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	t := doc.Find("table.construction").Eq(0)
	a, _ := t.Find("a.abortNow").Attr("onclick")
	r := regexp.MustCompile(`cancelbuilding\((\d+),\s?(\d+),`)
	m := r.FindStringSubmatch(a)
	if len(m) < 3 {
		return "", 0, 0, errors.New("unable to find techid/listid")
	}
	techID, _ = strconv.ParseInt(m[1], 10, 64)
	listID, _ = strconv.ParseInt(m[2], 10, 64)
	return
}

func extractCancelResearchInfosV7(pageHTML []byte) (token string, techID, listID int64, err error) {
	r1 := regexp.MustCompile(`cancelLinkresearch[^?]+\?page=ingame&component=overview&modus=2&token=(\w+)&action=cancel`)
	m1 := r1.FindSubmatch(pageHTML)
	if len(m1) < 2 {
		return "", 0, 0, errors.New("unable to find token")
	}
	token = string(m1[1])
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	t := doc.Find("table.construction").Eq(1)
	a, _ := t.Find("a.abortNow").Attr("onclick")
	r := regexp.MustCompile(`cancelresearch\((\d+),\s?(\d+),`)
	m := r.FindStringSubmatch(a)
	if len(m) < 3 {
		return "", 0, 0, errors.New("unable to find techid/listid")
	}
	techID, _ = strconv.ParseInt(m[1], 10, 64)
	listID, _ = strconv.ParseInt(m[2], 10, 64)
	return
}

func extractResourceSettingsFromDocV7(doc *goquery.Document) (ResourceSettings, error) {
	bodyID := extractBodyIDFromDocV6(doc)
	if bodyID == "overview" {
		return ResourceSettings{}, ErrInvalidPlanetID
	}
	vals := make([]int64, 0)
	doc.Find("option").Each(func(i int, s *goquery.Selection) {
		_, selectedExists := s.Attr("selected")
		if selectedExists {
			a, _ := s.Attr("value")
			val, _ := strconv.ParseInt(a, 10, 64)
			vals = append(vals, val)
		}
	})
	if len(vals) != 7 {
		return ResourceSettings{}, errors.New("failed to find all resource settings")
	}

	res := ResourceSettings{}
	res.MetalMine = vals[0]
	res.CrystalMine = vals[1]
	res.DeuteriumSynthesizer = vals[2]
	res.SolarPlant = vals[3]
	res.FusionReactor = vals[4]
	res.SolarSatellite = vals[5]
	res.Crawler = vals[6]

	return res, nil
}

func extractOverviewProductionFromDocV7(doc *goquery.Document) ([]Quantifiable, error) {
	res := make([]Quantifiable, 0)
	active := doc.Find("table.construction").Eq(2)
	href, _ := active.Find("td a").Attr("href")
	m := regexp.MustCompile(`openTech=(\d+)`).FindStringSubmatch(href)
	if len(m) == 0 {
		return []Quantifiable{}, nil
	}
	idInt, _ := strconv.ParseInt(m[1], 10, 64)
	activeID := ID(idInt)
	activeNbr, _ := strconv.ParseInt(active.Find("div.shipSumCount").Text(), 10, 64)
	res = append(res, Quantifiable{ID: activeID, Nbr: activeNbr})
	active.Parent().Find("table.queue td").Each(func(i int, s *goquery.Selection) {
		img := s.Find("img")
		alt := img.AttrOr("alt", "")
		m := regexp.MustCompile(`techId_(\d+)`).FindStringSubmatch(alt)
		if len(m) == 0 {
			return
		}
		idInt, _ := strconv.ParseInt(m[1], 10, 64)
		activeID := ID(idInt)
		activeNbr := ParseInt(s.Text())
		res = append(res, Quantifiable{ID: activeID, Nbr: activeNbr})
	})
	return res, nil
}

func extractOverviewShipSumCountdownFromBytesV7(pageHTML []byte) int64 {
	var shipSumCountdown int64
	shipSumCountdownMatch := regexp.MustCompile(`var restTimeship2 = (\d+);`).FindSubmatch(pageHTML)
	if len(shipSumCountdownMatch) > 0 {
		shipSumCountdown = int64(toInt(shipSumCountdownMatch[1]))
	}
	return shipSumCountdown
}

func extractCharacterClassFromDocV7(doc *goquery.Document) (CharacterClass, error) {
	characterClassDiv := doc.Find("div#characterclass a div")
	if characterClassDiv.HasClass("miner") {
		return Collector, nil
	} else if characterClassDiv.HasClass("warrior") {
		return General, nil
	} else if characterClassDiv.HasClass("explorer") {
		return Discoverer, nil
	}
	return 0, errors.New("character class not found")
}

func extractExpeditionMessagesFromDocV7(doc *goquery.Document, location *time.Location) ([]ExpeditionMessage, int64, error) {
	msgs := make([]ExpeditionMessage, 0)
	nbPage, _ := strconv.ParseInt(doc.Find("ul.pagination li").Last().AttrOr("data-page", "1"), 10, 64)
	doc.Find("li.msg").Each(func(i int, s *goquery.Selection) {
		if idStr, exists := s.Attr("data-msg-id"); exists {
			if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
				msg := ExpeditionMessage{ID: id}
				msg.CreatedAt, _ = time.ParseInLocation("02.01.2006 15:04:05", s.Find(".msg_date").Text(), location)
				msg.Coordinate = extractCoordV6(s.Find(".msg_title a").Text())
				msg.Coordinate.Type = PlanetType
				msg.Content, _ = s.Find("span.msg_content").Html()
				msg.Content = strings.TrimSpace(msg.Content)
				msgs = append(msgs, msg)
			}
		}
	})
	return msgs, nbPage, nil
}

func extractMarketplaceMessagesFromDocV7(doc *goquery.Document, location *time.Location) ([]MarketplaceMessage, int64, error) {
	msgs := make([]MarketplaceMessage, 0)
	tab, _ := strconv.ParseInt(doc.Find("ul.pagination li").Last().AttrOr("data-tab", ""), 10, 64)
	nbPage, _ := strconv.ParseInt(doc.Find("ul.pagination li").Last().AttrOr("data-page", "1"), 10, 64)
	doc.Find("li.msg").Each(func(i int, s *goquery.Selection) {
		if idStr, exists := s.Attr("data-msg-id"); exists {
			if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
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
					marketTransactionID, _ = strconv.ParseInt(marketTransactionIDStr, 10, 64)
				}
				msg := MarketplaceMessage{ID: id}
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

