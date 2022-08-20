package extractor

import (
	"github.com/PuerkitoBio/goquery"
	v6 "github.com/alaingilbert/ogame/pkg/extractor/v6"
	v7 "github.com/alaingilbert/ogame/pkg/extractor/v7"
	v9 "github.com/alaingilbert/ogame/pkg/extractor/v9"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"net/url"
	"time"
)

// Extractor ...
type Extractor interface {
	GetLanguage() string
	SetLanguage(lang string)
	GetLocation() *time.Location
	SetLocation(loc *time.Location)
	GetLifeformEnabled() bool
	SetLifeformEnabled(lifeformEnabled bool)
	ExtractActiveItems(pageHTML []byte) ([]ogame.ActiveItem, error)
	ExtractLifeformEnabled(pageHTML []byte) bool
	ExtractAdmiral(pageHTML []byte) bool
	ExtractAjaxChatToken(pageHTML []byte) (string, error)
	ExtractAllResources(pageHTML []byte) (map[ogame.CelestialID]ogame.Resources, error)
	ExtractAttacks(pageHTML []byte, ownCoords []ogame.Coordinate) ([]ogame.AttackEvent, error)
	ExtractAuction(pageHTML []byte) (ogame.Auction, error)
	ExtractBuffActivation(pageHTML []byte) (string, []ogame.Item, error)
	ExtractCancelBuildingInfos(pageHTML []byte) (token string, techID, listID int64, err error)
	ExtractCancelLfBuildingInfos(pageHTML []byte) (token string, id, listID int64, err error)
	ExtractCancelFleetToken(pageHTML []byte, fleetID ogame.FleetID) (string, error)
	ExtractCancelResearchInfos(pageHTML []byte) (token string, techID, listID int64, err error)
	ExtractCelestial(pageHTML []byte, v any) (ogame.Celestial, error)
	ExtractCelestials(pageHTML []byte) ([]ogame.Celestial, error)
	ExtractCharacterClass(pageHTML []byte) (ogame.CharacterClass, error)
	ExtractCombatReportMessagesSummary(pageHTML []byte) ([]ogame.CombatReportSummary, int64)
	ExtractCommander(pageHTML []byte) bool
	ExtractConstructions(pageHTML []byte) (buildingID ogame.ID, buildingCountdown int64, researchID ogame.ID, researchCountdown int64)
	ExtractCoord(v string) (coord ogame.Coordinate)
	ExtractDMCosts(pageHTML []byte) (ogame.DMCosts, error)
	ExtractDefense(pageHTML []byte) (ogame.DefensesInfos, error)
	ExtractDestroyRockets(pageHTML []byte) (abm, ipm int64, token string, err error)
	ExtractEmpire(pageHTML []byte) ([]ogame.EmpireCelestial, error)
	ExtractEmpireJSON(pageHTML []byte) (any, error)
	ExtractEngineer(pageHTML []byte) bool
	ExtractEspionageReport(pageHTML []byte) (ogame.EspionageReport, error)
	ExtractEspionageReportMessageIDs(pageHTML []byte) ([]ogame.EspionageReportSummary, int64)
	ExtractExpeditionMessages(pageHTML []byte) ([]ogame.ExpeditionMessage, int64, error)
	ExtractFacilities(pageHTML []byte) (ogame.Facilities, error)
	ExtractFederation(pageHTML []byte) url.Values
	ExtractFleet1Ships(pageHTML []byte) ogame.ShipsInfos
	ExtractFleetDeutSaveFactor(pageHTML []byte) float64
	ExtractFleets(pageHTML []byte) (res []ogame.Fleet)
	ExtractFleetsFromEventList(pageHTML []byte) []ogame.Fleet
	ExtractGalaxyInfos(pageHTML []byte, botPlayerName string, botPlayerID, botPlayerRank int64) (ogame.SystemInfos, error)
	ExtractGeologist(pageHTML []byte) bool
	ExtractHiddenFields(pageHTML []byte) (fields url.Values)
	ExtractHighscore(pageHTML []byte) (ogame.Highscore, error)
	ExtractIPM(pageHTML []byte) (duration, max int64, token string)
	ExtractIsInVacation(pageHTML []byte) bool
	ExtractIsMobile(pageHTML []byte) bool
	ExtractJumpGate(pageHTML []byte) (ogame.ShipsInfos, string, []ogame.MoonID, int64)
	ExtractMarketplaceMessages(pageHTML []byte) ([]ogame.MarketplaceMessage, int64, error)
	ExtractMoon(pageHTML []byte, v any) (ogame.Moon, error)
	ExtractMoons(pageHTML []byte) []ogame.Moon
	ExtractOGameTimestampFromBytes(pageHTML []byte) int64
	ExtractOfferOfTheDay(pageHTML []byte) (int64, string, ogame.PlanetResources, ogame.Multiplier, error)
	ExtractOgameTimestamp(pageHTML []byte) int64
	ExtractOverviewProduction(pageHTML []byte) ([]ogame.Quantifiable, int64, error)
	ExtractOverviewShipSumCountdownFromBytes(pageHTML []byte) int64
	ExtractPhalanx(pageHTML []byte) ([]ogame.Fleet, error)
	ExtractPlanet(pageHTML []byte, v any) (ogame.Planet, error)
	ExtractPlanetCoordinate(pageHTML []byte) (ogame.Coordinate, error)
	ExtractPlanetID(pageHTML []byte) (ogame.CelestialID, error)
	ExtractPlanetType(pageHTML []byte) (ogame.CelestialType, error)
	ExtractPlanets(pageHTML []byte) []ogame.Planet
	ExtractPreferences(pageHTML []byte) ogame.Preferences
	ExtractPreferencesShowActivityMinutes(pageHTML []byte) bool
	ExtractPremiumToken(pageHTML []byte, days int64) (token string, err error)
	ExtractProduction(pageHTML []byte) ([]ogame.Quantifiable, int64, error)
	ExtractResearch(pageHTML []byte) ogame.Researches
	ExtractResourceSettings(pageHTML []byte) (ogame.ResourceSettings, error)
	ExtractResources(pageHTML []byte) ogame.Resources
	ExtractResourcesBuildings(pageHTML []byte) (ogame.ResourcesBuildings, error)
	ExtractResourcesDetails(pageHTML []byte) (out ogame.ResourcesDetails, err error)
	ExtractResourcesDetailsFromFullPage(pageHTML []byte) ogame.ResourcesDetails
	ExtractResourcesProductions(pageHTML []byte) (ogame.Resources, error)
	ExtractServerTime(pageHTML []byte) (time.Time, error)
	ExtractShips(pageHTML []byte) (ogame.ShipsInfos, error)
	ExtractSlots(pageHTML []byte) ogame.Slots
	ExtractSpioAnz(pageHTML []byte) int64
	ExtractTechnocrat(pageHTML []byte) bool
	ExtractTechs(pageHTML []byte) (ogame.ResourcesBuildings, ogame.Facilities, ogame.ShipsInfos, ogame.DefensesInfos, ogame.Researches, error)
	ExtractUserInfos(pageHTML []byte) (ogame.UserInfos, error)
	ExtractActivateAutofocusFromDoc(doc *goquery.Document) bool
	ExtractAdmiralFromDoc(doc *goquery.Document) bool
	ExtractAnimatedOverviewFromDoc(doc *goquery.Document) bool
	ExtractAnimatedSlidersFromDoc(doc *goquery.Document) bool
	ExtractAttacksFromDoc(doc *goquery.Document, ownCoords []ogame.Coordinate) ([]ogame.AttackEvent, error)
	ExtractAuctioneerNotificationsFromDoc(doc *goquery.Document) bool
	ExtractBodyIDFromDoc(doc *goquery.Document) string
	ExtractCelestialFromDoc(doc *goquery.Document, v any) (ogame.Celestial, error)
	ExtractCelestialsFromDoc(doc *goquery.Document) ([]ogame.Celestial, error)
	ExtractCharacterClassFromDoc(doc *goquery.Document) (ogame.CharacterClass, error)
	ExtractCombatReportMessagesFromDoc(doc *goquery.Document) ([]ogame.CombatReportSummary, int64)
	ExtractCommanderFromDoc(doc *goquery.Document) bool
	ExtractDefenseFromDoc(doc *goquery.Document) (ogame.DefensesInfos, error)
	ExtractDisableChatBarFromDoc(doc *goquery.Document) bool
	ExtractDisableOutlawWarningFromDoc(doc *goquery.Document) bool
	ExtractEconomyNotificationsFromDoc(doc *goquery.Document) bool
	ExtractEngineerFromDoc(doc *goquery.Document) bool
	ExtractEspionageReportFromDoc(doc *goquery.Document) (ogame.EspionageReport, error)
	ExtractEspionageReportMessageIDsFromDoc(doc *goquery.Document) ([]ogame.EspionageReportSummary, int64)
	ExtractEventsShowFromDoc(doc *goquery.Document) int64
	ExtractExpeditionMessagesFromDoc(doc *goquery.Document) ([]ogame.ExpeditionMessage, int64, error)
	ExtractFacilitiesFromDoc(doc *goquery.Document) (ogame.Facilities, error)
	ExtractFleet1ShipsFromDoc(doc *goquery.Document) (s ogame.ShipsInfos)
	ExtractFleetsFromDoc(doc *goquery.Document) (res []ogame.Fleet)
	ExtractFleetsFromEventListFromDoc(doc *goquery.Document) []ogame.Fleet
	ExtractGeologistFromDoc(doc *goquery.Document) bool
	ExtractHiddenFieldsFromDoc(doc *goquery.Document) url.Values
	ExtractHighscoreFromDoc(doc *goquery.Document) (ogame.Highscore, error)
	ExtractIPMFromDoc(doc *goquery.Document) (duration, max int64, token string)
	ExtractIsInVacationFromDoc(doc *goquery.Document) bool
	ExtractIsMobileFromDoc(doc *goquery.Document) bool
	ExtractMobileVersionFromDoc(doc *goquery.Document) bool
	ExtractMoonFromDoc(doc *goquery.Document, v any) (ogame.Moon, error)
	ExtractMoonsFromDoc(doc *goquery.Document) []ogame.Moon
	ExtractMsgResultsPerPageFromDoc(doc *goquery.Document) int64
	ExtractNotifAccountFromDoc(doc *goquery.Document) bool
	ExtractNotifAllianceBroadcastsFromDoc(doc *goquery.Document) bool
	ExtractNotifAllianceMessagesFromDoc(doc *goquery.Document) bool
	ExtractNotifAuctionsFromDoc(doc *goquery.Document) bool
	ExtractNotifBuildListFromDoc(doc *goquery.Document) bool
	ExtractNotifForeignEspionageFromDoc(doc *goquery.Document) bool
	ExtractNotifFriendlyFleetActivitiesFromDoc(doc *goquery.Document) bool
	ExtractNotifHostileFleetActivitiesFromDoc(doc *goquery.Document) bool
	ExtractOGameSessionFromDoc(doc *goquery.Document) string
	ExtractOfferOfTheDayFromDoc(doc *goquery.Document) (price int64, importToken string, planetResources ogame.PlanetResources, multiplier ogame.Multiplier, err error)
	ExtractOgameTimestampFromDoc(doc *goquery.Document) int64
	ExtractOverviewProductionFromDoc(doc *goquery.Document) ([]ogame.Quantifiable, error)
	ExtractPlanetFromDoc(doc *goquery.Document, v any) (ogame.Planet, error)
	ExtractPlanetIDFromDoc(doc *goquery.Document) (ogame.CelestialID, error)
	ExtractPlanetTypeFromDoc(doc *goquery.Document) (ogame.CelestialType, error)
	ExtractPlanetsFromDoc(doc *goquery.Document) []ogame.Planet
	ExtractPopopsCombatreportFromDoc(doc *goquery.Document) bool
	ExtractPopupsNoticesFromDoc(doc *goquery.Document) bool
	ExtractPreferencesFromDoc(doc *goquery.Document) ogame.Preferences
	ExtractPreserveSystemOnPlanetChangeFromDoc(doc *goquery.Document) bool
	ExtractProductionFromDoc(doc *goquery.Document) ([]ogame.Quantifiable, error)
	ExtractResearchFromDoc(doc *goquery.Document) ogame.Researches
	ExtractResourceSettingsFromDoc(doc *goquery.Document) (ogame.ResourceSettings, error)
	ExtractResourcesBuildingsFromDoc(doc *goquery.Document) (ogame.ResourcesBuildings, error)
	ExtractResourcesDetailsFromFullPageFromDoc(doc *goquery.Document) ogame.ResourcesDetails
	ExtractResourcesFromDoc(doc *goquery.Document) ogame.Resources
	ExtractResourcesProductionsFromDoc(doc *goquery.Document) (ogame.Resources, error)
	ExtractServerTimeFromDoc(doc *goquery.Document) (time.Time, error)
	ExtractShipsFromDoc(doc *goquery.Document) (ogame.ShipsInfos, error)
	ExtractShowActivityMinutesFromDoc(doc *goquery.Document) bool
	ExtractShowDetailOverlayFromDoc(doc *goquery.Document) bool
	ExtractShowOldDropDownsFromDoc(doc *goquery.Document) bool
	ExtractSlotsFromDoc(doc *goquery.Document) ogame.Slots
	ExtractSortOrderFromDoc(doc *goquery.Document) int64
	ExtractSortSettingFromDoc(doc *goquery.Document) int64
	ExtractSpioAnzFromDoc(doc *goquery.Document) int64
	ExtractSpioReportPicturesFromDoc(doc *goquery.Document) bool
	ExtractTechnocratFromDoc(doc *goquery.Document) bool
}

// Compile time checks to ensure type satisfies Extractor interface
var _ Extractor = (*v6.Extractor)(nil)
var _ Extractor = (*v7.Extractor)(nil)
var _ Extractor = (*v9.Extractor)(nil)
