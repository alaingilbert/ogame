package parser

import (
	"bytes"
	"errors"
	"github.com/alaingilbert/ogame/pkg/utils"
	"strings"
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
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(p.content))
		if err != nil {
			doc, err = goquery.NewDocumentFromReader(strings.NewReader(""))
			if err != nil {
				panic(err)
			}
		}
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
type AllianceOverviewTabAjaxPage struct{ Page }

type FullPage struct{ Page }
type OverviewPage struct{ FullPage }
type PreferencesPage struct{ FullPage }
type SuppliesPage struct{ FullPage }
type ResourcesSettingsPage struct{ FullPage }
type ResearchPage struct{ FullPage }
type FacilitiesPage struct{ FullPage }
type ShipyardPage struct{ FullPage }
type FleetDispatchPage struct{ FullPage }
type DefensesPage struct{ FullPage }
type MovementPage struct{ FullPage }
type LfBuildingsPage struct{ FullPage }
type LfResearchPage struct{ FullPage }
type LfBonusesPage struct{ FullPage }

type FullPagePages interface {
	OverviewPage |
		PreferencesPage |
		SuppliesPage |
		ResourcesSettingsPage |
		LfBonusesPage |
		FacilitiesPage |
		LfBuildingsPage |
		LfResearchPage |
		//TraderOverviewPageContent |
		//TraderResourcesPageContent |
		ResearchPage |
		ShipyardPage |
		FleetDispatchPage |
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
		JumpGateAjaxPage |
		AllianceOverviewTabAjaxPage
}

type IFullPage interface {
	ExtractLifeformTypeFromDoc() ogame.LifeformType
	ExtractOGameSession() string
	ExtractIsInVacation() bool
	ExtractPlanets() []ogame.Planet
	ExtractPlanetID() (ogame.CelestialID, error)
	ExtractPlanetCoordinate() (ogame.Coordinate, error)
	ExtractAjaxChatToken() (string, error)
	ExtractToken() (string, error)
	ExtractCharacterClass() (ogame.CharacterClass, error)
	ExtractCommander() bool
	ExtractAdmiral() bool
	ExtractEngineer() bool
	ExtractGeologist() bool
	ExtractTechnocrat() bool
	ExtractColonies() (int64, int64)
	ExtractServerTime() (time.Time, error)
	ExtractResources() ogame.Resources
}

func AutoParseFullPage(e extractor.Extractor, pageHTML []byte) (out IFullPage) {
	fullPage := FullPage{Page{e: e, content: pageHTML}}
	if bytes.Contains(pageHTML, []byte(`currentPage = "overview";`)) || bytes.Contains(pageHTML, []byte(`currentPage = "intro";`)) {
		out = &OverviewPage{fullPage}
	} else if bytes.Contains(pageHTML, []byte(`currentPage = "preferences";`)) {
		out = &PreferencesPage{fullPage}
	} else if bytes.Contains(pageHTML, []byte(`currentPage = "research";`)) {
		out = &ResearchPage{fullPage}
	} else if bytes.Contains(pageHTML, []byte(`currentPage = "lfbonuses";`)) {
		out = &LfBonusesPage{fullPage}
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
			return utils.Ptr(T(OverviewPage{fullPage})), nil
		}
	case DefensesPage:
		if isDefensesPage(e, pageHTML) {
			return utils.Ptr(T(DefensesPage{fullPage})), nil
		}
	case ShipyardPage:
		if bytes.Contains(pageHTML, []byte(`currentPage = "shipyard";`)) {
			return utils.Ptr(T(ShipyardPage{fullPage})), nil
		}
	case FleetDispatchPage:
		if bytes.Contains(pageHTML, []byte(`currentPage = "fleetdispatch";`)) {
			return utils.Ptr(T(FleetDispatchPage{fullPage})), nil
		}
	case ResearchPage:
		return utils.Ptr(T(ResearchPage{fullPage})), nil
	case LfBonusesPage:
		return utils.Ptr(T(LfBonusesPage{fullPage})), nil
	case FacilitiesPage:
		return utils.Ptr(T(FacilitiesPage{fullPage})), nil
	case LfBuildingsPage:
		return utils.Ptr(T(LfBuildingsPage{fullPage})), nil
	case LfResearchPage:
		return utils.Ptr(T(LfResearchPage{fullPage})), nil
	case SuppliesPage:
		return utils.Ptr(T(SuppliesPage{fullPage})), nil
	case ResourcesSettingsPage:
		return utils.Ptr(T(ResourcesSettingsPage{fullPage})), nil
	case PreferencesPage:
		return utils.Ptr(T(PreferencesPage{fullPage})), nil
	case MovementPage:
		return utils.Ptr(T(MovementPage{fullPage})), nil
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
	case AllianceOverviewTabAjaxPage:
		return T(AllianceOverviewTabAjaxPage{page}), nil
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
