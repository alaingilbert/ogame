package wrapper

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	err2 "errors"
	"fmt"
	"github.com/alaingilbert/ogame/pkg/extractor"
	"github.com/alaingilbert/ogame/pkg/extractor/v6"
	"github.com/alaingilbert/ogame/pkg/extractor/v7"
	"github.com/alaingilbert/ogame/pkg/extractor/v71"
	"github.com/alaingilbert/ogame/pkg/extractor/v8"
	"github.com/alaingilbert/ogame/pkg/extractor/v874"
	"github.com/alaingilbert/ogame/pkg/extractor/v9"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/parser"
	"github.com/alaingilbert/ogame/pkg/taskRunner"
	"github.com/alaingilbert/ogame/pkg/utils"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"math"
	"mime/multipart"
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
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	version "github.com/hashicorp/go-version"
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
	ctx                   context.Context
	cancelCtx             context.CancelFunc
	stateChangeCallbacks  []func(locked bool, actor string)
	quiet                 bool
	Player                ogame.UserInfos
	CachedPreferences     ogame.Preferences
	isVacationModeEnabled bool
	researches            *ogame.Researches
	planets               []Planet
	planetsMu             sync.RWMutex
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
	server                Server
	serverData            ServerData
	location              *time.Location
	serverURL             string
	client                *OGameClient
	logger                *log.Logger
	chatCallbacks         []func(msg ogame.ChatMsg)
	wsCallbacks           map[string]func(msg []byte)
	auctioneerCallbacks   []func(any)
	interceptorCallbacks  []func(method, url string, params, payload url.Values, pageHTML []byte)
	closeChatCh           chan struct{}
	ws                    *websocket.Conn
	taskRunnerInst        *taskRunner.TaskRunner[*Prioritize]
	loginWrapper          func(func() (bool, error)) error
	getServerDataWrapper  func(func() (ServerData, error)) (ServerData, error)
	loginProxyTransport   http.RoundTripper
	extractor             extractor.Extractor
	apiNewHostname        string
	characterClass        ogame.CharacterClass
	hasCommander          bool
	hasAdmiral            bool
	hasEngineer           bool
	hasGeologist          bool
	hasTechnocrat         bool
	captchaCallback       CaptchaCallback
}

// CaptchaCallback ...
type CaptchaCallback func(question, icons []byte) (int64, error)

const defaultUserAgent = "" +
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) " +
	"AppleWebKit/537.36 (KHTML, like Gecko) " +
	"Chrome/104.0.0.0 " +
	"Safari/537.36"

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
	CookiesFilename string
	Client          *OGameClient
	CaptchaCallback CaptchaCallback
}

// Lobby constants
const (
	Lobby         = "lobby"
	LobbyPioneers = "lobby-pioneers"
)

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
	return b.client.WithTransport(b.loginProxyTransport, func(client IHttpClient) error {
		return ValidateAccount(client, b.ctx, b.lobby, code)
	})
}

// New creates a new instance of OGame wrapper.
func New(universe, username, password, lang string) (*OGame, error) {
	b, err := NewNoLogin(username, password, "", "", universe, lang, "", 0, nil)
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
	b, err := NewNoLogin(params.Username, params.Password, params.OTPSecret, params.BearerToken, params.Universe, params.Lang, params.CookiesFilename, params.PlayerID, params.Client)
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
func NewNoLogin(username, password, otpSecret, bearerToken, universe, lang, cookiesFilename string, playerID int64, client *OGameClient) (*OGame, error) {
	b := new(OGame)
	b.getServerDataWrapper = DefaultGetServerDataWrapper
	b.loginWrapper = DefaultLoginWrapper
	b.Enable()
	b.quiet = false
	b.logger = log.New(os.Stdout, "", 0)

	b.Universe = universe
	b.SetOGameCredentials(username, password, otpSecret, bearerToken)
	b.setOGameLobby(Lobby)
	b.language = lang
	b.playerID = playerID

	b.extractor = v874.NewExtractor()

	if client == nil {
		jar, err := cookiejar.New(&cookiejar.Options{
			Filename:              cookiesFilename,
			PersistSessionCookies: true,
		})
		if err != nil {
			return nil, err
		}

		// Ensure we remove any cookies that would set the mobile view
		cookies := jar.AllCookies()
		for _, c := range cookies {
			if c.Name == "device" {
				jar.RemoveCookie(c)
			}
		}

		b.client = NewOGameClient()
		b.client.Jar = jar
		b.client.SetUserAgent(defaultUserAgent)
	} else {
		b.client = client
	}

	factory := func() *Prioritize { return &Prioritize{bot: b} }
	b.taskRunnerInst = taskRunner.NewTaskRunner(b.ctx, factory)

	b.wsCallbacks = make(map[string]func([]byte))

	return b, nil
}

func findServer(universe, lang string, servers []Server) (out Server, found bool) {
	for _, s := range servers {
		if s.Name == universe && s.Language == lang {
			return s, true
		}
	}
	return
}

func findAccount(universe, lang string, playerID int64, accounts []Account, servers []Server) (Account, Server, error) {
	if lang == "ba" {
		lang = "yu"
	}
	var acc Account
	server, found := findServer(universe, lang, servers)
	if !found {
		return Account{}, Server{}, fmt.Errorf("server %s, %s not found", universe, lang)
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
		return Account{}, Server{}, ogame.ErrAccountNotFound
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

// Return either or not the bot logged in using the provided bearer token.
func (b *OGame) loginWithBearerToken(token string) (bool, error) {
	if token == "" {
		err := b.login()
		return false, err
	}
	b.bearerToken = token
	server, userAccount, err := b.loginPart1(token)
	if err2.Is(err, context.Canceled) {
		return false, err
	}
	if err == ogame.ErrAccountBlocked {
		return false, err
	}
	if err != nil {
		err := b.login()
		return false, err
	}

	if err := b.loginPart2(server); err != nil {
		return false, err
	}

	page, err := getPage[parser.OverviewPage](b, SkipRetry)
	if err != nil {
		if err == ogame.ErrNotLogged {
			b.debug("get login link")
			loginLink, err := GetLoginLink(b.client, b.ctx, b.lobby, userAccount, token)
			if err != nil {
				return true, err
			}
			pageHTML, err := execLoginLink(b, loginLink)
			if err != nil {
				return true, err
			}
			page, err := getPage[parser.OverviewPage](b, SkipRetry)
			if err != nil {
				if err == ogame.ErrNotLogged {
					err := b.login()
					return false, err
				}
			}
			b.debug("login using existing cookies")
			if err := b.loginPart3(userAccount, page); err != nil {
				return false, err
			}
			if err := b.client.Jar.(*cookiejar.Jar).Save(); err != nil {
				return false, err
			}
			for _, fn := range b.interceptorCallbacks {
				fn(http.MethodGet, loginLink, nil, nil, pageHTML)
			}
			return true, nil
		}
		return false, err
	}
	b.debug("login using existing cookies")
	if err := b.loginPart3(userAccount, page); err != nil {
		return false, err
	}
	return true, nil
}

// Return either or not the bot logged in using the existing cookies.
func (b *OGame) loginWithExistingCookies() (bool, error) {
	token := ""
	if b.bearerToken != "" {
		token = b.bearerToken
	} else {
		cookies := b.client.Jar.(*cookiejar.Jar).AllCookies()
		for _, c := range cookies {
			if c.Name == TokenCookieName {
				token = c.Value
				break
			}
		}
	}
	return b.loginWithBearerToken(token)
}

// TelegramSolver ...
func TelegramSolver(tgBotToken string, tgChatID int64) CaptchaCallback {
	return func(question, icons []byte) (int64, error) {
		tgBot, _ := tgbotapi.NewBotAPI(tgBotToken)
		keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("0", "0"),
			tgbotapi.NewInlineKeyboardButtonData("1", "1"),
			tgbotapi.NewInlineKeyboardButtonData("2", "2"),
			tgbotapi.NewInlineKeyboardButtonData("3", "3"),
		))
		questionImgOrig, _ := png.Decode(bytes.NewReader(question))
		bounds := questionImgOrig.Bounds()
		upLeft := image.Point{X: 0, Y: 0}
		lowRight := bounds.Max
		img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})
		for y := 0; y < lowRight.Y; y++ {
			for x := 0; x < lowRight.X; x++ {
				c := questionImgOrig.At(x, y)
				r, g, b, _ := c.RGBA()
				img.Set(x, y, color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: 255})
			}
		}
		buf := bytes.NewBuffer(nil)
		_ = png.Encode(buf, img)
		questionImg := tgbotapi.FileBytes{Name: "question", Bytes: buf.Bytes()}
		iconsImg := tgbotapi.FileBytes{Name: "icons", Bytes: icons}
		_, _ = tgBot.Send(tgbotapi.NewPhotoUpload(tgChatID, questionImg))
		_, _ = tgBot.Send(tgbotapi.NewPhotoUpload(tgChatID, iconsImg))
		msg := tgbotapi.NewMessage(tgChatID, "Pick one")
		msg.ReplyMarkup = keyboard
		_, _ = tgBot.Send(msg)
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		updates, _ := tgBot.GetUpdatesChan(u)
		for update := range updates {
			if update.CallbackQuery != nil {
				_, _ = tgBot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data))
				_, _ = tgBot.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "got "+update.CallbackQuery.Data))
				v, err := utils.ParseI64(update.CallbackQuery.Data)
				if err != nil {
					return 0, err
				}
				return v, nil
			}
		}
		return 0, errors.New("failed to get answer")
	}
}

// NinjaSolver direct integration of ogame.ninja captcha auto solver service
func NinjaSolver(apiKey string) CaptchaCallback {
	return func(question, icons []byte) (int64, error) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("question", "question.png")
		_, _ = io.Copy(part, bytes.NewReader(question))
		part1, _ := writer.CreateFormFile("icons", "icons.png")
		_, _ = io.Copy(part1, bytes.NewReader(icons))
		_ = writer.Close()

		req, _ := http.NewRequest(http.MethodPost, "https://www.ogame.ninja/api/v1/captcha/solve", body)
		req.Header.Add("Content-Type", writer.FormDataContentType())
		req.Header.Set("NJA_API_KEY", apiKey)
		resp, _ := http.DefaultClient.Do(req)
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			by, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return 0, errors.New("failed to auto solve captcha: " + err.Error())
			}
			return 0, errors.New("failed to auto solve captcha: " + string(by))
		}
		by, _ := ioutil.ReadAll(resp.Body)
		var answerJson struct {
			Answer int64 `json:"answer"`
		}
		if err := json.Unmarshal(by, &answerJson); err != nil {
			return 0, errors.New("failed to auto solve captcha: " + err.Error())
		}
		return answerJson.Answer, nil
	}
}

