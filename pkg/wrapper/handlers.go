package wrapper

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/alaingilbert/ogame/pkg/gameforge"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/utils"
	echo "github.com/labstack/echo/v4"
)

// APIResp ...
type APIResp struct {
	Status  string
	Code    int
	Message string
	Result  any
}

// SuccessResp ...
func SuccessResp(data any) APIResp {
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
	return c.JSON(http.StatusOK, map[string]any{
		"version": version,
		"commit":  commit,
		"date":    date,
	})
}

// TasksHandler return how many tasks are queued in the heap.
func TasksHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	return c.JSON(http.StatusOK, SuccessResp(bot.GetTasks()))
}

// GetServerHandler ...
func GetServerHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	return c.JSON(http.StatusOK, SuccessResp(bot.GetServer()))
}

// GetServerDataHandler ...
func GetServerDataHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	return c.JSON(http.StatusOK, SuccessResp(bot.cache.serverData))
}

// SetUserAgentHandler deprecated
func SetUserAgentHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, ErrorResp(http.StatusBadRequest, "deprecated"))
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
	pageHTML, _ := bot.GetPageContent(c.Request().Form)
	return c.JSON(http.StatusOK, SuccessResp(pageHTML))
}

// LoginHandler ...
func LoginHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	if _, _, err := bot.LoginWithExistingCookies(); err != nil {
		if err == gameforge.ErrBadCredentials {
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
	return c.JSON(http.StatusOK, SuccessResp(bot.cache.serverData.Speed))
}

// GetUniverseSpeedFleetHandler ...
func GetUniverseSpeedFleetHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	return c.JSON(http.StatusOK, SuccessResp(bot.cache.serverData.SpeedFleet))
}

// ServerVersionHandler ...
func ServerVersionHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	return c.JSON(http.StatusOK, SuccessResp(bot.cache.serverData.Version))
}

// ServerTimeHandler ...
func ServerTimeHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	serverTime, _ := bot.ServerTime()
	return c.JSON(http.StatusOK, SuccessResp(serverTime))
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

func IsUnderAttackByIDHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := utils.ParseI64(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	isUnderAttack, err := bot.IsUnderAttack(ChangePlanet(ogame.CelestialID(planetID)))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(isUnderAttack))
}

// IsVacationModeHandler ...
func IsVacationModeHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	isVacationMode := bot.cache.isVacationModeEnabled
	return c.JSON(http.StatusOK, SuccessResp(isVacationMode))
}

// GetUserInfosHandler ...
func GetUserInfosHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	userInfo, _ := bot.GetUserInfos()
	return c.JSON(http.StatusOK, SuccessResp(userInfo))
}

// GetCharacterClassHandler ...
func GetCharacterClassHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	return c.JSON(http.StatusOK, SuccessResp(bot.CharacterClass()))
}

// HasCommanderHandler ...
func HasCommanderHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	hasCommander := bot.cache.hasCommander
	return c.JSON(http.StatusOK, SuccessResp(hasCommander))
}

// HasAdmiralHandler ...
func HasAdmiralHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	hasAdmiral := bot.cache.hasAdmiral
	return c.JSON(http.StatusOK, SuccessResp(hasAdmiral))
}

// HasEngineerHandler ...
func HasEngineerHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	hasEngineer := bot.cache.hasEngineer
	return c.JSON(http.StatusOK, SuccessResp(hasEngineer))
}

// HasGeologistHandler ...
func HasGeologistHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	hasGeologist := bot.cache.hasGeologist
	return c.JSON(http.StatusOK, SuccessResp(hasGeologist))
}

// HasTechnocratHandler ...
func HasTechnocratHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	hasTechnocrat := bot.cache.hasTechnocrat
	return c.JSON(http.StatusOK, SuccessResp(hasTechnocrat))
}

// GetEspionageReportMessagesHandler ...
func GetEspionageReportMessagesHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	report, err := bot.GetEspionageReportMessages(-1)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(report))
}

