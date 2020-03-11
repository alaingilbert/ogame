package ogame

import (
	"bytes"
	"compress/gzip"
	"container/heap"
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/hashicorp/go-version"
	cookiejar "github.com/orirawlings/persistent-cookiejar"
	"github.com/pkg/errors"
	lua "github.com/yuin/gopher-lua"
	"golang.org/x/net/proxy"
	"golang.org/x/net/websocket"
)

// OGame is a client for ogame.org. It is safe for concurrent use by
// multiple goroutines (thread-safe)
type OGame struct {
	sync.Mutex
	isEnabledAtom         int32  // atomic, prevent auto re login if we manually logged out
	isLoggedInAtom        int32  // atomic, prevent auto re login if we manually logged out
	isConnectedAtom       int32  // atomic, either or not communication between the bot and OGame is possible
	lockedAtom            int32  // atomic, bot state locked/unlocked
	chatConnectedAtom     int32  // atomic, either or not the chat is connected
	state                 string // keep name of the function that currently lock the bot
	stateChangeCallbacks  []func(locked bool, actor string)
	quiet                 bool
	Player                UserInfos
	CachedPreferences     Preferences
	isVacationModeEnabled bool
	researches            *Researches
	Planets               []Planet
	ajaxChatToken         string
	Universe              string
	Username              string
	password              string
	language              string
	playerID              int64
	lobby                 string
	ogameSession          string
	sessionChatCounter    int64
	server                Server
	serverData            ServerData
	location              *time.Location
	serverURL             string
	Client                *OGameClient
	logger                *log.Logger
	chatCallbacks         []func(msg ChatMsg)
	auctioneerCallbacks   []func(packet []byte)
	interceptorCallbacks  []func(method, url string, params, payload url.Values, pageHTML []byte)
	closeChatCh           chan struct{}
	chatRetry             *ExponentialBackoff
	ws                    *websocket.Conn
	tasks                 priorityQueue
	tasksLock             sync.Mutex
	tasksPushCh           chan *item
	tasksPopCh            chan struct{}
	loginWrapper          func(func() error) error
	loginProxyTransport   http.RoundTripper
	bytesUploaded         int64
	bytesDownloaded       int64
	extractor             Extractor
	apiNewHostname        string
	characterClass        CharacterClass
	hasCommander          bool
	hasAdmiral            bool
	hasEngineer           bool
	hasGeologist          bool
	hasTechnocrat         bool
}

// Preferences ...
type Preferences struct {
	SpioAnz                      int64
	DisableChatBar               bool // no-mobile
	DisableOutlawWarning         bool
	MobileVersion                bool
	ShowOldDropDowns             bool
	ActivateAutofocus            bool
	EventsShow                   int64 // Hide: 1, Above the content: 2, Below the content: 3
	SortSetting                  int64 // Order of emergence: 0, Coordinates: 1, Alphabet: 2, Size: 3, Used fields: 4
	SortOrder                    int64 // Up: 0, Down: 1
	ShowDetailOverlay            bool
	AnimatedSliders              bool // no-mobile
	AnimatedOverview             bool // no-mobile
	PopupsNotices                bool // no-mobile
	PopopsCombatreport           bool // no-mobile
	SpioReportPictures           bool
	MsgResultsPerPage            int64 // 10, 25, 50
	AuctioneerNotifications      bool
	EconomyNotifications         bool
	ShowActivityMinutes          bool
	PreserveSystemOnPlanetChange bool

	// Mobile only
	Notifications struct {
		BuildList               bool
		FriendlyFleetActivities bool
		HostileFleetActivities  bool
		ForeignEspionage        bool
		AllianceBroadcasts      bool
		AllianceMessages        bool
		Auctions                bool
		Account                 bool
	}
}

const defaultUserAgent = "" +
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) " +
	"AppleWebKit/537.36 (KHTML, like Gecko) " +
	"Chrome/79.0.3945.130 " +
	"Safari/537.36"

type options struct {
	SkipInterceptor bool
}

// Option functions to be passed to public interface to change behaviors
type Option func(*options)

// SkipInterceptor option to skip html interceptors
func SkipInterceptor(opt *options) {
	opt.SkipInterceptor = true
}

// CelestialID represent either a PlanetID or a MoonID
type CelestialID int64

// Params parameters for more fine-grained initialization
type Params struct {
	Username        string
	Password        string
	Universe        string
	Lang            string
	PlayerID        int64
	AutoLogin       bool
	Proxy           string
	ProxyUsername   string
	ProxyPassword   string
	ProxyType       string
	ProxyLoginOnly  bool
	Lobby           string
	APINewHostname  string
	CookiesFilename string
}

// New creates a new instance of OGame wrapper.
func New(universe, username, password, lang string) (*OGame, error) {
	b := NewNoLogin(username, password, universe, lang, "", 0)
	if err := b.LoginWithExistingCookies(); err != nil {
		return nil, err
	}
	return b, nil
}

// NewWithParams create a new OGame instance with full control over the possible parameters
func NewWithParams(params Params) (*OGame, error) {
	b := NewNoLogin(params.Username, params.Password, params.Universe, params.Lang, params.CookiesFilename, params.PlayerID)
	b.setOGameLobby(params.Lobby)
	b.apiNewHostname = params.APINewHostname
	if params.Proxy != "" {
		if err := b.SetProxy(params.Proxy, params.ProxyUsername, params.ProxyPassword, params.ProxyType, params.ProxyLoginOnly); err != nil {
			return nil, err
		}
	}
	if params.AutoLogin {
		if err := b.LoginWithExistingCookies(); err != nil {
			return nil, err
		}
	}
	return b, nil
}

// NewNoLogin does not auto login.
func NewNoLogin(username, password, universe, lang, cookiesFilename string, playerID int64) *OGame {
	b := new(OGame)
	b.loginWrapper = DefaultLoginWrapper
	b.Enable()
	b.quiet = false
	b.logger = log.New(os.Stdout, "", 0)

	b.Universe = universe
	b.SetOGameCredentials(username, password)
	b.setOGameLobby("lobby")
	b.language = lang
	b.playerID = playerID

	b.extractor = NewExtractorV71()

	jar, _ := cookiejar.New(&cookiejar.Options{
		Filename:              cookiesFilename,
		PersistSessionCookies: true,
	})
	b.Client = NewOGameClient()
	b.Client.Jar = jar
	b.Client.UserAgent = defaultUserAgent

	b.tasks = make(priorityQueue, 0)
	heap.Init(&b.tasks)
	b.tasksPushCh = make(chan *item, 100)
	b.tasksPopCh = make(chan struct{}, 100)
	b.taskRunner()

	return b
}

// Server ogame information for their servers
type Server struct {
	Language      string
	Number        int64
	Name          string
	PlayerCount   int64
	PlayersOnline int64
	Opened        string
	StartDate     string
	EndDate       *string
	ServerClosed  int64
	Prefered      int64
	SignupClosed  int64
	Settings      struct {
		AKS                      int64
		FleetSpeed               int64
		WreckField               int64
		ServerLabel              string
		EconomySpeed             int64
		PlanetFields             int64
		UniverseSize             int64 // Nb of galaxies
		ServerCategory           string
		EspionageProbeRaids      int64
		PremiumValidationGift    int64
		DebrisFieldFactorShips   int64
		DebrisFieldFactorDefence int64
	}
}

// ogame cookie name for php session id
const phpSessionIDCookieName = "PHPSESSID"

func getPhpSessionID(b *OGame, username, password string) (string, error) {
	payload := url.Values{
		"kid":                   {""},
		"language":              {"en"},
		"autologin":             {"false"},
		"credentials[email]":    {username},
		"credentials[password]": {password},
	}
	req, err := http.NewRequest("POST", "https://"+b.lobby+".ogame.gameforge.com/api/users", strings.NewReader(payload.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := b.doReqWithLoginProxyTransport(req)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			b.error(err)
		}
	}()

	if resp.StatusCode >= 500 {
		return "", errors.New("OGame server error code : " + resp.Status)
	}

	if resp.StatusCode != 200 {
		by, err := ioutil.ReadAll(resp.Body)
		b.error(resp.StatusCode, string(by), err)
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
		Number   int64
	}
	ID         int64 // player ID
	Name       string
	LastPlayed string
	Blocked    bool
	Details    []struct {
		Type  string
		Title string
		Value interface{} // Can be string or int
	}
	Sitting struct {
		Shared       bool
		EndTime      *string
		CooldownTime *string
	}
}

func getUserAccounts(b *OGame) ([]account, error) {
	var userAccounts []account
	req, err := http.NewRequest("GET", "https://"+b.lobby+".ogame.gameforge.com/api/users/me/accounts", nil)
	if err != nil {
		return userAccounts, err
	}
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	resp, err := b.Client.Do(req)
	if err != nil {
		return userAccounts, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			b.error(err)
		}
	}()
	by, err := readBody(b, resp)
	if err != nil {
		return userAccounts, err
	}
	b.bytesUploaded += req.ContentLength
	if err := json.Unmarshal(by, &userAccounts); err != nil {
		return userAccounts, err
	}
	return userAccounts, nil
}

func getServers(b *OGame) ([]Server, error) {
	var servers []Server
	req, err := http.NewRequest("GET", "https://"+b.lobby+".ogame.gameforge.com/api/servers", nil)
	if err != nil {
		return servers, err
	}
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	resp, err := b.Client.Do(req)
	if err != nil {
		return servers, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			b.error(err)
		}
	}()
	by, err := readBody(b, resp)
	if err != nil {
		return servers, err
	}
	b.bytesUploaded += req.ContentLength
	if err := json.Unmarshal(by, &servers); err != nil {
		return servers, err
	}
	return servers, nil
}

func findAccount(universe, lang string, playerID int64, accounts []account, servers []Server) (account, Server, error) {
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
			if playerID != 0 {
				if a.ID == playerID {
					acc = a
					break
				}
			} else {
				acc = a
				break
			}
		}
	}
	if server.Number == 0 {
		return account{}, Server{}, fmt.Errorf("server %s, %s not found", universe, lang)
	}
	if acc.ID == 0 {
		return account{}, Server{}, ErrAccountNotFound
	}
	return acc, server, nil
}

func execLoginLink(b *OGame, loginLink string) ([]byte, error) {
	req, err := http.NewRequest("GET", loginLink, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	b.debug("login to universe")
	resp, err := b.doReqWithLoginProxyTransport(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			b.error(err)
		}
	}()
	b.bytesUploaded += req.ContentLength
	return readBody(b, resp)
}

func readBody(b *OGame, resp *http.Response) ([]byte, error) {
	n := int64(0)
	defer func() {
		b.bytesDownloaded += n
	}()
	isGzip := false
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		isGzip = true
		n = resp.ContentLength
		var err error
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return []byte{}, err
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}
	by, err := ioutil.ReadAll(reader)
	if err != nil {
		return []byte{}, err
	}
	if !isGzip {
		n = int64(len(by))
	}
	return by, nil
}

func getLoginLink(b *OGame, userAccount account) (string, error) {
	ogURL := fmt.Sprintf("https://"+b.lobby+".ogame.gameforge.com/api/users/me/loginLink?id=%d&server[language]=%s&server[number]=%d",
		userAccount.ID, userAccount.Server.Language, userAccount.Server.Number)
	req, err := http.NewRequest("GET", ogURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	resp, err := b.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			b.error(err)
		}
	}()
	by, err := readBody(b, resp)
	if err != nil {
		return "", err
	}
	b.bytesUploaded += req.ContentLength
	var loginLink struct {
		URL string
	}
	if err := json.Unmarshal(by, &loginLink); err != nil {
		return "", err
	}
	return loginLink.URL, nil
}

// ServerData represent api result from https://s157-ru.ogame.gameforge.com/api/serverData.xml
type ServerData struct {
	Name                          string  `xml:"name"`                          // Europa
	Number                        int64   `xml:"number"`                        // 157
	Language                      string  `xml:"language"`                      // ru
	Timezone                      string  `xml:"timezone"`                      // Europe/Moscow
	TimezoneOffset                string  `xml:"timezoneOffset"`                // +03:00
	Domain                        string  `xml:"domain"`                        // s157-ru.ogame.gameforge.com
	Version                       string  `xml:"version"`                       // 6.8.8-pl2
	Speed                         int64   `xml:"speed"`                         // 6
	SpeedFleet                    int64   `xml:"speedFleet"`                    // 6
	Galaxies                      int64   `xml:"galaxies"`                      // 4
	Systems                       int64   `xml:"systems"`                       // 499
	ACS                           bool    `xml:"aCS"`                           // 1
	RapidFire                     bool    `xml:"rapidFire"`                     // 1
	DefToTF                       bool    `xml:"defToTF"`                       // 0
	DebrisFactor                  float64 `xml:"debrisFactor"`                  // 0.5
	DebrisFactorDef               float64 `xml:"debrisFactorDef"`               // 0
	RepairFactor                  float64 `xml:"repairFactor"`                  // 0.7
	NewbieProtectionLimit         int64   `xml:"newbieProtectionLimit"`         // 500000
	NewbieProtectionHigh          int64   `xml:"newbieProtectionHigh"`          // 50000
	TopScore                      int64   `xml:"topScore"`                      // 60259362
	BonusFields                   int64   `xml:"bonusFields"`                   // 30
	DonutGalaxy                   bool    `xml:"donutGalaxy"`                   // 1
	DonutSystem                   bool    `xml:"donutSystem"`                   // 1
	WfEnabled                     bool    `xml:"wfEnabled"`                     // 1 (WreckField)
	WfMinimumRessLost             int64   `xml:"wfMinimumRessLost"`             // 150000
	WfMinimumLossPercentage       int64   `xml:"wfMinimumLossPercentage"`       // 5
	WfBasicPercentageRepairable   int64   `xml:"wfBasicPercentageRepairable"`   // 45
	GlobalDeuteriumSaveFactor     float64 `xml:"globalDeuteriumSaveFactor"`     // 0.5
	Bashlimit                     int64   `xml:"bashlimit"`                     // 0
	ProbeCargo                    int64   `xml:"probeCargo"`                    // 5
	ResearchDurationDivisor       int64   `xml:"researchDurationDivisor"`       // 2
	DarkMatterNewAcount           int64   `xml:"darkMatterNewAcount"`           // 8000
	CargoHyperspaceTechMultiplier int64   `xml:"cargoHyperspaceTechMultiplier"` // 5
}

// gets the server data from xml api
func (b *OGame) getServerData() (ServerData, error) {
	var serverData ServerData
	req, err := http.NewRequest("GET", "https://s"+strconv.FormatInt(b.server.Number, 10)+"-"+b.server.Language+".ogame.gameforge.com/api/serverData.xml", nil)
	if err != nil {
		return serverData, err
	}
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	resp, err := b.Client.Do(req)
	if err != nil {
		return serverData, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			b.error(err)
		}
	}()
	by, err := readBody(b, resp)
	if err != nil {
		return serverData, err
	}
	b.bytesUploaded += req.ContentLength
	if err := xml.Unmarshal(by, &serverData); err != nil {
		return serverData, err
	}
	return serverData, nil
}

func (b *OGame) loginWithExistingCookies() error {
	cookies := b.Client.Jar.(*cookiejar.Jar).AllCookies()
	found := false
	for _, c := range cookies {
		if c.Name == phpSessionIDCookieName {
			found = true
			break
		}
	}
	if !found {
		return b.login()
	}
	server, userAccount, err := b.loginPart1()
	if err != nil {
		return b.login()
	}

	if err := b.loginPart2(server, userAccount); err != nil {
		return err
	}

	pageHTML, err := b.getPage(OverviewPage, CelestialID(0))
	if err != nil {
		return err
	}
	b.debug("login using existing cookies")
	if err := b.loginPart3(userAccount, pageHTML); err != nil {
		return err
	}
	return nil
}

func (b *OGame) login() error {
	b.debug("get session")
	if _, err := getPhpSessionID(b, b.Username, b.password); err != nil {
		return err
	}

	server, userAccount, err := b.loginPart1()
	if err != nil {
		return err
	}

	b.debug("get login link")
	loginLink, err := getLoginLink(b, userAccount)
	if err != nil {
		return err
	}
	pageHTML, err := execLoginLink(b, loginLink)
	if err != nil {
		return err
	}

	if err := b.loginPart2(server, userAccount); err != nil {
		return err
	}
	if err := b.loginPart3(userAccount, pageHTML); err != nil {
		return err
	}

	if err := b.Client.Jar.(*cookiejar.Jar).Save(); err != nil {
		return err
	}
	for _, fn := range b.interceptorCallbacks {
		fn("GET", loginLink, nil, nil, pageHTML)
	}
	return nil
}

func (b *OGame) loginPart1() (server Server, userAccount account, err error) {
	b.debug("get user accounts")
	accounts, err := getUserAccounts(b)
	if err != nil {
		return
	}
	b.debug("get servers")
	servers, err := getServers(b)
	if err != nil {
		return
	}
	b.debug("find account & server for universe")
	userAccount, server, err = findAccount(b.Universe, b.language, b.playerID, accounts, servers)
	if err != nil {
		return
	}
	if userAccount.Blocked {
		return server, userAccount, ErrAccountBlocked
	}
	b.debug("Players online: " + strconv.FormatInt(server.PlayersOnline, 10) + ", Players: " + strconv.FormatInt(server.PlayerCount, 10))
	return
}

func (b *OGame) loginPart2(server Server, userAccount account) error {
	atomic.StoreInt32(&b.isLoggedInAtom, 1) // At this point, we are logged in
	atomic.StoreInt32(&b.isConnectedAtom, 1)
	// Get server data
	start := time.Now()
	b.server = server
	serverData, err := b.getServerData()
	if err != nil {
		return err
	}
	b.serverData = serverData
	b.language = userAccount.Server.Language
	b.serverURL = "https://s" + strconv.FormatInt(server.Number, 10) + "-" + server.Language + ".ogame.gameforge.com"
	b.debug("get server data", time.Since(start))
	return nil
}

