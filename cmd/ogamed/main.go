package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/alaingilbert/ogame"
	"github.com/labstack/echo"
	"gopkg.in/urfave/cli.v2"
)

var bot *ogame.OGame

func main() {
	app := cli.App{}
	app.Authors = []*cli.Author{
		{"Alain Gilbert", "alain.gilbert.15@gmail.com"},
	}
	app.Name = "ogamed"
	app.Usage = "ogame deamon service"
	app.Version = "0.0.0"
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
	host := c.String("host")
	port := c.Int("port")
	var err error
	bot, err = ogame.New(universe, username, password, language)
	if err != nil {
		return err
	}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Debug = false
	e.GET("/", home)
	e.POST("/bot/set-user-agent", setUserAgent)
	e.GET("/bot/server-url", serverURL)
	e.POST("/bot/page-content", pageContent)
	e.GET("/bot/login", login)
	e.GET("/bot/logout", logout)
	e.GET("/bot/server/speed", getUniverseSpeed)
	e.GET("/bot/server/version", serverVersion)
	e.GET("/bot/server/time", serverTime)
	e.GET("/bot/is-under-attack", isUnderAttack)
	e.GET("/bot/user-infos", getUserInfos)
	e.POST("/bot/send-message", sendMessage)
	e.GET("/bot/fleets", getFleets)
	e.POST("/bot/fleets/:fleetID/cancel", cancelFleet)
	e.GET("/bot/attacks", getAttacks)
	e.GET("/bot/galaxy-infos/:galaxy/:system", galaxyInfos)
	e.GET("/bot/get-research", getResearch)
	e.GET("/bot/planets", getPlanets)
	e.GET("/bot/planets/:planetID", getPlanet)
	e.GET("/bot/planets/:galaxy/:system/:position", getPlanetByCoord)
	e.GET("/bot/planets/:planetID/resource-settings", getResourceSettings)
	e.POST("/bot/planets/:planetID/resource-settings", setResourceSettings)
	e.GET("/bot/planets/:planetID/resources-buildings", getResourcesBuildings)
	e.GET("/bot/planets/:planetID/defence", getDefense)
	e.GET("/bot/planets/:planetID/ships", getShips)
	e.GET("/bot/planets/:planetID/facilities", getFacilities)
	e.POST("/bot/planets/:planetID/build/:ogameID/:nbr", build)
	e.POST("/bot/planets/:planetID/build/cancelable/:ogameID", buildCancelable)
	e.POST("/bot/planets/:planetID/build/production/:ogameID/:nbr", buildProduction)
	e.POST("/bot/planets/:planetID/build/building/:ogameID", buildBuilding)
	e.POST("/bot/planets/:planetID/build/technology/:ogameID", buildTechnology)
	e.POST("/bot/planets/:planetID/build/defence/:ogameID/:nbr", buildDefense)
	e.POST("/bot/planets/:planetID/build/ships/:ogameID/:nbr", buildShips)
	e.GET("/bot/planets/:planetID/production", getProduction)
	e.GET("/bot/planets/:planetID/constructions", constructionsBeingBuilt)
	e.POST("/bot/planets/:planetID/cancel-building", cancelBuilding)
	e.POST("/bot/planets/:planetID/cancel-research", cancelResearch)
	e.GET("/bot/planets/:planetID/resources", getResources)
	e.POST("/bot/planets/:planetID/send-fleet", sendFleet)
	return e.Start(host + ":" + strconv.Itoa(port))
}

type apiResp struct {
	Status  string
	Code    int
	Message string
	Result  interface{}
}

func successResp(data interface{}) apiResp {
	return apiResp{Status: "ok", Code: 200, Result: data}
}

func errorResp(code int, message string) apiResp {
	return apiResp{Status: "error", Code: code, Message: message}
}

func home(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

// curl 127.0.0.1:1234/bot/set-user-agent -d 'userAgent="New user agent"'
func setUserAgent(c echo.Context) error {
	userAgent := c.Request().PostFormValue("userAgent")
	bot.SetUserAgent(userAgent)
	return c.JSON(http.StatusOK, successResp(nil))
}

func serverURL(c echo.Context) error {
	return c.JSON(http.StatusOK, successResp(bot.ServerURL()))
}

// curl 127.0.0.1:1234/bot/page-content -d 'page=overview&cp=123'
func pageContent(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, err.Error()))
	}
	return c.JSON(http.StatusOK, successResp(bot.GetPageContent(c.Request().Form)))
}

