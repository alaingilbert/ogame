package ogame

import (
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
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/Sirupsen/logrus"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

// Wrapper ...
type Wrapper interface {
	SetUserAgent(newUserAgent string)
	ServerURL() string
	GetPageContent(url.Values) string
	PostPageContent(url.Values, url.Values) string
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
	GalaxyInfos(galaxy, system int) ([]PlanetInfos, error)
	GetResearch() Researches
	GetPlanets() []Planet
	GetPlanetByCoord(Coordinate) (Planet, error)
	GetPlanet(PlanetID) (Planet, error)
	GetEspionageReportMessages() ([]EspionageReportSummary, error)
	GetEspionageReport(msgID int) (EspionageReport, error)
	DeleteMessage(msgID int) error

	// Planet specific functions
	GetResourceSettings(PlanetID) (ResourceSettings, error)
	SetResourceSettings(PlanetID, ResourceSettings) error
	GetResourcesBuildings(PlanetID) (ResourcesBuildings, error)
	GetDefense(PlanetID) (DefensesInfos, error)
	GetShips(PlanetID) (ShipsInfos, error)
	GetFacilities(PlanetID) (Facilities, error)
	Build(planetID PlanetID, id ID, nbr int) error
	BuildCancelable(PlanetID, ID) error
	BuildProduction(planetID PlanetID, id ID, nbr int) error
	BuildBuilding(planetID PlanetID, buildingID ID) error
	BuildTechnology(planetID PlanetID, technologyID ID) error
	BuildDefense(planetID PlanetID, defenseID ID, nbr int) error
	BuildShips(planetID PlanetID, shipID ID, nbr int) error
	GetProduction(PlanetID) ([]Quantifiable, error)
	ConstructionsBeingBuilt(PlanetID) (buildingID ID, buildingCountdown int, researchID ID, researchCountdown int)
	CancelBuilding(PlanetID) error
	CancelResearch(PlanetID) error
	GetResources(PlanetID) (Resources, error)
	SendFleet(planetID PlanetID, ships []Quantifiable, speed Speed, where Coordinate, mission MissionID, resources Resources) (FleetID, error)
	//GetResourcesProductionRatio(PlanetID) (float64, error)
	GetResourcesProductions(PlanetID) (Resources, error)
}

const defaultUserAgent = "" +
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) " +
	"AppleWebKit/537.36 (KHTML, like Gecko) " +
	"Chrome/51.0.2704.103 " +
	"Safari/537.36"

// ErrNotLogged ...
var ErrNotLogged = errors.New("not logged")

// ErrBadCredentials ...
var ErrBadCredentials = errors.New("bad credentials")

// ErrInvalidPlanetID ...
var ErrInvalidPlanetID = errors.New("invalid planet id")

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

// Technology ...
type Technology interface {
	GetID() ID
	GetName() string
	GetBaseCost() Resources
	GetIncreaseFactor() float64
	GetRequirements() map[ID]int
	IsAvailable(ResourcesBuildings, Facilities, Researches, int) bool
	GetLevel(ResourcesBuildings, Facilities, Researches) int
	GetPrice(level int) Resources
	ConstructionTime(level, universeSpeed int, facilities Facilities) int
}

// Building ...
type Building interface {
	GetID() ID
	GetName() string
	GetBaseCost() Resources
	GetIncreaseFactor() float64
	GetRequirements() map[ID]int
	IsAvailable(ResourcesBuildings, Facilities, Researches, int) bool
	GetLevel(ResourcesBuildings, Facilities, Researches) int
	GetPrice(level int) Resources
	ConstructionTime(level, universeSpeed int, facilities Facilities) int
}

// Ship ...
type Ship interface {
	GetID() ID
	GetName() string
	GetRequirements() map[ID]int
	IsAvailable(ResourcesBuildings, Facilities, Researches, int) bool
	GetPrice(int) Resources
	GetStructuralIntegrity() int
	GetShieldPower() int
	GetWeaponPower() int
	GetCargoCapacity() int
	GetBaseSpeed() int
	GetSpeed(researches Researches) int
	GetFuelConsumption() int
	GetRapidfireFrom() map[ID]int
	GetRapidfireAgainst() map[ID]int
}

// Defense ...
type Defense interface {
	GetID() ID
	GetName() string
	GetPrice(int) Resources
	GetRequirements() map[ID]int
	IsAvailable(ResourcesBuildings, Facilities, Researches, int) bool
	GetStructuralIntegrity() int
	GetShieldPower() int
	GetWeaponPower() int
	GetRapidfireFrom() map[ID]int
}

// ObjsStruct ...
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

// Objs ...
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

// Defenses ...
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

// Ships ...
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

// Buildings ...
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

// Technologies ...
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
	location           *time.Location
	universeSpeed      int
	universeSpeedFleet int
	donutGalaxy        bool
	donutSystem        bool
	ogameVersion       string
	serverURL          string
	client             *ogameClient
	logger             *log.Logger
}

// Params ...
type Params struct {
	Universe  string
	Username  string
	Password  string
	AutoLogin bool
	Proxy     string
}

// New creates a new instance of OGame wrapper.
func New(universe, username, password string) (*OGame, error) {
	b := NewNoLogin(universe, username, password)

	if err := b.Login(); err != nil {
		return nil, err
	}

	return b, nil
}

// NewWithParams ...
func NewWithParams(params Params) (*OGame, error) {
	b := NewNoLogin(params.Universe, params.Username, params.Password)

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
func NewNoLogin(universe, username, password string) *OGame {
	b := new(OGame)
	b.quiet = false
	b.logger = log.New(os.Stdout, "", 0)

	b.Universe = universe
	b.Username = username
	b.password = password

	jar, _ := cookiejar.New(nil)
	b.client = &ogameClient{}
	b.client.Jar = jar
	b.client.UserAgent = defaultUserAgent

	return b
}

// Quiet ...
func (b *OGame) Quiet(quiet bool) {
	b.quiet = quiet
}

// SetLogger ...
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

type server struct {
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
		if cookie.Name == "PHPSESSID" {
			return cookie.Value, nil
		}
	}

	return "", errors.New("PHPSESSID not found")
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
	req.AddCookie(&http.Cookie{Name: "PHPSESSID", Value: phpSessionID})
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

func getServers(client *ogameClient) ([]server, error) {
	var servers []server
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

func findAccountByName(universe string, accounts []account, servers []server) (account, server, error) {
	for _, a := range accounts {
		for _, s := range servers {
			if universe == s.Name && a.Server.Language == s.Language && a.Server.Number == s.Number {
				return a, s, nil
			}
		}
	}
	return account{}, server{}, fmt.Errorf("server %s not found", universe)
}

func getLoginLink(client *ogameClient, userAccount account, phpSessionID string) (string, error) {
	ogURL := fmt.Sprintf("https://lobby-api.ogame.gameforge.com/users/me/loginLink?id=%d&server[language]=%s&server[number]=%d",
		userAccount.ID, userAccount.Server.Language, userAccount.Server.Number)
	req, err := http.NewRequest("GET", ogURL, nil)
	if err != nil {
		return "", err
	}
	req.AddCookie(&http.Cookie{Name: "PHPSESSID", Value: phpSessionID})
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
	userAccount, server, err := findAccountByName(b.Universe, accounts, servers)
	if err != nil {
		return err
	}
	b.debug("Players online: " + strconv.Itoa(server.PlayersOnline) + ", Players: " + strconv.Itoa(server.PlayerCount))
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

	by, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	pageHTML := string(by)

	b.debug("extract informations from html")
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(pageHTML))
	if err != nil {
		return err
	}
	b.ogameSession, _ = doc.Find("meta[name=ogame-session]").Attr("content")
	if b.ogameSession == "" {
		return errors.New("bad credentials")
	}

	serverTime, _ := extractServerTime(pageHTML)
	b.location = serverTime.Location()
	b.universeSpeed, _ = strconv.Atoi(doc.Find("meta[name=ogame-universe-speed]").AttrOr("content", "1"))
	b.universeSpeedFleet, _ = strconv.Atoi(doc.Find("meta[name=ogame-universe-speed-fleet]").AttrOr("content", "1"))
	b.donutGalaxy, _ = strconv.ParseBool(doc.Find("meta[name=ogame-donut-galaxy]").AttrOr("content", "1"))
	b.donutSystem, _ = strconv.ParseBool(doc.Find("meta[name=ogame-donut-system]").AttrOr("content", "1"))
	b.ogameVersion = doc.Find("meta[name=ogame-version]").AttrOr("content", "")

	b.Player, _ = extractUserInfos(pageHTML, b.language)
	b.Planets = extractPlanets(pageHTML, b)

	return nil
}

func (b *OGame) logout() {
	b.getPageContent(url.Values{"page": {"logout"}})
}

func isLogged(pageHTML string) bool {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(pageHTML))
	if err != nil {
		return false
	}
	ogameSession, _ := doc.Find("meta[name=ogame-session]").Attr("content")
	if ogameSession == "" {
		return false
	}
	return true
}

