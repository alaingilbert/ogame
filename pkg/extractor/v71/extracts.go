package v71

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	v6 "github.com/alaingilbert/ogame/pkg/extractor/v6"
	v7 "github.com/alaingilbert/ogame/pkg/extractor/v7"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/alaingilbert/clockwork"
	"golang.org/x/net/html"
)

type resourcesResp struct {
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
		Population struct {
			Amount  float64 `json:"amount"`
			Tooltip string  `json:"tooltip"`
		} `json:"population"`
		Food struct {
			Amount  float64 `json:"amount"`
			Tooltip string  `json:"tooltip"`
		} `json:"food"`
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

func extractResourcesDetails(pageHTML []byte) (out ogame.ResourcesDetails, err error) {
	var res resourcesResp
	if err = json.Unmarshal(pageHTML, &res); err != nil {
		if v6.IsLogged(pageHTML) {
			return out, ogame.ErrInvalidPlanetID
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
	out.Population.Available = int64(res.Resources.Population.Amount)
	out.Food.Available = int64(res.Resources.Food.Amount)
	metalDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(res.Resources.Metal.Tooltip))
	crystalDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(res.Resources.Crystal.Tooltip))
	deuteriumDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(res.Resources.Deuterium.Tooltip))
	darkmatterDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(res.Resources.Darkmatter.Tooltip))
	energyDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(res.Resources.Energy.Tooltip))
	populationDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(res.Resources.Population.Tooltip))
	foodDoc, _ := goquery.NewDocumentFromReader(strings.NewReader(res.Resources.Food.Tooltip))
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
	out.Food.Available = utils.ParseInt(foodDoc.Find("table tr").Eq(0).Find("td").Eq(0).Text())
	out.Food.StorageCapacity = utils.ParseInt(foodDoc.Find("table tr").Eq(1).Find("td").Eq(0).Text())
	out.Food.Overproduction = utils.ParseInt(foodDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	out.Food.ConsumedIn = utils.ParseInt(foodDoc.Find("table tr").Eq(3).Find("td").Eq(0).Text())
	out.Food.TimeTillFoodRunsOut = utils.ParseInt(foodDoc.Find("table tr").Eq(4).Find("td").Eq(0).Text())
	out.Population.Available = utils.ParseInt(populationDoc.Find("table tr").Eq(0).Find("td").Eq(0).Text())
	out.Population.T2Lifeforms = utils.ParseInt(populationDoc.Find("table tr").Eq(1).Find("td").Eq(0).Text())
	out.Population.T3Lifeforms = utils.ParseInt(populationDoc.Find("table tr").Eq(2).Find("td").Eq(0).Text())
	out.Population.LivingSpace = utils.ParseInt(populationDoc.Find("table tr").Eq(3).Find("td").Eq(0).Text())
	out.Population.Satisfied = utils.ParseInt(populationDoc.Find("table tr").Eq(4).Find("td").Eq(0).Text())
	out.Population.Hungry, _ = strconv.ParseFloat(populationDoc.Find("table tr").Eq(5).Find("td").Eq(0).Text(), 64)
	out.Population.GrowthRate, _ = strconv.ParseFloat(strings.TrimPrefix(populationDoc.Find("table tr").Eq(6).Find("td").Eq(0).Text(), "±"), 64)
	out.Population.BunkerSpace = utils.ParseInt(populationDoc.Find("table tr").Eq(7).Find("td").Eq(0).Text())
	return
}

