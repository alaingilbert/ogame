package ogame

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"time"
)

// ExtractorV8 ...
type ExtractorV8 struct {
	ExtractorV71
}

// NewExtractorV8 ...
func NewExtractorV8() *ExtractorV8 {
	return &ExtractorV8{}
}

// ExtractIsInVacation ...
func (e ExtractorV8) ExtractIsInVacation(pageHTML []byte) bool {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractIsInVacationFromDoc(doc)
}

// ExtractIsInVacationFromDoc ...
func (e ExtractorV8) ExtractIsInVacationFromDoc(doc *goquery.Document) bool {
	return extractIsInVacationFromDocV8(doc)
}

// ExtractEspionageReport ...
func (e ExtractorV8) ExtractEspionageReport(pageHTML []byte, location *time.Location) (EspionageReport, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractEspionageReportFromDoc(doc, location)
}

// ExtractEspionageReportFromDoc ...
func (e ExtractorV8) ExtractEspionageReportFromDoc(doc *goquery.Document, location *time.Location) (EspionageReport, error) {
	return extractEspionageReportFromDocV8(doc, location)
}
