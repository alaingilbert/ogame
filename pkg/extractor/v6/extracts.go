package v6

import (
	"bytes"
	"encoding/json"
	"errors"
	"math"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/alaingilbert/clockwork"
	lua "github.com/yuin/gopher-lua"
	"golang.org/x/net/html"
)

func extractUpgradeToken(pageHTML []byte) (string, error) {
	rgx := regexp.MustCompile(`var upgradeEndpoint = ".+&token=([^&]+)&`)
	m := rgx.FindSubmatch(pageHTML)
	if len(m) != 2 {
		return "", errors.New("unable to find form token")
	}
	return string(m[1]), nil
}

func extractTearDownButtonEnabledFromDoc(doc *goquery.Document) bool {
	return !doc.Find("a.demolish_link div").HasClass("demolish_img_disabled")
}

func extractIsInVacationFromDoc(doc *goquery.Document) bool {
	href := doc.Find("div#advice-bar a").AttrOr("href", "")
	if href == "" {
		return false
	}
	u, _ := url.Parse(href)
	q := u.Query()
	if q.Get("page") == "preferences" && q.Get("selectedTab") == "3" && q.Get("openGroup") == "0" {
		return true
	}
	return false
}

func extractResourcesFromDoc(doc *goquery.Document) ogame.Resources {
	res := ogame.Resources{}
	metalDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(doc.Find("li#metal_box").AttrOr("title", "")))
	crystalDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(doc.Find("li#crystal_box").AttrOr("title", "")))
	deuteriumDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(doc.Find("li#deuterium_box").AttrOr("title", "")))
	energyDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(doc.Find("li#energy_box").AttrOr("title", "")))
	darkmatterDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(doc.Find("li#darkmatter_box").AttrOr("title", "")))
	res.Metal = utils.ParseInt(metalDoc.Find("table tr").Eq(0).Find("td").Eq(0).Text())
	res.Crystal = utils.ParseInt(crystalDoc.Find("table tr").Eq(0).Find("td").Eq(0).Text())
	res.Deuterium = utils.ParseInt(deuteriumDoc.Find("table tr").Eq(0).Find("td").Eq(0).Text())
	res.Energy = utils.ParseInt(energyDoc.Find("table tr").Eq(0).Find("td").Eq(0).Text())
	res.Darkmatter = utils.ParseInt(darkmatterDoc.Find("table tr").Eq(0).Find("td").Eq(0).Text())
	return res
}

