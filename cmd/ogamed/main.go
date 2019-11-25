package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	"github.com/alaingilbert/ogame"
	"github.com/labstack/echo"
	cli "gopkg.in/urfave/cli.v2"
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
		&cli.StringFlag{
			Name:    "proxy",
			Usage:   "Proxy Url",
			Value:   "",
			EnvVars: []string{"OGAMED_PROXY"},
		},
	}
	app.Action = start
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

var port int

func start(c *cli.Context) error {
	universe := c.String("universe")
	username := c.String("username")
	password := c.String("password")
	language := c.String("language")
	host := c.String("host")
	port = c.Int("port")
	proxy := c.String("proxy")

	params := ogame.Params{
		Universe:  universe,
		Username:  username,
		Password:  password,
		Lang:      language,
		AutoLogin: true,
		Proxy:     proxy,
	}

	//bot, err := ogame.New(universe, username, password, language)
	bot, err := ogame.NewWithParams(params)
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
	e.HideBanner = true
	e.HidePort = true
	e.Debug = true

	///////////////////////////////////
	tmp, _ := template.New("").Funcs(templateFuncs).ParseGlob("templates/*.html")

	t := &Template{
		templates: tmp,
	}
	e.Renderer = t
	///////////////////////////////////

	e.GET("/", ogame.HomeHandler)
	e.GET("/bot/server", ogame.GetServerHandler)
	e.POST("/bot/set-user-agent", ogame.SetUserAgentHandler)
	e.GET("/bot/server-url", ogame.ServerURLHandler)
	e.GET("/bot/language", ogame.GetLanguageHandler)
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

	e.GET("/game/index.php", getFromGame)
	e.POST("/game/index.php", postToGame)
	e.GET("/game/allianceInfo.php", getAlliancePageContent)
	e.GET("/api/*", getStatic)
	e.GET("/cdn/*", getStatic)
	e.GET("/headerCache/*", getStatic)
	e.GET("/favicon.ico", getStatic)
	e.GET("/game/sw.js", getStatic)

	return e.Start(host + ":" + strconv.Itoa(port))
}

////////////////////////////////////////////////////////////////////////////////////////////////

var t *Template

// Template for HTML
type Template struct {
	templates *template.Template
}

var templateFuncs = template.FuncMap{
	"rangeStruct": RangeStructer,
	"add":         add,
	"fleetsave":   fleetsave,
}

// Render for Templates
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// RangeStructer takes the first argument, which must be a struct, and
// returns the value of each field in a slice. It will return nil
// if there are no arguments or first argument is not a struct
func RangeStructer(args ...interface{}) []interface{} {
	if len(args) == 0 {
		return nil
	}

	v := reflect.ValueOf(args[0])
	if v.Kind() != reflect.Struct {
		return nil
	}

	out := make([]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		out[i] = v.Field(i).Interface()
	}

	return out
}

func add(x int, y int) int {
	return x + y
}

// EnableFleetsaver Exported Variable Variable to Enable FleetSaver
var EnableFleetsaver = false

func fleetsave() bool {
	return EnableFleetsaver
}