func isAjaxPage(page string) bool {
	return page == "fetchEventbox" ||
		page == "fetchResources" ||
		page == "galaxyContent" ||
		page == "eventList" ||
		page == "ajaxChat"
}

func isPartialPage(vals url.Values) bool {
	page := vals.Get("page")
	ajax := vals.Get("ajax")

	if page == "techtree" {
		return true
	}

	if ajax == "1" {
		return true
	}

	return false
}

func (b *OGame) postPageContent(vals, payload url.Values) string {
	finalURL := b.serverURL + "/game/index.php?" + vals.Encode()
	req, err := http.NewRequest("POST", finalURL, strings.NewReader(payload.Encode()))
	if err != nil {
		b.error(err)
		return ""
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")

	resp, err := b.client.Do(req)
	if err != nil {
		b.error(err)
		return ""
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
		return ""
	}

	return string(body)
}

func (b *OGame) getPageContent(vals url.Values) string {
	if b.serverURL == "" {
		logrus.Error("serverURL is empty")
		return ""
	}

	finalURL := b.serverURL + "/game/index.php?" + vals.Encode()
	page := vals.Get("page")
	var pageHTML string

	b.withRetry(func() error {
		req, err := http.NewRequest("GET", finalURL, nil)
		if err != nil {
			return err
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
		pageHTML = string(by)

		if page != "logout" && !isAjaxPage(page) && !isPartialPage(vals) && !isLogged(pageHTML) {
			return ErrNotLogged
		}

		return nil
	})

	if page == "overview" {
		b.Player, _ = extractUserInfos(pageHTML, b.language)
		b.Planets = extractPlanets(pageHTML, b)
	} else if isAjaxPage(page) {
	} else {
		b.Planets = extractPlanets(pageHTML, b)
	}
	return pageHTML
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
		if err := json.Unmarshal([]byte(pageJSON), v); err != nil {
			return ErrNotLogged
		}
		return nil
	})
}

// extract universe speed from html calculation
// pageHTML := b.getPageContent(url.Values{"page": {"techtree"}, "tab": {"2"}, "techID": {"1"}})
func extractUniverseSpeed(pageHTML string) int {
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(pageHTML))
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

var planetInfosRgx = regexp.MustCompile(`([^\[]+) \[(\d+):(\d+):(\d+)]([\d.]+)km \((\d+)/(\d+)\)(?:de )?([-\d]+).+C\s*(?:bis|to|à|a) ([-\d]+).+C`)

func extractPlanets(pageHTML string, b *OGame) []Planet {
	res := make([]Planet, 0)
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(pageHTML))
	doc.Find("div.smallplanet").Each(func(i int, s *goquery.Selection) {
		el, _ := s.Attr("id")
		id, err := strconv.Atoi(strings.TrimPrefix(el, "planet-"))
		if err != nil {
			return
		}

		planetName := s.Find("span.planet-name").Text()
		planetKoords := s.Find("span.planet-koords").Text()
		planetPic, _ := s.Find("img.planetPic").Attr("src")

		txt, _ := s.Find("a.planetlink").Attr("title")
		p := bluemonday.StrictPolicy()
		m1 := planetInfosRgx.FindStringSubmatch(p.Sanitize(txt))
		if len(m1) < 10 {
			b.error("failed to parse planet infos: " + txt)
			return
		}

		planet := Planet{}
		planet.ogame = b
		planet.Img = planetPic
		planet.ID = PlanetID(id)
		planet.Name = planetName
		planet.Coordinate = extractCoord(planetKoords)
		planet.Diameter = parseInt(m1[5])
		planet.Fields.Built, _ = strconv.Atoi(m1[6])
		planet.Fields.Total, _ = strconv.Atoi(m1[7])
		planet.Temperature.Min, _ = strconv.Atoi(m1[8])
		planet.Temperature.Max, _ = strconv.Atoi(m1[9])

		res = append(res, planet)
	})
	return res
}

func (b *OGame) getPlanets() []Planet {
	pageHTML := b.getPageContent(url.Values{"page": {"overview"}})
	return extractPlanets(pageHTML, b)
}

func (b *OGame) getPlanetByCoord(coord Coordinate) (Planet, error) {
	planets := b.getPlanets()
	for _, planet := range planets {
		if planet.Coordinate.Equal(coord) {
			return planet, nil
		}
	}
	return Planet{}, errors.New("invalid planet coordinate")
}

func (b *OGame) getPlanet(planetID PlanetID) (Planet, error) {
	planetIDStr := planetID.String()
	pageHTML := b.getPageContent(url.Values{"page": {"overview"}, "cp": {planetIDStr}})
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(pageHTML))
	s := doc.Find("div#planet-" + planetIDStr)
	if len(s.Nodes) > 0 { // planet
		title, _ := s.Find("a").Attr("title")
		root, err := html.Parse(strings.NewReader(title))
		if err != nil {
			return Planet{}, err
		}
		txt := goquery.NewDocumentFromNode(root).Text()
		m := planetInfosRgx.FindStringSubmatch(txt)

		res := Planet{}
		res.ogame = b
		res.Img, _ = s.Find("img").Attr("src")
		res.ID = planetID
		res.Name = m[1]
		res.Coordinate.Galaxy, _ = strconv.Atoi(m[2])
		res.Coordinate.System, _ = strconv.Atoi(m[3])
		res.Coordinate.Position, _ = strconv.Atoi(m[4])
		res.Diameter, _ = strconv.Atoi(m[5])
		res.Fields.Built, _ = strconv.Atoi(m[6])
		res.Fields.Total, _ = strconv.Atoi(m[7])
		res.Temperature.Min, _ = strconv.Atoi(m[8])
		res.Temperature.Max, _ = strconv.Atoi(m[9])
		return res, nil
	}
	return Planet{}, errors.New("failed to find planetID")
}

func (b *OGame) serverVersion() string {
	return b.ogameVersion
}

func extractServerTime(pageHTML string) (time.Time, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(pageHTML))
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