func extractResourcesDetailsFromFullPageFromDoc(doc *goquery.Document) ogame.ResourcesDetails {
	out := ogame.ResourcesDetails{}
	out.Metal.Available = utils.ParseInt(doc.Find("span#resources_metal").Text())
	out.Crystal.Available = utils.ParseInt(doc.Find("span#resources_crystal").Text())
	out.Deuterium.Available = utils.ParseInt(doc.Find("span#resources_deuterium").Text())
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

func extractHiddenFieldsFromDoc(doc *goquery.Document) url.Values {
	fields := url.Values{}
	doc.Find("input[type=hidden]").Each(func(i int, s *goquery.Selection) {
		name, _ := s.Attr("name")
		value, _ := s.Attr("value")
		fields.Add(name, value)
	})
	return fields
}

func ExtractBodyIDFromDoc(doc *goquery.Document) string {
	bodyID := doc.Find("body").AttrOr("id", "")
	if bodyID == "ingamepage" {
		pageHTML, _ := doc.Html()
		m := regexp.MustCompile(`var currentPage = "([^"]+)";`).FindStringSubmatch(pageHTML)
		if len(m) == 2 {
			return m[1]
		}
	}
	return bodyID
}

func extractCelestialByIDFromDoc(doc *goquery.Document, celestialID ogame.CelestialID) (ogame.Celestial, error) {
	celestials := extractCelestialsFromDoc(doc)
	for _, celestial := range celestials {
		if celestial.GetID() == celestialID {
			return celestial, nil
		}
	}
	return nil, errors.New("invalid celestial id")
}

func extractCelestialByCoordFromDoc(doc *goquery.Document, coord ogame.Coordinate) (ogame.Celestial, error) {
	celestials := extractCelestialsFromDoc(doc)
	for _, celestial := range celestials {
		if celestial.GetCoordinate().Equal(coord) {
			return celestial, nil
		}
	}
	return nil, errors.New("invalid coordinate")
}

func extractCelestialsFromDoc(doc *goquery.Document) []ogame.Celestial {
	celestials := make([]ogame.Celestial, 0)
	planets := extractPlanetsFromDoc(doc)
	for _, planet := range planets {
		celestials = append(celestials, planet)
		if planet.Moon != nil {
			celestials = append(celestials, *planet.Moon)
		}
	}
	return celestials
}

func extractPlanetsFromDoc(doc *goquery.Document) []ogame.Planet {
	res := make([]ogame.Planet, 0)
	doc.Find("div.smallplanet").Each(func(i int, s *goquery.Selection) {
		planet, err := extractPlanetFromSelection(s)
		if err != nil {
			return
		}
		res = append(res, planet)
	})
	return res
}

func extractMoonsFromDoc(doc *goquery.Document) []ogame.Moon {
	res := make([]ogame.Moon, 0)
	doc.Find("a.moonlink").Each(func(i int, s *goquery.Selection) {
		moon, err := extractMoonFromSelection(s)
		if err != nil {
			return
		}
		res = append(res, moon)
	})
	return res
}

func extractOgameTimestampFromDoc(doc *goquery.Document) int64 {
	ogameTimestamp := utils.DoParseI64(doc.Find("meta[name=ogame-timestamp]").AttrOr("content", "0"))
	return ogameTimestamp
}

type CelestialTypes interface {
	ogame.Planet | ogame.Moon
}

func extractPlanetMoonFromDoc[T CelestialTypes](doc *goquery.Document, v any) (T, error) {
	var zero T
	celestial, err := extractCelestialFromDoc(doc, v)
	if err != nil {
		return zero, err
	}
	if typed, ok := celestial.(T); ok {
		return typed, nil
	}
	return zero, errors.New("not found")
}

func extractPlanetFromDoc(doc *goquery.Document, v any) (ogame.Planet, error) {
	return extractPlanetMoonFromDoc[ogame.Planet](doc, v)
}

func extractMoonFromDoc(doc *goquery.Document, v any) (ogame.Moon, error) {
	return extractPlanetMoonFromDoc[ogame.Moon](doc, v)
}

func extractCelestialFromDoc(doc *goquery.Document, v any) (ogame.Celestial, error) {
	switch vv := v.(type) {
	case ogame.Celestial:
		return extractCelestialByIDFromDoc(doc, vv.GetID())
	case ogame.PlanetID:
		return extractCelestialByIDFromDoc(doc, vv.Celestial())
	case ogame.MoonID:
		return extractCelestialByIDFromDoc(doc, vv.Celestial())
	case ogame.CelestialID:
		return extractCelestialByIDFromDoc(doc, vv)
	case int:
		return extractCelestialByIDFromDoc(doc, ogame.CelestialID(vv))
	case int32:
		return extractCelestialByIDFromDoc(doc, ogame.CelestialID(vv))
	case int64:
		return extractCelestialByIDFromDoc(doc, ogame.CelestialID(vv))
	case float32:
		return extractCelestialByIDFromDoc(doc, ogame.CelestialID(vv))
	case float64:
		return extractCelestialByIDFromDoc(doc, ogame.CelestialID(vv))
	case lua.LNumber:
		return extractCelestialByIDFromDoc(doc, ogame.CelestialID(vv))
	case ogame.Coordinate:
		return extractCelestialByCoordFromDoc(doc, vv)
	case string:
		coord, err := ogame.ParseCoord(vv)
		if err != nil {
			return nil, err
		}
		return extractCelestialByCoordFromDoc(doc, coord)
	default:
		return nil, ErrUnsupportedType
	}
}

var ErrUnsupportedType = errors.New("unsupported type")

func extractResourcesBuildingsFromDoc(doc *goquery.Document) (ogame.ResourcesBuildings, error) {
	doc.Find("span.textlabel").Remove()
	bodyID := ExtractBodyIDFromDoc(doc)
	if bodyID == "overview" {
		return ogame.ResourcesBuildings{}, ogame.ErrInvalidPlanetID
	}
	res := ogame.ResourcesBuildings{}
	res.MetalMine = utils.GetNbr(doc, "supply1")
	res.CrystalMine = utils.GetNbr(doc, "supply2")
	res.DeuteriumSynthesizer = utils.GetNbr(doc, "supply3")
	res.SolarPlant = utils.GetNbr(doc, "supply4")
	res.FusionReactor = utils.GetNbr(doc, "supply12")
	res.SolarSatellite = utils.GetNbr(doc, "supply212")
	res.MetalStorage = utils.GetNbr(doc, "supply22")
	res.CrystalStorage = utils.GetNbr(doc, "supply23")
	res.DeuteriumTank = utils.GetNbr(doc, "supply24")
	return res, nil
}

func extractDefenseFromDoc(doc *goquery.Document) (ogame.DefensesInfos, error) {
	bodyID := ExtractBodyIDFromDoc(doc)
	if bodyID == "overview" {
		return ogame.DefensesInfos{}, ogame.ErrInvalidPlanetID
	}
	doc.Find("span.textlabel").Remove()
	res := ogame.DefensesInfos{}
	res.RocketLauncher = utils.GetNbr(doc, "defense401")
	res.LightLaser = utils.GetNbr(doc, "defense402")
	res.HeavyLaser = utils.GetNbr(doc, "defense403")
	res.GaussCannon = utils.GetNbr(doc, "defense404")
	res.IonCannon = utils.GetNbr(doc, "defense405")
	res.PlasmaTurret = utils.GetNbr(doc, "defense406")
	res.SmallShieldDome = utils.GetNbr(doc, "defense407")
	res.LargeShieldDome = utils.GetNbr(doc, "defense408")
	res.AntiBallisticMissiles = utils.GetNbr(doc, "defense502")
	res.InterplanetaryMissiles = utils.GetNbr(doc, "defense503")
	return res, nil
}

func extractShipsFromDoc(doc *goquery.Document) (ogame.ShipsInfos, error) {
	doc.Find("span.textlabel").Remove()
	bodyID := ExtractBodyIDFromDoc(doc)
	if bodyID == "overview" {
		return ogame.ShipsInfos{}, ogame.ErrInvalidPlanetID
	}
	res := ogame.ShipsInfos{}
	res.LightFighter = utils.GetNbrShips(doc, "military204")
	res.HeavyFighter = utils.GetNbrShips(doc, "military205")
	res.Cruiser = utils.GetNbrShips(doc, "military206")
	res.Battleship = utils.GetNbrShips(doc, "military207")
	res.Battlecruiser = utils.GetNbrShips(doc, "military215")
	res.Bomber = utils.GetNbrShips(doc, "military211")
	res.Destroyer = utils.GetNbrShips(doc, "military213")
	res.Deathstar = utils.GetNbrShips(doc, "military214")
	res.SmallCargo = utils.GetNbrShips(doc, "civil202")
	res.LargeCargo = utils.GetNbrShips(doc, "civil203")
	res.ColonyShip = utils.GetNbrShips(doc, "civil208")
	res.Recycler = utils.GetNbrShips(doc, "civil209")
	res.EspionageProbe = utils.GetNbrShips(doc, "civil210")
	res.SolarSatellite = utils.GetNbrShips(doc, "civil212")

	return res, nil
}

func extractFacilitiesFromDoc(doc *goquery.Document) (ogame.Facilities, error) {
	doc.Find("span.textlabel").Remove()
	bodyID := ExtractBodyIDFromDoc(doc)
	if bodyID == "overview" {
		return ogame.Facilities{}, ogame.ErrInvalidPlanetID
	}
	res := ogame.Facilities{}
	res.RoboticsFactory = utils.GetNbr(doc, "station14")
	res.Shipyard = utils.GetNbr(doc, "station21")
	res.ResearchLab = utils.GetNbr(doc, "station31")
	res.AllianceDepot = utils.GetNbr(doc, "station34")
	res.MissileSilo = utils.GetNbr(doc, "station44")
	res.NaniteFactory = utils.GetNbr(doc, "station15")
	res.Terraformer = utils.GetNbr(doc, "station33")
	res.SpaceDock = utils.GetNbr(doc, "station36")
	res.LunarBase = utils.GetNbr(doc, "station41")
	res.SensorPhalanx = utils.GetNbr(doc, "station42")
	res.JumpGate = utils.GetNbr(doc, "station43")
	return res, nil
}

func extractResearchFromDoc(doc *goquery.Document) ogame.Researches {
	doc.Find("span.textlabel").Remove()
	res := ogame.Researches{}
	res.EnergyTechnology = utils.GetNbr(doc, "research113")
	res.LaserTechnology = utils.GetNbr(doc, "research120")
	res.IonTechnology = utils.GetNbr(doc, "research121")
	res.HyperspaceTechnology = utils.GetNbr(doc, "research114")
	res.PlasmaTechnology = utils.GetNbr(doc, "research122")
	res.CombustionDrive = utils.GetNbr(doc, "research115")
	res.ImpulseDrive = utils.GetNbr(doc, "research117")
	res.HyperspaceDrive = utils.GetNbr(doc, "research118")
	res.EspionageTechnology = utils.GetNbr(doc, "research106")
	res.ComputerTechnology = utils.GetNbr(doc, "research108")
	res.Astrophysics = utils.GetNbr(doc, "research124")
	res.IntergalacticResearchNetwork = utils.GetNbr(doc, "research123")
	res.GravitonTechnology = utils.GetNbr(doc, "research199")
	res.WeaponsTechnology = utils.GetNbr(doc, "research109")
	res.ShieldingTechnology = utils.GetNbr(doc, "research110")
	res.ArmourTechnology = utils.GetNbr(doc, "research111")
	return res
}

func ExtractOGameSessionFromDoc(doc *goquery.Document) string {
	sessionMeta := doc.Find("meta[name=ogame-session]")
	if sessionMeta.Size() == 0 {
		r := regexp.MustCompile(`var session = "([^"]+)";`)
		scripts := doc.Find("script")
		for i := 0; i < scripts.Size(); i++ {
			scriptText := scripts.Eq(i).Text()
			m := r.FindStringSubmatch(scriptText)
			if len(m) == 2 {
				return m[1]
			}
		}
	}
	return sessionMeta.AttrOr("content", "")
}

func extractAttacksFromDoc(doc *goquery.Document, clock clockwork.Clock, ownCoords []ogame.Coordinate) ([]ogame.AttackEvent, error) {
	attacks := make([]*ogame.AttackEvent, 0)
	out := make([]ogame.AttackEvent, 0)
	if doc.Find("body").Size() == 1 && ExtractOGameSessionFromDoc(doc) != "" && doc.Find("div#eventListWrap").Size() == 0 {
		return out, ogame.ErrEventsBoxNotDisplayed
	} else if doc.Find("div#eventListWrap").Size() == 0 {
		return out, ogame.ErrNotLogged
	}

	allianceAttacks := make(map[int64]*ogame.AttackEvent)

	tmp := func(i int, s *goquery.Selection) {
		classes, _ := s.Attr("class")
		partner := strings.Contains(classes, "partnerInfo")

		td := s.Find("td.countDown")
		isHostile := td.HasClass("hostile") || td.Find("span.hostile").Size() > 0
		if !isHostile {
			return
		}
		missionTypeInt := utils.DoParseI64(s.AttrOr("data-mission-type", ""))
		arrivalTimeInt := utils.DoParseI64(s.AttrOr("data-arrival-time", ""))
		missionType := ogame.MissionID(missionTypeInt)
		if missionType != ogame.Attack && missionType != ogame.GroupedAttack && missionType != ogame.Destroy &&
			missionType != ogame.MissileAttack && missionType != ogame.Spy {
			return
		}
		attack := &ogame.AttackEvent{}
		attack.MissionType = missionType
		if missionType == ogame.Attack || missionType == ogame.MissileAttack || missionType == ogame.Spy || missionType == ogame.Destroy || missionType == ogame.GroupedAttack {
			linkSendMail := s.Find("a.sendMail")
			attack.AttackerID = utils.DoParseI64(linkSendMail.AttrOr("data-playerid", ""))
			attack.AttackerName = linkSendMail.AttrOr("title", "")
			if attack.AttackerID != 0 {
				coordsOrigin := strings.TrimSpace(s.Find("td.coordsOrigin").Text())
				attack.Origin = ExtractCoord(coordsOrigin)
				attack.Origin.Type = ogame.PlanetType
				if s.Find("td.originFleet figure").HasClass("moon") {
					attack.Origin.Type = ogame.MoonType
				}
			}
		}
		if missionType == ogame.MissileAttack {
			attack.Missiles = utils.ParseInt(s.Find("td.detailsFleet span").First().Text())
		}

		// Get ships infos if available
		if movement, exists := s.Find("td.icon_movement span").Attr("title"); exists {
			root, err := html.Parse(strings.NewReader(movement))
			if err != nil {
				return
			}
			attack.Ships = new(ogame.ShipsInfos)
			q := goquery.NewDocumentFromNode(root)
			q.Find("tr").Each(func(i int, s *goquery.Selection) {
				name := s.Find("td").Eq(0).Text()
				nbrTxt := s.Find("td").Eq(1).Text()
				nbr := utils.ParseInt(nbrTxt)
				if name != "" && nbr > 0 {
					attack.Ships.Set(ogame.ShipName2ID(name), nbr)
				} else if nbrTxt == "?" {
					attack.Ships.Set(ogame.ShipName2ID(name), -1)
				}
			})
		}

		rgx := regexp.MustCompile(`union(\d+)`)
		classesArr := strings.Split(classes, " ")
		for _, c := range classesArr {
			m := rgx.FindStringSubmatch(c)
			if len(m) == 2 {
				attack.UnionID = utils.DoParseI64(m[1])
			}
		}

		destCoords := strings.TrimSpace(s.Find("td.destCoords").Text())
		attack.Destination = ExtractCoord(destCoords)
		attack.Destination.Type = ogame.PlanetType
		if s.Find("td.destFleet figure").HasClass("moon") {
			attack.Destination.Type = ogame.MoonType
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
				if allianceAttack.Origin.Equal(ogame.Coordinate{}) {
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
	doc.Find("tr.allianceAttack").Each(tmp)
	doc.Find("tr.eventFleet").Each(tmp)

	for _, a := range attacks {
		out = append(out, *a)
	}

	return out, nil
}

func extractOfferOfTheDayFromDoc(doc *goquery.Document) (price int64, importToken string, planetResources ogame.PlanetResources, multiplier ogame.Multiplier, err error) {
	s := doc.Find("div.js_import_price")
	if s.Size() == 0 {
		err = errors.New("failed to extract offer of the day price")
		return
	}
	price = utils.ParseInt(s.Text())
	script := doc.Find("script").Text()
	m := regexp.MustCompile(`var importToken\s?=\s?"([^"]*)";`).FindSubmatch([]byte(script))
	if len(m) != 2 {
		err = errors.New("failed to extract offer of the day import token")
		return
	}
	importToken = string(m[1])
	m = regexp.MustCompile(`var planetResources\s?=\s?({[^;]*});`).FindSubmatch([]byte(script))
	if len(m) != 2 {
		err = errors.New("failed to extract offer of the day raw planet resources")
		return
	}
	if err = json.Unmarshal(m[1], &planetResources); err != nil {
		return
	}
	m = regexp.MustCompile(`var multiplier\s?=\s?({[^;]*});`).FindSubmatch([]byte(script))
	if len(m) != 2 {
		err = errors.New("failed to extract offer of the day raw multiplier")
		return
	}
	if err = json.Unmarshal(m[1], &multiplier); err != nil {
		return
	}
	return
}

func extractProductionFromDoc(doc *goquery.Document) ([]ogame.Quantifiable, error) {
	res := make([]ogame.Quantifiable, 0)
	active := doc.Find("table.construction")
	href, _ := active.Find("td a").Attr("href")
	m := regexp.MustCompile(`openTech=(\d+)`).FindStringSubmatch(href)
	if len(m) == 0 {
		return []ogame.Quantifiable{}, nil
	}
	idInt := utils.DoParseI64(m[1])
	activeID := ogame.ID(idInt)
	activeNbr := utils.DoParseI64(active.Find("div.shipSumCount").Text())
	res = append(res, ogame.Quantifiable{ID: activeID, Nbr: activeNbr})
	doc.Find("div#pqueue ul li").Each(func(i int, s *goquery.Selection) {
		link := s.Find("a")
		itemIDstr, exists := link.Attr("ref")
		if !exists {
			href := link.AttrOr("href", "")
			m := regexp.MustCompile(`openTech=(\d+)`).FindStringSubmatch(href)
			if len(m) > 0 {
				itemIDstr = m[1]
			} else {
				src := s.Find("img").AttrOr("src", "")
				if strings.HasSuffix(src, "fb4e438cabd12ef1b0500a0f41abc1.jpg") {
					itemIDstr = utils.FI64(ogame.AntiBallisticMissilesID)
				} else if strings.HasSuffix(src, "36221e9493458b9fcc776bf350983e.jpg") {
					itemIDstr = utils.FI64(ogame.InterplanetaryMissilesID)
				}
			}
		}
		itemID := utils.DoParseI64(itemIDstr)
		itemNbr := utils.ParseInt(s.Find("span.number").Text())
		res = append(res, ogame.Quantifiable{ID: ogame.ID(itemID), Nbr: itemNbr})
	})
	return res, nil
}

func extractOverviewProductionFromDoc(doc *goquery.Document) ([]ogame.Quantifiable, error) {
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
		link := s.Find("a")
		href := link.AttrOr("href", "")
		m := regexp.MustCompile(`openTech=(\d+)`).FindStringSubmatch(href)
		if len(m) == 0 {
			return
		}
		idInt := utils.DoParseI64(m[1])
		activeID := ogame.ID(idInt)
		activeNbr := utils.ParseInt(link.Text())
		res = append(res, ogame.Quantifiable{ID: activeID, Nbr: activeNbr})
	})
	return res, nil
}

func extractFleet1ShipsFromDoc(doc *goquery.Document) (s ogame.ShipsInfos) {
	onclick := doc.Find("a#sendall").AttrOr("onclick", "")
	matches := regexp.MustCompile(`setMaxIntInput\("form\[name=shipsChosen]", (.+)\); checkShips`).FindStringSubmatch(onclick)
	if len(matches) == 0 {
		return
	}
	m := matches[1]
	var res map[ogame.ID]int64
	if err := json.Unmarshal([]byte(m), &res); err != nil {
		return
	}
	for k, v := range res {
		s.Set(k, v)
	}
	return
}

func extractFleetDispatchACSFromDoc(doc *goquery.Document) []ogame.ACSValues {
	out := make([]ogame.ACSValues, 0)
	doc.Find("select[name=acsValues] option").Each(func(i int, s *goquery.Selection) {
		acsValues := s.AttrOr("value", "")
		m := regexp.MustCompile(`\d+#\d+#\d+#\d+#.*#(\d+)`).FindStringSubmatch(acsValues)
		if len(m) == 2 {
			optUnionID := utils.DoParseI64(m[1])
			out = append(out, ogame.ACSValues{ACSValues: acsValues, Union: optUnionID})
		}
	})
	return out
}

func extractEspionageReportMessageIDsFromDoc(doc *goquery.Document) ([]ogame.EspionageReportSummary, int64) {
	msgs := make([]ogame.EspionageReportSummary, 0)
	nbPage := utils.DoParseI64(doc.Find("ul.pagination li").Last().AttrOr("data-page", "1"))
	doc.Find("li.msg").Each(func(i int, s *goquery.Selection) {
		if idStr, exists := s.Attr("data-msg-id"); exists {
			if id, err := utils.ParseI64(idStr); err == nil {
				messageType := ogame.Report
				if s.Find("span.espionageDefText").Size() > 0 {
					messageType = ogame.Action
				}
				report := ogame.EspionageReportSummary{ID: id, Type: messageType}
				report.From = s.Find("span.msg_sender").Text()
				spanLink := s.Find("span.msg_title a")
				targetStr := spanLink.Text()
				report.Target = ExtractCoord(targetStr)
				report.Target.Type = ogame.PlanetType
				if spanLink.Find("figure").HasClass("moon") {
					report.Target.Type = ogame.MoonType
				}
				if messageType == ogame.Report {
					s.Find("div.compacting").Each(func(i int, s *goquery.Selection) {
						if regexp.MustCompile(`%`).MatchString(s.Text()) {
							report.LootPercentage, _ = strconv.ParseFloat(regexp.MustCompile(`: (\d+)%`).FindStringSubmatch(s.Text())[1], 64)
							report.LootPercentage /= 100
						}
					})
				}
				msgs = append(msgs, report)

			}
		}
	})
	return msgs, nbPage
}

func extractCombatReportMessagesFromDoc(doc *goquery.Document) ([]ogame.CombatReportSummary, int64) {
	msgs := make([]ogame.CombatReportSummary, 0)
	nbPage := utils.DoParseI64(doc.Find("ul.pagination li").Last().AttrOr("data-page", "1"))
	doc.Find("li.msg").Each(func(i int, s *goquery.Selection) {
		if idStr, exists := s.Attr("data-msg-id"); exists {
			if id, err := utils.ParseI64(idStr); err == nil {
				report := ogame.CombatReportSummary{ID: id}
				report.Destination = ExtractCoord(s.Find("div.msg_head a").Text())
				if s.Find("div.msg_head figure").HasClass("planet") {
					report.Destination.Type = ogame.PlanetType
				} else if s.Find("div.msg_head figure").HasClass("moon") {
					report.Destination.Type = ogame.MoonType
				} else {
					report.Destination.Type = ogame.PlanetType
				}
				resTitle := s.Find("span.msg_content div.combatLeftSide span").Eq(1).AttrOr("title", "")
				m := regexp.MustCompile(`([\d.,]+)<br/>\D*([\d.,]+)<br/>\D*([\d.,]+)`).FindStringSubmatch(resTitle)
				if len(m) == 4 {
					report.Metal = utils.ParseInt(m[1])
					report.Crystal = utils.ParseInt(m[2])
					report.Deuterium = utils.ParseInt(m[3])
				}
				debrisFieldTitle := s.Find("span.msg_content div.combatLeftSide span").Eq(2).AttrOr("title", "0")
				report.DebrisField = utils.ParseInt(debrisFieldTitle)
				resText := s.Find("span.msg_content div.combatLeftSide span").Eq(1).Text()
				m = regexp.MustCompile(`[\d.,]+\D*([\d.,]+)`).FindStringSubmatch(resText)
				if len(m) == 2 {
					report.Loot = utils.ParseInt(m[1])
				}
				msgDate, _ := time.Parse("02.01.2006 15:04:05", s.Find("span.msg_date").Text())
				report.CreatedAt = msgDate

				link := s.Find("div.msg_actions a span.icon_attack").Parent().AttrOr("href", "")
				m = regexp.MustCompile(`page=fleet1&galaxy=(\d+)&system=(\d+)&position=(\d+)&type=(\d+)&`).FindStringSubmatch(link)
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
	activity := doc.Find("div.detail_txt").Eq(1).Find("font")
	if len(activity.Text()) == 2 {
		report.LastActivity = utils.ParseInt(activity.Text())
	}

	// CounterEspionage
	ceTxt := doc.Find("div.detail_txt").Eq(1).Text()
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

func extractResourcesProductionsFromDoc(doc *goquery.Document) (ogame.Resources, error) {
	res := ogame.Resources{}
	selector := "table.listOfResourceSettingsPerPlanet tr.summary td span"
	el := doc.Find(selector)
	res.Metal = utils.ParseInt(el.Eq(0).AttrOr("title", "0"))
	res.Crystal = utils.ParseInt(el.Eq(1).AttrOr("title", "0"))
	res.Deuterium = utils.ParseInt(el.Eq(2).AttrOr("title", "0"))
	res.Energy = utils.ParseInt(el.Eq(3).AttrOr("title", "0"))
	return res, nil
}

func extractPreferencesFromDoc(doc *goquery.Document) ogame.Preferences {
	prefs := ogame.Preferences{
		SpioAnz:                      extractSpioAnzFromDoc(doc),
		DisableChatBar:               extractDisableChatBarFromDoc(doc),
		DisableOutlawWarning:         extractDisableOutlawWarningFromDoc(doc),
		MobileVersion:                extractMobileVersionFromDoc(doc),
		ShowOldDropDowns:             extractShowOldDropDownsFromDoc(doc),
		ActivateAutofocus:            extractActivateAutofocusFromDoc(doc),
		EventsShow:                   extractEventsShowFromDoc(doc),
		SortSetting:                  extractSortSettingFromDoc(doc),
		SortOrder:                    extractSortOrderFromDoc(doc),
		ShowDetailOverlay:            extractShowDetailOverlayFromDoc(doc),
		AnimatedSliders:              extractAnimatedSlidersFromDoc(doc),
		AnimatedOverview:             extractAnimatedOverviewFromDoc(doc),
		PopupsNotices:                extractPopupsNoticesFromDoc(doc),
		PopopsCombatreport:           extractPopopsCombatreportFromDoc(doc),
		SpioReportPictures:           extractSpioReportPicturesFromDoc(doc),
		MsgResultsPerPage:            extractMsgResultsPerPageFromDoc(doc),
		AuctioneerNotifications:      extractAuctioneerNotificationsFromDoc(doc),
		EconomyNotifications:         extractEconomyNotificationsFromDoc(doc),
		ShowActivityMinutes:          extractShowActivityMinutesFromDoc(doc),
		PreserveSystemOnPlanetChange: extractPreserveSystemOnPlanetChangeFromDoc(doc),
		UrlaubsModus:                 extractUrlaubsModus(doc),
	}
	if prefs.MobileVersion {
		prefs.Notifications.BuildList = extractNotifBuildListFromDoc(doc)
		prefs.Notifications.FriendlyFleetActivities = extractNotifFriendlyFleetActivitiesFromDoc(doc)
		prefs.Notifications.HostileFleetActivities = extractNotifHostileFleetActivitiesFromDoc(doc)
		prefs.Notifications.ForeignEspionage = extractNotifForeignEspionageFromDoc(doc)
		prefs.Notifications.AllianceBroadcasts = extractNotifAllianceBroadcastsFromDoc(doc)
		prefs.Notifications.AllianceMessages = extractNotifAllianceMessagesFromDoc(doc)
		prefs.Notifications.Auctions = extractNotifAuctionsFromDoc(doc)
		prefs.Notifications.Account = extractNotifAccountFromDoc(doc)
	}
	return prefs
}

func extractResourceSettingsFromDoc(doc *goquery.Document) (ogame.ResourceSettings, string, error) {
	bodyID := ExtractBodyIDFromDoc(doc)
	if bodyID == "overview" {
		return ogame.ResourceSettings{}, "", ogame.ErrInvalidPlanetID
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
	if len(vals) != 6 {
		return ogame.ResourceSettings{}, "", errors.New("failed to find all resource settings")
	}

	res := ogame.ResourceSettings{}
	res.MetalMine = vals[0]
	res.CrystalMine = vals[1]
	res.DeuteriumSynthesizer = vals[2]
	res.SolarPlant = vals[3]
	res.FusionReactor = vals[4]
	res.SolarSatellite = vals[5]

	token, exists := doc.Find("form input[name=token]").Attr("value")
	if !exists {
		return ogame.ResourceSettings{}, "", errors.New("unable to find token")
	}

	return res, token, nil
}

func extractFleetsFromEventListFromDoc(doc *goquery.Document) []ogame.Fleet {
	type Tmp struct {
		fleet ogame.Fleet
		res   ogame.Resources
	}
	tmp := make([]Tmp, 0)
	res := make([]ogame.Fleet, 0)
	doc.Find("tr.eventFleet").Each(func(i int, s *goquery.Selection) {
		fleet := ogame.Fleet{}

		movement := s.Find("td span.tooltip").AttrOr("title", "")
		if movement == "" {
			return
		}

		root, _ := html.Parse(strings.NewReader(movement))
		doc2 := goquery.NewDocumentFromNode(root)
		doc2.Find("tr").Each(func(i int, s *goquery.Selection) {
			if i == 0 {
				return
			}
			name := s.Find("td").Eq(0).Text()
			nbr := utils.ParseInt(s.Find("td").Eq(1).Text())
			if name != "" && nbr > 0 {
				fleet.Ships.Set(ogame.ShipName2ID(name), nbr)
			}
		})
		fleet.Origin = ExtractCoord(doc.Find("td.coordsOrigin").Text())
		fleet.Destination = ExtractCoord(doc.Find("td.destCoords").Text())

		res := ogame.Resources{}
		trs := doc2.Find("tr")
		res.Metal = utils.ParseInt(trs.Eq(trs.Size() - 3).Find("td").Eq(1).Text())
		res.Crystal = utils.ParseInt(trs.Eq(trs.Size() - 2).Find("td").Eq(1).Text())
		res.Deuterium = utils.ParseInt(trs.Eq(trs.Size() - 1).Find("td").Eq(1).Text())

		tmp = append(tmp, Tmp{fleet: fleet, res: res})
	})

	for _, t := range tmp {
		res = append(res, t.fleet)
	}

	return res
}

func extractIPMFromDoc(doc *goquery.Document) (duration, max int64, token string) {
	duration = utils.DoParseI64(doc.Find("span#timer").AttrOr("data-duration", "0"))
	max = utils.DoParseI64(doc.Find("input[name=anz]").AttrOr("data-max", "0"))
	token = doc.Find("input[name=token]").AttrOr("value", "")
	return
}

func extractFleetsFromDoc(doc *goquery.Document, location *time.Location, lifeformEnabled bool) (res []ogame.Fleet) {
	res = make([]ogame.Fleet, 0)
	script := doc.Find("body script").Text()
	doc.Find("div.fleetDetails").Each(func(i int, s *goquery.Selection) {
		originText := s.Find("span.originCoords a").Text()
		origin := ExtractCoord(originText)
		origin.Type = ogame.PlanetType
		if s.Find("span.originPlanet figure").HasClass("moon") {
			origin.Type = ogame.MoonType
		}

		destText := s.Find("span.destinationCoords a").Text()
		dest := ExtractCoord(destText)
		dest.Type = ogame.PlanetType
		if s.Find("span.destinationPlanet figure").HasClass("moon") {
			dest.Type = ogame.MoonType
		} else if s.Find("span.destinationPlanet figure").HasClass("tf") {
			dest.Type = ogame.DebrisType
		}

		id := utils.DoParseI64(s.Find("a.openCloseDetails").AttrOr("data-mission-id", "0"))

		timerID := s.Find("span.timer").AttrOr("id", "")
		m := regexp.MustCompile(`getElementByIdWithCache\("` + timerID + `"\),\s*(\d+),`).FindStringSubmatch(script)
		var arriveIn int64
		if len(m) == 2 {
			arriveIn = utils.DoParseI64(m[1])
		}

		timerNextID := s.Find("span.nextTimer").AttrOr("id", "")
		m = regexp.MustCompile(`getElementByIdWithCache\("` + timerNextID + `"\),\s*(\d+)\s*\);`).FindStringSubmatch(script)
		var backIn int64
		if len(m) == 2 {
			backIn = utils.DoParseI64(m[1])
		}

		missionType := utils.DoParseI64(s.AttrOr("data-mission-type", ""))
		returnFlight, _ := strconv.ParseBool(s.AttrOr("data-return-flight", ""))
		inDeepSpace := s.Find("span.fleetDetailButton a").HasClass("fleet_icon_forward_end")
		arrivalTime := utils.DoParseI64(s.AttrOr("data-arrival-time", ""))
		endTime := utils.DoParseI64(s.Find("a.openCloseDetails").AttrOr("data-end-time", ""))

		trs := s.Find("table.fleetinfo tr")
		shipment := ogame.Resources{}
		metalTrOffset := 3
		crystalTrOffset := 2
		DeuteriumTrOffset := 1
		if lifeformEnabled {
			metalTrOffset = 4
			crystalTrOffset = 3
			DeuteriumTrOffset = 2
		}
		shipment.Metal = utils.ParseInt(trs.Eq(trs.Size() - metalTrOffset).Find("td").Eq(1).Text())
		shipment.Crystal = utils.ParseInt(trs.Eq(trs.Size() - crystalTrOffset).Find("td").Eq(1).Text())
		shipment.Deuterium = utils.ParseInt(trs.Eq(trs.Size() - DeuteriumTrOffset).Find("td").Eq(1).Text())

		fedAttackHref := s.Find("span.fedAttack a").AttrOr("href", "")
		fedAttackURL, _ := url.Parse(fedAttackHref)
		fedAttackQuery := fedAttackURL.Query()
		targetPlanetID := utils.DoParseI64(fedAttackQuery.Get("target"))
		unionID := utils.DoParseI64(fedAttackQuery.Get("union"))

		fleet := ogame.Fleet{}
		fleet.ID = ogame.FleetID(id)
		fleet.Origin = origin
		fleet.Destination = dest
		fleet.Mission = ogame.MissionID(missionType)
		fleet.ReturnFlight = returnFlight
		fleet.InDeepSpace = inDeepSpace
		fleet.Resources = shipment
		fleet.TargetPlanetID = targetPlanetID
		fleet.UnionID = unionID
		fleet.ArrivalTime = time.Unix(endTime, 0)
		fleet.BackTime = time.Unix(arrivalTime, 0)

		var startTimeString string
		var startTimeStringExists bool
		if !returnFlight {
			fleet.ArriveIn = arriveIn
			fleet.BackIn = backIn
			startTimeString, startTimeStringExists = s.Find("div.origin img").Attr("title")
		} else {
			fleet.ArriveIn = -1
			fleet.BackIn = arriveIn
			startTimeString, startTimeStringExists = s.Find("div.destination img").Attr("title")
		}

		var startTime time.Time
		if startTimeStringExists {
			startTimeArray := strings.Split(startTimeString, ":| ")
			if len(startTimeArray) == 2 {
				startTime, _ = time.ParseInLocation("02.01.2006<br>15:04:05", startTimeArray[1], location)
			}
		}
		fleet.StartTime = startTime.Local()

		for i := 1; i < trs.Size()-5; i++ {
			tds := trs.Eq(i).Find("td")
			name := strings.ToLower(strings.Trim(strings.TrimSpace(tds.Eq(0).Text()), ":"))
			qty := utils.ParseInt(tds.Eq(1).Text())
			shipID := ogame.ShipName2ID(name)
			fleet.Ships.Set(shipID, qty)
		}

		res = append(res, fleet)
	})
	return
}

func extractSlotsFromDoc(doc *goquery.Document) ogame.Slots {
	slots := ogame.Slots{}
	page := ExtractBodyIDFromDoc(doc)
	if page == "movement" {
		slots.InUse = utils.ParseInt(doc.Find("span.fleetSlots > span.current").Text())
		slots.Total = utils.ParseInt(doc.Find("span.fleetSlots > span.all").Text())
		slots.ExpInUse = utils.ParseInt(doc.Find("span.expSlots > span.current").Text())
		slots.ExpTotal = utils.ParseInt(doc.Find("span.expSlots > span.all").Text())
	} else if page == "fleetdispatch" || page == "fleet1" {
		r := regexp.MustCompile(`(\d+)/(\d+)`)
		txt := doc.Find("div#slots>div").Eq(0).Text()
		m := r.FindStringSubmatch(txt)
		if len(m) == 3 {
			slots.InUse = utils.DoParseI64(m[1])
			slots.Total = utils.DoParseI64(m[2])
		}
		txt = doc.Find("div#slots>div").Eq(1).Text()
		m = r.FindStringSubmatch(txt)
		if len(m) == 3 {
			slots.ExpInUse = utils.DoParseI64(m[1])
			slots.ExpTotal = utils.DoParseI64(m[2])
		}
	}
	return slots
}

func extractServerTimeFromDoc(doc *goquery.Document) (time.Time, error) {
	txt := doc.Find("li.OGameClock").First().Text()
	serverTime, err := time.Parse("02.01.2006 15:04:05", txt)
	if err != nil {
		return time.Time{}, err
	}

	u1 := time.Now().UTC().Unix()
	u2 := serverTime.Unix()
	n := int(math.Round(float64(u2-u1)/900)) * 900 // u2-u1 should be close to 0, round to nearest 15min difference

	serverTime = serverTime.Add(time.Duration(-n) * time.Second).In(time.FixedZone("OGT", n))

	return serverTime, nil
}

func extractSpioAnzFromDoc(doc *goquery.Document) int64 {
	out := utils.DoParseI64(doc.Find("input[name=spio_anz]").AttrOr("value", "1"))
	return out
}

func extractDisableChatBarFromDoc(doc *goquery.Document) bool {
	_, exists := doc.Find("input[name=disableChatBar]").Attr("checked")
	return exists
}

func extractDisableOutlawWarningFromDoc(doc *goquery.Document) bool {
	_, exists := doc.Find("input[name=disableOutlawWarning]").Attr("checked")
	return exists
}

func extractMobileVersionFromDoc(doc *goquery.Document) bool {
	_, exists := doc.Find("input[name=mobileVersion]").Attr("checked")
	return exists
}

func extractUrlaubsModus(doc *goquery.Document) bool {
	_, exists := doc.Find("input[name=urlaubs_modus]").Attr("checked")
	return exists
}

func extractShowOldDropDownsFromDoc(doc *goquery.Document) bool {
	_, exists := doc.Find("input[name=showOldDropDowns]").Attr("checked")
	return exists
}

func extractActivateAutofocusFromDoc(doc *goquery.Document) bool {
	_, exists := doc.Find("input[name=activateAutofocus]").Attr("checked")
	return exists
}

func extractEventsShowFromDoc(doc *goquery.Document) int64 {
	return utils.DoParseI64(doc.Find("select[name=eventsShow] option[selected]").AttrOr("value", "1"))
}

func extractSortSettingFromDoc(doc *goquery.Document) int64 {
	return utils.DoParseI64(doc.Find("select#sortSetting option[selected]").AttrOr("value", "0"))
}

func extractSortOrderFromDoc(doc *goquery.Document) int64 {
	return utils.DoParseI64(doc.Find("select#sortOrder option[selected]").AttrOr("value", "0"))
}

func extractShowDetailOverlayFromDoc(doc *goquery.Document) bool {
	_, exists := doc.Find("input[name=showDetailOverlay]").Attr("checked")
	return exists
}

func extractAnimatedSlidersFromDoc(doc *goquery.Document) bool {
	_, exists := doc.Find("input[name=animatedSliders]").Attr("checked")
	return exists
}

func extractAnimatedOverviewFromDoc(doc *goquery.Document) bool {
	_, exists := doc.Find("input[name=animatedOverview]").Attr("checked")
	return exists
}

func extractPopupsNoticesFromDoc(doc *goquery.Document) bool {
	_, exists := doc.Find(`input[name="popups[notices]"]`).Attr("checked")
	return exists
}

func extractPopopsCombatreportFromDoc(doc *goquery.Document) bool {
	_, exists := doc.Find(`input[name="popups[combatreport]"]`).Attr("checked")
	return exists
}

func extractSpioReportPicturesFromDoc(doc *goquery.Document) bool {
	_, exists := doc.Find("input[name=spioReportPictures]").Attr("checked")
	return exists
}

func extractMsgResultsPerPageFromDoc(doc *goquery.Document) int64 {
	return utils.DoParseI64(doc.Find("select[name=msgResultsPerPage] option[selected]").AttrOr("value", "10"))
}

func extractAuctioneerNotificationsFromDoc(doc *goquery.Document) bool {
	_, exists := doc.Find("input[name=auctioneerNotifications]").Attr("checked")
	return exists
}

func extractEconomyNotificationsFromDoc(doc *goquery.Document) bool {
	_, exists := doc.Find("input[name=economyNotifications]").Attr("checked")
	return exists
}

func extractShowActivityMinutesFromDoc(doc *goquery.Document) bool {
	_, exists := doc.Find("input[name=showActivityMinutes]").Attr("checked")
	return exists
}

func extractPreserveSystemOnPlanetChangeFromDoc(doc *goquery.Document) bool {
	_, exists := doc.Find("input[name=preserveSystemOnPlanetChange]").Attr("checked")
	return exists
}

func extractNotifBuildListFromDoc(doc *goquery.Document) bool {
	_, exists := doc.Find(`input[name="notifications[buildList]"]`).Attr("checked")
	return exists
}

func extractNotifFriendlyFleetActivitiesFromDoc(doc *goquery.Document) bool {
	_, exists := doc.Find(`input[name="notifications[friendlyFleetActivities]"]`).Attr("checked")
	return exists
}

func extractNotifHostileFleetActivitiesFromDoc(doc *goquery.Document) bool {
	_, exists := doc.Find(`input[name="notifications[hostileFleetActivities]"]`).Attr("checked")
	return exists
}

func extractNotifForeignEspionageFromDoc(doc *goquery.Document) bool {
	_, exists := doc.Find(`input[name="notifications[foreignEspionage]"]`).Attr("checked")
	return exists
}

func extractNotifAllianceBroadcastsFromDoc(doc *goquery.Document) bool {
	_, exists := doc.Find(`input[name="notifications[allianceBroadcasts]"]`).Attr("checked")
	return exists
}

func extractNotifAllianceMessagesFromDoc(doc *goquery.Document) bool {
	_, exists := doc.Find(`input[name="notifications[allianceMessages]"]`).Attr("checked")
	return exists
}

func extractNotifAuctionsFromDoc(doc *goquery.Document) bool {
	_, exists := doc.Find(`input[name="notifications[auctions]"]`).Attr("checked")
	return exists
}

func extractNotifAccountFromDoc(doc *goquery.Document) bool {
	_, exists := doc.Find(`input[name="notifications[account]"]`).Attr("checked")
	return exists
}

func extractCommanderFromDoc(doc *goquery.Document) bool {
	return doc.Find("div#officers a.commander").HasClass("on")
}

func extractAdmiralFromDoc(doc *goquery.Document) bool {
	return doc.Find("div#officers a.admiral").HasClass("on")
}

func extractEngineerFromDoc(doc *goquery.Document) bool {
	return doc.Find("div#officers a.engineer").HasClass("on")
}

func extractGeologistFromDoc(doc *goquery.Document) bool {
	return doc.Find("div#officers a.geologist").HasClass("on")
}

func extractTechnocratFromDoc(doc *goquery.Document) bool {
	return doc.Find("div#officers a.technocrat").HasClass("on")
}

func extractAbandonInformation(doc *goquery.Document) (string, string) {
	abandonToken := doc.Find("form#planetMaintenanceDelete input[name=abandon]").AttrOr("value", "")
	token := doc.Find("form#planetMaintenanceDelete input[name=token]").AttrOr("value", "")
	return abandonToken, token
}

func extractPlanetCoordinate(pageHTML []byte) (ogame.Coordinate, error) {
	m := regexp.MustCompile(`<meta name="ogame-planet-coordinates" content="(\d+):(\d+):(\d+)"/>`).FindSubmatch(pageHTML)
	if len(m) == 0 {
		return ogame.Coordinate{}, errors.New("planet coordinate not found")
	}
	galaxy := utils.DoParseI64(string(m[1]))
	system := utils.DoParseI64(string(m[2]))
	position := utils.DoParseI64(string(m[3]))
	planetType, _ := extractPlanetType(pageHTML)
	return ogame.Coordinate{galaxy, system, position, planetType}, nil
}

func extractTearDownToken(pageHTML []byte) (string, error) {
	m := regexp.MustCompile(`modus=3&token=([^&]+)&`).FindSubmatch(pageHTML)
	if len(m) != 2 {
		return "", errors.New("unable to find tear down token")
	}
	return string(m[1]), nil
}

func extractPlanetID(pageHTML []byte) (ogame.CelestialID, error) {
	m := regexp.MustCompile(`<meta name="ogame-planet-id" content="(\d+)"/>`).FindSubmatch(pageHTML)
	if len(m) == 0 {
		return 0, errors.New("planet id not found")
	}
	planetID := utils.DoParseI64(string(m[1]))
	return ogame.CelestialID(planetID), nil
}

func extractPlanetIDFromDoc(doc *goquery.Document) (ogame.CelestialID, error) {
	planetID := utils.DoParseI64(doc.Find("meta[name=ogame-planet-id]").AttrOr("content", "0"))
	if planetID == 0 {
		return 0, errors.New("planet id not found")
	}
	return ogame.CelestialID(planetID), nil
}

func extractOverviewShipSumCountdownFromBytes(pageHTML []byte) int64 {
	var shipSumCountdown int64
	shipSumCountdownMatch := regexp.MustCompile(`getElementByIdWithCache\('shipSumCount7'\),\d+,\d+,(\d+),`).FindSubmatch(pageHTML)
	if len(shipSumCountdownMatch) > 0 {
		shipSumCountdown = int64(utils.ToInt(shipSumCountdownMatch[1]))
	}
	return shipSumCountdown
}

func extractOGameTimestampFromBytes(pageHTML []byte) int64 {
	m := regexp.MustCompile(`<meta name="ogame-timestamp" content="(\d+)"/>`).FindSubmatch(pageHTML)
	if len(m) != 2 {
		return 0
	}
	ts := utils.DoParseI64(string(m[1]))
	return ts
}

func extractPlanetType(pageHTML []byte) (ogame.CelestialType, error) {
	m := regexp.MustCompile(`<meta name="ogame-planet-type" content="(\w+)"/>`).FindSubmatch(pageHTML)
	if len(m) == 0 {
		return 0, errors.New("planet type not found")
	}
	if bytes.Equal(m[1], []byte("planet")) {
		return ogame.PlanetType, nil
	} else if bytes.Equal(m[1], []byte("moon")) {
		return ogame.MoonType, nil
	}
	return 0, errors.New("invalid planet type : " + string(m[1]))
}

func extractPlanetTypeFromDoc(doc *goquery.Document) (ogame.CelestialType, error) {
	planetType := doc.Find("meta[name=ogame-planet-type]").AttrOr("content", "")
	if planetType == "" {
		return 0, errors.New("planet type not found")
	}
	if planetType == "planet" {
		return ogame.PlanetType, nil
	} else if planetType == "moon" {
		return ogame.MoonType, nil
	}
	return 0, errors.New("invalid planet type : " + planetType)
}

func extractAjaxChatToken(pageHTML []byte) (string, error) {
	r1 := regexp.MustCompile(`ajaxChatToken\s?=\s?['"](\w+)['"]`)
	m1 := r1.FindSubmatch(pageHTML)
	if len(m1) < 2 {
		return "", errors.New("unable to find token")
	}
	token := string(m1[1])
	return token, nil
}

func extractUserInfos(pageHTML []byte, lang string) (ogame.UserInfos, error) {
	playerIDRgx := regexp.MustCompile(`<meta name="ogame-player-id" content="(\d+)"/>`)
	playerNameRgx := regexp.MustCompile(`<meta name="ogame-player-name" content="([^"]+)"/>`)
	txtContent := regexp.MustCompile(`textContent\[7]\s?=\s?"([^"]+)"`)
	playerIDGroups := playerIDRgx.FindSubmatch(pageHTML)
	playerNameGroups := playerNameRgx.FindSubmatch(pageHTML)
	subHTMLGroups := txtContent.FindSubmatch(pageHTML)
	if len(playerIDGroups) < 2 {
		return ogame.UserInfos{}, errors.New("cannot find player id")
	}
	if len(playerNameGroups) < 2 {
		return ogame.UserInfos{}, errors.New("cannot find player name")
	}
	if len(subHTMLGroups) < 2 {
		return ogame.UserInfos{}, errors.New("cannot find sub html")
	}
	res := ogame.UserInfos{}
	res.PlayerID = int64(utils.ToInt(playerIDGroups[1]))
	res.PlayerName = string(playerNameGroups[1])
	html2 := []byte(strings.ReplaceAll(string(subHTMLGroups[1]), ",", "."))

	infosRgx := regexp.MustCompile(`([\d\\.]+) \(Place ([\d.]+) of ([\d.]+)\)`)
	switch lang {
	case "fr":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(Place ([\d.]+) sur ([\d.]+)\)`)
	case "hu":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(Helyez\\u00e9s \\/ J\\u00e1t\\u00e9kosok: ([\d.]+) \\/ ([\d.]+)\)`)
	case "si":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(Mesto ([\d.]+) od ([\d.]+)\)`)
	case "sk":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(Umiestnenie v rebr\\u00ed\\u010dku: ([\d.]+) z ([\d.]+)\)`)
	case "no":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(Plass ([\d.]+) av ([\d.]+)\)`)
	case "hr":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(Mjesto ([\d.]+) od ([\d.]+)\)`)
	case "gr":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(\\u039a\\u03b1\\u03c4\\u03ac\\u03c4\\u03b1\\u03be\\u03b7 ([\d.]+) \\u03b1\\u03c0\\u03cc ([\d.]+)\)`)
	case "tw":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(([\d.]+) \\u4eba\\u4e2d\\u7684\\u7b2c ([\d.]+) \\u4f4d\)`)
	case "cz":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(Pozice ([\d.]+) z ([\d.]+)\)`)
	case "de":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(Platz ([\d.]+) von ([\d.]+)\)`)
	case "es":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(Posici\\u00f3n ([\d.]+) de ([\d.]+)\)`)
	case "ar":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(Lugar ([\d.]+) de ([\d.]+)\)`)
	case "mx":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(Lugar ([\d.]+) de ([\d.]+)\)`)
	case "br":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(Posi\\u00e7\\u00e3o ([\d.]+) de ([\d.]+)\)`)
	case "it":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(Posizione ([\d.]+) su ([\d.]+)\)`)
	case "jp":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(([\d.]+)\\u4eba\\u4e2d([\d.]+)\\u4f4d\)`)
	case "pl":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(Miejsce ([\d.]+) z ([\d.]+)\)`)
	case "tr":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(([\d.]+) oyuncu i\\u00e7inde ([\d.]+)\. s\\u0131rada\)`)
	case "pt":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(Posi\\u00e7\\u00e3o ([\d.]+) de ([\d.]+)\)`)
	case "nl":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(Plaats ([\d.]+) van ([\d.]+)\)`)
	case "dk":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(Placering ([\d.]+) af ([\d.]+)\)`)
	case "ro":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(Locul ([\d.]+) din ([\d.]+)\)`)
	case "fi":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(Sijoitus ([\d.]+) kaikista pelaajista ([\d.]+)\)`)
	case "ba":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(Mjesto ([\d.]+) od ([\d.]+)\)`)
	case "ru":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(\\u041c\\u0435\\u0441\\u0442\\u043e ([\d.]+) \\u0438\\u0437 ([\d.]+)\)`)
	}
	// pl: 0 (Miejsce 5.872 z 5.875)
	// fr: 0 (Place 3.197 sur 3.348)
	// de: 0 (Platz 2.979 von 2.980)
	// jp: 0 (73人中72位)
	// pt: 0 (Posição 1.861 de 1.862
	infos := infosRgx.FindSubmatch(html2)
	if len(infos) < 4 {
		return ogame.UserInfos{}, errors.New("cannot find infos in sub html")
	}
	res.Points = utils.ParseInt(string(infos[1]))
	res.Rank = utils.ParseInt(string(infos[2]))
	res.Total = utils.ParseInt(string(infos[3]))
	if lang == "tr" || lang == "jp" {
		res.Rank = utils.ParseInt(string(infos[3]))
		res.Total = utils.ParseInt(string(infos[2]))
	}
	honourPointsRgx := regexp.MustCompile(`textContent\[9]\s?=\s?"([^"]+)"`)
	honourPointsGroups := honourPointsRgx.FindSubmatch(pageHTML)
	if len(honourPointsGroups) < 2 {
		return ogame.UserInfos{}, errors.New("cannot find honour points")
	}
	res.HonourPoints = utils.ParseInt(string(honourPointsGroups[1]))
	return res, nil
}

