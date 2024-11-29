package v11_15_0

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/alaingilbert/ogame/pkg/extractor/v11_13_0"
	"github.com/alaingilbert/ogame/pkg/ogame"
)

// Extractor ...
type Extractor struct {
	v11_13_0.Extractor
}

// NewExtractor ...
func NewExtractor() *Extractor {
	return &Extractor{}
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

// ExtractEspionageReportMessageIDs ...
func (e *Extractor) ExtractEspionageReportMessageIDs(pageHTML []byte) ([]ogame.EspionageReportSummary, int64, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractEspionageReportMessageIDsFromDoc(doc)
}

// ExtractEspionageReportMessageIDsFromDoc ...
func (e *Extractor) ExtractEspionageReportMessageIDsFromDoc(doc *goquery.Document) ([]ogame.EspionageReportSummary, int64, error) {
	return extractEspionageReportMessageIDsFromDoc(doc)
}

// ExtractCombatReportMessagesSummary ...
func (e *Extractor) ExtractCombatReportMessagesSummary(pageHTML []byte) ([]ogame.CombatReportSummary, int64, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractCombatReportMessagesFromDoc(doc)
}

// ExtractCombatReportMessagesFromDoc ...
func (e *Extractor) ExtractCombatReportMessagesFromDoc(doc *goquery.Document) ([]ogame.CombatReportSummary, int64, error) {
	return extractCombatReportMessagesFromDoc(doc)
}

// ExtractExpeditionMessages ...
func (e *Extractor) ExtractExpeditionMessages(pageHTML []byte) ([]ogame.ExpeditionMessage, int64, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractExpeditionMessagesFromDoc(doc)
}

// ExtractExpeditionMessagesFromDoc ...
func (e *Extractor) ExtractExpeditionMessagesFromDoc(doc *goquery.Document) ([]ogame.ExpeditionMessage, int64, error) {
	return extractExpeditionMessagesFromDoc(doc, e.GetLocation())
}

// ExtractLfBonuses ...
func (e *Extractor) ExtractLfBonuses(pageHTML []byte) (ogame.LfBonuses, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractLfBonusesFromDoc(doc)
}

// ExtractLfBonusesFromDoc ...
func (e *Extractor) ExtractLfBonusesFromDoc(doc *goquery.Document) (ogame.LfBonuses, error) {
	return extractLfBonusesFromDoc(doc)
}

// ExtractAllianceClass ...
func (e Extractor) ExtractAllianceClass(pageHTML []byte) (ogame.AllianceClass, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractAllianceClassFromDoc(doc)
}

// ExtractAllianceClassFromDoc ...
func (e Extractor) ExtractAllianceClassFromDoc(doc *goquery.Document) (ogame.AllianceClass, error) {
	return extractAllianceClassFromDoc(doc)
}

// ExtractPhalanxNewToken ...
func (e *Extractor) ExtractPhalanxNewToken(pageHTML []byte) (string, error) {
	return extractPhalanxNewToken(pageHTML)
}
