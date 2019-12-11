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
