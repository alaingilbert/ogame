package ogame

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/labstack/echo"
)

// APIResp ...
type APIResp struct {
	Status  string
	Code    int
	Message string
	Result  interface{}
}

// SuccessResp ...
func SuccessResp(data interface{}) APIResp {
	return APIResp{Status: "ok", Code: 200, Result: data}
}

// ErrorResp ...
func ErrorResp(code int, message string) APIResp {
	return APIResp{Status: "error", Code: code, Message: message}
}

// HomeHandler ...
func HomeHandler(c echo.Context) error {
	version := c.Get("version").(string)
	commit := c.Get("commit").(string)
	date := c.Get("date").(string)
	return c.JSON(http.StatusOK, map[string]interface{}{
		"version": version,
		"commit":  commit,
		"date":    date,
	})
}

// GetServerHandler ...
func GetServerHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	return c.JSON(http.StatusOK, SuccessResp(bot.GetServer()))
}

// SetUserAgentHandler ...
// curl 127.0.0.1:1234/bot/set-user-agent -d 'userAgent="New user agent"'
func SetUserAgentHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	userAgent := c.Request().PostFormValue("userAgent")
	bot.SetUserAgent(userAgent)
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// ServerURLHandler ...
func ServerURLHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	return c.JSON(http.StatusOK, SuccessResp(bot.ServerURL()))
}

// GetLanguageHandler ...
func GetLanguageHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	return c.JSON(http.StatusOK, SuccessResp(bot.GetLanguage()))
}

// PageContentHandler ...
// curl 127.0.0.1:1234/bot/page-content -d 'page=overview&cp=123'
func PageContentHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	if err := c.Request().ParseForm(); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(bot.GetPageContent(c.Request().Form)))
}

// LoginHandler ...
func LoginHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	if err := bot.Login(); err != nil {
		if err == ErrBadCredentials {
			return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// LogoutHandler ...
func LogoutHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	bot.Logout()
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// GetUsernameHandler ...
func GetUsernameHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	return c.JSON(http.StatusOK, SuccessResp(bot.GetUsername()))
}

// GetUniverseNameHandler ...
func GetUniverseNameHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	return c.JSON(http.StatusOK, SuccessResp(bot.GetUniverseName()))
}

// GetUniverseSpeedHandler ...
func GetUniverseSpeedHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	return c.JSON(http.StatusOK, SuccessResp(bot.serverData.Speed))
}

// GetUniverseSpeedFleetHandler ...
func GetUniverseSpeedFleetHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	return c.JSON(http.StatusOK, SuccessResp(bot.serverData.SpeedFleet))
}

// ServerVersionHandler ...
func ServerVersionHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	return c.JSON(http.StatusOK, SuccessResp(bot.serverData.Version))
}

// ServerTimeHandler ...
func ServerTimeHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	return c.JSON(http.StatusOK, SuccessResp(bot.ServerTime()))
}

// IsUnderAttackHandler ...
func IsUnderAttackHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	isUnderAttack, err := bot.IsUnderAttack()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(isUnderAttack))
}

// GetUserInfosHandler ...
func GetUserInfosHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	return c.JSON(http.StatusOK, SuccessResp(bot.GetUserInfos()))
}

// GetEspionageReportMessagesHandler ...
func GetEspionageReportMessagesHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	report, err := bot.GetEspionageReportMessages()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(report))
}

// GetEspionageReportHandler ...
func GetEspionageReportHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	msgID, err := strconv.ParseInt(c.Param("msgid"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid msgid id"))
	}
	espionageReport, err := bot.GetEspionageReport(msgID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(espionageReport))
}

