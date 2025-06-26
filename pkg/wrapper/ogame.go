package wrapper

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/alaingilbert/mtx"
	"github.com/alaingilbert/ogame/pkg/device"
	"github.com/alaingilbert/ogame/pkg/exponentialBackoff"
	"github.com/alaingilbert/ogame/pkg/extractor"
	"github.com/alaingilbert/ogame/pkg/extractor/v12_0_0"
	v6 "github.com/alaingilbert/ogame/pkg/extractor/v6"
	"github.com/alaingilbert/ogame/pkg/gameforge"
	"github.com/alaingilbert/ogame/pkg/httpclient"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/parser"
	"github.com/alaingilbert/ogame/pkg/taskRunner"
	"github.com/alaingilbert/ogame/pkg/utils"
	cookiejar "github.com/orirawlings/persistent-cookiejar"
	"golang.org/x/net/proxy"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// OGame is a client for ogame.org. It is safe for concurrent use by
// multiple goroutines (thread-safe)
type OGame struct {
	sync.Mutex
	isEnabledAtom        atomic.Bool // atomic, prevent auto re login if we manually logged out
	isLoggedInAtom       atomic.Bool // atomic, prevent auto re login if we manually logged out
	isConnectedAtom      atomic.Bool // atomic, either or not communication between the bot and OGame is possible
	lockedAtom           atomic.Bool // atomic, bot state locked/unlocked
	chatConnectedAtom    atomic.Bool // atomic, either or not the chat is connected
	state                string      // keep name of the function that currently lock the bot
	parentCtx            context.Context
	ctx                  context.Context
	cancelCtx            context.CancelFunc
	stateChangeCallbacks []func(locked bool, actor string)
	quiet                bool
	universe             string
	username             string
	password             string
	otpSecret            string
	bearerToken          string
	language             string
	playerID             int64
	lobby                string
	server               gameforge.Server
	logger               *log.Logger
	chatCallbacks        []func(msg ogame.ChatMsg)
	wsCallbacks          mtx.RWMutexMap[string, func([]byte)]
	auctioneerCallbacks  []func(any)
	interceptorCallbacks []func(method, url string, params, payload url.Values, pageHTML []byte)
	closeChatCtx         context.Context
	closeChatCancel      context.CancelFunc
	ws                   *websocket.Conn
	taskRunnerInst       *taskRunner.TaskRunner[*Prioritize]
	loginWrapper         func(LoginFn) error
	loginProxyTransport  http.RoundTripper
	extractor            extractor.Extractor
	apiNewHostname       string
	captchaCallback      gameforge.CaptchaSolver
	device               *device.Device
	cache                struct {
		serverData            ServerData
		location              *time.Location
		player                ogame.UserInfos
		CachedPreferences     ogame.Preferences
		researches            *ogame.Researches
		lfBonuses             *ogame.LfBonuses
		characterClass        ogame.CharacterClass
		allianceClass         *ogame.AllianceClass
		planets               mtx.RWMutex[[]Planet]
		ogameSession          string
		token                 string
		ajaxChatToken         string
		serverURL             string
		coloniesCount         int64
		coloniesPossible      int64
		planetID              ogame.CelestialID
		isVacationModeEnabled bool
		hasCommander          bool
		hasAdmiral            bool
		hasEngineer           bool
		hasGeologist          bool
		hasTechnocrat         bool
	}
}

// Params parameters for more fine-grained initialization
type Params struct {
	Ctx            context.Context
	Username       string
	Password       string
	BearerToken    string // Gameforge auth bearer token
	OTPSecret      string
	Universe       string
	Lang           string
	PlayerID       int64
	AutoLogin      bool
	Proxy          string
	ProxyUsername  string
	ProxyPassword  string
	ProxyType      string
	ProxyLoginOnly bool
	TLSConfig      *tls.Config
	Lobby          string
	APINewHostname string
	Device         *device.Device
	CaptchaSolver  gameforge.CaptchaSolver
	Logger         *log.Logger
	Quiet          bool
}

// New creates a new instance of OGame wrapper.
func New(deviceInst *device.Device, universe, username, password, lang string) (*OGame, error) {
	return newWithParams(Params{
		Universe:  universe,
		Username:  username,
		Password:  password,
		Lang:      lang,
		Device:    deviceInst,
		AutoLogin: true,
	})
}

// NewNoLogin creates a new instance of OGame wrapper, does not auto-login.
func NewNoLogin(deviceInst *device.Device, universe, username, password, lang string) (*OGame, error) {
	return newWithParams(Params{
		Universe:  universe,
		Username:  username,
		Password:  password,
		Lang:      lang,
		Device:    deviceInst,
		AutoLogin: false,
	})
}

// NewWithParams create a new OGame instance with full control over the possible parameters
func NewWithParams(params Params) (*OGame, error) {
	return newWithParams(params)
}

func newWithParams(params Params) (*OGame, error) {
	if params.Device == nil {
		return nil, errors.New("no device defined")
	}
	if params.Ctx == nil {
		params.Ctx = context.Background()
	}
	if params.Logger == nil {
		params.Logger = log.New(os.Stdout, "", 0)
	}

	b := new(OGame)
	b.parentCtx = params.Ctx
	b.device = params.Device
	b.loginWrapper = DefaultLoginWrapper
	b.enable()
	b.quiet = params.Quiet
	b.logger = params.Logger

	b.universe = params.Universe
	b.setOGameCredentials(params.Username, params.Password, params.OTPSecret, params.BearerToken)
	b.setOGameLobby(params.Lobby)
	b.language = params.Lang
	b.playerID = params.PlayerID

	ext := v12_0_0.NewExtractor()
	ext.SetLanguage(params.Lang)
	ext.SetLocation(time.UTC)
	b.extractor = ext

	factory := func() *Prioritize { return &Prioritize{bot: b} }
	b.taskRunnerInst = taskRunner.NewTaskRunner(params.Ctx, factory)

	b.wsCallbacks.Store(make(map[string]func([]byte)))

	b.captchaCallback = params.CaptchaSolver
	b.apiNewHostname = params.APINewHostname
	if params.Proxy != "" {
		if err := b.setProxy(params.Proxy, params.ProxyUsername, params.ProxyPassword, params.ProxyType, params.ProxyLoginOnly, params.TLSConfig); err != nil {
			return nil, err
		}
	}
	if params.AutoLogin {
		if _, _, err := b.LoginWithExistingCookies(); err != nil {
			return nil, err
		}
	}
	return b, nil
}

const PLATFORM = gameforge.OGAME

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
	SpeedFleetPeaceful            int64   `xml:"speedFleetPeaceful"`            // 1
	SpeedFleetWar                 int64   `xml:"speedFleetWar"`                 // 1
	SpeedFleetHolding             int64   `xml:"speedFleetHolding"`             // 1
	Galaxies                      int64   `xml:"galaxies"`                      // 4
	Systems                       int64   `xml:"systems"`                       // 499
	ACS                           bool    `xml:"acs"`                           // 1
	RapidFire                     bool    `xml:"rapidFire"`                     // 1
	DefToTF                       bool    `xml:"defToTF"`                       // 0
	DebrisFactor                  float64 `xml:"debrisFactor"`                  // 0.5
	DebrisFactorDef               float64 `xml:"debrisFactorDef"`               // 0
	RepairFactor                  float64 `xml:"repairFactor"`                  // 0.7
	NewbieProtectionLimit         int64   `xml:"newbieProtectionLimit"`         // 500000
	NewbieProtectionHigh          int64   `xml:"newbieProtectionHigh"`          // 50000
	TopScore                      float64 `xml:"topScore"`                      // 60259362 / 1.0363090034999E+17
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
	SpeedFleet                    int64   `xml:"speedFleet"`                    // 6 // Deprecated in 8.1.0
	FleetIgnoreEmptySystems       bool    `xml:"fleetIgnoreEmptySystems"`       // 1
	FleetIgnoreInactiveSystems    bool    `xml:"fleetIgnoreInactiveSystems"`    // 1
}

// getServerData gets the server data from xml api
func getServerData(ctx context.Context, client gameforge.HttpClient, serverNumber int64, serverLang string) (ServerData, error) {
	var serverData ServerData
	serverDataURL := fmt.Sprintf("https://s%d-%s.ogame.gameforge.com/api/serverData.xml", serverNumber, serverLang)
	req, err := http.NewRequest(http.MethodGet, serverDataURL, nil)
	if err != nil {
		return serverData, err
	}
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return serverData, err
	}
	defer resp.Body.Close()
	by, err := io.ReadAll(resp.Body)
	if err != nil {
		return serverData, err
	}
	if err := xml.Unmarshal(by, &serverData); err != nil {
		return serverData, fmt.Errorf("failed to xml unmarshal %s : %w", serverDataURL, err)
	}
	serverData.SpeedFleetWar = max(serverData.SpeedFleetWar, 1)
	serverData.SpeedFleetPeaceful = max(serverData.SpeedFleetPeaceful, 1)
	serverData.SpeedFleetHolding = max(serverData.SpeedFleetHolding, 1)
	serverData.SpeedFleet = utils.Or(serverData.SpeedFleet, serverData.SpeedFleetPeaceful)
	return serverData, nil
}

