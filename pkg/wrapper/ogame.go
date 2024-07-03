package wrapper

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alaingilbert/ogame/pkg/device"
	"github.com/alaingilbert/ogame/pkg/gameforge"
	"github.com/alaingilbert/ogame/pkg/wrapper/solvers"
	err2 "github.com/pkg/errors"
	"io"
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

	"github.com/alaingilbert/ogame/pkg/exponentialBackoff"
	"github.com/alaingilbert/ogame/pkg/extractor"
	v6 "github.com/alaingilbert/ogame/pkg/extractor/v6"
	v874 "github.com/alaingilbert/ogame/pkg/extractor/v874"
	"github.com/alaingilbert/ogame/pkg/httpclient"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/parser"
	"github.com/alaingilbert/ogame/pkg/taskRunner"
	"github.com/alaingilbert/ogame/pkg/utils"

	"github.com/PuerkitoBio/goquery"
	version "github.com/hashicorp/go-version"
	cookiejar "github.com/orirawlings/persistent-cookiejar"
	lua "github.com/yuin/gopher-lua"
	"golang.org/x/net/proxy"
	"golang.org/x/net/websocket"
)

// OGame is a client for ogame.org. It is safe for concurrent use by
// multiple goroutines (thread-safe)
type OGame struct {
	sync.Mutex
	isEnabledAtom         atomic.Bool // atomic, prevent auto re login if we manually logged out
	isLoggedInAtom        atomic.Bool // atomic, prevent auto re login if we manually logged out
	isConnectedAtom       atomic.Bool // atomic, either or not communication between the bot and OGame is possible
	lockedAtom            atomic.Bool // atomic, bot state locked/unlocked
	chatConnectedAtom     atomic.Bool // atomic, either or not the chat is connected
	state                 string      // keep name of the function that currently lock the bot
	ctx                   context.Context
	cancelCtx             context.CancelFunc
	stateChangeCallbacks  []func(locked bool, actor string)
	quiet                 bool
	Player                ogame.UserInfos
	CachedPreferences     ogame.Preferences
	isVacationModeEnabled bool
	researches            *ogame.Researches
	lfBonuses             *ogame.LfBonuses
	planets               []Planet
	planetsMu             sync.RWMutex
	token                 string
	ajaxChatToken         string
	Universe              string
	Username              string
	password              string
	otpSecret             string
	bearerToken           string
	language              string
	playerID              int64
	lobby                 string
	ogameSession          string
	sessionChatCounter    int64
	server                gameforge.Server
	serverData            gameforge.ServerData
	serverVersion         *version.Version
	location              *time.Location
	serverURL             string
	logger                *log.Logger
	chatCallbacks         []func(msg ogame.ChatMsg)
	wsCallbacks           map[string]func(msg []byte)
	auctioneerCallbacks   []func(any)
	interceptorCallbacks  []func(method, url string, params, payload url.Values, pageHTML []byte)
	closeChatCh           chan struct{}
	ws                    *websocket.Conn
	taskRunnerInst        *taskRunner.TaskRunner[*Prioritize]
	loginWrapper          func(func() (bool, error)) error
	getServerDataWrapper  func(func() (gameforge.ServerData, error)) (gameforge.ServerData, error)
	loginProxyTransport   http.RoundTripper
	extractor             extractor.Extractor
	apiNewHostname        string
	characterClass        ogame.CharacterClass
	allianceClass         *ogame.AllianceClass
	hasCommander          bool
	hasAdmiral            bool
	hasEngineer           bool
	hasGeologist          bool
	hasTechnocrat         bool
	captchaCallback       solvers.CaptchaCallback
	device                *device.Device
	coloniesCount         int64
	coloniesPossible      int64
}

// Params parameters for more fine-grained initialization
type Params struct {
	Username        string
	Password        string
	BearerToken     string // Gameforge auth bearer token
	OTPSecret       string
	Universe        string
	Lang            string
	PlayerID        int64
	AutoLogin       bool
	Proxy           string
	ProxyUsername   string
	ProxyPassword   string
	ProxyType       string
	ProxyLoginOnly  bool
	TLSConfig       *tls.Config
	Lobby           string
	APINewHostname  string
	Device          *device.Device
	CaptchaCallback solvers.CaptchaCallback
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
		return gameforge.ValidateAccount(client, b.ctx, b.lobby, code)
	})
}

// New creates a new instance of OGame wrapper.
func New(deviceInst *device.Device, universe, username, password, lang string) (*OGame, error) {
	b, err := NewNoLogin(username, password, "", "", universe, lang, 0, deviceInst)
	if err != nil {
		return nil, err
	}
	if _, err := b.LoginWithExistingCookies(); err != nil {
		return nil, err
	}
	return b, nil
}

// NewWithParams create a new OGame instance with full control over the possible parameters
func NewWithParams(params Params) (*OGame, error) {
	if params.Device == nil {
		return nil, errors.New("no device defined")
	}
	b, err := NewNoLogin(params.Username, params.Password, params.OTPSecret, params.BearerToken, params.Universe, params.Lang, params.PlayerID, params.Device)
	if err != nil {
		return nil, err
	}
	b.captchaCallback = params.CaptchaCallback
	b.setOGameLobby(params.Lobby)
	b.apiNewHostname = params.APINewHostname
	if params.Proxy != "" {
		if err := b.SetProxy(params.Proxy, params.ProxyUsername, params.ProxyPassword, params.ProxyType, params.ProxyLoginOnly, params.TLSConfig); err != nil {
			return nil, err
		}
	}
	if params.AutoLogin {
		if params.BearerToken != "" {
			if _, err := b.LoginWithBearerToken(params.BearerToken); err != nil {
				return nil, err
			}
		} else {
			if _, err := b.LoginWithExistingCookies(); err != nil {
				return nil, err
			}
		}
	}
	return b, nil
}

// NewNoLogin does not auto login.
func NewNoLogin(username, password, otpSecret, bearerToken, universe, lang string, playerID int64, device *device.Device) (*OGame, error) {
	b := new(OGame)
	b.device = device
	b.getServerDataWrapper = DefaultGetServerDataWrapper
	b.loginWrapper = DefaultLoginWrapper
	b.Enable()
	b.quiet = false
	b.logger = log.New(os.Stdout, "", 0)

	b.Universe = universe
	b.SetOGameCredentials(username, password, otpSecret, bearerToken)
	b.setOGameLobby(gameforge.Lobby)
	b.language = lang
	b.playerID = playerID

	ext := v874.NewExtractor()
	ext.SetLanguage(lang)
	ext.SetLocation(time.UTC)
	b.extractor = ext

	factory := func() *Prioritize { return &Prioritize{bot: b} }
	b.taskRunnerInst = taskRunner.NewTaskRunner(context.Background(), factory)

	b.wsCallbacks = make(map[string]func([]byte))

	return b, nil
}

