package v6

import (
	"bytes"
	"errors"
	"net/url"
	"time"

	"github.com/alaingilbert/ogame/pkg/ogame"

	"github.com/PuerkitoBio/goquery"
	"github.com/alaingilbert/clockwork"
)

// Extractor ...
type Extractor struct {
	loc             *time.Location
	lang            string
	lifeformEnabled bool
}

// NewExtractor ...
func NewExtractor() *Extractor {
	//loc := time.UTC
	//lang := "en"
	return &Extractor{}
}

func (e *Extractor) SetLocation(loc *time.Location)          { e.loc = loc }
func (e *Extractor) SetLanguage(lang string)                 { e.lang = lang }
func (e *Extractor) SetLifeformEnabled(lifeformEnabled bool) { e.lifeformEnabled = lifeformEnabled }
func (e *Extractor) GetLifeformEnabled() bool                { return e.lifeformEnabled }
func (e *Extractor) GetLocation() *time.Location {
	if e.loc == nil {
		return time.UTC
	}
	return e.loc
}
func (e *Extractor) GetLanguage() string {
	if e.lang == "" {
		return "en"
	}
	return e.lang
}

// ExtractTechnologyDetails ...
func (e *Extractor) ExtractTechnologyDetails(pageHTML []byte) (out ogame.TechnologyDetails, err error) {
	panic("implement me")
}

// ExtractTechnologyDetailsFromDoc ...
func (e *Extractor) ExtractTechnologyDetailsFromDoc(doc *goquery.Document) (ogame.TechnologyDetails, error) {
	panic("implement me")
}

func (e *Extractor) ExtractCancelLfBuildingInfos(pageHTML []byte) (token string, id, listID int64, err error) {
	panic("implement me")
}

// ExtractActiveItems ...
func (e *Extractor) ExtractActiveItems(pageHTML []byte) ([]ogame.ActiveItem, error) {
	panic("implement me")
}

// ExtractPremiumToken ...
func (e *Extractor) ExtractPremiumToken(pageHTML []byte, days int64) (string, error) {
	panic("implement me")
}

// ExtractTechs ...
func (e *Extractor) ExtractTechs(pageHTML []byte) (ogame.ResourcesBuildings, ogame.Facilities, ogame.ShipsInfos, ogame.DefensesInfos, ogame.Researches, ogame.LfBuildings, ogame.LfResearches, error) {
	panic("implement me")
}

// ExtractDestroyRockets ...
func (e *Extractor) ExtractDestroyRockets(pageHTML []byte) (abm, ipm int64, token string, err error) {
	panic("implement me")
}

// ExtractCancelFleetToken ...
func (e *Extractor) ExtractCancelFleetToken(pageHTML []byte, fleetID ogame.FleetID) (string, error) {
	panic("implement me")
}

// ExtractMarketplaceMessages ...
func (e *Extractor) ExtractMarketplaceMessages(pageHTML []byte) ([]ogame.MarketplaceMessage, int64, error) {
	panic("implement me")
}

// ExtractExpeditionMessages ...
func (e *Extractor) ExtractExpeditionMessages(pageHTML []byte) ([]ogame.ExpeditionMessage, int64, error) {
	panic("implement me")
}

// ExtractExpeditionMessagesFromDoc ...
func (e *Extractor) ExtractExpeditionMessagesFromDoc(doc *goquery.Document) ([]ogame.ExpeditionMessage, int64, error) {
	panic("implement me")
}

// ExtractTearDownButtonEnabled ...
func (e *Extractor) ExtractTearDownButtonEnabled(pageHTML []byte) bool {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractTearDownButtonEnabledFromDoc(doc)
}

// ExtractUpgradeToken ...
func (e *Extractor) ExtractUpgradeToken(pageHTML []byte) (string, error) {
	return extractUpgradeToken(pageHTML)
}

// ExtractLifeformEnabled ...
func (e *Extractor) ExtractLifeformEnabled(pageHTML []byte) bool {
	return bytes.Contains(pageHTML, []byte(`lifeformEnabled":true`))
}

// ExtractIsInVacation ...
func (e *Extractor) ExtractIsInVacation(pageHTML []byte) bool {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractIsInVacationFromDoc(doc)
}

// ExtractPlanets ...
func (e *Extractor) ExtractPlanets(pageHTML []byte) []ogame.Planet {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractPlanetsFromDoc(doc)
}

// ExtractPlanet ...
func (e *Extractor) ExtractPlanet(pageHTML []byte, v any) (ogame.Planet, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractPlanetFromDoc(doc, v)
}

