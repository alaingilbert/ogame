package ogame

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
	"golang.org/x/net/websocket"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Wrapper all available functions to control ogame bot
type Wrapper interface {
	GetSession() string
	AddAccount(number int, lang string) (NewAccount, error)
	GetServer() Server
	SetUserAgent(newUserAgent string)
	ServerURL() string
	GetLanguage() string
	GetPageContent(url.Values) []byte
	PostPageContent(url.Values, url.Values) []byte
	Login() error
	Logout()
	GetUsername() string
	GetUniverseName() string
	GetUniverseSpeed() int
	GetUniverseSpeedFleet() int
	IsDonutGalaxy() bool
	IsDonutSystem() bool
	ServerVersion() string
	ServerTime() time.Time
	IsUnderAttack() bool
	GetUserInfos() UserInfos
	SendMessage(playerID int, message string) error
	GetFleets() []Fleet
	CancelFleet(FleetID) error
	GetAttacks() []AttackEvent
	GalaxyInfos(galaxy, system int) (SystemInfos, error)
	GetResearch() Researches
	GetCachedPlanets() []Planet
	GetPlanets() []Planet
	GetPlanet(PlanetID) (Planet, error)
	GetPlanetByCoord(Coordinate) (Planet, error)
	GetMoons(MoonID) []Moon
	GetMoon(MoonID) (Moon, error)
	GetMoonByCoord(Coordinate) (Moon, error)
	GetEspionageReportMessages() ([]EspionageReportSummary, error)
	GetEspionageReport(msgID int) (EspionageReport, error)
	DeleteMessage(msgID int) error
	FlightTime(origin, destination Coordinate, speed Speed, ships ShipsInfos) (secs, fuel int)
	RegisterChatCallback(func(ChatMsg))

	// Planet or Moon functions
	GetResources(CelestialID) (Resources, error)
	SendFleet(celestialID CelestialID, ships []Quantifiable, speed Speed, where Coordinate, destType DestinationType, mission MissionID, resources Resources) (FleetID, error)
	Build(celestialID CelestialID, id ID, nbr int) error
	BuildCancelable(CelestialID, ID) error
	BuildProduction(celestialID CelestialID, id ID, nbr int) error
	BuildBuilding(celestialID CelestialID, buildingID ID) error
	BuildDefense(celestialID CelestialID, defenseID ID, nbr int) error
	BuildShips(celestialID CelestialID, shipID ID, nbr int) error
	CancelBuilding(CelestialID) error
	ConstructionsBeingBuilt(CelestialID) (buildingID ID, buildingCountdown int, researchID ID, researchCountdown int)
	GetProduction(CelestialID) ([]Quantifiable, error)
	GetFacilities(CelestialID) (Facilities, error)
	GetDefense(CelestialID) (DefensesInfos, error)
	GetShips(CelestialID) (ShipsInfos, error)

	// Planet specific functions
	GetResourceSettings(PlanetID) (ResourceSettings, error)
	SetResourceSettings(PlanetID, ResourceSettings) error
	GetResourcesBuildings(PlanetID) (ResourcesBuildings, error)
	BuildTechnology(planetID PlanetID, technologyID ID) error
	CancelResearch(PlanetID) error
	//GetResourcesProductionRatio(PlanetID) (float64, error)
	GetResourcesProductions(PlanetID) (Resources, error)

	// Moon specific functions
	Phalanx(MoonID, Coordinate) ([]Fleet, error)
}

const defaultUserAgent = "" +
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) " +
	"AppleWebKit/537.36 (KHTML, like Gecko) " +
	"Chrome/51.0.2704.103 " +
	"Safari/537.36"

// ErrNotLogged returned when the bot is not logged
var ErrNotLogged = errors.New("not logged")

// ErrBadCredentials returned when the provided credentials are invalid
var ErrBadCredentials = errors.New("bad credentials")

// ErrInvalidPlanetID returned when a planet id is invalid
var ErrInvalidPlanetID = errors.New("invalid planet id")

// Send fleet errors
var (
	ErrNoShipSelected = errors.New("no ships to send")

	ErrUninhabitedPlanet                  = errors.New("uninhabited planet")
	ErrNoDebrisField                      = errors.New("no debris field")
	ErrPlayerInVacationMode               = errors.New("player in vacation mode")
	ErrAdminOrGM                          = errors.New("admin or GM")
	ErrNoAstrophysics                     = errors.New("you have to research Astrophysics first")
	ErrNoobProtection                     = errors.New("noob protection")
	ErrPlayerTooStrong                    = errors.New("this planet can not be attacked as the player is to strong")
	ErrNoMoonAvailable                    = errors.New("no moon available")
	ErrNoRecyclerAvailable                = errors.New("no recycler available")
	ErrNoEventsRunning                    = errors.New("there are currently no events running")
	ErrPlanetAlreadyReservecForRelocation = errors.New("this planet has already been reserved for a relocation")
)

// All ogame objects
var (
	AllianceDepot                = newAllianceDepot() // Buildings
	CrystalMine                  = newCrystalMine()
	CrystalStorage               = newCrystalStorage()
	DeuteriumSynthesizer         = newDeuteriumSynthesizer()
	DeuteriumTank                = newDeuteriumTank()
	FusionReactor                = newFusionReactor()
	MetalMine                    = newMetalMine()
	MetalStorage                 = newMetalStorage()
	MissileSilo                  = newMissileSilo()
	NaniteFactory                = newNaniteFactory()
	ResearchLab                  = newResearchLab()
	RoboticsFactory              = newRoboticsFactory()
	SeabedDeuteriumDen           = newSeabedDeuteriumDen()
	ShieldedMetalDen             = newShieldedMetalDen()
	Shipyard                     = newShipyard()
	SolarPlant                   = newSolarPlant()
	SpaceDock                    = newSpaceDock()
	LunarBase                    = newLunarBase()
	SensorPhalanx                = newSensorPhalanx()
	JumpGate                     = newJumpGate()
	Terraformer                  = newTerraformer()
	UndergroundCrystalDen        = newUndergroundCrystalDen()
	SolarSatellite               = newSolarSatellite()
	AntiBallisticMissiles        = newAntiBallisticMissiles() // Defense
	GaussCannon                  = newGaussCannon()
	HeavyLaser                   = newHeavyLaser()
	InterplanetaryMissiles       = newInterplanetaryMissiles()
	IonCannon                    = newIonCannon()
	LargeShieldDome              = newLargeShieldDome()
	LightLaser                   = newLightLaser()
	PlasmaTurret                 = newPlasmaTurret()
	RocketLauncher               = newRocketLauncher()
	SmallShieldDome              = newSmallShieldDome()
	Battlecruiser                = newBattlecruiser() // Ships
	Battleship                   = newBattleship()
	Bomber                       = newBomber()
	ColonyShip                   = newColonyShip()
	Cruiser                      = newCruiser()
	Deathstar                    = newDeathstar()
	Destroyer                    = newDestroyer()
	EspionageProbe               = newEspionageProbe()
	HeavyFighter                 = newHeavyFighter()
	LargeCargo                   = newLargeCargo()
	LightFighter                 = newLightFighter()
	Recycler                     = newRecycler()
	SmallCargo                   = newSmallCargo()
	ArmourTechnology             = newArmourTechnology() // Technologies
	Astrophysics                 = newAstrophysics()
	CombustionDrive              = newCombustionDrive()
	ComputerTechnology           = newComputerTechnology()
	EnergyTechnology             = newEnergyTechnology()
	EspionageTechnology          = newEspionageTechnology()
	GravitonTechnology           = newGravitonTechnology()
	HyperspaceDrive              = newHyperspaceDrive()
	HyperspaceTechnology         = newHyperspaceTechnology()
	ImpulseDrive                 = newImpulseDrive()
	IntergalacticResearchNetwork = newIntergalacticResearchNetwork()
	IonTechnology                = newIonTechnology()
	LaserTechnology              = newLaserTechnology()
	PlasmaTechnology             = newPlasmaTechnology()
	ShieldingTechnology          = newShieldingTechnology()
	WeaponsTechnology            = newWeaponsTechnology()
)

// CelestialID represent either a PlanetID or a MoonID
type CelestialID int

// BaseOgameObj base interface for all ogame objects (buildings, technologies, ships, defenses)
type BaseOgameObj interface {
	GetID() ID
	GetName() string
	ConstructionTime(nbr, universeSpeed int, facilities Facilities) time.Duration
	GetRequirements() map[ID]int
	GetPrice(int) Resources
	IsAvailable(ResourcesBuildings, Facilities, Researches, int) bool
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
	GetCargoCapacity() int
	GetSpeed(researches Researches) int
	GetFuelConsumption() int
	GetRapidfireAgainst() map[ID]int
}

// Defense interface implemented by all defenses units
type Defense interface {
	DefenderObj
}

// ObjsStruct structure containing all possible ogame objects
type ObjsStruct struct {
	AllianceDepot                *allianceDepot
	CrystalMine                  *crystalMine
	CrystalStorage               *crystalStorage
	DeuteriumSynthesizer         *deuteriumSynthesizer
	DeuteriumTank                *deuteriumTank
	FusionReactor                *fusionReactor
	MetalMine                    *metalMine
	MetalStorage                 *metalStorage
	MissileSilo                  *missileSilo
	NaniteFactory                *naniteFactory
	ResearchLab                  *researchLab
	RoboticsFactory              *roboticsFactory
	SeabedDeuteriumDen           *seabedDeuteriumDen
	ShieldedMetalDen             *shieldedMetalDen
	Shipyard                     *shipyard
	SolarPlant                   *solarPlant
	SpaceDock                    *spaceDock
	LunarBase                    *lunarBase
	SensorPhalanx                *sensorPhalanx
	JumpGate                     *jumpGate
	Terraformer                  *terraformer
	UndergroundCrystalDen        *undergroundCrystalDen
	SolarSatellite               *solarSatellite
	AntiBallisticMissiles        *antiBallisticMissiles
	GaussCannon                  *gaussCannon
	HeavyLaser                   *heavyLaser
	InterplanetaryMissiles       *interplanetaryMissiles
	IonCannon                    *ionCannon
	LargeShieldDome              *largeShieldDome
	LightLaser                   *lightLaser
	PlasmaTurret                 *plasmaTurret
	RocketLauncher               *rocketLauncher
	SmallShieldDome              *smallShieldDome
	Battlecruiser                *battlecruiser
	Battleship                   *battleship
	Bomber                       *bomber
	ColonyShip                   *colonyShip
	Cruiser                      *cruiser
	Deathstar                    *deathstar
	Destroyer                    *destroyer
	EspionageProbe               *espionageProbe
	HeavyFighter                 *heavyFighter
	LargeCargo                   *largeCargo
	LightFighter                 *lightFighter
	Recycler                     *recycler
	SmallCargo                   *smallCargo
	ArmourTechnology             *armourTechnology
	Astrophysics                 *astrophysics
	CombustionDrive              *combustionDrive
	ComputerTechnology           *computerTechnology
	EnergyTechnology             *energyTechnology
	EspionageTechnology          *espionageTechnology
	GravitonTechnology           *gravitonTechnology
	HyperspaceDrive              *hyperspaceDrive
	HyperspaceTechnology         *hyperspaceTechnology
	ImpulseDrive                 *impulseDrive
	IntergalacticResearchNetwork *intergalacticResearchNetwork
	IonTechnology                *ionTechnology
	LaserTechnology              *laserTechnology
	PlasmaTechnology             *plasmaTechnology
	ShieldingTechnology          *shieldingTechnology
	WeaponsTechnology            *weaponsTechnology
}

// Objs all ogame objects
var Objs = ObjsStruct{
	AllianceDepot:                AllianceDepot,
	CrystalMine:                  CrystalMine,
	CrystalStorage:               CrystalStorage,
	DeuteriumSynthesizer:         DeuteriumSynthesizer,
	DeuteriumTank:                DeuteriumTank,
	FusionReactor:                FusionReactor,
	MetalMine:                    MetalMine,
	MetalStorage:                 MetalStorage,
	MissileSilo:                  MissileSilo,
	NaniteFactory:                NaniteFactory,
	ResearchLab:                  ResearchLab,
	RoboticsFactory:              RoboticsFactory,
	SeabedDeuteriumDen:           SeabedDeuteriumDen,
	ShieldedMetalDen:             ShieldedMetalDen,
	Shipyard:                     Shipyard,
	SolarPlant:                   SolarPlant,
	SpaceDock:                    SpaceDock,
	LunarBase:                    LunarBase,
	SensorPhalanx:                SensorPhalanx,
	JumpGate:                     JumpGate,
	Terraformer:                  Terraformer,
	UndergroundCrystalDen:        UndergroundCrystalDen,
	SolarSatellite:               SolarSatellite,
	AntiBallisticMissiles:        AntiBallisticMissiles,
	GaussCannon:                  GaussCannon,
	HeavyLaser:                   HeavyLaser,
	InterplanetaryMissiles:       InterplanetaryMissiles,
	IonCannon:                    IonCannon,
	LargeShieldDome:              LargeShieldDome,
	LightLaser:                   LightLaser,
	PlasmaTurret:                 PlasmaTurret,
	RocketLauncher:               RocketLauncher,
	SmallShieldDome:              SmallShieldDome,
	Battlecruiser:                Battlecruiser,
	Battleship:                   Battleship,
	Bomber:                       Bomber,
	ColonyShip:                   ColonyShip,
	Cruiser:                      Cruiser,
	Deathstar:                    Deathstar,
	Destroyer:                    Destroyer,
	EspionageProbe:               EspionageProbe,
	HeavyFighter:                 HeavyFighter,
	LargeCargo:                   LargeCargo,
	LightFighter:                 LightFighter,
	Recycler:                     Recycler,
	SmallCargo:                   SmallCargo,
	ArmourTechnology:             ArmourTechnology,
	Astrophysics:                 Astrophysics,
	CombustionDrive:              CombustionDrive,
	ComputerTechnology:           ComputerTechnology,
	EnergyTechnology:             EnergyTechnology,
	EspionageTechnology:          EspionageTechnology,
	GravitonTechnology:           GravitonTechnology,
	HyperspaceDrive:              HyperspaceDrive,
	HyperspaceTechnology:         HyperspaceTechnology,
	ImpulseDrive:                 ImpulseDrive,
	IntergalacticResearchNetwork: IntergalacticResearchNetwork,
	IonTechnology:                IonTechnology,
	LaserTechnology:              LaserTechnology,
	PlasmaTechnology:             PlasmaTechnology,
	ShieldingTechnology:          ShieldingTechnology,
	WeaponsTechnology:            WeaponsTechnology,
}

// Defenses array of all defenses objects
var Defenses = []Defense{
	AntiBallisticMissiles,
	GaussCannon,
	HeavyLaser,
	InterplanetaryMissiles,
	IonCannon,
	LargeShieldDome,
	LightLaser,
	PlasmaTurret,
	RocketLauncher,
	SmallShieldDome,
}

// Ships array of all ships objects
var Ships = []Ship{
	Battlecruiser,
	Battleship,
	Bomber,
	ColonyShip,
	Cruiser,
	Deathstar,
	Destroyer,
	EspionageProbe,
	HeavyFighter,
	LargeCargo,
	LightFighter,
	Recycler,
	SmallCargo,
	SolarSatellite,
}

// Buildings array of all buildings/facilities objects
var Buildings = []Building{
	AllianceDepot,
	CrystalMine,
	CrystalStorage,
	DeuteriumSynthesizer,
	DeuteriumTank,
	FusionReactor,
	MetalMine,
	MetalStorage,
	MissileSilo,
	NaniteFactory,
	ResearchLab,
	RoboticsFactory,
	SeabedDeuteriumDen,
	ShieldedMetalDen,
	Shipyard,
	SolarPlant,
	SpaceDock,
	Terraformer,
	UndergroundCrystalDen,
	SolarSatellite,
}

// Technologies array of all technologies objects
var Technologies = []Technology{
	ArmourTechnology,
	Astrophysics,
	CombustionDrive,
	ComputerTechnology,
	EnergyTechnology,
	EspionageTechnology,
	GravitonTechnology,
	HyperspaceDrive,
	HyperspaceTechnology,
	ImpulseDrive,
	IntergalacticResearchNetwork,
	IonTechnology,
	LaserTechnology,
	PlasmaTechnology,
	ShieldingTechnology,
	WeaponsTechnology,
}

