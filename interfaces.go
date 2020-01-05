package ogame

import (
	"net/http"
	"net/url"
	"time"
)

// Wrapper all available functions to control ogame bot
type Wrapper interface {
	IsV7() bool
	GetExtractor() Extractor
	SetOGameCredentials(username, password string)
	SetLoginProxy(proxy, username, password string) error
	SetProxy(proxy, username, password string) error
	SetSocks5Proxy(proxy, username, password string) error
	SetLoginWrapper(func(func() error) error)
	GetClient() *OGameClient
	Enable()
	Disable()
	IsEnabled() bool
	Quiet(bool)
	Tx(clb func(tx *Prioritize) error) error
	Begin() *Prioritize
	BeginNamed(name string) *Prioritize
	WithPriority(priority int) *Prioritize
	GetPublicIP() (string, error)
	OnStateChange(clb func(locked bool, actor string))
	GetState() (bool, string)
	IsLocked() bool
	GetSession() string
	AddAccount(number int, lang string) (NewAccount, error)
	GetServer() Server
	GetServerData() ServerData
	SetUserAgent(newUserAgent string)
	ServerURL() string
	GetLanguage() string
	GetPageContent(url.Values) []byte
	GetAlliancePageContent(url.Values) []byte
	PostPageContent(url.Values, url.Values) []byte
	Login() error
	Logout()
	IsLoggedIn() bool
	IsConnected() bool
	GetUsername() string
	GetUniverseName() string
	GetUniverseSpeed() int64
	GetUniverseSpeedFleet() int64
	GetResearchSpeed() int64
	SetResearchSpeed(int64)
	GetNbSystems() int64
	SetNbSystems(int64)
	IsDonutGalaxy() bool
	IsDonutSystem() bool
	FleetDeutSaveFactor() float64
	ServerVersion() string
	ServerTime() time.Time
	Location() *time.Location
	IsUnderAttack() (bool, error)
	GetUserInfos() UserInfos
	SendMessage(playerID int64, message string) error
	SendMessageAlliance(associationID int64, message string) error
	ReconnectChat() bool
	GetFleets() ([]Fleet, Slots)
	GetFleetsFromEventList() []Fleet
	CancelFleet(FleetID) error
	GetAttacks() ([]AttackEvent, error)
	GetAttacksUsing(CelestialID) ([]AttackEvent, error)
	GalaxyInfos(galaxy, system int64) (SystemInfos, error)
	GetCachedResearch() Researches
	GetResearch() Researches
	GetCachedPlanets() []Planet
	GetCachedMoons() []Moon
	GetCachedCelestials() []Celestial
	GetCachedCelestial(interface{}) Celestial
	GetCachedPlayer() UserInfos
	GetCachedPreferences() Preferences
	IsVacationModeEnabled() bool
	GetPlanets() []Planet
	GetPlanet(interface{}) (Planet, error)
	GetMoons() []Moon
	GetMoon(interface{}) (Moon, error)
	GetCelestial(interface{}) (Celestial, error)
	GetCelestials() ([]Celestial, error)
	Abandon(interface{}) error
	GetEspionageReportMessages() ([]EspionageReportSummary, error)
	GetEspionageReportFor(Coordinate) (EspionageReport, error)
	GetEspionageReport(msgID int64) (EspionageReport, error)
	GetCombatReportSummaryFor(Coordinate) (CombatReportSummary, error)
	//GetCombatReport(msgID int) (CombatReport, error)
	DeleteMessage(msgID int64) error
	DeleteAllMessagesFromTab(tabID int64) error
	Distance(origin, destination Coordinate) int64
	FlightTime(origin, destination Coordinate, speed Speed, ships ShipsInfos) (secs, fuel int64)
	RegisterChatCallback(func(ChatMsg))
	RegisterAuctioneerCallback(func([]byte))
	RegisterHTMLInterceptor(func(method, url string, params, payload url.Values, pageHTML []byte))
	GetSlots() Slots
	BuyOfferOfTheDay() error
	BytesDownloaded() int64
	BytesUploaded() int64
	CreateUnion(fleet Fleet) (int64, error)
	GetEmpire(nbr int64) (interface{}, error)
	HeadersForPage(url string) (http.Header, error)

	// Planet or Moon functions
	GetResources(CelestialID) (Resources, error)
	GetResourcesDetails(CelestialID) (ResourcesDetails, error)
	SendFleet(celestialID CelestialID, ships []Quantifiable, speed Speed, where Coordinate, mission MissionID, resources Resources, expeditiontime, unionID int64) (Fleet, error)
	EnsureFleet(celestialID CelestialID, ships []Quantifiable, speed Speed, where Coordinate, mission MissionID, resources Resources, expeditiontime, unionID int64) (Fleet, error)
	Build(celestialID CelestialID, id ID, nbr int64) error
	BuildCancelable(CelestialID, ID) error
	BuildProduction(celestialID CelestialID, id ID, nbr int64) error
	BuildBuilding(celestialID CelestialID, buildingID ID) error
	BuildDefense(celestialID CelestialID, defenseID ID, nbr int64) error
	BuildShips(celestialID CelestialID, shipID ID, nbr int64) error
	CancelBuilding(CelestialID) error
	TearDown(celestialID CelestialID, id ID) error
	ConstructionsBeingBuilt(CelestialID) (buildingID ID, buildingCountdown int64, researchID ID, researchCountdown int64)
	GetProduction(CelestialID) ([]Quantifiable, int64, error)
	GetFacilities(CelestialID) (Facilities, error)
	GetDefense(CelestialID) (DefensesInfos, error)
	GetShips(CelestialID) (ShipsInfos, error)
	GetResourcesBuildings(CelestialID) (ResourcesBuildings, error)
	CancelResearch(CelestialID) error
	BuildTechnology(celestialID CelestialID, technologyID ID) error

	// Planet specific functions
	GetResourceSettings(PlanetID) (ResourceSettings, error)
	SetResourceSettings(PlanetID, ResourceSettings) error
	SendIPM(PlanetID, Coordinate, int64, ID) (int64, error)
	//GetResourcesProductionRatio(PlanetID) (float64, error)
	GetResourcesProductions(PlanetID) (Resources, error)
	GetResourcesProductionsLight(ResourcesBuildings, Researches, ResourceSettings, Temperature) Resources

	// Moon specific functions
	Phalanx(MoonID, Coordinate) ([]Fleet, error)
	UnsafePhalanx(MoonID, Coordinate) ([]Fleet, error)
	JumpGate(origin, dest MoonID, ships ShipsInfos) (bool, int64, error)
}

// BaseOgameObj base interface for all ogame objects (buildings, technologies, ships, defenses)
type BaseOgameObj interface {
	GetID() ID
	GetName() string
	ConstructionTime(nbr, universeSpeed int64, facilities Facilities) time.Duration
	GetRequirements() map[ID]int64
	GetPrice(int64) Resources
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
	GetCargoCapacity(techs Researches, probeRaids bool) int64
	GetSpeed(Researches) int64
	GetFuelConsumption() int64
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
}
