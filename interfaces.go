package ogame

import (
	"crypto/tls"
	"github.com/alaingilbert/ogame/pkg/taskRunner"
	"net/http"
	"net/url"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// Prioritizable ...
type Prioritizable interface {
	Abandon(any) error
	ActivateItem(string, CelestialID) error
	Begin() Prioritizable
	BeginNamed(name string) Prioritizable
	BuyMarketplace(itemID int64, celestialID CelestialID) error
	BuyOfferOfTheDay() error
	CancelFleet(FleetID) error
	CollectAllMarketplaceMessages() error
	CollectMarketplaceMessage(MarketplaceMessage) error
	CreateUnion(fleet Fleet, unionUsers []string) (int64, error)
	DeleteAllMessagesFromTab(tabID MessagesTabID) error
	DeleteMessage(msgID int64) error
	DoAuction(bid map[CelestialID]Resources) error
	Done()
	FlightTime(origin, destination Coordinate, speed Speed, ships ShipsInfos, mission MissionID) (secs, fuel int64)
	GalaxyInfos(galaxy, system int64, opts ...Option) (SystemInfos, error)
	GetActiveItems(CelestialID) ([]ActiveItem, error)
	GetAllResources() (map[CelestialID]Resources, error)
	GetAttacks(...Option) ([]AttackEvent, error)
	GetAuction() (Auction, error)
	GetCachedResearch() Researches
	GetCelestial(any) (Celestial, error)
	GetCelestials() ([]Celestial, error)
	GetCombatReportSummaryFor(Coordinate) (CombatReportSummary, error)
	GetDMCosts(CelestialID) (DMCosts, error)
	GetEmpire(CelestialType) ([]EmpireCelestial, error)
	GetEmpireJSON(nbr int64) (any, error)
	GetEspionageReport(msgID int64) (EspionageReport, error)
	GetEspionageReportFor(Coordinate) (EspionageReport, error)
	GetEspionageReportMessages() ([]EspionageReportSummary, error)
	GetExpeditionMessageAt(time.Time) (ExpeditionMessage, error)
	GetExpeditionMessages() ([]ExpeditionMessage, error)
	GetFleets(...Option) ([]Fleet, Slots)
	GetFleetsFromEventList() []Fleet
	GetItems(CelestialID) ([]Item, error)
	GetMoon(any) (Moon, error)
	GetMoons() []Moon
	GetPageContent(url.Values) ([]byte, error)
	GetPlanet(any) (Planet, error)
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
	OfferBuyMarketplace(itemID any, quantity, priceType, price, priceRange int64, celestialID CelestialID) error
	OfferSellMarketplace(itemID any, quantity, priceType, price, priceRange int64, celestialID CelestialID) error
	PostPageContent(url.Values, url.Values) ([]byte, error)
	RecruitOfficer(typ, days int64) error
	SendMessage(playerID int64, message string) error
	SendMessageAlliance(associationID int64, message string) error
	ServerTime() time.Time
	SetInitiator(initiator string) Prioritizable
	SetVacationMode() error
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
	EnsureFleet(celestialID CelestialID, ships []Quantifiable, speed Speed, where Coordinate, mission MissionID, resources Resources, holdingTime, unionID int64) (Fleet, error)
	GetDefense(CelestialID, ...Option) (DefensesInfos, error)
	GetFacilities(CelestialID, ...Option) (Facilities, error)
	GetProduction(CelestialID) ([]Quantifiable, int64, error)
	GetResources(CelestialID) (Resources, error)
	GetResourcesBuildings(CelestialID, ...Option) (ResourcesBuildings, error)
	GetResourcesDetails(CelestialID) (ResourcesDetails, error)
	GetShips(CelestialID, ...Option) (ShipsInfos, error)
	GetTechs(celestialID CelestialID) (ResourcesBuildings, Facilities, ShipsInfos, DefensesInfos, Researches, error)
	SendFleet(celestialID CelestialID, ships []Quantifiable, speed Speed, where Coordinate, mission MissionID, resources Resources, holdingTime, unionID int64) (Fleet, error)
	TearDown(celestialID CelestialID, id ID) error

	// Planet specific functions
	DestroyRockets(PlanetID, int64, int64) error
	GetResourceSettings(PlanetID, ...Option) (ResourceSettings, error)
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
	AddAccount(number int, lang string) (*AddAccountRes, error)
	BytesDownloaded() int64
	BytesUploaded() int64
	CharacterClass() CharacterClass
	ConstructionTime(id ID, nbr int64, facilities Facilities) time.Duration
	Disable()
	Distance(origin, destination Coordinate) int64
	Enable()
	FleetDeutSaveFactor() float64
	GetCachedCelestial(any) Celestial
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
	RegisterChatCallback(func(ChatMsg))
	RegisterHTMLInterceptor(func(method, url string, params, payload url.Values, pageHTML []byte))
	RegisterWSCallback(string, func([]byte))
	RemoveWSCallback(string)
	ServerURL() string
	ServerVersion() string
	SetClient(*OGameClient)
	SetGetServerDataWrapper(func(func() (ServerData, error)) (ServerData, error))
	SetLoginWrapper(func(func() (bool, error)) error)
	SetOGameCredentials(username, password, otpSecret, bearerToken string)
	SetProxy(proxyAddress, username, password, proxyType string, loginOnly bool, config *tls.Config) error
	SetUserAgent(newUserAgent string)
	ValidateAccount(code string) error
	WithPriority(priority taskRunner.Priority) Prioritizable
}

// BaseOgameObj base interface for all ogame objects (buildings, technologies, ships, defenses)
type BaseOgameObj interface {
	ConstructionTime(nbr, universeSpeed int64, facilities Facilities, hasTechnocrat, isDiscoverer bool) time.Duration
	GetID() ID
	GetName() string
	GetPrice(int64) Resources
	GetRequirements() map[ID]int64
	IsAvailable(CelestialType, LazyResourcesBuildings, LazyFacilities, LazyResearches, int64, CharacterClass) bool
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
	GetRapidfireAgainst() map[ID]int64
	GetRapidfireFrom() map[ID]int64
	GetShieldPower(Researches) int64
	GetStructuralIntegrity(Researches) int64
	GetWeaponPower(Researches) int64
}

// Ship interface implemented by all ships units
type Ship interface {
	DefenderObj
	GetCargoCapacity(techs Researches, probeRaids, isCollector, isPioneers bool) int64
	GetFuelConsumption(techs Researches, fleetDeutSaveFactor float64, isGeneral bool) int64
	GetSpeed(techs Researches, isCollector, isGeneral bool) int64
}

// Defense interface implemented by all defenses units
type Defense interface {
	DefenderObj
}

type ICelestial interface {
	CelestialID() CelestialID
	Name() string
	Diameter() int64
	Fields() Fields
	Coordinate() Coordinate
	Img() string
}

// Celestial ...
type Celestial interface {
	ActivateItem(string) error
	Build(id ID, nbr int64) error
	BuildBuilding(buildingID ID) error
	BuildDefense(defenseID ID, nbr int64) error
	BuildTechnology(technologyID ID) error
	CancelBuilding() error
	CancelResearch() error
	ConstructionsBeingBuilt() (ID, int64, ID, int64)
	EnsureFleet([]Quantifiable, Speed, Coordinate, MissionID, Resources, int64, int64) (Fleet, error)
	GetCoordinate() Coordinate
	GetDefense(...Option) (DefensesInfos, error)
	GetDiameter() int64
	GetFacilities(...Option) (Facilities, error)
	GetFields() Fields
	GetID() CelestialID
	GetItems() ([]Item, error)
	GetName() string
	GetProduction() ([]Quantifiable, int64, error)
	GetResources() (Resources, error)
	GetResourcesBuildings(...Option) (ResourcesBuildings, error)
	GetResourcesDetails() (ResourcesDetails, error)
	GetShips(...Option) (ShipsInfos, error)
	GetType() CelestialType
	SendFleet([]Quantifiable, Speed, Coordinate, MissionID, Resources, int64, int64) (Fleet, error)
	TearDown(buildingID ID) error
}

// Extractor ...
type Extractor interface {
	ExtractActiveItems(pageHTML []byte) ([]ActiveItem, error)
	ExtractAdmiral(pageHTML []byte) bool
	ExtractAjaxChatToken(pageHTML []byte) (string, error)
	ExtractAllResources(pageHTML []byte) (map[CelestialID]Resources, error)
	ExtractAttacks(pageHTML []byte, ownCoords []Coordinate) ([]AttackEvent, error)
	ExtractAuction(pageHTML []byte) (Auction, error)
	ExtractBuffActivation(pageHTML []byte) (string, []Item, error)
	ExtractCancelBuildingInfos(pageHTML []byte) (token string, techID, listID int64, err error)
	ExtractCancelFleetToken(pageHTML []byte, fleetID FleetID) (string, error)
	ExtractCancelResearchInfos(pageHTML []byte) (token string, techID, listID int64, err error)
	ExtractCelestial(pageHTML []byte, v any) (ICelestial, error)
	ExtractCelestials(pageHTML []byte) ([]ICelestial, error)
	ExtractCharacterClass(pageHTML []byte) (CharacterClass, error)
	ExtractCombatReportMessagesSummary(pageHTML []byte) ([]CombatReportSummary, int64)
	ExtractCommander(pageHTML []byte) bool
	ExtractConstructions(pageHTML []byte) (buildingID ID, buildingCountdown int64, researchID ID, researchCountdown int64)
	ExtractCoord(v string) (coord Coordinate)
	ExtractDMCosts(pageHTML []byte) (DMCosts, error)
	ExtractDefense(pageHTML []byte) (DefensesInfos, error)
	ExtractDestroyRockets(pageHTML []byte) (abm, ipm int64, token string, err error)
	ExtractEmpire(pageHTML []byte) ([]EmpireCelestial, error)
	ExtractEmpireJSON(pageHTML []byte) (any, error)
	ExtractEngineer(pageHTML []byte) bool
	ExtractEspionageReport(pageHTML []byte, location *time.Location) (EspionageReport, error)
	ExtractEspionageReportMessageIDs(pageHTML []byte) ([]EspionageReportSummary, int64)
	ExtractExpeditionMessages(pageHTML []byte, location *time.Location) ([]ExpeditionMessage, int64, error)
	ExtractFacilities(pageHTML []byte) (Facilities, error)
	ExtractFederation(pageHTML []byte) url.Values
	ExtractFleet1Ships(pageHTML []byte) ShipsInfos
	ExtractFleetDeutSaveFactor(pageHTML []byte) float64
	ExtractFleets(pageHTML []byte, location *time.Location) (res []Fleet)
	ExtractFleetsFromEventList(pageHTML []byte) []Fleet
	ExtractGalaxyInfos(pageHTML []byte, botPlayerName string, botPlayerID, botPlayerRank int64) (SystemInfos, error)
	ExtractGeologist(pageHTML []byte) bool
	ExtractHiddenFields(pageHTML []byte) (fields url.Values)
	ExtractHighscore(pageHTML []byte) (Highscore, error)
	ExtractIPM(pageHTML []byte) (duration, max int64, token string)
	ExtractIsInVacation(pageHTML []byte) bool
	ExtractIsMobile(pageHTML []byte) bool
	ExtractJumpGate(pageHTML []byte) (ShipsInfos, string, []MoonID, int64)
	ExtractMarketplaceMessages(pageHTML []byte, location *time.Location) ([]MarketplaceMessage, int64, error)
	ExtractMoon(pageHTML []byte, v any) (ExtractorMoon, error)
	ExtractMoons(pageHTML []byte) []ExtractorMoon
	ExtractOGameTimestampFromBytes(pageHTML []byte) int64
	ExtractOfferOfTheDay(pageHTML []byte) (int64, string, PlanetResources, Multiplier, error)
	ExtractOgameTimestamp(pageHTML []byte) int64
	ExtractOverviewProduction(pageHTML []byte) ([]Quantifiable, int64, error)
	ExtractOverviewShipSumCountdownFromBytes(pageHTML []byte) int64
	ExtractPhalanx(pageHTML []byte) ([]Fleet, error)
	ExtractPlanet(pageHTML []byte, v any) (ExtractorPlanet, error)
	ExtractPlanetCoordinate(pageHTML []byte) (Coordinate, error)
	ExtractPlanetID(pageHTML []byte) (CelestialID, error)
	ExtractPlanetType(pageHTML []byte) (CelestialType, error)
	ExtractPlanets(pageHTML []byte) []ExtractorPlanet
	ExtractPreferences(pageHTML []byte) Preferences
	ExtractPreferencesShowActivityMinutes(pageHTML []byte) bool
	ExtractPremiumToken(pageHTML []byte, days int64) (token string, err error)
	ExtractProduction(pageHTML []byte) ([]Quantifiable, int64, error)
	ExtractResearch(pageHTML []byte) Researches
	ExtractResourceSettings(pageHTML []byte) (ResourceSettings, error)
	ExtractResources(pageHTML []byte) Resources
	ExtractResourcesBuildings(pageHTML []byte) (ResourcesBuildings, error)
	ExtractResourcesDetails(pageHTML []byte) (out ResourcesDetails, err error)
	ExtractResourcesDetailsFromFullPage(pageHTML []byte) ResourcesDetails
	ExtractResourcesProductions(pageHTML []byte) (Resources, error)
	ExtractServerTime(pageHTML []byte) (time.Time, error)
	ExtractShips(pageHTML []byte) (ShipsInfos, error)
	ExtractSlots(pageHTML []byte) Slots
	ExtractSpioAnz(pageHTML []byte) int64
	ExtractTechnocrat(pageHTML []byte) bool
	ExtractTechs(pageHTML []byte) (ResourcesBuildings, Facilities, ShipsInfos, DefensesInfos, Researches, error)
	ExtractUserInfos(pageHTML []byte, lang string) (UserInfos, error)
	ExtractActivateAutofocusFromDoc(doc *goquery.Document) bool
	ExtractAdmiralFromDoc(doc *goquery.Document) bool
	ExtractAnimatedOverviewFromDoc(doc *goquery.Document) bool
	ExtractAnimatedSlidersFromDoc(doc *goquery.Document) bool
	ExtractAttacksFromDoc(doc *goquery.Document, ownCoords []Coordinate) ([]AttackEvent, error)
	ExtractAuctioneerNotificationsFromDoc(doc *goquery.Document) bool
	ExtractBodyIDFromDoc(doc *goquery.Document) string
	ExtractCelestialFromDoc(doc *goquery.Document, v any) (ICelestial, error)
	ExtractCelestialsFromDoc(doc *goquery.Document) ([]ICelestial, error)
	ExtractCharacterClassFromDoc(doc *goquery.Document) (CharacterClass, error)
	ExtractCombatReportMessagesFromDoc(doc *goquery.Document) ([]CombatReportSummary, int64)
	ExtractCommanderFromDoc(doc *goquery.Document) bool
	ExtractDefenseFromDoc(doc *goquery.Document) (DefensesInfos, error)
	ExtractDisableChatBarFromDoc(doc *goquery.Document) bool
	ExtractDisableOutlawWarningFromDoc(doc *goquery.Document) bool
	ExtractEconomyNotificationsFromDoc(doc *goquery.Document) bool
	ExtractEngineerFromDoc(doc *goquery.Document) bool
	ExtractEspionageReportFromDoc(doc *goquery.Document, location *time.Location) (EspionageReport, error)
	ExtractEspionageReportMessageIDsFromDoc(doc *goquery.Document) ([]EspionageReportSummary, int64)
	ExtractEventsShowFromDoc(doc *goquery.Document) int64
	ExtractExpeditionMessagesFromDoc(doc *goquery.Document, location *time.Location) ([]ExpeditionMessage, int64, error)
	ExtractFacilitiesFromDoc(doc *goquery.Document) (Facilities, error)
	ExtractFleet1ShipsFromDoc(doc *goquery.Document) (s ShipsInfos)
	ExtractFleetsFromDoc(doc *goquery.Document, location *time.Location) (res []Fleet)
	ExtractFleetsFromEventListFromDoc(doc *goquery.Document) []Fleet
	ExtractGeologistFromDoc(doc *goquery.Document) bool
	ExtractHiddenFieldsFromDoc(doc *goquery.Document) url.Values
	ExtractHighscoreFromDoc(doc *goquery.Document) (Highscore, error)
	ExtractIPMFromDoc(doc *goquery.Document) (duration, max int64, token string)
	ExtractIsInVacationFromDoc(doc *goquery.Document) bool
	ExtractIsMobileFromDoc(doc *goquery.Document) bool
	ExtractMobileVersionFromDoc(doc *goquery.Document) bool
	ExtractMoonFromDoc(doc *goquery.Document, v any) (ExtractorMoon, error)
	ExtractMoonsFromDoc(doc *goquery.Document) []ExtractorMoon
	ExtractMsgResultsPerPageFromDoc(doc *goquery.Document) int64
	ExtractNotifAccountFromDoc(doc *goquery.Document) bool
	ExtractNotifAllianceBroadcastsFromDoc(doc *goquery.Document) bool
	ExtractNotifAllianceMessagesFromDoc(doc *goquery.Document) bool
	ExtractNotifAuctionsFromDoc(doc *goquery.Document) bool
	ExtractNotifBuildListFromDoc(doc *goquery.Document) bool
	ExtractNotifForeignEspionageFromDoc(doc *goquery.Document) bool
	ExtractNotifFriendlyFleetActivitiesFromDoc(doc *goquery.Document) bool
	ExtractNotifHostileFleetActivitiesFromDoc(doc *goquery.Document) bool
	ExtractOGameSessionFromDoc(doc *goquery.Document) string
	ExtractOfferOfTheDayFromDoc(doc *goquery.Document) (price int64, importToken string, planetResources PlanetResources, multiplier Multiplier, err error)
	ExtractOgameTimestampFromDoc(doc *goquery.Document) int64
	ExtractOverviewProductionFromDoc(doc *goquery.Document) ([]Quantifiable, error)
	ExtractPlanetFromDoc(doc *goquery.Document, v any) (ExtractorPlanet, error)
	ExtractPlanetIDFromDoc(doc *goquery.Document) (CelestialID, error)
	ExtractPlanetTypeFromDoc(doc *goquery.Document) (CelestialType, error)
	ExtractPlanetsFromDoc(doc *goquery.Document) []ExtractorPlanet
	ExtractPopopsCombatreportFromDoc(doc *goquery.Document) bool
	ExtractPopupsNoticesFromDoc(doc *goquery.Document) bool
	ExtractPreferencesFromDoc(doc *goquery.Document) Preferences
	ExtractPreserveSystemOnPlanetChangeFromDoc(doc *goquery.Document) bool
	ExtractProductionFromDoc(doc *goquery.Document) ([]Quantifiable, error)
	ExtractResearchFromDoc(doc *goquery.Document) Researches
	ExtractResourceSettingsFromDoc(doc *goquery.Document) (ResourceSettings, error)
	ExtractResourcesBuildingsFromDoc(doc *goquery.Document) (ResourcesBuildings, error)
	ExtractResourcesDetailsFromFullPageFromDoc(doc *goquery.Document) ResourcesDetails
	ExtractResourcesFromDoc(doc *goquery.Document) Resources
	ExtractResourcesProductionsFromDoc(doc *goquery.Document) (Resources, error)
	ExtractServerTimeFromDoc(doc *goquery.Document) (time.Time, error)
	ExtractShipsFromDoc(doc *goquery.Document) (ShipsInfos, error)
	ExtractShowActivityMinutesFromDoc(doc *goquery.Document) bool
	ExtractShowDetailOverlayFromDoc(doc *goquery.Document) bool
	ExtractShowOldDropDownsFromDoc(doc *goquery.Document) bool
	ExtractSlotsFromDoc(doc *goquery.Document) Slots
	ExtractSortOrderFromDoc(doc *goquery.Document) int64
	ExtractSortSettingFromDoc(doc *goquery.Document) int64
	ExtractSpioAnzFromDoc(doc *goquery.Document) int64
	ExtractSpioReportPicturesFromDoc(doc *goquery.Document) bool
	ExtractTechnocratFromDoc(doc *goquery.Document) bool
}