// GetClientWithProxy ...
func GetClientWithProxy(proxyAddr, proxyUsername, proxyPassword, proxyType string, config *tls.Config) (*http.Client, error) {
	var err error
	client := &http.Client{}
	client.Transport, err = getTransport(proxyAddr, proxyUsername, proxyPassword, proxyType, config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (b *OGame) validateAccount(code string) error {
	return b.device.GetClient().WithTransport(b.loginProxyTransport, func(client *httpclient.Client) error {
		return gameforge.ValidateAccount(b.ctx, client, PLATFORM, b.lobby, code)
	})
}

func (b *OGame) execInterceptorCallbacks(method, url string, params, payload url.Values, pageHTML []byte) {
	for _, fn := range b.interceptorCallbacks {
		fn(method, url, params, payload, pageHTML)
	}
}

// V11 IntroBypass
func (b *OGame) introBypass(page *parser.OverviewPage) error {
	if bytes.Contains(page.GetContent(), []byte(`currentPage = "intro";`)) {
		b.debug("bypassing intro page")
		vals := url.Values{
			"page":      {"ingame"},
			"component": {"intro"},
			"action":    {"continueToClassSelection"},
		}
		payload := url.Values{
			"username":  {b.cache.player.PlayerName},
			"isVeteran": {"1"},
		}
		if _, err := b.postPageContent(vals, payload); err != nil {
			return err
		}
	}
	return nil
}

func postSessions(b *OGame) (bearerToken string, err error) {
	b.debug("post sessions")
	client := b.device.GetClient()
	if err := client.WithTransport(b.loginProxyTransport, func(client *httpclient.Client) error {
		gf, _ := gameforge.New(&gameforge.Config{
			Ctx:      b.ctx,
			Device:   b.device,
			Platform: PLATFORM,
			Lobby:    b.lobby,
			Solver:   b.captchaCallback,
		})
		res, err := gf.Login(&gameforge.LoginParams{
			Username:  b.username,
			Password:  b.password,
			OtpSecret: b.otpSecret,
		})
		if err != nil {
			return err
		}
		bearerToken = res.Token
		return err
	}); err != nil {
		return "", err
	}

	// put in cookie jar so that we can re-login reusing the cookies
	appendCookie(client, &http.Cookie{
		Name:   gameforge.TokenCookieName,
		Value:  bearerToken,
		Path:   "/",
		Domain: ".gameforge.com",
	})
	b.bearerToken = bearerToken
	return bearerToken, nil
}

func appendCookie(client *httpclient.Client, cookie *http.Cookie) {
	u, err := url.Parse("https://gameforge.com")
	if err != nil {
		panic(err)
	}
	cookies := client.Jar.Cookies(u)
	cookies = append(cookies, cookie)
	client.Jar.SetCookies(u, cookies)
}

func convertCelestialGeneric[T, U any](b *OGame, moonsIn []T, convertFn func(*OGame, T) U) []U {
	out := make([]U, len(moonsIn))
	for i, moon := range moonsIn {
		out[i] = convertFn(b, moon)
	}
	return out
}

func convertPlanets(b *OGame, planetsIn []ogame.Planet) []Planet {
	return convertCelestialGeneric(b, planetsIn, convertPlanet)
}

func convertMoons(b *OGame, moonsIn []ogame.Moon) []Moon {
	return convertCelestialGeneric(b, moonsIn, convertMoon)
}

func convertCelestials(b *OGame, celestials []ogame.Celestial) []Celestial {
	return convertCelestialGeneric(b, celestials, convertCelestial)
}

func convertPlanet(b *OGame, planet ogame.Planet) Planet {
	newPlanet := Planet{ogame: b, Planet: planet}
	if planet.Moon != nil {
		moon := convertMoon(b, *planet.Moon)
		newPlanet.Moon = &moon
	}
	return newPlanet
}

func convertMoon(b *OGame, moonIn ogame.Moon) Moon {
	return Moon{ogame: b, Moon: moonIn}
}

func convertCelestial(b *OGame, celestial ogame.Celestial) Celestial {
	switch v := celestial.(type) {
	case ogame.Planet:
		return convertPlanet(b, v)
	case ogame.Moon:
		return convertMoon(b, v)
	case *ogame.Moon:
		return convertMoon(b, *v)
	}
	return nil
}

func (b *OGame) cacheFullPageInfo(page parser.IFullPage) {
	b.cache.planets.Store(convertPlanets(b, page.ExtractPlanets()))
	b.cache.isVacationModeEnabled = page.ExtractIsInVacation()
	b.cache.token, _ = page.ExtractToken()
	b.cache.ajaxChatToken, _ = page.ExtractAjaxChatToken()
	b.cache.characterClass, _ = page.ExtractCharacterClass()
	b.cache.hasCommander = page.ExtractCommander()
	b.cache.hasAdmiral = page.ExtractAdmiral()
	b.cache.hasEngineer = page.ExtractEngineer()
	b.cache.hasGeologist = page.ExtractGeologist()
	b.cache.hasTechnocrat = page.ExtractTechnocrat()
	b.cache.coloniesCount, b.cache.coloniesPossible = page.ExtractColonies()
	b.cache.planetID, _ = page.ExtractPlanetID()

	switch castedPage := page.(type) {
	case *parser.OverviewPage:
		if playerInfo, err := castedPage.ExtractUserInfos(); err == nil {
			b.cache.player = playerInfo
		}
	case *parser.PreferencesPage:
		b.cache.CachedPreferences = castedPage.ExtractPreferences()
	case *parser.ResearchPage:
		researches := castedPage.ExtractResearch()
		b.cache.researches = &researches
	case *parser.LfBonusesPage:
		if bonuses, err := castedPage.ExtractLfBonuses(); err == nil {
			b.cache.lfBonuses = &bonuses
		}
	}
}

// LoginFn ...
type LoginFn func() (bool, bool, error)

// DefaultLoginWrapper ...
var DefaultLoginWrapper = func(loginFn LoginFn) error {
	_, _, err := loginFn()
	return err
}

func (b *OGame) setOGameLobby(lobby string) {
	if lobby != gameforge.LobbyPioneers {
		lobby = gameforge.Lobby
	}
	b.lobby = lobby
}

// execute a request using the login proxy transport if set
func (b *OGame) doReqWithLoginProxyTransport(req *http.Request) (resp *http.Response, err error) {
	req = req.WithContext(b.ctx)
	_ = b.device.GetClient().WithTransport(b.loginProxyTransport, func(client *httpclient.Client) error {
		resp, err = client.Do(req)
		return nil
	})
	return
}

func getTransport(proxy, username, password, proxyType string, config *tls.Config) (http.RoundTripper, error) {
	var err error
	transport := http.DefaultTransport.(*http.Transport).Clone()
	if proxyType == "socks5" {
		transport, err = getSocks5Transport(proxy, username, password)
	} else if proxyType == "http" {
		transport, err = getProxyTransport(proxy, username, password)
	}
	if transport != nil {
		transport.TLSClientConfig = config
	}
	return transport, err
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

func (b *OGame) setProxy(proxyAddress, username, password, proxyType string, loginOnly bool, config *tls.Config) error {
	client := b.device.GetClient()
	proxyType = utils.Or(proxyType, "socks5")
	transport := http.DefaultTransport
	var loginTransport http.RoundTripper
	if proxyAddress != "" {
		proxyTransport, err := getTransport(proxyAddress, username, password, proxyType, config)
		if err != nil {
			return err
		}
		loginTransport = proxyTransport
		if !loginOnly {
			transport = proxyTransport
		}
	}
	b.loginProxyTransport = loginTransport
	client.SetTransport(transport)
	return nil
}

func (b *OGame) connectChat(chatRetry *exponentialBackoff.ExponentialBackoff, host, port string, sessionChatCounter *int64) {
	b.connectChatV8(chatRetry, host, port, sessionChatCounter)
}

// Socket IO v3 timestamp encoding
// https://github.com/unshiftio/yeast/blob/28d15f72fc5a4273592bc209056c328a54e2b522/index.js#L17
// fmt.Println(yeast(time.Now().UnixNano() / 1000000))
func yeast(num int64) (encoded string) {
	alphabet := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_"
	length := int64(len(alphabet))
	for num > 0 {
		encoded = string(alphabet[int(num%length)]) + encoded
		num = int64(math.Floor(float64(num / length)))
	}
	return
}

func getWebsocket(host, port string) (*websocket.Conn, error) {
	token := yeast(time.Now().UnixNano() / 1000000)
	req, err := http.NewRequest(http.MethodGet, "https://"+host+":"+port+"/socket.io/?EIO=4&transport=polling&t="+token, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get socket.io token: %w", err)
	}
	defer resp.Body.Close()
	by, _ := io.ReadAll(resp.Body)
	m := regexp.MustCompile(`"sid":"([^"]+)"`).FindSubmatch(by)
	if len(m) != 2 {
		return nil, fmt.Errorf("failed to get websocket sid: %s", string(by))
	}
	sid := string(m[1])

	origin := "https://" + host + ":" + port + "/"
	wssURL := "wss://" + host + ":" + port + "/socket.io/?EIO=4&transport=websocket&sid=" + sid
	ws, err := websocket.Dial(wssURL, "", origin)
	if err != nil {
		return nil, fmt.Errorf("failed to dial websocket: %w", err)
	}
	return ws, nil
}

func (b *OGame) connectChatV8(chatRetry *exponentialBackoff.ExponentialBackoff, host, port string, sessionChatCounter *int64) {
	ws, err := getWebsocket(host, port)
	if err != nil {
		b.error("failed to dial websocket:", err)
		return
	}
	defer ws.Close()
	b.ws = ws
	chatRetry.Reset()
	_ = websocket.Message.Send(ws, "2probe")

	// Recv msgs
	for {
		if b.closeChatCtx.Err() != nil {
			return
		}

		var buf string
		if err := ws.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
			b.error("failed to set read deadline:", err)
		}
		if err := websocket.Message.Receive(ws, &buf); err != nil {
			var ne net.Error
			if err == io.EOF {
				b.error("chat eof:", err)
				break
			} else if errors.Is(err, net.ErrClosed) {
				break
			} else if errors.As(err, &ne) && ne.Timeout() {
				continue
			} else {
				b.error("chat unexpected error", err)
				// connection reset by peer
				break
			}
		}
		b.wsCallbacks.Each(func(_ string, clb func(msg []byte)) {
			go clb([]byte(buf))
		})
		if buf == "3probe" {
			_ = websocket.Message.Send(ws, "5")
			_ = websocket.Message.Send(ws, "40/chat,")
			_ = websocket.Message.Send(ws, "40/auctioneer,")
		} else if buf == "2" {
			_ = websocket.Message.Send(ws, "3")
		} else if regexp.MustCompile(`40/auctioneer,{"sid":"[^"]+"}`).MatchString(buf) {
			b.debug("got auctioneer sid")
		} else if regexp.MustCompile(`40/chat,{"sid":"[^"]+"}`).MatchString(buf) {
			b.debug("got chat sid")
			_ = websocket.Message.Send(ws, `42/chat,`+utils.FI64(*sessionChatCounter)+`["authorize","`+b.cache.ogameSession+`"]`)
			*sessionChatCounter++
		} else if regexp.MustCompile(`43/chat,\d+\[true]`).MatchString(buf) {
			b.debug("chat connected")
		} else if regexp.MustCompile(`43/chat,\d+\[false]`).MatchString(buf) {
			b.error("Failed to connect to chat")
		} else if strings.HasPrefix(buf, `42/chat,["chat",`) {
			payload := strings.TrimPrefix(buf, `42/chat,["chat",`)
			payload = strings.TrimSuffix(payload, `]`)
			var chatMsg ogame.ChatMsg
			if err := json.Unmarshal([]byte(payload), &chatMsg); err != nil {
				b.error("Unable to unmarshal chat payload", err, payload)
				continue
			}
			for _, clb := range b.chatCallbacks {
				clb(chatMsg)
			}
		} else if regexp.MustCompile(`^\d+/auctioneer`).MatchString(buf) {
			pck, err := processAuctioneerMessage(buf)
			if err != nil {
				b.error(err.Error())
				continue
			}
			for _, clb := range b.auctioneerCallbacks {
				clb(pck)
			}
		} else {
			b.error("unknown message received:", buf)
			select {
			case <-time.After(time.Second):
			}
		}
	}
}

func processAuctioneerMessage(buf string) (any, error) {
	// 42/auctioneer,["timeLeft","<span style=\"color:#99CC00;\"><b>approx. 30m</b></span> remaining until the auction ends"] // every minute
	// 42/auctioneer,["timeLeft","Next auction in:<br />\n<span class=\"nextAuction\" id=\"nextAuction\">117</span>"]
	// 42/auctioneer,["new bid",{"player":{"id":219657,"name":"Payback","link":"https://s129-en.ogame.gameforge.com/game/index.php?page=ingame&component=galaxy&galaxy=2&system=146"},"sum":5000,"price":6000,"bids":5,"auctionId":"42894"}]
	// 42/auctioneer,["new auction",{"info":"<span style=\"color:#99CC00;\"><b>approx. 35m</b></span> remaining until the auction ends","item":{"uuid":"0968999df2fe956aa4a07aea74921f860af7d97f","image":"55d4b1750985e4843023d7d0acd2b9bafb15f0b7","rarity":"rare"},"oldAuction":{"item":{"uuid":"3c9f85221807b8d593fa5276cdf7af9913c4a35d","imageSmall":"286f3eaf6072f55d8858514b159d1df5f16a5654","rarity":"common"},"time":"20.05.2021 08:42:07","bids":5,"sum":5000,"player":{"id":219657,"name":"Payback","link":"http://s129-en.ogame.gameforge.com/game/index.php?page=ingame&component=galaxy&galaxy=2&system=146"}},"auctionId":42895}]
	// 42/auctioneer,["auction finished",{"sum":5000,"player":{"id":219657,"name":"Payback","link":"http://s129-en.ogame.gameforge.com/game/index.php?page=ingame&component=galaxy&galaxy=2&system=146"},"bids":5,"info":"Next auction in:<br />\n<span class=\"nextAuction\" id=\"nextAuction\">1072</span>","time":"08:42"}]
	parts := strings.SplitN(buf, ",", 2)
	msg := parts[1]
	var pck any = msg
	var out []any
	_ = json.Unmarshal([]byte(msg), &out)
	if len(out) == 0 {
		return nil, fmt.Errorf("unknown message received: %s", buf)
	}
	name, ok := out[0].(string)
	if !ok {
		return pck, nil
	}
	arg := out[1]
	if name == "new bid" {
		if firstArg, ok := arg.(map[string]any); ok {
			auctionID := utils.DoParseI64(utils.DoCastStr(firstArg["auctionId"]))
			pck1 := ogame.AuctioneerNewBid{
				Sum:       int64(utils.DoCastF64(firstArg["sum"])),
				Price:     int64(utils.DoCastF64(firstArg["price"])),
				Bids:      int64(utils.DoCastF64(firstArg["bids"])),
				AuctionID: auctionID,
			}
			if player, ok := firstArg["player"].(map[string]any); ok {
				pck1.Player.ID = int64(utils.DoCastF64(player["id"]))
				pck1.Player.Name = utils.DoCastStr(player["name"])
				pck1.Player.Link = utils.DoCastStr(player["link"])
			}
			pck = pck1
		}
	} else if name == "timeLeft" {
		if timeLeftMsg, ok := arg.(string); ok {
			if strings.Contains(timeLeftMsg, "color:") {
				if doc, err := goquery.NewDocumentFromReader(strings.NewReader(timeLeftMsg)); err == nil {
					rgx := regexp.MustCompile(`\d+`)
					txt := rgx.FindString(doc.Find("b").Text())
					approx := utils.DoParseI64(txt)
					pck = ogame.AuctioneerTimeRemaining{Approx: approx * 60}
				}
			} else if strings.Contains(timeLeftMsg, "nextAuction") {
				if doc, err := goquery.NewDocumentFromReader(strings.NewReader(timeLeftMsg)); err == nil {
					rgx := regexp.MustCompile(`\d+`)
					txt := rgx.FindString(doc.Find("span").Text())
					secs := utils.DoParseI64(txt)
					pck = ogame.AuctioneerNextAuction{Secs: secs}
				}
			}
		}
	} else if name == "new auction" {
		if firstArg, ok := arg.(map[string]any); ok {
			pck1 := ogame.AuctioneerNewAuction{
				AuctionID: int64(utils.DoCastF64(firstArg["auctionId"])),
			}
			if infoMsg, ok := firstArg["info"].(string); ok {
				if doc, err := goquery.NewDocumentFromReader(strings.NewReader(infoMsg)); err == nil {
					rgx := regexp.MustCompile(`\d+`)
					txt := rgx.FindString(doc.Find("b").Text())
					approx := utils.DoParseI64(txt)
					pck1.Approx = approx * 60
				}
			}
			pck = pck1
		}
	} else if name == "auction finished" {
		if firstArg, ok := arg.(map[string]any); ok {
			pck1 := ogame.AuctioneerAuctionFinished{
				Sum:  int64(utils.DoCastF64(firstArg["sum"])),
				Bids: int64(utils.DoCastF64(firstArg["bids"])),
			}
			if player, ok := firstArg["player"].(map[string]any); ok {
				pck1.Player.ID = int64(utils.DoCastF64(player["id"]))
				pck1.Player.Name = utils.DoCastStr(player["name"])
				pck1.Player.Link = utils.DoCastStr(player["link"])
			}
			pck = pck1
		}
	}
	return pck, nil
}

func (b *OGame) logout() error {
	_, err := b.getPage(LogoutPageName)
	if err != nil {
		return err
	}
	if err := b.device.GetClient().Jar.(*cookiejar.Jar).Save(); err != nil {
		return err
	}
	b.softLogout()
	return nil
}

// Simulate closing the browser without logging out
func (b *OGame) softLogout() {
	if b.isLoggedInAtom.CompareAndSwap(true, false) {
		b.closeChatCancel()
	}
}

// IsKnowFullPage ...
func IsKnowFullPage(vals url.Values) bool {
	page := getPageName(vals)
	return utils.InArr(page, []string{
		OverviewPageName,
		TraderOverviewPageName,
		ResearchPageName,
		ShipyardPageName,
		GalaxyPageName,
		AlliancePageName,
		PremiumPageName,
		ShopPageName,
		RewardsPageName,
		ResourceSettingsPageName,
		MovementPageName,
		HighscorePageName,
		BuddiesPageName,
		PreferencesPageName,
		MessagesPageName,
		ChatPageName,
		DefensesPageName,
		SuppliesPageName,
		LfBuildingsPageName,
		LfResearchPageName,
		FacilitiesPageName,
		FleetdispatchPageName,
		LfBonusesPageName,
	})
}

func IsEmpirePage(vals url.Values) bool {
	return vals.Get("page") == "standalone" && vals.Get("component") == "empire"
}

// IsAjaxPage either the requested page is a partial/ajax page
func IsAjaxPage(vals url.Values) bool {
	page := getPageName(vals)
	ajax := vals.Get("ajax")
	asJson := vals.Get("asJson")
	return vals.Get("page") == "ajax" ||
		ajax == "1" ||
		asJson == "1" ||
		utils.InArr(page, []string{
			FetchEventboxAjaxPageName,
			FetchResourcesAjaxPageName,
			GalaxyContentAjaxPageName,
			EventListAjaxPageName,
			AjaxChatAjaxPageName,
			NoticesAjaxPageName,
			RepairlayerAjaxPageName,
			TechtreeAjaxPageName,
			PhalanxAjaxPageName,
			ShareReportOverlayAjaxPageName,
			JumpgatelayerAjaxPageName,
			FederationlayerAjaxPageName,
			UnionchangeAjaxPageName,
			ChangenickAjaxPageName,
			PlanetlayerAjaxPageName,
			TraderlayerAjaxPageName,
			PlanetRenameAjaxPageName,
			RightmenuAjaxPageName,
			AllianceOverviewAjaxPageName,
			SupportAjaxPageName,
			BuffActivationAjaxPageName,
			AuctioneerAjaxPageName,
			HighscoreContentAjaxPageName,
			LfResearchLayerPageName,
			LfResearchResetLayerPageName,
		})
}

func canParseEventBox(by []byte) bool {
	return json.Unmarshal(by, &eventboxResp{}) == nil
}

func canParseSystemInfos(by []byte) bool {
	return json.Unmarshal(by, &ogame.SystemInfos{}) == nil
}

func canParseNewSystemInfos(by []byte) bool {
	var success struct{ Success bool }
	return json.Unmarshal(by, &success) == nil
}

func (b *OGame) preRequestChecks() error {
	if b.parentCtx.Err() != nil {
		return b.parentCtx.Err()
	}
	if !b.IsEnabled() {
		return ogame.ErrBotInactive
	}
	if !b.IsLoggedIn() {
		return ogame.ErrBotLoggedOut
	}
	if b.cache.serverURL == "" {
		return errors.New("serverURL is empty")
	}
	return nil
}

func (b *OGame) execRequest(method, finalURL string, payload, vals url.Values) ([]byte, error) {
	var body io.Reader
	if method == http.MethodPost {
		body = strings.NewReader(payload.Encode())
	}

	req, err := http.NewRequest(method, finalURL, body)
	if err != nil {
		return []byte{}, err
	}

	if method == http.MethodPost {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	if IsAjaxPage(vals) {
		req.Header.Add("X-Requested-With", "XMLHttpRequest")
	}

	req = req.WithContext(b.ctx)
	var resp *http.Response
	if vals.Get("component") == "support" {
		err = b.device.GetClient().WithTransport(b.loginProxyTransport, func(client *httpclient.Client) error {
			resp, err = client.Do(req)
			return err
		})
	} else {
		resp, err = b.device.GetClient().Do(req)
	}
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b.error("bad status | vals:" + vals.Encode() + " | " + resp.Status)
	}

	if resp.StatusCode >= http.StatusInternalServerError {
		return []byte{}, nil
	}
	return io.ReadAll(resp.Body)
}

func getPageName(vals url.Values) string {
	page := vals.Get("page")
	component := vals.Get("component")
	if page == "ingame" ||
		(page == "componentOnly" && component == FetchEventboxAjaxPageName) ||
		(page == "componentOnly" && component == EventListAjaxPageName && vals.Get("action") != "fetchEventBox") {
		page = component
	}
	return page
}

func getOptions(opts ...Option) (out Options) {
	for _, opt := range opts {
		if opt != nil {
			opt(&out)
		}
	}
	return
}

func setCPParam(b *OGame, vals url.Values, cfg Options) {
	celestial, _ := b.getCachedCelestial(cfg.ChangePlanet)
	if vals.Get("cp") == "" &&
		cfg.ChangePlanet != 0 &&
		celestial != nil {
		vals.Set("cp", utils.FI64(cfg.ChangePlanet))
	}
}

func detectLoggedOut(method, page string, vals url.Values, pageHTML []byte) bool {
	if vals.Get("allianceId") != "" {
		return false
	}
	switch method {
	case http.MethodGet:
		return (page != LogoutPageName && (IsKnowFullPage(vals) || page == "") && !IsAjaxPage(vals) && !v6.IsLogged(pageHTML)) ||
			(page == EventListAjaxPageName && !bytes.Contains(pageHTML, []byte("eventListWrap"))) ||
			(page == FetchEventboxAjaxPageName && !canParseEventBox(pageHTML))

	case http.MethodPost:
		return (page == GalaxyContentAjaxPageName && !canParseSystemInfos(pageHTML)) ||
			(page == GalaxyAjaxPageName && !canParseNewSystemInfos(pageHTML))
	}
	return false
}

func constructFinalURL(b *OGame, vals url.Values) string {
	finalURL := b.cache.serverURL + "/game/index.php?" + vals.Encode()

	allianceID := vals.Get("allianceId")
	if allianceID != "" {
		finalURL = b.cache.serverURL + "/game/allianceInfo.php?allianceId=" + allianceID
	}
	return finalURL
}

func retryPolicyFromConfig(b *OGame, cfg Options) func(func() error) error {
	return utils.Ternary(cfg.SkipRetry, b.withoutRetry, b.withRetry)
}

func (b *OGame) getPageContent(vals url.Values, opts ...Option) ([]byte, error) {
	return b.pageContent(http.MethodGet, vals, nil, opts...)
}

func (b *OGame) postPageContent(vals, payload url.Values, opts ...Option) ([]byte, error) {
	return b.pageContent(http.MethodPost, vals, payload, opts...)
}

func (b *OGame) pageContent(method string, vals, payload url.Values, opts ...Option) ([]byte, error) {
	cfg := getOptions(opts...)

	if err := b.preRequestChecks(); err != nil {
		return []byte{}, err
	}

	setCPParam(b, vals, cfg)

	alterPayload(method, b, vals, payload)

	finalURL := constructFinalURL(b, vals)

	page := getPageName(vals)
	var pageHTMLBytes []byte

	clb := func() (err error) {
		if method == http.MethodPost || vals.Get("component") == "support" {
			// Needs to be inside the withRetry, so if we need to re-login the redirect is back for the login call
			// Prevent redirect (301) https://stackoverflow.com/a/38150816/4196220
			// Content returned on "support" endpoint:
			// <script>document.location.href='https://ogame.support.gameforge.com/index.php?fld=en&sso=login&key=100000-XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX&origin=sXXX-en.ogame.gameforge.com'</script>
			client := b.device.GetClient()
			client.CheckRedirect = func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }
			defer func() { client.CheckRedirect = nil }()
		}

		if err = applyDelay(b, cfg.Delay); err != nil {
			return err
		}

		pageHTMLBytes, err = b.execRequest(method, finalURL, payload, vals)
		if err != nil {
			return err
		}

		if detectLoggedOut(method, page, vals, pageHTMLBytes) {
			b.error("Err not logged on page : ", page)
			saveNotLoggedHTML(page, pageHTMLBytes)
			b.isConnectedAtom.Store(false)
			return ogame.ErrNotLogged
		}

		return nil
	}

	retryPolicy := retryPolicyFromConfig(b, cfg)
	if err := retryPolicy(clb); err != nil {
		b.error(err)
		return []byte{}, err
	}

	if err := processResponseHTML(method, b, pageHTMLBytes, page, payload, vals, cfg.SkipCacheFullPage); err != nil {
		return []byte{}, err
	}

	if !cfg.SkipInterceptor {
		go func() {
			b.execInterceptorCallbacks(method, finalURL, vals, payload, pageHTMLBytes)
		}()
	}

	return pageHTMLBytes, nil
}

// Save html when we detect a "notLogged", only keep last 20 files, delete others
func saveNotLoggedHTML(page string, pageHTMLBytes []byte) {
	if home, err := os.UserHomeDir(); err == nil {
		type FileInfo struct {
			Name    string
			ModTime time.Time
		}
		notLoggedPath := filepath.Join(home, ".ogame", "not_logged")
		if err := os.MkdirAll(notLoggedPath, 0755); err == nil {
			if entries, err := os.ReadDir(notLoggedPath); err == nil {
				// Create a slice to hold file info
				fileInfos := make([]FileInfo, 0, len(entries))
				for _, entry := range entries {
					file, err := entry.Info()
					if err != nil {
						continue
					}
					if !file.IsDir() {
						fileInfos = append(fileInfos, FileInfo{
							Name:    file.Name(),
							ModTime: file.ModTime(),
						})
					}
				}
				sort.Slice(fileInfos, func(i, j int) bool { return fileInfos[i].ModTime.After(fileInfos[j].ModTime) })
				if len(fileInfos) > 20 {
					for _, file := range fileInfos[20:] {
						_ = os.Remove(filepath.Join(notLoggedPath, file.Name))
					}
				}
				filename := fmt.Sprintf("not_logged_%s_%s_.html", page, time.Now().Format("2006-01-02_15:04:05"))
				_ = os.WriteFile(filepath.Join(notLoggedPath, filename), pageHTMLBytes, 0644)
			}
		}
	}
}

func applyDelay(b *OGame, delay time.Duration) error {
	if delay > 0 {
		select {
		case <-time.After(delay):
		case <-b.parentCtx.Done():
			return b.parentCtx.Err()
		case <-b.ctx.Done():
			return ogame.ErrBotInactive
		}
	}
	return nil
}

func alterPayload(method string, b *OGame, vals, payload url.Values) {
	switch method {
	case http.MethodPost:
		page := vals.Get("page")
		if page == "ingame" {
			page = vals.Get("component")
		}
		if page == "ajaxChat" && payload.Get("mode") == "1" {
			payload.Set("token", b.cache.ajaxChatToken)
		} else if (page == "fleetdispatch" && vals.Get("action") == "miniFleet") ||
			(page == "movement" && vals.Get("action") == "recallFleetAjax") {
			payload.Set("token", b.cache.token)
		}
	}
}

func processResponseHTML(method string, b *OGame, pageHTML []byte, page string, payload, vals url.Values, SkipCacheFullPage bool) error {
	extractNewAjaxToken := func() {
		var res struct {
			NewAjaxToken string `json:"newAjaxToken"`
		}
		if err := json.Unmarshal(pageHTML, &res); err == nil {
			if res.NewAjaxToken != "" {
				b.cache.token = res.NewAjaxToken
			}
		}
	}
	isAjax := IsAjaxPage(vals)
	switch method {
	case http.MethodGet:
		if !isAjax && !IsEmpirePage(vals) && v6.IsLogged(pageHTML) {
			if !SkipCacheFullPage {
				parsedFullPage := parser.AutoParseFullPage(b.extractor, pageHTML)
				b.cacheFullPageInfo(parsedFullPage)
			}
		} else if vals.Get("page") == "ajax" && vals.Get("component") == "lfbonuses" {
			if bonuses, err := b.extractor.ExtractLfBonuses(pageHTML); err == nil {
				b.cache.lfBonuses = &bonuses
			}
		} else if isAjax && vals.Get("component") == "alliance" && vals.Get("tab") == "overview" && vals.Get("action") == "fetchOverview" {
			if !SkipCacheFullPage {
				var res parser.AllianceOverviewTabRes
				if err := json.Unmarshal(pageHTML, &res); err == nil {
					allianceClass, _ := b.extractor.ExtractAllianceClass([]byte(res.Content.AllianceAllianceOverview))
					b.cache.allianceClass = &allianceClass
					b.cache.token = res.NewAjaxToken
				}
			}
		} else if isAjax {
			extractNewAjaxToken()
		}

	case http.MethodPost:
		if page == PreferencesPageName {
			b.cache.token, _ = b.extractor.ExtractToken(pageHTML)
			if prefs, err := b.extractor.ExtractPreferences(pageHTML); err == nil {
				b.cache.CachedPreferences = prefs
			}
		} else if page == "ajaxChat" && (payload.Get("mode") == "1" || payload.Get("mode") == "3") {
			if err := extractNewChatToken(b, pageHTML); err != nil {
				return err
			}
		} else if isAjax {
			extractNewAjaxToken()
		}
	}
	return nil
}

func extractNewChatToken(b *OGame, pageHTMLBytes []byte) error {
	var res ChatPostResp
	if err := json.Unmarshal(pageHTMLBytes, &res); err != nil {
		return err
	}
	b.cache.ajaxChatToken = res.NewToken
	return nil
}

type eventboxResp struct {
	Hostile  int
	Neutral  int
	Friendly int
}

func (b *OGame) withoutRetry(fn func() error) error {
	return fn()
}

func (b *OGame) withRetry(fn func() error) error {
	maxRetry := 10
	eb := exponentialBackoff.New(b.ctx, 60)
	for range eb.Iterator() {
		err := fn()
		if err == nil {
			return nil
		}
		// If we manually logged out, do not try to auto re login.
		if err := b.parentCtx.Err(); err != nil {
			return err
		}
		if !b.IsEnabled() {
			return ogame.ErrBotInactive
		}
		if !b.IsLoggedIn() {
			return ogame.ErrBotLoggedOut
		}
		maxRetry--
		if maxRetry <= 0 {
			return fmt.Errorf("failed to execute callback: %w", err)
		}
		b.error(err.Error())
		if errors.Is(err, ogame.ErrNotLogged) {
			if _, _, loginErr := b.wrapLoginWithExistingCookies(); loginErr != nil {
				b.error(loginErr.Error()) // log error
				var accountBlockedError *gameforge.AccountBlockedError
				if errors.Is(loginErr, gameforge.ErrAccountNotFound) ||
					errors.As(loginErr, &accountBlockedError) ||
					errors.Is(loginErr, gameforge.ErrBadCredentials) ||
					errors.Is(loginErr, gameforge.ErrOTPRequired) ||
					errors.Is(loginErr, gameforge.ErrOTPInvalid) {
					return loginErr
				}
			}
		}
	}
	return ogame.ErrBotInactive
}

func (b *OGame) getPageJSON(vals url.Values, v any) error {
	pageJSON, err := b.getPageContent(vals)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(pageJSON, v); err != nil {
		return ogame.ErrNotLogged
	}
	return nil
}

func (b *OGame) constructionTime(id ogame.ID, nbr int64, facilities ogame.Facilities) time.Duration {
	obj := ogame.Objs.ByID(id)
	if obj == nil {
		return 0
	}
	lfBonuses, _ := b.getCachedLfBonuses()
	return obj.ConstructionTime(nbr, b.getUniverseSpeed(), facilities, lfBonuses, b.cache.characterClass, b.cache.hasTechnocrat)
}

func (b *OGame) setOGameCredentials(username, password, otpSecret, bearerToken string) {
	b.username = username
	b.password = password
	b.otpSecret = otpSecret
	b.bearerToken = bearerToken
}

func (b *OGame) enable() {
	b.ctx, b.cancelCtx = context.WithCancel(b.parentCtx)
	b.isEnabledAtom.Store(true)
	b.stateChanged(false, "Enable")
}

func (b *OGame) disable() {
	b.isEnabledAtom.Store(false)
	b.cancelCtx()
	b.stateChanged(false, "Disable")
}

func (b *OGame) isEnabled() bool {
	return b.isEnabledAtom.Load()
}

func (b *OGame) isLoggedIn() bool {
	return b.isLoggedInAtom.Load()
}

func (b *OGame) isConnected() bool {
	return b.isConnectedAtom.Load()
}

func (b *OGame) isCollector() bool {
	return b.cache.characterClass == ogame.Collector
}

func (b *OGame) isGeneral() bool {
	return b.cache.characterClass == ogame.General
}

func (b *OGame) isDiscoverer() bool {
	return b.cache.characterClass == ogame.Discoverer
}

func (b *OGame) getUniverseSpeed() int64 {
	return b.cache.serverData.Speed
}

func (b *OGame) getUniverseSpeedFleet() int64 {
	return b.cache.serverData.SpeedFleet
}

func (b *OGame) isDonutGalaxy() bool {
	return b.cache.serverData.DonutGalaxy
}

func (b *OGame) isDonutSystem() bool {
	return b.cache.serverData.DonutSystem
}

func (b *OGame) fetchEventbox() (res eventboxResp, err error) {
	err = b.getPageJSON(url.Values{"page": {FetchEventboxAjaxPageName}}, &res)
	return
}

func (b *OGame) isUnderAttack(opts ...Option) (bool, error) {
	attacks, err := b.getAttacks(opts...)
	return len(attacks) > 0, err
}

func (b *OGame) setVacationMode() error {
	vals := url.Values{"page": {"ingame"}, "component": {"preferences"}}
	pageHTML, err := b.getPageContent(vals)
	if err != nil {
		return err
	}
	rgx := regexp.MustCompile(`type='hidden' name='token' value='(\w+)'`)
	m := rgx.FindSubmatch(pageHTML)
	if len(m) < 2 {
		return errors.New("unable to find token")
	}
	token := string(m[1])
	payload := url.Values{"mode": {"save"}, "selectedTab": {"0"}, "urlaubs_modus": {"on"}, "token": {token}}
	_, err = b.postPageContent(vals, payload)
	return err
}

func (b *OGame) setPreferencesLang(lang string) error {
	vals := url.Values{"page": {"ingame"}, "component": {"preferences"}}
	pageHTML, err := b.getPageContent(vals)
	if err != nil {
		return err
	}
	rgx := regexp.MustCompile(`type='hidden' name='token' value='(\w+)'`)
	m := rgx.FindSubmatch(pageHTML)
	if len(m) < 2 {
		return errors.New("unable to find token")
	}
	token := string(m[1])
	payload := url.Values{
		"mode":        {"save"},
		"selectedTab": {"0"},
		"token":       {token},
	}
	payload.Set("language", lang)
	_, err = b.postPageContent(vals, payload)
	return err
}

func (b *OGame) setPreferences(p ogame.Preferences) error {
	vals := url.Values{"page": {"ingame"}, "component": {"preferences"}}
	pageHTML, err := b.getPageContent(vals)
	if err != nil {
		return err
	}
	rgx := regexp.MustCompile(`type='hidden' name='token' value='(\w+)'`)
	m := rgx.FindSubmatch(pageHTML)
	if len(m) < 2 {
		return errors.New("unable to find token")
	}
	token := string(m[1])
	payload := url.Values{
		"mode":        {"save"},
		"selectedTab": {"0"},
		"token":       {token},
	}

	setValue := func(ok bool, k, v string) {
		if ok {
			payload.Set(k, v)
		}
	}

	setValue(p.ShowOldDropDowns, "showOldDropDowns", "on")
	setValue(p.SpioReportPictures, "spioReportPictures", "on")
	setValue(p.ActivateAutofocus, "activateAutofocus", "on")
	setValue(p.ShowDetailOverlay, "showDetailOverlay", "on")
	setValue(p.AnimatedSliders, "animatedSliders", "on")
	setValue(p.AnimatedOverview, "animatedOverview", "on")
	setValue(p.PopupsNotices, "popups[notices]", "on")
	setValue(p.PopopsCombatreport, "popups[combatreport]", "on")
	setValue(p.AuctioneerNotifications, "auctioneerNotifications", "on")
	setValue(p.EconomyNotifications, "economyNotifications", "on")
	setValue(p.ShowActivityMinutes, "showActivityMinutes", "1")
	setValue(p.PreserveSystemOnPlanetChange, "preserveSystemOnPlanetChange", "1")
	setValue(p.DisableOutlawWarning, "disableOutlawWarning", "on")
	setValue(p.DiscoveryWarningEnabled, "discoveryWarningEnabled", "1")
	payload.Set("msgResultsPerPage", utils.FI64(p.MsgResultsPerPage))
	payload.Set("spySystemAutomaticQuantity", utils.FI64(p.SpySystemAutomaticQuantity))
	payload.Set("spySystemTargetPlanetTypes", utils.FI64(p.SpySystemTargetPlanetTypes))
	payload.Set("spySystemTargetPlayerTypes", utils.FI64(p.SpySystemTargetPlayerTypes))
	payload.Set("spySystemIgnoreSpiedInLastXMinutes", utils.FI64(p.SpySystemIgnoreSpiedInLastXMinutes))
	payload.Set("settings_sort", utils.FI64(p.SortSetting))
	payload.Set("settings_order", utils.FI64(p.SortOrder))
	payload.Set("spio_anz", utils.FI64(p.SpioAnz))
	payload.Set("eventsShow", utils.FI64(p.EventsShow))
	payload.Set("language", p.Language)

	_, err = b.postPageContent(vals, payload)
	return err
}

func (b *OGame) getPlanets() (out []Planet, err error) {
	page, err := getPage[parser.OverviewPage](b)
	if err != nil {
		return
	}
	return convertPlanets(b, page.ExtractPlanets()), nil
}

func (b *OGame) getPlanet(v IntoCelestial) (Planet, error) {
	page, err := getPage[parser.OverviewPage](b)
	if err != nil {
		return Planet{}, err
	}
	planet, err := page.ExtractPlanet(v)
	if err != nil {
		return Planet{}, err
	}
	return convertPlanet(b, planet), nil
}

func (b *OGame) getMoons() (out []Moon, err error) {
	page, err := getPage[parser.OverviewPage](b)
	if err != nil {
		return
	}
	return convertMoons(b, page.ExtractMoons()), nil
}

func (b *OGame) getMoon(v IntoCelestial) (Moon, error) {
	page, err := getPage[parser.OverviewPage](b)
	if err != nil {
		return Moon{}, err
	}
	moon, err := page.ExtractMoon(v)
	if err != nil {
		return Moon{}, err
	}
	cMoon := convertMoon(b, moon)
	return cMoon, nil
}

func (b *OGame) getCelestials() ([]Celestial, error) {
	page, err := getPage[parser.OverviewPage](b)
	if err != nil {
		return nil, err
	}
	celestials, err := page.ExtractCelestials()
	if err != nil {
		return nil, err
	}
	return convertCelestials(b, celestials), nil
}

func (b *OGame) getCelestial(v IntoCelestial) (Celestial, error) {
	page, err := getPage[parser.OverviewPage](b)
	if err != nil {
		return nil, err
	}
	celestial, err := page.ExtractCelestial(v)
	if err != nil {
		return nil, err
	}
	return convertCelestial(b, celestial), nil
}

func (b *OGame) recruitOfficer(typ, days int64) error {
	if typ != 2 && typ != 3 && typ != 4 && typ != 5 && typ != 6 {
		return errors.New("invalid officer type")
	}
	if days != 7 && days != 90 {
		return errors.New("invalid days")
	}
	pageHTML, err := b.getPageContent(url.Values{"page": {"premium"}, "ajax": {"1"}, "type": {utils.FI64(typ)}})
	if err != nil {
		return err
	}
	token, err := b.extractor.ExtractPremiumToken(pageHTML, days)
	if err != nil {
		return err
	}
	if _, err := b.getPageContent(url.Values{"page": {"premium"}, "buynow": {"1"},
		"type": {utils.FI64(typ)}, "days": {utils.FI64(days)},
		"token": {token}}); err != nil {
		return err
	}
	return nil
}

func (b *OGame) abandon(v IntoPlanet) error {
	page, err := getPage[parser.OverviewPage](b)
	if err != nil {
		return err
	}
	planet, err := page.ExtractPlanet(v)
	if err != nil {
		return errors.New("invalid parameter")
	}
	pageHTML, _ := b.getPage(PlanetlayerPageName, ChangePlanet(planet.GetID()))
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return err
	}
	abandonToken, token := b.extractor.ExtractAbandonInformation(doc)
	payload := url.Values{
		"abandon":  {abandonToken},
		"token":    {token},
		"password": {b.password},
	}
	pageHTML, err = b.postPageContent(url.Values{
		"page":      {"ingame"},
		"component": {"overview"},
		"action":    {"planetGiveup"},
		"ajax":      {"1"},
		"asJson":    {"1"},
	}, payload)
	var res struct {
		NewAjaxToken string `json:"newAjaxToken"`
	}
	_ = json.Unmarshal(pageHTML, &res)
	return err
}

