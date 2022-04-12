package main

import (
	"bytes"
	"compress/gzip"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	//badger "github.com/dgraph-io/badger/v3"
	"github.com/faunX/ogame"
	"github.com/faunX/ogame/cmd/ogamed/database"
	"github.com/faunX/ogame/cmd/ogamed/ogb"
	"github.com/faunX/ogame/cmd/ogamed/skv"
	"github.com/faunX/ogame/cmd/ogamed/webserver"
	"github.com/faunX/ogame/cmd/ogamed/webserver/bindata"
	"github.com/inhies/go-bytesize"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/urfave/cli.v2"
)

var version = "0.0.0"
var commit = ""
var date = ""

var (
	LogWarn  = log.New(os.Stderr, "[ warn  ]", log.Ltime|log.Lshortfile)
	LogInfo  = log.New(os.Stderr, "[ info  ]", log.Ltime|log.Lshortfile)
	LogDebug = log.New(os.Stderr, "[ debug ]", log.Ltime|log.Lshortfile)
	LogError = log.New(os.Stderr, "[ error ]", log.Ltime|log.Lshortfile)
)

var myLogger *log.Logger = log.New(os.Stdout, "", 0)

//var BadgerDB *badger.DB

func main() {
	LogInfo.Print("Hello")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	database.InitDatabase()

	go webserver.HandleAttackCh()

	// Open the Badger database located in the /tmp/badger directory.
	// It will be created if it doesn't exist.
	//var err error
	// opt := badger.DefaultOptions("./storage.db")
	// opt.NumMemtables = 1
	// opt.ValueLogFileSize = 1 << 20
	// BadgerDB, err = badger.Open(opt)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer BadgerDB.Close()
	store, err := skv.Open("./sessions.db")
	if err != nil {
		panic(err)
	}
	webserver.WebKVStore = store
	// put: encodes value with gob and updates the boltdb
	//err := svk.Put(sessionId, info)

	// get: fetches from boltdb and does gob decode
	//err := svk.Get(sessionId, &info)

	// delete: seeks in boltdb and deletes the record
	//err := svk.Delete(sessionId)

	// close the store
	defer store.Close()

	app := cli.App{}
	app.Authors = []*cli.Author{
		{Name: "Alain Gilbert", Email: "alain.gilbert.15@gmail.com"},
	}
	app.Name = "ogamed2"
	app.Usage = "ogame deamon service"
	app.Version = version
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "universe",
			Usage:   "Universe name",
			Aliases: []string{"u"},
			EnvVars: []string{"OGAMED_UNIVERSE"},
		},
		&cli.StringFlag{
			Name:    "playerid",
			Usage:   "playerid of ogame server",
			Value:   "0",
			Aliases: []string{"id"},
			EnvVars: []string{"OGAMED_PLAYERID"},
		},
		&cli.StringFlag{
			Name:    "username",
			Usage:   "Email address to login on ogame",
			Aliases: []string{"e"},
			EnvVars: []string{"OGAMED_USERNAME"},
		},
		&cli.StringFlag{
			Name:    "password",
			Usage:   "Password to login on ogame",
			Aliases: []string{"p"},
			EnvVars: []string{"OGAMED_PASSWORD"},
		},
		&cli.StringFlag{
			Name:    "language",
			Usage:   "Language to login on ogame",
			Value:   "en",
			Aliases: []string{"l"},
			EnvVars: []string{"OGAMED_LANGUAGE"},
		},
		&cli.StringFlag{
			Name:    "host",
			Usage:   "HTTP host",
			Value:   "127.0.0.1",
			EnvVars: []string{"OGAMED_HOST"},
		},
		&cli.IntFlag{
			Name:    "port",
			Usage:   "HTTP port",
			Value:   8080,
			EnvVars: []string{"OGAMED_PORT"},
		},
		&cli.BoolFlag{
			Name:    "auto-login",
			Usage:   "Login when process starts",
			Value:   true,
			EnvVars: []string{"OGAMED_AUTO_LOGIN"},
		},
		&cli.StringFlag{
			Name:    "proxy",
			Usage:   "Proxy address",
			Value:   "",
			EnvVars: []string{"OGAMED_PROXY"},
		},
		&cli.StringFlag{
			Name:    "proxy-username",
			Usage:   "Proxy username",
			Value:   "",
			EnvVars: []string{"OGAMED_PROXY_USERNAME"},
		},
		&cli.StringFlag{
			Name:    "proxy-password",
			Usage:   "Proxy password",
			Value:   "",
			EnvVars: []string{"OGAMED_PROXY_PASSWORD"},
		},
		&cli.StringFlag{
			Name:    "proxy-type",
			Usage:   "Proxy type (socks5/http)",
			Value:   "socks5",
			EnvVars: []string{"OGAMED_PROXY_TYPE"},
		},
		&cli.BoolFlag{
			Name:    "proxy-login-only",
			Usage:   "Proxy login requests only",
			Value:   false,
			EnvVars: []string{"OGAMED_PROXY_LOGIN_ONLY"},
		},
		&cli.StringFlag{
			Name:    "lobby",
			Usage:   "Lobby to use (lobby | lobby-pioneers)",
			Value:   "lobby",
			EnvVars: []string{"OGAMED_PROXY_PASSWORD"},
		},
		&cli.StringFlag{
			Name:    "api-new-hostname",
			Usage:   "New OGame Hostname eg: https://someuniverse.example.com",
			Value:   "",
			EnvVars: []string{"OGAMED_NEW_HOSTNAME"},
		},
		&cli.StringFlag{
			Name:    "basic-auth-username",
			Usage:   "Basic auth username eg: admin",
			Value:   "",
			EnvVars: []string{"OGAMED_AUTH_USERNAME"},
		},
		&cli.StringFlag{
			Name:    "basic-auth-password",
			Usage:   "Basic auth password eg: secret",
			Value:   "",
			EnvVars: []string{"OGAMED_AUTH_PASSWORD"},
		},
		&cli.StringFlag{
			Name:    "enable-tls",
			Usage:   "Enable TLS. Needs key.pem and cert.pem",
			Value:   "false",
			EnvVars: []string{"OGAMED_ENABLE_TLS"},
		},
		&cli.StringFlag{
			Name:    "tls-key-file",
			Usage:   "Path to key.pem",
			Value:   "~/.ogame/key.pem",
			EnvVars: []string{"OGAMED_TLS_CERTFILE"},
		},
		&cli.StringFlag{
			Name:    "tls-cert-file",
			Usage:   "Path to cert.pem",
			Value:   "~/.ogame/cert.pem",
			EnvVars: []string{"OGAMED_TLS_KEYFILE"},
		},
		&cli.StringFlag{
			Name:    "cookies-filename",
			Usage:   "Path cookies file",
			Value:   "",
			EnvVars: []string{"OGAMED_COOKIES_FILENAME"},
		},
		&cli.BoolFlag{
			Name:    "cors-enabled",
			Usage:   "Enable CORS",
			Value:   true,
			EnvVars: []string{"CORS_ENABLED"},
		},
		&cli.StringFlag{
			Name:    "nja-api-key",
			Usage:   "Ninja API key",
			Value:   "",
			EnvVars: []string{"NJA_API_KEY"},
		},
		&cli.StringFlag{
			Name:    "telegram-token",
			Usage:   "Telegram Token",
			Value:   "",
			EnvVars: []string{"OGAMED_TELEGRAM_TOKEN"},
		},
		&cli.StringFlag{
			Name:    "telegram-chatid",
			Usage:   "Telegram Chat ID",
			Value:   "",
			EnvVars: []string{"OGAMED_TELEGRAM_CHATID"},
		},
	}
	app.Action = start
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func start(c *cli.Context) error {
	universe := c.String("universe")
	playerid := c.Int64("playerid")
	username := c.String("username")
	password := c.String("password")
	language := c.String("language")
	autoLogin := c.Bool("auto-login")
	//autoLogin = true
	host := c.String("host")
	port := c.Int("port")
	proxyAddr := c.String("proxy")
	proxyUsername := c.String("proxy-username")
	proxyPassword := c.String("proxy-password")
	proxyType := c.String("proxy-type")
	proxyLoginOnly := c.Bool("proxy-login-only")
	lobby := c.String("lobby")
	apiNewHostname := c.String("api-new-hostname")
	enableTLS := c.Bool("enable-tls")
	tlsKeyFile := c.String("tls-key-file")
	tlsCertFile := c.String("tls-cert-file")
	basicAuthUsername := c.String("basic-auth-username")
	basicAuthPassword := c.String("basic-auth-password")
	cookiesFilename := c.String("cookies-filename")
	corsEnabled := c.Bool("cors-enabled")
	njaApiKey := c.String("nja-api-key")
	telegramToken := c.String("telegram-token")
	telegramChatID := c.Int64("telegram-chatid")

	params := ogame.Params{
		Universe:        universe,
		PlayerID:        playerid,
		Username:        username,
		Password:        password,
		Lang:            language,
		AutoLogin:       autoLogin,
		Proxy:           proxyAddr,
		ProxyUsername:   proxyUsername,
		ProxyPassword:   proxyPassword,
		ProxyType:       proxyType,
		ProxyLoginOnly:  proxyLoginOnly,
		Lobby:           lobby,
		APINewHostname:  apiNewHostname,
		CookiesFilename: cookiesFilename,
	}

	if njaApiKey != "" {
		params.CaptchaCallback = ogame.NinjaSolver(njaApiKey)
	}

	if telegramToken != "" && telegramChatID != 0 {
		params.CaptchaCallback = ogame.TelegramSolver(telegramToken, telegramChatID)
	}

	bot, err := ogame.NewWithParams(params)
	if err != nil {
		return err
	}

	database := ogb.New()
	var dataFile []byte
	var fileMutex sync.Mutex
	// _, err = bot.LoginWithExistingCookies()
	// if err != nil {
	// 	log.Println(err)
	// 	return err
	// }
	servers, err := bot.GetServers()
	accounts, err := bot.GetUserAccounts()
	account, _, err := bot.FindAccount(universe, language, 0, accounts, servers)
	if err != nil {
		log.Println(err)
	}

	if len(dataFile) == 0 {
		// Load old Data

		log.Println("read Data " + strconv.FormatInt(account.ID, 10) + "-" + universe + "-" + language + " file")

		dataFile, err = os.ReadFile(strconv.FormatInt(account.ID, 10) + "-" + universe + "-" + language)
		if err != nil {
			// No Data File found
			log.Println(err)
		} else {
			// Decompress
			var bufRead bytes.Buffer
			gunzipWrite(&bufRead, dataFile)
			err := json.Unmarshal(bufRead.Bytes(), &database)

			//err := json.Unmarshal(dataFile, &database)
			if err != nil {
				log.Println(err)
			}
		}
	}

	sseData = make(chan *ogb.Database)
	go updateSSEdata(database, bot)
	// initialize global variable
	sseChannel = SSEChannel{
		Clients:  make([]chan string, 0),
		Notifier: make(chan string),
	}

	// done signal to go routine.
	broadcastDone := make(chan interface{})
	defer close(broadcastDone)

	// run our broadcaster go routine.
	go broadcaster(broadcastDone)

	bot.RegisterHTMLInterceptor(func(method string, urlString string, params url.Values, payload url.Values, pageHTML []byte) {
		database.Traffic.In = bot.BytesDownloaded()
		database.Traffic.Out = bot.BytesUploaded()
		database.Traffic.RPS = int64(bot.Client.GetRPS())

		database.Lock()
		defer database.Unlock()

		e := bot.GetExtractor()
		uri, _ := url.Parse(urlString)
		var page string

		if params.Has("cp") {
			tmpCP, _ := strconv.ParseInt(params.Get("cp"), 10, 64)
			if tmpCP != 0 {
				database.LastActiveCelestialID = ogame.CelestialID(tmpCP)
			}
		}

		if method == "POST" {
			//log.Println(params)
			//log.Println(payload)

			if params.Get("component") == "galaxyContent" {
				systemInfos, err := e.ExtractGalaxyInfos(pageHTML, bot.Player.PlayerName, bot.Player.PlayerID, bot.Player.Rank)
				//var planetInfos []ogame.PlanetInfos = make([]ogame.PlanetInfos, 15)
				// i := 0
				// systemInfos.Each(func(p *ogame.PlanetInfos) {
				// 	var planetInfo ogame.PlanetInfos
				// 	data, _ := json.Marshal(p)
				// 	json.Unmarshal(data, &planetInfo)
				// 	if planetInfo.Player.ID != 0 {
				// 		planetInfos[i] = planetInfo
				// 	}
				// 	i++
				// })

				if err == nil {
					database.Galaxy[strconv.FormatInt(systemInfos.Galaxy(), 10)+":"+strconv.FormatInt(systemInfos.System(), 10)] = struct {
						SystemInfos ogame.SystemInfos "json:\"systemInfos\""
						//PlanetInfos []ogame.PlanetInfos "json:\"planetInfo\""
					}{
						SystemInfos: systemInfos,
						//PlanetInfos: planetInfos,
					}
				}
			}

			if params.Get("component") == "fleetdispatch" && params.Get("action") == "sendFleet" && params.Get("ajax") == "1" && params.Get("asJson") == "1" {
				token := payload.Get("token")

				fleetSpeed, _ := strconv.ParseInt(payload.Get("speed"), 10, 64)

				fleetResources := ogame.Resources{}
				fleetResources.Metal, _ = strconv.ParseInt(payload.Get("metal"), 10, 64)
				fleetResources.Crystal, _ = strconv.ParseInt(payload.Get("crystal"), 10, 64)
				fleetResources.Deuterium, _ = strconv.ParseInt(payload.Get("deuterium"), 10, 64)

				mission, _ := strconv.ParseInt(payload.Get("mission"), 10, 64)
				fleetMission := ogame.MissionID(mission)

				fleetHoldingtime, _ := strconv.ParseInt(payload.Get("holdingtime"), 10, 64)

				// Destination
				fleetDestination := ogame.Coordinate{}
				fleetDestination.Galaxy, _ = strconv.ParseInt(payload.Get("galaxy"), 10, 64)
				fleetDestination.System, _ = strconv.ParseInt(payload.Get("system"), 10, 64)
				fleetDestination.Position, _ = strconv.ParseInt(payload.Get("position"), 10, 64)
				fleetDestinationType, _ := strconv.ParseInt(payload.Get("type"), 10, 64)
				switch fleetDestinationType {
				case ogame.PlanetType.Int64():
					fleetDestination.Type = ogame.PlanetType
				case ogame.MoonType.Int64():
					fleetDestination.Type = ogame.MoonType
				case ogame.DebrisType.Int64():
					fleetDestination.Type = ogame.DebrisType
				}

				//fleetUnionID := strconv.ParseInt(payload.Get("union"), 10, 64)

				var fleetOrigincelestialID ogame.CelestialID
				if params.Has("cp") {
					fleetOriginID, _ := strconv.ParseInt(params.Get("cp"), 10, 64)
					fleetOrigincelestialID = ogame.CelestialID(fleetOriginID)
				} else {
					fleetOrigincelestialID = database.LastActiveCelestialID
				}

				originCelestial := bot.GetCachedCelestialByID(fleetOrigincelestialID)
				fleetOriginCoords := originCelestial.GetCoordinate()

				var fleetShips ogame.ShipsInfos
				for _, s := range ogame.Ships {

					if payload.Has("am" + strconv.FormatInt(s.GetID().Int64(), 10)) {
						nbr, _ := strconv.ParseInt(payload.Get("am"+strconv.FormatInt(s.GetID().Int64(), 10)), 10, 64)
						fleetShips.Set(s.GetID(), nbr)
					}
				}
				sentFleet := ogb.SentFleet{}

				sentFleet.OriginID = fleetOrigincelestialID
				sentFleet.OriginCoords = fleetOriginCoords
				sentFleet.DestinationCoords = fleetDestination
				sentFleet.Mission = fleetMission
				sentFleet.Speed = ogame.Speed(fleetSpeed)
				sentFleet.HoldingTime = fleetHoldingtime
				sentFleet.Ships = fleetShips
				sentFleet.Resources = fleetResources
				sentFleet.Token = token

				var resStruct struct {
					Success           bool          `json:"success"`
					Message           string        `json:"message"`
					FleetSendingToken string        `json:"fleetSendingToken"`
					Components        []interface{} `json:"components"`
					RedirectURL       string        `json:"redirectUrl"`
					Errors            []struct {
						Message string `json:"message"`
						Error   int64  `json:"error"`
					} `json:"errors"`
				}
				if err := json.Unmarshal(pageHTML, &resStruct); err == nil {
					if resStruct.Success == true {
						database.SentFleets = append(database.SentFleets, sentFleet)
					}
				}
				return
			}

			if params.Get("page") == "ajax" && (params.Get("component") == "traderauctioneer" || params.Get("component") == "traderimportexport") && params.Get("action") == "refreshPlanet" && params.Get("asJson") == "1" && params.Get("ajax") == "1" {
				planetID, _ := strconv.ParseInt(payload.Get("planetId"), 10, 64)
				database.Activities[ogame.CelestialID(planetID)] = time.Now()
			}
		} else {
			uri, _ = url.Parse(urlString)
			if uri.Query().Get("page") == "ingame" {
				page = uri.Query().Get("component")
			}

			timestamp := time.Unix(e.ExtractOGameTimestampFromBytes(pageHTML), 0)
			celestialID, err := e.ExtractPlanetID(pageHTML)
			if err == nil {
				database.LastActiveCelestialID = celestialID
				database.Activities[celestialID] = timestamp
			}

			if celestialID == 0 {
				celestialID = database.LastActiveCelestialID
			}

			if uri.Query().Get("page") == ogame.FetchTechs {
				cp, _ := strconv.ParseInt(uri.Query().Get("cp"), 10, 64)
				celestialID = ogame.CelestialID(cp)
				if cp == 0 {
					cp = int64(database.LastActiveCelestialID)
				}
				if cp == 0 {
					return
				}
				_, ok := database.ShipsInfos[celestialID]
				if ok {
					delete(database.ShipsInfos, ogame.CelestialID(0))
				}
				database.ResourcesBuildings[celestialID], database.Facilities[celestialID], database.ShipsInfos[celestialID], database.DefensesInfos[celestialID], database.Researches, _ = e.ExtractTechs(pageHTML)
				return
			}

			if uri.Query().Get("page") == ogame.FetchResourcesAjaxPage {
				cp, _ := strconv.ParseInt(uri.Query().Get("cp"), 10, 64)
				celestialID = ogame.CelestialID(cp)
				database.ResourcesDetails[celestialID], _ = e.ExtractResourcesDetails(pageHTML)
				return
			}

			planets := e.ExtractPlanets(pageHTML, bot)
			if len(planets) > 0 {
				database.Planets = planets
			}
			delete(database.ResourcesDetails, 0)
			database.Celestials = bot.GetCachedPlanets()
			if celestialID != 0 && ogame.IsKnowFullPage(params) {
				database.ResourcesDetails[celestialID] = e.ExtractResourcesDetailsFromFullPage(pageHTML)
			}
			attacks := make([]ogame.AttackEvent, 0)
			attacks, err = e.ExtractAttacks(pageHTML)
			if err == nil {
				database.AttackEvents = attacks
				for _, a := range attacks {
					_, b := database.AttackEventsHistory[ogame.FleetID(a.ID)]
					if !b {
						if len(webserver.OnAttackCh) >= 10 {
							<-webserver.OnAttackCh
						}
						webserver.OnAttackCh <- a
						database.AttackEventsHistory[ogame.FleetID(a.ID)] = a
					}

				}
			}
			eventListFleet := make([]ogame.Fleet, 0)
			eventListFleet = e.ExtractFleetsFromEventList(pageHTML)
			if len(eventListFleet) > 0 {
				database.EventFleets = eventListFleet
				for _, f := range eventListFleet {
					database.EventFleetsHistory[f.ID] = f
				}
			}

			if page == ogame.SuppliesPage || page == ogame.FacilitiesPage {
				buildingID, buildingCountDown, _, _ := e.ExtractConstructions(pageHTML)
				if buildingCountDown > 0 {
					database.Constructions[celestialID] = ogb.BrainConstruction{
						Quantifiable: ogame.Quantifiable{ID: buildingID, Nbr: buildingCountDown},
						FinishAt:     timestamp.Add(time.Duration(buildingCountDown) * time.Second),
					}
				}
			}

			switch page {
			case ogame.OverviewPage:
				buildingID, buildingCountDown, researchID, researchCountdown := e.ExtractConstructions(pageHTML)

				if buildingCountDown > 0 {
					database.Constructions[celestialID] = ogb.BrainConstruction{
						Quantifiable: ogame.Quantifiable{ID: buildingID, Nbr: buildingCountDown},
						FinishAt:     timestamp.Add(time.Duration(buildingCountDown) * time.Second),
					}
				}

				if researchCountdown > 0 {
					database.ResearchInProgress = ogb.BrainConstruction{
						Quantifiable: ogame.Quantifiable{ID: researchID, Nbr: researchCountdown},
						FinishAt:     timestamp.Add(time.Duration(researchCountdown) * time.Second),
					}
				}

				productions, productionCountDown, err := e.ExtractProduction(pageHTML)
				if err == nil && productionCountDown > 0 {
					database.Productions[celestialID] = ogb.BrainProduction{
						Productions: productions,
						FinishAt:    timestamp.Add(time.Duration(productionCountDown) * time.Second),
					}
				}
				break
			case ogame.SuppliesPage:
				res, err := e.ExtractResourcesBuildings(pageHTML)
				if err == nil {
					database.ResourcesBuildings[celestialID] = res
				} else {
				}
				break
			case ogame.FacilitiesPage:
				fac, err := e.ExtractFacilities(pageHTML)
				if err == nil {
					database.Facilities[celestialID] = fac
				}
				break
			case ogame.ShipyardPage:
				productions, productionCountDown, err := e.ExtractProduction(pageHTML)
				if err == nil && productionCountDown > 0 {
					database.Productions[celestialID] = ogb.BrainProduction{
						Productions: productions,
						FinishAt:    timestamp.Add(time.Duration(productionCountDown) * time.Second),
					}
				}
				shipyard, err := e.ExtractShips(pageHTML)
				if err == nil {
					database.ShipsInfos[celestialID] = shipyard
				}
				break
			case ogame.DefensesPage:
				productions, productionCountDown, err := e.ExtractProduction(pageHTML)
				if err == nil && productionCountDown > 0 {
					database.Productions[celestialID] = ogb.BrainProduction{
						Productions: productions,
						FinishAt:    timestamp.Add(time.Duration(productionCountDown) * time.Second),
					}
				}
				defenses, err := e.ExtractDefense(pageHTML)
				if err == nil {
					database.DefensesInfos[celestialID] = defenses
				}
				break
			case ogame.MovementPage:
				database.Movements = e.ExtractFleets(pageHTML, bot.Location())
				database.Slots = e.ExtractSlots(pageHTML)
				break
			case ogame.ResearchPage:
				database.Researches = e.ExtractResearch(pageHTML)
				_, _, researchID, researchCountdown := e.ExtractConstructions(pageHTML)
				if researchCountdown > 0 {
					database.ResearchInProgress = ogb.BrainConstruction{
						Quantifiable: ogame.Quantifiable{ID: researchID, Nbr: researchCountdown},
						FinishAt:     timestamp.Add(time.Duration(researchCountdown) * time.Second),
					}
				}
				break
			case ogame.FleetdispatchPage:
				tokenM := regexp.MustCompile(`var fleetSendingToken = "([^"]+)";`).FindSubmatch(pageHTML)
				if bot.IsV8() {
					tokenM = regexp.MustCompile(`var token = "([^"]+)";`).FindSubmatch(pageHTML)
				}
				if len(tokenM) != 2 {
					return
				}
				ships := e.ExtractFleet1Ships(pageHTML)
				ships.SolarSatellite = database.ShipsInfos[celestialID].SolarSatellite
				ships.Crawler = database.ShipsInfos[celestialID].Crawler
				database.ShipsInfos[celestialID] = ships
				doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
				researchFleet1 := e.ExtractFleet1ResearchesFromDoc(doc)

				database.Researches.WeaponsTechnology = researchFleet1.WeaponsTechnology
				database.Researches.ShieldingTechnology = researchFleet1.ShieldingTechnology
				database.Researches.ArmourTechnology = researchFleet1.ArmourTechnology
				database.Researches.HyperspaceTechnology = researchFleet1.HyperspaceTechnology
				database.Researches.CombustionDrive = researchFleet1.CombustionDrive
				database.Researches.ImpulseDrive = researchFleet1.ImpulseDrive
				database.Researches.HyperspaceDrive = researchFleet1.HyperspaceDrive
				break
			}

		}

		// Get Messages
		if uri.Query().Get("page") == "messages" {
			tab := uri.Query().Get("tab")
			if uri.Query().Get("tabid") != "" {
				tab = uri.Query().Get("tabid")
			}

			if method == "POST" {
				/*
					messageId	"-1"
					tabid	"22"
					action	"107"
					pagination	"2"
					ajax	"1"
				*/
				tab = payload.Get("tabid")
			}

			switch tab {
			// Fleets Tab
			case "20": // Espionage Reports
				if uri.Query().Get("messageId") != "" {
					m, err := e.ExtractEspionageReport(pageHTML, bot.Location())
					if err == nil {
						_, exists := database.EspionageReports[m.ID]
						if !exists {
							database.EspionageReports[m.ID] = m
						}
					}
				} else {
					newMessages, _ := e.ExtractEspionageReportMessageIDs(pageHTML)
					if err == nil {
						for _, m := range newMessages {
							_, exists := database.EspionageReportSummary[m.ID]
							if !exists {
								database.EspionageReportSummary[m.ID] = m
							}
						}
					}
				}
				break
			case "21": // Combat Reports
				newMessages, _ := e.ExtractCombatReportMessagesSummary(pageHTML)
				if err == nil {
					for _, m := range newMessages {
						_, exists := database.EspionageReportSummary[m.ID]
						if !exists {
							database.CombatReportSummary[m.ID] = m
						}
					}
				}
				if payload.Get("messageId") != "" || uri.Query().Get("messageId") != "" {
					id, _ := strconv.ParseInt(payload.Get("messageId"), 10, 64)
					if id == 0 {
						id, _ = strconv.ParseInt(uri.Query().Get("messageId"), 10, 64)
					}
					if id == -1 {
						return
					}
					fullCombatReport, err := e.ExtractFullCombatReport(pageHTML)
					if err == nil {
						_, exists := database.FullCombatReports[id]
						if !exists {
							log.Println("Add Full Combat Report to Datatbase")
							database.FullCombatReports[id] = fullCombatReport
						}
					} else {
						log.Println(err)
					}
				}
				break
			case "22": // Expedition Reports
				messages, _, err := e.ExtractExpeditionMessages(pageHTML, bot.Location())
				if err == nil {
					for _, m := range messages {
						_, exists := database.ExpeditionMessages[m.ID]
						if !exists {
							database.ExpeditionMessages[m.ID] = m
						}
					}
				}
				break
			default: // Transports
				//newMessages, newNbPage, _ := b.extractor.ExtractMessages(pageHTML, b.location)
				messages, _, err := e.ExtractMessages(pageHTML, bot.Location())
				if err == nil {
					for _, m := range messages {
						_, exists := database.Messages[m.ID]
						if !exists {
							tabInt, _ := strconv.ParseInt(tab, 10, 64)
							newMessage := struct {
								Tabid int64
								ogame.Message
							}{
								Tabid:   tabInt,
								Message: m,
							}
							database.Messages[m.ID] = newMessage
						}
					}

				}
				break
			}

		}
		// End Get Messages

		data, err := json.MarshalIndent(database, "", "  ")
		if err != nil {
			log.Println(err)
		}
		//os.WriteFile(strconv.FormatInt(bot.GetCachedPlayer().PlayerID, 10)+"-"+bot.GetUniverseName()+"-"+bot.GetLanguage(), data, 644)

		// compress
		fileMutex.Lock()
		defer fileMutex.Unlock()
		var bufWrite bytes.Buffer
		gzipWrite(&bufWrite, data)
		os.WriteFile(strconv.FormatInt(bot.GetCachedPlayer().PlayerID, 10)+"-"+bot.GetUniverseName()+"-"+bot.GetLanguage(), bufWrite.Bytes(), 644)
	})

	e := webserver.Start()

	if corsEnabled {
		e.Use(middleware.CORS())
	}
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx.Set("bot", bot)
			ctx.Set("ogb", database)
			//ctx.Set("badgerDB", BadgerDB)
			ctx.Set("version", version)
			ctx.Set("commit", commit)
			ctx.Set("date", date)
			ctx.Set("database", database)
			return next(ctx)
		}
	})
	if len(basicAuthUsername) > 0 && len(basicAuthPassword) > 0 {
		log.Println("Enable Basic Auth")
		e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
			// Be careful to use constant time comparison to prevent timing attacks
			if subtle.ConstantTimeCompare([]byte(username), []byte(basicAuthUsername)) == 1 &&
				subtle.ConstantTimeCompare([]byte(password), []byte(basicAuthPassword)) == 1 {
				return true, nil
			}
			return false, nil
		}))
	}
	e.HideBanner = true
	e.HidePort = true
	e.Debug = false
	e.GET("/", ogame.HomeHandler)
	e.GET("/tasks", ogame.TasksHandler)

	// CAPTCHA Handler
	e.GET("/bot/captcha", ogame.GetCaptchaHandler)
	e.GET("/bot/captcha/icons/:challengeID", ogame.GetCaptchaImgHandler)
	e.GET("/bot/captcha/question/:challengeID", ogame.GetCaptchaTextHandler)
	e.POST("/bot/captcha/solve", ogame.GetCaptchaSolverHandler)

	e.GET("/bot/server", ogame.GetServerHandler)
	e.GET("/bot/server-data", ogame.GetServerDataHandler)
	e.POST("/bot/set-user-agent", ogame.SetUserAgentHandler)
	e.GET("/bot/server-url", ogame.ServerURLHandler)
	e.GET("/bot/language", ogame.GetLanguageHandler)
	e.GET("/bot/empire/type/:typeID", ogame.GetEmpireHandler)
	e.POST("/bot/page-content", ogame.PageContentHandler)
	e.GET("/bot/login", ogame.LoginHandler)
	e.GET("/bot/logout", ogame.LogoutHandler)
	e.GET("/bot/username", ogame.GetUsernameHandler)
	e.GET("/bot/universe-name", ogame.GetUniverseNameHandler)
	e.GET("/bot/server/speed", ogame.GetUniverseSpeedHandler)
	e.GET("/bot/server/speed-fleet", ogame.GetUniverseSpeedFleetHandler)
	e.GET("/bot/server/version", ogame.ServerVersionHandler)
	e.GET("/bot/server/time", ogame.ServerTimeHandler)
	e.GET("/bot/is-under-attack", ogame.IsUnderAttackHandler)
	e.GET("/bot/is-vacation-mode", ogame.IsVacationModeHandler)
	e.GET("/bot/user-infos", ogame.GetUserInfosHandler)
	e.GET("/bot/character-class", ogame.GetCharacterClassHandler)
	e.GET("/bot/has-commander", ogame.HasCommanderHandler)
	e.GET("/bot/has-admiral", ogame.HasAdmiralHandler)
	e.GET("/bot/has-engineer", ogame.HasEngineerHandler)
	e.GET("/bot/has-geologist", ogame.HasGeologistHandler)
	e.GET("/bot/has-technocrat", ogame.HasTechnocratHandler)
	e.GET("/bot/get-messages", ogame.GetMessagesHandler)
	e.POST("/bot/send-message", ogame.SendMessageHandler)
	e.GET("/bot/fleets", ogame.GetFleetsHandler)
	e.GET("/bot/fleets/slots", ogame.GetSlotsHandler)
	e.POST("/bot/fleets/:fleetID/cancel", ogame.CancelFleetHandler)
	e.GET("/bot/espionage-report/:msgid", ogame.GetEspionageReportHandler)
	e.GET("/bot/espionage-report/:galaxy/:system/:position", ogame.GetEspionageReportForHandler)
	e.GET("/bot/espionage-report", ogame.GetEspionageReportMessagesHandler)
	e.POST("/bot/delete-report/:messageID", ogame.DeleteMessageHandler)
	e.POST("/bot/delete-all-espionage-reports", ogame.DeleteEspionageMessagesHandler)
	e.POST("/bot/delete-all-reports/:tabIndex", ogame.DeleteMessagesFromTabHandler)
	e.GET("/bot/attacks", ogame.GetAttacksHandler)
	e.GET("/bot/get-auction", ogame.GetAuctionHandler)
	e.POST("/bot/do-auction", ogame.DoAuctionHandler)
	e.GET("/bot/galaxy-infos/:galaxy/:system", ogame.GalaxyInfosHandler)
	e.GET("/bot/get-research", ogame.GetResearchHandler)
	e.GET("/bot/buy-offer-of-the-day", ogame.BuyOfferOfTheDayHandler)
	e.GET("/bot/price/:ogameID/:nbr", ogame.GetPriceHandler)
	e.GET("/bot/moons", ogame.GetMoonsHandler)
	e.GET("/bot/moons/:moonID", ogame.GetMoonHandler)
	e.GET("/bot/moons/:galaxy/:system/:position", ogame.GetMoonByCoordHandler)
	e.GET("/bot/celestials/:celestialID/items", ogame.GetCelestialItemsHandler)
	e.GET("/bot/celestials/:celestialID/items/:itemRef/activate", ogame.ActivateCelestialItemHandler)
	e.GET("/bot/celestials/:celestialID/techs", ogame.TechsHandler)
	e.GET("/bot/planets", ogame.GetPlanetsHandler)
	e.GET("/bot/planets/:planetID", ogame.GetPlanetHandler)
	e.GET("/bot/planets/:galaxy/:system/:position", ogame.GetPlanetByCoordHandler)
	e.GET("/bot/planets/:planetID/resources-details", ogame.GetResourcesDetailsHandler)
	e.GET("/bot/planets/:planetID/resource-settings", ogame.GetResourceSettingsHandler)
	e.POST("/bot/planets/:planetID/resource-settings", ogame.SetResourceSettingsHandler)
	e.GET("/bot/planets/:planetID/resources-buildings", ogame.GetResourcesBuildingsHandler)
	e.GET("/bot/planets/:planetID/defence", ogame.GetDefenseHandler)
	e.GET("/bot/planets/:planetID/ships", ogame.GetShipsHandler)
	e.GET("/bot/planets/:planetID/facilities", ogame.GetFacilitiesHandler)
	e.POST("/bot/planets/:planetID/build/:ogameID/:nbr", ogame.BuildHandler)
	e.POST("/bot/planets/:planetID/build/cancelable/:ogameID", ogame.BuildCancelableHandler)
	e.POST("/bot/planets/:planetID/build/production/:ogameID/:nbr", ogame.BuildProductionHandler)
	e.POST("/bot/planets/:planetID/build/building/:ogameID", ogame.BuildBuildingHandler)
	e.POST("/bot/planets/:planetID/build/technology/:ogameID", ogame.BuildTechnologyHandler)
	e.POST("/bot/planets/:planetID/build/defence/:ogameID/:nbr", ogame.BuildDefenseHandler)
	e.POST("/bot/planets/:planetID/build/ships/:ogameID/:nbr", ogame.BuildShipsHandler)
	e.POST("/bot/planets/:planetID/teardown/:ogameID", ogame.TeardownHandler)
	e.GET("/bot/planets/:planetID/production", ogame.GetProductionHandler)
	e.GET("/bot/planets/:planetID/constructions", ogame.ConstructionsBeingBuiltHandler)
	e.POST("/bot/planets/:planetID/cancel-building", ogame.CancelBuildingHandler)
	e.POST("/bot/planets/:planetID/cancel-research", ogame.CancelResearchHandler)
	e.GET("/bot/planets/:planetID/resources", ogame.GetResourcesHandler)
	e.POST("/bot/planets/:planetID/flighttime", ogame.FlightTimeHandler)
	e.POST("/bot/planets/:planetID/send-fleet", ogame.SendFleetHandler)
	e.GET("/bot/planets/:planetID/abandon", webserver.GetAbandonHandler)
	e.POST("/bot/planets/:planetID/send-ipm", ogame.SendIPMHandler)
	e.GET("/bot/moons/:moonID/phalanx/:galaxy/:system/:position", ogame.PhalanxHandler)
	e.POST("/bot/moons/:moonID/jump-gate", ogame.JumpGateHandler)
	e.GET("/game/allianceInfo.php", webserver.GetAlliancePageContentHandler2) // Example: //game/allianceInfo.php?allianceId=500127

	e.GET("/bot/vacation", webserver.GetVacationModeHandler)
	e.GET("/bot/empireJSON", webserver.GetEmpireFromGameHandler)

	// Get/Post Page Content
	//e.GET("/game/index.php", ogame.GetFromGameHandler)
	e.GET("/game/index.php", webserver.GetFromGameHandler2)

	e.POST("/game/index.php", webserver.PostToGameHandler2)

	// For AntiGame plugin
	// Static content
	e.GET("/cdn/*", ogame.GetStaticHandler)
	e.GET("/assets/css/*", ogame.GetStaticHandler)
	e.GET("/headerCache/*", ogame.GetStaticHandler)
	e.GET("/favicon.ico", ogame.GetStaticHandler)
	e.GET("/game/sw.js", ogame.GetStaticHandler)
	// JSON API
	/*
		/api/serverData.xml
		/api/localization.xml
		/api/players.xml
		/api/universe.xml
	*/
	e.GET("/api/*", ogame.GetStaticHandler)
	e.HEAD("/api/*", ogame.GetStaticHEADHandler) // AntiGame uses this to check if the cached XML files need to be refreshed

	// Own Pages
	e.GET("/empire", webserver.EmpireHandler)
	e.GET("/ip", func(c echo.Context) error {
		bot := c.Get("bot").(*ogame.OGame)
		publicip, err := bot.GetPublicIP()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ogame.ErrorResp(500, err.Error()))
		}
		return c.JSON(http.StatusOK, ogame.SuccessResp(publicip))
	})
	e.GET("/scripts", webserver.GetScriptsHandler)
	e.POST("/scripts/save", webserver.PostNewScriptHandler)
	e.POST("/scripts/run", webserver.RunScriptHandler)

	e.GET("/settings", webserver.SettingsHandler)
	e.POST("/settings", webserver.SettingsHandler)
	e.GET("/features", webserver.FeatureDefenderHandler)

	e.GET("/toggle-manual-mode", webserver.GetToggleManualModeHandler)
	e.POST("/toggle-manual-mode", webserver.GetToggleManualModeHandler)

	e.GET("/sse", sseHandler)

	//http://127.0.0.1:8080/public/html/ -> /webserver/bindata/
	assetHandler := http.FileServer(bindata.GetContent("assets"))
	e.GET("/public/*", echo.WrapHandler(http.StripPrefix("/", assetHandler)))
	jsHandler := http.FileServer(bindata.GetContent("js"))
	e.GET("/js/*", echo.WrapHandler(http.StripPrefix("/", jsHandler)))

	if enableTLS {
		log.Println("Enable TLS Support")
		return e.StartTLS(host+":"+strconv.Itoa(port), tlsCertFile, tlsKeyFile)
	}
	log.Println("Disable TLS Support")
	return e.Start(host + ":" + strconv.Itoa(port))
}

