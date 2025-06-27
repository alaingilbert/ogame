package v7

import (
	"bytes"
	"github.com/alaingilbert/ogame/pkg/extractor/v6"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/alaingilbert/clockwork"
)

// Extractor ...
type Extractor struct {
	v6.Extractor
}

// NewExtractor ...
func NewExtractor() *Extractor {
	return &Extractor{}
}

// ExtractPremiumToken ...
func (e Extractor) ExtractPremiumToken(pageHTML []byte, days int64) (string, error) {
	return extractPremiumToken(pageHTML, days)
}

// ExtractResourcesDetailsFromFullPage ...
func (e Extractor) ExtractResourcesDetailsFromFullPage(pageHTML []byte) (ogame.ResourcesDetails, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return ogame.ResourcesDetails{}, err
	}
	return e.ExtractResourcesDetailsFromFullPageFromDoc(doc), nil
}

// ExtractResourcesDetailsFromFullPageFromDoc ...
func (e Extractor) ExtractResourcesDetailsFromFullPageFromDoc(doc *goquery.Document) ogame.ResourcesDetails {
	return extractResourcesDetailsFromFullPageFromDoc(doc)
}

// ExtractExpeditionMessages ...
func (e Extractor) ExtractExpeditionMessages(pageHTML []byte) ([]ogame.ExpeditionMessage, int64, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return nil, 0, err
	}
	return e.ExtractExpeditionMessagesFromDoc(doc)
}

// ExtractMarketplaceMessages ...
func (e Extractor) ExtractMarketplaceMessages(pageHTML []byte) ([]ogame.MarketplaceMessage, int64, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return nil, 0, err
	}
	return e.ExtractMarketplaceMessagesFromDoc(doc, e.GetLocation())
}

// ExtractDefense ...
func (e Extractor) ExtractDefense(pageHTML []byte) (ogame.DefensesInfos, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return ogame.DefensesInfos{}, err
	}
	return e.ExtractDefenseFromDoc(doc)
}

// ExtractFacilities ...
func (e Extractor) ExtractFacilities(pageHTML []byte) (ogame.Facilities, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return ogame.Facilities{}, err
	}
	return e.ExtractFacilitiesFromDoc(doc)
}

// ExtractResearch ...
func (e Extractor) ExtractResearch(pageHTML []byte) (ogame.Researches, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return ogame.Researches{}, err
	}
	return e.ExtractResearchFromDoc(doc), nil
}

// ExtractShips ...
func (e Extractor) ExtractShips(pageHTML []byte) (ogame.ShipsInfos, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return ogame.ShipsInfos{}, err
	}
	return e.ExtractShipsFromDoc(doc)
}

// ExtractResourceSettings ...
func (e Extractor) ExtractResourceSettings(pageHTML []byte) (ogame.ResourceSettings, string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return ogame.ResourceSettings{}, "", err
	}
	return e.ExtractResourceSettingsFromDoc(doc)
}

// ExtractCharacterClass ...
func (e Extractor) ExtractCharacterClass(pageHTML []byte) (ogame.CharacterClass, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return ogame.NoClass, err
	}
	return e.ExtractCharacterClassFromDoc(doc)
}

// ExtractResourcesBuildings ...
func (e Extractor) ExtractResourcesBuildings(pageHTML []byte) (ogame.ResourcesBuildings, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return ogame.ResourcesBuildings{}, err
	}
	return e.ExtractResourcesBuildingsFromDoc(doc)
}

// ExtractResourcesDetails ...
func (e Extractor) ExtractResourcesDetails(pageHTML []byte) (out ogame.ResourcesDetails, err error) {
	return extractResourcesDetails(pageHTML)
}

// ExtractConstructions ...
func (e Extractor) ExtractConstructions(pageHTML []byte) (ogame.Constructions, error) {
	return ExtractConstructions(pageHTML, clockwork.NewRealClock()), nil
}

// ExtractFleet1Ships ...
func (e Extractor) ExtractFleet1Ships(pageHTML []byte) (ogame.ShipsInfos, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return ogame.ShipsInfos{}, err
	}
	return e.ExtractFleet1ShipsFromDoc(doc)
}

// ExtractCombatReportMessagesSummary ...
func (e Extractor) ExtractCombatReportMessagesSummary(pageHTML []byte) ([]ogame.CombatReportSummary, int64, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return nil, 0, err
	}
	return e.ExtractCombatReportMessagesFromDoc(doc)
}

// ExtractIPM ...
func (e Extractor) ExtractIPM(pageHTML []byte) (duration, max int64, token string, err error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return 0, 0, "", err
	}
	return e.ExtractIPMFromDoc(doc)
}

