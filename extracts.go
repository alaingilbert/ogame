package ogame

import (
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/alaingilbert/clockwork"
)

// Extractor ...
type Extractor interface {
	ExtractIsInVacation(pageHTML []byte) bool
	ExtractPlanets(pageHTML []byte, b *OGame) []Planet
	ExtractPlanet(pageHTML []byte, v interface{}, b *OGame) (Planet, error)
	ExtractPlanetByCoord(pageHTML []byte, b *OGame, coord Coordinate) (Planet, error)
	ExtractMoons(pageHTML []byte, b *OGame) []Moon
	ExtractMoon(pageHTML []byte, b *OGame, v interface{}) (Moon, error)
	ExtractMoonByCoord(pageHTML []byte, b *OGame, coord Coordinate) (Moon, error)
	ExtractCelestials(pageHTML []byte, b *OGame) ([]Celestial, error)
	ExtractCelestial(pageHTML []byte, b *OGame, v interface{}) (Celestial, error)
	ExtractServerTime(pageHTML []byte) (time.Time, error)
	ExtractFleetsFromEventList(pageHTML []byte) []Fleet
	ExtractIPM(pageHTML []byte) (duration, max int64, token string)
	ExtractFleets(pageHTML []byte) (res []Fleet)
	ExtractSlots(pageHTML []byte) Slots
	ExtractOgameTimestamp(pageHTML []byte) int64
	ExtractResources(pageHTML []byte) Resources
	ExtractResourcesDetailsFromFullPage(pageHTML []byte) ResourcesDetails
	ExtractResourceSettings(pageHTML []byte) (ResourceSettings, error)
	ExtractAttacks(pageHTML []byte) ([]AttackEvent, error)
	ExtractOfferOfTheDay(pageHTML []byte) (int64, string, PlanetResources, Multiplier, error)
	ExtractResourcesBuildings(pageHTML []byte) (ResourcesBuildings, error)
	ExtractDefense(pageHTML []byte) (DefensesInfos, error)
	ExtractShips(pageHTML []byte) (ShipsInfos, error)
	ExtractFacilities(pageHTML []byte) (Facilities, error)
	ExtractResearch(pageHTML []byte) Researches
	ExtractProduction(pageHTML []byte) ([]Quantifiable, int64, error)
	ExtractOverviewProduction(pageHTML []byte) ([]Quantifiable, int64, error)
	ExtractFleet1Ships(pageHTML []byte) ShipsInfos
	ExtractEspionageReportMessageIDs(pageHTML []byte) ([]EspionageReportSummary, int64)
	ExtractCombatReportMessagesSummary(pageHTML []byte) ([]CombatReportSummary, int64)
	ExtractEspionageReport(pageHTML []byte, location *time.Location) (EspionageReport, error)
	ExtractResourcesProductions(pageHTML []byte) (Resources, error)
	ExtractPreferences(pageHTML []byte) Preferences
	ExtractSpioAnz(pageHTML []byte) int64
	ExtractPreferencesShowActivityMinutes(pageHTML []byte) bool
	ExtractHiddenFields(pageHTML []byte) (fields url.Values)
	ExtractHiddenFieldsFromDoc(doc *goquery.Document) url.Values
	ExtractBodyIDFromDoc(doc *goquery.Document) string
	ExtractIsInVacationFromDoc(doc *goquery.Document) bool
	ExtractPlanetsFromDoc(doc *goquery.Document, b *OGame) []Planet
	ExtractPlanetByIDFromDoc(doc *goquery.Document, b *OGame, planetID PlanetID) (Planet, error)
	ExtractCelestialByIDFromDoc(doc *goquery.Document, b *OGame, celestialID CelestialID) (Celestial, error)
	ExtractPlanetByCoordFromDoc(doc *goquery.Document, b *OGame, coord Coordinate) (Planet, error)
	ExtractOgameTimestampFromDoc(doc *goquery.Document) int64
	ExtractResourcesFromDoc(doc *goquery.Document) Resources
	ExtractResourcesDetailsFromFullPageFromDoc(doc *goquery.Document) ResourcesDetails
	ExtractPlanetFromDoc(doc *goquery.Document, v interface{}, b *OGame) (Planet, error)
	ExtractMoonsFromDoc(doc *goquery.Document, b *OGame) []Moon
	ExtractMoonFromDoc(doc *goquery.Document, b *OGame, v interface{}) (Moon, error)
	ExtractMoonByCoordFromDoc(doc *goquery.Document, b *OGame, coord Coordinate) (Moon, error)
	ExtractMoonByIDFromDoc(doc *goquery.Document, b *OGame, moonID MoonID) (Moon, error)
	ExtractCelestialsFromDoc(doc *goquery.Document, b *OGame) ([]Celestial, error)
	ExtractCelestialFromDoc(doc *goquery.Document, b *OGame, v interface{}) (Celestial, error)
	ExtractResourcesBuildingsFromDoc(doc *goquery.Document) (ResourcesBuildings, error)
	ExtractDefenseFromDoc(doc *goquery.Document) (DefensesInfos, error)
	ExtractShipsFromDoc(doc *goquery.Document) (ShipsInfos, error)
	ExtractFacilitiesFromDoc(doc *goquery.Document) (Facilities, error)
	ExtractResearchFromDoc(doc *goquery.Document) Researches
	ExtractOGameSessionFromDoc(doc *goquery.Document) string
	ExtractAttacksFromDoc(doc *goquery.Document, clock clockwork.Clock) ([]AttackEvent, error)
	ExtractOfferOfTheDayFromDoc(doc *goquery.Document) (price int64, importToken string, planetResources PlanetResources, multiplier Multiplier, err error)
	ExtractProductionFromDoc(doc *goquery.Document) ([]Quantifiable, error)
	ExtractOverviewProductionFromDoc(doc *goquery.Document) ([]Quantifiable, error)
	ExtractFleet1ShipsFromDoc(doc *goquery.Document) (s ShipsInfos)
	ExtractEspionageReportMessageIDsFromDoc(doc *goquery.Document) ([]EspionageReportSummary, int64)
	ExtractCombatReportMessagesFromDoc(doc *goquery.Document) ([]CombatReportSummary, int64)
	ExtractEspionageReportFromDoc(doc *goquery.Document, location *time.Location) (EspionageReport, error)
	ExtractResourcesProductionsFromDoc(doc *goquery.Document) (Resources, error)
	ExtractPreferencesFromDoc(doc *goquery.Document) Preferences
	ExtractResourceSettingsFromDoc(doc *goquery.Document) (ResourceSettings, error)
	ExtractFleetsFromEventListFromDoc(doc *goquery.Document) []Fleet
	ExtractIPMFromDoc(doc *goquery.Document) (duration, max int64, token string)
	ExtractFleetsFromDoc(doc *goquery.Document) (res []Fleet)
	ExtractSlotsFromDoc(doc *goquery.Document) Slots
	ExtractServerTimeFromDoc(doc *goquery.Document) (time.Time, error)
	ExtractSpioAnzFromDoc(doc *goquery.Document) int64
	ExtractDisableChatBarFromDoc(doc *goquery.Document) bool
	ExtractDisableOutlawWarningFromDoc(doc *goquery.Document) bool
	ExtractMobileVersionFromDoc(doc *goquery.Document) bool
	ExtractShowOldDropDownsFromDoc(doc *goquery.Document) bool
	ExtractActivateAutofocusFromDoc(doc *goquery.Document) bool
	ExtractEventsShowFromDoc(doc *goquery.Document) int64
	ExtractSortSettingFromDoc(doc *goquery.Document) int64
	ExtractSortOrderFromDoc(doc *goquery.Document) int64
	ExtractShowDetailOverlayFromDoc(doc *goquery.Document) bool
	ExtractAnimatedSlidersFromDoc(doc *goquery.Document) bool
	ExtractAnimatedOverviewFromDoc(doc *goquery.Document) bool
	ExtractPopupsNoticesFromDoc(doc *goquery.Document) bool
	ExtractPopopsCombatreportFromDoc(doc *goquery.Document) bool
	ExtractSpioReportPicturesFromDoc(doc *goquery.Document) bool
	ExtractMsgResultsPerPageFromDoc(doc *goquery.Document) int64
	ExtractAuctioneerNotificationsFromDoc(doc *goquery.Document) bool
	ExtractEconomyNotificationsFromDoc(doc *goquery.Document) bool
	ExtractShowActivityMinutesFromDoc(doc *goquery.Document) bool
	ExtractPreserveSystemOnPlanetChangeFromDoc(doc *goquery.Document) bool
	ExtractNotifBuildListFromDoc(doc *goquery.Document) bool
	ExtractNotifFriendlyFleetActivitiesFromDoc(doc *goquery.Document) bool
	ExtractNotifHostileFleetActivitiesFromDoc(doc *goquery.Document) bool
	ExtractNotifForeignEspionageFromDoc(doc *goquery.Document) bool
	ExtractNotifAllianceBroadcastsFromDoc(doc *goquery.Document) bool
	ExtractNotifAllianceMessagesFromDoc(doc *goquery.Document) bool
	ExtractNotifAuctionsFromDoc(doc *goquery.Document) bool
	ExtractNotifAccountFromDoc(doc *goquery.Document) bool
	ExtractPlanetCoordinate(pageHTML []byte) (Coordinate, error)
	ExtractPlanetID(pageHTML []byte) (CelestialID, error)
	ExtractOverviewShipSumCountdownFromBytes(pageHTML []byte) int64
	ExtractOGameTimestampFromBytes(pageHTML []byte) int64
	ExtractPlanetType(pageHTML []byte) (CelestialType, error)
	ExtractAjaxChatToken(pageHTML []byte) (string, error)
	ExtractUserInfos(pageHTML []byte, lang string) (UserInfos, error)
	ExtractResourcesDetails(pageHTML []byte) (out ResourcesDetails, err error)
	ExtractCoord(v string) (coord Coordinate)
	ExtractGalaxyInfos(pageHTML []byte, botPlayerName string, botPlayerID, botPlayerRank int64) (SystemInfos, error)
	ExtractPhalanx(pageHTML []byte) ([]Fleet, error)
	ExtractJumpGate(pageHTML []byte) (ShipsInfos, string, []MoonID, int64)
	ExtractFederation(pageHTML []byte) url.Values
	ExtractConstructions(pageHTML []byte) (buildingID ID, buildingCountdown int64, researchID ID, researchCountdown int64)
	ExtractFleetDeutSaveFactor(pageHTML []byte) float64
	ExtractCancelBuildingInfos(pageHTML []byte) (token string, techID, listID int64, err error)
	ExtractCancelResearchInfos(pageHTML []byte) (token string, techID, listID int64, err error)
	ExtractEmpire(pageHTML []byte, nbr int64) (interface{}, error)
	ExtractCharacterClass(pageHTML []byte) (CharacterClass, error)
	ExtractCharacterClassFromDoc(doc *goquery.Document) (CharacterClass, error)
	ExtractCommander(pageHTML []byte) bool
	ExtractAdmiral(pageHTML []byte) bool
	ExtractEngineer(pageHTML []byte) bool
	ExtractGeologist(pageHTML []byte) bool
	ExtractTechnocrat(pageHTML []byte) bool
	ExtractCommanderFromDoc(doc *goquery.Document) bool
	ExtractAdmiralFromDoc(doc *goquery.Document) bool
	ExtractEngineerFromDoc(doc *goquery.Document) bool
	ExtractGeologistFromDoc(doc *goquery.Document) bool
	ExtractTechnocratFromDoc(doc *goquery.Document) bool
	ExtractAuction(pageHTML []byte) (Auction, error)
	ExtractHighscore(pageHTML []byte) (Highscore, error)
	ExtractHighscoreFromDoc(doc *goquery.Document) (Highscore, error)
	ExtractAllResources(pageHTML []byte) (map[CelestialID]Resources, error)
	ExtractDMCosts(pageHTML []byte) (DMCosts, error)
	ExtractBuffActivation(pageHTML []byte) (string, []Item, error)
	ExtractIsMobile(pageHTML []byte) bool
	ExtractIsMobileFromDoc(doc *goquery.Document) bool
}

// Compile time checks to ensure type satisfies Extractor interface
var _ Extractor = ExtractorV6{}
var _ Extractor = (*ExtractorV6)(nil)
var _ Extractor = ExtractorV7{}
var _ Extractor = (*ExtractorV7)(nil)

// extract universe speed from html calculation
// pageHTML := b.getPageContent(url.Values{"page": {"techtree"}, "tab": {"2"}, "techID": {"1"}})
func extractUniverseSpeed(pageHTML []byte) int64 {
	return extractUniverseSpeedV6(pageHTML)
}