// ExtractMoons ...
func (e *Extractor) ExtractMoons(pageHTML []byte) []ogame.Moon {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractMoonsFromDoc(doc)
}

// ExtractMoon ...
func (e *Extractor) ExtractMoon(pageHTML []byte, v any) (ogame.Moon, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractMoonFromDoc(doc, v)
}

// ExtractCelestials ...
func (e *Extractor) ExtractCelestials(pageHTML []byte) ([]ogame.Celestial, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractCelestialsFromDoc(doc)
}

// ExtractCelestial ...
func (e *Extractor) ExtractCelestial(pageHTML []byte, v any) (ogame.Celestial, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractCelestialFromDoc(doc, v)
}

// ExtractServerTime ...
func (e *Extractor) ExtractServerTime(pageHTML []byte) (time.Time, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractServerTimeFromDoc(doc)
}

// ExtractFleetsFromEventList ...
func (e *Extractor) ExtractFleetsFromEventList(pageHTML []byte) []ogame.Fleet {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractFleetsFromEventListFromDoc(doc)
}

// ExtractIPM ...
func (e *Extractor) ExtractIPM(pageHTML []byte) (duration, max int64, token string) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractIPMFromDoc(doc)
}

// ExtractFleets ...
func (e *Extractor) ExtractFleets(pageHTML []byte) (res []ogame.Fleet) {
	return e.extractFleets(pageHTML, e.loc)
}

func (e *Extractor) extractFleets(pageHTML []byte, location *time.Location) (res []ogame.Fleet) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.extractFleetsFromDoc(doc, location)
}

// ExtractSlots ...
func (e *Extractor) ExtractSlots(pageHTML []byte) (ogame.Slots, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractSlotsFromDoc(doc)
}

// ExtractOgameTimestamp ...
func (e *Extractor) ExtractOgameTimestamp(pageHTML []byte) int64 {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractOgameTimestampFromDoc(doc)
}

// ExtractResources ...
func (e *Extractor) ExtractResources(pageHTML []byte) ogame.Resources {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractResourcesFromDoc(doc)
}

// ExtractResourcesDetailsFromFullPage ...
func (e *Extractor) ExtractResourcesDetailsFromFullPage(pageHTML []byte) ogame.ResourcesDetails {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractResourcesDetailsFromFullPageFromDoc(doc)
}

// ExtractResourceSettings ...
func (e *Extractor) ExtractResourceSettings(pageHTML []byte) (ogame.ResourceSettings, string, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractResourceSettingsFromDoc(doc)
}

// ExtractAttacks ...
func (e *Extractor) ExtractAttacks(pageHTML []byte, ownCoords []ogame.Coordinate) ([]ogame.AttackEvent, error) {
	return e.extractAttacks(pageHTML, clockwork.NewRealClock(), ownCoords)
}

func (e *Extractor) extractAttacks(pageHTML []byte, clock clockwork.Clock, ownCoords []ogame.Coordinate) ([]ogame.AttackEvent, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.extractAttacksFromDoc(doc, clock, ownCoords)
}

// ExtractOfferOfTheDay ...
func (e *Extractor) ExtractOfferOfTheDay(pageHTML []byte) (int64, string, ogame.PlanetResources, ogame.Multiplier, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractOfferOfTheDayFromDoc(doc)
}

// ExtractResourcesBuildings ...
func (e *Extractor) ExtractResourcesBuildings(pageHTML []byte) (ogame.ResourcesBuildings, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractResourcesBuildingsFromDoc(doc)
}

// ExtractDefense ...
func (e *Extractor) ExtractDefense(pageHTML []byte) (ogame.DefensesInfos, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractDefenseFromDoc(doc)
}

// ExtractShips ...
func (e *Extractor) ExtractShips(pageHTML []byte) (ogame.ShipsInfos, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractShipsFromDoc(doc)
}

// ExtractFacilities ...
func (e *Extractor) ExtractFacilities(pageHTML []byte) (ogame.Facilities, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractFacilitiesFromDoc(doc)
}

// ExtractTearDownToken ...
func (e *Extractor) ExtractTearDownToken(pageHTML []byte) (string, error) {
	return extractTearDownToken(pageHTML)
}

// ExtractResearch ...
func (e *Extractor) ExtractResearch(pageHTML []byte) ogame.Researches {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractResearchFromDoc(doc)
}

// ExtractProduction extracts ships/defenses production from the shipyard page
func (e *Extractor) ExtractProduction(pageHTML []byte) ([]ogame.Quantifiable, int64, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	shipSumCountdown := e.ExtractOverviewShipSumCountdownFromBytes(pageHTML)
	production, err := e.ExtractProductionFromDoc(doc)
	return production, shipSumCountdown, err
}

