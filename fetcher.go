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

func (b *OGame) getPage(page string, celestialID CelestialID) ([]byte, error) {
	vals := url.Values{"page": {page}}
	if b.serverData.Version[0] == '7' {
		if page == DefensePage {
			page = "defenses"
		} else if page == ResourcesPage {
			page = "supplies"
		} else if page == StationPage {
			page = "facilities"
		} else if page == Fleet1Page {
			page = "fleetdispatch"
		}
		vals = url.Values{"page": {"ingame"}, "component": {page}}
	}
	if celestialID != 0 {
		vals.Add("cp", strconv.Itoa(int(celestialID)))
	}
	return b.getPageContent(vals)
}