// GetEspionageReportForHandler ...
func GetEspionageReportForHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	galaxy, err := strconv.ParseInt(c.Param("galaxy"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid galaxy"))
	}
	system, err := strconv.ParseInt(c.Param("system"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid system"))
	}
	position, err := strconv.ParseInt(c.Param("position"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid position"))
	}
	planet, err := bot.GetEspionageReportFor(Coordinate{Type: PlanetType, Galaxy: galaxy, System: system, Position: position})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(planet))
}

// SendMessageHandler ...
// curl 127.0.0.1:1234/bot/send-message -d 'playerID=123&message="Sup boi!"'
func SendMessageHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	playerID, err := strconv.ParseInt(c.Request().PostFormValue("playerID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	message := c.Request().PostFormValue("message")
	if err := bot.SendMessage(playerID, message); err != nil {
		if err.Error() == "invalid parameters" {
			return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// GetFleetsHandler ...
func GetFleetsHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	fleets, _ := bot.GetFleets()
	return c.JSON(http.StatusOK, SuccessResp(fleets))
}

// GetSlotsHandler ...
func GetSlotsHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	slots := bot.GetSlots()
	return c.JSON(http.StatusOK, SuccessResp(slots))
}

// CancelFleetHandler ...
func CancelFleetHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	fleetID, err := strconv.ParseInt(c.Param("fleetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(bot.CancelFleet(FleetID(fleetID))))
}

// GetAttacksHandler ...
func GetAttacksHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	attacks, err := bot.GetAttacks()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(attacks))
}

// GalaxyInfosHandler ...
func GalaxyInfosHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	galaxy, err := strconv.ParseInt(c.Param("galaxy"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	system, err := strconv.ParseInt(c.Param("system"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	res, err := bot.GalaxyInfos(galaxy, system)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(res))
}

// GetResearchHandler ...
func GetResearchHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	return c.JSON(http.StatusOK, SuccessResp(bot.GetResearch()))
}

// GetPlanetsHandler ...
func GetPlanetsHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	return c.JSON(http.StatusOK, SuccessResp(bot.GetPlanets()))
}

// GetPlanetHandler ...
func GetPlanetHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	planet, err := bot.GetPlanet(PlanetID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(planet))
}

// GetPlanetByCoordHandler ...
func GetPlanetByCoordHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	galaxy, err := strconv.ParseInt(c.Param("galaxy"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid galaxy"))
	}
	system, err := strconv.ParseInt(c.Param("system"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid system"))
	}
	position, err := strconv.ParseInt(c.Param("position"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid position"))
	}
	planet, err := bot.GetPlanet(Coordinate{Type: PlanetType, Galaxy: galaxy, System: system, Position: position})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(planet))
}

// GetResourceSettingsHandler ...
func GetResourceSettingsHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	res, err := bot.GetResourceSettings(PlanetID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(res))
}

// SetResourceSettingsHandler ...
// curl 127.0.0.1:1234/bot/planets/123/resource-settings -d 'metalMine=100&crystalMine=100&deuteriumSynthesizer=100&solarPlant=100&fusionReactor=100&solarSatellite=100'
func SetResourceSettingsHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	metalMine, err := strconv.ParseInt(c.Request().PostFormValue("metalMine"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid metalMine"))
	}
	crystalMine, err := strconv.ParseInt(c.Request().PostFormValue("crystalMine"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid crystalMine"))
	}
	deuteriumSynthesizer, err := strconv.ParseInt(c.Request().PostFormValue("deuteriumSynthesizer"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid deuteriumSynthesizer"))
	}
	solarPlant, err := strconv.ParseInt(c.Request().PostFormValue("solarPlant"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid solarPlant"))
	}
	fusionReactor, err := strconv.ParseInt(c.Request().PostFormValue("fusionReactor"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid fusionReactor"))
	}
	solarSatellite, err := strconv.ParseInt(c.Request().PostFormValue("solarSatellite"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid solarSatellite"))
	}
	crawler, err := strconv.ParseInt(c.Request().PostFormValue("crawler"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid crawler"))
	}
	settings := ResourceSettings{
		MetalMine:            metalMine,
		CrystalMine:          crystalMine,
		DeuteriumSynthesizer: deuteriumSynthesizer,
		SolarPlant:           solarPlant,
		FusionReactor:        fusionReactor,
		SolarSatellite:       solarSatellite,
		Crawler:              crawler,
	}
	if err := bot.SetResourceSettings(PlanetID(planetID), settings); err != nil {
		if err == ErrInvalidPlanetID {
			return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// GetResourcesBuildingsHandler ...
func GetResourcesBuildingsHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	res, err := bot.GetResourcesBuildings(CelestialID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(res))
}

// GetDefenseHandler ...
func GetDefenseHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	res, err := bot.GetDefense(CelestialID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(res))
}

// GetShipsHandler ...
func GetShipsHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	res, err := bot.GetShips(CelestialID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(res))
}

// GetFacilitiesHandler ...
func GetFacilitiesHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	res, err := bot.GetFacilities(CelestialID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(res))
}

// BuildHandler ...
func BuildHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	ogameID, err := strconv.ParseInt(c.Param("ogameID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ogame id"))
	}
	nbr, err := strconv.ParseInt(c.Param("nbr"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid nbr"))
	}
	if err := bot.Build(CelestialID(planetID), ID(ogameID), nbr); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// BuildCancelableHandler ...
func BuildCancelableHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	ogameID, err := strconv.ParseInt(c.Param("ogameID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ogame id"))
	}
	if err := bot.BuildCancelable(CelestialID(planetID), ID(ogameID)); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// BuildProductionHandler ...
func BuildProductionHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	ogameID, err := strconv.ParseInt(c.Param("ogameID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ogame id"))
	}
	nbr, err := strconv.ParseInt(c.Param("nbr"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid nbr"))
	}
	if err := bot.BuildProduction(CelestialID(planetID), ID(ogameID), nbr); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// BuildBuildingHandler ...
func BuildBuildingHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	ogameID, err := strconv.ParseInt(c.Param("ogameID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ogame id"))
	}
	if err := bot.BuildBuilding(CelestialID(planetID), ID(ogameID)); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// BuildTechnologyHandler ...
func BuildTechnologyHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	ogameID, err := strconv.ParseInt(c.Param("ogameID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ogame id"))
	}
	if err := bot.BuildTechnology(CelestialID(planetID), ID(ogameID)); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// BuildDefenseHandler ...
func BuildDefenseHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	ogameID, err := strconv.ParseInt(c.Param("ogameID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ogame id"))
	}
	nbr, err := strconv.ParseInt(c.Param("nbr"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid nbr"))
	}
	if err := bot.BuildDefense(CelestialID(planetID), ID(ogameID), nbr); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// BuildShipsHandler ...
func BuildShipsHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	ogameID, err := strconv.ParseInt(c.Param("ogameID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ogame id"))
	}
	nbr, err := strconv.ParseInt(c.Param("nbr"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid nbr"))
	}
	if err := bot.BuildShips(CelestialID(planetID), ID(ogameID), nbr); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// GetProductionHandler ...
func GetProductionHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	res, _, err := bot.GetProduction(CelestialID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(res))
}

// ConstructionsBeingBuiltHandler ...
func ConstructionsBeingBuiltHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	buildingID, buildingCountdown, researchID, researchCountdown := bot.ConstructionsBeingBuilt(CelestialID(planetID))
	return c.JSON(http.StatusOK, SuccessResp(
		struct {
			BuildingID        int64
			BuildingCountdown int64
			ResearchID        int64
			ResearchCountdown int64
		}{
			BuildingID:        int64(buildingID),
			BuildingCountdown: buildingCountdown,
			ResearchID:        int64(researchID),
			ResearchCountdown: researchCountdown,
		},
	))
}

// CancelBuildingHandler ...
func CancelBuildingHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	if err := bot.CancelBuilding(CelestialID(planetID)); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// CancelResearchHandler ...
func CancelResearchHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	if err := bot.CancelResearch(CelestialID(planetID)); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// GetResourcesHandler ...
func GetResourcesHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	res, err := bot.GetResources(CelestialID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(res))
}

