package wrapper

import (
	"crypto/tls"
	"github.com/alaingilbert/ogame/pkg/device"
	"github.com/alaingilbert/ogame/pkg/gameforge"
	"net/http"
	"net/url"
	"time"

	"github.com/alaingilbert/ogame/pkg/extractor"
	"github.com/alaingilbert/ogame/pkg/httpclient"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/taskRunner"
)

// Celestial superset of ogame.Celestial.
// Add methods that can be called for a planet or moon.
type Celestial interface {
	ogame.Celestial
	ActivateItem(string) error
	Build(id ogame.ID, nbr int64) error
	BuildBuilding(buildingID ogame.ID) error
	BuildDefense(defenseID ogame.ID, nbr int64) error
	BuildTechnology(technologyID ogame.ID) error
	CancelBuilding() error
	CancelLfBuilding() error
	CancelResearch() error
	ConstructionsBeingBuilt() (ogame.Constructions, error)
	EnsureFleet(ogame.ShipsInfos, ogame.Speed, ogame.Coordinate, ogame.MissionID, ogame.Resources, int64, int64) (ogame.Fleet, error)
	GetDefense(...Option) (ogame.DefensesInfos, error)
	GetFacilities(...Option) (ogame.Facilities, error)
	GetItems() ([]ogame.Item, error)
	GetLfBuildings(...Option) (ogame.LfBuildings, error)
	GetLfResearch(...Option) (ogame.LfResearches, error)
	GetProduction() ([]ogame.Quantifiable, int64, error)
	GetResources() (ogame.Resources, error)
	GetResourcesBuildings(...Option) (ogame.ResourcesBuildings, error)
	GetResourcesDetails() (ogame.ResourcesDetails, error)
	GetShips(...Option) (ogame.ShipsInfos, error)
	GetTechs() (ogame.ResourcesBuildings, ogame.Facilities, ogame.ShipsInfos, ogame.DefensesInfos, ogame.Researches, ogame.LfBuildings, ogame.LfResearches, error)
	SendFleet(ogame.ShipsInfos, ogame.Speed, ogame.Coordinate, ogame.MissionID, ogame.Resources, int64, int64) (ogame.Fleet, error)
	TearDown(buildingID ogame.ID) error
}