// ExtractOverviewProduction extracts ships/defenses (partial) production from the overview page
func (e *Extractor) ExtractOverviewProduction(pageHTML []byte) ([]ogame.Quantifiable, int64, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	shipSumCountdown := e.ExtractOverviewShipSumCountdownFromBytes(pageHTML)
	production, err := e.ExtractOverviewProductionFromDoc(doc)
	return production, shipSumCountdown, err
}

// ExtractFleet1Ships ...
func (e *Extractor) ExtractFleet1Ships(pageHTML []byte) ogame.ShipsInfos {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractFleet1ShipsFromDoc(doc)
}

// ExtractEspionageReportMessageIDs ...
func (e *Extractor) ExtractEspionageReportMessageIDs(pageHTML []byte) ([]ogame.EspionageReportSummary, int64, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractEspionageReportMessageIDsFromDoc(doc)
}

// ExtractCombatReportMessagesSummary ...
func (e *Extractor) ExtractCombatReportMessagesSummary(pageHTML []byte) ([]ogame.CombatReportSummary, int64, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractCombatReportMessagesFromDoc(doc)
}

// ExtractEspionageReport ...
func (e *Extractor) ExtractEspionageReport(pageHTML []byte) (ogame.EspionageReport, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractEspionageReportFromDoc(doc)
}

// ExtractResourcesProductions ...
func (e *Extractor) ExtractResourcesProductions(pageHTML []byte) (ogame.Resources, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractResourcesProductionsFromDoc(doc)
}

// ExtractPreferences ...
func (e *Extractor) ExtractPreferences(pageHTML []byte) ogame.Preferences {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractPreferencesFromDoc(doc)
}

// ExtractSpioAnz ...
func (e *Extractor) ExtractSpioAnz(pageHTML []byte) int64 {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractSpioAnzFromDoc(doc)
}

// ExtractPreferencesShowActivityMinutes ...
func (e *Extractor) ExtractPreferencesShowActivityMinutes(pageHTML []byte) bool {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractShowActivityMinutesFromDoc(doc)
}

// ExtractHiddenFields utils function to extract hidden input from a page
func (e *Extractor) ExtractHiddenFields(pageHTML []byte) (fields url.Values) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractHiddenFieldsFromDoc(doc)
}

// ExtractCommander ...
func (e *Extractor) ExtractCommander(pageHTML []byte) bool {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractCommanderFromDoc(doc)
}

// ExtractAdmiral ...
func (e *Extractor) ExtractAdmiral(pageHTML []byte) bool {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractAdmiralFromDoc(doc)
}

// ExtractEngineer ...
func (e *Extractor) ExtractEngineer(pageHTML []byte) bool {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractEngineerFromDoc(doc)
}

// ExtractGeologist ...
func (e *Extractor) ExtractGeologist(pageHTML []byte) bool {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractGeologistFromDoc(doc)
}

// ExtractTechnocrat ...
func (e *Extractor) ExtractTechnocrat(pageHTML []byte) bool {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractTechnocratFromDoc(doc)
}

// ExtractOGameSession ...
func (e *Extractor) ExtractOGameSession(pageHTML []byte) string {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractOGameSessionFromDoc(doc)
}

// <Extract from doc> ---------------------------------------------------------

// ExtractPlanetTypeFromDoc extracts planet type from doc
func (e *Extractor) ExtractPlanetTypeFromDoc(doc *goquery.Document) (ogame.CelestialType, error) {
	return extractPlanetTypeFromDoc(doc)
}

// ExtractPlanetIDFromDoc extracts planet id from doc
func (e *Extractor) ExtractPlanetIDFromDoc(doc *goquery.Document) (ogame.CelestialID, error) {
	return extractPlanetIDFromDoc(doc)
}

// ExtractHiddenFieldsFromDoc utils function to extract hidden input from a page
func (e *Extractor) ExtractHiddenFieldsFromDoc(doc *goquery.Document) url.Values {
	return extractHiddenFieldsFromDoc(doc)
}

// ExtractBodyIDFromDoc ...
func (e *Extractor) ExtractBodyIDFromDoc(doc *goquery.Document) string {
	return ExtractBodyIDFromDoc(doc)
}

// ExtractIsInVacationFromDoc ...
func (e *Extractor) ExtractIsInVacationFromDoc(doc *goquery.Document) bool {
	return extractIsInVacationFromDoc(doc)
}

// ExtractTearDownButtonEnabledFromDoc ...
func (e *Extractor) ExtractTearDownButtonEnabledFromDoc(doc *goquery.Document) bool {
	return extractTearDownButtonEnabledFromDoc(doc)
}