func findAccount(universe, lang string, playerID int64, accounts []gameforge.Account, servers []gameforge.Server) (gameforge.Account, gameforge.Server, error) {
	if lang == "ba" {
		lang = "yu"
	}
	var acc gameforge.Account
	server, found := gameforge.FindServer(universe, lang, servers)
	if !found {
		return gameforge.Account{}, gameforge.Server{}, fmt.Errorf("server %s, %s not found", universe, lang)
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
	if acc.ID == 0 {
		return gameforge.Account{}, gameforge.Server{}, ogame.ErrAccountNotFound
	}
	return acc, server, nil
}

func execLoginLink(b *OGame, loginLink string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, loginLink, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	b.debug("login to universe")
	resp, err := b.doReqWithLoginProxyTransport(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return utils.ReadBody(resp)
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
			"username":  {b.Player.PlayerName},
			"isVeteran": {"1"},
		}
		if _, err := b.postPageContent(vals, payload); err != nil {
			return err
		}
	}
	return nil
}

func postSessions(b *OGame) (out *gameforge.GFLoginRes, err error) {
	client := b.device.GetClient()
	if err := client.WithTransport(b.loginProxyTransport, func(client *httpclient.Client) error {
		var challengeID string
		maxTry := uint(3)
		for {
			params := &gameforge.GfLoginParams{
				Ctx:         b.ctx,
				Device:      b.device,
				Lobby:       b.lobby,
				Username:    b.Username,
				Password:    b.password,
				OtpSecret:   b.otpSecret,
				ChallengeID: challengeID,
			}
			out, err = gameforge.GFLogin(params)
			var captchaErr *gameforge.CaptchaRequiredError
			if errors.As(err, &captchaErr) {
				captchaCallback := b.captchaCallback
				if maxTry == 0 || captchaCallback == nil {
					return err
				}
				maxTry--
				challengeID = captchaErr.ChallengeID
				if err := solveCaptcha(client, b.ctx, challengeID, captchaCallback); err != nil {
					return err
				}
				continue
			} else if err != nil {
				return err
			}
			break
		}
		return nil
	}); err != nil {
		return nil, err
	}

	bearerToken := out.Token

	// put in cookie jar so that we can re-login reusing the cookies
	appendCookie(client, &http.Cookie{
		Name:   gameforge.TokenCookieName,
		Value:  bearerToken,
		Path:   "/",
		Domain: ".gameforge.com",
	})
	b.bearerToken = bearerToken
	return out, nil
}

func appendCookie(client *httpclient.Client, cookie *http.Cookie) {
	u, _ := url.Parse("https://gameforge.com")
	cookies := client.Jar.Cookies(u)
	cookies = append(cookies, cookie)
	client.Jar.SetCookies(u, cookies)
}

func solveCaptcha(client httpclient.IHttpClient, ctx context.Context, challengeID string, captchaCallback solvers.CaptchaCallback) error {
	questionRaw, iconsRaw, err := gameforge.StartCaptchaChallenge(client, ctx, challengeID)
	if err != nil {
		return errors.New("failed to start captcha challenge: " + err.Error())
	}
	answer, err := captchaCallback(questionRaw, iconsRaw)
	if err != nil {
		return errors.New("failed to get answer for captcha challenge: " + err.Error())
	}
	if err := gameforge.SolveChallenge(client, ctx, challengeID, answer); err != nil {
		return errors.New("failed to solve captcha challenge: " + err.Error())
	}
	return err
}

func convertPlanets(b *OGame, planetsIn []ogame.Planet) []Planet {
	out := make([]Planet, 0)
	for _, planet := range planetsIn {
		out = append(out, convertPlanet(b, planet))
	}
	return out
}

func convertPlanet(b *OGame, planet ogame.Planet) Planet {
	newPlanet := Planet{ogame: b, Planet: planet}
	if planet.Moon != nil {
		moon := convertMoon(b, *planet.Moon)
		newPlanet.Moon = &moon
	}
	return newPlanet
}

func convertMoons(b *OGame, moonsIn []ogame.Moon) []Moon {
	out := make([]Moon, 0)
	for _, moon := range moonsIn {
		cMoon := convertMoon(b, moon)
		out = append(out, cMoon)
	}
	return out
}

func convertMoon(b *OGame, moonIn ogame.Moon) Moon {
	return Moon{ogame: b, Moon: moonIn}
}

func convertCelestials(b *OGame, celestials []ogame.Celestial) []Celestial {
	out := make([]Celestial, 0)
	for _, celestial := range celestials {
		out = append(out, convertCelestial(b, celestial))
	}
	return out
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
	b.planetsMu.Lock()
	b.planets = convertPlanets(b, page.ExtractPlanets())
	b.planetsMu.Unlock()
	b.isVacationModeEnabled = page.ExtractIsInVacation()
	b.token, _ = page.ExtractToken()
	b.ajaxChatToken, _ = page.ExtractAjaxChatToken()
	b.characterClass, _ = page.ExtractCharacterClass()
	b.hasCommander = page.ExtractCommander()
	b.hasAdmiral = page.ExtractAdmiral()
	b.hasEngineer = page.ExtractEngineer()
	b.hasGeologist = page.ExtractGeologist()
	b.hasTechnocrat = page.ExtractTechnocrat()
	b.coloniesCount, b.coloniesPossible = page.ExtractColonies()

	switch castedPage := page.(type) {
	case *parser.OverviewPage:
		var err error
		b.Player, err = castedPage.ExtractUserInfos()
		if err != nil {
			b.error(err)
		}
	case *parser.PreferencesPage:
		b.CachedPreferences = castedPage.ExtractPreferences()
	case *parser.ResearchPage:
		researches := castedPage.ExtractResearch()
		b.researches = &researches
	case *parser.LfBonusesPage:
		if bonuses, err := castedPage.ExtractLfBonuses(); err == nil {
			b.lfBonuses = &bonuses
		}
	}
}

// DefaultGetServerDataWrapper ...
var DefaultGetServerDataWrapper = func(getServerDataFn func() (gameforge.ServerData, error)) (gameforge.ServerData, error) {
	return getServerDataFn()
}

// DefaultLoginWrapper ...
var DefaultLoginWrapper = func(loginFn func() (bool, error)) error {
	_, err := loginFn()
	return err
}

// GetExtractor gets extractor object
func (b *OGame) GetExtractor() extractor.Extractor {
	return b.extractor
}

// SetOGameCredentials sets ogame credentials for the bot
func (b *OGame) SetOGameCredentials(username, password, otpSecret, bearerToken string) {
	b.Username = username
	b.password = password
	b.otpSecret = otpSecret
	b.bearerToken = bearerToken
}

func (b *OGame) setOGameLobby(lobby string) {
	if lobby != gameforge.LobbyPioneers {
		lobby = gameforge.Lobby
	}
	b.lobby = lobby
}

// SetGetServerDataWrapper ...
func (b *OGame) SetGetServerDataWrapper(newWrapper func(func() (gameforge.ServerData, error)) (gameforge.ServerData, error)) {
	b.getServerDataWrapper = newWrapper
}