func postSessions(b *OGame, lobby, username, password, otpSecret string) (out *GFLoginRes, err error) {
	if err := b.client.WithTransport(b.loginProxyTransport, func(client IHttpClient) error {
		var challengeID string
		tried := false
		for {
			out, err = GFLogin(client, b.ctx, lobby, username, password, otpSecret, challengeID)
			var captchaErr *CaptchaRequiredError
			if errors.As(err, &captchaErr) {
				if tried || b.captchaCallback == nil {
					return err
				}
				tried = true

				questionRaw, iconsRaw, err := StartCaptchaChallenge(client, b.ctx, captchaErr.ChallengeID)
				if err != nil {
					return errors.New("failed to start captcha challenge: " + err.Error())
				}
				answer, err := b.captchaCallback(questionRaw, iconsRaw)
				if err != nil {
					return errors.New("failed to get answer for captcha challenge: " + err.Error())
				}
				if err := SolveChallenge(client, b.ctx, captchaErr.ChallengeID, answer); err != nil {
					return errors.New("failed to solve captcha challenge: " + err.Error())
				}
				challengeID = captchaErr.ChallengeID
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

	// put in cookie jar so that we can re-login reusing the cookies
	u, _ := url.Parse("https://gameforge.com")
	cookies := b.client.Jar.Cookies(u)
	cookie := &http.Cookie{
		Name:   TokenCookieName,
		Value:  out.Token,
		Path:   "/",
		Domain: ".gameforge.com",
	}
	cookies = append(cookies, cookie)
	b.client.Jar.SetCookies(u, cookies)
	fmt.Println("bearer", out.Token)
	b.bearerToken = out.Token
	return out, nil
}

func (b *OGame) login() error {
	b.debug("post sessions")
	postSessionsRes, err := postSessions(b, b.lobby, b.Username, b.password, b.otpSecret)
	if err != nil {
		return err
	}

	server, userAccount, err := b.loginPart1(postSessionsRes.Token)
	if err != nil {
		return err
	}

	b.debug("get login link")
	loginLink, err := GetLoginLink(b.client, b.ctx, b.lobby, userAccount, postSessionsRes.Token)
	if err != nil {
		return err
	}
	pageHTML, err := execLoginLink(b, loginLink)
	if err != nil {
		return err
	}

	if err := b.loginPart2(server); err != nil {
		return err
	}
	page, err := parser.ParsePage[parser.OverviewPage](b.extractor, pageHTML)
	if err != nil {
		return err
	}
	if err := b.loginPart3(userAccount, page); err != nil {
		return err
	}

	if err := b.client.Jar.(*cookiejar.Jar).Save(); err != nil {
		return err
	}
	for _, fn := range b.interceptorCallbacks {
		fn(http.MethodGet, loginLink, nil, nil, pageHTML)
	}
	return nil
}

func (b *OGame) loginPart1(token string) (server Server, userAccount Account, err error) {
	b.debug("get user accounts")
	accounts, err := GetUserAccounts(b.client, b.ctx, b.lobby, token)
	if err != nil {
		return
	}
	b.debug("get servers")
	servers, err := GetServers(b.lobby, b.client, b.ctx)
	if err != nil {
		return
	}
	b.debug("find account & server for universe")
	userAccount, server, err = findAccount(b.Universe, b.language, b.playerID, accounts, servers)
	if err != nil {
		return
	}
	if userAccount.Blocked {
		return server, userAccount, ogame.ErrAccountBlocked
	}
	b.debug("Players online: " + utils.FI64(server.PlayersOnline) + ", Players: " + utils.FI64(server.PlayerCount))
	return
}

func (b *OGame) loginPart2(server Server) error {
	atomic.StoreInt32(&b.isLoggedInAtom, 1) // At this point, we are logged in
	atomic.StoreInt32(&b.isConnectedAtom, 1)
	// Get server data
	start := time.Now()
	b.server = server
	serverData, err := b.getServerDataWrapper(func() (ServerData, error) {
		return GetServerData(b.client, b.ctx, b.server.Number, b.server.Language)
	})
	if err != nil {
		return err
	}
	if serverData.SpeedFleetWar == 0 {
		serverData.SpeedFleetWar = 1
	}
	if serverData.SpeedFleetPeaceful == 0 {
		serverData.SpeedFleetPeaceful = 1
	}
	if serverData.SpeedFleetHolding == 0 {
		serverData.SpeedFleetHolding = 1
	}
	if serverData.SpeedFleet == 0 {
		serverData.SpeedFleet = serverData.SpeedFleetPeaceful
	}
	b.serverData = serverData
	lang := server.Language
	if server.Language == "yu" {
		lang = "ba"
	}
	b.language = lang
	b.serverURL = "https://s" + utils.FI64(server.Number) + "-" + lang + ".ogame.gameforge.com"
	b.debug("get server data", time.Since(start))
	return nil
}

func (b *OGame) loginPart3(userAccount Account, page parser.OverviewPage) error {
	if ogVersion, err := version.NewVersion(b.serverData.Version); err == nil {
		if ogVersion.GreaterThanOrEqual(version.Must(version.NewVersion("9.0.0"))) {
			b.extractor = v9.NewExtractor()
		} else if ogVersion.GreaterThanOrEqual(version.Must(version.NewVersion("8.7.4-pl3"))) {
			b.extractor = v874.NewExtractor()
		} else if ogVersion.GreaterThanOrEqual(version.Must(version.NewVersion("8.0.0"))) {
			b.extractor = v8.NewExtractor()
		} else if ogVersion.GreaterThanOrEqual(version.Must(version.NewVersion("7.1.0-rc0"))) {
			b.extractor = v71.NewExtractor()
		} else if ogVersion.GreaterThanOrEqual(version.Must(version.NewVersion("7.0.0-rc0"))) {
			b.extractor = v7.NewExtractor()
		}
		b.extractor.SetLanguage(b.language)
		b.extractor.SetLocation(b.location)
		b.extractor.SetLifeformEnabled(page.ExtractLifeformEnabled())
	} else {
		b.error("failed to parse ogame version: " + err.Error())
	}

	b.sessionChatCounter = 1

	b.debug("logged in as " + userAccount.Name + " on " + b.Universe + "-" + b.language)

	b.debug("extract information from html")
	b.ogameSession = page.ExtractOGameSession()
	if b.ogameSession == "" {
		return ogame.ErrBadCredentials
	}

	serverTime, _ := page.ExtractServerTime()
	b.location = serverTime.Location()

	b.cacheFullPageInfo(page)

	_, _ = b.getPage(PreferencesPageName) // Will update preferences cached values

	// Extract chat host and port
	m := regexp.MustCompile(`var nodeUrl\s?=\s?"https:\\/\\/([^:]+):(\d+)\\/socket.io\\/socket.io.js"`).FindSubmatch(page.GetContent())
	chatHost := string(m[1])
	chatPort := string(m[2])

	if atomic.CompareAndSwapInt32(&b.chatConnectedAtom, 0, 1) {
		b.closeChatCh = make(chan struct{})
		go func(b *OGame) {
			defer atomic.StoreInt32(&b.chatConnectedAtom, 0)
			chatRetry := NewExponentialBackoff(60)
		LOOP:
			for {
				select {
				case <-b.closeChatCh:
					break LOOP
				default:
					b.connectChat(chatRetry, chatHost, chatPort)
					chatRetry.Wait()
				}
			}
		}(b)
	} else {
		b.ReconnectChat()
	}

	return nil
}

func convertPlanets(b *OGame, planetsIn []ogame.Planet) []Planet {
	out := make([]Planet, 0)
	for _, planet := range planetsIn {
		out = append(out, convertPlanet(b, planet))
	}
	return out
}

func convertPlanet(b *OGame, planet ogame.Planet) Planet {
	newPlanet := Planet{
		ogame: b,
		Planet: ogame.Planet{
			Img:         planet.Img,
			ID:          planet.ID,
			Name:        planet.Name,
			Diameter:    planet.Diameter,
			Coordinate:  planet.Coordinate,
			Fields:      planet.Fields,
			Temperature: planet.Temperature,
		},
	}
	if planet.Moon != nil {
		newPlanet.Moon = convertMoon(b, *planet.Moon)
	}
	return newPlanet
}

func convertMoons(b *OGame, moonsIn []ogame.Moon) []Moon {
	out := make([]Moon, 0)
	for _, moon := range moonsIn {
		tmp := convertMoon(b, moon)
		out = append(out, *tmp)
	}
	return out
}

func convertMoon(b *OGame, moonIn ogame.Moon) *Moon {
	return &Moon{
		ogame: b,
		Moon: ogame.Moon{
			ID:         moonIn.ID,
			Img:        moonIn.Img,
			Name:       moonIn.Name,
			Diameter:   moonIn.Diameter,
			Coordinate: moonIn.Coordinate,
			Fields:     moonIn.Fields,
		},
	}
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
	b.ajaxChatToken, _ = page.ExtractAjaxChatToken()
	b.characterClass, _ = page.ExtractCharacterClass()
	b.hasCommander = page.ExtractCommander()
	b.hasAdmiral = page.ExtractAdmiral()
	b.hasEngineer = page.ExtractEngineer()
	b.hasGeologist = page.ExtractGeologist()
	b.hasTechnocrat = page.ExtractTechnocrat()

	switch castedPage := page.(type) {
	case parser.OverviewPage:
		b.Player, _ = castedPage.ExtractUserInfos()
	case parser.PreferencesPage:
		b.CachedPreferences = castedPage.ExtractPreferences()
	case parser.ResearchPage:
		researches := castedPage.ExtractResearch()
		b.researches = &researches
	}
}

// DefaultGetServerDataWrapper ...
var DefaultGetServerDataWrapper = func(getServerDataFn func() (ServerData, error)) (ServerData, error) {
	return getServerDataFn()
}

// DefaultLoginWrapper ...
var DefaultLoginWrapper = func(loginFn func() (bool, error)) error {
	_, err := loginFn()
	return err
}

func (b *OGame) wrapLoginWithBearerToken(token string) (useToken bool, err error) {
	fn := func() (bool, error) {
		useToken, err = b.loginWithBearerToken(token)
		return useToken, err
	}
	return useToken, b.loginWrapper(fn)
}

func (b *OGame) wrapLoginWithExistingCookies() (useCookies bool, err error) {
	fn := func() (bool, error) {
		useCookies, err = b.loginWithExistingCookies()
		return useCookies, err
	}
	return useCookies, b.loginWrapper(fn)
}

func (b *OGame) wrapLogin() error {
	return b.loginWrapper(func() (bool, error) { return false, b.login() })
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
	if lobby != LobbyPioneers {
		lobby = Lobby
	}
	b.lobby = lobby
}

// SetGetServerDataWrapper ...
func (b *OGame) SetGetServerDataWrapper(newWrapper func(func() (ServerData, error)) (ServerData, error)) {
	b.getServerDataWrapper = newWrapper
}

// SetLoginWrapper ...
func (b *OGame) SetLoginWrapper(newWrapper func(func() (bool, error)) error) {
	b.loginWrapper = newWrapper
}

// execute a request using the login proxy transport if set
func (b *OGame) doReqWithLoginProxyTransport(req *http.Request) (resp *http.Response, err error) {
	req = req.WithContext(b.ctx)
	_ = b.client.WithTransport(b.loginProxyTransport, func(client IHttpClient) error {
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
	if proxyType == "" {
		proxyType = "socks5"
	}
	if proxyAddress == "" {
		b.loginProxyTransport = nil
		b.client.SetTransport(http.DefaultTransport)
		return nil
	}
	transport, err := getTransport(proxyAddress, username, password, proxyType, config)
	b.loginProxyTransport = transport
	if loginOnly {
		b.client.SetTransport(http.DefaultTransport)
	} else {
		b.client.SetTransport(transport)
	}
	return err
}

// SetProxy this will change the bot http transport object.
// proxyType can be "http" or "socks5".
// An empty proxyAddress will reset the client transport to default value.
func (b *OGame) SetProxy(proxyAddress, username, password, proxyType string, loginOnly bool, config *tls.Config) error {
	return b.setProxy(proxyAddress, username, password, proxyType, loginOnly, config)
}

func (b *OGame) connectChat(chatRetry *ExponentialBackoff, host, port string) {
	if b.IsV8() || b.IsV9() {
		b.connectChatV8(chatRetry, host, port)
	} else {
		b.connectChatV7(chatRetry, host, port)
	}
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

func (b *OGame) connectChatV8(chatRetry *ExponentialBackoff, host, port string) {
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
	by, _ := ioutil.ReadAll(resp.Body)
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
		for _, clb := range b.wsCallbacks {
			go clb([]byte(buf))
		}
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
			time.Sleep(time.Second)
		}
	}
}

func (b *OGame) connectChatV7(chatRetry *ExponentialBackoff, host, port string) {
	req, err := http.NewRequest(http.MethodGet, "https://"+host+":"+port+"/socket.io/1/?t="+utils.FI64(time.Now().UnixNano()/int64(time.Millisecond)), nil)
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
		n, err := b.ws.Read(buf)
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
		for _, clb := range b.wsCallbacks {
			go clb(buf[0:n])
		}
		msg := bytes.Trim(buf, "\x00")
		if bytes.Equal(msg, []byte("1::")) {
			_, _ = b.ws.Write([]byte("1::/chat"))       // subscribe to chat events
			_, _ = b.ws.Write([]byte("1::/auctioneer")) // subscribe to auctioneer events
		} else if bytes.Equal(msg, []byte("1::/chat")) {
			authMsg := `5:` + utils.FI64(b.sessionChatCounter) + `+:/chat:{"name":"authorize","args":["` + b.ogameSession + `"]}`
			_, _ = b.ws.Write([]byte(authMsg))
			b.sessionChatCounter++
		} else if bytes.Equal(msg, []byte("2::")) {
			_, _ = b.ws.Write([]byte("2::"))
		} else if regexp.MustCompile(`\d+::/auctioneer`).Match(msg) {
			// 5::/auctioneer:{"name":"timeLeft","args":["Next auction in:<br />\n<span class=\"nextAuction\" id=\"nextAuction\">598</span>"]}
			// 5::/auctioneer:{"name":"timeLeft","args":["<span style=\"color:#FFA500;\"><b>approx. 10m</b></span> remaining until the auction ends"]} // every minute
			// 5::/auctioneer:{"name":"new auction","args":[{"info":"<span style=\"color:#99CC00;\"><b>approx. 45m</b></span> remaining until the auction ends","item":{"uuid":"118d34e685b5d1472267696d1010a393a59aed03","image":"bdb4508609de1df58bf4a6108fff73078c89f777","rarity":"rare"},"oldAuction":{"item":{"uuid":"8a4f9e8309e1078f7f5ced47d558d30ae15b4a1b","imageSmall":"014827f6d1d5b78b1edd0d4476db05639e7d9367","rarity":"rare"},"time":"06.01.2021 17:35:05","bids":1,"sum":1000,"player":{"id":111106,"name":"Governor Skat","link":"http://s152-en.ogame.gameforge.com/game/index.php?page=ingame&component=galaxy&galaxy=1&system=218"}},"auctionId":18550}]}
			// 5::/auctioneer:{"name":"new bid","args":[{"player":{"id":106734,"name":"Someone","link":"https://s152-en.ogame.gameforge.com/game/index.php?page=ingame&component=galaxy&galaxy=4&system=116"},"sum":2000,"price":3000,"bids":2,"auctionId":"13355"}]}
			// 5::/auctioneer:{"name":"auction finished","args":[{"sum":2000,"player":{"id":106734,"name":"Someone","link":"http://s152-en.ogame.gameforge.com/game/index.php?page=ingame&component=galaxy&galaxy=4&system=116"},"bids":2,"info":"Next auction in:<br />\n<span class=\"nextAuction\" id=\"nextAuction\">1390</span>","time":"06:36"}]}
			msg = bytes.TrimPrefix(msg, []byte("5::/auctioneer:"))
			var pck any = string(msg)
			var out map[string]any
			_ = json.Unmarshal(msg, &out)
			if args, ok := out["args"].([]any); ok {
				if len(args) > 0 {
					if name, ok := out["name"].(string); ok && name == "new bid" {
						if firstArg, ok := args[0].(map[string]any); ok {
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
					} else if name, ok := out["name"].(string); ok && name == "timeLeft" {
						if timeLeftMsg, ok := args[0].(string); ok {
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
					} else if name, ok := out["name"].(string); ok && name == "new auction" {
						if firstArg, ok := args[0].(map[string]any); ok {
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
					} else if name, ok := out["name"].(string); ok && name == "auction finished" {
						if firstArg, ok := args[0].(map[string]any); ok {
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
			}
			for _, clb := range b.auctioneerCallbacks {
				clb(pck)
			}
		} else if regexp.MustCompile(`6::/chat:\d+\+\[true]`).Match(msg) {
			b.debug("chat connected")
		} else if regexp.MustCompile(`6::/chat:\d+\+\[false]`).Match(msg) {
			b.error("Failed to connect to chat")
		} else if bytes.HasPrefix(msg, []byte("5::/chat:")) {
			payload := bytes.TrimPrefix(msg, []byte("5::/chat:"))
			var chatPayload ogame.ChatPayload
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
	_ = websocket.Message.Send(b.ws, "1::/chat")
	return true
}

func (b *OGame) logout() {
	_, _ = b.getPage(LogoutPageName)
	_ = b.client.Jar.(*cookiejar.Jar).Save()
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

// IsKnowFullPage ...
func IsKnowFullPage(vals url.Values) bool {
	page := getPageName(vals)
	return page == OverviewPageName ||
		page == TraderOverviewPageName ||
		page == ResearchPageName ||
		page == ShipyardPageName ||
		page == GalaxyPageName ||
		page == AlliancePageName ||
		page == PremiumPageName ||
		page == ShopPageName ||
		page == RewardsPageName ||
		page == ResourceSettingsPageName ||
		page == MovementPageName ||
		page == HighscorePageName ||
		page == BuddiesPageName ||
		page == PreferencesPageName ||
		page == MessagesPageName ||
		page == ChatPageName ||
		page == DefensesPageName ||
		page == SuppliesPageName ||
		page == FacilitiesPageName ||
		page == FleetdispatchPageName
}

func IsEmpirePage(vals url.Values) bool {
	return vals.Get("page") == "standalone" && vals.Get("component") == "empire"
}

// IsAjaxPage either the requested page is a partial/ajax page
func IsAjaxPage(vals url.Values) bool {
	page := getPageName(vals)
	ajax := vals.Get("ajax")
	asJson := vals.Get("asJson")
	return page == FetchEventboxAjaxPageName ||
		page == FetchResourcesAjaxPageName ||
		page == GalaxyContentAjaxPageName ||
		page == EventListAjaxPageName ||
		page == AjaxChatAjaxPageName ||
		page == NoticesAjaxPageName ||
		page == RepairlayerAjaxPageName ||
		page == TechtreeAjaxPageName ||
		page == PhalanxAjaxPageName ||
		page == ShareReportOverlayAjaxPageName ||
		page == JumpgatelayerAjaxPageName ||
		page == FederationlayerAjaxPageName ||
		page == UnionchangeAjaxPageName ||
		page == ChangenickAjaxPageName ||
		page == PlanetlayerAjaxPageName ||
		page == TraderlayerAjaxPageName ||
		page == PlanetRenameAjaxPageName ||
		page == RightmenuAjaxPageName ||
		page == AllianceOverviewAjaxPageName ||
		page == SupportAjaxPageName ||
		page == BuffActivationAjaxPageName ||
		page == AuctioneerAjaxPageName ||
		page == HighscoreContentAjaxPageName ||
		ajax == "1" ||
		asJson == "1"
}

func canParseEventBox(by []byte) bool {
	err := json.Unmarshal(by, &eventboxResp{})
	return err == nil
}

func canParseSystemInfos(by []byte) bool {
	err := json.Unmarshal(by, &ogame.SystemInfos{})
	return err == nil
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
	resp, err := b.client.Do(req)
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
		(page == "componentOnly" && component == FetchEventboxAjaxPageName) ||
		(page == "componentOnly" && component == EventListAjaxPageName && vals.Get("action") != "fetchEventBox") {
		page = component
	}
	return page
}

func getOptions(opts ...Option) (out Options) {
	for _, opt := range opts {
		opt(&out)
	}
	return
}

func setCPParam(b *OGame, vals url.Values, cfg Options) {
	if vals.Get("cp") == "" &&
		cfg.ChangePlanet != 0 &&
		b.getCachedCelestial(cfg.ChangePlanet) != nil {
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
		return page == GalaxyContentAjaxPageName && !canParseSystemInfos(pageHTML)
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
	retryPolicy := b.withRetry
	if cfg.SkipRetry {
		retryPolicy = b.withoutRetry
	}
	return retryPolicy
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
			b.client.CheckRedirect = func(req *http.Request, via []*http.Request) error { return http.ErrUseLastResponse }
			defer func() { b.client.CheckRedirect = nil }()
		}

		pageHTMLBytes, err = b.execRequest(method, finalURL, payload, vals)
		if err != nil {
			return err
		}

		if detectLoggedOut(method, page, vals, pageHTMLBytes) {
			b.error("Err not logged on page : ", page)
			atomic.StoreInt32(&b.isConnectedAtom, 0)
			return ogame.ErrNotLogged
		}

		return nil
	}

	retryPolicy := retryPolicyFromConfig(b, cfg)
	if err := retryPolicy(clb); err != nil {
		b.error(err)
		return []byte{}, err
	}

	if err := processResponseHTML(method, b, pageHTMLBytes, page, payload, vals); err != nil {
		return []byte{}, err
	}

	if !cfg.SkipInterceptor {
		go func() {
			for _, fn := range b.interceptorCallbacks {
				fn(method, finalURL, vals, payload, pageHTMLBytes)
			}
		}()
	}

	return pageHTMLBytes, nil
}

func alterPayload(method string, b *OGame, vals, payload url.Values) {
	switch method {
	case http.MethodPost:
		if vals.Get("page") == "ajaxChat" && payload.Get("mode") == "1" {
			payload.Set("token", b.ajaxChatToken)
		}
	}
}

func processResponseHTML(method string, b *OGame, pageHTML []byte, page string, payload, vals url.Values) error {
	switch method {
	case http.MethodGet:
		if !IsAjaxPage(vals) && !IsEmpirePage(vals) && v6.IsLogged(pageHTML) {
			parsedFullPage := parser.AutoParseFullPage(b.extractor, pageHTML)
			b.cacheFullPageInfo(parsedFullPage)
		}

	case http.MethodPost:
		if page == PreferencesPageName {
			b.CachedPreferences = b.extractor.ExtractPreferences(pageHTML)
		} else if page == "ajaxChat" && (payload.Get("mode") == "1" || payload.Get("mode") == "3") {
			if err := extractNewChatToken(b, pageHTML); err != nil {
				return err
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
			return errors.Wrap(err, ogame.ErrFailedExecuteCallback.Error())
		}

		if retryErr := retry(err); retryErr != nil {
			return retryErr
		}

		if err == ogame.ErrNotLogged {
			if _, loginErr := b.wrapLoginWithExistingCookies(); loginErr != nil {
				b.error(loginErr.Error()) // log error
				if loginErr == ogame.ErrAccountNotFound ||
					loginErr == ogame.ErrAccountBlocked ||
					loginErr == ogame.ErrBadCredentials ||
					loginErr == ogame.ErrOTPRequired ||
					loginErr == ogame.ErrOTPInvalid {
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
	return obj.ConstructionTime(nbr, b.getUniverseSpeed(), facilities, b.hasTechnocrat, b.isDiscoverer())
}

func (b *OGame) enable() {
	b.ctx, b.cancelCtx = context.WithCancel(context.Background())
	atomic.StoreInt32(&b.isEnabledAtom, 1)
	b.stateChanged(false, "Enable")
}

func (b *OGame) disable() {
	atomic.StoreInt32(&b.isEnabledAtom, 0)
	b.cancelCtx()
	b.stateChanged(false, "Disable")
}

func (b *OGame) isEnabled() bool {
	return atomic.LoadInt32(&b.isEnabledAtom) == 1
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

func (b *OGame) isUnderAttack() (bool, error) {
	res, err := b.fetchEventbox()
	return res.Hostile > 0, err
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

func (b *OGame) getPlanets() []Planet {
	page, _ := getPage[parser.OverviewPage](b)
	return convertPlanets(b, page.ExtractPlanets())
}

func (b *OGame) getPlanet(v any) (Planet, error) {
	page, _ := getPage[parser.OverviewPage](b)
	planet, err := page.ExtractPlanet(v)
	if err != nil {
		return Planet{}, err
	}
	return convertPlanet(b, planet), nil
}

func (b *OGame) getMoons() []Moon {
	page, _ := getPage[parser.OverviewPage](b)
	return convertMoons(b, page.ExtractMoons())
}

func (b *OGame) getMoon(v any) (Moon, error) {
	page, _ := getPage[parser.OverviewPage](b)
	moon, err := page.ExtractMoon(v)
	if err != nil {
		return Moon{}, err
	}
	cMoon := convertMoon(b, moon)
	return *cMoon, nil
}

func (b *OGame) getCelestials() ([]Celestial, error) {
	page, _ := getPage[parser.OverviewPage](b)
	celestials, err := page.ExtractCelestials()
	if err != nil {
		return nil, err
	}
	return convertCelestials(b, celestials), nil
}

func (b *OGame) getCelestial(v any) (Celestial, error) {
	page, _ := getPage[parser.OverviewPage](b)
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

func (b *OGame) abandon(v any) error {
	page, _ := getPage[parser.OverviewPage](b)
	planet, err := page.ExtractPlanet(v)
	if err != nil {
		return errors.New("invalid parameter")
	}
	pageHTML, _ := b.getPage(PlanetlayerPageName, ChangePlanet(planet.GetID()))
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	abandonToken := doc.Find("form#planetMaintenanceDelete input[name=abandon]").AttrOr("value", "")
	token := doc.Find("form#planetMaintenanceDelete input[name=token]").AttrOr("value", "")
	payload := url.Values{
		"abandon":  {abandonToken},
		"token":    {token},
		"password": {b.password},
	}
	_, err = b.postPageContent(url.Values{
		"page":      {"ingame"},
		"component": {"overview"},
		"action":    {"planetGiveup"},
		"ajax":      {"1"},
		"asJson":    {"1"},
	}, payload)
	return err
}

func (b *OGame) serverTime() time.Time {
	page, err := getPage[parser.OverviewPage](b)
	serverTime, err := page.ExtractServerTime()
	if err != nil {
		b.error(err.Error())
	}
	return serverTime
}

func (b *OGame) getUserInfos() ogame.UserInfos {
	page, err := getPage[parser.OverviewPage](b)
	userInfos, err := page.ExtractUserInfos()
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
	page, _ := getPage[parser.MovementPage](b, opts...)
	fleets := page.ExtractFleets()
	slots := page.ExtractSlots()
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

func (b *OGame) getSlots() ogame.Slots {
	pageHTML, _ := b.getPage(FleetdispatchPageName)
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

func findSlowestSpeed(ships ogame.ShipsInfos, techs ogame.Researches, isCollector, isGeneral bool) int64 {
	var minSpeed int64 = math.MaxInt64
	for _, ship := range ogame.Ships {
		if ship.GetID() == ogame.SolarSatelliteID || ship.GetID() == ogame.CrawlerID {
			continue
		}
		shipSpeed := ship.GetSpeed(techs, isCollector, isGeneral)
		if ships.ByID(ship.GetID()) > 0 && shipSpeed < minSpeed {
			minSpeed = shipSpeed
		}
	}
	return minSpeed
}

func calcFuel(ships ogame.ShipsInfos, dist, duration int64, universeSpeedFleet, fleetDeutSaveFactor float64, techs ogame.Researches, isCollector, isGeneral bool) (fuel int64) {
	tmpFn := func(baseFuel, nbr, shipSpeed int64) float64 {
		tmpSpeed := (35000 / (float64(duration)*universeSpeedFleet - 10)) * math.Sqrt(float64(dist)*10/float64(shipSpeed))
		return float64(baseFuel*nbr*dist) / 35000 * math.Pow(tmpSpeed/10+1, 2)
	}
	tmpFuel := 0.0
	for _, ship := range ogame.Ships {
		if ship.GetID() == ogame.SolarSatelliteID || ship.GetID() == ogame.CrawlerID {
			continue
		}
		nbr := ships.ByID(ship.GetID())
		if nbr > 0 {
			tmpFuel += tmpFn(ship.GetFuelConsumption(techs, fleetDeutSaveFactor, isGeneral), nbr, ship.GetSpeed(techs, isCollector, isGeneral))
		}
	}
	fuel = int64(1 + math.Round(tmpFuel))
	return
}

// CalcFlightTime ...
func CalcFlightTime(origin, destination ogame.Coordinate, universeSize, nbSystems int64, donutGalaxy, donutSystem bool,
	fleetDeutSaveFactor, speed float64, universeSpeedFleet int64, ships ogame.ShipsInfos, techs ogame.Researches, characterClass ogame.CharacterClass) (secs, fuel int64) {
	if !ships.HasShips() {
		return
	}
	isCollector := characterClass == ogame.Collector
	isGeneral := characterClass == ogame.General
	s := speed
	v := float64(findSlowestSpeed(ships, techs, isCollector, isGeneral))
	a := float64(universeSpeedFleet)
	d := float64(Distance(origin, destination, universeSize, nbSystems, donutGalaxy, donutSystem))
	secs = int64(math.Round(((3500/s)*math.Sqrt(d*10/v) + 10) / a))
	fuel = calcFuel(ships, int64(d), secs, float64(universeSpeedFleet), fleetDeutSaveFactor, techs, isCollector, isGeneral)
	return
}

// CalcFlightTime calculates the flight time and the fuel consumption
func (b *OGame) CalcFlightTime(origin, destination ogame.Coordinate, speed float64, ships ogame.ShipsInfos, missionID ogame.MissionID) (secs, fuel int64) {
	return CalcFlightTime(origin, destination, b.serverData.Galaxies, b.serverData.Systems, b.serverData.DonutGalaxy,
		b.serverData.DonutSystem, b.serverData.GlobalDeuteriumSaveFactor, speed, GetFleetSpeedForMission(b.serverData, missionID), ships,
		b.GetCachedResearch(), b.characterClass)
}

// getPhalanx makes 3 calls to ogame server (2 validation, 1 scan)
func (b *OGame) getPhalanx(moonID ogame.MoonID, coord ogame.Coordinate) ([]ogame.Fleet, error) {
	res := make([]ogame.Fleet, 0)

	// Get moon facilities html page (first call to ogame server)
	moonFacilitiesHTML, _ := b.getPage(FacilitiesPageName, ChangePlanet(moonID.Celestial()))

	// Extract bunch of infos from the html
	moon, err := b.extractor.ExtractMoon(moonFacilitiesHTML, moonID)
	if err != nil {
		return res, errors.New("moon not found")
	}
	resources := b.extractor.ExtractResources(moonFacilitiesHTML)
	moonFacilities, _ := b.extractor.ExtractFacilities(moonFacilitiesHTML)
	phalanxLvl := moonFacilities.SensorPhalanx

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

	// Run the phalanx scan (second & third calls to ogame server)
	return b.getUnsafePhalanx(moonID, coord)
}

// getUnsafePhalanx ...
func (b *OGame) getUnsafePhalanx(moonID ogame.MoonID, coord ogame.Coordinate) ([]ogame.Fleet, error) {
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

	vals := url.Values{
		"page":     {PhalanxAjaxPageName},
		"galaxy":   {utils.FI64(coord.Galaxy)},
		"system":   {utils.FI64(coord.System)},
		"position": {utils.FI64(coord.Position)},
		"ajax":     {"1"},
		"token":    {planetInfos.OverlayToken},
	}
	page, _ := getAjaxPage[parser.PhalanxAjaxPage](b, vals, ChangePlanet(moonID.Celestial()))
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
	resp, err := b.client.Do(req)
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

func (b *OGame) jumpGateDestinations(originMoonID ogame.MoonID) ([]ogame.MoonID, int64, error) {
	pageHTML, _ := b.getPage(JumpgatelayerPageName, ChangePlanet(originMoonID.Celestial()))
	_, _, dests, wait := b.extractor.ExtractJumpGate(pageHTML)
	if wait > 0 {
		return dests, wait, fmt.Errorf("jump gate is in recharge mode for %d seconds", wait)
	}
	return dests, wait, nil
}

func (b *OGame) executeJumpGate(originMoonID, destMoonID ogame.MoonID, ships ogame.ShipsInfos) (bool, int64, error) {
	pageHTML, _ := b.getPage(JumpgatelayerPageName, ChangePlanet(originMoonID.Celestial()))
	availShips, token, dests, wait := b.extractor.ExtractJumpGate(pageHTML)
	if wait > 0 {
		return false, wait, fmt.Errorf("jump gate is in recharge mode for %d seconds", wait)
	}

	// Validate destination moon id
	if !moonIDInSlice(destMoonID, dests) {
		return false, 0, errors.New("destination moon id invalid")
	}

	payload := url.Values{"token": {token}, "zm": {utils.FI64(destMoonID)}}

	// Add ships to payload
	for _, s := range ogame.Ships {
		// Get the min between what is available and what we want
		nbr := int64(math.Min(float64(ships.ByID(s.GetID())), float64(availShips.ByID(s.GetID()))))
		if nbr > 0 {
			payload.Add("ship_"+utils.FI64(s.GetID()), utils.FI64(nbr))
		}
	}

	if _, err := b.postPageContent(url.Values{"page": {"jumpgate_execute"}}, payload); err != nil {
		return false, 0, err
	}
	return true, 0, nil
}

func (b *OGame) getEmpire(celestialType ogame.CelestialType) (out []ogame.EmpireCelestial, err error) {
	var planetType int
	if celestialType == ogame.PlanetType {
		planetType = 0
	} else if celestialType == ogame.MoonType {
		planetType = 1
	} else {
		return out, errors.New("invalid celestial type")
	}
	vals := url.Values{"page": {"standalone"}, "component": {"empire"}, "planetType": {strconv.Itoa(planetType)}}
	pageHTMLBytes, err := b.getPageContent(vals)
	if err != nil {
		return out, err
	}
	return b.extractor.ExtractEmpire(pageHTMLBytes)
}

func (b *OGame) getEmpireJSON(nbr int64) (any, error) {
	// Valid URLs:
	// /game/index.php?page=standalone&component=empire&planetType=0
	// /game/index.php?page=standalone&component=empire&planetType=1
	vals := url.Values{"page": {"standalone"}, "component": {"empire"}, "planetType": {utils.FI64(nbr)}}
	pageHTMLBytes, err := b.getPageContent(vals)
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
	page, _ := getPage[parser.OverviewPage](b, ChangePlanet(celestialID))
	return page.ExtractDMCosts()
}

func (b *OGame) useDM(typ string, celestialID ogame.CelestialID) error {
	if typ != "buildings" && typ != "research" && typ != "shipyard" {
		return fmt.Errorf("invalid type %s", typ)
	}
	page, _ := getPage[parser.OverviewPage](b, ChangePlanet(celestialID))
	costs, err := page.ExtractDMCosts()
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
	getToken := func(pageHTML []byte) (string, error) {
		m := regexp.MustCompile(`var token = "([^"]+)"`).FindSubmatch(pageHTML)
		if len(m) != 2 {
			return "", errors.New("unable to find token")
		}
		return string(m[1]), nil
	}
	token, _ := getToken(pageHTML)

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
	page, _ := getPage[parser.OverviewPage](b, ChangePlanet(celestialID))
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

		payload.Add("bid[planets]["+utils.FI64(celestialID)+"][metal]", utils.FI64(metalNeeded))
		payload.Add("bid[planets]["+utils.FI64(celestialID)+"][crystal]", utils.FI64(crystalNeeded))
		payload.Add("bid[planets]["+utils.FI64(celestialID)+"][deuterium]", utils.FI64(deuteriumNeeded))
	}
	return payload
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
	payload := calcResources(price, planetResources, multiplier)
	payload.Add("action", "trade")
	payload.Add("bid[honor]", "0")
	payload.Add("token", importToken)
	payload.Add("ajax", "1")
	pageHTML1, err := b.postPageContent(url.Values{"page": {"ajax"}, "component": {"traderimportexport"}, "ajax": {"1"}, "action": {"trade"}, "asJson": {"1"}}, payload)
	if err != nil {
		return err
	}
	// {"message":"You have bought a container.","error":false,"item":{"uuid":"40f6c78e11be01ad3389b7dccd6ab8efa9347f3c","itemText":"You have purchased 1 KRAKEN Bronze.","bargainText":"The contents of the container not appeal to you? For 500 Dark Matter you can exchange the container for another random container of the same quality. You can only carry out this exchange 2 times per daily offer.","bargainCost":500,"bargainCostText":"Costs: 500 Dark Matter","tooltip":"KRAKEN Bronze|Reduces the building time of buildings currently under construction by <b>30m<\/b>.<br \/><br \/>\nDuration: now<br \/><br \/>\nPrice: --- <br \/>\nIn Inventory: 1","image":"98629d11293c9f2703592ed0314d99f320f45845","amount":1,"rarity":"common"},"newToken":"07eefc14105db0f30cb331a8b7af0bfe"}
	var tmp struct {
		Message      string
		Error        bool
		NewAjaxToken string
	}
	if err := json.Unmarshal(pageHTML1, &tmp); err != nil {
		return err
	}
	if tmp.Error {
		return errors.New(tmp.Message)
	}

	payload2 := url.Values{"action": {"takeItem"}, "token": {tmp.NewAjaxToken}, "ajax": {"1"}}
	pageHTML2, _ := b.postPageContent(url.Values{"page": {"ajax"}, "component": {"traderimportexport"}, "ajax": {"1"}, "action": {"takeItem"}, "asJson": {"1"}}, payload2)
	var tmp2 struct {
		Message      string
		Error        bool
		NewAjaxToken string
	}
	if err := json.Unmarshal(pageHTML2, &tmp2); err != nil {
		return err
	}
	if tmp2.Error {
		return errors.New(tmp2.Message)
	}
	// {"error":false,"message":"You have accepted the offer and put the item in your inventory.","item":{"name":"Bronze Deuterium Booster","image":"f0e514af79d0808e334e9b6b695bf864b861bdfa","imageLarge":"c7c2837a0b341d37383d6a9d8f8986f500db7bf9","title":"Bronze Deuterium Booster|+10% more Deuterium Synthesizer harvest on one planet<br \/><br \/>\nDuration: 1w<br \/><br \/>\nPrice: --- <br \/>\nIn Inventory: 134","effect":"+10% more Deuterium Synthesizer harvest on one planet","ref":"d9fa5f359e80ff4f4c97545d07c66dbadab1d1be","rarity":"common","amount":134,"amount_free":134,"amount_bought":0,"category":["d8d49c315fa620d9c7f1f19963970dea59a0e3be","e71139e15ee5b6f472e2c68a97aa4bae9c80e9da"],"currency":"dm","costs":"2500","isReduced":false,"buyable":false,"canBeActivated":true,"canBeBoughtAndActivated":false,"isAnUpgrade":false,"isCharacterClassItem":false,"hasEnoughCurrency":true,"cooldown":0,"duration":604800,"durationExtension":null,"totalTime":null,"timeLeft":null,"status":null,"extendable":false,"firstStatus":"effecting","toolTip":"Bronze Deuterium Booster|+10% more Deuterium Synthesizer harvest on one planet&lt;br \/&gt;&lt;br \/&gt;\nDuration: 1w&lt;br \/&gt;&lt;br \/&gt;\nPrice: --- &lt;br \/&gt;\nIn Inventory: 134","buyTitle":"This item is currently unavailable for purchase.","activationTitle":"Activate","moonOnlyItem":false,"newOffer":false,"noOfferMessage":"There are no further offers today. Please come again tomorrow."},"newToken":"dec779714b893be9b39c0bedf5738450","components":[],"newAjaxToken":"e20cf0a6ca0e9b43a81ccb8fe7e7e2e3"}

	return nil
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
	ownCoords := make([]ogame.Coordinate, 0)
	planets := b.GetCachedPlanets()
	for _, planet := range planets {
		ownCoords = append(ownCoords, planet.Coordinate)
		if planet.Moon != nil {
			ownCoords = append(ownCoords, planet.Moon.Coordinate)
		}
	}
	out, err = page.ExtractAttacks(ownCoords)
	if err != nil {
		return
	}
	fixAttackEvents(out, planets)
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
	vals := url.Values{"page": {"ingame"}, "component": {"galaxyContent"}, "ajax": {"1"}}
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
	if res.Tmpgalaxy != galaxy || res.Tmpsystem != system {
		return ogame.SystemInfos{}, errors.New("not enough deuterium")
	}
	return res, err
}

func (b *OGame) getResourceSettings(planetID ogame.PlanetID, options ...Option) (ogame.ResourceSettings, error) {
	options = append(options, ChangePlanet(planetID.Celestial()))
	page, _ := getPage[parser.ResourcesSettingsPage](b, options...)
	return page.ExtractResourceSettings()
}

func (b *OGame) setResourceSettings(planetID ogame.PlanetID, settings ogame.ResourceSettings) error {
	pageHTML, _ := b.getPage(ResourceSettingsPageName, ChangePlanet(planetID.Celestial()))
	doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	bodyID := b.extractor.ExtractBodyIDFromDoc(doc)
	if bodyID == "overview" {
		return ogame.ErrInvalidPlanetID
	}
	token, exists := doc.Find("form input[name=token]").Attr("value")
	if !exists {
		return errors.New("unable to find token")
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
	resp, err := b.client.PostForm(url2, payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func (b *OGame) getCachedResearch() ogame.Researches {
	if b.researches == nil {
		return b.getResearch()
	}
	return *b.researches
}

func (b *OGame) getResearch() ogame.Researches {
	page, _ := getPage[parser.ResearchPage](b)
	researches := page.ExtractResearch()
	b.researches = &researches
	return researches
}

func (b *OGame) getResourcesBuildings(celestialID ogame.CelestialID, options ...Option) (ogame.ResourcesBuildings, error) {
	options = append(options, ChangePlanet(celestialID))
	page, _ := getPage[parser.SuppliesPage](b, options...)
	return page.ExtractResourcesBuildings()
}

func (b *OGame) getDefense(celestialID ogame.CelestialID, options ...Option) (ogame.DefensesInfos, error) {
	options = append(options, ChangePlanet(celestialID))
	page, _ := getPage[parser.DefensesPage](b, options...)
	return page.ExtractDefense()
}

func (b *OGame) getShips(celestialID ogame.CelestialID, options ...Option) (ogame.ShipsInfos, error) {
	options = append(options, ChangePlanet(celestialID))
	page, _ := getPage[parser.ShipyardPage](b, options...)
	return page.ExtractShips()
}

func (b *OGame) getFacilities(celestialID ogame.CelestialID, options ...Option) (ogame.Facilities, error) {
	options = append(options, ChangePlanet(celestialID))
	page, _ := getPage[parser.FacilitiesPage](b, options...)
	return page.ExtractFacilities()
}

func (b *OGame) getTechs(celestialID ogame.CelestialID) (ogame.ResourcesBuildings, ogame.Facilities, ogame.ShipsInfos, ogame.DefensesInfos, ogame.Researches, error) {
	vals := url.Values{"page": {FetchTechsName}}
	page, _ := getAjaxPage[parser.FetchTechsAjaxPage](b, vals, ChangePlanet(celestialID))
	return page.ExtractTechs()
}

func (b *OGame) getProduction(celestialID ogame.CelestialID) ([]ogame.Quantifiable, int64, error) {
	page, _ := getPage[parser.ShipyardPage](b, ChangePlanet(celestialID))
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

func getToken(b *OGame, page string, celestialID ogame.CelestialID) (string, error) {
	pageHTML, _ := b.getPage(page, ChangePlanet(celestialID))
	rgx := regexp.MustCompile(`var upgradeEndpoint = ".+&token=([^&]+)&`)
	m := rgx.FindSubmatch(pageHTML)
	if len(m) != 2 {
		return "", errors.New("unable to find form token")
	}
	return string(m[1]), nil
}

func getDemolishToken(b *OGame, page string, celestialID ogame.CelestialID) (string, error) {
	pageHTML, _ := b.getPage(page, ChangePlanet(celestialID))
	m := regexp.MustCompile(`modus=3&token=([^&]+)&`).FindSubmatch(pageHTML)
	if len(m) != 2 {
		return "", errors.New("unable to find form token")
	}
	return string(m[1]), nil
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

	token, err := getDemolishToken(b, page, celestialID)
	if err != nil {
		return err
	}

	pageHTML, _ := b.getPageContent(url.Values{
		"page":       {"ingame"},
		"component":  {"technologydetails"},
		"ajax":       {"1"},
		"action":     {"getDetails"},
		"technology": {utils.FI64(id)},
		"cp":         {utils.FI64(celestialID)},
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
		"type":      {utils.FI64(id)},
		"cp":        {utils.FI64(celestialID)},
	}
	_, err = b.getPageContent(params)
	return err
}

func (b *OGame) build(celestialID ogame.CelestialID, id ogame.ID, nbr int64) error {
	var page string
	if id.IsDefense() {
		page = DefensesPageName
	} else if id.IsShip() {
		page = ShipyardPageName
	} else if id.IsLfBuilding() {
		page = LfbuildingsPageName
	} else if id.IsBuilding() {
		page = SuppliesPageName
	} else if id.IsTech() {
		page = ResearchPageName
	} else {
		return errors.New("invalid id " + id.String())
	}
	vals := url.Values{
		"page":      {"ingame"},
		"component": {page},
		"modus":     {"1"},
		"type":      {utils.FI64(id)},
		"cp":        {utils.FI64(celestialID)},
	}

	token, err := getToken(b, page, celestialID)
	if err != nil {
		return err
	}
	vals.Add("token", token)

	if id.IsDefense() || id.IsShip() {
		var maximumNbr int64 = 99999
		var err error
		var token string
		for nbr > 0 {
			tmp := int64(math.Min(float64(nbr), float64(maximumNbr)))
			vals.Set("menge", utils.FI64(tmp))
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

	_, err = b.getPageContent(vals)
	return err
}

func (b *OGame) buildCancelable(celestialID ogame.CelestialID, id ogame.ID) error {
	if !id.IsBuilding() && !id.IsTech() {
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
	if !technologyID.IsTech() {
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

func (b *OGame) constructionsBeingBuilt(celestialID ogame.CelestialID) (ogame.ID, int64, ogame.ID, int64) {
	page, _ := getPage[parser.OverviewPage](b, ChangePlanet(celestialID))
	return page.ExtractConstructions()
}

func (b *OGame) cancel(token string, techID, listID int64) error {
	_, _ = b.getPageContent(url.Values{"page": {"ingame"}, "component": {"overview"}, "modus": {"2"}, "token": {token},
		"type": {utils.FI64(techID)}, "listid": {utils.FI64(listID)}, "action": {"cancel"}})
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
		// components??
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

	duration, max, token := page.ExtractIPM()
	if max == 0 {
		return 0, errors.New("no missile available")
	}
	if nbr > max {
		nbr = max
	}
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
		"missilePrimaryTarget": {},
	}
	if priority != 0 {
		payload.Add("missilePrimaryTarget", utils.FI64(priority))
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

func (b *OGame) sendFleet(celestialID ogame.CelestialID, ships []ogame.Quantifiable, speed ogame.Speed, where ogame.Coordinate,
	mission ogame.MissionID, resources ogame.Resources, holdingTime, unionID int64, ensure bool) (ogame.Fleet, error) {

	// Get existing fleet, so we can ensure new fleet ID is greater
	initialFleets, slots := b.getFleets()
	maxInitialFleetID := ogame.FleetID(0)
	for _, f := range initialFleets {
		if f.ID > maxInitialFleetID {
			maxInitialFleetID = f.ID
		}
	}

	if slots.InUse == slots.Total {
		return ogame.Fleet{}, ogame.ErrAllSlotsInUse
	}

	if mission == ogame.Expedition {
		if slots.ExpInUse == slots.ExpTotal {
			return ogame.Fleet{}, ogame.ErrAllSlotsInUse
		}
	}

	// Page 1 : get to fleet page
	pageHTML, err := b.getPage(FleetdispatchPageName, ChangePlanet(celestialID))
	if err != nil {
		return ogame.Fleet{}, err
	}

	fleet1Doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(pageHTML))
	fleet1BodyID := b.extractor.ExtractBodyIDFromDoc(fleet1Doc)
	if fleet1BodyID != FleetdispatchPageName {
		now := time.Now().Unix()
		b.error(ogame.ErrInvalidPlanetID.Error()+", planetID:", celestialID, ", ts: ", now)
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
				return ogame.Fleet{}, fmt.Errorf("not enough ships to send, %s", ogame.Objs.ByID(ship.ID).GetName())
			}
			atLeastOneShipSelected = true
		}
	}
	if !atLeastOneShipSelected {
		return ogame.Fleet{}, ogame.ErrNoShipSelected
	}

	payload := b.extractor.ExtractHiddenFieldsFromDoc(fleet1Doc)
	for _, s := range ships {
		if s.ID.IsFlyableShip() && s.Nbr > 0 {
			payload.Set("am"+utils.FI64(s.ID), utils.FI64(s.Nbr))
		}
	}

	tokenM := regexp.MustCompile(`var fleetSendingToken = "([^"]+)";`).FindSubmatch(pageHTML)
	if b.IsV8() || b.IsV9() {
		tokenM = regexp.MustCompile(`var token = "([^"]+)";`).FindSubmatch(pageHTML)
	}
	if len(tokenM) != 2 {
		return ogame.Fleet{}, errors.New("token not found")
	}

	payload.Set("token", string(tokenM[1]))
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
		fleet1Doc.Find("select[name=acsValues] option").Each(func(i int, s *goquery.Selection) {
			acsValues := s.AttrOr("value", "")
			m := regexp.MustCompile(`\d+#\d+#\d+#\d+#.*#(\d+)`).FindStringSubmatch(acsValues)
			if len(m) == 2 {
				optUnionID := utils.DoParseI64(m[1])
				if unionID == optUnionID {
					found = true
					payload.Add("acsValues", acsValues)
					payload.Add("union", m[1])
					mission = ogame.GroupedAttack
				}
			}
		})
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

	cargo := ogame.ShipsInfos{}.FromQuantifiables(ships).Cargo(b.getCachedResearch(), b.server.Settings.EspionageProbeRaids == 1, b.isCollector(), b.IsPioneers())
	newResources := ogame.Resources{}
	if resources.Total() > cargo {
		newResources.Deuterium = int64(math.Min(float64(resources.Deuterium), float64(cargo)))
		cargo -= newResources.Deuterium
		newResources.Crystal = int64(math.Min(float64(resources.Crystal), float64(cargo)))
		cargo -= newResources.Crystal
		newResources.Metal = int64(math.Min(float64(resources.Metal), float64(cargo)))
	} else {
		newResources = resources
	}

	newResources.Metal = utils.MaxInt(newResources.Metal, 0)
	newResources.Crystal = utils.MaxInt(newResources.Crystal, 0)
	newResources.Deuterium = utils.MaxInt(newResources.Deuterium, 0)

	// Page 3 : select coord, mission, speed
	if b.IsV8() || b.IsV9() {
		payload.Set("token", checkRes.NewAjaxToken)
	}
	payload.Set("speed", strconv.FormatInt(int64(speed), 10))
	payload.Set("crystal", utils.FI64(newResources.Crystal))
	payload.Set("deuterium", utils.FI64(newResources.Deuterium))
	payload.Set("metal", utils.FI64(newResources.Metal))
	payload.Set("mission", utils.FI64(mission))
	payload.Set("prioMetal", "1")
	payload.Set("prioCrystal", "2")
	payload.Set("prioDeuterium", "3")
	payload.Set("retreatAfterDefenderRetreat", "0")
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
	movementHTML, _ := b.getPage(MovementPageName)
	movementDoc, _ := goquery.NewDocumentFromReader(bytes.NewReader(movementHTML))
	originCoords, _ := b.extractor.ExtractPlanetCoordinate(movementHTML)
	fleets := b.extractor.ExtractFleetsFromDoc(movementDoc)
	if len(fleets) > 0 {
		max := ogame.Fleet{}
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
		return ogame.Fleet{}, ogame.ErrAllSlotsInUse
	}

	if mission == ogame.Expedition {
		if slots.ExpInUse == slots.ExpTotal {
			return ogame.Fleet{}, ogame.ErrAllSlotsInUse
		}
	}

	now := time.Now().Unix()
	b.error(errors.New("could not find new fleet ID").Error()+", planetID:", celestialID, ", ts: ", now)
	return ogame.Fleet{}, errors.New("could not find new fleet ID")
}

func (b *OGame) getPageMessages(page int64, tabid ogame.MessagesTabID) ([]byte, error) {
	payload := url.Values{
		"messageId":  {"-1"},
		"tabid":      {utils.FI64(tabid)},
		"action":     {"107"},
		"pagination": {utils.FI64(page)},
		"ajax":       {"1"},
	}
	return b.postPageContent(url.Values{"page": {"messages"}}, payload)
}

func (b *OGame) getEspionageReportMessages() ([]ogame.EspionageReportSummary, error) {
	var page int64 = 1
	var nbPage int64 = 1
	msgs := make([]ogame.EspionageReportSummary, 0)
	for page <= nbPage {
		pageHTML, _ := b.getPageMessages(page, EspionageMessagesTabID)
		newMessages, newNbPage := b.extractor.ExtractEspionageReportMessageIDs(pageHTML)
		msgs = append(msgs, newMessages...)
		nbPage = newNbPage
		page++
	}
	return msgs, nil
}

func (b *OGame) getCombatReportMessages() ([]ogame.CombatReportSummary, error) {
	var page int64 = 1
	var nbPage int64 = 1
	msgs := make([]ogame.CombatReportSummary, 0)
	for page <= nbPage {
		pageHTML, _ := b.getPageMessages(page, CombatReportsMessagesTabID)
		newMessages, newNbPage := b.extractor.ExtractCombatReportMessagesSummary(pageHTML)
		msgs = append(msgs, newMessages...)
		nbPage = newNbPage
		page++
	}
	return msgs, nil
}

func (b *OGame) getExpeditionMessages() ([]ogame.ExpeditionMessage, error) {
	var page int64 = 1
	var nbPage int64 = 1
	msgs := make([]ogame.ExpeditionMessage, 0)
	for page <= nbPage {
		pageHTML, _ := b.getPageMessages(page, ExpeditionsMessagesTabID)
		newMessages, newNbPage, _ := b.extractor.ExtractExpeditionMessages(pageHTML)
		msgs = append(msgs, newMessages...)
		nbPage = newNbPage
		page++
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
	return b.getMarketplaceMessages(MarketplacePurchasesMessagesTabID)
}

func (b *OGame) getMarketplaceSalesMessages() ([]ogame.MarketplaceMessage, error) {
	return b.getMarketplaceMessages(MarketplaceSalesMessagesTabID)
}

// tabID 26: purchases, 27: sales
func (b *OGame) getMarketplaceMessages(tabID ogame.MessagesTabID) ([]ogame.MarketplaceMessage, error) {
	var page int64 = 1
	var nbPage int64 = 1
	msgs := make([]ogame.MarketplaceMessage, 0)
	for page <= nbPage {
		pageHTML, _ := b.getPageMessages(page, tabID)
		newMessages, newNbPage, _ := b.extractor.ExtractMarketplaceMessages(pageHTML)
		msgs = append(msgs, newMessages...)
		nbPage = newNbPage
		page++
	}
	return msgs, nil
}

func (b *OGame) getExpeditionMessageAt(t time.Time) (ogame.ExpeditionMessage, error) {
	var page int64 = 1
	var nbPage int64 = 1
LOOP:
	for page <= nbPage {
		pageHTML, _ := b.getPageMessages(page, ExpeditionsMessagesTabID)
		newMessages, newNbPage, _ := b.extractor.ExtractExpeditionMessages(pageHTML)
		for _, m := range newMessages {
			if m.CreatedAt.Unix() == t.Unix() {
				return m, nil
			}
			if m.CreatedAt.Unix() < t.Unix() {
				break LOOP
			}
		}
		nbPage = newNbPage
		page++
	}
	return ogame.ExpeditionMessage{}, errors.New("expedition message not found for " + t.String())
}

func (b *OGame) getCombatReportFor(coord ogame.Coordinate) (ogame.CombatReportSummary, error) {
	var page int64 = 1
	var nbPage int64 = 1
	for page <= nbPage {
		pageHTML, err := b.getPageMessages(page, CombatReportsMessagesTabID)
		if err != nil {
			return ogame.CombatReportSummary{}, err
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
	return ogame.CombatReportSummary{}, errors.New("combat report not found for " + coord.String())
}

func (b *OGame) getEspionageReport(msgID int64) (ogame.EspionageReport, error) {
	pageHTML, _ := b.getPageContent(url.Values{"page": {"messages"}, "messageId": {utils.FI64(msgID)}, "tabid": {"20"}, "ajax": {"1"}})
	return b.extractor.ExtractEspionageReport(pageHTML)
}

func (b *OGame) getEspionageReportFor(coord ogame.Coordinate) (ogame.EspionageReport, error) {
	var page int64 = 1
	var nbPage int64 = 1
	for page <= nbPage {
		pageHTML, err := b.getPageMessages(page, EspionageMessagesTabID)
		if err != nil {
			return ogame.EspionageReport{}, err
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
	return ogame.EspionageReport{}, errors.New("espionage report not found for " + coord.String())
}

func (b *OGame) getDeleteMessagesToken() (string, error) {
	pageHTML, _ := b.getPageContent(url.Values{"page": {"messages"}, "tab": {"20"}, "ajax": {"1"}})
	tokenM := regexp.MustCompile(`name='token' value='([^']+)'`).FindSubmatch(pageHTML)
	if len(tokenM) != 2 {
		return "", errors.New("token not found")
	}
	return string(tokenM[1]), nil
}

func (b *OGame) deleteMessage(msgID int64) error {
	token, err := b.getDeleteMessagesToken()
	if err != nil {
		return err
	}
	payload := url.Values{
		"messageId": {utils.FI64(msgID)},
		"action":    {"103"},
		"ajax":      {"1"},
		"token":     {token},
	}
	by, err := b.postPageContent(url.Values{"page": {"messages"}}, payload)
	if err != nil {
		return err
	}

	var res map[string]any
	if err := json.Unmarshal(by, &res); err != nil {
		return errors.New("unable to find message id " + utils.FI64(msgID))
	}
	if val, ok := res[utils.FI64(msgID)]; ok {
		if valB, ok := val.(bool); !ok || !valB {
			return errors.New("unable to find message id " + utils.FI64(msgID))
		}
	} else {
		return errors.New("unable to find message id " + utils.FI64(msgID))
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
	_, err = b.postPageContent(url.Values{"page": {"messages"}}, payload)
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
	researches := b.getResearch()
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

func (b *OGame) addAccount(number int, lang string) (*AddAccountRes, error) {
	accountGroup := fmt.Sprintf("%s_%d", lang, number)
	return AddAccount(b.client, b.ctx, b.lobby, accountGroup, b.bearerToken)
}

func (b *OGame) getCachedCelestial(v any) Celestial {
	switch vv := v.(type) {
	case Celestial:
		return vv
	case Planet:
		return vv
	case Moon:
		return vv
	case ogame.CelestialID:
		return b.GetCachedCelestialByID(vv)
	case ogame.PlanetID:
		return b.GetCachedCelestialByID(vv.Celestial())
	case ogame.MoonID:
		return b.GetCachedCelestialByID(vv.Celestial())
	case int:
		return b.GetCachedCelestialByID(ogame.CelestialID(vv))
	case int32:
		return b.GetCachedCelestialByID(ogame.CelestialID(vv))
	case int64:
		return b.GetCachedCelestialByID(ogame.CelestialID(vv))
	case float32:
		return b.GetCachedCelestialByID(ogame.CelestialID(vv))
	case float64:
		return b.GetCachedCelestialByID(ogame.CelestialID(vv))
	case lua.LNumber:
		return b.GetCachedCelestialByID(ogame.CelestialID(vv))
	case ogame.Coordinate:
		return b.GetCachedCelestialByCoord(vv)
	case string:
		coord, err := ogame.ParseCoord(vv)
		if err != nil {
			return nil
		}
		return b.GetCachedCelestialByCoord(coord)
	}
	return nil
}

// GetCachedCelestialByID return celestial from cached value
func (b *OGame) GetCachedCelestialByID(celestialID ogame.CelestialID) Celestial {
	for _, p := range b.GetCachedPlanets() {
		if p.ID.Celestial() == celestialID {
			return p
		}
		if p.Moon != nil && p.Moon.ID.Celestial() == celestialID {
			return *p.Moon
		}
	}
	return nil
}

// GetCachedCelestialByCoord return celestial from cached value
func (b *OGame) GetCachedCelestialByCoord(coord ogame.Coordinate) Celestial {
	for _, p := range b.GetCachedPlanets() {
		if p.GetCoordinate().Equal(coord) {
			return p
		}
		if p.Moon != nil && p.Moon.GetCoordinate().Equal(coord) {
			return *p.Moon
		}
	}
	return nil
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

func (b *OGame) getTasks() (out taskRunner.TasksOverview) {
	return b.taskRunnerInst.GetTasks()
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
	return b.client
}

// SetClient set the http client used by the bot
func (b *OGame) SetClient(client *OGameClient) {
	b.client = client
}

// GetLoginClient get the http client used by the bot for login operations
func (b *OGame) GetLoginClient() *OGameClient {
	return b.client
}

// GetPublicIP get the public IP used by the bot
func (b *OGame) GetPublicIP() (string, error) {
	return b.getPublicIP()
}

// ValidateAccount validate a gameforge account
func (b *OGame) ValidateAccount(code string) error {
	return b.validateAccount(code)
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
func (b *OGame) AddAccount(number int, lang string) (*AddAccountRes, error) {
	return b.addAccount(number, lang)
}

// WithPriority ...
func (b *OGame) WithPriority(priority taskRunner.Priority) Prioritizable {
	return b.taskRunnerInst.WithPriority(priority)
}

// Begin start a transaction. Once this function is called, "Done" must be called to release the lock.
func (b *OGame) Begin() Prioritizable {
	return b.WithPriority(taskRunner.Normal).Begin()
}

// BeginNamed begins a new transaction with a name. "Done" must be called to release the lock.
func (b *OGame) BeginNamed(name string) Prioritizable {
	return b.WithPriority(taskRunner.Normal).BeginNamed(name)
}

// SetInitiator ...
func (b *OGame) SetInitiator(initiator string) Prioritizable {
	return nil
}

// Done ...
func (b *OGame) Done() {}

// Tx locks the bot during the transaction and ensure the lock is released afterward
func (b *OGame) Tx(clb func(tx Prioritizable) error) error {
	return b.WithPriority(taskRunner.Normal).Tx(clb)
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
	b.client.SetUserAgent(newUserAgent)
}

// LoginWithBearerToken to ogame server reusing existing token
func (b *OGame) LoginWithBearerToken(token string) (bool, error) {
	return b.WithPriority(taskRunner.Normal).LoginWithBearerToken(token)
}

// LoginWithExistingCookies to ogame server reusing existing cookies
func (b *OGame) LoginWithExistingCookies() (bool, error) {
	return b.WithPriority(taskRunner.Normal).LoginWithExistingCookies()
}

// Login to ogame server
// Can fails with BadCredentialsError
func (b *OGame) Login() error {
	return b.WithPriority(taskRunner.Normal).Login()
}

// Logout the bot from ogame server
func (b *OGame) Logout() { b.WithPriority(taskRunner.Normal).Logout() }

// BytesDownloaded returns the amount of bytes downloaded
func (b *OGame) BytesDownloaded() int64 {
	return b.client.BytesDownloaded()
}

// BytesUploaded returns the amount of bytes uploaded
func (b *OGame) BytesUploaded() int64 {
	return b.client.BytesUploaded()
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

// IsPioneers either or not the bot use lobby-pioneers
func (b *OGame) IsPioneers() bool {
	return b.lobby == LobbyPioneers
}

// IsDonutGalaxy shortcut to get ogame galaxy donut config
func (b *OGame) IsDonutGalaxy() bool {
	return b.isDonutGalaxy()
}

// IsDonutSystem shortcut to get ogame system donut config
func (b *OGame) IsDonutSystem() bool {
	return b.isDonutSystem()
}

// ConstructionTime get duration to build something
func (b *OGame) ConstructionTime(id ogame.ID, nbr int64, facilities ogame.Facilities) time.Duration {
	return b.constructionTime(id, nbr, facilities)
}

// FleetDeutSaveFactor returns the fleet deut save factor
func (b *OGame) FleetDeutSaveFactor() float64 {
	return b.serverData.GlobalDeuteriumSaveFactor
}

// GetPageContent gets the html for a specific ogame page
func (b *OGame) GetPageContent(vals url.Values) ([]byte, error) {
	return b.WithPriority(taskRunner.Normal).GetPageContent(vals)
}

// PostPageContent make a post request to ogame server
// This is useful when simulating a web browser
func (b *OGame) PostPageContent(vals, payload url.Values) ([]byte, error) {
	return b.WithPriority(taskRunner.Normal).PostPageContent(vals, payload)
}

// IsUnderAttack returns true if the user is under attack, false otherwise
func (b *OGame) IsUnderAttack() (bool, error) {
	return b.WithPriority(taskRunner.Normal).IsUnderAttack()
}

// GetCachedPlayer returns cached player infos
func (b *OGame) GetCachedPlayer() ogame.UserInfos {
	return b.Player
}

// GetCachedPreferences returns cached preferences
func (b *OGame) GetCachedPreferences() ogame.Preferences {
	return b.CachedPreferences
}

// SetVacationMode puts account in vacation mode
func (b *OGame) SetVacationMode() error {
	return b.WithPriority(taskRunner.Normal).SetVacationMode()
}

// IsVacationModeEnabled returns either or not the bot is in vacation mode
func (b *OGame) IsVacationModeEnabled() bool {
	return b.isVacationModeEnabled
}

// GetPlanets returns the user planets
func (b *OGame) GetPlanets() []Planet {
	return b.WithPriority(taskRunner.Normal).GetPlanets()
}

// GetCachedPlanets return planets from cached value
func (b *OGame) GetCachedPlanets() []Planet {
	b.planetsMu.RLock()
	defer b.planetsMu.RUnlock()
	return b.planets
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
func (b *OGame) GetCachedCelestial(v any) Celestial {
	return b.getCachedCelestial(v)
}

// GetPlanet gets infos for planetID
// Fails if planetID is invalid
func (b *OGame) GetPlanet(v any) (Planet, error) {
	return b.WithPriority(taskRunner.Normal).GetPlanet(v)
}

// GetMoons returns the user moons
func (b *OGame) GetMoons() []Moon {
	return b.WithPriority(taskRunner.Normal).GetMoons()
}

// GetMoon gets infos for moonID
func (b *OGame) GetMoon(v any) (Moon, error) {
	return b.WithPriority(taskRunner.Normal).GetMoon(v)
}

// GetCelestials get the player's planets & moons
func (b *OGame) GetCelestials() ([]Celestial, error) {
	return b.WithPriority(taskRunner.Normal).GetCelestials()
}

// RecruitOfficer recruit an officer.
// Typ 2: Commander, 3: Admiral, 4: Engineer, 5: Geologist, 6: Technocrat
// Days: 7 or 90
func (b *OGame) RecruitOfficer(typ, days int64) error {
	return b.WithPriority(taskRunner.Normal).RecruitOfficer(typ, days)
}

// Abandon a planet
func (b *OGame) Abandon(v any) error {
	return b.WithPriority(taskRunner.Normal).Abandon(v)
}

// GetCelestial get the player's planet/moon using the coordinate
func (b *OGame) GetCelestial(v any) (Celestial, error) {
	return b.WithPriority(taskRunner.Normal).GetCelestial(v)
}

// ServerVersion returns OGame version
func (b *OGame) ServerVersion() string {
	return b.serverData.Version
}

// ServerTime returns server time
// Timezone is OGT (OGame Time zone)
func (b *OGame) ServerTime() time.Time {
	return b.WithPriority(taskRunner.Normal).ServerTime()
}

// Location returns bot Time zone.
func (b *OGame) Location() *time.Location {
	return b.location
}

// GetUserInfos gets the user information
func (b *OGame) GetUserInfos() ogame.UserInfos {
	return b.WithPriority(taskRunner.Normal).GetUserInfos()
}

// SendMessage sends a message to playerID
func (b *OGame) SendMessage(playerID int64, message string) error {
	return b.WithPriority(taskRunner.Normal).SendMessage(playerID, message)
}

// SendMessageAlliance sends a message to associationID
func (b *OGame) SendMessageAlliance(associationID int64, message string) error {
	return b.WithPriority(taskRunner.Normal).SendMessageAlliance(associationID, message)
}

// GetFleets get the player's own fleets activities
func (b *OGame) GetFleets(opts ...Option) ([]ogame.Fleet, ogame.Slots) {
	return b.WithPriority(taskRunner.Normal).GetFleets(opts...)
}

// GetFleetsFromEventList get the player's own fleets activities
func (b *OGame) GetFleetsFromEventList() []ogame.Fleet {
	return b.WithPriority(taskRunner.Normal).GetFleetsFromEventList()
}

// CancelFleet cancel a fleet
func (b *OGame) CancelFleet(fleetID ogame.FleetID) error {
	return b.WithPriority(taskRunner.Normal).CancelFleet(fleetID)
}

// GetAttacks get enemy fleets attacking you
func (b *OGame) GetAttacks(opts ...Option) ([]ogame.AttackEvent, error) {
	return b.WithPriority(taskRunner.Normal).GetAttacks(opts...)
}

// GalaxyInfos get information of all planets and moons of a solar system
func (b *OGame) GalaxyInfos(galaxy, system int64, options ...Option) (ogame.SystemInfos, error) {
	return b.WithPriority(taskRunner.Normal).GalaxyInfos(galaxy, system, options...)
}

// GetResourceSettings gets the resources settings for specified planetID
func (b *OGame) GetResourceSettings(planetID ogame.PlanetID, options ...Option) (ogame.ResourceSettings, error) {
	return b.WithPriority(taskRunner.Normal).GetResourceSettings(planetID, options...)
}

// SetResourceSettings set the resources settings on a planet
func (b *OGame) SetResourceSettings(planetID ogame.PlanetID, settings ogame.ResourceSettings) error {
	return b.WithPriority(taskRunner.Normal).SetResourceSettings(planetID, settings)
}

// GetResourcesBuildings gets the resources buildings levels
func (b *OGame) GetResourcesBuildings(celestialID ogame.CelestialID, options ...Option) (ogame.ResourcesBuildings, error) {
	return b.WithPriority(taskRunner.Normal).GetResourcesBuildings(celestialID, options...)
}

// GetDefense gets all the defenses units information of a planet
// Fails if planetID is invalid
func (b *OGame) GetDefense(celestialID ogame.CelestialID, options ...Option) (ogame.DefensesInfos, error) {
	return b.WithPriority(taskRunner.Normal).GetDefense(celestialID, options...)
}

// GetShips gets all ships units information of a planet
func (b *OGame) GetShips(celestialID ogame.CelestialID, options ...Option) (ogame.ShipsInfos, error) {
	return b.WithPriority(taskRunner.Normal).GetShips(celestialID, options...)
}

// GetFacilities gets all facilities information of a planet
func (b *OGame) GetFacilities(celestialID ogame.CelestialID, options ...Option) (ogame.Facilities, error) {
	return b.WithPriority(taskRunner.Normal).GetFacilities(celestialID, options...)
}

// GetProduction get what is in the production queue.
// (ships & defense being built)
func (b *OGame) GetProduction(celestialID ogame.CelestialID) ([]ogame.Quantifiable, int64, error) {
	return b.WithPriority(taskRunner.Normal).GetProduction(celestialID)
}

// GetCachedResearch returns cached researches
func (b *OGame) GetCachedResearch() ogame.Researches {
	return b.WithPriority(taskRunner.Normal).GetCachedResearch()
}

// GetResearch gets the player researches information
func (b *OGame) GetResearch() ogame.Researches {
	return b.WithPriority(taskRunner.Normal).GetResearch()
}

// GetSlots gets the player current and total slots information
func (b *OGame) GetSlots() ogame.Slots {
	return b.WithPriority(taskRunner.Normal).GetSlots()
}

// Build builds any ogame objects (building, technology, ship, defence)
func (b *OGame) Build(celestialID ogame.CelestialID, id ogame.ID, nbr int64) error {
	return b.WithPriority(taskRunner.Normal).Build(celestialID, id, nbr)
}

// TearDown tears down any ogame building
func (b *OGame) TearDown(celestialID ogame.CelestialID, id ogame.ID) error {
	return b.WithPriority(taskRunner.Normal).TearDown(celestialID, id)
}

// BuildCancelable builds any cancelable ogame objects (building, technology)
func (b *OGame) BuildCancelable(celestialID ogame.CelestialID, id ogame.ID) error {
	return b.WithPriority(taskRunner.Normal).BuildCancelable(celestialID, id)
}

// BuildProduction builds any line production ogame objects (ship, defence)
func (b *OGame) BuildProduction(celestialID ogame.CelestialID, id ogame.ID, nbr int64) error {
	return b.WithPriority(taskRunner.Normal).BuildProduction(celestialID, id, nbr)
}

// BuildBuilding ensure what is being built is a building
func (b *OGame) BuildBuilding(celestialID ogame.CelestialID, buildingID ogame.ID) error {
	return b.WithPriority(taskRunner.Normal).BuildBuilding(celestialID, buildingID)
}

// BuildDefense builds a defense unit
func (b *OGame) BuildDefense(celestialID ogame.CelestialID, defenseID ogame.ID, nbr int64) error {
	return b.WithPriority(taskRunner.Normal).BuildDefense(celestialID, defenseID, nbr)
}

// BuildShips builds a ship unit
func (b *OGame) BuildShips(celestialID ogame.CelestialID, shipID ogame.ID, nbr int64) error {
	return b.WithPriority(taskRunner.Normal).BuildShips(celestialID, shipID, nbr)
}

// ConstructionsBeingBuilt returns the building & research being built, and the time remaining (secs)
func (b *OGame) ConstructionsBeingBuilt(celestialID ogame.CelestialID) (ogame.ID, int64, ogame.ID, int64) {
	return b.WithPriority(taskRunner.Normal).ConstructionsBeingBuilt(celestialID)
}

// CancelBuilding cancel the construction of a building on a specified planet
func (b *OGame) CancelBuilding(celestialID ogame.CelestialID) error {
	return b.WithPriority(taskRunner.Normal).CancelBuilding(celestialID)
}

// CancelLfBuilding cancel the construction of a lifeform building on a specified planet
func (b *OGame) CancelLfBuilding(celestialID ogame.CelestialID) error {
	return b.WithPriority(taskRunner.Normal).CancelLfBuilding(celestialID)
}

// CancelResearch cancel the research
func (b *OGame) CancelResearch(celestialID ogame.CelestialID) error {
	return b.WithPriority(taskRunner.Normal).CancelResearch(celestialID)
}

// BuildTechnology ensure that we're trying to build a technology
func (b *OGame) BuildTechnology(celestialID ogame.CelestialID, technologyID ogame.ID) error {
	return b.WithPriority(taskRunner.Normal).BuildTechnology(celestialID, technologyID)
}

// GetResources gets user resources
func (b *OGame) GetResources(celestialID ogame.CelestialID) (ogame.Resources, error) {
	return b.WithPriority(taskRunner.Normal).GetResources(celestialID)
}

// GetResourcesDetails gets user resources
func (b *OGame) GetResourcesDetails(celestialID ogame.CelestialID) (ogame.ResourcesDetails, error) {
	return b.WithPriority(taskRunner.Normal).GetResourcesDetails(celestialID)
}

// GetTechs gets a celestial supplies/facilities/ships/researches
func (b *OGame) GetTechs(celestialID ogame.CelestialID) (ogame.ResourcesBuildings, ogame.Facilities, ogame.ShipsInfos, ogame.DefensesInfos, ogame.Researches, error) {
	return b.WithPriority(taskRunner.Normal).GetTechs(celestialID)
}

// SendFleet sends a fleet
func (b *OGame) SendFleet(celestialID ogame.CelestialID, ships []ogame.Quantifiable, speed ogame.Speed, where ogame.Coordinate,
	mission ogame.MissionID, resources ogame.Resources, holdingTime, unionID int64) (ogame.Fleet, error) {
	return b.WithPriority(taskRunner.Normal).SendFleet(celestialID, ships, speed, where, mission, resources, holdingTime, unionID)
}

// EnsureFleet either sends all the requested ships or fail
func (b *OGame) EnsureFleet(celestialID ogame.CelestialID, ships []ogame.Quantifiable, speed ogame.Speed, where ogame.Coordinate,
	mission ogame.MissionID, resources ogame.Resources, holdingTime, unionID int64) (ogame.Fleet, error) {
	return b.WithPriority(taskRunner.Normal).EnsureFleet(celestialID, ships, speed, where, mission, resources, holdingTime, unionID)
}

// DestroyRockets destroys anti-ballistic & inter-planetary missiles
func (b *OGame) DestroyRockets(planetID ogame.PlanetID, abm, ipm int64) error {
	return b.WithPriority(taskRunner.Normal).DestroyRockets(planetID, abm, ipm)
}

// SendIPM sends IPM
func (b *OGame) SendIPM(planetID ogame.PlanetID, coord ogame.Coordinate, nbr int64, priority ogame.ID) (int64, error) {
	return b.WithPriority(taskRunner.Normal).SendIPM(planetID, coord, nbr, priority)
}

// GetCombatReportSummaryFor gets the latest combat report for a given coordinate
func (b *OGame) GetCombatReportSummaryFor(coord ogame.Coordinate) (ogame.CombatReportSummary, error) {
	return b.WithPriority(taskRunner.Normal).GetCombatReportSummaryFor(coord)
}

// GetEspionageReportFor gets the latest espionage report for a given coordinate
func (b *OGame) GetEspionageReportFor(coord ogame.Coordinate) (ogame.EspionageReport, error) {
	return b.WithPriority(taskRunner.Normal).GetEspionageReportFor(coord)
}

// GetExpeditionMessages gets the expedition messages
func (b *OGame) GetExpeditionMessages() ([]ogame.ExpeditionMessage, error) {
	return b.WithPriority(taskRunner.Normal).GetExpeditionMessages()
}

// GetExpeditionMessageAt gets the expedition message for time t
func (b *OGame) GetExpeditionMessageAt(t time.Time) (ogame.ExpeditionMessage, error) {
	return b.WithPriority(taskRunner.Normal).GetExpeditionMessageAt(t)
}

// CollectAllMarketplaceMessages collect all marketplace messages
func (b *OGame) CollectAllMarketplaceMessages() error {
	return b.WithPriority(taskRunner.Normal).CollectAllMarketplaceMessages()
}

// CollectMarketplaceMessage collect marketplace message
func (b *OGame) CollectMarketplaceMessage(msg ogame.MarketplaceMessage) error {
	return b.WithPriority(taskRunner.Normal).CollectMarketplaceMessage(msg)
}

// GetEspionageReportMessages gets the summary of each espionage reports
func (b *OGame) GetEspionageReportMessages() ([]ogame.EspionageReportSummary, error) {
	return b.WithPriority(taskRunner.Normal).GetEspionageReportMessages()
}

// GetEspionageReport gets a detailed espionage report
func (b *OGame) GetEspionageReport(msgID int64) (ogame.EspionageReport, error) {
	return b.WithPriority(taskRunner.Normal).GetEspionageReport(msgID)
}

// DeleteMessage deletes a message from the mail box
func (b *OGame) DeleteMessage(msgID int64) error {
	return b.WithPriority(taskRunner.Normal).DeleteMessage(msgID)
}

// DeleteAllMessagesFromTab deletes all messages from a tab in the mail box
func (b *OGame) DeleteAllMessagesFromTab(tabID ogame.MessagesTabID) error {
	return b.WithPriority(taskRunner.Normal).DeleteAllMessagesFromTab(tabID)
}

// GetResourcesProductions gets the planet resources production
func (b *OGame) GetResourcesProductions(planetID ogame.PlanetID) (ogame.Resources, error) {
	return b.WithPriority(taskRunner.Normal).GetResourcesProductions(planetID)
}

// GetResourcesProductionsLight gets the planet resources production
func (b *OGame) GetResourcesProductionsLight(resBuildings ogame.ResourcesBuildings, researches ogame.Researches,
	resSettings ogame.ResourceSettings, temp ogame.Temperature) ogame.Resources {
	return b.WithPriority(taskRunner.Normal).GetResourcesProductionsLight(resBuildings, researches, resSettings, temp)
}

// FlightTime calculate flight time and fuel needed
func (b *OGame) FlightTime(origin, destination ogame.Coordinate, speed ogame.Speed, ships ogame.ShipsInfos, missionID ogame.MissionID) (secs, fuel int64) {
	return b.WithPriority(taskRunner.Normal).FlightTime(origin, destination, speed, ships, missionID)
}

// Distance return distance between two coordinates
func (b *OGame) Distance(origin, destination ogame.Coordinate) int64 {
	return Distance(origin, destination, b.serverData.Galaxies, b.serverData.Systems, b.serverData.DonutGalaxy, b.serverData.DonutSystem)
}

// RegisterWSCallback ...
func (b *OGame) RegisterWSCallback(id string, fn func(msg []byte)) {
	b.Lock()
	defer b.Unlock()
	b.wsCallbacks[id] = fn
}

// RemoveWSCallback ...
func (b *OGame) RemoveWSCallback(id string) {
	b.Lock()
	defer b.Unlock()
	delete(b.wsCallbacks, id)
}

// RegisterChatCallback register a callback that is called when chat messages are received
func (b *OGame) RegisterChatCallback(fn func(msg ogame.ChatMsg)) {
	b.chatCallbacks = append(b.chatCallbacks, fn)
}

// RegisterAuctioneerCallback register a callback that is called when auctioneer packets are received
func (b *OGame) RegisterAuctioneerCallback(fn func(packet any)) {
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
func (b *OGame) Phalanx(moonID ogame.MoonID, coord ogame.Coordinate) ([]ogame.Fleet, error) {
	return b.WithPriority(taskRunner.Normal).Phalanx(moonID, coord)
}

// UnsafePhalanx same as Phalanx but does not perform any input validation.
func (b *OGame) UnsafePhalanx(moonID ogame.MoonID, coord ogame.Coordinate) ([]ogame.Fleet, error) {
	return b.WithPriority(taskRunner.Normal).UnsafePhalanx(moonID, coord)
}

// JumpGateDestinations returns available destinations for jump gate.
func (b *OGame) JumpGateDestinations(origin ogame.MoonID) (moonIDs []ogame.MoonID, rechargeCountdown int64, err error) {
	return b.WithPriority(taskRunner.Normal).JumpGateDestinations(origin)
}

// JumpGate sends ships through a jump gate.
func (b *OGame) JumpGate(origin, dest ogame.MoonID, ships ogame.ShipsInfos) (success bool, rechargeCountdown int64, err error) {
	return b.WithPriority(taskRunner.Normal).JumpGate(origin, dest, ships)
}

// BuyOfferOfTheDay buys the offer of the day.
func (b *OGame) BuyOfferOfTheDay() error {
	return b.WithPriority(taskRunner.Normal).BuyOfferOfTheDay()
}

// CreateUnion creates a union
func (b *OGame) CreateUnion(fleet ogame.Fleet, users []string) (int64, error) {
	return b.WithPriority(taskRunner.Normal).CreateUnion(fleet, users)
}

// HeadersForPage gets the headers for a specific ogame page
func (b *OGame) HeadersForPage(url string) (http.Header, error) {
	return b.WithPriority(taskRunner.Normal).HeadersForPage(url)
}

// GetEmpire gets all planets/moons information resources/supplies/facilities/ships/researches
func (b *OGame) GetEmpire(celestialType ogame.CelestialType) ([]ogame.EmpireCelestial, error) {
	return b.WithPriority(taskRunner.Normal).GetEmpire(celestialType)
}

// GetEmpireJSON retrieves JSON from Empire page (Commander only).
func (b *OGame) GetEmpireJSON(nbr int64) (any, error) {
	return b.WithPriority(taskRunner.Normal).GetEmpireJSON(nbr)
}

// CharacterClass returns the bot character class
func (b *OGame) CharacterClass() ogame.CharacterClass {
	return b.characterClass
}

// GetAuction ...
func (b *OGame) GetAuction() (ogame.Auction, error) {
	return b.WithPriority(taskRunner.Normal).GetAuction()
}

// DoAuction ...
func (b *OGame) DoAuction(bid map[ogame.CelestialID]ogame.Resources) error {
	return b.WithPriority(taskRunner.Normal).DoAuction(bid)
}

// Highscore ...
func (b *OGame) Highscore(category, typ, page int64) (ogame.Highscore, error) {
	return b.WithPriority(taskRunner.Normal).Highscore(category, typ, page)
}

// GetAllResources gets the resources of all planets and moons
func (b *OGame) GetAllResources() (map[ogame.CelestialID]ogame.Resources, error) {
	return b.WithPriority(taskRunner.Normal).GetAllResources()
}

// GetTasks return how many tasks are queued in the heap.
func (b *OGame) GetTasks() taskRunner.TasksOverview {
	return b.getTasks()
}

// GetDMCosts returns fast build with DM information
func (b *OGame) GetDMCosts(celestialID ogame.CelestialID) (ogame.DMCosts, error) {
	return b.WithPriority(taskRunner.Normal).GetDMCosts(celestialID)
}

// UseDM use dark matter to fast build
func (b *OGame) UseDM(typ string, celestialID ogame.CelestialID) error {
	return b.WithPriority(taskRunner.Normal).UseDM(typ, celestialID)
}

// GetItems get all items information
func (b *OGame) GetItems(celestialID ogame.CelestialID) ([]ogame.Item, error) {
	return b.WithPriority(taskRunner.Normal).GetItems(celestialID)
}

// GetActiveItems ...
func (b *OGame) GetActiveItems(celestialID ogame.CelestialID) ([]ogame.ActiveItem, error) {
	return b.WithPriority(taskRunner.Normal).GetActiveItems(celestialID)
}

// ActivateItem activate an item
func (b *OGame) ActivateItem(ref string, celestialID ogame.CelestialID) error {
	return b.WithPriority(taskRunner.Normal).ActivateItem(ref, celestialID)
}

// BuyMarketplace buy an item on the marketplace
func (b *OGame) BuyMarketplace(itemID int64, celestialID ogame.CelestialID) error {
	return b.WithPriority(taskRunner.Normal).BuyMarketplace(itemID, celestialID)
}

// OfferSellMarketplace sell offer on marketplace
func (b *OGame) OfferSellMarketplace(itemID any, quantity, priceType, price, priceRange int64, celestialID ogame.CelestialID) error {
	return b.WithPriority(taskRunner.Normal).OfferSellMarketplace(itemID, quantity, priceType, price, priceRange, celestialID)
}

// OfferBuyMarketplace buy offer on marketplace
func (b *OGame) OfferBuyMarketplace(itemID any, quantity, priceType, price, priceRange int64, celestialID ogame.CelestialID) error {
	return b.WithPriority(taskRunner.Normal).OfferBuyMarketplace(itemID, quantity, priceType, price, priceRange, celestialID)
}