func (b *OGame) loginPart3(userAccount account, pageHTML []byte) error {
	if ogVersion, err := version.NewVersion(b.serverData.Version); err == nil {
		if ogVersion.GreaterThanOrEqual(version.Must(version.NewVersion("7.1.0-rc0"))) {
			b.extractor = NewExtractorV71()
		} else if ogVersion.GreaterThanOrEqual(version.Must(version.NewVersion("7.0.0-rc0"))) {
			b.extractor = NewExtractorV7()
		}
	} else {
		b.error("failed to parse ogame version: " + err.Error())
	}

	b.sessionChatCounter = 1

	b.debug("logged in as " + userAccount.Name + " on " + b.Universe + "-" + b.language)

	b.debug("extract information from html")
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return err
	}
	b.ogameSession = b.extractor.ExtractOGameSessionFromDoc(doc)
	if b.ogameSession == "" {
		return ErrBadCredentials
	}

	serverTime, _ := b.extractor.ExtractServerTime(pageHTML)
	b.location = serverTime.Location()

	b.cacheFullPageInfo("overview", pageHTML)

	_, _ = b.getPage(PreferencesPage, CelestialID(0)) // Will update preferences cached values

	// Extract chat host and port
	m := regexp.MustCompile(`var nodeUrl\s?=\s?"https:\\/\\/([^:]+):(\d+)\\/socket.io\\/socket.io.js"`).FindSubmatch(pageHTML)
	chatHost := string(m[1])
	chatPort := string(m[2])

	if atomic.CompareAndSwapInt32(&b.chatConnectedAtom, 0, 1) {
		b.closeChatCh = make(chan struct{})
		go func(b *OGame) {
			defer atomic.StoreInt32(&b.chatConnectedAtom, 0)
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

func (b *OGame) cacheFullPageInfo(page string, pageHTML []byte) {
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	b.Planets = b.extractor.ExtractPlanetsFromDoc(doc, b)
	b.isVacationModeEnabled = b.extractor.ExtractIsInVacationFromDoc(doc)
	b.ajaxChatToken, _ = b.extractor.ExtractAjaxChatToken(pageHTML)
	b.characterClass, _ = b.extractor.ExtractCharacterClassFromDoc(doc)
	b.hasCommander = b.extractor.ExtractCommanderFromDoc(doc)
	b.hasAdmiral = b.extractor.ExtractAdmiralFromDoc(doc)
	b.hasEngineer = b.extractor.ExtractEngineerFromDoc(doc)
	b.hasGeologist = b.extractor.ExtractGeologistFromDoc(doc)
	b.hasTechnocrat = b.extractor.ExtractTechnocratFromDoc(doc)

	if page == "overview" {
		b.Player, _ = b.extractor.ExtractUserInfos(pageHTML, b.language)
	} else if page == "preferences" {
		b.CachedPreferences = b.extractor.ExtractPreferencesFromDoc(doc)
	}
}

// DefaultLoginWrapper ...
var DefaultLoginWrapper = func(loginFn func() error) error {
	return loginFn()
}

func (b *OGame) wrapLoginWithExistingCookies() error {
	return b.loginWrapper(b.loginWithExistingCookies)
}

func (b *OGame) wrapLogin() error {
	return b.loginWrapper(b.login)
}

// GetExtractor gets extractor object
func (b *OGame) GetExtractor() Extractor {
	return b.extractor
}

// SetOGameCredentials sets ogame credentials for the bot
func (b *OGame) SetOGameCredentials(username, password string) {
	b.Username = username
	b.password = password
}

func (b *OGame) setOGameLobby(lobby string) {
	if lobby != "lobby-pioneers" {
		lobby = "lobby"
	}
	b.lobby = lobby
}

// SetLoginWrapper ...
func (b *OGame) SetLoginWrapper(newWrapper func(func() error) error) {
	b.loginWrapper = newWrapper
}

// execute a request using the login proxy transport if set
func (b *OGame) doReqWithLoginProxyTransport(req *http.Request) (resp *http.Response, err error) {
	if b.loginProxyTransport != nil {
		oldTransport := b.Client.Transport
		b.Client.Transport = b.loginProxyTransport
		resp, err = b.Client.Do(req)
		b.Client.Transport = oldTransport
	} else {
		resp, err = b.Client.Do(req)
	}
	return
}

// Creates a proxy http transport with optional basic auth
func getProxyTransport(proxy, username, password string) (*http.Transport, error) {
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		return nil, err
	}
	t := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	if username != "" || password != "" {
		basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
		t.ProxyConnectHeader = http.Header{"Proxy-Authorization": {basicAuth}}
	}
	return t, nil
}

func getSocks5Transport(proxyAddress, username, password string) (*http.Transport, error) {
	var auth *proxy.Auth
	if username != "" || password != "" {
		auth = &proxy.Auth{User: username, Password: password}
	}
	dialer, err := proxy.SOCKS5("tcp", proxyAddress, auth, proxy.Direct)
	if err != nil {
		return nil, err
	}
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		},
	}
	return transport, nil
}

func (b *OGame) setProxy(proxyAddress, username, password, proxyType string, loginOnly bool) error {
	if proxyType == "" {
		proxyType = "socks5"
	}
	if proxyAddress == "" {
		b.loginProxyTransport = nil
		b.Client.Transport = http.DefaultTransport
		return nil
	}
	var err error
	transport := http.DefaultTransport
	if proxyType == "socks5" {
		transport, err = getSocks5Transport(proxyAddress, username, password)
	} else if proxyType == "http" {
		transport, err = getProxyTransport(proxyAddress, username, password)
	}
	if loginOnly {
		b.loginProxyTransport = transport
		b.Client.Transport = http.DefaultTransport
	} else {
		b.loginProxyTransport = transport
		b.Client.Transport = transport
	}
	return err
}

// SetProxy this will change the bot http transport object.
// proxyType can be "http" or "socks5".
// An empty proxyAddress will reset the client transport to default value.
func (b *OGame) SetProxy(proxyAddress, username, password, proxyType string, loginOnly bool) error {
	return b.setProxy(proxyAddress, username, password, proxyType, loginOnly)
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
	defer func() {
		if err := resp.Body.Close(); err != nil {
			b.error(err)
		}
	}()
	b.chatRetry.Reset()
	by, _ := ioutil.ReadAll(resp.Body)
	token := strings.Split(string(by), ":")[0]

	origin := "https://" + host + ":" + port + "/"
	wssURL := "wss://" + host + ":" + port + "/socket.io/1/websocket/" + token
	b.ws, err = websocket.Dial(wssURL, "", origin)
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
		if err := b.ws.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
			b.error("failed to set read deadline:", err)
		}
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
				// connection reset by peer
				break
			}
		}
		msg := bytes.Trim(buf, "\x00")
		if bytes.Equal(msg, []byte("1::")) {
			_, _ = b.ws.Write([]byte("1::/chat"))
		} else if bytes.Equal(msg, []byte("1::/chat")) {
			authMsg := `5:` + strconv.FormatInt(b.sessionChatCounter, 10) + `+:/chat:{"name":"authorize","args":["` + b.ogameSession + `"]}`
			_, _ = b.ws.Write([]byte(authMsg))
			b.sessionChatCounter++
		} else if bytes.Equal(msg, []byte("2::")) {
			_, _ = b.ws.Write([]byte("2::"))
		} else if regexp.MustCompile(`\d+::/auctioneer`).Match(msg) {
			for _, clb := range b.auctioneerCallbacks {
				clb(msg)
			}
		} else if regexp.MustCompile(`6::/chat:\d+\+\[true]`).Match(msg) {
			b.debug("chat connected")
		} else if regexp.MustCompile(`6::/chat:\d+\+\[false]`).Match(msg) {
			b.error("Failed to connect to chat")
		} else if bytes.HasPrefix(msg, []byte("5::/chat:")) {
			payload := bytes.TrimPrefix(msg, []byte("5::/chat:"))
			var chatPayload ChatPayload
			if err := json.Unmarshal(payload, &chatPayload); err != nil {
				b.error("Unable to unmarshal chat payload", err, payload)
				continue
			}
			for _, chatMsg := range chatPayload.Args {
				for _, clb := range b.chatCallbacks {
					clb(chatMsg)
				}
			}
		} else {
			b.error("unknown message received:", string(buf))
			time.Sleep(time.Second)
		}
	}
}

// ReconnectChat ...
func (b *OGame) ReconnectChat() bool {
	if b.ws == nil {
		return false
	}
	_, _ = b.ws.Write([]byte("1::/chat"))
	return true
}

// ChatPayload ...
type ChatPayload struct {
	Name string    `json:"name"`
	Args []ChatMsg `json:"args"`
}

// ChatMsg ...
type ChatMsg struct {
	SenderID      int64  `json:"senderId"`
	SenderName    string `json:"senderName"`
	AssociationID int64  `json:"associationId"`
	Text          string `json:"text"`
	ID            int64  `json:"id"`
	Date          int64  `json:"date"`
}

func (m ChatMsg) String() string {
	return "\n" +
		"     Sender ID: " + strconv.FormatInt(m.SenderID, 10) + "\n" +
		"   Sender name: " + m.SenderName + "\n" +
		"Association ID: " + strconv.FormatInt(m.AssociationID, 10) + "\n" +
		"          Text: " + m.Text + "\n" +
		"            ID: " + strconv.FormatInt(m.ID, 10) + "\n" +
		"          Date: " + strconv.FormatInt(m.Date, 10)
}

func (b *OGame) logout() {
	_, _ = b.getPage(LogoutPage, CelestialID(0))
	if atomic.CompareAndSwapInt32(&b.isLoggedInAtom, 1, 0) {
		select {
		case <-b.closeChatCh:
		default:
			close(b.closeChatCh)
			if b.ws != nil {
				_ = b.ws.Close()
			}
		}
	}
}

func isLogged(pageHTML []byte) bool {
	return len(regexp.MustCompile(`<meta name="ogame-session" content="\w+"/>`).FindSubmatch(pageHTML)) == 1 ||
		len(regexp.MustCompile(`var session = "\w+"`).FindSubmatch(pageHTML)) == 1
}

// IsKnowFullPage ...
func IsKnowFullPage(vals url.Values) bool {
	page := vals.Get("page")
	if page == "ingame" {
		page = vals.Get("component")
	}
	return page == OverviewPage ||
		page == ResourcesPage ||
		page == StationPage ||
		page == TraderOverviewPage ||
		page == ResearchPage ||
		page == ShipyardPage ||
		page == DefensePage ||
		page == Fleet1Page ||
		page == GalaxyPage ||
		page == AlliancePage ||
		page == PremiumPage ||
		page == ShopPage ||
		page == RewardsPage ||
		page == ResourceSettingsPage ||
		page == MovementPage ||
		page == HighscorePage ||
		page == BuddiesPage ||
		page == PreferencesPage ||
		page == MessagesPage ||
		page == ChatPage ||

		page == DefensesPage ||
		page == SuppliesPage ||
		page == FacilitiesPage ||
		page == FleetdispatchPage
}

// IsAjaxPage either the requested page is a partial/ajax page
func IsAjaxPage(vals url.Values) bool {
	page := vals.Get("page")
	if page == "ingame" {
		page = vals.Get("component")
	}
	ajax := vals.Get("ajax")
	asJson := vals.Get("asJson")
	return page == FetchEventboxAjaxPage ||
		page == FetchResourcesAjaxPage ||
		page == GalaxyContentAjaxPage ||
		page == EventListAjaxPage ||
		page == AjaxChatAjaxPage ||
		page == NoticesAjaxPage ||
		page == RepairlayerAjaxPage ||
		page == TechtreeAjaxPage ||
		page == PhalanxAjaxPage ||
		page == ShareReportOverlayAjaxPage ||
		page == JumpgatelayerAjaxPage ||
		page == FederationlayerAjaxPage ||
		page == UnionchangeAjaxPage ||
		page == ChangenickAjaxPage ||
		page == PlanetlayerAjaxPage ||
		page == TraderlayerAjaxPage ||
		page == PlanetRenameAjaxPage ||
		page == RightmenuAjaxPage ||
		page == AllianceOverviewAjaxPage ||
		page == SupportAjaxPage ||
		page == BuffActivationAjaxPage ||
		page == AuctioneerAjaxPage ||
		page == HighscoreContentAjaxPage ||
		ajax == "1" ||
		asJson == "1"
}

func canParseEventBox(by []byte) bool {
	err := json.Unmarshal(by, &eventboxResp{})
	return err == nil
}

func canParseSystemInfos(by []byte) bool {
	err := json.Unmarshal(by, &SystemInfos{})
	return err == nil
}

func (b *OGame) preRequestChecks() error {
	if !b.IsEnabled() {
		return ErrBotInactive
	}
	if !b.IsLoggedIn() {
		return ErrBotLoggedOut
	}
	if b.serverURL == "" {
		return errors.New("serverURL is empty")
	}
	return nil
}

func (b *OGame) execRequest(method, finalURL string, payload, vals url.Values) ([]byte, error) {
	var req *http.Request
	var err error
	if method == "GET" {
		req, err = http.NewRequest(method, finalURL, nil)
	} else {
		req, err = http.NewRequest(method, finalURL, strings.NewReader(payload.Encode()))
	}
	if err != nil {
		return []byte{}, err
	}

	if method == "POST" {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	if IsAjaxPage(vals) {
		req.Header.Add("X-Requested-With", "XMLHttpRequest")
	}

	resp, err := b.Client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			b.error(err)
		}
	}()

	if resp.StatusCode >= 500 {
		return []byte{}, err
	}
	by, err := readBody(b, resp)
	if err != nil {
		return []byte{}, err
	}
	b.bytesUploaded += req.ContentLength
	return by, nil
}

func (b *OGame) getPageContent(vals url.Values, opts ...Option) ([]byte, error) {
	var cfg options
	for _, opt := range opts {
		opt(&cfg)
	}

	if err := b.preRequestChecks(); err != nil {
		return []byte{}, err
	}

	finalURL := b.serverURL + "/game/index.php?" + vals.Encode()

	allianceID := vals.Get("allianceId")
	if allianceID != "" {
		finalURL = b.serverURL + "/game/allianceInfo.php?allianceID=" + allianceID
	}

	page := vals.Get("page")
	if page == "ingame" ||
		(page == "componentOnly" && vals.Get("component") == "fetchEventbox") ||
		(page == "componentOnly" && vals.Get("component") == "eventList" && vals.Get("action") != "fetchEventBox") {
		page = vals.Get("component")
	}
	var pageHTMLBytes []byte

	if err := b.withRetry(func() (err error) {
		pageHTMLBytes, err = b.execRequest("GET", finalURL, nil, vals)
		if err != nil {
			return err
		}

		if allianceID != "" {
			return nil
		}
		if (page != LogoutPage && (IsKnowFullPage(vals) || page == "") && !IsAjaxPage(vals) && !isLogged(pageHTMLBytes)) ||
			(page == "eventList" && !bytes.Contains(pageHTMLBytes, []byte("eventListWrap"))) ||
			(page == "fetchEventbox" && !canParseEventBox(pageHTMLBytes)) {
			b.error("Err not logged on page : ", page)
			atomic.StoreInt32(&b.isConnectedAtom, 0)
			return ErrNotLogged
		}

		return nil
	}); err != nil {
		b.error(err)
		return []byte{}, err
	}

	if !IsAjaxPage(vals) && isLogged(pageHTMLBytes) {
		b.cacheFullPageInfo(page, pageHTMLBytes)
	}

	if !cfg.SkipInterceptor {
		go func() {
			for _, fn := range b.interceptorCallbacks {
				fn("GET", finalURL, vals, nil, pageHTMLBytes)
			}
		}()
	}

	return pageHTMLBytes, nil
}

func (b *OGame) postPageContent(vals, payload url.Values, opts ...Option) ([]byte, error) {
	var cfg options
	for _, opt := range opts {
		opt(&cfg)
	}

	if err := b.preRequestChecks(); err != nil {
		return []byte{}, err
	}

	if vals.Get("page") == "ajaxChat" && payload.Get("mode") == "1" {
		payload.Set("token", b.ajaxChatToken)
	}

	finalURL := b.serverURL + "/game/index.php?" + vals.Encode()
	page := vals.Get("page")
	if page == "ingame" {
		page = vals.Get("component")
	}
	var pageHTMLBytes []byte

	if err := b.withRetry(func() (err error) {
		// Needs to be inside the withRetry, so if we need to re-login the redirect is back for the login call
		// Prevent redirect (301) https://stackoverflow.com/a/38150816/4196220
		b.Client.CheckRedirect = func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }
		defer func() { b.Client.CheckRedirect = nil }()

		pageHTMLBytes, err = b.execRequest("POST", finalURL, payload, vals)
		if err != nil {
			return err
		}

		if page == "galaxyContent" && !canParseSystemInfos(pageHTMLBytes) {
			b.error("Err not logged on page : ", page)
			atomic.StoreInt32(&b.isConnectedAtom, 0)
			return ErrNotLogged
		}

		return nil
	}); err != nil {
		b.error(err)
		return []byte{}, err
	}

	if page == "preferences" {
		b.CachedPreferences = b.extractor.ExtractPreferences(pageHTMLBytes)
	} else if page == "ajaxChat" && (payload.Get("mode") == "1" || payload.Get("mode") == "3") {
		var res ChatPostResp
		if err := json.Unmarshal(pageHTMLBytes, &res); err != nil {
			return []byte{}, err
		}
		b.ajaxChatToken = res.NewToken
	}

	if !cfg.SkipInterceptor {
		go func() {
			for _, fn := range b.interceptorCallbacks {
				fn("POST", finalURL, vals, payload, pageHTMLBytes)
			}
		}()
	}

	return pageHTMLBytes, nil
}

func (b *OGame) getAlliancePageContent(vals url.Values) ([]byte, error) {
	if err := b.preRequestChecks(); err != nil {
		return []byte{}, err
	}
	finalURL := b.serverURL + "/game/allianceInfo.php?" + vals.Encode()
	return b.execRequest("GET", finalURL, nil, vals)
}

type eventboxResp struct {
	Hostile  int
	Neutral  int
	Friendly int
}

func (b *OGame) withRetry(fn func() error) error {
	maxRetry := 10
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
		err := fn()
		if err == nil {
			break
		}
		// If we manually logged out, do not try to auto re login.
		if !b.IsEnabled() {
			return ErrBotInactive
		}
		if !b.IsLoggedIn() {
			return ErrBotLoggedOut
		}
		maxRetry--
		if maxRetry <= 0 {
			return errors.Wrap(err, ErrFailedExecuteCallback.Error())
		}

		retry(err)

		if err == ErrNotLogged {
			if loginErr := b.wrapLogin(); loginErr != nil {
				b.error(loginErr.Error()) // log error
				if loginErr == ErrAccountNotFound ||
					loginErr == ErrAccountBlocked {
					return loginErr
				}
			}
		}
	}
	return nil
}

func (b *OGame) getPageJSON(vals url.Values, v interface{}) error {
	pageJSON, err := b.getPageContent(vals)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(pageJSON, v); err != nil {
		return ErrNotLogged
	}
	return nil
}

func (b *OGame) enable() {
	atomic.StoreInt32(&b.isEnabledAtom, 1)
	b.stateChanged(false, "Enable")
}

func (b *OGame) disable() {
	atomic.StoreInt32(&b.isEnabledAtom, 0)
	b.stateChanged(false, "Disable")
}

func (b *OGame) isEnabled() bool {
	return atomic.LoadInt32(&b.isEnabledAtom) == 1
}

func (b *OGame) isCollector() bool {
	return b.characterClass == Collector
}

func (b *OGame) isGeneral() bool {
	return b.characterClass == General
}

func (b *OGame) isDiscoverer() bool {
	return b.characterClass == Discoverer
}

func (b *OGame) getUniverseSpeed() int64 {
	return b.serverData.Speed
}

func (b *OGame) getUniverseSpeedFleet() int64 {
	return b.serverData.SpeedFleet
}

func (b *OGame) isDonutGalaxy() bool {
	return b.serverData.DonutGalaxy
}

func (b *OGame) isDonutSystem() bool {
	return b.serverData.DonutSystem
}

