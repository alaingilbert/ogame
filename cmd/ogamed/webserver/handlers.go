package webserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/dustin/go-humanize"
	"github.com/faunX/ogame"
	"github.com/faunX/ogame/cmd/ogamed/ogb"
	"github.com/labstack/echo"
	"golang.org/x/net/html"
)

func EmpireHandler(c echo.Context) error {
	bot := c.Get("bot").(*ogame.OGame)
	database := c.Get("database").(*ogb.Ogb)
	cache, _ := json.Marshal(database)
	db := ogb.New()
	json.Unmarshal(cache, &db)

	db.Celestials = bot.GetCachedPlanets()

	var ResourcesIn map[ogame.CelestialID]ogame.Resources = map[ogame.CelestialID]ogame.Resources{}
	var ResourcesSum ogame.Resources
	delete(db.ResourcesDetails, 0)
	for k, v := range db.ResourcesDetails {
		secs := time.Now().Sub(db.Activities[k])
		ResourcesIn[k] = v.AvailableIn(secs.Seconds())
		ResourcesSum = ResourcesSum.Add(ResourcesIn[k])
	}
	ResourcesIn[0] = ResourcesSum

	var sumResBuildings ogame.ResourcesBuildings

	for k, _ := range db.ResourcesBuildings {
		sumResBuildings.MetalMine = sumResBuildings.MetalMine + db.ResourcesBuildings[k].MetalMine
		sumResBuildings.CrystalMine = sumResBuildings.CrystalMine + db.ResourcesBuildings[k].CrystalMine
		sumResBuildings.DeuteriumSynthesizer = sumResBuildings.DeuteriumSynthesizer + db.ResourcesBuildings[k].DeuteriumSynthesizer
		sumResBuildings.SolarPlant = sumResBuildings.SolarPlant + db.ResourcesBuildings[k].SolarPlant
		sumResBuildings.FusionReactor = sumResBuildings.FusionReactor + db.ResourcesBuildings[k].FusionReactor
		sumResBuildings.SolarSatellite = sumResBuildings.SolarSatellite + db.ResourcesBuildings[k].SolarSatellite
		sumResBuildings.MetalStorage = sumResBuildings.MetalStorage + db.ResourcesBuildings[k].MetalStorage
		sumResBuildings.CrystalStorage = sumResBuildings.CrystalStorage + db.ResourcesBuildings[k].CrystalStorage
		sumResBuildings.DeuteriumTank = sumResBuildings.DeuteriumTank + db.ResourcesBuildings[k].DeuteriumTank
	}
	db.ResourcesBuildings[0] = sumResBuildings

	var sumFacilities ogame.Facilities
	for k, _ := range db.Facilities {
		sumFacilities.RoboticsFactory = sumFacilities.RoboticsFactory + db.Facilities[k].RoboticsFactory
		sumFacilities.Shipyard = sumFacilities.Shipyard + db.Facilities[k].Shipyard
		sumFacilities.ResearchLab = sumFacilities.ResearchLab + db.Facilities[k].ResearchLab
		sumFacilities.AllianceDepot = sumFacilities.AllianceDepot + db.Facilities[k].AllianceDepot
		sumFacilities.MissileSilo = sumFacilities.MissileSilo + db.Facilities[k].MissileSilo
		sumFacilities.NaniteFactory = sumFacilities.NaniteFactory + db.Facilities[k].NaniteFactory
		sumFacilities.Terraformer = sumFacilities.Terraformer + db.Facilities[k].Terraformer
		sumFacilities.SpaceDock = sumFacilities.SpaceDock + db.Facilities[k].SpaceDock

		sumFacilities.LunarBase = sumFacilities.LunarBase + db.Facilities[k].LunarBase
		sumFacilities.SensorPhalanx = sumFacilities.SensorPhalanx + db.Facilities[k].SensorPhalanx
		sumFacilities.JumpGate = sumFacilities.JumpGate + db.Facilities[k].JumpGate
	}
	db.Facilities[0] = sumFacilities

	var sumDefenses ogame.DefensesInfos
	for k, _ := range db.DefensesInfos {
		sumDefenses.RocketLauncher = sumDefenses.RocketLauncher + db.DefensesInfos[k].RocketLauncher
		sumDefenses.LightLaser = sumDefenses.LightLaser + db.DefensesInfos[k].LightLaser
		sumDefenses.HeavyLaser = sumDefenses.HeavyLaser + db.DefensesInfos[k].HeavyLaser
		sumDefenses.GaussCannon = sumDefenses.GaussCannon + db.DefensesInfos[k].GaussCannon
		sumDefenses.IonCannon = sumDefenses.IonCannon + db.DefensesInfos[k].IonCannon
		sumDefenses.PlasmaTurret = sumDefenses.PlasmaTurret + db.DefensesInfos[k].PlasmaTurret
		sumDefenses.SmallShieldDome = sumDefenses.SmallShieldDome + db.DefensesInfos[k].SmallShieldDome
		sumDefenses.LargeShieldDome = sumDefenses.LargeShieldDome + db.DefensesInfos[k].LargeShieldDome
		sumDefenses.AntiBallisticMissiles = sumDefenses.AntiBallisticMissiles + db.DefensesInfos[k].AntiBallisticMissiles
		sumDefenses.InterplanetaryMissiles = sumDefenses.InterplanetaryMissiles + db.DefensesInfos[k].InterplanetaryMissiles
	}
	db.DefensesInfos[0] = sumDefenses

	var sumShips ogame.ShipsInfos
	for k, _ := range db.ShipsInfos {
		sumShips.SmallCargo = sumShips.SmallCargo + db.ShipsInfos[k].SmallCargo
		sumShips.LargeCargo = sumShips.LargeCargo + db.ShipsInfos[k].LargeCargo
		sumShips.LightFighter = sumShips.LightFighter + db.ShipsInfos[k].LightFighter
		sumShips.HeavyFighter = sumShips.HeavyFighter + db.ShipsInfos[k].HeavyFighter
		sumShips.Cruiser = sumShips.Cruiser + db.ShipsInfos[k].Cruiser
		sumShips.Battleship = sumShips.Battleship + db.ShipsInfos[k].Battleship
		sumShips.ColonyShip = sumShips.ColonyShip + db.ShipsInfos[k].ColonyShip
		sumShips.Recycler = sumShips.Recycler + db.ShipsInfos[k].Recycler
		sumShips.EspionageProbe = sumShips.EspionageProbe + db.ShipsInfos[k].EspionageProbe
		sumShips.Bomber = sumShips.Bomber + db.ShipsInfos[k].Bomber
		sumShips.SolarSatellite = sumShips.SolarSatellite + db.ShipsInfos[k].SolarSatellite
		sumShips.Destroyer = sumShips.Destroyer + db.ShipsInfos[k].Destroyer
		sumShips.Deathstar = sumShips.Deathstar + db.ShipsInfos[k].Deathstar
		sumShips.Battlecruiser = sumShips.Battlecruiser + db.ShipsInfos[k].Battlecruiser
		sumShips.Crawler = sumShips.Crawler + db.ShipsInfos[k].Crawler
		sumShips.Reaper = sumShips.Reaper + db.ShipsInfos[k].Reaper
		sumShips.Pathfinder = sumShips.Pathfinder + db.ShipsInfos[k].Pathfinder
	}
	db.ShipsInfos[0] = sumShips

	obj := struct {
		Bot             *ogame.OGame
		DB              *ogb.Ogb
		ObjsStruct      ogame.ObjsStruct
		PlanetBuildings []ogame.Building
		MoonBuildings   []ogame.Building
		Buildings       []ogame.Building
		Ships           []ogame.Ship
		Defenses        []ogame.Defense
		Technologies    []ogame.Technology
		ResourcesIn     map[ogame.CelestialID]ogame.Resources
	}{
		bot,
		db,
		ogame.Objs,
		ogame.PlanetBuildings,
		ogame.MoonBuildings,
		ogame.Buildings,
		ogame.Ships,
		ogame.Defenses,
		ogame.Technologies,
		ResourcesIn,
	}

	return c.Render(http.StatusOK, "empire", obj)
}