func (b *OGame) serverTime() (out time.Time, err error) {
	page, err := getPage[parser.OverviewPage](b)
	if err != nil {
		return
	}
	return page.ExtractServerTime()
}

func (b *OGame) getUserInfos() (out ogame.UserInfos, err error) {
	page, err := getPage[parser.OverviewPage](b)
	if err != nil {
		return
	}
	return page.ExtractUserInfos()
}

// ChatPostResp ...
type ChatPostResp struct {
	Status   string `json:"status"`
	ID       int    `json:"id"`
	SenderID int    `json:"senderId"`
	TargetID int    `json:"targetId"`
	Text     string `json:"text"`
	Date     int64  `json:"date"`
	NewToken string `json:"newToken"`
}

func (b *OGame) sendMessage(id int64, message string, isPlayer bool) error {
	payload := url.Values{
		"text":  {message + "\n"},
		"ajax":  {"1"},
		"token": {b.cache.ajaxChatToken},
	}
	if isPlayer {
		payload.Set("playerId", utils.FI64(id))
		payload.Set("mode", "1")
	} else {
		payload.Set("associationId", utils.FI64(id))
		payload.Set("mode", "3")
	}
	bodyBytes, err := b.postPageContent(url.Values{"page": {"ajaxChat"}}, payload)
	if err != nil {
		return err
	}
	if strings.Contains(string(bodyBytes), "INVALID_PARAMETERS") {
		return errors.New("invalid parameters")
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(bodyBytes)))
	if err != nil {
		return err
	}
	if doc.Find("title").Text() == "OGame Lobby" {
		return ogame.ErrNotLogged
	}
	var res ChatPostResp
	if err := json.Unmarshal(bodyBytes, &res); err != nil {
		return err
	}
	b.cache.ajaxChatToken = res.NewToken
	return nil
}

func (b *OGame) getFleetsFromEventList() ([]ogame.Fleet, error) {
	pageHTML, err := b.getPageContent(url.Values{"eventList": {"movement"}, "ajax": {"1"}})
	if err != nil {
		return nil, err
	}
	return b.extractor.ExtractFleetsFromEventList(pageHTML)
}

func (b *OGame) getFleets(opts ...Option) ([]ogame.Fleet, ogame.Slots, error) {
	page, err := getPage[parser.MovementPage](b, opts...)
	if err != nil {
		return []ogame.Fleet{}, ogame.Slots{}, err
	}
	fleets, err := page.ExtractFleets()
	if err != nil {
		return []ogame.Fleet{}, ogame.Slots{}, err
	}
	slots, err := page.ExtractSlots()
	if err != nil {
		return []ogame.Fleet{}, ogame.Slots{}, err
	}
	return fleets, slots, nil
}

func (b *OGame) cancelFleet(fleetID ogame.FleetID) error {
	page, err := getPage[parser.MovementPage](b)
	if err != nil {
		return err
	}
	token, err := page.ExtractCancelFleetToken(fleetID)
	if err != nil {
		return err
	}
	if _, err = b.getPageContent(url.Values{"page": {"ingame"}, "component": {"movement"}, "return": {fleetID.String()}, "token": {token}}); err != nil {
		return err
	}
	return nil
}

func (b *OGame) getLastFleetFor(origin, destination ogame.Coordinate, mission ogame.MissionID) (ogame.Fleet, error) {
	page, err := getPage[parser.MovementPage](b)
	if err != nil {
		return ogame.Fleet{}, err
	}
	fleets, err := page.ExtractFleets()
	if err != nil {
		return ogame.Fleet{}, err
	}
	return getLastFleetFor(fleets, origin, destination, mission)
}

func getLastFleetFor(fleets []ogame.Fleet, origin, destination ogame.Coordinate, mission ogame.MissionID) (maxV ogame.Fleet, err error) {
	if len(fleets) > 0 {
		for i, fleet := range fleets {
			if fleet.ID > maxV.ID &&
				fleet.Origin.Equal(origin) &&
				fleet.Destination.Equal(destination) &&
				fleet.Mission == mission &&
				!fleet.ReturnFlight {
				maxV = fleets[i]
			}
		}
		if maxV.ID > 0 {
			return maxV, nil
		}
	}
	return maxV, errors.New("could not find fleet")
}

func (b *OGame) getFleetDispatch(celestialID ogame.CelestialID, options ...Option) (out ogame.FleetDispatchInfos, err error) {
	options = append(options, ChangePlanet(celestialID))
	page, err := getPage[parser.FleetDispatchPage](b, options...)
	if err != nil {
		return
	}
	ships, err := page.ExtractShips()
	if err != nil {
		return
	}
	slots, err := page.ExtractSlots()
	if err != nil {
		return
	}
	acsValues := page.ExtractAcsValues()
	out.Resources = page.ExtractResources()
	out.Ships = ships
	out.Slots = slots
	out.ACSValues = acsValues
	return
}

func (b *OGame) getSlots() (out ogame.Slots, err error) {
	pageHTML, err := b.getPage(FleetdispatchPageName)
	if err != nil {
		return
	}
	return b.extractor.ExtractSlots(pageHTML)
}

// Returns the distance between two galaxy
func galaxyDistance(galaxy1, galaxy2, universeSize int64, donutGalaxy bool) int64 {
	distance := math.Abs(float64(galaxy2 - galaxy1))
	if donutGalaxy {
		distance = min(distance, float64(universeSize)-distance)
	}
	return int64(20_000 * distance)
}

func systemDistance(nbSystems, system1, system2 int64, donutSystem bool) int64 {
	distance := math.Abs(float64(system2 - system1))
	if donutSystem {
		distance = min(distance, float64(nbSystems)-distance)
	}
	return int64(distance)
}

// Returns the distance between two systems
// https://ogame.fandom.com/wiki/Distance
func flightSystemDistance(nbSystems, system1, system2, systemsSkip int64, donutSystem bool) (distance int64) {
	dist := max(systemDistance(nbSystems, system1, system2, donutSystem)-systemsSkip, 0)
	return 2_700 + 95*dist
}

// Returns the distance between two planets
func planetDistance(planet1, planet2 int64) (distance int64) {
	return int64(1_000 + 5*math.Abs(float64(planet2-planet1)))
}

// Distance returns the distance between two coordinates
func Distance(c1, c2 ogame.Coordinate, universeSize, nbSystems, systemsSkip int64, donutGalaxy, donutSystem bool) (distance int64) {
	if c1.Galaxy != c2.Galaxy {
		return galaxyDistance(c1.Galaxy, c2.Galaxy, universeSize, donutGalaxy)
	}
	if c1.System != c2.System {
		return flightSystemDistance(nbSystems, c1.System, c2.System, systemsSkip, donutSystem)
	}
	if c1.Position != c2.Position {
		return planetDistance(c1.Position, c2.Position)
	}
	return 5
}

func calcFuel(ships ogame.ShipsInfos, dist, duration, holdingTime int64, universeSpeedFleet, fleetDeutSaveFactor float64, techs ogame.Researches,
	lfBonuses ogame.LfBonuses, characterClass ogame.CharacterClass, allianceClass ogame.AllianceClass) (fuel int64) {
	var holdingCosts int64
	tmpFuel := 0.0
	for shipID, nb := range ships.IterFlyable() {
		ship := ogame.Objs.GetShip(shipID)
		getFuelConsumption := ship.GetFuelConsumption(techs, lfBonuses, characterClass, fleetDeutSaveFactor)
		speed := ship.GetSpeed(techs, lfBonuses, characterClass, allianceClass)
		holdingCosts += getFuelConsumption * nb * holdingTime
		shipSpeedValue := (35_000 / max(0.5, float64(duration)*universeSpeedFleet-10)) * math.Sqrt(float64(dist)*10/float64(speed))
		tmpFuel += max(float64(getFuelConsumption*nb*dist)/35_000*math.Pow(shipSpeedValue/10+1, 2), 1)
	}
	fuel = int64(1 + math.Round(tmpFuel))
	if holdingTime > 0 {
		fuel += max(int64(float64(holdingCosts)/10.0), 1)
	}
	return
}

// CalcFlightTime ...
// Systems that are empty/inactive can be skipped for distance calculation
// (server settings: fleetIgnoreEmptySystems, fleetIgnoreInactiveSystems)
// https://board.en.ogame.gameforge.com/index.php?thread/838751-flight-time-consumption-ignores-empty-inactive-systems
// speed: 1 -> 100% | 0.5 -> 50% | 0.05 -> 5%
func CalcFlightTime(origin, destination ogame.Coordinate, universeSize, nbSystems int64, donutGalaxy, donutSystem bool,
	fleetDeutSaveFactor, speed float64, universeSpeedFleet int64, ships ogame.ShipsInfos, techs ogame.Researches, lfBonuses ogame.LfBonuses,
	characterClass ogame.CharacterClass, allianceClass ogame.AllianceClass, systemsSkip, holdingTime int64) (secs, fuel int64) {
	if !ships.HasShips() {
		return
	}
	v := ships.Speed(techs, lfBonuses, characterClass, allianceClass)
	secs = CalcFlightTimeWithBaseSpeed(origin, destination, universeSize, nbSystems, donutGalaxy, donutSystem, speed, v, universeSpeedFleet, systemsSkip)
	d := float64(Distance(origin, destination, universeSize, nbSystems, systemsSkip, donutGalaxy, donutSystem))
	fuel = calcFuel(ships, int64(d), secs, holdingTime, float64(universeSpeedFleet), fleetDeutSaveFactor, techs, lfBonuses, characterClass, allianceClass)
	return
}