// ExtractIPMFromDoc ...
func (e Extractor) ExtractIPMFromDoc(doc *goquery.Document) (duration, max int64, token string, err error) {
	return extractIPMFromDoc(doc)
}

// ExtractEspionageReport ...
func (e Extractor) ExtractEspionageReport(pageHTML []byte) (ogame.EspionageReport, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return ogame.EspionageReport{}, err
	}
	return e.ExtractEspionageReportFromDoc(doc)
}

// ExtractOverviewProduction extracts ships/defenses (partial) production from the overview page
func (e Extractor) ExtractOverviewProduction(pageHTML []byte) ([]ogame.Quantifiable, int64, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return nil, 0, err
	}
	shipSumCountdown := e.ExtractOverviewShipSumCountdownFromBytes(pageHTML)
	production, err := e.ExtractOverviewProductionFromDoc(doc)
	return production, shipSumCountdown, err
}

// ExtractOverviewShipSumCountdownFromBytes ...
func (e Extractor) ExtractOverviewShipSumCountdownFromBytes(pageHTML []byte) int64 {
	return extractOverviewShipSumCountdownFromBytes(pageHTML)
}

// ExtractOverviewProductionFromDoc extracts ships/defenses (partial) production from the overview page
func (e Extractor) ExtractOverviewProductionFromDoc(doc *goquery.Document) ([]ogame.Quantifiable, error) {
	return extractOverviewProductionFromDoc(doc)
}

// ExtractFleet1ShipsFromDoc ...
func (e Extractor) ExtractFleet1ShipsFromDoc(doc *goquery.Document) (ogame.ShipsInfos, error) {
	return extractFleet1ShipsFromDoc(doc)
}

// ExtractResourceSettingsFromDoc ...
func (e Extractor) ExtractResourceSettingsFromDoc(doc *goquery.Document) (ogame.ResourceSettings, string, error) {
	return extractResourceSettingsFromDoc(doc)
}

// ExtractDefenseFromDoc ...
func (e Extractor) ExtractDefenseFromDoc(doc *goquery.Document) (ogame.DefensesInfos, error) {
	return extractDefenseFromDoc(doc)
}

// ExtractExpeditionMessagesFromDoc ...
func (e Extractor) ExtractExpeditionMessagesFromDoc(doc *goquery.Document) ([]ogame.ExpeditionMessage, int64, error) {
	return extractExpeditionMessagesFromDoc(doc, e.GetLocation())
}

// ExtractMarketplaceMessagesFromDoc ...
func (e Extractor) ExtractMarketplaceMessagesFromDoc(doc *goquery.Document, location *time.Location) ([]ogame.MarketplaceMessage, int64, error) {
	return extractMarketplaceMessagesFromDoc(doc, location)
}

// ExtractFacilitiesFromDoc ...
func (e Extractor) ExtractFacilitiesFromDoc(doc *goquery.Document) (ogame.Facilities, error) {
	return ExtractFacilitiesFromDoc(doc)
}

// ExtractResearchFromDoc ...
func (e Extractor) ExtractResearchFromDoc(doc *goquery.Document) ogame.Researches {
	return extractResearchFromDoc(doc)
}

// ExtractShipsFromDoc ...
func (e Extractor) ExtractShipsFromDoc(doc *goquery.Document) (ogame.ShipsInfos, error) {
	return extractShipsFromDoc(doc)
}

// ExtractResourcesBuildingsFromDoc ...
func (e Extractor) ExtractResourcesBuildingsFromDoc(doc *goquery.Document) (ogame.ResourcesBuildings, error) {
	return extractResourcesBuildingsFromDoc(doc)
}

// ExtractCombatReportMessagesFromDoc ...
func (e Extractor) ExtractCombatReportMessagesFromDoc(doc *goquery.Document) ([]ogame.CombatReportSummary, int64, error) {
	return extractCombatReportMessagesFromDoc(doc)
}

// ExtractEspionageReportFromDoc ...
func (e Extractor) ExtractEspionageReportFromDoc(doc *goquery.Document) (ogame.EspionageReport, error) {
	return extractEspionageReportFromDoc(doc, e.GetLocation())
}

// ExtractCancelBuildingInfos ...
func (e Extractor) ExtractCancelBuildingInfos(pageHTML []byte) (token string, techID, listID int64, err error) {
	return extractCancelBuildingInfos(pageHTML)
}

// ExtractCancelResearchInfos ...
func (e Extractor) ExtractCancelResearchInfos(pageHTML []byte) (token string, techID, listID int64, err error) {
	return extractCancelResearchInfos(pageHTML)
}

// ExtractCharacterClassFromDoc ...
func (e Extractor) ExtractCharacterClassFromDoc(doc *goquery.Document) (ogame.CharacterClass, error) {
	return extractCharacterClassFromDoc(doc)
}