// OGame is a client for ogame.org. It is safe for concurrent use by
// multiple goroutines (thread-safe)
type OGame struct {
	sync.Mutex
	quiet              bool
	Player             UserInfos
	Planets            []Planet
	Universe           string
	Username           string
	password           string
	language           string
	ogameSession       string
	sessionChatCounter int
	server             Server
	location           *time.Location
	universeSpeed      int
	universeSize       int
	universeSpeedFleet int
	donutGalaxy        bool
	donutSystem        bool
	ogameVersion       string
	serverURL          string
	client             *ogameClient
	logger             *log.Logger
	chatCallbacks      []func(msg ChatMsg)
	closeChatCh        chan struct{}
	chatConnected      int32
	chatRetry          *ExponentialBackoff
	ws                 *websocket.Conn
}

// Params parameters for more fine-grained initialization
type Params struct {
	Universe  string
	Username  string
	Password  string
	Lang      string
	AutoLogin bool
	Proxy     string
}

// New creates a new instance of OGame wrapper.
func New(universe, username, password, lang string) (*OGame, error) {
	b := NewNoLogin(universe, username, password, lang)

	if err := b.Login(); err != nil {
		return nil, err
	}

	return b, nil
}

// NewWithParams create a new OGame instance with full control over the possible parameters
func NewWithParams(params Params) (*OGame, error) {
	b := NewNoLogin(params.Universe, params.Username, params.Password, params.Lang)

	if params.Proxy != "" {
		proxyURL, err := url.Parse(params.Proxy)
		if err != nil {
			return nil, err
		}
		b.client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	}

	if params.AutoLogin {
		if err := b.Login(); err != nil {
			return nil, err
		}
	}

	return b, nil
}

// NewNoLogin does not auto login.
func NewNoLogin(universe, username, password, lang string) *OGame {
	b := new(OGame)
	b.quiet = false
	b.logger = log.New(os.Stdout, "", 0)

	b.Universe = universe
	b.Username = username
	b.password = password
	b.language = lang

	jar, _ := cookiejar.New(nil)
	b.client = &ogameClient{}
	b.client.Jar = jar
	b.client.UserAgent = defaultUserAgent

	return b
}

// Quiet mode will not show any informative output
func (b *OGame) Quiet(quiet bool) {
	b.quiet = quiet
}

// SetLogger set a custom logger for the bot
func (b *OGame) SetLogger(logger *log.Logger) {
	b.logger = logger
}

// Terminal styling constants
const (
	knrm = "\x1B[0m"
	kred = "\x1B[31m"
	//kgrn = "\x1B[32m"
	kyel = "\x1B[33m"
	//kblu = "\x1B[34m"
	kmag = "\x1B[35m"
	kcyn = "\x1B[36m"
	kwht = "\x1B[37m"
)

func (b *OGame) log(prefix, color string, v ...interface{}) {
	if !b.quiet {
		_, f, l, _ := runtime.Caller(2)
		args := append([]interface{}{fmt.Sprintf(color+"%s"+knrm+" [%s:%d]", prefix, filepath.Base(f), l)}, v...)
		b.logger.Println(args...)
	}
}

func (b *OGame) trace(v ...interface{}) {
	b.log("TRAC", kwht, v...)
}

func (b *OGame) info(v ...interface{}) {
	b.log("INFO", kcyn, v...)
}

func (b *OGame) warn(v ...interface{}) {
	b.log("WARN", kyel, v...)
}

func (b *OGame) error(v ...interface{}) {
	b.log("ERRO", kred, v...)
}

func (b *OGame) critical(v ...interface{}) {
	b.log("CRIT", kred, v...)
}

func (b *OGame) debug(v ...interface{}) {
	b.log("DEBU", kmag, v...)
}

func (b *OGame) println(v ...interface{}) {
	b.log("PRIN", kwht, v...)
}

// Server ogame information for their servers
type Server struct {
	Language      string
	Number        int
	Name          string
	PlayerCount   int
	PlayersOnline int
	Opened        string
	StartDate     string
	EndDate       *string
	ServerClosed  int
	Prefered      int
	SignupClosed  int
	Settings      struct {
		AKS                      int
		FleetSpeed               int
		WreckField               int
		ServerLabel              string
		EconomySpeed             int
		PlanetFields             int
		UniverseSize             int // Nb of galaxies
		ServerCategory           string
		EspionageProbeRaids      int
		PremiumValidationGift    int
		DebrisFieldFactorShips   int
		DebrisFieldFactorDefence int
	}
}

// ogame cookie name for php session id
const phpSessionIDCookieName = "PHPSESSID"

func getPhpSessionID(client *ogameClient, username, password string) (string, error) {
	payload := url.Values{
		"kid":                   {""},
		"language":              {"en"},
		"autologin":             {"false"},
		"credentials[email]":    {username},
		"credentials[password]": {password},
	}
	req, err := http.NewRequest("POST", "https://lobby-api.ogame.gameforge.com/users", strings.NewReader(payload.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return "", errors.New("OGame server error code : " + resp.Status)
	}

	if resp.StatusCode != 200 {
		return "", ErrBadCredentials
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == phpSessionIDCookieName {
			return cookie.Value, nil
		}
	}

	return "", errors.New(phpSessionIDCookieName + " not found")
}

type account struct {
	Server struct {
		Language string
		Number   int
	}
	ID         int
	Name       string
	LastPlayed string
	Blocked    bool
	Details    []struct {
		Type  string
		Title string
		Value string
	}
	Sitting struct {
		Shared       bool
		EndTime      *int
		CooldownTime *int
	}
}

func getUserAccounts(client *ogameClient, phpSessionID string) ([]account, error) {
	var userAccounts []account
	req, err := http.NewRequest("GET", "https://lobby-api.ogame.gameforge.com/users/me/accounts", nil)
	if err != nil {
		return userAccounts, err
	}
	req.AddCookie(&http.Cookie{Name: phpSessionIDCookieName, Value: phpSessionID})
	resp, err := client.Do(req)
	if err != nil {
		return userAccounts, err
	}
	defer resp.Body.Close()
	by, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return userAccounts, err
	}
	if err := json.Unmarshal(by, &userAccounts); err != nil {
		return userAccounts, err
	}
	return userAccounts, nil
}

func getServers(client *ogameClient) ([]Server, error) {
	var servers []Server
	req, err := http.NewRequest("GET", "https://lobby-api.ogame.gameforge.com/servers", nil)
	if err != nil {
		return servers, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return servers, err
	}
	defer resp.Body.Close()
	by, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return servers, err
	}
	if err := json.Unmarshal(by, &servers); err != nil {
		return servers, err
	}
	return servers, nil
}

func findAccountByName(universe, lang string, accounts []account, servers []Server) (account, Server, error) {
	var server Server
	var acc account
	for _, s := range servers {
		if s.Name == universe && s.Language == lang {
			server = s
			break
		}
	}
	for _, a := range accounts {
		if a.Server.Language == server.Language && a.Server.Number == server.Number {
			acc = a
			break
		}
	}
	if server.Number == 0 {
		return account{}, Server{}, fmt.Errorf("server %s, %s not found", universe, lang)
	}
	if acc.ID == 0 {
		return account{}, Server{}, errors.New("account not found")
	}
	return acc, server, nil
}

func getLoginLink(client *ogameClient, userAccount account, phpSessionID string) (string, error) {
	ogURL := fmt.Sprintf("https://lobby-api.ogame.gameforge.com/users/me/loginLink?id=%d&server[language]=%s&server[number]=%d",
		userAccount.ID, userAccount.Server.Language, userAccount.Server.Number)
	req, err := http.NewRequest("GET", ogURL, nil)
	if err != nil {
		return "", err
	}
	req.AddCookie(&http.Cookie{Name: phpSessionIDCookieName, Value: phpSessionID})
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	by, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var loginLink struct {
		URL string
	}
	if err := json.Unmarshal(by, &loginLink); err != nil {
		return "", err
	}
	return loginLink.URL, nil
}

func (b *OGame) login() error {
	jar, _ := cookiejar.New(nil)
	b.client.Jar = jar

	b.debug("get session")
	phpSessionID, err := getPhpSessionID(b.client, b.Username, b.password)
	if err != nil {
		return err
	}
	b.debug("get user accounts")
	accounts, err := getUserAccounts(b.client, phpSessionID)
	if err != nil {
		return err
	}
	b.debug("get servers")
	servers, err := getServers(b.client)
	if err != nil {
		return err
	}
	b.debug("find account & server for universe")
	userAccount, server, err := findAccountByName(b.Universe, b.language, accounts, servers)
	if err != nil {
		return err
	}
	if userAccount.Blocked {
		return errors.New("your account is banned")
	}
	b.debug("Players online: " + strconv.Itoa(server.PlayersOnline) + ", Players: " + strconv.Itoa(server.PlayerCount))
	b.server = server
	b.language = userAccount.Server.Language
	b.debug("get login link")
	loginLink, err := getLoginLink(b.client, userAccount, phpSessionID)
	if err != nil {
		return err
	}

	r := regexp.MustCompile(`(https://.+\.ogame\.gameforge\.com)/game`)
	res := r.FindStringSubmatch(loginLink)
	if len(res) != 2 {
		return errors.New("failed to get server url")
	}
	b.serverURL = res[1]

	req, err := http.NewRequest("GET", loginLink, nil)
	if err != nil {
		return err
	}
	b.debug("login to universe")
	resp, err := b.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	pageHTML, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	b.debug("extract information from html")
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return err
	}
	b.ogameSession, _ = doc.Find("meta[name=ogame-session]").Attr("content")
	if b.ogameSession == "" {
		return errors.New("bad credentials")
	}
	b.sessionChatCounter = 1

	serverTime, _ := extractServerTime(pageHTML)
	b.location = serverTime.Location()
	b.universeSize = server.Settings.UniverseSize
	b.universeSpeed, _ = strconv.Atoi(doc.Find("meta[name=ogame-universe-speed]").AttrOr("content", "1"))
	b.universeSpeedFleet, _ = strconv.Atoi(doc.Find("meta[name=ogame-universe-speed-fleet]").AttrOr("content", "1"))
	b.donutGalaxy, _ = strconv.ParseBool(doc.Find("meta[name=ogame-donut-galaxy]").AttrOr("content", "1"))
	b.donutSystem, _ = strconv.ParseBool(doc.Find("meta[name=ogame-donut-system]").AttrOr("content", "1"))
	b.ogameVersion = doc.Find("meta[name=ogame-version]").AttrOr("content", "")

	b.Player, _ = extractUserInfos(pageHTML, b.language)
	b.Planets = extractPlanets(pageHTML, b)

	// Extract chat host and port
	m := regexp.MustCompile(`var nodeUrl="https:\\/\\/([^:]+):(\d+)\\/socket.io\\/socket.io.js";`).FindSubmatch(pageHTML)
	chatHost := string(m[1])
	chatPort := string(m[2])

	if atomic.CompareAndSwapInt32(&b.chatConnected, 0, 1) {
		b.closeChatCh = make(chan struct{})
		go func(b *OGame) {
			defer atomic.StoreInt32(&b.chatConnected, 0)
			b.chatRetry = NewExponentialBackoff(60)
		LOOP:
			for {
				select {
				case <-b.closeChatCh:
					break LOOP
				default:
					b.connectChat(chatHost, chatPort)
					b.chatRetry.Wait()
				}
			}
		}(b)
	}

	return nil
}

func (b *OGame) connectChat(host, port string) {
	req, err := http.NewRequest("GET", "https://"+host+":"+port+"/socket.io/1/?t="+strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10), nil)
	if err != nil {
		b.error("failed to create request:", err)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		b.error("failed to get socket.io token:", err)
		return
	}
	defer resp.Body.Close()
	b.chatRetry.Reset()
	by, _ := ioutil.ReadAll(resp.Body)
	token := strings.Split(string(by), ":")[0]

	origin := "https://" + host + ":" + port + "/"
	url := "wss://" + host + ":" + port + "/socket.io/1/websocket/" + token
	b.ws, err = websocket.Dial(url, "", origin)
	if err != nil {
		b.error("failed to dial websocket:", err)
		return
	}

	// Recv msgs
LOOP:
	for {
		select {
		case <-b.closeChatCh:
			break LOOP
		default:
		}

		var buf = make([]byte, 1024*1024)
		b.ws.SetReadDeadline(time.Now().Add(time.Second))
		if _, err = b.ws.Read(buf); err != nil {
			if err == io.EOF {
				b.error("chat eof:", err)
				break
			} else if strings.HasSuffix(err.Error(), "use of closed network connection") {
				break
			} else if strings.HasSuffix(err.Error(), "i/o timeout") {
				continue
			} else {
				b.error("chat unexpected error", err)
			}
		}
		msg := bytes.Trim(buf, "\x00")
		if bytes.Equal(msg, []byte("1::")) {
			b.ws.Write([]byte("1::/chat"))
		} else if bytes.Equal(msg, []byte("1::/chat")) {
			authMsg := `5:` + strconv.Itoa(b.sessionChatCounter) + `+:/chat:{"name":"authorize","args":["` + b.ogameSession + `"]}`
			b.ws.Write([]byte(authMsg))
			b.sessionChatCounter++
		} else if bytes.Equal(msg, []byte("2::")) {
			b.ws.Write([]byte("2::"))
		} else if regexp.MustCompile(`6::/chat:\d+\+\[true]`).Match(msg) {
			b.debug("chat connected")
		} else if regexp.MustCompile(`6::/chat:\d+\+\[false]`).Match(msg) {
			b.error("Failed to connect to chat")
		} else if bytes.HasPrefix(msg, []byte("5::/chat:")) {
			payload := bytes.TrimPrefix(msg, []byte("5::/chat:"))
			var chatPayload ChatPayload
			if err := json.Unmarshal([]byte(payload), &chatPayload); err != nil {
				b.error("Unable to unmarshal chat payload", err, payload)
				continue
			}
			for _, chatMsg := range chatPayload.Args {
				for _, clb := range b.chatCallbacks {
					clb(chatMsg)
				}
			}
		}
	}
}

// ChatPayload ...
type ChatPayload struct {
	Name string    `json:"name"`
	Args []ChatMsg `json:"args"`
}

// ChatMsg ...
type ChatMsg struct {
	SenderID      int    `json:"senderId"`
	SenderName    string `json:"senderName"`
	AssociationID int    `json:"associationId"`
	Text          string `json:"text"`
	ID            int    `json:"id"`
	Date          int    `json:"date"`
}

func (m ChatMsg) String() string {
	return "\n" +
		"     Sender ID: " + strconv.Itoa(m.SenderID) + "\n" +
		"   Sender name: " + m.SenderName + "\n" +
		"Association ID: " + strconv.Itoa(m.AssociationID) + "\n" +
		"          Text: " + m.Text + "\n" +
		"            ID: " + strconv.Itoa(m.ID) + "\n" +
		"          Date: " + strconv.Itoa(m.Date)
}

func (b *OGame) logout() {
	b.getPageContent(url.Values{"page": {"logout"}})
	select {
	case <-b.closeChatCh:
	default:
		close(b.closeChatCh)
		b.ws.Close()
	}
}

func isLogged(pageHTML []byte) bool {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return false
	}
	ogameSession := doc.Find("meta[name=ogame-session]").AttrOr("content", "")
	return ogameSession != ""
}

// IsAjaxPage either the requested page is a partial/ajax page
func IsAjaxPage(vals url.Values) bool {
	page := vals.Get("page")
	ajax := vals.Get("ajax")
	return page == "fetchEventbox" ||
		page == "fetchResources" ||
		page == "galaxyContent" ||
		page == "eventList" ||
		page == "ajaxChat" ||
		page == "notices" ||
		page == "repairlayer" ||
		page == "techtree" ||
		page == "phalanx" ||
		page == "shareReportOverlay" ||
		page == "jumpgatelayer" ||
		ajax == "1"
}

func (b *OGame) postPageContent(vals, payload url.Values) []byte {
	finalURL := b.serverURL + "/game/index.php?" + vals.Encode()
	req, err := http.NewRequest("POST", finalURL, strings.NewReader(payload.Encode()))
	if err != nil {
		b.error(err)
		return []byte{}
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")

	// Prevent redirect (301) https://stackoverflow.com/a/38150816/4196220
	b.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	defer func() {
		b.client.CheckRedirect = nil
	}()

	resp, err := b.client.Do(req)
	if err != nil {
		b.error(err)
		return []byte{}
	}
	defer resp.Body.Close()

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		defer reader.Close()
	default:
		reader = resp.Body
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		b.error(err)
		return []byte{}
	}

	return body
}

