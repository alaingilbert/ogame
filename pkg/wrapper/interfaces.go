package wrapper

import (
	"crypto/tls"
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
	ConstructionsBeingBuilt() (ogame.ID, int64, ogame.ID, int64, ogame.ID, int64)
	EnsureFleet([]ogame.Quantifiable, ogame.Speed, ogame.Coordinate, ogame.MissionID, ogame.Resources, int64, int64) (ogame.Fleet, error)
	GetDefense(...Option) (ogame.DefensesInfos, error)
	GetFacilities(...Option) (ogame.Facilities, error)
	GetItems() ([]ogame.Item, error)
	GetLfBuildings(...Option) (ogame.LfBuildings, error)
	GetProduction() ([]ogame.Quantifiable, int64, error)
	GetResources() (ogame.Resources, error)
	GetResourcesBuildings(...Option) (ogame.ResourcesBuildings, error)
	GetResourcesDetails() (ogame.ResourcesDetails, error)
	GetShips(...Option) (ogame.ShipsInfos, error)
	GetTechs() (ogame.ResourcesBuildings, ogame.Facilities, ogame.ShipsInfos, ogame.DefensesInfos, ogame.Researches, ogame.LfBuildings, error)
	SendFleet([]ogame.Quantifiable, ogame.Speed, ogame.Coordinate, ogame.MissionID, ogame.Resources, int64, int64) (ogame.Fleet, error)
	TearDown(buildingID ogame.ID) error
}

// Prioritizable list of all actions that needs to communicate with ogame server.
// These actions can also be prioritized.
type Prioritizable interface {
	Abandon(any) error
	ActivateItem(string, ogame.CelestialID) error
	Begin() Prioritizable
	BeginNamed(name string) Prioritizable
	BuyMarketplace(itemID int64, celestialID ogame.CelestialID) error
	BuyOfferOfTheDay() error
	CancelFleet(ogame.FleetID) error
	CollectAllMarketplaceMessages() error
	CollectMarketplaceMessage(ogame.MarketplaceMessage) error
	CreateUnion(fleet ogame.Fleet, unionUsers []string) (int64, error)
	DeleteAllMessagesFromTab(tabID ogame.MessagesTabID) error
	DeleteMessage(msgID int64) error
	DoAuction(bid map[ogame.CelestialID]ogame.Resources) error
	Done()
	FlightTime(origin, destination ogame.Coordinate, speed ogame.Speed, ships ogame.ShipsInfos, mission ogame.MissionID) (secs, fuel int64)
	GalaxyInfos(galaxy, system int64, opts ...Option) (ogame.SystemInfos, error)
	GetActiveItems(ogame.CelestialID) ([]ogame.ActiveItem, error)
	GetAllResources() (map[ogame.CelestialID]ogame.Resources, error)
	GetAttacks(...Option) ([]ogame.AttackEvent, error)
	GetAuction() (ogame.Auction, error)
	GetCachedResearch() ogame.Researches
	GetCelestial(any) (Celestial, error)
	GetCelestials() ([]Celestial, error)
	GetCombatReportSummaryFor(ogame.Coordinate) (ogame.CombatReportSummary, error)
	GetDMCosts(ogame.CelestialID) (ogame.DMCosts, error)
	GetEmpire(ogame.CelestialType) ([]ogame.EmpireCelestial, error)
	GetEmpireJSON(nbr int64) (any, error)
	GetEspionageReport(msgID int64) (ogame.EspionageReport, error)
	GetEspionageReportFor(ogame.Coordinate) (ogame.EspionageReport, error)
	GetEspionageReportMessages() ([]ogame.EspionageReportSummary, error)
	GetExpeditionMessageAt(time.Time) (ogame.ExpeditionMessage, error)
	GetExpeditionMessages() ([]ogame.ExpeditionMessage, error)
	GetFleets(...Option) ([]ogame.Fleet, ogame.Slots)
	GetFleetsFromEventList() []ogame.Fleet
	GetItems(ogame.CelestialID) ([]ogame.Item, error)
	GetMoon(any) (Moon, error)
	GetMoons() []Moon
	GetPageContent(url.Values) ([]byte, error)
	GetPlanet(any) (Planet, error)
	GetPlanets() []Planet
	GetResearch() ogame.Researches
	GetSlots() ogame.Slots
	GetUserInfos() ogame.UserInfos
	HeadersForPage(url string) (http.Header, error)
	Highscore(category, typ, page int64) (ogame.Highscore, error)
	IsUnderAttack() (bool, error)
	Login() error
	LoginWithBearerToken(token string) (bool, error)
	LoginWithExistingCookies() (bool, error)
	Logout()
	OfferBuyMarketplace(itemID any, quantity, priceType, price, priceRange int64, celestialID ogame.CelestialID) error
	OfferSellMarketplace(itemID any, quantity, priceType, price, priceRange int64, celestialID ogame.CelestialID) error
	PostPageContent(url.Values, url.Values) ([]byte, error)
	RecruitOfficer(typ, days int64) error
	SendMessage(playerID int64, message string) error
	SendMessageAlliance(associationID int64, message string) error
	ServerTime() time.Time
	SetInitiator(initiator string) Prioritizable
	SetVacationMode() error
	Tx(clb func(tx Prioritizable) error) error
	UseDM(string, ogame.CelestialID) error

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
	ConstructionsBeingBuilt(ogame.CelestialID) (buildingID ogame.ID, buildingCountdown int64, researchID ogame.ID, researchCountdown int64, lfBuildingID ogame.ID, lfBuildingCountdown int64)
	EnsureFleet(celestialID ogame.CelestialID, ships []ogame.Quantifiable, speed ogame.Speed, where ogame.Coordinate, mission ogame.MissionID, resources ogame.Resources, holdingTime, unionID int64) (ogame.Fleet, error)
	GetDefense(ogame.CelestialID, ...Option) (ogame.DefensesInfos, error)
	GetFacilities(ogame.CelestialID, ...Option) (ogame.Facilities, error)
	GetLfBuildings(ogame.CelestialID, ...Option) (ogame.LfBuildings, error)
	GetProduction(ogame.CelestialID) ([]ogame.Quantifiable, int64, error)
	GetResources(ogame.CelestialID) (ogame.Resources, error)
	GetResourcesBuildings(ogame.CelestialID, ...Option) (ogame.ResourcesBuildings, error)
	GetResourcesDetails(ogame.CelestialID) (ogame.ResourcesDetails, error)
	GetShips(ogame.CelestialID, ...Option) (ogame.ShipsInfos, error)
	GetTechs(celestialID ogame.CelestialID) (ogame.ResourcesBuildings, ogame.Facilities, ogame.ShipsInfos, ogame.DefensesInfos, ogame.Researches, ogame.LfBuildings, error)
	SendFleet(celestialID ogame.CelestialID, ships []ogame.Quantifiable, speed ogame.Speed, where ogame.Coordinate, mission ogame.MissionID, resources ogame.Resources, holdingTime, unionID int64) (ogame.Fleet, error)
	TearDown(celestialID ogame.CelestialID, id ogame.ID) error

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
	Phalanx(ogame.MoonID, ogame.Coordinate) ([]ogame.Fleet, error)
	UnsafePhalanx(ogame.MoonID, ogame.Coordinate) ([]ogame.Fleet, error)
}

