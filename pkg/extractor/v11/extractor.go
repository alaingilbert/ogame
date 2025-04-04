package v11

import (
	"github.com/PuerkitoBio/goquery"
	v104 "github.com/alaingilbert/ogame/pkg/extractor/v104"
	"github.com/alaingilbert/ogame/pkg/ogame"
)

// Extractor ...
type Extractor struct {
	v104.Extractor
}

// NewExtractor ...
func NewExtractor() *Extractor {
	return &Extractor{}
}

// ExtractResourceSettings ...
func (e Extractor) ExtractResourceSettings(pageHTML []byte) (ogame.ResourceSettings, string, error) {
	return extractResourceSettingsFromPage(pageHTML)
}

// ExtractCancelBuildingInfos ...
func (e Extractor) ExtractCancelBuildingInfos(pageHTML []byte) (token string, techID, listID int64, err error) {
	return extractCancelBuildingInfos(pageHTML)
}

// ExtractCancelResearchInfos ...
func (e Extractor) ExtractCancelResearchInfos(pageHTML []byte) (token string, techID, listID int64, err error) {
	return extractCancelResearchInfos(pageHTML)
}

// ExtractCancelLfBuildingInfos ...
func (e Extractor) ExtractCancelLfBuildingInfos(pageHTML []byte) (token string, id, listID int64, err error) {
	return extractCancelLfBuildingInfos(pageHTML)
}

// ExtractLifeformTypeFromDoc ...
func (e Extractor) ExtractLifeformTypeFromDoc(doc *goquery.Document) ogame.LifeformType {
	return extractLifeformTypeFromDoc(doc)
}

// ExtractJumpGate return the available ships to send, form token, possible moon IDs and wait time (if any)
// given a jump gate popup html.
func (e *Extractor) ExtractJumpGate(pageHTML []byte) (ogame.ShipsInfos, string, []ogame.MoonID, int64, error) {
	return extractJumpGate(pageHTML)
}

// ExtractPreferencesFromDoc ...
func (e *Extractor) ExtractPreferencesFromDoc(doc *goquery.Document) ogame.Preferences {
	return extractPreferencesFromDoc(doc)
}
