package ogame

import (
	"net/http"
	"net/url"
	"time"
)

// Prioritizable ...
type Prioritizable interface {
	Abandon(interface{}) error
	ActivateItem(string, CelestialID) error
	Begin() Prioritizable
	BeginNamed(name string) Prioritizable
	BuyMarketplace(itemID int64, celestialID CelestialID) error
	BuyOfferOfTheDay() error
	CancelFleet(FleetID) error
	CollectAllMarketplaceMessages() error
	CollectMarketplaceMessage(MarketplaceMessage) error
	CreateUnion(fleet Fleet, unionUsers []string) (int64, error)
	DoAuction(bid map[CelestialID]Resources) error
	Done()
	DeleteAllMessagesFromTab(tabID int64) error
	DeleteMessage(msgID int64) error
	FlightTime(origin, destination Coordinate, speed Speed, ships ShipsInfos) (secs, fuel int64)
	GalaxyInfos(galaxy, system int64, opts ...Option) (SystemInfos, error)
	GetAllResources() (map[CelestialID]Resources, error)
	GetAttacks(...Option) ([]AttackEvent, error)
	GetAuction() (Auction, error)
	GetCachedResearch() Researches
	GetCelestial(interface{}) (Celestial, error)
	GetCelestials() ([]Celestial, error)
	GetCombatReportSummaryFor(Coordinate) (CombatReportSummary, error)
	GetDMCosts(CelestialID) (DMCosts, error)
	GetEmpire(nbr int64) (interface{}, error)
	GetEspionageReport(msgID int64) (EspionageReport, error)
	GetEspionageReportFor(Coordinate) (EspionageReport, error)
	GetEspionageReportMessages() ([]EspionageReportSummary, error)
	GetExpeditionMessageAt(time.Time) (ExpeditionMessage, error)
	GetExpeditionMessages() ([]ExpeditionMessage, error)
	GetFleets(...Option) ([]Fleet, Slots)
	GetFleetsFromEventList() []Fleet
	GetItems(CelestialID) ([]Item, error)
	GetMoon(interface{}) (Moon, error)
	GetMoons() []Moon
	GetPageContent(url.Values) ([]byte, error)
	GetPlanet(interface{}) (Planet, error)
	GetPlanets() []Planet
	GetResearch() Researches
	GetSlots() Slots
	GetUserInfos() UserInfos
	HeadersForPage(url string) (http.Header, error)
	Highscore(category, typ, page int64) (Highscore, error)
	IsUnderAttack() (bool, error)
	Login() error
	LoginWithExistingCookies() (bool, error)
	Logout()
	OfferBuyMarketplace(itemID interface{}, quantity, priceType, price, priceRange int64, celestialID CelestialID) error
	OfferSellMarketplace(itemID interface{}, quantity, priceType, price, priceRange int64, celestialID CelestialID) error
	PostPageContent(url.Values, url.Values) ([]byte, error)
	SendMessage(playerID int64, message string) error
	SendMessageAlliance(associationID int64, message string) error
	ServerTime() time.Time
	Tx(clb func(tx Prioritizable) error) error
	UseDM(string, CelestialID) error

	// Planet or Moon functions
	Build(celestialID CelestialID, id ID, nbr int64) error
	BuildBuilding(celestialID CelestialID, buildingID ID) error
	BuildCancelable(CelestialID, ID) error
	BuildDefense(celestialID CelestialID, defenseID ID, nbr int64) error
	BuildProduction(celestialID CelestialID, id ID, nbr int64) error
	BuildShips(celestialID CelestialID, shipID ID, nbr int64) error
	BuildTechnology(celestialID CelestialID, technologyID ID) error
	CancelBuilding(CelestialID) error
	CancelResearch(CelestialID) error
	ConstructionsBeingBuilt(CelestialID) (buildingID ID, buildingCountdown int64, researchID ID, researchCountdown int64)
	EnsureFleet(celestialID CelestialID, ships []Quantifiable, speed Speed, where Coordinate, mission MissionID, resources Resources, expeditiontime, unionID int64) (Fleet, error)
	GetDefense(CelestialID) (DefensesInfos, error)
	GetFacilities(CelestialID) (Facilities, error)
	GetProduction(CelestialID) ([]Quantifiable, int64, error)
	GetResources(CelestialID) (Resources, error)
	GetResourcesBuildings(CelestialID) (ResourcesBuildings, error)
	GetResourcesDetails(CelestialID) (ResourcesDetails, error)
	GetShips(CelestialID) (ShipsInfos, error)
	SendFleet(celestialID CelestialID, ships []Quantifiable, speed Speed, where Coordinate, mission MissionID, resources Resources, expeditiontime, unionID int64) (Fleet, error)
	TearDown(celestialID CelestialID, id ID) error

	// Planet specific functions
	GetResourceSettings(PlanetID) (ResourceSettings, error)
	GetResourcesProductions(PlanetID) (Resources, error)
	GetResourcesProductionsLight(ResourcesBuildings, Researches, ResourceSettings, Temperature) Resources
	SendIPM(PlanetID, Coordinate, int64, ID) (int64, error)
	SetResourceSettings(PlanetID, ResourceSettings) error

	// Moon specific functions
	JumpGate(origin, dest MoonID, ships ShipsInfos) (bool, int64, error)
	JumpGateDestinations(origin MoonID) ([]MoonID, int64, error)
	Phalanx(MoonID, Coordinate) ([]Fleet, error)
	UnsafePhalanx(MoonID, Coordinate) ([]Fleet, error)
}

