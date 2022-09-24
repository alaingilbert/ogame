package v9

import (
	"bytes"

	"github.com/PuerkitoBio/goquery"
	"github.com/alaingilbert/clockwork"
	v874 "github.com/alaingilbert/ogame/pkg/extractor/v874"
	"github.com/alaingilbert/ogame/pkg/ogame"
)

// Extractor ...
type Extractor struct {
	v874.Extractor
}

// NewExtractor ...
func NewExtractor() *Extractor {
	return &Extractor{}
}

// ExtractCancelLfBuildingInfos ...
func (e *Extractor) ExtractCancelLfBuildingInfos(pageHTML []byte) (token string, id, listID int64, err error) {
	return extractCancelLfBuildingInfos(pageHTML)
}

// ExtractCancelResearchInfos ...
func (e *Extractor) ExtractCancelResearchInfos(pageHTML []byte) (token string, id, listID int64, err error) {
	return extractCancelResearchInfos(pageHTML, e.GetLifeformEnabled())
}

// ExtractEmpire ...
func (e *Extractor) ExtractEmpire(pageHTML []byte) ([]ogame.EmpireCelestial, error) {
	return extractEmpire(pageHTML)
}

// ExtractOverviewProduction extracts ships/defenses (partial) production from the overview page
func (e *Extractor) ExtractOverviewProduction(pageHTML []byte) ([]ogame.Quantifiable, int64, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	shipSumCountdown := e.ExtractOverviewShipSumCountdownFromBytes(pageHTML)
	production, err := e.ExtractOverviewProductionFromDoc(doc)
	return production, shipSumCountdown, err
}

// ExtractOverviewProductionFromDoc extracts ships/defenses (partial) production from the overview page
func (e *Extractor) ExtractOverviewProductionFromDoc(doc *goquery.Document) ([]ogame.Quantifiable, error) {
	return extractOverviewProductionFromDoc(doc)
}

// ExtractEspionageReport ...
func (e *Extractor) ExtractEspionageReport(pageHTML []byte) (ogame.EspionageReport, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractEspionageReportFromDoc(doc)
}

// ExtractEspionageReportFromDoc ...
func (e *Extractor) ExtractEspionageReportFromDoc(doc *goquery.Document) (ogame.EspionageReport, error) {
	return extractEspionageReportFromDoc(doc, e.GetLocation())
}

// ExtractResources ...
func (e *Extractor) ExtractResources(pageHTML []byte) ogame.Resources {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractResourcesFromDoc(doc)
}

// ExtractResourcesFromDoc ...
func (e *Extractor) ExtractResourcesFromDoc(doc *goquery.Document) ogame.Resources {
	return extractResourcesFromDoc(doc)
}

// ExtractResourcesDetailsFromFullPage ...
func (e *Extractor) ExtractResourcesDetailsFromFullPage(pageHTML []byte) ogame.ResourcesDetails {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractResourcesDetailsFromFullPageFromDoc(doc)
}

// ExtractResourcesDetailsFromFullPageFromDoc ...
func (e *Extractor) ExtractResourcesDetailsFromFullPageFromDoc(doc *goquery.Document) ogame.ResourcesDetails {
	return extractResourcesDetailsFromFullPageFromDoc(doc)
}

// ExtractConstructions ...
func (e *Extractor) ExtractConstructions(pageHTML []byte) (buildingID ogame.ID, buildingCountdown int64, researchID ogame.ID, researchCountdown int64, lfBuildingID ogame.ID, lfBuildingCountdown int64, lfTechID ogame.ID, lfTechCountdown int64) {
	return ExtractConstructions(pageHTML, clockwork.NewRealClock())
}

// ExtractResourceSettings ...
func (e *Extractor) ExtractResourceSettings(pageHTML []byte) (ogame.ResourceSettings, string, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractResourceSettingsFromDoc(doc)
}

// ExtractLfBuildings ...
func (e *Extractor) ExtractLfBuildings(pageHTML []byte) (ogame.LfBuildings, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractLfBuildingsFromDoc(doc)
}

// ExtractLfBuildingsFromDoc ...
func (e *Extractor) ExtractLfBuildingsFromDoc(doc *goquery.Document) (ogame.LfBuildings, error) {
	return extractLfBuildingsFromDoc(doc)
}

// ExtractLfTechs ...
func (e *Extractor) ExtractLfTechs(pageHTML []byte) (ogame.LfTechs, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractLfTechsFromDoc(doc)
}

// ExtractLfTechsFromDoc ...
func (e *Extractor) ExtractLfTechsFromDoc(doc *goquery.Document) (ogame.LfTechs, error) {
	return extractLfTechsFromDoc(doc)
}