func (b *OGame) fetchEventbox() (res eventboxResp, err error) {
	err = b.getPageJSON(url.Values{"page": {"fetchEventbox"}}, &res)
	return
}

func (b *OGame) isUnderAttack() (bool, error) {
	res, err := b.fetchEventbox()
	return res.Hostile > 0, err
}

type resourcesResp struct {
	Metal struct {
		Resources struct {
			ActualFormat string
			Actual       int64
			Max          int64
			Production   float64
		}
		Tooltip string
		Class   string
	}
	Crystal struct {
		Resources struct {
			ActualFormat string
			Actual       int64
			Max          int64
			Production   float64
		}
		Tooltip string
		Class   string
	}
	Deuterium struct {
		Resources struct {
			ActualFormat string
			Actual       int64
			Max          int64
			Production   float64
		}
		Tooltip string
		Class   string
	}
	Energy struct {
		Resources struct {
			ActualFormat string
			Actual       int64
		}
		Tooltip string
		Class   string
	}
	Darkmatter struct {
		Resources struct {
			ActualFormat string
			Actual       int64
		}
		String  string
		Tooltip string
	}
	HonorScore int64
}

func (b *OGame) getPlanets() []Planet {
	pageHTML, _ := b.getPage(OverviewPage, CelestialID(0))
	return b.extractor.ExtractPlanets(pageHTML, b)
}

func (b *OGame) getPlanet(v interface{}) (Planet, error) {
	pageHTML, _ := b.getPage(OverviewPage, CelestialID(0))
	return b.extractor.ExtractPlanet(pageHTML, v, b)
}

func (b *OGame) getMoons() []Moon {
	pageHTML, _ := b.getPage(OverviewPage, CelestialID(0))
	return b.extractor.ExtractMoons(pageHTML, b)
}

func (b *OGame) getMoon(v interface{}) (Moon, error) {
	pageHTML, _ := b.getPage(OverviewPage, CelestialID(0))
	return b.extractor.ExtractMoon(pageHTML, b, v)
}

func (b *OGame) getCelestials() ([]Celestial, error) {
	pageHTML, _ := b.getPage(OverviewPage, CelestialID(0))
	return b.extractor.ExtractCelestials(pageHTML, b)
}

func (b *OGame) getCelestial(v interface{}) (Celestial, error) {
	pageHTML, _ := b.getPage(OverviewPage, CelestialID(0))
	return b.extractor.ExtractCelestial(pageHTML, b, v)
}

func (b *OGame) abandon(v interface{}) error {
	pageHTML, _ := b.getPage(OverviewPage, CelestialID(0))
	var planetID PlanetID
	if coordStr, ok := v.(string); ok {
		coord, err := ParseCoord(coordStr)
		if err != nil {
			return err
		}
		planet, err := b.extractor.ExtractPlanetByCoord(pageHTML, b, coord)
		if err != nil {
			return err
		}
		planetID = planet.ID
	} else if coord, ok := v.(Coordinate); ok {
		planet, err := b.extractor.ExtractPlanetByCoord(pageHTML, b, coord)
		if err != nil {
			return err
		}
		planetID = planet.ID
	} else if planet, ok := v.(Planet); ok {
		planetID = planet.ID
	} else if id, ok := v.(PlanetID); ok {
		planetID = id
	} else if id, ok := v.(int); ok {
		planetID = PlanetID(id)
	} else if id, ok := v.(int32); ok {
		planetID = PlanetID(id)
	} else if id, ok := v.(int64); ok {
		planetID = PlanetID(id)
	} else if id, ok := v.(float32); ok {
		planetID = PlanetID(id)
	} else if id, ok := v.(float64); ok {
		planetID = PlanetID(id)
	} else if id, ok := v.(lua.LNumber); ok {
		planetID = PlanetID(id)
	} else {
		return errors.New("invalid parameter")
	}
	planets := b.extractor.ExtractPlanets(pageHTML, b)
	found := false
	for _, planet := range planets {
		if planet.ID == planetID {
			found = true
			break
		}
	}
	if !found {
		return errors.New("invalid planet id")
	}
	pageHTML, _ = b.getPage(PlanetlayerPage, planetID.Celestial())
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	abandonToken := doc.Find("form#planetMaintenanceDelete input[name=abandon]").AttrOr("value", "")
	token := doc.Find("form#planetMaintenanceDelete input[name=token]").AttrOr("value", "")
	payload := url.Values{
		"abandon":  {abandonToken},
		"token":    {token},
		"password": {b.password},
	}
	_, err := b.postPageContent(url.Values{"page": {"planetGiveup"}}, payload)
	return err
}

func (b *OGame) serverTime() time.Time {
	pageHTML, _ := b.getPage(OverviewPage, CelestialID(0))
	serverTime, err := b.extractor.ExtractServerTime(pageHTML)
	if err != nil {
		b.error(err.Error())
	}
	return serverTime
}

func (b *OGame) getUserInfos() UserInfos {
	pageHTML, _ := b.getPage(OverviewPage, CelestialID(0))
	userInfos, err := b.extractor.ExtractUserInfos(pageHTML, b.language)
	if err != nil {
		b.error(err)
	}
	return userInfos
}

// ChatPostResp ...
type ChatPostResp struct {
	Status   string `json:"status"`
	ID       int    `json:"id"`
	SenderID int    `json:"senderId"`
	TargetID int    `json:"targetId"`
	Text     string `json:"text"`
	Date     int    `json:"date"`
	NewToken string `json:"newToken"`
}

func (b *OGame) sendMessage(id int64, message string, isPlayer bool) error {
	payload := url.Values{
		"text":  {message + "\n"},
		"ajax":  {"1"},
		"token": {b.ajaxChatToken},
	}
	if isPlayer {
		payload.Set("playerId", strconv.FormatInt(id, 10))
		payload.Set("mode", "1")
	} else {
		payload.Set("associationId", strconv.FormatInt(id, 10))
		payload.Set("mode", "3")
	}
	bobyBytes, err := b.postPageContent(url.Values{"page": {"ajaxChat"}}, payload)
	if err != nil {
		return err
	}
	if strings.Contains(string(bobyBytes), "INVALID_PARAMETERS") {
		return errors.New("invalid parameters")
	}
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(string(bobyBytes)))
	if doc.Find("title").Text() == "OGame Lobby" {
		return ErrNotLogged
	}
	var res ChatPostResp
	if err := json.Unmarshal(bobyBytes, &res); err != nil {
		return err
	}
	b.ajaxChatToken = res.NewToken
	return nil
}

func (b *OGame) getFleetsFromEventList() []Fleet {
	pageHTML, _ := b.getPageContent(url.Values{"eventList": {"movement"}, "ajax": {"1"}})
	return b.extractor.ExtractFleetsFromEventList(pageHTML)
}

func (b *OGame) getFleets(opts ...Option) ([]Fleet, Slots) {
	pageHTML, _ := b.getPage(MovementPage, CelestialID(0), opts...)
	fleets := b.extractor.ExtractFleets(pageHTML)
	slots := b.extractor.ExtractSlots(pageHTML)
	return fleets, slots
}

func (b *OGame) cancelFleet(fleetID FleetID) error {
	_, _ = b.getPageContent(url.Values{"page": {"movement"}, "return": {fleetID.String()}})
	return nil
}

// Slots ...
type Slots struct {
	InUse    int64
	Total    int64
	ExpInUse int64
	ExpTotal int64
}

func (b *OGame) getSlots() Slots {
	pageHTML, _ := b.getPage(Fleet1Page, CelestialID(0))
	return b.extractor.ExtractSlots(pageHTML)
}

// Returns the distance between two galaxy
func galaxyDistance(galaxy1, galaxy2, universeSize int64, donutGalaxy bool) (distance int64) {
	if !donutGalaxy {
		return int64(20000 * math.Abs(float64(galaxy2-galaxy1)))
	}
	if galaxy1 > galaxy2 {
		galaxy1, galaxy2 = galaxy2, galaxy1
	}
	val := math.Min(float64(galaxy2-galaxy1), float64((galaxy1+universeSize)-galaxy2))
	return int64(20000 * val)
}

func systemDistance(nbSystems, system1, system2 int64, donutSystem bool) (distance int64) {
	if !donutSystem {
		return int64(math.Abs(float64(system2 - system1)))
	}
	if system1 > system2 {
		system1, system2 = system2, system1
	}
	return int64(math.Min(float64(system2-system1), float64((system1+nbSystems)-system2)))
}

// Returns the distance between two systems
func flightSystemDistance(nbSystems, system1, system2 int64, donutSystem bool) (distance int64) {
	return 2700 + 95*systemDistance(nbSystems, system1, system2, donutSystem)
}

// Returns the distance between two planets
func planetDistance(planet1, planet2 int64) (distance int64) {
	return int64(1000 + 5*math.Abs(float64(planet2-planet1)))
}

// Distance returns the distance between two coordinates
func Distance(c1, c2 Coordinate, universeSize, nbSystems int64, donutGalaxy, donutSystem bool) (distance int64) {
	if c1.Galaxy != c2.Galaxy {
		return galaxyDistance(c1.Galaxy, c2.Galaxy, universeSize, donutGalaxy)
	}
	if c1.System != c2.System {
		return flightSystemDistance(nbSystems, c1.System, c2.System, donutSystem)
	}
	if c1.Position != c2.Position {
		return planetDistance(c1.Position, c2.Position)
	}
	return 5
}

func findSlowestSpeed(ships ShipsInfos, techs Researches, isCollector, isGeneral bool) int64 {
	var minSpeed int64 = math.MaxInt64
	for _, ship := range Ships {
		if ship.GetID() == SolarSatelliteID || ship.GetID() == CrawlerID {
			continue
		}
		shipSpeed := ship.GetSpeed(techs, isCollector, isGeneral)
		if ships.ByID(ship.GetID()) > 0 && shipSpeed < minSpeed {
			minSpeed = shipSpeed
		}
	}
	return minSpeed
}

func calcFuel(ships ShipsInfos, dist, duration int64, universeSpeedFleet, fleetDeutSaveFactor float64, techs Researches, isCollector, isGeneral bool) (fuel int64) {
	tmpFn := func(baseFuel, nbr, shipSpeed int64) float64 {
		tmpSpeed := (35000 / (float64(duration)*universeSpeedFleet - 10)) * math.Sqrt(float64(dist)*10/float64(shipSpeed))
		return float64(baseFuel*nbr*dist) / 35000 * math.Pow(tmpSpeed/10+1, 2)
	}
	tmpFuel := 0.0
	for _, ship := range Ships {
		if ship.GetID() == SolarSatelliteID || ship.GetID() == CrawlerID {
			continue
		}
		nbr := ships.ByID(ship.GetID())
		if nbr > 0 {
			tmpFuel += tmpFn(ship.GetFuelConsumption(techs), nbr, ship.GetSpeed(techs, isCollector, isGeneral))
		}
	}
	fuel = int64(1 + math.Floor(tmpFuel*fleetDeutSaveFactor))
	return
}

func calcFlightTime(origin, destination Coordinate, universeSize, nbSystems int64, donutGalaxy, donutSystem bool,
	fleetDeutSaveFactor, speed float64, universeSpeedFleet int64, ships ShipsInfos, techs Researches, characterClass CharacterClass) (secs, fuel int64) {
	if !ships.HasShips() {
		return
	}
	isCollector := characterClass == Collector
	isGeneral := characterClass == General
	s := speed
	v := float64(findSlowestSpeed(ships, techs, isCollector, isGeneral))
	a := float64(universeSpeedFleet)
	d := float64(Distance(origin, destination, universeSize, nbSystems, donutGalaxy, donutSystem))
	secs = int64(math.Round(((3500/s)*math.Sqrt(d*10/v) + 10) / a))
	fuel = calcFuel(ships, int64(d), secs, float64(universeSpeedFleet), fleetDeutSaveFactor, techs, isCollector, isGeneral)
	return
}

// getPhalanx makes 3 calls to ogame server (2 validation, 1 scan)
func (b *OGame) getPhalanx(moonID MoonID, coord Coordinate) ([]Fleet, error) {
	res := make([]Fleet, 0)

	// Get moon facilities html page (first call to ogame server)
	moonFacilitiesHTML, _ := b.getPage(StationPage, moonID.Celestial())

	// Extract bunch of infos from the html
	moon, err := b.extractor.ExtractMoon(moonFacilitiesHTML, b, moonID)
	if err != nil {
		return res, errors.New("moon not found")
	}
	resources := b.extractor.ExtractResources(moonFacilitiesHTML)
	moonFacilities, _ := b.extractor.ExtractFacilities(moonFacilitiesHTML)
	phalanxLvl := moonFacilities.SensorPhalanx

	// Ensure we have the resources to scan the planet
	if resources.Deuterium < SensorPhalanx.ScanConsumption() {
		return res, errors.New("not enough deuterium")
	}

	// Verify that coordinate is in phalanx range
	phalanxRange := SensorPhalanx.GetRange(phalanxLvl)
	if moon.Coordinate.Galaxy != coord.Galaxy ||
		systemDistance(b.serverData.Systems, moon.Coordinate.System, coord.System, b.serverData.DonutSystem) > phalanxRange {
		return res, errors.New("coordinate not in phalanx range")
	}

	// Get galaxy planets information, verify coordinate is valid planet (second call to ogame server)
	planetInfos, _ := b.galaxyInfos(coord.Galaxy, coord.System)
	target := planetInfos.Position(coord.Position)
	if target == nil {
		return res, errors.New("invalid planet coordinate")
	}
	// Ensure you are not scanning your own planet
	if target.Player.ID == b.Player.PlayerID {
		return res, errors.New("cannot scan own planet")
	}

	// Run the phalanx scan (third call to ogame server)
	return b.getUnsafePhalanx(moonID, coord)
}

// getUnsafePhalanx ...
func (b *OGame) getUnsafePhalanx(moonID MoonID, coord Coordinate) ([]Fleet, error) {
	pageHTML, _ := b.getPageContent(url.Values{
		"page":     {"phalanx"},
		"galaxy":   {strconv.FormatInt(coord.Galaxy, 10)},
		"system":   {strconv.FormatInt(coord.System, 10)},
		"position": {strconv.FormatInt(coord.Position, 10)},
		"ajax":     {"1"},
		"cp":       {strconv.FormatInt(int64(moonID), 10)},
	})
	return b.extractor.ExtractPhalanx(pageHTML)
}

func moonIDInSlice(needle MoonID, haystack []MoonID) bool {
	for _, element := range haystack {
		if needle == element {
			return true
		}
	}
	return false
}

func (b *OGame) headersForPage(url string) (http.Header, error) {
	if !b.IsEnabled() {
		return nil, ErrBotInactive
	}
	if !b.IsLoggedIn() {
		return nil, ErrBotLoggedOut
	}

	if b.serverURL == "" {
		err := errors.New("serverURL is empty")
		b.error(err)
		return nil, err
	}

	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}

	finalURL := b.serverURL + url

	req, err := http.NewRequest("HEAD", finalURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := b.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 500 {
		return nil, err
	}

	return resp.Header, err
}

func (b *OGame) jumpGateDestinations(originMoonID MoonID) ([]MoonID, int64, error) {
	pageHTML, _ := b.getPage(JumpgatelayerPage, originMoonID.Celestial())
	_, _, dests, wait := b.extractor.ExtractJumpGate(pageHTML)
	if wait > 0 {
		return dests, wait, fmt.Errorf("jump gate is in recharge mode for %d seconds", wait)
	}
	return dests, wait, nil
}

func (b *OGame) executeJumpGate(originMoonID, destMoonID MoonID, ships ShipsInfos) (bool, int64, error) {
	pageHTML, _ := b.getPage(JumpgatelayerPage, originMoonID.Celestial())
	availShips, token, dests, wait := b.extractor.ExtractJumpGate(pageHTML)
	if wait > 0 {
		return false, wait, fmt.Errorf("jump gate is in recharge mode for %d seconds", wait)
	}

	// Validate destination moon id
	if !moonIDInSlice(destMoonID, dests) {
		return false, 0, errors.New("destination moon id invalid")
	}

	payload := url.Values{"token": {token}, "zm": {strconv.FormatInt(int64(destMoonID), 10)}}

	// Add ships to payload
	for _, s := range Ships {
		// Get the min between what is available and what we want
		nbr := int64(math.Min(float64(ships.ByID(s.GetID())), float64(availShips.ByID(s.GetID()))))
		if nbr > 0 {
			payload.Add("ship_"+strconv.FormatInt(int64(s.GetID()), 10), strconv.FormatInt(nbr, 10))
		}
	}

	if _, err := b.postPageContent(url.Values{"page": {"jumpgate_execute"}}, payload); err != nil {
		return false, 0, err
	}
	return true, 0, nil
}

func (b *OGame) getEmpire(nbr int64) (interface{}, error) {
	// Valid URLs:
	// /game/index.php?page=standalone&component=empire&planetType=0
	// /game/index.php?page=standalone&component=empire&planetType=1
	vals := url.Values{"page": {"standalone"}, "component": {"empire"}, "planetType": {strconv.FormatInt(nbr, 10)}}
	pageHTMLBytes, err := b.getPageContent(vals)
	if err != nil {
		return nil, err
	}
	// Replace the Ogame hostname with our custom hostname
	pageHTML := strings.Replace(string(pageHTMLBytes), b.serverURL, b.apiNewHostname, -1)
	return b.extractor.ExtractEmpire([]byte(pageHTML), nbr)
}

func (b *OGame) createUnion(fleet Fleet, allUnionUsers []UserInfos) (int64, error) {
	if fleet.ID == 0 {
		return 0, errors.New("invalid fleet id")
	}
	pageHTML, _ := b.getPageContent(url.Values{"page": {"federationlayer"}, "union": {"0"}, "fleet": {strconv.FormatInt(int64(fleet.ID), 10)}, "target": {strconv.FormatInt(fleet.TargetPlanetID, 10)}, "ajax": {"1"}})
	payload := b.extractor.ExtractFederation(pageHTML)

	payload.Del("unionUsers")

	var unionUsers string
	for _, uu := range allUnionUsers {
		if unionUsers == "" {
			unionUsers += uu.PlayerName
		} else {
			unionUsers += ";" + uu.PlayerName
		}
	}

	payload.Add("unionUsers", unionUsers)

	by, err := b.postPageContent(url.Values{"page": {"unionchange"}, "ajax": {"1"}}, payload)
	if err != nil {
		return 0, err
	}
	var res struct {
		FleetID  int64
		UnionID  int64
		TargetID int64
		Errorbox struct {
			Type   string
			Text   string
			Failed int64
		}
	}
	if err := json.Unmarshal(by, &res); err != nil {
		return 0, err
	}
	if res.Errorbox.Failed != 0 {
		return 0, errors.New(res.Errorbox.Text)
	}
	return res.UnionID, nil
}

