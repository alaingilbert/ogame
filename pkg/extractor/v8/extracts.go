package v8

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	v71 "github.com/alaingilbert/ogame/pkg/extractor/v71"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/utils"
)

func extractIsInVacationFromDoc(doc *goquery.Document) bool {
	var isVacation bool
	for _, s := range doc.Find("div#advice-bar a").EachIter() {
		href := s.AttrOr("href", "")
		if href == "" {
			continue
		}
		u, _ := url.Parse(href)
		q := u.Query()
		if q.Get("component") == "preferences" && q.Get("selectedTab") == "3" && q.Get("openGroup") == "0" {
			isVacation = true
			continue
		}
	}
	return isVacation
}

func extractEspionageReportFromDoc(doc *goquery.Document, location *time.Location) (ogame.EspionageReport, error) {
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
	characterClassStr := doc.Find("div.detail_txt").Eq(1).Find("span span").First().Text()
	characterClassStr = strings.TrimSpace(characterClassStr)
	report.CharacterClass = v71.GetCharacterClass(characterClassStr)

	report.AllianceClass = ogame.NoAllianceClass
	allianceClassSpan := doc.Find("div.detail_txt").Eq(2).Find("span.alliance_class")
	if allianceClassSpan.HasClass("trader") {
		report.AllianceClass = ogame.Trader
	} else if allianceClassSpan.HasClass("warrior") {
		report.AllianceClass = ogame.Warrior
	} else if allianceClassSpan.HasClass("researcher") {
		report.AllianceClass = ogame.Researcher
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
		report.LastActivity = utils.ParseInt(activity.Text())
	}

	// CounterEspionage
	ceTxt := doc.Find("div.detail_txt").Eq(2).Text()
	m1 := regexp.MustCompile(`(\d+)%`).FindStringSubmatch(ceTxt)
	if len(m1) == 2 {
		report.CounterEspionage = utils.DoParseI64(m1[1])
	}

	hasError := false
	resourcesFound := false
	for _, s := range doc.Find("ul.detail_list").EachIter() {
		dataType := s.AttrOr("data-type", "")
		if dataType == "resources" && !resourcesFound {
			resourcesFound = true
			report.Metal = utils.ParseInt(s.Find("li").Eq(0).AttrOr("title", "0"))
			report.Crystal = utils.ParseInt(s.Find("li").Eq(1).AttrOr("title", "0"))
			report.Deuterium = utils.ParseInt(s.Find("li").Eq(2).AttrOr("title", "0"))
			report.Energy = utils.ParseInt(s.Find("li").Eq(3).AttrOr("title", "0"))
		} else if dataType == "buildings" {
			report.HasBuildingsInformation = s.Find("li.detail_list_fail").Size() == 0
			for _, s2 := range s.Find("li.detail_list_el").EachIter() {
				img := s2.Find("img")
				if img.Size() == 0 {
					hasError = true
					break
				}
				imgClass := img.AttrOr("class", "")
				r := regexp.MustCompile(`building(\d+)`)
				buildingID := utils.DoParseI64(r.FindStringSubmatch(imgClass)[1])
				level := utils.Ptr(utils.ParseInt(s2.Find("span.fright").Text()))
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
			}
		} else if dataType == "research" {
			report.HasResearchesInformation = s.Find("li.detail_list_fail").Size() == 0
			for _, s2 := range s.Find("li.detail_list_el").EachIter() {
				img := s2.Find("img")
				if img.Size() == 0 {
					hasError = true
					break
				}
				imgClass := img.AttrOr("class", "")
				r := regexp.MustCompile(`research(\d+)`)
				researchID := utils.DoParseI64(r.FindStringSubmatch(imgClass)[1])
				level := utils.Ptr(utils.ParseInt(s2.Find("span.fright").Text()))
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
			}
		} else if dataType == "ships" {
			report.HasFleetInformation = s.Find("li.detail_list_fail").Size() == 0
			for _, s2 := range s.Find("li.detail_list_el").EachIter() {
				img := s2.Find("img")
				if img.Size() == 0 {
					hasError = true
					break
				}
				imgClass := img.AttrOr("class", "")
				r := regexp.MustCompile(`tech(\d+)`)
				shipID := utils.DoParseI64(r.FindStringSubmatch(imgClass)[1])
				level := utils.Ptr(utils.ParseInt(s2.Find("span.fright").Text()))
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
			}
		} else if dataType == "defense" {
			report.HasDefensesInformation = s.Find("li.detail_list_fail").Size() == 0
			for _, s2 := range s.Find("li.detail_list_el").EachIter() {
				img := s2.Find("img")
				if img.Size() == 0 {
					hasError = true
					break
				}
				imgClass := img.AttrOr("class", "")
				r := regexp.MustCompile(`defense(\d+)`)
				defenceID := utils.DoParseI64(r.FindStringSubmatch(imgClass)[1])
				level := utils.Ptr(utils.ParseInt(s2.Find("span.fright").Text()))
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
			}
		}
	}
	if hasError {
		return report, ogame.ErrDeactivateHidePictures
	}
	return report, nil
}
