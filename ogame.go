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
	"net/http/cookiejar"
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
	"github.com/pkg/errors"
	"github.com/yuin/gopher-lua"
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
	ogameSession          string
	sessionChatCounter    int
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
	loginProxyTransport   *http.Transport
	bytesUploaded         int64
	bytesDownloaded       int64
}

// Preferences ...
type Preferences struct {
	SpioAnz                      int
	DisableChatBar               bool // no-mobile
	DisableOutlawWarning         bool
	MobileVersion                bool
	ShowOldDropDowns             bool
	ActivateAutofocus            bool
	EventsShow                   int // Hide: 1, Above the content: 2, Below the content: 3
	SortSetting                  int // Order of emergence: 0, Coordinates: 1, Alphabet: 2, Size: 3, Used fields: 4
	SortOrder                    int // Up: 0, Down: 1
	ShowDetailOverlay            bool
	AnimatedSliders              bool // no-mobile
	AnimatedOverview             bool // no-mobile
	PopupsNotices                bool // no-mobile
	PopopsCombatreport           bool // no-mobile
	SpioReportPictures           bool
	MsgResultsPerPage            int // 10, 25, 50
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
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) " +
	"AppleWebKit/537.36 (KHTML, like Gecko) " +
	"Chrome/51.0.2704.103 " +
	"Safari/537.36"

// CelestialID represent either a PlanetID or a MoonID
type CelestialID int