func (b *OGame) highscore(category, typ, page int64) (out Highscore, err error) {
	if category < 1 || category > 2 {
		return out, errors.New("category must be in [1, 2] (1:player, 2:alliance)")
	}
	if typ < 0 || typ > 7 {
		return out, errors.New("category must be in [0, 7] (0:Total, 1:Economy, 2:Research, 3:Military, 4:Military Built, 5:Military Destroyed, 6:Military Lost, 7:Honor)")
	}
	if page < 1 {
		return out, errors.New("page must be greater than or equal to 1")
	}
	vals := url.Values{
		"page":     {HighscoreContentAjaxPage},
		"category": {strconv.FormatInt(category, 10)},
		"type":     {strconv.FormatInt(typ, 10)},
		"site":     {strconv.FormatInt(page, 10)},
	}
	payload := url.Values{}
	pageHTML, _ := b.postPageContent(vals, payload)
	return b.extractor.ExtractHighscore(pageHTML)
}

func (b *OGame) getAllResources() (map[CelestialID]Resources, error) {
	vals := url.Values{
		"page": {"traderOverview"},
	}
	payload := url.Values{
		"show": {"auctioneer"},
		"ajax": {"1"},
	}
	pageHTML, _ := b.postPageContent(vals, payload)
	return b.extractor.ExtractAllResources(pageHTML)
}

func (b *OGame) getDMCosts(celestialID CelestialID) (DMCosts, error) {
	pageHTML, _ := b.getPage(OverviewPage, celestialID)
	return b.extractor.ExtractDMCosts(pageHTML)
}

func (b *OGame) useDM(typ string, celestialID CelestialID) error {
	if typ != "buildings" && typ != "research" && typ != "shipyard" {
		return fmt.Errorf("invalid type %s", typ)
	}
	pageHTML, _ := b.getPage(OverviewPage, celestialID)
	costs, err := b.extractor.ExtractDMCosts(pageHTML)
	if err != nil {
		return err
	}
	var buyAndActivate, token string
	switch typ {
	case "buildings":
		buyAndActivate, token = costs.Buildings.BuyAndActivateToken, costs.Buildings.Token
	case "research":
		buyAndActivate, token = costs.Research.BuyAndActivateToken, costs.Research.Token
	case "shipyard":
		buyAndActivate, token = costs.Shipyard.BuyAndActivateToken, costs.Shipyard.Token
	}
	params := url.Values{
		"page":           {"inventory"},
		"buyAndActivate": {buyAndActivate},
	}
	payload := url.Values{
		"ajax":         {"1"},
		"token":        {token},
		"referrerPage": {"ingame"},
	}
	if _, err := b.postPageContent(params, payload); err != nil {
		return err
	}
	return nil
}

func (b *OGame) getItems(celestialID CelestialID) (items []Item, err error) {
	params := url.Values{"page": {"buffActivation"}, "ajax": {"1"}, "type": {"1"}}
	if celestialID != 0 {
		params.Set("cp", strconv.FormatInt(int64(celestialID), 10))
	}
	pageHTML, _ := b.getPageContent(params)
	_, items, err = b.extractor.ExtractBuffActivation(pageHTML)
	return
}

func (b *OGame) activateItem(ref string, celestialID CelestialID) error {
	params := url.Values{"page": {"buffActivation"}, "ajax": {"1"}, "type": {"1"}}
	if celestialID != 0 {
		params.Set("cp", strconv.FormatInt(int64(celestialID), 10))
	}
	pageHTML, _ := b.getPageContent(params)
	token, _, err := b.extractor.ExtractBuffActivation(pageHTML)
	if err != nil {
		return err
	}
	params = url.Values{"page": {"inventory"}}
	payload := url.Values{
		"ajax":         {"1"},
		"token":        {token},
		"referrerPage": {"ingame"},
		"item":         {ref},
	}
	var res struct {
		Message string `json:"message"`
		Error   bool   `json:"error"`
	}
	by, err := b.postPageContent(params, payload)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(by, &res); err != nil {
		return err
	}
	if res.Error {
		return errors.New(res.Message)
	}
	return err
}

func (b *OGame) getAuction(celestialID CelestialID) (Auction, error) {
	payload := url.Values{"show": {"auctioneer"}, "ajax": {"1"}}
	if celestialID != 0 {
		payload.Set("cp", strconv.FormatInt(int64(celestialID), 10))
	}
	auctionHTML, err := b.postPageContent(url.Values{"page": {"traderOverview"}}, payload)
	if err != nil {
		return Auction{}, err
	}
	return b.extractor.ExtractAuction(auctionHTML)
}

func (b *OGame) doAuction(celestialID CelestialID, bid map[CelestialID]Resources) error {
	// Get fresh token (among others)
	auction, err := b.getAuction(celestialID)
	if err != nil {
		return err
	}

	if auction.HasFinished {
		return errors.New("auction completed")
	}

	payload := url.Values{}
	for auctionCelestialIDString, _ := range auction.Resources {
		payload.Set("bid[planets]["+auctionCelestialIDString+"][metal]", "0")
		payload.Set("bid[planets]["+auctionCelestialIDString+"][crystal]", "0")
		payload.Set("bid[planets]["+auctionCelestialIDString+"][deuterium]", "0")
	}
	for celestialID, resources := range bid {
		payload.Set("bid[planets]["+strconv.FormatInt(int64(celestialID), 10)+"][metal]", strconv.FormatInt(resources.Metal, 10))
		payload.Set("bid[planets]["+strconv.FormatInt(int64(celestialID), 10)+"][crystal]", strconv.FormatInt(resources.Crystal, 10))
		payload.Set("bid[planets]["+strconv.FormatInt(int64(celestialID), 10)+"][deuterium]", strconv.FormatInt(resources.Deuterium, 10))
	}

	payload.Add("bid[honor]", "0")
	payload.Add("token", auction.Token)
	payload.Add("ajax", "1")

	if celestialID != 0 {
		payload.Set("cp", strconv.FormatInt(int64(celestialID), 10))
	}

	auctionHTML, err := b.postPageContent(url.Values{"page": {"auctioneer"}}, payload)
	if err != nil {
		return err
	}

	/*
		Example return from postPageContent on page:auctioneer :
		{
		  "error": false,
		  "message": "Your bid has been accepted.",
		  "planetResources": {
		    "$planetID": {
		      "metal": $metal,
		      "crystal": $crystal,
		      "deuterium": $deuterium
		    },
		    "31434289": {
		      "metal": 5202955.0986408,
		      "crystal": 2043854.5003197,
		      "deuterium": 1552571.3257004
		    }
		    <...>
		  },
		  "honor": 10107,
		  "newToken": "940387sf93e28fbf47b24920c510db38"
		}
	*/

	var jsonObj map[string]interface{}
	if err := json.Unmarshal(auctionHTML, &jsonObj); err != nil {
		return err
	}
	if jsonObj["error"] == true {
		return errors.New(jsonObj["message"].(string))
	}
	return nil
}

type planetResource struct {
	Input struct {
		Metal     int64
		Crystal   int64
		Deuterium int64
	}
	Output struct {
		Metal     int64
		Crystal   int64
		Deuterium int64
	}
	IsMoon        bool
	ImageFileName string
	Name          string
	// OtherPlanet   string // can be null or apparently number (cannot unmarshal number into Go struct field planetResource.OtherPlanet of type string)
}

// PlanetResources ...
type PlanetResources map[CelestialID]planetResource

// Multiplier ...
type Multiplier struct {
	Metal     float64
	Crystal   float64
	Deuterium float64
	Honor     float64
}

func calcResources(price int64, planetResources PlanetResources, multiplier Multiplier) url.Values {
	sortedCelestialIDs := make([]CelestialID, 0)
	for celestialID := range planetResources {
		sortedCelestialIDs = append(sortedCelestialIDs, celestialID)
	}
	sort.Slice(sortedCelestialIDs, func(i, j int) bool {
		return int64(sortedCelestialIDs[i]) < int64(sortedCelestialIDs[j])
	})

	payload := url.Values{}
	remaining := price
	for celestialID, res := range planetResources {
		metalNeeded := res.Input.Metal
		if remaining < int64(float64(metalNeeded)*multiplier.Metal) {
			metalNeeded = int64(math.Ceil(float64(remaining) / multiplier.Metal))
		}
		remaining -= int64(float64(metalNeeded) * multiplier.Metal)

		crystalNeeded := res.Input.Crystal
		if remaining < int64(float64(crystalNeeded)*multiplier.Crystal) {
			crystalNeeded = int64(math.Ceil(float64(remaining) / multiplier.Crystal))
		}
		remaining -= int64(float64(crystalNeeded) * multiplier.Crystal)

		deuteriumNeeded := res.Input.Deuterium
		if remaining < int64(float64(deuteriumNeeded)*multiplier.Deuterium) {
			deuteriumNeeded = int64(math.Ceil(float64(remaining) / multiplier.Deuterium))
		}
		remaining -= int64(float64(deuteriumNeeded) * multiplier.Deuterium)

		payload.Add("bid[planets]["+strconv.FormatInt(int64(celestialID), 10)+"][metal]", strconv.FormatInt(metalNeeded, 10))
		payload.Add("bid[planets]["+strconv.FormatInt(int64(celestialID), 10)+"][crystal]", strconv.FormatInt(crystalNeeded, 10))
		payload.Add("bid[planets]["+strconv.FormatInt(int64(celestialID), 10)+"][deuterium]", strconv.FormatInt(deuteriumNeeded, 10))
	}
	return payload
}

func (b *OGame) buyOfferOfTheDay() error {
	pageHTML, err := b.postPageContent(url.Values{"page": {"traderOverview"}}, url.Values{"show": {"importexport"}, "ajax": {"1"}})
	if err != nil {
		return err
	}

	price, importToken, planetResources, multiplier, err := b.extractor.ExtractOfferOfTheDay(pageHTML)
	if err != nil {
		return err
	}
	payload := calcResources(price, planetResources, multiplier)
	payload.Add("action", "trade")
	payload.Add("bid[honor]", "0")
	payload.Add("token", importToken)
	payload.Add("ajax", "1")
	pageHTML1, err := b.postPageContent(url.Values{"page": {"import"}}, payload)
	if err != nil {
		return err
	}
	// {"message":"You have bought a container.","error":false,"item":{"uuid":"40f6c78e11be01ad3389b7dccd6ab8efa9347f3c","itemText":"You have purchased 1 KRAKEN Bronze.","bargainText":"The contents of the container not appeal to you? For 500 Dark Matter you can exchange the container for another random container of the same quality. You can only carry out this exchange 2 times per daily offer.","bargainCost":500,"bargainCostText":"Costs: 500 Dark Matter","tooltip":"KRAKEN Bronze|Reduces the building time of buildings currently under construction by <b>30m<\/b>.<br \/><br \/>\nDuration: now<br \/><br \/>\nPrice: --- <br \/>\nIn Inventory: 1","image":"98629d11293c9f2703592ed0314d99f320f45845","amount":1,"rarity":"common"},"newToken":"07eefc14105db0f30cb331a8b7af0bfe"}
	var tmp struct {
		Message  string
		Error    bool
		NewToken string
	}
	if err := json.Unmarshal(pageHTML1, &tmp); err != nil {
		return err
	}
	if tmp.Error {
		return errors.New(tmp.Message)
	}

	payload2 := url.Values{"action": {"takeItem"}, "token": {tmp.NewToken}, "ajax": {"1"}}
	pageHTML2, err := b.postPageContent(url.Values{"page": {"import"}}, payload2)
	var tmp2 struct {
		Message  string
		Error    bool
		NewToken string
	}
	if err := json.Unmarshal(pageHTML2, &tmp2); err != nil {
		return err
	}
	if tmp2.Error {
		return errors.New(tmp2.Message)
	}
	// {"message":"You have accepted the offer and put the item in your inventory.","error":false,"item":{"name":"KRAKEN Bronze","image":"bc4e2315f7db4286ba72a424a32c920e78af8e27","imageLarge":"98629d11293c9f2703592ed0314d99f320f45845","title":"KRAKEN Bronze|Reduces the building time of buildings currently under construction by <b>30m<\/b>.<br \/><br \/>\nDuration: now<br \/><br \/>\nPrice: --- <br \/>\nIn Inventory: 2","effect":"Reduces the building time of buildings currently under construction by <b>30m<\/b>.","ref":"40f6c78e11be01ad3389b7dccd6ab8efa9347f3c","rarity":"common","amount":2,"amount_free":2,"amount_bought":0,"category":["d8d49c315fa620d9c7f1f19963970dea59a0e3be","dc9ec90e5a2163cc063b8bb3e9fe392782f565c8"],"currency":"dm","costs":"3000","isReduced":false,"buyable":false,"canBeActivated":false,"canBeBoughtAndActivated":false,"isAnUpgrade":false,"hasEnoughCurrency":true,"cooldown":0,"duration":0,"durationExtension":null,"totalTime":null,"timeLeft":null,"status":null,"extendable":false,"firstStatus":"effecting","toolTip":"KRAKEN Bronze|Reduces the building time of buildings currently under construction by &lt;b&gt;30m&lt;\/b&gt;.&lt;br \/&gt;&lt;br \/&gt;\nDuration: now&lt;br \/&gt;&lt;br \/&gt;\nPrice: --- &lt;br \/&gt;\nIn Inventory: 2","buyTitle":"This item is currently unavailable for purchase.","activationTitle":"There is no facility currently being built whose construction time can be shortened.","moonOnlyItem":false,"newOffer":false,"noOfferMessage":"There are no further offers today. Please come again tomorrow."},"newToken":"68198ffde0837211de8421b1c6447448"}

	return nil
}

// Hack fix: When moon name is >12, the moon image disappear from the EventsBox
// and attacks are detected on planet instead.
func fixAttackEvents(attacks []AttackEvent, planets []Planet) {
	for i, attack := range attacks {
		if len(attack.DestinationName) > 12 {
			for _, planet := range planets {
				if attack.Destination.Equal(planet.Coordinate) &&
					planet.Moon != nil &&
					attack.DestinationName != planet.Name &&
					attack.DestinationName == planet.Moon.Name {
					attacks[i].Destination.Type = MoonType
				}
			}
		}
	}
}

func (b *OGame) getAttacks(celestialID CelestialID) (out []AttackEvent, err error) {
	params := url.Values{"page": {"eventList"}, "ajax": {"1"}}
	if b.IsV7() {
		params = url.Values{"page": {"componentOnly"}, "component": {"eventList"}, "ajax": {"1"}}
	}
	if celestialID != 0 {
		params.Set("cp", strconv.FormatInt(int64(celestialID), 10))
	}
	pageHTML, err := b.getPageContent(params)
	if err != nil {
		return
	}
	out, err = b.extractor.ExtractAttacks(pageHTML)
	if err != nil {
		return
	}
	planets := b.GetCachedPlanets()
	fixAttackEvents(out, planets)
	return
}

func (b *OGame) galaxyInfos(galaxy, system int64, options ...Option) (SystemInfos, error) {
	if galaxy < 0 || galaxy > b.server.Settings.UniverseSize {
		return SystemInfos{}, fmt.Errorf("galaxy must be within [0, %d]", b.server.Settings.UniverseSize)
	}
	if system < 0 || system > b.serverData.Systems {
		return SystemInfos{}, errors.New("system must be within [0, " + strconv.FormatInt(b.serverData.Systems, 10) + "]")
	}
	payload := url.Values{
		"galaxy": {strconv.FormatInt(galaxy, 10)},
		"system": {strconv.FormatInt(system, 10)},
	}
	var res SystemInfos
	vals := url.Values{"page": {"galaxyContent"}, "ajax": {"1"}}
	if b.IsV7() {
		vals = url.Values{"page": {"ingame"}, "component": {"galaxyContent"}, "ajax": {"1"}}
	}
	pageHTML, err := b.postPageContent(vals, payload, options...)
	if err != nil {
		return res, err
	}
	return b.extractor.ExtractGalaxyInfos(pageHTML, b.Player.PlayerName, b.Player.PlayerID, b.Player.Rank)
}

func (b *OGame) getResourceSettings(planetID PlanetID) (ResourceSettings, error) {
	pageHTML, _ := b.getPage(ResourceSettingsPage, planetID.Celestial())
	return b.extractor.ExtractResourceSettings(pageHTML)
}

func (b *OGame) setResourceSettings(planetID PlanetID, settings ResourceSettings) error {
	pageHTML, _ := b.getPage(ResourceSettingsPage, planetID.Celestial())
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	bodyID := b.extractor.ExtractBodyIDFromDoc(doc)
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
		"last1":        {strconv.FormatInt(settings.MetalMine, 10)},
		"last2":        {strconv.FormatInt(settings.CrystalMine, 10)},
		"last3":        {strconv.FormatInt(settings.DeuteriumSynthesizer, 10)},
		"last4":        {strconv.FormatInt(settings.SolarPlant, 10)},
		"last12":       {strconv.FormatInt(settings.FusionReactor, 10)},
		"last212":      {strconv.FormatInt(settings.SolarSatellite, 10)},
		"last217":      {strconv.FormatInt(settings.Crawler, 10)},
	}
	url2 := b.serverURL + "/game/index.php?page=resourceSettings"
	resp, err := b.Client.PostForm(url2, payload)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			b.error(err)
		}
	}()
	return nil
}

func getNbr(doc *goquery.Document, name string) int64 {
	div := doc.Find("div." + name)
	level := div.Find("span.level")
	level.Children().Remove()
	return int64(ParseInt(level.Text()))
}

func getNbrShips(doc *goquery.Document, name string) int64 {
	div := doc.Find("div." + name)
	title := div.AttrOr("title", "")
	if title == "" {
		title = div.Find("a").AttrOr("title", "")
	}
	m := regexp.MustCompile(`.+\(([\d.,]+)\)`).FindStringSubmatch(title)
	if len(m) != 2 {
		return 0
	}
	return ParseInt(m[1])
}

func (b *OGame) getCachedResearch() Researches {
	if b.researches == nil {
		return b.getResearch()
	}
	return *b.researches
}

func (b *OGame) getResearch() Researches {
	pageHTML, _ := b.getPage(ResearchPage, CelestialID(0))
	researches := b.extractor.ExtractResearch(pageHTML)
	b.researches = &researches
	return researches
}

func (b *OGame) getResourcesBuildings(celestialID CelestialID) (ResourcesBuildings, error) {
	pageHTML, _ := b.getPage(ResourcesPage, celestialID)
	return b.extractor.ExtractResourcesBuildings(pageHTML)
}

func (b *OGame) getDefense(celestialID CelestialID) (DefensesInfos, error) {
	pageHTML, _ := b.getPage(DefensePage, celestialID)
	return b.extractor.ExtractDefense(pageHTML)
}

func (b *OGame) getShips(celestialID CelestialID) (ShipsInfos, error) {
	pageHTML, _ := b.getPage(ShipyardPage, celestialID)
	return b.extractor.ExtractShips(pageHTML)
}

func (b *OGame) getFacilities(celestialID CelestialID) (Facilities, error) {
	pageHTML, _ := b.getPage(StationPage, celestialID)
	return b.extractor.ExtractFacilities(pageHTML)
}

func (b *OGame) getProduction(celestialID CelestialID) ([]Quantifiable, int64, error) {
	pageHTML, _ := b.getPage(ShipyardPage, celestialID)
	return b.extractor.ExtractProduction(pageHTML)
}

