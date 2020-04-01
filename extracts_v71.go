package ogame

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/alaingilbert/clockwork"
	"golang.org/x/net/html"
)

type resourcesRespV71 struct {
	Resources struct {
		Metal struct {
			Amount  float64 `json:"amount"`
			Storage float64 `json:"storage"`
			Tooltip string  `json:"tooltip"`
		} `json:"metal"`
		Crystal struct {
			Amount  float64 `json:"amount"`
			Storage float64 `json:"storage"`
			Tooltip string  `json:"tooltip"`
		} `json:"crystal"`
		Deuterium struct {
			Amount  float64 `json:"amount"`
			Storage float64 `json:"storage"`
			Tooltip string  `json:"tooltip"`
		} `json:"deuterium"`
		Energy struct {
			Amount  float64 `json:"amount"`
			Tooltip string  `json:"tooltip"`
		} `json:"energy"`
		Darkmatter struct {
			Amount  float64 `json:"amount"`
			Tooltip string  `json:"tooltip"`
		} `json:"darkmatter"`
	} `json:"resources"`
	HonorScore int64 `json:"honorScore"`
	Techs      struct {
		Num1 struct {
			TechID     int64 `json:"techId"`
			Production struct {
				Metal     float64 `json:"metal"`
				Crystal   float64 `json:"crystal"`
				Deuterium float64 `json:"deuterium"`
				Energy    float64 `json:"energy"`
			} `json:"production"`
			Consumption struct {
				Metal     float64 `json:"metal"`
				Crystal   float64 `json:"crystal"`
				Deuterium float64 `json:"deuterium"`
				Energy    float64 `json:"energy"`
			} `json:"consumption"`
		} `json:"1"`
		Num2 struct {
			TechID     int64 `json:"techId"`
			Production struct {
				Metal     float64 `json:"metal"`
				Crystal   float64 `json:"crystal"`
				Deuterium float64 `json:"deuterium"`
				Energy    float64 `json:"energy"`
			} `json:"production"`
			Consumption struct {
				Metal     float64 `json:"metal"`
				Crystal   float64 `json:"crystal"`
				Deuterium float64 `json:"deuterium"`
				Energy    float64 `json:"energy"`
			} `json:"consumption"`
		} `json:"2"`
		Num3 struct {
			TechID     int64 `json:"techId"`
			Production struct {
				Metal     float64 `json:"metal"`
				Crystal   float64 `json:"crystal"`
				Deuterium float64 `json:"deuterium"`
				Energy    float64 `json:"energy"`
			} `json:"production"`
			Consumption struct {
				Metal     float64 `json:"metal"`
				Crystal   float64 `json:"crystal"`
				Deuterium float64 `json:"deuterium"`
				Energy    float64 `json:"energy"`
			} `json:"consumption"`
		} `json:"3"`
		Num4 struct {
			TechID     int64 `json:"techId"`
			Production struct {
				Metal     float64 `json:"metal"`
				Crystal   float64 `json:"crystal"`
				Deuterium float64 `json:"deuterium"`
				Energy    float64 `json:"energy"`
			} `json:"production"`
			Consumption struct {
				Metal     float64 `json:"metal"`
				Crystal   float64 `json:"crystal"`
				Deuterium float64 `json:"deuterium"`
				Energy    float64 `json:"energy"`
			} `json:"consumption"`
		} `json:"4"`
		Num12 struct {
			TechID     int64 `json:"techId"`
			Production struct {
				Metal     float64 `json:"metal"`
				Crystal   float64 `json:"crystal"`
				Deuterium float64 `json:"deuterium"`
				Energy    float64 `json:"energy"`
			} `json:"production"`
			Consumption struct {
				Metal     float64 `json:"metal"`
				Crystal   float64 `json:"crystal"`
				Deuterium float64 `json:"deuterium"`
				Energy    float64 `json:"energy"`
			} `json:"consumption"`
		} `json:"12"`
		Num212 struct {
			TechID     int64 `json:"techId"`
			Production struct {
				Metal     float64 `json:"metal"`
				Crystal   float64 `json:"crystal"`
				Deuterium float64 `json:"deuterium"`
				Energy    float64 `json:"energy"`
			} `json:"production"`
			Consumption struct {
				Metal     float64 `json:"metal"`
				Crystal   float64 `json:"crystal"`
				Deuterium float64 `json:"deuterium"`
				Energy    float64 `json:"energy"`
			} `json:"consumption"`
		} `json:"212"`
		Num217 struct {
			TechID     int64 `json:"techId"`
			Production struct {
				Metal     float64 `json:"metal"`
				Crystal   float64 `json:"crystal"`
				Deuterium float64 `json:"deuterium"`
				Energy    float64 `json:"energy"`
			} `json:"production"`
			Consumption struct {
				Metal     float64 `json:"metal"`
				Crystal   float64 `json:"crystal"`
				Deuterium float64 `json:"deuterium"`
				Energy    float64 `json:"energy"`
			} `json:"consumption"`
		} `json:"217"`
	} `json:"techs"`
}

