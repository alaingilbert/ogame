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
func (e *Extractor) ExtractConstructions(pageHTML []byte) (ogame.Constructions, error) {
	return ExtractConstructions(pageHTML, clockwork.NewRealClock())
}

// ExtractProduction extracts ships/defenses production from the shipyard page
func (e *Extractor) ExtractProduction(pageHTML []byte) ([]ogame.Quantifiable, int64, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return nil, 0, err
	}
	shipSumCountdown := e.ExtractOverviewShipSumCountdownFromBytes(pageHTML)
	production, err := e.ExtractProductionFromDoc(doc)
	return production, shipSumCountdown, err
}

// ExtractOverviewShipSumCountdownFromBytes ...
func (e Extractor) ExtractOverviewShipSumCountdownFromBytes(pageHTML []byte) int64 {
	return extractOverviewShipSumCountdownFromBytes(pageHTML)
}

// ExtractFleets ...
func (e *Extractor) ExtractFleets(pageHTML []byte) ([]ogame.Fleet, error) {
	return e.extractFleets(pageHTML, e.GetLocation())
}

func (e *Extractor) extractFleets(pageHTML []byte, location *time.Location) ([]ogame.Fleet, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return nil, err
	}
	return e.extractFleetsFromDoc(doc, location), nil
}

// ExtractFleetsFromDoc ...
func (e *Extractor) ExtractFleetsFromDoc(doc *goquery.Document) (res []ogame.Fleet) {
	return e.extractFleetsFromDoc(doc, e.GetLocation())
}

func (e *Extractor) extractFleetsFromDoc(doc *goquery.Document, location *time.Location) []ogame.Fleet {
	return extractFleetsFromDoc(doc, location, e.GetLifeformEnabled())
}