// Params parameters for more fine-grained initialization
type Params struct {
	Universe       string
	Username       string
	Password       string
	Lang           string
	AutoLogin      bool
	Proxy          string
	ProxyUsername  string
	ProxyPassword  string
	Socks5Address  string
	Socks5Username string
	Socks5Password string
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
		if err := b.SetProxy(params.Proxy, params.ProxyUsername, params.ProxyPassword); err != nil {
			return nil, err
		}
	}
	if params.Socks5Address != "" {
		if err := b.SetSocks5Proxy(params.Socks5Address, params.Socks5Username, params.Socks5Password); err != nil {
			return nil, err
		}
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
	b.loginWrapper = DefaultLoginWrapper
	b.Enable()
	b.quiet = false
	b.logger = log.New(os.Stdout, "", 0)

	b.Universe = universe
	b.SetOGameCredentials(username, password)
	b.language = lang

	jar, _ := cookiejar.New(nil)
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

func getPhpSessionID(b *OGame, username, password string) (string, error) {
	payload := url.Values{
		"kid":                   {""},
		"language":              {"en"},
		"autologin":             {"false"},
		"credentials[email]":    {username},
		"credentials[password]": {password},
	}
	req, err := http.NewRequest("POST", "https://lobby.ogame.gameforge.com/api/users", strings.NewReader(payload.Encode()))
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
		Number   int
	}
	ID         int
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

func getUserAccounts(b *OGame, phpSessionID string) ([]account, error) {
	var userAccounts []account
	req, err := http.NewRequest("GET", "https://lobby.ogame.gameforge.com/api/users/me/accounts", nil)
	if err != nil {
		return userAccounts, err
	}
	req.AddCookie(&http.Cookie{Name: phpSessionIDCookieName, Value: phpSessionID})
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
	req, err := http.NewRequest("GET", "https://lobby.ogame.gameforge.com/api/servers", nil)
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

func getLoginLink(b *OGame, userAccount account, phpSessionID string) (string, error) {
	ogURL := fmt.Sprintf("https://lobby.ogame.gameforge.com/api/users/me/loginLink?id=%d&server[language]=%s&server[number]=%d",
		userAccount.ID, userAccount.Server.Language, userAccount.Server.Number)
	req, err := http.NewRequest("GET", ogURL, nil)
	if err != nil {
		return "", err
	}
	req.AddCookie(&http.Cookie{Name: phpSessionIDCookieName, Value: phpSessionID})
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
	Name                          string  // Europa
	Number                        int     // 157
	Language                      string  // ru
	Timezone                      string  // Europe/Moscow
	TimezoneOffset                string  // +03:00
	Domain                        string  // s157-ru.ogame.gameforge.com
	Version                       string  // 6.8.8-pl2
	Speed                         int     // 6
	SpeedFleet                    int     // 6
	Galaxies                      int     // 4
	Systems                       int     // 499
	ACS                           bool    // 1
	RapidFire                     bool    // 1
	DefToTF                       bool    // 0
	DebrisFactor                  float64 // 0.5
	DebrisFactorDef               float64 // 0
	RepairFactor                  float64 // 0.7
	NewbieProtectionLimit         int     // 500000
	NewbieProtectionHigh          int     // 50000
	TopScore                      int     // 60259362
	BonusFields                   int     // 30
	DonutGalaxy                   bool    // 1
	DonutSystem                   bool    // 1
	WfEnabled                     bool    // 1 (WreckField)
	WfMinimumRessLost             int     // 150000
	WfMinimumLossPercentage       int     // 5
	WfBasicPercentageRepairable   int     // 45
	GlobalDeuteriumSaveFactor     float64 // 0.5
	Bashlimit                     int     // 0
	ProbeCargo                    bool    // 5
	ResearchDurationDivisor       int     // 2
	DarkMatterNewAcount           int     // 8000
	CargoHyperspaceTechMultiplier int     // 5
}

// gets the server data from xml api
func (b *OGame) getServerData() (ServerData, error) {
	var serverData ServerData
	req, err := http.NewRequest("GET", "https://s"+strconv.Itoa(b.server.Number)+"-"+b.server.Language+".ogame.gameforge.com/api/serverData.xml", nil)
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

func (b *OGame) login() error {
	jar, _ := cookiejar.New(nil)
	b.Client.Jar = jar

	b.debug("get session")
	phpSessionID, err := getPhpSessionID(b, b.Username, b.password)
	if err != nil {
		return err
	}
	b.debug("get user accounts")
	accounts, err := getUserAccounts(b, phpSessionID)
	if err != nil {
		return err
	}
	b.debug("get servers")
	servers, err := getServers(b)
	if err != nil {
		return err
	}
	b.debug("find account & server for universe")
	userAccount, server, err := findAccountByName(b.Universe, b.language, accounts, servers)
	if err != nil {
		return err
	}
	if userAccount.Blocked {
		return ErrAccountBlocked
	}
	b.debug("Players online: " + strconv.Itoa(server.PlayersOnline) + ", Players: " + strconv.Itoa(server.PlayerCount))
	b.server = server
	b.language = userAccount.Server.Language
	b.debug("get login link")
	loginLink, err := getLoginLink(b, userAccount, phpSessionID)
	if err != nil {
		return err
	}

	r := regexp.MustCompile(`(https://.+\.ogame\.gameforge\.com)/game`)
	res := r.FindStringSubmatch(loginLink)
	if len(res) != 2 {
		return errors.New("failed to get server url")
	}
	b.serverURL = res[1]

	pageHTML, err := execLoginLink(b, loginLink)
	b.debug("extract information from html")
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return err
	}
	b.ogameSession = ExtractOGameSessionFromDoc(doc)
	if b.ogameSession == "" {
		return errors.New("bad credentials")
	}

	// Get server data
	start := time.Now()
	serverData, err := b.getServerData()
	if err != nil {
		return err
	}
	b.serverData = serverData
	b.debug("get server data", time.Since(start))

	atomic.StoreInt32(&b.isLoggedInAtom, 1) // At this point, we are logged in
	atomic.StoreInt32(&b.isConnectedAtom, 1)
	b.sessionChatCounter = 1

	serverTime, _ := extractServerTime(pageHTML)
	b.location = serverTime.Location()

	b.cacheFullPageInfo("overview", pageHTML)

	for _, fn := range b.interceptorCallbacks {
		fn("GET", loginLink, nil, nil, pageHTML)
	}

	_, _ = b.getPageContent(url.Values{"page": {"preferences"}}) // Will update preferences cached values

	// Extract chat host and port
	m := regexp.MustCompile(`var nodeUrl="https:\\/\\/([^:]+):(\d+)\\/socket.io\\/socket.io.js";`).FindSubmatch(pageHTML)
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
	b.Planets = ExtractPlanets(pageHTML, b)
	b.isVacationModeEnabled = ExtractIsInVacation(pageHTML)
	b.ajaxChatToken, _ = ExtractAjaxChatToken(pageHTML)
	if page == "overview" {
		b.Player, _ = ExtractUserInfos(pageHTML, b.language)
	} else if page == "preferences" {
		b.CachedPreferences = ExtractPreferences(pageHTML)
	}
}

// DefaultLoginWrapper ...
var DefaultLoginWrapper = func(loginFn func() error) error {
	return loginFn()
}

func (b *OGame) wrapLogin() error {
	return b.loginWrapper(b.login)
}

// SetOGameCredentials sets ogame credentials for the bot
func (b *OGame) SetOGameCredentials(username, password string) {
	b.Username = username
	b.password = password
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

// SetLoginProxy set the proxy to use for login requests
func (b *OGame) SetLoginProxy(proxy, username, password string) error {
	transport, err := getProxyTransport(proxy, username, password)
	if err != nil {
		return err
	}
	b.loginProxyTransport = transport
	return nil
}

// SetProxy this will change the bot http transport object
func (b *OGame) SetProxy(proxy, username, password string) error {
	t, err := getProxyTransport(proxy, username, password)
	if err != nil {
		return err
	}
	b.Client.Transport = t
	return nil
}

// SetSocks5Proxy this will change the bot http transport object
func (b *OGame) SetSocks5Proxy(socks5Address, socks5Username, socks5Password string) error {
	var auth *proxy.Auth
	if socks5Username != "" || socks5Password != "" {
		auth = &proxy.Auth{User: socks5Username, Password: socks5Password}
	}
	dialer, err := proxy.SOCKS5("tcp", socks5Address, auth, proxy.Direct)
	if err != nil {
		return err
	}
	b.Client.Transport = &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		},
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
			authMsg := `5:` + strconv.Itoa(b.sessionChatCounter) + `+:/chat:{"name":"authorize","args":["` + b.ogameSession + `"]}`
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
	_, _ = b.getPageContent(url.Values{"page": {"logout"}})
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
	return len(regexp.MustCompile(`<meta name="ogame-session" content="\w+"/>`).FindSubmatch(pageHTML)) == 1
}

// IsKnowFullPage ...
func IsKnowFullPage(vals url.Values) bool {
	page := vals.Get("page")
	return page == "overview" ||
		page == "resources" ||
		page == "station" ||
		page == "traderOverview" ||
		page == "research" ||
		page == "shipyard" ||
		page == "defense" ||
		page == "fleet1" ||
		page == "galaxy" ||
		page == "alliance" ||
		page == "premium" ||
		page == "shop" ||
		page == "rewards" ||
		page == "resourceSettings" ||
		page == "movement" ||
		page == "highscore" ||
		page == "buddies" ||
		page == "preferences" ||
		page == "messages" ||
		page == "chat"
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
		page == "federationlayer" ||
		page == "unionchange" ||
		page == "changenick" ||
		page == "planetlayer" ||
		page == "traderlayer" ||
		page == "planetRename" ||
		page == "rightmenu" ||
		page == "allianceOverview" ||
		page == "support" ||
		page == "buffActivation" ||
		ajax == "1"
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

func (b *OGame) getPageContent(vals url.Values) ([]byte, error) {
	if err := b.preRequestChecks(); err != nil {
		return []byte{}, err
	}

	finalURL := b.serverURL + "/game/index.php?" + vals.Encode()
	page := vals.Get("page")
	var pageHTMLBytes []byte

	if err := b.withRetry(func() (err error) {
		pageHTMLBytes, err = b.execRequest("GET", finalURL, nil, vals)
		if err != nil {
			return err
		}

		if (page != "logout" && (IsKnowFullPage(vals) || page == "") && !IsAjaxPage(vals) && !isLogged(pageHTMLBytes)) ||
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

	go func() {
		for _, fn := range b.interceptorCallbacks {
			fn("GET", finalURL, vals, nil, pageHTMLBytes)
		}
	}()

	return pageHTMLBytes, nil
}

func (b *OGame) postPageContent(vals, payload url.Values) ([]byte, error) {
	if err := b.preRequestChecks(); err != nil {
		return []byte{}, err
	}

	if vals.Get("page") == "ajaxChat" && payload.Get("mode") == "1" {
		payload.Set("token", b.ajaxChatToken)
	}

	finalURL := b.serverURL + "/game/index.php?" + vals.Encode()
	page := vals.Get("page")
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
		b.CachedPreferences = ExtractPreferences(pageHTMLBytes)
	} else if page == "ajaxChat" && (payload.Get("mode") == "1" || payload.Get("mode") == "3") {
		var res ChatPostResp
		if err := json.Unmarshal(pageHTMLBytes, &res); err != nil {
			return []byte{}, err
		}
		b.ajaxChatToken = res.NewToken
	}

	go func() {
		for _, fn := range b.interceptorCallbacks {
			fn("POST", finalURL, vals, payload, pageHTMLBytes)
		}
	}()

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
			return ErrFailedExecuteCallback
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

func (b *OGame) getUniverseSpeed() int {
	return b.serverData.Speed
}

func (b *OGame) getUniverseSpeedFleet() int {
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

func (b *OGame) getPlanets() []Planet {
	pageHTML, _ := b.getPageContent(url.Values{"page": {"overview"}})
	return ExtractPlanets(pageHTML, b)
}

func (b *OGame) getPlanet(v interface{}) (Planet, error) {
	pageHTML, _ := b.getPageContent(url.Values{"page": {"overview"}})
	return ExtractPlanet(pageHTML, v, b)
}

func (b *OGame) getMoons() []Moon {
	pageHTML, _ := b.getPageContent(url.Values{"page": {"overview"}})
	return ExtractMoons(pageHTML, b)
}

func (b *OGame) getMoon(v interface{}) (Moon, error) {
	pageHTML, _ := b.getPageContent(url.Values{"page": {"overview"}})
	return ExtractMoon(pageHTML, b, v)
}

func (b *OGame) getCelestials() ([]Celestial, error) {
	pageHTML, _ := b.getPageContent(url.Values{"page": {"overview"}})
	return ExtractCelestials(pageHTML, b)
}

func (b *OGame) getCelestial(v interface{}) (Celestial, error) {
	pageHTML, _ := b.getPageContent(url.Values{"page": {"overview"}})
	return ExtractCelestial(pageHTML, b, v)
}

func (b *OGame) abandon(v interface{}) error {
	pageHTML, _ := b.getPageContent(url.Values{"page": {"overview"}})
	var planetID PlanetID
	if coordStr, ok := v.(string); ok {
		coord, err := ParseCoord(coordStr)
		if err != nil {
			return err
		}
		planet, err := ExtractPlanetByCoord(pageHTML, b, coord)
		if err != nil {
			return err
		}
		planetID = planet.ID
	} else if coord, ok := v.(Coordinate); ok {
		planet, err := ExtractPlanetByCoord(pageHTML, b, coord)
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
	planets := ExtractPlanets(pageHTML, b)
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
	pageHTML, _ = b.getPageContent(url.Values{"page": {"planetlayer"}, "cp": {strconv.Itoa(int(planetID))}})
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
	pageHTML, _ := b.getPageContent(url.Values{"page": {"overview"}})
	serverTime, err := extractServerTime(pageHTML)
	if err != nil {
		b.error(err.Error())
	}
	return serverTime
}

func (b *OGame) getUserInfos() UserInfos {
	pageHTML, _ := b.getPageContent(url.Values{"page": {"overview"}})
	userInfos, err := ExtractUserInfos(pageHTML, b.language)
	if err != nil {
		b.error(err)
	}
	return userInfos
}

type ChatPostResp struct {
	Status   string `json:"status"`
	ID       int    `json:"id"`
	SenderID int    `json:"senderId"`
	TargetID int    `json:"targetId"`
	Text     string `json:"text"`
	Date     int    `json:"date"`
	NewToken string `json:"newToken"`
}

func (b *OGame) sendMessage(id int, message string, isPlayer bool) error {
	payload := url.Values{
		"text":  {message + "\n"},
		"ajax":  {"1"},
		"token": {b.ajaxChatToken},
	}
	if isPlayer {
		payload.Set("playerId", strconv.Itoa(id))
		payload.Set("mode", "1")
	} else {
		payload.Set("associationId", strconv.Itoa(id))
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
	return ExtractFleetsFromEventList(pageHTML)
}

func (b *OGame) getFleets() ([]Fleet, Slots) {
	pageHTML, _ := b.getPageContent(url.Values{"page": {"movement"}})
	fleets := ExtractFleets(pageHTML)
	slots := ExtractSlots(pageHTML)
	return fleets, slots
}

func (b *OGame) cancelFleet(fleetID FleetID) error {
	_, _ = b.getPageContent(url.Values{"page": {"movement"}, "return": {fleetID.String()}})
	return nil
}

// Slots ...
type Slots struct {
	InUse    int
	Total    int
	ExpInUse int
	ExpTotal int
}

func (b *OGame) getSlots() Slots {
	pageHTML, _ := b.getPageContent(url.Values{"page": {"fleet1"}})
	return ExtractSlots(pageHTML)
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

func systemDistance(nbSystems, system1, system2 int, donutSystem bool) (distance int) {
	if !donutSystem {
		return int(math.Abs(float64(system2 - system1)))
	}
	if system1 > system2 {
		system1, system2 = system2, system1
	}
	return int(math.Min(float64(system2-system1), float64((system1+nbSystems)-system2)))
}

// Returns the distance between two systems
func flightSystemDistance(nbSystems, system1, system2 int, donutSystem bool) (distance int) {
	return 2700 + 95*systemDistance(nbSystems, system1, system2, donutSystem)
}

// Returns the distance between two planets
func planetDistance(planet1, planet2 int) (distance int) {
	return int(1000 + 5*math.Abs(float64(planet2-planet1)))
}

// Distance returns the distance between two coordinates
func Distance(c1, c2 Coordinate, universeSize, nbSystems int, donutGalaxy, donutSystem bool) (distance int) {
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

func findSlowestSpeed(ships ShipsInfos, techs Researches) int {
	minSpeed := math.MaxInt32
	for _, ship := range Ships {
		if ship.GetID() == SolarSatelliteID {
			continue
		}
		shipSpeed := ship.GetSpeed(techs)
		if ships.ByID(ship.GetID()) > 0 && shipSpeed < minSpeed {
			minSpeed = shipSpeed
		}
	}
	return minSpeed
}

func calcFuel(ships ShipsInfos, dist, duration int, universeSpeedFleet, fleetDeutSaveFactor float64, techs Researches) (fuel int) {
	tmpFn := func(baseFuel, nbr, shipSpeed int) float64 {
		tmpSpeed := (35000 / (float64(duration)*universeSpeedFleet - 10)) * math.Sqrt(float64(dist)*10/float64(shipSpeed))
		return float64(baseFuel*nbr*dist) / 35000 * math.Pow(tmpSpeed/10+1, 2)
	}
	tmpFuel := 0.0
	for _, ship := range Ships {
		if ship.GetID() == SolarSatelliteID {
			continue
		}
		nbr := ships.ByID(ship.GetID())
		if nbr > 0 {
			tmpFuel += tmpFn(ship.GetFuelConsumption(), nbr, ship.GetSpeed(techs))
		}
	}
	fuel = int(1 + math.Floor(tmpFuel*fleetDeutSaveFactor))
	return
}

func calcFlightTime(origin, destination Coordinate, universeSize, nbSystems int, donutGalaxy, donutSystem bool,
	fleetDeutSaveFactor, speed float64, universeSpeedFleet int, ships ShipsInfos, techs Researches) (secs, fuel int) {
	if !ships.HasShips() {
		return
	}
	s := speed
	v := float64(findSlowestSpeed(ships, techs))
	a := float64(universeSpeedFleet)
	d := float64(Distance(origin, destination, universeSize, nbSystems, donutGalaxy, donutSystem))
	secs = int(math.Round(((3500/s)*math.Sqrt(d*10/v) + 10) / a))
	fuel = calcFuel(ships, int(d), secs, float64(universeSpeedFleet), fleetDeutSaveFactor, techs)
	return
}

// getPhalanx makes 3 calls to ogame server (2 validation, 1 scan)
func (b *OGame) getPhalanx(moonID MoonID, coord Coordinate) ([]Fleet, error) {
	res := make([]Fleet, 0)

	// Get moon facilities html page (first call to ogame server)
	moonFacilitiesHTML, _ := b.getPageContent(url.Values{"page": {"station"}, "cp": {strconv.Itoa(int(moonID))}})

	// Extract bunch of infos from the html
	moon, err := ExtractMoon(moonFacilitiesHTML, b, moonID)
	if err != nil {
		return res, errors.New("moon not found")
	}
	resources := ExtractResources(moonFacilitiesHTML)
	moonFacilities, _ := ExtractFacilities(moonFacilitiesHTML)
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
		"galaxy":   {strconv.Itoa(coord.Galaxy)},
		"system":   {strconv.Itoa(coord.System)},
		"position": {strconv.Itoa(coord.Position)},
		"ajax":     {"1"},
		"cp":       {strconv.Itoa(int(moonID))},
	})
	return extractPhalanx(pageHTML)
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
	pageHTML, _ := b.getPageContent(url.Values{"page": {"jumpgatelayer"}, "cp": {strconv.Itoa(int(originMoonID))}})
	availShips, token, dests, wait := extractJumpGate(pageHTML)
	if wait > 0 {
		return fmt.Errorf("jump gate is in recharge mode for %d seconds", wait)
	}

	// Validate destination moon id
	if !moonIDInSlice(destMoonID, dests) {
		return errors.New("destination moon id invalid")
	}

	payload := url.Values{"token": {token}, "zm": {strconv.Itoa(int(destMoonID))}}

	// Add ships to payload
	for _, s := range Ships {
		// Get the min between what is available and what we want
		nbr := int(math.Min(float64(ships.ByID(s.GetID())), float64(availShips.ByID(s.GetID()))))
		if nbr > 0 {
			payload.Add("ship_"+strconv.Itoa(int(s.GetID())), strconv.Itoa(nbr))
		}
	}

	if _, err := b.postPageContent(url.Values{"page": {"jumpgate_execute"}}, payload); err != nil {
		return err
	}
	return nil
}

func (b *OGame) createUnion(fleet Fleet) (int, error) {
	if fleet.ID == 0 {
		return 0, errors.New("invalid fleet id")
	}
	pageHTML, _ := b.getPageContent(url.Values{"page": {"federationlayer"}, "union": {"0"}, "fleet": {strconv.Itoa(int(fleet.ID))}, "target": {strconv.Itoa(fleet.TargetPlanetID)}, "ajax": {"1"}})
	payload := ExtractFederation(pageHTML)
	by, err := b.postPageContent(url.Values{"page": {"unionchange"}, "ajax": {"1"}}, payload)
	if err != nil {
		return 0, err
	}
	var res struct {
		FleetID  int
		UnionID  int
		TargetID int
		Errorbox struct {
			Type   string
			Text   string
			Failed int
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

func calcResources(price int, planetResources PlanetResources, multiplier Multiplier) url.Values {
	sortedCelestialIDs := make([]CelestialID, 0)
	for celestialID := range planetResources {
		sortedCelestialIDs = append(sortedCelestialIDs, celestialID)
	}
	sort.Slice(sortedCelestialIDs, func(i, j int) bool {
		return int(sortedCelestialIDs[i]) < int(sortedCelestialIDs[j])
	})

	payload := url.Values{}
	remaining := price
	for celestialID, res := range planetResources {
		metalNeeded := res.Input.Metal
		if remaining < int(float64(metalNeeded)*multiplier.Metal) {
			metalNeeded = int(math.Ceil(float64(remaining) / multiplier.Metal))
		}
		remaining -= int(float64(metalNeeded) * multiplier.Metal)

		crystalNeeded := res.Input.Crystal
		if remaining < int(float64(crystalNeeded)*multiplier.Crystal) {
			crystalNeeded = int(math.Ceil(float64(remaining) / multiplier.Crystal))
		}
		remaining -= int(float64(crystalNeeded) * multiplier.Crystal)

		deuteriumNeeded := res.Input.Deuterium
		if remaining < int(float64(deuteriumNeeded)*multiplier.Deuterium) {
			deuteriumNeeded = int(math.Ceil(float64(remaining) / multiplier.Deuterium))
		}
		remaining -= int(float64(deuteriumNeeded) * multiplier.Deuterium)

		payload.Add("bid[planets]["+strconv.Itoa(int(celestialID))+"][metal]", strconv.Itoa(metalNeeded))
		payload.Add("bid[planets]["+strconv.Itoa(int(celestialID))+"][crystal]", strconv.Itoa(crystalNeeded))
		payload.Add("bid[planets]["+strconv.Itoa(int(celestialID))+"][deuterium]", strconv.Itoa(deuteriumNeeded))
	}
	return payload
}

func (b *OGame) buyOfferOfTheDay() error {
	pageHTML, err := b.postPageContent(url.Values{"page": {"traderOverview"}}, url.Values{"show": {"importexport"}, "ajax": {"1"}})
	if err != nil {
		return err
	}

	price, importToken, planetResources, multiplier, err := ExtractOfferOfTheDay(pageHTML)
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

func (b *OGame) getAttacks(celestialID CelestialID) (out []AttackEvent, err error) {
	params := url.Values{"page": {"eventList"}, "ajax": {"1"}}
	if celestialID != 0 {
		params.Set("cp", strconv.Itoa(int(celestialID)))
	}
	pageHTML, err := b.getPageContent(params)
	if err != nil {
		return
	}
	return ExtractAttacks(pageHTML)
}

func (b *OGame) galaxyInfos(galaxy, system int) (SystemInfos, error) {
	if galaxy < 0 || galaxy > b.server.Settings.UniverseSize {
		return SystemInfos{}, fmt.Errorf("galaxy must be within [0, %d]", b.server.Settings.UniverseSize)
	}
	if system < 0 || system > b.serverData.Systems {
		return SystemInfos{}, errors.New("system must be within [0, " + strconv.Itoa(b.serverData.Systems) + "]")
	}
	payload := url.Values{
		"galaxy": {strconv.Itoa(galaxy)},
		"system": {strconv.Itoa(system)},
	}
	var res SystemInfos
	pageHTML, err := b.postPageContent(url.Values{"page": {"galaxyContent"}, "ajax": {"1"}}, payload)
	if err != nil {
		return res, err
	}
	return ExtractGalaxyInfos(pageHTML, b.Player.PlayerName, b.Player.PlayerID, b.Player.Rank)
}

func (b *OGame) getResourceSettings(planetID PlanetID) (ResourceSettings, error) {
	pageHTML, _ := b.getPageContent(url.Values{"page": {"resourceSettings"}, "cp": {planetID.String()}})
	return ExtractResourceSettings(pageHTML)
}

func (b *OGame) setResourceSettings(planetID PlanetID, settings ResourceSettings) error {
	pageHTML, _ := b.getPageContent(url.Values{"page": {"resourceSettings"}, "cp": {planetID.String()}})
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	bodyID := ExtractBodyIDFromDoc(doc)
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

func getNbr(doc *goquery.Document, name string) int {
	div := doc.Find("div." + name)
	level := div.Find("span.level")
	level.Children().Remove()
	return ParseInt(level.Text())
}

func getNbrShips(doc *goquery.Document, name string) int {
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
	pageHTML, _ := b.getPageContent(url.Values{"page": {"research"}})
	researches := ExtractResearch(pageHTML)
	b.researches = &researches
	return researches
}

func (b *OGame) getResourcesBuildings(celestialID CelestialID) (ResourcesBuildings, error) {
	pageHTML, _ := b.getPageContent(url.Values{"page": {"resources"}, "cp": {strconv.Itoa(int(celestialID))}})
	return ExtractResourcesBuildings(pageHTML)
}

func (b *OGame) getDefense(celestialID CelestialID) (DefensesInfos, error) {
	pageHTML, _ := b.getPageContent(url.Values{"page": {"defense"}, "cp": {strconv.Itoa(int(celestialID))}})
	return ExtractDefense(pageHTML)
}

func (b *OGame) getShips(celestialID CelestialID) (ShipsInfos, error) {
	pageHTML, _ := b.getPageContent(url.Values{"page": {"shipyard"}, "cp": {strconv.Itoa(int(celestialID))}})
	return ExtractShips(pageHTML)
}

func (b *OGame) getFacilities(celestialID CelestialID) (Facilities, error) {
	pageHTML, _ := b.getPageContent(url.Values{"page": {"station"}, "cp": {strconv.Itoa(int(celestialID))}})
	return ExtractFacilities(pageHTML)
}

func (b *OGame) getProduction(celestialID CelestialID) ([]Quantifiable, error) {
	pageHTML, _ := b.getPageContent(url.Values{"page": {"shipyard"}, "cp": {strconv.Itoa(int(celestialID))}})
	return ExtractProduction(pageHTML)
}

func getToken(b *OGame, page string, celestialID CelestialID) (string, error) {
	pageHTML, _ := b.getPageContent(url.Values{"page": {page}, "cp": {strconv.Itoa(int(celestialID))}})
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
	pageHTML, _ := b.getPageContent(url.Values{"page": {page}, "cp": {strconv.Itoa(int(celestialID))}})
	m := regexp.MustCompile(`modus=3&token=([^&]+)&`).FindSubmatch(pageHTML)
	if len(m) != 2 {
		return "", errors.New("unable to find form token")
	}
	return string(m[1]), nil
}

func (b *OGame) tearDown(celestialID CelestialID, id ID) error {
	var page string
	if id.IsResourceBuilding() {
		page = "resources"
	} else if id.IsFacility() {
		page = "station"
	} else {
		return errors.New("invalid id " + id.String())
	}

	token, err := getDemolishToken(b, page, celestialID)
	if err != nil {
		return err
	}

	pageHTML, _ := b.getPageContent(url.Values{"page": {page}, "ajax": {"1"}, "type": {strconv.Itoa(int(celestialID))}, "cp": {strconv.Itoa(int(celestialID))}})
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return err
	}
	imgDisabled := doc.Find("a.demolish_link div").HasClass("demolish_img_disabled")
	if imgDisabled {
		return errors.New("tear down button is disabled")
	}

	params := url.Values{
		"page":  {page},
		"modus": {"3"},
		"token": {token},
		"type":  {strconv.Itoa(int(id))},
	}
	_, err = b.getPageContent(params)
	return err
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
		token, err := getToken(b, page, celestialID)
		if err != nil {
			return err
		}
		payload.Add("token", token)
	}

	if id.IsDefense() || id.IsShip() {
		maximumNbr := 99999
		var err error
		var token string
		for nbr > 0 {
			tmp := int(math.Min(float64(nbr), float64(maximumNbr)))
			payload.Set("menge", strconv.Itoa(tmp))
			_, err = b.postPageContent(url.Values{"page": {page}, "cp": {strconv.Itoa(int(celestialID))}}, payload)
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

	_, err := b.postPageContent(url.Values{"page": {page}, "cp": {strconv.Itoa(int(celestialID))}}, payload)
	return err
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

func (b *OGame) buildTechnology(celestialID CelestialID, technologyID ID) error {
	if !technologyID.IsTech() {
		return errors.New("invalid technology id " + technologyID.String())
	}
	return b.buildCancelable(celestialID, technologyID)
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
	return b.buildProduction(celestialID, shipID, nbr)
}

func (b *OGame) constructionsBeingBuilt(celestialID CelestialID) (ID, int, ID, int) {
	pageHTML, _ := b.getPageContent(url.Values{"page": {"overview"}, "cp": {strconv.Itoa(int(celestialID))}})
	return ExtractConstructions(pageHTML)
}

func (b *OGame) cancel(token string, techID, listID int) error {
	_, _ = b.getPageContent(url.Values{"page": {"overview"}, "modus": {"2"}, "token": {token},
		"techid": {strconv.Itoa(techID)}, "listid": {strconv.Itoa(listID)}})
	return nil
}

func (b *OGame) cancelBuilding(celestialID CelestialID) error {
	pageHTML, _ := b.getPageContent(url.Values{"page": {"overview"}, "cp": {strconv.Itoa(int(celestialID))}})
	token, techID, listID, _ := extractCancelBuildingInfos(pageHTML)
	return b.cancel(token, techID, listID)
}

func (b *OGame) cancelResearch(celestialID CelestialID) error {
	pageHTML, _ := b.getPageContent(url.Values{"page": {"overview"}, "cp": {strconv.Itoa(int(celestialID))}})
	token, techID, listID, _ := extractCancelResearchInfos(pageHTML)
	return b.cancel(token, techID, listID)
}

func (b *OGame) fetchResources(celestialID CelestialID) (ResourcesDetails, error) {
	pageJSON, _ := b.getPageContent(url.Values{"page": {"fetchResources"}, "cp": {strconv.Itoa(int(celestialID))}})
	return ExtractResourcesDetails(pageJSON)
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

func (b *OGame) sendIPM(planetID PlanetID, coord Coordinate, nbr int, priority ID) (int, error) {
	if priority != 0 && (!priority.IsDefense() || priority == AntiBallisticMissilesID || priority == InterplanetaryMissilesID) {
		return 0, errors.New("invalid target id")
	}
	pageHTML, err := b.getPageContent(url.Values{
		"page":       {"missileattacklayer"},
		"galaxy":     {strconv.Itoa(coord.Galaxy)},
		"system":     {strconv.Itoa(coord.System)},
		"position":   {strconv.Itoa(coord.Position)},
		"planetType": {strconv.Itoa(int(coord.Type))},
		"cp":         {strconv.Itoa(int(planetID))},
	})
	if err != nil {
		return 0, err
	}
	duration, max, token := ExtractIPM(pageHTML)
	if max == 0 {
		return 0, errors.New("no missile available")
	}
	if nbr > max {
		nbr = max
	}
	payload := url.Values{
		"galaxy":     {strconv.Itoa(coord.Galaxy)},
		"system":     {strconv.Itoa(coord.System)},
		"position":   {strconv.Itoa(coord.Position)},
		"planetType": {strconv.Itoa(int(coord.Type))},
		"token":      {token},
		"anz":        {strconv.Itoa(nbr)},
		"pziel":      {},
	}
	if priority != 0 {
		payload.Add("pziel", strconv.Itoa(int(priority)))
	}
	by, err := b.postPageContent(url.Values{"page": {"missileattack_execute"}}, payload)
	if err != nil {
		return 0, err
	}
	// {"status":false,"errorbox":{"type":"fadeBox","text":"Target doesn`t exist!","failed":1}}
	var resp struct {
		Status   bool
		ErrorBox struct {
			Type   string
			Text   string
			Failed int
		}
	}
	if err := json.Unmarshal(by, &resp); err != nil {
		return 0, err
	}
	if resp.ErrorBox.Failed == 1 {
		return 0, errors.New(resp.ErrorBox.Text)
	}
	fmt.Println(string(by))

	return duration, nil
}

func (b *OGame) sendFleet(celestialID CelestialID, ships []Quantifiable, speed Speed, where Coordinate,
	mission MissionID, resources Resources, expeditiontime, unionID int, ensure bool) (Fleet, error) {

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
	pageHTML, err := b.getPageContent(url.Values{"page": {"fleet1"}, "cp": {strconv.Itoa(int(celestialID))}})
	if err != nil {
		return Fleet{}, err
	}

	fleet1Doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	fleet1BodyID := ExtractBodyIDFromDoc(fleet1Doc)
	if fleet1BodyID != "fleet1" {
		now := time.Now().Unix()
		b.error(ErrInvalidPlanetID.Error()+", planetID:", celestialID, ", ts: ", now)
		return Fleet{}, ErrInvalidPlanetID
	}

	if ExtractIsInVacationFromDoc(fleet1Doc) {
		return Fleet{}, ErrAccountInVacationMode
	}

	// Ensure we're not trying to attack/spy ourselves
	destinationIsMyOwnPlanet := false
	myCelestials, _ := ExtractCelestialsFromDoc(fleet1Doc, b)
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

	availableShips := ExtractFleet1ShipsFromDoc(fleet1Doc)

	if !ensure {
		atLeastOneShipSelected := false
		for _, ship := range ships {
			if ship.Nbr > 0 && availableShips.ByID(ship.ID) > 0 {
				atLeastOneShipSelected = true
				break
			}
		}
		if !atLeastOneShipSelected {
			return Fleet{}, ErrNoShipSelected
		}
	} else {
		for _, ship := range ships {
			if ship.Nbr > availableShips.ByID(ship.ID) {
				return Fleet{}, ErrNotEnoughShips
			}
		}
	}

	payload := ExtractHiddenFieldsFromDoc(fleet1Doc)
	cs := false       // ColonyShip flag for fleet check
	recycler := false // Recycler flag for fleet check
	for _, s := range ships {
		if s.Nbr > 0 {
			if s.ID == ColonyShipID {
				cs = true
			} else if s.ID == RecyclerID {
				recycler = true
			}
			payload.Add("am"+strconv.Itoa(int(s.ID)), strconv.Itoa(s.Nbr))
		}
	}

	// Page 2 : select ships
	pageHTML, err = b.postPageContent(url.Values{"page": {"fleet2"}}, payload)
	if err != nil {
		return Fleet{}, err
	}
	fleet2Doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	fleet2BodyID := ExtractBodyIDFromDoc(fleet2Doc)
	if fleet2BodyID != "fleet2" {
		now := time.Now().Unix()
		b.error(errors.New("unknown error").Error()+", planetID:", celestialID, ", ts: ", now)
		return Fleet{}, errors.New("unknown error")
	}

	payload = ExtractHiddenFieldsFromDoc(fleet2Doc)
	payload.Add("speed", strconv.Itoa(int(speed)))
	payload.Add("galaxy", strconv.Itoa(where.Galaxy))
	payload.Add("system", strconv.Itoa(where.System))
	payload.Add("position", strconv.Itoa(where.Position))
	if mission == RecycleDebrisField {
		where.Type = DebrisType // Send to debris field
	} else if mission == Colonize || mission == Expedition {
		where.Type = PlanetType
	}
	payload.Add("type", strconv.Itoa(int(where.Type)))

	if unionID != 0 {
		found := false
		fleet2Doc.Find("select[name=acsValues] option").Each(func(i int, s *goquery.Selection) {
			acsValues := s.AttrOr("value", "")
			m := regexp.MustCompile(`\d+#\d+#\d+#\d+#.*#(\d+)`).FindStringSubmatch(acsValues)
			if len(m) == 2 {
				optUnionID, _ := strconv.Atoi(m[1])
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
		"galaxy": {strconv.Itoa(where.Galaxy)},
		"system": {strconv.Itoa(where.System)},
		"planet": {strconv.Itoa(where.Position)},
		"type":   {strconv.Itoa(int(where.Type))},
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
		return Fleet{}, ErrPlanetAlreadyReservecForRelocation
	}

	// Page 3 : select coord, mission, speed
	pageHTML, err = b.postPageContent(url.Values{"page": {"fleet3"}}, payload)
	if err != nil {
		return Fleet{}, err
	}

	fleet3Doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	fleet3BodyID := ExtractBodyIDFromDoc(fleet3Doc)
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

	payload = ExtractHiddenFieldsFromDoc(fleet3Doc)
	var finalShips ShipsInfos
	for k, v := range payload {
		var shipID int
		if n, err := fmt.Sscanf(k, "am%d", &shipID); err == nil && n == 1 {
			nbr, _ := strconv.Atoi(v[0])
			finalShips.Set(ID(shipID), nbr)
		}
	}
	deutConsumption := ParseInt(fleet3Doc.Find("div#roundup span#consumption").Text())
	resourcesAvailable := ExtractResourcesFromDoc(fleet3Doc)
	if deutConsumption > resourcesAvailable.Deuterium {
		return Fleet{}, fmt.Errorf("not enough deuterium, avail: %d, need: %d", resourcesAvailable.Deuterium, deutConsumption)
	}
	// finalCargo := ParseInt(fleet3Doc.Find("#maxresources").Text())
	baseCargo := finalShips.Cargo(Researches{})
	if b.GetServer().Settings.EspionageProbeRaids != 1 {
		baseCargo += finalShips.EspionageProbe * EspionageProbe.BaseCargoCapacity
	}
	if deutConsumption > baseCargo {
		return Fleet{}, fmt.Errorf("not enough cargo capacity, avail: %d, need: %d", baseCargo, deutConsumption)
	}
	payload.Add("crystal", strconv.Itoa(resources.Crystal))
	payload.Add("deuterium", strconv.Itoa(resources.Deuterium))
	payload.Add("metal", strconv.Itoa(resources.Metal))
	payload.Set("mission", strconv.Itoa(int(mission)))
	if mission == Expedition {
		payload.Set("expeditiontime", strconv.Itoa(expeditiontime))
	}

	// Page 4 : send the fleet
	_, _ = b.postPageContent(url.Values{"page": {"movement"}}, payload)

	// Page 5
	movementHTML, _ := b.getPageContent(url.Values{"page": {"movement"}})
	movementDoc, _ := goquery.NewDocumentFromReader(bytes.NewReader(movementHTML))
	originCoords, _ := ExtractPlanetCoordinate(movementHTML)
	fleets := ExtractFleetsFromDoc(movementDoc)
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

	slots = ExtractSlotsFromDoc(movementDoc)
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
	ID           int
	Origin       *Coordinate
	Destination  Coordinate
	AttackerName string
	DefenderName string
	Loot         int
	Metal        int
	Crystal      int
	Deuterium    int
	CreatedAt    time.Time
}

// EspionageReportSummary summary of espionage report
type EspionageReportSummary struct {
	ID             int
	Type           EspionageReportType
	From           string
	Target         Coordinate
	LootPercentage float64
}

func (b *OGame) getPageMessages(page, tabid int) ([]byte, error) {
	payload := url.Values{
		"messageId":  {"-1"},
		"tabid":      {strconv.Itoa(tabid)},
		"action":     {"107"},
		"pagination": {strconv.Itoa(page)},
		"ajax":       {"1"},
	}
	return b.postPageContent(url.Values{"page": {"messages"}}, payload)
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

func (b *OGame) getCombatReportMessages() ([]CombatReportSummary, error) {
	tabid := 21
	page := 1
	nbPage := 1
	msgs := make([]CombatReportSummary, 0)
	for page <= nbPage {
		pageHTML, _ := b.getPageMessages(page, tabid)
		newMessages, newNbPage := extractCombatReportMessagesSummary(pageHTML)
		msgs = append(msgs, newMessages...)
		nbPage = newNbPage
		page++
	}
	return msgs, nil
}

func (b *OGame) getCombatReportFor(coord Coordinate) (CombatReportSummary, error) {
	tabid := 21
	page := 1
	nbPage := 1
	for page <= nbPage {
		pageHTML, err := b.getPageMessages(page, tabid)
		if err != nil {
			return CombatReportSummary{}, err
		}
		newMessages, newNbPage := extractCombatReportMessagesSummary(pageHTML)
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

func (b *OGame) getEspionageReport(msgID int) (EspionageReport, error) {
	pageHTML, _ := b.getPageContent(url.Values{"page": {"messages"}, "messageId": {strconv.Itoa(msgID)}, "tabid": {"20"}, "ajax": {"1"}})
	return extractEspionageReport(pageHTML, b.location)
}

func (b *OGame) getEspionageReportFor(coord Coordinate) (EspionageReport, error) {
	tabid := 20
	page := 1
	nbPage := 1
	for page <= nbPage {
		pageHTML, err := b.getPageMessages(page, tabid)
		if err != nil {
			return EspionageReport{}, err
		}
		newMessages, newNbPage := extractEspionageReportMessageIDs(pageHTML)
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

func (b *OGame) deleteMessage(msgID int) error {
	payload := url.Values{
		"messageId": {strconv.Itoa(msgID)},
		"action":    {"103"},
		"ajax":      {"1"},
	}
	by, err := b.postPageContent(url.Values{"page": {"messages"}}, payload)
	if err != nil {
		return err
	}

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
	req, err := http.NewRequest("PUT", "https://lobby.ogame.gameforge.com/api/users/me/accounts", strings.NewReader(string(jsonPayloadBytes)))
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
func (b *OGame) GetResearchSpeed() int {
	return b.serverData.ResearchDurationDivisor
}

// Deprecated: SetResearchSpeed sets the research speed
func (b *OGame) SetResearchSpeed(newSpeed int) {
	b.serverData.ResearchDurationDivisor = newSpeed
}

// GetNbSystems gets the number of systems
func (b *OGame) GetNbSystems() int {
	return b.serverData.Systems
}

// Deprecated: SetNbSystems sets the number of speed
func (b *OGame) SetNbSystems(newNbSystems int) {
	b.serverData.Systems = newNbSystems
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

// GetCachedNbProbes returns cached number of probes from preferences
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
func (b *OGame) SendMessage(playerID int, message string) error {
	return b.WithPriority(Normal).SendMessage(playerID, message)
}

// SendMessageAlliance sends a message to associationID
func (b *OGame) SendMessageAlliance(associationID int, message string) error {
	return b.WithPriority(Normal).SendMessageAlliance(associationID, message)
}

// GetFleets get the player's own fleets activities
func (b *OGame) GetFleets() ([]Fleet, Slots) {
	return b.WithPriority(Normal).GetFleets()
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
func (b *OGame) GalaxyInfos(galaxy, system int) (SystemInfos, error) {
	return b.WithPriority(Normal).GalaxyInfos(galaxy, system)
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
func (b *OGame) GetProduction(celestialID CelestialID) ([]Quantifiable, error) {
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
func (b *OGame) Build(celestialID CelestialID, id ID, nbr int) error {
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
func (b *OGame) BuildProduction(celestialID CelestialID, id ID, nbr int) error {
	return b.WithPriority(Normal).BuildProduction(celestialID, id, nbr)
}

// BuildBuilding ensure what is being built is a building
func (b *OGame) BuildBuilding(celestialID CelestialID, buildingID ID) error {
	return b.WithPriority(Normal).BuildBuilding(celestialID, buildingID)
}

// BuildDefense builds a defense unit
func (b *OGame) BuildDefense(celestialID CelestialID, defenseID ID, nbr int) error {
	return b.WithPriority(Normal).BuildDefense(celestialID, defenseID, nbr)
}

// BuildShips builds a ship unit
func (b *OGame) BuildShips(celestialID CelestialID, shipID ID, nbr int) error {
	return b.WithPriority(Normal).BuildShips(celestialID, shipID, nbr)
}

// ConstructionsBeingBuilt returns the building & research being built, and the time remaining (secs)
func (b *OGame) ConstructionsBeingBuilt(celestialID CelestialID) (ID, int, ID, int) {
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
	mission MissionID, resources Resources, expeditiontime, unionID int) (Fleet, error) {
	return b.WithPriority(Normal).SendFleet(celestialID, ships, speed, where, mission, resources, expeditiontime, unionID)
}

// EnsureFleet either sends all the requested ships or fail
func (b *OGame) EnsureFleet(celestialID CelestialID, ships []Quantifiable, speed Speed, where Coordinate,
	mission MissionID, resources Resources, expeditiontime, unionID int) (Fleet, error) {
	return b.WithPriority(Normal).EnsureFleet(celestialID, ships, speed, where, mission, resources, expeditiontime, unionID)
}

// SendIPM sends IPM
func (b *OGame) SendIPM(planetID PlanetID, coord Coordinate, nbr int, priority ID) (int, error) {
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
func (b *OGame) GetEspionageReport(msgID int) (EspionageReport, error) {
	return b.WithPriority(Normal).GetEspionageReport(msgID)
}

// DeleteMessage deletes a message from the mail box
func (b *OGame) DeleteMessage(msgID int) error {
	return b.WithPriority(Normal).DeleteMessage(msgID)
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
func (b *OGame) FlightTime(origin, destination Coordinate, speed Speed, ships ShipsInfos) (secs, fuel int) {
	return b.WithPriority(Normal).FlightTime(origin, destination, speed, ships)
}

// Distance return distance between two coordinates
func (b *OGame) Distance(origin, destination Coordinate) int {
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

// JumpGate sends ships through a jump gate.
func (b *OGame) JumpGate(origin, dest MoonID, ships ShipsInfos) error {
	return b.WithPriority(Normal).JumpGate(origin, dest, ships)
}

// BuyOfferOfTheDay buys the offer of the day.
func (b *OGame) BuyOfferOfTheDay() error {
	return b.WithPriority(Normal).BuyOfferOfTheDay()
}

// CreateUnion creates a union
func (b *OGame) CreateUnion(fleet Fleet) (int, error) {
	return b.WithPriority(Normal).CreateUnion(fleet)
}