// ExtractPlanetsFromDoc ...
func (e *Extractor) ExtractPlanetsFromDoc(doc *goquery.Document) []ogame.Planet {
	return extractPlanetsFromDoc(doc)
}

// ExtractOgameTimestampFromDoc ...
func (e *Extractor) ExtractOgameTimestampFromDoc(doc *goquery.Document) int64 {
	return extractOgameTimestampFromDoc(doc)
}

// ExtractResourcesFromDoc ...
func (e *Extractor) ExtractResourcesFromDoc(doc *goquery.Document) ogame.Resources {
	return extractResourcesFromDoc(doc)
}

// ExtractResourcesDetailsFromFullPageFromDoc ...
func (e *Extractor) ExtractResourcesDetailsFromFullPageFromDoc(doc *goquery.Document) ogame.ResourcesDetails {
	return extractResourcesDetailsFromFullPageFromDoc(doc)
}

// ExtractPlanetFromDoc ...
func (e *Extractor) ExtractPlanetFromDoc(doc *goquery.Document, v any) (ogame.Planet, error) {
	return extractPlanetFromDoc(doc, v)
}

// ExtractMoonsFromDoc ...
func (e *Extractor) ExtractMoonsFromDoc(doc *goquery.Document) []ogame.Moon {
	return extractMoonsFromDoc(doc)
}

// ExtractMoonFromDoc ...
func (e *Extractor) ExtractMoonFromDoc(doc *goquery.Document, v any) (ogame.Moon, error) {
	return extractMoonFromDoc(doc, v)
}

// ExtractCelestialsFromDoc ...
func (e *Extractor) ExtractCelestialsFromDoc(doc *goquery.Document) ([]ogame.Celestial, error) {
	return extractCelestialsFromDoc(doc), nil
}

// ExtractCelestialFromDoc ...
func (e *Extractor) ExtractCelestialFromDoc(doc *goquery.Document, v any) (ogame.Celestial, error) {
	return extractCelestialFromDoc(doc, v)
}

// ExtractResourcesBuildingsFromDoc ...
func (e *Extractor) ExtractResourcesBuildingsFromDoc(doc *goquery.Document) (ogame.ResourcesBuildings, error) {
	return extractResourcesBuildingsFromDoc(doc)
}

// ExtractDefenseFromDoc ...
func (e *Extractor) ExtractDefenseFromDoc(doc *goquery.Document) (ogame.DefensesInfos, error) {
	return extractDefenseFromDoc(doc)
}

// ExtractShipsFromDoc ...
func (e *Extractor) ExtractShipsFromDoc(doc *goquery.Document) (ogame.ShipsInfos, error) {
	return extractShipsFromDoc(doc)
}

// ExtractFacilitiesFromDoc ...
func (e *Extractor) ExtractFacilitiesFromDoc(doc *goquery.Document) (ogame.Facilities, error) {
	return extractFacilitiesFromDoc(doc)
}

// ExtractResearchFromDoc ...
func (e *Extractor) ExtractResearchFromDoc(doc *goquery.Document) ogame.Researches {
	return extractResearchFromDoc(doc)
}

// ExtractOGameSessionFromDoc ...
func (e *Extractor) ExtractOGameSessionFromDoc(doc *goquery.Document) string {
	return ExtractOGameSessionFromDoc(doc)
}

// ExtractAttacksFromDoc ...
func (e *Extractor) ExtractAttacksFromDoc(doc *goquery.Document, ownCoords []ogame.Coordinate) ([]ogame.AttackEvent, error) {
	return e.extractAttacksFromDoc(doc, clockwork.NewRealClock(), ownCoords)
}

func (e *Extractor) extractAttacksFromDoc(doc *goquery.Document, clock clockwork.Clock, ownCoords []ogame.Coordinate) ([]ogame.AttackEvent, error) {
	return extractAttacksFromDoc(doc, clock, ownCoords)
}

// ExtractOfferOfTheDayFromDoc ...
func (e *Extractor) ExtractOfferOfTheDayFromDoc(doc *goquery.Document) (price int64, importToken string, planetResources ogame.PlanetResources, multiplier ogame.Multiplier, err error) {
	return extractOfferOfTheDayFromDoc(doc)
}

// ExtractProductionFromDoc extracts ships/defenses production from the shipyard page
func (e *Extractor) ExtractProductionFromDoc(doc *goquery.Document) ([]ogame.Quantifiable, error) {
	return extractProductionFromDoc(doc)
}

