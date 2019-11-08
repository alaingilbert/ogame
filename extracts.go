package ogame

import (
	"bytes"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/alaingilbert/clockwork"
)

type extractor interface {
	ExtractIsInVacation(pageHTML []byte) bool
	ExtractPlanets(pageHTML []byte, b *OGame) []Planet
	ExtractPlanet(pageHTML []byte, v interface{}, b *OGame) (Planet, error)
	ExtractPlanetByCoord(pageHTML []byte, b *OGame, coord Coordinate) (Planet, error)
	ExtractMoons(pageHTML []byte, b *OGame) []Moon
	ExtractMoon(pageHTML []byte, b *OGame, v interface{}) (Moon, error)
	ExtractMoonByCoord(pageHTML []byte, b *OGame, coord Coordinate) (Moon, error)
}

// ExtractIsInVacation ...
func ExtractIsInVacation(pageHTML []byte) bool {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return ExtractIsInVacationFromDoc(doc)
}

// ExtractPlanets ...
func ExtractPlanets(pageHTML []byte, b *OGame) []Planet {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return ExtractPlanetsFromDoc(doc, b)
}

// ExtractPlanet ...
func ExtractPlanet(pageHTML []byte, v interface{}, b *OGame) (Planet, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return ExtractPlanetFromDoc(doc, v, b)
}

// ExtractPlanetByCoord ...
func ExtractPlanetByCoord(pageHTML []byte, b *OGame, coord Coordinate) (Planet, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return ExtractPlanetByCoordFromDoc(doc, b, coord)
}

// ExtractMoons ...
func ExtractMoons(pageHTML []byte, b *OGame) []Moon {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return ExtractMoonsFromDoc(doc, b)
}

// ExtractMoon ...
func ExtractMoon(pageHTML []byte, b *OGame, v interface{}) (Moon, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return ExtractMoonFromDoc(doc, b, v)
}

// ExtractMoonByCoord ...
func ExtractMoonByCoord(pageHTML []byte, b *OGame, coord Coordinate) (Moon, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return ExtractMoonByCoordFromDoc(doc, b, coord)
}

// ExtractCelestials ...
func ExtractCelestials(pageHTML []byte, b *OGame) ([]Celestial, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return ExtractCelestialsFromDoc(doc, b)
}

// ExtractCelestial ...
func ExtractCelestial(pageHTML []byte, b *OGame, v interface{}) (Celestial, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return ExtractCelestialFromDoc(doc, b, v)
}

func extractServerTime(pageHTML []byte) (time.Time, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return extractServerTimeFromDoc(doc)
}

// ExtractFleetsFromEventList ...
func ExtractFleetsFromEventList(pageHTML []byte) []Fleet {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return ExtractFleetsFromEventListFromDoc(doc)
}

// ExtractIPM ...
func ExtractIPM(pageHTML []byte) (duration, max int, token string) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return ExtractIPMFromDoc(doc)
}

// ExtractFleets ...
func ExtractFleets(pageHTML []byte) (res []Fleet) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return ExtractFleetsFromDoc(doc)
}

// ExtractSlots ...
func ExtractSlots(pageHTML []byte) Slots {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return ExtractSlotsFromDoc(doc)
}

// ExtractOgameTimestamp ...
func ExtractOgameTimestamp(pageHTML []byte) int64 {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return ExtractOgameTimestampFromDoc(doc)
}

// ExtractResources ...
func ExtractResources(pageHTML []byte) Resources {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return ExtractResourcesFromDoc(doc)
}

// ExtractResourcesDetailsFromFullPage ...
func ExtractResourcesDetailsFromFullPage(pageHTML []byte) ResourcesDetails {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return ExtractResourcesDetailsFromFullPageFromDoc(doc)
}

// ExtractResourceSettings ...
func ExtractResourceSettings(pageHTML []byte) (ResourceSettings, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return ExtractResourceSettingsFromDoc(doc)
}

// ExtractAttacks ...
func ExtractAttacks(pageHTML []byte) ([]AttackEvent, error) {
	return extractAttacks(pageHTML, clockwork.NewRealClock())
}

func extractAttacks(pageHTML []byte, clock clockwork.Clock) ([]AttackEvent, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return ExtractAttacksFromDoc(doc, clock)
}