// CalcFlightTimeWithBaseSpeed ...
// baseSpeed is the speed of the slowest ship in a fleet
// speed: 1 -> 100% | 0.5 -> 50% | 0.05 -> 5%
// https://ogame.fandom.com/wiki/Distance
func CalcFlightTimeWithBaseSpeed(origin, destination ogame.Coordinate, universeSize, nbSystems int64, donutGalaxy, donutSystem bool, speed float64, baseSpeed, universeSpeedFleet, systemsSkip int64) (secs int64) {
	d := Distance(origin, destination, universeSize, nbSystems, systemsSkip, donutGalaxy, donutSystem)
	return CalcFlightTimeWithBaseSpeedDistance(speed, baseSpeed, d, universeSpeedFleet)
}

// CalcFlightTimeWithBaseSpeedDistance ...
func CalcFlightTimeWithBaseSpeedDistance(speed float64, baseSpeed, distance, universeSpeedFleet int64) (secs int64) {
	s := speed
	v := float64(baseSpeed)
	a := float64(universeSpeedFleet)
	d := float64(distance)
	return int64(math.Round(((3_500/s)*math.Sqrt(d*10/v) + 10) / a))
}

// CalcFlightTime calculates the flight time and the fuel consumption
func (b *OGame) CalcFlightTime(origin, destination ogame.Coordinate, speed float64, ships ogame.ShipsInfos, missionID ogame.MissionID, holdingTime int64) (secs, fuel int64) {
	serverData := b.cache.serverData
	lfBonuses, _ := b.GetCachedLfBonuses()
	allianceClass, _ := b.GetCachedAllianceClass()
	fleetIgnoreEmptySystems := b.cache.serverData.FleetIgnoreEmptySystems
	fleetIgnoreInactiveSystems := b.cache.serverData.FleetIgnoreInactiveSystems
	var systemsSkip int64
	if fleetIgnoreEmptySystems || fleetIgnoreInactiveSystems {
		opts := make([]Option, 0)
		if originCelestial, err := b.GetCachedCelestial(origin); err == nil {
			opts = append(opts, ChangePlanet(originCelestial.GetID()))
		}
		res, _ := b.CheckTarget(ships, destination, opts...)
		if fleetIgnoreEmptySystems {
			systemsSkip += res.EmptySystems
		}
		if fleetIgnoreInactiveSystems {
			systemsSkip += res.InactiveSystems
		}
	}
	return CalcFlightTime(origin, destination, serverData.Galaxies, serverData.Systems, serverData.DonutGalaxy,
		serverData.DonutSystem, serverData.GlobalDeuteriumSaveFactor, speed, GetFleetSpeedForMission(serverData, missionID), ships,
		b.GetCachedResearch(), lfBonuses, b.cache.characterClass, allianceClass, systemsSkip, holdingTime)
}

// getPhalanx makes 3 calls to ogame server (2 validation, 1 scan)
func (b *OGame) getPhalanx(moonID ogame.MoonID, coord ogame.Coordinate) ([]ogame.PhalanxFleet, error) {
	res := make([]ogame.PhalanxFleet, 0)

	// Get moon facilities html page (first call to ogame server)
	moonFacilitiesHTML, _ := b.getPage(FacilitiesPageName, ChangePlanet(moonID.Celestial()))

	// Extract a bunch of infos from the html
	moon, err := b.extractor.ExtractMoon(moonFacilitiesHTML, moonID)
	if err != nil {
		return res, errors.New("moon not found")
	}
	resources, err := b.extractor.ExtractResources(moonFacilitiesHTML)
	if err != nil {
		return res, err
	}
	moonFacilities, _ := b.extractor.ExtractFacilities(moonFacilitiesHTML)
	phalanxLvl := moonFacilities.SensorPhalanx

	if phalanxLvl == 0 {
		return res, errors.New("no sensor phalanx on this moon")
	}

	// Ensure we have the resources to scan the planet
	if resources.Deuterium < ogame.SensorPhalanx.ScanConsumption() {
		return res, ogame.ErrNotEnoughDeuterium
	}

	// Verify that coordinate is in phalanx range
	phalanxRange := ogame.SensorPhalanx.GetRange(phalanxLvl, b.isDiscoverer())
	if moon.GetCoordinate().Galaxy != coord.Galaxy ||
		systemDistance(b.cache.serverData.Systems, moon.GetCoordinate().System, coord.System, b.cache.serverData.DonutSystem) > phalanxRange {
		return res, errors.New("coordinate not in phalanx range")
	}

	// Get galaxy planets information, verify coordinate is valid planet (call to ogame server)
	planetInfos, _ := b.galaxyInfos(coord.Galaxy, coord.System)
	target := planetInfos.Position(coord.Position)
	if target == nil {
		return nil, errors.New("invalid planet coordinate")
	}
	// Ensure you are not scanning your own planet
	if target.Player.ID == b.cache.player.PlayerID {
		return nil, errors.New("cannot scan own planet")
	}

	// Run the phalanx scan (second & third calls to ogame server)
	return b.getUnsafePhalanx(moonID, coord)
}

// getUnsafePhalanx ...
func (b *OGame) getUnsafePhalanx(moonID ogame.MoonID, coord ogame.Coordinate) ([]ogame.PhalanxFleet, error) {
	vals := url.Values{
		"page":     {PhalanxAjaxPageName},
		"galaxy":   {utils.FI64(coord.Galaxy)},
		"system":   {utils.FI64(coord.System)},
		"position": {utils.FI64(coord.Position)},
		"ajax":     {"1"},
		"token":    {b.cache.token},
	}
	page, err := getAjaxPage[parser.PhalanxAjaxPage](b, vals, ChangePlanet(moonID.Celestial()))
	if err != nil {
		return []ogame.PhalanxFleet{}, err
	}
	b.cache.token, _ = page.ExtractPhalanxNewToken()
	return page.ExtractPhalanx()
}

func (b *OGame) headersForPage(url string) (http.Header, error) {
	if !b.IsEnabled() {
		return nil, ogame.ErrBotInactive
	}
	if !b.IsLoggedIn() {
		return nil, ogame.ErrBotLoggedOut
	}

	if b.cache.serverURL == "" {
		err := errors.New("serverURL is empty")
		b.error(err)
		return nil, err
	}

	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}

	finalURL := b.cache.serverURL + url

	req, err := http.NewRequest(http.MethodHead, finalURL, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(b.ctx)
	resp, err := b.device.GetClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode >= http.StatusInternalServerError {
		return nil, fmt.Errorf("gameforme internal server error : %s", resp.Status)
	}

	return resp.Header, err
}

func (b *OGame) getJumpGatePage(originMoonID ogame.MoonID) (parser.JumpGateAjaxPage, error) {
	vals := url.Values{"page": {"ajax"}, "component": {"jumpgate"}, "overlay": {"1"}, "ajax": {"1"}}
	return getAjaxPage[parser.JumpGateAjaxPage](b, vals, ChangePlanet(originMoonID.Celestial()))
}

func (b *OGame) jumpGateDestinations(originMoonID ogame.MoonID) ([]ogame.MoonID, int64, error) {
	page, err := b.getJumpGatePage(originMoonID)
	if err != nil {
		return nil, 0, err
	}
	_, _, dests, wait, err := page.ExtractJumpGate()
	if err != nil {
		return nil, 0, err
	}
	if wait > 0 {
		return dests, wait, fmt.Errorf("jump gate is in recharge mode for %d seconds", wait)
	}
	return dests, wait, nil
}

func (b *OGame) executeJumpGate(originMoonID, destMoonID ogame.MoonID, ships ogame.ShipsInfos) (bool, int64, error) {
	page, err := b.getJumpGatePage(originMoonID)
	if err != nil {
		return false, 0, err
	}
	availShips, token, dests, wait, err := page.ExtractJumpGate()
	if err != nil {
		return false, 0, err
	}
	if wait > 0 {
		return false, wait, fmt.Errorf("jump gate is in recharge mode for %d seconds", wait)
	}

	// Validate destination moon id
	if !utils.InArr(destMoonID, dests) {
		return false, 0, errors.New("destination moon id invalid")
	}

	payload := url.Values{"token": {token}, "targetSpaceObjectId": {utils.FI64(destMoonID)}}

	// Add ships to payload
	for _, s := range ogame.Ships {
		// Get the min between what is available and what we want
		nbr := min(ships.ByID(s.GetID()), availShips.ByID(s.GetID()))
		if nbr > 0 {
			payload.Add("ship_"+utils.FI64(s.GetID()), utils.FI64(nbr))
		}
	}

	if _, err := b.postPageContent(url.Values{"page": {"componentOnly"}, "component": {"jumpgate"}, "action": {"executeJump"}, "asJson": {"1"}}, payload); err != nil {
		return false, 0, err
	}
	return true, 0, nil
}

func (b *OGame) getEmpireHtml(celestialType ogame.CelestialType) ([]byte, error) {
	var planetType int
	switch celestialType {
	case ogame.PlanetType:
		planetType = 0
	case ogame.MoonType:
		planetType = 1
	default:
		return nil, errors.New("invalid celestial type")
	}
	vals := url.Values{"page": {"standalone"}, "component": {"empire"}, "planetType": {strconv.Itoa(planetType)}}
	pageHTMLBytes, err := b.getPageContent(vals)
	if err != nil {
		return nil, err
	}
	return pageHTMLBytes, nil
}

func (b *OGame) getEmpire(celestialType ogame.CelestialType) (out []ogame.EmpireCelestial, err error) {
	pageHTMLBytes, err := b.getEmpireHtml(celestialType)
	if err != nil {
		return out, err
	}
	return b.extractor.ExtractEmpire(pageHTMLBytes)
}

func (b *OGame) getEmpireJSON(celestialType ogame.CelestialType) (any, error) {
	// Valid URLs:
	// /game/index.php?page=standalone&component=empire&planetType=0
	// /game/index.php?page=standalone&component=empire&planetType=1
	pageHTMLBytes, err := b.getEmpireHtml(celestialType)
	if err != nil {
		return nil, err
	}
	// Replace the Ogame hostname with our custom hostname
	pageHTML := strings.Replace(string(pageHTMLBytes), b.cache.serverURL, b.apiNewHostname, -1)
	return b.extractor.ExtractEmpireJSON([]byte(pageHTML))
}

