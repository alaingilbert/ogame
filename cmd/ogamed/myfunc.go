package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	"github.com/alaingilbert/ogame"
	"github.com/labstack/echo"
)

// Get Data for Planet View
func getPlanetView(c echo.Context) error {
	bot := c.Get("bot").(*ogame.OGame)
	allianceID := c.QueryParam("allianceId")
	vals := url.Values{"allianceId": {allianceID}}

	p := bot.GetCachedPlanets()

	fmt.Printf("Coordinates: %s", p[0].Coordinate)

	return c.HTML(http.StatusOK, string(vals.Encode()))
}

var t *Template

// Template for HTML
type Template struct {
	templates *template.Template
}

var templateFuncs = template.FuncMap{
	"add": add,
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

func add(x int64, y int64) int64 {
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

type extPlanet struct {
	Planet             ogame.Planet
	ResourcesDetails   ogame.ResourcesDetails
	ResourcesBuildings ogame.ResourcesBuildings
	Facilities         ogame.Facilities
	Technology         ogame.Researches
	DefensesInfos      ogame.DefensesInfos
	ShipsInfos         ogame.ShipsInfos
	Buildings          []ogame.Building
	Ships              []ogame.Ship
	Defenses           []ogame.Defense
	Techs              []ogame.Technology
	ConstructionQueue  []ogame.Quantifiable
}

// All Planets
var planets []extPlanet

func initial(bot *ogame.OGame) {
	planets = make([]extPlanet, 1)

	for _, p := range bot.GetPlanets() {
		ResourcesBuildings, _ := p.GetResourcesBuildings()
		Facilities, _ := p.GetFacilities()
		ResourcesDetails, _ := p.GetResourcesDetails()
		ShipsInfos, _ := p.GetShips()
		DefensesInfos, _ := p.GetDefense()
		Technology := bot.GetResearch()
		var ConstructionQueue []ogame.Quantifiable

		//planets[0].Buildings[0].GetID()
		//planets[0].Buildings[0].GetLevel

		//planets[0].ResourcesBuildings.ByID()
		//planets[0].Buildings[0].GetPrice()

		planets = append(planets, extPlanet{p, ResourcesDetails, ResourcesBuildings, Facilities, Technology, DefensesInfos, ShipsInfos, ogame.Buildings, ogame.Ships, ogame.Defenses, ogame.Technologies, ConstructionQueue})

	}
}

func htmlPlanetView(c echo.Context) error {
	//bot := c.Get("bot").(*ogame.OGame)
	planetID, _ := strconv.Atoi(c.Param("planetID"))

	planet := ogame.PlanetID(planetID)

	var selectedplanet extPlanet

	for i := 0; i < len(planets); i++ {
		if planets[i].Planet.ID == planet {
			selectedplanet = planets[i]
			id, _ := strconv.Atoi(c.QueryParam("id"))
			nbr, _ := strconv.Atoi(c.QueryParam("nbr"))
			if len(c.QueryParam("id")) > 0 {
				planets[i].ConstructionQueue = append(planets[i].ConstructionQueue, ogame.Quantifiable{ogame.ID(id), nbr})
			}
		}
	}

	return c.Render(http.StatusOK, "planetview", selectedplanet)
}