func extractUserInfos(pageHTML, lang string) (UserInfos, error) {
	playerIDRgx := regexp.MustCompile(`playerId="(\w+)"`)
	playerNameRgx := regexp.MustCompile(`playerName="([^"]+)"`)
	txtContent := regexp.MustCompile(`textContent\[7]="([^"]+)"`)
	playerIDGroups := playerIDRgx.FindStringSubmatch(pageHTML)
	playerNameGroups := playerNameRgx.FindStringSubmatch(pageHTML)
	subHTMLGroups := txtContent.FindStringSubmatch(pageHTML)
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
	res.PlayerID, _ = strconv.Atoi(playerIDGroups[1])
	res.PlayerName = playerNameGroups[1]
	html2 := subHTMLGroups[1]
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html2))

	infosRgx := regexp.MustCompile(`([\d\\.]+) \(Place ([\d.]+) of ([\d.]+)\)`)
	switch lang {
	case "fr":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(Place ([\d.]+) sur ([\d.]+)\)`)
	case "de":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(Platz ([\d.]+) von ([\d.]+)\)`)
	case "es":
		infosRgx = regexp.MustCompile(`([\d\\.]+) \(Lugar ([\d.]+) de ([\d.]+)\)`)
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
	honourPointsGroups := honourPointsRgx.FindStringSubmatch(pageHTML)
	if len(honourPointsGroups) < 2 {
		return UserInfos{}, errors.New("cannot find honour points")
	}
	res.HonourPoints = parseInt(honourPointsGroups[1])
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

func extractFleets(pageHTML string) (res []Fleet) {
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(pageHTML))
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
		shipment.Metal, _ = strconv.Atoi(strings.Trim(trs.Eq(trs.Size()-3).Find("td").Eq(1).Text(), "\n\t\r "))
		shipment.Crystal, _ = strconv.Atoi(strings.Trim(trs.Eq(trs.Size()-2).Find("td").Eq(1).Text(), "\n\t\r "))
		shipment.Deuterium, _ = strconv.Atoi(strings.Trim(trs.Eq(trs.Size()-1).Find("td").Eq(1).Text(), "\n\t\r "))

		fleet := Fleet{}
		fleet.ID = FleetID(id)
		fleet.Origin = origin
		fleet.Destination = dest
		fleet.Mission = MissionID(missionType)
		fleet.ReturnFlight = returnFlight
		fleet.Resources = shipment
		fleet.ArriveIn = secs

		for i := 1; i < trs.Size()-5; i++ {
			name := strings.ToLower(strings.Trim(trs.Eq(i).Find("td").Eq(0).Text(), "\n\t\r :"))
			qty := parseInt(strings.Trim(trs.Eq(i).Find("td").Eq(1).Text(), "\n\t\r "))
			shipID, _ := parseShip(name)
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
	val := 0.0
	if galaxy1 < galaxy2 {
		val = math.Min(float64(galaxy2-galaxy1), float64((galaxy1+universeSize)-galaxy2))
	} else {
		val = math.Min(float64(galaxy1-galaxy2), float64((galaxy2+universeSize)-galaxy1))
	}
	return int(20000 * val)
}

// Returns the distance between two systems
func systemDistance(system1, system2 int, donutSystem bool) (distance int) {
	if !donutSystem {
		return int(2700 + 95*math.Abs(float64(system2-system1)))
	}
	systemSize := 499
	val := 0.0
	if system1 < system2 {
		val = math.Min(float64(system2-system1), float64((system1+systemSize)-system2))
	} else {
		val = math.Min(float64(system1-system2), float64((system2+systemSize)-system1))
	}
	return int(2700 + 95*val)
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
		return systemDistance(c1.System, c2.System, donutSystem)
	}
	return planetDistance(c1.Position, c2.Position)
}

func findSlowestSpeed(ships ShipsInfos) int {
	minSpeed := math.MaxInt32
	if ships.LightFighter > 0 && LightFighter.BaseSpeed < minSpeed {
		minSpeed = LightFighter.BaseSpeed
	}
	if ships.HeavyFighter > 0 && HeavyFighter.BaseSpeed < minSpeed {
		minSpeed = HeavyFighter.BaseSpeed
	}
	if ships.Cruiser > 0 && Cruiser.BaseSpeed < minSpeed {
		minSpeed = Cruiser.BaseSpeed
	}
	if ships.Battleship > 0 && Battleship.BaseSpeed < minSpeed {
		minSpeed = Battleship.BaseSpeed
	}
	if ships.Battlecruiser > 0 && Battlecruiser.BaseSpeed < minSpeed {
		minSpeed = Battlecruiser.BaseSpeed
	}
	if ships.Bomber > 0 && Bomber.BaseSpeed < minSpeed {
		minSpeed = Bomber.BaseSpeed
	}
	if ships.Destroyer > 0 && Destroyer.BaseSpeed < minSpeed {
		minSpeed = Destroyer.BaseSpeed
	}
	if ships.Deathstar > 0 && Deathstar.BaseSpeed < minSpeed {
		minSpeed = Deathstar.BaseSpeed
	}
	if ships.SmallCargo > 0 && SmallCargo.BaseSpeed < minSpeed {
		minSpeed = SmallCargo.BaseSpeed
	}
	if ships.LargeCargo > 0 && LargeCargo.BaseSpeed < minSpeed {
		minSpeed = LargeCargo.BaseSpeed
	}
	if ships.ColonyShip > 0 && ColonyShip.BaseSpeed < minSpeed {
		minSpeed = ColonyShip.BaseSpeed
	}
	if ships.Recycler > 0 && Recycler.BaseSpeed < minSpeed {
		minSpeed = Recycler.BaseSpeed
	}
	if ships.EspionageProbe > 0 && EspionageProbe.BaseSpeed < minSpeed {
		minSpeed = EspionageProbe.BaseSpeed
	}
	return minSpeed
}

func calcFuel(ships ShipsInfos, dist int, speed float64) (fuel int) {
	tmpFn := func(baseFuel int) int {
		return int(1 + math.Round(((float64(baseFuel)*float64(dist))/35000)*math.Pow(speed+1, 2)))
	}
	if ships.LightFighter > 0 {
		fuel += tmpFn(LightFighter.FuelConsumption) * ships.LightFighter
	}
	if ships.HeavyFighter > 0 {
		fuel += tmpFn(HeavyFighter.FuelConsumption) * ships.LightFighter
	}
	if ships.Cruiser > 0 {
		fuel += tmpFn(Cruiser.FuelConsumption) * ships.LightFighter
	}
	if ships.Battleship > 0 {
		fuel += tmpFn(Battleship.FuelConsumption) * ships.LightFighter
	}
	if ships.Battlecruiser > 0 {
		fuel += tmpFn(Battlecruiser.FuelConsumption) * ships.LightFighter
	}
	if ships.Bomber > 0 {
		fuel += tmpFn(Bomber.FuelConsumption) * ships.LightFighter
	}
	if ships.Destroyer > 0 {
		fuel += tmpFn(Destroyer.FuelConsumption) * ships.LightFighter
	}
	if ships.Deathstar > 0 {
		fuel += tmpFn(Deathstar.FuelConsumption) * ships.LightFighter
	}
	if ships.SmallCargo > 0 {
		fuel += tmpFn(SmallCargo.FuelConsumption) * ships.LightFighter
	}
	if ships.LargeCargo > 0 {
		fuel += tmpFn(LargeCargo.FuelConsumption) * ships.LightFighter
	}
	if ships.ColonyShip > 0 {
		fuel += tmpFn(ColonyShip.FuelConsumption) * ships.LightFighter
	}
	if ships.Recycler > 0 {
		fuel += tmpFn(Recycler.FuelConsumption) * ships.LightFighter
	}
	if ships.EspionageProbe > 0 {
		fuel += tmpFn(EspionageProbe.FuelConsumption) * ships.LightFighter
	}
	return
}

func calcFlightTime(origin, destination Coordinate, universeSize int, donutGalaxy, donutSystem bool, speed float64,
	universeSpeedFleet int, ships ShipsInfos, techs Researches) (secs, fuel int) {
	baseSpeed := float64(findSlowestSpeed(ships))
	s := speed
	v := baseSpeed + (baseSpeed*0.2)*6
	a := float64(universeSpeedFleet)
	d := float64(distance(origin, destination, universeSize, donutGalaxy, donutSystem))
	secs = int(((10 + (3500 / s)) * math.Sqrt((10*d)/v)) / a)
	fuel = calcFuel(ships, int(d), speed)
	return
}

