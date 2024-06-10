package v11_15_0

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/PuerkitoBio/goquery"
	v6 "github.com/alaingilbert/ogame/pkg/extractor/v6"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/utils"
	"golang.org/x/net/html"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func extractCombatReportMessagesFromDoc(doc *goquery.Document) ([]ogame.CombatReportSummary, int64, error) {
	msgs := make([]ogame.CombatReportSummary, 0)
	doc.Find(".msg").Each(func(i int, s *goquery.Selection) {
		if idStr, exists := s.Attr("data-msg-id"); exists {
			if id, err := utils.ParseI64(idStr); err == nil {
				rawMessageData := s.Find("div.rawMessageData")
				resultStr := rawMessageData.AttrOr("data-raw-result", "")
				var result struct {
					Loot struct {
						Percentage int64
						Resources  []struct {
							Resource string
							Amount   int64
						}
					}
				}
				_ = json.Unmarshal([]byte(resultStr), &result)

				report := ogame.CombatReportSummary{ID: id}
				report.Destination = v6.ExtractCoord(s.Find("div.msgHead a").Text())
				report.Destination.Type = ogame.PlanetType
				if s.Find("div.msgHead figure").HasClass("moon") {
					report.Destination.Type = ogame.MoonType
				}
				apiKeyTitle := s.Find("button.icon_apikey").AttrOr("title", "")
				m := regexp.MustCompile(`'(cr-[^']+)'`).FindStringSubmatch(apiKeyTitle)
				if len(m) == 2 {
					report.APIKey = m[1]
				}

				for _, resource := range result.Loot.Resources {
					res := resource.Resource
					if ogame.IsStrDeuterium(res) {
						report.Deuterium = resource.Amount
					} else if ogame.IsStrCrystal(res) {
						report.Crystal = resource.Amount
					} else if ogame.IsStrMetal(res) {
						report.Metal = resource.Amount
					}
				}

				debrisFieldTitle := s.Find("span.msg_content div.combatLeftSide span").Eq(2).AttrOr("title", "0")
				report.DebrisField = utils.ParseInt(debrisFieldTitle)
				resText := s.Find("span.msg_content div.combatLeftSide span").Eq(1).Text()
				m = regexp.MustCompile(`[\d.,]+[^\d]*([\d.,]+)`).FindStringSubmatch(resText)
				if len(m) == 2 {
					report.Loot = utils.ParseInt(m[1])
				}
				msgDate, _ := time.Parse("02.01.2006 15:04:05", s.Find("div.msgDate").Text())
				report.CreatedAt = msgDate

				link := s.Find("message-footer.msg_actions button.msgAttackBtn").AttrOr("onclick", "")
				m = regexp.MustCompile(`page=ingame&component=fleetdispatch&galaxy=(\d+)&system=(\d+)&position=(\d+)&type=(\d+)&`).FindStringSubmatch(link)
				if len(m) != 5 {
					return
				}
				galaxy := utils.DoParseI64(m[1])
				system := utils.DoParseI64(m[2])
				position := utils.DoParseI64(m[3])
				planetType := utils.DoParseI64(m[4])
				report.Origin = &ogame.Coordinate{Galaxy: galaxy, System: system, Position: position, Type: ogame.CelestialType(planetType)}
				if report.Origin.Equal(report.Destination) {
					report.Origin = nil
				}

				msgs = append(msgs, report)
			}
		}
	})
	return msgs, 1, nil
}

