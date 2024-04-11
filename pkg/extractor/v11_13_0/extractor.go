package v11_13_0

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/alaingilbert/clockwork"
	"github.com/alaingilbert/ogame/pkg/extractor/v11_9_0"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"time"
)

// Extractor ...
type Extractor struct {
	v11_9_0.Extractor
}

// NewExtractor ...
func NewExtractor() *Extractor {
	return &Extractor{}
}

// ExtractConstructions ...
func (e *Extractor) ExtractConstructions(pageHTML []byte) (buildingID ogame.ID, buildingCountdown int64, researchID ogame.ID, researchCountdown int64, lfBuildingID ogame.ID, lfBuildingCountdown int64, lfResearchID ogame.ID, lfResearchCountdown int64) {
	return ExtractConstructions(pageHTML, clockwork.NewRealClock())
}

// ExtractProduction extracts ships/defenses production from the shipyard page
func (e *Extractor) ExtractProduction(pageHTML []byte) ([]ogame.Quantifiable, int64, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	shipSumCountdown := e.ExtractOverviewShipSumCountdownFromBytes(pageHTML)
	production, err := e.ExtractProductionFromDoc(doc)
	return production, shipSumCountdown, err
}

// ExtractOverviewShipSumCountdownFromBytes ...
func (e Extractor) ExtractOverviewShipSumCountdownFromBytes(pageHTML []byte) int64 {
	return extractOverviewShipSumCountdownFromBytes(pageHTML)
}

// ExtractFleets ...
func (e *Extractor) ExtractFleets(pageHTML []byte) (res []ogame.Fleet) {
	return e.extractFleets(pageHTML, e.GetLocation())
}

func (e *Extractor) extractFleets(pageHTML []byte, location *time.Location) (res []ogame.Fleet) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.extractFleetsFromDoc(doc, location)
}

// ExtractFleetsFromDoc ...
func (e *Extractor) ExtractFleetsFromDoc(doc *goquery.Document) (res []ogame.Fleet) {
	return e.extractFleetsFromDoc(doc, e.GetLocation())
}

func (e *Extractor) extractFleetsFromDoc(doc *goquery.Document, location *time.Location) (res []ogame.Fleet) {
	return extractFleetsFromDoc(doc, location, e.GetLifeformEnabled())
}
