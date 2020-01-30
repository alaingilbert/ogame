package ogame

import (
	"bytes"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// ExtractorV71 ...
type ExtractorV71 struct {
	ExtractorV7
}

// NewExtractorV71 ...
func NewExtractorV71() *ExtractorV71 {
	return &ExtractorV71{}
}

// ExtractFacilitiesFromDoc ...
func (e ExtractorV71) ExtractFacilitiesFromDoc(doc *goquery.Document) (Facilities, error) {
	return extractFacilitiesFromDocV71(doc)
}

// ExtractFacilities ...
func (e ExtractorV71) ExtractFacilities(pageHTML []byte) (Facilities, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractFacilitiesFromDoc(doc)
}

// ExtractResourcesDetails ...
func (e ExtractorV71) ExtractResourcesDetails(pageHTML []byte) (out ResourcesDetails, err error) {
	return extractResourcesDetailsV71(pageHTML)
}

// ExtractEspionageReport ...
func (e ExtractorV71) ExtractEspionageReport(pageHTML []byte, location *time.Location) (EspionageReport, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractEspionageReportFromDoc(doc, location)
}

// ExtractEspionageReportFromDoc ...
func (e ExtractorV71) ExtractEspionageReportFromDoc(doc *goquery.Document, location *time.Location) (EspionageReport, error) {
	return extractEspionageReportFromDocV71(doc, location)
}

// ExtractIPM ...
func (e ExtractorV71) ExtractIPM(pageHTML []byte) (duration int64, max int64, token string) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractIPMFromDoc(doc)
}

// ExtractIPMFromDoc ...
func (e ExtractorV71) ExtractIPMFromDoc(doc *goquery.Document) (duration int64, max int64, token string) {
	return extractIPMFromDocV71(doc)
}

// ExtractProduction extracts ships/defenses production from the shipyard page
func (e ExtractorV71) ExtractProduction(pageHTML []byte) ([]Quantifiable, int64, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	shipSumCountdown := e.ExtractOverviewShipSumCountdownFromBytes(pageHTML)
	production, err := e.ExtractProductionFromDoc(doc)
	return production, shipSumCountdown, err
}

// ExtractProductionFromDoc extracts ships/defenses production from the shipyard page
func (e ExtractorV71) ExtractProductionFromDoc(doc *goquery.Document) ([]Quantifiable, error) {
	return extractProductionFromDocV71(doc)
}

// ExtractHighscore ...
func (e ExtractorV71) ExtractHighscore(pageHTML []byte) (Highscore, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractHighscoreFromDoc(doc)
}

// ExtractHighscoreFromDoc ...
func (e ExtractorV71) ExtractHighscoreFromDoc(doc *goquery.Document) (Highscore, error) {
	return extractHighscoreFromDocV71(doc)
}