func (b *OGame) getPageContent(vals url.Values) []byte {
	if b.serverURL == "" {
		b.error("serverURL is empty")
		return []byte{}
	}

	finalURL := b.serverURL + "/game/index.php?" + vals.Encode()
	page := vals.Get("page")
	var pageHTMLBytes []byte

	b.withRetry(func() error {
		req, err := http.NewRequest("GET", finalURL, nil)
		if err != nil {
			return err
		}

		if IsAjaxPage(vals) {
			req.Header.Add("X-Requested-With", "XMLHttpRequest")
		}

		resp, err := b.client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 500 {
			return err
		}

		by, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		pageHTMLBytes = by

		if page != "logout" && !IsAjaxPage(vals) && !isLogged(pageHTMLBytes) {
			b.error("Err not logged on page : ", page)
			return ErrNotLogged
		}

		return nil
	})

	if page == "overview" {
		b.Player, _ = extractUserInfos(pageHTMLBytes, b.language)
		b.Planets = extractPlanets(pageHTMLBytes, b)
	} else if IsAjaxPage(vals) {
	} else {
		b.Planets = extractPlanets(pageHTMLBytes, b)
	}

	return pageHTMLBytes
}

type eventboxResp struct {
	Hostile  int
	Neutral  int
	Friendly int
}

func (b *OGame) withRetry(fn func() error) {
	retryInterval := 1
	retry := func(err error) {
		b.error(err.Error())
		time.Sleep(time.Duration(retryInterval) * time.Second)
		retryInterval *= 2
		if retryInterval > 60 {
			retryInterval = 60
		}
	}

	for {
		if err := fn(); err != nil {
			if err == ErrNotLogged {
				retry(err)
				if loginErr := b.login(); loginErr != nil {
					b.error(loginErr.Error())
					continue
				}
				continue
			} else {
				retry(err)
				continue
			}
		}
		break
	}
}

func (b *OGame) getPageJSON(vals url.Values, v interface{}) {
	b.withRetry(func() error {
		pageJSON := b.getPageContent(vals)
		if err := json.Unmarshal(pageJSON, v); err != nil {
			return ErrNotLogged
		}
		return nil
	})
}

// extract universe speed from html calculation
// pageHTML := b.getPageContent(url.Values{"page": {"techtree"}, "tab": {"2"}, "techID": {"1"}})
func extractUniverseSpeed(pageHTML []byte) int {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	spans := doc.Find("span.undermark")
	level := parseInt(spans.Eq(0).Text())
	val := parseInt(spans.Eq(1).Text())
	metalProduction := int(math.Floor(30 * float64(level) * math.Pow(1.1, float64(level))))
	universeSpeed := val / metalProduction
	return universeSpeed
}

func (b *OGame) getUniverseSpeed() int {
	return b.universeSpeed
}

func (b *OGame) getUniverseSpeedFleet() int {
	return b.universeSpeedFleet
}

func (b *OGame) isDonutGalaxy() bool {
	return b.donutGalaxy
}

func (b *OGame) isDonutSystem() bool {
	return b.donutSystem
}

func (b *OGame) fetchEventbox() (res eventboxResp) {
	b.getPageJSON(url.Values{"page": {"fetchEventbox"}}, &res)
	return
}

func (b *OGame) isUnderAttack() bool {
	return b.fetchEventbox().Hostile > 0
}

type resourcesResp struct {
	Metal struct {
		Resources struct {
			ActualFormat string
			Actual       int
			Max          int
			Production   float64
		}
		Tooltip string
		Class   string
	}
	Crystal struct {
		Resources struct {
			ActualFormat string
			Actual       int
			Max          int
			Production   float64
		}
		Tooltip string
		Class   string
	}
	Deuterium struct {
		Resources struct {
			ActualFormat string
			Actual       int
			Max          int
			Production   float64
		}
		Tooltip string
		Class   string
	}
	Energy struct {
		Resources struct {
			ActualFormat string
			Actual       int
		}
		Tooltip string
		Class   string
	}
	Darkmatter struct {
		Resources struct {
			ActualFormat string
			Actual       int
		}
		String  string
		Tooltip string
	}
	HonorScore int
}

func extractPlanetFromSelection(s *goquery.Selection, b *OGame) (Planet, error) {
	el, _ := s.Attr("id")
	id, err := strconv.Atoi(strings.TrimPrefix(el, "planet-"))
	if err != nil {
		return Planet{}, err
	}

	title, _ := s.Find("a.planetlink").Attr("title")
	root, err := html.Parse(strings.NewReader(title))
	if err != nil {
		return Planet{}, err
	}

	txt := goquery.NewDocumentFromNode(root).Text()
	planetInfosRgx := regexp.MustCompile(`([^\[]+) \[(\d+):(\d+):(\d+)]([\d.]+)km \((\d+)/(\d+)\)(?:de )?([-\d]+).+C\s*(?:bis|para|to|à|a) ([-\d]+).+C`)
	m := planetInfosRgx.FindStringSubmatch(txt)
	if len(m) < 10 {
		return Planet{}, errors.New("failed to parse planet infos: " + txt)
	}

	res := Planet{}
	res.ogame = b
	res.Img = s.Find("img.planetPic").AttrOr("src", "")
	res.ID = PlanetID(id)
	res.Name = m[1]
	res.Coordinate.Galaxy, _ = strconv.Atoi(m[2])
	res.Coordinate.System, _ = strconv.Atoi(m[3])
	res.Coordinate.Position, _ = strconv.Atoi(m[4])
	res.Diameter = parseInt(m[5])
	res.Fields.Built, _ = strconv.Atoi(m[6])
	res.Fields.Total, _ = strconv.Atoi(m[7])
	res.Temperature.Min, _ = strconv.Atoi(m[8])
	res.Temperature.Max, _ = strconv.Atoi(m[9])

	res.Moon, _ = extractMoonFromPlanetSelection(s, b)

	return res, nil
}

func extractMoonFromPlanetSelection(s *goquery.Selection, b *OGame) (*Moon, error) {
	moonLink := s.Find("a.moonlink")
	moon, err := extractMoonFromSelection(moonLink, b)
	if err != nil {
		return nil, err
	}
	return &moon, nil
}

func extractMoonFromSelection(moonLink *goquery.Selection, b *OGame) (Moon, error) {
	href, found := moonLink.Attr("href")
	if !found {
		return Moon{}, errors.New("no moon found")
	}
	m := regexp.MustCompile(`&cp=(\d+)`).FindStringSubmatch(href)
	id, _ := strconv.Atoi(m[1])
	title, _ := moonLink.Attr("title")
	root, err := html.Parse(strings.NewReader(title))
	if err != nil {
		return Moon{}, err
	}
	txt := goquery.NewDocumentFromNode(root).Text()
	moonInfosRgx := regexp.MustCompile(`([^\[]+) \[(\d+):(\d+):(\d+)]([\d.]+)km \((\d+)/(\d+)\)`)
	mm := moonInfosRgx.FindStringSubmatch(txt)
	if len(mm) < 8 {
		return Moon{}, errors.New("failed to parse moon infos: " + txt)
	}
	moon := Moon{}
	moon.ogame = b
	moon.ID = MoonID(id)
	moon.Name = mm[1]
	moon.Coordinate.Galaxy, _ = strconv.Atoi(mm[2])
	moon.Coordinate.System, _ = strconv.Atoi(mm[3])
	moon.Coordinate.Position, _ = strconv.Atoi(mm[4])
	moon.Diameter = parseInt(mm[5])
	moon.Fields.Built, _ = strconv.Atoi(mm[6])
	moon.Fields.Total, _ = strconv.Atoi(mm[7])
	moon.Img = moonLink.Find("img.icon-moon").AttrOr("src", "")
	return moon, nil
}

func extractPlanets(pageHTML []byte, b *OGame) []Planet {
	res := make([]Planet, 0)
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	doc.Find("div.smallplanet").Each(func(i int, s *goquery.Selection) {
		planet, err := extractPlanetFromSelection(s, b)
		if err != nil {
			b.error(err)
			return
		}
		res = append(res, planet)
	})
	return res
}

func extractPlanet(pageHTML []byte, planetID PlanetID, b *OGame) (Planet, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	s := doc.Find("div#planet-" + planetID.String())
	if len(s.Nodes) > 0 { // planet
		return extractPlanetFromSelection(s, b)
	}
	return Planet{}, errors.New("failed to find planetID")
}

func extractPlanetByCoord(pageHTML []byte, b *OGame, coord Coordinate) (Planet, error) {
	planets := extractPlanets(pageHTML, b)
	for _, planet := range planets {
		if planet.Coordinate.Equal(coord) {
			return planet, nil
		}
	}
	return Planet{}, errors.New("invalid planet coordinate")
}

func extractMoons(pageHTML []byte, b *OGame) []Moon {
	res := make([]Moon, 0)
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	doc.Find("a.moonlink").Each(func(i int, s *goquery.Selection) {
		moon, err := extractMoonFromSelection(s, b)
		if err != nil {
			return
		}
		res = append(res, moon)
	})
	return res
}

func extractMoon(pageHTML []byte, b *OGame, moonID MoonID) (Moon, error) {
	moons := extractMoons(pageHTML, b)
	for _, moon := range moons {
		if moon.ID == moonID {
			return moon, nil
		}
	}
	return Moon{}, errors.New("moon not found")
}

func extractMoonByCoord(pageHTML []byte, b *OGame, coord Coordinate) (Moon, error) {
	moons := extractMoons(pageHTML, b)
	for _, moon := range moons {
		if moon.Coordinate.Equal(coord) {
			return moon, nil
		}
	}
	return Moon{}, errors.New("invalid moon coordinate")
}

func (b *OGame) getPlanets() []Planet {
	pageHTML := b.getPageContent(url.Values{"page": {"overview"}})
	return extractPlanets(pageHTML, b)
}

func (b *OGame) getPlanet(planetID PlanetID) (Planet, error) {
	pageHTML := b.getPageContent(url.Values{"page": {"overview"}})
	return extractPlanet(pageHTML, planetID, b)
}

func (b *OGame) getPlanetByCoord(coord Coordinate) (Planet, error) {
	pageHTML := b.getPageContent(url.Values{"page": {"overview"}})
	return extractPlanetByCoord(pageHTML, b, coord)
}

func (b *OGame) getMoons() []Moon {
	pageHTML := b.getPageContent(url.Values{"page": {"overview"}})
	return extractMoons(pageHTML, b)
}

func (b *OGame) getMoon(moonID MoonID) (Moon, error) {
	pageHTML := b.getPageContent(url.Values{"page": {"overview"}})
	return extractMoon(pageHTML, b, moonID)
}

func (b *OGame) getMoonByCoord(coord Coordinate) (Moon, error) {
	pageHTML := b.getPageContent(url.Values{"page": {"overview"}})
	return extractMoonByCoord(pageHTML, b, coord)
}

func (b *OGame) serverVersion() string {
	return b.ogameVersion
}

func extractServerTime(pageHTML []byte) (time.Time, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return time.Time{}, err
	}
	txt := doc.Find("li.OGameClock").First().Text()
	serverTime, err := time.Parse("02.01.2006 15:04:05", txt)
	if err != nil {
		return time.Time{}, err
	}

	u1 := time.Now().UTC().Unix()
	u2 := serverTime.Unix()
	n := int(math.Round(float64(u2-u1)/15)) * 15

	serverTime = serverTime.Add(time.Duration(-n) * time.Second).In(time.FixedZone("OGT", n))

	return serverTime, nil
}

func (b *OGame) serverTime() time.Time {
	pageHTML := b.getPageContent(url.Values{"page": {"overview"}})
	serverTime, err := extractServerTime(pageHTML)
	if err != nil {
		b.error(err.Error())
	}
	return serverTime
}

func name2id(name string) ID {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	name, _, _ = transform.String(t, name)
	reg, _ := regexp.Compile("[^a-zA-Z]+")
	processedString := strings.ToLower(reg.ReplaceAllString(name, ""))
	nameMap := map[string]ID{
		// en
		"lightfighter":   LightFighterID,
		"heavyfighter":   HeavyFighterID,
		"cruiser":        CruiserID,
		"battleship":     BattleshipID,
		"battlecruiser":  BattlecruiserID,
		"bomber":         BomberID,
		"destroyer":      DestroyerID,
		"deathstar":      DeathstarID,
		"smallcargo":     SmallCargoID,
		"largecargo":     LargeCargoID,
		"colonyship":     ColonyShipID,
		"recycler":       RecyclerID,
		"espionageprobe": EspionageProbeID,
		"solarsatellite": SolarSatelliteID,

		// de
		"leichterjager":      LightFighterID,
		"schwererjager":      HeavyFighterID,
		"kreuzer":            CruiserID,
		"schlachtschiff":     BattleshipID,
		"schlachtkreuzer":    BattlecruiserID,
		"zerstorer":          DestroyerID,
		"todesstern":         DeathstarID,
		"kleinertransporter": SmallCargoID,
		"groertransporter":   LargeCargoID,
		"kolonieschiff":      ColonyShipID,
		"spionagesonde":      EspionageProbeID,
		"solarsatellit":      SolarSatelliteID,
		// "bomber":             BomberID,
		// "recycler":           RecyclerID,

		// es
		"cazadorligero":      LightFighterID,
		"cazadorpesado":      HeavyFighterID,
		"crucero":            CruiserID,
		"navedebatalla":      BattleshipID,
		"acorazado":          BattlecruiserID,
		"bombardero":         BomberID,
		"destructor":         DestroyerID,
		"estrelladelamuerte": DeathstarID,
		"navepequenadecarga": SmallCargoID,
		"navegrandedecarga":  LargeCargoID,
		"colonizador":        ColonyShipID,
		"reciclador":         RecyclerID,
		"sondadeespionaje":   EspionageProbeID,
		"satelitesolar":      SolarSatelliteID,

		// fr
		"chasseurleger":          LightFighterID,
		"chasseurlourd":          HeavyFighterID,
		"croiseur":               CruiserID,
		"vaisseaudebataille":     BattleshipID,
		"traqueur":               BattlecruiserID,
		"bombardier":             BomberID,
		"destructeur":            DestroyerID,
		"etoiledelamort":         DeathstarID,
		"petittransporteur":      SmallCargoID,
		"grandtransporteur":      LargeCargoID,
		"vaisseaudecolonisation": ColonyShipID,
		"recycleur":              RecyclerID,
		"sondedespionnage":       EspionageProbeID,
		"satellitesolaire":       SolarSatelliteID,

		// br
		"cacaligeiro":       LightFighterID,
		"cacapesado":        HeavyFighterID,
		"cruzador":          CruiserID,
		"navedebatalha":     BattleshipID,
		"interceptador":     BattlecruiserID,
		"bombardeiro":       BomberID,
		"destruidor":        DestroyerID,
		"estreladamorte":    DeathstarID,
		"cargueiropequeno":  SmallCargoID,
		"cargueirogrande":   LargeCargoID,
		"navecolonizadora":  ColonyShipID,
		"sondadeespionagem": EspionageProbeID,
		//"reciclador":        RecyclerID,
		//"satelitesolar":     SolarSatelliteID,

		// jp
		"軽戦闘機":      LightFighterID,
		"重戦闘機":      HeavyFighterID,
		"巡洋艦":       CruiserID,
		"バトルシップ":    BattleshipID,
		"大型戦艦":      BattlecruiserID,
		"爆撃機":       BomberID,
		"デストロイヤー":   DestroyerID,
		"デススター":     DeathstarID,
		"小型輸送機":     SmallCargoID,
		"大型輸送機":     LargeCargoID,
		"コロニーシップ":   ColonyShipID,
		"残骸回収船":     RecyclerID,
		"偵察機":       EspionageProbeID,
		"ソーラーサテライト": SolarSatelliteID,
	}
	return nameMap[processedString]
}

