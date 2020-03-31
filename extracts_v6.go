package ogame

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

	"github.com/PuerkitoBio/goquery"
	"github.com/alaingilbert/clockwork"
	lua "github.com/yuin/gopher-lua"
	"golang.org/x/net/html"
)

func extractIsInVacationFromDocV6(doc *goquery.Document) bool {
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

func extractPlanetsFromDocV6(doc *goquery.Document, b *OGame) []Planet {
	res := make([]Planet, 0)
	doc.Find("div.smallplanet").Each(func(i int, s *goquery.Selection) {
		planet, err := extractPlanetFromSelectionV6(s, b)
		if err != nil {
			b.error(err)
			return
		}
		res = append(res, planet)
	})
	return res
}

func extractPlanetByIDFromDocV6(doc *goquery.Document, b *OGame, planetID PlanetID) (Planet, error) {
	planets := extractPlanetsFromDocV6(doc, b)
	for _, planet := range planets {
		if planet.ID == planetID {
			return planet, nil
		}
	}
	return Planet{}, errors.New("invalid planet id")
}

func extractResourcesFromDocV6(doc *goquery.Document) Resources {
	res := Resources{}
	metalDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(doc.Find("li#metal_box").AttrOr("title", "")))
	crystalDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(doc.Find("li#crystal_box").AttrOr("title", "")))
	deuteriumDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(doc.Find("li#deuterium_box").AttrOr("title", "")))
	energyDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(doc.Find("li#energy_box").AttrOr("title", "")))
	darkmatterDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(doc.Find("li#darkmatter_box").AttrOr("title", "")))
	res.Metal = ParseInt(metalDoc.Find("table tr").Eq(0).Find("td").Eq(0).Text())
	res.Crystal = ParseInt(crystalDoc.Find("table tr").Eq(0).Find("td").Eq(0).Text())
	res.Deuterium = ParseInt(deuteriumDoc.Find("table tr").Eq(0).Find("td").Eq(0).Text())
	res.Energy = ParseInt(energyDoc.Find("table tr").Eq(0).Find("td").Eq(0).Text())
	res.Darkmatter = ParseInt(darkmatterDoc.Find("table tr").Eq(0).Find("td").Eq(0).Text())
	return res
}

