package v12_0_0

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/alaingilbert/clockwork"
	"github.com/alaingilbert/ogame/pkg/extractor/v11_15_0"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"time"
)

// Extractor ...
type Extractor struct {
	v11_15_0.Extractor
}

// NewExtractor ...
func NewExtractor() *Extractor {
	return &Extractor{}
}

// ExtractServerTime ...
func (e *Extractor) ExtractServerTime(pageHTML []byte) (time.Time, error) {
	clock := clockwork.NewRealClock()
	return e.extractServerTime(pageHTML, clock)
}

func (e *Extractor) extractServerTime(pageHTML []byte, clock clockwork.Clock) (time.Time, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return extractServerTimeFromDoc(doc, clock)
}

// ExtractServerTimeFromDoc ...
func (e *Extractor) ExtractServerTimeFromDoc(doc *goquery.Document) (time.Time, error) {
	clock := clockwork.NewRealClock()
	return extractServerTimeFromDoc(doc, clock)
}

// ExtractHighscore ...
func (e *Extractor) ExtractHighscore(pageHTML []byte) (ogame.Highscore, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractHighscoreFromDoc(doc)
}

// ExtractHighscoreFromDoc ...
func (e *Extractor) ExtractHighscoreFromDoc(doc *goquery.Document) (ogame.Highscore, error) {
	return extractHighscoreFromDoc(doc)
}