func IsLogged(pageHTML []byte) bool {
	return len(regexp.MustCompile(`<meta name="ogame-session" content="\w+"/>`).FindSubmatch(pageHTML)) == 1 ||
		len(regexp.MustCompile(`var session = "\w+"`).FindSubmatch(pageHTML)) == 1
}

func extractResourcesDetails(pageHTML []byte) (out ogame.ResourcesDetails, err error) {
	var res ogame.ResourcesResp
	if err = json.Unmarshal(pageHTML, &res); err != nil {
		if IsLogged(pageHTML) {
			return out, ogame.ErrInvalidPlanetID
		}
		return
	}
	out.Metal.Available = res.Metal.Resources.Actual
	out.Metal.StorageCapacity = res.Metal.Resources.Max
	out.Crystal.Available = res.Crystal.Resources.Actual
	out.Crystal.StorageCapacity = res.Crystal.Resources.Max
	out.Deuterium.Available = res.Deuterium.Resources.Actual
	out.Deuterium.StorageCapacity = res.Deuterium.Resources.Max
	out.Energy.Available = res.Energy.Resources.Actual
	out.Darkmatter.Available = res.Darkmatter.Resources.Actual
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

func ExtractCoord(v string) (coord ogame.Coordinate) {
	coordRgx := regexp.MustCompile(`\[(\d+):(\d+):(\d+)]`)
	m := coordRgx.FindStringSubmatch(v)
	if len(m) == 4 {
		coord.Galaxy = utils.DoParseI64(m[1])
		coord.System = utils.DoParseI64(m[2])
		coord.Position = utils.DoParseI64(m[3])
	}
	return
}

func extractGalaxyInfos(pageHTML []byte, botPlayerName string, botPlayerID, botPlayerRank int64) (ogame.SystemInfos, error) {
	prefixedNumRgx := regexp.MustCompile(`.*: ([\d.,]+)`)

	extractActivity := func(activityDiv *goquery.Selection) int64 {
		var activity int64
		if activityDiv != nil {
			activityDivClass := activityDiv.AttrOr("class", "")
			if strings.Contains(activityDivClass, "minute15") {
				activity = 15
			} else if strings.Contains(activityDivClass, "showMinutes") {
				activity = utils.DoParseI64(strings.TrimSpace(activityDiv.Text()))
			}
		}
		return activity
	}

	var tmp struct {
		Galaxy string
	}
	var res ogame.SystemInfos
	if err := json.Unmarshal(pageHTML, &tmp); err != nil {
		return res, ogame.ErrNotLogged
	}

	overlayTokenRgx := regexp.MustCompile(`data-overlay-token="([^"]+)"`)
	m := overlayTokenRgx.FindStringSubmatch(tmp.Galaxy)
	if len(m) == 2 {
		res.OverlayToken = m[1]
	}

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(tmp.Galaxy))
	res.Tmpgalaxy = utils.ParseInt(doc.Find("table").AttrOr("data-galaxy", "0"))
	res.Tmpsystem = utils.ParseInt(doc.Find("table").AttrOr("data-system", "0"))
	isVacationMode := doc.Find("div#warning").Length() == 1
	if isVacationMode {
		return res, ogame.ErrAccountInVacationMode
	}
	isMobile := doc.Find("span.fright span#filter_empty").Length() == 0
	if isMobile {
		return res, ogame.ErrMobileView
	}
	doc.Find("tr.row").Each(func(i int, s *goquery.Selection) {
		classes, _ := s.Attr("class")
		if !strings.Contains(classes, "empty_filter") {
			position := s.Find("td.position").Text()

			tooltips := s.Find("div.htmlTooltip")
			planetTooltip := tooltips.First()
			planetName := planetTooltip.Find("h1").Find("span").Text()
			planetImg, _ := planetTooltip.Find("img").Attr("src")
			coordsRaw := planetTooltip.Find("span#pos-planet").Text()

			metalTxt := s.Find("div#debris" + position + " ul.ListLinks li").First().Text()
			crystalTxt := s.Find("div#debris" + position + " ul.ListLinks li").Eq(1).Text()
			recyclersTxt := s.Find("div#debris" + position + " ul.ListLinks li").Eq(2).Text()

			planetInfos := new(ogame.PlanetInfos)
			planetInfos.ID = utils.DoParseI64(s.Find("td.colonized").AttrOr("data-planet-id", ""))

			moonID := utils.DoParseI64(s.Find("td.moon").AttrOr("data-moon-id", ""))
			moonSize := utils.DoParseI64(strings.Split(s.Find("td.moon span#moonsize").Text(), " ")[0])
			if moonID > 0 {
				planetInfos.Moon = new(ogame.MoonInfos)
				planetInfos.Moon.ID = moonID
				planetInfos.Moon.Diameter = moonSize
				planetInfos.Moon.Activity = extractActivity(s.Find("td.moon div.activity"))
			}

			allianceSpan := s.Find("span.allytagwrapper")
			if allianceSpan.Size() > 0 {
				longID, _ := allianceSpan.Attr("rel")
				planetInfos.Alliance = new(ogame.AllianceInfos)
				planetInfos.Alliance.Name = allianceSpan.Find("h1").Text()
				planetInfos.Alliance.ID = utils.DoParseI64(strings.TrimPrefix(longID, "alliance"))
				planetInfos.Alliance.Rank = utils.DoParseI64(allianceSpan.Find("ul.ListLinks li").First().Find("a").Text())
				planetInfos.Alliance.Member = utils.ParseInt(prefixedNumRgx.FindStringSubmatch(allianceSpan.Find("ul.ListLinks li").Eq(1).Text())[1])
			}

			if len(prefixedNumRgx.FindStringSubmatch(metalTxt)) > 0 {
				planetInfos.Debris.Metal = utils.ParseInt(prefixedNumRgx.FindStringSubmatch(metalTxt)[1])
				planetInfos.Debris.Crystal = utils.ParseInt(prefixedNumRgx.FindStringSubmatch(crystalTxt)[1])
				planetInfos.Debris.RecyclersNeeded = utils.ParseInt(prefixedNumRgx.FindStringSubmatch(recyclersTxt)[1])
			}

			planetInfos.Activity = extractActivity(s.Find("td:not(.moon) div.activity"))
			planetInfos.Name = planetName
			planetInfos.Img = planetImg
			planetInfos.Inactive = strings.Contains(classes, "inactive_filter")
			planetInfos.StrongPlayer = s.Find("span.status_abbr_strong").Size() > 0
			planetInfos.Newbie = strings.Contains(classes, "newbie_filter")
			planetInfos.Vacation = strings.Contains(classes, "vacation_filter")
			planetInfos.HonorableTarget = s.Find("span.status_abbr_honorableTarget").Size() > 0
			planetInfos.Administrator = s.Find("span.status_abbr_admin").Size() > 0
			planetInfos.Banned = s.Find("td.playername a span.status_abbr_banned").Size() > 0
			tdPlayername := s.Find("td.playername span")
			planetInfos.Player.IsBandit = tdPlayername.HasClass("rank_bandit1") || tdPlayername.HasClass("rank_bandit2") || tdPlayername.HasClass("rank_bandit3")
			planetInfos.Player.IsStarlord = tdPlayername.HasClass("rank_starlord1") || tdPlayername.HasClass("rank_starlord2") || tdPlayername.HasClass("rank_starlord3")
			planetInfos.Coordinate = ExtractCoord(coordsRaw)
			planetInfos.Coordinate.Type = ogame.PlanetType
			planetInfos.Date = time.Now()

			var playerID int64
			var playerName string
			var playerRank int64
			if len(tooltips.Nodes) > 1 {
				tooltips.Each(func(i int, s *goquery.Selection) {
					idAttr, _ := s.Attr("id")
					if strings.HasPrefix(idAttr, "player") {
						playerID = utils.DoParseI64(regexp.MustCompile(`player(\d+)`).FindStringSubmatch(idAttr)[1])
						playerName = s.Find("h1").Find("span").Text()
						playerRank = utils.DoParseI64(s.Find("li.rank").Find("a").Text())
					}
				})
			}
			if playerName == "" {
				playerName := strings.TrimSpace(s.Find("td.playername").Find("span").Text())
				if playerName == "" {
					planetInfos.Destroyed = true
				}
			}

			if !planetInfos.Destroyed && playerID == 0 {
				playerID = botPlayerID
				playerName = botPlayerName
				playerRank = botPlayerRank
			}

			planetInfos.Player.ID = playerID
			planetInfos.Player.Name = playerName
			planetInfos.Player.Rank = playerRank

			res.Tmpplanets[i] = planetInfos
		}
	})

	debris16Div := doc.Find("div#debris16")
	if debris16Div.Size() > 0 {
		lis := debris16Div.Find("ul.ListLinks li")
		metalTxt := lis.First().Text()
		crystalTxt := lis.Eq(1).Text()
		pathfindersTxt := lis.Eq(2).Text()
		res.ExpeditionDebris.Metal = utils.ParseInt(prefixedNumRgx.FindStringSubmatch(metalTxt)[1])
		res.ExpeditionDebris.Crystal = utils.ParseInt(prefixedNumRgx.FindStringSubmatch(crystalTxt)[1])
		res.ExpeditionDebris.PathfindersNeeded = utils.ParseInt(prefixedNumRgx.FindStringSubmatch(pathfindersTxt)[1])
	}

	debris17Div := doc.Find("div#debris17")
	if debris17Div.Size() > 0 {
		lis := debris17Div.Find("ul.ListLinks li")
		darkmatterTxt := lis.First().Text()
		darkmatterMatches := prefixedNumRgx.FindStringSubmatch(darkmatterTxt)
		if len(darkmatterMatches) == 2 {
			res.Events.Darkmatter = utils.ParseInt(darkmatterMatches[1])
		}
	}

	planet17Div := doc.Find("div#planet17")
	if planet17Div.Size() > 0 {
		res.Events.HasAsteroid = true
	}

	return res, nil
}

