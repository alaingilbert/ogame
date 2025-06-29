package v71

import (
	"bytes"

	"github.com/PuerkitoBio/goquery"
	"github.com/alaingilbert/clockwork"
	v7 "github.com/alaingilbert/ogame/pkg/extractor/v7"
	"github.com/alaingilbert/ogame/pkg/ogame"
)

// Extractor ...
type Extractor struct {
	v7.Extractor
}

// NewExtractor ...
func NewExtractor() *Extractor {
	return &Extractor{}
}

// ExtractCancelFleetToken ...
func (e *Extractor) ExtractCancelFleetToken(pageHTML []byte, fleetID ogame.FleetID) (string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return "", err
	}
	return e.ExtractCancelFleetTokenFromDoc(doc, fleetID)
}

// ExtractCancelFleetTokenFromDoc ...
func (e *Extractor) ExtractCancelFleetTokenFromDoc(doc *goquery.Document, fleetID ogame.FleetID) (string, error) {
	return extractCancelFleetTokenFromDoc(doc, fleetID)
}

// ExtractFacilitiesFromDoc ...
func (e *Extractor) ExtractFacilitiesFromDoc(doc *goquery.Document) (ogame.Facilities, error) {
	return extractFacilitiesFromDoc(doc)
}

// ExtractFacilities ...
func (e *Extractor) ExtractFacilities(pageHTML []byte) (ogame.Facilities, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return ogame.Facilities{}, err
	}
	return e.ExtractFacilitiesFromDoc(doc)
}

// ExtractResourcesDetails ...
func (e *Extractor) ExtractResourcesDetails(pageHTML []byte) (out ogame.ResourcesDetails, err error) {
	return extractResourcesDetails(pageHTML)
}

// ExtractTechs ...
func (e *Extractor) ExtractTechs(pageHTML []byte) (ogame.Techs, error) {
	return extractTechs(pageHTML)
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

// ExtractDestroyRockets ...
func (e *Extractor) ExtractDestroyRockets(pageHTML []byte) (abm, ipm int64, token string, err error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return
	}
	return e.ExtractDestroyRocketsFromDoc(doc)
}

// ExtractDestroyRocketsFromDoc ...
func (e *Extractor) ExtractDestroyRocketsFromDoc(doc *goquery.Document) (abm, ipm int64, token string, err error) {
	return extractDestroyRocketsFromDoc(doc)
}

// ExtractIPM ...
func (e *Extractor) ExtractIPM(pageHTML []byte) (duration int64, max int64, token string, err error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return 0, 0, "", err
	}
	return e.ExtractIPMFromDoc(doc)
}

// ExtractIPMFromDoc ...
func (e *Extractor) ExtractIPMFromDoc(doc *goquery.Document) (duration int64, max int64, token string, err error) {
	return extractIPMFromDoc(doc)
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

// ExtractProductionFromDoc extracts ships/defenses production from the shipyard page
func (e *Extractor) ExtractProductionFromDoc(doc *goquery.Document) ([]ogame.Quantifiable, error) {
	return extractProductionFromDoc(doc)
}

// ExtractHighscore ...
func (e *Extractor) ExtractHighscore(pageHTML []byte) (ogame.Highscore, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return ogame.Highscore{}, err
	}
	return e.ExtractHighscoreFromDoc(doc)
}

// ExtractHighscoreFromDoc ...
func (e *Extractor) ExtractHighscoreFromDoc(doc *goquery.Document) (ogame.Highscore, error) {
	return extractHighscoreFromDoc(doc)
}

// ExtractAllResources ...
func (e *Extractor) ExtractAllResources(pageHTML []byte) (map[ogame.CelestialID]ogame.Resources, error) {
	return extractAllResources(pageHTML)
}

// ExtractAttacksFromDoc ...
func (e *Extractor) ExtractAttacksFromDoc(doc *goquery.Document, ownCoords []ogame.Coordinate) ([]ogame.AttackEvent, error) {
	return e.extractAttacksFromDoc(doc, clockwork.NewRealClock(), ownCoords)
}

// ExtractAttacks ...
func (e *Extractor) ExtractAttacks(pageHTML []byte, ownCoords []ogame.Coordinate) ([]ogame.AttackEvent, error) {
	return e.extractAttacks(pageHTML, clockwork.NewRealClock(), ownCoords)
}

func (e *Extractor) extractAttacks(pageHTML []byte, clock clockwork.Clock, ownCoords []ogame.Coordinate) ([]ogame.AttackEvent, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return nil, err
	}
	return e.extractAttacksFromDoc(doc, clock, ownCoords)
}

func (e *Extractor) extractAttacksFromDoc(doc *goquery.Document, clock clockwork.Clock, ownCoords []ogame.Coordinate) ([]ogame.AttackEvent, error) {
	return extractAttacksFromDoc(doc, clock, ownCoords)
}

// ExtractDMCosts ...
func (e *Extractor) ExtractDMCosts(pageHTML []byte) (ogame.DMCosts, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return ogame.DMCosts{}, err
	}
	return e.ExtractDMCostsFromDoc(doc)
}

// ExtractDMCostsFromDoc ...
func (e *Extractor) ExtractDMCostsFromDoc(doc *goquery.Document) (ogame.DMCosts, error) {
	return extractDMCostsFromDoc(doc)
}

// ExtractBuffActivation ...
func (e *Extractor) ExtractBuffActivation(pageHTML []byte) (string, []ogame.Item, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return "", nil, err
	}
	return e.ExtractBuffActivationFromDoc(doc)
}

// ExtractBuffActivationFromDoc ...
func (e *Extractor) ExtractBuffActivationFromDoc(doc *goquery.Document) (string, []ogame.Item, error) {
	return extractBuffActivationFromDoc(doc)
}

// ExtractActiveItems ...
func (e *Extractor) ExtractActiveItems(pageHTML []byte) ([]ogame.ActiveItem, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return nil, err
	}
	return e.ExtractActiveItemsFromDoc(doc)
}

// ExtractActiveItemsFromDoc ...
func (e *Extractor) ExtractActiveItemsFromDoc(doc *goquery.Document) ([]ogame.ActiveItem, error) {
	return extractActiveItemsFromDoc(doc)
}

// ExtractIsMobile ...
func (e *Extractor) ExtractIsMobile(pageHTML []byte) (bool, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return false, err
	}
	return e.ExtractIsMobileFromDoc(doc), nil
}

// ExtractIsMobileFromDoc ...
func (e *Extractor) ExtractIsMobileFromDoc(doc *goquery.Document) bool {
	return extractIsMobileFromDoc(doc)
}