// Wrapper all available functions to control ogame bot
type Wrapper interface {
	Prioritizable
	AddAccount(number int, lang string) (NewAccount, error)
	BytesDownloaded() int64
	BytesUploaded() int64
	CharacterClass() CharacterClass
	Disable()
	Distance(origin, destination Coordinate) int64
	Enable()
	FleetDeutSaveFactor() float64
	GetAlliancePageContent(url.Values) ([]byte, error)
	GetCachedCelestial(interface{}) Celestial
	GetCachedCelestials() []Celestial
	GetCachedMoons() []Moon
	GetCachedPlanets() []Planet
	GetCachedPlayer() UserInfos
	GetCachedPreferences() Preferences
	GetClient() *OGameClient
	GetExtractor() Extractor
	GetLanguage() string
	GetNbSystems() int64
	GetPublicIP() (string, error)
	GetResearchSpeed() int64
	GetServer() Server
	GetServerData() ServerData
	GetSession() string
	GetState() (bool, string)
	GetTasks() TasksOverview
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
	IsVacationModeEnabled() bool
	IsV7() bool
	Location() *time.Location
	OnStateChange(clb func(locked bool, actor string))
	Quiet(bool)
	ReconnectChat() bool
	RegisterAuctioneerCallback(func([]byte))
	RegisterChatCallback(func(ChatMsg))
	RegisterHTMLInterceptor(func(method, url string, params, payload url.Values, pageHTML []byte))
	RegisterWSCallback(string, func([]byte))
	RemoveWSCallback(string)
	ServerURL() string
	ServerVersion() string
	SetLoginWrapper(func(func() (bool, error)) error)
	SetOGameCredentials(username, password string)
	SetProxy(proxyAddress, username, password, proxyType string, loginOnly bool) error
	SetUserAgent(newUserAgent string)
	WithPriority(priority int) Prioritizable
}

// BaseOgameObj base interface for all ogame objects (buildings, technologies, ships, defenses)
type BaseOgameObj interface {
	ConstructionTime(nbr, universeSpeed int64, facilities Facilities, hasTechnocrat, isDiscoverer bool) time.Duration
	GetID() ID
	GetName() string
	GetPrice(int64) Resources
	GetRequirements() map[ID]int64
	IsAvailable(CelestialType, LazyResourcesBuildings, LazyFacilities, LazyResearches, int64) bool
}

// Levelable base interface for all levelable ogame objects (buildings, technologies)
type Levelable interface {
	BaseOgameObj
	GetLevel(LazyResourcesBuildings, LazyFacilities, LazyResearches) int64
}

// Technology interface that all technologies implement
type Technology interface {
	Levelable
}

// Building interface that all buildings implement
type Building interface {
	Levelable
	DeconstructionPrice(lvl int64, techs Researches) Resources
}

// DefenderObj base interface for all defensive units (ships, defenses)
type DefenderObj interface {
	BaseOgameObj
	GetStructuralIntegrity(Researches) int64
	GetShieldPower(Researches) int64
	GetWeaponPower(Researches) int64
	GetRapidfireFrom() map[ID]int64
	GetRapidfireAgainst() map[ID]int64
}

// Ship interface implemented by all ships units
type Ship interface {
	DefenderObj
	GetCargoCapacity(techs Researches, probeRaids, isCollector bool) int64
	GetSpeed(techs Researches, isCollector, isGeneral bool) int64
	GetFuelConsumption(techs Researches, fleetDeutSaveFactor float64, isGeneral bool) int64
}

// Defense interface implemented by all defenses units
type Defense interface {
	DefenderObj
}

// Celestial ...
type Celestial interface {
	GetID() CelestialID
	GetType() CelestialType
	GetName() string
	GetDiameter() int64
	GetCoordinate() Coordinate
	GetFields() Fields
	GetResources() (Resources, error)
	GetResourcesDetails() (ResourcesDetails, error)
	GetFacilities() (Facilities, error)
	SendFleet([]Quantifiable, Speed, Coordinate, MissionID, Resources, int64, int64) (Fleet, error)
	EnsureFleet([]Quantifiable, Speed, Coordinate, MissionID, Resources, int64, int64) (Fleet, error)
	GetDefense() (DefensesInfos, error)
	GetShips() (ShipsInfos, error)
	BuildDefense(defenseID ID, nbr int64) error
	ConstructionsBeingBuilt() (ID, int64, ID, int64)
	GetProduction() ([]Quantifiable, int64, error)
	GetResourcesBuildings() (ResourcesBuildings, error)
	Build(id ID, nbr int64) error
	BuildBuilding(buildingID ID) error
	BuildTechnology(technologyID ID) error
	CancelResearch() error
	CancelBuilding() error
	TearDown(buildingID ID) error
	GetItems() ([]Item, error)
	ActivateItem(string) error
}