type planetTechsResp struct {
	Num1   int64 `json:"1"`
	Num2   int64 `json:"2"`
	Num3   int64 `json:"3"`
	Num4   int64 `json:"4"`
	Num12  int64 `json:"12"`
	Num14  int64 `json:"14"`
	Num15  int64 `json:"15"`
	Num21  int64 `json:"21"`
	Num22  int64 `json:"22"`
	Num23  int64 `json:"23"`
	Num24  int64 `json:"24"`
	Num31  int64 `json:"31"`
	Num33  int64 `json:"33"`
	Num34  int64 `json:"34"`
	Num36  int64 `json:"36"`
	Num41  int64 `json:"41"`
	Num42  int64 `json:"42"`
	Num43  int64 `json:"43"`
	Num44  int64 `json:"44"`
	Num106 int64 `json:"106"`
	Num108 int64 `json:"108"`
	Num109 int64 `json:"109"`
	Num110 int64 `json:"110"`
	Num111 int64 `json:"111"`
	Num113 int64 `json:"113"`
	Num114 int64 `json:"114"`
	Num115 int64 `json:"115"`
	Num117 int64 `json:"117"`
	Num118 int64 `json:"118"`
	Num120 int64 `json:"120"`
	Num121 int64 `json:"121"`
	Num122 int64 `json:"122"`
	Num123 int64 `json:"123"`
	Num124 int64 `json:"124"`
	Num199 int64 `json:"199"`
	Num202 int64 `json:"202"`
	Num203 int64 `json:"203"`
	Num204 int64 `json:"204"`
	Num205 int64 `json:"205"`
	Num206 int64 `json:"206"`
	Num207 int64 `json:"207"`
	Num208 int64 `json:"208"`
	Num209 int64 `json:"209"`
	Num210 int64 `json:"210"`
	Num211 int64 `json:"211"`
	Num212 int64 `json:"212"`
	Num213 int64 `json:"213"`
	Num214 int64 `json:"214"`
	Num215 int64 `json:"215"`
	Num217 int64 `json:"217"`
	Num218 int64 `json:"218"`
	Num219 int64 `json:"219"`
	Num401 int64 `json:"401"`
	Num402 int64 `json:"402"`
	Num403 int64 `json:"403"`
	Num404 int64 `json:"404"`
	Num405 int64 `json:"405"`
	Num406 int64 `json:"406"`
	Num407 int64 `json:"407"`
	Num408 int64 `json:"408"`
	Num502 int64 `json:"502"`
	Num503 int64 `json:"503"`

	// LFbuildings
	Num11101 int64 `json:"11101"`
	Num11102 int64 `json:"11102"`
	Num11103 int64 `json:"11103"`
	Num11104 int64 `json:"11104"`
	Num11105 int64 `json:"11105"`
	Num11106 int64 `json:"11106"`
	Num11107 int64 `json:"11107"`
	Num11108 int64 `json:"11108"`
	Num11109 int64 `json:"11109"`
	Num11110 int64 `json:"11110"`
	Num11111 int64 `json:"11111"`
	Num11112 int64 `json:"11112"`
	Num12101 int64 `json:"12101"`
	Num12102 int64 `json:"12102"`
	Num12103 int64 `json:"12103"`
	Num12104 int64 `json:"12104"`
	Num12105 int64 `json:"12105"`
	Num12106 int64 `json:"12106"`
	Num12107 int64 `json:"12107"`
	Num12108 int64 `json:"12108"`
	Num12109 int64 `json:"12109"`
	Num12110 int64 `json:"12110"`
	Num12111 int64 `json:"12111"`
	Num12112 int64 `json:"12112"`
	Num13101 int64 `json:"13101"`
	Num13102 int64 `json:"13102"`
	Num13103 int64 `json:"13103"`
	Num13104 int64 `json:"13104"`
	Num13105 int64 `json:"13105"`
	Num13106 int64 `json:"13106"`
	Num13107 int64 `json:"13107"`
	Num13108 int64 `json:"13108"`
	Num13109 int64 `json:"13109"`
	Num13110 int64 `json:"13110"`
	Num13111 int64 `json:"13111"`
	Num13112 int64 `json:"13112"`
	Num14101 int64 `json:"14101"`
	Num14102 int64 `json:"14102"`
	Num14103 int64 `json:"14103"`
	Num14104 int64 `json:"14104"`
	Num14105 int64 `json:"14105"`
	Num14106 int64 `json:"14106"`
	Num14107 int64 `json:"14107"`
	Num14108 int64 `json:"14108"`
	Num14109 int64 `json:"14109"`
	Num14110 int64 `json:"14110"`
	Num14111 int64 `json:"14111"`
	Num14112 int64 `json:"14112"`
}