// GetAlliancePageContentHandler ...
func GetAlliancePageContentHandler2(c echo.Context) error {
	bot := c.Get("bot").(*ogame.OGame)
	allianceID := c.QueryParam("allianceId")
	vals := url.Values{"allianceId": {allianceID}}

	var pageHTML []byte
	locked, state := bot.GetState()
	if locked && state == "Manuel Mode" {
		pageHTML, _ = tx.GetAlliancePageContent(vals)
	} else {
		pageHTML, _ = bot.GetAlliancePageContent(vals)
	}

	return c.HTML(http.StatusOK, string(pageHTML))
}

var lastActiveCelestialID ogame.CelestialID
var lastActiveCelestialIDMu sync.RWMutex

// GetFromGameHandler ...
func GetFromGameHandler2(c echo.Context) error {
	bot := c.Get("bot").(*ogame.OGame)
	ogbI := c.Get("database").(*ogb.Ogb)

	db := ogbI.GetDatabase()

	vals := url.Values{"page": {"ingame"}, "component": {"overview"}}
	if len(c.QueryParams()) > 0 {
		vals = c.QueryParams()
	}

	//Build
	//page=ingame&component=supplies&modus=1&type=1&menge=1&token=f31ec93db6bb080c5669f2f19978b7ab

	//Cancel
	//page=ingame&component=supplies&modus=2&token=a402d9355937abf4c7e60f71af6006be&action=cancel&type=1&listid=13311770

	// Destroy
	//page=ingame&component=supplies&modus=3&token=93d712277a9f5ec3f18537e04db4a395&type=2
	if vals.Get("page") == "ingame" && (vals.Get("component") == "supplies" || vals.Get("component") == "facilities" || vals.Get("component") == "research") &&
		vals.Get("modus") == "1" {
		tmpID, _ := strconv.ParseInt(vals.Get("type"), 10, 64)
		id := ogame.ID(tmpID)
		nbr, _ := strconv.ParseInt(vals.Get("menge"), 10, 64)

		var b ogb.BrainQueueType
		b.CelestialID = ogbI.LastActiveCelestialID
		b.ID = id
		b.Nbr = nbr
		//ogbI.Database.Lock()
		//ogbI.Database.AddToBrainQueue(b)
		//ogbI.Database.Unlock()

		// Ensure Building
		/*
			Not Working in Manual Mode
			err := bot.Build(b.CelestialID, id, nbr)
			if err == nil {
				vals.Del("modus")
				vals.Del("type")
				vals.Del("menge")
			}
		*/

		//data, _ := json.Marshal(ogbI.Database.BrainQueue)
		//return c.HTMLBlob(http.StatusOK, data)
	}

	var pageHTML []byte
	locked, state := bot.GetState()
	if locked && state == "Manuel Mode" {
		var err error
		pageHTML, err = tx.GetPageContent(vals)
		if err != nil {
			log.Println(err)
		}
	} else {
		pageHTML, _ = bot.GetPageContent(vals)
	}
	pageHTML = ogame.ReplaceHostname(bot, pageHTML, c.Request())

	if ogame.IsKnowFullPage(vals) {
		pageHTML = ogame.HTMLCleaner(bot, c.Request().Method, c.Request().URL.String(), c.QueryParams(), nil, pageHTML)

		doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))

		ok := false
		doc.Find("meta").Each(func(i int, q *goquery.Selection) {
			v := q.AttrOr("name", "")
			if v == "ogame-session" {

			}
			k := q.AttrOr("content", "")
			if k != "" {
				ok = true
			}
		})

		if ok {
			stateON := "off"
			if state == "Manuel Mode" {
				stateON = "on"
			}

			button := `
				<div style="position: fixed; right: 10px; top: 10px; z-index: 3001;">
					<button id="manuelModeBtn" name="manuelModeBtn" style="height: 20px; width: 130px; display: block; margin-bottom: 3px;">Manual mode (` + stateON + `)</button>
				</div>`

			script := `<script>
	
			manuelModeBtn.onclick = function() {
				var formData = new FormData();
				//formData.append('csrf', 'POBpUTkNsaDiuEQ9oiPdnncsVsw9ginl');
				$.ajax({
					url: "/toggle-manual-mode", data: formData, type: 'POST', processData: false, contentType: false,
					success: function(res) {
						$(manuelModeBtn).text("Manual mode (" + (res ? "on" : "off") + ")");
					},
					error: function(req) { console.log(req.responseText); },
				});
			};
	
			function inIframe () {
				try {
					return window.self !== window.top;
				} catch (e) {
					return true;
				}
			}
			</script>`

			doc.Find("body").AppendHtml(button)
			doc.Find("body").AppendHtml(script)
			txtHTML, _ := doc.Html()
			pageHTML = []byte(txtHTML)
		}
	}

	if vals.Get("page") == "ingame" && vals.Get("component") == "overview" {
		doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
		planetDetails, err := doc.Find("#planetDetails table tbody").Html()
		if err != nil {
			return c.HTMLBlob(http.StatusOK, pageHTML)
		}
		e := bot.GetExtractor()
		//celestialID, _ := e.ExtractPlanetID(pageHTML)
		resDetails := e.ExtractResourcesDetailsFromFullPage(pageHTML)

		storageCapacity := ogame.Resources{
			Metal:     resDetails.Metal.StorageCapacity,
			Crystal:   resDetails.Crystal.StorageCapacity,
			Deuterium: resDetails.Deuterium.StorageCapacity,
		}

		//res := db.GetResources(celestialID)

		storageTimes := GetStorageTime(storageCapacity, resDetails)

		var metalProductionColor = "undermark"
		if resDetails.Metal.CurrentProduction == 0 {
			metalProductionColor = "overmark"
		}
		var metalStorageColor = ""
		if resDetails.Metal.StorageCapacity <= resDetails.Metal.Available {
			metalStorageColor = "overmark"
		}
		metalFill := float64(resDetails.Metal.Available) / float64(resDetails.Metal.StorageCapacity)
		if metalFill >= 0.9 && metalFill < 1 {
			metalStorageColor = "middlemark"
		}
		planetStorageDetails := ""
		planetStorageDetails = `<table width="100%" cellspacing="0" cellpadding="0"><tbody>
		<tr>
		<td class="desc">Metal:</td>
		<td class="data"><span class="` + metalProductionColor + `">(+` + humanize.Comma(resDetails.Metal.CurrentProduction) + `) </span><span class="` + metalStorageColor + `">` + humanize.Comma(resDetails.Metal.Available) + ` / ` + humanize.Comma(resDetails.Metal.StorageCapacity) + `</span></td></tr>
		<tr><td class="desc"></td><td class="data"><span class="notago_overview_storagetime" style="font-size: 10px;color: #7cff0f;vertical-align: top;">` + time.Now().Add(storageTimes.Metal).Format("02.01. 15:04:05") + `</span></td></tr>`

		var crystalProductionColor = "undermark"
		if resDetails.Crystal.CurrentProduction == 0 {
			crystalProductionColor = "overmark"
		}
		var crystalStorageColor = ""
		if resDetails.Crystal.StorageCapacity <= resDetails.Crystal.Available {
			crystalStorageColor = "overmark"
		}
		crystalFill := float64(resDetails.Crystal.Available) / float64(resDetails.Crystal.StorageCapacity) * 100
		if crystalFill >= 0.9 && crystalFill < 1 {
			crystalStorageColor = "middlemark"
		}
		planetStorageDetails += `<tr>
		<td class="desc">Crystal:</td>
		<td class="data"><span class=" ` + crystalProductionColor + `">(+` + humanize.Comma(resDetails.Crystal.CurrentProduction) + `) </span><span class="` + crystalStorageColor + `">` + humanize.Comma(resDetails.Crystal.Available) + ` / ` + humanize.Comma(resDetails.Crystal.StorageCapacity) + `</td></tr>
		<tr><td class="desc"></td><td class="data"><span class="notago_overview_storagetime" style="font-size: 10px;color: #7cff0f;vertical-align: top;">` + time.Now().Add(storageTimes.Crystal).Format("02.01. 15:04:05") + `</span></td></tr>`

		var deuteriumProductionColor = "undermark"
		if resDetails.Deuterium.CurrentProduction == 0 {
			deuteriumProductionColor = "overmark"
		}
		var deuteriumStorageColor = ""
		if resDetails.Deuterium.StorageCapacity <= resDetails.Deuterium.Available {
			deuteriumStorageColor = "overmark"
		}
		deuteriumFill := float64(resDetails.Deuterium.Available) / float64(resDetails.Deuterium.StorageCapacity)
		if deuteriumFill >= 0.9 && deuteriumFill < 1 {
			deuteriumStorageColor = "middlemark"
		}
		planetStorageDetails += `<tr>
		<td class="desc">Deuterium:</td>
		<td class="data"><span class="` + deuteriumProductionColor + `">(+` + humanize.Comma(resDetails.Deuterium.CurrentProduction) + `) </span><span class="` + deuteriumStorageColor + `">` + humanize.Comma(resDetails.Deuterium.Available) + ` / ` + humanize.Comma(resDetails.Deuterium.StorageCapacity) + `</td></tr>
		<tr><td class="desc"></td><td class="data"><span class="notago_overview_storagetime" style="font-size: 10px;color: #7cff0f;vertical-align: top;">` + time.Now().Add(storageTimes.Deuterium).Format("02.01. 15:04:05") + `</span></td></tr>`

		planetDetails = planetDetails + "</tbody></table>"
		doc.Find("#planetDetails table tbody").SetHtml(planetStorageDetails + planetDetails)
		doc.Find("#planetDetails").SetAttr("style", "height: 194px")
		doc.Find("#planetdata").SetAttr("style", "height: 217px;margin-top: 6px")

		pageString, _ := doc.Html()
		pageHTML = []byte(pageString)

	}

	if vals.Get("page") == "ingame" && vals.Get("component") == "technologydetails" && vals.Get("ajax") == "1" {
		lastPlanet := bot.GetCachedCelestialByID(db.LastActiveCelestialID)
		if lastPlanet.GetCoordinate().IsPlanet() {
			id, _ := strconv.ParseInt(vals.Get("technology"), 10, 64)
			obj := ogame.Objs.ByID(ogame.ID(id))
			obj.GetID()

			technologydetails := struct {
				Target  string `json:"target"`
				Content struct {
					Technologydetails string `json:"technologydetails"`
				} `json:"content"`
				Files struct {
					Js  []string `json:"js"`
					Css []string `json:"css"`
				} `json:"files"`
				Page struct {
					StateObj interface{} `json:"stateObj"`
					Title    string      `json:"title"`
					Url      string      `json:"url"`
				} `json:"page"`
				ServerTime   int64  `json:"serverTime"`
				NewAjaxToken string `json:"newAjaxToken"`
			}{}
			json.Unmarshal(pageHTML, &technologydetails)
			// Ships Start
			lastActiveCelestialIDMu.RLock()
			//res, _ := bot.GetResourcesDetails(lastActiveCelestialID)
			res := db.ResourcesDetails[db.LastActiveCelestialID]
			lastActiveCelestialIDMu.RUnlock()

			if obj.GetID().IsShip() || obj.GetID().IsDefense() {
				s := strings.ReplaceAll(``+technologydetails.Content.Technologydetails+``, "\\n", "")
				s = strings.ReplaceAll(``+s+``, "\\", "")

				node, _ := html.Parse(bytes.NewReader([]byte(s)))
				doc := goquery.NewDocumentFromNode(node)

				max := res.Available().Div(obj.GetPrice(1))
				doc.Find("div.build_amount input").SetAttr("min", "0")
				doc.Find("div.build_amount input").SetAttr("max", strconv.FormatInt(max, 10))
				doc.Find("div.build_amount input").SetAttr("onfocus", `clearInput(this);"", "0"`)
				doc.Find("div.build_amount input").SetAttr("onkeyup", `checkIntInput(this, 1, `+strconv.FormatInt(max, 10)+`);event.stopPropagation();`)
				doc.Find("div.build_amount").AppendHtml("<button class=\"maximum\">[max. " + strconv.FormatInt(max, 10) + "]</button>")

				var err error
				technologydetails.Content.Technologydetails, err = doc.Html()
				if err != nil {
					log.Printf("Error occured %s", err.Error())
				}
				pageHTML, _ = json.Marshal(technologydetails)
			}
			// Ships End

			doc, _ := goquery.NewDocumentFromReader(bytes.NewReader([]byte(technologydetails.Content.Technologydetails)))
			level, _ := strconv.ParseInt(doc.Find("span.level").AttrOr("data-value", "0"), 10, 64)
			if level == 0 {
				level++
			}
			// Add
			resTime := GetResourceTime(obj.GetPrice(level), db.ResourcesDetails[db.LastActiveCelestialID])
			if resTime.Time.Seconds() > 0 {
				doc.Find("div.information ul").AppendHtml("<li><strong> Missing Resources: " + resTime.Name + `</strong> <time class="value tooltip" datetime="` + resTime.Time.String() + `" title="">` + resTime.Time.String() + `</time></li>`)
			} else {
				doc.Find("div.information ul").AppendHtml(`<li><strong>Can be built: </strong><time class="value tooltip" datetime="now" title="">now</time></li>`)
			}
			doc.Find("input#build_amount").SetAttr("onfocus", "clearInput(this);")

			//doc.Find("script").Remove()

			doc.Find("script").SetHtml(`if(document.getElementById("build_amount") !== null) {
				document.getElementById("build_amount").focus();
			}
			var lastBuildingSlot = {"showWarning":false, "slotWarning":""};`)
			//`var lastBuildingSlot = {"showWarning":false,"slotWarning":"Dieses Gebu00e4ude wird den letzten Bauplatz verbrauchen. Baue deinen Terraformer aus oder kaufe ein Planetenfelditem (z.B. <a href="http://10.156.176.2:8080/game/index.php?page=shop#page=shop&amp;category=dc9ec90e5a2163cc063b8bb3e9fe392782f565c8&amp;item=04e58444d6d0beb57b3e998edc34c60f8318825a" target="_parent" title="Planetenfelder Gold|+15 zusu00e4tzliche Felder auf einem Planeten&lt;br /&gt; &lt;br /&gt; Laufzeit: Permanent&lt;br /&gt; &lt;br /&gt; Preis: 300.000 Dunkle Materie&lt;br /&gt; Im Inventar: 0" class="tooltipHTML itemLink">Planetenfelder Gold</a>), um weitere Plu00e4tze zu bekommen.<br />Mu00f6chtest du das Gebu00e4ude wirklich bauen?"};`

			technologydetails.Content.Technologydetails, _ = doc.Html()

			pageHTML, _ = json.MarshalIndent(technologydetails, "", "  ")
		}
	}
	return c.HTMLBlob(http.StatusOK, pageHTML)
}