/*
type loop struct {
	quitCh chan bool
	Rerun  chan bool
}

func (l *loop) New() *loop {
	l.quitCh = make(chan bool)
	l.Rerun = make(chan bool)
	return l
}

func (l *loop) loop() {
	for {
		next := nextRandom()
		select {
		case rerun := <-l.Rerun:
			fmt.Println(rerun)
			fmt.Println("Rerun loop")
		case <-time.After(next):
			fmt.Println("time over")
		}
		fmt.Println("Wait " + next.String())
	}
}

// GetNextRandom ...
func nextRandom() time.Duration {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return time.Duration(r1.Int63n(10)+10) * time.Second
}

type test struct {
	MyMap map[string]string
}

func (b *test) NewTest() *test {
	var t test
	t.MyMap = make(map[string]string)
	return &t
}
*/

var sseData chan *ogb.Database

func updateSSEdata(o *ogb.Ogb, b *ogame.OGame) {
	for {
		//locked, actor := b.OnStateChange()
		data := o.GetDatabase()
		sseData <- data
		// var buf bytes.Buffer
		// enc := json.NewEncoder(&buf)
		// enc.Encode(data)
		//fmt.Println(buf.String())
		//sseChannel.Notifier <- buf.String()

		time.Sleep(1 * time.Second)
	}
}

