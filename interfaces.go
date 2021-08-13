package ogame

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/alaingilbert/clockwork"
)

// Prioritizable ...
type Prioritizable interface {
	RecruitOfficer(typ, days int64) error
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
	FlightTime(origin, destination Coordinate, speed Speed, ships ShipsInfos, mission MissionID) (secs, fuel int64)
	GalaxyInfos(galaxy, system int64, opts ...Option) (SystemInfos, error)
	GetAlliancePageContent(url.Values) ([]byte, error)
	GetAllResources() (map[CelestialID]Resources, error)
	GetAttacks(...Option) ([]AttackEvent, error)
	GetAuction() (Auction, error)
	GetCachedResearch() Researches
	GetCelestial(interface{}) (Celestial, error)
	GetCelestials() ([]Celestial, error)
	GetCombatReportSummaryFor(Coordinate) (CombatReportSummary, error)
	GetDMCosts(CelestialID) (DMCosts, error)
	GetEmpire(CelestialType) ([]EmpireCelestial, error)
	GetEmpireJSON(nbr int64) (interface{}, error)
	GetEspionageReport(msgID int64) (EspionageReport, error)
	GetEspionageReportFor(Coordinate) (EspionageReport, error)
	GetEspionageReportMessages() ([]EspionageReportSummary, error)
	GetExpeditionMessageAt(time.Time) (ExpeditionMessage, error)
	GetExpeditionMessages() ([]ExpeditionMessage, error)
	GetFleets(...Option) ([]Fleet, Slots)
	GetFleetsFromEventList(...Option) []Fleet
	GetItems(CelestialID) ([]Item, error)
	GetActiveItems(CelestialID) ([]ActiveItem, error)
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
	LoginWithBearerToken(token string) (bool, error)
	LoginWithExistingCookies() (bool, error)
	Logout()
	OfferBuyMarketplace(itemID interface{}, quantity, priceType, price, priceRange int64, celestialID CelestialID) error
	OfferSellMarketplace(itemID interface{}, quantity, priceType, price, priceRange int64, celestialID CelestialID) error
	PostPageContent(url.Values, url.Values) ([]byte, error)
	SendMessage(playerID int64, message string) error
	SendMessageAlliance(associationID int64, message string) error
	ServerTime() time.Time
	SetInitiator(initiator string) Prioritizable
	Tx(clb func(tx Prioritizable) error) error
	UseDM(string, CelestialID) error
	GetCachedData() Data

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
	EnsureFleet(celestialID CelestialID, ships []Quantifiable, speed Speed, where Coordinate, mission MissionID, resources Resources, holdingTime, unionID int64) (Fleet, error)
	GetDefense(CelestialID, ...Option) (DefensesInfos, error)
	GetFacilities(CelestialID, ...Option) (Facilities, error)
	GetProduction(CelestialID) ([]Quantifiable, int64, error)
	GetResources(CelestialID) (Resources, error)
	GetResourcesBuildings(CelestialID, ...Option) (ResourcesBuildings, error)
	GetResourcesDetails(CelestialID) (ResourcesDetails, error)
	GetTechs(celestialID CelestialID) (ResourcesBuildings, Facilities, ShipsInfos, DefensesInfos, Researches, error)
	GetShips(CelestialID, ...Option) (ShipsInfos, error)
	SendFleet(celestialID CelestialID, ships []Quantifiable, speed Speed, where Coordinate, mission MissionID, resources Resources, holdingTime, unionID int64) (Fleet, error)
	TearDown(celestialID CelestialID, id ID) error

	// Planet specific functions
	GetResourceSettings(PlanetID, ...Option) (ResourceSettings, error)
	GetResourcesProductions(PlanetID) (Resources, error)
	GetResourcesProductionsLight(ResourcesBuildings, Researches, ResourceSettings, Temperature) Resources
	DestroyRockets(PlanetID, int64, int64) error
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
	ValidateAccount(code string) error
	AddAccount(number int, lang string) (NewAccount, error)
	BytesDownloaded() int64
	BytesUploaded() int64
	IsPioneers() bool
	CharacterClass() CharacterClass
	Disable()
	Distance(origin, destination Coordinate) int64
	Enable()
	FleetDeutSaveFactor() float64
	GetCachedCelestial(interface{}) Celestial
	GetCachedCelestials() []Celestial
	GetCachedMoons() []Moon
	GetCachedPlanets() []Planet
	GetCachedPlayer() UserInfos
	GetCachedPreferences() Preferences
	GetClient() *OGameClient
	SetClient(*OGameClient)
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
	ConstructionTime(id ID, nbr int64, facilities Facilities) time.Duration
	IsEnabled() bool
	IsLocked() bool
	IsLoggedIn() bool
	IsVacationModeEnabled() bool
	IsV7() bool
	Location() *time.Location
	OnStateChange(clb func(locked bool, actor string))
	Quiet(bool)
	ReconnectChat() bool
	RegisterAuctioneerCallback(func(interface{}))
	RegisterChatCallback(func(ChatMsg))
	RegisterHTMLInterceptor(func(method, url string, params, payload url.Values, pageHTML []byte))
	RegisterWSCallback(string, func([]byte))
	RemoveWSCallback(string)
	ServerURL() string
	ServerVersion() string
	SetLoginWrapper(func(func() (bool, error)) error)
	SetOGameCredentials(username, password, otpSecret, bearerToken string)
	SetProxy(proxyAddress, username, password, proxyType string, loginOnly bool, config *tls.Config) error
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
	GetCargoCapacity(techs Researches, probeRaids, isCollector, isPioneers bool) int64
	GetSpeed(techs Researches, isCollector, isGeneral bool) int64
	GetFuelConsumption(techs Researches, fleetDeutSaveFactor float64, isGeneral bool) int64
	GetFuelCapacity() int64
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
	GetFacilities(...Option) (Facilities, error)
	SendFleet([]Quantifiable, Speed, Coordinate, MissionID, Resources, int64, int64) (Fleet, error)
	EnsureFleet([]Quantifiable, Speed, Coordinate, MissionID, Resources, int64, int64) (Fleet, error)
	GetDefense(...Option) (DefensesInfos, error)
	GetShips(...Option) (ShipsInfos, error)
	BuildDefense(defenseID ID, nbr int64) error
	ConstructionsBeingBuilt() (ID, int64, ID, int64)
	GetProduction() ([]Quantifiable, int64, error)
	GetResourcesBuildings(...Option) (ResourcesBuildings, error)
	Build(id ID, nbr int64) error
	BuildBuilding(buildingID ID) error
	BuildTechnology(technologyID ID) error
	CancelResearch() error
	CancelBuilding() error
	TearDown(buildingID ID) error
	GetItems() ([]Item, error)
	ActivateItem(string) error
}