func extractResourcesDetailsFromFullPageFromDocV6(doc *goquery.Document) ResourcesDetails {
	out := ResourcesDetails{}
	out.Metal.Available = ParseInt(doc.Find("span#resources_metal").Text())
	out.Crystal.Available = ParseInt(doc.Find("span#resources_crystal").Text())
	out.Deuterium.Available = ParseInt(doc.Find("span#resources_deuterium").Text())
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

func extractHiddenFieldsFromDocV6(doc *goquery.Document) url.Values {
	fields := url.Values{}
	doc.Find("input[type=hidden]").Each(func(i int, s *goquery.Selection) {
		name, _ := s.Attr("name")
		value, _ := s.Attr("value")
		fields.Add(name, value)
	})
	return fields
}

func extractBodyIDFromDocV6(doc *goquery.Document) string {
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

func extractCelestialByIDFromDocV6(doc *goquery.Document, b *OGame, celestialID CelestialID) (Celestial, error) {
	planets := extractPlanetsFromDocV6(doc, b)
	for _, planet := range planets {
		if planet.ID.Celestial() == celestialID {
			return planet, nil
		}
		if planet.Moon != nil && planet.Moon.ID.Celestial() == celestialID {
			return planet.Moon, nil
		}
	}
	return Planet{}, errors.New("invalid celestial id")
}

func extractPlanetByCoordFromDocV6(doc *goquery.Document, b *OGame, coord Coordinate) (Planet, error) {
	planets := extractPlanetsFromDocV6(doc, b)
	for _, planet := range planets {
		if planet.Coordinate.Equal(coord) {
			return planet, nil
		}
	}
	return Planet{}, errors.New("invalid planet coordinate")
}

func extractOgameTimestampFromDocV6(doc *goquery.Document) int64 {
	ogameTimestamp, _ := strconv.ParseInt(doc.Find("meta[name=ogame-timestamp]").AttrOr("content", "0"), 10, 64)
	return ogameTimestamp
}

func extractPlanetFromDocV6(doc *goquery.Document, v interface{}, b *OGame) (Planet, error) {
	if coordStr, ok := v.(string); ok {
		coord, err := ParseCoord(coordStr)
		if err != nil {
			return Planet{}, err
		}
		return extractPlanetByCoordFromDocV6(doc, b, coord)
	} else if coord, ok := v.(Coordinate); ok {
		return extractPlanetByCoordFromDocV6(doc, b, coord)
	} else if planetID, ok := v.(PlanetID); ok {
		return extractPlanetByIDFromDocV6(doc, b, planetID)
	} else if id, ok := v.(int); ok {
		return extractPlanetByIDFromDocV6(doc, b, PlanetID(id))
	} else if id, ok := v.(int32); ok {
		return extractPlanetByIDFromDocV6(doc, b, PlanetID(id))
	} else if id, ok := v.(int64); ok {
		return extractPlanetByIDFromDocV6(doc, b, PlanetID(id))
	} else if id, ok := v.(float32); ok {
		return extractPlanetByIDFromDocV6(doc, b, PlanetID(id))
	} else if id, ok := v.(float64); ok {
		return extractPlanetByIDFromDocV6(doc, b, PlanetID(id))
	} else if id, ok := v.(lua.LNumber); ok {
		return extractPlanetByIDFromDocV6(doc, b, PlanetID(id))
	}
	return Planet{}, errors.New("failed to find planet")
}

func extractMoonsFromDocV6(doc *goquery.Document, b *OGame) []Moon {
	res := make([]Moon, 0)
	doc.Find("a.moonlink").Each(func(i int, s *goquery.Selection) {
		moon, err := extractMoonFromSelectionV6(s, b)
		if err != nil {
			return
		}
		res = append(res, moon)
	})
	return res
}

func extractMoonFromDocV6(doc *goquery.Document, b *OGame, v interface{}) (Moon, error) {
	if coordStr, ok := v.(string); ok {
		coord, err := ParseCoord(coordStr)
		if err != nil {
			return Moon{}, err
		}
		return extractMoonByCoordFromDocV6(doc, b, coord)
	} else if coord, ok := v.(Coordinate); ok {
		return extractMoonByCoordFromDocV6(doc, b, coord)
	} else if moonID, ok := v.(MoonID); ok {
		return extractMoonByIDFromDocV6(doc, b, moonID)
	} else if id, ok := v.(int); ok {
		return extractMoonByIDFromDocV6(doc, b, MoonID(id))
	} else if id, ok := v.(int32); ok {
		return extractMoonByIDFromDocV6(doc, b, MoonID(id))
	} else if id, ok := v.(int64); ok {
		return extractMoonByIDFromDocV6(doc, b, MoonID(id))
	} else if id, ok := v.(float32); ok {
		return extractMoonByIDFromDocV6(doc, b, MoonID(id))
	} else if id, ok := v.(float64); ok {
		return extractMoonByIDFromDocV6(doc, b, MoonID(id))
	} else if id, ok := v.(lua.LNumber); ok {
		return extractMoonByIDFromDocV6(doc, b, MoonID(id))
	}
	return Moon{}, errors.New("moon not found")
}

func extractMoonByCoordFromDocV6(doc *goquery.Document, b *OGame, coord Coordinate) (Moon, error) {
	moons := extractMoonsFromDocV6(doc, b)
	for _, moon := range moons {
		if moon.Coordinate.Equal(coord) {
			return moon, nil
		}
	}
	return Moon{}, errors.New("invalid moon coordinate")
}

func extractMoonByIDFromDocV6(doc *goquery.Document, b *OGame, moonID MoonID) (Moon, error) {
	moons := extractMoonsFromDocV6(doc, b)
	for _, moon := range moons {
		if moon.ID == moonID {
			return moon, nil
		}
	}
	return Moon{}, errors.New("invalid moon id")
}

func extractCelestialsFromDocV6(doc *goquery.Document, b *OGame) ([]Celestial, error) {
	celestials := make([]Celestial, 0)
	planets := extractPlanetsFromDocV6(doc, b)
	for _, planet := range planets {
		celestials = append(celestials, planet)
		if planet.Moon != nil {
			celestials = append(celestials, planet.Moon)
		}
	}
	return celestials, nil
}

func extractCelestialFromDocV6(doc *goquery.Document, b *OGame, v interface{}) (Celestial, error) {
	if planetID, ok := v.(PlanetID); ok {
		return extractPlanetByIDFromDocV6(doc, b, planetID)
	} else if moonID, ok := v.(MoonID); ok {
		return extractMoonByIDFromDocV6(doc, b, moonID)
	} else if celestialID, ok := v.(CelestialID); ok {
		return extractCelestialByIDFromDocV6(doc, b, celestialID)
	} else if id, ok := v.(int); ok {
		return extractCelestialByIDFromDocV6(doc, b, CelestialID(id))
	} else if id, ok := v.(int32); ok {
		return extractCelestialByIDFromDocV6(doc, b, CelestialID(id))
	} else if id, ok := v.(int64); ok {
		return extractCelestialByIDFromDocV6(doc, b, CelestialID(id))
	} else if id, ok := v.(float32); ok {
		return extractCelestialByIDFromDocV6(doc, b, CelestialID(id))
	} else if id, ok := v.(float64); ok {
		return extractCelestialByIDFromDocV6(doc, b, CelestialID(id))
	} else if id, ok := v.(lua.LNumber); ok {
		return extractCelestialByIDFromDocV6(doc, b, CelestialID(id))
	} else if coord, ok := v.(Coordinate); ok {
		if coord.Type == PlanetType {
			return extractPlanetByCoordFromDocV6(doc, b, coord)
		} else if coord.Type == MoonType {
			return extractMoonByCoordFromDocV6(doc, b, coord)
		}
	} else if coordStr, ok := v.(string); ok {
		coord, err := ParseCoord(coordStr)
		if err != nil {
			return nil, err
		}
		if coord.Type == PlanetType {
			return extractPlanetByCoordFromDocV6(doc, b, coord)
		} else if coord.Type == MoonType {
			return extractMoonByCoordFromDocV6(doc, b, coord)
		}
	}
	return nil, errors.New("celestial not found")
}

func extractResourcesBuildingsFromDocV6(doc *goquery.Document) (ResourcesBuildings, error) {
	doc.Find("span.textlabel").Remove()
	bodyID := extractBodyIDFromDocV6(doc)
	if bodyID == "overview" {
		return ResourcesBuildings{}, ErrInvalidPlanetID
	}
	res := ResourcesBuildings{}
	res.MetalMine = getNbr(doc, "supply1")
	res.CrystalMine = getNbr(doc, "supply2")
	res.DeuteriumSynthesizer = getNbr(doc, "supply3")
	res.SolarPlant = getNbr(doc, "supply4")
	res.FusionReactor = getNbr(doc, "supply12")
	res.SolarSatellite = getNbr(doc, "supply212")
	res.MetalStorage = getNbr(doc, "supply22")
	res.CrystalStorage = getNbr(doc, "supply23")
	res.DeuteriumTank = getNbr(doc, "supply24")
	return res, nil
}

func extractDefenseFromDocV6(doc *goquery.Document) (DefensesInfos, error) {
	bodyID := extractBodyIDFromDocV6(doc)
	if bodyID == "overview" {
		return DefensesInfos{}, ErrInvalidPlanetID
	}
	doc.Find("span.textlabel").Remove()
	res := DefensesInfos{}
	res.RocketLauncher = getNbr(doc, "defense401")
	res.LightLaser = getNbr(doc, "defense402")
	res.HeavyLaser = getNbr(doc, "defense403")
	res.GaussCannon = getNbr(doc, "defense404")
	res.IonCannon = getNbr(doc, "defense405")
	res.PlasmaTurret = getNbr(doc, "defense406")
	res.SmallShieldDome = getNbr(doc, "defense407")
	res.LargeShieldDome = getNbr(doc, "defense408")
	res.AntiBallisticMissiles = getNbr(doc, "defense502")
	res.InterplanetaryMissiles = getNbr(doc, "defense503")
	return res, nil
}

func extractShipsFromDocV6(doc *goquery.Document) (ShipsInfos, error) {
	doc.Find("span.textlabel").Remove()
	bodyID := extractBodyIDFromDocV6(doc)
	if bodyID == "overview" {
		return ShipsInfos{}, ErrInvalidPlanetID
	}
	res := ShipsInfos{}
	res.LightFighter = getNbrShips(doc, "military204")
	res.HeavyFighter = getNbrShips(doc, "military205")
	res.Cruiser = getNbrShips(doc, "military206")
	res.Battleship = getNbrShips(doc, "military207")
	res.Battlecruiser = getNbrShips(doc, "military215")
	res.Bomber = getNbrShips(doc, "military211")
	res.Destroyer = getNbrShips(doc, "military213")
	res.Deathstar = getNbrShips(doc, "military214")
	res.SmallCargo = getNbrShips(doc, "civil202")
	res.LargeCargo = getNbrShips(doc, "civil203")
	res.ColonyShip = getNbrShips(doc, "civil208")
	res.Recycler = getNbrShips(doc, "civil209")
	res.EspionageProbe = getNbrShips(doc, "civil210")
	res.SolarSatellite = getNbrShips(doc, "civil212")

	return res, nil
}

func extractFacilitiesFromDocV6(doc *goquery.Document) (Facilities, error) {
	doc.Find("span.textlabel").Remove()
	bodyID := extractBodyIDFromDocV6(doc)
	if bodyID == "overview" {
		return Facilities{}, ErrInvalidPlanetID
	}
	res := Facilities{}
	res.RoboticsFactory = getNbr(doc, "station14")
	res.Shipyard = getNbr(doc, "station21")
	res.ResearchLab = getNbr(doc, "station31")
	res.AllianceDepot = getNbr(doc, "station34")
	res.MissileSilo = getNbr(doc, "station44")
	res.NaniteFactory = getNbr(doc, "station15")
	res.Terraformer = getNbr(doc, "station33")
	res.SpaceDock = getNbr(doc, "station36")
	res.LunarBase = getNbr(doc, "station41")
	res.SensorPhalanx = getNbr(doc, "station42")
	res.JumpGate = getNbr(doc, "station43")
	return res, nil
}

func extractResearchFromDocV6(doc *goquery.Document) Researches {
	doc.Find("span.textlabel").Remove()
	res := Researches{}
	res.EnergyTechnology = getNbr(doc, "research113")
	res.LaserTechnology = getNbr(doc, "research120")
	res.IonTechnology = getNbr(doc, "research121")
	res.HyperspaceTechnology = getNbr(doc, "research114")
	res.PlasmaTechnology = getNbr(doc, "research122")
	res.CombustionDrive = getNbr(doc, "research115")
	res.ImpulseDrive = getNbr(doc, "research117")
	res.HyperspaceDrive = getNbr(doc, "research118")
	res.EspionageTechnology = getNbr(doc, "research106")
	res.ComputerTechnology = getNbr(doc, "research108")
	res.Astrophysics = getNbr(doc, "research124")
	res.IntergalacticResearchNetwork = getNbr(doc, "research123")
	res.GravitonTechnology = getNbr(doc, "research199")
	res.WeaponsTechnology = getNbr(doc, "research109")
	res.ShieldingTechnology = getNbr(doc, "research110")
	res.ArmourTechnology = getNbr(doc, "research111")
	return res
}

func extractOGameSessionFromDocV6(doc *goquery.Document) string {
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

func extractAttacksFromDocV6(doc *goquery.Document, clock clockwork.Clock) ([]AttackEvent, error) {
	attacks := make([]*AttackEvent, 0)
	out := make([]AttackEvent, 0)
	if doc.Find("body").Size() == 1 && extractOGameSessionFromDocV6(doc) != "" && doc.Find("div#eventListWrap").Size() == 0 {
		return out, ErrEventsBoxNotDisplayed
	} else if doc.Find("div#eventListWrap").Size() == 0 {
		return out, ErrNotLogged
	}

	allianceAttacks := make(map[int64]*AttackEvent)

	tmp := func(i int, s *goquery.Selection) {
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
		if missionType != Attack && missionType != GroupedAttack && missionType != Destroy &&
			missionType != MissileAttack && missionType != Spy {
			return
		}
		attack := &AttackEvent{}
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
	doc.Find("tr.allianceAttack").Each(tmp)
	doc.Find("tr.eventFleet").Each(tmp)

	for _, a := range attacks {
		out = append(out, *a)
	}

	return out, nil
}

func extractOfferOfTheDayFromDocV6(doc *goquery.Document) (price int64, importToken string, planetResources PlanetResources, multiplier Multiplier, err error) {
	s := doc.Find("div.js_import_price")
	if s.Size() == 0 {
		err = errors.New("failed to extract offer of the day price")
		return
	}
	price = ParseInt(s.Text())
	script := doc.Find("script").Text()
	m := regexp.MustCompile(`var importToken="([^"]*)";`).FindSubmatch([]byte(script))
	if len(m) != 2 {
		err = errors.New("failed to extract offer of the day import token")
		return
	}
	importToken = string(m[1])
	m = regexp.MustCompile(`var planetResources=({[^;]*});`).FindSubmatch([]byte(script))
	if len(m) != 2 {
		err = errors.New("failed to extract offer of the day raw planet resources")
		return
	}
	if err = json.Unmarshal(m[1], &planetResources); err != nil {
		return
	}
	m = regexp.MustCompile(`var multiplier=({[^;]*});`).FindSubmatch([]byte(script))
	if len(m) != 2 {
		err = errors.New("failed to extract offer of the day raw multiplier")
		return
	}
	if err = json.Unmarshal(m[1], &multiplier); err != nil {
		return
	}
	return
}

func extractProductionFromDocV6(doc *goquery.Document) ([]Quantifiable, error) {
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
					itemIDstr = strconv.FormatInt(int64(AntiBallisticMissilesID), 10)
				} else if strings.HasSuffix(src, "36221e9493458b9fcc776bf350983e.jpg") {
					itemIDstr = strconv.FormatInt(int64(InterplanetaryMissilesID), 10)
				}
			}
		}
		itemID, _ := strconv.ParseInt(itemIDstr, 10, 64)
		itemNbr := ParseInt(s.Find("span.number").Text())
		res = append(res, Quantifiable{ID: ID(itemID), Nbr: itemNbr})
	})
	return res, nil
}