func extractUserInfos(pageHTML []byte, lang string) (UserInfos, error) {
	playerIDRgx := regexp.MustCompile(`playerId="(\d+)"`)
	playerNameRgx := regexp.MustCompile(`playerName="([^"]+)"`)
	txtContent := regexp.MustCompile(`textContent\[7]="([^"]+)"`)
	playerIDGroups := playerIDRgx.FindSubmatch(pageHTML)
	playerNameGroups := playerNameRgx.FindSubmatch(pageHTML)
	subHTMLGroups := txtContent.FindSubmatch(pageHTML)
	if len(playerIDGroups) < 2 {
		return UserInfos{}, errors.New("cannot find player id")
	}
	if len(playerNameGroups) < 2 {
		return UserInfos{}, errors.New("cannot find player name")
	}
	if len(subHTMLGroups) < 2 {
		return UserInfos{}, errors.New("cannot find sub html")
	}
	res := UserInfos{}
	res.PlayerID = toInt(playerIDGroups[1])
	res.PlayerName = string(playerNameGroups[1])
	html2 := subHTMLGroups[1]
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(html2))

	infosRgx := regexp.MustCompile(`([\d\\.]+) \(Place ([\d.]+) of ([\d.]+)\)`)
	switch lang {
	case "fr":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(Place ([\d.]+) sur ([\d.]+)\)`)
	case "de":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(Platz ([\d.]+) von ([\d.]+)\)`)
	case "es":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(Lugar ([\d.]+) de ([\d.]+)\)`)
	case "br":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(Posi\\u00e7\\u00e3o ([\d.]+) de ([\d.]+)\)`)
	case "jp":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(([\d.]+)\\u4eba\\u4e2d([\d.]+)\\u4f4d\)`)
	}
	// fr: 0 (Place 3.197 sur 3.348)
	// de: 0 (Platz 2.979 von 2.980)
	// jp: 0 (73人中72位)
	infos := infosRgx.FindStringSubmatch(doc.Text())
	if len(infos) < 4 {
		return UserInfos{}, errors.New("cannot find infos in sub html")
	}
	res.Points = parseInt(infos[1])
	res.Rank = parseInt(infos[2])
	res.Total = parseInt(infos[3])
	honourPointsRgx := regexp.MustCompile(`textContent\[9]="([^"]+)"`)
	honourPointsGroups := honourPointsRgx.FindSubmatch(pageHTML)
	if len(honourPointsGroups) < 2 {
		return UserInfos{}, errors.New("cannot find honour points")
	}
	res.HonourPoints = parseInt(string(honourPointsGroups[1]))
	return res, nil
}

func (b *OGame) getUserInfos() UserInfos {
	pageHTML := b.getPageContent(url.Values{"page": {"overview"}})
	userInfos, err := extractUserInfos(pageHTML, b.language)
	if err != nil {
		b.error(err)
	}
	return userInfos
}

