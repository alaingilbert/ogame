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
			Name:    "socks5-proxy",
			Usage:   "Socks5 proxy address",
			Value:   "",
			EnvVars: []string{"OGAMED_SOCKS5_PROXY"},
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
	socks5ProxyAddr := c.String("socks5-proxy")
	proxyUsername := c.String("proxy-username")
	proxyPassword := c.String("proxy-password")
	lobby := c.String("lobby")
	apiNewHostname := c.String("api-new-hostname")
	enableTLS := c.Bool("enable-tls")
	tlsKeyFile := c.String("tls-key-file")
	tlsCertFile := c.String("tls-cert-file")
	basicAuthUsername := c.String("basic-auth-username")
	basicAuthPassword := c.String("basic-auth-password")

	bot, err := ogame.NewWithParams(ogame.Params{
		Universe:       universe,
		Username:       username,
		Password:       password,
		Lang:           language,
		AutoLogin:      autoLogin,
		Proxy:          proxyAddr,
		ProxyUsername:  proxyUsername,
		ProxyPassword:  proxyPassword,
		Socks5Address:  socks5ProxyAddr,
		Socks5Username: proxyUsername,
		Socks5Password: proxyPassword,
		Lobby:          lobby,
		APINewHostname: apiNewHostname,
	})
	if err != nil {
		return err
	}

	e := echo.New()
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
	e.GET("/bot/server", ogame.GetServerHandler)
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
	e.GET("/bot/user-infos", ogame.GetUserInfosHandler)
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
	e.GET("/bot/galaxy-infos/:galaxy/:system", ogame.GalaxyInfosHandler)
	e.GET("/bot/get-research", ogame.GetResearchHandler)
	e.GET("/bot/planets", ogame.GetPlanetsHandler)
	e.GET("/bot/planets/:planetID", ogame.GetPlanetHandler)
	e.GET("/bot/planets/:galaxy/:system/:position", ogame.GetPlanetByCoordHandler)
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
	e.GET("/bot/planets/:planetID/production", ogame.GetProductionHandler)
	e.GET("/bot/planets/:planetID/constructions", ogame.ConstructionsBeingBuiltHandler)
	e.POST("/bot/planets/:planetID/cancel-building", ogame.CancelBuildingHandler)
	e.POST("/bot/planets/:planetID/cancel-research", ogame.CancelResearchHandler)
	e.GET("/bot/planets/:planetID/resources", ogame.GetResourcesHandler)
	e.POST("/bot/planets/:planetID/send-fleet", ogame.SendFleetHandler)
	e.GET("/game/allianceInfo.php", ogame.GetAlliancePageContentHandler) // Example: //game/allianceInfo.php?allianceId=500127
	e.POST("/bot/planets/:planetID/send-ipm", ogame.SendIPMHandler)

	// Get/Post Page Content
	e.GET("/game/index.php", ogame.GetFromGameHandler)
	e.POST("/game/index.php", ogame.PostToGameHandler)

	// For AntiGame plugin
	// Static content
	e.GET("/cdn/*", ogame.GetStaticHandler)
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
	} else {
		log.Println("Disable TLS Support")
		return e.Start(host + ":" + strconv.Itoa(port))
	}
}