func login(c echo.Context) error {
	if err := bot.Login(); err != nil {
		if err == ogame.ErrBadCredentials {
			return c.JSON(http.StatusBadRequest, errorResp(400, err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, errorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, successResp(nil))
}

func logout(c echo.Context) error {
	bot.Logout()
	return c.JSON(http.StatusOK, successResp(nil))
}

func getUniverseSpeed(c echo.Context) error {
	return c.JSON(http.StatusOK, successResp(bot.GetUniverseSpeed()))
}

func serverVersion(c echo.Context) error {
	return c.JSON(http.StatusOK, successResp(bot.ServerVersion()))
}

func serverTime(c echo.Context) error {
	return c.JSON(http.StatusOK, successResp(bot.ServerTime()))
}

func isUnderAttack(c echo.Context) error {
	return c.JSON(http.StatusOK, successResp(bot.IsUnderAttack()))
}

func getUserInfos(c echo.Context) error {
	return c.JSON(http.StatusOK, successResp(bot.GetUserInfos()))
}

// curl 127.0.0.1:1234/bot/send-message -d 'playerID=123&message="Sup boi!"'
func sendMessage(c echo.Context) error {
	playerID, err := strconv.Atoi(c.Request().PostFormValue("playerID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, err.Error()))
	}
	message := c.Request().PostFormValue("message")
	if err := bot.SendMessage(playerID, message); err != nil {
		if err.Error() == "invalid parameters" {
			return c.JSON(http.StatusBadRequest, errorResp(400, err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, errorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, successResp(nil))
}

func getFleets(c echo.Context) error {
	return c.JSON(http.StatusOK, successResp(bot.GetFleets()))
}

func cancelFleet(c echo.Context) error {
	fleetID, err := strconv.Atoi(c.Param("fleetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, err.Error()))
	}
	return c.JSON(http.StatusOK, successResp(bot.CancelFleet(ogame.FleetID(fleetID))))
}

func getAttacks(c echo.Context) error {
	return c.JSON(http.StatusOK, successResp(bot.GetAttacks()))
}

func galaxyInfos(c echo.Context) error {
	galaxy, err := strconv.Atoi(c.Param("galaxy"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, err.Error()))
	}
	system, err := strconv.Atoi(c.Param("system"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, err.Error()))
	}
	res, err := bot.GalaxyInfos(galaxy, system)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, successResp(res))
}

func getResearch(c echo.Context) error {
	return c.JSON(http.StatusOK, successResp(bot.GetResearch()))
}

func getPlanets(c echo.Context) error {
	return c.JSON(http.StatusOK, successResp(bot.GetPlanets()))
}

func getPlanet(c echo.Context) error {
	planetID, err := strconv.Atoi(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid planet id"))
	}
	planet, err := bot.GetPlanet(ogame.PlanetID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, successResp(planet))
}

func getPlanetByCoord(c echo.Context) error {
	galaxy, err := strconv.Atoi(c.Param("galaxy"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid galaxy"))
	}
	system, err := strconv.Atoi(c.Param("system"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid system"))
	}
	position, err := strconv.Atoi(c.Param("position"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid position"))
	}
	planet, err := bot.GetPlanetByCoord(ogame.Coordinate{Galaxy: galaxy, System: system, Position: position})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, successResp(planet))
}

func getResourceSettings(c echo.Context) error {
	planetID, err := strconv.Atoi(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid planet id"))
	}
	res, err := bot.GetResourceSettings(ogame.PlanetID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, successResp(res))
}

// curl 127.0.0.1:1234/bot/planets/123/resource-settings -d 'metalMine=100&crystalMine=100&deuteriumSynthesizer=100&solarPlant=100&fusionReactor=100&solarSatellite=100'
func setResourceSettings(c echo.Context) error {
	planetID, err := strconv.Atoi(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid planet id"))
	}
	metalMine, err := strconv.Atoi(c.Request().PostFormValue("metalMine"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid metalMine"))
	}
	crystalMine, err := strconv.Atoi(c.Request().PostFormValue("crystalMine"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid crystalMine"))
	}
	deuteriumSynthesizer, err := strconv.Atoi(c.Request().PostFormValue("deuteriumSynthesizer"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid deuteriumSynthesizer"))
	}
	solarPlant, err := strconv.Atoi(c.Request().PostFormValue("solarPlant"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid solarPlant"))
	}
	fusionReactor, err := strconv.Atoi(c.Request().PostFormValue("fusionReactor"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid fusionReactor"))
	}
	solarSatellite, err := strconv.Atoi(c.Request().PostFormValue("solarSatellite"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid solarSatellite"))
	}
	settings := ogame.ResourceSettings{
		MetalMine:            metalMine,
		CrystalMine:          crystalMine,
		DeuteriumSynthesizer: deuteriumSynthesizer,
		SolarPlant:           solarPlant,
		FusionReactor:        fusionReactor,
		SolarSatellite:       solarSatellite,
	}
	if err := bot.SetResourceSettings(ogame.PlanetID(planetID), settings); err != nil {
		if err == ogame.ErrInvalidPlanetID {
			return c.JSON(http.StatusBadRequest, errorResp(400, err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, errorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, successResp(nil))
}

func getResourcesBuildings(c echo.Context) error {
	planetID, err := strconv.Atoi(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid planet id"))
	}
	res, err := bot.GetResourcesBuildings(ogame.PlanetID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, successResp(res))
}

func getDefense(c echo.Context) error {
	planetID, err := strconv.Atoi(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid planet id"))
	}
	res, err := bot.GetDefense(ogame.PlanetID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, successResp(res))
}

func getShips(c echo.Context) error {
	planetID, err := strconv.Atoi(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid planet id"))
	}
	res, err := bot.GetShips(ogame.PlanetID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, successResp(res))
}

func getFacilities(c echo.Context) error {
	planetID, err := strconv.Atoi(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid planet id"))
	}
	res, err := bot.GetFacilities(ogame.PlanetID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, successResp(res))
}

func build(c echo.Context) error {
	planetID, err := strconv.Atoi(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid planet id"))
	}
	ogameID, err := strconv.Atoi(c.Param("ogameID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid ogame id"))
	}
	nbr, err := strconv.Atoi(c.Param("nbr"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid nbr"))
	}
	if err := bot.Build(ogame.PlanetID(planetID), ogame.ID(ogameID), nbr); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, successResp(nil))
}

func buildCancelable(c echo.Context) error {
	planetID, err := strconv.Atoi(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid planet id"))
	}
	ogameID, err := strconv.Atoi(c.Param("ogameID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid ogame id"))
	}
	if err := bot.BuildCancelable(ogame.PlanetID(planetID), ogame.ID(ogameID)); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, successResp(nil))
}

func buildProduction(c echo.Context) error {
	planetID, err := strconv.Atoi(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid planet id"))
	}
	ogameID, err := strconv.Atoi(c.Param("ogameID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid ogame id"))
	}
	nbr, err := strconv.Atoi(c.Param("nbr"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid nbr"))
	}
	if err := bot.BuildProduction(ogame.PlanetID(planetID), ogame.ID(ogameID), nbr); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, successResp(nil))
}

func buildBuilding(c echo.Context) error {
	planetID, err := strconv.Atoi(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid planet id"))
	}
	ogameID, err := strconv.Atoi(c.Param("ogameID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid ogame id"))
	}
	if err := bot.BuildBuilding(ogame.PlanetID(planetID), ogame.ID(ogameID)); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, successResp(nil))
}

func buildTechnology(c echo.Context) error {
	planetID, err := strconv.Atoi(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid planet id"))
	}
	ogameID, err := strconv.Atoi(c.Param("ogameID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid ogame id"))
	}
	if err := bot.BuildTechnology(ogame.PlanetID(planetID), ogame.ID(ogameID)); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, successResp(nil))
}

func buildDefense(c echo.Context) error {
	planetID, err := strconv.Atoi(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid planet id"))
	}
	ogameID, err := strconv.Atoi(c.Param("ogameID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid ogame id"))
	}
	nbr, err := strconv.Atoi(c.Param("nbr"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid nbr"))
	}
	if err := bot.BuildDefense(ogame.PlanetID(planetID), ogame.ID(ogameID), nbr); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, successResp(nil))
}

func buildShips(c echo.Context) error {
	planetID, err := strconv.Atoi(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid planet id"))
	}
	ogameID, err := strconv.Atoi(c.Param("ogameID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid ogame id"))
	}
	nbr, err := strconv.Atoi(c.Param("nbr"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid nbr"))
	}
	if err := bot.BuildShips(ogame.PlanetID(planetID), ogame.ID(ogameID), nbr); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, successResp(nil))
}

func getProduction(c echo.Context) error {
	planetID, err := strconv.Atoi(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid planet id"))
	}
	res, err := bot.GetProduction(ogame.PlanetID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, successResp(res))
}

func constructionsBeingBuilt(c echo.Context) error {
	planetID, err := strconv.Atoi(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid planet id"))
	}
	buildingID, buildingCountdown, researchID, researchCountdown := bot.ConstructionsBeingBuilt(ogame.PlanetID(planetID))
	return c.JSON(http.StatusOK, successResp(
		struct {
			BuildingID        int
			BuildingCountdown int
			ResearchID        int
			ResearchCountdown int
		}{
			BuildingID:        int(buildingID),
			BuildingCountdown: buildingCountdown,
			ResearchID:        int(researchID),
			ResearchCountdown: researchCountdown,
		},
	))
}

func cancelBuilding(c echo.Context) error {
	planetID, err := strconv.Atoi(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid planet id"))
	}
	if err := bot.CancelBuilding(ogame.PlanetID(planetID)); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, successResp(nil))
}

func cancelResearch(c echo.Context) error {
	planetID, err := strconv.Atoi(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid planet id"))
	}
	if err := bot.CancelResearch(ogame.PlanetID(planetID)); err != nil {
		return c.JSON(http.StatusInternalServerError, errorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, successResp(nil))
}

func getResources(c echo.Context) error {
	planetID, err := strconv.Atoi(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid planet id"))
	}
	res, err := bot.GetResources(ogame.PlanetID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, successResp(res))
}

// curl 127.0.0.1:1234/bot/planets/123/send-fleet -d 'ships="203,1"&ships="204,10"&speed=10&galaxy=1&system=1&position=1&mission=3&metal=1&crystal=2&deuterium=3'
func sendFleet(c echo.Context) error {
	planetID, err := strconv.Atoi(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(400, "invalid planet id"))
	}

	c.Request().ParseForm()

	var ships []ogame.Quantifiable
	where := ogame.Coordinate{}
	mission := ogame.Transport
	payload := ogame.Resources{}
	speed := ogame.HundredPercent
	for key, values := range c.Request().PostForm {
		switch key {
		case "ships":
			for _, s := range values {
				a := strings.Split(s, ",")
				shipID, err := strconv.Atoi(a[0])
				if err != nil || !ogame.IsShipID(shipID) {
					return c.JSON(http.StatusBadRequest, errorResp(400, "invalid ship id "+a[0]))
				}
				nbr, err := strconv.Atoi(a[1])
				if err != nil || nbr < 0 {
					return c.JSON(http.StatusBadRequest, errorResp(400, "invalid nbr "+a[1]))
				}
				ships = append(ships, ogame.Quantifiable{ID: ogame.ID(shipID), Nbr: nbr})
			}
		case "speed":
			speedInt, err := strconv.Atoi(values[0])
			if err != nil || speedInt < 0 || speedInt > 10 {
				return c.JSON(http.StatusBadRequest, errorResp(400, "invalid speed"))
			}
			speed = ogame.Speed(speedInt)
		case "galaxy":
			galaxy, err := strconv.Atoi(values[0])
			if err != nil {
				return c.JSON(http.StatusBadRequest, errorResp(400, "invalid galaxy"))
			}
			where.Galaxy = galaxy
		case "system":
			system, err := strconv.Atoi(values[0])
			if err != nil {
				return c.JSON(http.StatusBadRequest, errorResp(400, "invalid system"))
			}
			where.System = system
		case "position":
			position, err := strconv.Atoi(values[0])
			if err != nil {
				return c.JSON(http.StatusBadRequest, errorResp(400, "invalid position"))
			}
			where.Position = position
		case "mission":
			missionInt, err := strconv.Atoi(values[0])
			if err != nil {
				return c.JSON(http.StatusBadRequest, errorResp(400, "invalid mission"))
			}
			mission = ogame.MissionID(missionInt)
		case "metal":
			metal, err := strconv.Atoi(values[0])
			if err != nil || metal < 0 {
				return c.JSON(http.StatusBadRequest, errorResp(400, "invalid metal"))
			}
			payload.Metal = metal
		case "crystal":
			crystal, err := strconv.Atoi(values[0])
			if err != nil || crystal < 0 {
				return c.JSON(http.StatusBadRequest, errorResp(400, "invalid crystal"))
			}
			payload.Crystal = crystal
		case "deuterium":
			deuterium, err := strconv.Atoi(values[0])
			if err != nil || deuterium < 0 {
				return c.JSON(http.StatusBadRequest, errorResp(400, "invalid deuterium"))
			}
			payload.Deuterium = deuterium
		}
	}

	fleetID, err := bot.SendFleet(ogame.PlanetID(planetID), ships, speed, where, mission, payload)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, successResp(fleetID))
}