func extractPhalanx(pageHTML []byte) ([]ogame.Fleet, error) {
	res := make([]ogame.Fleet, 0)
	var ogameTimestamp int64
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	eventFleet := doc.Find("div.eventFleet")
	if eventFleet.Size() == 0 {
		txt := strings.TrimSpace(doc.Find("div#phalanxEventContent").Text())
		// TODO: 'fleet' and 'deuterium' won't work in other languages
		if strings.Contains(txt, "fleet") {
			return res, nil
		} else if strings.Contains(txt, "deuterium") {
			return res, errors.New(strings.TrimSpace(txt))
		}
		return res, errors.New(txt)
	}

	m := regexp.MustCompile(`var mytime = ([0-9]+)`).FindSubmatch(pageHTML)
	if len(m) > 0 {
		ogameTimestamp = utils.DoParseI64(string(m[1]))
	}

	eventFleet.Each(func(i int, s *goquery.Selection) {
		mission := utils.DoParseI64(s.AttrOr("data-mission-type", "0"))
		returning, _ := strconv.ParseBool(s.AttrOr("data-return-flight", "false"))
		arrivalTime := utils.DoParseI64(s.AttrOr("data-arrival-time", "0"))
		arriveIn := arrivalTime - ogameTimestamp
		if arriveIn < 0 {
			arriveIn = 0
		}
		originFleetFigure := s.Find("li.originFleet figure")
		originTxt := s.Find("li.coordsOrigin a").Text()
		destTxt := s.Find("li.destCoords a").Text()

		fleet := ogame.Fleet{}

		if movement, exists := s.Find("li.detailsFleet span").Attr("title"); exists {
			root, err := html.Parse(strings.NewReader(movement))
			if err != nil {
				return
			}
			doc2 := goquery.NewDocumentFromNode(root)
			doc2.Find("tr").Each(func(i int, s *goquery.Selection) {
				if i == 0 {
					return
				}
				name := s.Find("td").Eq(0).Text()
				nbr := utils.ParseInt(s.Find("td").Eq(1).Text())
				if name != "" && nbr > 0 {
					fleet.Ships.Set(ogame.ShipName2ID(name), nbr)
				}
			})
		}

		fleet.Mission = ogame.MissionID(mission)
		fleet.ReturnFlight = returning
		fleet.ArriveIn = arriveIn
		fleet.ArrivalTime = time.Unix(arrivalTime, 0)
		fleet.Origin = ExtractCoord(originTxt)
		fleet.Origin.Type = ogame.PlanetType
		if originFleetFigure.HasClass("moon") {
			fleet.Origin.Type = ogame.MoonType
		}
		fleet.Destination = ExtractCoord(destTxt)
		fleet.Destination.Type = ogame.PlanetType
		res = append(res, fleet)
	})
	return res, nil
}

