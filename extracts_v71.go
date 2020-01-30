package ogame

import (
	"encoding/json"
	"errors"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
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
			report.HasBuildings = s.Find("li.detail_list_fail").Size() == 0
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
			report.HasResearches = s.Find("li.detail_list_fail").Size() == 0
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
			report.HasFleet = s.Find("li.detail_list_fail").Size() == 0
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
			report.HasDefenses = s.Find("li.detail_list_fail").Size() == 0
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
	script := doc.Find("script").First().Text()
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

	changeSiteSize := doc.Find("select.changeSite option").Size()
	out.NbPage = MaxInt(int64(changeSiteSize)-1, 0)

	doc.Find("#ranks tbody tr").Each(func(i int, s *goquery.Selection) {
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
