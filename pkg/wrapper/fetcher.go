package wrapper

import (
	"net/url"

	"github.com/alaingilbert/ogame/pkg/parser"
)

// Page names
const (
	OverviewPageName         = "overview"
	PreferencesPageName      = "preferences"
	ResourceSettingsPageName = "resourceSettings"
	DefensesPageName         = "defenses"
	LfbuildingsPageName      = "lfbuildings"
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

func getPage[T parser.FullPagePages](b *OGame, opts ...Option) (T, error) {
	var zero T
	var pageName string
	switch any(zero).(type) {
	case parser.OverviewPage:
		pageName = OverviewPageName
	case parser.SuppliesPage:
		pageName = SuppliesPageName
	case parser.DefensesPage:
		pageName = DefensesPageName
	case parser.ResearchPage:
		pageName = ResearchPageName
	case parser.LfBuildingsPage:
		pageName = LfbuildingsPageName
	case parser.ShipyardPage:
		pageName = ShipyardPageName
	case parser.ResourcesSettingsPage:
		pageName = ResourceSettingsPageName
	case parser.FacilitiesPage:
		pageName = FacilitiesPageName
	case parser.MovementPage:
		pageName = MovementPageName
	case parser.PreferencesPage:
		pageName = PreferencesPageName
	default:
		panic("not implemented")
	}
	pageHTML, err := b.getPage(pageName, opts...)
	if err != nil {
		return zero, err
	}
	return parser.ParsePage[T](b.extractor, pageHTML)
}

func getAjaxPage[T parser.AjaxPagePages](b *OGame, vals url.Values, opts ...Option) (T, error) {
	var zero T
	switch any(zero).(type) {
	case parser.EventListAjaxPage:
	case parser.MissileAttackLayerAjaxPage:
	case parser.FetchTechsAjaxPage:
	case parser.RocketlayerAjaxPage:
	case parser.PhalanxAjaxPage:
	default:
		panic("not implemented")
	}
	pageHTML, err := b.getPageContent(vals, opts...)
	if err != nil {
		return zero, err
	}
	return parser.ParseAjaxPage[T](b.extractor, pageHTML)
}