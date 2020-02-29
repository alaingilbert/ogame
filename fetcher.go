package ogame

import (
	"net/url"
	"strconv"
)

// Page names
const (
	// V6 full pages
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
	TraderOverviewPage   = "traderOverview"
	GalaxyPage           = "galaxy"
	AlliancePage         = "alliance"
	PremiumPage          = "premium"
	ShopPage             = "shop"
	RewardsPage          = "rewards"
	HighscorePage        = "highscore"
	BuddiesPage          = "buddies"
	MessagesPage         = "messages"
	ChatPage             = "chat"

	// V6 ajax pages
	FetchEventboxAjaxPage      = "fetchEventbox"
	FetchResourcesAjaxPage     = "fetchResources"
	GalaxyContentAjaxPage      = "galaxyContent"
	EventListAjaxPage          = "eventList"
	AjaxChatAjaxPage           = "ajaxChat"
	NoticesAjaxPage            = "notices"
	RepairlayerAjaxPage        = "repairlayer"
	TechtreeAjaxPage           = "techtree"
	PhalanxAjaxPage            = "phalanx"
	ShareReportOverlayAjaxPage = "shareReportOverlay"
	JumpgatelayerAjaxPage      = "jumpgatelayer"
	FederationlayerAjaxPage    = "federationlayer"
	UnionchangeAjaxPage        = "unionchange"
	ChangenickAjaxPage         = "changenick"
	PlanetlayerAjaxPage        = "planetlayer"
	TraderlayerAjaxPage        = "traderlayer"
	PlanetRenameAjaxPage       = "planetRename"
	RightmenuAjaxPage          = "rightmenu"
	AllianceOverviewAjaxPage   = "allianceOverview"
	SupportAjaxPage            = "support"
	BuffActivationAjaxPage     = "buffActivation"
	AuctioneerAjaxPage         = "auctioneer"
	HighscoreContentAjaxPage   = "highscoreContent"

	// V7 pages
	DefensesPage      = "defenses"
	SuppliesPage      = "supplies"
	FacilitiesPage    = "facilities"
	FleetdispatchPage = "fleetdispatch"
)

var pageV7Mapping = map[string]string{
	DefensePage:   DefensesPage,
	ResourcesPage: SuppliesPage,
	StationPage:   FacilitiesPage,
	Fleet1Page:    FleetdispatchPage,
}

func (b *OGame) getPage(page string, celestialID CelestialID, opts ...Option) ([]byte, error) {
	vals := url.Values{"page": {page}}
	if b.IsV7() && page != FetchResourcesPage {
		if newPage, ok := pageV7Mapping[page]; ok {
			page = newPage
		}
		vals = url.Values{"page": {"ingame"}, "component": {page}}
	}
	if celestialID != 0 {
		vals.Add("cp", strconv.FormatInt(int64(celestialID), 10))
	}
	return b.getPageContent(vals, opts...)
}