// ExtractOfferOfTheDay ...
func ExtractOfferOfTheDay(pageHTML []byte) (int, string, PlanetResources, Multiplier, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return ExtractOfferOfTheDayFromDoc(doc)
}

// ExtractResourcesBuildings ...
func ExtractResourcesBuildings(pageHTML []byte) (ResourcesBuildings, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return ExtractResourcesBuildingsFromDoc(doc)
}

// ExtractDefense ...
func ExtractDefense(pageHTML []byte) (DefensesInfos, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return ExtractDefenseFromDoc(doc)
}

// ExtractShips ...
func ExtractShips(pageHTML []byte) (ShipsInfos, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return ExtractShipsFromDoc(doc)
}

// ExtractFacilities ...
func ExtractFacilities(pageHTML []byte) (Facilities, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return ExtractFacilitiesFromDoc(doc)
}

// ExtractResearch ...
func ExtractResearch(pageHTML []byte) Researches {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return ExtractResearchFromDoc(doc)
}

// ExtractProduction extracts ships/defenses production from the shipyard page
func ExtractProduction(pageHTML []byte) ([]Quantifiable, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return ExtractProductionFromDoc(doc)
}

// ExtractOverviewProduction extracts ships/defenses (partial) production from the overview page
func ExtractOverviewProduction(pageHTML []byte) ([]Quantifiable, int, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	shipSumCountdown := ExtractOverviewShipSumCountdownFromBytes(pageHTML)
	production, err := ExtractOverviewProductionFromDoc(doc)
	return production, shipSumCountdown, err
}

// ExtractFleet1Ships ...
func ExtractFleet1Ships(pageHTML []byte) ShipsInfos {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return ExtractFleet1ShipsFromDoc(doc)
}

func extractEspionageReportMessageIDs(pageHTML []byte) ([]EspionageReportSummary, int) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return extractEspionageReportMessageIDsFromDoc(doc)
}

func extractCombatReportMessagesSummary(pageHTML []byte) ([]CombatReportSummary, int) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return extractCombatReportMessagesFromDoc(doc)
}

func extractEspionageReport(pageHTML []byte, location *time.Location) (EspionageReport, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return extractEspionageReportFromDoc(doc, location)
}

// ExtractResourcesProductions ...
func ExtractResourcesProductions(pageHTML []byte) (Resources, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return ExtractResourcesProductionsFromDoc(doc)
}

// ExtractPreferences ...
func ExtractPreferences(pageHTML []byte) Preferences {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return ExtractPreferencesFromDoc(doc)
}

// ExtractSpioAnz ...
func ExtractSpioAnz(pageHTML []byte) int {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return ExtractSpioAnzFromDoc(doc)
}

// ExtractNbProbes ...
func ExtractPreferencesShowActivityMinutes(pageHTML []byte) bool {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return ExtractShowActivityMinutesFromDoc(doc)
}

// ExtractHiddenFields utils function to extract hidden input from a page
func ExtractHiddenFields(pageHTML []byte) (fields url.Values) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return ExtractHiddenFieldsFromDoc(doc)
}

// <Extract from doc> ---------------------------------------------------------

// ExtractHiddenFieldsFromDoc utils function to extract hidden input from a page
func ExtractHiddenFieldsFromDoc(doc *goquery.Document) url.Values {
	return extractHiddenFieldsFromDocV6(doc)
}

// ExtractBodyIDFromDoc ...
func ExtractBodyIDFromDoc(doc *goquery.Document) string {
	return extractBodyIDFromDocV6(doc)
}

// ExtractIsInVacationFromDoc ...
func ExtractIsInVacationFromDoc(doc *goquery.Document) bool {
	return extractIsInVacationFromDocV6(doc)
}

// ExtractPlanetsFromDoc ...
func ExtractPlanetsFromDoc(doc *goquery.Document, b *OGame) []Planet {
	return extractPlanetsFromDocV6(doc, b)
}

// ExtractPlanetByIDFromDoc ...
func ExtractPlanetByIDFromDoc(doc *goquery.Document, b *OGame, planetID PlanetID) (Planet, error) {
	return extractPlanetByIDFromDocV6(doc, b, planetID)
}

