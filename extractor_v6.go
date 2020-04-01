package ogame

import (
	"bytes"
	"errors"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/alaingilbert/clockwork"
)

// ExtractorV6 ...
type ExtractorV6 struct {
}

// NewExtractorV6 ...
func NewExtractorV6() *ExtractorV6 {
	return &ExtractorV6{}
}

// ExtractIsInVacation ...
func (e ExtractorV6) ExtractIsInVacation(pageHTML []byte) bool {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractIsInVacationFromDoc(doc)
}

// ExtractPlanets ...
func (e ExtractorV6) ExtractPlanets(pageHTML []byte, b *OGame) []Planet {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractPlanetsFromDoc(doc, b)
}

// ExtractPlanet ...
func (e ExtractorV6) ExtractPlanet(pageHTML []byte, v interface{}, b *OGame) (Planet, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractPlanetFromDoc(doc, v, b)
}

// ExtractPlanetByCoord ...
func (e ExtractorV6) ExtractPlanetByCoord(pageHTML []byte, b *OGame, coord Coordinate) (Planet, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractPlanetByCoordFromDoc(doc, b, coord)
}

// ExtractMoons ...
func (e ExtractorV6) ExtractMoons(pageHTML []byte, b *OGame) []Moon {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractMoonsFromDoc(doc, b)
}

// ExtractMoon ...
func (e ExtractorV6) ExtractMoon(pageHTML []byte, b *OGame, v interface{}) (Moon, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractMoonFromDoc(doc, b, v)
}

// ExtractMoonByCoord ...
func (e ExtractorV6) ExtractMoonByCoord(pageHTML []byte, b *OGame, coord Coordinate) (Moon, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractMoonByCoordFromDoc(doc, b, coord)
}

// ExtractCelestials ...
func (e ExtractorV6) ExtractCelestials(pageHTML []byte, b *OGame) ([]Celestial, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractCelestialsFromDoc(doc, b)
}

// ExtractCelestial ...
func (e ExtractorV6) ExtractCelestial(pageHTML []byte, b *OGame, v interface{}) (Celestial, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractCelestialFromDoc(doc, b, v)
}

// ExtractServerTime ...
func (e ExtractorV6) ExtractServerTime(pageHTML []byte) (time.Time, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractServerTimeFromDoc(doc)
}

// ExtractFleetsFromEventList ...
func (e ExtractorV6) ExtractFleetsFromEventList(pageHTML []byte) []Fleet {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractFleetsFromEventListFromDoc(doc)
}

// ExtractIPM ...
func (e ExtractorV6) ExtractIPM(pageHTML []byte) (duration, max int64, token string) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractIPMFromDoc(doc)
}

// ExtractFleets ...
func (e ExtractorV6) ExtractFleets(pageHTML []byte) (res []Fleet) {
	return e.extractFleets(pageHTML, clockwork.NewRealClock())
}

func (e ExtractorV6) extractFleets(pageHTML []byte, clock clockwork.Clock) (res []Fleet) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.extractFleetsFromDoc(doc, clock)
}

// ExtractSlots ...
func (e ExtractorV6) ExtractSlots(pageHTML []byte) Slots {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractSlotsFromDoc(doc)
}

// ExtractOgameTimestamp ...
func (e ExtractorV6) ExtractOgameTimestamp(pageHTML []byte) int64 {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractOgameTimestampFromDoc(doc)
}

// ExtractResources ...
func (e ExtractorV6) ExtractResources(pageHTML []byte) Resources {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractResourcesFromDoc(doc)
}

// ExtractResourcesDetailsFromFullPage ...
func (e ExtractorV6) ExtractResourcesDetailsFromFullPage(pageHTML []byte) ResourcesDetails {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractResourcesDetailsFromFullPageFromDoc(doc)
}

// ExtractResourceSettings ...
func (e ExtractorV6) ExtractResourceSettings(pageHTML []byte) (ResourceSettings, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractResourceSettingsFromDoc(doc)
}