// PostToGameHandler ...
func PostToGameHandler2(c echo.Context) error {
	bot := c.Get("bot").(*ogame.OGame)
	vals := url.Values{"page": {"ingame"}, "component": {"overview"}}
	if len(c.QueryParams()) > 0 {
		vals = c.QueryParams()
	}
	payload, _ := c.FormParams()

	var pageHTML []byte
	locked, state := bot.GetState()
	if locked && state == "Manuel Mode" {
		pageHTML, _ = tx.PostPageContent(vals, payload)
	} else {
		pageHTML, _ = bot.PostPageContent(vals, payload)
	}
	pageHTML = ogame.ReplaceHostname(bot, pageHTML, c.Request())
	return c.HTMLBlob(http.StatusOK, pageHTML)
}

func galaxyExtender(b *ogame.OGame, pageHTML []byte) []byte {

	var galaxy interface{}
	json.Unmarshal(pageHTML, &galaxy)

	return pageHTML
}

type resourceTime struct {
	Name string
	Time time.Duration
}

func GetResourceTime(p ogame.Resources, d ogame.ResourcesDetails) resourceTime {
	var metal, crystal, deuterium float64
	neededResources := p.Sub(d.Available())
	if neededResources.Metal != 0 {
		metal = math.Ceil(float64(neededResources.Metal) / (float64(d.Metal.CurrentProduction)) * 3600)
	}
	if neededResources.Crystal != 0 {
		crystal = math.Ceil(float64(neededResources.Crystal) / (float64(d.Crystal.CurrentProduction)) * 3600)
	}
	if neededResources.Deuterium != 0 {
		deuterium = math.Ceil(float64(neededResources.Deuterium) / (float64(d.Deuterium.CurrentProduction)) * 3600)
	}
	/*
		return ogame.Resources{
			Metal:     int64(metal),
			Crystal:   int64(crystal),
			Deuterium: int64(deuterium),
		}
	*/

	max := ogame.MaxInt(int64(metal), int64(crystal), int64(deuterium))
	name := ""
	switch float64(max) {
	case metal:
		if metal > 0 {
			name = "metal"
		}
		break
	case crystal:
		if crystal > 0 {
			name = "crystal"
		}
		break
	case deuterium:
		if deuterium > 0 {
			name = "deuterium"
		}
		break
	}

	return resourceTime{
		Name: name,
		Time: time.Duration(max) * time.Second,
	}
}

