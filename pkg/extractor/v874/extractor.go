package v874

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/alaingilbert/ogame/pkg/extractor/v8"
	"github.com/alaingilbert/ogame/pkg/ogame"
)

// Extractor ...
type Extractor struct {
	v8.Extractor
}

// NewExtractor ...
func NewExtractor() *Extractor {
	return &Extractor{}
}

// ExtractOfferOfTheDay ...
func (e *Extractor) ExtractOfferOfTheDay(pageHTML []byte) (int64, string, ogame.PlanetResources, ogame.Multiplier, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractOfferOfTheDayFromDoc(doc)
}

// ExtractOfferOfTheDayFromDoc ...
func (e *Extractor) ExtractOfferOfTheDayFromDoc(doc *goquery.Document) (price int64, importToken string, planetResources ogame.PlanetResources, multiplier ogame.Multiplier, err error) {
	return extractOfferOfTheDayFromDocV874(doc)
}

// ExtractAuction ...
func (e *Extractor) ExtractAuction(pageHTML []byte) (ogame.Auction, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return extractAuctionFromDocV874(doc)
}

// ExtractBuffActivation ...
func (e *Extractor) ExtractBuffActivation(pageHTML []byte) (string, []ogame.Item, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractBuffActivationFromDoc(doc)
}

// ExtractBuffActivationFromDoc ...
func (e *Extractor) ExtractBuffActivationFromDoc(doc *goquery.Document) (string, []ogame.Item, error) {
	return extractBuffActivationFromDocV874(doc)
}