func extractJumpGate(pageHTML []byte) (ogame.ShipsInfos, string, []ogame.MoonID, int64) {
	m := regexp.MustCompile(`\$\("#cooldown"\), (\d+),`).FindSubmatch(pageHTML)
	ships := ogame.ShipsInfos{}
	var destinations []ogame.MoonID
	if len(m) > 0 {
		waitTime := int64(utils.ToInt(m[1]))
		return ships, "", destinations, waitTime
	}
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	for _, s := range ogame.Ships {
		ships.Set(s.GetID(), utils.ParseInt(doc.Find("input#ship_"+utils.FI64(s.GetID())).AttrOr("rel", "0")))
	}
	token := doc.Find("input[name=token]").AttrOr("value", "")

	doc.Find("select[name=zm] option").Each(func(i int, s *goquery.Selection) {
		moonID := utils.ParseInt(s.AttrOr("value", "0"))
		if moonID > 0 {
			destinations = append(destinations, ogame.MoonID(moonID))
		}
	})

	return ships, token, destinations, 0
}

func extractFederation(pageHTML []byte) url.Values {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	payload := extractHiddenFieldsFromDoc(doc)
	groupName := doc.Find("input#groupNameInput").AttrOr("value", "")
	doc.Find("ul#participantselect li").Each(func(i int, s *goquery.Selection) {
		payload.Add("unionUsers", s.Text())
	})
	payload.Add("groupname", groupName)
	return payload
}

