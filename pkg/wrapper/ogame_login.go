package wrapper

import (
	"context"
	"errors"
	"fmt"
	"github.com/alaingilbert/clockwork"
	"github.com/alaingilbert/ogame/pkg/exponentialBackoff"
	"github.com/alaingilbert/ogame/pkg/extractor"
	v10 "github.com/alaingilbert/ogame/pkg/extractor/v10"
	v104 "github.com/alaingilbert/ogame/pkg/extractor/v104"
	v11 "github.com/alaingilbert/ogame/pkg/extractor/v11"
	"github.com/alaingilbert/ogame/pkg/extractor/v11_13_0"
	"github.com/alaingilbert/ogame/pkg/extractor/v11_15_0"
	"github.com/alaingilbert/ogame/pkg/extractor/v11_9_0"
	"github.com/alaingilbert/ogame/pkg/extractor/v12_0_0"
	v7 "github.com/alaingilbert/ogame/pkg/extractor/v7"
	v71 "github.com/alaingilbert/ogame/pkg/extractor/v71"
	v8 "github.com/alaingilbert/ogame/pkg/extractor/v8"
	v874 "github.com/alaingilbert/ogame/pkg/extractor/v874"
	v9 "github.com/alaingilbert/ogame/pkg/extractor/v9"
	"github.com/alaingilbert/ogame/pkg/gameforge"
	"github.com/alaingilbert/ogame/pkg/httpclient"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/parser"
	"github.com/alaingilbert/ogame/pkg/utils"
	"github.com/hashicorp/go-version"
	cookiejar "github.com/orirawlings/persistent-cookiejar"
	"net/http"
	"regexp"
	"time"
)

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