// Wrapper all available functions to control ogame bot
type Wrapper interface {
	Prioritizable
	AddAccount(number int, lang string) (*AddAccountRes, error)
	BytesDownloaded() int64
	BytesUploaded() int64
	CharacterClass() ogame.CharacterClass
	ConstructionTime(id ogame.ID, nbr int64, facilities ogame.Facilities) time.Duration
	Disable()
	Distance(origin, destination ogame.Coordinate) int64
	Enable()
	FleetDeutSaveFactor() float64
	GetCachedCelestial(any) Celestial
	GetCachedCelestials() []Celestial
	GetCachedMoons() []Moon
	GetCachedPlanets() []Planet
	GetCachedPlayer() ogame.UserInfos
	GetCachedPreferences() ogame.Preferences
	GetClient() *httpclient.Client
	GetExtractor() extractor.Extractor
	GetLanguage() string
	GetNbSystems() int64
	GetPublicIP() (string, error)
	GetResearchSpeed() int64
	GetServer() Server
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
	IsV7() bool
	IsV9() bool
	IsVacationModeEnabled() bool
	Location() *time.Location
	OnStateChange(clb func(locked bool, actor string))
	Quiet(bool)
	ReconnectChat() bool
	RegisterAuctioneerCallback(func(any))
	RegisterChatCallback(func(ogame.ChatMsg))
	RegisterHTMLInterceptor(func(method, url string, params, payload url.Values, pageHTML []byte))
	RegisterWSCallback(string, func([]byte))
	RemoveWSCallback(string)
	ServerURL() string
	ServerVersion() string
	SetClient(*httpclient.Client)
	SetGetServerDataWrapper(func(func() (ServerData, error)) (ServerData, error))
	SetLoginWrapper(func(func() (bool, error)) error)
	SetOGameCredentials(username, password, otpSecret, bearerToken string)
	SetProxy(proxyAddress, username, password, proxyType string, loginOnly bool, config *tls.Config) error
	SetUserAgent(newUserAgent string)
	ValidateAccount(code string) error
	WithPriority(priority taskRunner.Priority) Prioritizable
}
