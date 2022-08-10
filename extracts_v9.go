package ogame

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func extractOverviewProductionFromDocV9(doc *goquery.Document) ([]Quantifiable, error) {
	res := make([]Quantifiable, 0)
	active := doc.Find("table.construction").Eq(4)
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
		activeID := ShipName2ID(alt)
		if !activeID.IsSet() {
			activeID = DefenceName2ID(alt)
			if !activeID.IsSet() {
				return
			}
		}
		activeNbr := ParseInt(s.Text())
		res = append(res, Quantifiable{ID: activeID, Nbr: activeNbr})
	})
	return res, nil
}

func extractResourcesFromDocV9(doc *goquery.Document) Resources {
	return extractResourcesDetailsFromFullPageFromDocV9(doc).Available()
}

func extractResourcesDetailsFromFullPageFromDocV9(doc *goquery.Document) ResourcesDetails {
	out := ResourcesDetails{}
	metalDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(doc.Find("div#metal_box").AttrOr("title", "")))
	crystalDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(doc.Find("div#crystal_box").AttrOr("title", "")))
	deuteriumDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(doc.Find("div#deuterium_box").AttrOr("title", "")))
	energyDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(doc.Find("div#energy_box").AttrOr("title", "")))
	darkmatterDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(doc.Find("div#darkmatter_box").AttrOr("title", "")))
	out.Metal.Available = ParseInt(metalDoc.Find("table tr").Eq(0).Find("td").Eq(0).Text())
	out.Metal.StorageCapacity = ParseInt(metalDoc.Find("table tr").Eq(1).Find("td").Eq(0).Text())
	out.Metal.CurrentProduction = ParseInt(metalDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	out.Crystal.Available = ParseInt(crystalDoc.Find("table tr").Eq(0).Find("td").Eq(0).Text())
	out.Crystal.StorageCapacity = ParseInt(crystalDoc.Find("table tr").Eq(1).Find("td").Eq(0).Text())
	out.Crystal.CurrentProduction = ParseInt(crystalDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	out.Deuterium.Available = ParseInt(deuteriumDoc.Find("table tr").Eq(0).Find("td").Eq(0).Text())
	out.Deuterium.StorageCapacity = ParseInt(deuteriumDoc.Find("table tr").Eq(1).Find("td").Eq(0).Text())
	out.Deuterium.CurrentProduction = ParseInt(deuteriumDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	out.Energy.Available = ParseInt(energyDoc.Find("table tr").Eq(0).Find("td").Eq(0).Text())
	out.Energy.CurrentProduction = ParseInt(energyDoc.Find("table tr").Eq(1).Find("td").Eq(0).Text())
	out.Energy.Consumption = ParseInt(energyDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	out.Darkmatter.Available = ParseInt(darkmatterDoc.Find("table tr").Eq(0).Find("td").Eq(0).Text())
	out.Darkmatter.Purchased = ParseInt(darkmatterDoc.Find("table tr").Eq(1).Find("td").Eq(0).Text())
	out.Darkmatter.Found = ParseInt(darkmatterDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	return out
}

func extractEspionageReportFromDocV9(doc *goquery.Document, location *time.Location) (EspionageReport, error) {
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
	characterClassStr := doc.Find("div.detail_txt").Eq(1).Find("span span").First().Text()
	characterClassStr = strings.TrimSpace(characterClassStr)
	report.CharacterClass = getCharacterClass(characterClassStr)

	report.AllianceClass = NoAllianceClass
	allianceClassSpan := doc.Find("div.detail_txt").Eq(2).Find("span.alliance_class")
	if allianceClassSpan.HasClass("trader") {
		report.AllianceClass = Trader
	} else if allianceClassSpan.HasClass("warrior") {
		report.AllianceClass = Warrior
	} else if allianceClassSpan.HasClass("researcher") {
		report.AllianceClass = Researcher
	}

	// Bandit, Starlord
	banditstarlord := doc.Find("div.detail_txt").First().Find("span")
	if banditstarlord.HasClass("honorRank") {
		report.IsBandit = banditstarlord.HasClass("rank_bandit1") || banditstarlord.HasClass("rank_bandit2") || banditstarlord.HasClass("rank_bandit3")
		report.IsStarlord = banditstarlord.HasClass("rank_starlord1") || banditstarlord.HasClass("rank_starlord2") || banditstarlord.HasClass("rank_starlord3")
	}

	honorableFound := doc.Find("div.detail_txt").First().Find("span.status_abbr_honorableTarget")
	report.HonorableTarget = honorableFound.Length() > 0

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
	activity := doc.Find("div.detail_txt").Eq(3).Find("font")
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
	resourcesFound := false
	buildingsFound := false
	doc.Find("ul.detail_list").Each(func(i int, s *goquery.Selection) {
		dataType := s.AttrOr("data-type", "")
		if dataType == "resources" && !resourcesFound {
			resourcesFound = true
			report.Metal = ParseInt(s.Find("li").Eq(0).AttrOr("title", "0"))
			report.Crystal = ParseInt(s.Find("li").Eq(1).AttrOr("title", "0"))
			report.Deuterium = ParseInt(s.Find("li").Eq(2).AttrOr("title", "0"))
			report.Energy = ParseInt(s.Find("li").Eq(3).AttrOr("title", "0"))
		} else if dataType == "buildings" && !buildingsFound {
			buildingsFound = true
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