// GetEspionageReportHandler ...
func GetEspionageReportHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	msgID, err := utils.ParseI64(c.Param("msgid"))
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
	galaxy, err := utils.ParseI64(c.Param("galaxy"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid galaxy"))
	}
	system, err := utils.ParseI64(c.Param("system"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid system"))
	}
	position, err := utils.ParseI64(c.Param("position"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid position"))
	}
	planet, err := bot.GetEspionageReportFor(ogame.Coordinate{Type: ogame.PlanetType, Galaxy: galaxy, System: system, Position: position})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(planet))
}

// SendMessageHandler ...
// curl 127.0.0.1:1234/bot/send-message -d 'playerID=123&message="Sup boi!"'
func SendMessageHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	playerID, err := utils.ParseI64(c.Request().PostFormValue("playerID"))
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
	fleets, _, _ := bot.GetFleets()
	return c.JSON(http.StatusOK, SuccessResp(fleets))
}

// GetSlotsHandler ...
func GetSlotsHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	slots, _ := bot.GetSlots()
	return c.JSON(http.StatusOK, SuccessResp(slots))
}

// CancelFleetHandler ...
func CancelFleetHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	fleetID, err := utils.ParseI64(c.Param("fleetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(bot.CancelFleet(ogame.FleetID(fleetID))))
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
	galaxy, err := utils.ParseI64(c.Param("galaxy"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	system, err := utils.ParseI64(c.Param("system"))
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
	researches, _ := bot.GetResearch()
	return c.JSON(http.StatusOK, SuccessResp(researches))
}

// BuyOfferOfTheDayHandler ...
func BuyOfferOfTheDayHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	if err := bot.BuyOfferOfTheDay(); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// GetMoonsHandler ...
func GetMoonsHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	moons, _ := bot.GetMoons()
	return c.JSON(http.StatusOK, SuccessResp(moons))
}

// GetMoonHandler ...
func GetMoonHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	moonID, err := utils.ParseI64(c.Param("moonID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid moon id"))
	}
	moon, err := bot.GetMoon(moonID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid moon id"))
	}
	return c.JSON(http.StatusOK, SuccessResp(moon))
}

// GetMoonByCoordHandler ...
func GetMoonByCoordHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	galaxy, err := utils.ParseI64(c.Param("galaxy"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid galaxy"))
	}
	system, err := utils.ParseI64(c.Param("system"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid system"))
	}
	position, err := utils.ParseI64(c.Param("position"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid position"))
	}
	planet, err := bot.GetMoon(ogame.Coordinate{Type: ogame.MoonType, Galaxy: galaxy, System: system, Position: position})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(planet))
}

// GetPlanetsHandler ...
func GetPlanetsHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planets, _ := bot.GetPlanets()
	return c.JSON(http.StatusOK, SuccessResp(planets))
}

// CelestialAbandonHandler ...
func CelestialAbandonHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	celestialID, err := utils.ParseI64(c.Param("celestialID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid celestial id"))
	}
	err = bot.Abandon(celestialID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	_, err = bot.GetCelestial(celestialID)
	if err != nil {
		return c.JSON(http.StatusOK, SuccessResp(struct {
			CelestialID int64
			Result      string
		}{
			CelestialID: celestialID,
			Result:      "succeed",
		}))
	} else {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "Celestial could not be deleted"))
	}

}

// GetCelestialItemsHandler ...
func GetCelestialItemsHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	celestialID, err := utils.ParseI64(c.Param("celestialID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid celestial id"))
	}
	items, err := bot.GetItems(ogame.CelestialID(celestialID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(items))
}

// ActivateCelestialItemHandler ...
func ActivateCelestialItemHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	celestialID, err := utils.ParseI64(c.Param("celestialID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid celestial id"))
	}
	ref := c.Param("itemRef")
	if err := bot.ActivateItem(ref, ogame.CelestialID(celestialID)); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// GetPlanetHandler ...
func GetPlanetHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := utils.ParseI64(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	planet, err := bot.GetPlanet(ogame.PlanetID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(planet))
}

