package main

import (
	"crypto/subtle"
	"log"
	"os"
	"strconv"

	"github.com/alaingilbert/ogame"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/urfave/cli.v2"
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
			EnvVars: []string{"OGAMED_PROXY_PASSWORD"},
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
	cookiesFilename := c.String("cookies-filename")
	corsEnabled := c.Bool("cors-enabled")
	njaApiKey := c.String("nja-api-key")

	params := ogame.Params{
		Universe:        universe,
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

	bot, err := ogame.NewWithParams(params)
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
	e.POST("/bot/planets/:planetID/send-ipm", ogame.SendIPMHandler)
	e.GET("/bot/moons/:moonID/phalanx/:galaxy/:system/:position", ogame.PhalanxHandler)
	e.POST("/bot/moons/:moonID/jump-gate", ogame.JumpGateHandler)
	e.GET("/game/allianceInfo.php", ogame.GetAlliancePageContentHandler) // Example: //game/allianceInfo.php?allianceId=500127

	// Get/Post Page Content
	e.GET("/game/index.php", ogame.GetFromGameHandler)
	e.POST("/game/index.php", ogame.PostToGameHandler)

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

	if enableTLS {
		log.Println("Enable TLS Support")
		return e.StartTLS(host+":"+strconv.Itoa(port), tlsCertFile, tlsKeyFile)
	}
	log.Println("Disable TLS Support")
	return e.Start(host + ":" + strconv.Itoa(port))
}
