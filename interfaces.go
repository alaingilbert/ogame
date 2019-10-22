package ogame

import (
	"net/url"
	"time"
)

// Wrapper all available functions to control ogame bot
type Wrapper interface {
	SetLoginProxy(proxy, username, password string) error
	SetLoginWrapper(func(func() error) error)
	GetClient() *OGameClient
	Enable()
	Disable()
	IsEnabled() bool
	Quiet(bool)
	Tx(clb func(tx *Prioritize) error) error
	Begin() *Prioritize
	WithPriority(priority int) *Prioritize
	GetPublicIP() (string, error)
	OnStateChange(clb func(locked bool, actor string))
	GetState() (bool, string)
	IsLocked() bool
	GetSession() string
	AddAccount(number int, lang string) (NewAccount, error)
	GetServer() Server
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
	GetUniverseSpeed() int
	GetUniverseSpeedFleet() int
	GetResearchSpeed() int
	SetResearchSpeed(int)
	GetNbSystems() int
	SetNbSystems(int)
	IsDonutGalaxy() bool
	IsDonutSystem() bool
	FleetDeutSaveFactor() float64
	ServerVersion() string
	ServerTime() time.Time
	Location() *time.Location
	IsUnderAttack() bool
	GetUserInfos() UserInfos
	SendMessage(playerID int, message string) error
	SendMessageAlliance(associationID int, message string) error
	ReconnectChat() bool
	GetFleets() ([]Fleet, Slots)
	GetFleetsFromEventList() []Fleet
	CancelFleet(FleetID) error
	GetAttacks() []AttackEvent
	GalaxyInfos(galaxy, system int) (SystemInfos, error)
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
	GetEspionageReport(msgID int) (EspionageReport, error)
	GetCombatReportSummaryFor(Coordinate) (CombatReportSummary, error)
	//GetCombatReport(msgID int) (CombatReport, error)
	DeleteMessage(msgID int) error
	Distance(origin, destination Coordinate) int
	FlightTime(origin, destination Coordinate, speed Speed, ships ShipsInfos) (secs time.Duration, fuel int)
	RegisterChatCallback(func(ChatMsg))
	RegisterAuctioneerCallback(func([]byte))
	RegisterHTMLInterceptor(func(method, url string, params, payload url.Values, pageHTML []byte))
	GetSlots() Slots
	BuyOfferOfTheDay() error
	BytesDownloaded() int64
	BytesUploaded() int64
	CreateUnion(fleet Fleet) (int, error)

	// Planet or Moon functions
	GetResources(CelestialID) (Resources, error)
	GetResourcesDetails(CelestialID) (ResourcesDetails, error)
	SendFleet(celestialID CelestialID, ships []Quantifiable, speed Speed, where Coordinate, mission MissionID, resources Resources, expeditiontime, unionID int) (Fleet, error)
	EnsureFleet(celestialID CelestialID, ships []Quantifiable, speed Speed, where Coordinate, mission MissionID, resources Resources, expeditiontime, unionID int) (Fleet, error)
	Build(celestialID CelestialID, id ID, nbr int) error
	BuildCancelable(CelestialID, ID) error
	BuildProduction(celestialID CelestialID, id ID, nbr int) error
	BuildBuilding(celestialID CelestialID, buildingID ID) error
	BuildDefense(celestialID CelestialID, defenseID ID, nbr int) error
	BuildShips(celestialID CelestialID, shipID ID, nbr int) error
	CancelBuilding(CelestialID) error
	TearDown(celestialID CelestialID, id ID) error
	ConstructionsBeingBuilt(CelestialID) (buildingID ID, buildingCountdown int, researchID ID, researchCountdown int)
	GetProduction(CelestialID) ([]Quantifiable, error)
	GetFacilities(CelestialID) (Facilities, error)
	GetDefense(CelestialID) (DefensesInfos, error)
	GetShips(CelestialID) (ShipsInfos, error)
	GetResourcesBuildings(CelestialID) (ResourcesBuildings, error)
	CancelResearch(CelestialID) error
	BuildTechnology(celestialID CelestialID, technologyID ID) error

	// Planet specific functions
	GetResourceSettings(PlanetID) (ResourceSettings, error)
	SetResourceSettings(PlanetID, ResourceSettings) error
	SendIPM(PlanetID, Coordinate, int, ID) (int, error)
	//GetResourcesProductionRatio(PlanetID) (float64, error)
	GetResourcesProductions(PlanetID) (Resources, error)
	GetResourcesProductionsLight(ResourcesBuildings, Researches, ResourceSettings, Temperature) Resources

	// Moon specific functions
	Phalanx(MoonID, Coordinate) ([]Fleet, error)
	UnsafePhalanx(MoonID, Coordinate) ([]Fleet, error)
	JumpGate(origin, dest MoonID, ships ShipsInfos) error
}

// BaseOgameObj base interface for all ogame objects (buildings, technologies, ships, defenses)
type BaseOgameObj interface {
	GetID() ID
	GetName() string
	ConstructionTime(nbr, universeSpeed int, facilities Facilities) time.Duration
	GetRequirements() map[ID]int
	GetPrice(int) Resources
	IsAvailable(CelestialType, ResourcesBuildings, Facilities, Researches, int) bool
}

// Levelable base interface for all levelable ogame objects (buildings, technologies)
type Levelable interface {
	BaseOgameObj
	GetLevel(ResourcesBuildings, Facilities, Researches) int
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
	GetStructuralIntegrity(Researches) int
	GetShieldPower(Researches) int
	GetWeaponPower(Researches) int
	GetRapidfireFrom() map[ID]int
}

// Ship interface implemented by all ships units
type Ship interface {
	DefenderObj
	GetCargoCapacity(Researches) int
	GetSpeed(Researches) int
	GetFuelConsumption() int
	GetRapidfireAgainst() map[ID]int
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
	SendFleet([]Quantifiable, Speed, Coordinate, MissionID, Resources, int, int) (Fleet, error)
	EnsureFleet([]Quantifiable, Speed, Coordinate, MissionID, Resources, int, int) (Fleet, error)
	GetDefense() (DefensesInfos, error)
	GetShips() (ShipsInfos, error)
	BuildDefense(defenseID ID, nbr int) error
	ConstructionsBeingBuilt() (ID, int, ID, int)
	GetProduction() ([]Quantifiable, error)
	GetResourcesBuildings() (ResourcesBuildings, error)
	Build(id ID, nbr int) error
	BuildBuilding(buildingID ID) error
	BuildTechnology(technologyID ID) error
	CancelResearch() error
	CancelBuilding() error
	TearDown(buildingID ID) error
}