func extractTechs(pageHTML []byte) (supplies ogame.ResourcesBuildings, facilities ogame.Facilities, shipsInfos ogame.ShipsInfos, defenses ogame.DefensesInfos, researches ogame.Researches, lfBuildings ogame.LfBuildings, err error) {
	var res planetTechsResp
	if err = json.Unmarshal(pageHTML, &res); err != nil {
		if v6.IsLogged(pageHTML) {
			return supplies, facilities, shipsInfos, defenses, researches, lfBuildings, ogame.ErrInvalidPlanetID
		}
		return
	}
	supplies = ogame.ResourcesBuildings{
		MetalMine:            res.Num1,
		CrystalMine:          res.Num2,
		DeuteriumSynthesizer: res.Num3,
		SolarPlant:           res.Num4,
		FusionReactor:        res.Num12,
		SolarSatellite:       res.Num212,
		MetalStorage:         res.Num22,
		CrystalStorage:       res.Num23,
		DeuteriumTank:        res.Num24,
	}
	facilities = ogame.Facilities{
		RoboticsFactory: res.Num14,
		Shipyard:        res.Num21,
		ResearchLab:     res.Num31,
		AllianceDepot:   res.Num34,
		MissileSilo:     res.Num44,
		NaniteFactory:   res.Num15,
		Terraformer:     res.Num33,
		SpaceDock:       res.Num36,
		LunarBase:       res.Num41,
		SensorPhalanx:   res.Num42,
		JumpGate:        res.Num43,
	}
	shipsInfos = ogame.ShipsInfos{
		LightFighter:   res.Num204,
		HeavyFighter:   res.Num205,
		Cruiser:        res.Num206,
		Battleship:     res.Num207,
		Battlecruiser:  res.Num215,
		Bomber:         res.Num211,
		Destroyer:      res.Num213,
		Deathstar:      res.Num214,
		SmallCargo:     res.Num202,
		LargeCargo:     res.Num203,
		ColonyShip:     res.Num208,
		Recycler:       res.Num209,
		EspionageProbe: res.Num210,
		SolarSatellite: res.Num212,
		Crawler:        res.Num217,
		Reaper:         res.Num218,
		Pathfinder:     res.Num219,
	}
	defenses = ogame.DefensesInfos{
		RocketLauncher:         res.Num401,
		LightLaser:             res.Num402,
		HeavyLaser:             res.Num403,
		GaussCannon:            res.Num404,
		IonCannon:              res.Num405,
		PlasmaTurret:           res.Num406,
		SmallShieldDome:        res.Num407,
		LargeShieldDome:        res.Num408,
		AntiBallisticMissiles:  res.Num502,
		InterplanetaryMissiles: res.Num503,
	}
	researches = ogame.Researches{
		EnergyTechnology:             res.Num113,
		LaserTechnology:              res.Num120,
		IonTechnology:                res.Num121,
		HyperspaceTechnology:         res.Num114,
		PlasmaTechnology:             res.Num122,
		CombustionDrive:              res.Num115,
		ImpulseDrive:                 res.Num117,
		HyperspaceDrive:              res.Num118,
		EspionageTechnology:          res.Num106,
		ComputerTechnology:           res.Num108,
		Astrophysics:                 res.Num124,
		IntergalacticResearchNetwork: res.Num123,
		GravitonTechnology:           res.Num199,
		WeaponsTechnology:            res.Num109,
		ShieldingTechnology:          res.Num110,
		ArmourTechnology:             res.Num111,
	}
	lfBuildings = ogame.LfBuildings{
		ResidentialSector:          res.Num11101,
		BiosphereFarm:              res.Num11102,
		ResearchCentre:             res.Num11103,
		AcademyOfSciences:          res.Num11104,
		NeuroCalibrationCentre:     res.Num11105,
		HighEnergySmelting:         res.Num11106,
		FoodSilo:                   res.Num11107,
		FusionPoweredProduction:    res.Num11108,
		Skyscraper:                 res.Num11109,
		BiotechLab:                 res.Num11110,
		Metropolis:                 res.Num11111,
		PlanetaryShield:            res.Num11112,
		MeditationEnclave:          res.Num12101,
		CrystalFarm:                res.Num12102,
		RuneTechnologium:           res.Num12103,
		RuneForge:                  res.Num12104,
		Oriktorium:                 res.Num12105,
		MagmaForge:                 res.Num12106,
		DisruptionChamber:          res.Num12107,
		Megalith:                   res.Num12108,
		CrystalRefinery:            res.Num12109,
		DeuteriumSynthesiser:       res.Num12110,
		MineralResearchCentre:      res.Num12111,
		MetalRecyclingPlant:        res.Num12112,
		AssemblyLine:               res.Num13101,
		FusionCellFactory:          res.Num13102,
		RoboticsResearchCentre:     res.Num13103,
		UpdateNetwork:              res.Num13104,
		QuantumComputerCentre:      res.Num13105,
		AutomatisedAssemblyCentre:  res.Num13106,
		HighPerformanceTransformer: res.Num13107,
		MicrochipAssemblyLine:      res.Num13108,
		ProductionAssemblyHall:     res.Num13109,
		HighPerformanceSynthesiser: res.Num13110,
		ChipMassProduction:         res.Num13111,
		NanoRepairBots:             res.Num13112,
		Sanctuary:                  res.Num14101,
		AntimatterCondenser:        res.Num14102,
		VortexChamber:              res.Num14103,
		HallsOfRealisation:         res.Num14104,
		ForumOfTranscendence:       res.Num14105,
		AntimatterConvector:        res.Num14106,
		CloningLaboratory:          res.Num14107,
		ChrysalisAccelerator:       res.Num14108,
		BioModifier:                res.Num14109,
		PsionicModulator:           res.Num14110,
		ShipManufacturingHall:      res.Num14111,
		SupraRefractor:             res.Num14112,
	}
	return
}