// IsV7 ...
func (b *OGame) IsV7() bool {
	return len(b.ServerVersion()) > 0 && b.ServerVersion()[0] == '7'
}

func getToken(b *OGame, page string, celestialID CelestialID) (string, error) {
	pageHTML, _ := b.getPage(page, celestialID)
	if b.IsV7() {
		rgx := regexp.MustCompile(`var upgradeEndpoint = ".+&token=([^&]+)&`)
		m := rgx.FindSubmatch(pageHTML)
		if len(m) != 2 {
			return "", errors.New("unable to find form token")
		}
		return string(m[1]), nil
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return "", err
	}
	token, exists := doc.Find("form").Find("input[name=token]").Attr("value")
	if !exists {
		return "", errors.New("unable to find form token")
	}
	return token, nil
}

func getDemolishToken(b *OGame, page string, celestialID CelestialID) (string, error) {
	pageHTML, _ := b.getPage(page, celestialID)
	m := regexp.MustCompile(`modus=3&token=([^&]+)&`).FindSubmatch(pageHTML)
	if len(m) != 2 {
		return "", errors.New("unable to find form token")
	}
	return string(m[1]), nil
}

func (b *OGame) tearDown(celestialID CelestialID, id ID) error {
	var page string
	if id.IsResourceBuilding() {
		page = "supplies"
	} else if id.IsFacility() {
		page = "facilities"
	} else {
		return errors.New("invalid id " + id.String())
	}

	token, err := getDemolishToken(b, page, celestialID)
	if err != nil {
		return err
	}

	pageHTML, _ := b.getPageContent(url.Values{
		"page":       {"ingame"},
		"component":  {"technologydetails"},
		"ajax":       {"1"},
		"action":     {"getDetails"},
		"technology": {strconv.FormatInt(int64(id), 10)},
		"cp":         {strconv.FormatInt(int64(celestialID), 10)},
	})

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return err
	}
	imgDisabled := doc.Find("a.demolish_link div").HasClass("demolish_img_disabled")
	if imgDisabled {
		return errors.New("tear down button is disabled")
	}

	params := url.Values{
		"page":      {"ingame"},
		"component": {page},
		"modus":     {"3"},
		"token":     {token},
		"type":      {strconv.FormatInt(int64(id), 10)},
		"cp":        {strconv.FormatInt(int64(celestialID), 10)},
	}
	_, err = b.getPageContent(params)
	return err
}

func (b *OGame) build(celestialID CelestialID, id ID, nbr int64) error {
	if b.IsV7() {
		return b.buildV7(celestialID, id, nbr)
	}
	return b.buildV6(celestialID, id, nbr)
}

func (b *OGame) buildV6(celestialID CelestialID, id ID, nbr int64) error {
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
		"type":  {strconv.FormatInt(int64(id), 10)},
	}

	// Techs don't have a token
	if !id.IsTech() {
		token, err := getToken(b, page, celestialID)
		if err != nil {
			return err
		}
		payload.Add("token", token)
	}

	if id.IsDefense() || id.IsShip() {
		var maximumNbr int64 = 99999
		var err error
		var token string
		for nbr > 0 {
			tmp := int64(math.Min(float64(nbr), float64(maximumNbr)))
			payload.Set("menge", strconv.FormatInt(tmp, 10))
			_, err = b.postPageContent(url.Values{"page": {page}, "cp": {strconv.FormatInt(int64(celestialID), 10)}}, payload)
			if err != nil {
				break
			}
			token, err = getToken(b, page, celestialID)
			if err != nil {
				break
			}
			payload.Set("token", token)
			nbr -= maximumNbr
		}
		return err
	}

	_, err := b.postPageContent(url.Values{"page": {page}, "cp": {strconv.FormatInt(int64(celestialID), 10)}}, payload)
	return err
}

func (b *OGame) buildV7(celestialID CelestialID, id ID, nbr int64) error {
	var page string
	if id.IsDefense() {
		page = DefensesPage
	} else if id.IsShip() {
		page = ShipyardPage
	} else if id.IsBuilding() {
		page = SuppliesPage
	} else if id.IsTech() {
		page = ResearchPage
	} else {
		return errors.New("invalid id " + id.String())
	}
	vals := url.Values{
		"page":      {"ingame"},
		"component": {page},
		"modus":     {"1"},
		"type":      {strconv.FormatInt(int64(id), 10)},
		"cp":        {strconv.FormatInt(int64(celestialID), 10)},
	}

	// Techs don't have a token
	if !id.IsTech() {
		token, err := getToken(b, page, celestialID)
		if err != nil {
			return err
		}
		vals.Add("token", token)
	}

	if id.IsDefense() || id.IsShip() {
		var maximumNbr int64 = 99999
		var err error
		var token string
		for nbr > 0 {
			tmp := int64(math.Min(float64(nbr), float64(maximumNbr)))
			vals.Set("menge", strconv.FormatInt(tmp, 10))
			_, err = b.getPageContent(vals)
			if err != nil {
				break
			}
			token, err = getToken(b, page, celestialID)
			if err != nil {
				break
			}
			vals.Set("token", token)
			nbr -= maximumNbr
		}
		return err
	}

	_, err := b.getPageContent(vals)
	return err
}

func (b *OGame) buildCancelable(celestialID CelestialID, id ID) error {
	if !id.IsBuilding() && !id.IsTech() {
		return errors.New("invalid id " + id.String())
	}
	return b.build(celestialID, id, 0)
}