// ExtractOverviewProductionFromDoc extracts ships/defenses (partial) production from the overview page
func (e *Extractor) ExtractOverviewProductionFromDoc(doc *goquery.Document) ([]ogame.Quantifiable, error) {
	return extractOverviewProductionFromDoc(doc)
}

// ExtractFleet1ShipsFromDoc ...
func (e *Extractor) ExtractFleet1ShipsFromDoc(doc *goquery.Document) (s ogame.ShipsInfos) {
	return extractFleet1ShipsFromDoc(doc)
}

// ExtractFleetDispatchACSFromDoc ...
func (e *Extractor) ExtractFleetDispatchACSFromDoc(doc *goquery.Document) []ogame.ACSValues {
	return extractFleetDispatchACSFromDoc(doc)
}

// ExtractEspionageReportMessageIDsFromDoc ...
func (e *Extractor) ExtractEspionageReportMessageIDsFromDoc(doc *goquery.Document) ([]ogame.EspionageReportSummary, int64, error) {
	return extractEspionageReportMessageIDsFromDoc(doc)
}

// ExtractCombatReportMessagesFromDoc ...
func (e *Extractor) ExtractCombatReportMessagesFromDoc(doc *goquery.Document) ([]ogame.CombatReportSummary, int64, error) {
	return extractCombatReportMessagesFromDoc(doc)
}

// ExtractEspionageReportFromDoc ...
func (e *Extractor) ExtractEspionageReportFromDoc(doc *goquery.Document) (ogame.EspionageReport, error) {
	return extractEspionageReportFromDoc(doc, e.loc)
}

// ExtractResourcesProductionsFromDoc ...
func (e *Extractor) ExtractResourcesProductionsFromDoc(doc *goquery.Document) (ogame.Resources, error) {
	return extractResourcesProductionsFromDoc(doc)
}

// ExtractPreferencesFromDoc ...
func (e *Extractor) ExtractPreferencesFromDoc(doc *goquery.Document) ogame.Preferences {
	return extractPreferencesFromDoc(doc)
}

// ExtractResourceSettingsFromDoc ...
func (e *Extractor) ExtractResourceSettingsFromDoc(doc *goquery.Document) (ogame.ResourceSettings, string, error) {
	return extractResourceSettingsFromDoc(doc)
}

// ExtractFleetsFromEventListFromDoc ...
func (e *Extractor) ExtractFleetsFromEventListFromDoc(doc *goquery.Document) []ogame.Fleet {
	return extractFleetsFromEventListFromDoc(doc)
}

// ExtractIPMFromDoc ...
func (e *Extractor) ExtractIPMFromDoc(doc *goquery.Document) (duration, max int64, token string) {
	return extractIPMFromDoc(doc)
}

// ExtractFleetsFromDoc ...
func (e *Extractor) ExtractFleetsFromDoc(doc *goquery.Document) (res []ogame.Fleet) {
	return e.extractFleetsFromDoc(doc, e.loc)
}

func (e *Extractor) extractFleetsFromDoc(doc *goquery.Document, location *time.Location) (res []ogame.Fleet) {
	return extractFleetsFromDoc(doc, location, e.lifeformEnabled)
}

// ExtractSlotsFromDoc extract fleet slots from page "fleet1"
// page "movement" redirect to "fleet1" when there is no fleet
func (e *Extractor) ExtractSlotsFromDoc(doc *goquery.Document) (ogame.Slots, error) {
	return extractSlotsFromDoc(doc)
}

// ExtractServerTimeFromDoc ...
func (e *Extractor) ExtractServerTimeFromDoc(doc *goquery.Document) (time.Time, error) {
	return extractServerTimeFromDoc(doc)
}

// ExtractSpioAnzFromDoc ...
func (e *Extractor) ExtractSpioAnzFromDoc(doc *goquery.Document) int64 {
	return extractSpioAnzFromDoc(doc)
}

// ExtractDisableChatBarFromDoc ...
func (e *Extractor) ExtractDisableChatBarFromDoc(doc *goquery.Document) bool {
	return extractDisableChatBarFromDoc(doc)
}

// ExtractDisableOutlawWarningFromDoc ...
func (e *Extractor) ExtractDisableOutlawWarningFromDoc(doc *goquery.Document) bool {
	return extractDisableOutlawWarningFromDoc(doc)
}

// ExtractMobileVersionFromDoc ...
func (e *Extractor) ExtractMobileVersionFromDoc(doc *goquery.Document) bool {
	return extractMobileVersionFromDoc(doc)
}

// ExtractShowOldDropDownsFromDoc ...
func (e *Extractor) ExtractShowOldDropDownsFromDoc(doc *goquery.Document) bool {
	return extractShowOldDropDownsFromDoc(doc)
}