// ExtractCelestialByIDFromDoc ...
func ExtractCelestialByIDFromDoc(doc *goquery.Document, b *OGame, celestialID CelestialID) (Celestial, error) {
	return extractCelestialByIDFromDocV6(doc, b, celestialID)
}

// ExtractPlanetByCoordFromDoc ...
func ExtractPlanetByCoordFromDoc(doc *goquery.Document, b *OGame, coord Coordinate) (Planet, error) {
	return extractPlanetByCoordFromDocV6(doc, b, coord)
}

// ExtractOgameTimestampFromDoc ...
func ExtractOgameTimestampFromDoc(doc *goquery.Document) int64 {
	return extractOgameTimestampFromDocV6(doc)
}

// ExtractResourcesFromDoc ...
func ExtractResourcesFromDoc(doc *goquery.Document) Resources {
	return extractResourcesFromDocV6(doc)
}

// ExtractResourcesDetailsFromFullPageFromDoc ...
func ExtractResourcesDetailsFromFullPageFromDoc(doc *goquery.Document) ResourcesDetails {
	return extractResourcesDetailsFromFullPageFromDocV6(doc)
}

// ExtractPlanetFromDoc ...
func ExtractPlanetFromDoc(doc *goquery.Document, v interface{}, b *OGame) (Planet, error) {
	return extractPlanetFromDocV6(doc, v, b)
}

// ExtractMoonsFromDoc ...
func ExtractMoonsFromDoc(doc *goquery.Document, b *OGame) []Moon {
	return extractMoonsFromDocV6(doc, b)
}

// ExtractMoonFromDoc ...
func ExtractMoonFromDoc(doc *goquery.Document, b *OGame, v interface{}) (Moon, error) {
	return extractMoonFromDocV6(doc, b, v)
}

// ExtractMoonByCoordFromDoc ...
func ExtractMoonByCoordFromDoc(doc *goquery.Document, b *OGame, coord Coordinate) (Moon, error) {
	return extractMoonByCoordFromDocV6(doc, b, coord)
}

// ExtractMoonByIDFromDoc ...
func ExtractMoonByIDFromDoc(doc *goquery.Document, b *OGame, moonID MoonID) (Moon, error) {
	return extractMoonByIDFromDocV6(doc, b, moonID)
}

// ExtractCelestialsFromDoc ...
func ExtractCelestialsFromDoc(doc *goquery.Document, b *OGame) ([]Celestial, error) {
	return extractCelestialsFromDocV6(doc, b)
}

// ExtractCelestialFromDoc ...
func ExtractCelestialFromDoc(doc *goquery.Document, b *OGame, v interface{}) (Celestial, error) {
	return extractCelestialFromDocV6(doc, b, v)
}

// ExtractResourcesBuildingsFromDoc ...
func ExtractResourcesBuildingsFromDoc(doc *goquery.Document) (ResourcesBuildings, error) {
	return extractResourcesBuildingsFromDocV6(doc)
}

// ExtractDefenseFromDoc ...
func ExtractDefenseFromDoc(doc *goquery.Document) (DefensesInfos, error) {
	return extractDefenseFromDocV6(doc)
}

// ExtractShipsFromDoc ...
func ExtractShipsFromDoc(doc *goquery.Document) (ShipsInfos, error) {
	return extractShipsFromDocV6(doc)
}

// ExtractFacilitiesFromDoc ...
func ExtractFacilitiesFromDoc(doc *goquery.Document) (Facilities, error) {
	return extractFacilitiesFromDocV6(doc)
}

// ExtractResearchFromDoc ...
func ExtractResearchFromDoc(doc *goquery.Document) Researches {
	return extractResearchFromDocV6(doc)
}

// ExtractOGameSessionFromDoc ...
func ExtractOGameSessionFromDoc(doc *goquery.Document) string {
	return extractOGameSessionFromDocV6(doc)
}

// ExtractAttacksFromDoc ...
func ExtractAttacksFromDoc(doc *goquery.Document, clock clockwork.Clock) ([]AttackEvent, error) {
	return extractAttacksFromDocV6(doc, clock)
}

type planetResource struct {
	Input struct {
		Metal     int
		Crystal   int
		Deuterium int
	}
	Output struct {
		Metal     int
		Crystal   int
		Deuterium int
	}
	IsMoon        bool
	ImageFileName string
	Name          string
	OtherPlanet   string
}