func (b *OGame) buildProduction(celestialID CelestialID, id ID, nbr int64) error {
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

func (b *OGame) buildTechnology(celestialID CelestialID, technologyID ID) error {
	if !technologyID.IsTech() {
		return errors.New("invalid technology id " + technologyID.String())
	}
	return b.buildCancelable(celestialID, technologyID)
}

func (b *OGame) buildDefense(celestialID CelestialID, defenseID ID, nbr int64) error {
	if !defenseID.IsDefense() {
		return errors.New("invalid defense id " + defenseID.String())
	}
	return b.buildProduction(celestialID, ID(defenseID), nbr)
}

func (b *OGame) buildShips(celestialID CelestialID, shipID ID, nbr int64) error {
	if !shipID.IsShip() {
		return errors.New("invalid ship id " + shipID.String())
	}
	return b.buildProduction(celestialID, shipID, nbr)
}

func (b *OGame) constructionsBeingBuilt(celestialID CelestialID) (ID, int64, ID, int64) {
	pageHTML, _ := b.getPage(OverviewPage, celestialID)
	return b.extractor.ExtractConstructions(pageHTML)
}

func (b *OGame) cancel(token string, techID, listID int64) error {
	if b.IsV7() {
		_, _ = b.getPageContent(url.Values{"page": {"ingame"}, "component": {"overview"}, "modus": {"2"}, "token": {token},
			"type": {strconv.FormatInt(techID, 10)}, "listid": {strconv.FormatInt(listID, 10)}, "action": {"cancel"}})
	} else {
		_, _ = b.getPageContent(url.Values{"page": {"overview"}, "modus": {"2"}, "token": {token},
			"techid": {strconv.FormatInt(techID, 10)}, "listid": {strconv.FormatInt(listID, 10)}})
	}
	return nil
}

func (b *OGame) cancelBuilding(celestialID CelestialID) error {
	pageHTML, err := b.getPage(OverviewPage, celestialID)
	if err != nil {
		return err
	}
	token, techID, listID, _ := b.extractor.ExtractCancelBuildingInfos(pageHTML)
	return b.cancel(token, techID, listID)
}

func (b *OGame) cancelResearch(celestialID CelestialID) error {
	pageHTML, err := b.getPage(OverviewPage, celestialID)
	if err != nil {
		return err
	}
	token, techID, listID, _ := b.extractor.ExtractCancelResearchInfos(pageHTML)
	return b.cancel(token, techID, listID)
}

func (b *OGame) fetchResources(celestialID CelestialID) (ResourcesDetails, error) {
	pageJSON, err := b.getPage(FetchResourcesPage, celestialID)
	if err != nil {
		return ResourcesDetails{}, err
	}
	return b.extractor.ExtractResourcesDetails(pageJSON)
}

func (b *OGame) getResources(celestialID CelestialID) (Resources, error) {
	res, err := b.fetchResources(celestialID)
	if err != nil {
		return Resources{}, err
	}
	return Resources{
		Metal:      res.Metal.Available,
		Crystal:    res.Crystal.Available,
		Deuterium:  res.Deuterium.Available,
		Energy:     res.Energy.Available,
		Darkmatter: res.Darkmatter.Available,
	}, nil
}

func (b *OGame) getResourcesDetails(celestialID CelestialID) (ResourcesDetails, error) {
	return b.fetchResources(celestialID)
}

func (b *OGame) sendIPM(planetID PlanetID, coord Coordinate, nbr int64, priority ID) (int64, error) {
	if priority != 0 && (!priority.IsDefense() || priority == AntiBallisticMissilesID || priority == InterplanetaryMissilesID) {
		return 0, errors.New("invalid defense target id")
	}
	vals := url.Values{
		"page":       {"missileattacklayer"},
		"galaxy":     {strconv.FormatInt(coord.Galaxy, 10)},
		"system":     {strconv.FormatInt(coord.System, 10)},
		"position":   {strconv.FormatInt(coord.Position, 10)},
		"planetType": {strconv.FormatInt(int64(coord.Type), 10)},
		"cp":         {strconv.FormatInt(int64(planetID), 10)},
	}
	if b.IsV7() {
		vals = url.Values{
			"page":       {"ajax"},
			"component":  {"missileattacklayer"},
			"galaxy":     {strconv.FormatInt(coord.Galaxy, 10)},
			"system":     {strconv.FormatInt(coord.System, 10)},
			"position":   {strconv.FormatInt(coord.Position, 10)},
			"planetType": {strconv.FormatInt(int64(coord.Type), 10)},
			"cp":         {strconv.FormatInt(int64(planetID), 10)},
		}
	}
	pageHTML, err := b.getPageContent(vals)
	if err != nil {
		return 0, err
	}
	duration, max, token := b.extractor.ExtractIPM(pageHTML)
	if max == 0 {
		return 0, errors.New("no missile available")
	}
	if nbr > max {
		nbr = max
	}
	payload := url.Values{
		"galaxy":     {strconv.FormatInt(coord.Galaxy, 10)},
		"system":     {strconv.FormatInt(coord.System, 10)},
		"position":   {strconv.FormatInt(coord.Position, 10)},
		"planetType": {strconv.FormatInt(int64(coord.Type), 10)},
		"token":      {token},
		"anz":        {strconv.FormatInt(nbr, 10)},
		"pziel":      {},
	}
	if priority != 0 {
		payload.Add("pziel", strconv.FormatInt(int64(priority), 10))
	}
	params := url.Values{"page": {"missileattack_execute"}}
	if b.IsV7() {
		params = url.Values{
			"page":      {"ajax"},
			"component": {"missileattacklayer"},
			"action":    {"sendMissiles"},
			"ajax":      {"1"},
			"asJson":    {"1"},
		}
		payload = url.Values{
			"galaxy":               {strconv.FormatInt(coord.Galaxy, 10)},
			"system":               {strconv.FormatInt(coord.System, 10)},
			"position":             {strconv.FormatInt(coord.Position, 10)},
			"type":                 {strconv.FormatInt(int64(coord.Type), 10)},
			"token":                {token},
			"missileCount":         {strconv.FormatInt(nbr, 10)},
			"missilePrimaryTarget": {},
		}
		if priority != 0 {
			payload.Add("missilePrimaryTarget", strconv.FormatInt(int64(priority), 10))
		}
	}
	by, err := b.postPageContent(params, payload)
	if err != nil {
		return 0, err
	}
	// {"status":false,"errorbox":{"type":"fadeBox","text":"Target doesn`t exist!","failed":1}} // OgameV6
	// {"status":true,"rockets":0,"errorbox":{"type":"fadeBox","text":"25 raketten zijn gelanceerd!","failed":0},"components":[]} // OgameV7
	var resp struct {
		Status   bool
		Rockets  int64
		ErrorBox struct {
			Type   string
			Text   string
			Failed int
		}
		// components??
	}
	if err := json.Unmarshal(by, &resp); err != nil {
		return 0, err
	}
	if resp.ErrorBox.Failed == 1 {
		return 0, errors.New(resp.ErrorBox.Text)
	}

	return duration, nil
}

func (b *OGame) sendFleet(celestialID CelestialID, ships []Quantifiable, speed Speed, where Coordinate,
	mission MissionID, resources Resources, expeditiontime, unionID int64, ensure bool) (Fleet, error) {
	if b.IsV7() {
		return b.sendFleetV7(celestialID, ships, speed, where, mission, resources, expeditiontime, unionID, ensure)
	}
	return b.sendFleetV6(celestialID, ships, speed, where, mission, resources, expeditiontime, unionID, ensure)
}

// CheckTargetResponse ...
type CheckTargetResponse struct {
	Status string `json:"status"`
	Orders struct {
		Num1  bool `json:"1"`
		Num2  bool `json:"2"`
		Num3  bool `json:"3"`
		Num4  bool `json:"4"`
		Num5  bool `json:"5"`
		Num6  bool `json:"6"`
		Num7  bool `json:"7"`
		Num8  bool `json:"8"`
		Num9  bool `json:"9"`
		Num15 bool `json:"15"`
	} `json:"orders"`
	TargetInhabited           bool   `json:"targetInhabited"`
	TargetIsStrong            bool   `json:"targetIsStrong"`
	TargetIsOutlaw            bool   `json:"targetIsOutlaw"`
	TargetIsBuddyOrAllyMember bool   `json:"targetIsBuddyOrAllyMember"`
	TargetPlayerID            int    `json:"targetPlayerId"`
	TargetPlayerName          string `json:"targetPlayerName"`
	TargetPlayerColorClass    string `json:"targetPlayerColorClass"`
	TargetPlayerRankIcon      string `json:"targetPlayerRankIcon"`
	PlayerIsOutlaw            bool   `json:"playerIsOutlaw"`
	TargetPlanet              struct {
		Galaxy   int    `json:"galaxy"`
		System   int    `json:"system"`
		Position int    `json:"position"`
		Type     int    `json:"type"`
		Name     string `json:"name"`
	} `json:"targetPlanet"`
	TargetOk   bool          `json:"targetOk"`
	Components []interface{} `json:"components"`
}

func (b *OGame) sendFleetV7(celestialID CelestialID, ships []Quantifiable, speed Speed, where Coordinate,
	mission MissionID, resources Resources, expeditiontime, unionID int64, ensure bool) (Fleet, error) {

	// Get existing fleet, so we can ensure new fleet ID is greater
	initialFleets, slots := b.getFleets()
	maxInitialFleetID := FleetID(0)
	for _, f := range initialFleets {
		if f.ID > maxInitialFleetID {
			maxInitialFleetID = f.ID
		}
	}

	if slots.InUse == slots.Total {
		return Fleet{}, ErrAllSlotsInUse
	}

	if mission == Expedition {
		if slots.ExpInUse == slots.ExpTotal {
			return Fleet{}, ErrAllSlotsInUse
		}
	}

	// Page 1 : get to fleet page
	pageHTML, err := b.getPage(Fleet1Page, celestialID)
	if err != nil {
		return Fleet{}, err
	}

	fleet1Doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	fleet1BodyID := b.extractor.ExtractBodyIDFromDoc(fleet1Doc)
	if fleet1BodyID != FleetdispatchPage {
		now := time.Now().Unix()
		b.error(ErrInvalidPlanetID.Error()+", planetID:", celestialID, ", ts: ", now)
		return Fleet{}, ErrInvalidPlanetID
	}

	if b.extractor.ExtractIsInVacationFromDoc(fleet1Doc) {
		return Fleet{}, ErrAccountInVacationMode
	}

	// Ensure we're not trying to attack/spy ourselves
	destinationIsMyOwnPlanet := false
	myCelestials, _ := b.extractor.ExtractCelestialsFromDoc(fleet1Doc, b)
	for _, c := range myCelestials {
		if c.GetCoordinate().Equal(where) && c.GetID() == celestialID {
			return Fleet{}, errors.New("origin and destination are the same")
		}
		if c.GetCoordinate().Equal(where) {
			destinationIsMyOwnPlanet = true
			break
		}
	}
	if destinationIsMyOwnPlanet {
		switch mission {
		case Spy:
			return Fleet{}, errors.New("you cannot spy yourself")
		case Attack:
			return Fleet{}, errors.New("you cannot attack yourself")
		}
	}

	availableShips := b.extractor.ExtractFleet1ShipsFromDoc(fleet1Doc)

	atLeastOneShipSelected := false
	if !ensure {
		for i := range ships {
			avail := availableShips.ByID(ships[i].ID)
			ships[i].Nbr = int64(math.Min(float64(ships[i].Nbr), float64(avail)))
			if ships[i].Nbr > 0 {
				atLeastOneShipSelected = true
			}
		}
	} else {
		for _, ship := range ships {
			if ship.Nbr > availableShips.ByID(ship.ID) {
				return Fleet{}, ErrNotEnoughShips
			}
			atLeastOneShipSelected = true
		}
	}
	if !atLeastOneShipSelected {
		return Fleet{}, ErrNoShipSelected
	}

	payload := b.extractor.ExtractHiddenFieldsFromDoc(fleet1Doc)
	for _, s := range ships {
		if s.Nbr > 0 {
			payload.Set("am"+strconv.FormatInt(int64(s.ID), 10), strconv.FormatInt(s.Nbr, 10))
		}
	}

	tokenM := regexp.MustCompile(`var fleetSendingToken = "([^"]+)";`).FindSubmatch(pageHTML)
	if len(tokenM) != 2 {
		return Fleet{}, errors.New("token not found")
	}

	payload.Set("token", string(tokenM[1]))
	payload.Set("galaxy", strconv.FormatInt(where.Galaxy, 10))
	payload.Set("system", strconv.FormatInt(where.System, 10))
	payload.Set("position", strconv.FormatInt(where.Position, 10))
	if mission == RecycleDebrisField {
		where.Type = DebrisType // Send to debris field
	} else if mission == Colonize || mission == Expedition {
		where.Type = PlanetType
	}
	payload.Set("type", strconv.FormatInt(int64(where.Type), 10))
	payload.Set("union", "0")

	if unionID != 0 {
		found := false
		fleet1Doc.Find("select[name=acsValues] option").Each(func(i int, s *goquery.Selection) {
			acsValues := s.AttrOr("value", "")
			m := regexp.MustCompile(`\d+#\d+#\d+#\d+#.*#(\d+)`).FindStringSubmatch(acsValues)
			if len(m) == 2 {
				optUnionID, _ := strconv.ParseInt(m[1], 10, 64)
				if unionID == optUnionID {
					found = true
					payload.Add("acsValues", acsValues)
					payload.Add("union", m[1])
					mission = GroupedAttack
				}
			}
		})
		if !found {
			return Fleet{}, ErrUnionNotFound
		}
	}

	// Check
	by1, err := b.postPageContent(url.Values{"page": {"ingame"}, "component": {"fleetdispatch"}, "action": {"checkTarget"}, "ajax": {"1"}, "asJson": {"1"}}, payload)
	if err != nil {
		b.error(err.Error())
		return Fleet{}, err
	}
	var checkRes CheckTargetResponse
	if err := json.Unmarshal(by1, &checkRes); err != nil {
		b.error(err.Error())
		return Fleet{}, err
	}

	if !checkRes.TargetOk {
		return Fleet{}, errors.New("target is not ok")
	}

	cargo := ShipsInfos{}.FromQuantifiables(ships).Cargo(b.getCachedResearch(), b.server.Settings.EspionageProbeRaids == 1, b.isCollector())
	newResources := Resources{}
	if resources.Total() > cargo {
		newResources.Deuterium = int64(math.Min(float64(resources.Deuterium), float64(cargo)))
		cargo -= newResources.Deuterium
		newResources.Crystal = int64(math.Min(float64(resources.Crystal), float64(cargo)))
		cargo -= newResources.Crystal
		newResources.Metal = int64(math.Min(float64(resources.Metal), float64(cargo)))
		cargo -= newResources.Metal
	} else {
		newResources = resources
	}

	newResources.Metal = MaxInt(newResources.Metal, 0)
	newResources.Crystal = MaxInt(newResources.Crystal, 0)
	newResources.Deuterium = MaxInt(newResources.Deuterium, 0)

	// Page 3 : select coord, mission, speed
	payload.Set("speed", strconv.FormatInt(int64(speed), 10))
	payload.Set("crystal", strconv.FormatInt(newResources.Crystal, 10))
	payload.Set("deuterium", strconv.FormatInt(newResources.Deuterium, 10))
	payload.Set("metal", strconv.FormatInt(newResources.Metal, 10))
	payload.Set("mission", strconv.FormatInt(int64(mission), 10))
	if mission == Expedition {
		if expeditiontime <= 0 {
			expeditiontime = 1
		}
		payload.Set("holdingtime", strconv.FormatInt(expeditiontime, 10))
	}

	// Page 4 : send the fleet
	res, _ := b.postPageContent(url.Values{"page": {"ingame"}, "component": {"fleetdispatch"}, "action": {"sendFleet"}, "ajax": {"1"}, "asJson": {"1"}}, payload)
	// {"success":true,"message":"Your fleet has been successfully sent.","redirectUrl":"https:\/\/s801-en.ogame.gameforge.com\/game\/index.php?page=ingame&component=fleetdispatch","components":[]}
	// Insufficient resources. (4060)
	// {"success":false,"errors":[{"message":"Not enough cargo space!","error":4029}],"fleetSendingToken":"b4786751c6d5e64e56d8eb94807fbf88","components":[]}
	// {"success":false,"errors":[{"message":"Fleet launch failure: The fleet could not be launched. Please try again later.","error":4047}],"fleetSendingToken":"1507c7228b206b4a298dec1d34a5a207","components":[]} // bad token I think
	// {"success":false,"errors":[{"message":"Recyclers must be sent to recycle this debris field!","error":4013}],"fleetSendingToken":"b826ff8c3d4e04066c28d10399b32ab8","components":[]}
	// {"success":false,"errors":[{"message":"Error, no ships available","error":4059}],"fleetSendingToken":"b369e37ce34bb64e3a59fa26bd8d5602","components":[]}
	// {"success":false,"errors":[{"message":"You have to select a valid target.","error":4049}],"fleetSendingToken":"19218f446d0985dfd79e03c3ec008514","components":[]} // colonize debris field
	// {"success":false,"errors":[{"message":"Planet is already inhabited!","error":4053}],"fleetSendingToken":"3281f9ad5b4cba6c0c26a24d3577bd4c","components":[]}
	// {"success":false,"errors":[{"message":"Colony ships must be sent to colonise this planet!","error":4038}],"fleetSendingToken":"8700c275a055c59ca276a7f66c81b205","components":[]}
	// fetch("https://s801-en.ogame.gameforge.com/game/index.php?page=ingame&component=fleetdispatch&action=sendFleet&ajax=1&asJson=1", {"credentials":"include","headers":{"content-type":"application/x-www-form-urlencoded; charset=UTF-8","sec-fetch-mode":"cors","sec-fetch-site":"same-origin","x-requested-with":"XMLHttpRequest"},"body":"token=414847e59344881d5c71303023735ab8&am209=1&am202=10&galaxy=9&system=297&position=7&type=2&metal=0&crystal=0&deuterium=0&prioMetal=1&prioCrystal=2&prioDeuterium=3&mission=8&speed=1&retreatAfterDefenderRetreat=0&union=0&holdingtime=0","method":"POST","mode":"cors"}).then(res => res.json()).then(r => console.log(r));

	var resStruct struct {
		Success           bool          `json:"success"`
		Message           string        `json:"message"`
		FleetSendingToken string        `json:"fleetSendingToken"`
		Components        []interface{} `json:"components"`
		RedirectURL       string        `json:"redirectUrl"`
		Errors            []struct {
			Message string `json:"message"`
			Error   int64  `json:"error"`
		} `json:"errors"`
	}
	if err := json.Unmarshal(res, &resStruct); err != nil {
		return Fleet{}, errors.New("failed to unmarshal response: " + err.Error())
	}

	if len(resStruct.Errors) > 0 {
		return Fleet{}, errors.New(resStruct.Errors[0].Message + " (" + strconv.FormatInt(resStruct.Errors[0].Error, 10) + ")")
	}

	// Page 5
	movementHTML, _ := b.getPage(MovementPage, CelestialID(0))
	movementDoc, _ := goquery.NewDocumentFromReader(bytes.NewReader(movementHTML))
	originCoords, _ := b.extractor.ExtractPlanetCoordinate(movementHTML)
	fleets := b.extractor.ExtractFleetsFromDoc(movementDoc)
	if len(fleets) > 0 {
		max := Fleet{}
		for i, fleet := range fleets {
			if fleet.ID > max.ID &&
				fleet.Origin.Equal(originCoords) &&
				fleet.Destination.Equal(where) &&
				fleet.Mission == mission &&
				!fleet.ReturnFlight {
				max = fleets[i]
			}
		}
		if max.ID > maxInitialFleetID {
			return max, nil
		}
	}

	slots = b.extractor.ExtractSlotsFromDoc(movementDoc)
	if slots.InUse == slots.Total {
		return Fleet{}, ErrAllSlotsInUse
	}

	if mission == Expedition {
		if slots.ExpInUse == slots.ExpTotal {
			return Fleet{}, ErrAllSlotsInUse
		}
	}

	now := time.Now().Unix()
	b.error(errors.New("could not find new fleet ID").Error()+", planetID:", celestialID, ", ts: ", now)
	return Fleet{}, errors.New("could not find new fleet ID")
}

func (b *OGame) sendFleetV6(celestialID CelestialID, ships []Quantifiable, speed Speed, where Coordinate,
	mission MissionID, resources Resources, expeditiontime, unionID int64, ensure bool) (Fleet, error) {

	// Get existing fleet, so we can ensure new fleet ID is greater
	initialFleets, slots := b.getFleets()
	maxInitialFleetID := FleetID(0)
	for _, f := range initialFleets {
		if f.ID > maxInitialFleetID {
			maxInitialFleetID = f.ID
		}
	}

	if slots.InUse == slots.Total {
		return Fleet{}, ErrAllSlotsInUse
	}

	if mission == Expedition {
		if slots.ExpInUse == slots.ExpTotal {
			return Fleet{}, ErrAllSlotsInUse
		}
	}

	// Page 1 : get to fleet page
	pageHTML, err := b.getPage(Fleet1Page, celestialID)
	if err != nil {
		return Fleet{}, err
	}

	fleet1Doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	fleet1BodyID := b.extractor.ExtractBodyIDFromDoc(fleet1Doc)
	if fleet1BodyID != "fleet1" {
		now := time.Now().Unix()
		b.error(ErrInvalidPlanetID.Error()+", planetID:", celestialID, ", ts: ", now)
		return Fleet{}, ErrInvalidPlanetID
	}

	if b.extractor.ExtractIsInVacationFromDoc(fleet1Doc) {
		return Fleet{}, ErrAccountInVacationMode
	}

	// Ensure we're not trying to attack/spy ourselves
	destinationIsMyOwnPlanet := false
	myCelestials, _ := b.extractor.ExtractCelestialsFromDoc(fleet1Doc, b)
	for _, c := range myCelestials {
		if c.GetCoordinate().Equal(where) && c.GetID() == celestialID {
			return Fleet{}, errors.New("origin and destination are the same")
		}
		if c.GetCoordinate().Equal(where) {
			destinationIsMyOwnPlanet = true
			break
		}
	}
	if destinationIsMyOwnPlanet {
		switch mission {
		case Spy:
			return Fleet{}, errors.New("you cannot spy yourself")
		case Attack:
			return Fleet{}, errors.New("you cannot attack yourself")
		}
	}

	availableShips := b.extractor.ExtractFleet1ShipsFromDoc(fleet1Doc)

	atLeastOneShipSelected := false
	if !ensure {
		for _, ship := range ships {
			if ship.Nbr > 0 && availableShips.ByID(ship.ID) > 0 {
				atLeastOneShipSelected = true
				break
			}
		}
	} else {
		for _, ship := range ships {
			if ship.Nbr > availableShips.ByID(ship.ID) {
				return Fleet{}, ErrNotEnoughShips
			}
			atLeastOneShipSelected = true
		}
	}
	if !atLeastOneShipSelected {
		return Fleet{}, ErrNoShipSelected
	}

	payload := b.extractor.ExtractHiddenFieldsFromDoc(fleet1Doc)
	cs := false       // ColonyShip flag for fleet check
	recycler := false // Recycler flag for fleet check
	for _, s := range ships {
		if s.Nbr > 0 {
			if s.ID == ColonyShipID {
				cs = true
			} else if s.ID == RecyclerID {
				recycler = true
			}
			payload.Add("am"+strconv.FormatInt(int64(s.ID), 10), strconv.FormatInt(s.Nbr, 10))
		}
	}

	// Page 2 : select ships
	pageHTML, err = b.postPageContent(url.Values{"page": {"fleet2"}}, payload)
	if err != nil {
		return Fleet{}, err
	}
	fleet2Doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	fleet2BodyID := b.extractor.ExtractBodyIDFromDoc(fleet2Doc)
	if fleet2BodyID != "fleet2" {
		now := time.Now().Unix()
		b.error(errors.New("unknown error").Error()+", planetID:", celestialID, ", ts: ", now)
		return Fleet{}, errors.New("unknown error")
	}

	payload = b.extractor.ExtractHiddenFieldsFromDoc(fleet2Doc)
	payload.Add("speed", strconv.FormatInt(int64(speed), 10))
	payload.Add("galaxy", strconv.FormatInt(where.Galaxy, 10))
	payload.Add("system", strconv.FormatInt(where.System, 10))
	payload.Add("position", strconv.FormatInt(where.Position, 10))
	if mission == RecycleDebrisField {
		where.Type = DebrisType // Send to debris field
	} else if mission == Colonize || mission == Expedition {
		where.Type = PlanetType
	}
	payload.Add("type", strconv.FormatInt(int64(where.Type), 10))

	if unionID != 0 {
		found := false
		fleet2Doc.Find("select[name=acsValues] option").Each(func(i int, s *goquery.Selection) {
			acsValues := s.AttrOr("value", "")
			m := regexp.MustCompile(`\d+#\d+#\d+#\d+#.*#(\d+)`).FindStringSubmatch(acsValues)
			if len(m) == 2 {
				optUnionID, _ := strconv.ParseInt(m[1], 10, 64)
				if unionID == optUnionID {
					found = true
					payload.Add("acsValues", acsValues)
					payload.Add("union", m[1])
					mission = GroupedAttack
				}
			}
		})
		if !found {
			return Fleet{}, ErrUnionNotFound
		}
	}

	// Check
	fleetCheckPayload := url.Values{
		"galaxy": {strconv.FormatInt(where.Galaxy, 10)},
		"system": {strconv.FormatInt(where.System, 10)},
		"planet": {strconv.FormatInt(where.Position, 10)},
		"type":   {strconv.FormatInt(int64(where.Type), 10)},
	}
	if cs {
		fleetCheckPayload.Add("cs", "1")
	}
	if recycler {
		fleetCheckPayload.Add("recycler", "1")
	}
	by1, err := b.postPageContent(url.Values{"page": {"fleetcheck"}, "ajax": {"1"}, "espionage": {"0"}}, fleetCheckPayload)
	if err != nil {
		return Fleet{}, err
	}
	switch string(by1) {
	case "1":
		return Fleet{}, ErrUninhabitedPlanet
	case "1d":
		return Fleet{}, ErrNoDebrisField
	case "2":
		return Fleet{}, ErrPlayerInVacationMode
	case "3":
		return Fleet{}, ErrAdminOrGM
	case "4":
		return Fleet{}, ErrNoAstrophysics
	case "5":
		return Fleet{}, ErrNoobProtection
	case "6":
		return Fleet{}, ErrPlayerTooStrong
	case "10":
		return Fleet{}, ErrNoMoonAvailable
	case "11":
		return Fleet{}, ErrNoRecyclerAvailable
	case "15":
		return Fleet{}, ErrNoEventsRunning
	case "16":
		return Fleet{}, ErrPlanetAlreadyReservedForRelocation
	}

	// Page 3 : select coord, mission, speed
	pageHTML, err = b.postPageContent(url.Values{"page": {"fleet3"}}, payload)
	if err != nil {
		return Fleet{}, err
	}

	fleet3Doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	fleet3BodyID := b.extractor.ExtractBodyIDFromDoc(fleet3Doc)
	if fleet3BodyID != "fleet3" {
		now := time.Now().Unix()
		b.error(errors.New("unknown error").Error()+", planetID:", celestialID, ", ts: ", now)
		return Fleet{}, errors.New("unknown error")
	}

	if mission == Spy && fleet3Doc.Find("li#button6").HasClass("off") {
		return Fleet{}, errors.New("target cannot be spied (button disabled)")
	} else if mission == Attack && fleet3Doc.Find("li#button1").HasClass("off") {
		return Fleet{}, errors.New("target cannot be attacked (button disabled)")
	} else if mission == Transport && fleet3Doc.Find("li#button3").HasClass("off") {
		return Fleet{}, errors.New("cannot send transport (button disabled)")
	} else if mission == Park && fleet3Doc.Find("li#button4").HasClass("off") {
		return Fleet{}, errors.New("cannot send deployment (button disabled)")
	} else if mission == Colonize && fleet3Doc.Find("li#button7").HasClass("off") {
		return Fleet{}, errors.New("cannot send colonisation (button disabled)")
	} else if mission == Expedition && fleet3Doc.Find("li#button15").HasClass("off") {
		return Fleet{}, errors.New("cannot send expedition (button disabled)")
	} else if mission == RecycleDebrisField && fleet3Doc.Find("li#button8").HasClass("off") {
		return Fleet{}, errors.New("cannot recycle (button disabled)")
		//} else if mission == Transport && fleet3Doc.Find("li#button5").HasClass("off") {
		//	return Fleet{}, errors.New("cannot acs defend (button disabled)")
	} else if mission == GroupedAttack && fleet3Doc.Find("li#button2").HasClass("off") {
		return Fleet{}, errors.New("cannot acs attack (button disabled)")
	} else if mission == Destroy && fleet3Doc.Find("li#button9").HasClass("off") {
		return Fleet{}, errors.New("cannot destroy (button disabled)")
	}

	payload = b.extractor.ExtractHiddenFieldsFromDoc(fleet3Doc)
	var finalShips ShipsInfos
	for k, v := range payload {
		var shipID int
		if n, err := fmt.Sscanf(k, "am%d", &shipID); err == nil && n == 1 {
			nbr, _ := strconv.ParseInt(v[0], 10, 64)
			finalShips.Set(ID(shipID), nbr)
		}
	}
	deutConsumption := ParseInt(fleet3Doc.Find("div#roundup span#consumption").Text())
	resourcesAvailable := b.extractor.ExtractResourcesFromDoc(fleet3Doc)
	if deutConsumption > resourcesAvailable.Deuterium {
		return Fleet{}, fmt.Errorf("not enough deuterium, avail: %d, need: %d", resourcesAvailable.Deuterium, deutConsumption)
	}
	// finalCargo := ParseInt(fleet3Doc.Find("#maxresources").Text())
	baseCargo := finalShips.Cargo(Researches{}, b.GetServer().Settings.EspionageProbeRaids == 1, b.characterClass == Collector)
	if b.GetServer().Settings.EspionageProbeRaids != 1 {
		baseCargo += finalShips.EspionageProbe * EspionageProbe.BaseCargoCapacity
	}
	if deutConsumption > baseCargo {
		return Fleet{}, fmt.Errorf("not enough cargo capacity for fuel, avail: %d, need: %d", baseCargo, deutConsumption)
	}
	payload.Add("crystal", strconv.FormatInt(resources.Crystal, 10))
	payload.Add("deuterium", strconv.FormatInt(resources.Deuterium, 10))
	payload.Add("metal", strconv.FormatInt(resources.Metal, 10))
	payload.Set("mission", strconv.FormatInt(int64(mission), 10))
	if mission == Expedition {
		payload.Set("expeditiontime", strconv.FormatInt(expeditiontime, 10))
	}

	// Page 4 : send the fleet
	_, _ = b.postPageContent(url.Values{"page": {"movement"}}, payload)

	// Page 5
	movementHTML, _ := b.getPage(MovementPage, CelestialID(0))
	movementDoc, _ := goquery.NewDocumentFromReader(bytes.NewReader(movementHTML))
	originCoords, _ := b.extractor.ExtractPlanetCoordinate(movementHTML)
	fleets := b.extractor.ExtractFleetsFromDoc(movementDoc)
	if len(fleets) > 0 {
		max := Fleet{}
		for i, fleet := range fleets {
			if fleet.ID > max.ID &&
				fleet.Origin.Equal(originCoords) &&
				fleet.Destination.Equal(where) &&
				fleet.Mission == mission &&
				!fleet.ReturnFlight {
				max = fleets[i]
			}
		}
		if max.ID > maxInitialFleetID {
			return max, nil
		}
	}

	slots = b.extractor.ExtractSlotsFromDoc(movementDoc)
	if slots.InUse == slots.Total {
		return Fleet{}, ErrAllSlotsInUse
	}

	if mission == Expedition {
		if slots.ExpInUse == slots.ExpTotal {
			return Fleet{}, ErrAllSlotsInUse
		}
	}

	now := time.Now().Unix()
	b.error(errors.New("could not find new fleet ID").Error()+", planetID:", celestialID, ", ts: ", now)
	return Fleet{}, errors.New("could not find new fleet ID")
}

// EspionageReportType type of espionage report (action or report)
type EspionageReportType int

// Action message received when an enemy is seen naer your planet
const Action EspionageReportType = 0

// Report message received when you spied on someone
const Report EspionageReportType = 1

// CombatReportSummary summary of combat report
type CombatReportSummary struct {
	ID           int64
	APIKey       string
	Origin       *Coordinate
	Destination  Coordinate
	AttackerName string
	DefenderName string
	Loot         int64
	Metal        int64
	Crystal      int64
	Deuterium    int64
	DebrisField  int64
	CreatedAt    time.Time
}

// EspionageReportSummary summary of espionage report
type EspionageReportSummary struct {
	ID             int64
	Type           EspionageReportType
	From           string
	Target         Coordinate
	LootPercentage float64
}

func (b *OGame) getPageMessages(page, tabid int64) ([]byte, error) {
	payload := url.Values{
		"messageId":  {"-1"},
		"tabid":      {strconv.FormatInt(tabid, 10)},
		"action":     {"107"},
		"pagination": {strconv.FormatInt(page, 10)},
		"ajax":       {"1"},
	}
	return b.postPageContent(url.Values{"page": {"messages"}}, payload)
}

func (b *OGame) getEspionageReportMessages() ([]EspionageReportSummary, error) {
	var tabid int64 = 20
	var page int64 = 1
	var nbPage int64 = 1
	msgs := make([]EspionageReportSummary, 0)
	for page <= nbPage {
		pageHTML, _ := b.getPageMessages(page, tabid)
		newMessages, newNbPage := b.extractor.ExtractEspionageReportMessageIDs(pageHTML)
		msgs = append(msgs, newMessages...)
		nbPage = newNbPage
		page++
	}
	return msgs, nil
}

func (b *OGame) getCombatReportMessages() ([]CombatReportSummary, error) {
	var tabid int64 = 21
	var page int64 = 1
	var nbPage int64 = 1
	msgs := make([]CombatReportSummary, 0)
	for page <= nbPage {
		pageHTML, _ := b.getPageMessages(page, tabid)
		newMessages, newNbPage := b.extractor.ExtractCombatReportMessagesSummary(pageHTML)
		msgs = append(msgs, newMessages...)
		nbPage = newNbPage
		page++
	}
	return msgs, nil
}

func (b *OGame) getCombatReportFor(coord Coordinate) (CombatReportSummary, error) {
	var tabid int64 = 21
	var page int64 = 1
	var nbPage int64 = 1
	for page <= nbPage {
		pageHTML, err := b.getPageMessages(page, tabid)
		if err != nil {
			return CombatReportSummary{}, err
		}
		newMessages, newNbPage := b.extractor.ExtractCombatReportMessagesSummary(pageHTML)
		for _, m := range newMessages {
			if m.Destination.Equal(coord) {
				return m, nil
			}
		}
		nbPage = newNbPage
		page++
	}
	return CombatReportSummary{}, errors.New("combat report not found for " + coord.String())
}

func (b *OGame) getEspionageReport(msgID int64) (EspionageReport, error) {
	pageHTML, _ := b.getPageContent(url.Values{"page": {"messages"}, "messageId": {strconv.FormatInt(msgID, 10)}, "tabid": {"20"}, "ajax": {"1"}})
	return b.extractor.ExtractEspionageReport(pageHTML, b.location)
}

func (b *OGame) getEspionageReportFor(coord Coordinate) (EspionageReport, error) {
	var tabid int64 = 20
	var page int64 = 1
	var nbPage int64 = 1
	for page <= nbPage {
		pageHTML, err := b.getPageMessages(page, tabid)
		if err != nil {
			return EspionageReport{}, err
		}
		newMessages, newNbPage := b.extractor.ExtractEspionageReportMessageIDs(pageHTML)
		for _, m := range newMessages {
			if m.Target.Equal(coord) {
				return b.getEspionageReport(m.ID)
			}
		}
		nbPage = newNbPage
		page++
	}
	return EspionageReport{}, errors.New("espionage report not found for " + coord.String())
}

func (b *OGame) deleteMessage(msgID int64) error {
	payload := url.Values{
		"messageId": {strconv.FormatInt(msgID, 10)},
		"action":    {"103"},
		"ajax":      {"1"},
	}
	by, err := b.postPageContent(url.Values{"page": {"messages"}}, payload)
	if err != nil {
		return err
	}

	var res map[string]bool
	if err := json.Unmarshal(by, &res); err != nil {
		return errors.New("unable to find message id " + strconv.FormatInt(msgID, 10))
	}
	if val, ok := res[strconv.FormatInt(msgID, 10)]; !ok || !val {
		return errors.New("unable to find message id " + strconv.FormatInt(msgID, 10))
	}
	return nil
}

func (b *OGame) deleteAllMessagesFromTab(tabID int64) error {
	/*
		Request URL: https://$ogame/game/index.php?page=messages
		Request Method: POST

		tabid: 20 => Espionage
		tabid: 21 => Combat Reports
		tabid: 22 => Expeditions
		tabid: 23 => Unions/Transport
		tabid: 24 => Other

		E.g. :

		tabid=24&messageId=-1&action=103&ajax=1

		tabid: 24
		messageId: -1
		action: 103
		ajax: 1
	*/
	payload := url.Values{
		"tabid":     {strconv.FormatInt(tabID, 10)},
		"messageId": {strconv.FormatInt(-1, 10)},
		"action":    {"103"},
		"ajax":      {"1"},
	}
	_, err := b.postPageContent(url.Values{"page": {"messages"}}, payload)
	return err
}

func energyProduced(temp Temperature, resourcesBuildings ResourcesBuildings, resSettings ResourceSettings, energyTechnology int64) int64 {
	energyProduced := int64(float64(SolarPlant.Production(resourcesBuildings.SolarPlant)) * (float64(resSettings.SolarPlant) / 100))
	energyProduced += int64(float64(FusionReactor.Production(energyTechnology, resourcesBuildings.FusionReactor)) * (float64(resSettings.FusionReactor) / 100))
	energyProduced += int64(float64(SolarSatellite.Production(temp, resourcesBuildings.SolarSatellite, false)) * (float64(resSettings.SolarSatellite) / 100))
	return energyProduced
}

func energyNeeded(resourcesBuildings ResourcesBuildings, resSettings ResourceSettings) int64 {
	energyNeeded := int64(float64(MetalMine.EnergyConsumption(resourcesBuildings.MetalMine)) * (float64(resSettings.MetalMine) / 100))
	energyNeeded += int64(float64(CrystalMine.EnergyConsumption(resourcesBuildings.CrystalMine)) * (float64(resSettings.CrystalMine) / 100))
	energyNeeded += int64(float64(DeuteriumSynthesizer.EnergyConsumption(resourcesBuildings.DeuteriumSynthesizer)) * (float64(resSettings.DeuteriumSynthesizer) / 100))
	return energyNeeded
}

func productionRatio(temp Temperature, resourcesBuildings ResourcesBuildings, resSettings ResourceSettings, energyTechnology int64) float64 {
	energyProduced := energyProduced(temp, resourcesBuildings, resSettings, energyTechnology)
	energyNeeded := energyNeeded(resourcesBuildings, resSettings)
	ratio := 1.0
	if energyNeeded > energyProduced {
		ratio = float64(energyProduced) / float64(energyNeeded)
	}
	return ratio
}

func getProductions(resBuildings ResourcesBuildings, resSettings ResourceSettings, researches Researches, universeSpeed int64,
	temp Temperature, globalRatio float64) Resources {
	energyProduced := energyProduced(temp, resBuildings, resSettings, researches.EnergyTechnology)
	energyNeeded := energyNeeded(resBuildings, resSettings)
	metalSetting := float64(resSettings.MetalMine) / 100
	crystalSetting := float64(resSettings.CrystalMine) / 100
	deutSetting := float64(resSettings.DeuteriumSynthesizer) / 100
	return Resources{
		Metal:     MetalMine.Production(universeSpeed, metalSetting, globalRatio, researches.PlasmaTechnology, resBuildings.MetalMine),
		Crystal:   CrystalMine.Production(universeSpeed, crystalSetting, globalRatio, researches.PlasmaTechnology, resBuildings.CrystalMine),
		Deuterium: DeuteriumSynthesizer.Production(universeSpeed, temp.Mean(), deutSetting, globalRatio, resBuildings.DeuteriumSynthesizer) - FusionReactor.GetFuelConsumption(universeSpeed, globalRatio, resBuildings.FusionReactor),
		Energy:    energyProduced - energyNeeded,
	}
}

func (b *OGame) getResourcesProductions(planetID PlanetID) (Resources, error) {
	planet, _ := b.getPlanet(planetID)
	resBuildings, _ := b.getResourcesBuildings(planetID.Celestial())
	researches := b.getResearch()
	universeSpeed := b.serverData.Speed
	resSettings, _ := b.getResourceSettings(planetID)
	ratio := productionRatio(planet.Temperature, resBuildings, resSettings, researches.EnergyTechnology)
	productions := getProductions(resBuildings, resSettings, researches, universeSpeed, planet.Temperature, ratio)
	return productions, nil
}

func (b *OGame) getResourcesProductionsLight(resBuildings ResourcesBuildings, researches Researches,
	resSettings ResourceSettings, temp Temperature) Resources {
	universeSpeed := b.serverData.Speed
	ratio := productionRatio(temp, resBuildings, resSettings, researches.EnergyTechnology)
	productions := getProductions(resBuildings, resSettings, researches, universeSpeed, temp, ratio)
	return productions
}

func (b *OGame) getPublicIP() (string, error) {
	var res struct {
		IP string `json:"ip"`
	}
	req, err := http.NewRequest("GET", "https://jsonip.com/", nil)
	if err != nil {
		return "", err
	}
	resp, err := b.doReqWithLoginProxyTransport(req)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			b.error(err)
		}
	}()
	by, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(by, &res); err != nil {
		return "", err
	}
	return res.IP, nil
}

