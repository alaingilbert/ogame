package ogame

import (
	"bytes"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/alaingilbert/clockwork"
)

// ExtractorV7 ...
type ExtractorV7 struct {
	ExtractorV6
}

// NewExtractorV7 ...
func NewExtractorV7() *ExtractorV7 {
	return &ExtractorV7{}
}

// ExtractPremiumToken ...
func (e ExtractorV7) ExtractPremiumToken(pageHTML []byte, days int64) (string, error) {
	return extractPremiumTokenV7(pageHTML, days)
}

// ExtractResourcesDetailsFromFullPage ...
func (e ExtractorV7) ExtractResourcesDetailsFromFullPage(pageHTML []byte) ResourcesDetails {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractResourcesDetailsFromFullPageFromDoc(doc)
}

// ExtractResourcesDetailsFromFullPageFromDoc ...
func (e ExtractorV7) ExtractResourcesDetailsFromFullPageFromDoc(doc *goquery.Document) ResourcesDetails {
	return extractResourcesDetailsFromFullPageFromDocV7(doc)
}

// ExtractExpeditionMessages ...
func (e ExtractorV7) ExtractExpeditionMessages(pageHTML []byte, location *time.Location) ([]ExpeditionMessage, int64, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractExpeditionMessagesFromDoc(doc, location)
}

// ExtractMarketplaceMessages ...
func (e ExtractorV7) ExtractMarketplaceMessages(pageHTML []byte, location *time.Location) ([]MarketplaceMessage, int64, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractMarketplaceMessagesFromDoc(doc, location)
}

// ExtractDefense ...
func (e ExtractorV7) ExtractDefense(pageHTML []byte) (DefensesInfos, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractDefenseFromDoc(doc)
}

// ExtractFacilities ...
func (e ExtractorV7) ExtractFacilities(pageHTML []byte) (Facilities, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractFacilitiesFromDoc(doc)
}

// ExtractResearch ...
func (e ExtractorV7) ExtractResearch(pageHTML []byte) Researches {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractResearchFromDoc(doc)
}

// ExtractShips ...
func (e ExtractorV7) ExtractShips(pageHTML []byte) (ShipsInfos, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractShipsFromDoc(doc)
}

// ExtractResourceSettings ...
func (e ExtractorV7) ExtractResourceSettings(pageHTML []byte) (ResourceSettings, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractResourceSettingsFromDoc(doc)
}

// ExtractCharacterClass ...
func (e ExtractorV7) ExtractCharacterClass(pageHTML []byte) (CharacterClass, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractCharacterClassFromDoc(doc)
}

// ExtractResourcesBuildings ...
func (e ExtractorV7) ExtractResourcesBuildings(pageHTML []byte) (ResourcesBuildings, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractResourcesBuildingsFromDoc(doc)
}

// ExtractResourcesDetails ...
func (e ExtractorV7) ExtractResourcesDetails(pageHTML []byte) (out ResourcesDetails, err error) {
	return extractResourcesDetailsV7(pageHTML)
}

// ExtractConstructions ...
func (e ExtractorV7) ExtractConstructions(pageHTML []byte) (buildingID ID, buildingCountdown int64, researchID ID, researchCountdown int64) {
	return extractConstructionsV7(pageHTML, clockwork.NewRealClock())
}

// ExtractFleet1Ships ...
func (e ExtractorV7) ExtractFleet1Ships(pageHTML []byte) ShipsInfos {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractFleet1ShipsFromDoc(doc)
}

// ExtractCombatReportMessagesSummary ...
func (e ExtractorV7) ExtractCombatReportMessagesSummary(pageHTML []byte) ([]CombatReportSummary, int64) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractCombatReportMessagesFromDoc(doc)
}

// ExtractIPM ...
func (e ExtractorV7) ExtractIPM(pageHTML []byte) (duration, max int64, token string) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractIPMFromDoc(doc)
}

// ExtractIPMFromDoc ...
func (e ExtractorV7) ExtractIPMFromDoc(doc *goquery.Document) (duration, max int64, token string) {
	return extractIPMFromDocV7(doc)
}

// ExtractEspionageReport ...
func (e ExtractorV7) ExtractEspionageReport(pageHTML []byte, location *time.Location) (EspionageReport, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractEspionageReportFromDoc(doc, location)
}