type SSEChannel struct {
	Clients  []chan string
	Notifier chan string
}

var sseChannel SSEChannel

// sseHandler ...
func sseHandler(c echo.Context) error {
	bot := c.Get("bot").(*ogame.OGame)

	// c.Response().Header().Set("Content-Type", "text/event-stream")
	// c.Response().Header().Set("Cache-Control", "no-cache")
	// c.Response().Header().Set("Connection", "keepalive")

	// var buf bytes.Buffer
	// enc := json.NewEncoder(&buf)
	// enc.Encode(<-sseData)

	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	//c.Response().Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := c.Response().Writer.(http.Flusher)
	if !ok {
		http.Error(c.Response().Writer, "Connection doesnot support streaming", http.StatusBadRequest)
		return nil
	}

	//sseChan := make(chan string)
	//sseChannel.Clients = append(sseChannel.Clients, sseChan)

	//log.Printf("Connecting... %d", len(sseChannel.Clients))

	//defer fmt.Println("Closing channel.")
	timeout := time.Now().Add(60 * time.Second)
	for {
		select {
		case data := <-sseData:

			tasks := int64(bot.Client.GetRPS())
			in := data.Traffic.In
			out := data.Traffic.Out

			time.Sleep(1 * time.Second)
			data2 := struct {
				Traffic struct {
					In  string `json:"In"`
					Out string `json:"Out"`
					RPS int64  `json:"RPS"`
				} `json:"traffic"`
			}{}
			data2.Traffic.In = bytesize.New(float64(in)).String()
			data2.Traffic.Out = bytesize.New(float64(out)).String()
			data2.Traffic.RPS = tasks
			//fmt.Printf("%s", b)

			by, _ := json.Marshal(data2)

			fmt.Fprintf(c.Response().Writer, "data: %v \n\n", string(by))
			flusher.Flush()
			if time.Now().After(timeout) {
				//log.Println("Timeout Reached")
				return nil
			}
			// case data := <-sseChan:
			// 	fmt.Fprintf(c.Response().Writer, "data: %v \n\n", data)
			// 	flusher.Flush()
			// 	if time.Now().After(timeout) {
			// 		log.Println("Timeout Reached")
			// 		return
			// 	}
		}
	}

	//return c.String(http.StatusOK, "data: "+buf.String()+"\n\n")
}

