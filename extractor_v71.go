package ogame

import (
	"bytes"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/alaingilbert/clockwork"
)

// ExtractorV71 ...
type ExtractorV71 struct {
	ExtractorV7
}

// NewExtractorV71 ...
func NewExtractorV71() *ExtractorV71 {
	return &ExtractorV71{}
}

// ExtractFacilitiesFromDoc ...
func (e ExtractorV71) ExtractFacilitiesFromDoc(doc *goquery.Document) (Facilities, error) {
	return extractFacilitiesFromDocV71(doc)
}

// ExtractFacilities ...
func (e ExtractorV71) ExtractFacilities(pageHTML []byte) (Facilities, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractFacilitiesFromDoc(doc)
}

// ExtractResourcesDetails ...
func (e ExtractorV71) ExtractResourcesDetails(pageHTML []byte) (out ResourcesDetails, err error) {
	return extractResourcesDetailsV71(pageHTML)
}

// ExtractEspionageReport ...
func (e ExtractorV71) ExtractEspionageReport(pageHTML []byte, location *time.Location) (EspionageReport, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractEspionageReportFromDoc(doc, location)
}

// ExtractEspionageReportFromDoc ...
func (e ExtractorV71) ExtractEspionageReportFromDoc(doc *goquery.Document, location *time.Location) (EspionageReport, error) {
	return extractEspionageReportFromDocV71(doc, location)
}

// ExtractIPM ...
func (e ExtractorV71) ExtractIPM(pageHTML []byte) (duration int64, max int64, token string) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractIPMFromDoc(doc)
}

// ExtractIPMFromDoc ...
func (e ExtractorV71) ExtractIPMFromDoc(doc *goquery.Document) (duration int64, max int64, token string) {
	return extractIPMFromDocV71(doc)
}

// ExtractProduction extracts ships/defenses production from the shipyard page
func (e ExtractorV71) ExtractProduction(pageHTML []byte) ([]Quantifiable, int64, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	shipSumCountdown := e.ExtractOverviewShipSumCountdownFromBytes(pageHTML)
	production, err := e.ExtractProductionFromDoc(doc)
	return production, shipSumCountdown, err
}

// ExtractProductionFromDoc extracts ships/defenses production from the shipyard page
func (e ExtractorV71) ExtractProductionFromDoc(doc *goquery.Document) ([]Quantifiable, error) {
	return extractProductionFromDocV71(doc)
}

// ExtractHighscore ...
func (e ExtractorV71) ExtractHighscore(pageHTML []byte) (Highscore, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractHighscoreFromDoc(doc)
}

// ExtractHighscoreFromDoc ...
func (e ExtractorV71) ExtractHighscoreFromDoc(doc *goquery.Document) (Highscore, error) {
	return extractHighscoreFromDocV71(doc)
}

// ExtractAllResources ...
func (e ExtractorV71) ExtractAllResources(pageHTML []byte) (map[CelestialID]Resources, error) {
	return extractAllResourcesV71(pageHTML)
}

// ExtractAttacksFromDoc ...
func (e ExtractorV71) ExtractAttacksFromDoc(doc *goquery.Document, clock clockwork.Clock) ([]AttackEvent, error) {
	return extractAttacksFromDocV71(doc, clock)
}

// ExtractAttacks ...
func (e ExtractorV71) ExtractAttacks(pageHTML []byte) ([]AttackEvent, error) {
	return e.extractAttacks(pageHTML, clockwork.NewRealClock())
}

func (e ExtractorV71) extractAttacks(pageHTML []byte, clock clockwork.Clock) ([]AttackEvent, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractAttacksFromDoc(doc, clock)
}

// DMCost ...
type DMCost struct {
	Cost                int64
	CanBuy              bool  // Either or not we have enough DM
	Complete            bool  // false means we will halve the time, true will complete
	OGameID             ID    // What we are going to build
	Nbr                 int64 // Either the amount of ships/defences or the building/research level
	BuyAndActivateToken string
	Token               string
}

// String ...
func (d DMCost) String() string {
	return "\n" +
		"               Cost: " + strconv.FormatInt(d.Cost, 10) + "\n" +
		"             CanBuy: " + strconv.FormatBool(d.CanBuy) + "\n" +
		"           Complete: " + strconv.FormatBool(d.Complete) + "\n" +
		"            OGameID: " + strconv.FormatInt(int64(d.OGameID), 10) + "\n" +
		"                Nbr: " + strconv.FormatInt(d.Nbr, 10) + "\n" +
		"BuyAndActivateToken: " + d.BuyAndActivateToken + "\n" +
		"              Token: " + d.Token
}

// DMCosts ...
type DMCosts struct {
	Buildings DMCost
	Research  DMCost
	Shipyard  DMCost
}

// String ...
func (d DMCosts) String() string {
	return "\n" +
		"Buildings:" + d.Buildings.String() + "\n" +
		"Research:" + d.Research.String() + "\n" +
		"Shipyard:" + d.Shipyard.String()
}

// ExtractDMCosts ...
func (e ExtractorV71) ExtractDMCosts(pageHTML []byte) (DMCosts, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractDMCostsFromDoc(doc)
}

// ExtractDMCostsFromDoc ...
func (e ExtractorV71) ExtractDMCostsFromDoc(doc *goquery.Document) (DMCosts, error) {
	return extractDMCostsFromDocV71(doc)
}

// ExtractBuffActivation ...
func (e ExtractorV71) ExtractBuffActivation(pageHTML []byte) (string, []Item, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractBuffActivationFromDoc(doc)
}

// ExtractBuffActivationFromDoc ...
func (e ExtractorV71) ExtractBuffActivationFromDoc(doc *goquery.Document) (string, []Item, error) {
	return extractBuffActivationFromDocV71(doc)
}

// ExtractIsMobile ...
func (e ExtractorV71) ExtractIsMobile(pageHTML []byte) bool {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	return e.ExtractIsMobileFromDoc(doc)
}

// ExtractIsMobileFromDoc ...
func (e ExtractorV71) ExtractIsMobileFromDoc(doc *goquery.Document) bool {
	return extractIsMobileFromDocV71(doc)
}