// ExtractAttacks ...
func (e ExtractorV6) ExtractAttacks(pageHTML []byte) ([]AttackEvent, error) {
	return e.extractAttacks(pageHTML, clockwork.NewRealClock())
}

func (e ExtractorV6) extractAttacks(pageHTML []byte, clock clockwork.Clock) ([]AttackEvent, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractAttacksFromDoc(doc, clock)
}

// ExtractOfferOfTheDay ...
func (e ExtractorV6) ExtractOfferOfTheDay(pageHTML []byte) (int64, string, PlanetResources, Multiplier, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractOfferOfTheDayFromDoc(doc)
}

// ExtractResourcesBuildings ...
func (e ExtractorV6) ExtractResourcesBuildings(pageHTML []byte) (ResourcesBuildings, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractResourcesBuildingsFromDoc(doc)
}

// ExtractDefense ...
func (e ExtractorV6) ExtractDefense(pageHTML []byte) (DefensesInfos, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractDefenseFromDoc(doc)
}

// ExtractShips ...
func (e ExtractorV6) ExtractShips(pageHTML []byte) (ShipsInfos, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractShipsFromDoc(doc)
}

// ExtractFacilities ...
func (e ExtractorV6) ExtractFacilities(pageHTML []byte) (Facilities, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractFacilitiesFromDoc(doc)
}

// ExtractResearch ...
func (e ExtractorV6) ExtractResearch(pageHTML []byte) Researches {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractResearchFromDoc(doc)
}

// ExtractProduction extracts ships/defenses production from the shipyard page
func (e ExtractorV6) ExtractProduction(pageHTML []byte) ([]Quantifiable, int64, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	shipSumCountdown := e.ExtractOverviewShipSumCountdownFromBytes(pageHTML)
	production, err := e.ExtractProductionFromDoc(doc)
	return production, shipSumCountdown, err
}

// ExtractOverviewProduction extracts ships/defenses (partial) production from the overview page
func (e ExtractorV6) ExtractOverviewProduction(pageHTML []byte) ([]Quantifiable, int64, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	shipSumCountdown := e.ExtractOverviewShipSumCountdownFromBytes(pageHTML)
	production, err := e.ExtractOverviewProductionFromDoc(doc)
	return production, shipSumCountdown, err
}

// ExtractFleet1Ships ...
func (e ExtractorV6) ExtractFleet1Ships(pageHTML []byte) ShipsInfos {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractFleet1ShipsFromDoc(doc)
}

// ExtractEspionageReportMessageIDs ...
func (e ExtractorV6) ExtractEspionageReportMessageIDs(pageHTML []byte) ([]EspionageReportSummary, int64) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractEspionageReportMessageIDsFromDoc(doc)
}

// ExtractCombatReportMessagesSummary ...
func (e ExtractorV6) ExtractCombatReportMessagesSummary(pageHTML []byte) ([]CombatReportSummary, int64) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractCombatReportMessagesFromDoc(doc)
}

// ExtractEspionageReport ...
func (e ExtractorV6) ExtractEspionageReport(pageHTML []byte, location *time.Location) (EspionageReport, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractEspionageReportFromDoc(doc, location)
}

// ExtractResourcesProductions ...
func (e ExtractorV6) ExtractResourcesProductions(pageHTML []byte) (Resources, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractResourcesProductionsFromDoc(doc)
}

// ExtractPreferences ...
func (e ExtractorV6) ExtractPreferences(pageHTML []byte) Preferences {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractPreferencesFromDoc(doc)
}

// ExtractSpioAnz ...
func (e ExtractorV6) ExtractSpioAnz(pageHTML []byte) int64 {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractSpioAnzFromDoc(doc)
}

// ExtractPreferencesShowActivityMinutes ...
func (e ExtractorV6) ExtractPreferencesShowActivityMinutes(pageHTML []byte) bool {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractShowActivityMinutesFromDoc(doc)
}

// ExtractHiddenFields utils function to extract hidden input from a page
func (e ExtractorV6) ExtractHiddenFields(pageHTML []byte) (fields url.Values) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractHiddenFieldsFromDoc(doc)
}