// Prioritizable list of all actions that needs to communicate with ogame server.
// These actions can also be prioritized.
type Prioritizable interface {
	Abandon(IntoPlanet) error
	ActivateItem(string, ogame.CelestialID) error
	Begin() Prioritizable
	BeginNamed(name string) Prioritizable
	BuyMarketplace(itemID int64, celestialID ogame.CelestialID) error
	BuyOfferOfTheDay() error
	BuyResetTree(planetID ogame.PlanetID, tier int64) error
	CancelFleet(ogame.FleetID) error
	CheckTarget(ogame.ShipsInfos, ogame.Coordinate, ...Option) (CheckTargetResponse, error)
	CollectAllMarketplaceMessages() error
	CollectMarketplaceMessage(ogame.MarketplaceMessage) error
	CreateUnion(fleet ogame.Fleet, unionUsers []string) (int64, error)
	DeleteAllMessagesFromTab(tabID ogame.MessagesTabID) error
	DeleteMessage(msgID int64) error
	DoAuction(bid map[ogame.CelestialID]ogame.Resources) error
	Done()
	FastFlightTime(origin, destination ogame.Coordinate, speed ogame.Speed, ships ogame.ShipsInfos, mission ogame.MissionID, holdingTime int64) (secs, fuel int64)
	FlightTime(origin, destination ogame.Coordinate, speed ogame.Speed, ships ogame.ShipsInfos, mission ogame.MissionID, holdingTime int64) (secs, fuel int64)
	FreeResetTree(planetID ogame.PlanetID, tier int64) error
	GalaxyInfos(galaxy, system int64, opts ...Option) (ogame.SystemInfos, error)
	GetActiveItems(ogame.CelestialID) ([]ogame.ActiveItem, error)
	GetAllResources() (map[ogame.CelestialID]ogame.Resources, error)
	GetAttacks(...Option) ([]ogame.AttackEvent, error)
	GetAuction() (ogame.Auction, error)
	GetAvailableDiscoveries(...Option) (int64, error)
	GetCachedAllianceClass() (ogame.AllianceClass, error)
	GetCachedLfBonuses() (ogame.LfBonuses, error)
	GetCachedResearch() ogame.Researches
	GetCelestial(IntoCelestial) (Celestial, error)
	GetCelestials() ([]Celestial, error)
	GetCombatReportSummaryForFleet(ogame.FleetID) (ogame.CombatReportSummary, error)
	GetCombatReportSummaryFor(ogame.Coordinate) (ogame.CombatReportSummary, error)
	GetDMCosts(ogame.CelestialID) (ogame.DMCosts, error)
	GetEmpire(ogame.CelestialType) ([]ogame.EmpireCelestial, error)
	GetEmpireJSON(ogame.CelestialType) (any, error)
	GetEspionageReport(msgID int64) (ogame.EspionageReport, error)
	GetEspionageReportFor(ogame.Coordinate) (ogame.EspionageReport, error)
	GetEspionageReportMessages(maxPage int64) ([]ogame.EspionageReportSummary, error)
	GetExpeditionMessageAt(time.Time) (ogame.ExpeditionMessage, error)
	GetExpeditionMessages(maxPage int64) ([]ogame.ExpeditionMessage, error)
	GetFleetDispatch(ogame.CelestialID, ...Option) (ogame.FleetDispatchInfos, error)
	GetFleets(...Option) ([]ogame.Fleet, ogame.Slots, error)
	GetFleetsFromEventList() ([]ogame.Fleet, error)
	GetItems(ogame.CelestialID) ([]ogame.Item, error)
	GetLfBonuses() (ogame.LfBonuses, error)
	GetMoon(IntoMoon) (Moon, error)
	GetMoons() ([]Moon, error)
	GetPageContent(url.Values) ([]byte, error)
	GetPlanet(IntoPlanet) (Planet, error)
	GetPlanets() ([]Planet, error)
	GetPositionsAvailableForDiscoveryFleet(galaxy int64, system int64, opts ...Option) ([]ogame.Coordinate, error)
	GetResearch() (ogame.Researches, error)
	GetSlots() (ogame.Slots, error)
	GetUserInfos() (ogame.UserInfos, error)
	HeadersForPage(url string) (http.Header, error)
	Highscore(category, typ, page int64) (ogame.Highscore, error)
	IsUnderAttack(opts ...Option) (bool, error)
	Login() error
	LoginWithBearerToken(token string) (bool, bool, error)
	LoginWithExistingCookies() (bool, bool, error)
	Logout() error
	OfferBuyMarketplace(itemID any, quantity, priceType, price, priceRange int64, celestialID ogame.CelestialID) error
	OfferSellMarketplace(itemID any, quantity, priceType, price, priceRange int64, celestialID ogame.CelestialID) error
	PostPageContent(url.Values, url.Values) ([]byte, error)
	RecruitOfficer(typ, days int64) error
	SelectLfResearchArtifacts(planetID ogame.PlanetID, slotNumber int64, techID ogame.ID) error
	SelectLfResearchRandom(planetID ogame.PlanetID, slotNumber int64) error
	SelectLfResearchSelect(planetID ogame.PlanetID, slotNumber int64) error
	SendMessage(playerID int64, message string) error
	SendMessageAlliance(associationID int64, message string) error
	ServerTime() (time.Time, error)
	SetInitiator(initiator string) Prioritizable
	SetPreferences(ogame.Preferences) error
	SetPreferencesLang(lang string) error
	SetVacationMode() error
	Tx(clb func(tx Prioritizable) error) error
	TxNamed(name string, clb func(Prioritizable) error) error
	UseDM(ogame.DMType, ogame.CelestialID) error

	// Planet or Moon functions
	Build(celestialID ogame.CelestialID, id ogame.ID, nbr int64) error
	BuildBuilding(celestialID ogame.CelestialID, buildingID ogame.ID) error
	BuildCancelable(ogame.CelestialID, ogame.ID) error
	BuildDefense(celestialID ogame.CelestialID, defenseID ogame.ID, nbr int64) error
	BuildProduction(celestialID ogame.CelestialID, id ogame.ID, nbr int64) error
	BuildShips(celestialID ogame.CelestialID, shipID ogame.ID, nbr int64) error
	BuildTechnology(celestialID ogame.CelestialID, technologyID ogame.ID) error
	CancelBuilding(ogame.CelestialID) error
	CancelLfBuilding(ogame.CelestialID) error
	CancelResearch(ogame.CelestialID) error
	ConstructionsBeingBuilt(ogame.CelestialID) (ogame.Constructions, error)
	EnsureFleet(celestialID ogame.CelestialID, ships ogame.ShipsInfos, speed ogame.Speed, where ogame.Coordinate, mission ogame.MissionID, resources ogame.Resources, holdingTime, unionID int64) (ogame.Fleet, error)
	FastMiniFleetSpy(coordinate ogame.Coordinate, nbShips int64, opts ...Option) (ogame.MinifleetResponse, error)
	GetDefense(ogame.CelestialID, ...Option) (ogame.DefensesInfos, error)
	GetFacilities(ogame.CelestialID, ...Option) (ogame.Facilities, error)
	GetLfBuildings(ogame.CelestialID, ...Option) (ogame.LfBuildings, error)
	GetLfResearch(ogame.CelestialID, ...Option) (ogame.LfResearches, error)
	GetLfResearchDetails(ogame.CelestialID, ...Option) (ogame.LfResearchDetails, error)
	GetProduction(ogame.CelestialID) ([]ogame.Quantifiable, int64, error)
	GetResources(ogame.CelestialID) (ogame.Resources, error)
	GetResourcesBuildings(ogame.CelestialID, ...Option) (ogame.ResourcesBuildings, error)
	GetResourcesDetails(ogame.CelestialID) (ogame.ResourcesDetails, error)
	GetShips(ogame.CelestialID, ...Option) (ogame.ShipsInfos, error)
	GetTechs(celestialID ogame.CelestialID) (ogame.ResourcesBuildings, ogame.Facilities, ogame.ShipsInfos, ogame.DefensesInfos, ogame.Researches, ogame.LfBuildings, ogame.LfResearches, error)
	MiniFleetSpy(coordinate ogame.Coordinate, nbShips int64, opts ...Option) (ogame.Fleet, error)
	SendDiscoveryFleet(ogame.CelestialID, ogame.Coordinate, ...Option) error
	SendDiscoveryFleet2(ogame.CelestialID, ogame.Coordinate, ...Option) (ogame.Fleet, error)
	SendFleet(celestialID ogame.CelestialID, ships ogame.ShipsInfos, speed ogame.Speed, where ogame.Coordinate, mission ogame.MissionID, resources ogame.Resources, holdingTime, unionID int64) (ogame.Fleet, error)
	SendSystemDiscoveryFleet(celestialID ogame.CelestialID, galaxy, system int64, options ...Option) ([]ogame.Coordinate, error)
	TearDown(celestialID ogame.CelestialID, id ogame.ID) error
	TechnologyDetails(celestialID ogame.CelestialID, id ogame.ID) (ogame.TechnologyDetails, error)

	// Planet specific functions
	DestroyRockets(ogame.PlanetID, int64, int64) error
	GetResourceSettings(ogame.PlanetID, ...Option) (ogame.ResourceSettings, error)
	GetResourcesProductions(ogame.PlanetID) (ogame.Resources, error)
	GetResourcesProductionsLight(ogame.ResourcesBuildings, ogame.Researches, ogame.ResourceSettings, ogame.Temperature) ogame.Resources
	SendIPM(ogame.PlanetID, ogame.Coordinate, int64, ogame.ID) (int64, error)
	SetResourceSettings(ogame.PlanetID, ogame.ResourceSettings) error

	// Moon specific functions
	JumpGate(origin, dest ogame.MoonID, ships ogame.ShipsInfos) (bool, int64, error)
	JumpGateDestinations(origin ogame.MoonID) ([]ogame.MoonID, int64, error)
	Phalanx(ogame.MoonID, ogame.Coordinate) ([]ogame.PhalanxFleet, error)
	UnsafePhalanx(ogame.MoonID, ogame.Coordinate) ([]ogame.PhalanxFleet, error)
}