func extractAttacks(pageHTML string) []AttackEvent {
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(pageHTML))
	attacks := make([]AttackEvent, 0)
	tmp := func(i int, s *goquery.Selection) {
		classes, _ := s.Attr("class")
		if strings.Contains(classes, "partnerInfo") {
			return
		}
		missionTypeStr, _ := s.Attr("data-mission-type")
		arrivalTimeStr, _ := s.Attr("data-arrival-time")
		missionTypeInt, _ := strconv.Atoi(missionTypeStr)
		arrivalTimeInt, _ := strconv.Atoi(arrivalTimeStr)
		missionType := MissionID(missionTypeInt)
		if missionType != Attack && missionType != GroupedAttack && missionType != MissileAttack {
			return
		}
		attack := AttackEvent{}
		attack.MissionType = missionType
		if missionType == Attack || missionType == MissileAttack {
			coordsOrigin := strings.Trim(s.Find("td.coordsOrigin").Text(), " \r\t\n")
			attack.Origin = extractCoord(coordsOrigin)
		}
		destCoords := strings.Trim(s.Find("td.destCoords").Text(), " \r\t\n")
		attack.Destination = extractCoord(destCoords)

		attack.ArrivalTime = time.Unix(int64(arrivalTimeInt), 0)

		if missionType == Attack || missionType == MissileAttack {
			attackerIDStr, _ := s.Find("a.sendMail").Attr("data-playerid")
			attack.AttackerID, _ = strconv.Atoi(attackerIDStr)
		}

		if missionType == MissileAttack {
			missilesStr := s.Find("td.detailsFleet span").First().Text()
			attack.Missiles, _ = strconv.Atoi(missilesStr)
		}

		attacks = append(attacks, attack)
	}
	doc.Find("tr.eventFleet").Each(tmp)
	doc.Find("tr.llianceAttack").Each(tmp)
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
	by, _ := ioutil.ReadAll(resp.Body)
	pageHTML := string(by)

	return extractAttacks(pageHTML)
}

func extractGalaxyInfos(pageHTML, lang string) ([]PlanetInfos, error) {
	var tmp struct {
		Galaxy string
	}
	json.Unmarshal([]byte(pageHTML), &tmp)
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(tmp.Galaxy))
	res := make([]PlanetInfos, 0)
	doc.Find("tr.row").Each(func(i int, s *goquery.Selection) {
		classes, _ := s.Attr("class")
		if !strings.Contains(classes, "empty_filter") {
			activity := 0
			activityDiv := s.Find("div.activity")
			if activityDiv != nil {
				activityRaw := strings.Trim(activityDiv.Text(), " \r\t\n")
				if activityRaw != "" {
					activity, _ = strconv.Atoi(activityRaw)
				}
			}

			position := s.Find("td.position").Text()

			tooltips := s.Find("div.htmlTooltip")
			planetTooltip := tooltips.First()
			planetName := planetTooltip.Find("h1").Find("span").Text()
			planetImg, _ := planetTooltip.Find("img").Attr("src")
			coordsRaw := planetTooltip.Find("span#pos-planet").Text()

			metalRgx := regexp.MustCompile(`Metal: ([\d.]+)`)
			crystalRgx := regexp.MustCompile(`Crystal: ([\d.]+)`)
			recyclersRgx := regexp.MustCompile(`Recyclers needed: ([\d.]+)`)
			switch lang {
			case "de":
				metalRgx = regexp.MustCompile(`Metall: ([\d.]+)`)
				crystalRgx = regexp.MustCompile(`Kristall: ([\d.]+)`)
				recyclersRgx = regexp.MustCompile(`Benötigte Recycler: ([\d.]+)`)
			case "es":
				metalRgx = regexp.MustCompile(`Metal: ([\d.]+)`)
				crystalRgx = regexp.MustCompile(`Cristal: ([\d.]+)`)
				recyclersRgx = regexp.MustCompile(`Se necesitan recicladores: ([\d.]+)`)
			case "fr":
				metalRgx = regexp.MustCompile(`Métal: ([\d.]+)`)
				crystalRgx = regexp.MustCompile(`Cristal: ([\d.]+)`)
				recyclersRgx = regexp.MustCompile(`Recycleurs nécessaires: ([\d.]+)`)
			}

			metalTxt := s.Find("div#debris" + position + " ul.ListLinks li").First().Text()
			crystalTxt := s.Find("div#debris" + position + " ul.ListLinks li").Eq(1).Text()
			recyclersTxt := s.Find("div#debris" + position + " ul.ListLinks li").Eq(2).Text()

			planetInfos := PlanetInfos{}

			allianceSpan := s.Find("span.allytagwrapper")
			if allianceSpan.Size() > 0 {
				longID, _ := allianceSpan.Attr("rel")
				planetInfos.Alliance.Name = allianceSpan.Find("h1").Text()
				planetInfos.Alliance.ID, _ = strconv.Atoi(strings.TrimPrefix(longID, "alliance"))
				planetInfos.Alliance.Rank, _ = strconv.Atoi(allianceSpan.Find("ul.ListLinks li").First().Find("a").Text())
				planetInfos.Alliance.Member = parseInt(strings.TrimPrefix(allianceSpan.Find("ul.ListLinks li").Eq(1).Text(), "Member: "))
			}

			if len(metalRgx.FindStringSubmatch(metalTxt)) > 0 {
				planetInfos.Debris.Metal = parseInt(metalRgx.FindStringSubmatch(metalTxt)[1])
				planetInfos.Debris.Crystal = parseInt(crystalRgx.FindStringSubmatch(crystalTxt)[1])
				planetInfos.Debris.RecyclersNeeded = parseInt(recyclersRgx.FindStringSubmatch(recyclersTxt)[1])
			}

			planetInfos.Activity = activity
			planetInfos.Name = planetName
			planetInfos.Img = planetImg
			planetInfos.Inactive = strings.Contains(classes, "inactive_filter")
			planetInfos.StrongPlayer = strings.Contains(classes, "strong_filter")
			planetInfos.Vacation = strings.Contains(classes, "vacation_filter")
			planetInfos.HonorableTarget = s.Find("span.status_abbr_honorableTarget").Size() > 0
			planetInfos.Administrator = s.Find("span.status_abbr_admin").Size() > 0
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
				playerName = strings.Trim(s.Find("td.playername").Find("span").Text(), " \r\t\n")
			}

			planetInfos.Player.ID = playerID
			planetInfos.Player.Name = playerName
			planetInfos.Player.Rank = playerRank

			res = append(res, planetInfos)
		}
	})
	return res, nil
}

func (b *OGame) galaxyInfos(galaxy, system int) ([]PlanetInfos, error) {
	finalURL := b.serverURL + "/game/index.php?page=galaxyContent&ajax=1"
	payload := url.Values{
		"galaxy": {strconv.Itoa(galaxy)},
		"system": {strconv.Itoa(system)},
	}
	req, err := http.NewRequest("POST", finalURL, strings.NewReader(payload.Encode()))
	if err != nil {
		return []PlanetInfos{}, err
	}
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	resp, err := b.client.Do(req)
	if err != nil {
		return []PlanetInfos{}, err
	}
	defer resp.Body.Close()
	by, _ := ioutil.ReadAll(resp.Body)
	pageHTML := string(by)
	return extractGalaxyInfos(pageHTML, b.language)
}

func (b *OGame) getResourceSettings(planetID PlanetID) (ResourceSettings, error) {
	planetIDStr := planetID.String()
	pageHTML := b.getPageContent(url.Values{"page": {"resourceSettings"}, "cp": {planetIDStr}})
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(pageHTML))
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
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(pageHTML))
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

func getNbr(doc *goquery.Document, name string) (int, error) {
	div := doc.Find("div." + name)
	level := div.Find("span.level")
	return strconv.Atoi(strings.Trim(level.Contents().Text(), " \r\t\n"))
}