type storageTimes struct {
	Metal     time.Duration
	Crystal   time.Duration
	Deuterium time.Duration
}

func GetStorageTime(p ogame.Resources, d ogame.ResourcesDetails) storageTimes {
	var metal, crystal, deuterium float64
	neededResources := p.Sub(d.Available())
	if neededResources.Metal != 0 {
		metal = math.Ceil(float64(neededResources.Metal) / (float64(d.Metal.CurrentProduction)) * 3600)
	}
	if neededResources.Crystal != 0 {
		crystal = math.Ceil(float64(neededResources.Crystal) / (float64(d.Crystal.CurrentProduction)) * 3600)
	}
	if neededResources.Deuterium != 0 {
		deuterium = math.Ceil(float64(neededResources.Deuterium) / (float64(d.Deuterium.CurrentProduction)) * 3600)
	}

	return storageTimes{
		Metal:     time.Duration(int64(metal)) * time.Second,
		Crystal:   time.Duration(int64(crystal)) * time.Second,
		Deuterium: time.Duration(int64(deuterium)) * time.Second,
	}
}

func FormatSince(t time.Time) string {
	const (
		Decisecond = 100 * time.Millisecond
		Day        = 24 * time.Hour
	)
	ts := time.Since(t)
	sign := time.Duration(1)
	if ts < 0 {
		sign = -1
		ts = -ts
	}
	ts += +Decisecond / 2
	d := sign * (ts / Day)
	ts = ts % Day
	h := ts / time.Hour
	ts = ts % time.Hour
	m := ts / time.Minute
	ts = ts % time.Minute
	s := ts / time.Second
	ts = ts % time.Second
	f := ts / Decisecond
	return fmt.Sprintf("%dd%dh%dm%d.%ds", d, h, m, s, f)
}

