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

	FetchTechsName         = "fetchTechs"
	FetchResourcesPageName = "fetchResources"

	// ajax pages
	RocketlayerPageName            = "rocketlayer"
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

func (b *OGame) getPage(page string, opts ...Option) ([]byte, error) {
	vals := url.Values{"page": {"ingame"}, "component": {page}}
	if page == FetchResourcesPageName || page == FetchTechsName {
		vals = url.Values{"page": {page}}
	}
	return b.getPageContent(vals, opts...)
}

func getPage[T FullPagePages](b *OGame, opts ...Option) (T, error) {
	var zero T
	var pageName string
	switch any(zero).(type) {
	case OverviewPage:
		pageName = OverviewPageName
	case SuppliesPage:
		pageName = SuppliesPageName
	case DefensesPage:
		pageName = DefensesPageName
	case ResearchPage:
		pageName = ResearchPageName
	case ShipyardPage:
		pageName = ShipyardPageName
	case ResourcesSettingsPage:
		pageName = ResourceSettingsPageName
	case FacilitiesPage:
		pageName = FacilitiesPageName
	case MovementPage:
		pageName = MovementPageName
	case PreferencesPage:
		pageName = PreferencesPageName
	default:
		panic("not implemented")
	}
	pageHTML, err := b.getPage(pageName, opts...)
	if err != nil {
		return zero, err
	}
	return ParsePage[T](b, pageHTML)
}

func getAjaxPage[T AjaxPagePages](b *OGame, vals url.Values, opts ...Option) (T, error) {
	var zero T
	switch any(zero).(type) {
	case EventListAjaxPage:
	case MissileAttackLayerAjaxPage:
	case FetchTechsAjaxPage:
	case RocketlayerAjaxPage:
	case PhalanxAjaxPage:
	default:
		panic("not implemented")
	}
	pageHTML, err := b.getPageContent(vals, opts...)
	if err != nil {
		return zero, err
	}
	return ParseAjaxPage[T](b, pageHTML)
}
