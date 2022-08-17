package ogame

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strings"
	"time"
)

func extractEmpireV9(pageHTML []byte) ([]EmpireCelestial, error) {
	var out []EmpireCelestial
	raw, err := extractEmpireJSON(pageHTML)
	if err != nil {
		return nil, err
	}
	j, ok := raw.(map[string]any)
	if !ok {
		return nil, errors.New("failed to parse json")
	}
	planetsRaw, ok := j["planets"].([]any)
	if !ok {
		return nil, errors.New("failed to parse json")
	}
	for _, planetRaw := range planetsRaw {
		planet, ok := planetRaw.(map[string]any)
		if !ok {
			return nil, errors.New("failed to parse json")
		}

		var tempMin, tempMax int64
		temperatureStr := doCastStr(planet["temperature"])
		m := temperatureRgx.FindStringSubmatch(temperatureStr)
		if len(m) == 3 {
			tempMin = DoParseI64(m[1])
			tempMax = DoParseI64(m[2])
		}
		mm := diameterRgx.FindStringSubmatch(doCastStr(planet["diameter"]))
		energyStr := doCastStr(planet["energy"])
		energyDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(energyStr))
		energy := ParseInt(energyDoc.Find("div span").Text())
		celestialType := CelestialType(doCastF64(planet["type"]))
		out = append(out, EmpireCelestial{
			Name:     doCastStr(planet["name"]),
			ID:       CelestialID(doCastF64(planet["id"])),
			Diameter: ParseInt(mm[1]),
			Img:      doCastStr(planet["image"]),
			Type:     celestialType,
			Fields: Fields{
				Built: DoParseI64(doCastStr(planet["fieldUsed"])),
				Total: DoParseI64(doCastStr(planet["fieldMax"])),
			},
			Temperature: Temperature{
				Min: tempMin,
				Max: tempMax,
			},
			Coordinate: Coordinate{
				Galaxy:   int64(doCastF64(planet["galaxy"])),
				System:   int64(doCastF64(planet["system"])),
				Position: int64(doCastF64(planet["position"])),
				Type:     celestialType,
			},
			Resources: Resources{
				Metal:     int64(doCastF64(planet["metal"])),
				Crystal:   int64(doCastF64(planet["crystal"])),
				Deuterium: int64(doCastF64(planet["deuterium"])),
				Energy:    energy,
			},
			Supplies: ResourcesBuildings{
				MetalMine:            int64(doCastF64(planet["1"])),
				CrystalMine:          int64(doCastF64(planet["2"])),
				DeuteriumSynthesizer: int64(doCastF64(planet["3"])),
				SolarPlant:           int64(doCastF64(planet["4"])),
				FusionReactor:        int64(doCastF64(planet["12"])),
				SolarSatellite:       int64(doCastF64(planet["212"])),
				MetalStorage:         int64(doCastF64(planet["22"])),
				CrystalStorage:       int64(doCastF64(planet["23"])),
				DeuteriumTank:        int64(doCastF64(planet["24"])),
			},
			Facilities: Facilities{
				RoboticsFactory: int64(doCastF64(planet["14"])),
				Shipyard:        int64(doCastF64(planet["21"])),
				ResearchLab:     int64(doCastF64(planet["31"])),
				AllianceDepot:   int64(doCastF64(planet["34"])),
				MissileSilo:     int64(doCastF64(planet["44"])),
				NaniteFactory:   int64(doCastF64(planet["15"])),
				Terraformer:     int64(doCastF64(planet["33"])),
				SpaceDock:       int64(doCastF64(planet["36"])),
				LunarBase:       int64(doCastF64(planet["41"])),
				SensorPhalanx:   int64(doCastF64(planet["42"])),
				JumpGate:        int64(doCastF64(planet["43"])),
			},
			Defenses: DefensesInfos{
				RocketLauncher:         int64(doCastF64(planet["401"])),
				LightLaser:             int64(doCastF64(planet["402"])),
				HeavyLaser:             int64(doCastF64(planet["403"])),
				GaussCannon:            int64(doCastF64(planet["404"])),
				IonCannon:              int64(doCastF64(planet["405"])),
				PlasmaTurret:           int64(doCastF64(planet["406"])),
				SmallShieldDome:        int64(doCastF64(planet["407"])),
				LargeShieldDome:        int64(doCastF64(planet["408"])),
				AntiBallisticMissiles:  int64(doCastF64(planet["502"])),
				InterplanetaryMissiles: int64(doCastF64(planet["503"])),
			},
			Researches: Researches{
				EnergyTechnology:             int64(doCastF64(planet["113"])),
				LaserTechnology:              int64(doCastF64(planet["120"])),
				IonTechnology:                int64(doCastF64(planet["121"])),
				HyperspaceTechnology:         int64(doCastF64(planet["114"])),
				PlasmaTechnology:             int64(doCastF64(planet["122"])),
				CombustionDrive:              int64(doCastF64(planet["115"])),
				ImpulseDrive:                 int64(doCastF64(planet["117"])),
				HyperspaceDrive:              int64(doCastF64(planet["118"])),
				EspionageTechnology:          int64(doCastF64(planet["106"])),
				ComputerTechnology:           int64(doCastF64(planet["108"])),
				Astrophysics:                 int64(doCastF64(planet["124"])),
				IntergalacticResearchNetwork: int64(doCastF64(planet["123"])),
				GravitonTechnology:           int64(doCastF64(planet["199"])),
				WeaponsTechnology:            int64(doCastF64(planet["109"])),
				ShieldingTechnology:          int64(doCastF64(planet["110"])),
				ArmourTechnology:             int64(doCastF64(planet["111"])),
			},
			Ships: ShipsInfos{
				LightFighter:   int64(doCastF64(planet["204"])),
				HeavyFighter:   int64(doCastF64(planet["205"])),
				Cruiser:        int64(doCastF64(planet["206"])),
				Battleship:     int64(doCastF64(planet["207"])),
				Battlecruiser:  int64(doCastF64(planet["215"])),
				Bomber:         int64(doCastF64(planet["211"])),
				Destroyer:      int64(doCastF64(planet["213"])),
				Deathstar:      int64(doCastF64(planet["214"])),
				SmallCargo:     int64(doCastF64(planet["202"])),
				LargeCargo:     int64(doCastF64(planet["203"])),
				ColonyShip:     int64(doCastF64(planet["208"])),
				Recycler:       int64(doCastF64(planet["209"])),
				EspionageProbe: int64(doCastF64(planet["210"])),
				SolarSatellite: int64(doCastF64(planet["212"])),
				Crawler:        int64(doCastF64(planet["217"])),
				Reaper:         int64(doCastF64(planet["218"])),
				Pathfinder:     int64(doCastF64(planet["219"])),
			},
		})
	}
	return out, nil
}