type PlanetResources map[CelestialID]planetResource

type Multiplier struct {
	Metal     float64
	Crystal   float64
	Deuterium float64
	Honor     float64
}

// ExtractOfferOfTheDayFromDoc ...
func ExtractOfferOfTheDayFromDoc(doc *goquery.Document) (price int, importToken string, planetResources PlanetResources, multiplier Multiplier, err error) {
	return extractOfferOfTheDayFromDocV6(doc)
}

// ExtractProductionFromDoc extracts ships/defenses production from the shipyard page
func ExtractProductionFromDoc(doc *goquery.Document) ([]Quantifiable, error) {
	return extractProductionFromDocV6(doc)
}

// ExtractOverviewProductionFromDoc extracts ships/defenses (partial) production from the overview page
func ExtractOverviewProductionFromDoc(doc *goquery.Document) ([]Quantifiable, error) {
	return extractOverviewProductionFromDocV6(doc)
}

// ExtractFleet1ShipsFromDoc ...
func ExtractFleet1ShipsFromDoc(doc *goquery.Document) (s ShipsInfos) {
	return extractFleet1ShipsFromDocV6(doc)
}

func extractEspionageReportMessageIDsFromDoc(doc *goquery.Document) ([]EspionageReportSummary, int) {
	return extractEspionageReportMessageIDsFromDocV6(doc)
}

func extractCombatReportMessagesFromDoc(doc *goquery.Document) ([]CombatReportSummary, int) {
	return extractCombatReportMessagesFromDocV6(doc)
}

func extractEspionageReportFromDoc(doc *goquery.Document, location *time.Location) (EspionageReport, error) {
	return extractEspionageReportFromDocV6(doc, location)
}

// ExtractResourcesProductionsFromDoc ...
func ExtractResourcesProductionsFromDoc(doc *goquery.Document) (Resources, error) {
	return extractResourcesProductionsFromDocV6(doc)
}

// ExtractPreferencesFromDoc ...
func ExtractPreferencesFromDoc(doc *goquery.Document) Preferences {
	return extractPreferencesFromDocV6(doc)
}

// ExtractResourceSettingsFromDoc ...
func ExtractResourceSettingsFromDoc(doc *goquery.Document) (ResourceSettings, error) {
	return extractResourceSettingsFromDocV6(doc)
}

// ExtractFleetsFromEventListFromDoc ...
func ExtractFleetsFromEventListFromDoc(doc *goquery.Document) []Fleet {
	return extractFleetsFromEventListFromDocV6(doc)
}

// ExtractIPMFromDoc ...
func ExtractIPMFromDoc(doc *goquery.Document) (duration, max int, token string) {
	return extractIPMFromDocV6(doc)
}

// ExtractFleetsFromDoc ...
func ExtractFleetsFromDoc(doc *goquery.Document) (res []Fleet) {
	return extractFleetsFromDocV6(doc)
}

// ExtractSlotsFromDoc extract fleet slots from page "fleet1"
// page "movement" redirect to "fleet1" when there is no fleet
func ExtractSlotsFromDoc(doc *goquery.Document) Slots {
	return extractSlotsFromDocV6(doc)
}

func extractServerTimeFromDoc(doc *goquery.Document) (time.Time, error) {
	return extractServerTimeFromDocV6(doc)
}

// ExtractSpioAnzFromDoc ...
func ExtractSpioAnzFromDoc(doc *goquery.Document) int {
	return extractSpioAnzFromDocV6(doc)
}

// ExtractDisableChatBarFromDoc ...
func ExtractDisableChatBarFromDoc(doc *goquery.Document) bool {
	return extractDisableChatBarFromDocV6(doc)
}

// ExtractDisableOutlawWarningFromDoc ...
func ExtractDisableOutlawWarningFromDoc(doc *goquery.Document) bool {
	return extractDisableOutlawWarningFromDocV6(doc)
}

// ExtractMobileVersionFromDoc ...
func ExtractMobileVersionFromDoc(doc *goquery.Document) bool {
	return extractMobileVersionFromDocV6(doc)
}

// ExtractShowOldDropDownsFromDoc ...
func ExtractShowOldDropDownsFromDoc(doc *goquery.Document) bool {
	return extractShowOldDropDownsFromDocV6(doc)
}