// ExtractCommander ...
func (e ExtractorV6) ExtractCommander(pageHTML []byte) bool {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractCommanderFromDoc(doc)
}

// ExtractAdmiral ...
func (e ExtractorV6) ExtractAdmiral(pageHTML []byte) bool {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractAdmiralFromDoc(doc)
}

// ExtractEngineer ...
func (e ExtractorV6) ExtractEngineer(pageHTML []byte) bool {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractEngineerFromDoc(doc)
}

// ExtractGeologist ...
func (e ExtractorV6) ExtractGeologist(pageHTML []byte) bool {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractGeologistFromDoc(doc)
}

// ExtractTechnocrat ...
func (e ExtractorV6) ExtractTechnocrat(pageHTML []byte) bool {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractTechnocratFromDoc(doc)
}

// ExtractOGameSession ...
func (e ExtractorV6) ExtractOGameSession(pageHTML []byte) string {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractOGameSessionFromDoc(doc)
}

// <Extract from doc> ---------------------------------------------------------

// ExtractHiddenFieldsFromDoc utils function to extract hidden input from a page
func (e ExtractorV6) ExtractHiddenFieldsFromDoc(doc *goquery.Document) url.Values {
	return extractHiddenFieldsFromDocV6(doc)
}

// ExtractBodyIDFromDoc ...
func (e ExtractorV6) ExtractBodyIDFromDoc(doc *goquery.Document) string {
	return extractBodyIDFromDocV6(doc)
}

// ExtractIsInVacationFromDoc ...
func (e ExtractorV6) ExtractIsInVacationFromDoc(doc *goquery.Document) bool {
	return extractIsInVacationFromDocV6(doc)
}

// ExtractPlanetsFromDoc ...
func (e ExtractorV6) ExtractPlanetsFromDoc(doc *goquery.Document, b *OGame) []Planet {
	return extractPlanetsFromDocV6(doc, b)
}

// ExtractPlanetByIDFromDoc ...
func (e ExtractorV6) ExtractPlanetByIDFromDoc(doc *goquery.Document, b *OGame, planetID PlanetID) (Planet, error) {
	return extractPlanetByIDFromDocV6(doc, b, planetID)
}

// ExtractCelestialByIDFromDoc ...
func (e ExtractorV6) ExtractCelestialByIDFromDoc(doc *goquery.Document, b *OGame, celestialID CelestialID) (Celestial, error) {
	return extractCelestialByIDFromDocV6(doc, b, celestialID)
}

// ExtractPlanetByCoordFromDoc ...
func (e ExtractorV6) ExtractPlanetByCoordFromDoc(doc *goquery.Document, b *OGame, coord Coordinate) (Planet, error) {
	return extractPlanetByCoordFromDocV6(doc, b, coord)
}

// ExtractOgameTimestampFromDoc ...
func (e ExtractorV6) ExtractOgameTimestampFromDoc(doc *goquery.Document) int64 {
	return extractOgameTimestampFromDocV6(doc)
}

// ExtractResourcesFromDoc ...
func (e ExtractorV6) ExtractResourcesFromDoc(doc *goquery.Document) Resources {
	return extractResourcesFromDocV6(doc)
}

// ExtractResourcesDetailsFromFullPageFromDoc ...
func (e ExtractorV6) ExtractResourcesDetailsFromFullPageFromDoc(doc *goquery.Document) ResourcesDetails {
	return extractResourcesDetailsFromFullPageFromDocV6(doc)
}

// ExtractPlanetFromDoc ...
func (e ExtractorV6) ExtractPlanetFromDoc(doc *goquery.Document, v interface{}, b *OGame) (Planet, error) {
	return extractPlanetFromDocV6(doc, v, b)
}

// ExtractMoonsFromDoc ...
func (e ExtractorV6) ExtractMoonsFromDoc(doc *goquery.Document, b *OGame) []Moon {
	return extractMoonsFromDocV6(doc, b)
}

