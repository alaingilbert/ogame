package v8

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/alaingilbert/ogame/pkg/extractor/v71"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"time"
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
func (e *Extractor) ExtractIsInVacation(pageHTML []byte) (bool, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return false, err
	}
	return e.ExtractIsInVacationFromDoc(doc), nil
}

// ExtractIsInVacationFromDoc ...
func (e *Extractor) ExtractIsInVacationFromDoc(doc *goquery.Document) bool {
	return extractIsInVacationFromDoc(doc)
}

// ExtractAttackBlock ...
func (e *Extractor) ExtractAttackBlock(pageHTML []byte) (bool, time.Time, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return false, time.Time{}, err
	}
	attackBlockActivated, blockedUntil := e.ExtractAttackBlockFromDoc(doc)
	return attackBlockActivated, blockedUntil, nil
}

// ExtractAttackBlockFromDoc ...
func (e *Extractor) ExtractAttackBlockFromDoc(doc *goquery.Document) (bool, time.Time) {
	return extractAttackBlockFromDoc(doc)
}

// ExtractEspionageReport ...
func (e *Extractor) ExtractEspionageReport(pageHTML []byte) (ogame.EspionageReport, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return ogame.EspionageReport{}, err
	}
	return e.ExtractEspionageReportFromDoc(doc)
}

// ExtractEspionageReportFromDoc ...
func (e *Extractor) ExtractEspionageReportFromDoc(doc *goquery.Document) (ogame.EspionageReport, error) {
	return extractEspionageReportFromDoc(doc, e.GetLocation())
}