// GetPriceHandler ...
func GetPriceHandler(c echo.Context) error {
	ogameID, err := strconv.ParseInt(c.Param("ogameID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ogameID"))
	}
	nbr, err := strconv.ParseInt(c.Param("nbr"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid nbr"))
	}
	ogameObj := Objs.ByID(ID(ogameID))
	if ogameObj != nil {
		price := ogameObj.GetPrice(nbr)
		return c.JSON(http.StatusOK, SuccessResp(price))
	}
	return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ogameID"))
}

// SendFleetHandler ...
// curl 127.0.0.1:1234/bot/planets/123/send-fleet -d 'ships="203,1"&ships="204,10"&speed=10&galaxy=1&system=1&type=1&position=1&mission=3&metal=1&crystal=2&deuterium=3'
func SendFleetHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}

	if err := c.Request().ParseForm(); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid form"))
	}

	var ships []Quantifiable
	where := Coordinate{Type: PlanetType}
	mission := Transport
	var duration int64
	var unionID int64
	payload := Resources{}
	speed := HundredPercent
	for key, values := range c.Request().PostForm {
		switch key {
		case "ships":
			for _, s := range values {
				a := strings.Split(s, ",")
				shipID, err := strconv.ParseInt(a[0], 10, 64)
				if err != nil || !IsShipID(shipID) {
					return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ship id "+a[0]))
				}
				nbr, err := strconv.ParseInt(a[1], 10, 64)
				if err != nil || nbr < 0 {
					return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid nbr "+a[1]))
				}
				ships = append(ships, Quantifiable{ID: ID(shipID), Nbr: nbr})
			}
		case "speed":
			speedInt, err := strconv.ParseInt(values[0], 10, 64)
			if err != nil || speedInt < 0 || speedInt > 10 {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid speed"))
			}
			speed = Speed(speedInt)
		case "galaxy":
			galaxy, err := strconv.ParseInt(values[0], 10, 64)
			if err != nil {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid galaxy"))
			}
			where.Galaxy = galaxy
		case "system":
			system, err := strconv.ParseInt(values[0], 10, 64)
			if err != nil {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid system"))
			}
			where.System = system
		case "position":
			position, err := strconv.ParseInt(values[0], 10, 64)
			if err != nil {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid position"))
			}
			where.Position = position
		case "type":
			t, err := strconv.ParseInt(values[0], 10, 64)
			if err != nil {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid type"))
			}
			where.Type = CelestialType(t)
		case "mission":
			missionInt, err := strconv.ParseInt(values[0], 10, 64)
			if err != nil {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid mission"))
			}
			mission = MissionID(missionInt)
		case "duration":
			duration, err = strconv.ParseInt(values[0], 10, 64)
			if err != nil {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid duration"))
			}
		case "union":
			unionID, err = strconv.ParseInt(values[0], 10, 64)
			if err != nil {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid union id"))
			}
		case "metal":
			metal, err := strconv.ParseInt(values[0], 10, 64)
			if err != nil || metal < 0 {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid metal"))
			}
			payload.Metal = metal
		case "crystal":
			crystal, err := strconv.ParseInt(values[0], 10, 64)
			if err != nil || crystal < 0 {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid crystal"))
			}
			payload.Crystal = crystal
		case "deuterium":
			deuterium, err := strconv.ParseInt(values[0], 10, 64)
			if err != nil || deuterium < 0 {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid deuterium"))
			}
			payload.Deuterium = deuterium
		}
	}

	fleet, err := bot.SendFleet(CelestialID(planetID), ships, speed, where, mission, payload, duration, unionID)
	if err != nil &&
		(err == ErrInvalidPlanetID ||
			err == ErrNoShipSelected ||
			err == ErrUninhabitedPlanet ||
			err == ErrNoDebrisField ||
			err == ErrPlayerInVacationMode ||
			err == ErrAdminOrGM ||
			err == ErrNoAstrophysics ||
			err == ErrNoobProtection ||
			err == ErrPlayerTooStrong ||
			err == ErrNoMoonAvailable ||
			err == ErrNoRecyclerAvailable ||
			err == ErrNoEventsRunning ||
			err == ErrPlanetAlreadyReservecForRelocation) {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(fleet))
}

