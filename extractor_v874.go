package ogame

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
)

// ExtractorV874 ...
type ExtractorV874 struct {
	ExtractorV8
}

// NewExtractorV874 ...
func NewExtractorV874() *ExtractorV874 {
	return &ExtractorV874{}
}

// ExtractOfferOfTheDay ...
func (e ExtractorV874) ExtractOfferOfTheDay(pageHTML []byte) (int64, string, PlanetResources, Multiplier, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractOfferOfTheDayFromDoc(doc)
}

// ExtractOfferOfTheDayFromDoc ...
func (e ExtractorV874) ExtractOfferOfTheDayFromDoc(doc *goquery.Document) (price int64, importToken string, planetResources PlanetResources, multiplier Multiplier, err error) {
	return extractOfferOfTheDayFromDocV874(doc)
}

// ExtractAuction ...
func (e ExtractorV874) ExtractAuction(pageHTML []byte) (Auction, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return extractAuctionFromDocV874(doc)
}

// ExtractBuffActivation ...
func (e ExtractorV874) ExtractBuffActivation(pageHTML []byte) (string, []Item, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractBuffActivationFromDoc(doc)
}

// ExtractBuffActivationFromDoc ...
func (e ExtractorV874) ExtractBuffActivationFromDoc(doc *goquery.Document) (string, []Item, error) {
	return extractBuffActivationFromDocV874(doc)
}