func (b *OGame) createUnion(fleet ogame.Fleet, unionUsers []string) (int64, error) {
	if fleet.ID == 0 {
		return 0, errors.New("invalid fleet id")
	}
	pageHTML, err := b.getPageContent(url.Values{"page": {"federationlayer"}, "union": {"0"}, "fleet": {utils.FI64(fleet.ID)}, "target": {utils.FI64(fleet.TargetPlanetID)}, "ajax": {"1"}})
	if err != nil {
		return 0, err
	}
	payload, err := b.extractor.ExtractFederation(pageHTML)
	if err != nil {
		return 0, err
	}

	payloadUnionUsers := payload["unionUsers"]
	for _, user := range payloadUnionUsers {
		if user != "" {
			unionUsers = append(unionUsers, user)
		}
	}
	payload.Set("unionUsers", strings.Join(unionUsers, ";"))

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

func (b *OGame) highscore(category, typ, page int64) (out ogame.Highscore, err error) {
	if category < 1 || category > 2 {
		return out, errors.New("category must be in [1, 2] (1:player, 2:alliance)")
	}
	if typ < 0 || typ > 11 {
		msg := "typ must be in [0, 11] (0:Total, 1:Economy, 2:Research, 3:Military, 4:Military Built, " +
			"5:Military Destroyed, 6:Military Lost, 7:Honor, 8:Lifeform, 9:Lifeform Economy, 10:Lifeform Technology, " +
			"11:Lifeform Discoveries)"
		return out, errors.New(msg)
	}
	if page < 1 {
		return out, errors.New("page must be greater than or equal to 1")
	}
	vals := url.Values{
		"page":     {HighscoreContentAjaxPageName},
		"category": {utils.FI64(category)},
		"type":     {utils.FI64(typ)},
		"site":     {utils.FI64(page)},
	}
	payload := url.Values{}
	pageHTML, err := b.postPageContent(vals, payload)
	if err != nil {
		return out, err
	}
	return b.extractor.ExtractHighscore(pageHTML)
}

func (b *OGame) getAllResources() (map[ogame.CelestialID]ogame.Resources, error) {
	vals := url.Values{
		"page":      {"ajax"},
		"component": {"traderauctioneer"},
	}
	payload := url.Values{
		"show": {"auctioneer"},
		"ajax": {"1"},
	}
	pageHTML, err := b.postPageContent(vals, payload)
	if err != nil {
		return nil, err
	}
	return b.extractor.ExtractAllResources(pageHTML)
}

func (b *OGame) getDMCosts(celestialID ogame.CelestialID) (ogame.DMCosts, error) {
	page, err := getPage[parser.OverviewPage](b, ChangePlanet(celestialID))
	if err != nil {
		return ogame.DMCosts{}, err
	}
	return page.ExtractDMCosts()
}

func (b *OGame) useDM(typ ogame.DMType, celestialID ogame.CelestialID) error {
	if !typ.IsValid() {
		return fmt.Errorf("invalid type %s", typ)
	}
	page, err := getPage[parser.OverviewPage](b, ChangePlanet(celestialID))
	if err != nil {
		return err
	}
	costs, err := page.ExtractDMCosts()
	if err != nil {
		return err
	}
	var buyAndActivate, token string
	switch typ {
	case ogame.BuildingsDmType:
		buyAndActivate, token = costs.Buildings.BuyAndActivateToken, costs.Buildings.Token
	case ogame.ResearchDmType:
		buyAndActivate, token = costs.Research.BuyAndActivateToken, costs.Research.Token
	case ogame.ShipyardDmType:
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

// marketItemType 3 -> offer buy
// marketItemType 4 -> offer sell
// itemID 1 -> metal
// itemID 2 -> crystal
// itemID 3 -> deuterium
// itemID 204 -> light fighter
// itemID <HASH> -> item
func (b *OGame) offerMarketplace(marketItemType int64, itemID any, quantity, priceType, price, priceRange int64, celestialID ogame.CelestialID) error {
	params := url.Values{"page": {"ingame"}, "component": {"marketplace"}, "tab": {"create_offer"}, "action": {"submitOffer"}, "asJson": {"1"}}
	const (
		shipsItemType = iota + 1
		resourcesItemType
		itemItemType
	)
	var itemIDPayload string
	var itemType int64
	if itemIDStr, ok := itemID.(string); ok {
		if len(itemIDStr) == 40 {
			itemType = itemItemType
			itemIDPayload = itemIDStr
		} else {
			return errors.New("invalid itemID string")
		}
	} else if itemIDInt64, ok := itemID.(int64); ok {
		if itemIDInt64 >= 1 && itemIDInt64 <= 3 {
			itemType = resourcesItemType
			itemIDPayload = utils.FI64(itemIDInt64)
		} else if ogame.ID(itemIDInt64).IsShip() {
			itemType = shipsItemType
			itemIDPayload = utils.FI64(itemIDInt64)
		} else {
			return errors.New("invalid itemID int64")
		}
	} else if itemIDInt, ok := itemID.(int); ok {
		if itemIDInt >= 1 && itemIDInt <= 3 {
			itemType = resourcesItemType
			itemIDPayload = strconv.Itoa(itemIDInt)
		} else if ogame.ID(itemIDInt).IsShip() {
			itemType = shipsItemType
			itemIDPayload = strconv.Itoa(itemIDInt)
		} else {
			return errors.New("invalid itemID int")
		}
	} else if itemIDID, ok := itemID.(ogame.ID); ok {
		if itemIDID.IsShip() {
			itemType = shipsItemType
			itemIDPayload = utils.FI64(itemIDID)
		} else {
			return errors.New("invalid itemID ID")
		}
	} else {
		return errors.New("invalid itemID type")
	}

	vals := url.Values{
		"page":      {"ingame"},
		"component": {"marketplace"},
		"tab":       {"create_offer"},
	}
	pageHTML, err := b.getPageContent(vals)
	if err != nil {
		return err
	}
	token, _ := b.extractor.ExtractToken(pageHTML)

	payload := url.Values{
		"marketItemType": {utils.FI64(marketItemType)},
		"itemType":       {utils.FI64(itemType)},
		"itemId":         {itemIDPayload},
		"quantity":       {utils.FI64(quantity)},
		"priceType":      {utils.FI64(priceType)},
		"price":          {utils.FI64(price)},
		"priceRange":     {utils.FI64(priceRange)},
		"token":          {token},
	}
	var res struct {
		Status  string       `json:"status"`
		Message string       `json:"message"`
		Errors  []OGameError `json:"errors"`
	}
	by, err := b.postPageContent(params, payload, ChangePlanet(celestialID))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(by, &res); err != nil {
		return err
	}
	if len(res.Errors) > 0 {
		return errors.New(utils.FI64(res.Errors[0].Error) + " : " + res.Errors[0].Message)
	}
	return err
}

func (b *OGame) buyMarketplace(itemID int64, celestialID ogame.CelestialID) (err error) {
	params := url.Values{"page": {"ingame"}, "component": {"marketplace"}, "tab": {"buying"}, "action": {"acceptRequest"}, "asJson": {"1"}}
	payload := url.Values{
		"marketItemId": {utils.FI64(itemID)},
	}
	var res struct {
		Status  string       `json:"status"`
		Message string       `json:"message"`
		Errors  []OGameError `json:"errors"`
	}
	by, err := b.postPageContent(params, payload, ChangePlanet(celestialID))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(by, &res); err != nil {
		return err
	}
	if len(res.Errors) > 0 {
		return errors.New(utils.FI64(res.Errors[0].Error) + " : " + res.Errors[0].Message)
	}
	return err
}

func (b *OGame) getItems(celestialID ogame.CelestialID) (items []ogame.Item, err error) {
	params := url.Values{"page": {"ajax"}, "component": {"buffactivation"}, "ajax": {"1"}, "type": {"1"}}
	pageHTML, err := b.getPageContent(params, ChangePlanet(celestialID))
	if err != nil {
		return nil, err
	}
	_, items, err = b.extractor.ExtractBuffActivation(pageHTML)
	return
}

func (b *OGame) getActiveItems(celestialID ogame.CelestialID) (items []ogame.ActiveItem, err error) {
	page, err := getPage[parser.OverviewPage](b, ChangePlanet(celestialID))
	if err != nil {
		return []ogame.ActiveItem{}, err
	}
	return page.ExtractActiveItems()
}

type MessageSuccess struct {
	Buff          string `json:"buff"`
	Status        string `json:"status"`
	Duration      int    `json:"duration"`
	Extendable    bool   `json:"extendable"`
	TotalDuration int    `json:"totalDuration"`
	Tooltip       string `json:"tooltip"`
	Reload        bool   `json:"reload"`
	BuffID        string `json:"buffId"`
	Item          struct {
		Name                    string   `json:"name"`
		Image                   string   `json:"image"`
		ImageLarge              string   `json:"imageLarge"`
		Title                   string   `json:"title"`
		Effect                  string   `json:"effect"`
		Ref                     string   `json:"ref"`
		Rarity                  string   `json:"rarity"`
		Amount                  int      `json:"amount"`
		AmountFree              int      `json:"amount_free"`
		AmountBought            int      `json:"amount_bought"`
		Category                []string `json:"category"`
		Currency                string   `json:"currency"`
		Costs                   string   `json:"costs"`
		IsReduced               bool     `json:"isReduced"`
		Buyable                 bool     `json:"buyable"`
		CanBeActivated          bool     `json:"canBeActivated"`
		CanBeBoughtAndActivated bool     `json:"canBeBoughtAndActivated"`
		IsAnUpgrade             bool     `json:"isAnUpgrade"`
		IsCharacterClassItem    bool     `json:"isCharacterClassItem"`
		HasEnoughCurrency       bool     `json:"hasEnoughCurrency"`
		Cooldown                int      `json:"cooldown"`
		Duration                int      `json:"duration"`
		DurationExtension       any      `json:"durationExtension"`
		TotalTime               int      `json:"totalTime"`
		TimeLeft                int      `json:"timeLeft"`
		Status                  string   `json:"status"`
		Extendable              bool     `json:"extendable"`
		FirstStatus             string   `json:"firstStatus"`
		ToolTip                 string   `json:"toolTip"`
		BuyTitle                string   `json:"buyTitle"`
		ActivationTitle         string   `json:"activationTitle"`
		MoonOnlyItem            bool     `json:"moonOnlyItem"`
	} `json:"item"`
	Message string `json:"message"`
}

func (b *OGame) activateItem(ref string, celestialID ogame.CelestialID) error {
	params := url.Values{
		"page":         {"componentOnly"},
		"component":    {"itemactions"},
		"asJson":       {"1"},
		"itemUuid":     {ref},
		"action":       {"activate"},
		"token":        {b.cache.token},
		"referrerPage": {"ingame"},
		"_":            {strconv.FormatInt(time.Now().UnixMilli(), 10)},
	}
	params.Add("itemUuid", ref)
	pageHTML, err := b.getPageContent(params, ChangePlanet(celestialID))
	if err != nil {
		return err
	}
	var responseStruct struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(pageHTML, &responseStruct); err != nil {
		return err
	}
	if responseStruct.Status == "failure" {
		return errors.New("failed to activate item")
	}
	params = url.Values{"page": {"ajax"}, "component": {"buffactivation"}, "ajax": {"1"}}
	payload := url.Values{"type": {ref}}
	if _, err := b.postPageContent(params, payload); err != nil {
		return err
	}
	return nil
}

func (b *OGame) getAuction(celestialID ogame.CelestialID) (ogame.Auction, error) {
	payload := url.Values{"show": {"auctioneer"}, "ajax": {"1"}}
	auctionHTML, err := b.postPageContent(url.Values{"page": {"ajax"}, "component": {"traderauctioneer"}}, payload, ChangePlanet(celestialID))
	if err != nil {
		return ogame.Auction{}, err
	}
	return b.extractor.ExtractAuction(auctionHTML)
}

func (b *OGame) doAuction(celestialID ogame.CelestialID, bid map[ogame.CelestialID]ogame.Resources) error {
	// Get fresh token (among others)
	auction, err := b.getAuction(celestialID)
	if err != nil {
		return err
	}

	if auction.HasFinished {
		return errors.New("auction completed")
	}

	payload := url.Values{}
	for auctionCelestialIDString := range auction.Resources {
		payload.Set("bid[planets]["+auctionCelestialIDString+"][metal]", "0")
		payload.Set("bid[planets]["+auctionCelestialIDString+"][crystal]", "0")
		payload.Set("bid[planets]["+auctionCelestialIDString+"][deuterium]", "0")
	}
	for celestialID, resources := range bid {
		payload.Set("bid[planets]["+utils.FI64(celestialID)+"][metal]", utils.FI64(resources.Metal))
		payload.Set("bid[planets]["+utils.FI64(celestialID)+"][crystal]", utils.FI64(resources.Crystal))
		payload.Set("bid[planets]["+utils.FI64(celestialID)+"][deuterium]", utils.FI64(resources.Deuterium))
	}

	payload.Add("bid[honor]", "0")
	payload.Add("token", auction.Token)
	payload.Add("ajax", "1")

	if celestialID != 0 {
		payload.Set("cp", utils.FI64(celestialID))
	}

	auctionHTML, err := b.postPageContent(url.Values{"page": {"ajax"}, "component": {"traderauctioneer"}, "ajax": {"1"}, "action": {"submitBid"}, "asJson": {"1"}}, payload)
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

	var jsonObj map[string]any
	if err := json.Unmarshal(auctionHTML, &jsonObj); err != nil {
		return err
	}
	if jsonObj["error"] == true {
		return errors.New(jsonObj["message"].(string))
	}
	return nil
}

func calcResources(price int64, planetResources ogame.PlanetResources, multiplier ogame.Multiplier) url.Values {
	sortedCelestialIDs := make([]ogame.CelestialID, 0)
	for celestialID := range planetResources {
		sortedCelestialIDs = append(sortedCelestialIDs, celestialID)
	}
	sort.Slice(sortedCelestialIDs, func(i, j int) bool {
		return int64(sortedCelestialIDs[i]) < int64(sortedCelestialIDs[j])
	})

	payload := url.Values{}
	remaining := price
	multMetal := multiplier.Metal
	multCrystal := multiplier.Crystal
	multDeuterium := multiplier.Deuterium
	for celestialID, res := range planetResources {
		metalNeeded := res.Input.Metal
		if remaining < int64(float64(metalNeeded)*multMetal) {
			metalNeeded = int64(math.Ceil(float64(remaining) / multMetal))
		}
		remaining -= int64(float64(metalNeeded) * multMetal)

		crystalNeeded := res.Input.Crystal
		if remaining < int64(float64(crystalNeeded)*multCrystal) {
			crystalNeeded = int64(math.Ceil(float64(remaining) / multCrystal))
		}
		remaining -= int64(float64(crystalNeeded) * multCrystal)

		deuteriumNeeded := res.Input.Deuterium
		if remaining < int64(float64(deuteriumNeeded)*multDeuterium) {
			deuteriumNeeded = int64(math.Ceil(float64(remaining) / multDeuterium))
		}
		remaining -= int64(float64(deuteriumNeeded) * multDeuterium)

		payload.Add("bid[planets]["+utils.FI64(celestialID)+"][metal]", utils.FI64(metalNeeded))
		payload.Add("bid[planets]["+utils.FI64(celestialID)+"][crystal]", utils.FI64(crystalNeeded))
		payload.Add("bid[planets]["+utils.FI64(celestialID)+"][deuterium]", utils.FI64(deuteriumNeeded))
	}
	return payload
}

func (b *OGame) traderImportExportTrade(price int64, importToken string, planetResources ogame.PlanetResources, multiplier ogame.Multiplier) (string, error) {
	payload := calcResources(price, planetResources, multiplier)
	payload.Add("action", "trade")
	payload.Add("bid[honor]", "0")
	payload.Add("token", importToken)
	payload.Add("ajax", "1")
	pageHTML1, err := b.postPageContent(url.Values{"page": {"ajax"}, "component": {"traderimportexport"}, "ajax": {"1"}, "action": {"trade"}, "asJson": {"1"}}, payload)
	if err != nil {
		return "", err
	}
	// {"message":"You have bought a container.","error":false,"item":{"uuid":"40f6c78e11be01ad3389b7dccd6ab8efa9347f3c","itemText":"You have purchased 1 KRAKEN Bronze.","bargainText":"The contents of the container not appeal to you? For 500 Dark Matter you can exchange the container for another random container of the same quality. You can only carry out this exchange 2 times per daily offer.","bargainCost":500,"bargainCostText":"Costs: 500 Dark Matter","tooltip":"KRAKEN Bronze|Reduces the building time of buildings currently under construction by <b>30m<\/b>.<br \/><br \/>\nDuration: now<br \/><br \/>\nPrice: --- <br \/>\nIn Inventory: 1","image":"98629d11293c9f2703592ed0314d99f320f45845","amount":1,"rarity":"common"},"newToken":"07eefc14105db0f30cb331a8b7af0bfe"}
	var result struct {
		Message      string
		Error        bool
		NewAjaxToken string
	}
	if err := json.Unmarshal(pageHTML1, &result); err != nil {
		return "", err
	}
	if result.Error {
		return "", errors.New(result.Message)
	}
	return result.NewAjaxToken, nil
}

func (b *OGame) traderImportExportTakeItem(token string) error {
	payload := url.Values{"action": {"takeItem"}, "token": {token}, "ajax": {"1"}}
	pageHTML, err := b.postPageContent(url.Values{"page": {"ajax"}, "component": {"traderimportexport"}, "ajax": {"1"}, "action": {"takeItem"}, "asJson": {"1"}}, payload)
	if err != nil {
		return err
	}
	var result struct {
		Message      string
		Error        bool
		NewAjaxToken string
	}
	if err := json.Unmarshal(pageHTML, &result); err != nil {
		return err
	}
	if result.Error {
		return errors.New(result.Message)
	}
	// {"error":false,"message":"You have accepted the offer and put the item in your inventory.","item":{"name":"Bronze Deuterium Booster","image":"f0e514af79d0808e334e9b6b695bf864b861bdfa","imageLarge":"c7c2837a0b341d37383d6a9d8f8986f500db7bf9","title":"Bronze Deuterium Booster|+10% more Deuterium Synthesizer harvest on one planet<br \/><br \/>\nDuration: 1w<br \/><br \/>\nPrice: --- <br \/>\nIn Inventory: 134","effect":"+10% more Deuterium Synthesizer harvest on one planet","ref":"d9fa5f359e80ff4f4c97545d07c66dbadab1d1be","rarity":"common","amount":134,"amount_free":134,"amount_bought":0,"category":["d8d49c315fa620d9c7f1f19963970dea59a0e3be","e71139e15ee5b6f472e2c68a97aa4bae9c80e9da"],"currency":"dm","costs":"2500","isReduced":false,"buyable":false,"canBeActivated":true,"canBeBoughtAndActivated":false,"isAnUpgrade":false,"isCharacterClassItem":false,"hasEnoughCurrency":true,"cooldown":0,"duration":604800,"durationExtension":null,"totalTime":null,"timeLeft":null,"status":null,"extendable":false,"firstStatus":"effecting","toolTip":"Bronze Deuterium Booster|+10% more Deuterium Synthesizer harvest on one planet&lt;br \/&gt;&lt;br \/&gt;\nDuration: 1w&lt;br \/&gt;&lt;br \/&gt;\nPrice: --- &lt;br \/&gt;\nIn Inventory: 134","buyTitle":"This item is currently unavailable for purchase.","activationTitle":"Activate","moonOnlyItem":false,"newOffer":false,"noOfferMessage":"There are no further offers today. Please come again tomorrow."},"newToken":"dec779714b893be9b39c0bedf5738450","components":[],"newAjaxToken":"e20cf0a6ca0e9b43a81ccb8fe7e7e2e3"}
	return nil
}

func (b *OGame) buyOfferOfTheDay() error {
	pageHTML, err := b.postPageContent(url.Values{"page": {"ajax"}, "component": {"traderimportexport"}}, url.Values{"show": {"importexport"}, "ajax": {"1"}})
	if err != nil {
		return err
	}
	price, importToken, planetResources, multiplier, err := b.extractor.ExtractOfferOfTheDay(pageHTML)
	if err != nil {
		return err
	}
	newAjaxToken, err := b.traderImportExportTrade(price, importToken, planetResources, multiplier)
	if err != nil {
		return err
	}
	return b.traderImportExportTakeItem(newAjaxToken)
}

// Hack fix: When moon name is >12, the moon image disappear from the EventsBox
// and attacks are detected on planet instead.
func fixAttackEvents(attacks []ogame.AttackEvent, planets []Planet) {
	for i, attack := range attacks {
		if len(attack.DestinationName) > 12 {
			for _, planet := range planets {
				if attack.Destination.Equal(planet.Coordinate) &&
					planet.Moon != nil &&
					attack.DestinationName != planet.Name &&
					attack.DestinationName == planet.Moon.Name {
					attacks[i].Destination.Type = ogame.MoonType
				}
			}
		}
	}
}

func (b *OGame) getAttacks(opts ...Option) (out []ogame.AttackEvent, err error) {
	vals := url.Values{"page": {"componentOnly"}, "component": {EventListAjaxPageName}, "ajax": {"1"}}
	page, err := getAjaxPage[parser.EventListAjaxPage](b, vals, opts...)
	if err != nil {
		return
	}
	planets := b.getCachedPlanets()
	ownCoords := getOwnCoordinates(planets)
	out, err = page.ExtractAttacks(ownCoords)
	if err != nil {
		return
	}
	fixAttackEvents(out, planets)
	return
}

func getOwnCoordinates(planets []Planet) (ownCoords []ogame.Coordinate) {
	for _, planet := range planets {
		ownCoords = append(ownCoords, planet.Coordinate)
		if planet.Moon != nil {
			ownCoords = append(ownCoords, planet.Moon.Coordinate)
		}
	}
	return
}

func (b *OGame) galaxyInfos(galaxy, system int64, opts ...Option) (ogame.SystemInfos, error) {
	cfg := getOptions(opts...)
	var res ogame.SystemInfos
	if galaxy < 1 || galaxy > b.server.OGameSettings().UniverseSize {
		return res, fmt.Errorf("galaxy must be within [1, %d]", b.server.OGameSettings().UniverseSize)
	}
	if system < 1 || system > b.cache.serverData.Systems {
		return res, errors.New("system must be within [1, " + utils.FI64(b.cache.serverData.Systems) + "]")
	}
	payload := url.Values{
		"galaxy": {utils.FI64(galaxy)},
		"system": {utils.FI64(system)},
	}
	vals := url.Values{"page": {"ingame"}, "component": {"galaxy"}, "action": {"fetchGalaxyContent"}, "ajax": {"1"}, "asJson": {"1"}}
	pageHTML, err := b.postPageContent(vals, payload, opts...)
	if err != nil {
		return res, err
	}
	player := b.cache.player
	res, err = b.extractor.ExtractGalaxyInfos(pageHTML, player.PlayerName, player.PlayerID, player.Rank)
	if err != nil {
		if cfg.DebugGalaxy {
			fmt.Println(string(pageHTML))
		}
		return res, err
	}
	if res.Galaxy() != galaxy || res.System() != system {
		return ogame.SystemInfos{}, ogame.ErrNotEnoughDeuterium
	}
	return res, err
}

func (b *OGame) getGalaxyPage(galaxy int64, system int64, opts ...Option) (*GalaxyPageContent, error) {
	// Get galaxy page content for the desired system.
	by, err := b.postPageContent(url.Values{
		"page":      {"ingame"},
		"component": {"galaxy"},
		"action":    {"fetchGalaxyContent"},
		"ajax":      {"1"},
		"asJson":    {"1"},
	}, url.Values{
		"galaxy": {strconv.Itoa(int(galaxy))},
		"system": {strconv.Itoa(int(system))},
	}, opts...)
	if err != nil {
		return nil, err
	}
	// Parse the json result, only defining the type for the GalaxyContent (Position and AvailableMissions properties).
	var res GalaxyPageContent
	if err = json.Unmarshal(by, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (b *OGame) getResourceSettings(planetID ogame.PlanetID, options ...Option) (ogame.ResourceSettings, error) {
	options = append(options, ChangePlanet(planetID.Celestial()))
	pageHTML, _ := b.getPage(ResourceSettingsPageName, options...)
	settings, _, err := b.extractor.ExtractResourceSettings(pageHTML)
	return settings, err
}

func (b *OGame) setResourceSettings(planetID ogame.PlanetID, settings ogame.ResourceSettings) error {
	pageHTML, _ := b.getPage(ResourceSettingsPageName, ChangePlanet(planetID.Celestial()))
	_, token, err := b.extractor.ExtractResourceSettings(pageHTML)
	if err != nil {
		return err
	}
	payload := url.Values{
		"saveSettings": {"1"},
		"token":        {token},
		"last1":        {utils.FI64(settings.MetalMine)},
		"last2":        {utils.FI64(settings.CrystalMine)},
		"last3":        {utils.FI64(settings.DeuteriumSynthesizer)},
		"last4":        {utils.FI64(settings.SolarPlant)},
		"last12":       {utils.FI64(settings.FusionReactor)},
		"last212":      {utils.FI64(settings.SolarSatellite)},
		"last217":      {utils.FI64(settings.Crawler)},
	}
	url2 := b.cache.serverURL + "/game/index.php?page=resourceSettings"
	resp, err := b.device.GetClient().PostForm(url2, payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (b *OGame) getCachedResearch() ogame.Researches {
	if b.cache.researches == nil {
		researches, _ := b.getResearch()
		return researches
	}
	return *b.cache.researches
}

func (b *OGame) getResearch() (out ogame.Researches, err error) {
	page, err := getPage[parser.ResearchPage](b)
	if err != nil {
		return
	}
	researches := page.ExtractResearch()
	b.cache.researches = &researches
	return researches, nil
}

func (b *OGame) getCachedLfBonuses() (out ogame.LfBonuses, err error) {
	if b.cache.lfBonuses == nil {
		return b.getLfBonuses()
	}
	return *b.cache.lfBonuses, nil
}

func (b *OGame) getLfBonuses() (out ogame.LfBonuses, err error) {
	page, err := getPage[parser.LfBonusesPage](b)
	if err != nil {
		return
	}
	bonuses, err := page.ExtractLfBonuses()
	if err != nil {
		return
	}
	b.cache.lfBonuses = &bonuses
	return bonuses, nil
}

func (b *OGame) getCachedAllianceClass() (out ogame.AllianceClass, err error) {
	allianceClass := b.cache.allianceClass
	if allianceClass == nil {
		return b.getAllianceClass()
	}
	return *allianceClass, nil
}

func (b *OGame) getAllianceClass() (out ogame.AllianceClass, err error) {
	pageHTML, err := b.getPage("alliance")
	if err != nil {
		return
	}
	token, err := b.extractor.ExtractToken(pageHTML)
	if err != nil {
		return
	}
	allianceClass := ogame.NoAllianceClass
	if !bytes.Contains(pageHTML, []byte("createNewAlliance")) {
		vals := url.Values{"page": {"ingame"}, "component": {"alliance"}, "tab": {"overview"}, "action": {"fetchOverview"}, "ajax": {"1"}, "token": {token}}
		pageHTML, err = b.getPageContent(vals, SkipCacheFullPage)
		if err == nil && len(pageHTML) > 0 {
			var res parser.AllianceOverviewTabRes
			if err = json.Unmarshal(pageHTML, &res); err == nil {
				allianceClass, _ = b.extractor.ExtractAllianceClass([]byte(res.Content.AllianceAllianceOverview))
				b.cache.token = res.NewAjaxToken
			}
		}
	}
	b.cache.allianceClass = &allianceClass
	return allianceClass, nil
}

func (b *OGame) getResourcesBuildings(celestialID ogame.CelestialID, options ...Option) (ogame.ResourcesBuildings, error) {
	options = append(options, ChangePlanet(celestialID))
	page, err := getPage[parser.SuppliesPage](b, options...)
	if err != nil {
		return ogame.ResourcesBuildings{}, err
	}
	return page.ExtractResourcesBuildings()
}

func (b *OGame) getLfBuildings(celestialID ogame.CelestialID, options ...Option) (ogame.LfBuildings, error) {
	options = append(options, ChangePlanet(celestialID))
	page, err := getPage[parser.LfBuildingsPage](b, options...)
	if err != nil {
		return ogame.LfBuildings{}, err
	}
	return page.ExtractLfBuildings()
}

func (b *OGame) getLfResearch(celestialID ogame.CelestialID, options ...Option) (ogame.LfResearches, error) {
	options = append(options, ChangePlanet(celestialID))
	page, err := getPage[parser.LfResearchPage](b, options...)
	if err != nil {
		return ogame.LfResearches{}, err
	}
	return page.ExtractLfResearch()
}

func (b *OGame) getLfResearchDetails(celestialID ogame.CelestialID, options ...Option) (ogame.LfResearchDetails, error) {
	options = append(options, ChangePlanet(celestialID))
	page, err := getPage[parser.LfResearchPage](b, options...)
	if err != nil {
		return ogame.LfResearchDetails{}, err
	}
	lfResearch, err := page.ExtractLfResearch()
	if err != nil {
		return ogame.LfResearchDetails{}, err
	}
	slots := page.ExtractLfSlots()
	collected, limit := page.ExtractArtefacts()
	out := ogame.LfResearchDetails{
		LfResearches:       lfResearch,
		Slots:              slots,
		ArtefactsCollected: collected,
		ArtefactsLimit:     limit,
	}
	return out, nil
}

func (b *OGame) getDefense(celestialID ogame.CelestialID, options ...Option) (ogame.DefensesInfos, error) {
	options = append(options, ChangePlanet(celestialID))
	page, err := getPage[parser.DefensesPage](b, options...)
	if err != nil {
		return ogame.DefensesInfos{}, err
	}
	return page.ExtractDefense()
}

func (b *OGame) getShips(celestialID ogame.CelestialID, options ...Option) (ogame.ShipsInfos, error) {
	options = append(options, ChangePlanet(celestialID))
	page, err := getPage[parser.ShipyardPage](b, options...)
	if err != nil {
		return ogame.ShipsInfos{}, err
	}
	return page.ExtractShips()
}

func (b *OGame) getFacilities(celestialID ogame.CelestialID, options ...Option) (ogame.Facilities, error) {
	options = append(options, ChangePlanet(celestialID))
	page, err := getPage[parser.FacilitiesPage](b, options...)
	if err != nil {
		return ogame.Facilities{}, err
	}
	return page.ExtractFacilities()
}

func (b *OGame) getTechs(celestialID ogame.CelestialID) (ogame.Techs, error) {
	vals := url.Values{"page": {FetchTechsName}}
	page, err := getAjaxPage[parser.FetchTechsAjaxPage](b, vals, ChangePlanet(celestialID))
	if err != nil {
		return ogame.Techs{}, err
	}
	return page.ExtractTechs()
}

func (b *OGame) getProduction(celestialID ogame.CelestialID) ([]ogame.Quantifiable, int64, error) {
	page, err := getPage[parser.ShipyardPage](b, ChangePlanet(celestialID))
	if err != nil {
		return []ogame.Quantifiable{}, 0, err
	}
	return page.ExtractProduction()
}

func (b *OGame) technologyDetails(celestialID ogame.CelestialID, id ogame.ID) (ogame.TechnologyDetails, error) {
	pageHTML, err := b.getPageContent(url.Values{
		"page":       {"ingame"},
		"component":  {"technologydetails"},
		"ajax":       {"1"},
		"action":     {"getDetails"},
		"technology": {utils.FI64(id)},
		"cp":         {utils.FI64(celestialID)},
	})
	if err != nil {
		return ogame.TechnologyDetails{}, err
	}
	return b.extractor.ExtractTechnologyDetails(pageHTML)
}

func getToken(b *OGame, page string, celestialID ogame.CelestialID) (string, error) {
	pageHTML, _ := b.getPage(page, ChangePlanet(celestialID))
	return b.extractor.ExtractToken(pageHTML)
}

func (b *OGame) tearDown(celestialID ogame.CelestialID, id ogame.ID) error {
	var page string
	if id.IsResourceBuilding() {
		page = "supplies"
	} else if id.IsFacility() {
		page = "facilities"
	} else {
		return errors.New("invalid id " + id.String())
	}

	pageHTML, _ := b.getPage(page, ChangePlanet(celestialID))
	token, err := b.extractor.ExtractToken(pageHTML)
	if err != nil {
		return err
	}

	pageHTML, err = b.getPageContent(url.Values{
		"page":       {"ingame"},
		"component":  {"technologydetails"},
		"ajax":       {"1"},
		"action":     {"getDetails"},
		"technology": {utils.FI64(id)},
		"cp":         {utils.FI64(celestialID)},
	})
	if err != nil {
		return err
	}

	var jsonContent struct {
		Target  string `json:"target"`
		Content struct {
			Technologydetails string `json:"technologydetails"`
		} `json:"content"`
		Files struct {
			Js  []string `json:"js"`
			CSS []string `json:"css"`
		} `json:"files"`
		Page struct {
			StateObj string `json:"stateObj"`
			Title    string `json:"title"`
			URL      string `json:"url"`
		} `json:"page"`
		ServerTime   int    `json:"serverTime"`
		NewAjaxToken string `json:"newAjaxToken"`
	}

	if err := json.Unmarshal(pageHTML, &jsonContent); err != nil {
		return err
	}

	if ok, err := b.extractor.ExtractTearDownButtonEnabled([]byte(jsonContent.Content.Technologydetails)); err != nil || !ok {
		return errors.New("tear down button is disabled")
	}

	vals := url.Values{
		"page":      {"componentOnly"},
		"component": {"buildlistactions"},
		"action":    {"scheduleEntry"},
		"asJson":    {"1"},
	}
	payload := url.Values{
		"technologyId": {utils.FI64(id)},
		"mode":         {"3"},
		"token":        {token},
	}
	_, err = b.postPageContent(vals, payload)
	return err
}

var ErrBuild = errors.New("failed to build")

func (b *OGame) build(celestialID ogame.CelestialID, id ogame.ID, nbr int64) error {
	var page string
	if id.IsDefense() {
		page = DefensesPageName
	} else if id.IsShip() {
		page = ShipyardPageName
	} else if id.IsLfBuilding() {
		page = LfBuildingsPageName
	} else if id.IsLfTech() {
		page = LfResearchPageName
	} else if id.IsBuilding() {
		page = SuppliesPageName
	} else if id.IsTech() {
		page = ResearchPageName
	} else {
		return errors.New("invalid id " + id.String())
	}

	token, err := getToken(b, page, celestialID)
	if err != nil {
		return err
	}

	vals := url.Values{
		"page":      {"componentOnly"},
		"component": {"buildlistactions"},
		"action":    {"scheduleEntry"},
		"asJson":    {"1"},
	}

	var amount int64 = 1
	if id.IsShip() || id.IsDefense() {
		var maximumNbr int64 = 99999
		amount = min(nbr, maximumNbr)
	}

	payload := url.Values{
		"technologyId": {utils.FI64(id)},
		"amount":       {utils.FI64(amount)},
		"mode":         {"1"},
		"token":        {token},
		"planetId":     {utils.FI64(celestialID)},
	}

	var responseStruct struct {
		JsServerlang string        `json:"js_serverlang"`
		JsServerid   string        `json:"js_serverid"`
		Status       string        `json:"status"`
		Errors       []OGameError  `json:"errors"`
		Components   []interface{} `json:"components"`
		NewAjaxToken string        `json:"newAjaxToken"`
	}

	by, err := b.postPageContent(vals, payload)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(by, &responseStruct); err != nil {
		return err
	}
	if responseStruct.Status == "failure" {
		errInst := ErrBuild
		if len(responseStruct.Errors) > 0 {
			errStruct := responseStruct.Errors[0]
			errInst = fmt.Errorf("%w : %s (%d)", errInst, errStruct.Message, errStruct.Error)
		}
		return errInst
	}
	return nil
}

func (b *OGame) buildCancelable(celestialID ogame.CelestialID, id ogame.ID) error {
	if !id.IsBuilding() && !id.IsTech() && !id.IsLfBuilding() && !id.IsLfTech() {
		return errors.New("invalid id " + id.String())
	}
	return b.build(celestialID, id, 0)
}

func (b *OGame) buildProduction(celestialID ogame.CelestialID, id ogame.ID, nbr int64) error {
	if !id.IsDefense() && !id.IsShip() {
		return errors.New("invalid id " + id.String())
	}
	return b.build(celestialID, id, nbr)
}

func (b *OGame) buildBuilding(celestialID ogame.CelestialID, buildingID ogame.ID) error {
	if !buildingID.IsBuilding() {
		return errors.New("invalid building id " + buildingID.String())
	}
	return b.buildCancelable(celestialID, buildingID)
}

func (b *OGame) buildTechnology(celestialID ogame.CelestialID, technologyID ogame.ID) error {
	if !technologyID.IsTech() && !technologyID.IsLfTech() {
		return errors.New("invalid technology id " + technologyID.String())
	}
	return b.buildCancelable(celestialID, technologyID)
}

func (b *OGame) buildDefense(celestialID ogame.CelestialID, defenseID ogame.ID, nbr int64) error {
	if !defenseID.IsDefense() {
		return errors.New("invalid defense id " + defenseID.String())
	}
	return b.buildProduction(celestialID, defenseID, nbr)
}

func (b *OGame) buildShips(celestialID ogame.CelestialID, shipID ogame.ID, nbr int64) error {
	if !shipID.IsShip() {
		return errors.New("invalid ship id " + shipID.String())
	}
	return b.buildProduction(celestialID, shipID, nbr)
}

func (b *OGame) constructionsBeingBuilt(celestialID ogame.CelestialID) (ogame.Constructions, error) {
	page, err := getPage[parser.OverviewPage](b, ChangePlanet(celestialID))
	if err != nil {
		return ogame.Constructions{}, err
	}
	return page.ExtractConstructions()
}

func (b *OGame) cancel(token string, techID, listID int64) error {
	_, err := b.postPageContent(url.Values{"page": {"componentOnly"}, "component": {"buildlistactions"}, "action": {"cancelEntry"}, "asJson": {"1"}},
		url.Values{"technologyId": {utils.FI64(techID)}, "listId": {utils.FI64(listID)}, "token": {token}})
	if err != nil {
		return err
	}
	return nil
}

func (b *OGame) cancelBuilding(celestialID ogame.CelestialID) error {
	page, err := getPage[parser.OverviewPage](b, ChangePlanet(celestialID))
	if err != nil {
		return err
	}
	token, techID, listID, _ := page.ExtractCancelBuildingInfos()
	return b.cancel(token, techID, listID)
}

func (b *OGame) cancelLfBuilding(celestialID ogame.CelestialID) error {
	page, err := getPage[parser.OverviewPage](b, ChangePlanet(celestialID))
	if err != nil {
		return err
	}
	token, id, listID, _ := page.ExtractCancelLfBuildingInfos()
	return b.cancel(token, id, listID)
}

func (b *OGame) cancelResearch(celestialID ogame.CelestialID) error {
	page, err := getPage[parser.OverviewPage](b, ChangePlanet(celestialID))
	if err != nil {
		return err
	}
	token, techID, listID, _ := page.ExtractCancelResearchInfos()
	return b.cancel(token, techID, listID)
}

func (b *OGame) fetchResources(celestialID ogame.CelestialID) (ogame.ResourcesDetails, error) {
	pageJSON, err := b.getPage(FetchResourcesPageName, ChangePlanet(celestialID))
	if err != nil {
		return ogame.ResourcesDetails{}, err
	}
	return b.extractor.ExtractResourcesDetails(pageJSON)
}

func (b *OGame) getResources(celestialID ogame.CelestialID) (ogame.Resources, error) {
	res, err := b.fetchResources(celestialID)
	if err != nil {
		return ogame.Resources{}, err
	}
	return ogame.Resources{
		Metal:      res.Metal.Available,
		Crystal:    res.Crystal.Available,
		Deuterium:  res.Deuterium.Available,
		Energy:     res.Energy.Available,
		Darkmatter: res.Darkmatter.Available,
		Population: res.Population.Available,
		Food:       res.Food.Available,
	}, nil
}

func (b *OGame) getResourcesDetails(celestialID ogame.CelestialID) (ogame.ResourcesDetails, error) {
	return b.fetchResources(celestialID)
}

func (b *OGame) destroyRockets(planetID ogame.PlanetID, abm, ipm int64) error {
	vals := url.Values{
		"page":      {"ajax"},
		"component": {RocketlayerPageName},
		"overlay":   {"1"},
	}
	page, err := getAjaxPage[parser.RocketlayerAjaxPage](b, vals, ChangePlanet(planetID.Celestial()))
	if err != nil {
		return err
	}
	maxABM, maxIPM, token, err := page.ExtractDestroyRockets()
	if err != nil {
		return err
	}
	if maxABM == 0 && maxIPM == 0 {
		return errors.New("no missile to destroy")
	}
	if abm > maxABM {
		abm = maxABM
	}
	if ipm > maxIPM {
		ipm = maxIPM
	}
	params := url.Values{
		"page":      {"ajax"},
		"component": {"rocketlayer"},
		"action":    {"destroy"},
		"ajax":      {"1"},
		"asJson":    {"1"},
	}
	payload := url.Values{
		"interceptorMissile":    {utils.FI64(abm)},
		"interplanetaryMissile": {utils.FI64(ipm)},
		"token":                 {token},
	}
	by, err := b.postPageContent(params, payload)
	if err != nil {
		return err
	}
	// {"status":"success","message":"The following missiles have been destroyed:\nInterplanetary missiles: 1\nAnti-ballistic missiles: 2","components":[],"newAjaxToken":"ec306346888f14e38c4248aa78e56610"}
	var resp struct {
		Status       string `json:"status"`
		Message      string `json:"message"`
		NewAjaxToken string `json:"newAjaxToken"`
	}
	if err := json.Unmarshal(by, &resp); err != nil {
		return err
	}
	if resp.Status != "success" {
		return errors.New(resp.Message)
	}

	return nil
}

func (b *OGame) sendIPM(planetID ogame.PlanetID, coord ogame.Coordinate, nbr int64, priority ogame.ID) (int64, error) {
	if !priority.IsValidIPMTarget() {
		return 0, errors.New("invalid defense target id")
	}
	vals := url.Values{
		"page":       {"ajax"},
		"component":  {"missileattacklayer"},
		"galaxy":     {utils.FI64(coord.Galaxy)},
		"system":     {utils.FI64(coord.System)},
		"position":   {utils.FI64(coord.Position)},
		"planetType": {utils.FI64(coord.Type)},
	}
	page, err := getAjaxPage[parser.MissileAttackLayerAjaxPage](b, vals, ChangePlanet(planetID.Celestial()))
	if err != nil {
		return 0, err
	}

	duration, maxV, token, err := page.ExtractIPM()
	if err != nil {
		return 0, err
	}
	if maxV == 0 {
		return 0, errors.New("no missile available")
	}
	nbr = min(nbr, maxV)
	params := url.Values{
		"page":      {"ajax"},
		"component": {"missileattacklayer"},
		"action":    {"sendMissiles"},
		"ajax":      {"1"},
		"asJson":    {"1"},
	}
	payload := url.Values{
		"galaxy":               {utils.FI64(coord.Galaxy)},
		"system":               {utils.FI64(coord.System)},
		"position":             {utils.FI64(coord.Position)},
		"type":                 {utils.FI64(coord.Type)},
		"token":                {token},
		"missileCount":         {utils.FI64(nbr)},
		"missilePrimaryTarget": {utils.FI64(priority)},
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
	Errors          []OGameError `json:"errors"`
	TargetOk        bool         `json:"targetOk"`
	Components      []any        `json:"components"`
	EmptySystems    int64        `json:"emptySystems"`
	InactiveSystems int64        `json:"inactiveSystems"`
	NewAjaxToken    string       `json:"newAjaxToken"`
}

func (b *OGame) checkTarget(ships ogame.ShipsInfos, where ogame.Coordinate, opts ...Option) (out CheckTargetResponse, err error) {
	payload := url.Values{}
	for shipID, nb := range ships.IterFlyable() {
		payload.Set("am"+utils.FI64(shipID), utils.FI64(nb))
	}
	payload.Set("token", b.cache.token)
	payload.Set("galaxy", utils.FI64(where.Galaxy))
	payload.Set("system", utils.FI64(where.System))
	payload.Set("position", utils.FI64(where.Position))
	payload.Set("type", utils.FI64(where.Type))
	payload.Set("union", "0")
	by, err := b.postPageContent(url.Values{"page": {"ingame"}, "component": {"fleetdispatch"}, "action": {"checkTarget"}, "ajax": {"1"}, "asJson": {"1"}}, payload, opts...)
	if err != nil {
		return out, err
	}
	if err := json.Unmarshal(by, &out); err != nil {
		return out, err
	}
	return
}

func (b *OGame) sendFleet(celestialID ogame.CelestialID, ships ogame.ShipsInfos, speed ogame.Speed, where ogame.Coordinate,
	mission ogame.MissionID, resources ogame.Resources, holdingTime, unionID int64, ensure bool) (ogame.Fleet, error) {
	zeroFleet := ogame.MakeFleet()

	// Get existing fleet, so we can ensure new fleet ID is greater
	initialFleets, slots, err := b.getFleets()
	if err != nil {
		return zeroFleet, err
	}
	maxInitialFleetID := ogame.FleetID(0)
	for _, f := range initialFleets {
		maxInitialFleetID = max(maxInitialFleetID, f.ID)
	}

	if slots.IsAllSlotsInUse(mission) {
		return zeroFleet, ogame.ErrAllSlotsInUse
	}

	// Page 1 : get to fleet page
	pageHTML, err := b.getPage(FleetdispatchPageName, ChangePlanet(celestialID))
	if err != nil {
		return zeroFleet, err
	}

	fleet1Doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return zeroFleet, err
	}
	fleet1BodyID := b.extractor.ExtractBodyIDFromDoc(fleet1Doc)
	if fleet1BodyID != FleetdispatchPageName {
		return zeroFleet, ogame.ErrInvalidPlanetID
	}

	if attackBlockActivated, blockedUntil := b.extractor.ExtractAttackBlockFromDoc(fleet1Doc); attackBlockActivated {
		return zeroFleet, ogame.NewAttackBlockActivatedErr(blockedUntil)
	}

	if b.extractor.ExtractIsInVacationFromDoc(fleet1Doc) {
		return zeroFleet, ogame.ErrAccountInVacationMode
	}

	// Ensure we're not trying to attack/spy ourselves
	myCelestials, _ := b.extractor.ExtractCelestialsFromDoc(fleet1Doc)
	for _, c := range myCelestials {
		if c.GetCoordinate().Equal(where) {
			if c.GetID() == celestialID {
				return zeroFleet, errors.New("origin and destination are the same")
			} else if mission == ogame.Spy {
				return zeroFleet, errors.New("you cannot spy yourself")
			} else if mission == ogame.Attack {
				return zeroFleet, errors.New("you cannot attack yourself")
			}
			break
		}
	}

	availableShips, err := b.extractor.ExtractFleet1ShipsFromDoc(fleet1Doc)
	if err != nil {
		return zeroFleet, err
	}

	atLeastOneShipSelected := false
	for shipID, nb := range ships.IterFlyable() {
		avail := availableShips.ByID(shipID)
		if ensure && nb > avail {
			return zeroFleet, fmt.Errorf("%w, %s (%d > %d)", ogame.ErrNotEnoughShips, ogame.Objs.ByID(shipID).GetName(), nb, avail)
		}
		nb = min(nb, avail)
		if nb > 0 {
			atLeastOneShipSelected = true
		}
	}
	if !atLeastOneShipSelected {
		return zeroFleet, ogame.ErrNoShipSelected
	}

	payload := b.extractor.ExtractHiddenFieldsFromDoc(fleet1Doc)
	payload.Del("expeditionFleetTemplateId")
	for shipID, nb := range ships.IterFlyable() {
		payload.Set("am"+utils.FI64(shipID), utils.FI64(nb))
	}

	token, err := b.extractor.ExtractToken(pageHTML)
	if err != nil {
		return zeroFleet, err
	}

	payload.Set("token", token)
	payload.Set("galaxy", utils.FI64(where.Galaxy))
	payload.Set("system", utils.FI64(where.System))
	payload.Set("position", utils.FI64(where.Position))
	if mission == ogame.RecycleDebrisField {
		where.Type = ogame.DebrisType // Send to debris field
	} else if mission == ogame.Colonize || mission == ogame.Expedition {
		where.Type = ogame.PlanetType
	}
	payload.Set("type", utils.FI64(where.Type))
	payload.Set("union", "0")

	if unionID != 0 {
		acsArr := b.extractor.ExtractFleetDispatchACSFromDoc(fleet1Doc)
		acs := utils.Find(acsArr, func(v ogame.ACSValues) bool { return v.Union == unionID })
		if acs == nil {
			return zeroFleet, ogame.ErrUnionNotFound
		}
		payload.Add("acsValues", acs.ACSValues)
		payload.Add("union", utils.FI64(acs.Union))
		mission = ogame.GroupedAttack
	}

	// Check
	by1, err := b.postPageContent(url.Values{"page": {"ingame"}, "component": {"fleetdispatch"}, "action": {"checkTarget"}, "ajax": {"1"}, "asJson": {"1"}}, payload)
	if err != nil {
		return zeroFleet, err
	}
	var checkRes CheckTargetResponse
	if err := json.Unmarshal(by1, &checkRes); err != nil {
		return zeroFleet, err
	}

	if !checkRes.TargetOk {
		if len(checkRes.Errors) > 0 {
			return zeroFleet, errors.New(checkRes.Errors[0].Message + " (" + strconv.Itoa(checkRes.Errors[0].Error) + ")")
		}
		return zeroFleet, errors.New("target is not ok")
	}

	lfBonuses, err := b.getCachedLfBonuses()
	if err != nil {
		return zeroFleet, err
	}
	serverData := b.getServerData()
	multiplier := float64(serverData.CargoHyperspaceTechMultiplier) / 100.0
	cargo := ships.Cargo(b.getCachedResearch(), lfBonuses, b.cache.characterClass, multiplier, b.server.OGameSettings().ProbeRaidsEnabled())
	newResources := ogame.Resources{}
	if resources.Total() > cargo {
		newResources.Deuterium = utils.Clamp(resources.Deuterium, 0, cargo)
		cargo -= newResources.Deuterium
		newResources.Crystal = utils.Clamp(resources.Crystal, 0, cargo)
		cargo -= newResources.Crystal
		newResources.Metal = utils.Clamp(resources.Metal, 0, cargo)
	} else {
		newResources = resources
	}

	// Page 3 : select coord, mission, speed
	payload.Set("token", checkRes.NewAjaxToken)
	payload.Set("speed", utils.FI64(int64(speed)))
	payload.Set("crystal", utils.FI64(newResources.Crystal))
	payload.Set("deuterium", utils.FI64(newResources.Deuterium))
	payload.Set("metal", utils.FI64(newResources.Metal))
	payload.Set("mission", utils.FI64(mission))
	payload.Set("prioMetal", "1")
	payload.Set("prioCrystal", "2")
	payload.Set("prioDeuterium", "3")
	payload.Set("retreatAfterDefenderRetreat", "0")
	payload.Set("lootFoodOnAttack", "0")
	if mission == ogame.ParkInThatAlly || mission == ogame.Expedition {
		if mission == ogame.Expedition { // Expedition 1 to 18
			holdingTime = utils.Clamp(holdingTime, 1, 18)
		} else { // ParkInThatAlly 0, 1, 2, 4, 8, 16, 32
			holdingTime = utils.Clamp(holdingTime, 0, 32)
		}
		payload.Set("holdingtime", utils.FI64(holdingTime))
	}

	// Page 4 : send the fleet
	res, err := b.postPageContent(url.Values{"page": {"ingame"}, "component": {"fleetdispatch"}, "action": {"sendFleet"}, "ajax": {"1"}, "asJson": {"1"}}, payload)
	if err != nil {
		return zeroFleet, err
	}
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
		Success           bool   `json:"success"`
		Message           string `json:"message"`
		FleetSendingToken string `json:"fleetSendingToken"`
		Components        []any  `json:"components"`
		RedirectURL       string `json:"redirectUrl"`
		Errors            []struct {
			Message string `json:"message"`
			Error   int64  `json:"error"`
		} `json:"errors"`
	}
	if err := json.Unmarshal(res, &resStruct); err != nil {
		return zeroFleet, errors.New("failed to unmarshal response: " + err.Error())
	}

	if len(resStruct.Errors) > 0 {
		return zeroFleet, errors.New(resStruct.Errors[0].Message + " (" + utils.FI64(resStruct.Errors[0].Error) + ")")
	}

	// Page 5
	page, err := getPage[parser.MovementPage](b)
	if err != nil {
		return zeroFleet, err
	}
	originCoords, _ := page.ExtractPlanetCoordinate()
	fleets, err := page.ExtractFleets()
	if err != nil {
		return zeroFleet, err
	}
	if maxV, err := getLastFleetFor(fleets, originCoords, where, mission); err == nil && maxV.ID > maxInitialFleetID {
		return maxV, nil
	}

	slots, _ = page.ExtractSlots()
	if slots.IsAllSlotsInUse(mission) {
		return zeroFleet, ogame.ErrAllSlotsInUse
	}

	return zeroFleet, errors.New("could not find new fleet ID")
}

func (b *OGame) fastMiniFleetSpy(coord ogame.Coordinate, shipCount int64, options ...Option) (ogame.MinifleetResponse, error) {
	vals := url.Values{
		"page":      {"ingame"},
		"component": {"fleetdispatch"},
		"action":    {"miniFleet"},
		"ajax":      {"1"},
		"asJson":    {"1"},
	}
	payload := url.Values{
		"mission":   {utils.FI64(ogame.Spy)},
		"galaxy":    {utils.FI64(coord.Galaxy)},
		"system":    {utils.FI64(coord.System)},
		"position":  {utils.FI64(coord.Position)},
		"type":      {utils.FI64(coord.Type)},
		"shipCount": {utils.FI64(shipCount)},
		"token":     {b.cache.token},
	}
	var res ogame.MinifleetResponse
	pageHTML, err := b.postPageContent(vals, payload, options...)
	if err != nil {
		return res, err
	}
	if err := json.Unmarshal(pageHTML, &res); err != nil {
		return res, err
	}
	if !res.Response.Success {
		msg := res.Response.Message
		rgx := regexp.MustCompile(`\d{2}\.\d{2}\.\d{4} \d{2}:\d{2}:\d{2}`)
		if match := rgx.FindString(msg); match != "" {
			if blockedUntil, err := time.Parse("02.01.2006 15:04:05", match); err == nil {
				return res, ogame.NewAttackBlockActivatedErr(blockedUntil)
			}
		}
		return res, errors.New(msg)
	}
	return res, nil
}

func (b *OGame) miniFleetSpy(coord ogame.Coordinate, shipCount int64, options ...Option) (ogame.Fleet, error) {
	fleet := ogame.MakeFleet()
	if _, err := b.fastMiniFleetSpy(coord, shipCount, options...); err != nil {
		return fleet, err
	}
	by, err := b.getPageContent(url.Values{"page": {"componentOnly"}, "component": {EventListAjaxPageName}, "ajax": {"1"}})
	if err != nil {
		return fleet, err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(by))
	if err != nil {
		return fleet, err
	}
	type FleetRow struct {
		ArrivalTs int64
		FleetID   ogame.FleetID
	}
	rows := make([]FleetRow, 0)
	for _, s := range doc.Find("tr.eventFleet").EachIter() {
		arrivalTs := utils.DoParseI64(s.AttrOr("data-arrival-time", ""))
		fleetID := ogame.FleetID(utils.DoParseI64(s.Find("a.recallFleet").AttrOr("data-fleet-id", "")))
		rows = append(rows, FleetRow{ArrivalTs: arrivalTs, FleetID: fleetID})
	}
	maxRow := FleetRow{}
	for _, row := range rows {
		if row.FleetID > maxRow.FleetID {
			maxRow = row
		}
	}
	if maxRow.FleetID == 0 {
		return fleet, errors.New("could not find fleet ID")
	}
	currCelestial, _ := b.getCachedCelestial(b.cache.planetID)
	fleet.Mission = ogame.Spy
	fleet.ID = maxRow.FleetID
	fleet.Origin = currCelestial.GetCoordinate()
	fleet.Destination = coord
	fleet.Ships = ogame.ShipsInfos{EspionageProbe: shipCount}
	fleet.StartTime = time.Now()
	fleet.ArrivalTime = time.Unix(maxRow.ArrivalTs, 0)
	fleet.ArriveIn = int64(time.Until(fleet.ArrivalTime).Seconds())
	fleet.BackTime = fleet.ArrivalTime.Add(time.Duration(2*fleet.ArriveIn) * time.Second)
	fleet.BackIn = int64(time.Until(fleet.BackTime).Seconds())
	return fleet, nil
}

func (b *OGame) getPageMessages(page int64, tabid ogame.MessagesTabID) ([]byte, error) {
	payload := url.Values{
		"activeSubTab": {utils.FI64(tabid)},
		"showTrash":    {"false"},
	}
	return b.postPageContent(url.Values{"page": {"componentOnly"}, "component": {"messages"}, "asJson": {"1"}, "action": {"getMessagesList"}}, payload)
}

func (b *OGame) getEspionageReportMessages(maxPage int64) ([]ogame.EspionageReportSummary, error) {
	return getMessages(b, maxPage, EspionageMessagesTabID, b.extractor.ExtractEspionageReportMessageIDs)
}

func (b *OGame) getCombatReportMessages(maxPage int64) ([]ogame.CombatReportSummary, error) {
	return getMessages(b, maxPage, CombatReportsMessagesTabID, b.extractor.ExtractCombatReportMessagesSummary)
}

func (b *OGame) getExpeditionMessages(maxPage int64) ([]ogame.ExpeditionMessage, error) {
	return getMessages(b, maxPage, ExpeditionsMessagesTabID, b.extractor.ExtractExpeditionMessages)
}

func getMessages[T any](b *OGame, maxPage int64, tabID ogame.MessagesTabID, extractor func([]byte) ([]T, int64, error)) ([]T, error) {
	var res struct {
		ServerLang   string `json:"js_serverlang"`
		ServerID     string `json:"js_serverid"`
		Status       string `json:"status"`
		Messages     []any  `json:"messages"`
		Components   []any  `json:"components"`
		NewAjaxToken string `json:"newAjaxToken"`
	}
	pageHTML, _ := b.getPageMessages(1, tabID)
	_ = json.Unmarshal(pageHTML, &res)
	msgs := make([]T, 0)
	for _, m := range res.Messages {
		doc := ([]byte)(m.(string))
		newMessage, _, _ := extractor(doc)
		msgs = append(msgs, newMessage...)
	}
	return msgs, nil
}

func (b *OGame) collectAllMarketplaceMessages() error {
	purchases, _ := b.getMarketplacePurchasesMessages()
	sales, _ := b.getMarketplaceSalesMessages()
	msgs := make([]ogame.MarketplaceMessage, 0)
	msgs = append(msgs, purchases...)
	msgs = append(msgs, sales...)
	newToken := ""
	var err error
	for _, msg := range msgs {
		if msg.MarketTransactionID != 0 {
			newToken, err = b.collectMarketplaceMessage(msg, newToken)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type collectMarketplaceResponse struct {
	MarketTransactionID int    `json:"marketTransactionId"`
	Status              string `json:"status"`
	Message             string `json:"message"`
	StatusMessage       string `json:"statusMessage"`
	NewToken            string `json:"newToken"`
	Components          []any  `json:"components"`
}

func (b *OGame) collectMarketplaceMessage(msg ogame.MarketplaceMessage, newToken string) (string, error) {
	params := url.Values{
		"page":                {"componentOnly"},
		"component":           {"marketplace"},
		"marketTransactionId": {utils.FI64(msg.MarketTransactionID)},
		"token":               {msg.Token},
		"asJson":              {"1"},
	}
	if msg.Type == 26 { // purchase
		params.Set("action", "collectItem")
	} else if msg.Type == 27 { // sale
		params.Set("action", "collectPrice")
	}
	payload := url.Values{
		"newToken": {newToken},
	}
	by, err := b.postPageContent(params, payload)
	var res collectMarketplaceResponse
	if err := json.Unmarshal(by, &res); err != nil {
		return "", errors.New("failed to unmarshal json response: " + err.Error())
	}
	return res.NewToken, err
}

func (b *OGame) getMarketplacePurchasesMessages() ([]ogame.MarketplaceMessage, error) {
	return b.getMarketplaceMessages(-1, MarketplacePurchasesMessagesTabID)
}

func (b *OGame) getMarketplaceSalesMessages() ([]ogame.MarketplaceMessage, error) {
	return b.getMarketplaceMessages(-1, MarketplaceSalesMessagesTabID)
}

// tabID 26: purchases, 27: sales
func (b *OGame) getMarketplaceMessages(maxPage int64, tabID ogame.MessagesTabID) ([]ogame.MarketplaceMessage, error) {
	return getMessages(b, maxPage, tabID, b.extractor.ExtractMarketplaceMessages)
}

func (b *OGame) getExpeditionMessageAt(t time.Time) (ogame.ExpeditionMessage, error) {
	pageHTML, _ := b.getPageMessages(1, ExpeditionsMessagesTabID)
	var res struct {
		ServerLang   string `json:"js_serverlang"`
		ServerID     string `json:"js_serverid"`
		Status       string `json:"status"`
		Messages     []any  `json:"messages"`
		Components   []any  `json:"components"`
		NewAjaxToken string `json:"newAjaxToken"`
	}
	_ = json.Unmarshal(pageHTML, &res)
	newMessages := make([]ogame.ExpeditionMessage, 0)
	for _, m := range res.Messages {
		doc := ([]byte)(m.(string))
		newMessage, _, _ := b.extractor.ExtractExpeditionMessages(doc)
		newMessages = append(newMessages, newMessage...)
	}
	for _, m := range newMessages {
		if m.CreatedAt.Unix() == t.Unix() {
			return m, nil
		}
		if m.CreatedAt.Unix() < t.Unix() {
			break
		}
	}
	return ogame.ExpeditionMessage{}, errors.New("expedition message not found for " + t.String())
}

func (b *OGame) getCombatReportSummaries() ([]ogame.CombatReportSummary, error) {
	pageHTML, err := b.getPageMessages(1, CombatReportsMessagesTabID)
	if err != nil {
		return nil, err
	}
	var res struct {
		ServerLang   string `json:"js_serverlang"`
		ServerID     string `json:"js_serverid"`
		Status       string `json:"status"`
		Messages     []any  `json:"messages"`
		Components   []any  `json:"components"`
		NewAjaxToken string `json:"newAjaxToken"`
	}
	_ = json.Unmarshal(pageHTML, &res)
	newMessages := make([]ogame.CombatReportSummary, 0)
	for i, m := range res.Messages {
		if i > 40 {
			break
		}
		doc := ([]byte)(m.(string))
		newMessage, _, _ := b.extractor.ExtractCombatReportMessagesSummary(doc)
		newMessages = append(newMessages, newMessage...)
	}
	return newMessages, nil
}

func (b *OGame) getCombatReportForFleet(fleetID ogame.FleetID) (ogame.CombatReportSummary, error) {
	newMessages, err := b.getCombatReportSummaries()
	if err != nil {
		return ogame.CombatReportSummary{}, err
	}
	for _, m := range newMessages {
		if m.FleetID == fleetID {
			return m, nil
		}
	}
	return ogame.CombatReportSummary{}, errors.New("combat report not found for " + fleetID.String())
}

func (b *OGame) getCombatReportFor(coord ogame.Coordinate) (ogame.CombatReportSummary, error) {
	newMessages, err := b.getCombatReportSummaries()
	if err != nil {
		return ogame.CombatReportSummary{}, err
	}
	for _, m := range newMessages {
		if m.Destination.Equal(coord) {
			return m, nil
		}
	}
	return ogame.CombatReportSummary{}, errors.New("combat report not found for " + coord.String())
}

func (b *OGame) getEspionageReport(msgID int64) (ogame.EspionageReport, error) {
	pageHTML, err := b.getPageContent(url.Values{"page": {"componentOnly"}, "component": {"messagedetails"}, "messageId": {utils.FI64(msgID)}})
	if err != nil {
		return ogame.EspionageReport{}, err
	}
	return b.extractor.ExtractEspionageReport(pageHTML)
}

func (b *OGame) getEspionageReportFor(coord ogame.Coordinate) (ogame.EspionageReport, error) {
	pageHTML, err := b.getPageMessages(1, EspionageMessagesTabID)
	if err != nil {
		return ogame.EspionageReport{}, err
	}
	var res struct {
		ServerLang   string `json:"js_serverlang"`
		ServerID     string `json:"js_serverid"`
		Status       string `json:"status"`
		Messages     []any  `json:"messages"`
		Components   []any  `json:"components"`
		NewAjaxToken string `json:"newAjaxToken"`
	}
	_ = json.Unmarshal(pageHTML, &res)
	newMessages := make([]ogame.EspionageReportSummary, 0)
	for i, m := range res.Messages {
		if i > 40 {
			break
		}
		doc := ([]byte)(m.(string))
		newMessage, _, _ := b.extractor.ExtractEspionageReportMessageIDs(doc)
		newMessages = append(newMessages, newMessage...)
	}
	for _, m := range newMessages {
		if m.Target.Equal(coord) {
			return b.getEspionageReport(m.ID)
		}
	}
	return ogame.EspionageReport{}, errors.New("espionage report not found for " + coord.String())
}

func (b *OGame) getDeleteMessagesToken() (string, error) {
	var tmp struct {
		NewAjaxToken string
	}
	pageJson, err := b.getPageMessages(1, EspionageMessagesTabID)
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(pageJson, &tmp); err != nil {
		return "", err
	}
	return tmp.NewAjaxToken, nil
}

func (b *OGame) deleteMessage(msgID int64) error {
	token, err := b.getDeleteMessagesToken()
	if err != nil {
		return err
	}
	vals := url.Values{
		"page":      {"componentOnly"},
		"component": {"messages"},
		"asJson":    {"1"},
		"action":    {"flagDeleted"},
	}
	payload := url.Values{
		"token":        {token},
		"messageIds[]": {utils.FI64(msgID)},
	}
	by, err := b.postPageContent(vals, payload)
	if err != nil {
		return err
	}

	var res struct {
		Status       string `json:"status"`
		Message      string `json:"message"`
		NewAjaxToken string `json:"newAjaxToken"`
	}
	if err := json.Unmarshal(by, &res); err != nil {
		return errors.New("unable to find message id " + utils.FI64(msgID))
	}
	if res.Status != "success" {
		return errors.New("unable to find message id " + utils.FI64(msgID) + " : " + res.Message)
	}
	return nil
}

const (
	EspionageMessagesTabID            ogame.MessagesTabID = 20
	CombatReportsMessagesTabID        ogame.MessagesTabID = 21
	ExpeditionsMessagesTabID          ogame.MessagesTabID = 22
	UnionsTransportMessagesTabID      ogame.MessagesTabID = 23
	OtherMessagesTabID                ogame.MessagesTabID = 24
	MarketplacePurchasesMessagesTabID ogame.MessagesTabID = 26
	MarketplaceSalesMessagesTabID     ogame.MessagesTabID = 27
)

func (b *OGame) deleteAllMessagesFromTab(tabID ogame.MessagesTabID) error {
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
	token, err := b.getDeleteMessagesToken()
	if err != nil {
		return err
	}
	payload := url.Values{
		"tabid":     {utils.FI64(tabID)},
		"messageId": {utils.FI64(-1)},
		"action":    {"103"},
		"ajax":      {"1"},
		"token":     {token},
	}
	pageHTML, err := b.postPageContent(url.Values{"page": {"messages"}}, payload)
	var res struct {
		Status       string `json:"status"`
		Message      string `json:"message"`
		NewAjaxToken string `json:"newAjaxToken"`
	}
	_ = json.Unmarshal(pageHTML, &res)
	return err
}

func energyProduced(temp ogame.Temperature, resourcesBuildings ogame.ResourcesBuildings, resSettings ogame.ResourceSettings, energyTechnology int64) int64 {
	energyProduced := int64(float64(ogame.SolarPlant.Production(resourcesBuildings.SolarPlant)) * (float64(resSettings.SolarPlant) / 100))
	energyProduced += int64(float64(ogame.FusionReactor.Production(energyTechnology, resourcesBuildings.FusionReactor)) * (float64(resSettings.FusionReactor) / 100))
	energyProduced += int64(float64(ogame.SolarSatellite.Production(temp, resourcesBuildings.SolarSatellite, false)) * (float64(resSettings.SolarSatellite) / 100))
	return energyProduced
}

func energyNeeded(resourcesBuildings ogame.ResourcesBuildings, resSettings ogame.ResourceSettings) int64 {
	energyNeeded := int64(float64(ogame.MetalMine.EnergyConsumption(resourcesBuildings.MetalMine)) * (float64(resSettings.MetalMine) / 100))
	energyNeeded += int64(float64(ogame.CrystalMine.EnergyConsumption(resourcesBuildings.CrystalMine)) * (float64(resSettings.CrystalMine) / 100))
	energyNeeded += int64(float64(ogame.DeuteriumSynthesizer.EnergyConsumption(resourcesBuildings.DeuteriumSynthesizer)) * (float64(resSettings.DeuteriumSynthesizer) / 100))
	return energyNeeded
}

func productionRatio(temp ogame.Temperature, resourcesBuildings ogame.ResourcesBuildings, resSettings ogame.ResourceSettings, energyTechnology int64) float64 {
	energyProduced := energyProduced(temp, resourcesBuildings, resSettings, energyTechnology)
	energyNeeded := energyNeeded(resourcesBuildings, resSettings)
	ratio := 1.0
	if energyNeeded > energyProduced {
		ratio = float64(energyProduced) / float64(energyNeeded)
	}
	return ratio
}

func getProductions(resBuildings ogame.ResourcesBuildings, resSettings ogame.ResourceSettings, researches ogame.Researches, universeSpeed int64,
	temp ogame.Temperature, globalRatio float64) ogame.Resources {
	energyProduced := energyProduced(temp, resBuildings, resSettings, researches.EnergyTechnology)
	energyNeeded := energyNeeded(resBuildings, resSettings)
	metalSetting := float64(resSettings.MetalMine) / 100
	crystalSetting := float64(resSettings.CrystalMine) / 100
	deutSetting := float64(resSettings.DeuteriumSynthesizer) / 100
	return ogame.Resources{
		Metal:     ogame.MetalMine.Production(universeSpeed, metalSetting, globalRatio, researches.PlasmaTechnology, resBuildings.MetalMine),
		Crystal:   ogame.CrystalMine.Production(universeSpeed, crystalSetting, globalRatio, researches.PlasmaTechnology, resBuildings.CrystalMine),
		Deuterium: ogame.DeuteriumSynthesizer.Production(universeSpeed, temp.Mean(), deutSetting, globalRatio, researches.PlasmaTechnology, resBuildings.DeuteriumSynthesizer) - ogame.FusionReactor.GetFuelConsumption(universeSpeed, float64(resSettings.FusionReactor)/100, resBuildings.FusionReactor),
		Energy:    energyProduced - energyNeeded,
	}
}

func (b *OGame) getResourcesProductions(planetID ogame.PlanetID) (ogame.Resources, error) {
	planet, _ := b.getPlanet(planetID)
	resBuildings, _ := b.getResourcesBuildings(planetID.Celestial())
	researches, _ := b.getResearch()
	universeSpeed := b.cache.serverData.Speed
	resSettings, _ := b.getResourceSettings(planetID)
	ratio := productionRatio(planet.Temperature, resBuildings, resSettings, researches.EnergyTechnology)
	productions := getProductions(resBuildings, resSettings, researches, universeSpeed, planet.Temperature, ratio)
	return productions, nil
}

func getResourcesProductionsLight(resBuildings ogame.ResourcesBuildings, researches ogame.Researches,
	resSettings ogame.ResourceSettings, temp ogame.Temperature, universeSpeed int64) ogame.Resources {
	ratio := productionRatio(temp, resBuildings, resSettings, researches.EnergyTechnology)
	productions := getProductions(resBuildings, resSettings, researches, universeSpeed, temp, ratio)
	return productions
}

func (b *OGame) getPublicIP() (string, error) {
	var res struct {
		IP string `json:"ip"`
	}
	req, err := http.NewRequest(http.MethodGet, "https://jsonip.com/", nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	resp, err := b.doReqWithLoginProxyTransport(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	by, err := io.ReadAll(resp.Body)
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
	if b.lockedAtom.CompareAndSwap(false, true) {
		b.state = lockedBy
		b.stateChanged(true, lockedBy)
	}
}

func (b *OGame) botUnlock(unlockedBy string) {
	b.Unlock()
	if b.lockedAtom.CompareAndSwap(true, false) {
		b.state = unlockedBy
		b.stateChanged(false, unlockedBy)
	}
}

func (b *OGame) addAccount(number int, lang string) (*gameforge.AddAccountResponse, error) {
	accountGroup := fmt.Sprintf("%s_%d", lang, number)
	return gameforge.AddAccount(b.ctx, b.device, PLATFORM, b.lobby, accountGroup, b.bearerToken)
}

func (b *OGame) getCachedCelestial(v IntoCelestial) (Celestial, error) {
	getCachedCelestialByID := b.getCachedCelestialByID
	getCachedCelestialByCoord := b.getCachedCelestialByCoord
	switch vv := v.(type) {
	case Celestial:
		return vv, nil
	case Planet:
		return vv, nil
	case Moon:
		return vv, nil
	case ogame.CelestialID:
		return getCachedCelestialByID(vv)
	case ogame.PlanetID:
		return getCachedCelestialByID(vv.Celestial())
	case ogame.MoonID:
		return getCachedCelestialByID(vv.Celestial())
	case int:
		return getCachedCelestialByID(ogame.CelestialID(vv))
	case int32:
		return getCachedCelestialByID(ogame.CelestialID(vv))
	case int64:
		return getCachedCelestialByID(ogame.CelestialID(vv))
	case float32:
		return getCachedCelestialByID(ogame.CelestialID(vv))
	case float64:
		return getCachedCelestialByID(ogame.CelestialID(vv))
	case HasCoordinate:
		coordinate, err := ConvertIntoCoordinate(b, vv)
		if err != nil {
			return nil, err
		}
		return getCachedCelestialByCoord(coordinate)
	case ogame.Coordinate:
		return getCachedCelestialByCoord(vv)
	case string:
		coord, err := ogame.ParseCoord(vv)
		if err != nil {
			return nil, err
		}
		return getCachedCelestialByCoord(coord)
	}
	return nil, ErrIntoCelestial
}

var ErrIntoCelestial = errors.New("unable to find celestial")

// getCachedCelestialByID return celestial from cached value
func (b *OGame) getCachedCelestialByID(celestialID ogame.CelestialID) (Celestial, error) {
	return b.getCelestialByPredicateFn(func(c Celestial) bool { return c.GetID() == celestialID })
}

// getCachedCelestialByCoord return celestial from cached value
func (b *OGame) getCachedCelestialByCoord(coord ogame.Coordinate) (Celestial, error) {
	return b.getCelestialByPredicateFn(func(c Celestial) bool { return c.GetCoordinate().Equal(coord) })
}

func (b *OGame) getCelestialByPredicateFn(clb func(Celestial) bool) (Celestial, error) {
	celestials := b.getCachedCelestials()
	for _, c := range celestials {
		if clb(c) {
			return c, nil
		}
	}
	return nil, ErrIntoCelestial
}

func (b *OGame) getCachedPlanets() []Planet {
	return b.cache.planets.Load()
}

func (b *OGame) getCachedMoons() []Moon {
	var moons []Moon
	for _, p := range b.getCachedPlanets() {
		if p.Moon != nil {
			moons = append(moons, *p.Moon)
		}
	}
	return moons
}

func (b *OGame) getCachedCelestials() []Celestial {
	celestials := make([]Celestial, 0)
	for _, p := range b.getCachedPlanets() {
		celestials = append(celestials, p)
		if p.Moon != nil {
			celestials = append(celestials, *p.Moon)
		}
	}
	return celestials
}

func (b *OGame) getCachedPlanet(v IntoPlanet) (Planet, error) {
	if c, err := b.getCachedCelestial(v); err == nil {
		if planet, ok := c.(Planet); ok {
			return planet, nil
		}
	}
	return Planet{}, errors.New("invalid planet")
}

func (b *OGame) getCachedMoon(v IntoMoon) (Moon, error) {
	if c, err := b.getCachedCelestial(v); err == nil {
		if moon, ok := c.(Moon); ok {
			return moon, nil
		}
	}
	return Moon{}, errors.New("invalid moon")
}

func (b *OGame) getTasks() (out taskRunner.TasksOverview) {
	return b.taskRunnerInst.GetTasks()
}

func (b *OGame) sendDiscoveryFleet(celestialID ogame.CelestialID, coord ogame.Coordinate, options ...Option) error {
	options = append(options, ChangePlanet(celestialID))
	// Check if the sendDiscoveryFleet button is available for the target.
	// This checks for the envoys technology, if the planet has enough resources, if there's fleet slots available and if there's no cooldown on the position.
	galaxyPage, err := b.getGalaxyPage(coord.Galaxy, coord.System, options...)
	if err != nil {
		return err
	}
	for _, position := range galaxyPage.System.GalaxyContent {
		if position.Position == coord.Position {
			for _, availableMission := range position.AvailableMissions {
				if availableMission.MissionType == ogame.SearchForLifeforms {
					errMsg := "can't send discovery mission"
					if canSend, ok := availableMission.CanSend.(bool); ok {
						if !canSend {
							return errors.New(errMsg)
						}
					} else if canSendStr, ok := availableMission.CanSend.(string); ok {
						return errors.New(errMsg + ": " + canSendStr)
					}
				}
			}
		}
	}

	// Send fleet.
	res, err := b.postPageContent(url.Values{
		"page":      {"ingame"},
		"component": {"fleetdispatch"},
		"action":    {"sendDiscoveryFleet"},
		"ajax":      {"1"},
		"asJson":    {"1"},
	}, url.Values{
		"galaxy":   {utils.FI64(coord.Galaxy)},
		"system":   {utils.FI64(coord.System)},
		"position": {utils.FI64(coord.Position)},
		"token":    {galaxyPage.Token},
	})
	if err != nil {
		return err
	}

	var resStruct struct {
		Response struct {
			Success bool   `json:"success"`
			Message string `json:"message"`
		} `json:"response"`
		Components   []any  `json:"components"`
		NewAjaxToken string `json:"newAjaxToken"`
	}
	if err := json.Unmarshal(res, &resStruct); err != nil {
		return errors.New("failed to unmarshal response: " + err.Error())
	}
	if !resStruct.Response.Success {
		return errors.New(resStruct.Response.Message)
	}
	return nil
}

func (b *OGame) sendDiscoveryFleet2(celestialID ogame.CelestialID, coord ogame.Coordinate, options ...Option) (ogame.Fleet, error) {
	if err := b.sendDiscoveryFleet(celestialID, coord, options...); err != nil {
		return ogame.Fleet{}, err
	}
	select {
	case <-time.After(utils.RandMs(250, 500)):
	case <-b.parentCtx.Done():
		return ogame.Fleet{}, b.parentCtx.Err()
	case <-b.ctx.Done():
		return ogame.Fleet{}, ogame.ErrBotInactive
	}
	c, err := b.getCachedCelestial(celestialID)
	if err != nil {
		return ogame.Fleet{}, err
	}
	fleet, err := b.getLastFleetFor(c.GetCoordinate(), coord, ogame.SearchForLifeforms)
	if err != nil {
		return ogame.Fleet{}, err
	}
	return fleet, nil
}

func (b *OGame) sendSystemDiscoveryFleet(celestialID ogame.CelestialID, galaxy, system int64, options ...Option) ([]ogame.Coordinate, error) {
	options = append(options, ChangePlanet(celestialID))
	galaxyPage, err := b.getGalaxyPage(galaxy, system, options...)
	if err != nil {
		return nil, err
	}
	if _, ok := galaxyPage.System.CanSendSystemDiscovery.(string); ok {
		return nil, errors.New("can't send system discovery")
	}
	// Send fleets.
	res, err := b.postPageContent(url.Values{
		"page":      {"ingame"},
		"component": {"fleetdispatch"},
		"action":    {"sendSystemDiscoveryFleet"},
		"asJson":    {"1"},
	}, url.Values{
		"galaxy": {utils.FI64(galaxy)},
		"system": {utils.FI64(system)},
		"token":  {galaxyPage.Token},
	})
	if err != nil {
		return nil, err
	}
	var resStruct struct {
		Response struct {
			Message           string `json:"message"`
			ShipsSent         int    `json:"shipsSent"`
			SentToCoordinates []struct {
				Galaxy   int64 `json:"galaxy"`
				System   int64 `json:"system"`
				Position int64 `json:"position"`
			} `json:"sentToCoordinates"`
			Discovery struct {
				CanSendDiscovery any    `json:"canSendDiscovery"` // `true` or `"Maximum number of fleets reached."`
				DiscoveryCount   string `json:"discoveryCount"`
				GalaxyHeader     struct {
					LocaGalaxyLifeformDiscoveryCount string `json:"LOCA_GALAXY_LIFEFORM_DISCOVERY_COUNT"`
				} `json:"galaxyHeader"`
			} `json:"discovery"`
			Success bool `json:"success"`
		} `json:"response"`
		NewAjaxToken string `json:"newAjaxToken"`
	}
	if err := json.Unmarshal(res, &resStruct); err != nil {
		return nil, errors.New("failed to unmarshal response: " + err.Error())
	}
	if !resStruct.Response.Success {
		return nil, errors.New(resStruct.Response.Message)
	}
	coordinates := make([]ogame.Coordinate, 0)
	for _, resCoord := range resStruct.Response.SentToCoordinates {
		coord := ogame.NewPlanetCoordinate(resCoord.Galaxy, resCoord.System, resCoord.Position)
		coordinates = append(coordinates, coord)
	}
	return coordinates, nil
}

func (b *OGame) getAvailableDiscoveries(opts ...Option) (int64, error) {
	// Return the amount of available discoveries.
	pageHTML, err := b.getPageContent(url.Values{
		"page":      {"ingame"},
		"component": {"galaxy"},
	}, opts...)
	if err != nil {
		return 0, err
	}
	return b.extractor.ExtractAvailableDiscoveries(pageHTML)
}

type GalaxyPageContent struct {
	System struct {
		CanSendSystemDiscovery any `json:"canSendSystemDiscovery"`
		GalaxyContent          []struct {
			Position          int64 `json:"position"`
			AvailableMissions []struct {
				CanSend     any             `json:"canSend,omitempty"`
				MissionType ogame.MissionID `json:"missionType,omitempty"`
			} `json:"availableMissions"`
		} `json:"galaxyContent"`
	} `json:"system"`
	Token string `json:"token"`
}

func (b *OGame) getPositionsAvailableForDiscoveryFleet(galaxy int64, system int64, opts ...Option) ([]ogame.Coordinate, error) {
	galaxyPage, err := b.getGalaxyPage(galaxy, system, opts...)
	if err != nil {
		return nil, err
	}

	// Loop through all AvailableMissions in all positions, add those positions that are available.
	availablePositions := make([]ogame.Coordinate, 0)
	for _, position := range galaxyPage.System.GalaxyContent {
		for _, availableMission := range position.AvailableMissions {
			if availableMission.CanSend == true {
				availablePositions = append(availablePositions, ogame.Coordinate{Galaxy: galaxy, System: system, Position: position.Position, Type: ogame.PlanetType})
			}
		}
	}

	return availablePositions, nil
}

func (b *OGame) selectLfResearchSelect(planetID ogame.PlanetID, slotNumber int64) error {
	return b.selectLfResearch(planetID, slotNumber, "select", ogame.NoID)
}

func (b *OGame) selectLfResearchRandom(planetID ogame.PlanetID, slotNumber int64) error {
	return b.selectLfResearch(planetID, slotNumber, "random", ogame.NoID)
}

func (b *OGame) selectLfResearchArtifacts(planetID ogame.PlanetID, slotNumber int64, techID ogame.ID) error {
	return b.selectLfResearch(planetID, slotNumber, "selectArtifacts", techID)
}

func (b *OGame) selectLfResearch(planetID ogame.PlanetID, slotNumber int64, action string, techID ogame.ID) error {
	if slotNumber < 1 || slotNumber > 18 {
		return errors.New("invalid slot number")
	}
	vals := url.Values{
		"page":      {"ingame"},
		"component": {"lfresearch"},
		"action":    {action},
		"asJson":    {"1"},
		"planetId":  {planetID.String()},
	}
	payload := url.Values{
		"token":      {b.cache.token},
		"slotNumber": {utils.FI64(slotNumber)},
	}
	if techID.IsSet() {
		// Ensure techID is valid for the selected slotNumber
		techIdx := slotNumber - 1
		humanTech := ogame.HumansTechnologiesIDs[techIdx]
		rocktalTech := ogame.RocktalTechnologiesIDs[techIdx]
		mechasTech := ogame.MechasTechnologiesIDs[techIdx]
		kaeleshTech := ogame.KaeleshTechnologiesIDs[techIdx]
		if !utils.InArr(techID, []ogame.ID{humanTech, rocktalTech, mechasTech, kaeleshTech}) {
			return errors.New("invalid tech id for slot")
		}
		payload.Set("technologyId", utils.FI64(techID.Int64()))
	}
	if _, err := b.postPageContent(vals, payload); err != nil {
		return err
	}
	return nil
}

func (b *OGame) freeResetTree(planetID ogame.PlanetID, tier int64) error {
	return b.resetTree(planetID, tier, "freeResetTree")
}

func (b *OGame) buyResetTree(planetID ogame.PlanetID, tier int64) error {
	return b.resetTree(planetID, tier, "buyResetTree")
}

func (b *OGame) resetTree(planetID ogame.PlanetID, tier int64, action string) error {
	if tier < 1 || tier > 3 {
		return errors.New("invalid tier")
	}
	vals := url.Values{
		"page":      {"ingame"},
		"component": {"lfresearch"},
		"action":    {action},
		"asJson":    {"1"},
		"planetId":  {planetID.String()},
	}
	payload := url.Values{
		"token": {b.cache.token},
		"tier":  {utils.FI64(tier)},
	}
	by, err := b.postPageContent(vals, payload)
	if err != nil {
		return err
	}
	var res struct {
		Status string       `json:"status"`
		Errors []OGameError `json:"errors"`
	}
	if err := json.Unmarshal(by, &res); err != nil {
		return err
	}
	if res.Status == "failure" {
		var ogameErr OGameError
		if len(res.Errors) > 0 {
			ogameErr = res.Errors[0]
		}
		return fmt.Errorf("failed to reset tree for tier%d: %s (#%d)", tier, ogameErr.Message, ogameErr.Error)
	}
	return nil
}

// OGameError ogame struct for errors
type OGameError struct {
	Message string `json:"message"`
	Error   int    `json:"error"`
}

func (b *OGame) reconnectChat() bool {
	if b.ws != nil {
		_ = websocket.Message.Send(b.ws, "1::/chat")
		return true
	}
	return false
}

func (b *OGame) setLoginWrapper(newWrapper func(LoginFn) error) {
	b.loginWrapper = newWrapper
}

func (b *OGame) distance(origin, destination ogame.Coordinate) int64 {
	serverData := b.cache.serverData
	return Distance(origin, destination, serverData.Galaxies, serverData.Systems, 0, serverData.DonutGalaxy, serverData.DonutSystem)
}

func (b *OGame) systemDistance(system1, system2 int64) int64 {
	serverData := b.cache.serverData
	return systemDistance(serverData.Systems, system1, system2, serverData.DonutSystem)
}

func (b *OGame) getSession() string {
	return b.cache.ogameSession
}

func (b *OGame) withPriority(priority taskRunner.Priority) *Prioritize {
	return b.taskRunnerInst.WithPriority(priority)
}

func (b *OGame) getDevice() *device.Device {
	return b.device
}

func (b *OGame) getClient() *httpclient.Client {
	return b.device.GetClient()
}

func (b *OGame) setClient(client *httpclient.Client) {
	b.device.SetClient(client)
}

func (b *OGame) onStateChange(clb func(locked bool, actor string)) {
	b.stateChangeCallbacks = append(b.stateChangeCallbacks, clb)
}

func (b *OGame) getState() (bool, string) {
	return b.lockedAtom.Load(), b.state
}

func (b *OGame) isLocked() bool {
	return b.lockedAtom.Load()
}

func (b *OGame) getServer() gameforge.Server {
	return b.server
}

func (b *OGame) planetID() ogame.CelestialID {
	return b.cache.planetID
}

func (b *OGame) serverURL() string {
	return b.cache.serverURL
}

func (b *OGame) getLanguage() string {
	return b.language
}

func (b *OGame) bytesDownloaded() int64 {
	return b.device.GetClient().BytesDownloaded()
}

func (b *OGame) bytesUploaded() int64 {
	return b.device.GetClient().BytesUploaded()
}

func (b *OGame) getUniverseName() string {
	return b.universe
}

func (b *OGame) getUsername() string {
	return b.username
}

func (b *OGame) isPioneers() bool {
	return b.lobby == gameforge.LobbyPioneers
}

func (b *OGame) getCachedPreferences() ogame.Preferences {
	return b.cache.CachedPreferences
}

func (b *OGame) isVacationModeEnabled() bool {
	return b.cache.isVacationModeEnabled
}

func (b *OGame) location() *time.Location {
	return b.cache.location
}

func (b *OGame) getCachedToken() string {
	return b.cache.token
}

func (b *OGame) getServerData() ServerData {
	return b.cache.serverData
}

func (b *OGame) getResearchSpeed() int64 {
	return b.cache.serverData.ResearchDurationDivisor
}

func (b *OGame) getNbSystems() int64 {
	return b.cache.serverData.Systems
}

func (b *OGame) fleetDeutSaveFactor() float64 {
	return b.cache.serverData.GlobalDeuteriumSaveFactor
}

func (b *OGame) serverVersion() string {
	return b.cache.serverData.Version
}

func (b *OGame) characterClass() ogame.CharacterClass {
	return b.cache.characterClass
}

func (b *OGame) countColonies() (int64, int64) {
	return b.cache.coloniesCount, b.cache.coloniesPossible
}

func (b *OGame) getExtractor() extractor.Extractor {
	return b.extractor
}

func (b *OGame) setAllianceClass(allianceClass ogame.AllianceClass) {
	b.cache.allianceClass = &allianceClass
}

func (b *OGame) setResearches(researches ogame.Researches) {
	b.cache.researches = &researches
}

func (b *OGame) setLfBonuses(lfBonuses ogame.LfBonuses) {
	b.cache.lfBonuses = &lfBonuses
}

func (b *OGame) registerWSCallback(id string, fn func(msg []byte)) {
	b.wsCallbacks.Insert(id, fn)
}

func (b *OGame) removeWSCallback(id string) {
	b.wsCallbacks.Delete(id)
}

func (b *OGame) registerChatCallback(fn func(msg ogame.ChatMsg)) {
	b.chatCallbacks = append(b.chatCallbacks, fn)
}

func (b *OGame) registerAuctioneerCallback(fn func(packet any)) {
	b.auctioneerCallbacks = append(b.auctioneerCallbacks, fn)
}

func (b *OGame) registerHTMLInterceptor(fn func(method, url string, params, payload url.Values, pageHTML []byte)) {
	b.interceptorCallbacks = append(b.interceptorCallbacks, fn)
}

func (b *OGame) setInitiator(_ string) Prioritizable {
	return nil
}

func (b *OGame) done() {}