// GetPlanetByCoordHandler ...
func GetPlanetByCoordHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	galaxy, err := utils.ParseI64(c.Param("galaxy"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid galaxy"))
	}
	system, err := utils.ParseI64(c.Param("system"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid system"))
	}
	position, err := utils.ParseI64(c.Param("position"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid position"))
	}
	planet, err := bot.GetPlanet(ogame.Coordinate{Type: ogame.PlanetType, Galaxy: galaxy, System: system, Position: position})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(planet))
}

// GetResourcesDetailsHandler ...
func GetResourcesDetailsHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := utils.ParseI64(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	resources, err := bot.GetResourcesDetails(ogame.CelestialID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(resources))
}

// GetResourceSettingsHandler ...
func GetResourceSettingsHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := utils.ParseI64(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	res, err := bot.GetResourceSettings(ogame.PlanetID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(res))
}

// SetResourceSettingsHandler ...
// curl 127.0.0.1:1234/bot/planets/123/resource-settings -d 'metalMine=100&crystalMine=100&deuteriumSynthesizer=100&solarPlant=100&fusionReactor=100&solarSatellite=100'
func SetResourceSettingsHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := utils.ParseI64(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	metalMine, err := utils.ParseI64(c.Request().PostFormValue("metalMine"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid metalMine"))
	}
	crystalMine, err := utils.ParseI64(c.Request().PostFormValue("crystalMine"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid crystalMine"))
	}
	deuteriumSynthesizer, err := utils.ParseI64(c.Request().PostFormValue("deuteriumSynthesizer"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid deuteriumSynthesizer"))
	}
	solarPlant, err := utils.ParseI64(c.Request().PostFormValue("solarPlant"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid solarPlant"))
	}
	fusionReactor, err := utils.ParseI64(c.Request().PostFormValue("fusionReactor"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid fusionReactor"))
	}
	solarSatellite, err := utils.ParseI64(c.Request().PostFormValue("solarSatellite"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid solarSatellite"))
	}
	crawler, err := utils.ParseI64(c.Request().PostFormValue("crawler"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid crawler"))
	}
	settings := ogame.ResourceSettings{
		MetalMine:            metalMine,
		CrystalMine:          crystalMine,
		DeuteriumSynthesizer: deuteriumSynthesizer,
		SolarPlant:           solarPlant,
		FusionReactor:        fusionReactor,
		SolarSatellite:       solarSatellite,
		Crawler:              crawler,
	}
	if err := bot.SetResourceSettings(ogame.PlanetID(planetID), settings); err != nil {
		if err == ogame.ErrInvalidPlanetID {
			return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// GetLfBuildingsHandler ...
func GetLfBuildingsHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := utils.ParseI64(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	res, err := bot.GetLfBuildings(ogame.CelestialID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(res))
}

// GetLfResearchHandler ...
func GetLfResearchHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := utils.ParseI64(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	res, err := bot.GetLfResearch(ogame.CelestialID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(res))
}

// GetResourcesBuildingsHandler ...
func GetResourcesBuildingsHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := utils.ParseI64(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	res, err := bot.GetResourcesBuildings(ogame.CelestialID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(res))
}

// GetDefenseHandler ...
func GetDefenseHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := utils.ParseI64(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	res, err := bot.GetDefense(ogame.CelestialID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(res))
}

// GetShipsHandler ...
func GetShipsHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := utils.ParseI64(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	res, err := bot.GetShips(ogame.CelestialID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(res))
}

// GetFacilitiesHandler ...
func GetFacilitiesHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := utils.ParseI64(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	res, err := bot.GetFacilities(ogame.CelestialID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(res))
}

// BuildHandler ...
func BuildHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := utils.ParseI64(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	ogameID, err := utils.ParseI64(c.Param("ogameID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ogame id"))
	}
	nbr, err := utils.ParseI64(c.Param("nbr"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid nbr"))
	}
	if err := bot.Build(ogame.CelestialID(planetID), ogame.ID(ogameID), nbr); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// BuildCancelableHandler ...
func BuildCancelableHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := utils.ParseI64(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	ogameID, err := utils.ParseI64(c.Param("ogameID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ogame id"))
	}
	if err := bot.BuildCancelable(ogame.CelestialID(planetID), ogame.ID(ogameID)); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// BuildProductionHandler ...
func BuildProductionHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := utils.ParseI64(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	ogameID, err := utils.ParseI64(c.Param("ogameID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ogame id"))
	}
	nbr, err := utils.ParseI64(c.Param("nbr"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid nbr"))
	}
	if err := bot.BuildProduction(ogame.CelestialID(planetID), ogame.ID(ogameID), nbr); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// BuildBuildingHandler ...
func BuildBuildingHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := utils.ParseI64(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	ogameID, err := utils.ParseI64(c.Param("ogameID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ogame id"))
	}
	if err := bot.BuildBuilding(ogame.CelestialID(planetID), ogame.ID(ogameID)); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// BuildTechnologyHandler ...
func BuildTechnologyHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := utils.ParseI64(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	ogameID, err := utils.ParseI64(c.Param("ogameID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ogame id"))
	}
	if err := bot.BuildTechnology(ogame.CelestialID(planetID), ogame.ID(ogameID)); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// BuildDefenseHandler ...
func BuildDefenseHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := utils.ParseI64(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	ogameID, err := utils.ParseI64(c.Param("ogameID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ogame id"))
	}
	nbr, err := utils.ParseI64(c.Param("nbr"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid nbr"))
	}
	if err := bot.BuildDefense(ogame.CelestialID(planetID), ogame.ID(ogameID), nbr); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// BuildShipsHandler ...
func BuildShipsHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := utils.ParseI64(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	ogameID, err := utils.ParseI64(c.Param("ogameID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ogame id"))
	}
	nbr, err := utils.ParseI64(c.Param("nbr"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid nbr"))
	}
	if err := bot.BuildShips(ogame.CelestialID(planetID), ogame.ID(ogameID), nbr); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// GetProductionHandler ...
func GetProductionHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := utils.ParseI64(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	res, _, err := bot.GetProduction(ogame.CelestialID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(res))
}

// ConstructionsBeingBuiltHandler ...
func ConstructionsBeingBuiltHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := utils.ParseI64(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	constructions, _ := bot.ConstructionsBeingBuilt(ogame.CelestialID(planetID))
	return c.JSON(http.StatusOK, SuccessResp(
		struct {
			BuildingID          int64
			BuildingCountdown   int64
			ResearchID          int64
			ResearchCountdown   int64
			LfBuildingID        int64
			LfBuildingCountdown int64
			LfResearchID        int64
			LfResearchCountdown int64
		}{
			BuildingID:          int64(constructions.Building.ID),
			BuildingCountdown:   int64(constructions.Building.Countdown.Seconds()),
			ResearchID:          int64(constructions.Research.ID),
			ResearchCountdown:   int64(constructions.Research.Countdown.Seconds()),
			LfBuildingID:        int64(constructions.LfBuilding.ID),
			LfBuildingCountdown: int64(constructions.LfBuilding.Countdown.Seconds()),
			LfResearchID:        int64(constructions.LfResearch.ID),
			LfResearchCountdown: int64(constructions.LfResearch.Countdown.Seconds()),
		},
	))
}

// CancelBuildingHandler ...
func CancelBuildingHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := utils.ParseI64(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	if err := bot.CancelBuilding(ogame.CelestialID(planetID)); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// CancelResearchHandler ...
func CancelResearchHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := utils.ParseI64(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	if err := bot.CancelResearch(ogame.CelestialID(planetID)); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// GetResourcesHandler ...
func GetResourcesHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := utils.ParseI64(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	res, err := bot.GetResources(ogame.CelestialID(planetID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(res))
}

// GetRequirementsHandler ...
func GetRequirementsHandler(c echo.Context) error {
	ogameID, err := utils.ParseI64(c.Param("ogameID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ogameID"))
	}
	ogameObj := ogame.Objs.ByID(ogame.ID(ogameID))
	if ogameObj != nil {
		requirements := ogameObj.GetRequirements()
		return c.JSON(http.StatusOK, SuccessResp(requirements))
	}
	return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ogameID"))
}

// GetPriceHandler ...
func GetPriceHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	ogameID, err := utils.ParseI64(c.Param("ogameID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ogameID"))
	}
	nbr, err := utils.ParseI64(c.Param("nbr"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid nbr"))
	}
	ogameObj := ogame.Objs.ByID(ogame.ID(ogameID))
	if ogameObj != nil {
		lfBonuses, _ := bot.GetCachedLfBonuses()
		price := ogameObj.GetPrice(nbr, lfBonuses)
		return c.JSON(http.StatusOK, SuccessResp(price))
	}
	return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ogameID"))
}

// SendFleetHandler ...
// curl 127.0.0.1:1234/bot/planets/123/send-fleet -d 'ships=203,1&ships=204,10&speed=10&galaxy=1&system=1&type=1&position=1&mission=3&metal=1&crystal=2&deuterium=3'
func SendFleetHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := utils.ParseI64(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}

	if err := c.Request().ParseForm(); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid form"))
	}

	var ships ogame.ShipsInfos
	where := ogame.Coordinate{Type: ogame.PlanetType}
	mission := ogame.Transport
	var duration int64
	var unionID int64
	payload := ogame.Resources{}
	speed := ogame.HundredPercent
	for key, values := range c.Request().PostForm {
		switch key {
		case "ships":
			for _, s := range values {
				a := strings.Split(s, ",")
				shipID, err := utils.ParseI64(a[0])
				if err != nil || !ogame.ID(shipID).IsShip() {
					return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ship id "+a[0]))
				}
				nbr, err := utils.ParseI64(a[1])
				if err != nil || nbr < 0 {
					return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid nbr "+a[1]))
				}
				ships.Set(ogame.ID(shipID), nbr)
			}
		case "speed":
			speedInt, err := utils.ParseI64(values[0])
			if err != nil || speedInt < 0 || speedInt > 10 {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid speed"))
			}
			speed = ogame.Speed(speedInt)
		case "galaxy":
			galaxy, err := utils.ParseI64(values[0])
			if err != nil {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid galaxy"))
			}
			where.Galaxy = galaxy
		case "system":
			system, err := utils.ParseI64(values[0])
			if err != nil {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid system"))
			}
			where.System = system
		case "position":
			position, err := utils.ParseI64(values[0])
			if err != nil {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid position"))
			}
			where.Position = position
		case "type":
			t, err := utils.ParseI64(values[0])
			if err != nil {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid type"))
			}
			where.Type = ogame.CelestialType(t)
		case "mission":
			missionInt, err := utils.ParseI64(values[0])
			if err != nil {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid mission"))
			}
			mission = ogame.MissionID(missionInt)
		case "duration":
			duration, err = utils.ParseI64(values[0])
			if err != nil {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid duration"))
			}
		case "union":
			unionID, err = utils.ParseI64(values[0])
			if err != nil {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid union id"))
			}
		case "metal":
			metal, err := utils.ParseI64(values[0])
			if err != nil || metal < 0 {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid metal"))
			}
			payload.Metal = metal
		case "crystal":
			crystal, err := utils.ParseI64(values[0])
			if err != nil || crystal < 0 {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid crystal"))
			}
			payload.Crystal = crystal
		case "deuterium":
			deuterium, err := utils.ParseI64(values[0])
			if err != nil || deuterium < 0 {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid deuterium"))
			}
			payload.Deuterium = deuterium
		}
	}

	fleet, err := bot.SendFleet(ogame.CelestialID(planetID), ships, speed, where, mission, payload, duration, unionID)
	if err != nil &&
		(err == ogame.ErrInvalidPlanetID ||
			err == ogame.ErrNoShipSelected ||
			err == ogame.ErrUninhabitedPlanet ||
			err == ogame.ErrNoDebrisField ||
			err == ogame.ErrPlayerInVacationMode ||
			err == ogame.ErrAdminOrGM ||
			err == ogame.ErrNoAstrophysics ||
			err == ogame.ErrNoobProtection ||
			err == ogame.ErrPlayerTooStrong ||
			err == ogame.ErrNoMoonAvailable ||
			err == ogame.ErrNoRecyclerAvailable ||
			err == ogame.ErrNoEventsRunning ||
			err == ogame.ErrPlanetAlreadyReservedForRelocation) {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(fleet))
}

// SendDiscoveryHandler ...
// curl 127.0.0.1:1234/bot/planets/123/send-discovery -d 'galaxy=1&system=1&type=1&position=1'
func SendDiscoveryHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := utils.ParseI64(c.Param("planetID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}

	if err := c.Request().ParseForm(); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid form"))
	}

	where := ogame.Coordinate{Type: ogame.PlanetType}
	for key, values := range c.Request().PostForm {
		switch key {
		case "galaxy":
			galaxy, err := utils.ParseI64(values[0])
			if err != nil {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid galaxy"))
			}
			where.Galaxy = galaxy
		case "system":
			system, err := utils.ParseI64(values[0])
			if err != nil {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid system"))
			}
			where.System = system
		case "position":
			position, err := utils.ParseI64(values[0])
			if err != nil {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid position"))
			}
			where.Position = position
		}
	}

	if err := bot.SendDiscoveryFleet(ogame.CelestialID(planetID), where); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(true))
}

// GetAlliancePageContentHandler ...
func GetAlliancePageContentHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	allianceID := c.QueryParam("allianceId")
	vals := url.Values{"allianceId": {allianceID}}
	pageHTML, _ := bot.GetPageContent(vals)
	pageHTML = removeCookiesBanner(pageHTML)
	return c.HTML(http.StatusOK, string(pageHTML))
}

func replaceHostname(bot *OGame, html []byte) []byte {
	serverURLBytes := []byte(bot.cache.serverURL)
	apiNewHostnameBytes := []byte(bot.apiNewHostname)
	escapedServerURL := bytes.Replace(serverURLBytes, []byte("/"), []byte(`\/`), -1)
	doubleEscapedServerURL := bytes.Replace(serverURLBytes, []byte("/"), []byte("\\\\\\/"), -1)
	escapedAPINewHostname := bytes.Replace(apiNewHostnameBytes, []byte("/"), []byte(`\/`), -1)
	doubleEscapedAPINewHostname := bytes.Replace(apiNewHostnameBytes, []byte("/"), []byte("\\\\\\/"), -1)
	html = bytes.Replace(html, serverURLBytes, apiNewHostnameBytes, -1)
	html = bytes.Replace(html, escapedServerURL, escapedAPINewHostname, -1)
	html = bytes.Replace(html, doubleEscapedServerURL, doubleEscapedAPINewHostname, -1)
	return html
}

// GetStaticHandler ...
func GetStaticHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)

	newURL := bot.cache.serverURL + c.Request().URL.String()
	req, err := http.NewRequest(http.MethodGet, newURL, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	resp, err := bot.device.GetClient().Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}

	// Copy the original HTTP headers to our client
	for k, vv := range resp.Header { // duplicate headers are acceptable in HTTP spec, so add all of them individually: https://stackoverflow.com/questions/4371328/are-duplicate-http-response-headers-acceptable
		k = http.CanonicalHeaderKey(k)
		if k != "Content-Length" && k != "Content-Encoding" { // https://github.com/alaingilbert/ogame/pull/80#issuecomment-674559853
			for _, v := range vv {
				c.Response().Header().Add(k, v)
			}
		}
	}

	if strings.Contains(c.Request().URL.String(), ".xml") {
		body = replaceHostname(bot, body)
		return c.Blob(http.StatusOK, "application/xml", body)
	}

	contentType := http.DetectContentType(body)
	if strings.Contains(newURL, ".css") {
		contentType = "text/css"
	} else if strings.Contains(newURL, ".js") {
		contentType = "application/javascript"
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
	pageHTML, _ := bot.GetPageContent(vals)
	pageHTML = replaceHostname(bot, pageHTML)
	pageHTML = removeCookiesBanner(pageHTML)
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
	pageHTML, _ := bot.PostPageContent(vals, payload)
	pageHTML = replaceHostname(bot, pageHTML)
	pageHTML = removeCookiesBanner(pageHTML)
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
	nbr, err := utils.ParseI64(c.Param("typeID"))
	if err != nil || nbr > 1 {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid typeID"))
	}
	var celestialType ogame.CelestialType
	switch nbr {
	case 0:
		celestialType = ogame.PlanetType
	case 1:
		celestialType = ogame.MoonType
	}
	getEmpire, err := bot.GetEmpireJSON(celestialType)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(getEmpire))
}

// DeleteMessageHandler ...
func DeleteMessageHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	messageID, err := utils.ParseI64(c.Param("messageID"))
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
	tabIndex, err := utils.ParseI64(c.Param("tabIndex"))
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
	if err := bot.DeleteAllMessagesFromTab(ogame.MessagesTabID(tabIndex)); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "Unable to delete message from tab "+utils.FI64(tabIndex)))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// SendIPMHandler ...
func SendIPMHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	ipmAmount, err := utils.ParseI64(c.Param("ipmAmount"))
	if err != nil || ipmAmount < 1 {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ipmAmount"))
	}
	planetID, err := utils.ParseI64(c.Param("planetID"))
	if err != nil || planetID < 1 {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	galaxy, err := utils.ParseI64(c.Request().PostFormValue("galaxy"))
	if err != nil || galaxy < 1 || galaxy > bot.cache.serverData.Galaxies {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid galaxy"))
	}
	system, err := utils.ParseI64(c.Request().PostFormValue("system"))
	if err != nil || system < 1 || system > bot.cache.serverData.Systems {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid system"))
	}
	position, err := utils.ParseI64(c.Request().PostFormValue("position"))
	if err != nil || position < 1 || position > 15 {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid position"))
	}
	planetTypeInt, err := utils.ParseI64(c.Param("type"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	planetType := ogame.CelestialType(planetTypeInt)
	if planetType != ogame.PlanetType && planetType != ogame.MoonType { // only accept planet/moon types
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid type"))
	}
	priority := utils.DoParseI64(c.Request().PostFormValue("priority"))
	coord := ogame.Coordinate{Type: planetType, Galaxy: galaxy, System: system, Position: position}
	duration, err := bot.SendIPM(ogame.PlanetID(planetID), coord, ipmAmount, ogame.ID(priority))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(duration))
}

// TeardownHandler ...
func TeardownHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	planetID, err := utils.ParseI64(c.Param("planetID"))
	if err != nil || planetID < 0 {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid planet id"))
	}
	ogameID, err := utils.ParseI64(c.Param("ogameID"))
	if err != nil || planetID < 0 {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ogame id"))
	}
	if err = bot.TearDown(ogame.CelestialID(planetID), ogame.ID(ogameID)); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// GetAuctionHandler ...
func GetAuctionHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	auction, err := bot.GetAuction()
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "could not open auction page"))
	}
	return c.JSON(http.StatusOK, SuccessResp(auction))
}

// DoAuctionHandler (`celestialID=metal:crystal:deuterium` eg: `123456=123:456:789`)
func DoAuctionHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	bid := make(map[ogame.CelestialID]ogame.Resources)
	if err := c.Request().ParseForm(); err != nil { // Required for PostForm, not for PostFormValue
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid form"))
	}
	for key, values := range c.Request().PostForm {
		for _, s := range values {
			var metal, crystal, deuterium int64
			if n, err := fmt.Sscanf(s, "%d:%d:%d", &metal, &crystal, &deuterium); err != nil || n != 3 {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid bid format"))
			}
			celestialIDInt, err := utils.ParseI64(key)
			if err != nil {
				return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid celestial ID"))
			}
			bid[ogame.CelestialID(celestialIDInt)] = ogame.Resources{Metal: metal, Crystal: crystal, Deuterium: deuterium}
		}
	}
	if err := bot.DoAuction(bid); err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(nil))
}

// PhalanxHandler ...
func PhalanxHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	moonID, err := utils.ParseI64(c.Param("moonID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid moon id"))
	}
	galaxy, err := utils.ParseI64(c.Param("galaxy"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid galaxy"))
	}
	system, err := utils.ParseI64(c.Param("system"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid system"))
	}
	position, err := utils.ParseI64(c.Param("position"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid position"))
	}
	coord := ogame.Coordinate{Type: ogame.PlanetType, Galaxy: galaxy, System: system, Position: position}
	fleets, err := bot.Phalanx(ogame.MoonID(moonID), coord)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(fleets))
}

// JumpGateHandler ...
func JumpGateHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	if err := c.Request().ParseForm(); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid form"))
	}
	moonOriginID, err := utils.ParseI64(c.Param("moonID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid origin moon id"))
	}
	moonDestinationID, err := utils.ParseI64(c.Request().PostFormValue("moonDestination"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid destination moon id"))
	}
	var ships ogame.ShipsInfos
	for key, values := range c.Request().PostForm {
		switch key {
		case "ships":
			for _, s := range values {
				a := strings.Split(s, ",")
				shipID, err := utils.ParseI64(a[0])
				if err != nil || !ogame.ID(shipID).IsShip() {
					return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid ship id "+a[0]))
				}
				nbr, err := utils.ParseI64(a[1])
				if err != nil || nbr < 0 {
					return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid nbr "+a[1]))
				}
				ships.Set(ogame.ID(shipID), nbr)
			}
		}
	}
	success, rechargeCountdown, err := bot.JumpGate(ogame.MoonID(moonOriginID), ogame.MoonID(moonDestinationID), ships)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(map[string]any{
		"success":           success,
		"rechargeCountdown": rechargeCountdown,
	}))
}