// ExtractActivateAutofocusFromDoc ...
func (e *Extractor) ExtractActivateAutofocusFromDoc(doc *goquery.Document) bool {
	return extractActivateAutofocusFromDoc(doc)
}

// ExtractEventsShowFromDoc ...
func (e *Extractor) ExtractEventsShowFromDoc(doc *goquery.Document) int64 {
	return extractEventsShowFromDoc(doc)
}

// ExtractSortSettingFromDoc ...
func (e *Extractor) ExtractSortSettingFromDoc(doc *goquery.Document) int64 {
	return extractSortSettingFromDoc(doc)
}

// ExtractSortOrderFromDoc ...
func (e *Extractor) ExtractSortOrderFromDoc(doc *goquery.Document) int64 {
	return extractSortOrderFromDoc(doc)
}

// ExtractShowDetailOverlayFromDoc ...
func (e *Extractor) ExtractShowDetailOverlayFromDoc(doc *goquery.Document) bool {
	return extractShowDetailOverlayFromDoc(doc)
}

// ExtractAnimatedSlidersFromDoc ...
func (e *Extractor) ExtractAnimatedSlidersFromDoc(doc *goquery.Document) bool {
	return extractAnimatedSlidersFromDoc(doc)
}

// ExtractAnimatedOverviewFromDoc ...
func (e *Extractor) ExtractAnimatedOverviewFromDoc(doc *goquery.Document) bool {
	return extractAnimatedOverviewFromDoc(doc)
}

// ExtractPopupsNoticesFromDoc ...
func (e *Extractor) ExtractPopupsNoticesFromDoc(doc *goquery.Document) bool {
	return extractPopupsNoticesFromDoc(doc)
}

// ExtractPopopsCombatreportFromDoc ...
func (e *Extractor) ExtractPopopsCombatreportFromDoc(doc *goquery.Document) bool {
	return extractPopopsCombatreportFromDoc(doc)
}

// ExtractSpioReportPicturesFromDoc ...
func (e *Extractor) ExtractSpioReportPicturesFromDoc(doc *goquery.Document) bool {
	return extractSpioReportPicturesFromDoc(doc)
}

// ExtractMsgResultsPerPageFromDoc ...
func (e *Extractor) ExtractMsgResultsPerPageFromDoc(doc *goquery.Document) int64 {
	return extractMsgResultsPerPageFromDoc(doc)
}

// ExtractAuctioneerNotificationsFromDoc ...
func (e *Extractor) ExtractAuctioneerNotificationsFromDoc(doc *goquery.Document) bool {
	return extractAuctioneerNotificationsFromDoc(doc)
}

// ExtractEconomyNotificationsFromDoc ...
func (e *Extractor) ExtractEconomyNotificationsFromDoc(doc *goquery.Document) bool {
	return extractEconomyNotificationsFromDoc(doc)
}

// ExtractShowActivityMinutesFromDoc ...
func (e *Extractor) ExtractShowActivityMinutesFromDoc(doc *goquery.Document) bool {
	return extractShowActivityMinutesFromDoc(doc)
}

// ExtractPreserveSystemOnPlanetChangeFromDoc ...
func (e *Extractor) ExtractPreserveSystemOnPlanetChangeFromDoc(doc *goquery.Document) bool {
	return extractPreserveSystemOnPlanetChangeFromDoc(doc)
}

// ExtractNotifBuildListFromDoc ...
func (e *Extractor) ExtractNotifBuildListFromDoc(doc *goquery.Document) bool {
	return extractNotifBuildListFromDoc(doc)
}

// ExtractNotifFriendlyFleetActivitiesFromDoc ...
func (e *Extractor) ExtractNotifFriendlyFleetActivitiesFromDoc(doc *goquery.Document) bool {
	return extractNotifFriendlyFleetActivitiesFromDoc(doc)
}

// ExtractNotifHostileFleetActivitiesFromDoc ...
func (e *Extractor) ExtractNotifHostileFleetActivitiesFromDoc(doc *goquery.Document) bool {
	return extractNotifHostileFleetActivitiesFromDoc(doc)
}

// ExtractNotifForeignEspionageFromDoc ...
func (e *Extractor) ExtractNotifForeignEspionageFromDoc(doc *goquery.Document) bool {
	return extractNotifForeignEspionageFromDoc(doc)
}

// ExtractNotifAllianceBroadcastsFromDoc ...
func (e *Extractor) ExtractNotifAllianceBroadcastsFromDoc(doc *goquery.Document) bool {
	return extractNotifAllianceBroadcastsFromDoc(doc)
}