func logHTTPRequest(w http.ResponseWriter, r *http.Request) {
	buf := new(strings.Builder)
	if _, err := io.Copy(buf, r.Body); err != nil {
		fmt.Printf("Error: %v", err)
	}
	method := r.Method

	logMsg := fmt.Sprintf("Method: %v, Body: %v", method, buf.String())
	fmt.Println(logMsg)
	sseChannel.Notifier <- logMsg
}

func trafficSSE(done <-chan interface{}, bot *ogame.OGame) {
	for {
		select {
		case <-done:
			return
		default:
			tasks := bot.GetTasks()
			in := bot.BytesDownloaded()
			out := bot.BytesUploaded()

			time.Sleep(1 * time.Second)
			data2 := struct {
				Traffic struct {
					In  int64 `json:"In"`
					Out int64 `json:"Out"`
					RPS int64 `json:"RPS"`
				} `json:"traffic"`
			}{}

			data2.Traffic.In = tasks.Total
			data2.Traffic.Out = in
			data2.Traffic.RPS = out

			var buf bytes.Buffer
			enc := json.NewEncoder(&buf)
			//enc.Encode(data)
			enc.Encode(data2)

			for _, channel := range sseChannel.Clients {
				channel <- buf.String()
			}
		}
		time.Sleep(1 * time.Second)
	}

}