// ExtractMoonFromDoc ...
func (e ExtractorV6) ExtractMoonFromDoc(doc *goquery.Document, b *OGame, v interface{}) (Moon, error) {
	return extractMoonFromDocV6(doc, b, v)
}

// ExtractMoonByCoordFromDoc ...
func (e ExtractorV6) ExtractMoonByCoordFromDoc(doc *goquery.Document, b *OGame, coord Coordinate) (Moon, error) {
	return extractMoonByCoordFromDocV6(doc, b, coord)
}

// ExtractMoonByIDFromDoc ...
func (e ExtractorV6) ExtractMoonByIDFromDoc(doc *goquery.Document, b *OGame, moonID MoonID) (Moon, error) {
	return extractMoonByIDFromDocV6(doc, b, moonID)
}

// ExtractCelestialsFromDoc ...
func (e ExtractorV6) ExtractCelestialsFromDoc(doc *goquery.Document, b *OGame) ([]Celestial, error) {
	return extractCelestialsFromDocV6(doc, b)
}

// ExtractCelestialFromDoc ...
func (e ExtractorV6) ExtractCelestialFromDoc(doc *goquery.Document, b *OGame, v interface{}) (Celestial, error) {
	return extractCelestialFromDocV6(doc, b, v)
}

// ExtractResourcesBuildingsFromDoc ...
func (e ExtractorV6) ExtractResourcesBuildingsFromDoc(doc *goquery.Document) (ResourcesBuildings, error) {
	return extractResourcesBuildingsFromDocV6(doc)
}

// ExtractDefenseFromDoc ...
func (e ExtractorV6) ExtractDefenseFromDoc(doc *goquery.Document) (DefensesInfos, error) {
	return extractDefenseFromDocV6(doc)
}

// ExtractShipsFromDoc ...
func (e ExtractorV6) ExtractShipsFromDoc(doc *goquery.Document) (ShipsInfos, error) {
	return extractShipsFromDocV6(doc)
}

// ExtractFacilitiesFromDoc ...
func (e ExtractorV6) ExtractFacilitiesFromDoc(doc *goquery.Document) (Facilities, error) {
	return extractFacilitiesFromDocV6(doc)
}

// ExtractResearchFromDoc ...
func (e ExtractorV6) ExtractResearchFromDoc(doc *goquery.Document) Researches {
	return extractResearchFromDocV6(doc)
}

// ExtractOGameSessionFromDoc ...
func (e ExtractorV6) ExtractOGameSessionFromDoc(doc *goquery.Document) string {
	return extractOGameSessionFromDocV6(doc)
}

// ExtractAttacksFromDoc ...
func (e ExtractorV6) ExtractAttacksFromDoc(doc *goquery.Document, clock clockwork.Clock) ([]AttackEvent, error) {
	return extractAttacksFromDocV6(doc, clock)
}

// ExtractOfferOfTheDayFromDoc ...
func (e ExtractorV6) ExtractOfferOfTheDayFromDoc(doc *goquery.Document) (price int64, importToken string, planetResources PlanetResources, multiplier Multiplier, err error) {
	return extractOfferOfTheDayFromDocV6(doc)
}

// ExtractProductionFromDoc extracts ships/defenses production from the shipyard page
func (e ExtractorV6) ExtractProductionFromDoc(doc *goquery.Document) ([]Quantifiable, error) {
	return extractProductionFromDocV6(doc)
}

// ExtractOverviewProductionFromDoc extracts ships/defenses (partial) production from the overview page
func (e ExtractorV6) ExtractOverviewProductionFromDoc(doc *goquery.Document) ([]Quantifiable, error) {
	return extractOverviewProductionFromDocV6(doc)
}

// ExtractFleet1ShipsFromDoc ...
func (e ExtractorV6) ExtractFleet1ShipsFromDoc(doc *goquery.Document) (s ShipsInfos) {
	return extractFleet1ShipsFromDocV6(doc)
}

// ExtractEspionageReportMessageIDsFromDoc ...
func (e ExtractorV6) ExtractEspionageReportMessageIDsFromDoc(doc *goquery.Document) ([]EspionageReportSummary, int64) {
	return extractEspionageReportMessageIDsFromDocV6(doc)
}

