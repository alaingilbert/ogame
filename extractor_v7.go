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

// ExtractResourcesBuildings ...
func (e ExtractorV7) ExtractResourcesBuildings(pageHTML []byte) (ResourcesBuildings, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractResourcesBuildingsFromDoc(doc)
}

func (e ExtractorV7) ExtractResourcesDetails(pageHTML []byte) (out ResourcesDetails, err error) {
	return extractResourcesDetailsV7(pageHTML)
}

// ExtractConstructions ...
func (e ExtractorV7) ExtractConstructions(pageHTML []byte) (buildingID ID, buildingCountdown int, researchID ID, researchCountdown int) {
	return extractConstructionsV7(pageHTML, clockwork.NewRealClock())
}

// ExtractFleet1Ships ...
func (e ExtractorV7) ExtractFleet1Ships(pageHTML []byte) ShipsInfos {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractFleet1ShipsFromDoc(doc)
}

// ExtractCombatReportMessagesSummary ...
func (e ExtractorV7) ExtractCombatReportMessagesSummary(pageHTML []byte) ([]CombatReportSummary, int) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractCombatReportMessagesFromDoc(doc)
}

// ExtractEspionageReport ...
func (e ExtractorV7) ExtractEspionageReport(pageHTML []byte, location *time.Location) (EspionageReport, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractEspionageReportFromDoc(doc, location)
}

// ExtractFleet1ShipsFromDoc ...
func (e ExtractorV7) ExtractFleet1ShipsFromDoc(doc *goquery.Document) (s ShipsInfos) {
	return extractFleet1ShipsFromDocV7(doc)
}

// ExtractDefenseFromDoc ...
func (e ExtractorV7) ExtractDefenseFromDoc(doc *goquery.Document) (DefensesInfos, error) {
	return extractDefenseFromDocV7(doc)
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
func (e ExtractorV7) ExtractCombatReportMessagesFromDoc(doc *goquery.Document) ([]CombatReportSummary, int) {
	return extractCombatReportMessagesFromDocV7(doc)
}

// ExtractEspionageReportFromDoc ...
func (e ExtractorV7) ExtractEspionageReportFromDoc(doc *goquery.Document, location *time.Location) (EspionageReport, error) {
	return extractEspionageReportFromDocV7(doc, location)
}

// ExtractCancelBuildingInfos ...
func (e ExtractorV7) ExtractCancelBuildingInfos(pageHTML []byte) (token string, techID, listID int, err error) {
	return extractCancelBuildingInfosV7(pageHTML)
}

// ExtractCancelResearchInfos ...
func (e ExtractorV7) ExtractCancelResearchInfos(pageHTML []byte) (token string, techID, listID int, err error) {
	return extractCancelResearchInfosV7(pageHTML)
}
