package ogame

import (
	"net/url"
	"strconv"
)

// Page names
const (
	OverviewPage         = "overview"
	PreferencesPage      = "preferences"
	ResourceSettingsPage = "resourceSettings"
	DefensesPage         = "defenses"
	SuppliesPage         = "supplies"
	FacilitiesPage       = "facilities"
	FleetdispatchPage    = "fleetdispatch"
	ShipyardPage         = "shipyard"
	MovementPage         = "movement"
	ResearchPage         = "research"
	PlanetlayerPage      = "planetlayer"
	LogoutPage           = "logout"
	JumpgatelayerPage    = "jumpgatelayer"
	FetchResourcesPage   = "fetchResources"
	FetchTechs           = "fetchTechs"
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

	// ajax pages
	FetchEventboxAjaxPage      = "fetchEventbox"
	FetchResourcesAjaxPage     = "fetchResources"
	FetchTechsAjaxPage = "fetchTechs"
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
)

func (b *OGame) getPage(page string, celestialID CelestialID, opts ...Option) ([]byte, error) {
	vals := url.Values{"page": {"ingame"}, "component": {page}}
	if page == FetchResourcesPage || page == FetchTechs {
		vals = url.Values{"page": {page}}
	}
	if celestialID != 0 {
		vals.Add("cp", strconv.FormatInt(int64(celestialID), 10))
	}
	return b.getPageContent(vals, opts...)
}