func extractConstructions(pageHTML []byte) (buildingID ogame.ID, buildingCountdown int64, researchID ogame.ID, researchCountdown int64, lfBuildingID ogame.ID, lfBuildingCountdown int64) {
	buildingCountdownMatch := regexp.MustCompile(`getElementByIdWithCache\("Countdown"\),(\d+),`).FindSubmatch(pageHTML)
	if len(buildingCountdownMatch) > 0 {
		buildingCountdown = int64(utils.ToInt(buildingCountdownMatch[1]))
		buildingIDInt := utils.ToInt(regexp.MustCompile(`onclick="cancelProduction\((\d+),`).FindSubmatch(pageHTML)[1])
		buildingID = ogame.ID(buildingIDInt)
	}
	researchCountdownMatch := regexp.MustCompile(`getElementByIdWithCache\("researchCountdown"\),(\d+),`).FindSubmatch(pageHTML)
	if len(researchCountdownMatch) > 0 {
		researchCountdown = int64(utils.ToInt(researchCountdownMatch[1]))
		researchIDInt := utils.ToInt(regexp.MustCompile(`onclick="cancelResearch\((\d+),`).FindSubmatch(pageHTML)[1])
		researchID = ogame.ID(researchIDInt)
	}
	return
}

func extractFleetDeutSaveFactor(pageHTML []byte) float64 {
	factor := 1.0
	m := regexp.MustCompile(`var fleetDeutSaveFactor=([+-]?([0-9]*[.])?[0-9]+);`).FindSubmatch(pageHTML)
	if len(m) > 0 {
		factor, _ = strconv.ParseFloat(string(m[1]), 64)
	}
	return factor
}