// ExtractCombatReportMessagesFromDoc ...
func (e ExtractorV6) ExtractCombatReportMessagesFromDoc(doc *goquery.Document) ([]CombatReportSummary, int64) {
	return extractCombatReportMessagesFromDocV6(doc)
}

// ExtractEspionageReportFromDoc ...
func (e ExtractorV6) ExtractEspionageReportFromDoc(doc *goquery.Document, location *time.Location) (EspionageReport, error) {
	return extractEspionageReportFromDocV6(doc, location)
}

// ExtractResourcesProductionsFromDoc ...
func (e ExtractorV6) ExtractResourcesProductionsFromDoc(doc *goquery.Document) (Resources, error) {
	return extractResourcesProductionsFromDocV6(doc)
}

// ExtractPreferencesFromDoc ...
func (e ExtractorV6) ExtractPreferencesFromDoc(doc *goquery.Document) Preferences {
	return extractPreferencesFromDocV6(doc)
}

// ExtractResourceSettingsFromDoc ...
func (e ExtractorV6) ExtractResourceSettingsFromDoc(doc *goquery.Document) (ResourceSettings, error) {
	return extractResourceSettingsFromDocV6(doc)
}

// ExtractFleetsFromEventListFromDoc ...
func (e ExtractorV6) ExtractFleetsFromEventListFromDoc(doc *goquery.Document) []Fleet {
	return extractFleetsFromEventListFromDocV6(doc)
}

// ExtractIPMFromDoc ...
func (e ExtractorV6) ExtractIPMFromDoc(doc *goquery.Document) (duration, max int64, token string) {
	return extractIPMFromDocV6(doc)
}

// ExtractFleetsFromDoc ...
func (e ExtractorV6) ExtractFleetsFromDoc(doc *goquery.Document) (res []Fleet) {
	return e.extractFleetsFromDoc(doc, clockwork.NewRealClock())
}

func (e ExtractorV6) extractFleetsFromDoc(doc *goquery.Document, clock clockwork.Clock) (res []Fleet) {
	return extractFleetsFromDocV6(doc, clock)
}

// ExtractSlotsFromDoc extract fleet slots from page "fleet1"
// page "movement" redirect to "fleet1" when there is no fleet
func (e ExtractorV6) ExtractSlotsFromDoc(doc *goquery.Document) Slots {
	return extractSlotsFromDocV6(doc)
}

// ExtractServerTimeFromDoc ...
func (e ExtractorV6) ExtractServerTimeFromDoc(doc *goquery.Document) (time.Time, error) {
	return extractServerTimeFromDocV6(doc)
}

// ExtractSpioAnzFromDoc ...
func (e ExtractorV6) ExtractSpioAnzFromDoc(doc *goquery.Document) int64 {
	return extractSpioAnzFromDocV6(doc)
}

// ExtractDisableChatBarFromDoc ...
func (e ExtractorV6) ExtractDisableChatBarFromDoc(doc *goquery.Document) bool {
	return extractDisableChatBarFromDocV6(doc)
}

// ExtractDisableOutlawWarningFromDoc ...
func (e ExtractorV6) ExtractDisableOutlawWarningFromDoc(doc *goquery.Document) bool {
	return extractDisableOutlawWarningFromDocV6(doc)
}

// ExtractMobileVersionFromDoc ...
func (e ExtractorV6) ExtractMobileVersionFromDoc(doc *goquery.Document) bool {
	return extractMobileVersionFromDocV6(doc)
}

// ExtractShowOldDropDownsFromDoc ...
func (e ExtractorV6) ExtractShowOldDropDownsFromDoc(doc *goquery.Document) bool {
	return extractShowOldDropDownsFromDocV6(doc)
}

// ExtractActivateAutofocusFromDoc ...
func (e ExtractorV6) ExtractActivateAutofocusFromDoc(doc *goquery.Document) bool {
	return extractActivateAutofocusFromDocV6(doc)
}

