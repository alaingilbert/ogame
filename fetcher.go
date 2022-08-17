package ogame

import (
	"net/url"
)

// Page names
const (
	OverviewPageName         = "overview"
	PreferencesPageName      = "preferences"
	ResourceSettingsPageName = "resourceSettings"
	DefensesPageName         = "defenses"
	SuppliesPageName         = "supplies"
	FacilitiesPageName       = "facilities"
	FleetdispatchPageName    = "fleetdispatch"
	ShipyardPageName         = "shipyard"
	MovementPageName         = "movement"
	ResearchPageName         = "research"
	PlanetlayerPageName      = "planetlayer"
	LogoutPageName           = "logout"
	JumpgatelayerPageName    = "jumpgatelayer"
	FetchResourcesPageName   = "fetchResources"
	FetchTechsName           = "fetchTechs"
	TraderOverviewPageName   = "traderOverview"
	GalaxyPageName           = "galaxy"
	AlliancePageName         = "alliance"
	PremiumPageName          = "premium"
	ShopPageName             = "shop"
	RewardsPageName          = "rewards"
	HighscorePageName        = "highscore"
	BuddiesPageName          = "buddies"
	MessagesPageName         = "messages"
	ChatPageName             = "chat"

	// ajax pages
	FetchEventboxAjaxPageName      = "fetchEventbox"
	FetchResourcesAjaxPageName     = "fetchResources"
	GalaxyContentAjaxPageName      = "galaxyContent"
	EventListAjaxPageName          = "eventList"
	AjaxChatAjaxPageName           = "ajaxChat"
	NoticesAjaxPageName            = "notices"
	RepairlayerAjaxPageName        = "repairlayer"
	TechtreeAjaxPageName           = "techtree"
	PhalanxAjaxPageName            = "phalanx"
	ShareReportOverlayAjaxPageName = "shareReportOverlay"
	JumpgatelayerAjaxPageName      = "jumpgatelayer"
	FederationlayerAjaxPageName    = "federationlayer"
	UnionchangeAjaxPageName        = "unionchange"
	ChangenickAjaxPageName         = "changenick"
	PlanetlayerAjaxPageName        = "planetlayer"
	TraderlayerAjaxPageName        = "traderlayer"
	PlanetRenameAjaxPageName       = "planetRename"
	RightmenuAjaxPageName          = "rightmenu"
	AllianceOverviewAjaxPageName   = "allianceOverview"
	SupportAjaxPageName            = "support"
	BuffActivationAjaxPageName     = "buffActivation"
	AuctioneerAjaxPageName         = "auctioneer"
	HighscoreContentAjaxPageName   = "highscoreContent"
)

func (b *OGame) getPage(page string, celestialID CelestialID, opts ...Option) ([]byte, error) {
	vals := url.Values{"page": {"ingame"}, "component": {page}}
	if page == FetchResourcesPageName || page == FetchTechsName {
		vals = url.Values{"page": {page}}
	}
	if celestialID != 0 {
		vals.Add("cp", FI64(celestialID))
	}
	return b.getPageContent(vals, opts...)
}