type TBotSettings struct {
	Credentials struct {
		Universe      string `json:"Universe"`
		Email         string `json:"Email"`
		Password      string `json:"Password"`
		Language      string `json:"Language"`
		LobbyPioneers bool   `json:"LobbyPioneers"`
		BasicAuth     struct {
			Username string `json:"Username"`
			Password string `json:"Password"`
		} `json:"BasicAuth"`
	} `json:"Credentials"`
	General struct {
		UserAgent      string `json:"UserAgent"`
		Host           string `json:"Host"`
		Port           string `json:"Port"`
		NewAPIHostname string `json:"NewApiHostname"`
		Proxy          struct {
			Enabled   bool   `json:"Enabled"`
			Address   string `json:"Address"`
			Type      string `json:"Type"`
			Username  string `json:"Username"`
			Password  string `json:"Password"`
			LoginOnly bool   `json:"LoginOnly"`
		} `json:"Proxy"`
		CaptchaAPIKey    string `json:"CaptchaAPIKey"`
		CustomTitle      string `json:"CustomTitle"`
		SlotsToLeaveFree int    `json:"SlotsToLeaveFree"`
	} `json:"General"`
	SleepMode struct {
		Active                  bool   `json:"Active"`
		GoToSleep               string `json:"GoToSleep"`
		WakeUp                  string `json:"WakeUp"`
		PreventIfThereAreFleets bool   `json:"PreventIfThereAreFleets"`
		TelegramMessenger       struct {
			Active bool `json:"Active"`
		} `json:"TelegramMessenger"`
		AutoFleetSave struct {
			Active      bool `json:"Active"`
			OnlyMoons   bool `json:"OnlyMoons"`
			ForceUnsafe bool `json:"ForceUnsafe"`
			DeutToLeave int  `json:"DeutToLeave"`
			Recall      bool `json:"Recall"`
		} `json:"AutoFleetSave"`
	} `json:"SleepMode"`
	Defender struct {
		Active           bool `json:"Active"`
		CheckIntervalMin int  `json:"CheckIntervalMin"`
		CheckIntervalMax int  `json:"CheckIntervalMax"`
		IgnoreProbes     bool `json:"IgnoreProbes"`
		IgnoreWeakAttack bool `json:"IgnoreWeakAttack"`
		WeakAttackRatio  int  `json:"WeakAttackRatio"`
		Autofleet        struct {
			Active            bool `json:"Active"`
			TelegramMessenger struct {
				Active bool `json:"Active"`
			} `json:"TelegramMessenger"`
		} `json:"Autofleet"`
		WhiteList   []int `json:"WhiteList"`
		SpyAttacker struct {
			Active bool `json:"Active"`
			Probes int  `json:"Probes"`
		} `json:"SpyAttacker"`
		MessageAttacker struct {
			Active   bool     `json:"Active"`
			Messages []string `json:"Messages"`
		} `json:"MessageAttacker"`
		TelegramMessenger struct {
			Active bool `json:"Active"`
		} `json:"TelegramMessenger"`
		Alarm struct {
			Active bool `json:"Active"`
		} `json:"Alarm"`
	} `json:"Defender"`
	Brain struct {
		Active   bool `json:"Active"`
		AutoMine struct {
			Active                  bool `json:"Active"`
			MaxMetalMine            int  `json:"MaxMetalMine"`
			MaxCrystalMine          int  `json:"MaxCrystalMine"`
			MaxDeuteriumSynthetizer int  `json:"MaxDeuteriumSynthetizer"`
			MaxSolarPlant           int  `json:"MaxSolarPlant"`
			MaxFusionReactor        int  `json:"MaxFusionReactor"`
			DepositHours            int  `json:"DepositHours"`
			MaxMetalStorage         int  `json:"MaxMetalStorage"`
			MaxCrystalStorage       int  `json:"MaxCrystalStorage"`
			MaxDeuteriumTank        int  `json:"MaxDeuteriumTank"`
			MaxRoboticsFactory      int  `json:"MaxRoboticsFactory"`
			MaxShipyard             int  `json:"MaxShipyard"`
			MaxResearchLab          int  `json:"MaxResearchLab"`
			MaxMissileSilo          int  `json:"MaxMissileSilo"`
			MaxNaniteFactory        int  `json:"MaxNaniteFactory"`
			MaxLunarBase            int  `json:"MaxLunarBase"`
			MaxLunarShipyard        int  `json:"MaxLunarShipyard"`
			MaxLunarRoboticsFactory int  `json:"MaxLunarRoboticsFactory"`
			MaxSensorPhalanx        int  `json:"MaxSensorPhalanx"`
			MaxJumpGate             int  `json:"MaxJumpGate"`
			Trasports               struct {
				Active      bool   `json:"Active"`
				CargoType   string `json:"CargoType"`
				DeutToLeave int    `json:"DeutToLeave"`
				Origin      struct {
					Galaxy   int    `json:"Galaxy"`
					System   int    `json:"System"`
					Position int    `json:"Position"`
					Type     string `json:"Type"`
				} `json:"Origin"`
			} `json:"Trasports"`
			RandomOrder                            bool          `json:"RandomOrder"`
			Exclude                                []interface{} `json:"Exclude"`
			PrioritizeRobotsAndNanitesOnNewPlanets bool          `json:"PrioritizeRobotsAndNanitesOnNewPlanets"`
			CheckIntervalMin                       int           `json:"CheckIntervalMin"`
			CheckIntervalMax                       int           `json:"CheckIntervalMax"`
		} `json:"AutoMine"`
		AutoResearch struct {
			Active                          bool `json:"Active"`
			MaxEnergyTechnology             int  `json:"MaxEnergyTechnology"`
			MaxLaserTechnology              int  `json:"MaxLaserTechnology"`
			MaxIonTechnology                int  `json:"MaxIonTechnology"`
			MaxHyperspaceTechnology         int  `json:"MaxHyperspaceTechnology"`
			MaxPlasmaTechnology             int  `json:"MaxPlasmaTechnology"`
			MaxCombustionDrive              int  `json:"MaxCombustionDrive"`
			MaxImpulseDrive                 int  `json:"MaxImpulseDrive"`
			MaxHyperspaceDrive              int  `json:"MaxHyperspaceDrive"`
			MaxEspionageTechnology          int  `json:"MaxEspionageTechnology"`
			MaxComputerTechnology           int  `json:"MaxComputerTechnology"`
			MaxAstrophysics                 int  `json:"MaxAstrophysics"`
			MaxIntergalacticResearchNetwork int  `json:"MaxIntergalacticResearchNetwork"`
			MaxWeaponsTechnology            int  `json:"MaxWeaponsTechnology"`
			MaxShieldingTechnology          int  `json:"MaxShieldingTechnology"`
			MaxArmourTechnology             int  `json:"MaxArmourTechnology"`
			Target                          struct {
				Galaxy   int    `json:"Galaxy"`
				System   int    `json:"System"`
				Position int    `json:"Position"`
				Type     string `json:"Type"`
			} `json:"Target"`
			Trasports struct {
				Active      bool   `json:"Active"`
				CargoType   string `json:"CargoType"`
				DeutToLeave int    `json:"DeutToLeave"`
				Origin      struct {
					Galaxy   int    `json:"Galaxy"`
					System   int    `json:"System"`
					Position int    `json:"Position"`
					Type     string `json:"Type"`
				} `json:"Origin"`
			} `json:"Trasports"`
			CheckIntervalMin int `json:"CheckIntervalMin"`
			CheckIntervalMax int `json:"CheckIntervalMax"`
		} `json:"AutoResearch"`
		AutoCargo struct {
			Active                  bool          `json:"Active"`
			ExcludeMoons            bool          `json:"ExcludeMoons"`
			CargoType               string        `json:"CargoType"`
			RandomOrder             bool          `json:"RandomOrder"`
			MaxCargosToBuild        int           `json:"MaxCargosToBuild"`
			MaxCargosToKeep         int           `json:"MaxCargosToKeep"`
			SkipIfIncomingTransport bool          `json:"SkipIfIncomingTransport"`
			Exclude                 []interface{} `json:"Exclude"`
			CheckIntervalMin        int           `json:"CheckIntervalMin"`
			CheckIntervalMax        int           `json:"CheckIntervalMax"`
		} `json:"AutoCargo"`
		AutoRepatriate struct {
			Active           bool `json:"Active"`
			ExcludeMoons     bool `json:"ExcludeMoons"`
			MinimumResources int  `json:"MinimumResources"`
			LeaveDeut        struct {
				OnlyOnMoons bool `json:"OnlyOnMoons"`
				DeutToLeave int  `json:"DeutToLeave"`
			} `json:"LeaveDeut"`
			Target struct {
				Galaxy   int    `json:"Galaxy"`
				System   int    `json:"System"`
				Position int    `json:"Position"`
				Type     string `json:"Type"`
			} `json:"Target"`
			CargoType               string `json:"CargoType"`
			RandomOrder             bool   `json:"RandomOrder"`
			SkipIfIncomingTransport bool   `json:"SkipIfIncomingTransport"`
			Exclude                 []struct {
				Galaxy   int    `json:"Galaxy"`
				System   int    `json:"System"`
				Position int    `json:"Position"`
				Type     string `json:"Type"`
			} `json:"Exclude"`
			CheckIntervalMin int `json:"CheckIntervalMin"`
			CheckIntervalMax int `json:"CheckIntervalMax"`
		} `json:"AutoRepatriate"`
		BuyOfferOfTheDay struct {
			Active           bool `json:"Active"`
			CheckIntervalMin int  `json:"CheckIntervalMin"`
			CheckIntervalMax int  `json:"CheckIntervalMax"`
		} `json:"BuyOfferOfTheDay"`
	} `json:"Brain"`
	Expeditions struct {
		Active              bool `json:"Active"`
		AutoSendExpeditions struct {
			Active                         bool   `json:"Active"`
			Comment                        string `json:"_comment"`
			Comment2                       string `json:"_comment2"`
			MainShip                       string `json:"MainShip"`
			MinCargosToSend                int    `json:"MinCargosToSend"`
			WaitForAllExpeditions          bool   `json:"WaitForAllExpeditions"`
			SplitExpeditionsBetweenSystems bool   `json:"SplitExpeditionsBetweenSystems"`
			RandomizeOrder                 bool   `json:"RandomizeOrder"`
			FuelToCarry                    int    `json:"FuelToCarry"`
			Origin                         []struct {
				Galaxy   int    `json:"Galaxy"`
				System   int    `json:"System"`
				Position int    `json:"Position"`
				Type     string `json:"Type"`
			} `json:"Origin"`
		} `json:"AutoSendExpeditions"`
	} `json:"Expeditions"`
	AutoHarvest struct {
		Active           bool `json:"Active"`
		HarvestOwnDF     bool `json:"HarvestOwnDF"`
		HarvestDeepSpace bool `json:"HarvestDeepSpace"`
		MinimumResources int  `json:"MinimumResources"`
		CheckIntervalMin int  `json:"CheckIntervalMin"`
		CheckIntervalMax int  `json:"CheckIntervalMax"`
	} `json:"AutoHarvest"`
	TelegramMessenger struct {
		Active bool   `json:"Active"`
		API    string `json:"API"`
		ChatID string `json:"ChatId"`
	} `json:"TelegramMessenger"`
}

