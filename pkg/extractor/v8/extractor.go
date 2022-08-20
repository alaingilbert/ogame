package v8

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/alaingilbert/ogame/pkg/extractor/v71"
	"github.com/alaingilbert/ogame/pkg/ogame"
)

// Extractor ...
type Extractor struct {
	v71.Extractor
}

// NewExtractor ...
func NewExtractor() *Extractor {
	return &Extractor{}
}

// ExtractIsInVacation ...
func (e *Extractor) ExtractIsInVacation(pageHTML []byte) bool {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractIsInVacationFromDoc(doc)
}

// ExtractIsInVacationFromDoc ...
func (e *Extractor) ExtractIsInVacationFromDoc(doc *goquery.Document) bool {
	return extractIsInVacationFromDocV8(doc)
}

// ExtractEspionageReport ...
func (e *Extractor) ExtractEspionageReport(pageHTML []byte) (ogame.EspionageReport, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractEspionageReportFromDoc(doc)
}

// ExtractEspionageReportFromDoc ...
func (e *Extractor) ExtractEspionageReportFromDoc(doc *goquery.Document) (ogame.EspionageReport, error) {
	return extractEspionageReportFromDocV8(doc, e.GetLocation())
}