// Extractor ...
type Extractor interface {
	ExtractIsInVacation(pageHTML []byte) bool
	ExtractPlanets(pageHTML []byte, b *OGame) []Planet
	ExtractPlanet(pageHTML []byte, v interface{}, b *OGame) (Planet, error)
	ExtractPlanetByCoord(pageHTML []byte, b *OGame, coord Coordinate) (Planet, error)
	ExtractPremiumToken(pageHTML []byte, days int64) (token string, err error)
	ExtractMoons(pageHTML []byte, b *OGame) []Moon
	ExtractMoon(pageHTML []byte, b *OGame, v interface{}) (Moon, error)
	ExtractMoonByCoord(pageHTML []byte, b *OGame, coord Coordinate) (Moon, error)
	ExtractCelestials(pageHTML []byte, b *OGame) ([]Celestial, error)
	ExtractCelestial(pageHTML []byte, b *OGame, v interface{}) (Celestial, error)
	ExtractServerTime(pageHTML []byte) (time.Time, error)
	ExtractFleetsFromEventList(pageHTML []byte) []Fleet
	ExtractDestroyRockets(pageHTML []byte) (abm, ipm int64, token string, err error)
	ExtractIPM(pageHTML []byte) (duration, max int64, token string)
	ExtractFleets(pageHTML []byte, location *time.Location) (res []Fleet)
	ExtractSlots(pageHTML []byte) Slots
	ExtractOgameTimestamp(pageHTML []byte) int64
	ExtractResources(pageHTML []byte) Resources
	ExtractResourcesDetailsFromFullPage(pageHTML []byte) ResourcesDetails
	ExtractResourceSettings(pageHTML []byte) (ResourceSettings, error)
	ExtractAttacks(pageHTML []byte) ([]AttackEvent, error)
	ExtractOfferOfTheDay(pageHTML []byte) (int64, string, PlanetResources, Multiplier, error)
	ExtractResourcesBuildings(pageHTML []byte) (ResourcesBuildings, error)
	ExtractExpeditionMessages(pageHTML []byte, location *time.Location) ([]ExpeditionMessage, int64, error)
	ExtractMarketplaceMessages(pageHTML []byte, location *time.Location) ([]MarketplaceMessage, int64, error)
	ExtractDefense(pageHTML []byte) (DefensesInfos, error)
	ExtractShips(pageHTML []byte) (ShipsInfos, error)
	ExtractFacilities(pageHTML []byte) (Facilities, error)
	ExtractResearch(pageHTML []byte) Researches
	ExtractProduction(pageHTML []byte) ([]Quantifiable, int64, error)
	ExtractOverviewProduction(pageHTML []byte) ([]Quantifiable, int64, error)
	ExtractFleet1Ships(pageHTML []byte) ShipsInfos
	ExtractEspionageReportMessageIDs(pageHTML []byte) ([]EspionageReportSummary, int64)
	ExtractCombatReportMessagesSummary(pageHTML []byte) ([]CombatReportSummary, int64)
	ExtractEspionageReport(pageHTML []byte, location *time.Location) (EspionageReport, error)
	ExtractResourcesProductions(pageHTML []byte) (Resources, error)
	ExtractPreferences(pageHTML []byte) Preferences
	ExtractSpioAnz(pageHTML []byte) int64
	ExtractPreferencesShowActivityMinutes(pageHTML []byte) bool
	ExtractHiddenFields(pageHTML []byte) (fields url.Values)
	ExtractHiddenFieldsFromDoc(doc *goquery.Document) url.Values
	ExtractBodyIDFromDoc(doc *goquery.Document) string
	ExtractIsInVacationFromDoc(doc *goquery.Document) bool
	ExtractPlanetsFromDoc(doc *goquery.Document, b *OGame) []Planet
	ExtractPlanetByIDFromDoc(doc *goquery.Document, b *OGame, planetID PlanetID) (Planet, error)
	ExtractCelestialByIDFromDoc(doc *goquery.Document, b *OGame, celestialID CelestialID) (Celestial, error)
	ExtractPlanetByCoordFromDoc(doc *goquery.Document, b *OGame, coord Coordinate) (Planet, error)
	ExtractOgameTimestampFromDoc(doc *goquery.Document) int64
	ExtractResourcesFromDoc(doc *goquery.Document) Resources
	ExtractResourcesDetailsFromFullPageFromDoc(doc *goquery.Document) ResourcesDetails
	ExtractPlanetFromDoc(doc *goquery.Document, v interface{}, b *OGame) (Planet, error)
	ExtractMoonsFromDoc(doc *goquery.Document, b *OGame) []Moon
	ExtractMoonFromDoc(doc *goquery.Document, b *OGame, v interface{}) (Moon, error)
	ExtractMoonByCoordFromDoc(doc *goquery.Document, b *OGame, coord Coordinate) (Moon, error)
	ExtractMoonByIDFromDoc(doc *goquery.Document, b *OGame, moonID MoonID) (Moon, error)
	ExtractCelestialsFromDoc(doc *goquery.Document, b *OGame) ([]Celestial, error)
	ExtractCelestialFromDoc(doc *goquery.Document, b *OGame, v interface{}) (Celestial, error)
	ExtractResourcesBuildingsFromDoc(doc *goquery.Document) (ResourcesBuildings, error)
	ExtractDefenseFromDoc(doc *goquery.Document) (DefensesInfos, error)
	ExtractShipsFromDoc(doc *goquery.Document) (ShipsInfos, error)
	ExtractFacilitiesFromDoc(doc *goquery.Document) (Facilities, error)
	ExtractResearchFromDoc(doc *goquery.Document) Researches
	ExtractOGameSessionFromDoc(doc *goquery.Document) string
	ExtractAttacksFromDoc(doc *goquery.Document, clock clockwork.Clock) ([]AttackEvent, error)
	ExtractOfferOfTheDayFromDoc(doc *goquery.Document) (price int64, importToken string, planetResources PlanetResources, multiplier Multiplier, err error)
	ExtractProductionFromDoc(doc *goquery.Document) ([]Quantifiable, error)
	ExtractOverviewProductionFromDoc(doc *goquery.Document) ([]Quantifiable, error)
	ExtractFleet1ShipsFromDoc(doc *goquery.Document) (s ShipsInfos)
	ExtractEspionageReportMessageIDsFromDoc(doc *goquery.Document) ([]EspionageReportSummary, int64)
	ExtractCombatReportMessagesFromDoc(doc *goquery.Document) ([]CombatReportSummary, int64)
	ExtractExpeditionMessagesFromDoc(doc *goquery.Document, location *time.Location) ([]ExpeditionMessage, int64, error)
	ExtractEspionageReportFromDoc(doc *goquery.Document, location *time.Location) (EspionageReport, error)
	ExtractResourcesProductionsFromDoc(doc *goquery.Document) (Resources, error)
	ExtractPreferencesFromDoc(doc *goquery.Document) Preferences
	ExtractResourceSettingsFromDoc(doc *goquery.Document) (ResourceSettings, error)
	ExtractFleetsFromEventListFromDoc(doc *goquery.Document) []Fleet
	ExtractIPMFromDoc(doc *goquery.Document) (duration, max int64, token string)
	ExtractFleetsFromDoc(doc *goquery.Document, location *time.Location) (res []Fleet)
	ExtractSlotsFromDoc(doc *goquery.Document) Slots
	ExtractServerTimeFromDoc(doc *goquery.Document) (time.Time, error)
	ExtractSpioAnzFromDoc(doc *goquery.Document) int64
	ExtractDisableChatBarFromDoc(doc *goquery.Document) bool
	ExtractDisableOutlawWarningFromDoc(doc *goquery.Document) bool
	ExtractMobileVersionFromDoc(doc *goquery.Document) bool
	ExtractShowOldDropDownsFromDoc(doc *goquery.Document) bool
	ExtractActivateAutofocusFromDoc(doc *goquery.Document) bool
	ExtractEventsShowFromDoc(doc *goquery.Document) int64
	ExtractSortSettingFromDoc(doc *goquery.Document) int64
	ExtractSortOrderFromDoc(doc *goquery.Document) int64
	ExtractShowDetailOverlayFromDoc(doc *goquery.Document) bool
	ExtractAnimatedSlidersFromDoc(doc *goquery.Document) bool
	ExtractAnimatedOverviewFromDoc(doc *goquery.Document) bool
	ExtractPopupsNoticesFromDoc(doc *goquery.Document) bool
	ExtractPopopsCombatreportFromDoc(doc *goquery.Document) bool
	ExtractSpioReportPicturesFromDoc(doc *goquery.Document) bool
	ExtractMsgResultsPerPageFromDoc(doc *goquery.Document) int64
	ExtractAuctioneerNotificationsFromDoc(doc *goquery.Document) bool
	ExtractEconomyNotificationsFromDoc(doc *goquery.Document) bool
	ExtractShowActivityMinutesFromDoc(doc *goquery.Document) bool
	ExtractPreserveSystemOnPlanetChangeFromDoc(doc *goquery.Document) bool
	ExtractNotifBuildListFromDoc(doc *goquery.Document) bool
	ExtractNotifFriendlyFleetActivitiesFromDoc(doc *goquery.Document) bool
	ExtractNotifHostileFleetActivitiesFromDoc(doc *goquery.Document) bool
	ExtractNotifForeignEspionageFromDoc(doc *goquery.Document) bool
	ExtractNotifAllianceBroadcastsFromDoc(doc *goquery.Document) bool
	ExtractNotifAllianceMessagesFromDoc(doc *goquery.Document) bool
	ExtractNotifAuctionsFromDoc(doc *goquery.Document) bool
	ExtractNotifAccountFromDoc(doc *goquery.Document) bool
	ExtractPlanetCoordinate(pageHTML []byte) (Coordinate, error)
	ExtractPlanetID(pageHTML []byte) (CelestialID, error)
	ExtractPlanetIDFromDoc(doc *goquery.Document) (CelestialID, error)
	ExtractOverviewShipSumCountdownFromBytes(pageHTML []byte) int64
	ExtractOGameTimestampFromBytes(pageHTML []byte) int64
	ExtractPlanetType(pageHTML []byte) (CelestialType, error)
	ExtractPlanetTypeFromDoc(doc *goquery.Document) (CelestialType, error)
	ExtractAjaxChatToken(pageHTML []byte) (string, error)
	ExtractCancelFleetToken(pageHTML []byte, fleetID FleetID) (string, error)
	ExtractUserInfos(pageHTML []byte, lang string) (UserInfos, error)
	ExtractResourcesDetails(pageHTML []byte) (out ResourcesDetails, err error)
	ExtractTechs(pageHTML []byte) (ResourcesBuildings, Facilities, ShipsInfos, DefensesInfos, Researches, error)
	ExtractCoord(v string) (coord Coordinate)
	ExtractGalaxyInfos(pageHTML []byte, botPlayerName string, botPlayerID, botPlayerRank int64) (SystemInfos, error)
	ExtractPhalanx(pageHTML []byte) ([]Fleet, error)
	ExtractJumpGate(pageHTML []byte) (ShipsInfos, string, []MoonID, int64)
	ExtractFederation(pageHTML []byte) url.Values
	ExtractConstructions(pageHTML []byte) (buildingID ID, buildingCountdown int64, researchID ID, researchCountdown int64)
	ExtractFleetDeutSaveFactor(pageHTML []byte) float64
	ExtractCancelBuildingInfos(pageHTML []byte) (token string, techID, listID int64, err error)
	ExtractCancelResearchInfos(pageHTML []byte) (token string, techID, listID int64, err error)
	ExtractEmpire(pageHTML []byte) ([]EmpireCelestial, error)
	ExtractEmpireJSON(pageHTML []byte) (interface{}, error)
	ExtractCharacterClass(pageHTML []byte) (CharacterClass, error)
	ExtractCharacterClassFromDoc(doc *goquery.Document) (CharacterClass, error)
	ExtractCommander(pageHTML []byte) bool
	ExtractAdmiral(pageHTML []byte) bool
	ExtractEngineer(pageHTML []byte) bool
	ExtractGeologist(pageHTML []byte) bool
	ExtractTechnocrat(pageHTML []byte) bool
	ExtractCommanderFromDoc(doc *goquery.Document) bool
	ExtractAdmiralFromDoc(doc *goquery.Document) bool
	ExtractEngineerFromDoc(doc *goquery.Document) bool
	ExtractGeologistFromDoc(doc *goquery.Document) bool
	ExtractTechnocratFromDoc(doc *goquery.Document) bool
	ExtractAuction(pageHTML []byte) (Auction, error)
	ExtractHighscore(pageHTML []byte) (Highscore, error)
	ExtractHighscoreFromDoc(doc *goquery.Document) (Highscore, error)
	ExtractAllResources(pageHTML []byte) (map[CelestialID]Resources, error)
	ExtractDMCosts(pageHTML []byte) (DMCosts, error)
	ExtractBuffActivation(pageHTML []byte) (string, []Item, error)
	ExtractActiveItems(pageHTML []byte) ([]ActiveItem, error)
	ExtractIsMobile(pageHTML []byte) bool
	ExtractIsMobileFromDoc(doc *goquery.Document) bool
}