// ExtractActivateAutofocusFromDoc ...
func ExtractActivateAutofocusFromDoc(doc *goquery.Document) bool {
	return extractActivateAutofocusFromDocV6(doc)
}

// ExtractEventsShowFromDoc ...
func ExtractEventsShowFromDoc(doc *goquery.Document) int {
	return extractEventsShowFromDocV6(doc)
}

// ExtractSortSettingFromDoc ...
func ExtractSortSettingFromDoc(doc *goquery.Document) int {
	return extractSortSettingFromDocV6(doc)
}

// ExtractSortOrderFromDoc ...
func ExtractSortOrderFromDoc(doc *goquery.Document) int {
	return extractSortOrderFromDocV6(doc)
}

// ExtractShowDetailOverlayFromDoc ...
func ExtractShowDetailOverlayFromDoc(doc *goquery.Document) bool {
	return extractShowDetailOverlayFromDocV6(doc)
}

// ExtractAnimatedSlidersFromDoc ...
func ExtractAnimatedSlidersFromDoc(doc *goquery.Document) bool {
	return extractAnimatedSlidersFromDocV6(doc)
}

// ExtractAnimatedOverviewFromDoc ...
func ExtractAnimatedOverviewFromDoc(doc *goquery.Document) bool {
	return extractAnimatedOverviewFromDocV6(doc)
}

// ExtractPopupsNoticesFromDoc ...
func ExtractPopupsNoticesFromDoc(doc *goquery.Document) bool {
	return extractPopupsNoticesFromDocV6(doc)
}

// ExtractPopopsCombatreportFromDoc ...
func ExtractPopopsCombatreportFromDoc(doc *goquery.Document) bool {
	return extractPopopsCombatreportFromDocV6(doc)
}

// ExtractSpioReportPicturesFromDoc ...
func ExtractSpioReportPicturesFromDoc(doc *goquery.Document) bool {
	return extractSpioReportPicturesFromDocV6(doc)
}

// ExtractMsgResultsPerPageFromDoc ...
func ExtractMsgResultsPerPageFromDoc(doc *goquery.Document) int {
	return extractMsgResultsPerPageFromDocV6(doc)
}

// ExtractAuctioneerNotificationsFromDoc ...
func ExtractAuctioneerNotificationsFromDoc(doc *goquery.Document) bool {
	return extractAuctioneerNotificationsFromDocV6(doc)
}

// ExtractEconomyNotificationsFromDoc ...
func ExtractEconomyNotificationsFromDoc(doc *goquery.Document) bool {
	return extractEconomyNotificationsFromDocV6(doc)
}

// ExtractShowActivityMinutesFromDoc ...
func ExtractShowActivityMinutesFromDoc(doc *goquery.Document) bool {
	return extractShowActivityMinutesFromDocV6(doc)
}

// ExtractPreserveSystemOnPlanetChangeFromDoc ...
func ExtractPreserveSystemOnPlanetChangeFromDoc(doc *goquery.Document) bool {
	return extractPreserveSystemOnPlanetChangeFromDocV6(doc)
}

// ExtractNotifBuildListFromDoc ...
func ExtractNotifBuildListFromDoc(doc *goquery.Document) bool {
	return extractNotifBuildListFromDocV6(doc)
}

// ExtractNotifFriendlyFleetActivitiesFromDoc ...
func ExtractNotifFriendlyFleetActivitiesFromDoc(doc *goquery.Document) bool {
	return extractNotifFriendlyFleetActivitiesFromDocV6(doc)
}

// ExtractNotifHostileFleetActivitiesFromDoc ...
func ExtractNotifHostileFleetActivitiesFromDoc(doc *goquery.Document) bool {
	return extractNotifHostileFleetActivitiesFromDocV6(doc)
}

// ExtractNotifForeignEspionageFromDoc ...
func ExtractNotifForeignEspionageFromDoc(doc *goquery.Document) bool {
	return extractNotifForeignEspionageFromDocV6(doc)
}

// ExtractNotifAllianceBroadcastsFromDoc ...
func ExtractNotifAllianceBroadcastsFromDoc(doc *goquery.Document) bool {
	return extractNotifAllianceBroadcastsFromDocV6(doc)
}