// Get OGame Website
func getFromGame(c echo.Context) error {
	bot := c.Get("bot").(*ogame.OGame)
	//return c.JSON(http.StatusOK, ogame.SuccessResp(bot.GetServer()))

	vals := url.Values{"page": {"overview"}}

	if len(c.QueryParams()) > 0 {
		vals = c.QueryParams()
	}

	localserverurl := strings.Replace(bot.ServerURL(), "https://", "", -1)

	// log.Print(c.Request())
	// log.Print(c.Request().Header)
	// log.Print(c.Request().URL)
	// log.Print(c.Request().Host)
	// log.Print(c.Request().RequestURI)
	// log.Print(bot.ServerURL())

	byteArray := bot.GetPageContent(vals)

	//bytes.Replace(byteArray, []byte("s107-nl.ogame.gameforge.com"), []byte("localhost:4567"), -1)

	// Replace "s107-nl.ogame.gameforge.com" with "gemini.example.com"
	html := string(byteArray)
	/*
		unishort := strings.Split(localserverurl, ".")
		html = strings.Replace(html, "<meta name=\"ogame-universe\" content=\""+localserverurl+"\"/>", "<meta name=\"ogame-universe\" content=\""+unishort[0]+"\"/>", -1)
	*/
	html = strings.Replace(html, localserverurl, c.Request().Host, -1)
	html = strings.Replace(html, "<meta name=\"ogame-universe\" content=\""+c.Request().Host+"\"/>", "<meta name=\"ogame-universe\" content=\""+strings.Replace(bot.ServerURL(), "https://", "", -1)+"\"/>", -1)
	html = strings.Replace(html, "https", "http", -1)

	//html = strings.Replace(html, "s107-nl.ogame.gameforge.com", "127.0.0.1:8080", -1)
	//	html = strings.Replace(html, "\"/cdn", "\"https://s107-nl.ogame.gameforge.com/cdn", -1)

	// Todo: 15jan2019:
	// https://gf1.geo.gfsrv.net/cdn
	// https://gf2.geo.gfsrv.net/cdn
	// https://gf3.geo.gfsrv.net/cdn
	// Nieuwe Nginx URL aanmaken met caching

	return c.HTML(http.StatusOK, html)
}

// Post OGame Website
func postToGame(c echo.Context) error {
	bot := c.Get("bot").(*ogame.OGame)
	var localserverurl = strings.Replace(bot.ServerURL(), "https://", "", -1)
	vals := url.Values{"page": {"overview"}}

	if len(c.QueryParams()) > 0 {
		vals = c.QueryParams()
	}

	// Payload
	payload, _ := c.FormParams()

	// log.Println(c.Request())
	// log.Println(c.Request().Header)
	// log.Println(c.Request().URL)
	// log.Println(c.Request().Host)
	// log.Println(c.Request().RequestURI)
	// log.Println(bot.ServerURL())
	//log.Println(payload.Encode())
	//log.Println(vals.Encode())

	// Perform the post to the library
	byteArray := bot.PostPageContent(vals, payload)

	// Replace "s107-nl.ogame.gameforge.com" with "gemini.example.com"
	html := string(byteArray)
	html = strings.Replace(html, localserverurl, c.Request().Host, -1)
	html = strings.Replace(html, "https", "http", -1)
	return c.HTML(http.StatusOK, html)
}

// GetStatic Elements
func getStatic(c echo.Context) error {
	bot := c.Get("bot").(*ogame.OGame)
	var localserverurl = strings.Replace(bot.ServerURL(), "https://", "", -1)
	url := bot.ServerURL() + c.Request().URL.String()

	if len(c.QueryParams()) > 0 {
		url = url + "?" + c.QueryParams().Encode()
	}

	resp, err := bot.Client.Get(url)

	//resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	contentType := string(http.DetectContentType(body))

	if strings.Contains(url, ".css") {
		contentType = "text/css"
	}

	if strings.HasSuffix(url, ".js") {
		contentType = "text/javascript"
	}

	//if strings.Contains(c.Request().URL.String(), "localization.xml") || strings.Contains(c.Request().URL.String(), "serverData.xml") {

	if strings.Contains(c.Request().URL.String(), ".xml") {
		//log.Println(c.Request().URL.String())
		body2 := strings.Replace(string(body), localserverurl, c.Request().Host, -1)
		body2 = strings.Replace(string(body2), "https", "http", -1)
		return c.Blob(http.StatusOK, "application/xml", []byte(body2))
	}

	return c.Blob(http.StatusOK, contentType, body)
}

// Get Alliance Page
func getAlliancePageContent(c echo.Context) error {
	//bot := c.Get("bot").(*ogame.OGame)
	allianceID := c.QueryParam("allianceId")
	vals := url.Values{"allianceId": {allianceID}}

	// ogame.go:
	// func (b *OGame) getAlliancePageContent(vals url.Values) ([]byte, error) {
	// finalURL := b.serverURL + "/game/allianceInfo.php?" + vals.Encode()

	return c.HTML(http.StatusOK, string(vals.Encode()))
}