func (b *OGame) getResourcesBuildings(planetID PlanetID) (ResourcesBuildings, error) {
	planetIDStr := planetID.String()
	pageHTML := b.getPageContent(url.Values{"page": {"resources"}, "cp": {planetIDStr}})
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(pageHTML))
	doc.Find("span.textlabel").Remove()
	bodyID, _ := doc.Find("body").Attr("id")
	if bodyID == "overview" {
		return ResourcesBuildings{}, ErrInvalidPlanetID
	}
	res := ResourcesBuildings{}
	res.MetalMine, _ = getNbr(doc, "supply1")
	res.CrystalMine, _ = getNbr(doc, "supply2")
	res.DeuteriumSynthesizer, _ = getNbr(doc, "supply3")
	res.SolarPlant, _ = getNbr(doc, "supply4")
	res.FusionReactor, _ = getNbr(doc, "supply12")
	res.SolarSatellite, _ = getNbr(doc, "supply212")
	res.MetalStorage, _ = getNbr(doc, "supply22")
	res.CrystalStorage, _ = getNbr(doc, "supply23")
	res.DeuteriumTank, _ = getNbr(doc, "supply24")
	return res, nil
}

func extractDefense(pageHTML string) (DefensesInfos, error) {
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(pageHTML))
	bodyID, _ := doc.Find("body").Attr("id")
	if bodyID == "overview" {
		return DefensesInfos{}, ErrInvalidPlanetID
	}
	doc.Find("span.textlabel").Remove()
	res := DefensesInfos{}
	res.RocketLauncher, _ = getNbr(doc, "defense401")
	res.LightLaser, _ = getNbr(doc, "defense402")
	res.HeavyLaser, _ = getNbr(doc, "defense403")
	res.GaussCannon, _ = getNbr(doc, "defense404")
	res.IonCannon, _ = getNbr(doc, "defense405")
	res.PlasmaTurret, _ = getNbr(doc, "defense406")
	res.SmallShieldDome, _ = getNbr(doc, "defense407")
	res.LargeShieldDome, _ = getNbr(doc, "defense408")
	res.AntiBallisticMissiles, _ = getNbr(doc, "defense502")
	res.InterplanetaryMissiles, _ = getNbr(doc, "defense503")

	return res, nil
}

func (b *OGame) getDefense(planetID PlanetID) (DefensesInfos, error) {
	planetIDStr := planetID.String()
	pageHTML := b.getPageContent(url.Values{"page": {"defense"}, "cp": {planetIDStr}})
	return extractDefense(pageHTML)
}

func extractShips(pageHTML string) (ShipsInfos, error) {
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(pageHTML))
	doc.Find("span.textlabel").Remove()
	bodyID, _ := doc.Find("body").Attr("id")
	if bodyID == "overview" {
		return ShipsInfos{}, ErrInvalidPlanetID
	}
	res := ShipsInfos{}
	res.LightFighter, _ = getNbr(doc, "military204")
	res.HeavyFighter, _ = getNbr(doc, "military205")
	res.Cruiser, _ = getNbr(doc, "military206")
	res.Battleship, _ = getNbr(doc, "military207")
	res.Battlecruiser, _ = getNbr(doc, "military215")
	res.Bomber, _ = getNbr(doc, "military211")
	res.Destroyer, _ = getNbr(doc, "military213")
	res.Deathstar, _ = getNbr(doc, "military214")
	res.SmallCargo, _ = getNbr(doc, "civil202")
	res.LargeCargo, _ = getNbr(doc, "civil203")
	res.ColonyShip, _ = getNbr(doc, "civil208")
	res.Recycler, _ = getNbr(doc, "civil209")
	res.EspionageProbe, _ = getNbr(doc, "civil210")
	res.SolarSatellite, _ = getNbr(doc, "civil212")

	return res, nil
}

func (b *OGame) getShips(planetID PlanetID) (ShipsInfos, error) {
	planetIDStr := planetID.String()
	pageHTML := b.getPageContent(url.Values{"page": {"shipyard"}, "cp": {planetIDStr}})
	return extractShips(pageHTML)
}

func extractFacilities(pageHTML string) (Facilities, error) {
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(pageHTML))
	doc.Find("span.textlabel").Remove()
	bodyID, _ := doc.Find("body").Attr("id")
	if bodyID == "overview" {
		return Facilities{}, ErrInvalidPlanetID
	}
	res := Facilities{}
	res.RoboticsFactory, _ = getNbr(doc, "station14")
	res.Shipyard, _ = getNbr(doc, "station21")
	res.ResearchLab, _ = getNbr(doc, "station31")
	res.AllianceDepot, _ = getNbr(doc, "station34")
	res.MissileSilo, _ = getNbr(doc, "station44")
	res.NaniteFactory, _ = getNbr(doc, "station15")
	res.Terraformer, _ = getNbr(doc, "station33")
	res.SpaceDock, _ = getNbr(doc, "station36")

	return res, nil
}

func (b *OGame) getFacilities(planetID PlanetID) (Facilities, error) {
	planetIDStr := planetID.String()
	pageHTML := b.getPageContent(url.Values{"page": {"station"}, "cp": {planetIDStr}})
	return extractFacilities(pageHTML)
}

func extractProduction(pageHTML string) ([]Quantifiable, error) {
	res := make([]Quantifiable, 0)
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(pageHTML))
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
		itemIDstr, _ := s.Find("a").Attr("ref")
		itemID, _ := strconv.Atoi(itemIDstr)
		itemNbr := parseInt(s.Find("span.number").Text())
		res = append(res, Quantifiable{ID: ID(itemID), Nbr: itemNbr})
	})
	return res, nil
}

func (b *OGame) getProduction(planetID PlanetID) ([]Quantifiable, error) {
	planetIDStr := planetID.String()
	pageHTML := b.getPageContent(url.Values{"page": {"shipyard"}, "cp": {planetIDStr}})
	return extractProduction(pageHTML)
}

func extractResearch(pageHTML string) Researches {
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(pageHTML))
	doc.Find("span.textlabel").Remove()
	res := Researches{}
	res.EnergyTechnology, _ = getNbr(doc, "research113")
	res.LaserTechnology, _ = getNbr(doc, "research120")
	res.IonTechnology, _ = getNbr(doc, "research121")
	res.HyperspaceTechnology, _ = getNbr(doc, "research114")
	res.PlasmaTechnology, _ = getNbr(doc, "research122")
	res.CombustionDrive, _ = getNbr(doc, "research115")
	res.ImpulseDrive, _ = getNbr(doc, "research117")
	res.HyperspaceDrive, _ = getNbr(doc, "research118")
	res.EspionageTechnology, _ = getNbr(doc, "research106")
	res.ComputerTechnology, _ = getNbr(doc, "research108")
	res.Astrophysics, _ = getNbr(doc, "research124")
	res.IntergalacticResearchNetwork, _ = getNbr(doc, "research123")
	res.GravitonTechnology, _ = getNbr(doc, "research199")
	res.WeaponsTechnology, _ = getNbr(doc, "research109")
	res.ShieldingTechnology, _ = getNbr(doc, "research110")
	res.ArmourTechnology, _ = getNbr(doc, "research111")

	return res
}

func (b *OGame) getResearch() Researches {
	pageHTML := b.getPageContent(url.Values{"page": {"research"}})
	return extractResearch(pageHTML)
}