// ExtractEventsShowFromDoc ...
func (e ExtractorV6) ExtractEventsShowFromDoc(doc *goquery.Document) int64 {
	return extractEventsShowFromDocV6(doc)
}

// ExtractSortSettingFromDoc ...
func (e ExtractorV6) ExtractSortSettingFromDoc(doc *goquery.Document) int64 {
	return extractSortSettingFromDocV6(doc)
}

// ExtractSortOrderFromDoc ...
func (e ExtractorV6) ExtractSortOrderFromDoc(doc *goquery.Document) int64 {
	return extractSortOrderFromDocV6(doc)
}

// ExtractShowDetailOverlayFromDoc ...
func (e ExtractorV6) ExtractShowDetailOverlayFromDoc(doc *goquery.Document) bool {
	return extractShowDetailOverlayFromDocV6(doc)
}

// ExtractAnimatedSlidersFromDoc ...
func (e ExtractorV6) ExtractAnimatedSlidersFromDoc(doc *goquery.Document) bool {
	return extractAnimatedSlidersFromDocV6(doc)
}

// ExtractAnimatedOverviewFromDoc ...
func (e ExtractorV6) ExtractAnimatedOverviewFromDoc(doc *goquery.Document) bool {
	return extractAnimatedOverviewFromDocV6(doc)
}

// ExtractPopupsNoticesFromDoc ...
func (e ExtractorV6) ExtractPopupsNoticesFromDoc(doc *goquery.Document) bool {
	return extractPopupsNoticesFromDocV6(doc)
}

// ExtractPopopsCombatreportFromDoc ...
func (e ExtractorV6) ExtractPopopsCombatreportFromDoc(doc *goquery.Document) bool {
	return extractPopopsCombatreportFromDocV6(doc)
}

// ExtractSpioReportPicturesFromDoc ...
func (e ExtractorV6) ExtractSpioReportPicturesFromDoc(doc *goquery.Document) bool {
	return extractSpioReportPicturesFromDocV6(doc)
}

// ExtractMsgResultsPerPageFromDoc ...
func (e ExtractorV6) ExtractMsgResultsPerPageFromDoc(doc *goquery.Document) int64 {
	return extractMsgResultsPerPageFromDocV6(doc)
}

// ExtractAuctioneerNotificationsFromDoc ...
func (e ExtractorV6) ExtractAuctioneerNotificationsFromDoc(doc *goquery.Document) bool {
	return extractAuctioneerNotificationsFromDocV6(doc)
}

// ExtractEconomyNotificationsFromDoc ...
func (e ExtractorV6) ExtractEconomyNotificationsFromDoc(doc *goquery.Document) bool {
	return extractEconomyNotificationsFromDocV6(doc)
}

// ExtractShowActivityMinutesFromDoc ...
func (e ExtractorV6) ExtractShowActivityMinutesFromDoc(doc *goquery.Document) bool {
	return extractShowActivityMinutesFromDocV6(doc)
}

// ExtractPreserveSystemOnPlanetChangeFromDoc ...
func (e ExtractorV6) ExtractPreserveSystemOnPlanetChangeFromDoc(doc *goquery.Document) bool {
	return extractPreserveSystemOnPlanetChangeFromDocV6(doc)
}

// ExtractNotifBuildListFromDoc ...
func (e ExtractorV6) ExtractNotifBuildListFromDoc(doc *goquery.Document) bool {
	return extractNotifBuildListFromDocV6(doc)
}

// ExtractNotifFriendlyFleetActivitiesFromDoc ...
func (e ExtractorV6) ExtractNotifFriendlyFleetActivitiesFromDoc(doc *goquery.Document) bool {
	return extractNotifFriendlyFleetActivitiesFromDocV6(doc)
}

// ExtractNotifHostileFleetActivitiesFromDoc ...
func (e ExtractorV6) ExtractNotifHostileFleetActivitiesFromDoc(doc *goquery.Document) bool {
	return extractNotifHostileFleetActivitiesFromDocV6(doc)
}

