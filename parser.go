package ogame

import (
	"bytes"
	"errors"
)

var ErrParsePageType = errors.New("failed to parse requested page type")

func AutoParseFullPage(b *OGame, pageHTML []byte) IFullPage {
	fullPage := FullPage{Page{b: b, content: pageHTML}}
	if bytes.Contains(pageHTML, []byte(`currentPage = "overview";`)) {
		return OverviewPage{fullPage}
	} else if bytes.Contains(pageHTML, []byte(`currentPage = "preferences";`)) {
		return PreferencesPage{fullPage}
	} else if bytes.Contains(pageHTML, []byte(`currentPage = "research";`)) {
		return ResearchPage{fullPage}
	}
	return fullPage
}

// ParsePage given a pageHTML and an extractor for the game version this html represent,
// returns a page of type T
func ParsePage[T FullPagePages](b *OGame, pageHTML []byte) (T, error) {
	var zero T
	fullPage := FullPage{Page{b: b, content: pageHTML}}
	switch any(zero).(type) {
	case OverviewPage:
		if bytes.Contains(pageHTML, []byte(`currentPage = "overview";`)) {
			return T(OverviewPage{fullPage}), nil
		}
	case DefensesPage:
		if isDefensesPage(b.extractor, pageHTML) {
			return T(DefensesPage{fullPage}), nil
		}
	case ShipyardPage:
		if bytes.Contains(pageHTML, []byte(`currentPage = "shipyard";`)) {
			return T(ShipyardPage{fullPage}), nil
		}
	case ResearchPage:
		return T(ResearchPage{fullPage}), nil
	case FacilitiesPage:
		return T(FacilitiesPage{fullPage}), nil
	case SuppliesPage:
		return T(SuppliesPage{fullPage}), nil
	case ResourcesSettingsPage:
		return T(ResourcesSettingsPage{fullPage}), nil
	case PreferencesPage:
		return T(PreferencesPage{fullPage}), nil
	case MovementPage:
		return T(MovementPage{fullPage}), nil
	default:
		return zero, errors.New("page type not implemented")
	}
	return zero, ErrParsePageType
}

func ParseAjaxPage[T AjaxPagePages](b *OGame, pageHTML []byte) (T, error) {
	var zero T
	page := Page{b: b, content: pageHTML}
	switch any(zero).(type) {
	case EventListAjaxPage:
		return T(EventListAjaxPage{page}), nil
	case MissileAttackLayerAjaxPage:
		return T(MissileAttackLayerAjaxPage{page}), nil
	case RocketlayerAjaxPage:
		return T(RocketlayerAjaxPage{page}), nil
	case PhalanxAjaxPage:
		return T(PhalanxAjaxPage{page}), nil
	case FetchTechsAjaxPage:
		return T(FetchTechsAjaxPage{page}), nil
	}
	return zero, ErrParsePageType
}

func isDefensesPage(e Extractor, pageHTML []byte) bool {
	var target string
	switch e.(type) {
	case *ExtractorV6:
		target = `currentPage="defense";`
	default:
		target = `currentPage = "defenses";`
	}
	return bytes.Contains(pageHTML, []byte(target))
}