func extractOverviewProductionFromDocV9(doc *goquery.Document) ([]Quantifiable, error) {
	res := make([]Quantifiable, 0)
	active := doc.Find("table.construction").Eq(4)
	href, _ := active.Find("td a").Attr("href")
	m := regexp.MustCompile(`openTech=(\d+)`).FindStringSubmatch(href)
	if len(m) == 0 {
		return []Quantifiable{}, nil
	}
	idInt := DoParseI64(m[1])
	activeID := ID(idInt)
	activeNbr := DoParseI64(active.Find("div.shipSumCount").Text())
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
	report.ID = DoParseI64(doc.Find("div.detail_msg").AttrOr("data-msg-id", "0"))
	spanLink := doc.Find("span.msg_title a").First()
	txt := spanLink.Text()
	figure := spanLink.Find("figure").First()
	r := regexp.MustCompile(`([^\[]+) \[(\d+):(\d+):(\d+)]`)
	m := r.FindStringSubmatch(txt)
	if len(m) == 5 {
		report.Coordinate.Galaxy = DoParseI64(m[2])
		report.Coordinate.System = DoParseI64(m[3])
		report.Coordinate.Position = DoParseI64(m[4])
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
		report.CounterEspionage = DoParseI64(m1[1])
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
				buildingID := DoParseI64(r.FindStringSubmatch(imgClass)[1])
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
				researchID := DoParseI64(r.FindStringSubmatch(imgClass)[1])
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
				shipID := DoParseI64(r.FindStringSubmatch(imgClass)[1])
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
				defenceID := DoParseI64(r.FindStringSubmatch(imgClass)[1])
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