// GetAlliancePageContentHandler ...
func GetAlliancePageContentHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	allianceID := c.QueryParam("allianceId")
	vals := url.Values{"allianceId": {allianceID}}
	return c.HTML(http.StatusOK, string(bot.GetAlliancePageContent(vals)))
}

func replaceHostname(bot *OGame, html []byte) []byte {
	serverURLBytes := []byte(bot.serverURL)
	apiNewHostnameBytes := []byte(bot.apiNewHostname)
	escapedServerURL := bytes.Replace(serverURLBytes, []byte("/"), []byte(`\/`), -1)
	escapedAPINewHostname := bytes.Replace(apiNewHostnameBytes, []byte("/"), []byte(`\/`), -1)
	html = bytes.Replace(html, serverURLBytes, apiNewHostnameBytes, -1)
	html = bytes.Replace(html, escapedServerURL, escapedAPINewHostname, -1)
	return html
}

// GetStaticHandler ...
func GetStaticHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)

	newURL := bot.serverURL + c.Request().URL.String()
	resp, err := http.Get(newURL)
	if err != nil {
		bot.error(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		bot.error(err)
	}

	// Copy the original HTTP headers to our client
	for k, vv := range resp.Header { // duplicate headers are acceptable in HTTP spec, so add all of them individually: https://stackoverflow.com/questions/4371328/are-duplicate-http-response-headers-acceptable
		k = http.CanonicalHeaderKey(k)
		for _, v := range vv {
			c.Response().Header().Add(k, v)
		}
	}

	if strings.Contains(c.Request().URL.String(), ".xml") {
		body = replaceHostname(bot, body)
		return c.XMLBlob(http.StatusOK, body)
	}

	contentType := http.DetectContentType(body)
	if strings.Contains(newURL, ".css") {
		contentType = "text/css"
	} else if strings.Contains(newURL, ".js") {
		contentType = "text/javascript"
	} else if strings.Contains(newURL, ".gif") {
		contentType = "image/gif"
	}

	return c.Blob(http.StatusOK, contentType, body)
}

// GetFromGameHandler ...
func GetFromGameHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	vals := url.Values{"page": {"ingame"}, "component": {"overview"}}
	if len(c.QueryParams()) > 0 {
		vals = c.QueryParams()
	}
	pageHTML := bot.GetPageContent(vals)
	pageHTML = replaceHostname(bot, pageHTML)
	return c.HTMLBlob(http.StatusOK, pageHTML)
}

// PostToGameHandler ...
func PostToGameHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	vals := url.Values{"page": {"ingame"}, "component": {"overview"}}
	if len(c.QueryParams()) > 0 {
		vals = c.QueryParams()
	}
	payload, _ := c.FormParams()
	pageHTML := bot.PostPageContent(vals, payload)
	pageHTML = replaceHostname(bot, pageHTML)
	return c.HTMLBlob(http.StatusOK, pageHTML)
}