func extractCancelBuildingInfos(pageHTML []byte) (token string, techID, listID int64, err error) {
	r1 := regexp.MustCompile(`page=overview&modus=2&token=(\w+)&techid="\+cancelProduction_id\+"&listid="\+production_listid`)
	m1 := r1.FindSubmatch(pageHTML)
	if len(m1) < 2 {
		return "", 0, 0, errors.New("unable to find token")
	}
	token = string(m1[1])
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	t := doc.Find("table.construction").Eq(0)
	a, _ := t.Find("a.abortNow").Attr("onclick")
	r := regexp.MustCompile(`cancelProduction\((\d+),\s?(\d+),`)
	m := r.FindStringSubmatch(a)
	if len(m) < 3 {
		return "", 0, 0, errors.New("unable to find techid/listid")
	}
	techID = utils.DoParseI64(m[1])
	listID = utils.DoParseI64(m[2])
	return
}

func extractCancelResearchInfos(pageHTML []byte) (token string, techID, listID int64, err error) {
	r1 := regexp.MustCompile(`page=overview&modus=2&token=(\w+)"\+"&techid="\+id\+"&listid="\+listId`)
	m1 := r1.FindSubmatch(pageHTML)
	if len(m1) < 2 {
		return "", 0, 0, errors.New("unable to find token")
	}
	token = string(m1[1])
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	t := doc.Find("table.construction").Eq(1)
	a, _ := t.Find("a.abortNow").Attr("onclick")
	r := regexp.MustCompile(`cancelResearch\((\d+),\s?(\d+),`)
	m := r.FindStringSubmatch(a)
	if len(m) < 3 {
		return "", 0, 0, errors.New("unable to find techid/listid")
	}
	techID = utils.DoParseI64(m[1])
	listID = utils.DoParseI64(m[2])
	return
}

// ExtractUniverseSpeed extract universe speed from html calculation
// pageHTML := b.getPageContent(url.Values{"page": {"techtree"}, "tab": {"2"}, "techID": {"1"}})
func ExtractUniverseSpeed(pageHTML []byte) int64 {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	spans := doc.Find("span.undermark")
	level := utils.ParseInt(spans.Eq(0).Text())
	val := utils.ParseInt(spans.Eq(1).Text())
	metalProduction := int64(math.Floor(30 * float64(level) * math.Pow(1.1, float64(level))))
	universeSpeed := val / metalProduction
	return universeSpeed
}

var temperatureRgxStr = `([-\d]+).+C\s*(?:bis|-tól|para|to|à|至|a|～|do|ile|tot|og|до|až|til|la|έως|:sta)\s*([-\d]+).+C`
var TemperatureRgx = regexp.MustCompile(temperatureRgxStr)
var diameterRgxStr = `([\d.,]+)(?i)(?:km|км|公里|χμ)`
var DiameterRgx = regexp.MustCompile(diameterRgxStr)
var lifeformRgxStr = `(?:[^:]+:\s\D+)?`
var planetInfosRgx = regexp.MustCompile(`([^\[]+) \[(\d+):(\d+):(\d+)]` + lifeformRgxStr + diameterRgxStr + ` \((\d+)/(\d+)\)(?:de|da|od|mellem|от)?\s*` + temperatureRgxStr)
var moonInfosRgx = regexp.MustCompile(`([^\[]+) \[(\d+):(\d+):(\d+)]([\d.,]+)(?i)(?:km|км|χμ|公里) \((\d+)/(\d+)\)`)
var cpRgx = regexp.MustCompile(`&cp=(\d+)`)

func extractPlanetFromSelection(s *goquery.Selection) (ogame.Planet, error) {
	el, _ := s.Attr("id")
	id, err := utils.ParseI64(strings.TrimPrefix(el, "planet-"))
	if err != nil {
		return ogame.Planet{}, err
	}

	title, _ := s.Find("a.planetlink").Attr("title")
	root, err := html.Parse(strings.NewReader(title))
	if err != nil {
		return ogame.Planet{}, err
	}

	txt := goquery.NewDocumentFromNode(root).Text()
	m := planetInfosRgx.FindStringSubmatch(txt)
	if len(m) < 10 {
		return ogame.Planet{}, errors.New("failed to parse planet infos: " + txt)
	}

	res := ogame.Planet{}
	res.Img = s.Find("img.planetPic").AttrOr("src", "")
	res.ID = ogame.PlanetID(id)
	res.Name = strings.TrimSpace(m[1])
	res.Coordinate.Galaxy = utils.DoParseI64(m[2])
	res.Coordinate.System = utils.DoParseI64(m[3])
	res.Coordinate.Position = utils.DoParseI64(m[4])
	res.Coordinate.Type = ogame.PlanetType
	res.Diameter = utils.ParseInt(m[5])
	res.Fields.Built = utils.DoParseI64(m[6])
	res.Fields.Total = utils.DoParseI64(m[7])
	res.Temperature.Min = utils.DoParseI64(m[8])
	res.Temperature.Max = utils.DoParseI64(m[9])

	res.Moon, _ = extractMoonFromPlanetSelection(s)

	return res, nil
}

func extractMoonFromSelection(moonLink *goquery.Selection) (ogame.Moon, error) {
	href, found := moonLink.Attr("href")
	if !found {
		return ogame.Moon{}, errors.New("no moon found")
	}
	m := cpRgx.FindStringSubmatch(href)
	id := utils.DoParseI64(m[1])
	title, _ := moonLink.Attr("title")
	root, err := html.Parse(strings.NewReader(title))
	if err != nil {
		return ogame.Moon{}, err
	}
	txt := goquery.NewDocumentFromNode(root).Text()
	mm := moonInfosRgx.FindStringSubmatch(txt)
	if len(mm) < 8 {
		return ogame.Moon{}, errors.New("failed to parse moon infos: " + txt)
	}
	moon := ogame.Moon{}
	moon.ID = ogame.MoonID(id)
	moon.Name = strings.TrimSpace(mm[1])
	moon.Coordinate.Galaxy = utils.DoParseI64(mm[2])
	moon.Coordinate.System = utils.DoParseI64(mm[3])
	moon.Coordinate.Position = utils.DoParseI64(mm[4])
	moon.Coordinate.Type = ogame.MoonType
	moon.Diameter = utils.ParseInt(mm[5])
	moon.Fields.Built = utils.DoParseI64(mm[6])
	moon.Fields.Total = utils.DoParseI64(mm[7])
	moon.Img = moonLink.Find("img.icon-moon").AttrOr("src", "")
	return moon, nil
}

func extractMoonFromPlanetSelection(s *goquery.Selection) (*ogame.Moon, error) {
	moonLink := s.Find("a.moonlink")
	moon, err := extractMoonFromSelection(moonLink)
	if err != nil {
		return nil, err
	}
	return &moon, nil
}