func SettingsHandler(c echo.Context) error {
	//bot := c.Get("bot").(*ogame.OGame)

	if c.Request().Method == "GET" {
		dataFile, _ := os.ReadFile("settings.json")
		settingsHTML := `<!DOCTYPE html>
<html>
	<head>
		<script src="https://code.jquery.com/jquery-3.6.0.min.js"></script>
	</head>
	<body>
		<form id="idForm" action="/settings" method="post">
			<button type="button" id="submitButtonId" value="submitButtonId">Click Me!</button><br/>

			<textarea name="settings" cols="120" rows="50">` + string(dataFile) + `</textarea>
		</form>
		<script>
		/*
			Const click = document.getElementById('button')
				Function f1 {
					Alert(` + "`" + `${click.Value}` + "`" + `)
				}
		*/ 	
			// this is the id of the submit button
			$("#submitButtonId").click(function() {
				//alert("Clicked");
			
				var url = "/settings"; // the script where you handle the form input.
			
				$.ajax({
						type: "POST",
						url: url,
						data: $("#idForm").serialize(), // serializes the form's elements.
						success: function(data)
						{
							//alert(data); // show response from the php script.
						}
						});
			
				return false; // avoid to execute the actual submit of the form.
			});
		</script>
	</body>
</html>`
		return c.HTML(http.StatusOK, settingsHTML)
	}

	if c.Request().Method == "POST" {
		c.Request().ParseForm()
		settings := c.FormValue("settings")
		os.WriteFile("settings.json", []byte(settings), 0644)
		settingsHTML := `<!DOCTYPE html>
		<html>
			<head>
				<script src="https://code.jquery.com/jquery-3.6.0.min.js"></script>
			</head>
			<body>
				<form id="idForm" action="/settings" method="post">
					<button type="button" id="submitButtonId" value="submitButtonId">Click Me!</button><br/>
		
					<textarea name="settings" cols="40" rows="50">` + settings + `</textarea>
				</form>
				<script>
				/*
					Const click = document.getElementById('button')
						Function f1 {
							Alert(` + "`" + `${click.Value}` + "`" + `)
						}
				*/ 	
					// this is the id of the submit button
					$("#submitButtonId").click(function() {
						//alert("Clicked");
					
						var url = "/settings"; // the script where you handle the form input.
					
						$.ajax({
								type: "POST",
								url: url,
								data: $("#idForm").serialize(), // serializes the form's elements.
								success: function(data)
								{
									//alert(data); // show response from the php script.
								}
								});
					
						return false; // avoid to execute the actual submit of the form.
					});
				</script>
			</body>
		</html>`
		settingsHTML = ""
		return c.HTML(http.StatusOK, settingsHTML)
	}

	return c.HTML(http.StatusOK, "nothing here")
	//return c.JSON(http.StatusOK, ogame.SuccessResp("success"))

}