func extractResourcesDetailsV71(pageHTML []byte) (out ResourcesDetails, err error) {
	var res resourcesRespV71
	if err = json.Unmarshal(pageHTML, &res); err != nil {
		if isLogged(pageHTML) {
			return out, ErrInvalidPlanetID
		}
		return
	}
	out.Metal.Available = int64(res.Resources.Metal.Amount)
	out.Metal.StorageCapacity = int64(res.Resources.Metal.Storage)
	out.Crystal.Available = int64(res.Resources.Crystal.Amount)
	out.Crystal.StorageCapacity = int64(res.Resources.Crystal.Storage)
	out.Deuterium.Available = int64(res.Resources.Deuterium.Amount)
	out.Deuterium.StorageCapacity = int64(res.Resources.Deuterium.Storage)
	out.Energy.Available = int64(res.Resources.Energy.Amount)
	out.Darkmatter.Available = int64(res.Resources.Darkmatter.Amount)
	metalDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(res.Resources.Metal.Tooltip))
	crystalDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(res.Resources.Crystal.Tooltip))
	deuteriumDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(res.Resources.Deuterium.Tooltip))
	darkmatterDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(res.Resources.Darkmatter.Tooltip))
	energyDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(res.Resources.Energy.Tooltip))
	out.Metal.CurrentProduction = ParseInt(metalDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	out.Crystal.CurrentProduction = ParseInt(crystalDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	out.Deuterium.CurrentProduction = ParseInt(deuteriumDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	out.Energy.CurrentProduction = ParseInt(energyDoc.Find("table tr").Eq(1).Find("td").Eq(0).Text())
	out.Energy.Consumption = ParseInt(energyDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	out.Darkmatter.Purchased = ParseInt(darkmatterDoc.Find("table tr").Eq(1).Find("td").Eq(0).Text())
	out.Darkmatter.Found = ParseInt(darkmatterDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	return
}

func extractEspionageReportFromDocV71(doc *goquery.Document, location *time.Location) (EspionageReport, error) {
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
	resourcesFound := false
	doc.Find("ul.detail_list").Each(func(i int, s *goquery.Selection) {
		dataType := s.AttrOr("data-type", "")
		if dataType == "resources" && !resourcesFound {
			resourcesFound = true
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

func extractIPMFromDocV71(doc *goquery.Document) (duration, max int64, token string) {
	durationFloat, _ := strconv.ParseFloat(doc.Find("span#timer").AttrOr("data-duration", "0"), 64)
	duration = int64(math.Ceil(durationFloat))
	max, _ = strconv.ParseInt(doc.Find("input#missileCount").AttrOr("data-max", "0"), 10, 64)
	token = doc.Find("input[name=token]").AttrOr("value", "")
	return
}

func extractFacilitiesFromDocV71(doc *goquery.Document) (Facilities, error) {
	res, err := extractFacilitiesFromDocV7(doc)
	if err != nil {
		return Facilities{}, err
	}
	res.LunarBase = getNbrV7(doc, "moonbase")
	return res, nil
}

func extractProductionFromDocV71(doc *goquery.Document) ([]Quantifiable, error) {
	res := make([]Quantifiable, 0)
	active := doc.Find("table.construction")
	href, _ := active.Find("td a").Attr("href")
	m := regexp.MustCompile(`openTech=(\d+)`).FindStringSubmatch(href)
	if len(m) == 0 {
		return []Quantifiable{}, nil
	}
	idInt, _ := strconv.ParseInt(m[1], 10, 64)
	activeID := ID(idInt)
	activeNbr, _ := strconv.ParseInt(active.Find("div.shipSumCount").Text(), 10, 64)
	res = append(res, Quantifiable{ID: activeID, Nbr: activeNbr})
	doc.Find("table.queue td").Each(func(i int, s *goquery.Selection) {
		link := s.Find("img")
		alt := link.AttrOr("alt", "")
		m := regexp.MustCompile(`techId_(\d+)`).FindStringSubmatch(alt)
		if len(m) == 0 {
			return
		}
		itemID, _ := strconv.ParseInt(m[1], 10, 64)
		itemNbr := ParseInt(s.Text())
		res = append(res, Quantifiable{ID: ID(itemID), Nbr: itemNbr})
	})
	return res, nil
}

// Highscore ...
type Highscore struct {
	NbPage   int64
	CurrPage int64
	Category int64 // 1:Player, 2:Alliance
	Type     int64 // 0:Total, 1:Economy, 2:Research, 3:Military, 4:Military Built, 5:Military Destroyed, 6:Military Lost, 7:Honor
	Players  []HighscorePlayer
}

// HighscorePlayer ...
type HighscorePlayer struct {
	Position     int64
	ID           int64
	Name         string
	Score        int64
	AllianceID   int64
	HonourPoints int64
	Homeworld    Coordinate
}

func extractHighscoreFromDocV71(doc *goquery.Document) (out Highscore, err error) {
	s := doc.Selection
	isFullPage := doc.Find("#stat_list_content").Size() == 1
	if isFullPage {
		s = doc.Find("#stat_list_content")
	}

	script := s.Find("script").First().Text()
	m := regexp.MustCompile(`var site = (\d+);`).FindStringSubmatch(script)
	if len(m) != 2 {
		return out, errors.New("failed to find site")
	}
	out.CurrPage, _ = strconv.ParseInt(m[1], 10, 64)

	m = regexp.MustCompile(`var currentCategory = (\d+);`).FindStringSubmatch(script)
	if len(m) != 2 {
		return out, errors.New("failed to find currentCategory")
	}
	out.Category, _ = strconv.ParseInt(m[1], 10, 64)

	m = regexp.MustCompile(`var currentType = (\d+);`).FindStringSubmatch(script)
	if len(m) != 2 {
		return out, errors.New("failed to find currentType")
	}
	out.Type, _ = strconv.ParseInt(m[1], 10, 64)

	changeSiteSize := s.Find("select.changeSite option").Size()
	out.NbPage = MaxInt(int64(changeSiteSize)-1, 0)

	s.Find("#ranks tbody tr").Each(func(i int, s *goquery.Selection) {
		p := HighscorePlayer{}
		p.Position, _ = strconv.ParseInt(s.Find("td.position").Text(), 10, 64)
		p.ID, _ = strconv.ParseInt(s.Find("td.sendmsg a").AttrOr("data-playerid", "0"), 10, 64)
		p.Name = strings.TrimSpace(s.Find("span.playername").Text())
		tdName := s.Find("td.name")
		allyTag := tdName.Find("span.ally-tag")
		if allyTag != nil {
			href := allyTag.Find("a").AttrOr("href", "")
			m := regexp.MustCompile(`allianceId=(\d+)`).FindStringSubmatch(href)
			if len(m) == 2 {
				p.AllianceID, _ = strconv.ParseInt(m[1], 10, 64)
			}
			allyTag.Remove()
		}
		href := tdName.Find("a").AttrOr("href", "")
		m := regexp.MustCompile(`galaxy=(\d+)&system=(\d+)&position=(\d+)`).FindStringSubmatch(href)
		if len(m) != 4 {
			return
		}
		p.Homeworld.Type = PlanetType
		p.Homeworld.Galaxy, _ = strconv.ParseInt(m[1], 10, 64)
		p.Homeworld.System, _ = strconv.ParseInt(m[2], 10, 64)
		p.Homeworld.Position, _ = strconv.ParseInt(m[3], 10, 64)
		honorScoreSpan := s.Find("span.honorScore span")
		if honorScoreSpan == nil {
			return
		}
		p.HonourPoints = ParseInt(strings.TrimSpace(honorScoreSpan.Text()))
		p.Score = ParseInt(strings.TrimSpace(s.Find("td.score").Text()))
		out.Players = append(out.Players, p)
	})

	return
}

func extractAllResourcesV71(pageHTML []byte) (out map[CelestialID]Resources, err error) {
	out = make(map[CelestialID]Resources)
	m := regexp.MustCompile(`var planetResources=([^;]+);`).FindSubmatch(pageHTML)
	if len(m) != 2 {
		return out, errors.New("failed to get resources json")
	}
	var data map[string]struct {
		Input struct {
			Metal     int64
			Crystal   int64
			Deuterium int64
		}
	}
	if err := json.Unmarshal(m[1], &data); err != nil {
		return out, err
	}
	for k, v := range data {
		ki, _ := strconv.ParseInt(k, 10, 64)
		out[CelestialID(ki)] = Resources{Metal: v.Input.Metal, Crystal: v.Input.Crystal, Deuterium: v.Input.Deuterium}
	}
	return
}

func extractAttacksFromDocV71(doc *goquery.Document, clock clockwork.Clock) ([]AttackEvent, error) {
	attacks := make([]*AttackEvent, 0)
	out := make([]AttackEvent, 0)
	if doc.Find("body").Size() == 1 && extractOGameSessionFromDocV6(doc) != "" && doc.Find("div#eventListWrap").Size() == 0 {
		return out, ErrEventsBoxNotDisplayed
	} else if doc.Find("div#eventListWrap").Size() == 0 {
		return out, ErrNotLogged
	}

	allianceAttacks := make(map[int64]*AttackEvent)

	tmp := func(rowType string) func(int, *goquery.Selection) {
		return func(i int, s *goquery.Selection) {
			trIDAttr := s.AttrOr("id", "")
			r := regexp.MustCompile(`eventRow-(union)?(\d+)`)
			m := r.FindStringSubmatch(trIDAttr)
			var id int64
			if len(m) != 3 {
				classes := s.AttrOr("class", "")
				r = regexp.MustCompile(`unionunion(\d+)`)
				m = r.FindStringSubmatch(classes)
				if len(m) == 2 {
					id, _ = strconv.ParseInt(m[1], 10, 64)
				}
			} else {
				id, _ = strconv.ParseInt(m[2], 10, 64)
			}

			classes, _ := s.Attr("class")
			partner := strings.Contains(classes, "partnerInfo")

			td := s.Find("td.countDown")
			isHostile := td.HasClass("hostile") || td.Find("span.hostile").Size() > 0
			if !isHostile {
				return
			}
			missionTypeInt, _ := strconv.ParseInt(s.AttrOr("data-mission-type", ""), 10, 64)
			arrivalTimeInt, _ := strconv.ParseInt(s.AttrOr("data-arrival-time", ""), 10, 64)
			missionType := MissionID(missionTypeInt)
			if rowType == "allianceAttack" {
				missionType = GroupedAttack
			}
			if missionType != Attack && missionType != GroupedAttack && missionType != Destroy &&
				missionType != MissileAttack && missionType != Spy {
				return
			}
			attack := &AttackEvent{}
			attack.ID = id
			attack.MissionType = missionType
			if missionType == Attack || missionType == MissileAttack || missionType == Spy || missionType == Destroy || missionType == GroupedAttack {
				linkSendMail := s.Find("a.sendMail")
				attack.AttackerID, _ = strconv.ParseInt(linkSendMail.AttrOr("data-playerid", ""), 10, 64)
				attack.AttackerName = linkSendMail.AttrOr("title", "")
				if attack.AttackerID != 0 {
					coordsOrigin := strings.TrimSpace(s.Find("td.coordsOrigin").Text())
					attack.Origin = extractCoordV6(coordsOrigin)
					attack.Origin.Type = PlanetType
					if s.Find("td.originFleet figure").HasClass("moon") {
						attack.Origin.Type = MoonType
					}
				}
			}
			if missionType == MissileAttack {
				attack.Missiles = ParseInt(s.Find("td.detailsFleet span").First().Text())
			}

			// Get ships infos if available
			if movement, exists := s.Find("td.icon_movement span").Attr("title"); exists {
				root, err := html.Parse(strings.NewReader(movement))
				if err != nil {
					return
				}
				attack.Ships = new(ShipsInfos)
				q := goquery.NewDocumentFromNode(root)
				q.Find("tr").Each(func(i int, s *goquery.Selection) {
					name := s.Find("td").Eq(0).Text()
					nbrTxt := s.Find("td").Eq(1).Text()
					nbr := ParseInt(nbrTxt)
					if name != "" && nbr > 0 {
						attack.Ships.Set(name2id(name), nbr)
					} else if nbrTxt == "?" {
						attack.Ships.Set(name2id(name), -1)
					}
				})
			}

			rgx := regexp.MustCompile(`union(\d+)`)
			classesArr := strings.Split(classes, " ")
			for _, c := range classesArr {
				m := rgx.FindStringSubmatch(c)
				if len(m) == 2 {
					attack.UnionID, _ = strconv.ParseInt(m[1], 10, 64)
				}
			}

			destCoords := strings.TrimSpace(s.Find("td.destCoords").Text())
			attack.Destination = extractCoordV6(destCoords)
			attack.Destination.Type = PlanetType
			if s.Find("td.destFleet figure").HasClass("moon") {
				attack.Destination.Type = MoonType
			}
			attack.DestinationName = strings.TrimSpace(s.Find("td.destFleet").Text())

			attack.ArrivalTime = time.Unix(arrivalTimeInt, 0)
			attack.ArriveIn = int64(clock.Until(attack.ArrivalTime).Seconds())

			if attack.UnionID != 0 {
				if allianceAttack, ok := allianceAttacks[attack.UnionID]; ok {
					if attack.Ships != nil {
						allianceAttack.Ships.Add(*attack.Ships)
					}
					if allianceAttack.AttackerID == 0 {
						allianceAttack.AttackerID = attack.AttackerID
					}
					if allianceAttack.Origin.Equal(Coordinate{}) {
						allianceAttack.Origin = attack.Origin
					}
				} else {
					allianceAttacks[attack.UnionID] = attack
				}
			}

			if !partner {
				attacks = append(attacks, attack)
			}
		}
	}
	doc.Find("tr.allianceAttack").Each(tmp("allianceAttack"))
	doc.Find("tr.eventFleet").Each(tmp("eventFleet"))

	for _, a := range attacks {
		out = append(out, *a)
	}

	return out, nil
}

func extractDMCostsFromDocV71(doc *goquery.Document) (DMCosts, error) {
	tmp := func(s *goquery.Selection) (id ID, nbr, cost int64, canBuy, isComplete bool, buyAndActivate string, token string) {
		imgAlt := s.Find("img.queuePic").AttrOr("alt", "")
		if n, err := fmt.Sscanf(imgAlt, "techId_%d", &id); err != nil || n != 1 {
			return
		}
		r := regexp.MustCompile(`([\d\.]+)`)
		levelTxt := s.Find("span.level").Text()
		if levelTxt == "" {
			levelTxt = s.Find("div.shipSumCount").Text()
		}
		m := r.FindStringSubmatch(levelTxt)
		if len(m) != 2 {
			return
		}
		nbr = ParseInt(m[1])
		canBuy = !s.Find("span.dm_cost").HasClass("overmark")
		costTxt := s.Find("span.dm_cost").Text()
		m = r.FindStringSubmatch(costTxt)
		if len(m) != 2 {
			return
		}
		cost = ParseInt(m[1])
		token = s.Find("a.build-faster").AttrOr("token", "")
		linkRel := s.Find("a.build-faster").AttrOr("rel", "")
		r = regexp.MustCompile(`buyAndActivate=([^"]+)`)
		m = r.FindStringSubmatch(linkRel)
		if len(m) != 2 {
			return
		}
		buyAndActivate = m[1]
		isComplete = s.Find("a.build-faster div").First().HasClass("build-finish-img")
		return
	}
	out := DMCosts{}
	buildingsBox := doc.Find("#productionboxbuildingcomponent")
	researchBox := doc.Find("#productionboxresearchcomponent")
	shipyardBox := doc.Find("#productionboxshipyardcomponent")
	out.Buildings.OGameID, out.Buildings.Nbr, out.Buildings.Cost, out.Buildings.CanBuy, out.Buildings.Complete, out.Buildings.BuyAndActivateToken, out.Buildings.Token = tmp(buildingsBox)
	out.Research.OGameID, out.Research.Nbr, out.Research.Cost, out.Research.CanBuy, out.Research.Complete, out.Research.BuyAndActivateToken, out.Research.Token = tmp(researchBox)
	out.Shipyard.OGameID, out.Shipyard.Nbr, out.Shipyard.Cost, out.Shipyard.CanBuy, out.Shipyard.Complete, out.Shipyard.BuyAndActivateToken, out.Shipyard.Token = tmp(shipyardBox)
	return out, nil
}

func extractBuffActivationFromDocV71(doc *goquery.Document) (token string, items []Item, err error) {
	scriptTxt := doc.Find("script").Text()
	r := regexp.MustCompile(`activateToken = "([^"]+)"`)
	m := r.FindStringSubmatch(scriptTxt)
	if len(m) != 2 {
		err = errors.New("failed to find activate token")
		return
	}
	token = m[1]
	r = regexp.MustCompile(`items_inventory = ({[^\n]+});\n`)
	m = r.FindStringSubmatch(scriptTxt)
	if len(m) != 2 {
		err = errors.New("failed to find items inventory")
		return
	}
	var inventoryMap map[string]Item
	if err = json.Unmarshal([]byte(m[1]), &inventoryMap); err != nil {
		fmt.Println(err)
		return
	}
	for _, item := range inventoryMap {
		items = append(items, item)
	}
	return
}

func extractIsMobileFromDocV71(doc *goquery.Document) bool {
	r := regexp.MustCompile(`var isMobile = (true|false);`)
	scripts := doc.Find("script")
	for i := 0; i < scripts.Size(); i++ {
		scriptText := scripts.Eq(i).Text()
		m := r.FindStringSubmatch(scriptText)
		if len(m) == 2 {
			return m[1] == "true"
		}
	}
	return false
}
