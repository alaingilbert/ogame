package ogame

import (
	"bytes"
	"errors"
)

var ErrParsePageType = errors.New("failed to parse requested page type")

// ParsePage given a pageHTML and an extractor for the game version this html represent,
// returns a page of type T
func ParsePage[T FullPagePages](b *OGame, pageHTML []byte) (T, error) {
	var zero T
	switch any(zero).(type) {
	case OverviewPage:
		if bytes.Contains(pageHTML, []byte(`currentPage = "overview";`)) {
			return T(OverviewPage{FullPage{Page{b: b, content: pageHTML}}}), nil
		}
	case DefensesPage:
		c, err := ParseDefensesPageContent(b, pageHTML)
		return T(c), err
	case ShipyardPage:
		if bytes.Contains(pageHTML, []byte(`currentPage = "shipyard";`)) {
			return T(ShipyardPage{FullPage{Page{b: b, content: pageHTML}}}), nil
		}
	case ResearchPage:
		return T(ResearchPage{FullPage{Page{b: b, content: pageHTML}}}), nil
	case FacilitiesPage:
		return T(FacilitiesPage{FullPage{Page{b: b, content: pageHTML}}}), nil
	case SuppliesPage:
		return T(SuppliesPage{FullPage{Page{b: b, content: pageHTML}}}), nil
	case ResourcesSettingsPage:
		return T(ResourcesSettingsPage{FullPage{Page{b: b, content: pageHTML}}}), nil
	case MovementPage:
		return T(MovementPage{FullPage{Page{b: b, content: pageHTML}}}), nil
	default:
		return zero, errors.New("page type not implemented")
	}
	return zero, ErrParsePageType
}

func ParseAjaxPage[T AjaxPagePages](b *OGame, pageHTML []byte) (T, error) {
	var zero T
	switch any(zero).(type) {
	case EventListAjaxPage:
		return T(EventListAjaxPage{Page{b: b, content: pageHTML}}), nil
	case MissileAttackLayerAjaxPage:
		return T(MissileAttackLayerAjaxPage{Page{b: b, content: pageHTML}}), nil
	case RocketlayerAjaxPage:
		return T(RocketlayerAjaxPage{Page{b: b, content: pageHTML}}), nil
	case FetchTechsAjaxPage:
		return T(FetchTechsAjaxPage{Page{b: b, content: pageHTML}}), nil
	}
	return zero, ErrParsePageType
}

func ParseDefensesPageContent(b *OGame, pageHTML []byte) (out DefensesPage, err error) {
	var target string
	switch b.extractor.(type) {
	case *ExtractorV6:
		target = `currentPage="defense";`
	default:
		target = `currentPage = "defenses";`
	}
	if bytes.Contains(pageHTML, []byte(target)) {
		return DefensesPage{FullPage{Page{b: b, content: pageHTML}}}, nil
	}
	return out, ErrParsePageType
}