func (b *OGame) build(planetID PlanetID, id ID, nbr int) error {
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
	planetIDStr := planetID.String()
	payload := url.Values{
		"modus": {"1"},
		"type":  {strconv.Itoa(int(id))},
	}

	// Techs don't have a token
	if !id.IsTech() {
		pageHTML := b.getPageContent(url.Values{"page": {page}, "cp": {planetIDStr}})
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(pageHTML))
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

	url2 := b.serverURL + "/game/index.php?page=" + page + "&cp=" + planetIDStr
	resp, err := b.client.PostForm(url2, payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (b *OGame) buildCancelable(planetID PlanetID, id ID) error {
	if !id.IsBuilding() && !id.IsTech() {
		return errors.New("invalid id " + id.String())
	}
	return b.build(planetID, id, 0)
}

func (b *OGame) buildProduction(planetID PlanetID, id ID, nbr int) error {
	if !id.IsDefense() && !id.IsShip() {
		return errors.New("invalid id " + id.String())
	}
	return b.build(planetID, id, nbr)
}

func (b *OGame) buildBuilding(planetID PlanetID, buildingID ID) error {
	if !buildingID.IsBuilding() {
		return errors.New("invalid building id " + buildingID.String())
	}
	return b.buildCancelable(planetID, buildingID)
}

func (b *OGame) buildTechnology(planetID PlanetID, technologyID ID) error {
	if technologyID.IsTech() {
		return errors.New("invalid technology id " + technologyID.String())
	}
	return b.buildCancelable(planetID, technologyID)
}

func (b *OGame) buildDefense(planetID PlanetID, defenseID ID, nbr int) error {
	if !defenseID.IsDefense() {
		return errors.New("invalid defense id " + defenseID.String())
	}
	return b.buildProduction(planetID, ID(defenseID), nbr)
}

func (b *OGame) buildShips(planetID PlanetID, shipID ID, nbr int) error {
	if !shipID.IsShip() {
		return errors.New("invalid ship id " + shipID.String())
	}
	return b.buildProduction(planetID, ID(shipID), nbr)
}

func extractConstructions(pageHTML string) (buildingID ID, buildingCountdown int, researchID ID, researchCountdown int) {
	buildingCountdownMatch := regexp.MustCompile(`getElementByIdWithCache\("Countdown"\),(\d+),`).FindStringSubmatch(pageHTML)
	if len(buildingCountdownMatch) > 0 {
		buildingCountdown, _ = strconv.Atoi(buildingCountdownMatch[1])
		buildingIDInt, _ := strconv.Atoi(regexp.MustCompile(`onclick="cancelProduction\((\d+),`).FindStringSubmatch(pageHTML)[1])
		buildingID = ID(buildingIDInt)
	}
	researchCountdownMatch := regexp.MustCompile(`getElementByIdWithCache\("researchCountdown"\),(\d+),`).FindStringSubmatch(pageHTML)
	if len(researchCountdownMatch) > 0 {
		researchCountdown, _ = strconv.Atoi(researchCountdownMatch[1])
		researchIDInt, _ := strconv.Atoi(regexp.MustCompile(`onclick="cancelResearch\((\d+),`).FindStringSubmatch(pageHTML)[1])
		researchID = ID(researchIDInt)
	}
	return
}

func (b *OGame) constructionsBeingBuilt(planetID PlanetID) (ID, int, ID, int) {
	planetIDStr := planetID.String()
	pageHTML := b.getPageContent(url.Values{"page": {"overview"}, "cp": {planetIDStr}})
	return extractConstructions(pageHTML)
}

func extractCancelBuildingInfos(pageHTML string) (token string, techID, listID int, err error) {
	r1 := regexp.MustCompile(`page=overview&modus=2&token=(\w+)&techid="\+cancelProduction_id\+"&listid="\+production_listid`)
	m1 := r1.FindStringSubmatch(pageHTML)
	if len(m1) < 2 {
		return "", 0, 0, errors.New("unable to find token")
	}
	token = m1[1]
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(pageHTML))
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

func extractCancelResearchInfos(pageHTML string) (token string, techID, listID int, err error) {
	r1 := regexp.MustCompile(`page=overview&modus=2&token=(\w+)"\+"&techid="\+id\+"&listid="\+listId`)
	m1 := r1.FindStringSubmatch(pageHTML)
	if len(m1) < 2 {
		return "", 0, 0, errors.New("unable to find token")
	}
	token = m1[1]
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(pageHTML))
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

func (b *OGame) cancel(planetID PlanetID, token string, techID, listID int) error {
	finalURL := b.serverURL + "/game/index.php?page=overview&modus=2&token=" + token + "&techid=" + strconv.Itoa(techID) + "&listid=" + strconv.Itoa(listID)
	req, _ := http.NewRequest("GET", finalURL, nil)
	resp, _ := b.client.Do(req)
	defer resp.Body.Close()
	return nil
}

func (b *OGame) cancelBuilding(planetID PlanetID) error {
	planetIDStr := planetID.String()
	pageHTML := b.getPageContent(url.Values{"page": {"overview"}, "cp": {planetIDStr}})
	token, techID, listID, _ := extractCancelBuildingInfos(pageHTML)
	return b.cancel(planetID, token, techID, listID)
}

func (b *OGame) cancelResearch(planetID PlanetID) error {
	planetIDStr := planetID.String()
	pageHTML := b.getPageContent(url.Values{"page": {"overview"}, "cp": {planetIDStr}})
	token, techID, listID, _ := extractCancelResearchInfos(pageHTML)
	return b.cancel(planetID, token, techID, listID)
}

func (b *OGame) fetchResources(planetID PlanetID) (resourcesResp, error) {
	planetIDStr := planetID.String()
	pageJSON := b.getPageContent(url.Values{"page": {"fetchResources"}, "cp": {planetIDStr}})
	var res resourcesResp
	if err := json.Unmarshal([]byte(pageJSON), &res); err != nil {
		if isLogged(pageJSON) {
			return resourcesResp{}, ErrInvalidPlanetID
		}
		return resourcesResp{}, err
	}
	return res, nil
}

func (b *OGame) getResources(planetID PlanetID) (Resources, error) {
	res, err := b.fetchResources(planetID)
	return Resources{
		Metal:      res.Metal.Resources.Actual,
		Crystal:    res.Crystal.Resources.Actual,
		Deuterium:  res.Deuterium.Resources.Actual,
		Energy:     res.Energy.Resources.Actual,
		Darkmatter: res.Darkmatter.Resources.Actual,
	}, err
}

func (b *OGame) sendFleet(planetID PlanetID, ships []Quantifiable, speed Speed, where Coordinate, mission MissionID,
	resources Resources) (FleetID, error) {
	getHiddenFields := func(pageHTML string) map[string]string {
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(pageHTML))
		fields := make(map[string]string)
		doc.Find("input[type=hidden]").Each(func(i int, s *goquery.Selection) {
			name, _ := s.Attr("name")
			value, _ := s.Attr("value")
			fields[name] = value
		})
		return fields
	}
	planetIDStr := planetID.String()
	pageHTML := b.getPageContent(url.Values{"page": {"fleet1"}, "cp": {planetIDStr}})

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
	fleet2URL := b.serverURL + "/game/index.php?page=fleet2"
	fleet2Resp, err := b.client.PostForm(fleet2URL, payload)
	if err != nil {
		return 0, err
	}
	defer fleet2Resp.Body.Close()
	by, _ := ioutil.ReadAll(fleet2Resp.Body)
	pageHTML = string(by)

	payload = url.Values{}
	hidden = getHiddenFields(pageHTML)
	for k, v := range hidden {
		payload.Add(k, v)
	}
	payload.Add("speed", speed.String())
	payload.Add("galaxy", strconv.Itoa(where.Galaxy))
	payload.Add("system", strconv.Itoa(where.System))
	payload.Add("position", strconv.Itoa(where.Position))
	t := "1"
	if mission == RecycleDebrisField {
		// planet type: 1
		// debris type: 2
		// moon type: 3
		t = "2" // Send to debris field
	}
	payload.Add("type", t)

	// Check
	fleetCheckURL := b.serverURL + "/game/index.php?page=fleetcheck&ajax=1&espionage=0"
	fleetCheckPayload := url.Values{
		"galaxy": {strconv.Itoa(where.Galaxy)},
		"system": {strconv.Itoa(where.System)},
		"planet": {strconv.Itoa(where.Position)},
		"type":   {t},
	}
	fleetCheckResp, err := b.client.PostForm(fleetCheckURL, fleetCheckPayload)
	if err != nil {
		return 0, err
	}
	defer fleetCheckResp.Body.Close()
	by1, _ := ioutil.ReadAll(fleetCheckResp.Body)
	switch string(by1) {
	case "1":
		return 0, errors.New("uninhabited planet")
	case "1d":
		return 0, errors.New("no debris field")
	case "2":
		return 0, errors.New("player in vacation mode")
	case "3":
		return 0, errors.New("admin or GM")
	case "4":
		return 0, errors.New("you have to research Astrophysics first")
	case "5":
		return 0, errors.New("noob protection")
	case "6":
		return 0, errors.New("this planet can not be attacked as the player is to strong")
	case "10":
		return 0, errors.New("no moon available")
	case "11":
		return 0, errors.New("no recycler available")
	case "15":
		return 0, errors.New("there are currently no events running")
	case "16":
		return 0, errors.New("this planet has already been reserved for a relocation")
	}

	fleet3URL := b.serverURL + "/game/index.php?page=fleet3"
	fleet3Resp, err := b.client.PostForm(fleet3URL, payload)
	if err != nil {
		return 0, err
	}
	defer fleet3Resp.Body.Close()
	by, _ = ioutil.ReadAll(fleet3Resp.Body)
	pageHTML = string(by)

	payload = url.Values{}
	hidden = getHiddenFields(pageHTML)
	for k, v := range hidden {
		payload.Add(k, v)
	}
	payload.Add("crystal", strconv.Itoa(resources.Crystal))
	payload.Add("deuterium", strconv.Itoa(resources.Deuterium))
	payload.Add("metal", strconv.Itoa(resources.Metal))
	payload.Add("mission", mission.String())
	movementURL := b.serverURL + "/game/index.php?page=movement"
	movementResp, _ := b.client.PostForm(movementURL, payload)
	defer movementResp.Body.Close()

	movementHTML := b.getPageContent(url.Values{"page": {"movement"}})
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(movementHTML))
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

// EspionageReportType ...
type EspionageReportType int

// Action ...
const Action EspionageReportType = 0

// Report ...
const Report EspionageReportType = 1

// EspionageReportSummary ...
type EspionageReportSummary struct {
	ID     int
	Type   EspionageReportType
	From   string
	Target Coordinate
}

func extractEspionageReportMessageIDs(pageHTML string) ([]EspionageReportSummary, int) {
	msgs := make([]EspionageReportSummary, 0)
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(pageHTML))
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

func (b *OGame) getEspionageReportMessages() ([]EspionageReportSummary, error) {
	page := 1
	nbPage := 1
	msgs := make([]EspionageReportSummary, 0)
	finalURL := b.serverURL + "/game/index.php?page=messages"
	for page <= nbPage {
		payload := url.Values{
			"messageId":  {"-1"},
			"tabid":      {"20"},
			"action":     {"107"},
			"pagination": {strconv.Itoa(page)},
			"ajax":       {"1"},
		}
		req, err := http.NewRequest("POST", finalURL, strings.NewReader(payload.Encode()))
		if err != nil {
			return []EspionageReportSummary{}, err
		}
		req.Header.Add("X-Requested-With", "XMLHttpRequest")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
		resp, err := b.client.Do(req)
		if err != nil {
			return []EspionageReportSummary{}, err
		}
		by, _ := ioutil.ReadAll(resp.Body)
		pageHTML := string(by)
		resp.Body.Close()
		newMessages, newNbPage := extractEspionageReportMessageIDs(pageHTML)
		msgs = append(msgs, newMessages...)
		nbPage = newNbPage
		page++
	}
	return msgs, nil
}

// EspionageReport ...
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

func extractEspionageReport(pageHTML string, location *time.Location) (EspionageReport, error) {
	report := EspionageReport{}
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(pageHTML))
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
				l, _ := strconv.Atoi(s2.Find("span.fright").Text())
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
				}
			})
		} else if dataType == "research" {
			s.Find("li.detail_list_el").Each(func(i int, s2 *goquery.Selection) {
				imgClass := s2.Find("img").AttrOr("class", "")
				r := regexp.MustCompile(`research(\d+)`)
				researchID, _ := strconv.Atoi(r.FindStringSubmatch(imgClass)[1])
				l, _ := strconv.Atoi(s2.Find("span.fright").Text())
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
				l, _ := strconv.Atoi(s2.Find("span.fright").Text())
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
				//case SolarSatellite:
				//	report.SolarSatellite = level
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
				l, _ := strconv.Atoi(s2.Find("span.fright").Text())
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

func energyProduced(maxTemp int, resourcesBuildings ResourcesBuildings, resSettings ResourceSettings, energyTechnology, solarSatellite int) int {
	energyProduced := int(float64(SolarPlant.Production(resourcesBuildings.SolarPlant)) * (float64(resSettings.SolarPlant) / 100))
	energyProduced += int(float64(FusionReactor.Production(energyTechnology, resourcesBuildings.FusionReactor)) * (float64(resSettings.FusionReactor) / 100))
	energyProduced += int(float64(SolarSatellite.Production(maxTemp, solarSatellite)) * (float64(resSettings.SolarSatellite) / 100))
	return energyProduced
}

func energyNeeded(resourcesBuildings ResourcesBuildings, resSettings ResourceSettings) int {
	energyNeeded := int(float64(MetalMine.EnergyConsumption(resourcesBuildings.MetalMine)) * (float64(resSettings.MetalMine) / 100))
	energyNeeded += int(float64(CrystalMine.EnergyConsumption(resourcesBuildings.CrystalMine)) * (float64(resSettings.CrystalMine) / 100))
	energyNeeded += int(float64(DeuteriumSynthesizer.EnergyConsumption(resourcesBuildings.DeuteriumSynthesizer)) * (float64(resSettings.DeuteriumSynthesizer) / 100))
	return energyNeeded
}

func productionRatio(maxTemp int, resourcesBuildings ResourcesBuildings, resSettings ResourceSettings, energyTechnology int) float64 {
	energyProduced := energyProduced(maxTemp, resourcesBuildings, resSettings, energyTechnology, resourcesBuildings.SolarSatellite)
	energyNeeded := energyNeeded(resourcesBuildings, resSettings)
	ratio := 1.0
	if energyNeeded > energyProduced {
		ratio = float64(energyProduced) / float64(energyNeeded)
	}
	return ratio
}

func getProductions(resBuildings ResourcesBuildings, resSettings ResourceSettings, researches Researches, universeSpeed,
	maxTemp int, productionRatio float64) Resources {
	energyProduced := energyProduced(maxTemp, resBuildings, resSettings, researches.EnergyTechnology, resBuildings.SolarSatellite)
	energyNeeded := energyNeeded(resBuildings, resSettings)
	return Resources{
		Metal:     MetalMine.Production(universeSpeed, productionRatio, resBuildings.MetalMine),
		Crystal:   CrystalMine.Production(universeSpeed, productionRatio, resBuildings.CrystalMine),
		Deuterium: DeuteriumSynthesizer.Production(universeSpeed, maxTemp, productionRatio, resBuildings.DeuteriumSynthesizer),
		Energy:    energyProduced - energyNeeded,
	}
}

func extractResourcesProductions(pageHTML string) (Resources, error) {
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(pageHTML))
	res := Resources{}
	selector := "table.listOfResourceSettingsPerPlanet tr.summary td span"
	res.Metal = parseInt(strings.Trim(doc.Find(selector).Eq(0).AttrOr("title", "0"), "\n\t\r "))
	res.Crystal = parseInt(strings.Trim(doc.Find(selector).Eq(1).AttrOr("title", "0"), "\n\t\r "))
	res.Deuterium = parseInt(strings.Trim(doc.Find(selector).Eq(2).AttrOr("title", "0"), "\n\t\r "))
	res.Energy = parseInt(strings.Trim(doc.Find(selector).Eq(3).AttrOr("title", "0"), "\n\t\r "))
	return res, nil
}

func (b *OGame) getResourcesProductions(planetID PlanetID) (Resources, error) {
	planet, _ := b.getPlanet(planetID)
	resBuildings, _ := b.getResourcesBuildings(planetID)
	researches := b.getResearch()
	universeSpeed := b.getUniverseSpeed()
	resSettings, _ := b.getResourceSettings(planetID)
	ratio := productionRatio(planet.Temperature.Max, resBuildings, resSettings, researches.EnergyTechnology)
	productions := getProductions(resBuildings, resSettings, researches, universeSpeed, planet.Temperature.Max, ratio)
	return productions, nil
}

// ServerURL ...
func (b *OGame) ServerURL() string {
	return b.serverURL
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

// Logout ...
func (b *OGame) Logout() {
	b.Lock()
	defer b.Unlock()
	b.logout()
}

// GetUniverseName ...
func (b *OGame) GetUniverseName() string {
	return b.Universe
}

// GetUsername ...
func (b *OGame) GetUsername() string {
	return b.Username
}

// GetUniverseSpeed ...
func (b *OGame) GetUniverseSpeed() int {
	return b.getUniverseSpeed()
}

// GetUniverseSpeedFleet ...
func (b *OGame) GetUniverseSpeedFleet() int {
	return b.getUniverseSpeedFleet()
}

// IsDonutGalaxy ...
func (b *OGame) IsDonutGalaxy() bool {
	return b.isDonutGalaxy()
}

// IsDonutSystem ...
func (b *OGame) IsDonutSystem() bool {
	return b.isDonutSystem()
}

// GetPageContent ...
func (b *OGame) GetPageContent(vals url.Values) string {
	b.Lock()
	defer b.Unlock()
	return b.getPageContent(vals)
}

// PostPageContent ...
func (b *OGame) PostPageContent(vals, payload url.Values) string {
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

// GetPlanets returns the user planets ids
func (b *OGame) GetPlanets() []Planet {
	b.Lock()
	defer b.Unlock()
	return b.getPlanets()
}

// GetPlanet gets infos for planetID
// Fails if planetID is invalid
func (b *OGame) GetPlanet(planetID PlanetID) (Planet, error) {
	b.Lock()
	defer b.Unlock()
	return b.getPlanet(planetID)
}

// GetPlanetByCoord ...
func (b *OGame) GetPlanetByCoord(coord Coordinate) (Planet, error) {
	b.Lock()
	defer b.Unlock()
	return b.getPlanetByCoord(coord)
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

// GetUserInfos gets the user informations
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

// GetFleets ...
func (b *OGame) GetFleets() []Fleet {
	b.Lock()
	defer b.Unlock()
	return b.getFleets()
}

// CancelFleet ...
func (b *OGame) CancelFleet(fleetID FleetID) error {
	b.Lock()
	defer b.Unlock()
	return b.cancelFleet(fleetID)
}

// GetAttacks ...
func (b *OGame) GetAttacks() []AttackEvent {
	b.Lock()
	defer b.Unlock()
	return b.getAttacks()
}

// GalaxyInfos ...
func (b *OGame) GalaxyInfos(galaxy, system int) ([]PlanetInfos, error) {
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

// SetResourceSettings ...
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

// GetDefense ...
// Fails if planetID is invalid
func (b *OGame) GetDefense(planetID PlanetID) (DefensesInfos, error) {
	b.Lock()
	defer b.Unlock()
	return b.getDefense(planetID)
}

// GetShips ...
func (b *OGame) GetShips(planetID PlanetID) (ShipsInfos, error) {
	b.Lock()
	defer b.Unlock()
	return b.getShips(planetID)
}

// GetFacilities ...
func (b *OGame) GetFacilities(planetID PlanetID) (Facilities, error) {
	b.Lock()
	defer b.Unlock()
	return b.getFacilities(planetID)
}

// GetProduction get what is in the production queue.
// (ships & defense being built)
func (b *OGame) GetProduction(planetID PlanetID) ([]Quantifiable, error) {
	b.Lock()
	defer b.Unlock()
	return b.getProduction(planetID)
}

// GetResearch ...
func (b *OGame) GetResearch() Researches {
	b.Lock()
	defer b.Unlock()
	return b.getResearch()
}

// Build ...
func (b *OGame) Build(planetID PlanetID, id ID, nbr int) error {
	b.Lock()
	defer b.Unlock()
	return b.build(planetID, id, nbr)
}

// BuildCancelable ...
func (b *OGame) BuildCancelable(planetID PlanetID, id ID) error {
	b.Lock()
	defer b.Unlock()
	return b.buildCancelable(planetID, id)
}

// BuildProduction ...
func (b *OGame) BuildProduction(planetID PlanetID, id ID, nbr int) error {
	b.Lock()
	defer b.Unlock()
	return b.buildProduction(planetID, id, nbr)
}

// BuildBuilding ...
func (b *OGame) BuildBuilding(planetID PlanetID, buildingID ID) error {
	b.Lock()
	defer b.Unlock()
	return b.buildBuilding(planetID, buildingID)
}

// BuildDefense builds a defense unit
func (b *OGame) BuildDefense(planetID PlanetID, defenseID ID, nbr int) error {
	b.Lock()
	defer b.Unlock()
	return b.buildDefense(planetID, defenseID, nbr)
}

// BuildShips builds a ship unit
func (b *OGame) BuildShips(planetID PlanetID, shipID ID, nbr int) error {
	b.Lock()
	defer b.Unlock()
	return b.buildShips(planetID, shipID, nbr)
}

// ConstructionsBeingBuilt ...
func (b *OGame) ConstructionsBeingBuilt(planetID PlanetID) (ID, int, ID, int) {
	b.Lock()
	defer b.Unlock()
	return b.constructionsBeingBuilt(planetID)
}

// CancelBuilding ...
func (b *OGame) CancelBuilding(planetID PlanetID) error {
	b.Lock()
	defer b.Unlock()
	return b.cancelBuilding(planetID)
}

// CancelResearch ...
func (b *OGame) CancelResearch(planetID PlanetID) error {
	b.Lock()
	defer b.Unlock()
	return b.cancelResearch(planetID)
}

// BuildTechnology ...
func (b *OGame) BuildTechnology(planetID PlanetID, technologyID ID) error {
	b.Lock()
	defer b.Unlock()
	return b.buildTechnology(planetID, technologyID)
}

// GetResources gets user resources
func (b *OGame) GetResources(planetID PlanetID) (Resources, error) {
	b.Lock()
	defer b.Unlock()
	return b.getResources(planetID)
}

// SendFleet ...
func (b *OGame) SendFleet(planetID PlanetID, ships []Quantifiable, speed Speed, where Coordinate, mission MissionID,
	resources Resources) (FleetID, error) {
	b.Lock()
	defer b.Unlock()
	return b.sendFleet(planetID, ships, speed, where, mission, resources)
}

// GetEspionageReportMessages ...
func (b *OGame) GetEspionageReportMessages() ([]EspionageReportSummary, error) {
	b.Lock()
	defer b.Unlock()
	return b.getEspionageReportMessages()
}

// GetEspionageReport ...
func (b *OGame) GetEspionageReport(msgID int) (EspionageReport, error) {
	b.Lock()
	defer b.Unlock()
	return b.getEspionageReport(msgID)
}

// DeleteMessage ...
func (b *OGame) DeleteMessage(msgID int) error {
	b.Lock()
	defer b.Unlock()
	return b.deleteMessage(msgID)
}

// GetResourcesProductions ...
func (b *OGame) GetResourcesProductions(planetID PlanetID) (Resources, error) {
	b.Lock()
	defer b.Unlock()
	return b.getResourcesProductions(planetID)
}