// GetStaticHEADHandler ...
func GetStaticHEADHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	newURL := "/api/" + strings.Join(c.ParamValues(), "") // + "?" + c.QueryString()
	if len(c.QueryString()) > 0 {
		newURL = newURL + "?" + c.QueryString()
	}
	headers, err := bot.HeadersForPage(newURL)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	if len(headers) < 1 {
		return c.NoContent(http.StatusFailedDependency)
	}
	// Copy the original HTTP HEAD headers to our client
	for k, vv := range headers { // duplicate headers are acceptable in HTTP spec, so add all of them individually: https://stackoverflow.com/questions/4371328/are-duplicate-http-response-headers-acceptable
		k = http.CanonicalHeaderKey(k)
		for _, v := range vv {
			c.Response().Header().Add(k, v)
		}
	}
	return c.NoContent(http.StatusOK)
}

// GetEmpireHandler ...
func GetEmpireHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	nbr, err := strconv.ParseInt(c.Param("typeID"), 10, 64)
	if err != nil || nbr > 1 {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid typeID"))
	}
	getEmpire, err := bot.GetEmpire(nbr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(getEmpire))
}

// DeleteMessageHandler ...
func DeleteMessageHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	messageID, err := strconv.ParseInt(c.Param("messageID"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid message id"))
	}
	if err := bot.DeleteMessage(messageID); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// DeleteEspionageMessagesHandler ...
func DeleteEspionageMessagesHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	if err := bot.DeleteAllMessagesFromTab(20); err != nil { // 20 = Espionage Reports
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "Unable to delete Espionage Reports"))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// DeleteMessagesFromTabHandler ...
func DeleteMessagesFromTabHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	tabIndex, err := strconv.ParseInt(c.Param("tabIndex"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "must provide tabIndex"))
	}
	if tabIndex < 20 || tabIndex > 24 {
		/*
			tabid: 20 => Espionage
			tabid: 21 => Combat Reports
			tabid: 22 => Expeditions
			tabid: 23 => Unions/Transport
			tabid: 24 => Other
		*/
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid tabIndex provided"))
	}
	if err := bot.DeleteAllMessagesFromTab(tabIndex); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "Unable to delete message from tab "+strconv.FormatInt(tabIndex, 10)))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// SendIPMHandler ...
func SendIPMHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	ipmAmount, err := strconv.ParseInt(c.Param("ipmAmount"), 10, 64)
	if err != nil || ipmAmount < 1 {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ipmAmount"))
	}
	planetID, err := strconv.ParseInt(c.Param("planetID"), 10, 64)
	if err != nil || planetID < 1 {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	galaxy, err := strconv.ParseInt(c.Request().PostFormValue("galaxy"), 10, 64)
	if err != nil || galaxy < 1 || galaxy > bot.serverData.Galaxies {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid galaxy"))
	}
	system, err := strconv.ParseInt(c.Request().PostFormValue("system"), 10, 64)
	if err != nil || system < 1 || system > bot.serverData.Systems {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid system"))
	}
	position, err := strconv.ParseInt(c.Request().PostFormValue("position"), 10, 64)
	if err != nil || position < 1 || position > 15 {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid position"))
	}
	planetTypeInt, err := strconv.ParseInt(c.Param("type"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	planetType := CelestialType(planetTypeInt)
	if planetType != PlanetType && planetType != MoonType { // only accept planet/moon types
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid type"))
	}
	priority, _ := strconv.ParseInt(c.Request().PostFormValue("priority"), 10, 64)
	coord := Coordinate{Type: planetType, Galaxy: galaxy, System: system, Position: position}
	duration, err := bot.SendIPM(PlanetID(planetID), coord, ipmAmount, ID(priority))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(duration))
}

func GetStaticHandler(c echo.Context) error {
	// TODO: this function does HTTP HEAD request but its requests will not use any (SOCKS) proxy
	// TODO: also add proxy for execRequest() and / or GetStaticHEADHandler?

	bot := c.Get("bot").(*OGame)
	orighostname := "s" + strconv.FormatInt(bot.server.Number, 10) + "-" + bot.server.Language + ".ogame.gameforge.com"

	if (orighostname != "s117-en.ogame.gameforge.com") && (orighostname != "s107-nl.ogame.gameforge.com") {
		bot.debug(orighostname)
	}

	// url := "https://s117-en.ogame.gameforge.com" + c.Request().URL.String() // TODO: swap s117-en<...> with : ogame.Params.APIoriginalhostname ogame.Params.apinewhostname
	url := "https://" + orighostname + c.Request().URL.String()

	resp, err := http.Get(url)

	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			bot.error(err)
		}
	}()

	// Added 21 Dec 2019 for testing as replace of above manual request. Below one uses correct UserAgent and potential use of Proxy.
	// func (b *OGame) execRequest(method, finalURL string, payload, vals url.Values) ([]byte, error) {
	// _, err := bot.execRequest("GET", url , c.QueryParams
	// if _, err := bot.execRequest("GET", url , nil, nil); err != nil {
		// log.Print(err)
	// }

	// Copy the original HTTP headers to our client
	// for k, v := range resp.Header {
		// c.Response().Header().Set(k, strings.Join(v, "")) // strings.Join() to change []string into string
	// }
	for k, vv := range resp.Header { // duplicate headers are acceptable in HTTP spec, so add all of them individually: https://stackoverflow.com/questions/4371328/are-duplicate-http-response-headers-acceptable
		k = http.CanonicalHeaderKey(k)
		for _, v := range vv {
			c.Response().Header().Add(k, v)
		}
	}

	//if strings.Contains(c.Request().URL.String(), "localization.xml") || strings.Contains(c.Request().URL.String(), "serverData.xml") {
	if strings.Contains(c.Request().URL.String(), ".xml") {
		// body2 := strings.Replace(string(body), "s117-en.ogame.gameforge.com", "quantum.webfreakz.nl", -1) // TODO: swap s117-en<...> with : ogame.Params.APIoriginalhostname ogame.Params.apinewhostname

		body2 := strings.Replace(string(body), orighostname, bot.APInewhostname, -1)
		return c.XMLBlob(http.StatusOK, []byte(body2)) // Added 21 Dec 2019
		// return c.Blob(http.StatusOK, contentType, []byte(body2))
	}

	contentType := string(http.DetectContentType(body))

	if strings.Contains(url, ".css") {
		contentType = "text/css"
	} else if strings.Contains(url, ".js") {
		contentType = "text/javascript"
	} else if strings.Contains(url, ".gif") {
		contentType = "image/gif"
	}

	return c.Blob(http.StatusOK, contentType, body)
}

func GetFromGameHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)

	vals := url.Values{"page": {"overview"}}

	if len(c.QueryParams()) > 0 {
		vals = c.QueryParams()
	}

//	log.Print(c.Request())
//	log.Print(c.Request().Header)
//	log.Print(c.Request().URL)
//	log.Print(c.Request().Host)
//	log.Print(c.Request().RequestURI)

	pageHTML := bot.GetPageContent(vals)

	// Replace "s117-en.ogame.gameforge.com" with "quantum.webfreakz.nl"
	html := string(pageHTML)
	// html = strings.Replace(html, "s117-en.ogame.gameforge.com", "quantum.webfreakz.nl", -1)  // TODO: swap s117-en<...> with : ogame.Params.APIoriginalhostname ogame.Params.apinewhostname

	orighostname := "s" + strconv.FormatInt(bot.server.Number, 10) + "-" + bot.server.Language + ".ogame.gameforge.com"

	if (orighostname != "s117-en.ogame.gameforge.com") && (orighostname != "s107-nl.ogame.gameforge.com") {
		bot.debug(orighostname) // s117-en.ogame.gameforge.com
		log.Print(bot.serverURL) // https://s117-en.ogame.gameforge.com
	}

	html = strings.Replace(html, orighostname, bot.APInewhostname, -1)

//	html = strings.Replace(html, "\"/cdn", "\"https://s117-en.ogame.gameforge.com/cdn", -1)

	// TODO: update GetPageContent to also return http.Headers
	// // Copy the original HTTP headers to our client
	// for k, v := range headers {
		// k = http.CanonicalHeaderKey(k)
		// c.Response().Header().Set(k, strings.Join(v, ""))
	// }

	// TODO: remove <div id="banner_skyscraper"></div> which contains ads.
	// TODO: remove <div id="mmonetbar" class="mmoogame" style="display: block;"> contains GameForge game ads

	// doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))

	// if err != nil {
		// log.Print(err)
		// return c.HTML(http.StatusOK, html)
	// }

	// if doc.Find("div#mmonetbar").Size() > 0 {
		// doc.Find("div#mmonetbar").Remove()
	// }

	// if doc.Find("div#banner_skyscraper").Size() > 0 {
		// doc.Find("div#banner_skyscraper").Remove()
	// }

	return c.HTML(http.StatusOK, html)

	// Below stuff should work and remove the ads, but EventList stopped working?
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html)) // Shouldn't be here, ideally it is in ogame.go -> prioritize.go -> extractor -> ExtractHTMLWithoutAdvertisementBars() but could not get it to work.

	if err != nil {
		bot.debug(err)
		return c.HTML(http.StatusOK, html)
	}

	if doc.Find("div#mmonetbar").Size() > 0 {
		doc.Find("div#mmonetbar").Remove()
	}

	if doc.Find("div#banner_skyscraper").Size() > 0 {
		doc.Find("div#banner_skyscraper").Remove()
	}

	returnHTML, err := doc.Html()

	if err != nil {
		bot.debug(err)
		return c.HTML(http.StatusOK, html) // Return original HTML
	}

	return c.HTML(http.StatusOK, returnHTML) // Return filtered HTML (no ads / banners)
}

func PostToGameHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)

	var vals url.Values

	if len(c.QueryParams()) > 0 {
		vals = c.QueryParams()
	} else {
		vals = url.Values{"page": {"overview"}}
	}

	// Payload
	payload, _ := c.FormParams()

	// Perform the post to the library
	byteArray := bot.PostPageContent(vals, payload)

	// Replace "s117-en.ogame.gameforge.com" with "quantum.webfreakz.nl"
	html := string(byteArray)
	// html = strings.Replace(html, "s117-en.ogame.gameforge.com", "quantum.webfreakz.nl", -1)  // TODO: swap s117-en<...> with : ogame.Params.APIoriginalhostname ogame.Params.apinewhostname

	orighostname := "s" + strconv.FormatInt(bot.server.Number, 10) + "-" + bot.server.Language + ".ogame.gameforge.com"

	if (orighostname != "s117-en.ogame.gameforge.com") && (orighostname != "s107-nl.ogame.gameforge.com") {
		bot.debug(orighostname)
		log.Print(bot.serverURL) // https://s117-en.ogame.gameforge.com
	}

	html = strings.Replace(html, orighostname, bot.APInewhostname, -1)

//	html = strings.Replace(html, "\"/cdn", "\"https://s117-en.ogame.gameforge.com/cdn", -1)
	return c.HTML(http.StatusOK, html)
}

func GetAlliancePageContentHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)

	allianceId := c.QueryParam("allianceId")
	vals := url.Values{"allianceId": {allianceId}}

	return c.HTML(http.StatusOK, string(bot.GetAlliancePageContent(vals)))
}

func PostPageContentHandler(c echo.Context) error {
	// bot := c.Get("bot").(*OGame)

	if err := c.Request().ParseForm(); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}

	get := url.Values{}
	post := url.Values{}
	log.Print(c.Request().PostForm)
	log.Print(reflect.TypeOf(c.Request().PostForm))
//	for key, values := range c.Request().PostForm {
//		log.Print(reflect.TypeOf(values))
//		log.Print(key)
//		log.Print(values)
//		switch key {
//		case "url":
//			get, _ = url.Parse(values)
//		case "data":
//			post, _ = url.Parse(values)
//		}
//	}

	log.Print(c.Request().Form["url"])
	log.Print(c.Request().Form["data"])
	log.Print(get)
	log.Print(post)

//	return c.JSON(http.StatusOK, SuccessResp())//bot.PostPageContent(get, post)))
	return c.JSON(http.StatusBadRequest, ErrorResp(400, "kweenie"))
}

func GetPageContentHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)

	if err := c.Request().ParseForm(); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(bot.GetPageContent(c.Request().Form)))
}

func GetStaticHEADHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	url := "/api/" + strings.Join(c.ParamValues(), "") // + "?" + c.QueryString()
	log.Print(url)

	if len(c.ParamValues()) > 0 {
		log.Print(c.ParamValues())
	}

	if len(c.QueryString()) > 0 {
		url = url + "?" + c.QueryString()
		log.Print(c.QueryParams())
	}

	headers, err := bot.HeadersForPage(url)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}

	log.Print(headers)

	if len(headers) < 1 {
	     return c.NoContent(http.StatusFailedDependency) // Code: 424
	}

	// Copy the original HTTP HEAD headers to our client
	// for k, v := range headers {
		// k = http.CanonicalHeaderKey(k)
		// c.Response().Header().Set(k, strings.Join(v, "")) // strings.Join() to change []string into string
	// }
	for k, vv := range headers { // duplicate headers are acceptable in HTTP spec, so add all of them individually: https://stackoverflow.com/questions/4371328/are-duplicate-http-response-headers-acceptable
		k = http.CanonicalHeaderKey(k)
		for _, v := range vv {
			c.Response().Header().Add(k, v)
		}
	}

	c.Response().WriteHeader(http.StatusOK)

	log.Print(c.Response().Header())

	return c.NoContent(http.StatusOK)

	/*
	Original CURL -I to Gameforge for Localization.xml / Players.xml :

	Original:
	################################################################
	curl -I https://s117-en.ogame.gameforge.com/api/localization.xml
	HTTP/1.1 200 OK
	Date: Sun, 03 Mar 2019 19:36:42 GMT
	Server: Apache
	P3P: CP="This is not a P3P policy. It is a ninja, sneaking by Internet Explorer and escaping its iframe at lightning speed."
	Cache-Control: public
	Last-Modified: Sun, 03 Mar 2019 12:44:43 GMT
	Expires: Mon, 04 Mar 2019 12:44:43 GMT
	Connection: close
	Content-Type: application/xml
	################################################################
	curl -I https://s117-en.ogame.gameforge.com/api/players.xml
	HTTP/1.1 200 OK
	Date: Sun, 03 Mar 2019 19:36:48 GMT
	Server: Apache
	P3P: CP="This is not a P3P policy. It is a ninja, sneaking by Internet Explorer and escaping its iframe at lightning speed."
	Cache-Control: public
	Last-Modified: Sun, 03 Mar 2019 12:53:22 GMT
	Expires: Mon, 04 Mar 2019 12:53:22 GMT
	Connection: close
	Content-Type: application/xml
	################################################################

	No Gameforge hostname is specified in any of the returned HTTP Headers: we don't need to use string.Replace() to swap it with our own hostname.

	*/
}


func GetCombatReportMessagesHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)

	report, err := bot.getCombatReportMessages()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(report))
}

// FullCombatReport, takes JSON from CR message. There is no JSON in EspioageReport.
func GetCombatReportHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)

	msgID, err := strconv.ParseInt(c.Param("msgid"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid msgid id"))
	}
	combatReport, err := bot.GetFullCombatReport(msgID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(combatReport))
}

func GetCombatReportForHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)

	galaxy, err := strconv.ParseInt(c.Param("galaxy"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid galaxy"))
	}
	system, err := strconv.ParseInt(c.Param("system"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid system"))
	}
	position, err := strconv.ParseInt(c.Param("position"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid position"))
	}

	planet, err := bot.getCombatReportFor(Coordinate{Type: PlanetType, Galaxy: galaxy, System: system, Position: position})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(planet))
}





// GalaxyEspionage
func GalaxyEspionageHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)

log.SetFlags(log.LstdFlags | log.Lshortfile)
	// e.GET("/bot/planets/:celestialID/espionage/:galaxy/:system/:position/:type/:probecount", galaxyEspionage)

	// PlanetID
	planetID, err := strconv.ParseInt(c.Param("celestialID"), 10, 64)
	if err != nil || planetID < 0 {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}

	// Galaxy
	galaxy, err := strconv.ParseInt(c.Param("galaxy"), 10, 64)
	if err != nil || galaxy < 0 || galaxy > 9 {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid galaxy"))
	}

	// System
	system, err := strconv.ParseInt(c.Param("system"), 10, 64)
	if err != nil || system < 0 || system > 499 {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid system"))
	}

	// Position
	position, err := strconv.ParseInt(c.Param("position"), 10, 64)
	if err != nil || position < 0 || position > 15 {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid position"))
	}

	// Type
	planettype, err := strconv.ParseInt(c.Param("type"), 10, 64)
	if err != nil || planettype < 1 || planettype > 3 || planettype == 2 {
	// log.Print(planettype)
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid type"))
	}

	// Probecount
	probecount, err := strconv.ParseInt(c.Param("probecount"), 10, 64)
	if err != nil || position < 0 {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid probecount"))
	}

	// --------------------------------------------------------------------------------------

	coord := Coordinate{Type: PlanetType, Galaxy: galaxy, System: system, Position: position}
	celestialID := CelestialID(planetID)

	espionage, err := bot.GalaxyEspionage(celestialID, coord, probecount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}

	return c.JSON(http.StatusOK, SuccessResp(espionage))
}

func GetEventListHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)

	return c.JSON(http.StatusOK, SuccessResp(bot.GetFleetsFromEventList()))
}

// -- Start: Auction
func GetAuctionHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)

	auction, err := bot.GetAuction()

	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "could not open auction page"))
	}

	return c.JSON(http.StatusOK, SuccessResp(auction))
}

func DoAuctionHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	// Collect our parameters:
	// - planetID
	// - planetType
	// - bid
	// - resourceType

	// PlanetID
	planetID, err := strconv.ParseInt(c.Request().PostFormValue("planetid"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}

	// PlanetType
	planetType, err := strconv.ParseInt(c.Request().PostFormValue("planettype"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planettype"))
	}

	// Bid
	bid, err := strconv.ParseInt(c.Request().PostFormValue("bid"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid bid"))
	}

	// ResourceType
	resourceType, err := strconv.ParseInt(c.Request().PostFormValue("resourcetype"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid resourcetype"))
	}
	// --------------------------------------------

	auctionMsg, err := bot.DoAuction(planetID, planetType, bid, resourceType)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}

	return c.JSON(http.StatusOK, SuccessResp(auctionMsg))
}

