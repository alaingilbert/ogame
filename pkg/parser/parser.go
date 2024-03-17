package parser

import (
	"bytes"
	"errors"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/alaingilbert/ogame/pkg/extractor"
	v6 "github.com/alaingilbert/ogame/pkg/extractor/v6"
	"github.com/alaingilbert/ogame/pkg/ogame"
)

var ErrParsePageType = errors.New("failed to parse requested page type")

type Page struct {
	e       extractor.Extractor
	doc     *goquery.Document
	content []byte
}

func (p *Page) SetExtractor(ext extractor.Extractor) { p.e = ext }

func (p *Page) GetContent() []byte { return p.content }

func (p *Page) GetDoc() *goquery.Document {
	if p.doc == nil {
		doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(p.content))
		p.doc = doc
	}
	return p.doc
}

type EventListAjaxPage struct{ Page }
type MissileAttackLayerAjaxPage struct{ Page }
type FetchTechsAjaxPage struct{ Page }
type RocketlayerAjaxPage struct{ Page }
type PhalanxAjaxPage struct{ Page }
type JumpGateAjaxPage struct{ Page }

type FullPage struct{ Page }
type OverviewPage struct{ FullPage }
type PreferencesPage struct{ FullPage }
type SuppliesPage struct{ FullPage }
type ResourcesSettingsPage struct{ FullPage }
type ResearchPage struct{ FullPage }
type FacilitiesPage struct{ FullPage }
type ShipyardPage struct{ FullPage }
type DefensesPage struct{ FullPage }
type MovementPage struct{ FullPage }
type LfBuildingsPage struct{ FullPage }
type LfResearchPage struct{ FullPage }

type FullPagePages interface {
	OverviewPage |
		PreferencesPage |
		SuppliesPage |
		ResourcesSettingsPage |
		FacilitiesPage |
		LfBuildingsPage |
		LfResearchPage |
		//TraderOverviewPageContent |
		//TraderResourcesPageContent |
		ResearchPage |
		ShipyardPage |
		DefensesPage |
		//FleetDispatchPageContent |
		MovementPage
	//GalaxyPageContent |
	//AlliancePageContent |
	//PremiumPageContent |
	//ShopPageContent |
	//MessagesPageContent |
	//ChatPageContent |
	//CharacterClassSelectionPageContent |
	//BuddiesPageContent |
	//HighScorePageContent
}

type AjaxPagePages interface {
	EventListAjaxPage |
		MissileAttackLayerAjaxPage |
		FetchTechsAjaxPage |
		RocketlayerAjaxPage |
		PhalanxAjaxPage |
		JumpGateAjaxPage
}

type IFullPage interface {
	ExtractLifeformTypeFromDoc() ogame.LifeformType
	ExtractOGameSession() string
	ExtractIsInVacation() bool
	ExtractPlanets() []ogame.Planet
	ExtractPlanetCoordinate() (ogame.Coordinate, error)
	ExtractAjaxChatToken() (string, error)
	ExtractCharacterClass() (ogame.CharacterClass, error)
	ExtractCommander() bool
	ExtractAdmiral() bool
	ExtractEngineer() bool
	ExtractGeologist() bool
	ExtractTechnocrat() bool
	ExtractColonies() (int64, int64)
	ExtractServerTime() (time.Time, error)
}

func AutoParseFullPage(e extractor.Extractor, pageHTML []byte) (out IFullPage) {
	fullPage := FullPage{Page{e: e, content: pageHTML}}
	if bytes.Contains(pageHTML, []byte(`currentPage = "overview";`)) || bytes.Contains(pageHTML, []byte(`currentPage = "intro";`)) {
		out = &OverviewPage{fullPage}
	} else if bytes.Contains(pageHTML, []byte(`currentPage = "preferences";`)) {
		out = &PreferencesPage{fullPage}
	} else if bytes.Contains(pageHTML, []byte(`currentPage = "research";`)) {
		out = &ResearchPage{fullPage}
	} else {
		out = &fullPage
	}
	return out
}

// ParsePage given a pageHTML and an extractor for the game version this html represent,
// returns a page of type T
func ParsePage[T FullPagePages](e extractor.Extractor, pageHTML []byte) (*T, error) {
	var zero T
	fullPage := FullPage{Page{e: e, content: pageHTML}}
	switch any(zero).(type) {
	case OverviewPage:
		if bytes.Contains(pageHTML, []byte(`currentPage = "overview";`)) ||
			bytes.Contains(pageHTML, []byte(`currentPage = "intro";`)) {
			tt := T(OverviewPage{fullPage})
			return &tt, nil
		}
	case DefensesPage:
		if isDefensesPage(e, pageHTML) {
			tt := T(DefensesPage{fullPage})
			return &tt, nil
		}
	case ShipyardPage:
		if bytes.Contains(pageHTML, []byte(`currentPage = "shipyard";`)) {
			tt := T(ShipyardPage{fullPage})
			return &tt, nil
		}
	case ResearchPage:
		tt := T(ResearchPage{fullPage})
		return &tt, nil
	case FacilitiesPage:
		tt := T(FacilitiesPage{fullPage})
		return &tt, nil
	case LfBuildingsPage:
		tt := T(LfBuildingsPage{fullPage})
		return &tt, nil
	case LfResearchPage:
		tt := T(LfResearchPage{fullPage})
		return &tt, nil
	case SuppliesPage:
		tt := T(SuppliesPage{fullPage})
		return &tt, nil
	case ResourcesSettingsPage:
		tt := T(ResourcesSettingsPage{fullPage})
		return &tt, nil
	case PreferencesPage:
		tt := T(PreferencesPage{fullPage})
		return &tt, nil
	case MovementPage:
		tt := T(MovementPage{fullPage})
		return &tt, nil
	default:
		return &zero, errors.New("page type not implemented")
	}
	return &zero, ErrParsePageType
}

func ParseAjaxPage[T AjaxPagePages](e extractor.Extractor, pageHTML []byte) (T, error) {
	var zero T
	page := Page{e: e, content: pageHTML}
	switch any(zero).(type) {
	case EventListAjaxPage:
		return T(EventListAjaxPage{page}), nil
	case MissileAttackLayerAjaxPage:
		return T(MissileAttackLayerAjaxPage{page}), nil
	case RocketlayerAjaxPage:
		return T(RocketlayerAjaxPage{page}), nil
	case PhalanxAjaxPage:
		return T(PhalanxAjaxPage{page}), nil
	case JumpGateAjaxPage:
		return T(JumpGateAjaxPage{page}), nil
	case FetchTechsAjaxPage:
		return T(FetchTechsAjaxPage{page}), nil
	}
	return zero, ErrParsePageType
}

func isDefensesPage(e extractor.Extractor, pageHTML []byte) bool {
	var target string
	switch e.(type) {
	case *v6.Extractor:
		target = `currentPage="defense";`
	default:
		target = `currentPage = "defenses";`
	}
	return bytes.Contains(pageHTML, []byte(target))
}