type feature struct {
	enabled bool
	context.Context
	context.CancelFunc
}

var features map[string]feature = make(map[string]feature)

func FeatureDefenderHandler(c echo.Context) error {
	query := c.Request().URL.Query()
	enable := query.Get("enable")
	featureText := query.Get("feature")

	if v, ok := features[featureText]; ok {
		if v.enabled == true && enable == "disable" {
			v.CancelFunc()
			delete(features, featureText)
		}
	}

	if _, ok := features[featureText]; !ok {
		if enable == "enable" {
			ctx, cancel := context.WithCancel(context.Background())
			go featureDefenderFunc(ctx)
			features[featureText] = feature{
				enabled:    true,
				Context:    ctx,
				CancelFunc: cancel,
			}
		}
	}

	return c.JSON(http.StatusOK, ogame.SuccessResp("success"))
}

func featureDefenderFunc(ctx context.Context) {
	log.Println("Defender Enabled")
	for {
		select {
		case <-ctx.Done():
			log.Println("disable Defender")
			return
		case <-time.After(15 * time.Second):
			log.Println("Check for Attacks...")
		}
	}
}

var tx ogame.Prioritizable = &ogame.Prioritize{}

func GetToggleManualModeHandler(c echo.Context) error {
	bot := c.Get("bot").(*ogame.OGame)
	locked, state := bot.GetState()
	if locked && state == "Manuel Mode" {
		log.Println("Disable Manual Mode")
		// disable Manuel Mode
		tx.Done()
		return c.JSON(http.StatusOK, false)
	} else {
		log.Println("Enable Manual Mode")
		tx = bot.BeginNamed("Manuel Mode")
		return c.JSON(http.StatusOK, true)
	}
}

func GetVacationModeHandler(c echo.Context) error {
	bot := c.Get("bot").(*ogame.OGame)
	data := struct {
		VacationModeEnabled bool
	}{
		VacationModeEnabled: bot.IsVacationModeEnabled(),
	}
	return c.JSON(http.StatusOK, data)
}

func GetEmpireFromGameHandler(c echo.Context) error {
	bot := c.Get("bot").(*ogame.OGame)
	category, _ := strconv.ParseInt(c.QueryParam("type"), 10, 64)
	data, _ := bot.GetEmpireJSON(category)
	return c.JSON(http.StatusOK, data)
}

func GetAbandonHandler(c echo.Context) error {
	bot := c.Get("bot").(*ogame.OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ogame.ErrorResp(400, "invalid planet id"))
	}
	err = bot.Abandon(ogame.CelestialID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ogame.ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, err)
}
