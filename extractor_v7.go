package ogame

import (
	"bytes"

	"github.com/PuerkitoBio/goquery"
)

// ExtractorV7 ...
type ExtractorV7 struct {
	ExtractorV6
}

// NewExtractorV7 ...
func NewExtractorV7() *ExtractorV7 {
	return &ExtractorV7{}
}

// ExtractDefense ...
func (e ExtractorV7) ExtractDefense(pageHTML []byte) (DefensesInfos, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractDefenseFromDoc(doc)
}

// ExtractFacilities ...
func (e ExtractorV7) ExtractFacilities(pageHTML []byte) (Facilities, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractFacilitiesFromDoc(doc)
}

// ExtractResearch ...
func (e ExtractorV7) ExtractResearch(pageHTML []byte) Researches {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractResearchFromDoc(doc)
}

// ExtractDefenseFromDoc ...
func (e ExtractorV7) ExtractDefenseFromDoc(doc *goquery.Document) (DefensesInfos, error) {
	return extractDefenseFromDocV7(doc)
}

// ExtractFacilitiesFromDoc ...
func (e ExtractorV7) ExtractFacilitiesFromDoc(doc *goquery.Document) (Facilities, error) {
	return extractFacilitiesFromDocV7(doc)
}

// ExtractResearchFromDoc ...
func (e ExtractorV7) ExtractResearchFromDoc(doc *goquery.Document) Researches {
	return extractResearchFromDocV7(doc)
}
