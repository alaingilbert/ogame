package ogame

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"time"
)

// ExtractorV9 ...
type ExtractorV9 struct {
	ExtractorV874
}

// NewExtractorV9 ...
func NewExtractorV9() *ExtractorV9 {
	return &ExtractorV9{}
}

// ExtractEmpire ...
func (e ExtractorV9) ExtractEmpire(pageHTML []byte) ([]EmpireCelestial, error) {
	return extractEmpireV9(pageHTML)
}

// ExtractOverviewProduction extracts ships/defenses (partial) production from the overview page
func (e ExtractorV9) ExtractOverviewProduction(pageHTML []byte) ([]Quantifiable, int64, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	shipSumCountdown := e.ExtractOverviewShipSumCountdownFromBytes(pageHTML)
	production, err := e.ExtractOverviewProductionFromDoc(doc)
	return production, shipSumCountdown, err
}

// ExtractOverviewProductionFromDoc extracts ships/defenses (partial) production from the overview page
func (e ExtractorV9) ExtractOverviewProductionFromDoc(doc *goquery.Document) ([]Quantifiable, error) {
	return extractOverviewProductionFromDocV9(doc)
}

// ExtractEspionageReport ...
func (e ExtractorV9) ExtractEspionageReport(pageHTML []byte, location *time.Location) (EspionageReport, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractEspionageReportFromDoc(doc, location)
}

// ExtractEspionageReportFromDoc ...
func (e ExtractorV9) ExtractEspionageReportFromDoc(doc *goquery.Document, location *time.Location) (EspionageReport, error) {
	return extractEspionageReportFromDocV9(doc, location)
}

// ExtractResources ...
func (e ExtractorV9) ExtractResources(pageHTML []byte) Resources {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractResourcesFromDoc(doc)
}

// ExtractResourcesFromDoc ...
func (e ExtractorV9) ExtractResourcesFromDoc(doc *goquery.Document) Resources {
	return extractResourcesFromDocV9(doc)
}

// ExtractResourcesDetailsFromFullPage ...
func (e ExtractorV9) ExtractResourcesDetailsFromFullPage(pageHTML []byte) ResourcesDetails {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractResourcesDetailsFromFullPageFromDoc(doc)
}

// ExtractResourcesDetailsFromFullPageFromDoc ...
func (e ExtractorV9) ExtractResourcesDetailsFromFullPageFromDoc(doc *goquery.Document) ResourcesDetails {
	return extractResourcesDetailsFromFullPageFromDocV9(doc)
}