func (b *OGame) sendMessage(playerID int, message string) error {
	finalURL := b.serverURL + "/game/index.php?page=ajaxChat"
	payload := url.Values{
		"playerId": {strconv.Itoa(playerID)},
		"text":     {message + "\n"},
		"mode":     {"1"},
		"ajax":     {"1"},
	}
	req, err := http.NewRequest("POST", finalURL, strings.NewReader(payload.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	resp, err := b.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bobyBytes, _ := ioutil.ReadAll(resp.Body)
	if strings.Contains(string(bobyBytes), "INVALID_PARAMETERS") {
		return errors.New("invalid parameters")
	}
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(string(bobyBytes)))
	if doc.Find("title").Text() == "OGame Lobby" {
		return ErrNotLogged
	}
	return nil
}

func extractCoord(v string) (coord Coordinate) {
	coordRgx := regexp.MustCompile(`\[(\d+):(\d+):(\d+)]`)
	m := coordRgx.FindStringSubmatch(v)
	if len(m) == 4 {
		coord.Galaxy, _ = strconv.Atoi(m[1])
		coord.System, _ = strconv.Atoi(m[2])
		coord.Position, _ = strconv.Atoi(m[3])
	}
	return
}

func extractFleets(pageHTML []byte) (res []Fleet) {
	res = make([]Fleet, 0)
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	doc.Find("div.fleetDetails").Each(func(i int, s *goquery.Selection) {
		originText := s.Find("span.originCoords a").Text()
		origin := extractCoord(originText)

		destText := s.Find("span.destinationCoords a").Text()
		dest := extractCoord(destText)

		idStr, _ := s.Find("span.reversal").Attr("ref")
		id, _ := strconv.Atoi(idStr)

		missionTypeRaw, _ := s.Attr("data-mission-type")
		returnFlightRaw, _ := s.Attr("data-return-flight")
		arrivalTimeRaw, _ := s.Attr("data-arrival-time")
		missionType, _ := strconv.Atoi(missionTypeRaw)
		returnFlight, _ := strconv.ParseBool(returnFlightRaw)
		arrivalTime, _ := strconv.Atoi(arrivalTimeRaw)
		ogameTimestamp, _ := strconv.Atoi(doc.Find("meta[name=ogame-timestamp]").AttrOr("content", "0"))
		secs := arrivalTime - ogameTimestamp
		if secs < 0 {
			secs = 0
		}

		trs := s.Find("table.fleetinfo tr")
		shipment := Resources{}
		shipment.Metal = parseInt(trs.Eq(trs.Size() - 3).Find("td").Eq(1).Text())
		shipment.Crystal = parseInt(trs.Eq(trs.Size() - 2).Find("td").Eq(1).Text())
		shipment.Deuterium = parseInt(trs.Eq(trs.Size() - 1).Find("td").Eq(1).Text())

		fleet := Fleet{}
		fleet.ID = FleetID(id)
		fleet.Origin = origin
		fleet.Destination = dest
		fleet.Mission = MissionID(missionType)
		fleet.ReturnFlight = returnFlight
		fleet.Resources = shipment
		fleet.ArriveIn = secs

		for i := 1; i < trs.Size()-5; i++ {
			tds := trs.Eq(i).Find("td")
			name := strings.ToLower(strings.Trim(strings.TrimSpace(tds.Eq(0).Text()), ":"))
			qty := parseInt(tds.Eq(1).Text())
			shipID := name2id(name)
			fleet.Ships.Set(shipID, qty)
		}

		res = append(res, fleet)
	})
	return
}

func (b *OGame) getFleets() []Fleet {
	pageHTML := b.getPageContent(url.Values{"page": {"movement"}})
	return extractFleets(pageHTML)
}

func (b *OGame) cancelFleet(fleetID FleetID) error {
	b.getPageContent(url.Values{"page": {"movement"}, "return": {fleetID.String()}})
	return nil
}

// Returns the distance between two galaxy
func galaxyDistance(galaxy1, galaxy2, universeSize int, donutGalaxy bool) (distance int) {
	if !donutGalaxy {
		return int(20000 * math.Abs(float64(galaxy2-galaxy1)))
	}
	if galaxy1 > galaxy2 {
		galaxy1, galaxy2 = galaxy2, galaxy1
	}
	val := math.Min(float64(galaxy2-galaxy1), float64((galaxy1+universeSize)-galaxy2))
	return int(20000 * val)
}

func systemDistance(system1, system2 int, donutSystem bool) (distance int) {
	if !donutSystem {
		return int(math.Abs(float64(system2 - system1)))
	}
	systemSize := 499
	if system1 > system2 {
		system1, system2 = system2, system1
	}
	return int(math.Min(float64(system2-system1), float64((system1+systemSize)-system2)))
}

// Returns the distance between two systems
func flightSystemDistance(system1, system2 int, donutSystem bool) (distance int) {
	return 2700 + 95*systemDistance(system1, system2, donutSystem)
}

// Returns the distance between two planets
func planetDistance(planet1, planet2 int) (distance int) {
	return int(1000 + 5*math.Abs(float64(planet2-planet1)))
}

func distance(c1, c2 Coordinate, universeSize int, donutGalaxy, donutSystem bool) (distance int) {
	if c1.Galaxy != c2.Galaxy {
		return galaxyDistance(c1.Galaxy, c2.Galaxy, universeSize, donutGalaxy)
	}
	if c1.System != c2.System {
		return flightSystemDistance(c1.System, c2.System, donutSystem)
	}
	return planetDistance(c1.Position, c2.Position)
}

func findSlowestSpeed(ships ShipsInfos, techs Researches) int {
	minSpeed := math.MaxInt32
	for _, ship := range Ships {
		shipSpeed := ship.GetSpeed(techs)
		if ships.ByID(ship.GetID()) > 0 && shipSpeed < minSpeed {
			minSpeed = shipSpeed
		}
	}
	return minSpeed
}

func calcFuel(ships ShipsInfos, dist int, speed float64) (fuel int) {
	tmpFn := func(baseFuel int) int {
		return int(1 + math.Round(((float64(baseFuel)*float64(dist))/35000)*math.Pow(speed+1, 2)))
	}
	for _, ship := range Ships {
		nbr := ships.ByID(ship.GetID())
		if nbr > 0 {
			fuel += tmpFn(ship.GetFuelConsumption()) * nbr
		}
	}
	return
}

func calcFlightTime(origin, destination Coordinate, universeSize int, donutGalaxy, donutSystem bool, speed float64,
	universeSpeedFleet int, ships ShipsInfos, techs Researches) (secs, fuel int) {
	s := speed
	v := float64(findSlowestSpeed(ships, techs))
	a := float64(universeSpeedFleet)
	d := float64(distance(origin, destination, universeSize, donutGalaxy, donutSystem))
	secs = int(math.Round(((10 + (3500 / s)) * math.Sqrt((10*d)/v)) / a))
	fuel = calcFuel(ships, int(d), speed)
	return
}

func extractOgameTimestamp(pageHTML []byte) int {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	ogameTimestamp, _ := strconv.Atoi(doc.Find("meta[name=ogame-timestamp]").AttrOr("content", "0"))
	return ogameTimestamp
}

func extractResources(pageHTML []byte) Resources {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	res := Resources{}
	res.Metal = parseInt(doc.Find("span#resources_metal").Text())
	res.Crystal = parseInt(doc.Find("span#resources_crystal").Text())
	res.Deuterium = parseInt(doc.Find("span#resources_deuterium").Text())
	res.Energy = parseInt(doc.Find("span#resources_energy").Text())
	res.Darkmatter = parseInt(doc.Find("span#resources_darkmatter").Text())
	return res
}

func extractPhalanx(pageHTML []byte, ogameTimestamp int) ([]Fleet, error) {
	res := make([]Fleet, 0)
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	eventFleet := doc.Find("div.eventFleet")
	if eventFleet.Size() == 0 {
		txt := doc.Find("div#phalanxEventContent").Text()
		// TODO: 'fleet' and 'deuterium' won't work in other languages
		if strings.Contains(txt, "fleet") {
			return res, nil
		} else if strings.Contains(txt, "deuterium") {
			return res, errors.New(strings.TrimSpace(txt))
		}
		return res, errors.New(txt)
	}
	eventFleet.Each(func(i int, s *goquery.Selection) {
		mission, _ := strconv.Atoi(s.AttrOr("data-mission-type", "0"))
		returning, _ := strconv.ParseBool(s.AttrOr("data-return-flight", "false"))
		arrivalTime, _ := strconv.Atoi(s.AttrOr("data-arrival-time", "0"))
		arriveIn := arrivalTime - ogameTimestamp
		if arriveIn < 0 {
			arriveIn = 0
		}
		originTxt := s.Find("li.coordsOrigin a").Text()
		destTxt := s.Find("li.destCoords a").Text()

		fleet := Fleet{}

		if movement, exists := s.Find("li.detailsFleet span").Attr("title"); exists {
			root, err := html.Parse(strings.NewReader(movement))
			if err != nil {
				return
			}
			doc2 := goquery.NewDocumentFromNode(root)
			doc2.Find("tr").Each(func(i int, s *goquery.Selection) {
				if i == 0 {
					return
				}
				name := s.Find("td").Eq(0).Text()
				nbr := parseInt(s.Find("td").Eq(1).Text())
				if name != "" && nbr > 0 {
					fleet.Ships.Set(name2id(name), nbr)
				}
			})
		}

		fleet.Mission = MissionID(mission)
		fleet.ReturnFlight = returning
		fleet.ArriveIn = arriveIn
		fleet.Origin = extractCoord(originTxt)
		fleet.Destination = extractCoord(destTxt)
		res = append(res, fleet)
	})
	return res, nil
}

// getPhalanx makes 3 calls to ogame server (2 validation, 1 scan)
func (b *OGame) getPhalanx(moonID MoonID, coord Coordinate) ([]Fleet, error) {
	// Get moon facilities html page (first call to ogame server)
	moonFacilitiesHTML := b.getPageContent(url.Values{"page": {"station"}, "cp": {strconv.Itoa(int(moonID))}})

	// Extract bunch of infos from the html
	moon, err := extractMoon(moonFacilitiesHTML, b, moonID)
	if err != nil {
		return make([]Fleet, 0), errors.New("moon not found")
	}
	resources := extractResources(moonFacilitiesHTML)
	moonFacilities, _ := extractFacilities(moonFacilitiesHTML)
	ogameTimestamp := extractOgameTimestamp(moonFacilitiesHTML)
	phalanxLvl := moonFacilities.SensorPhalanx

	// Ensure we have the resources to scan the planet
	if resources.Deuterium < SensorPhalanx.ScanConsumption() {
		return make([]Fleet, 0), errors.New("not enough deuterium")
	}

	// Verify that coordinate is in phalanx range
	phalanxRange := SensorPhalanx.GetRange(phalanxLvl)
	if moon.Coordinate.Galaxy != coord.Galaxy ||
		systemDistance(moon.Coordinate.System, coord.System, b.donutSystem) > phalanxRange {
		return make([]Fleet, 0), errors.New("coordinate not in phalanx range")
	}

	// Get galaxy planets information, verify coordinate is valid planet (second call to ogame server)
	planetInfos, _ := b.galaxyInfos(coord.Galaxy, coord.System)
	target := planetInfos.Position(coord.Position)
	if target == nil {
		return make([]Fleet, 0), errors.New("invalid planet coordinate")
	}
	// Ensure you are not scanning your own planet
	if target.Player.ID == b.Player.PlayerID {
		return make([]Fleet, 0), errors.New("cannot scan own planet")
	}

	// Run the phalanx scan (third call to ogame server)
	finalURL := fmt.Sprintf(b.serverURL+"/game/index.php?page=phalanx&galaxy=%d&system=%d&position=%d&ajax=1",
		coord.Galaxy, coord.System, coord.Position)
	req, err := http.NewRequest("GET", finalURL, nil)
	if err != nil {
		b.error(err.Error())
		return []Fleet{}, err
	}
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	resp, err := b.client.Do(req)
	if err != nil {
		b.error(err.Error())
		return []Fleet{}, err
	}
	defer resp.Body.Close()
	pageHTML, _ := ioutil.ReadAll(resp.Body)

	return extractPhalanx(pageHTML, ogameTimestamp)
}

func extractJumpGate(pageHTML []byte) (ShipsInfos, string, []MoonID, int) {
	m := regexp.MustCompile(`\$\("#cooldown"\), (\d+),`).FindSubmatch(pageHTML)
	ships := ShipsInfos{}
	var destinations []MoonID
	if len(m) > 0 {
		waitTime := toInt(m[1])
		return ships, "", destinations, waitTime
	}
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	for _, s := range Ships {
		ships.Set(s.GetID(), parseInt(doc.Find("input#ship_"+strconv.Itoa(int(s.GetID()))).AttrOr("rel", "0")))
	}
	token := doc.Find("input[name=token]").AttrOr("value", "")

	doc.Find("select[name=zm] option").Each(func(i int, s *goquery.Selection) {
		moonID := parseInt(s.AttrOr("value", "0"))
		if moonID > 0 {
			destinations = append(destinations, MoonID(moonID))
		}
	})

	return ships, token, destinations, 0
}

func moonIDInSlice(needle MoonID, haystack []MoonID) bool {
	for _, element := range haystack {
		if needle == element {
			return true
		}
	}
	return false
}

func (b *OGame) executeJumpGate(originMoonID, destMoonID MoonID, ships ShipsInfos) error {
	pageHTML := b.getPageContent(url.Values{"page": {"jumpgatelayer"}, "cp": {strconv.Itoa(int(originMoonID))}})
	availShips, token, dests, wait := extractJumpGate(pageHTML)
	if wait > 0 {
		return fmt.Errorf("jump gate is in recharge mode for %d seconds", wait)
	}

	// Validate destination moon id
	if !moonIDInSlice(destMoonID, dests) {
		return errors.New("destination moon id invalid")
	}

	finalURL := b.serverURL + "/game/index.php?page=jumpgate_execute"
	payload := url.Values{
		"token": {token},
		"zm":    {strconv.Itoa(int(destMoonID))},
	}

	// Add ships to payload
	for _, s := range Ships {
		// Get the min between what is available and what we want
		nbr := int(math.Min(float64(ships.ByID(s.GetID())), float64(availShips.ByID(s.GetID()))))
		if nbr > 0 {
			payload.Add("ship_"+strconv.Itoa(int(s.GetID())), strconv.Itoa(nbr))
		}
	}

	req, err := http.NewRequest("POST", finalURL, strings.NewReader(payload.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	resp, err := b.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to execute jump gate : %s", resp.Status)
	}
	return nil
}

func extractAttacks(pageHTML []byte) []AttackEvent {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	attacks := make([]AttackEvent, 0)
	tmp := func(i int, s *goquery.Selection) {
		classes, _ := s.Attr("class")
		if strings.Contains(classes, "partnerInfo") {
			return
		}
		isHostile := s.Find("td.countDown.hostile").Size() > 0
		if !isHostile {
			return
		}
		missionTypeStr, _ := s.Attr("data-mission-type")
		arrivalTimeStr, _ := s.Attr("data-arrival-time")
		missionTypeInt, _ := strconv.Atoi(missionTypeStr)
		arrivalTimeInt, _ := strconv.Atoi(arrivalTimeStr)
		missionType := MissionID(missionTypeInt)
		if missionType != Attack && missionType != GroupedAttack &&
			missionType != MissileAttack && missionType != Spy {
			return
		}
		attack := AttackEvent{}
		attack.MissionType = missionType
		if missionType == Attack || missionType == MissileAttack || missionType == Spy {
			coordsOrigin := strings.TrimSpace(s.Find("td.coordsOrigin").Text())
			attack.Origin = extractCoord(coordsOrigin)
			attackerIDStr, _ := s.Find("a.sendMail").Attr("data-playerid")
			attack.AttackerID, _ = strconv.Atoi(attackerIDStr)
		}
		if missionType == MissileAttack {
			attack.Missiles = parseInt(s.Find("td.detailsFleet span").First().Text())
		}

		// Get ships infos if available
		if movement, exists := s.Find("td.icon_movement span").Attr("title"); exists {
			root, err := html.Parse(strings.NewReader(movement))
			if err != nil {
				return
			}
			attack.Ships = new(ShipsInfos)
			q := goquery.NewDocumentFromNode(root)
			q.Find("tr").Each(func(i int, s *goquery.Selection) {
				name := s.Find("td").Eq(0).Text()
				nbr := parseInt(s.Find("td").Eq(1).Text())
				if name != "" && nbr > 0 {
					attack.Ships.Set(name2id(name), nbr)
				}
			})
		}

		if s.Find("td.destFleet figure.planet").Size() == 1 {
			attack.DestinationType = PlanetDest
		}
		if s.Find("td.destFleet figure.moon").Size() == 1 {
			attack.DestinationType = MoonDest
		}

		destCoords := strings.TrimSpace(s.Find("td.destCoords").Text())
		attack.Destination = extractCoord(destCoords)

		attack.ArrivalTime = time.Unix(int64(arrivalTimeInt), 0)

		attacks = append(attacks, attack)
	}
	doc.Find("tr.eventFleet").Each(tmp)
	doc.Find("tr.allianceAttack").Each(tmp)

	return attacks
}

func (b *OGame) getAttacks() []AttackEvent {
	finalURL := b.serverURL + "/game/index.php?page=eventList&ajax=1"
	req, err := http.NewRequest("GET", finalURL, nil)
	if err != nil {
		b.error(err.Error())
		return []AttackEvent{}
	}
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	resp, err := b.client.Do(req)
	if err != nil {
		b.error(err.Error())
		return []AttackEvent{}
	}
	defer resp.Body.Close()
	pageHTML, _ := ioutil.ReadAll(resp.Body)

	return extractAttacks(pageHTML)
}

func extractGalaxyInfos(pageHTML []byte, botPlayerName string, botPlayerID, botPlayerRank int) (SystemInfos, error) {
	prefixedNumRgx := regexp.MustCompile(`.*: ([\d.]+)`)

	extractActivity := func(activityDiv *goquery.Selection) int {
		activity := 0
		if activityDiv != nil {
			activityDivClass := activityDiv.AttrOr("class", "")
			if strings.Contains(activityDivClass, "minute15") {
				activity = 15
			} else if strings.Contains(activityDivClass, "showMinutes") {
				activity, _ = strconv.Atoi(strings.TrimSpace(activityDiv.Text()))
			}
		}
		return activity
	}

	var tmp struct {
		Galaxy string
	}
	json.Unmarshal(pageHTML, &tmp)
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(tmp.Galaxy))
	var res SystemInfos
	res.galaxy = parseInt(doc.Find("table").AttrOr("data-galaxy", "0"))
	res.system = parseInt(doc.Find("table").AttrOr("data-system", "0"))
	doc.Find("tr.row").Each(func(i int, s *goquery.Selection) {
		classes, _ := s.Attr("class")
		if !strings.Contains(classes, "empty_filter") {
			position := s.Find("td.position").Text()

			tooltips := s.Find("div.htmlTooltip")
			planetTooltip := tooltips.First()
			planetName := planetTooltip.Find("h1").Find("span").Text()
			planetImg, _ := planetTooltip.Find("img").Attr("src")
			coordsRaw := planetTooltip.Find("span#pos-planet").Text()

			metalTxt := s.Find("div#debris" + position + " ul.ListLinks li").First().Text()
			crystalTxt := s.Find("div#debris" + position + " ul.ListLinks li").Eq(1).Text()
			recyclersTxt := s.Find("div#debris" + position + " ul.ListLinks li").Eq(2).Text()

			planetInfos := new(PlanetInfos)
			planetInfos.ID, _ = strconv.Atoi(s.Find("td.colonized").AttrOr("data-planet-id", ""))

			moonID, _ := strconv.Atoi(s.Find("td.moon").AttrOr("data-moon-id", ""))
			moonSize, _ := strconv.Atoi(strings.Split(s.Find("td.moon span#moonsize").Text(), " ")[0])
			if moonID > 0 {
				planetInfos.Moon = new(MoonInfos)
				planetInfos.Moon.ID = moonID
				planetInfos.Moon.Diameter = moonSize
				planetInfos.Moon.Activity = extractActivity(s.Find("td.moon div.activity"))
			}

			allianceSpan := s.Find("span.allytagwrapper")
			if allianceSpan.Size() > 0 {
				longID, _ := allianceSpan.Attr("rel")
				planetInfos.Alliance = new(AllianceInfos)
				planetInfos.Alliance.Name = allianceSpan.Find("h1").Text()
				planetInfos.Alliance.ID, _ = strconv.Atoi(strings.TrimPrefix(longID, "alliance"))
				planetInfos.Alliance.Rank, _ = strconv.Atoi(allianceSpan.Find("ul.ListLinks li").First().Find("a").Text())
				planetInfos.Alliance.Member = parseInt(prefixedNumRgx.FindStringSubmatch(allianceSpan.Find("ul.ListLinks li").Eq(1).Text())[1])
			}

			if len(prefixedNumRgx.FindStringSubmatch(metalTxt)) > 0 {
				planetInfos.Debris.Metal = parseInt(prefixedNumRgx.FindStringSubmatch(metalTxt)[1])
				planetInfos.Debris.Crystal = parseInt(prefixedNumRgx.FindStringSubmatch(crystalTxt)[1])
				planetInfos.Debris.RecyclersNeeded = parseInt(prefixedNumRgx.FindStringSubmatch(recyclersTxt)[1])
			}

			planetInfos.Activity = extractActivity(s.Find("td:not(.moon) div.activity"))
			planetInfos.Name = planetName
			planetInfos.Img = planetImg
			planetInfos.Inactive = strings.Contains(classes, "inactive_filter")
			planetInfos.StrongPlayer = strings.Contains(classes, "strong_filter")
			planetInfos.Newbie = strings.Contains(classes, "newbie_filter")
			planetInfos.Vacation = strings.Contains(classes, "vacation_filter")
			planetInfos.HonorableTarget = s.Find("span.status_abbr_honorableTarget").Size() > 0
			planetInfos.Administrator = s.Find("span.status_abbr_admin").Size() > 0
			planetInfos.Banned = s.Find("td.playername a span.status_abbr_banned").Size() > 0
			planetInfos.Coordinate = extractCoord(coordsRaw)

			var playerID int
			var playerName string
			var playerRank int
			if len(tooltips.Nodes) > 1 {
				tooltips.Each(func(i int, s *goquery.Selection) {
					idAttr, _ := s.Attr("id")
					if strings.HasPrefix(idAttr, "player") {
						playerID, _ = strconv.Atoi(regexp.MustCompile(`player(\d+)`).FindStringSubmatch(idAttr)[1])
						playerName = s.Find("h1").Find("span").Text()
						playerRank, _ = strconv.Atoi(s.Find("li.rank").Find("a").Text())
					}
				})
			} else {
				playerName = strings.TrimSpace(s.Find("td.playername").Find("span").Text())
			}

			if playerID == 0 {
				playerID = botPlayerID
				playerName = botPlayerName
				playerRank = botPlayerRank
			}

			planetInfos.Player.ID = playerID
			planetInfos.Player.Name = playerName
			planetInfos.Player.Rank = playerRank

			res.planets[i] = planetInfos
		}
	})
	return res, nil
}

func (b *OGame) galaxyInfos(galaxy, system int) (SystemInfos, error) {
	if galaxy < 0 || galaxy > b.server.Settings.UniverseSize {
		return SystemInfos{}, fmt.Errorf("galaxy must be within [0, %d]", b.server.Settings.UniverseSize)
	}
	if system < 0 || system > 499 {
		return SystemInfos{}, errors.New("system must be within [0, 499]")
	}
	finalURL := b.serverURL + "/game/index.php?page=galaxyContent&ajax=1"
	payload := url.Values{
		"galaxy": {strconv.Itoa(galaxy)},
		"system": {strconv.Itoa(system)},
	}
	req, err := http.NewRequest("POST", finalURL, strings.NewReader(payload.Encode()))
	if err != nil {
		return SystemInfos{}, err
	}
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	resp, err := b.client.Do(req)
	if err != nil {
		return SystemInfos{}, err
	}
	defer resp.Body.Close()
	pageHTML, _ := ioutil.ReadAll(resp.Body)
	return extractGalaxyInfos(pageHTML, b.Player.PlayerName, b.Player.PlayerID, b.Player.Rank)
}

func (b *OGame) getResourceSettings(planetID PlanetID) (ResourceSettings, error) {
	pageHTML := b.getPageContent(url.Values{"page": {"resourceSettings"}, "cp": {planetID.String()}})
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	bodyID, _ := doc.Find("body").Attr("id")
	if bodyID == "overview" {
		return ResourceSettings{}, ErrInvalidPlanetID
	}
	vals := make([]int, 0)
	doc.Find("option").Each(func(i int, s *goquery.Selection) {
		_, selectedExists := s.Attr("selected")
		if selectedExists {
			a, _ := s.Attr("value")
			val, _ := strconv.Atoi(a)
			vals = append(vals, val)
		}
	})
	if len(vals) != 6 {
		return ResourceSettings{}, errors.New("failed to find all resource settings")
	}

	res := ResourceSettings{}
	res.MetalMine = vals[0]
	res.CrystalMine = vals[1]
	res.DeuteriumSynthesizer = vals[2]
	res.SolarPlant = vals[3]
	res.FusionReactor = vals[4]
	res.SolarSatellite = vals[5]

	return res, nil
}

func (b *OGame) setResourceSettings(planetID PlanetID, settings ResourceSettings) error {
	pageHTML := b.getPageContent(url.Values{"page": {"resourceSettings"}, "cp": {planetID.String()}})
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	bodyID, _ := doc.Find("body").Attr("id")
	if bodyID == "overview" {
		return ErrInvalidPlanetID
	}
	token, exists := doc.Find("form input[name=token]").Attr("value")
	if !exists {
		return errors.New("unable to find token")
	}
	payload := url.Values{
		"saveSettings": {"1"},
		"token":        {token},
		"last1":        {strconv.Itoa(settings.MetalMine)},
		"last2":        {strconv.Itoa(settings.CrystalMine)},
		"last3":        {strconv.Itoa(settings.DeuteriumSynthesizer)},
		"last4":        {strconv.Itoa(settings.SolarPlant)},
		"last12":       {strconv.Itoa(settings.FusionReactor)},
		"last212":      {strconv.Itoa(settings.SolarSatellite)},
	}
	url2 := b.serverURL + "/game/index.php?page=resourceSettings"
	resp, err := b.client.PostForm(url2, payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func getNbr(doc *goquery.Document, name string) int {
	div := doc.Find("div." + name)
	level := div.Find("span.level")
	level.Children().Remove()
	return parseInt(level.Text())
}

func extractResourcesBuildings(pageHTML []byte) (ResourcesBuildings, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	doc.Find("span.textlabel").Remove()
	bodyID, _ := doc.Find("body").Attr("id")
	if bodyID == "overview" {
		return ResourcesBuildings{}, ErrInvalidPlanetID
	}
	res := ResourcesBuildings{}
	res.MetalMine = getNbr(doc, "supply1")
	res.CrystalMine = getNbr(doc, "supply2")
	res.DeuteriumSynthesizer = getNbr(doc, "supply3")
	res.SolarPlant = getNbr(doc, "supply4")
	res.FusionReactor = getNbr(doc, "supply12")
	res.SolarSatellite = getNbr(doc, "supply212")
	res.MetalStorage = getNbr(doc, "supply22")
	res.CrystalStorage = getNbr(doc, "supply23")
	res.DeuteriumTank = getNbr(doc, "supply24")
	return res, nil
}

func extractDefense(pageHTML []byte) (DefensesInfos, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	bodyID, _ := doc.Find("body").Attr("id")
	if bodyID == "overview" {
		return DefensesInfos{}, ErrInvalidPlanetID
	}
	doc.Find("span.textlabel").Remove()
	res := DefensesInfos{}
	res.RocketLauncher = getNbr(doc, "defense401")
	res.LightLaser = getNbr(doc, "defense402")
	res.HeavyLaser = getNbr(doc, "defense403")
	res.GaussCannon = getNbr(doc, "defense404")
	res.IonCannon = getNbr(doc, "defense405")
	res.PlasmaTurret = getNbr(doc, "defense406")
	res.SmallShieldDome = getNbr(doc, "defense407")
	res.LargeShieldDome = getNbr(doc, "defense408")
	res.AntiBallisticMissiles = getNbr(doc, "defense502")
	res.InterplanetaryMissiles = getNbr(doc, "defense503")

	return res, nil
}

func extractShips(pageHTML []byte) (ShipsInfos, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	doc.Find("span.textlabel").Remove()
	bodyID, _ := doc.Find("body").Attr("id")
	if bodyID == "overview" {
		return ShipsInfos{}, ErrInvalidPlanetID
	}
	res := ShipsInfos{}
	res.LightFighter = getNbr(doc, "military204")
	res.HeavyFighter = getNbr(doc, "military205")
	res.Cruiser = getNbr(doc, "military206")
	res.Battleship = getNbr(doc, "military207")
	res.Battlecruiser = getNbr(doc, "military215")
	res.Bomber = getNbr(doc, "military211")
	res.Destroyer = getNbr(doc, "military213")
	res.Deathstar = getNbr(doc, "military214")
	res.SmallCargo = getNbr(doc, "civil202")
	res.LargeCargo = getNbr(doc, "civil203")
	res.ColonyShip = getNbr(doc, "civil208")
	res.Recycler = getNbr(doc, "civil209")
	res.EspionageProbe = getNbr(doc, "civil210")
	res.SolarSatellite = getNbr(doc, "civil212")

	return res, nil
}

func extractFacilities(pageHTML []byte) (Facilities, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	doc.Find("span.textlabel").Remove()
	bodyID, _ := doc.Find("body").Attr("id")
	if bodyID == "overview" {
		return Facilities{}, ErrInvalidPlanetID
	}
	res := Facilities{}
	res.RoboticsFactory = getNbr(doc, "station14")
	res.Shipyard = getNbr(doc, "station21")
	res.ResearchLab = getNbr(doc, "station31")
	res.AllianceDepot = getNbr(doc, "station34")
	res.MissileSilo = getNbr(doc, "station44")
	res.NaniteFactory = getNbr(doc, "station15")
	res.Terraformer = getNbr(doc, "station33")
	res.SpaceDock = getNbr(doc, "station36")
	res.LunarBase = getNbr(doc, "station41")
	res.SensorPhalanx = getNbr(doc, "station42")
	res.JumpGate = getNbr(doc, "station43")
	return res, nil
}

func extractResearch(pageHTML []byte) Researches {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	doc.Find("span.textlabel").Remove()
	res := Researches{}
	res.EnergyTechnology = getNbr(doc, "research113")
	res.LaserTechnology = getNbr(doc, "research120")
	res.IonTechnology = getNbr(doc, "research121")
	res.HyperspaceTechnology = getNbr(doc, "research114")
	res.PlasmaTechnology = getNbr(doc, "research122")
	res.CombustionDrive = getNbr(doc, "research115")
	res.ImpulseDrive = getNbr(doc, "research117")
	res.HyperspaceDrive = getNbr(doc, "research118")
	res.EspionageTechnology = getNbr(doc, "research106")
	res.ComputerTechnology = getNbr(doc, "research108")
	res.Astrophysics = getNbr(doc, "research124")
	res.IntergalacticResearchNetwork = getNbr(doc, "research123")
	res.GravitonTechnology = getNbr(doc, "research199")
	res.WeaponsTechnology = getNbr(doc, "research109")
	res.ShieldingTechnology = getNbr(doc, "research110")
	res.ArmourTechnology = getNbr(doc, "research111")

	return res
}

func (b *OGame) getResearch() Researches {
	pageHTML := b.getPageContent(url.Values{"page": {"research"}})
	return extractResearch(pageHTML)
}

func (b *OGame) getResourcesBuildings(planetID PlanetID) (ResourcesBuildings, error) {
	pageHTML := b.getPageContent(url.Values{"page": {"resources"}, "cp": {planetID.String()}})
	return extractResourcesBuildings(pageHTML)
}

func (b *OGame) getDefense(celestialID CelestialID) (DefensesInfos, error) {
	pageHTML := b.getPageContent(url.Values{"page": {"defense"}, "cp": {strconv.Itoa(int(celestialID))}})
	return extractDefense(pageHTML)
}

func (b *OGame) getShips(celestialID CelestialID) (ShipsInfos, error) {
	pageHTML := b.getPageContent(url.Values{"page": {"shipyard"}, "cp": {strconv.Itoa(int(celestialID))}})
	return extractShips(pageHTML)
}

func (b *OGame) getFacilities(celestialID CelestialID) (Facilities, error) {
	pageHTML := b.getPageContent(url.Values{"page": {"station"}, "cp": {strconv.Itoa(int(celestialID))}})
	return extractFacilities(pageHTML)
}

func extractProduction(pageHTML []byte) ([]Quantifiable, error) {
	res := make([]Quantifiable, 0)
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	active := doc.Find("table.construction")
	href, _ := active.Find("td a").Attr("href")
	m := regexp.MustCompile(`openTech=(\d+)`).FindStringSubmatch(href)
	if len(m) == 0 {
		return []Quantifiable{}, nil
	}
	idInt, _ := strconv.Atoi(m[1])
	activeID := ID(idInt)
	activeNbr, _ := strconv.Atoi(active.Find("div.shipSumCount").Text())
	res = append(res, Quantifiable{ID: activeID, Nbr: activeNbr})
	doc.Find("div#pqueue ul li").Each(func(i int, s *goquery.Selection) {
		link := s.Find("a")
		itemIDstr, exists := link.Attr("ref")
		if !exists {
			href := link.AttrOr("href", "")
			m := regexp.MustCompile(`openTech=(\d+)`).FindStringSubmatch(href)
			if len(m) > 0 {
				itemIDstr = m[1]
			}
		}
		itemID, _ := strconv.Atoi(itemIDstr)
		itemNbr := parseInt(s.Find("span.number").Text())
		res = append(res, Quantifiable{ID: ID(itemID), Nbr: itemNbr})
	})
	return res, nil
}

func (b *OGame) getProduction(celestialID CelestialID) ([]Quantifiable, error) {
	pageHTML := b.getPageContent(url.Values{"page": {"shipyard"}, "cp": {strconv.Itoa(int(celestialID))}})
	return extractProduction(pageHTML)
}

func (b *OGame) build(celestialID CelestialID, id ID, nbr int) error {
	var page string
	if id.IsDefense() {
		page = "defense"
	} else if id.IsShip() {
		page = "shipyard"
	} else if id.IsBuilding() {
		page = "resources"
	} else if id.IsTech() {
		page = "research"
	} else {
		return errors.New("invalid id " + id.String())
	}
	payload := url.Values{
		"modus": {"1"},
		"type":  {strconv.Itoa(int(id))},
	}

	// Techs don't have a token
	if !id.IsTech() {
		pageHTML := b.getPageContent(url.Values{"page": {page}, "cp": {strconv.Itoa(int(celestialID))}})
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
		if err != nil {
			return err
		}
		token, exists := doc.Find("form").Find("input[name=token]").Attr("value")
		if !exists {
			return errors.New("unable to find form token")
		}
		payload.Add("token", token)
	}

	if id.IsDefense() || id.IsShip() {
		payload.Add("menge", strconv.Itoa(nbr))
	}

	url2 := b.serverURL + "/game/index.php?page=" + page + "&cp=" + strconv.Itoa(int(celestialID))
	resp, err := b.client.PostForm(url2, payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (b *OGame) buildCancelable(celestialID CelestialID, id ID) error {
	if !id.IsBuilding() && !id.IsTech() {
		return errors.New("invalid id " + id.String())
	}
	return b.build(celestialID, id, 0)
}

func (b *OGame) buildProduction(celestialID CelestialID, id ID, nbr int) error {
	if !id.IsDefense() && !id.IsShip() {
		return errors.New("invalid id " + id.String())
	}
	return b.build(celestialID, id, nbr)
}

func (b *OGame) buildBuilding(celestialID CelestialID, buildingID ID) error {
	if !buildingID.IsBuilding() {
		return errors.New("invalid building id " + buildingID.String())
	}
	return b.buildCancelable(celestialID, buildingID)
}

func (b *OGame) buildTechnology(planetID PlanetID, technologyID ID) error {
	if technologyID.IsTech() {
		return errors.New("invalid technology id " + technologyID.String())
	}
	return b.buildCancelable(CelestialID(planetID), technologyID)
}

func (b *OGame) buildDefense(celestialID CelestialID, defenseID ID, nbr int) error {
	if !defenseID.IsDefense() {
		return errors.New("invalid defense id " + defenseID.String())
	}
	return b.buildProduction(celestialID, ID(defenseID), nbr)
}

func (b *OGame) buildShips(celestialID CelestialID, shipID ID, nbr int) error {
	if !shipID.IsShip() {
		return errors.New("invalid ship id " + shipID.String())
	}
	return b.buildProduction(celestialID, ID(shipID), nbr)
}

func extractConstructions(pageHTML []byte) (buildingID ID, buildingCountdown int, researchID ID, researchCountdown int) {
	buildingCountdownMatch := regexp.MustCompile(`getElementByIdWithCache\("Countdown"\),(\d+),`).FindSubmatch(pageHTML)
	if len(buildingCountdownMatch) > 0 {
		buildingCountdown = toInt(buildingCountdownMatch[1])
		buildingIDInt := toInt(regexp.MustCompile(`onclick="cancelProduction\((\d+),`).FindSubmatch(pageHTML)[1])
		buildingID = ID(buildingIDInt)
	}
	researchCountdownMatch := regexp.MustCompile(`getElementByIdWithCache\("researchCountdown"\),(\d+),`).FindSubmatch(pageHTML)
	if len(researchCountdownMatch) > 0 {
		researchCountdown = toInt(researchCountdownMatch[1])
		researchIDInt := toInt(regexp.MustCompile(`onclick="cancelResearch\((\d+),`).FindSubmatch(pageHTML)[1])
		researchID = ID(researchIDInt)
	}
	return
}

func (b *OGame) constructionsBeingBuilt(celestialID CelestialID) (ID, int, ID, int) {
	pageHTML := b.getPageContent(url.Values{"page": {"overview"}, "cp": {strconv.Itoa(int(celestialID))}})
	return extractConstructions(pageHTML)
}

func extractCancelBuildingInfos(pageHTML []byte) (token string, techID, listID int, err error) {
	r1 := regexp.MustCompile(`page=overview&modus=2&token=(\w+)&techid="\+cancelProduction_id\+"&listid="\+production_listid`)
	m1 := r1.FindSubmatch(pageHTML)
	if len(m1) < 2 {
		return "", 0, 0, errors.New("unable to find token")
	}
	token = string(m1[1])
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	t := doc.Find("table.construction").Eq(0)
	a, _ := t.Find("a.abortNow").Attr("onclick")
	r := regexp.MustCompile(`cancelProduction\((\d+),\s?(\d+),`)
	m := r.FindStringSubmatch(a)
	if len(m) < 3 {
		return "", 0, 0, errors.New("unable to find techid/listid")
	}
	techID, _ = strconv.Atoi(m[1])
	listID, _ = strconv.Atoi(m[2])
	return
}

func extractCancelResearchInfos(pageHTML []byte) (token string, techID, listID int, err error) {
	r1 := regexp.MustCompile(`page=overview&modus=2&token=(\w+)"\+"&techid="\+id\+"&listid="\+listId`)
	m1 := r1.FindSubmatch(pageHTML)
	if len(m1) < 2 {
		return "", 0, 0, errors.New("unable to find token")
	}
	token = string(m1[1])
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	t := doc.Find("table.construction").Eq(1)
	a, _ := t.Find("a.abortNow").Attr("onclick")
	r := regexp.MustCompile(`cancelResearch\((\d+),\s?(\d+),`)
	m := r.FindStringSubmatch(a)
	if len(m) < 3 {
		return "", 0, 0, errors.New("unable to find techid/listid")
	}
	techID, _ = strconv.Atoi(m[1])
	listID, _ = strconv.Atoi(m[2])
	return
}

func (b *OGame) cancel(token string, techID, listID int) error {
	finalURL := b.serverURL + "/game/index.php?page=overview&modus=2&token=" + token + "&techid=" + strconv.Itoa(techID) + "&listid=" + strconv.Itoa(listID)
	req, _ := http.NewRequest("GET", finalURL, nil)
	resp, _ := b.client.Do(req)
	defer resp.Body.Close()
	return nil
}

func (b *OGame) cancelBuilding(celestialID CelestialID) error {
	pageHTML := b.getPageContent(url.Values{"page": {"overview"}, "cp": {strconv.Itoa(int(celestialID))}})
	token, techID, listID, _ := extractCancelBuildingInfos(pageHTML)
	return b.cancel(token, techID, listID)
}

func (b *OGame) cancelResearch(planetID PlanetID) error {
	pageHTML := b.getPageContent(url.Values{"page": {"overview"}, "cp": {planetID.String()}})
	token, techID, listID, _ := extractCancelResearchInfos(pageHTML)
	return b.cancel(token, techID, listID)
}

func (b *OGame) fetchResources(celestialID CelestialID) (resourcesResp, error) {
	pageJSON := b.getPageContent(url.Values{"page": {"fetchResources"}, "cp": {strconv.Itoa(int(celestialID))}})
	var res resourcesResp
	if err := json.Unmarshal(pageJSON, &res); err != nil {
		if isLogged(pageJSON) {
			return resourcesResp{}, ErrInvalidPlanetID
		}
		return resourcesResp{}, err
	}
	return res, nil
}

func (b *OGame) getResources(celestialID CelestialID) (Resources, error) {
	res, err := b.fetchResources(celestialID)
	return Resources{
		Metal:      res.Metal.Resources.Actual,
		Crystal:    res.Crystal.Resources.Actual,
		Deuterium:  res.Deuterium.Resources.Actual,
		Energy:     res.Energy.Resources.Actual,
		Darkmatter: res.Darkmatter.Resources.Actual,
	}, err
}

func extractFleet1Ships(pageHTML []byte) ShipsInfos {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	onclick := doc.Find("a#sendall").AttrOr("onclick", "")
	matches := regexp.MustCompile(`setMaxIntInput\("form\[name=shipsChosen]", (.+)\); checkShips`).FindStringSubmatch(onclick)
	if len(matches) == 0 {
		return ShipsInfos{}
	}
	m := matches[1]
	var res map[ID]int
	json.Unmarshal([]byte(m), &res)
	s := ShipsInfos{}
	for k, v := range res {
		s.Set(k, v)
	}
	return s
}

func (b *OGame) sendFleet(celestialID CelestialID, ships []Quantifiable, speed Speed, where Coordinate,
	destType DestinationType, mission MissionID, resources Resources) (FleetID, error) {
	getHiddenFields := func(pageHTML []byte) map[string]string {
		doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
		fields := make(map[string]string)
		doc.Find("input[type=hidden]").Each(func(i int, s *goquery.Selection) {
			name, _ := s.Attr("name")
			value, _ := s.Attr("value")
			fields[name] = value
		})
		return fields
	}

	// Page 1 : get to fleet page
	pageHTML := b.getPageContent(url.Values{"page": {"fleet1"}, "cp": {strconv.Itoa(int(celestialID))}})

	fleet1Doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	fleet1BodyID := fleet1Doc.Find("body").AttrOr("id", "")
	if fleet1BodyID != "fleet1" {
		return 0, ErrInvalidPlanetID
	}

	availableShips := extractFleet1Ships(pageHTML)

	atLeastOneShipSelected := false
	for _, ship := range ships {
		if ship.Nbr > 0 && availableShips.ByID(ship.ID) > 0 {
			atLeastOneShipSelected = true
			break
		}
	}
	if !atLeastOneShipSelected {
		return 0, ErrNoShipSelected
	}

	payload := url.Values{}
	hidden := getHiddenFields(pageHTML)
	for k, v := range hidden {
		payload.Add(k, v)
	}
	for _, s := range ships {
		if s.Nbr > 0 {
			payload.Add("am"+strconv.Itoa(int(s.ID)), strconv.Itoa(s.Nbr))
		}
	}

	// Page 2 : select ships
	fleet2URL := b.serverURL + "/game/index.php?page=fleet2"
	fleet2Resp, err := b.client.PostForm(fleet2URL, payload)
	if err != nil {
		return 0, err
	}
	defer fleet2Resp.Body.Close()
	pageHTML, _ = ioutil.ReadAll(fleet2Resp.Body)

	payload = url.Values{}
	hidden = getHiddenFields(pageHTML)
	for k, v := range hidden {
		payload.Add(k, v)
	}
	payload.Add("speed", strconv.Itoa(int(speed)))
	payload.Add("galaxy", strconv.Itoa(where.Galaxy))
	payload.Add("system", strconv.Itoa(where.System))
	payload.Add("position", strconv.Itoa(where.Position))
	t := destType
	if mission == RecycleDebrisField {
		t = DebrisDest // Send to debris field
	}
	payload.Add("type", strconv.Itoa(int(t)))

	// Check
	fleetCheckURL := b.serverURL + "/game/index.php?page=fleetcheck&ajax=1&espionage=0"
	fleetCheckPayload := url.Values{
		"galaxy": {strconv.Itoa(where.Galaxy)},
		"system": {strconv.Itoa(where.System)},
		"planet": {strconv.Itoa(where.Position)},
		"type":   {strconv.Itoa(int(t))},
	}
	fleetCheckResp, err := b.client.PostForm(fleetCheckURL, fleetCheckPayload)
	if err != nil {
		return 0, err
	}
	defer fleetCheckResp.Body.Close()
	by1, _ := ioutil.ReadAll(fleetCheckResp.Body)
	switch string(by1) {
	case "1":
		return 0, ErrUninhabitedPlanet
	case "1d":
		return 0, ErrNoDebrisField
	case "2":
		return 0, ErrPlayerInVacationMode
	case "3":
		return 0, ErrAdminOrGM
	case "4":
		return 0, ErrNoAstrophysics
	case "5":
		return 0, ErrNoobProtection
	case "6":
		return 0, ErrPlayerTooStrong
	case "10":
		return 0, ErrNoMoonAvailable
	case "11":
		return 0, ErrNoRecyclerAvailable
	case "15":
		return 0, ErrNoEventsRunning
	case "16":
		return 0, ErrPlanetAlreadyReservecForRelocation
	}

	// Page 3 : select coord, mission, speed
	fleet3URL := b.serverURL + "/game/index.php?page=fleet3"
	fleet3Resp, err := b.client.PostForm(fleet3URL, payload)
	if err != nil {
		return 0, err
	}
	defer fleet3Resp.Body.Close()
	pageHTML, _ = ioutil.ReadAll(fleet3Resp.Body)

	fleet3Doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	pageID := fleet3Doc.Find("body").AttrOr("id", "")
	if pageID == "fleet1" {
		return 0, errors.New("probably not enough space for deuterium")
	}

	payload = url.Values{}
	hidden = getHiddenFields(pageHTML)
	for k, v := range hidden {
		payload.Add(k, v)
	}
	payload.Add("crystal", strconv.Itoa(resources.Crystal))
	payload.Add("deuterium", strconv.Itoa(resources.Deuterium))
	payload.Add("metal", strconv.Itoa(resources.Metal))
	payload.Add("mission", strconv.Itoa(int(mission)))

	// Page 4 : send the fleet
	movementURL := b.serverURL + "/game/index.php?page=movement"
	movementResp, _ := b.client.PostForm(movementURL, payload)
	defer movementResp.Body.Close()

	// Page 5
	movementHTML := b.getPageContent(url.Values{"page": {"movement"}})
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(movementHTML))
	matches := make([]int, 0)
	originCoords, _ := doc.Find("meta[name=ogame-planet-coordinates]").Attr("content")
	doc.Find("div.fleetDetails").Each(func(i int, s *goquery.Selection) {
		origin := s.Find("span.originCoords").Text()
		dest := s.Find("span.destinationCoords").Text()
		reversalSpan := s.Find("span.reversal")
		if reversalSpan == nil {
			return
		}
		fleetIDStr, _ := reversalSpan.Attr("ref")
		fleetID, _ := strconv.Atoi(fleetIDStr)
		if dest == fmt.Sprintf("[%d:%d:%d]", where.Galaxy, where.System, where.Position) &&
			origin == fmt.Sprintf("[%s]", originCoords) {
			matches = append(matches, fleetID)
		}
	})
	if len(matches) > 0 {
		max := 0
		for _, v := range matches {
			if v > max {
				max = v
			}
		}
		return FleetID(max), nil
	}
	return 0, errors.New("could not find new fleet ID")
}

// EspionageReportType type of espionage report (action or report)
type EspionageReportType int

// Action message received when an enemy is seen naer your planet
const Action EspionageReportType = 0

// Report message received when you spied on someone
const Report EspionageReportType = 1

// CombatReportSummary summary of combat report
type CombatReportSummary struct {
	ID int
}

// EspionageReportSummary summary of espionage report
type EspionageReportSummary struct {
	ID     int
	Type   EspionageReportType
	From   string
	Target Coordinate
}

func extractCombatReportMessageIDs(pageHTML []byte) ([]CombatReportSummary, int) {
	msgs := make([]CombatReportSummary, 0)
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	nbPage, _ := strconv.Atoi(doc.Find("ul.pagination li").Last().AttrOr("data-page", "1"))
	doc.Find("li.msg").Each(func(i int, s *goquery.Selection) {
		if idStr, exists := s.Attr("data-msg-id"); exists {
			if id, err := strconv.Atoi(idStr); err == nil {
				report := CombatReportSummary{ID: id}
				msgs = append(msgs, report)

			}
		}
	})
	return msgs, nbPage
}

func extractEspionageReportMessageIDs(pageHTML []byte) ([]EspionageReportSummary, int) {
	msgs := make([]EspionageReportSummary, 0)
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	nbPage, _ := strconv.Atoi(doc.Find("ul.pagination li").Last().AttrOr("data-page", "1"))
	doc.Find("li.msg").Each(func(i int, s *goquery.Selection) {
		if idStr, exists := s.Attr("data-msg-id"); exists {
			if id, err := strconv.Atoi(idStr); err == nil {
				messageType := Report
				if s.Find("span.espionageDefText").Size() > 0 {
					messageType = Action
				}
				report := EspionageReportSummary{ID: id, Type: messageType}
				report.From = s.Find("span.msg_sender").Text()
				targetStr := s.Find("span.msg_title a").Text()
				report.Target = extractCoord(targetStr)
				msgs = append(msgs, report)

			}
		}
	})
	return msgs, nbPage
}

func (b *OGame) getPageMessages(page, tabid int) ([]byte, error) {
	finalURL := b.serverURL + "/game/index.php?page=messages"
	payload := url.Values{
		"messageId":  {"-1"},
		"tabid":      {strconv.Itoa(tabid)},
		"action":     {"107"},
		"pagination": {strconv.Itoa(page)},
		"ajax":       {"1"},
	}
	req, err := http.NewRequest("POST", finalURL, strings.NewReader(payload.Encode()))
	if err != nil {
		return []byte{}, err
	}
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	resp, err := b.client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	by, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return by, nil
}

func (b *OGame) getEspionageReportMessages() ([]EspionageReportSummary, error) {
	tabid := 20
	page := 1
	nbPage := 1
	msgs := make([]EspionageReportSummary, 0)
	for page <= nbPage {
		pageHTML, _ := b.getPageMessages(page, tabid)
		newMessages, newNbPage := extractEspionageReportMessageIDs(pageHTML)
		msgs = append(msgs, newMessages...)
		nbPage = newNbPage
		page++
	}
	return msgs, nil
}

// EspionageReport detailed espionage report
type EspionageReport struct {
	Resources
	MetalMine                    *int // ResourcesBuildings
	CrystalMine                  *int
	DeuteriumSynthesizer         *int
	SolarPlant                   *int
	FusionReactor                *int
	SolarSatellite               *int
	MetalStorage                 *int
	CrystalStorage               *int
	DeuteriumTank                *int
	RoboticsFactory              *int // Facilities
	Shipyard                     *int
	ResearchLab                  *int
	AllianceDepot                *int
	MissileSilo                  *int
	NaniteFactory                *int
	Terraformer                  *int
	SpaceDock                    *int
	LunarBase                    *int
	SensorPhalanx                *int
	JumpGate                     *int
	EnergyTechnology             *int // Researches
	LaserTechnology              *int
	IonTechnology                *int
	HyperspaceTechnology         *int
	PlasmaTechnology             *int
	CombustionDrive              *int
	ImpulseDrive                 *int
	HyperspaceDrive              *int
	EspionageTechnology          *int
	ComputerTechnology           *int
	Astrophysics                 *int
	IntergalacticResearchNetwork *int
	GravitonTechnology           *int
	WeaponsTechnology            *int
	ShieldingTechnology          *int
	ArmourTechnology             *int
	RocketLauncher               *int // Defenses
	LightLaser                   *int
	HeavyLaser                   *int
	GaussCannon                  *int
	IonCannon                    *int
	PlasmaTurret                 *int
	SmallShieldDome              *int
	LargeShieldDome              *int
	AntiBallisticMissiles        *int
	InterplanetaryMissiles       *int
	LightFighter                 *int // Fleets
	HeavyFighter                 *int
	Cruiser                      *int
	Battleship                   *int
	Battlecruiser                *int
	Bomber                       *int
	Destroyer                    *int
	Deathstar                    *int
	SmallCargo                   *int
	LargeCargo                   *int
	ColonyShip                   *int
	Recycler                     *int
	EspionageProbe               *int
	Coordinate                   Coordinate
	Type                         EspionageReportType
	Date                         time.Time
}

func extractEspionageReport(pageHTML []byte, location *time.Location) (EspionageReport, error) {
	report := EspionageReport{}
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	txt := doc.Find("span.msg_title a").First().Text()
	r := regexp.MustCompile(`([^\[]+) \[(\d+):(\d+):(\d+)]`)
	m := r.FindStringSubmatch(txt)
	report.Coordinate.Galaxy, _ = strconv.Atoi(m[2])
	report.Coordinate.System, _ = strconv.Atoi(m[3])
	report.Coordinate.Position, _ = strconv.Atoi(m[4])
	messageType := Report
	if doc.Find("span.espionageDefText").Size() > 0 {
		messageType = Action
	}
	report.Type = messageType
	msgDateRaw := doc.Find("span.msg_date").Text()
	msgDate, _ := time.ParseInLocation("02.01.2006 15:04:05", msgDateRaw, location)
	report.Date = msgDate.In(location)
	doc.Find("ul.detail_list").Each(func(i int, s *goquery.Selection) {
		dataType := s.AttrOr("data-type", "")
		if dataType == "resources" {
			report.Metal = parseInt(s.Find("li").Eq(0).AttrOr("title", "0"))
			report.Crystal = parseInt(s.Find("li").Eq(1).AttrOr("title", "0"))
			report.Deuterium = parseInt(s.Find("li").Eq(2).AttrOr("title", "0"))
			report.Energy = parseInt(s.Find("li").Eq(3).AttrOr("title", "0"))
		} else if dataType == "buildings" {
			s.Find("li.detail_list_el").Each(func(i int, s2 *goquery.Selection) {
				imgClass := s2.Find("img").AttrOr("class", "")
				r := regexp.MustCompile(`building(\d+)`)
				buildingID, _ := strconv.Atoi(r.FindStringSubmatch(imgClass)[1])
				l := parseInt(s2.Find("span.fright").Text())
				level := &l
				switch ID(buildingID) {
				case MetalMine.ID:
					report.MetalMine = level
				case CrystalMine.ID:
					report.CrystalMine = level
				case DeuteriumSynthesizer.ID:
					report.DeuteriumSynthesizer = level
				case SolarPlant.ID:
					report.SolarPlant = level
				case FusionReactor.ID:
					report.FusionReactor = level
				case MetalStorage.ID:
					report.MetalStorage = level
				case CrystalStorage.ID:
					report.CrystalStorage = level
				case DeuteriumTank.ID:
					report.DeuteriumTank = level
				case AllianceDepot.ID:
					report.AllianceDepot = level
				case RoboticsFactory.ID:
					report.RoboticsFactory = level
				case Shipyard.ID:
					report.Shipyard = level
				case ResearchLab.ID:
					report.ResearchLab = level
				case MissileSilo.ID:
					report.MissileSilo = level
				case NaniteFactory.ID:
					report.NaniteFactory = level
				case Terraformer.ID:
					report.Terraformer = level
				case SpaceDock.ID:
					report.SpaceDock = level
				case LunarBase.ID:
					report.LunarBase = level
				case SensorPhalanx.ID:
					report.SensorPhalanx = level
				case JumpGate.ID:
					report.JumpGate = level
				}
			})
		} else if dataType == "research" {
			s.Find("li.detail_list_el").Each(func(i int, s2 *goquery.Selection) {
				imgClass := s2.Find("img").AttrOr("class", "")
				r := regexp.MustCompile(`research(\d+)`)
				researchID, _ := strconv.Atoi(r.FindStringSubmatch(imgClass)[1])
				l := parseInt(s2.Find("span.fright").Text())
				level := &l
				switch ID(researchID) {
				case EspionageTechnology.ID:
					report.EspionageTechnology = level
				case ComputerTechnology.ID:
					report.ComputerTechnology = level
				case WeaponsTechnology.ID:
					report.WeaponsTechnology = level
				case ShieldingTechnology.ID:
					report.ShieldingTechnology = level
				case ArmourTechnology.ID:
					report.ArmourTechnology = level
				case EnergyTechnology.ID:
					report.EnergyTechnology = level
				case HyperspaceTechnology.ID:
					report.HyperspaceTechnology = level
				case CombustionDrive.ID:
					report.CombustionDrive = level
				case ImpulseDrive.ID:
					report.ImpulseDrive = level
				case HyperspaceDrive.ID:
					report.HyperspaceDrive = level
				case LaserTechnology.ID:
					report.LaserTechnology = level
				case IonTechnology.ID:
					report.IonTechnology = level
				case PlasmaTechnology.ID:
					report.PlasmaTechnology = level
				case IntergalacticResearchNetwork.ID:
					report.IntergalacticResearchNetwork = level
				case Astrophysics.ID:
					report.Astrophysics = level
				case GravitonTechnology.ID:
					report.GravitonTechnology = level
				}
			})
		} else if dataType == "ships" {
			s.Find("li.detail_list_el").Each(func(i int, s2 *goquery.Selection) {
				imgClass := s2.Find("img").AttrOr("class", "")
				r := regexp.MustCompile(`tech(\d+)`)
				shipID, _ := strconv.Atoi(r.FindStringSubmatch(imgClass)[1])
				l := parseInt(s2.Find("span.fright").Text())
				level := &l
				switch ID(shipID) {
				case SmallCargo.ID:
					report.SmallCargo = level
				case LargeCargo.ID:
					report.LargeCargo = level
				case LightFighter.ID:
					report.LightFighter = level
				case HeavyFighter.ID:
					report.HeavyFighter = level
				case Cruiser.ID:
					report.Cruiser = level
				case Battleship.ID:
					report.Battleship = level
				case ColonyShip.ID:
					report.ColonyShip = level
				case Recycler.ID:
					report.Recycler = level
				case EspionageProbe.ID:
					report.EspionageProbe = level
				case Bomber.ID:
					report.Bomber = level
				case SolarSatellite.ID:
					report.SolarSatellite = level
				case Destroyer.ID:
					report.Destroyer = level
				case Deathstar.ID:
					report.Deathstar = level
				case Battlecruiser.ID:
					report.Battlecruiser = level
				}
			})
		} else if dataType == "defense" {
			s.Find("li.detail_list_el").Each(func(i int, s2 *goquery.Selection) {
				imgClass := s2.Find("img").AttrOr("class", "")
				r := regexp.MustCompile(`defense(\d+)`)
				defenceID, _ := strconv.Atoi(r.FindStringSubmatch(imgClass)[1])
				l := parseInt(s2.Find("span.fright").Text())
				level := &l
				switch ID(defenceID) {
				case RocketLauncher.ID:
					report.RocketLauncher = level
				case LightLaser.ID:
					report.LightLaser = level
				case HeavyLaser.ID:
					report.HeavyLaser = level
				case GaussCannon.ID:
					report.GaussCannon = level
				case IonCannon.ID:
					report.IonCannon = level
				case PlasmaTurret.ID:
					report.PlasmaTurret = level
				case SmallShieldDome.ID:
					report.SmallShieldDome = level
				case LargeShieldDome.ID:
					report.LargeShieldDome = level
				case AntiBallisticMissiles.ID:
					report.AntiBallisticMissiles = level
				case InterplanetaryMissiles.ID:
					report.InterplanetaryMissiles = level
				}
			})
		}
	})
	return report, nil
}

func (b *OGame) getEspionageReport(msgID int) (EspionageReport, error) {
	pageHTML := b.getPageContent(url.Values{"page": {"messages"}, "messageId": {strconv.Itoa(msgID)}, "tabid": {"20"}, "ajax": {"1"}})
	return extractEspionageReport(pageHTML, b.location)
}

func (b *OGame) deleteMessage(msgID int) error {
	finalURL := b.serverURL + "/game/index.php?page=messages"
	payload := url.Values{
		"messageId": {strconv.Itoa(msgID)},
		"action":    {"103"},
		"ajax":      {"1"},
	}
	req, err := http.NewRequest("POST", finalURL, strings.NewReader(payload.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	resp, err := b.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	by, _ := ioutil.ReadAll(resp.Body)
	var res map[string]bool
	if err := json.Unmarshal(by, &res); err != nil {
		return errors.New("unable to find message id " + strconv.Itoa(msgID))
	}
	if val, ok := res[strconv.Itoa(msgID)]; !ok || !val {
		return errors.New("unable to find message id " + strconv.Itoa(msgID))
	}
	return nil
}

func energyProduced(temp Temperature, resourcesBuildings ResourcesBuildings, resSettings ResourceSettings, energyTechnology int) int {
	energyProduced := int(float64(SolarPlant.Production(resourcesBuildings.SolarPlant)) * (float64(resSettings.SolarPlant) / 100))
	energyProduced += int(float64(FusionReactor.Production(energyTechnology, resourcesBuildings.FusionReactor)) * (float64(resSettings.FusionReactor) / 100))
	energyProduced += int(float64(SolarSatellite.Production(temp, resourcesBuildings.SolarSatellite)) * (float64(resSettings.SolarSatellite) / 100))
	return energyProduced
}

func energyNeeded(resourcesBuildings ResourcesBuildings, resSettings ResourceSettings) int {
	energyNeeded := int(float64(MetalMine.EnergyConsumption(resourcesBuildings.MetalMine)) * (float64(resSettings.MetalMine) / 100))
	energyNeeded += int(float64(CrystalMine.EnergyConsumption(resourcesBuildings.CrystalMine)) * (float64(resSettings.CrystalMine) / 100))
	energyNeeded += int(float64(DeuteriumSynthesizer.EnergyConsumption(resourcesBuildings.DeuteriumSynthesizer)) * (float64(resSettings.DeuteriumSynthesizer) / 100))
	return energyNeeded
}

func productionRatio(temp Temperature, resourcesBuildings ResourcesBuildings, resSettings ResourceSettings, energyTechnology int) float64 {
	energyProduced := energyProduced(temp, resourcesBuildings, resSettings, energyTechnology)
	energyNeeded := energyNeeded(resourcesBuildings, resSettings)
	ratio := 1.0
	if energyNeeded > energyProduced {
		ratio = float64(energyProduced) / float64(energyNeeded)
	}
	return ratio
}

func getProductions(resBuildings ResourcesBuildings, resSettings ResourceSettings, researches Researches, universeSpeed int,
	temp Temperature, productionRatio float64) Resources {
	energyProduced := energyProduced(temp, resBuildings, resSettings, researches.EnergyTechnology)
	energyNeeded := energyNeeded(resBuildings, resSettings)
	return Resources{
		Metal:     MetalMine.Production(universeSpeed, productionRatio, researches.PlasmaTechnology, resBuildings.MetalMine),
		Crystal:   CrystalMine.Production(universeSpeed, productionRatio, researches.PlasmaTechnology, resBuildings.CrystalMine),
		Deuterium: DeuteriumSynthesizer.Production(universeSpeed, temp.Mean(), productionRatio, resBuildings.DeuteriumSynthesizer),
		Energy:    energyProduced - energyNeeded,
	}
}

func extractResourcesProductions(pageHTML []byte) (Resources, error) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	res := Resources{}
	selector := "table.listOfResourceSettingsPerPlanet tr.summary td span"
	el := doc.Find(selector)
	res.Metal = parseInt(el.Eq(0).AttrOr("title", "0"))
	res.Crystal = parseInt(el.Eq(1).AttrOr("title", "0"))
	res.Deuterium = parseInt(el.Eq(2).AttrOr("title", "0"))
	res.Energy = parseInt(el.Eq(3).AttrOr("title", "0"))
	return res, nil
}

func (b *OGame) getResourcesProductions(planetID PlanetID) (Resources, error) {
	planet, _ := b.getPlanet(planetID)
	resBuildings, _ := b.getResourcesBuildings(planetID)
	researches := b.getResearch()
	universeSpeed := b.getUniverseSpeed()
	resSettings, _ := b.getResourceSettings(planetID)
	ratio := productionRatio(planet.Temperature, resBuildings, resSettings, researches.EnergyTechnology)
	productions := getProductions(resBuildings, resSettings, researches, universeSpeed, planet.Temperature, ratio)
	return productions, nil
}

// GetSession get ogame session
func (b *OGame) GetSession() string {
	return b.ogameSession
}

// NewAccount response from creating a new account
type NewAccount struct {
	ID     int
	Server struct {
		Language string
		Number   int
	}
}

// AddAccount add a new account (server) to your list of accounts
func (b *OGame) AddAccount(number int, lang string) (NewAccount, error) {
	var payload struct {
		Language string `json:"language"`
		Number   int    `json:"number"`
	}
	payload.Language = lang
	payload.Number = number
	jsonPayloadBytes, err := json.Marshal(&payload)
	var newAccount NewAccount
	if err != nil {
		return newAccount, err
	}
	req, err := http.NewRequest("PUT", "https://lobby-api.ogame.gameforge.com/users/me/accounts", strings.NewReader(string(jsonPayloadBytes)))
	if err != nil {
		return newAccount, err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := b.client.Do(req)
	if err != nil {
		return newAccount, err
	}
	defer resp.Body.Close()
	by, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return newAccount, err
	}
	if err := json.Unmarshal(by, &newAccount); err != nil {
		return newAccount, err
	}
	return newAccount, nil
}

// GetServer get ogame server information that the bot is connected to
func (b *OGame) GetServer() Server {
	return b.server
}

// ServerURL get the ogame server specific url
func (b *OGame) ServerURL() string {
	return b.serverURL
}

// GetLanguage get ogame server language
func (b *OGame) GetLanguage() string {
	return b.language
}

// SetUserAgent change the user-agent used by the http client
func (b *OGame) SetUserAgent(newUserAgent string) {
	b.client.UserAgent = newUserAgent
}

// Login to ogame server
// Can fails with BadCredentialsError
func (b *OGame) Login() error {
	b.Lock()
	defer b.Unlock()
	return b.login()
}

// Logout the bot from ogame server
func (b *OGame) Logout() {
	b.Lock()
	defer b.Unlock()
	b.logout()
}

// GetUniverseName get the name of the universe the bot is playing into
func (b *OGame) GetUniverseName() string {
	return b.Universe
}

// GetUsername get the username that was used to login on ogame server
func (b *OGame) GetUsername() string {
	return b.Username
}

// GetUniverseSpeed shortcut to get ogame universe speed
func (b *OGame) GetUniverseSpeed() int {
	return b.getUniverseSpeed()
}

// GetUniverseSpeedFleet shortcut to get ogame universe speed fleet
func (b *OGame) GetUniverseSpeedFleet() int {
	return b.getUniverseSpeedFleet()
}

// IsDonutGalaxy shortcut to get ogame galaxy donut config
func (b *OGame) IsDonutGalaxy() bool {
	return b.isDonutGalaxy()
}

// IsDonutSystem shortcut to get ogame system donut config
func (b *OGame) IsDonutSystem() bool {
	return b.isDonutSystem()
}

// GetPageContent gets the html for a specific ogame page
func (b *OGame) GetPageContent(vals url.Values) []byte {
	b.Lock()
	defer b.Unlock()
	return b.getPageContent(vals)
}

// PostPageContent make a post request to ogame server
// This is useful when simulating a web browser
func (b *OGame) PostPageContent(vals, payload url.Values) []byte {
	b.Lock()
	defer b.Unlock()
	return b.postPageContent(vals, payload)
}

// IsUnderAttack returns true if the user is under attack, false otherwise
func (b *OGame) IsUnderAttack() bool {
	b.Lock()
	defer b.Unlock()
	return b.isUnderAttack()
}

// GetPlanets returns the user planets
func (b *OGame) GetPlanets() []Planet {
	b.Lock()
	defer b.Unlock()
	return b.getPlanets()
}

// GetCachedPlanets return planets from cached value
func (b *OGame) GetCachedPlanets() []Planet {
	return b.Planets
}

// GetPlanet gets infos for planetID
// Fails if planetID is invalid
func (b *OGame) GetPlanet(planetID PlanetID) (Planet, error) {
	b.Lock()
	defer b.Unlock()
	return b.getPlanet(planetID)
}

// GetPlanetByCoord get the player's planet using the coordinate
func (b *OGame) GetPlanetByCoord(coord Coordinate) (Planet, error) {
	b.Lock()
	defer b.Unlock()
	return b.getPlanetByCoord(coord)
}

// GetMoons returns the user moons
func (b *OGame) GetMoons(moonID MoonID) []Moon {
	b.Lock()
	defer b.Unlock()
	return b.getMoons()
}

// GetMoon gets infos for moonID
func (b *OGame) GetMoon(moonID MoonID) (Moon, error) {
	b.Lock()
	defer b.Unlock()
	return b.getMoon(moonID)
}

// GetMoonByCoord get the player's moon using the coordinate
func (b *OGame) GetMoonByCoord(coord Coordinate) (Moon, error) {
	b.Lock()
	defer b.Unlock()
	return b.getMoonByCoord(coord)
}

// ServerVersion returns OGame version
func (b *OGame) ServerVersion() string {
	return b.serverVersion()
}

// ServerTime returns server time
// Timezone is OGT (OGame Time zone)
func (b *OGame) ServerTime() time.Time {
	b.Lock()
	defer b.Unlock()
	return b.serverTime()
}

// GetUserInfos gets the user information
func (b *OGame) GetUserInfos() UserInfos {
	b.Lock()
	defer b.Unlock()
	return b.getUserInfos()
}

// SendMessage sends a message to playerID
func (b *OGame) SendMessage(playerID int, message string) error {
	b.Lock()
	defer b.Unlock()
	return b.sendMessage(playerID, message)
}

// GetFleets get the player's own fleets activities
func (b *OGame) GetFleets() []Fleet {
	b.Lock()
	defer b.Unlock()
	return b.getFleets()
}

// CancelFleet cancel a fleet
func (b *OGame) CancelFleet(fleetID FleetID) error {
	b.Lock()
	defer b.Unlock()
	return b.cancelFleet(fleetID)
}

// GetAttacks get enemy fleets attacking you
func (b *OGame) GetAttacks() []AttackEvent {
	b.Lock()
	defer b.Unlock()
	return b.getAttacks()
}

// GalaxyInfos get information of all planets and moons of a solar system
func (b *OGame) GalaxyInfos(galaxy, system int) (SystemInfos, error) {
	b.Lock()
	defer b.Unlock()
	return b.galaxyInfos(galaxy, system)
}

// GetResourceSettings gets the resources settings for specified planetID
func (b *OGame) GetResourceSettings(planetID PlanetID) (ResourceSettings, error) {
	b.Lock()
	defer b.Unlock()
	return b.getResourceSettings(planetID)
}

// SetResourceSettings set the resources settings on a planet
func (b *OGame) SetResourceSettings(planetID PlanetID, settings ResourceSettings) error {
	b.Lock()
	defer b.Unlock()
	return b.setResourceSettings(planetID, settings)
}

// GetResourcesBuildings gets the resources buildings levels
func (b *OGame) GetResourcesBuildings(planetID PlanetID) (ResourcesBuildings, error) {
	b.Lock()
	defer b.Unlock()
	return b.getResourcesBuildings(planetID)
}

// GetDefense gets all the defenses units information of a planet
// Fails if planetID is invalid
func (b *OGame) GetDefense(celestialID CelestialID) (DefensesInfos, error) {
	b.Lock()
	defer b.Unlock()
	return b.getDefense(celestialID)
}

// GetShips gets all ships units information of a planet
func (b *OGame) GetShips(celestialID CelestialID) (ShipsInfos, error) {
	b.Lock()
	defer b.Unlock()
	return b.getShips(celestialID)
}

// GetFacilities gets all facilities information of a planet
func (b *OGame) GetFacilities(celestialID CelestialID) (Facilities, error) {
	b.Lock()
	defer b.Unlock()
	return b.getFacilities(celestialID)
}

// GetProduction get what is in the production queue.
// (ships & defense being built)
func (b *OGame) GetProduction(celestialID CelestialID) ([]Quantifiable, error) {
	b.Lock()
	defer b.Unlock()
	return b.getProduction(celestialID)
}

// GetResearch gets the player researches information
func (b *OGame) GetResearch() Researches {
	b.Lock()
	defer b.Unlock()
	return b.getResearch()
}

// Build builds any ogame objects (building, technology, ship, defence)
func (b *OGame) Build(celestialID CelestialID, id ID, nbr int) error {
	b.Lock()
	defer b.Unlock()
	return b.build(celestialID, id, nbr)
}

// BuildCancelable builds any cancelable ogame objects (building, technology)
func (b *OGame) BuildCancelable(celestialID CelestialID, id ID) error {
	b.Lock()
	defer b.Unlock()
	return b.buildCancelable(celestialID, id)
}

// BuildProduction builds any line production ogame objects (ship, defence)
func (b *OGame) BuildProduction(celestialID CelestialID, id ID, nbr int) error {
	b.Lock()
	defer b.Unlock()
	return b.buildProduction(celestialID, id, nbr)
}

// BuildBuilding ensure what is being built is a building
func (b *OGame) BuildBuilding(celestialID CelestialID, buildingID ID) error {
	b.Lock()
	defer b.Unlock()
	return b.buildBuilding(celestialID, buildingID)
}

// BuildDefense builds a defense unit
func (b *OGame) BuildDefense(celestialID CelestialID, defenseID ID, nbr int) error {
	b.Lock()
	defer b.Unlock()
	return b.buildDefense(celestialID, defenseID, nbr)
}

// BuildShips builds a ship unit
func (b *OGame) BuildShips(celestialID CelestialID, shipID ID, nbr int) error {
	b.Lock()
	defer b.Unlock()
	return b.buildShips(celestialID, shipID, nbr)
}

// ConstructionsBeingBuilt returns the building & research being built, and the time remaining (secs)
func (b *OGame) ConstructionsBeingBuilt(celestialID CelestialID) (ID, int, ID, int) {
	b.Lock()
	defer b.Unlock()
	return b.constructionsBeingBuilt(celestialID)
}

// CancelBuilding cancel the construction of a building on a specified planet
func (b *OGame) CancelBuilding(celestialID CelestialID) error {
	b.Lock()
	defer b.Unlock()
	return b.cancelBuilding(celestialID)
}

// CancelResearch cancel the research
func (b *OGame) CancelResearch(planetID PlanetID) error {
	b.Lock()
	defer b.Unlock()
	return b.cancelResearch(planetID)
}

// BuildTechnology ensure that we're trying to build a technology
func (b *OGame) BuildTechnology(planetID PlanetID, technologyID ID) error {
	b.Lock()
	defer b.Unlock()
	return b.buildTechnology(planetID, technologyID)
}

// GetResources gets user resources
func (b *OGame) GetResources(celestialID CelestialID) (Resources, error) {
	b.Lock()
	defer b.Unlock()
	return b.getResources(celestialID)
}

// SendFleet sends a fleet
func (b *OGame) SendFleet(celestialID CelestialID, ships []Quantifiable, speed Speed, where Coordinate,
	destType DestinationType, mission MissionID, resources Resources) (FleetID, error) {
	b.Lock()
	defer b.Unlock()
	return b.sendFleet(celestialID, ships, speed, where, destType, mission, resources)
}

// GetEspionageReportMessages gets the summary of each espionage reports
func (b *OGame) GetEspionageReportMessages() ([]EspionageReportSummary, error) {
	b.Lock()
	defer b.Unlock()
	return b.getEspionageReportMessages()
}

// GetEspionageReport gets a detailed espionage report
func (b *OGame) GetEspionageReport(msgID int) (EspionageReport, error) {
	b.Lock()
	defer b.Unlock()
	return b.getEspionageReport(msgID)
}

// DeleteMessage deletes a message from the mail box
func (b *OGame) DeleteMessage(msgID int) error {
	b.Lock()
	defer b.Unlock()
	return b.deleteMessage(msgID)
}

// GetResourcesProductions gets the planet resources production
func (b *OGame) GetResourcesProductions(planetID PlanetID) (Resources, error) {
	b.Lock()
	defer b.Unlock()
	return b.getResourcesProductions(planetID)
}

// FlightTime calculate flight time and fuel needed
func (b *OGame) FlightTime(origin, destination Coordinate, speed Speed, ships ShipsInfos) (secs, fuel int) {
	b.Lock()
	defer b.Unlock()
	return calcFlightTime(origin, destination, b.universeSize, b.donutGalaxy, b.donutSystem, float64(speed)/10, b.universeSpeedFleet, ships, Researches{})
}

// RegisterChatCallback register a callback that is called when chat messages are received
func (b *OGame) RegisterChatCallback(fn func(msg ChatMsg)) {
	b.chatCallbacks = append(b.chatCallbacks, fn)
}

// Phalanx scan a coordinate from a moon to get fleets information
// IMPORTANT: My account was instantly banned when I scanned an invalid coordinate.
// IMPORTANT: This function DOES validate that the coordinate is a valid planet in range of phalanx
// 			  and that you have enough deuterium.
func (b *OGame) Phalanx(moonID MoonID, coord Coordinate) ([]Fleet, error) {
	b.Lock()
	defer b.Unlock()
	return b.getPhalanx(moonID, coord)
}