// ar, Argentina           -> Recolector, General, Descubridor
// ba, Balkan              -> Sakupljač, General, Otkrivač
// br, Brasil              -> Coletor, General, Descobridor
// dk, Danmark             -> Samleren, Generalen, Opdageren
// de, Deutschland         -> Kollektor, General, Entdecker
// es, España              -> Recolector, General, Descubridor
// fr, France              -> Le collecteur, Général, L`explorateur
// hr, Hrvatska            -> Sakupljač, General, Otkrivač
// it, Italia              -> Collezionista, Generale, Esploratore
// hu, Magyarország        -> Gyűjtő, Tábornok, Felfedező
// mx, México              -> Recolector, General, Descubridor
// nl, Netherlands         -> Verzamelaar, Generaal, Ontdekker
// no, Norge               -> Collector, General, Discoverer (no i18n)
// pl, Polska              -> Zbieracz, Generał, Odkrywca
// pt, Portugal            -> Colecionador, General, Descobridor
// ro, Romania             -> Colecționarul, General, Exploratorul
// si, Slovenija           -> Zbiralec, Splošno, Odkritelj
// sk, Slovensko           -> Zberateľ, Generál, Objaviteľ
// fi, Suomi               -> Keräilijä, Komentaja, Löytäjä
// se, Sverige             -> Samlare, General, Upptäckare
// tr, Türkiye             -> Koleksiyoncu, General, Kaşif
// us, USA                 -> Collector, General, Discoverer
// en, United Kingdom      -> Collector, General, Discoverer
// cz, Če Republika        -> Sběratel, Generál, Průzkumník
// gr, Ελλάδα              -> Συλλέκτης, Στρατηγός, Εξερευνητής
// ru, Российс Федерация   -> Коллекционер, Генерал, Исследователь
// tw, 台灣                 -> 採礦師, 將軍, 探險家
// jp, 日本                 -> 回収船, 将軍, 探索船
func GetCharacterClass(characterClassStr string) ogame.CharacterClass {
	switch characterClassStr {
	case "Recolector", "Sakupljač", "Coletor", "Samleren", "Kollektor", "Le collecteur", "Collezionista", "Gyűjtő",
		"Verzamelaar", "Collector", "Zbieracz", "Colecionador", "Colecționarul", "Zbiralec", "Zberateľ", "Keräilijä",
		"Samlare", "Koleksiyoncu", "Sběratel", "Συλλέκτης", "Коллекционер", "採礦師", "回収船":
		return ogame.Collector
	case "General", "Generalen", "Général", "Generale", "Tábornok", "Generaal", "Generał", "Splošno", "Generál",
		"Komentaja", "Στρατηγός", "Генерал", "將軍", "将軍":
		return ogame.General
	case "Descubridor", "Otkrivač", "Descobridor", "Opdageren", "Entdecker", "L`explorateur", "Esploratore",
		"Felfedező", "Ontdekker", "Discoverer", "Odkrywca", "Exploratorul", "Odkritelj", "Objaviteľ", "Löytäjä",
		"Upptäckare", "Kaşif", "Průzkumník", "Εξερευνητής", "Исследователь", "探險家", "探索船":
		return ogame.Discoverer
	}
	return ogame.NoClass
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
	report.CharacterClass = GetCharacterClass(characterClassStr)

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
	resourcesFound := false
	doc.Find("ul.detail_list").Each(func(i int, s *goquery.Selection) {
		dataType := s.AttrOr("data-type", "")
		if dataType == "resources" && !resourcesFound {
			resourcesFound = true
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

func extractDestroyRocketsFromDoc(doc *goquery.Document) (abm, ipm int64, token string, err error) {
	scriptTxt := doc.Find("script").Text()
	r := regexp.MustCompile(`missileToken = "([^"]+)"`)
	m := r.FindStringSubmatch(scriptTxt)
	if len(m) != 2 {
		err = errors.New("failed to find missile token")
		return
	}
	token = m[1]
	abm = utils.DoParseI64(doc.Find("table tr").Eq(1).Find("td").Eq(1).Text())
	ipm = utils.DoParseI64(doc.Find("table tr").Eq(2).Find("td").Eq(1).Text())
	return
}

func extractIPMFromDoc(doc *goquery.Document) (duration, max int64, token string) {
	durationFloat, _ := strconv.ParseFloat(doc.Find("span#timer").AttrOr("data-duration", "0"), 64)
	duration = int64(math.Ceil(durationFloat))
	max = utils.DoParseI64(doc.Find("input#missileCount").AttrOr("data-max", "0"))
	token = doc.Find("input[name=token]").AttrOr("value", "")
	return
}

func extractFacilitiesFromDoc(doc *goquery.Document) (ogame.Facilities, error) {
	res, err := v7.ExtractFacilitiesFromDoc(doc)
	if err != nil {
		return ogame.Facilities{}, err
	}
	res.LunarBase = v7.GetNbr(doc, "moonbase")
	return res, nil
}

func extractCancelFleetTokenFromDoc(doc *goquery.Document, fleetID ogame.FleetID) (string, error) {
	href := doc.Find("div#fleet"+utils.FI64(fleetID)+" a.icon_link").AttrOr("href", "")
	m := regexp.MustCompile(`token=([^"]+)`).FindStringSubmatch(href)
	if len(m) != 2 {
		return "", errors.New("cancel fleet token not found")
	}
	token := m[1]
	return token, nil
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
	doc.Find("table.queue td").Each(func(i int, s *goquery.Selection) {
		link := s.Find("img")
		alt := link.AttrOr("alt", "")
		m := regexp.MustCompile(`techId_(\d+)`).FindStringSubmatch(alt)
		if len(m) == 0 {
			return
		}
		itemID := utils.DoParseI64(m[1])
		itemNbr := utils.ParseInt(s.Text())
		res = append(res, ogame.Quantifiable{ID: ogame.ID(itemID), Nbr: itemNbr})
	})
	return res, nil
}

func extractHighscoreFromDoc(doc *goquery.Document) (out ogame.Highscore, err error) {
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
	out.CurrPage = utils.DoParseI64(m[1])

	m = regexp.MustCompile(`var currentCategory = (\d+);`).FindStringSubmatch(script)
	if len(m) != 2 {
		return out, errors.New("failed to find currentCategory")
	}
	out.Category = utils.DoParseI64(m[1])

	m = regexp.MustCompile(`var currentType = (\d+);`).FindStringSubmatch(script)
	if len(m) != 2 {
		return out, errors.New("failed to find currentType")
	}
	out.Type = utils.DoParseI64(m[1])

	changeSiteSize := s.Find("select.changeSite option").Size()
	out.NbPage = utils.MaxInt(int64(changeSiteSize)-1, 0)

	s.Find("#ranks tbody tr").Each(func(i int, s *goquery.Selection) {
		p := ogame.HighscorePlayer{}
		p.Position = utils.DoParseI64(s.Find("td.position").Text())
		p.ID = utils.DoParseI64(s.Find("td.sendmsg a").AttrOr("data-playerid", "0"))
		p.Name = strings.TrimSpace(s.Find("span.playername").Text())
		tdName := s.Find("td.name")
		allyTag := tdName.Find("span.ally-tag")
		if allyTag != nil {
			href := allyTag.Find("a").AttrOr("href", "")
			m := regexp.MustCompile(`allianceId=(\d+)`).FindStringSubmatch(href)
			if len(m) == 2 {
				p.AllianceID = utils.DoParseI64(m[1])
			}
			allyTag.Remove()
		}
		href := tdName.Find("a").AttrOr("href", "")
		m := regexp.MustCompile(`galaxy=(\d+)&system=(\d+)&position=(\d+)`).FindStringSubmatch(href)
		if len(m) != 4 {
			return
		}
		p.Homeworld.Type = ogame.PlanetType
		p.Homeworld.Galaxy = utils.DoParseI64(m[1])
		p.Homeworld.System = utils.DoParseI64(m[2])
		p.Homeworld.Position = utils.DoParseI64(m[3])
		honorScoreSpan := s.Find("span.honorScore span")
		if honorScoreSpan == nil {
			return
		}
		p.HonourPoints = utils.ParseInt(strings.TrimSpace(honorScoreSpan.Text()))
		p.Score = utils.ParseInt(strings.TrimSpace(s.Find("td.score").Text()))
		shipsRgx := regexp.MustCompile(`([\d\.]+)`)
		shipsTitle := strings.TrimSpace(s.Find("td.score").AttrOr("title", "0"))
		shipsParts := shipsRgx.FindStringSubmatch(shipsTitle)
		if len(shipsParts) == 2 {
			p.Ships = utils.ParseInt(shipsParts[1])
		}
		out.Players = append(out.Players, p)
	})

	return
}

func extractAllResources(pageHTML []byte) (out map[ogame.CelestialID]ogame.Resources, err error) {
	out = make(map[ogame.CelestialID]ogame.Resources)
	m := regexp.MustCompile(`var planetResources\s?=\s?([^;]+);`).FindSubmatch(pageHTML)
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
		ki := utils.DoParseI64(k)
		out[ogame.CelestialID(ki)] = ogame.Resources{Metal: v.Input.Metal, Crystal: v.Input.Crystal, Deuterium: v.Input.Deuterium}
	}
	return
}

func extractAttacksFromDoc(doc *goquery.Document, clock clockwork.Clock, ownCoords []ogame.Coordinate) ([]ogame.AttackEvent, error) {
	attacks := make([]*ogame.AttackEvent, 0)
	out := make([]ogame.AttackEvent, 0)
	if doc.Find("body").Size() == 1 && v6.ExtractOGameSessionFromDoc(doc) != "" && doc.Find("div#eventListWrap").Size() == 0 {
		return out, ogame.ErrEventsBoxNotDisplayed
	} else if doc.Find("div#eventListWrap").Size() == 0 {
		return out, ogame.ErrNotLogged
	}

	allianceAttacks := make(map[int64]*ogame.AttackEvent)

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
					id = utils.DoParseI64(m[1])
				}
			} else {
				id = utils.DoParseI64(m[2])
			}

			classes, _ := s.Attr("class")
			partner := strings.Contains(classes, "partnerInfo")

			td := s.Find("td.countDown")
			isHostile := td.HasClass("hostile") || td.Find("span.hostile").Size() > 0
			isFriendly := td.HasClass("friendly") || td.Find("span.friendly").Size() > 0
			missionTypeInt := utils.DoParseI64(s.AttrOr("data-mission-type", ""))
			arrivalTimeInt := utils.DoParseI64(s.AttrOr("data-arrival-time", ""))
			missionType := ogame.MissionID(missionTypeInt)
			if rowType == "allianceAttack" {
				missionType = ogame.GroupedAttack
			}
			if missionType != ogame.Attack && missionType != ogame.GroupedAttack && missionType != ogame.Destroy &&
				missionType != ogame.MissileAttack && missionType != ogame.Spy {
				return
			}
			attack := &ogame.AttackEvent{}
			attack.ID = id
			attack.MissionType = missionType
			if missionType == ogame.Attack || missionType == ogame.MissileAttack || missionType == ogame.Spy || missionType == ogame.Destroy || missionType == ogame.GroupedAttack {
				linkSendMail := s.Find("a.sendMail")
				attack.AttackerID = utils.DoParseI64(linkSendMail.AttrOr("data-playerid", ""))
				attack.AttackerName = linkSendMail.AttrOr("title", "")
				if attack.AttackerID != 0 {
					coordsOrigin := strings.TrimSpace(s.Find("td.coordsOrigin").Text())
					attack.Origin = v6.ExtractCoord(coordsOrigin)
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
			attack.Destination = v6.ExtractCoord(destCoords)
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

				// People invite you to attack your own self, and ogame sees it as friendly.
				if isFriendly && attack.MissionType == ogame.GroupedAttack {
					found := false
					for _, ownCoord := range ownCoords {
						if attack.Destination.Equal(ownCoord) {
							found = true
							break
						}
					}
					if found {
						isHostile = true
					}
				}

				if !isHostile {
					return
				}
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

func extractDMCostsFromDoc(doc *goquery.Document) (ogame.DMCosts, error) {
	tmp := func(s *goquery.Selection) (id ogame.ID, nbr, cost int64, canBuy, isComplete bool, buyAndActivate string, token string) {
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
		nbr = utils.ParseInt(m[1])
		canBuy = !s.Find("span.dm_cost").HasClass("overmark")
		costTxt := s.Find("span.dm_cost").Text()
		m = r.FindStringSubmatch(costTxt)
		if len(m) != 2 {
			return
		}
		cost = utils.ParseInt(m[1])
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
	out := ogame.DMCosts{}
	buildingsBox := doc.Find("#productionboxbuildingcomponent")
	researchBox := doc.Find("#productionboxresearchcomponent")
	shipyardBox := doc.Find("#productionboxshipyardcomponent")
	out.Buildings.OGameID, out.Buildings.Nbr, out.Buildings.Cost, out.Buildings.CanBuy, out.Buildings.Complete, out.Buildings.BuyAndActivateToken, out.Buildings.Token = tmp(buildingsBox)
	out.Research.OGameID, out.Research.Nbr, out.Research.Cost, out.Research.CanBuy, out.Research.Complete, out.Research.BuyAndActivateToken, out.Research.Token = tmp(researchBox)
	out.Shipyard.OGameID, out.Shipyard.Nbr, out.Shipyard.Cost, out.Shipyard.CanBuy, out.Shipyard.Complete, out.Shipyard.BuyAndActivateToken, out.Shipyard.Token = tmp(shipyardBox)
	return out, nil
}

func extractBuffActivationFromDoc(doc *goquery.Document) (token string, items []ogame.Item, err error) {
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
	var inventoryMap map[string]ogame.Item
	if err = json.Unmarshal([]byte(m[1]), &inventoryMap); err != nil {
		fmt.Println(err)
		return
	}
	for _, item := range inventoryMap {
		items = append(items, item)
	}
	return
}

func extractActiveItemsFromDoc(doc *goquery.Document) (items []ogame.ActiveItem, err error) {
	doc.Find("ul.active_items div").Each(func(i int, s *goquery.Selection) {
		dataID := utils.ParseInt(s.AttrOr("data-id", ""))
		if dataID == 0 {
			return
		}
		durationDiv := s.Find("div.js_duration").First()
		aTitle := s.Find("a").AttrOr("title", "")
		imgSrc := s.Find("img").AttrOr("src", "")
		item := ogame.ActiveItem{}
		item.Ref = s.AttrOr("data-uuid", "")
		item.ID = dataID
		item.TotalDuration = utils.ParseInt(durationDiv.AttrOr("data-total-duration", ""))
		item.TimeRemaining = utils.ParseInt(durationDiv.Text())
		item.Name = strings.TrimSpace(strings.Split(aTitle, "|")[0])
		item.ImgSmall = imgSrc
		items = append(items, item)
	})
	return
}

func extractIsMobileFromDoc(doc *goquery.Document) bool {
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