// ExtractNotifAllianceMessagesFromDoc ...
func ExtractNotifAllianceMessagesFromDoc(doc *goquery.Document) bool {
	return extractNotifAllianceMessagesFromDocV6(doc)
}

// ExtractNotifAuctionsFromDoc ...
func ExtractNotifAuctionsFromDoc(doc *goquery.Document) bool {
	return extractNotifAuctionsFromDocV6(doc)
}

// ExtractNotifAccountFromDoc ...
func ExtractNotifAccountFromDoc(doc *goquery.Document) bool {
	return extractNotifAccountFromDocV6(doc)
}

// </ Extract from doc> -------------------------------------------------------

// <Works with []byte only> ---------------------------------------------------

// ExtractPlanetCoordinate extracts planet coordinate from html page
func ExtractPlanetCoordinate(pageHTML []byte) (Coordinate, error) {
	return extractPlanetCoordinateV6(pageHTML)
}

// ExtractPlanetID extracts planet id from html page
func ExtractPlanetID(pageHTML []byte) (CelestialID, error) {
	return extractPlanetIDV6(pageHTML)
}

// ExtractOverviewShipSumCountdownFromBytes ...
func ExtractOverviewShipSumCountdownFromBytes(pageHTML []byte) int {
	return extractOverviewShipSumCountdownFromBytesV6(pageHTML)
}

// ExtractOGameTimestampFromBytes extracts ogame timestamp from an html page
func ExtractOGameTimestampFromBytes(pageHTML []byte) int64 {
	return extractOGameTimestampFromBytesV6(pageHTML)
}

// ExtractPlanetType extracts planet type from html page
func ExtractPlanetType(pageHTML []byte) (CelestialType, error) {
	return extractPlanetTypeV6(pageHTML)
}

// ExtractAjaxChatToken ...
func ExtractAjaxChatToken(pageHTML []byte) (string, error) {
	return extractAjaxChatTokenV6(pageHTML)
}

// ExtractUserInfos ...
func ExtractUserInfos(pageHTML []byte, lang string) (UserInfos, error) {
	return extractUserInfosV6(pageHTML, lang)
}

func ExtractResourcesDetails(pageHTML []byte) (out ResourcesDetails, err error) {
	return extractResourcesDetailsV6(pageHTML)
}

// </Works with []byte only> --------------------------------------------------

// ExtractCoord ...
func ExtractCoord(v string) (coord Coordinate) {
	return extractCoordV6(v)
}

// ExtractGalaxyInfos ...
func ExtractGalaxyInfos(pageHTML []byte, botPlayerName string, botPlayerID, botPlayerRank int) (SystemInfos, error) {
	return extractGalaxyInfosV6(pageHTML, botPlayerName, botPlayerID, botPlayerRank)
}

func extractPhalanx(pageHTML []byte) ([]Fleet, error) {
	return extractPhalanxV6(pageHTML)
}

// Return the available ships to send, form token, possible moon IDs and wait time (if any)
// given a jump gate popup html.
func extractJumpGate(pageHTML []byte) (ShipsInfos, string, []MoonID, int) {
	return extractJumpGateV6(pageHTML)
}

// ExtractFederation ...
func ExtractFederation(pageHTML []byte) url.Values {
	return extractFederationV6(pageHTML)
}

// ExtractConstructions ...
func ExtractConstructions(pageHTML []byte) (buildingID ID, buildingCountdown int, researchID ID, researchCountdown int) {
	return extractConstructionsV6(pageHTML)
}

// ExtractFleetDeutSaveFactor extract fleet deut save factor
func ExtractFleetDeutSaveFactor(pageHTML []byte) float64 {
	return extractFleetDeutSaveFactorV6(pageHTML)
}

func extractCancelBuildingInfos(pageHTML []byte) (token string, techID, listID int, err error) {
	return extractCancelBuildingInfosV6(pageHTML)
}

func extractCancelResearchInfos(pageHTML []byte) (token string, techID, listID int, err error) {
	return extractCancelResearchInfosV6(pageHTML)
}

// extract universe speed from html calculation
// pageHTML := b.getPageContent(url.Values{"page": {"techtree"}, "tab": {"2"}, "techID": {"1"}})
func extractUniverseSpeed(pageHTML []byte) int {
	return extractUniverseSpeedV6(pageHTML)
}