func extractEspionageReportMessageIDsFromDoc(doc *goquery.Document) ([]ogame.EspionageReportSummary, int64, error) {
	msgs := make([]ogame.EspionageReportSummary, 0)
	doc.Find(".msg").Each(func(i int, s *goquery.Selection) {
		rawData := s.Find("div.rawMessageData").First()
		if idStr, exists := s.Attr("data-msg-id"); exists {
			if id, err := utils.ParseI64(idStr); err == nil {
				messageType := ogame.Report
				if s.Find("span.espionageDefText").Size() > 0 {
					messageType = ogame.Action
				}
				report := ogame.EspionageReportSummary{ID: id, Type: messageType}
				report.From = s.Find(".msgSender").Text()
				targetStr := rawData.AttrOr("data-raw-coordinates", "")
				report.Target, _ = ogame.ParseCoord(targetStr)
				report.Target.Type = ogame.PlanetType
				planetType := rawData.AttrOr("data-raw-targetplanettype", "1")
				if planetType == "3" {
					report.Target.Type = ogame.MoonType
				}
				if messageType == ogame.Report {
					s.Find(".lootPercentage").Each(func(i int, s *goquery.Selection) {
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
	return msgs, 1, nil
}

func extractEspionageReportFromDoc(doc *goquery.Document, location *time.Location) (ogame.EspionageReport, error) {
	report := ogame.EspionageReport{}
	report.ID = utils.DoParseI64(doc.Find("div.detail_msg").AttrOr("data-msg-id", "0"))
	rawMessageData := doc.Find("div.rawMessageData").First()
	txt := rawMessageData.AttrOr("data-raw-coordinates", "")
	report.Coordinate = ogame.DoParseCoord(txt)
	if rawMessageData.AttrOr("data-raw-targetPlanetType", "1") == "1" {
		report.Coordinate.Type = ogame.PlanetType
	} else {
		report.Coordinate.Type = ogame.MoonType
	}
	messageType := ogame.Report
	//if doc.Find("span.espionageDefText").Size() > 0 {
	//	messageType = ogame.Action
	//}
	report.Type = messageType

	msgDateRaw := doc.Find("span.msg_date").Text()
	msgDate, _ := time.ParseInLocation("02.01.2006 15:04:05", msgDateRaw, location)
	report.Date = msgDate.In(time.Local)

	report.Username = strings.TrimSpace(rawMessageData.AttrOr("data-raw-playername", ""))

	characterClassJsonStr := strings.TrimSpace(rawMessageData.AttrOr("data-raw-characterclass", ""))
	var characterClassStruct struct{ ID int }
	_ = json.Unmarshal([]byte(characterClassJsonStr), &characterClassStruct)
	switch characterClassStruct.ID {
	case 1:
		report.CharacterClass = ogame.Collector
	case 2:
		report.CharacterClass = ogame.General
	case 3:
		report.CharacterClass = ogame.Discoverer
	default:
		report.CharacterClass = ogame.NoClass
	}

	allianceClassJsonStr := strings.TrimSpace(rawMessageData.AttrOr("data-raw-allianceclass", ""))
	var allianceClassStruct struct{ ID int }
	_ = json.Unmarshal([]byte(allianceClassJsonStr), &allianceClassStruct)
	switch allianceClassStruct.ID {
	case 1:
		report.AllianceClass = ogame.Warrior
	case 2:
		report.AllianceClass = ogame.Trader
	case 3:
		report.AllianceClass = ogame.Researcher
	default:
		report.AllianceClass = ogame.NoAllianceClass
	}

	// Bandit, Starlord
	banditstarlord := doc.Find("span.honorRank").First()
	if banditstarlord.HasClass("honorRank") {
		report.IsBandit = banditstarlord.HasClass("rank_bandit1") || banditstarlord.HasClass("rank_bandit2") || banditstarlord.HasClass("rank_bandit3")
		report.IsStarlord = banditstarlord.HasClass("rank_starlord1") || banditstarlord.HasClass("rank_starlord2") || banditstarlord.HasClass("rank_starlord3")
	}

	report.HonorableTarget = doc.Find("span.status_abbr_honorableTarget").Length() > 0

	// IsInactive, IsLongInactive
	inactive := doc.Find("div.playerInfo").First().Find("span")
	if inactive.HasClass("status_abbr_longinactive") {
		report.IsInactive = true
		report.IsLongInactive = true
	} else if inactive.HasClass("status_abbr_inactive") {
		report.IsInactive = true
	}

	// APIKey
	apikey, _ := doc.Find("button.icon_apikey").Attr("title")
	apiDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(apikey))
	report.APIKey = apiDoc.Find("input").First().AttrOr("value", "")

	// Inactivity timer
	report.LastActivity = utils.ParseInt(rawMessageData.AttrOr("data-raw-activity", "-1"))
	if report.LastActivity == -1 {
		report.LastActivity = 0
	}

	// CounterEspionage
	report.CounterEspionage = utils.ParseInt(rawMessageData.AttrOr("data-raw-counterespionagechance", "0"))

	report.Metal = utils.DoParseI64(rawMessageData.AttrOr("data-raw-metal", "0"))
	report.Crystal = utils.DoParseI64(rawMessageData.AttrOr("data-raw-crystal", "0"))
	report.Deuterium = utils.DoParseI64(rawMessageData.AttrOr("data-raw-deuterium", "0"))
	report.Food = utils.DoParseI64(rawMessageData.AttrOr("data-raw-food", "0"))
	report.Population = utils.DoParseI64(rawMessageData.AttrOr("data-raw-population", "0"))
	report.Energy = utils.ParseInt(doc.Find("resource-icon.energy").Next().Text())

	report.HasBuildingsInformation = rawMessageData.AttrOr("data-raw-hiddenbuildings", "1") == ""
	if report.HasBuildingsInformation {
		buildingsStr := rawMessageData.AttrOr("data-raw-buildings", "{}")
		var buildingsStruct struct {
			MetalMine            *int64 `json:"1"`
			CrystalMine          *int64 `json:"2"`
			DeuteriumSynthesizer *int64 `json:"3"`
			SolarPlant           *int64 `json:"4"`
			FusionReactor        *int64 `json:"12"`
			MetalStorage         *int64 `json:"22"`
			CrystalStorage       *int64 `json:"23"`
			DeuteriumTank        *int64 `json:"24"`
			AllianceDepot        *int64 `json:"34"`
			RoboticsFactory      *int64 `json:"14"`
			Shipyard             *int64 `json:"21"`
			ResearchLab          *int64 `json:"31"`
			MissileSilo          *int64 `json:"44"`
			NaniteFactory        *int64 `json:"15"`
			Terraformer          *int64 `json:"33"`
			SpaceDock            *int64 `json:"36"`
			LunarBase            *int64 `json:"41"`
			SensorPhalanx        *int64 `json:"42"`
			JumpGate             *int64 `json:"43"`
		}
		_ = json.Unmarshal([]byte(buildingsStr), &buildingsStruct)
		report.MetalMine = buildingsStruct.MetalMine
		report.CrystalMine = buildingsStruct.CrystalMine
		report.DeuteriumSynthesizer = buildingsStruct.DeuteriumSynthesizer
		report.SolarPlant = buildingsStruct.SolarPlant
		report.FusionReactor = buildingsStruct.FusionReactor
		report.MetalStorage = buildingsStruct.MetalStorage
		report.CrystalStorage = buildingsStruct.CrystalStorage
		report.DeuteriumTank = buildingsStruct.DeuteriumTank
		report.AllianceDepot = buildingsStruct.AllianceDepot
		report.RoboticsFactory = buildingsStruct.RoboticsFactory
		report.Shipyard = buildingsStruct.Shipyard
		report.ResearchLab = buildingsStruct.ResearchLab
		report.MissileSilo = buildingsStruct.MissileSilo
		report.NaniteFactory = buildingsStruct.NaniteFactory
		report.Terraformer = buildingsStruct.Terraformer
		report.SpaceDock = buildingsStruct.SpaceDock
		report.LunarBase = buildingsStruct.LunarBase
		report.SensorPhalanx = buildingsStruct.SensorPhalanx
		report.JumpGate = buildingsStruct.JumpGate
	}
	report.HasResearchesInformation = rawMessageData.AttrOr("data-raw-hiddenresearch", "1") == ""
	if report.HasResearchesInformation {
		researchStr := rawMessageData.AttrOr("data-raw-research", "{}")
		var researchStruct struct {
			EspionageTechnology          *int64 `json:"106"`
			ComputerTechnology           *int64 `json:"108"`
			WeaponsTechnology            *int64 `json:"109"`
			ShieldingTechnology          *int64 `json:"110"`
			ArmourTechnology             *int64 `json:"111"`
			EnergyTechnology             *int64 `json:"113"`
			HyperspaceTechnology         *int64 `json:"114"`
			CombustionDrive              *int64 `json:"115"`
			ImpulseDrive                 *int64 `json:"117"`
			HyperspaceDrive              *int64 `json:"118"`
			LaserTechnology              *int64 `json:"120"`
			IonTechnology                *int64 `json:"121"`
			PlasmaTechnology             *int64 `json:"122"`
			IntergalacticResearchNetwork *int64 `json:"123"`
			Astrophysics                 *int64 `json:"124"`
			GravitonTechnology           *int64 `json:"199"`
		}
		_ = json.Unmarshal([]byte(researchStr), &researchStruct)
		report.EspionageTechnology = researchStruct.EspionageTechnology
		report.ComputerTechnology = researchStruct.ComputerTechnology
		report.WeaponsTechnology = researchStruct.WeaponsTechnology
		report.ShieldingTechnology = researchStruct.ShieldingTechnology
		report.ArmourTechnology = researchStruct.ArmourTechnology
		report.EnergyTechnology = researchStruct.EnergyTechnology
		report.HyperspaceTechnology = researchStruct.HyperspaceTechnology
		report.CombustionDrive = researchStruct.CombustionDrive
		report.ImpulseDrive = researchStruct.ImpulseDrive
		report.HyperspaceDrive = researchStruct.HyperspaceDrive
		report.LaserTechnology = researchStruct.LaserTechnology
		report.IonTechnology = researchStruct.IonTechnology
		report.PlasmaTechnology = researchStruct.PlasmaTechnology
		report.IntergalacticResearchNetwork = researchStruct.IntergalacticResearchNetwork
		report.Astrophysics = researchStruct.Astrophysics
		report.GravitonTechnology = researchStruct.GravitonTechnology
	}

	report.HasFleetInformation = rawMessageData.AttrOr("data-raw-hiddenships", "1") == ""
	if report.HasFleetInformation {
		fleetStr := rawMessageData.AttrOr("data-raw-fleet", "{}")
		var fleetStruct struct {
			SmallCargo     *int64 `json:"202"`
			LargeCargo     *int64 `json:"203"`
			LightFighter   *int64 `json:"204"`
			HeavyFighter   *int64 `json:"205"`
			Cruiser        *int64 `json:"206"`
			Battleship     *int64 `json:"207"`
			ColonyShip     *int64 `json:"208"`
			Recycler       *int64 `json:"209"`
			EspionageProbe *int64 `json:"210"`
			Bomber         *int64 `json:"211"`
			SolarSatellite *int64 `json:"212"`
			Destroyer      *int64 `json:"213"`
			Deathstar      *int64 `json:"214"`
			Battlecruiser  *int64 `json:"215"`
			Crawler        *int64 `json:"217"`
			Reaper         *int64 `json:"218"`
			Pathfinder     *int64 `json:"219"`
		}
		_ = json.Unmarshal([]byte(fleetStr), &fleetStruct)
		report.SmallCargo = fleetStruct.SmallCargo
		report.LargeCargo = fleetStruct.LargeCargo
		report.LightFighter = fleetStruct.LightFighter
		report.HeavyFighter = fleetStruct.HeavyFighter
		report.Cruiser = fleetStruct.Cruiser
		report.Battleship = fleetStruct.Battleship
		report.ColonyShip = fleetStruct.ColonyShip
		report.Recycler = fleetStruct.Recycler
		report.EspionageProbe = fleetStruct.EspionageProbe
		report.Bomber = fleetStruct.Bomber
		report.SolarSatellite = fleetStruct.SolarSatellite
		report.Destroyer = fleetStruct.Destroyer
		report.Deathstar = fleetStruct.Deathstar
		report.Battlecruiser = fleetStruct.Battlecruiser
		report.Crawler = fleetStruct.Crawler
		report.Reaper = fleetStruct.Reaper
		report.Pathfinder = fleetStruct.Pathfinder
	}

	report.HasDefensesInformation = rawMessageData.AttrOr("data-raw-hiddendef", "1") == ""
	if report.HasDefensesInformation {
		defStr := rawMessageData.AttrOr("data-raw-defense", "{}")
		var defStruct struct {
			RocketLauncher         *int64 `json:"401"`
			LightLaser             *int64 `json:"402"`
			HeavyLaser             *int64 `json:"403"`
			GaussCannon            *int64 `json:"404"`
			IonCannon              *int64 `json:"405"`
			PlasmaTurret           *int64 `json:"406"`
			SmallShieldDome        *int64 `json:"407"`
			LargeShieldDome        *int64 `json:"408"`
			AntiBallisticMissiles  *int64 `json:"502"`
			InterplanetaryMissiles *int64 `json:"503"`
		}
		_ = json.Unmarshal([]byte(defStr), &defStruct)
		report.RocketLauncher = defStruct.RocketLauncher
		report.LightLaser = defStruct.LightLaser
		report.HeavyLaser = defStruct.HeavyLaser
		report.GaussCannon = defStruct.GaussCannon
		report.IonCannon = defStruct.IonCannon
		report.PlasmaTurret = defStruct.PlasmaTurret
		report.SmallShieldDome = defStruct.SmallShieldDome
		report.LargeShieldDome = defStruct.LargeShieldDome
		report.AntiBallisticMissiles = defStruct.AntiBallisticMissiles
		report.InterplanetaryMissiles = defStruct.InterplanetaryMissiles
	}

	return report, nil
}

func extractExpeditionMessagesFromDoc(doc *goquery.Document, location *time.Location) ([]ogame.ExpeditionMessage, int64, error) {
	msgs := make([]ogame.ExpeditionMessage, 0)
	doc.Find(".msg").Each(func(i int, s *goquery.Selection) {
		if idStr, exists := s.Attr("data-msg-id"); exists {
			if id, err := utils.ParseI64(idStr); err == nil {
				msg := ogame.ExpeditionMessage{ID: id}
				msg.CreatedAt, _ = time.ParseInLocation("02.01.2006 15:04:05", s.Find(".msgDate").Text(), location)
				msg.Coordinate = v6.ExtractCoord(s.Find(".msgTitle a").Text())
				msg.Coordinate.Type = ogame.PlanetType
				msg.Content, _ = s.Find("div.msgContent").Html()
				msg.Content = strings.TrimSpace(msg.Content)

				var resStruct struct {
					Metal      int64 `json:"metal"`
					Crystal    int64 `json:"crystal"`
					Deuterium  int64 `json:"deuterium"`
					Darkmatter int64 `json:"darkMatter"`
				}
				resGained := s.Find("div.rawMessageData").AttrOr("data-raw-resourcesgained", "{}")
				_ = json.Unmarshal([]byte(resGained), &resStruct)

				msg.Resources.Metal = resStruct.Metal
				msg.Resources.Crystal = resStruct.Crystal
				msg.Resources.Deuterium = resStruct.Deuterium
				msg.Resources.Darkmatter = resStruct.Darkmatter

				type ShipInfo struct {
					Amount int64  `json:"amount"`
					Name   string `json:"name"`
				}
				var msgDataStruct struct {
					SmallCargo     ShipInfo `json:"202"`
					LargeCargo     ShipInfo `json:"203"`
					LightFighter   ShipInfo `json:"204"`
					HeavyFighter   ShipInfo `json:"205"`
					Cruiser        ShipInfo `json:"206"`
					Battleship     ShipInfo `json:"207"`
					EspionageProbe ShipInfo `json:"210"`
					Bomber         ShipInfo `json:"211"`
					Destroyer      ShipInfo `json:"213"`
					Battlecruiser  ShipInfo `json:"215"`
					Reaper         ShipInfo `json:"218"`
					Pathfinder     ShipInfo `json:"219"`
				}
				techGained := s.Find("div.rawMessageData").AttrOr("data-raw-technologiesgained", "{}")
				_ = json.Unmarshal([]byte(techGained), &msgDataStruct)
				msg.Ships.SmallCargo = msgDataStruct.SmallCargo.Amount
				msg.Ships.LargeCargo = msgDataStruct.LargeCargo.Amount
				msg.Ships.LightFighter = msgDataStruct.LightFighter.Amount
				msg.Ships.HeavyFighter = msgDataStruct.HeavyFighter.Amount
				msg.Ships.Cruiser = msgDataStruct.Cruiser.Amount
				msg.Ships.Battleship = msgDataStruct.Battleship.Amount
				msg.Ships.EspionageProbe = msgDataStruct.EspionageProbe.Amount
				msg.Ships.Bomber = msgDataStruct.Bomber.Amount
				msg.Ships.Destroyer = msgDataStruct.Destroyer.Amount
				msg.Ships.Battlecruiser = msgDataStruct.Battlecruiser.Amount
				msg.Ships.Reaper = msgDataStruct.Reaper.Amount
				msg.Ships.Pathfinder = msgDataStruct.Pathfinder.Amount

				msgs = append(msgs, msg)
			}
		}
	})
	return msgs, 1, nil
}

func extractLfBonusesFromDoc(doc *goquery.Document) (ogame.LfBonuses, error) {
	b := ogame.NewLfBonuses()
	doc.Find("bonus-item-content[data-toggable-target^=category]").Each(func(_ int, s *goquery.Selection) {
		category := s.AttrOr("data-toggable-target", "")
		if category == "categoryShips" || category == "categoryCostAndTime" {
			s.Find("inner-bonus-item-heading[data-toggable^=subcategory]").Each(func(_ int, g *goquery.Selection) {
				category, subcategory := extractCategories(g, category)
				if category != "" && subcategory != "" {
					assignBonusValue(g, b, category, subcategory)
				}
			})
		}
	})
	return *b, nil
}

func extractCategories(g *goquery.Selection, category string) (string, string) {
	c := strings.Replace(category, "category", "", 1)
	s, _ := g.Attr("data-toggable")
	v := "sub" + category
	return c, strings.Replace(s, v, "", 1)
}

func assignBonusValue(s *goquery.Selection, b *ogame.LfBonuses, category string, subcategory string) {
	switch category {
	case "Ships":
		extractShipStatBonus(s, b, subcategory)
	case "CostAndTime":
		extractCostReductionBonus(s, b, subcategory)
		extractTimeReductionBonus(s, b, subcategory)
	}
}

// Extracts cost reduction
func extractCostReductionBonus(s *goquery.Selection, l *ogame.LfBonuses, subcategory string) {
	i := utils.DoParseI64(subcategory)
	id := ogame.ID(i)
	s.Find("bonus-items").Each(func(_ int, s *goquery.Selection) {
		txt := s.Eq(0).Children().Eq(0).Contents().FilterFunction(func(i int, s *goquery.Selection) bool {
			return s.Nodes[0].Type == html.TextNode
		}).Text()
		costTimeBonus := l.CostTimeBonuses[id]
		costTimeBonus.Cost = extractBonusFromStringPercentage(txt)
		l.CostTimeBonuses[id] = costTimeBonus
	})
}

// Extracts time reduction
func extractTimeReductionBonus(s *goquery.Selection, l *ogame.LfBonuses, subcategory string) {
	i := utils.DoParseI64(subcategory)
	id := ogame.ID(i)
	s.Find("bonus-items").Each(func(_ int, s *goquery.Selection) {
		txt := s.Eq(0).Children().Eq(1).Contents().FilterFunction(func(i int, s *goquery.Selection) bool {
			return s.Nodes[0].Type == html.TextNode
		}).Text()
		costTimeBonus := l.CostTimeBonuses[id]
		costTimeBonus.Duration = extractBonusFromStringPercentage(txt)
		l.CostTimeBonuses[id] = costTimeBonus
	})
}

// Extracts ships stats fixed
func extractShipStatBonus(s *goquery.Selection, b *ogame.LfBonuses, subcategory string) {
	i := utils.DoParseI64(subcategory)
	id := ogame.ID(i)
	if !id.IsShip() {
		return
	}
	s.Find("bonus-items").Each(func(_ int, s *goquery.Selection) {
		extractFn := func(idx int) float64 {
			txt := s.Children().Eq(idx).Contents().FilterFunction(func(i int, s *goquery.Selection) bool {
				return s.Nodes[0].Type == html.TextNode
			}).Text()
			return extractBonusFromStringPercentage(txt)
		}
		shipBonus := ogame.LfShipBonus{
			ID:                  id,
			StructuralIntegrity: extractFn(0),
			ShieldPower:         extractFn(1),
			WeaponPower:         extractFn(2),
			Speed:               extractFn(3),
			CargoCapacity:       extractFn(4),
			FuelConsumption:     extractFn(5),
		}
		b.LfShipBonuses[id] = shipBonus
	})
	return
}

// Extract bonus value from a string with percentage sign [ex: 1.056% -> 0.01056]
func extractBonusFromStringPercentage(s string) float64 {
	v := strings.Replace(s, "%", "", 1)
	return extractBonusFromString(v) / 100.0
}

// Extract bonus value from a string [ex: 1.056]
func extractBonusFromString(s string) float64 {
	v := strings.TrimSpace(s)
	v = strings.Replace(v, ",", ".", 1)
	b, _ := strconv.ParseFloat(v, 64)
	return utils.RoundThousandth(b)
}

func extractAllianceClassFromDoc(doc *goquery.Document) (ogame.AllianceClass, error) {
	allianceClassTd := doc.Find("td.alliance_class").First()
	if allianceClassTd.HasClass("warrior") { // TODO: untested
		return ogame.Warrior, nil
	} else if allianceClassTd.HasClass("trader") {
		return ogame.Trader, nil
	} else if allianceClassTd.HasClass("explorer") {
		return ogame.Researcher, nil
	}
	return ogame.NoAllianceClass, errors.New("alliance class not found")
}

func extractPhalanxNewToken(pageHTML []byte) (string, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	token := doc.Find("a.refreshPhalanxLink").AttrOr("data-overlay-token", "")
	return token, nil
}