// ExtractNotifAllianceMessagesFromDoc ...
func (e *Extractor) ExtractNotifAllianceMessagesFromDoc(doc *goquery.Document) bool {
	return extractNotifAllianceMessagesFromDoc(doc)
}

// ExtractNotifAuctionsFromDoc ...
func (e *Extractor) ExtractNotifAuctionsFromDoc(doc *goquery.Document) bool {
	return extractNotifAuctionsFromDoc(doc)
}

// ExtractNotifAccountFromDoc ...
func (e *Extractor) ExtractNotifAccountFromDoc(doc *goquery.Document) bool {
	return extractNotifAccountFromDoc(doc)
}

// ExtractCharacterClassFromDoc ...
func (e *Extractor) ExtractCharacterClassFromDoc(doc *goquery.Document) (ogame.CharacterClass, error) {
	return 0, errors.New("character class not supported in ")
}

// ExtractCommanderFromDoc ...
func (e *Extractor) ExtractCommanderFromDoc(doc *goquery.Document) bool {
	return extractCommanderFromDoc(doc)
}

// ExtractAdmiralFromDoc ...
func (e *Extractor) ExtractAdmiralFromDoc(doc *goquery.Document) bool {
	return extractAdmiralFromDoc(doc)
}

// ExtractLifeformTypeFromDoc ...
func (e Extractor) ExtractLifeformTypeFromDoc(doc *goquery.Document) ogame.LifeformType {
	return ogame.NoneLfType
}

// ExtractEngineerFromDoc ...
func (e *Extractor) ExtractEngineerFromDoc(doc *goquery.Document) bool {
	return extractEngineerFromDoc(doc)
}

// ExtractGeologistFromDoc ...
func (e *Extractor) ExtractGeologistFromDoc(doc *goquery.Document) bool {
	return extractGeologistFromDoc(doc)
}

// ExtractTechnocratFromDoc ...
func (e *Extractor) ExtractTechnocratFromDoc(doc *goquery.Document) bool {
	return extractTechnocratFromDoc(doc)
}

// ExtractAbandonInformation ...
func (e *Extractor) ExtractAbandonInformation(doc *goquery.Document) (string, string) {
	return extractAbandonInformation(doc)
}

// </ Extract from doc> -------------------------------------------------------

// <Works with []byte only> ---------------------------------------------------

// ExtractPlanetCoordinate extracts planet coordinate from html page
func (e *Extractor) ExtractPlanetCoordinate(pageHTML []byte) (ogame.Coordinate, error) {
	return extractPlanetCoordinate(pageHTML)
}

// ExtractPlanetID extracts planet id from html page
func (e *Extractor) ExtractPlanetID(pageHTML []byte) (ogame.CelestialID, error) {
	return extractPlanetID(pageHTML)
}

// ExtractOverviewShipSumCountdownFromBytes ...
func (e *Extractor) ExtractOverviewShipSumCountdownFromBytes(pageHTML []byte) int64 {
	return extractOverviewShipSumCountdownFromBytes(pageHTML)
}

// ExtractOGameTimestampFromBytes extracts ogame timestamp from an html page
func (e *Extractor) ExtractOGameTimestampFromBytes(pageHTML []byte) int64 {
	return extractOGameTimestampFromBytes(pageHTML)
}

// ExtractPlanetType extracts planet type from html page
func (e *Extractor) ExtractPlanetType(pageHTML []byte) (ogame.CelestialType, error) {
	return extractPlanetType(pageHTML)
}

// ExtractAjaxChatToken ...
func (e *Extractor) ExtractAjaxChatToken(pageHTML []byte) (string, error) {
	return extractAjaxChatToken(pageHTML)
}

// ExtractUserInfos ...
func (e *Extractor) ExtractUserInfos(pageHTML []byte) (ogame.UserInfos, error) {
	return extractUserInfos(pageHTML, e.GetLanguage())
}

// ExtractResourcesDetails ...
func (e *Extractor) ExtractResourcesDetails(pageHTML []byte) (out ogame.ResourcesDetails, err error) {
	return extractResourcesDetails(pageHTML)
}

// </Works with []byte only> --------------------------------------------------

// ExtractCoord ...
func (e *Extractor) ExtractCoord(v string) (coord ogame.Coordinate) {
	return ExtractCoord(v)
}

// ExtractGalaxyInfos ...
func (e *Extractor) ExtractGalaxyInfos(pageHTML []byte, botPlayerName string, botPlayerID, botPlayerRank int64) (ogame.SystemInfos, error) {
	return extractGalaxyInfos(pageHTML, botPlayerName, botPlayerID, botPlayerRank)
}