// ExtractNotifForeignEspionageFromDoc ...
func (e ExtractorV6) ExtractNotifForeignEspionageFromDoc(doc *goquery.Document) bool {
	return extractNotifForeignEspionageFromDocV6(doc)
}

// ExtractNotifAllianceBroadcastsFromDoc ...
func (e ExtractorV6) ExtractNotifAllianceBroadcastsFromDoc(doc *goquery.Document) bool {
	return extractNotifAllianceBroadcastsFromDocV6(doc)
}

// ExtractNotifAllianceMessagesFromDoc ...
func (e ExtractorV6) ExtractNotifAllianceMessagesFromDoc(doc *goquery.Document) bool {
	return extractNotifAllianceMessagesFromDocV6(doc)
}

// ExtractNotifAuctionsFromDoc ...
func (e ExtractorV6) ExtractNotifAuctionsFromDoc(doc *goquery.Document) bool {
	return extractNotifAuctionsFromDocV6(doc)
}

// ExtractNotifAccountFromDoc ...
func (e ExtractorV6) ExtractNotifAccountFromDoc(doc *goquery.Document) bool {
	return extractNotifAccountFromDocV6(doc)
}

// ExtractCharacterClassFromDoc ...
func (e ExtractorV6) ExtractCharacterClassFromDoc(doc *goquery.Document) (CharacterClass, error) {
	return 0, errors.New("character class not supported in v6")
}

// ExtractCommanderFromDoc ...
func (e ExtractorV6) ExtractCommanderFromDoc(doc *goquery.Document) bool {
	return extractCommanderFromDocV6(doc)
}

// ExtractAdmiralFromDoc ...
func (e ExtractorV6) ExtractAdmiralFromDoc(doc *goquery.Document) bool {
	return extractAdmiralFromDocV6(doc)
}

// ExtractEngineerFromDoc ...
func (e ExtractorV6) ExtractEngineerFromDoc(doc *goquery.Document) bool {
	return extractEngineerFromDocV6(doc)
}

// ExtractGeologistFromDoc ...
func (e ExtractorV6) ExtractGeologistFromDoc(doc *goquery.Document) bool {
	return extractGeologistFromDocV6(doc)
}

// ExtractTechnocratFromDoc ...
func (e ExtractorV6) ExtractTechnocratFromDoc(doc *goquery.Document) bool {
	return extractTechnocratFromDocV6(doc)
}

// </ Extract from doc> -------------------------------------------------------

// <Works with []byte only> ---------------------------------------------------

// ExtractPlanetCoordinate extracts planet coordinate from html page
func (e ExtractorV6) ExtractPlanetCoordinate(pageHTML []byte) (Coordinate, error) {
	return extractPlanetCoordinateV6(pageHTML)
}

// ExtractPlanetID extracts planet id from html page
func (e ExtractorV6) ExtractPlanetID(pageHTML []byte) (CelestialID, error) {
	return extractPlanetIDV6(pageHTML)
}

// ExtractOverviewShipSumCountdownFromBytes ...
func (e ExtractorV6) ExtractOverviewShipSumCountdownFromBytes(pageHTML []byte) int64 {
	return extractOverviewShipSumCountdownFromBytesV6(pageHTML)
}

// ExtractOGameTimestampFromBytes extracts ogame timestamp from an html page
func (e ExtractorV6) ExtractOGameTimestampFromBytes(pageHTML []byte) int64 {
	return extractOGameTimestampFromBytesV6(pageHTML)
}

// ExtractPlanetType extracts planet type from html page
func (e ExtractorV6) ExtractPlanetType(pageHTML []byte) (CelestialType, error) {
	return extractPlanetTypeV6(pageHTML)
}

// ExtractAjaxChatToken ...
func (e ExtractorV6) ExtractAjaxChatToken(pageHTML []byte) (string, error) {
	return extractAjaxChatTokenV6(pageHTML)
}

// ExtractUserInfos ...
func (e ExtractorV6) ExtractUserInfos(pageHTML []byte, lang string) (UserInfos, error) {
	return extractUserInfosV6(pageHTML, lang)
}