func (b *OGame) stateChanged(locked bool, actor string) {
	for _, clb := range b.stateChangeCallbacks {
		clb(locked, actor)
	}
}

func (b *OGame) botLock(lockedBy string) {
	b.Lock()
	if atomic.CompareAndSwapInt32(&b.lockedAtom, 0, 1) {
		b.state = lockedBy
		b.stateChanged(true, lockedBy)
	}
}

func (b *OGame) botUnlock(unlockedBy string) {
	b.Unlock()
	if atomic.CompareAndSwapInt32(&b.lockedAtom, 1, 0) {
		b.state = unlockedBy
		b.stateChanged(false, unlockedBy)
	}
}

// NewAccount response from creating a new account
type NewAccount struct {
	ID     int
	Server struct {
		Language string
		Number   int
	}
}

func (b *OGame) addAccount(number int, lang string) (NewAccount, error) {
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
	req, err := http.NewRequest("PUT", "https://"+b.lobby+".ogame.gameforge.com/api/users/me/accounts", strings.NewReader(string(jsonPayloadBytes)))
	if err != nil {
		return newAccount, err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := b.Client.Do(req)
	if err != nil {
		return newAccount, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			b.error(err)
		}
	}()
	by, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return newAccount, err
	}
	b.bytesUploaded += req.ContentLength
	b.bytesDownloaded += int64(len(by))
	if err := json.Unmarshal(by, &newAccount); err != nil {
		return newAccount, err
	}
	return newAccount, nil
}

func (b *OGame) taskRunner() {
	go func() {
		for t := range b.tasksPushCh {
			b.tasksLock.Lock()
			heap.Push(&b.tasks, t)
			b.tasksLock.Unlock()
			b.tasksPopCh <- struct{}{}
		}
	}()
	go func() {
		for range b.tasksPopCh {
			b.tasksLock.Lock()
			task := heap.Pop(&b.tasks).(*item)
			b.tasksLock.Unlock()
			close(task.canBeProcessedCh)
			<-task.isDoneCh
		}
	}()
}

func (b *OGame) getCachedCelestial(v interface{}) Celestial {
	if celestial, ok := v.(Celestial); ok {
		return celestial
	} else if planet, ok := v.(Planet); ok {
		return planet
	} else if moon, ok := v.(Moon); ok {
		return moon
	} else if celestialID, ok := v.(CelestialID); ok {
		return b.GetCachedCelestialByID(celestialID)
	} else if planetID, ok := v.(PlanetID); ok {
		return b.GetCachedCelestialByID(planetID.Celestial())
	} else if moonID, ok := v.(MoonID); ok {
		return b.GetCachedCelestialByID(moonID.Celestial())
	} else if id, ok := v.(int); ok {
		return b.GetCachedCelestialByID(CelestialID(id))
	} else if id, ok := v.(int32); ok {
		return b.GetCachedCelestialByID(CelestialID(id))
	} else if id, ok := v.(int64); ok {
		return b.GetCachedCelestialByID(CelestialID(id))
	} else if id, ok := v.(float32); ok {
		return b.GetCachedCelestialByID(CelestialID(id))
	} else if id, ok := v.(float64); ok {
		return b.GetCachedCelestialByID(CelestialID(id))
	} else if id, ok := v.(lua.LNumber); ok {
		return b.GetCachedCelestialByID(CelestialID(id))
	} else if coord, ok := v.(Coordinate); ok {
		return b.GetCachedCelestialByCoord(coord)
	} else if coordStr, ok := v.(string); ok {
		coord, err := ParseCoord(coordStr)
		if err != nil {
			return nil
		}
		return b.GetCachedCelestialByCoord(coord)
	}
	return nil
}

// GetCachedCelestialByID return celestial from cached value
func (b *OGame) GetCachedCelestialByID(celestialID CelestialID) Celestial {
	for _, p := range b.Planets {
		if p.ID.Celestial() == celestialID {
			return p
		}
		if p.Moon != nil && p.Moon.ID.Celestial() == celestialID {
			return p.Moon
		}
	}
	return nil
}

// GetCachedCelestialByCoord return celestial from cached value
func (b *OGame) GetCachedCelestialByCoord(coord Coordinate) Celestial {
	for _, p := range b.Planets {
		if p.GetCoordinate().Equal(coord) {
			return p
		}
		if p.Moon != nil && p.Moon.GetCoordinate().Equal(coord) {
			return p.Moon
		}
	}
	return nil
}

func (b *OGame) fakeCall(name string, delay int) {
	fmt.Println("before", name)
	time.Sleep(time.Duration(delay) * time.Millisecond)
	fmt.Println("after", name)
}

// FakeCall used for debugging
func (b *OGame) FakeCall(priority int, name string, delay int) {
	b.WithPriority(priority).FakeCall(name, delay)
}

func (b *OGame) getCachedMoons() []Moon {
	var moons []Moon
	for _, p := range b.Planets {
		if p.Moon != nil {
			moons = append(moons, *p.Moon)
		}
	}
	return moons
}

func (b *OGame) getCachedCelestials() []Celestial {
	celestials := make([]Celestial, 0)
	for _, p := range b.Planets {
		celestials = append(celestials, p)
		if p.Moon != nil {
			celestials = append(celestials, p.Moon)
		}
	}
	return celestials
}

func (b *OGame) withPriority(priority int) *Prioritize {
	canBeProcessedCh := make(chan struct{})
	taskIsDoneCh := make(chan struct{})
	task := new(item)
	task.priority = priority
	task.canBeProcessedCh = canBeProcessedCh
	task.isDoneCh = taskIsDoneCh
	b.tasksPushCh <- task
	<-canBeProcessedCh
	return &Prioritize{bot: b, taskIsDoneCh: taskIsDoneCh}
}

// TasksOverview overview of tasks in heap
type TasksOverview struct {
	Low       int64
	Normal    int64
	Important int64
	Critical  int64
	Total     int64
}

func (b *OGame) getTasks() (out TasksOverview) {
	b.tasksLock.Lock()
	out.Total = int64(b.tasks.Len())
	for _, item := range b.tasks {
		switch item.priority {
		case Low:
			out.Low++
		case Normal:
			out.Normal++
		case Important:
			out.Important++
		case Critical:
			out.Critical++
		}
	}
	b.tasksLock.Unlock()
	return
}

// Public interface -----------------------------------------------------------

// Enable enables communications with OGame Server
func (b *OGame) Enable() {
	b.enable()
}

// Disable disables communications with OGame Server
func (b *OGame) Disable() {
	b.disable()
}

// IsEnabled returns true if the bot is enabled, otherwise false
func (b *OGame) IsEnabled() bool {
	return b.isEnabled()
}

// IsLoggedIn returns true if the bot is currently logged-in, otherwise false
func (b *OGame) IsLoggedIn() bool {
	return atomic.LoadInt32(&b.isLoggedInAtom) == 1
}

// IsConnected returns true if the bot is currently connected (communication between the bot and OGame is possible), otherwise false
func (b *OGame) IsConnected() bool {
	return atomic.LoadInt32(&b.isConnectedAtom) == 1
}

// GetClient get the http client used by the bot
func (b *OGame) GetClient() *OGameClient {
	return b.Client
}

// GetPublicIP get the public IP used by the bot
func (b *OGame) GetPublicIP() (string, error) {
	return b.getPublicIP()
}

// OnStateChange register a callback that is notified when the bot state changes
func (b *OGame) OnStateChange(clb func(locked bool, actor string)) {
	b.stateChangeCallbacks = append(b.stateChangeCallbacks, clb)
}

// GetState returns the current bot state
func (b *OGame) GetState() (bool, string) {
	return atomic.LoadInt32(&b.lockedAtom) == 1, b.state
}

// IsLocked returns either or not the bot is currently locked
func (b *OGame) IsLocked() bool {
	return atomic.LoadInt32(&b.lockedAtom) == 1
}

// GetSession get ogame session
func (b *OGame) GetSession() string {
	return b.ogameSession
}

// AddAccount add a new account (server) to your list of accounts
func (b *OGame) AddAccount(number int, lang string) (NewAccount, error) {
	return b.addAccount(number, lang)
}

// WithPriority ...
func (b *OGame) WithPriority(priority int) *Prioritize {
	return b.withPriority(priority)
}

// Begin start a transaction. Once this function is called, "Done" must be called to release the lock.
func (b *OGame) Begin() *Prioritize {
	return b.WithPriority(Normal).Begin()
}

// BeginNamed begins a new transaction with a name. "Done" must be called to release the lock.
func (b *OGame) BeginNamed(name string) *Prioritize {
	return b.WithPriority(Normal).BeginNamed(name)
}

// Tx locks the bot during the transaction and ensure the lock is released afterward
func (b *OGame) Tx(clb func(tx *Prioritize) error) error {
	return b.WithPriority(Normal).Tx(clb)
}

// GetServer get ogame server information that the bot is connected to
func (b *OGame) GetServer() Server {
	return b.server
}

// GetServerData get ogame server data information that the bot is connected to
func (b *OGame) GetServerData() ServerData {
	return b.serverData
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
	b.Client.UserAgent = newUserAgent
}

// LoginWithExistingCookies to ogame server reusing existing cookies
func (b *OGame) LoginWithExistingCookies() error {
	return b.WithPriority(Normal).LoginWithExistingCookies()
}

// Login to ogame server
// Can fails with BadCredentialsError
func (b *OGame) Login() error {
	return b.WithPriority(Normal).Login()
}

// Logout the bot from ogame server
func (b *OGame) Logout() { b.WithPriority(Normal).Logout() }