// ExtractPhalanx ...
func (e *Extractor) ExtractPhalanx(pageHTML []byte) ([]ogame.Fleet, error) {
	return extractPhalanx(pageHTML)
}

// ExtractJumpGate return the available ships to send, form token, possible moon IDs and wait time (if any)
// given a jump gate popup html.
func (e *Extractor) ExtractJumpGate(pageHTML []byte) (ogame.ShipsInfos, string, []ogame.MoonID, int64) {
	return extractJumpGate(pageHTML)
}

// ExtractFederation ...
func (e *Extractor) ExtractFederation(pageHTML []byte) url.Values {
	return extractFederation(pageHTML)
}

// ExtractConstructions ...
func (e *Extractor) ExtractConstructions(pageHTML []byte) (buildingID ogame.ID, buildingCountdown int64, researchID ogame.ID, researchCountdown int64, lfBuildingID ogame.ID, lfBuildingCountdown int64, lfResearchID ogame.ID, lfResearchCountdown int64) {
	return extractConstructions(pageHTML)
}

// ExtractFleetDeutSaveFactor extract fleet deut save factor
func (e *Extractor) ExtractFleetDeutSaveFactor(pageHTML []byte) float64 {
	return extractFleetDeutSaveFactor(pageHTML)
}

// ExtractCancelBuildingInfos ...
func (e *Extractor) ExtractCancelBuildingInfos(pageHTML []byte) (token string, techID, listID int64, err error) {
	return extractCancelBuildingInfos(pageHTML)
}

// ExtractCancelResearchInfos ...
func (e *Extractor) ExtractCancelResearchInfos(pageHTML []byte) (token string, techID, listID int64, err error) {
	return extractCancelResearchInfos(pageHTML)
}

// ExtractEmpire ...
func (e *Extractor) ExtractEmpire(pageHTML []byte) ([]ogame.EmpireCelestial, error) {
	return extractEmpire(pageHTML)
}

// ExtractEmpireJSON ...
func (e *Extractor) ExtractEmpireJSON(pageHTML []byte) (any, error) {
	return ExtractEmpireJSON(pageHTML)
}

// ExtractCharacterClass ...
func (e *Extractor) ExtractCharacterClass(pageHTML []byte) (ogame.CharacterClass, error) {
	return 0, errors.New("character class not supported in ")
}

// ExtractAuction ...
func (e *Extractor) ExtractAuction(pageHTML []byte) (ogame.Auction, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return extractAuctionFromDoc(doc)
}

// ExtractHighscore ...
func (e *Extractor) ExtractHighscore(pageHTML []byte) (ogame.Highscore, error) {
	panic("not implemented")
}

// ExtractHighscoreFromDoc ...
func (e *Extractor) ExtractHighscoreFromDoc(doc *goquery.Document) (ogame.Highscore, error) {
	panic("not implemented")
}

// ExtractAllResources ...
func (e *Extractor) ExtractAllResources(pageHTML []byte) (map[ogame.CelestialID]ogame.Resources, error) {
	panic("not implemented")
}

// ExtractDMCosts ...
func (e *Extractor) ExtractDMCosts(pageHTML []byte) (ogame.DMCosts, error) {
	panic("not implemented")
}

// ExtractBuffActivation ...
func (e *Extractor) ExtractBuffActivation(pageHTML []byte) (string, []ogame.Item, error) {
	panic("not implemented")
}

// ExtractIsMobile ...
func (e *Extractor) ExtractIsMobile(pageHTML []byte) bool {
	panic("not implemented")
}

// ExtractIsMobileFromDoc ...
func (e *Extractor) ExtractIsMobileFromDoc(doc *goquery.Document) bool {
	panic("not implemented")
}

// ExtractLfBuildings ...
func (e *Extractor) ExtractLfBuildings(pageHTML []byte) (ogame.LfBuildings, error) {
	panic("not implemented")
}

// ExtractLfBuildingsFromDoc ...
func (e *Extractor) ExtractLfBuildingsFromDoc(doc *goquery.Document) (ogame.LfBuildings, error) {
	panic("not implemented")
}

// ExtractLfResearch ...
func (e *Extractor) ExtractLfResearch(pageHTML []byte) (ogame.LfResearches, error) {
	panic("not implemented")
}

// ExtractLfResearchFromDoc ...
func (e *Extractor) ExtractLfResearchFromDoc(doc *goquery.Document) (ogame.LfResearches, error) {
	panic("not implemented")
}

// ExtractAvailableDiscoveries
func (e *Extractor) ExtractAvailableDiscoveries(pageHTML []byte) int64 {
	panic("not implemented")
}