// Compile time checks to ensure type satisfies Prioritizable interface
var _ Prioritizable = (*OGame)(nil)
var _ Prioritizable = (*Prioritize)(nil)

// Compile time checks to ensure type satisfies Wrapper interface
var _ Wrapper = (*OGame)(nil)

// Wrapper all available functions to control ogame bot
type Wrapper interface {
	Prioritizable
	AddAccount(number int, lang string) (*gameforge.AddAccountResponse, error)
	BytesDownloaded() int64
	BytesUploaded() int64
	CharacterClass() ogame.CharacterClass
	ConstructionTime(id ogame.ID, nbr int64, facilities ogame.Facilities) time.Duration
	CountColonies() (int64, int64)
	Disable()
	Distance(origin, destination ogame.Coordinate) int64
	Enable()
	FleetDeutSaveFactor() float64
	GetCachedAllianceClass() (ogame.AllianceClass, error)
	GetCachedCelestial(IntoCelestial) (Celestial, error)
	GetCachedCelestials() []Celestial
	GetCachedMoon(IntoMoon) (Moon, error)
	GetCachedMoons() []Moon
	GetCachedPlanet(IntoPlanet) (Planet, error)
	GetCachedPlanets() []Planet
	GetCachedPlayer() ogame.UserInfos
	GetCachedPreferences() ogame.Preferences
	GetCachedToken() string
	GetClient() *httpclient.Client
	GetDevice() *device.Device
	GetExtractor() extractor.Extractor
	GetLanguage() string
	GetNbSystems() int64
	GetPublicIP() (string, error)
	GetResearchSpeed() int64
	GetServer() gameforge.Server
	GetServerData() ServerData
	GetSession() string
	GetState() (bool, string)
	GetTasks() taskRunner.TasksOverview
	GetUniverseName() string
	GetUniverseSpeed() int64
	GetUniverseSpeedFleet() int64
	GetUsername() string
	IsConnected() bool
	IsDonutGalaxy() bool
	IsDonutSystem() bool
	IsEnabled() bool
	IsLocked() bool
	IsLoggedIn() bool
	IsPioneers() bool
	IsVacationModeEnabled() bool
	Location() *time.Location
	OnStateChange(clb func(locked bool, actor string))
	PlanetID() ogame.CelestialID
	Quiet(bool)
	ReconnectChat() bool
	RegisterAuctioneerCallback(func(any))
	RegisterChatCallback(func(ogame.ChatMsg))
	RegisterHTMLInterceptor(func(method, url string, params, payload url.Values, pageHTML []byte))
	RegisterWSCallback(string, func([]byte))
	RemoveWSCallback(string)
	ServerURL() string
	ServerVersion() string
	SetAllianceClass(ogame.AllianceClass)
	SetClient(*httpclient.Client)
	SetLfBonuses(lfBonuses ogame.LfBonuses)
	SetLoginWrapper(func(func() (bool, bool, error)) error)
	SetOGameCredentials(username, password, otpSecret, bearerToken string)
	SetProxy(proxyAddress, username, password, proxyType string, loginOnly bool, config *tls.Config) error
	SetResearches(ogame.Researches)
	SoftLogout()
	SystemDistance(system1, system2 int64) int64
	ValidateAccount(code string) error
	WithPriority(priority taskRunner.Priority) Prioritizable
}