func broadcaster(done <-chan interface{}) {
	fmt.Println("Broadcaster Started.")
	for {
		select {
		case <-done:
			return
		case data := <-sseData:

			data2 := struct {
				Traffic struct {
					In  int64 `json:"In"`
					Out int64 `json:"Out"`
					RPS int64 `json:"RPS"`
				} `json:"traffic"`
			}{}

			data2.Traffic.In = data.Traffic.In
			data2.Traffic.Out = data.Traffic.In
			data2.Traffic.RPS = data.Traffic.RPS

			var buf bytes.Buffer
			enc := json.NewEncoder(&buf)
			//enc.Encode(data)
			enc.Encode(data2)

			for _, channel := range sseChannel.Clients {
				channel <- buf.String()
			}
		case data := <-sseChannel.Notifier:
			for _, channel := range sseChannel.Clients {
				channel <- data
			}
		}
	}
}

// Write gzipped data to a Writer
func gzipWrite(w io.Writer, data []byte) error {
	// Write gzipped data to the client
	gw, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
	defer gw.Close()
	gw.Write(data)

	// Write gzipped data to the client
	// gr, err := gzip.NewReader(bytes.NewBuffer(data))
	// defer gr.Close()
	// _, err = io.Copy(w, gr)

	return err
}

// Write gunzipped data to a Writer
func gunzipWrite(w io.Writer, data []byte) error {
	// Write gzipped data to the client
	gr, err := gzip.NewReader(bytes.NewBuffer(data))
	defer gr.Close()
	data, err = ioutil.ReadAll(gr)
	if err != nil {
		return err
	}
	w.Write(data)
	return nil
}
