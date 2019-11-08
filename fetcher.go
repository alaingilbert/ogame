package ogame

import (
	"net/url"
	"strconv"
)

const (
	OverviewPage         = "overview"
	PreferencesPage      = "preferences"
	ResourceSettingsPage = "resourceSettings"
	ResourcesPage        = "resources"
	DefensePage          = "defense"
	ShipyardPage         = "shipyard"
	StationPage          = "station"
	MovementPage         = "movement"
	ResearchPage         = "research"
	PlanetlayerPage      = "planetlayer"
	LogoutPage           = "logout"
	Fleet1Page           = "fleet1"
	JumpgatelayerPage    = "jumpgatelayer"
	FetchResourcesPage   = "fetchResources"
)

var pageV7Mapping = map[string]string{
	DefensePage:   "defenses",
	ResourcesPage: "supplies",
	StationPage:   "facilities",
	Fleet1Page:    "fleetdispatch",
}

func (b *OGame) getPage(page string, celestialID CelestialID) ([]byte, error) {
	vals := url.Values{"page": {page}}
	if b.serverData.Version[0] == '7' {
		if newPage, ok := pageV7Mapping[page]; ok {
			page = newPage
		}
		vals = url.Values{"page": {"ingame"}, "component": {page}}
	}
	if celestialID != 0 {
		vals.Add("cp", strconv.Itoa(int(celestialID)))
	}
	return b.getPageContent(vals)
}