// Return either or not the bot logged in using the provided bearer token.
func (b *OGame) loginWithBearerToken(token string) (bool, error) {
	botLoginFn := b.login
	if token == "" {
		err := botLoginFn()
		return false, err
	}
	b.bearerToken = token
	server, userAccount, err := b.loginPart1(token)
	if errors.Is(err, context.Canceled) ||
		errors.Is(err, gameforge.ErrAccountBlocked) {
		return false, err
	} else if err != nil {
		err := botLoginFn()
		return false, err
	}

	if err := b.loginPart2(server); err != nil {
		return false, err
	}

	loginOpts := []Option{SkipRetry, SkipCacheFullPage}
	page, err := getPage[parser.OverviewPage](b, loginOpts...)
	if err != nil {
		if errors.Is(err, ogame.ErrNotLogged) {
			loginLink, pageHTML, err := b.getAndExecLoginLink(userAccount, token)
			if err != nil {
				return true, err
			}
			page, err := getPage[parser.OverviewPage](b, loginOpts...)
			if err != nil {
				if errors.Is(err, ogame.ErrNotLogged) {
					err := botLoginFn()
					return false, err
				}
				return false, err
			}
			b.debug("login using existing cookies")
			if err := b.loginPart3Tmp(userAccount, page, loginLink, pageHTML); err != nil {
				return false, err
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
	token := b.bearerToken
	if token == "" {
		token = b.getBearerTokenFromCookie()
	}
	return b.loginWithBearerToken(token)
}

func (b *OGame) getBearerTokenFromCookie() string {
	cookies := b.device.GetClient().Jar.(*cookiejar.Jar).AllCookies()
	for _, c := range cookies {
		if c.Name == gameforge.TokenCookieName {
			return c.Value
		}
	}
	return ""
}

func (b *OGame) login() error {
	b.debug("post sessions")
	postSessionsRes, err := postSessions(b)
	if err != nil {
		return err
	}
	token := postSessionsRes.Token

	server, userAccount, err := b.loginPart1(token)
	if err != nil {
		return err
	}

	loginLink, pageHTML, err := b.getAndExecLoginLink(userAccount, token)
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
	if err := b.loginPart3Tmp(userAccount, page, loginLink, pageHTML); err != nil {
		return err
	}
	return nil
}

func (b *OGame) getAndExecLoginLink(userAccount gameforge.Account, token string) (string, []byte, error) {
	b.debug("get login link")
	loginLink, err := gameforge.GetLoginLink(b.ctx, b.device, PLATFORM, b.lobby, userAccount, token)
	if err != nil {
		return "", nil, err
	}
	b.debug("login to universe")
	var pageHTML []byte
	err = b.device.GetClient().WithTransport(b.loginProxyTransport, func(client *httpclient.Client) error {
		pageHTML, err = gameforge.ExecLoginLink(b.ctx, client, loginLink)
		return err
	})
	if err != nil {
		return "", nil, err
	}
	return loginLink, pageHTML, nil
}

func (b *OGame) loginPart3Tmp(userAccount gameforge.Account, page *parser.OverviewPage, loginLink string, pageHTML []byte) error {
	if err := b.loginPart3(userAccount, page); err != nil {
		return err
	}
	if err := b.device.GetClient().Jar.(*cookiejar.Jar).Save(); err != nil {
		return err
	}
	b.execInterceptorCallbacks(http.MethodGet, loginLink, nil, nil, pageHTML)
	return nil
}

// Get user's accounts, get GF ogame servers, then find and return the server and userAccount that we asked to play in.
func (b *OGame) loginPart1(token string) (server gameforge.Server, userAccount gameforge.Account, err error) {
	client := b.device.GetClient()
	ctx := b.ctx
	lobby := b.lobby
	b.debug("find account & server for universe")
	userAccount, server, err = gameforge.GetServerAccount(ctx, client, PLATFORM, lobby, token, b.universe, b.language, b.playerID)
	if err != nil {
		return
	}
	if userAccount.Blocked {
		return server, userAccount, gameforge.ErrAccountBlocked
	}
	b.debug(fmt.Sprintf("Players online: %d, Players: %d", server.PlayersOnline, server.PlayerCount))
	return
}

func (b *OGame) loginPart2(server gameforge.Server) error {
	b.isLoggedInAtom.Store(true) // At this point, we are logged in
	b.isConnectedAtom.Store(true)
	// Get server data
	start := time.Now()
	b.server = server
	serverData, err := getServerData(b.ctx, b.device.GetClient(), b.server.Number, b.server.Language)
	if err != nil {
		return err
	}
	serverData.SpeedFleetWar = utils.MaxInt(serverData.SpeedFleetWar, 1)
	serverData.SpeedFleetPeaceful = utils.MaxInt(serverData.SpeedFleetPeaceful, 1)
	serverData.SpeedFleetHolding = utils.MaxInt(serverData.SpeedFleetHolding, 1)
	if serverData.SpeedFleet == 0 {
		serverData.SpeedFleet = serverData.SpeedFleetPeaceful
	}
	b.cache.serverData = serverData
	lang := server.Language
	if server.Language == "yu" {
		lang = "ba"
	}
	b.language = lang
	b.cache.serverURL = fmt.Sprintf("https://s%d-%s.ogame.gameforge.com", server.Number, lang)
	b.debug("get server data", time.Since(start))
	return nil
}

func (b *OGame) loginPart3(userAccount gameforge.Account, page *parser.OverviewPage) error {
	var ext extractor.Extractor = v12_0_0.NewExtractor()
	r := regexp.MustCompile(`(\d+\.\d+\.\d+)`)
	versionMatches := r.FindStringSubmatch(b.cache.serverData.Version)
	versionMatch := b.cache.serverData.Version
	if len(versionMatches) == 2 {
		versionMatch = versionMatches[1]
	}
	if ogVersion, err := version.NewVersion(versionMatch); err == nil {
		if isVGreaterThanOrEqual(ogVersion, "12.0.0") {
			ext = v12_0_0.NewExtractor()
		} else if isVGreaterThanOrEqual(ogVersion, "11.15.0") {
			ext = v11_15_0.NewExtractor()
		} else if isVGreaterThanOrEqual(ogVersion, "11.13.0") {
			ext = v11_13_0.NewExtractor()
		} else if isVGreaterThanOrEqual(ogVersion, "11.9.0") {
			ext = v11_9_0.NewExtractor()
		} else if isVGreaterThanOrEqual(ogVersion, "11.0.0-beta25") {
			ext = v11.NewExtractor()
		} else if isVGreaterThanOrEqual(ogVersion, "10.4.0-beta2") {
			ext = v104.NewExtractor()
		} else if isVGreaterThanOrEqual(ogVersion, "10.0.0") {
			ext = v10.NewExtractor()
		} else if isVGreaterThanOrEqual(ogVersion, "9.0.0") {
			ext = v9.NewExtractor()
		} else if isVGreaterThanOrEqual(ogVersion, "8.7.4-pl3") {
			ext = v874.NewExtractor()
		} else if isVGreaterThanOrEqual(ogVersion, "8.0.0") {
			ext = v8.NewExtractor()
		} else if isVGreaterThanOrEqual(ogVersion, "7.1.0-rc0") {
			ext = v71.NewExtractor()
		} else if isVGreaterThanOrEqual(ogVersion, "7.0.0") {
			ext = v7.NewExtractor()
		}
		ext.SetLanguage(b.language)
		ext.SetLifeformEnabled(page.ExtractLifeformEnabled())
	} else {
		b.error("failed to parse ogame version: " + err.Error())
	}

	b.sessionChatCounter = 1

	b.debug("logged in as " + userAccount.Name + " on " + b.universe + "-" + b.language)

	b.debug("extract information from html")
	b.cache.ogameSession = page.ExtractOGameSession()
	if b.cache.ogameSession == "" {
		return gameforge.ErrBadCredentials
	}

	serverTime, err := page.ExtractServerTime()
	if err != nil {
		b.error(err)
	}
	b.cache.location = serverTime.Location()

	ext.SetLocation(b.cache.location)
	b.extractor = ext

	preferencesPage, err := getPage[parser.PreferencesPage](b, SkipCacheFullPage)
	if err != nil {
		b.error(err)
	}
	b.cache.CachedPreferences = preferencesPage.ExtractPreferences()
	language := b.cache.serverData.Language
	if b.cache.CachedPreferences.Language != "" {
		language = b.cache.CachedPreferences.Language
	}
	ext.SetLanguage(language)
	b.extractor = ext

	page.SetExtractor(ext)

	b.cacheFullPageInfo(page)

	// Extract chat host and port
	m := regexp.MustCompile(`var nodeUrl\s?=\s?"https:\\/\\/([^:]+):(\d+)\\/socket.io\\/socket.io.js"`).FindSubmatch(page.GetContent())
	chatHost := string(m[1])
	chatPort := string(m[2])

	if b.chatConnectedAtom.CompareAndSwap(false, true) {
		b.closeChatCh = make(chan struct{})
		go func(b *OGame) {
			defer b.chatConnectedAtom.Store(false)
			chatRetry := exponentialBackoff.New(context.Background(), clockwork.NewRealClock(), 60)
			chatRetry.LoopForever(func() bool {
				select {
				case <-b.closeChatCh:
					return false
				default:
					b.connectChat(chatRetry, chatHost, chatPort)
				}
				return true
			})
		}(b)
	} else {
		b.ReconnectChat()
	}

	// V11 Intro bypass
	if err := b.introBypass(page); err != nil {
		b.error("failed to bypass intro:", err)
	}

	return nil
}