func extractOverviewProductionFromDocV6(doc *goquery.Document) ([]Quantifiable, error) {
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
		link := s.Find("a")
		href := link.AttrOr("href", "")
		m := regexp.MustCompile(`openTech=(\d+)`).FindStringSubmatch(href)
		if len(m) == 0 {
			return
		}
		idInt, _ := strconv.ParseInt(m[1], 10, 64)
		activeID := ID(idInt)
		activeNbr := ParseInt(link.Text())
		res = append(res, Quantifiable{ID: activeID, Nbr: activeNbr})
	})
	return res, nil
}

func extractFleet1ShipsFromDocV6(doc *goquery.Document) (s ShipsInfos) {
	onclick := doc.Find("a#sendall").AttrOr("onclick", "")
	matches := regexp.MustCompile(`setMaxIntInput\("form\[name=shipsChosen]", (.+)\); checkShips`).FindStringSubmatch(onclick)
	if len(matches) == 0 {
		return
	}
	m := matches[1]
	var res map[ID]int64
	if err := json.Unmarshal([]byte(m), &res); err != nil {
		return
	}
	for k, v := range res {
		s.Set(k, v)
	}
	return
}

func extractEspionageReportMessageIDsFromDocV6(doc *goquery.Document) ([]EspionageReportSummary, int64) {
	msgs := make([]EspionageReportSummary, 0)
	nbPage, _ := strconv.ParseInt(doc.Find("ul.pagination li").Last().AttrOr("data-page", "1"), 10, 64)
	doc.Find("li.msg").Each(func(i int, s *goquery.Selection) {
		if idStr, exists := s.Attr("data-msg-id"); exists {
			if id, err := strconv.ParseInt(idStr, 10, 64); err == nil {
				messageType := Report
				if s.Find("span.espionageDefText").Size() > 0 {
					messageType = Action
				}
				report := EspionageReportSummary{ID: id, Type: messageType}
				report.From = s.Find("span.msg_sender").Text()
				spanLink := s.Find("span.msg_title a")
				targetStr := spanLink.Text()
				report.Target = extractCoordV6(targetStr)
				report.Target.Type = PlanetType
				if spanLink.Find("figure").HasClass("moon") {
					report.Target.Type = MoonType
				}
				if messageType == Report {
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

func extractCombatReportMessagesFromDocV6(doc *goquery.Document) ([]CombatReportSummary, int64) {
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
				resTitle := s.Find("span.msg_content div.combatLeftSide span").Eq(1).AttrOr("title", "")
				m := regexp.MustCompile(`([\d.]+)<br/>[^\d]*([\d.]+)<br/>[^\d]*([\d.]+)`).FindStringSubmatch(resTitle)
				if len(m) == 4 {
					report.Metal = ParseInt(m[1])
					report.Crystal = ParseInt(m[2])
					report.Deuterium = ParseInt(m[3])
				}
				debrisFieldTitle := s.Find("span.msg_content div.combatLeftSide span").Eq(2).AttrOr("title", "0")
				report.DebrisField = ParseInt(debrisFieldTitle)
				resText := s.Find("span.msg_content div.combatLeftSide span").Eq(1).Text()
				m = regexp.MustCompile(`[\d.]+[^\d]*([\d.]+)`).FindStringSubmatch(resText)
				if len(m) == 2 {
					report.Loot = ParseInt(m[1])
				}
				msgDate, _ := time.Parse("02.01.2006 15:04:05", s.Find("span.msg_date").Text())
				report.CreatedAt = msgDate

				link := s.Find("div.msg_actions a span.icon_attack").Parent().AttrOr("href", "")
				m = regexp.MustCompile(`page=fleet1&galaxy=(\d+)&system=(\d+)&position=(\d+)&type=(\d+)&`).FindStringSubmatch(link)
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

func extractEspionageReportFromDocV6(doc *goquery.Document, location *time.Location) (EspionageReport, error) {
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
	activity := doc.Find("div.detail_txt").Eq(1).Find("font")
	if len(activity.Text()) == 2 {
		report.LastActivity = ParseInt(activity.Text())
	}

	// CounterEspionage
	ceTxt := doc.Find("div.detail_txt").Eq(1).Text()
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

func extractResourcesProductionsFromDocV6(doc *goquery.Document) (Resources, error) {
	res := Resources{}
	selector := "table.listOfResourceSettingsPerPlanet tr.summary td span"
	el := doc.Find(selector)
	res.Metal = ParseInt(el.Eq(0).AttrOr("title", "0"))
	res.Crystal = ParseInt(el.Eq(1).AttrOr("title", "0"))
	res.Deuterium = ParseInt(el.Eq(2).AttrOr("title", "0"))
	res.Energy = ParseInt(el.Eq(3).AttrOr("title", "0"))
	return res, nil
}

func extractPreferencesFromDocV6(doc *goquery.Document) Preferences {
	prefs := Preferences{
		SpioAnz:                      extractSpioAnzFromDocV6(doc),
		DisableChatBar:               extractDisableChatBarFromDocV6(doc),
		DisableOutlawWarning:         extractDisableOutlawWarningFromDocV6(doc),
		MobileVersion:                extractMobileVersionFromDocV6(doc),
		ShowOldDropDowns:             extractShowOldDropDownsFromDocV6(doc),
		ActivateAutofocus:            extractActivateAutofocusFromDocV6(doc),
		EventsShow:                   extractEventsShowFromDocV6(doc),
		SortSetting:                  extractSortSettingFromDocV6(doc),
		SortOrder:                    extractSortOrderFromDocV6(doc),
		ShowDetailOverlay:            extractShowDetailOverlayFromDocV6(doc),
		AnimatedSliders:              extractAnimatedSlidersFromDocV6(doc),
		AnimatedOverview:             extractAnimatedOverviewFromDocV6(doc),
		PopupsNotices:                extractPopupsNoticesFromDocV6(doc),
		PopopsCombatreport:           extractPopopsCombatreportFromDocV6(doc),
		SpioReportPictures:           extractSpioReportPicturesFromDocV6(doc),
		MsgResultsPerPage:            extractMsgResultsPerPageFromDocV6(doc),
		AuctioneerNotifications:      extractAuctioneerNotificationsFromDocV6(doc),
		EconomyNotifications:         extractEconomyNotificationsFromDocV6(doc),
		ShowActivityMinutes:          extractShowActivityMinutesFromDocV6(doc),
		PreserveSystemOnPlanetChange: extractPreserveSystemOnPlanetChangeFromDocV6(doc),
	}
	if prefs.MobileVersion {
		prefs.Notifications.BuildList = extractNotifBuildListFromDocV6(doc)
		prefs.Notifications.FriendlyFleetActivities = extractNotifFriendlyFleetActivitiesFromDocV6(doc)
		prefs.Notifications.HostileFleetActivities = extractNotifHostileFleetActivitiesFromDocV6(doc)
		prefs.Notifications.ForeignEspionage = extractNotifForeignEspionageFromDocV6(doc)
		prefs.Notifications.AllianceBroadcasts = extractNotifAllianceBroadcastsFromDocV6(doc)
		prefs.Notifications.AllianceMessages = extractNotifAllianceMessagesFromDocV6(doc)
		prefs.Notifications.Auctions = extractNotifAuctionsFromDocV6(doc)
		prefs.Notifications.Account = extractNotifAccountFromDocV6(doc)
	}
	return prefs
}

func extractResourceSettingsFromDocV6(doc *goquery.Document) (ResourceSettings, error) {
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
	if len(vals) != 6 {
		return ResourceSettings{}, errors.New("failed to find all resource settings")
	}

	res := ResourceSettings{}
	res.MetalMine = vals[0]
	res.CrystalMine = vals[1]
	res.DeuteriumSynthesizer = vals[2]
	res.SolarPlant = vals[3]
	res.FusionReactor = vals[4]
	res.SolarSatellite = vals[5]

	return res, nil
}

func extractFleetsFromEventListFromDocV6(doc *goquery.Document) []Fleet {
	type Tmp struct {
		fleet Fleet
		res   Resources
	}
	tmp := make([]Tmp, 0)
	res := make([]Fleet, 0)
	doc.Find("tr.eventFleet").Each(func(i int, s *goquery.Selection) {
		fleet := Fleet{}

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
			nbr := ParseInt(s.Find("td").Eq(1).Text())
			if name != "" && nbr > 0 {
				fleet.Ships.Set(name2id(name), nbr)
			}
		})
		fleet.Origin = extractCoordV6(doc.Find("td.coordsOrigin").Text())
		fleet.Destination = extractCoordV6(doc.Find("td.destCoords").Text())

		res := Resources{}
		trs := doc2.Find("tr")
		res.Metal = ParseInt(trs.Eq(trs.Size() - 3).Find("td").Eq(1).Text())
		res.Crystal = ParseInt(trs.Eq(trs.Size() - 2).Find("td").Eq(1).Text())
		res.Deuterium = ParseInt(trs.Eq(trs.Size() - 1).Find("td").Eq(1).Text())

		tmp = append(tmp, Tmp{fleet: fleet, res: res})
	})

	for _, t := range tmp {
		res = append(res, t.fleet)
	}

	return res
}

func extractIPMFromDocV6(doc *goquery.Document) (duration, max int64, token string) {
	duration, _ = strconv.ParseInt(doc.Find("span#timer").AttrOr("data-duration", "0"), 10, 64)
	max, _ = strconv.ParseInt(doc.Find("input[name=anz]").AttrOr("data-max", "0"), 10, 64)
	token = doc.Find("input[name=token]").AttrOr("value", "")
	return
}

func extractFleetsFromDocV6(doc *goquery.Document, clock clockwork.Clock) (res []Fleet) {
	res = make([]Fleet, 0)
	script := doc.Find("body script").Text()
	doc.Find("div.fleetDetails").Each(func(i int, s *goquery.Selection) {
		originText := s.Find("span.originCoords a").Text()
		origin := extractCoordV6(originText)
		origin.Type = PlanetType
		if s.Find("span.originPlanet figure").HasClass("moon") {
			origin.Type = MoonType
		}

		destText := s.Find("span.destinationCoords a").Text()
		dest := extractCoordV6(destText)
		dest.Type = PlanetType
		if s.Find("span.destinationPlanet figure").HasClass("moon") {
			dest.Type = MoonType
		} else if s.Find("span.destinationPlanet figure").HasClass("tf") {
			dest.Type = DebrisType
		}

		id, _ := strconv.ParseInt(s.Find("a.openCloseDetails").AttrOr("data-mission-id", "0"), 10, 64)

		timerID := s.Find("span.timer").AttrOr("id", "")
		m := regexp.MustCompile(`getElementByIdWithCache\("` + timerID + `"\),\s*(\d+),`).FindStringSubmatch(script)
		var arriveIn int64
		if len(m) == 2 {
			arriveIn, _ = strconv.ParseInt(m[1], 10, 64)
		}

		timerNextID := s.Find("span.nextTimer").AttrOr("id", "")
		m = regexp.MustCompile(`getElementByIdWithCache\("` + timerNextID + `"\),\s*(\d+)\s*\);`).FindStringSubmatch(script)
		var backIn int64
		if len(m) == 2 {
			backIn, _ = strconv.ParseInt(m[1], 10, 64)
		}

		missionType, _ := strconv.ParseInt(s.AttrOr("data-mission-type", ""), 10, 64)
		returnFlight, _ := strconv.ParseBool(s.AttrOr("data-return-flight", ""))
		arrivalTime, _ := strconv.ParseInt(s.AttrOr("data-arrival-time", ""), 10, 64)
		endTime, _ := strconv.ParseInt(s.Find("a.openCloseDetails").AttrOr("data-end-time", ""), 10, 64)
		ogameTimestamp, _ := strconv.ParseInt(doc.Find("meta[name=ogame-timestamp]").AttrOr("content", "0"), 10, 64)
		secs := arrivalTime - ogameTimestamp
		if secs < 0 {
			secs = 0
		}

		trs := s.Find("table.fleetinfo tr")
		shipment := Resources{}
		shipment.Metal = ParseInt(trs.Eq(trs.Size() - 3).Find("td").Eq(1).Text())
		shipment.Crystal = ParseInt(trs.Eq(trs.Size() - 2).Find("td").Eq(1).Text())
		shipment.Deuterium = ParseInt(trs.Eq(trs.Size() - 1).Find("td").Eq(1).Text())

		fedAttackHref := s.Find("span.fedAttack a").AttrOr("href", "")
		fedAttackURL, _ := url.Parse(fedAttackHref)
		fedAttackQuery := fedAttackURL.Query()
		targetPlanetID, _ := strconv.ParseInt(fedAttackQuery.Get("target"), 10, 64)
		unionID, _ := strconv.ParseInt(fedAttackQuery.Get("union"), 10, 64)

		fleet := Fleet{}
		fleet.ID = FleetID(id)
		fleet.Origin = origin
		fleet.Destination = dest
		fleet.Mission = MissionID(missionType)
		fleet.ReturnFlight = returnFlight
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
				startTime, _ = time.Parse("02.01.2006<br>15:04:05", startTimeArray[1])
			}
		}
		fleet.StartTime = startTime

		for i := 1; i < trs.Size()-5; i++ {
			tds := trs.Eq(i).Find("td")
			name := strings.ToLower(strings.Trim(strings.TrimSpace(tds.Eq(0).Text()), ":"))
			qty := ParseInt(tds.Eq(1).Text())
			shipID := name2id(name)
			fleet.Ships.Set(shipID, qty)
		}

		res = append(res, fleet)
	})
	return
}

func extractSlotsFromDocV6(doc *goquery.Document) Slots {
	slots := Slots{}
	page := extractBodyIDFromDocV6(doc)
	if page == "movement" {
		slots.InUse = ParseInt(doc.Find("span.fleetSlots > span.current").Text())
		slots.Total = ParseInt(doc.Find("span.fleetSlots > span.all").Text())
		slots.ExpInUse = ParseInt(doc.Find("span.expSlots > span.current").Text())
		slots.ExpTotal = ParseInt(doc.Find("span.expSlots > span.all").Text())
	} else if page == "fleet1" || page == "fleetdispatch" {
		r := regexp.MustCompile(`(\d+)/(\d+)`)
		txt := doc.Find("div#slots>div").Eq(0).Text()
		m := r.FindStringSubmatch(txt)
		if len(m) == 3 {
			slots.InUse, _ = strconv.ParseInt(m[1], 10, 64)
			slots.Total, _ = strconv.ParseInt(m[2], 10, 64)
		}
		txt = doc.Find("div#slots>div").Eq(1).Text()
		m = r.FindStringSubmatch(txt)
		if len(m) == 3 {
			slots.ExpInUse, _ = strconv.ParseInt(m[1], 10, 64)
			slots.ExpTotal, _ = strconv.ParseInt(m[2], 10, 64)
		}
	}
	return slots
}

func extractServerTimeFromDocV6(doc *goquery.Document) (time.Time, error) {
	txt := doc.Find("li.OGameClock").First().Text()
	serverTime, err := time.Parse("02.01.2006 15:04:05", txt)
	if err != nil {
		return time.Time{}, err
	}

	u1 := time.Now().UTC().Unix()
	u2 := serverTime.Unix()
	n := int(math.Round(float64(u2-u1)/15)) * 15

	serverTime = serverTime.Add(time.Duration(-n) * time.Second).In(time.FixedZone("OGT", n))

	return serverTime, nil
}

func extractSpioAnzFromDocV6(doc *goquery.Document) int64 {
	out, _ := strconv.ParseInt(doc.Find("input[name=spio_anz]").AttrOr("value", "1"), 10, 64)
	return out
}

func extractDisableChatBarFromDocV6(doc *goquery.Document) bool {
	_, exists := doc.Find("input[name=disableChatBar]").Attr("checked")
	return exists
}

func extractDisableOutlawWarningFromDocV6(doc *goquery.Document) bool {
	_, exists := doc.Find("input[name=disableOutlawWarning]").Attr("checked")
	return exists
}

func extractMobileVersionFromDocV6(doc *goquery.Document) bool {
	_, exists := doc.Find("input[name=mobileVersion]").Attr("checked")
	return exists
}

func extractShowOldDropDownsFromDocV6(doc *goquery.Document) bool {
	_, exists := doc.Find("input[name=showOldDropDowns]").Attr("checked")
	return exists
}

func extractActivateAutofocusFromDocV6(doc *goquery.Document) bool {
	_, exists := doc.Find("input[name=activateAutofocus]").Attr("checked")
	return exists
}

func extractEventsShowFromDocV6(doc *goquery.Document) int64 {
	val, _ := strconv.ParseInt(doc.Find("select[name=eventsShow] option[selected]").AttrOr("value", "1"), 10, 64)
	return val
}

func extractSortSettingFromDocV6(doc *goquery.Document) int64 {
	val, _ := strconv.ParseInt(doc.Find("select#sortSetting option[selected]").AttrOr("value", "0"), 10, 64)
	return val
}

func extractSortOrderFromDocV6(doc *goquery.Document) int64 {
	val, _ := strconv.ParseInt(doc.Find("select#sortOrder option[selected]").AttrOr("value", "0"), 10, 64)
	return val
}

func extractShowDetailOverlayFromDocV6(doc *goquery.Document) bool {
	_, exists := doc.Find("input[name=showDetailOverlay]").Attr("checked")
	return exists
}

func extractAnimatedSlidersFromDocV6(doc *goquery.Document) bool {
	_, exists := doc.Find("input[name=animatedSliders]").Attr("checked")
	return exists
}

func extractAnimatedOverviewFromDocV6(doc *goquery.Document) bool {
	_, exists := doc.Find("input[name=animatedOverview]").Attr("checked")
	return exists
}

func extractPopupsNoticesFromDocV6(doc *goquery.Document) bool {
	_, exists := doc.Find(`input[name="popups[notices]"]`).Attr("checked")
	return exists
}

func extractPopopsCombatreportFromDocV6(doc *goquery.Document) bool {
	_, exists := doc.Find(`input[name="popups[combatreport]"]`).Attr("checked")
	return exists
}

func extractSpioReportPicturesFromDocV6(doc *goquery.Document) bool {
	_, exists := doc.Find("input[name=spioReportPictures]").Attr("checked")
	return exists
}

func extractMsgResultsPerPageFromDocV6(doc *goquery.Document) int64 {
	val, _ := strconv.ParseInt(doc.Find("select[name=msgResultsPerPage] option[selected]").AttrOr("value", "10"), 10, 64)
	return val
}

func extractAuctioneerNotificationsFromDocV6(doc *goquery.Document) bool {
	_, exists := doc.Find("input[name=auctioneerNotifications]").Attr("checked")
	return exists
}

func extractEconomyNotificationsFromDocV6(doc *goquery.Document) bool {
	_, exists := doc.Find("input[name=economyNotifications]").Attr("checked")
	return exists
}

func extractShowActivityMinutesFromDocV6(doc *goquery.Document) bool {
	_, exists := doc.Find("input[name=showActivityMinutes]").Attr("checked")
	return exists
}

func extractPreserveSystemOnPlanetChangeFromDocV6(doc *goquery.Document) bool {
	_, exists := doc.Find("input[name=preserveSystemOnPlanetChange]").Attr("checked")
	return exists
}

func extractNotifBuildListFromDocV6(doc *goquery.Document) bool {
	_, exists := doc.Find(`input[name="notifications[buildList]"]`).Attr("checked")
	return exists
}

func extractNotifFriendlyFleetActivitiesFromDocV6(doc *goquery.Document) bool {
	_, exists := doc.Find(`input[name="notifications[friendlyFleetActivities]"]`).Attr("checked")
	return exists
}

func extractNotifHostileFleetActivitiesFromDocV6(doc *goquery.Document) bool {
	_, exists := doc.Find(`input[name="notifications[hostileFleetActivities]"]`).Attr("checked")
	return exists
}

func extractNotifForeignEspionageFromDocV6(doc *goquery.Document) bool {
	_, exists := doc.Find(`input[name="notifications[foreignEspionage]"]`).Attr("checked")
	return exists
}

func extractNotifAllianceBroadcastsFromDocV6(doc *goquery.Document) bool {
	_, exists := doc.Find(`input[name="notifications[allianceBroadcasts]"]`).Attr("checked")
	return exists
}

func extractNotifAllianceMessagesFromDocV6(doc *goquery.Document) bool {
	_, exists := doc.Find(`input[name="notifications[allianceMessages]"]`).Attr("checked")
	return exists
}

func extractNotifAuctionsFromDocV6(doc *goquery.Document) bool {
	_, exists := doc.Find(`input[name="notifications[auctions]"]`).Attr("checked")
	return exists
}

func extractNotifAccountFromDocV6(doc *goquery.Document) bool {
	_, exists := doc.Find(`input[name="notifications[account]"]`).Attr("checked")
	return exists
}

func extractCommanderFromDocV6(doc *goquery.Document) bool {
	return doc.Find("div#officers a.commander").HasClass("on")
}

func extractAdmiralFromDocV6(doc *goquery.Document) bool {
	return doc.Find("div#officers a.admiral").HasClass("on")
}

func extractEngineerFromDocV6(doc *goquery.Document) bool {
	return doc.Find("div#officers a.engineer").HasClass("on")
}

func extractGeologistFromDocV6(doc *goquery.Document) bool {
	return doc.Find("div#officers a.geologist").HasClass("on")
}

func extractTechnocratFromDocV6(doc *goquery.Document) bool {
	return doc.Find("div#officers a.technocrat").HasClass("on")
}

func extractPlanetCoordinateV6(pageHTML []byte) (Coordinate, error) {
	m := regexp.MustCompile(`<meta name="ogame-planet-coordinates" content="(\d+):(\d+):(\d+)"/>`).FindSubmatch(pageHTML)
	if len(m) == 0 {
		return Coordinate{}, errors.New("planet coordinate not found")
	}
	galaxy, _ := strconv.ParseInt(string(m[1]), 10, 64)
	system, _ := strconv.ParseInt(string(m[2]), 10, 64)
	position, _ := strconv.ParseInt(string(m[3]), 10, 64)
	planetType, _ := extractPlanetTypeV6(pageHTML)
	return Coordinate{galaxy, system, position, planetType}, nil
}

func extractPlanetIDV6(pageHTML []byte) (CelestialID, error) {
	m := regexp.MustCompile(`<meta name="ogame-planet-id" content="(\d+)"/>`).FindSubmatch(pageHTML)
	if len(m) == 0 {
		return 0, errors.New("planet id not found")
	}
	planetID, _ := strconv.ParseInt(string(m[1]), 10, 64)
	return CelestialID(planetID), nil
}

func extractOverviewShipSumCountdownFromBytesV6(pageHTML []byte) int64 {
	var shipSumCountdown int64
	shipSumCountdownMatch := regexp.MustCompile(`getElementByIdWithCache\('shipSumCount7'\),\d+,\d+,(\d+),`).FindSubmatch(pageHTML)
	if len(shipSumCountdownMatch) > 0 {
		shipSumCountdown = int64(toInt(shipSumCountdownMatch[1]))
	}
	return shipSumCountdown
}

func extractOGameTimestampFromBytesV6(pageHTML []byte) int64 {
	m := regexp.MustCompile(`<meta name="ogame-timestamp" content="(\d+)"/>`).FindSubmatch(pageHTML)
	if len(m) != 2 {
		return 0
	}
	ts, _ := strconv.ParseInt(string(m[1]), 10, 64)
	return ts
}

func extractPlanetTypeV6(pageHTML []byte) (CelestialType, error) {
	m := regexp.MustCompile(`<meta name="ogame-planet-type" content="(\w+)"/>`).FindSubmatch(pageHTML)
	if len(m) == 0 {
		return 0, errors.New("planet type not found")
	}
	if bytes.Equal(m[1], []byte("planet")) {
		return PlanetType, nil
	} else if bytes.Equal(m[1], []byte("moon")) {
		return MoonType, nil
	}
	return 0, errors.New("invalid planet type : " + string(m[1]))
}

func extractAjaxChatTokenV6(pageHTML []byte) (string, error) {
	r1 := regexp.MustCompile(`ajaxChatToken\s?=\s?['"](\w+)['"]`)
	m1 := r1.FindSubmatch(pageHTML)
	if len(m1) < 2 {
		return "", errors.New("unable to find token")
	}
	token := string(m1[1])
	return token, nil
}

func extractUserInfosV6(pageHTML []byte, lang string) (UserInfos, error) {
	playerIDRgx := regexp.MustCompile(`<meta name="ogame-player-id" content="(\d+)"/>`)
	playerNameRgx := regexp.MustCompile(`<meta name="ogame-player-name" content="([^"]+)"/>`)
	txtContent := regexp.MustCompile(`textContent\[7]\s?=\s?"([^"]+)"`)
	playerIDGroups := playerIDRgx.FindSubmatch(pageHTML)
	playerNameGroups := playerNameRgx.FindSubmatch(pageHTML)
	subHTMLGroups := txtContent.FindSubmatch(pageHTML)
	if len(playerIDGroups) < 2 {
		return UserInfos{}, errors.New("cannot find player id")
	}
	if len(playerNameGroups) < 2 {
		return UserInfos{}, errors.New("cannot find player name")
	}
	if len(subHTMLGroups) < 2 {
		return UserInfos{}, errors.New("cannot find sub html")
	}
	res := UserInfos{}
	res.PlayerID = int64(toInt(playerIDGroups[1]))
	res.PlayerName = string(playerNameGroups[1])
	html2 := subHTMLGroups[1]

	infosRgx := regexp.MustCompile(`([\d\\.]+) \(Place ([\d.]+) of ([\d.]+)\)`)
	switch lang {
	case "fr":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(Place ([\d.]+) sur ([\d.]+)\)`)
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
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(([\d.]+) \u4eba\u4e2d\u7684\u7b2c ([\d.]+) \u4f4d\)`)
	case "cz":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(Pozice ([\d.]+) z ([\d.]+)\)`)
	case "de":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(Platz ([\d.]+) von ([\d.]+)\)`)
	case "es":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(Lugar ([\d.]+) de ([\d.]+)\)`)
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
		return UserInfos{}, errors.New("cannot find infos in sub html")
	}
	res.Points = ParseInt(string(infos[1]))
	res.Rank = ParseInt(string(infos[2]))
	res.Total = ParseInt(string(infos[3]))
	if lang == "tr" || lang == "jp" {
		res.Rank = ParseInt(string(infos[3]))
		res.Total = ParseInt(string(infos[2]))
	}
	honourPointsRgx := regexp.MustCompile(`textContent\[9]\s?=\s?"([^"]+)"`)
	honourPointsGroups := honourPointsRgx.FindSubmatch(pageHTML)
	if len(honourPointsGroups) < 2 {
		return UserInfos{}, errors.New("cannot find honour points")
	}
	res.HonourPoints = ParseInt(string(honourPointsGroups[1]))
	return res, nil
}

func extractResourcesDetailsV6(pageHTML []byte) (out ResourcesDetails, err error) {
	var res resourcesResp
	if err = json.Unmarshal(pageHTML, &res); err != nil {
		if isLogged(pageHTML) {
			return out, ErrInvalidPlanetID
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
	out.Metal.CurrentProduction = ParseInt(metalDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	out.Crystal.CurrentProduction = ParseInt(crystalDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	out.Deuterium.CurrentProduction = ParseInt(deuteriumDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	out.Energy.CurrentProduction = ParseInt(energyDoc.Find("table tr").Eq(1).Find("td").Eq(0).Text())
	out.Energy.Consumption = ParseInt(energyDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	out.Darkmatter.Purchased = ParseInt(darkmatterDoc.Find("table tr").Eq(1).Find("td").Eq(0).Text())
	out.Darkmatter.Found = ParseInt(darkmatterDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	return
}

func extractCoordV6(v string) (coord Coordinate) {
	coordRgx := regexp.MustCompile(`\[(\d+):(\d+):(\d+)]`)
	m := coordRgx.FindStringSubmatch(v)
	if len(m) == 4 {
		coord.Galaxy, _ = strconv.ParseInt(m[1], 10, 64)
		coord.System, _ = strconv.ParseInt(m[2], 10, 64)
		coord.Position, _ = strconv.ParseInt(m[3], 10, 64)
	}
	return
}

func extractGalaxyInfosV6(pageHTML []byte, botPlayerName string, botPlayerID, botPlayerRank int64) (SystemInfos, error) {
	prefixedNumRgx := regexp.MustCompile(`.*: ([\d.]+)`)

	extractActivity := func(activityDiv *goquery.Selection) int64 {
		var activity int64
		if activityDiv != nil {
			activityDivClass := activityDiv.AttrOr("class", "")
			if strings.Contains(activityDivClass, "minute15") {
				activity = 15
			} else if strings.Contains(activityDivClass, "showMinutes") {
				activity, _ = strconv.ParseInt(strings.TrimSpace(activityDiv.Text()), 10, 64)
			}
		}
		return activity
	}

	var tmp struct {
		Galaxy string
	}
	var res SystemInfos
	if err := json.Unmarshal(pageHTML, &tmp); err != nil {
		return res, ErrNotLogged
	}
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(tmp.Galaxy))
	res.galaxy = ParseInt(doc.Find("table").AttrOr("data-galaxy", "0"))
	res.system = ParseInt(doc.Find("table").AttrOr("data-system", "0"))
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

			planetInfos := new(PlanetInfos)
			planetInfos.ID, _ = strconv.ParseInt(s.Find("td.colonized").AttrOr("data-planet-id", ""), 10, 64)

			moonID, _ := strconv.ParseInt(s.Find("td.moon").AttrOr("data-moon-id", ""), 10, 64)
			moonSize, _ := strconv.ParseInt(strings.Split(s.Find("td.moon span#moonsize").Text(), " ")[0], 10, 64)
			if moonID > 0 {
				planetInfos.Moon = new(MoonInfos)
				planetInfos.Moon.ID = moonID
				planetInfos.Moon.Diameter = moonSize
				planetInfos.Moon.Activity = extractActivity(s.Find("td.moon div.activity"))
			}

			allianceSpan := s.Find("span.allytagwrapper")
			if allianceSpan.Size() > 0 {
				longID, _ := allianceSpan.Attr("rel")
				planetInfos.Alliance = new(AllianceInfos)
				planetInfos.Alliance.Name = allianceSpan.Find("h1").Text()
				planetInfos.Alliance.ID, _ = strconv.ParseInt(strings.TrimPrefix(longID, "alliance"), 10, 64)
				planetInfos.Alliance.Rank, _ = strconv.ParseInt(allianceSpan.Find("ul.ListLinks li").First().Find("a").Text(), 10, 64)
				planetInfos.Alliance.Member = ParseInt(prefixedNumRgx.FindStringSubmatch(allianceSpan.Find("ul.ListLinks li").Eq(1).Text())[1])
			}

			if len(prefixedNumRgx.FindStringSubmatch(metalTxt)) > 0 {
				planetInfos.Debris.Metal = ParseInt(prefixedNumRgx.FindStringSubmatch(metalTxt)[1])
				planetInfos.Debris.Crystal = ParseInt(prefixedNumRgx.FindStringSubmatch(crystalTxt)[1])
				planetInfos.Debris.RecyclersNeeded = ParseInt(prefixedNumRgx.FindStringSubmatch(recyclersTxt)[1])
			}

			planetInfos.Activity = extractActivity(s.Find("td:not(.moon) div.activity"))
			planetInfos.Name = planetName
			planetInfos.Img = planetImg
			planetInfos.Inactive = strings.Contains(classes, "inactive_filter")
			planetInfos.StrongPlayer = strings.Contains(classes, "strong_filter")
			planetInfos.Newbie = strings.Contains(classes, "newbie_filter")
			planetInfos.Vacation = strings.Contains(classes, "vacation_filter")
			planetInfos.HonorableTarget = s.Find("span.status_abbr_honorableTarget").Size() > 0
			planetInfos.Administrator = s.Find("span.status_abbr_admin").Size() > 0
			planetInfos.Banned = s.Find("td.playername a span.status_abbr_banned").Size() > 0
			tdPlayername := s.Find("td.playername span")
			planetInfos.Player.IsBandit = tdPlayername.HasClass("rank_bandit1") || tdPlayername.HasClass("rank_bandit2") || tdPlayername.HasClass("rank_bandit3")
			planetInfos.Player.IsStarlord = tdPlayername.HasClass("rank_starlord1") || tdPlayername.HasClass("rank_starlord2") || tdPlayername.HasClass("rank_starlord3")
			planetInfos.Coordinate = extractCoordV6(coordsRaw)
			planetInfos.Coordinate.Type = PlanetType

			var playerID int64
			var playerName string
			var playerRank int64
			if len(tooltips.Nodes) > 1 {
				tooltips.Each(func(i int, s *goquery.Selection) {
					idAttr, _ := s.Attr("id")
					if strings.HasPrefix(idAttr, "player") {
						playerID, _ = strconv.ParseInt(regexp.MustCompile(`player(\d+)`).FindStringSubmatch(idAttr)[1], 10, 64)
						playerName = s.Find("h1").Find("span").Text()
						playerRank, _ = strconv.ParseInt(s.Find("li.rank").Find("a").Text(), 10, 64)
					}
				})
			}
			if playerName == "" {
				playerName := strings.TrimSpace(s.Find("td.playername").Find("span").Text())
				if playerName == "" {
					return
				}
			}

			if playerID == 0 {
				playerID = botPlayerID
				playerName = botPlayerName
				playerRank = botPlayerRank
			}

			planetInfos.Player.ID = playerID
			planetInfos.Player.Name = playerName
			planetInfos.Player.Rank = playerRank

			res.planets[i] = planetInfos
		}
	})

	debris16Div := doc.Find("div#debris16")
	if debris16Div.Size() > 0 {
		lis := debris16Div.Find("ul.ListLinks li")
		metalTxt := lis.First().Text()
		crystalTxt := lis.Eq(1).Text()
		pathfindersTxt := lis.Eq(2).Text()
		res.ExpeditionDebris.Metal = ParseInt(prefixedNumRgx.FindStringSubmatch(metalTxt)[1])
		res.ExpeditionDebris.Crystal = ParseInt(prefixedNumRgx.FindStringSubmatch(crystalTxt)[1])
		res.ExpeditionDebris.PathfindersNeeded = ParseInt(prefixedNumRgx.FindStringSubmatch(pathfindersTxt)[1])
	}

	debris17Div := doc.Find("div#debris17")
	if debris17Div.Size() > 0 {
		lis := debris17Div.Find("ul.ListLinks li")
		darkmatterTxt := lis.First().Text()
		darkmatterMatches := prefixedNumRgx.FindStringSubmatch(darkmatterTxt)
		if len(darkmatterMatches) == 2 {
			res.Events.Darkmatter = ParseInt(darkmatterMatches[1])
		}
	}

	planet17Div := doc.Find("div#planet17")
	if planet17Div.Size() > 0 {
		res.Events.HasAsteroid = true
	}

	return res, nil
}

func extractPhalanxV6(pageHTML []byte) ([]Fleet, error) {
	res := make([]Fleet, 0)
	var ogameTimestamp int64
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	eventFleet := doc.Find("div.eventFleet")
	if eventFleet.Size() == 0 {
		txt := doc.Find("div#phalanxEventContent").Text()
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
		ogameTimestamp, _ = strconv.ParseInt(string(m[1]), 10, 64)
	}

	eventFleet.Each(func(i int, s *goquery.Selection) {
		mission, _ := strconv.ParseInt(s.AttrOr("data-mission-type", "0"), 10, 64)
		returning, _ := strconv.ParseBool(s.AttrOr("data-return-flight", "false"))
		arrivalTime, _ := strconv.ParseInt(s.AttrOr("data-arrival-time", "0"), 10, 64)
		arriveIn := arrivalTime - ogameTimestamp
		if arriveIn < 0 {
			arriveIn = 0
		}
		originFleetFigure := s.Find("li.originFleet figure")
		originTxt := s.Find("li.coordsOrigin a").Text()
		destTxt := s.Find("li.destCoords a").Text()

		fleet := Fleet{}

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
				nbr := ParseInt(s.Find("td").Eq(1).Text())
				if name != "" && nbr > 0 {
					fleet.Ships.Set(name2id(name), nbr)
				}
			})
		}

		fleet.Mission = MissionID(mission)
		fleet.ReturnFlight = returning
		fleet.ArriveIn = arriveIn
		fleet.Origin = extractCoordV6(originTxt)
		fleet.Origin.Type = PlanetType
		if originFleetFigure.HasClass("moon") {
			fleet.Origin.Type = MoonType
		}
		fleet.Destination = extractCoordV6(destTxt)
		fleet.Destination.Type = PlanetType
		res = append(res, fleet)
	})
	return res, nil
}

func extractJumpGateV6(pageHTML []byte) (ShipsInfos, string, []MoonID, int64) {
	m := regexp.MustCompile(`\$\("#cooldown"\), (\d+),`).FindSubmatch(pageHTML)
	ships := ShipsInfos{}
	var destinations []MoonID
	if len(m) > 0 {
		waitTime := int64(toInt(m[1]))
		return ships, "", destinations, waitTime
	}
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	for _, s := range Ships {
		ships.Set(s.GetID(), ParseInt(doc.Find("input#ship_"+strconv.FormatInt(int64(s.GetID()), 10)).AttrOr("rel", "0")))
	}
	token := doc.Find("input[name=token]").AttrOr("value", "")

	doc.Find("select[name=zm] option").Each(func(i int, s *goquery.Selection) {
		moonID := ParseInt(s.AttrOr("value", "0"))
		if moonID > 0 {
			destinations = append(destinations, MoonID(moonID))
		}
	})

	return ships, token, destinations, 0
}

func extractFederationV6(pageHTML []byte) url.Values {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	payload := extractHiddenFieldsFromDocV6(doc)
	groupName := doc.Find("input#groupNameInput").AttrOr("value", "")
	doc.Find("ul#participantselect li").Each(func(i int, s *goquery.Selection) {
		payload.Add("unionUsers", s.Text())
	})
	payload.Add("groupname", groupName)
	return payload
}

func extractConstructionsV6(pageHTML []byte) (buildingID ID, buildingCountdown int64, researchID ID, researchCountdown int64) {
	buildingCountdownMatch := regexp.MustCompile(`getElementByIdWithCache\("Countdown"\),(\d+),`).FindSubmatch(pageHTML)
	if len(buildingCountdownMatch) > 0 {
		buildingCountdown = int64(toInt(buildingCountdownMatch[1]))
		buildingIDInt := toInt(regexp.MustCompile(`onclick="cancelProduction\((\d+),`).FindSubmatch(pageHTML)[1])
		buildingID = ID(buildingIDInt)
	}
	researchCountdownMatch := regexp.MustCompile(`getElementByIdWithCache\("researchCountdown"\),(\d+),`).FindSubmatch(pageHTML)
	if len(researchCountdownMatch) > 0 {
		researchCountdown = int64(toInt(researchCountdownMatch[1]))
		researchIDInt := toInt(regexp.MustCompile(`onclick="cancelResearch\((\d+),`).FindSubmatch(pageHTML)[1])
		researchID = ID(researchIDInt)
	}
	return
}

func extractFleetDeutSaveFactorV6(pageHTML []byte) float64 {
	factor := 1.0
	m := regexp.MustCompile(`var fleetDeutSaveFactor=([+-]?([0-9]*[.])?[0-9]+);`).FindSubmatch(pageHTML)
	if len(m) > 0 {
		factor, _ = strconv.ParseFloat(string(m[1]), 64)
	}
	return factor
}

func extractCancelBuildingInfosV6(pageHTML []byte) (token string, techID, listID int64, err error) {
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
	techID, _ = strconv.ParseInt(m[1], 10, 64)
	listID, _ = strconv.ParseInt(m[2], 10, 64)
	return
}

func extractCancelResearchInfosV6(pageHTML []byte) (token string, techID, listID int64, err error) {
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
	techID, _ = strconv.ParseInt(m[1], 10, 64)
	listID, _ = strconv.ParseInt(m[2], 10, 64)
	return
}

func extractUniverseSpeedV6(pageHTML []byte) int64 {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	spans := doc.Find("span.undermark")
	level := ParseInt(spans.Eq(0).Text())
	val := ParseInt(spans.Eq(1).Text())
	metalProduction := int64(math.Floor(30 * float64(level) * math.Pow(1.1, float64(level))))
	universeSpeed := val / metalProduction
	return universeSpeed
}

var planetInfosRgx = regexp.MustCompile(`([^\[]+) \[(\d+):(\d+):(\d+)]([\d.,]+)(?i)(?:km|км|公里|χμ) \((\d+)/(\d+)\)(?:de|da|od|mellem|от)?\s*([-\d]+).+C\s*(?:bis|para|to|à|至|a|～|do|ile|tot|og|до|až|til|la|έως)\s*([-\d]+).+C`)
var moonInfosRgx = regexp.MustCompile(`([^\[]+) \[(\d+):(\d+):(\d+)]([\d.]+)(?i)(?:km|км|χμ) \((\d+)/(\d+)\)`)
var cpRgx = regexp.MustCompile(`&cp=(\d+)`)

func extractPlanetFromSelectionV6(s *goquery.Selection, b *OGame) (Planet, error) {
	el, _ := s.Attr("id")
	id, err := strconv.ParseInt(strings.TrimPrefix(el, "planet-"), 10, 64)
	if err != nil {
		return Planet{}, err
	}

	title, _ := s.Find("a.planetlink").Attr("title")
	root, err := html.Parse(strings.NewReader(title))
	if err != nil {
		return Planet{}, err
	}

	txt := goquery.NewDocumentFromNode(root).Text()
	m := planetInfosRgx.FindStringSubmatch(txt)
	if len(m) < 10 {
		return Planet{}, errors.New("failed to parse planet infos: " + txt)
	}

	res := Planet{}
	res.ogame = b
	res.Img = s.Find("img.planetPic").AttrOr("src", "")
	res.ID = PlanetID(id)
	res.Name = strings.TrimSpace(m[1])
	res.Coordinate.Galaxy, _ = strconv.ParseInt(m[2], 10, 64)
	res.Coordinate.System, _ = strconv.ParseInt(m[3], 10, 64)
	res.Coordinate.Position, _ = strconv.ParseInt(m[4], 10, 64)
	res.Coordinate.Type = PlanetType
	res.Diameter = ParseInt(m[5])
	res.Fields.Built, _ = strconv.ParseInt(m[6], 10, 64)
	res.Fields.Total, _ = strconv.ParseInt(m[7], 10, 64)
	res.Temperature.Min, _ = strconv.ParseInt(m[8], 10, 64)
	res.Temperature.Max, _ = strconv.ParseInt(m[9], 10, 64)

	res.Moon, _ = extractMoonFromPlanetSelectionV6(s, b)

	return res, nil
}

func extractMoonFromPlanetSelectionV6(s *goquery.Selection, b *OGame) (*Moon, error) {
	moonLink := s.Find("a.moonlink")
	moon, err := extractMoonFromSelectionV6(moonLink, b)
	if err != nil {
		return nil, err
	}
	return &moon, nil
}

func extractMoonFromSelectionV6(moonLink *goquery.Selection, b *OGame) (Moon, error) {
	href, found := moonLink.Attr("href")
	if !found {
		return Moon{}, errors.New("no moon found")
	}
	m := cpRgx.FindStringSubmatch(href)
	id, _ := strconv.ParseInt(m[1], 10, 64)
	title, _ := moonLink.Attr("title")
	root, err := html.Parse(strings.NewReader(title))
	if err != nil {
		return Moon{}, err
	}
	txt := goquery.NewDocumentFromNode(root).Text()
	mm := moonInfosRgx.FindStringSubmatch(txt)
	if len(mm) < 8 {
		return Moon{}, errors.New("failed to parse moon infos: " + txt)
	}
	moon := Moon{}
	moon.ogame = b
	moon.ID = MoonID(id)
	moon.Name = strings.TrimSpace(mm[1])
	moon.Coordinate.Galaxy, _ = strconv.ParseInt(mm[2], 10, 64)
	moon.Coordinate.System, _ = strconv.ParseInt(mm[3], 10, 64)
	moon.Coordinate.Position, _ = strconv.ParseInt(mm[4], 10, 64)
	moon.Coordinate.Type = MoonType
	moon.Diameter = ParseInt(mm[5])
	moon.Fields.Built, _ = strconv.ParseInt(mm[6], 10, 64)
	moon.Fields.Total, _ = strconv.ParseInt(mm[7], 10, 64)
	moon.Img = moonLink.Find("img.icon-moon").AttrOr("src", "")
	return moon, nil
}

func extractEmpire(html string, nbr int64) (interface{}, error) {
	if nbr > 1 {
		return nil, errors.New("invalid number for Empire page")
	}
	m := regexp.MustCompile(`createImperiumHtml\("#mainWrapper",\s"#loading",\s(.*),\s\d+\s\);`).FindStringSubmatch(html)
	if len(m) != 2 {
		return nil, errors.New("regexp for Empire JSON did not match anything")
	}
	var empireJSON interface{}
	if err := json.Unmarshal([]byte(m[1]), &empireJSON); err != nil {
		return nil, err
	}
	return empireJSON, nil
}

// Auction ...
type Auction struct {
	HasFinished         bool
	Endtime             int64
	NumBids             int64
	CurrentBid          int64
	AlreadyBid          int64
	MinimumBid          int64
	DeficitBid          int64
	HighestBidder       string
	HighestBidderUserID int64
	CurrentItem         string
	CurrentItemLong     string
	Inventory           int64
	Token               string
	ResourceMultiplier  struct {
		Metal     float64
		Crystal   float64
		Deuterium float64
		Honor     int64
	}
	Resources map[string]interface{}
}

// String ...
func (a Auction) String() string {
	return "" +
		"  Has finished: " + strconv.FormatBool(a.HasFinished) + "\n" +
		"      End time: " + strconv.FormatInt(a.Endtime, 10) + "\n" +
		"      Num bids: " + strconv.FormatInt(a.NumBids, 10) + "\n" +
		"   Minimum bid: " + strconv.FormatInt(a.MinimumBid, 10) + "\n" +
		"Highest bidder: " + a.HighestBidder + " (" + strconv.FormatInt(a.HighestBidderUserID, 10) + ")" + "\n" +
		"  Current item: " + a.CurrentItem + " (" + a.CurrentItemLong + ")" + "\n" +
		"     Inventory: " + strconv.FormatInt(a.Inventory, 10) + "\n" +
		""
}

// ExtractAuction extract auction information from page "traderAuctioneer"
func extractAuctionFromDoc(doc *goquery.Document) (Auction, error) {
	auction := Auction{}
	auction.HasFinished = false

	// Detect if Auction has already finished
	nextAuction := doc.Find("#nextAuction")
	if nextAuction.Size() > 0 {
		// Find time until next auction starts
		auction.Endtime, _ = strconv.ParseInt(nextAuction.Text(), 10, 64)
		auction.HasFinished = true
	} else {
		endAtApprox := doc.Find("p.auction_info b").Text()
		m := regexp.MustCompile(`[^\d]+(\d+).*`).FindStringSubmatch(endAtApprox)
		if len(m) != 2 {
			return Auction{}, errors.New("failed to find end time approx")
		}
		endTimeMinutes, err := strconv.ParseInt(m[1], 10, 64)
		if err != nil {
			return Auction{}, errors.New("invalid end time approx: " + err.Error())
		}
		auction.Endtime = endTimeMinutes * 60
	}

	auction.HighestBidder = strings.TrimSpace(doc.Find("a.currentPlayer").Text())
	auction.HighestBidderUserID, _ = strconv.ParseInt(doc.Find("a.currentPlayer").AttrOr("data-player-id", ""), 10, 64)
	auction.NumBids, _ = strconv.ParseInt(doc.Find("div.numberOfBids").Text(), 10, 64)
	auction.CurrentBid = ParseInt(doc.Find("div.currentSum").Text())
	auction.Inventory, _ = strconv.ParseInt(doc.Find("span.level.amount").Text(), 10, 64)
	auction.CurrentItem = strings.ToLower(doc.Find("img").First().AttrOr("alt", ""))
	auction.CurrentItemLong = strings.ToLower(doc.Find("div.image_140px").First().Find("a").First().AttrOr("title", ""))
	multiplierRegex := regexp.MustCompile(`multiplier=([^;]+);`).FindStringSubmatch(doc.Text())
	if len(multiplierRegex) != 2 {
		return Auction{}, errors.New("failed to find auction multiplier")
	}
	if err := json.Unmarshal([]byte(multiplierRegex[1]), &auction.ResourceMultiplier); err != nil {
		return Auction{}, errors.New("failed to json parse auction multiplier: " + err.Error())
	}

	// Find auctioneer token
	tokenRegex := regexp.MustCompile(`auctioneerToken="([^"]+)";`).FindStringSubmatch(doc.Text())
	if len(tokenRegex) != 2 {
		return Auction{}, errors.New("failed to find auctioneer token")
	}
	auction.Token = tokenRegex[1]

	// Find Planet / Moon resources JSON
	planetMoonResources := regexp.MustCompile(`planetResources=([^;]+);`).FindStringSubmatch(doc.Text())
	if len(planetMoonResources) != 2 {
		return Auction{}, errors.New("failed to find planetResources")
	}
	if err := json.Unmarshal([]byte(planetMoonResources[1]), &auction.Resources); err != nil {
		return Auction{}, errors.New("failed to json unmarshal planetResources: " + err.Error())
	}

	// Find already-bid
	auction.AlreadyBid = ParseInt(doc.Find("table.table_ressources_sum tr td.auctionInfo.js_alreadyBidden").Text())

	// Find min-bid
	auction.MinimumBid = ParseInt(doc.Find("table.table_ressources_sum tr td.auctionInfo.js_price").Text())

	// Find deficit-bid
	auction.DeficitBid = ParseInt(doc.Find("table.table_ressources_sum tr td.auctionInfo.js_deficit").Text())

	// Note: Don't just bid the min-bid amount. It will keep doubling the total bid and grow exponentially...
	// DeficitBid is 1000 when another player has outbid you or if nobody has bid yet.
	// DeficitBid seems to be filled by Javascript in the browser. We're parsing it anyway. Correct Bid calculation would be:
	// bid = max(auction.DeficitBid, auction.MinimumBid - auction.AlreadyBid)

	return auction, nil
}