// ExtractResourcesDetails ...
func (e ExtractorV6) ExtractResourcesDetails(pageHTML []byte) (out ResourcesDetails, err error) {
	return extractResourcesDetailsV6(pageHTML)
}

// </Works with []byte only> --------------------------------------------------

// ExtractCoord ...
func (e ExtractorV6) ExtractCoord(v string) (coord Coordinate) {
	return extractCoordV6(v)
}

// ExtractGalaxyInfos ...
func (e ExtractorV6) ExtractGalaxyInfos(pageHTML []byte, botPlayerName string, botPlayerID, botPlayerRank int64) (SystemInfos, error) {
	return extractGalaxyInfosV6(pageHTML, botPlayerName, botPlayerID, botPlayerRank)
}

// ExtractPhalanx ...
func (e ExtractorV6) ExtractPhalanx(pageHTML []byte) ([]Fleet, error) {
	return extractPhalanxV6(pageHTML)
}

// ExtractJumpGate return the available ships to send, form token, possible moon IDs and wait time (if any)
// given a jump gate popup html.
func (e ExtractorV6) ExtractJumpGate(pageHTML []byte) (ShipsInfos, string, []MoonID, int64) {
	return extractJumpGateV6(pageHTML)
}

// ExtractFederation ...
func (e ExtractorV6) ExtractFederation(pageHTML []byte) url.Values {
	return extractFederationV6(pageHTML)
}

// ExtractConstructions ...
func (e ExtractorV6) ExtractConstructions(pageHTML []byte) (buildingID ID, buildingCountdown int64, researchID ID, researchCountdown int64) {
	return extractConstructionsV6(pageHTML)
}

// ExtractFleetDeutSaveFactor extract fleet deut save factor
func (e ExtractorV6) ExtractFleetDeutSaveFactor(pageHTML []byte) float64 {
	return extractFleetDeutSaveFactorV6(pageHTML)
}

// ExtractCancelBuildingInfos ...
func (e ExtractorV6) ExtractCancelBuildingInfos(pageHTML []byte) (token string, techID, listID int64, err error) {
	return extractCancelBuildingInfosV6(pageHTML)
}

// ExtractCancelResearchInfos ...
func (e ExtractorV6) ExtractCancelResearchInfos(pageHTML []byte) (token string, techID, listID int64, err error) {
	return extractCancelResearchInfosV6(pageHTML)
}

// ExtractEmpire ...
func (e ExtractorV6) ExtractEmpire(pageHTML []byte, nbr int64) (interface{}, error) {
	return extractEmpire(string(pageHTML), nbr)
}

// ExtractCharacterClass ...
func (e ExtractorV6) ExtractCharacterClass(pageHTML []byte) (CharacterClass, error) {
	return 0, errors.New("character class not supported in v6")
}

// ExtractAuction ...
func (e ExtractorV6) ExtractAuction(pageHTML []byte) (Auction, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return extractAuctionFromDoc(doc)
}

// ExtractHighscore ...
func (e ExtractorV6) ExtractHighscore(pageHTML []byte) (Highscore, error) {
	panic("not implemented")
}

// ExtractHighscoreFromDoc ...
func (e ExtractorV6) ExtractHighscoreFromDoc(doc *goquery.Document) (Highscore, error) {
	panic("not implemented")
}

// ExtractAllResources ...
func (e ExtractorV6) ExtractAllResources(pageHTML []byte) (map[CelestialID]Resources, error) {
	panic("not implemented")
}

// ExtractDMCosts ...
func (e ExtractorV6) ExtractDMCosts(pageHTML []byte) (DMCosts, error) {
	panic("not implemented")
}

// ExtractBuffActivation ...
func (e ExtractorV6) ExtractBuffActivation(pageHTML []byte) (string, []Item, error) {
	panic("not implemented")
}

// ExtractIsMobile ...
func (e ExtractorV6) ExtractIsMobile(pageHTML []byte) bool {
	panic("not implemented")
}

// ExtractIsMobileFromDoc ...
func (e ExtractorV6) ExtractIsMobileFromDoc(doc *goquery.Document) bool {
	panic("not implemented")
}
