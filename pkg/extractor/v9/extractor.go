package v9

import (
	"bytes"
	"encoding/json"
	"strings"

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

// ExtractTechnologyDetailsFromDoc ...
func (e *Extractor) ExtractTechnologyDetailsFromDoc(doc *goquery.Document) (ogame.TechnologyDetails, error) {
	return extractTechnologyDetailsFromDoc(doc)
}

type technologyDetailsStruct struct {
	Target  string `json:"target"`
	Content struct {
		Technologydetails string `json:"technologydetails"`
	} `json:"content"`
	Files struct {
		Js  []string `json:"js"`
		CSS []string `json:"css"`
	} `json:"files"`
	Page struct {
		StateObj string `json:"stateObj"`
		Title    string `json:"title"`
		URL      string `json:"url"`
	} `json:"page"`
	ServerTime int `json:"serverTime"`
}

// ExtractTechnologyDetails ...
func (e *Extractor) ExtractTechnologyDetails(pageHTML []byte) (out ogame.TechnologyDetails, err error) {
	var technologyDetails technologyDetailsStruct
	if err := json.Unmarshal(pageHTML, &technologyDetails); err != nil {
		return out, err
	}
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(technologyDetails.Content.Technologydetails))
	return e.ExtractTechnologyDetailsFromDoc(doc)
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
	return extractOverviewProductionFromDoc(doc, e.GetLifeformEnabled())
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
func (e *Extractor) ExtractConstructions(pageHTML []byte) (buildingID ogame.ID, buildingCountdown int64, researchID ogame.ID, researchCountdown int64, lfBuildingID ogame.ID, lfBuildingCountdown int64, lfResearchID ogame.ID, lfResearchCountdown int64) {
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

// ExtractLfResearch ...
func (e *Extractor) ExtractLfResearch(pageHTML []byte) (ogame.LfResearches, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractLfResearchFromDoc(doc)
}

// ExtractLfResearchFromDoc ...
func (e *Extractor) ExtractLfResearchFromDoc(doc *goquery.Document) (ogame.LfResearches, error) {
	return extractLfResearchFromDoc(doc)
}

// ExtractLfSlotsFromDoc ...
func (e *Extractor) ExtractLfSlotsFromDoc(doc *goquery.Document) [18]ogame.LfSlot {
	return extractLfSlotsFromDoc(doc)
}

// ExtractArtefactsFromDoc ...
func (e *Extractor) ExtractArtefactsFromDoc(doc *goquery.Document) (int64, int64) {
	return extractArtefactsFromDoc(doc)
}

// ExtractTearDownButtonEnabled ...
func (e *Extractor) ExtractTearDownButtonEnabled(pageHTML []byte) bool {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractTearDownButtonEnabledFromDoc(doc)
}

// ExtractTearDownButtonEnabledFromDoc ...
func (e *Extractor) ExtractTearDownButtonEnabledFromDoc(doc *goquery.Document) bool {
	return extractTearDownButtonEnabledFromDoc(doc)
}

// ExtractAvailableDiscoveries ...
func (e *Extractor) ExtractAvailableDiscoveries(pageHTML []byte) int64 {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.extractAvailableDiscoveriesFromDoc(doc)
}

// ExtractAvailableDiscoveriesFromDoc ...
func (e *Extractor) extractAvailableDiscoveriesFromDoc(doc *goquery.Document) int64 {
	return extractAvailableDiscoveriesFromDoc(doc)
}