// TechsHandler ...
func TechsHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	celestialID, err := utils.ParseI64(c.Param("celestialID"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, "invalid celestial id"))
	}
	techs, err := bot.GetTechs(ogame.CelestialID(celestialID))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResp(400, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(map[string]any{
		"supplies":     techs.ResourcesBuildings,
		"facilities":   techs.Facilities,
		"ships":        techs.ShipsInfos,
		"defenses":     techs.DefensesInfos,
		"researches":   techs.Researches,
		"lfbuildings":  techs.LfBuildings,
		"lfResearches": techs.LfResearches,
	}))
}

// GetCaptchaHandler ...
func GetCaptchaHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	gf, _ := gameforge.New(&gameforge.Config{Ctx: bot.ctx, Device: bot.device, Platform: PLATFORM, Lobby: bot.lobby})
	_, err := gf.Login(&gameforge.LoginParams{
		Username:  bot.username,
		Password:  bot.password,
		OtpSecret: bot.otpSecret,
	})
	var captchaErr *gameforge.CaptchaRequiredError
	if errors.As(err, &captchaErr) {
		questionRaw, iconsRaw, err := gameforge.StartChallenge(bot.ctx, bot.GetClient(), captchaErr.ChallengeID)
		if err != nil {
			return c.HTML(http.StatusOK, err.Error())
		}

		questionB64 := base64.StdEncoding.EncodeToString(questionRaw)
		iconsB64 := base64.StdEncoding.EncodeToString(iconsRaw)

		html := `<img style="background-color: black;" src="data:image/png;base64,` + questionB64 + `" /><br />
<img style="background-color: black;" src="data:image/png;base64,` + iconsB64 + `" /><br />
<form action="/bot/captcha/solve" method="POST">
	<input type="hidden" name="challenge_id" value="` + captchaErr.ChallengeID + `" />
	Enter 0,1,2 or 3 and press Enter <input type="number" name="answer" />
</form>` + captchaErr.ChallengeID

		return c.HTML(http.StatusOK, html)
	} else if err != nil {
		return c.HTML(http.StatusOK, err.Error())
	}
	return c.HTML(http.StatusOK, "no captcha found")
}