// ExtractOverviewProduction extracts ships/defenses (partial) production from the overview page
func (e ExtractorV7) ExtractOverviewProduction(pageHTML []byte) ([]Quantifiable, int64, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	shipSumCountdown := e.ExtractOverviewShipSumCountdownFromBytes(pageHTML)
	production, err := e.ExtractOverviewProductionFromDoc(doc)
	return production, shipSumCountdown, err
}

// ExtractOverviewShipSumCountdownFromBytes ...
func (e ExtractorV7) ExtractOverviewShipSumCountdownFromBytes(pageHTML []byte) int64 {
	return extractOverviewShipSumCountdownFromBytesV7(pageHTML)
}

// ExtractOverviewProductionFromDoc extracts ships/defenses (partial) production from the overview page
func (e ExtractorV7) ExtractOverviewProductionFromDoc(doc *goquery.Document) ([]Quantifiable, error) {
	return extractOverviewProductionFromDocV7(doc)
}

// ExtractFleet1ShipsFromDoc ...
func (e ExtractorV7) ExtractFleet1ShipsFromDoc(doc *goquery.Document) (s ShipsInfos) {
	return extractFleet1ShipsFromDocV7(doc)
}

// ExtractResourceSettingsFromDoc ...
func (e ExtractorV7) ExtractResourceSettingsFromDoc(doc *goquery.Document) (ResourceSettings, error) {
	return extractResourceSettingsFromDocV7(doc)
}

// ExtractDefenseFromDoc ...
func (e ExtractorV7) ExtractDefenseFromDoc(doc *goquery.Document) (DefensesInfos, error) {
	return extractDefenseFromDocV7(doc)
}

// ExtractExpeditionMessagesFromDoc ...
func (e ExtractorV7) ExtractExpeditionMessagesFromDoc(doc *goquery.Document, location *time.Location) ([]ExpeditionMessage, int64, error) {
	return extractExpeditionMessagesFromDocV7(doc, location)
}

// ExtractMarketplaceMessagesFromDoc ...
func (e ExtractorV7) ExtractMarketplaceMessagesFromDoc(doc *goquery.Document, location *time.Location) ([]MarketplaceMessage, int64, error) {
	return extractMarketplaceMessagesFromDocV7(doc, location)
}

// ExtractFacilitiesFromDoc ...
func (e ExtractorV7) ExtractFacilitiesFromDoc(doc *goquery.Document) (Facilities, error) {
	return extractFacilitiesFromDocV7(doc)
}

// ExtractResearchFromDoc ...
func (e ExtractorV7) ExtractResearchFromDoc(doc *goquery.Document) Researches {
	return extractResearchFromDocV7(doc)
}

// ExtractShipsFromDoc ...
func (e ExtractorV7) ExtractShipsFromDoc(doc *goquery.Document) (ShipsInfos, error) {
	return extractShipsFromDocV7(doc)
}

// ExtractResourcesBuildingsFromDoc ...
func (e ExtractorV7) ExtractResourcesBuildingsFromDoc(doc *goquery.Document) (ResourcesBuildings, error) {
	return extractResourcesBuildingsFromDocV7(doc)
}

// ExtractCombatReportMessagesFromDoc ...
func (e ExtractorV7) ExtractCombatReportMessagesFromDoc(doc *goquery.Document) ([]CombatReportSummary, int64) {
	return extractCombatReportMessagesFromDocV7(doc)
}

// ExtractEspionageReportFromDoc ...
func (e ExtractorV7) ExtractEspionageReportFromDoc(doc *goquery.Document, location *time.Location) (EspionageReport, error) {
	return extractEspionageReportFromDocV7(doc, location)
}

// ExtractCancelBuildingInfos ...
func (e ExtractorV7) ExtractCancelBuildingInfos(pageHTML []byte) (token string, techID, listID int64, err error) {
	return extractCancelBuildingInfosV7(pageHTML)
}

// ExtractCancelResearchInfos ...
func (e ExtractorV7) ExtractCancelResearchInfos(pageHTML []byte) (token string, techID, listID int64, err error) {
	return extractCancelResearchInfosV7(pageHTML)
}

// ExtractCharacterClassFromDoc ...
func (e ExtractorV7) ExtractCharacterClassFromDoc(doc *goquery.Document) (CharacterClass, error) {
	return extractCharacterClassFromDocV7(doc)
}
