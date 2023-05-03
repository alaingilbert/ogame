package main

import (
	"crypto/subtle"
	"github.com/alaingilbert/ogame/pkg/device"
	"github.com/alaingilbert/ogame/pkg/wrapper"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/urfave/cli.v2"
	"log"
	"os"
	"strconv"
)

var version = "0.0.0"
var commit = ""
var date = ""

func main() {
	app := cli.App{}
	app.Authors = []*cli.Author{
		{Name: "Alain Gilbert", Email: "alain.gilbert.15@gmail.com"},
	}
	app.Name = "ogamed"
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
			EnvVars: []string{"OGAMED_LOBBY"},
		},
		&cli.StringFlag{
			Name:    "api-new-hostname",
			Usage:   "New OGame Hostname eg: https://someuniverse.example.com",
			Value:   "http://127.0.0.1:8080",
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
		}, &cli.StringFlag{
			Name:    "device-name",
			Usage:   "Set the Device Name",
			Value:   "device_name",
			EnvVars: []string{"OGAMED_DEVICENAME"},
		},
	}
	app.Action = start
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func start(c *cli.Context) error {
	universe := c.String("universe")
	username := c.String("username")
	password := c.String("password")
	language := c.String("language")
	autoLogin := c.Bool("auto-login")
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
	corsEnabled := c.Bool("cors-enabled")
	njaApiKey := c.String("nja-api-key")
	deviceName := c.String("device-name")
	// TODO: put device config in flags & env variables
	deviceInst, err := device.NewBuilder(deviceName).
		SetOsName(device.Windows).
		SetBrowserName(device.Chrome).
		SetMemory(8).
		SetHardwareConcurrency(16).
		ScreenColorDepth(24).
		SetScreenWidth(1900).
		SetScreenHeight(900).
		SetTimezone("America/Los_Angeles").
		SetLanguages("en-US,en").
		Build()
	if err != nil {
		panic(err)
	}

	params := wrapper.Params{
		Device:         deviceInst,
		Universe:       universe,
		Username:       username,
		Password:       password,
		Lang:           language,
		AutoLogin:      autoLogin,
		Proxy:          proxyAddr,
		ProxyUsername:  proxyUsername,
		ProxyPassword:  proxyPassword,
		ProxyType:      proxyType,
		ProxyLoginOnly: proxyLoginOnly,
		Lobby:          lobby,
		APINewHostname: apiNewHostname,
	}
	if njaApiKey != "" {
		params.CaptchaCallback = wrapper.NinjaSolver(njaApiKey)
	}

	bot, err := wrapper.NewWithParams(params)
	if err != nil {
		return err
	}

	e := echo.New()
	if corsEnabled {
		e.Use(middleware.CORS())
	}
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx.Set("bot", bot)
			ctx.Set("version", version)
			ctx.Set("commit", commit)
			ctx.Set("date", date)
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
	e.GET("/", wrapper.HomeHandler)
	e.GET("/tasks", wrapper.TasksHandler)

	// CAPTCHA Handler
	e.GET("/bot/captcha", wrapper.GetCaptchaHandler)
	e.POST("/bot/captcha/solve", wrapper.GetCaptchaSolverHandler)
	e.GET("/bot/captcha/challenge", wrapper.GetCaptchaChallengeHandler)

	e.GET("/bot/ip", wrapper.GetPublicIPHandler)
	e.GET("/bot/server", wrapper.GetServerHandler)
	e.GET("/bot/server-data", wrapper.GetServerDataHandler)
	e.POST("/bot/set-user-agent", wrapper.SetUserAgentHandler)
	e.GET("/bot/server-url", wrapper.ServerURLHandler)
	e.GET("/bot/language", wrapper.GetLanguageHandler)
	e.GET("/bot/empire/type/:typeID", wrapper.GetEmpireHandler)
	e.POST("/bot/page-content", wrapper.PageContentHandler)
	e.GET("/bot/login", wrapper.LoginHandler)
	e.GET("/bot/logout", wrapper.LogoutHandler)
	e.GET("/bot/username", wrapper.GetUsernameHandler)
	e.GET("/bot/universe-name", wrapper.GetUniverseNameHandler)
	e.GET("/bot/server/speed", wrapper.GetUniverseSpeedHandler)
	e.GET("/bot/server/speed-fleet", wrapper.GetUniverseSpeedFleetHandler)
	e.GET("/bot/server/version", wrapper.ServerVersionHandler)
	e.GET("/bot/server/time", wrapper.ServerTimeHandler)
	e.GET("/bot/is-under-attack", wrapper.IsUnderAttackHandler)
	e.GET("/bot/is-vacation-mode", wrapper.IsVacationModeHandler)
	e.GET("/bot/user-infos", wrapper.GetUserInfosHandler)
	e.GET("/bot/character-class", wrapper.GetCharacterClassHandler)
	e.GET("/bot/has-commander", wrapper.HasCommanderHandler)
	e.GET("/bot/has-admiral", wrapper.HasAdmiralHandler)
	e.GET("/bot/has-engineer", wrapper.HasEngineerHandler)
	e.GET("/bot/has-geologist", wrapper.HasGeologistHandler)
	e.GET("/bot/has-technocrat", wrapper.HasTechnocratHandler)
	e.POST("/bot/send-message", wrapper.SendMessageHandler)
	e.GET("/bot/fleets", wrapper.GetFleetsHandler)
	e.GET("/bot/fleets/slots", wrapper.GetSlotsHandler)
	e.POST("/bot/fleets/:fleetID/cancel", wrapper.CancelFleetHandler)
	e.GET("/bot/espionage-report/:msgid", wrapper.GetEspionageReportHandler)
	e.GET("/bot/espionage-report/:galaxy/:system/:position", wrapper.GetEspionageReportForHandler)
	e.GET("/bot/espionage-report", wrapper.GetEspionageReportMessagesHandler)
	e.POST("/bot/delete-report/:messageID", wrapper.DeleteMessageHandler)
	e.POST("/bot/delete-all-espionage-reports", wrapper.DeleteEspionageMessagesHandler)
	e.POST("/bot/delete-all-reports/:tabIndex", wrapper.DeleteMessagesFromTabHandler)
	e.GET("/bot/attacks", wrapper.GetAttacksHandler)
	e.GET("/bot/get-auction", wrapper.GetAuctionHandler)
	e.POST("/bot/do-auction", wrapper.DoAuctionHandler)
	e.GET("/bot/galaxy-infos/:galaxy/:system", wrapper.GalaxyInfosHandler)
	e.GET("/bot/get-research", wrapper.GetResearchHandler)
	e.GET("/bot/buy-offer-of-the-day", wrapper.BuyOfferOfTheDayHandler)
	e.GET("/bot/price/:ogameID/:nbr", wrapper.GetPriceHandler)
	e.GET("/bot/requirements/:ogameID", wrapper.GetRequirementsHandler)
	e.GET("/bot/moons", wrapper.GetMoonsHandler)
	e.GET("/bot/moons/:moonID", wrapper.GetMoonHandler)
	e.GET("/bot/moons/:galaxy/:system/:position", wrapper.GetMoonByCoordHandler)
	e.GET("/bot/celestials/:celestialID/items", wrapper.GetCelestialItemsHandler)
	e.GET("/bot/celestials/:celestialID/items/:itemRef/activate", wrapper.ActivateCelestialItemHandler)
	e.GET("/bot/celestials/:celestialID/techs", wrapper.TechsHandler)
	e.GET("/bot/planets", wrapper.GetPlanetsHandler)
	e.GET("/bot/planets/:planetID", wrapper.GetPlanetHandler)
	e.GET("/bot/planets/:galaxy/:system/:position", wrapper.GetPlanetByCoordHandler)
	e.GET("/bot/planets/:planetID/resources-details", wrapper.GetResourcesDetailsHandler)
	e.GET("/bot/planets/:planetID/resource-settings", wrapper.GetResourceSettingsHandler)
	e.POST("/bot/planets/:planetID/resource-settings", wrapper.SetResourceSettingsHandler)
	e.GET("/bot/planets/:planetID/resources-buildings", wrapper.GetResourcesBuildingsHandler)
	e.GET("/bot/planets/:planetID/lifeform-buildings", wrapper.GetLfBuildingsHandler)
	e.GET("/bot/planets/:planetID/lifeform-techs", wrapper.GetLfResearchHandler)
	e.GET("/bot/planets/:planetID/defence", wrapper.GetDefenseHandler)
	e.GET("/bot/planets/:planetID/ships", wrapper.GetShipsHandler)
	e.GET("/bot/planets/:planetID/facilities", wrapper.GetFacilitiesHandler)
	e.POST("/bot/planets/:planetID/build/:ogameID/:nbr", wrapper.BuildHandler)
	e.POST("/bot/planets/:planetID/build/cancelable/:ogameID", wrapper.BuildCancelableHandler)
	e.POST("/bot/planets/:planetID/build/production/:ogameID/:nbr", wrapper.BuildProductionHandler)
	e.POST("/bot/planets/:planetID/build/building/:ogameID", wrapper.BuildBuildingHandler)
	e.POST("/bot/planets/:planetID/build/technology/:ogameID", wrapper.BuildTechnologyHandler)
	e.POST("/bot/planets/:planetID/build/defence/:ogameID/:nbr", wrapper.BuildDefenseHandler)
	e.POST("/bot/planets/:planetID/build/ships/:ogameID/:nbr", wrapper.BuildShipsHandler)
	e.POST("/bot/planets/:planetID/teardown/:ogameID", wrapper.TeardownHandler)
	e.GET("/bot/planets/:planetID/production", wrapper.GetProductionHandler)
	e.GET("/bot/planets/:planetID/constructions", wrapper.ConstructionsBeingBuiltHandler)
	e.POST("/bot/planets/:planetID/cancel-building", wrapper.CancelBuildingHandler)
	e.POST("/bot/planets/:planetID/cancel-research", wrapper.CancelResearchHandler)
	e.GET("/bot/planets/:planetID/resources", wrapper.GetResourcesHandler)
	e.POST("/bot/planets/:planetID/send-fleet", wrapper.SendFleetHandler)
	e.POST("/bot/planets/:planetID/send-ipm", wrapper.SendIPMHandler)
	e.GET("/bot/moons/:moonID/phalanx/:galaxy/:system/:position", wrapper.PhalanxHandler)
	e.POST("/bot/moons/:moonID/jump-gate", wrapper.JumpGateHandler)
	e.GET("/game/allianceInfo.php", wrapper.GetAlliancePageContentHandler) // Example: //game/allianceInfo.php?allianceId=500127

	// Get/Post Page Content
	e.GET("/game/index.php", wrapper.GetFromGameHandler)
	e.POST("/game/index.php", wrapper.PostToGameHandler)

	// For AntiGame plugin
	// Static content
	e.GET("/cdn/*", wrapper.GetStaticHandler)
	e.GET("/assets/css/*", wrapper.GetStaticHandler)
	e.GET("/headerCache/*", wrapper.GetStaticHandler)
	e.GET("/favicon.ico", wrapper.GetStaticHandler)
	e.GET("/game/sw.js", wrapper.GetStaticHandler)

	// JSON API
	/*
		/api/serverData.xml
		/api/localization.xml
		/api/players.xml
		/api/universe.xml
	*/
	e.GET("/api/*", wrapper.GetStaticHandler)
	e.HEAD("/api/*", wrapper.GetStaticHEADHandler) // AntiGame uses this to check if the cached XML files need to be refreshed

	if enableTLS {
		log.Println("Enable TLS Support")
		return e.StartTLS(host+":"+strconv.Itoa(port), tlsCertFile, tlsKeyFile)
	}
	log.Println("Disable TLS Support")
	return e.Start(host + ":" + strconv.Itoa(port))
}