// GetCaptchaSolverHandler ...
func GetCaptchaSolverHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	challengeID := c.Request().PostFormValue("challenge_id")
	answer := utils.DoParseI64(c.Request().PostFormValue("answer"))

	if err := gameforge.SolveChallenge(bot.ctx, bot.GetClient(), challengeID, answer); err != nil {
		bot.error(err)
	}

	if !bot.IsLoggedIn() {
		if err := bot.Login(); err != nil {
			bot.error(err)
		}
	}
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

// CaptchaChallenge ...
type CaptchaChallenge struct {
	ID       string
	Question string
	Icons    string
}

// GetCaptchaChallengeHandler ...
func GetCaptchaChallengeHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	gf, _ := gameforge.New(&gameforge.Config{Ctx: bot.ctx, Device: bot.device, Platform: PLATFORM, Lobby: bot.lobby})
	_, err := gf.Login(&gameforge.LoginParams{
		Username:  bot.username,
		Password:  bot.password,
		OtpSecret: bot.otpSecret,
	})
	var captchaErr *gameforge.CaptchaRequiredError
	if errors.As(err, &captchaErr) {
		questionRaw, iconsRaw, err := gameforge.StartChallenge(bot.ctx, bot.GetClient(), captchaErr.ChallengeID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
		}
		questionB64 := base64.StdEncoding.EncodeToString(questionRaw)
		iconsB64 := base64.StdEncoding.EncodeToString(iconsRaw)
		return c.JSON(http.StatusOK, SuccessResp(CaptchaChallenge{
			ID:       captchaErr.ChallengeID,
			Question: questionB64,
			Icons:    iconsB64,
		}))
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(CaptchaChallenge{}))
}

// GetPublicIPHandler ...
func GetPublicIPHandler(c echo.Context) error {
	bot := c.Get("bot").(*OGame)
	ip, err := bot.GetPublicIP()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResp(500, err.Error()))
	}
	return c.JSON(http.StatusOK, SuccessResp(ip))
}

func removeCookiesBanner(pageHTML []byte) []byte {
	regex := `<script[^>]*id="cookiebanner"[^>]*>[\s\S]*?</script>`
	re := regexp.MustCompile(regex)
	pageHTML = []byte(re.ReplaceAllString(string(pageHTML), ""))
	return pageHTML
}