// BytesDownloaded returns the amount of bytes downloaded
func (b *OGame) BytesDownloaded() int64 {
	return b.bytesDownloaded
}

// BytesUploaded returns the amount of bytes uploaded
func (b *OGame) BytesUploaded() int64 {
	return b.bytesUploaded
}

// GetUniverseName get the name of the universe the bot is playing into
func (b *OGame) GetUniverseName() string {
	return b.Universe
}

// GetUsername get the username that was used to login on ogame server
func (b *OGame) GetUsername() string {
	return b.Username
}

// GetResearchSpeed gets the research speed
func (b *OGame) GetResearchSpeed() int64 {
	return b.serverData.ResearchDurationDivisor
}

// GetNbSystems gets the number of systems
func (b *OGame) GetNbSystems() int64 {
	return b.serverData.Systems
}

// GetUniverseSpeed shortcut to get ogame universe speed
func (b *OGame) GetUniverseSpeed() int64 {
	return b.getUniverseSpeed()
}

// GetUniverseSpeedFleet shortcut to get ogame universe speed fleet
func (b *OGame) GetUniverseSpeedFleet() int64 {
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

// FleetDeutSaveFactor returns the fleet deut save factor
func (b *OGame) FleetDeutSaveFactor() float64 {
	return b.serverData.GlobalDeuteriumSaveFactor
}

// GetAlliancePageContent gets the html for a specific alliance page
func (b *OGame) GetAlliancePageContent(vals url.Values) []byte {
	return b.WithPriority(Normal).GetPageContent(vals)
}

// GetPageContent gets the html for a specific ogame page
func (b *OGame) GetPageContent(vals url.Values) []byte {
	return b.WithPriority(Normal).GetPageContent(vals)
}

// PostPageContent make a post request to ogame server
// This is useful when simulating a web browser
func (b *OGame) PostPageContent(vals, payload url.Values) []byte {
	return b.WithPriority(Normal).PostPageContent(vals, payload)
}

// IsUnderAttack returns true if the user is under attack, false otherwise
func (b *OGame) IsUnderAttack() (bool, error) {
	return b.WithPriority(Normal).IsUnderAttack()
}

// GetCachedPlayer returns cached player infos
func (b *OGame) GetCachedPlayer() UserInfos {
	return b.Player
}

// GetCachedPreferences returns cached preferences
func (b *OGame) GetCachedPreferences() Preferences {
	return b.CachedPreferences
}

// IsVacationModeEnabled returns either or not the bot is in vacation mode
func (b *OGame) IsVacationModeEnabled() bool {
	return b.isVacationModeEnabled
}

// GetPlanets returns the user planets
func (b *OGame) GetPlanets() []Planet {
	return b.WithPriority(Normal).GetPlanets()
}

// GetCachedPlanets return planets from cached value
func (b *OGame) GetCachedPlanets() []Planet {
	return b.Planets
}

// GetCachedMoons return moons from cached value
func (b *OGame) GetCachedMoons() []Moon {
	return b.getCachedMoons()
}

// GetCachedCelestials get all cached celestials
func (b *OGame) GetCachedCelestials() []Celestial {
	return b.getCachedCelestials()
}

// GetCachedCelestial return celestial from cached value
func (b *OGame) GetCachedCelestial(v interface{}) Celestial {
	return b.getCachedCelestial(v)
}

// GetPlanet gets infos for planetID
// Fails if planetID is invalid
func (b *OGame) GetPlanet(v interface{}) (Planet, error) {
	return b.WithPriority(Normal).GetPlanet(v)
}

// GetMoons returns the user moons
func (b *OGame) GetMoons() []Moon {
	return b.WithPriority(Normal).GetMoons()
}

// GetMoon gets infos for moonID
func (b *OGame) GetMoon(v interface{}) (Moon, error) {
	return b.WithPriority(Normal).GetMoon(v)
}

// GetCelestials get the player's planets & moons
func (b *OGame) GetCelestials() ([]Celestial, error) {
	return b.WithPriority(Normal).GetCelestials()
}

// Abandon a planet
func (b *OGame) Abandon(v interface{}) error {
	return b.WithPriority(Normal).Abandon(v)
}

// GetCelestial get the player's planet/moon using the coordinate
func (b *OGame) GetCelestial(v interface{}) (Celestial, error) {
	return b.WithPriority(Normal).GetCelestial(v)
}

// ServerVersion returns OGame version
func (b *OGame) ServerVersion() string {
	return b.serverData.Version
}

// ServerTime returns server time
// Timezone is OGT (OGame Time zone)
func (b *OGame) ServerTime() time.Time {
	return b.WithPriority(Normal).ServerTime()
}

// Location returns bot Time zone.
func (b *OGame) Location() *time.Location {
	return b.location
}

// GetUserInfos gets the user information
func (b *OGame) GetUserInfos() UserInfos {
	return b.WithPriority(Normal).GetUserInfos()
}

// SendMessage sends a message to playerID
func (b *OGame) SendMessage(playerID int64, message string) error {
	return b.WithPriority(Normal).SendMessage(playerID, message)
}

// SendMessageAlliance sends a message to associationID
func (b *OGame) SendMessageAlliance(associationID int64, message string) error {
	return b.WithPriority(Normal).SendMessageAlliance(associationID, message)
}

// GetFleets get the player's own fleets activities
func (b *OGame) GetFleets(opts ...Option) ([]Fleet, Slots) {
	return b.WithPriority(Normal).GetFleets(opts...)
}

// GetFleetsFromEventList get the player's own fleets activities
func (b *OGame) GetFleetsFromEventList() []Fleet {
	return b.WithPriority(Normal).GetFleetsFromEventList()
}

// CancelFleet cancel a fleet
func (b *OGame) CancelFleet(fleetID FleetID) error {
	return b.WithPriority(Normal).CancelFleet(fleetID)
}

// GetAttacks get enemy fleets attacking you
func (b *OGame) GetAttacks() ([]AttackEvent, error) {
	return b.WithPriority(Normal).GetAttacks()
}

// GetAttacksUsing get enemy fleets attacking you using a specific celestial to make the check
func (b *OGame) GetAttacksUsing(celestialID CelestialID) ([]AttackEvent, error) {
	return b.WithPriority(Normal).GetAttacksUsing(celestialID)
}

// GalaxyInfos get information of all planets and moons of a solar system
func (b *OGame) GalaxyInfos(galaxy, system int64, options ...Option) (SystemInfos, error) {
	return b.WithPriority(Normal).GalaxyInfos(galaxy, system, options...)
}

// GetResourceSettings gets the resources settings for specified planetID
func (b *OGame) GetResourceSettings(planetID PlanetID) (ResourceSettings, error) {
	return b.WithPriority(Normal).GetResourceSettings(planetID)
}

// SetResourceSettings set the resources settings on a planet
func (b *OGame) SetResourceSettings(planetID PlanetID, settings ResourceSettings) error {
	return b.WithPriority(Normal).SetResourceSettings(planetID, settings)
}

// GetResourcesBuildings gets the resources buildings levels
func (b *OGame) GetResourcesBuildings(celestialID CelestialID) (ResourcesBuildings, error) {
	return b.WithPriority(Normal).GetResourcesBuildings(celestialID)
}

// GetDefense gets all the defenses units information of a planet
// Fails if planetID is invalid
func (b *OGame) GetDefense(celestialID CelestialID) (DefensesInfos, error) {
	return b.WithPriority(Normal).GetDefense(celestialID)
}

// GetShips gets all ships units information of a planet
func (b *OGame) GetShips(celestialID CelestialID) (ShipsInfos, error) {
	return b.WithPriority(Normal).GetShips(celestialID)
}

// GetFacilities gets all facilities information of a planet
func (b *OGame) GetFacilities(celestialID CelestialID) (Facilities, error) {
	return b.WithPriority(Normal).GetFacilities(celestialID)
}

// GetProduction get what is in the production queue.
// (ships & defense being built)
func (b *OGame) GetProduction(celestialID CelestialID) ([]Quantifiable, int64, error) {
	return b.WithPriority(Normal).GetProduction(celestialID)
}

// GetCachedResearch returns cached researches
func (b *OGame) GetCachedResearch() Researches {
	return b.WithPriority(Normal).GetCachedResearch()
}

// GetResearch gets the player researches information
func (b *OGame) GetResearch() Researches {
	return b.WithPriority(Normal).GetResearch()
}

// GetSlots gets the player current and total slots information
func (b *OGame) GetSlots() Slots {
	return b.WithPriority(Normal).GetSlots()
}

// Build builds any ogame objects (building, technology, ship, defence)
func (b *OGame) Build(celestialID CelestialID, id ID, nbr int64) error {
	return b.WithPriority(Normal).Build(celestialID, id, nbr)
}

// TearDown tears down any ogame building
func (b *OGame) TearDown(celestialID CelestialID, id ID) error {
	return b.WithPriority(Normal).TearDown(celestialID, id)
}

// BuildCancelable builds any cancelable ogame objects (building, technology)
func (b *OGame) BuildCancelable(celestialID CelestialID, id ID) error {
	return b.WithPriority(Normal).BuildCancelable(celestialID, id)
}

// BuildProduction builds any line production ogame objects (ship, defence)
func (b *OGame) BuildProduction(celestialID CelestialID, id ID, nbr int64) error {
	return b.WithPriority(Normal).BuildProduction(celestialID, id, nbr)
}

// BuildBuilding ensure what is being built is a building
func (b *OGame) BuildBuilding(celestialID CelestialID, buildingID ID) error {
	return b.WithPriority(Normal).BuildBuilding(celestialID, buildingID)
}

// BuildDefense builds a defense unit
func (b *OGame) BuildDefense(celestialID CelestialID, defenseID ID, nbr int64) error {
	return b.WithPriority(Normal).BuildDefense(celestialID, defenseID, nbr)
}

// BuildShips builds a ship unit
func (b *OGame) BuildShips(celestialID CelestialID, shipID ID, nbr int64) error {
	return b.WithPriority(Normal).BuildShips(celestialID, shipID, nbr)
}

// ConstructionsBeingBuilt returns the building & research being built, and the time remaining (secs)
func (b *OGame) ConstructionsBeingBuilt(celestialID CelestialID) (ID, int64, ID, int64) {
	return b.WithPriority(Normal).ConstructionsBeingBuilt(celestialID)
}

// CancelBuilding cancel the construction of a building on a specified planet
func (b *OGame) CancelBuilding(celestialID CelestialID) error {
	return b.WithPriority(Normal).CancelBuilding(celestialID)
}

// CancelResearch cancel the research
func (b *OGame) CancelResearch(celestialID CelestialID) error {
	return b.WithPriority(Normal).CancelResearch(celestialID)
}

// BuildTechnology ensure that we're trying to build a technology
func (b *OGame) BuildTechnology(celestialID CelestialID, technologyID ID) error {
	return b.WithPriority(Normal).BuildTechnology(celestialID, technologyID)
}

// GetResources gets user resources
func (b *OGame) GetResources(celestialID CelestialID) (Resources, error) {
	return b.WithPriority(Normal).GetResources(celestialID)
}

// GetResourcesDetails gets user resources
func (b *OGame) GetResourcesDetails(celestialID CelestialID) (ResourcesDetails, error) {
	return b.WithPriority(Normal).GetResourcesDetails(celestialID)
}

// SendFleet sends a fleet
func (b *OGame) SendFleet(celestialID CelestialID, ships []Quantifiable, speed Speed, where Coordinate,
	mission MissionID, resources Resources, expeditiontime, unionID int64) (Fleet, error) {
	return b.WithPriority(Normal).SendFleet(celestialID, ships, speed, where, mission, resources, expeditiontime, unionID)
}

// EnsureFleet either sends all the requested ships or fail
func (b *OGame) EnsureFleet(celestialID CelestialID, ships []Quantifiable, speed Speed, where Coordinate,
	mission MissionID, resources Resources, expeditiontime, unionID int64) (Fleet, error) {
	return b.WithPriority(Normal).EnsureFleet(celestialID, ships, speed, where, mission, resources, expeditiontime, unionID)
}

// SendIPM sends IPM
func (b *OGame) SendIPM(planetID PlanetID, coord Coordinate, nbr int64, priority ID) (int64, error) {
	return b.WithPriority(Normal).SendIPM(planetID, coord, nbr, priority)
}

// GetCombatReportSummaryFor gets the latest combat report for a given coordinate
func (b *OGame) GetCombatReportSummaryFor(coord Coordinate) (CombatReportSummary, error) {
	return b.WithPriority(Normal).GetCombatReportSummaryFor(coord)
}

// GetEspionageReportFor gets the latest espionage report for a given coordinate
func (b *OGame) GetEspionageReportFor(coord Coordinate) (EspionageReport, error) {
	return b.WithPriority(Normal).GetEspionageReportFor(coord)
}

// GetEspionageReportMessages gets the summary of each espionage reports
func (b *OGame) GetEspionageReportMessages() ([]EspionageReportSummary, error) {
	return b.WithPriority(Normal).GetEspionageReportMessages()
}

// GetEspionageReport gets a detailed espionage report
func (b *OGame) GetEspionageReport(msgID int64) (EspionageReport, error) {
	return b.WithPriority(Normal).GetEspionageReport(msgID)
}

// DeleteMessage deletes a message from the mail box
func (b *OGame) DeleteMessage(msgID int64) error {
	return b.WithPriority(Normal).DeleteMessage(msgID)
}

// DeleteAllMessagesFromTab deletes all messages from a tab in the mail box
func (b *OGame) DeleteAllMessagesFromTab(tabID int64) error {
	return b.WithPriority(Normal).DeleteAllMessagesFromTab(tabID)
}

// GetResourcesProductions gets the planet resources production
func (b *OGame) GetResourcesProductions(planetID PlanetID) (Resources, error) {
	return b.WithPriority(Normal).GetResourcesProductions(planetID)
}

// GetResourcesProductionsLight gets the planet resources production
func (b *OGame) GetResourcesProductionsLight(resBuildings ResourcesBuildings, researches Researches,
	resSettings ResourceSettings, temp Temperature) Resources {
	return b.WithPriority(Normal).GetResourcesProductionsLight(resBuildings, researches, resSettings, temp)
}

// FlightTime calculate flight time and fuel needed
func (b *OGame) FlightTime(origin, destination Coordinate, speed Speed, ships ShipsInfos) (secs, fuel int64) {
	return b.WithPriority(Normal).FlightTime(origin, destination, speed, ships)
}

// Distance return distance between two coordinates
func (b *OGame) Distance(origin, destination Coordinate) int64 {
	return Distance(origin, destination, b.serverData.Galaxies, b.serverData.Systems, b.serverData.DonutGalaxy, b.serverData.DonutSystem)
}

// RegisterChatCallback register a callback that is called when chat messages are received
func (b *OGame) RegisterChatCallback(fn func(msg ChatMsg)) {
	b.chatCallbacks = append(b.chatCallbacks, fn)
}

// RegisterAuctioneerCallback register a callback that is called when auctioneer packets are received
func (b *OGame) RegisterAuctioneerCallback(fn func(packet []byte)) {
	b.auctioneerCallbacks = append(b.auctioneerCallbacks, fn)
}

// RegisterHTMLInterceptor ...
func (b *OGame) RegisterHTMLInterceptor(fn func(method, url string, params, payload url.Values, pageHTML []byte)) {
	b.interceptorCallbacks = append(b.interceptorCallbacks, fn)
}

// Phalanx scan a coordinate from a moon to get fleets information
// IMPORTANT: My account was instantly banned when I scanned an invalid coordinate.
// IMPORTANT: This function DOES validate that the coordinate is a valid planet in range of phalanx
// 			  and that you have enough deuterium.
func (b *OGame) Phalanx(moonID MoonID, coord Coordinate) ([]Fleet, error) {
	return b.WithPriority(Normal).Phalanx(moonID, coord)
}

// UnsafePhalanx same as Phalanx but does not perform any input validation.
func (b *OGame) UnsafePhalanx(moonID MoonID, coord Coordinate) ([]Fleet, error) {
	return b.WithPriority(Normal).UnsafePhalanx(moonID, coord)
}

// JumpGateDestinations returns available destinations for jump gate.
func (b *OGame) JumpGateDestinations(origin MoonID) (moonIDs []MoonID, rechargeCountdown int64, err error) {
	return b.WithPriority(Normal).JumpGateDestinations(origin)
}

// JumpGate sends ships through a jump gate.
func (b *OGame) JumpGate(origin, dest MoonID, ships ShipsInfos) (success bool, rechargeCountdown int64, err error) {
	return b.WithPriority(Normal).JumpGate(origin, dest, ships)
}

// BuyOfferOfTheDay buys the offer of the day.
func (b *OGame) BuyOfferOfTheDay() error {
	return b.WithPriority(Normal).BuyOfferOfTheDay()
}

// CreateUnion creates a union
func (b *OGame) CreateUnion(fleet Fleet, users []UserInfos) (int64, error) {
	return b.WithPriority(Normal).CreateUnion(fleet, users)
}

// HeadersForPage gets the headers for a specific ogame page
func (b *OGame) HeadersForPage(url string) (http.Header, error) {
	return b.WithPriority(Normal).HeadersForPage(url)
}

// GetEmpire retrieves JSON from Empire page (Commander only).
func (b *OGame) GetEmpire(nbr int64) (interface{}, error) {
	return b.WithPriority(Normal).GetEmpire(nbr)
}

// CharacterClass returns the bot character class
func (b *OGame) CharacterClass() CharacterClass {
	return b.characterClass
}

// GetAuction ...
func (b *OGame) GetAuction() (Auction, error) {
	return b.WithPriority(Normal).GetAuction()
}

// DoAuction ...
func (b *OGame) DoAuction(bid map[CelestialID]Resources) error {
	return b.WithPriority(Normal).DoAuction(bid)
}

// Highscore ...
func (b *OGame) Highscore(category, typ, page int64) (Highscore, error) {
	return b.WithPriority(Normal).Highscore(category, typ, page)
}

// GetAllResources gets the resources of all planets and moons
func (b *OGame) GetAllResources() (map[CelestialID]Resources, error) {
	return b.WithPriority(Normal).GetAllResources()
}

// GetTasks return how many tasks are queued in the heap.
func (b *OGame) GetTasks() TasksOverview {
	return b.getTasks()
}

// GetDMCosts returns fast build with DM information
func (b *OGame) GetDMCosts(celestialID CelestialID) (DMCosts, error) {
	return b.WithPriority(Normal).GetDMCosts(celestialID)
}

// UseDM use dark matter to fast build
func (b *OGame) UseDM(typ string, celestialID CelestialID) error {
	return b.WithPriority(Normal).UseDM(typ, celestialID)
}

// GetItems get all items information
func (b *OGame) GetItems(celestialID CelestialID) ([]Item, error) {
	return b.WithPriority(Normal).GetItems(celestialID)
}

// ActivateItem activate an item
func (b *OGame) ActivateItem(ref string, celestialID CelestialID) error {
	return b.WithPriority(Normal).ActivateItem(ref, celestialID)
}