func extractEmpire(pageHTML []byte) ([]ogame.EmpireCelestial, error) {
	var out []ogame.EmpireCelestial
	raw, err := ExtractEmpireJSON(pageHTML)
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
		temperatureStr := utils.DoCastStr(planet["temperature"])
		m := TemperatureRgx.FindStringSubmatch(temperatureStr)
		if len(m) == 3 {
			tempMin = utils.DoParseI64(m[1])
			tempMax = utils.DoParseI64(m[2])
		}
		mm := DiameterRgx.FindStringSubmatch(utils.DoCastStr(planet["diameter"]))
		energyStr := utils.DoCastStr(planet["energy"])
		energyDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(energyStr))
		energy := utils.ParseInt(energyDoc.Find("div span").Text())
		celestialType := ogame.CelestialType(utils.DoCastF64(planet["type"]))
		out = append(out, ogame.EmpireCelestial{
			Name:     utils.DoCastStr(planet["name"]),
			ID:       ogame.CelestialID(utils.DoCastF64(planet["id"])),
			Diameter: utils.ParseInt(mm[1]),
			Img:      utils.DoCastStr(planet["image"]),
			Type:     celestialType,
			Fields: ogame.Fields{
				Built: int64(utils.DoCastF64(planet["fieldUsed"])),
				Total: int64(utils.DoCastF64(planet["fieldMax"])),
			},
			Temperature: ogame.Temperature{
				Min: tempMin,
				Max: tempMax,
			},
			Coordinate: ogame.Coordinate{
				Galaxy:   int64(utils.DoCastF64(planet["galaxy"])),
				System:   int64(utils.DoCastF64(planet["system"])),
				Position: int64(utils.DoCastF64(planet["position"])),
				Type:     celestialType,
			},
			Resources: ogame.Resources{
				Metal:     int64(utils.DoCastF64(planet["metal"])),
				Crystal:   int64(utils.DoCastF64(planet["crystal"])),
				Deuterium: int64(utils.DoCastF64(planet["deuterium"])),
				Energy:    energy,
			},
			Supplies: ogame.ResourcesBuildings{
				MetalMine:            int64(utils.DoCastF64(planet["1"])),
				CrystalMine:          int64(utils.DoCastF64(planet["2"])),
				DeuteriumSynthesizer: int64(utils.DoCastF64(planet["3"])),
				SolarPlant:           int64(utils.DoCastF64(planet["4"])),
				FusionReactor:        int64(utils.DoCastF64(planet["12"])),
				SolarSatellite:       int64(utils.DoCastF64(planet["212"])),
				MetalStorage:         int64(utils.DoCastF64(planet["22"])),
				CrystalStorage:       int64(utils.DoCastF64(planet["23"])),
				DeuteriumTank:        int64(utils.DoCastF64(planet["24"])),
			},
			Facilities: ogame.Facilities{
				RoboticsFactory: int64(utils.DoCastF64(planet["14"])),
				Shipyard:        int64(utils.DoCastF64(planet["21"])),
				ResearchLab:     int64(utils.DoCastF64(planet["31"])),
				AllianceDepot:   int64(utils.DoCastF64(planet["34"])),
				MissileSilo:     int64(utils.DoCastF64(planet["44"])),
				NaniteFactory:   int64(utils.DoCastF64(planet["15"])),
				Terraformer:     int64(utils.DoCastF64(planet["33"])),
				SpaceDock:       int64(utils.DoCastF64(planet["36"])),
				LunarBase:       int64(utils.DoCastF64(planet["41"])),
				SensorPhalanx:   int64(utils.DoCastF64(planet["42"])),
				JumpGate:        int64(utils.DoCastF64(planet["43"])),
			},
			Defenses: ogame.DefensesInfos{
				RocketLauncher:         int64(utils.DoCastF64(planet["401"])),
				LightLaser:             int64(utils.DoCastF64(planet["402"])),
				HeavyLaser:             int64(utils.DoCastF64(planet["403"])),
				GaussCannon:            int64(utils.DoCastF64(planet["404"])),
				IonCannon:              int64(utils.DoCastF64(planet["405"])),
				PlasmaTurret:           int64(utils.DoCastF64(planet["406"])),
				SmallShieldDome:        int64(utils.DoCastF64(planet["407"])),
				LargeShieldDome:        int64(utils.DoCastF64(planet["408"])),
				AntiBallisticMissiles:  int64(utils.DoCastF64(planet["502"])),
				InterplanetaryMissiles: int64(utils.DoCastF64(planet["503"])),
			},
			Researches: ogame.Researches{
				EnergyTechnology:             int64(utils.DoCastF64(planet["113"])),
				LaserTechnology:              int64(utils.DoCastF64(planet["120"])),
				IonTechnology:                int64(utils.DoCastF64(planet["121"])),
				HyperspaceTechnology:         int64(utils.DoCastF64(planet["114"])),
				PlasmaTechnology:             int64(utils.DoCastF64(planet["122"])),
				CombustionDrive:              int64(utils.DoCastF64(planet["115"])),
				ImpulseDrive:                 int64(utils.DoCastF64(planet["117"])),
				HyperspaceDrive:              int64(utils.DoCastF64(planet["118"])),
				EspionageTechnology:          int64(utils.DoCastF64(planet["106"])),
				ComputerTechnology:           int64(utils.DoCastF64(planet["108"])),
				Astrophysics:                 int64(utils.DoCastF64(planet["124"])),
				IntergalacticResearchNetwork: int64(utils.DoCastF64(planet["123"])),
				GravitonTechnology:           int64(utils.DoCastF64(planet["199"])),
				WeaponsTechnology:            int64(utils.DoCastF64(planet["109"])),
				ShieldingTechnology:          int64(utils.DoCastF64(planet["110"])),
				ArmourTechnology:             int64(utils.DoCastF64(planet["111"])),
			},
			Ships: ogame.ShipsInfos{
				LightFighter:   int64(utils.DoCastF64(planet["204"])),
				HeavyFighter:   int64(utils.DoCastF64(planet["205"])),
				Cruiser:        int64(utils.DoCastF64(planet["206"])),
				Battleship:     int64(utils.DoCastF64(planet["207"])),
				Battlecruiser:  int64(utils.DoCastF64(planet["215"])),
				Bomber:         int64(utils.DoCastF64(planet["211"])),
				Destroyer:      int64(utils.DoCastF64(planet["213"])),
				Deathstar:      int64(utils.DoCastF64(planet["214"])),
				SmallCargo:     int64(utils.DoCastF64(planet["202"])),
				LargeCargo:     int64(utils.DoCastF64(planet["203"])),
				ColonyShip:     int64(utils.DoCastF64(planet["208"])),
				Recycler:       int64(utils.DoCastF64(planet["209"])),
				EspionageProbe: int64(utils.DoCastF64(planet["210"])),
				SolarSatellite: int64(utils.DoCastF64(planet["212"])),
				Crawler:        int64(utils.DoCastF64(planet["217"])),
				Reaper:         int64(utils.DoCastF64(planet["218"])),
				Pathfinder:     int64(utils.DoCastF64(planet["219"])),
			},
		})
	}
	return out, nil
}

func ExtractEmpireJSON(pageHTML []byte) (any, error) {
	m := regexp.MustCompile(`createImperiumHtml\("#mainWrapper",\s"#loading",\s(.*),\s\d+\s\);`).FindSubmatch(pageHTML)
	if len(m) != 2 {
		return nil, errors.New("regexp for Empire JSON did not match anything")
	}
	var empireJSON any
	if err := json.Unmarshal(m[1], &empireJSON); err != nil {
		return nil, err
	}
	return empireJSON, nil
}

// ExtractAuction extract auction information from page "traderAuctioneer"
func extractAuctionFromDoc(doc *goquery.Document) (ogame.Auction, error) {
	auction := ogame.Auction{}
	auction.HasFinished = false

	// Detect if Auction has already finished
	nextAuction := doc.Find("#nextAuction")
	if nextAuction.Size() > 0 {
		// Find time until next auction starts
		auction.Endtime = utils.DoParseI64(nextAuction.Text())
		auction.HasFinished = true
	} else {
		endAtApprox := doc.Find("p.auction_info b").Text()
		m := regexp.MustCompile(`[^\d]+(\d+).*`).FindStringSubmatch(endAtApprox)
		if len(m) != 2 {
			return ogame.Auction{}, errors.New("failed to find end time approx")
		}
		endTimeMinutes, err := utils.ParseI64(m[1])
		if err != nil {
			return ogame.Auction{}, errors.New("invalid end time approx: " + err.Error())
		}
		auction.Endtime = endTimeMinutes * 60
	}

	auction.HighestBidder = strings.TrimSpace(doc.Find("a.currentPlayer").Text())
	auction.HighestBidderUserID = utils.DoParseI64(doc.Find("a.currentPlayer").AttrOr("data-player-id", ""))
	auction.NumBids = utils.DoParseI64(doc.Find("div.numberOfBids").Text())
	auction.CurrentBid = utils.ParseInt(doc.Find("div.currentSum").Text())
	auction.Inventory = utils.DoParseI64(doc.Find("span.level.amount").Text())
	auction.CurrentItem = strings.ToLower(doc.Find("img").First().AttrOr("alt", ""))
	auction.CurrentItemLong = strings.ToLower(doc.Find("div.image_140px").First().Find("a").First().AttrOr("title", ""))
	multiplierRegex := regexp.MustCompile(`multiplier\s?=\s?([^;]+);`).FindStringSubmatch(doc.Text())
	if len(multiplierRegex) != 2 {
		return ogame.Auction{}, errors.New("failed to find auction multiplier")
	}
	if err := json.Unmarshal([]byte(multiplierRegex[1]), &auction.ResourceMultiplier); err != nil {
		return ogame.Auction{}, errors.New("failed to json parse auction multiplier: " + err.Error())
	}

	// Find auctioneer token
	tokenRegex := regexp.MustCompile(`auctioneerToken\s?=\s?"([^"]+)";`).FindStringSubmatch(doc.Text())
	if len(tokenRegex) != 2 {
		return ogame.Auction{}, errors.New("failed to find auctioneer token")
	}
	auction.Token = tokenRegex[1]

	// Find Planet / Moon resources JSON
	planetMoonResources := regexp.MustCompile(`planetResources\s?=\s?([^;]+);`).FindStringSubmatch(doc.Text())
	if len(planetMoonResources) != 2 {
		return ogame.Auction{}, errors.New("failed to find planetResources")
	}
	if err := json.Unmarshal([]byte(planetMoonResources[1]), &auction.Resources); err != nil {
		return ogame.Auction{}, errors.New("failed to json unmarshal planetResources: " + err.Error())
	}

	// Find already-bid
	m := regexp.MustCompile(`var playerBid\s?=\s?([^;]+);`).FindStringSubmatch(doc.Text())
	if len(m) != 2 {
		return ogame.Auction{}, errors.New("failed to get playerBid")
	}
	var alreadyBid int64
	if m[1] != "false" {
		alreadyBid = utils.DoParseI64(m[1])
	}
	auction.AlreadyBid = alreadyBid

	// Find min-bid
	auction.MinimumBid = utils.ParseInt(doc.Find("table.table_ressources_sum tr td.auctionInfo.js_price").Text())

	// Find deficit-bid
	auction.DeficitBid = utils.ParseInt(doc.Find("table.table_ressources_sum tr td.auctionInfo.js_deficit").Text())

	// Note: Don't just bid the min-bid amount. It will keep doubling the total bid and grow exponentially...
	// DeficitBid is 1000 when another player has outbid you or if nobody has bid yet.
	// DeficitBid seems to be filled by Javascript in the browser. We're parsing it anyway. Correct Bid calculation would be:
	// bid = max(auction.DeficitBid, auction.MinimumBid - auction.AlreadyBid)

	return auction, nil
}