// SetLoginWrapper ...
func (b *OGame) SetLoginWrapper(newWrapper func(func() (bool, error)) error) {
	b.loginWrapper = newWrapper
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
	if proxyType == "" {
		proxyType = "socks5"
	}
	if proxyAddress == "" {
		b.loginProxyTransport = nil
		client.SetTransport(http.DefaultTransport)
		return nil
	}
	transport, err := getTransport(proxyAddress, username, password, proxyType, config)
	b.loginProxyTransport = transport
	if loginOnly {
		client.SetTransport(http.DefaultTransport)
	} else {
		client.SetTransport(transport)
	}
	return err
}

func (b *OGame) connectChat(chatRetry *exponentialBackoff.ExponentialBackoff, host, port string) {
	b.connectChatV8(chatRetry, host, port)
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

func (b *OGame) connectChatV8(chatRetry *exponentialBackoff.ExponentialBackoff, host, port string) {
	token := yeast(time.Now().UnixNano() / 1000000)
	req, err := http.NewRequest(http.MethodGet, "https://"+host+":"+port+"/socket.io/?EIO=4&transport=polling&t="+token, nil)
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
	chatRetry.Reset()
	by, _ := io.ReadAll(resp.Body)
	m := regexp.MustCompile(`"sid":"([^"]+)"`).FindSubmatch(by)
	if len(m) != 2 {
		b.error("failed to get websocket sid:", err)
		return
	}
	sid := string(m[1])

	origin := "https://" + host + ":" + port + "/"
	wssURL := "wss://" + host + ":" + port + "/socket.io/?EIO=4&transport=websocket&sid=" + sid
	b.ws, err = websocket.Dial(wssURL, "", origin)
	if err != nil {
		b.error("failed to dial websocket:", err)
		return
	}
	_ = websocket.Message.Send(b.ws, "2probe")

	// Recv msgs
LOOP:
	for {
		select {
		case <-b.closeChatCh:
			break LOOP
		default:
		}

		var buf string
		if err := b.ws.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
			b.error("failed to set read deadline:", err)
		}
		err := websocket.Message.Receive(b.ws, &buf)
		if err != nil {
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
		b.Lock()
		for _, clb := range b.wsCallbacks {
			go clb([]byte(buf))
		}
		b.Unlock()
		if buf == "3probe" {
			_ = websocket.Message.Send(b.ws, "5")
			_ = websocket.Message.Send(b.ws, "40/chat,")
			_ = websocket.Message.Send(b.ws, "40/auctioneer,")
		} else if buf == "2" {
			_ = websocket.Message.Send(b.ws, "3")
		} else if regexp.MustCompile(`40/auctioneer,{"sid":"[^"]+"}`).MatchString(buf) {
			b.debug("got auctioneer sid")
		} else if regexp.MustCompile(`40/chat,{"sid":"[^"]+"}`).MatchString(buf) {
			b.debug("got chat sid")
			_ = websocket.Message.Send(b.ws, `42/chat,`+utils.FI64(b.sessionChatCounter)+`["authorize","`+b.ogameSession+`"]`)
			b.sessionChatCounter++
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
				b.error("unknown message received:", buf)
				continue
			}
			if name, ok := out[0].(string); ok {
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
							doc, _ := goquery.NewDocumentFromReader(strings.NewReader(timeLeftMsg))
							rgx := regexp.MustCompile(`\d+`)
							txt := rgx.FindString(doc.Find("b").Text())
							approx := utils.DoParseI64(txt)
							pck = ogame.AuctioneerTimeRemaining{Approx: approx * 60}
						} else if strings.Contains(timeLeftMsg, "nextAuction") {
							doc, _ := goquery.NewDocumentFromReader(strings.NewReader(timeLeftMsg))
							rgx := regexp.MustCompile(`\d+`)
							txt := rgx.FindString(doc.Find("span").Text())
							secs := utils.DoParseI64(txt)
							pck = ogame.AuctioneerNextAuction{Secs: secs}
						}
					}
				} else if name == "new auction" {
					if firstArg, ok := arg.(map[string]any); ok {
						pck1 := ogame.AuctioneerNewAuction{
							AuctionID: int64(utils.DoCastF64(firstArg["auctionId"])),
						}
						if infoMsg, ok := firstArg["info"].(string); ok {
							doc, _ := goquery.NewDocumentFromReader(strings.NewReader(infoMsg))
							rgx := regexp.MustCompile(`\d+`)
							txt := rgx.FindString(doc.Find("b").Text())
							approx := utils.DoParseI64(txt)
							pck1.Approx = approx * 60
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

// ReconnectChat ...
func (b *OGame) ReconnectChat() bool {
	if b.ws == nil {
		return false
	}
	_ = websocket.Message.Send(b.ws, "1::/chat")
	return true
}

func (b *OGame) logout() {
	_, _ = b.getPage(LogoutPageName)
	_ = b.device.GetClient().Jar.(*cookiejar.Jar).Save()
	if b.isLoggedInAtom.CompareAndSwap(true, false) {
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
	return ajax == "1" ||
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
	if !b.IsEnabled() {
		return ogame.ErrBotInactive
	}
	if !b.IsLoggedIn() {
		return ogame.ErrBotLoggedOut
	}
	if b.serverURL == "" {
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
	resp, err := b.device.GetClient().Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusInternalServerError {
		return []byte{}, err
	}
	by, err := utils.ReadBody(resp)
	if err != nil {
		return []byte{}, err
	}
	return by, nil
}

func getPageName(vals url.Values) string {
	page := vals.Get("page")
	component := vals.Get("component")
	if page == "ingame" ||
		page == "ajax" ||
		(page == "componentOnly" && component == FetchEventboxAjaxPageName) ||
		(page == "componentOnly" && component == EventListAjaxPageName && vals.Get("action") != "fetchEventBox") {
		page = component
	}
	return page
}

func getOptions(opts ...Option) (out Options) {
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(&out)
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
	finalURL := b.serverURL + "/game/index.php?" + vals.Encode()

	allianceID := vals.Get("allianceId")
	if allianceID != "" {
		finalURL = b.serverURL + "/game/allianceInfo.php?allianceId=" + allianceID
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
		if method == http.MethodPost {
			// Needs to be inside the withRetry, so if we need to re-login the redirect is back for the login call
			// Prevent redirect (301) https://stackoverflow.com/a/38150816/4196220
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

func applyDelay(b *OGame, delay time.Duration) error {
	if delay > 0 {
		select {
		case <-time.After(delay):
		case <-b.ctx.Done():
			return ogame.ErrBotInactive
		}
	}
	return nil
}

func alterPayload(method string, b *OGame, vals, payload url.Values) {
	switch method {
	case http.MethodPost:
		if vals.Get("page") == "ajaxChat" && payload.Get("mode") == "1" {
			payload.Set("token", b.ajaxChatToken)
		}
	}
}

func processResponseHTML(method string, b *OGame, pageHTML []byte, page string, payload, vals url.Values, SkipCacheFullPage bool) error {
	switch method {
	case http.MethodGet:
		if !IsAjaxPage(vals) && !IsEmpirePage(vals) && v6.IsLogged(pageHTML) {
			if !SkipCacheFullPage {
				parsedFullPage := parser.AutoParseFullPage(b.extractor, pageHTML)
				b.cacheFullPageInfo(parsedFullPage)
			}
		} else if IsAjaxPage(vals) && vals.Get("component") == "alliance" && vals.Get("tab") == "overview" && vals.Get("action") == "fetchOverview" {
			if !SkipCacheFullPage {
				var res parser.AllianceOverviewTabRes
				if err := json.Unmarshal(pageHTML, &res); err == nil {
					allianceClass, _ := b.extractor.ExtractAllianceClass([]byte(res.Content.AllianceAllianceOverview))
					b.allianceClass = &allianceClass
					b.token = res.NewAjaxToken
				}
			}
		} else if IsAjaxPage(vals) {
			var res struct {
				NewAjaxToken string `json:"newAjaxToken"`
			}
			if err := json.Unmarshal(pageHTML, &res); err == nil {
				if res.NewAjaxToken != "" {
					b.token = res.NewAjaxToken
				}
			}
		}

	case http.MethodPost:
		if page == PreferencesPageName {
			b.token, _ = b.extractor.ExtractToken(pageHTML)
			b.CachedPreferences = b.extractor.ExtractPreferences(pageHTML)
		} else if page == "ajaxChat" && (payload.Get("mode") == "1" || payload.Get("mode") == "3") {
			if err := extractNewChatToken(b, pageHTML); err != nil {
				return err
			}
		} else if IsAjaxPage(vals) {
			var res struct {
				NewAjaxToken string `json:"newAjaxToken"`
			}
			if err := json.Unmarshal(pageHTML, &res); err == nil {
				if res.NewAjaxToken != "" {
					b.token = res.NewAjaxToken
				}
			}
		}
	}
	return nil
}

func extractNewChatToken(b *OGame, pageHTMLBytes []byte) error {
	var res ChatPostResp
	if err := json.Unmarshal(pageHTMLBytes, &res); err != nil {
		return err
	}
	b.ajaxChatToken = res.NewToken
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
	retryInterval := 1
	retry := func(err error) error {
		b.error(err.Error())
		select {
		case <-time.After(time.Duration(retryInterval) * time.Second):
		case <-b.ctx.Done():
			return ogame.ErrBotInactive
		}
		retryInterval *= 2
		if retryInterval > 60 {
			retryInterval = 60
		}
		return nil
	}

	for {
		err := fn()
		if err == nil {
			break
		}
		// If we manually logged out, do not try to auto re login.
		if !b.IsEnabled() {
			return ogame.ErrBotInactive
		}
		if !b.IsLoggedIn() {
			return ogame.ErrBotLoggedOut
		}
		maxRetry--
		if maxRetry <= 0 {
			return err2.Wrap(err, ogame.ErrFailedExecuteCallback.Error())
		}

		if retryErr := retry(err); retryErr != nil {
			return retryErr
		}

		if errors.Is(err, ogame.ErrNotLogged) {
			if _, loginErr := b.wrapLoginWithExistingCookies(); loginErr != nil {
				b.error(loginErr.Error()) // log error
				if errors.Is(loginErr, ogame.ErrAccountNotFound) ||
					errors.Is(loginErr, ogame.ErrAccountBlocked) ||
					errors.Is(loginErr, ogame.ErrBadCredentials) ||
					errors.Is(loginErr, ogame.ErrOTPRequired) ||
					errors.Is(loginErr, ogame.ErrOTPInvalid) {
					return loginErr
				}
			}
		}
	}
	return nil
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
	return obj.ConstructionTime(nbr, b.getUniverseSpeed(), facilities, lfBonuses, b.characterClass, b.hasTechnocrat)
}

func (b *OGame) enable() {
	b.ctx, b.cancelCtx = context.WithCancel(context.Background())
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

func (b *OGame) isCollector() bool {
	return b.characterClass == ogame.Collector
}

func (b *OGame) isGeneral() bool {
	return b.characterClass == ogame.General
}

func (b *OGame) isDiscoverer() bool {
	return b.characterClass == ogame.Discoverer
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
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
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
		"token": {b.ajaxChatToken},
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
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(string(bodyBytes)))
	if doc.Find("title").Text() == "OGame Lobby" {
		return ogame.ErrNotLogged
	}
	var res ChatPostResp
	if err := json.Unmarshal(bodyBytes, &res); err != nil {
		return err
	}
	b.ajaxChatToken = res.NewToken
	return nil
}

func (b *OGame) getFleetsFromEventList() []ogame.Fleet {
	pageHTML, _ := b.getPageContent(url.Values{"eventList": {"movement"}, "ajax": {"1"}})
	return b.extractor.ExtractFleetsFromEventList(pageHTML)
}

func (b *OGame) getFleets(opts ...Option) ([]ogame.Fleet, ogame.Slots) {
	page, err := getPage[parser.MovementPage](b, opts...)
	if err != nil {
		return []ogame.Fleet{}, ogame.Slots{}
	}
	fleets := page.ExtractFleets()
	slots, _ := page.ExtractSlots()
	return fleets, slots
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
	page, _ := getPage[parser.MovementPage](b)
	fleets := page.ExtractFleets()
	return getLastFleetFor(fleets, origin, destination, mission)
}

func getLastFleetFor(fleets []ogame.Fleet, origin, destination ogame.Coordinate, mission ogame.MissionID) (ogame.Fleet, error) {
	if len(fleets) > 0 {
		maxV := ogame.MakeFleet()
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
	return ogame.Fleet{}, errors.New("could not find fleet")
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
	return utils.MinInt(system2-system1, (system1+nbSystems)-system2)
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
func Distance(c1, c2 ogame.Coordinate, universeSize, nbSystems int64, donutGalaxy, donutSystem bool) (distance int64) {
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

func findSlowestSpeed(ships ogame.ShipsInfos, techs ogame.Researches, lfBonuses ogame.LfBonuses, characterClass ogame.CharacterClass, allianceClass ogame.AllianceClass) int64 {
	var minSpeed int64 = math.MaxInt64
	for _, ship := range ogame.Ships {
		shipID := ship.GetID()
		if shipID == ogame.SolarSatelliteID || shipID == ogame.CrawlerID {
			continue
		}
		shipSpeed := ship.GetSpeed(techs, lfBonuses, characterClass, allianceClass)
		if ships.ByID(shipID) > 0 && shipSpeed < minSpeed {
			minSpeed = shipSpeed
		}
	}
	return minSpeed
}

func calcFuel(ships ogame.ShipsInfos, dist, duration int64, universeSpeedFleet, fleetDeutSaveFactor float64, techs ogame.Researches,
	lfBonuses ogame.LfBonuses, characterClass ogame.CharacterClass, allianceClass ogame.AllianceClass) (fuel int64) {
	tmpFn := func(baseFuel, nbr, shipSpeed int64) float64 {
		tmpSpeed := (35000 / (float64(duration)*universeSpeedFleet - 10)) * math.Sqrt(float64(dist)*10/float64(shipSpeed))
		return float64(baseFuel*nbr*dist) / 35000 * math.Pow(tmpSpeed/10+1, 2)
	}
	tmpFuel := 0.0
	for _, ship := range ogame.Ships {
		shipID := ship.GetID()
		if shipID == ogame.SolarSatelliteID || shipID == ogame.CrawlerID {
			continue
		}
		nbr := ships.ByID(shipID)
		if nbr > 0 {
			getFuelConsumption := ship.GetFuelConsumption(techs, lfBonuses, characterClass, fleetDeutSaveFactor)
			speed := ship.GetSpeed(techs, lfBonuses, characterClass, allianceClass)
			tmpFuel += tmpFn(getFuelConsumption, nbr, speed)
		}
	}
	fuel = int64(1 + math.Round(tmpFuel))
	return
}

// CalcFlightTime ...
// Systems that are empty/inactive can be skipped for distance calculation
// (server settings: fleetIgnoreEmptySystems, fleetIgnoreInactiveSystems)
// https://board.en.ogame.gameforge.com/index.php?thread/838751-flight-time-consumption-ignores-empty-inactive-systems
// speed: 1 -> 100% | 0.5 -> 50% | 0.05 -> 5%
func CalcFlightTime(origin, destination ogame.Coordinate, universeSize, nbSystems int64, donutGalaxy, donutSystem bool,
	fleetDeutSaveFactor, speed float64, universeSpeedFleet int64, ships ogame.ShipsInfos, techs ogame.Researches, lfBonuses ogame.LfBonuses,
	characterClass ogame.CharacterClass, allianceClass ogame.AllianceClass) (secs, fuel int64) {
	if !ships.HasShips() {
		return
	}
	v := findSlowestSpeed(ships, techs, lfBonuses, characterClass, allianceClass)
	secs = CalcFlightTimeWithBaseSpeed(origin, destination, universeSize, nbSystems, donutGalaxy, donutSystem, speed, v, universeSpeedFleet)
	d := float64(Distance(origin, destination, universeSize, nbSystems, donutGalaxy, donutSystem))
	fuel = calcFuel(ships, int64(d), secs, float64(universeSpeedFleet), fleetDeutSaveFactor, techs, lfBonuses, characterClass, allianceClass)
	return
}

// CalcFlightTimeWithBaseSpeed ...
// baseSpeed is the speed of the slowest ship in a fleet
// speed: 1 -> 100% | 0.5 -> 50% | 0.05 -> 5%
func CalcFlightTimeWithBaseSpeed(origin, destination ogame.Coordinate, universeSize, nbSystems int64, donutGalaxy, donutSystem bool, speed float64, baseSpeed, universeSpeedFleet int64) (secs int64) {
	s := speed
	v := float64(baseSpeed)
	a := float64(universeSpeedFleet)
	d := float64(Distance(origin, destination, universeSize, nbSystems, donutGalaxy, donutSystem))
	return int64(math.Round(((3500/s)*math.Sqrt(d*10/v) + 10) / a))
}

// CalcFlightTime calculates the flight time and the fuel consumption
func (b *OGame) CalcFlightTime(origin, destination ogame.Coordinate, speed float64, ships ogame.ShipsInfos, missionID ogame.MissionID) (secs, fuel int64) {
	lfBonuses, _ := b.GetCachedLfBonuses()
	allianceClass, _ := b.GetCachedAllianceClass()
	return CalcFlightTime(origin, destination, b.serverData.Galaxies, b.serverData.Systems, b.serverData.DonutGalaxy,
		b.serverData.DonutSystem, b.serverData.GlobalDeuteriumSaveFactor, speed, GetFleetSpeedForMission(b.serverData, missionID), ships,
		b.GetCachedResearch(), lfBonuses, b.characterClass, allianceClass)
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
	resources := b.extractor.ExtractResources(moonFacilitiesHTML)
	moonFacilities, _ := b.extractor.ExtractFacilities(moonFacilitiesHTML)
	phalanxLvl := moonFacilities.SensorPhalanx

	if phalanxLvl == 0 {
		return res, errors.New("no sensor phalanx on this moon")
	}

	// Ensure we have the resources to scan the planet
	if resources.Deuterium < ogame.SensorPhalanx.ScanConsumption() {
		return res, errors.New("not enough deuterium")
	}

	// Verify that coordinate is in phalanx range
	phalanxRange := ogame.SensorPhalanx.GetRange(phalanxLvl, b.isDiscoverer())
	if moon.GetCoordinate().Galaxy != coord.Galaxy ||
		systemDistance(b.serverData.Systems, moon.GetCoordinate().System, coord.System, b.serverData.DonutSystem) > phalanxRange {
		return res, errors.New("coordinate not in phalanx range")
	}

	// Get galaxy planets information, verify coordinate is valid planet (call to ogame server)
	planetInfos, _ := b.galaxyInfos(coord.Galaxy, coord.System)
	target := planetInfos.Position(coord.Position)
	if target == nil {
		return nil, errors.New("invalid planet coordinate")
	}
	// Ensure you are not scanning your own planet
	if target.Player.ID == b.Player.PlayerID {
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
		"token":    {b.token},
	}
	page, err := getAjaxPage[parser.PhalanxAjaxPage](b, vals, ChangePlanet(moonID.Celestial()))
	if err != nil {
		return []ogame.PhalanxFleet{}, err
	}
	b.token, _ = page.ExtractPhalanxNewToken()
	return page.ExtractPhalanx()
}

func moonIDInSlice(needle ogame.MoonID, haystack []ogame.MoonID) bool {
	for _, element := range haystack {
		if needle == element {
			return true
		}
	}
	return false
}

func (b *OGame) headersForPage(url string) (http.Header, error) {
	if !b.IsEnabled() {
		return nil, ogame.ErrBotInactive
	}
	if !b.IsLoggedIn() {
		return nil, ogame.ErrBotLoggedOut
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

	req = req.WithContext(b.ctx)
	resp, err := b.device.GetClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode >= 500 {
		return nil, err
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
	_, _, dests, wait := page.ExtractJumpGate()
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
	availShips, token, dests, wait := page.ExtractJumpGate()
	if wait > 0 {
		return false, wait, fmt.Errorf("jump gate is in recharge mode for %d seconds", wait)
	}

	// Validate destination moon id
	if !moonIDInSlice(destMoonID, dests) {
		return false, 0, errors.New("destination moon id invalid")
	}

	payload := url.Values{"token": {token}, "targetSpaceObjectId": {utils.FI64(destMoonID)}}

	// Add ships to payload
	for _, s := range ogame.Ships {
		// Get the min between what is available and what we want
		nbr := utils.MinInt(ships.ByID(s.GetID()), availShips.ByID(s.GetID()))
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
	pageHTML := strings.Replace(string(pageHTMLBytes), b.serverURL, b.apiNewHostname, -1)
	return b.extractor.ExtractEmpireJSON([]byte(pageHTML))
}

func (b *OGame) createUnion(fleet ogame.Fleet, unionUsers []string) (int64, error) {
	if fleet.ID == 0 {
		return 0, errors.New("invalid fleet id")
	}
	pageHTML, _ := b.getPageContent(url.Values{"page": {"federationlayer"}, "union": {"0"}, "fleet": {utils.FI64(fleet.ID)}, "target": {utils.FI64(fleet.TargetPlanetID)}, "ajax": {"1"}})
	payload := b.extractor.ExtractFederation(pageHTML)

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
	if typ < 0 || typ > 7 {
		return out, errors.New("typ must be in [0, 7] (0:Total, 1:Economy, 2:Research, 3:Military, 4:Military Built, 5:Military Destroyed, 6:Military Lost, 7:Honor)")
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
	pageHTML, _ := b.postPageContent(vals, payload)
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
	pageHTML, _ := b.postPageContent(vals, payload)
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
		Status  string `json:"status"`
		Message string `json:"message"`
		Errors  []struct {
			Message string `json:"message"`
			Error   int64  `json:"error"`
		} `json:"errors"`
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
		Status  string `json:"status"`
		Message string `json:"message"`
		Errors  []struct {
			Message string `json:"message"`
			Error   int64  `json:"error"`
		} `json:"errors"`
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
	params := url.Values{"page": {"buffActivation"}, "ajax": {"1"}, "type": {"1"}}
	pageHTML, _ := b.getPageContent(params, ChangePlanet(celestialID))
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
	params := url.Values{"page": {"buffActivation"}, "ajax": {"1"}, "type": {"1"}}
	pageHTML, _ := b.getPageContent(params, ChangePlanet(celestialID))
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
		Message  any    `json:"message"`
		Error    bool   `json:"error"`
		NewToken string `json:"newToken"`
	}
	by, err := b.postPageContent(params, payload)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(by, &res); err != nil {
		return err
	}
	if res.Error {
		if msg, ok := res.Message.(string); ok {
			return errors.New(msg)
		}
		return errors.New("unknown error")
	}
	return err
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
	pageHTML, _ := b.postPageContent(url.Values{"page": {"ajax"}, "component": {"traderimportexport"}, "ajax": {"1"}, "action": {"takeItem"}, "asJson": {"1"}}, payload)
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
	planets := b.GetCachedPlanets()
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
	if galaxy < 1 || galaxy > b.server.Settings.UniverseSize {
		return res, fmt.Errorf("galaxy must be within [1, %d]", b.server.Settings.UniverseSize)
	}
	if system < 1 || system > b.serverData.Systems {
		return res, errors.New("system must be within [1, " + utils.FI64(b.serverData.Systems) + "]")
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
	res, err = b.extractor.ExtractGalaxyInfos(pageHTML, b.Player.PlayerName, b.Player.PlayerID, b.Player.Rank)
	if err != nil {
		if cfg.DebugGalaxy {
			fmt.Println(string(pageHTML))
		}
		return res, err
	}
	if res.Galaxy() != galaxy || res.System() != system {
		return ogame.SystemInfos{}, errors.New("not enough deuterium")
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
	url2 := b.serverURL + "/game/index.php?page=resourceSettings"
	resp, err := b.device.GetClient().PostForm(url2, payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (b *OGame) getCachedResearch() ogame.Researches {
	if b.researches == nil {
		researches, _ := b.getResearch()
		return researches
	}
	return *b.researches
}

func (b *OGame) getResearch() (out ogame.Researches, err error) {
	page, err := getPage[parser.ResearchPage](b)
	if err != nil {
		return
	}
	researches := page.ExtractResearch()
	b.researches = &researches
	return researches, nil
}

func (b *OGame) getCachedLfBonuses() (out ogame.LfBonuses, err error) {
	if b.lfBonuses == nil {
		return b.getLfBonuses()
	}
	return *b.lfBonuses, nil
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
	b.lfBonuses = &bonuses
	return bonuses, nil
}

func (b *OGame) getCachedAllianceClass() (out ogame.AllianceClass, err error) {
	if b.allianceClass == nil {
		return b.getAllianceClass()
	}
	return *b.allianceClass, nil
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
	vals := url.Values{"page": {"ingame"}, "component": {"alliance"}, "tab": {"overview"}, "action": {"fetchOverview"}, "ajax": {"1"}, "token": {token}}
	pageHTML, err = b.getPageContent(vals, SkipCacheFullPage)
	if err != nil {
		return
	}
	if len(pageHTML) == 0 {
		tmp := ogame.NoAllianceClass
		b.allianceClass = &tmp
		return *b.allianceClass, nil
	}
	var res parser.AllianceOverviewTabRes
	if err = json.Unmarshal(pageHTML, &res); err != nil {
		return
	}
	allianceClass, _ := b.extractor.ExtractAllianceClass([]byte(res.Content.AllianceAllianceOverview))
	b.allianceClass = &allianceClass
	b.token = res.NewAjaxToken
	return *b.allianceClass, nil
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

func (b *OGame) getTechs(celestialID ogame.CelestialID) (ogame.ResourcesBuildings, ogame.Facilities, ogame.ShipsInfos, ogame.DefensesInfos, ogame.Researches, ogame.LfBuildings, ogame.LfResearches, error) {
	vals := url.Values{"page": {FetchTechsName}}
	page, err := getAjaxPage[parser.FetchTechsAjaxPage](b, vals, ChangePlanet(celestialID))
	if err != nil {
		return ogame.ResourcesBuildings{}, ogame.Facilities{}, ogame.ShipsInfos{}, ogame.DefensesInfos{}, ogame.Researches{}, ogame.LfBuildings{}, ogame.LfResearches{}, err
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

// IsV7 ...
func (b *OGame) IsV7() bool {
	return len(b.ServerVersion()) > 0 && b.ServerVersion()[0] == '7'
}

// IsV8 ...
func (b *OGame) IsV8() bool {
	return len(b.ServerVersion()) > 0 && b.ServerVersion()[0] == '8'
}

// IsV9 ...
func (b *OGame) IsV9() bool {
	return len(b.ServerVersion()) > 0 && b.ServerVersion()[0] == '9'
}

// IsV10 ...
func (b *OGame) IsV10() bool {
	return len(b.ServerVersion()) > 1 && b.ServerVersion()[:2] == "10"
}

// IsV104 ...
func (b *OGame) IsV104() bool {
	return len(b.ServerVersion()) > 3 && b.ServerVersion()[:4] == "10.4"
}

// IsV11 ...
func (b *OGame) IsV11() bool {
	return len(b.ServerVersion()) > 1 && b.ServerVersion()[:2] == "11"
}

// IsVGreaterThanOrEqual ...
func (b *OGame) IsVGreaterThanOrEqual(compareVersion string) bool {
	return isVGreaterThanOrEqual(b.serverVersion, compareVersion)
}

func isVGreaterThanOrEqual(v *version.Version, compareVersion string) bool {
	return v.GreaterThanOrEqual(version.Must(version.NewVersion(compareVersion)))
}

func (b *OGame) technologyDetails(celestialID ogame.CelestialID, id ogame.ID) (ogame.TechnologyDetails, error) {
	pageHTML, _ := b.getPageContent(url.Values{
		"page":       {"ingame"},
		"component":  {"technologydetails"},
		"ajax":       {"1"},
		"action":     {"getDetails"},
		"technology": {utils.FI64(id)},
		"cp":         {utils.FI64(celestialID)},
	})
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

	pageHTML, _ = b.getPageContent(url.Values{
		"page":       {"ingame"},
		"component":  {"technologydetails"},
		"ajax":       {"1"},
		"action":     {"getDetails"},
		"technology": {utils.FI64(id)},
		"cp":         {utils.FI64(celestialID)},
	})

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

	if !b.extractor.ExtractTearDownButtonEnabled([]byte(jsonContent.Content.Technologydetails)) {
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
		amount = utils.MinInt(nbr, maximumNbr)
	}

	payload := url.Values{
		"technologyId": {utils.FI64(id)},
		"amount":       {utils.FI64(amount)},
		"mode":         {"1"},
		"token":        {token},
		"planetId":     {utils.FI64(celestialID)},
	}

	var responseStruct struct {
		JsServerlang string `json:"js_serverlang"`
		JsServerid   string `json:"js_serverid"`
		Status       string `json:"status"`
		Errors       []struct {
			Message string `json:"message"`
			Error   int    `json:"error"`
		} `json:"errors"`
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

func (b *OGame) constructionsBeingBuilt(celestialID ogame.CelestialID) (ogame.ID, int64, ogame.ID, int64, ogame.ID, int64, ogame.ID, int64) {
	page, err := getPage[parser.OverviewPage](b, ChangePlanet(celestialID))
	if err != nil {
		return ogame.ID(0), 0, ogame.ID(0), 0, ogame.ID(0), 0, ogame.ID(0), 0
	}
	return page.ExtractConstructions()
}

func (b *OGame) cancel(token string, techID, listID int64) error {
	_, _ = b.postPageContent(url.Values{"page": {"componentOnly"}, "component": {"buildlistactions"}, "action": {"cancelEntry"}, "asJson": {"1"}},
		url.Values{"technologyId": {utils.FI64(techID)}, "listId": {utils.FI64(listID)}, "token": {token}})
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

	duration, maxV, token := page.ExtractIPM()
	if maxV == 0 {
		return 0, errors.New("no missile available")
	}
	nbr = utils.MinInt(nbr, maxV)
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
	Errors []struct {
		Message string `json:"message"`
		Error   int    `json:"error"`
	} `json:"errors"`
	TargetOk     bool   `json:"targetOk"`
	Components   []any  `json:"components"`
	NewAjaxToken string `json:"newAjaxToken"`
}

func (b *OGame) sendFleet(celestialID ogame.CelestialID, ships ogame.ShipsInfos, speed ogame.Speed, where ogame.Coordinate,
	mission ogame.MissionID, resources ogame.Resources, holdingTime, unionID int64, ensure bool) (ogame.Fleet, error) {

	// Get existing fleet, so we can ensure new fleet ID is greater
	initialFleets, slots := b.getFleets()
	maxInitialFleetID := ogame.FleetID(0)
	for _, f := range initialFleets {
		if f.ID > maxInitialFleetID {
			maxInitialFleetID = f.ID
		}
	}

	if slots.IsAllSlotsInUse(mission) {
		return ogame.Fleet{}, ogame.ErrAllSlotsInUse
	}

	// Page 1 : get to fleet page
	pageHTML, err := b.getPage(FleetdispatchPageName, ChangePlanet(celestialID))
	if err != nil {
		return ogame.Fleet{}, err
	}

	fleet1Doc, err := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	if err != nil {
		return ogame.Fleet{}, err
	}
	fleet1BodyID := b.extractor.ExtractBodyIDFromDoc(fleet1Doc)
	if fleet1BodyID != FleetdispatchPageName {
		b.error(ogame.ErrInvalidPlanetID.Error()+", planetID:", celestialID)
		return ogame.Fleet{}, ogame.ErrInvalidPlanetID
	}

	if b.extractor.ExtractIsInVacationFromDoc(fleet1Doc) {
		return ogame.Fleet{}, ogame.ErrAccountInVacationMode
	}

	// Ensure we're not trying to attack/spy ourselves
	destinationIsMyOwnPlanet := false
	myCelestials, _ := b.extractor.ExtractCelestialsFromDoc(fleet1Doc)
	for _, c := range myCelestials {
		if c.GetCoordinate().Equal(where) && c.GetID() == celestialID {
			return ogame.Fleet{}, errors.New("origin and destination are the same")
		}
		if c.GetCoordinate().Equal(where) {
			destinationIsMyOwnPlanet = true
			break
		}
	}
	if destinationIsMyOwnPlanet {
		switch mission {
		case ogame.Spy:
			return ogame.Fleet{}, errors.New("you cannot spy yourself")
		case ogame.Attack:
			return ogame.Fleet{}, errors.New("you cannot attack yourself")
		}
	}

	availableShips := b.extractor.ExtractFleet1ShipsFromDoc(fleet1Doc)

	atLeastOneShipSelected := false
	if !ensure {
		ships.EachFlyable(func(shipID ogame.ID, nb int64) {
			avail := availableShips.ByID(shipID)
			nb = utils.MinInt(nb, avail)
			if nb > 0 {
				atLeastOneShipSelected = true
			}
		})
	} else {
		var err1 error
		ships.EachFlyable(func(shipID ogame.ID, nb int64) {
			avail := availableShips.ByID(shipID)
			if nb > avail {
				err1 = fmt.Errorf("not enough ships to send, %s (%d > %d)", ogame.Objs.ByID(shipID).GetName(), nb, avail)
			}
			atLeastOneShipSelected = true
		})
		if err1 != nil {
			return ogame.Fleet{}, err1
		}
	}
	if !atLeastOneShipSelected {
		return ogame.Fleet{}, ogame.ErrNoShipSelected
	}

	payload := b.extractor.ExtractHiddenFieldsFromDoc(fleet1Doc)
	ships.EachFlyable(func(shipID ogame.ID, nb int64) {
		payload.Set("am"+utils.FI64(shipID), utils.FI64(nb))
	})

	token, err := b.extractor.ExtractToken(pageHTML)
	if err != nil {
		return ogame.Fleet{}, err
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
		found := false
		acsArr := b.extractor.ExtractFleetDispatchACSFromDoc(fleet1Doc)
		for _, acs := range acsArr {
			if unionID == acs.Union {
				found = true
				payload.Add("acsValues", acs.ACSValues)
				payload.Add("union", utils.FI64(acs.Union))
				mission = ogame.GroupedAttack
				break
			}
		}
		if !found {
			return ogame.Fleet{}, ogame.ErrUnionNotFound
		}
	}

	// Check
	by1, err := b.postPageContent(url.Values{"page": {"ingame"}, "component": {"fleetdispatch"}, "action": {"checkTarget"}, "ajax": {"1"}, "asJson": {"1"}}, payload)
	if err != nil {
		b.error(err.Error())
		return ogame.Fleet{}, err
	}
	var checkRes CheckTargetResponse
	if err := json.Unmarshal(by1, &checkRes); err != nil {
		b.error(err.Error())
		return ogame.Fleet{}, err
	}

	if !checkRes.TargetOk {
		if len(checkRes.Errors) > 0 {
			return ogame.Fleet{}, errors.New(checkRes.Errors[0].Message + " (" + strconv.Itoa(checkRes.Errors[0].Error) + ")")
		}
		return ogame.Fleet{}, errors.New("target is not ok")
	}

	lfBonuses, err := b.getCachedLfBonuses()
	if err != nil {
		return ogame.Fleet{}, err
	}
	multiplier := float64(b.GetServerData().CargoHyperspaceTechMultiplier) / 100.0
	cargo := ships.Cargo(b.getCachedResearch(), lfBonuses, b.characterClass, multiplier, b.server.ProbeRaidsEnabled())
	newResources := ogame.Resources{}
	if resources.Total() > cargo {
		newResources.Deuterium = utils.MinInt(resources.Deuterium, cargo)
		cargo -= newResources.Deuterium
		newResources.Crystal = utils.MinInt(resources.Crystal, cargo)
		cargo -= newResources.Crystal
		newResources.Metal = utils.MinInt(resources.Metal, cargo)
	} else {
		newResources = resources
	}

	newResources.Metal = utils.MaxInt(newResources.Metal, 0)
	newResources.Crystal = utils.MaxInt(newResources.Crystal, 0)
	newResources.Deuterium = utils.MaxInt(newResources.Deuterium, 0)

	// Page 3 : select coord, mission, speed
	payload.Set("token", checkRes.NewAjaxToken)
	payload.Set("speed", strconv.FormatInt(int64(speed), 10))
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
		} else if mission == ogame.ParkInThatAlly { // ParkInThatAlly 0, 1, 2, 4, 8, 16, 32
			holdingTime = utils.Clamp(holdingTime, 0, 32)
		}
		payload.Set("holdingtime", utils.FI64(holdingTime))
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
		return ogame.Fleet{}, errors.New("failed to unmarshal response: " + err.Error())
	}

	if len(resStruct.Errors) > 0 {
		return ogame.Fleet{}, errors.New(resStruct.Errors[0].Message + " (" + utils.FI64(resStruct.Errors[0].Error) + ")")
	}

	// Page 5
	page, _ := getPage[parser.MovementPage](b)
	originCoords, _ := page.ExtractPlanetCoordinate()
	fleets := page.ExtractFleets()
	if maxV, err := getLastFleetFor(fleets, originCoords, where, mission); err == nil && maxV.ID > maxInitialFleetID {
		return maxV, nil
	}

	slots, _ = page.ExtractSlots()
	if slots.InUse == slots.Total {
		return ogame.Fleet{}, ogame.ErrAllSlotsInUse
	}

	if mission == ogame.Expedition {
		if slots.ExpInUse == slots.ExpTotal {
			return ogame.Fleet{}, ogame.ErrAllSlotsInUse
		}
	}

	b.error(errors.New("could not find new fleet ID").Error()+", planetID:", celestialID)
	return ogame.Fleet{}, errors.New("could not find new fleet ID")
}

func (b *OGame) miniFleetSpy(coord ogame.Coordinate, shipCount int64) error {
	token := ""
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
		"type":      {"1"}, // ?
		"shipCount": {utils.FI64(shipCount)},
		"token":     {token},
	}
	pageHTML, err := b.postPageContent(vals, payload)
	if err != nil {
		return err
	}
	var res struct {
		Response struct {
			Message     string `json:"message"`
			Type        int    `json:"type"`
			Slots       int    `json:"slots"`
			Probes      int    `json:"probes"`
			Recyclers   int    `json:"recyclers"`
			Explorers   int    `json:"explorers"`
			Missiles    int    `json:"missiles"`
			ShipsSent   int    `json:"shipsSent"`
			Coordinates struct {
				Galaxy   int `json:"galaxy"`
				System   int `json:"system"`
				Position int `json:"position"`
			} `json:"coordinates"`
			PlanetType int  `json:"planetType"`
			Success    bool `json:"success"`
		} `json:"response"`
		NewAjaxToken string `json:"newAjaxToken"`
	}
	if err := json.Unmarshal(pageHTML, &res); err != nil {
		return err
	}
	if !res.Response.Success {
		return errors.New(res.Response.Message)
	}
	return nil
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

func (b *OGame) getCombatReportFor(coord ogame.Coordinate) (ogame.CombatReportSummary, error) {
	pageHTML, err := b.getPageMessages(1, CombatReportsMessagesTabID)
	if err != nil {
		return ogame.CombatReportSummary{}, err
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
	for _, m := range newMessages {
		if m.Destination.Equal(coord) {
			return m, nil
		}
	}
	return ogame.CombatReportSummary{}, errors.New("combat report not found for " + coord.String())
}

func (b *OGame) getEspionageReport(msgID int64) (ogame.EspionageReport, error) {
	pageHTML, _ := b.getPageContent(url.Values{"page": {"componentOnly"}, "component": {"messagedetails"}, "messageId": {utils.FI64(msgID)}})
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
		"token":     {token},
		"messageId": {utils.FI64(msgID)},
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
	universeSpeed := b.serverData.Speed
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
	by, err := utils.ReadBody(resp)
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

func (b *OGame) addAccount(number int, lang string) (*gameforge.AddAccountRes, error) {
	accountGroup := fmt.Sprintf("%s_%d", lang, number)
	return gameforge.AddAccount(b.device, b.ctx, b.lobby, accountGroup, b.bearerToken)
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
	case lua.LNumber:
		return getCachedCelestialByID(ogame.CelestialID(vv))
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
	b.planetsMu.RLock()
	defer b.planetsMu.RUnlock()
	return b.planets
}

func (b *OGame) getCachedMoons() []Moon {
	var moons []Moon
	for _, p := range b.GetCachedPlanets() {
		if p.Moon != nil {
			moons = append(moons, *p.Moon)
		}
	}
	return moons
}

func (b *OGame) getCachedCelestials() []Celestial {
	celestials := make([]Celestial, 0)
	for _, p := range b.GetCachedPlanets() {
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
	return Moon{}, errors.New("invalid planet")
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

func (b *OGame) getAvailableDiscoveries(opts ...Option) int64 {
	// Return the amount of available discoveries.
	pageHTML, _ := b.getPageContent(url.Values{
		"page":      {"ingame"},
		"component": {"galaxy"},
	}, opts...)
	return b.extractor.ExtractAvailableDiscoveries(pageHTML)
}

type GalaxyPageContent struct {
	System struct {
		GalaxyContent []struct {
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
		"token":      {b.token},
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
		"token": {b.token},
		"tier":  {utils.FI64(tier)},
	}
	if _, err := b.postPageContent(vals, payload); err != nil {
		return err
	}
	return nil
}
