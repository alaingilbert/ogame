package gameforge

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// TokenCookieName gameforge cookie name for token id
const (
	TokenCookieName         = "gf-token-production"
	ChallengeIDCookieName   = "gf-challenge-id"
	acceptEncodingHeaderKey = "Accept-Encoding"
	contentTypeHeaderKey    = "Content-Type"
	authorizationHeaderKey  = "Authorization"
	twoFactorHeaderKey      = "tnt-2fa-code"
	installationIDHeaderKey = "tnt-installation-id"
	applicationJson         = "application/json"
	gzipEncoding            = "gzip, deflate, br"
	challengeBaseURL        = "https://challenge.gameforge.com"
	imgDropChallengeBaseURL = "https://image-drop-challenge.gameforge.com"
	endpointLoc             = "en-GB"
)

type GameforgeClient interface {
	Login(params *GfLoginParams) (out *LoginResponse, err error)
	GetUserAccounts() ([]Account, error)
	GetServers() ([]Server, error)
	StartCaptchaChallenge(challengeID string) (questionRaw, iconsRaw []byte, err error)
	SolveChallenge(challengeID string, answer int64) error
	Register(email, password, challengeID, lang string) error
	RedeemCode(code string) error
	AddAccount(serverName, lang string) (*AddAccountResponse, error)
	ValidateAccount(code string) error
	GetServerAccount(serverName, lang string, playerID int64) (account Account, server Server, err error)
	GetLoginLink(userAccount Account) (string, error)
	ExecLoginLink(loginLink string) ([]byte, error)
}

// Compile time checks to ensure type satisfies IGameforge interface
var _ GameforgeClient = (*Gameforge)(nil)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Device is anything that can make http requests and solve the gameforge blackbox
type Device interface {
	HttpClient
	GetBlackbox() (string, error)
}

func Ternary[T any](predicate bool, a, b T) T {
	if predicate {
		return a
	}
	return b
}

// Or return "a" if it is non-zero otherwise "b"
func Or[T comparable](a, b T) (zero T) {
	return Ternary(a != zero, a, b)
}

type CaptchaRequiredError struct {
	ChallengeID string
}

func NewCaptchaRequiredError(challengeID string) *CaptchaRequiredError {
	return &CaptchaRequiredError{ChallengeID: challengeID}
}

func (e CaptchaRequiredError) Error() string {
	return fmt.Sprintf("captcha required, %s", e.ChallengeID)
}

type RegisterError struct{ ErrorString string }

func (e *RegisterError) Error() string { return e.ErrorString }

var (
	ErrEmailInvalid    = &RegisterError{"Please enter a valid email address."}
	ErrEmailUsed       = &RegisterError{"Failed to create new lobby, email already used."}
	ErrPasswordInvalid = &RegisterError{"Must contain at least 10 characters including at least one upper and lowercase letter and a number."}
)

// ErrBadCredentials returned when the provided credentials are invalid
var ErrBadCredentials = errors.New("bad credentials")

// ErrOTPRequired returned when the otp is required
var ErrOTPRequired = errors.New("otp required")

// ErrOTPInvalid returned when the otp is invalid
var ErrOTPInvalid = errors.New("otp invalid")

// ErrLoginLink returned when account is somewhat banned, cannot login for no apparent reason
var ErrLoginLink = errors.New("failed to get login link")

// ErrAccountNotFound returned when the account is not found
var ErrAccountNotFound = errors.New("account not found")

// ErrAccountBlocked returned when account is banned
var ErrAccountBlocked = errors.New("account is blocked")

type GfLoginParams struct {
	Username    string
	Password    string
	OtpSecret   string
	ChallengeID string
}

type gfLoginParams struct {
	*GfLoginParams
	Ctx      context.Context
	Device   Device
	platform Platform
	lobby    string
}

// CaptchaSolver the returned answer should be one of "0" "1" "2" "3"
type CaptchaSolver func(ctx context.Context, question, icons []byte) (int64, error)

func getChallengeURL(base, challengeID string) string {
	return fmt.Sprintf("%s/challenge/%s", base, challengeID)
}

const blackboxPrefix = "tra:"

type Platform string

const (
	OGAME   Platform = "ogame"
	IKARIAM Platform = "ikariam"
)

// Lobby constants
const (
	Lobby         = "lobby"
	LobbyPioneers = "lobby-pioneers"
)

func (p Platform) isValid() bool {
	return p == OGAME || p == IKARIAM
}

// Gameforge ...
type Gameforge struct {
	ctx               context.Context
	lobby             string
	platform          Platform
	device            Device
	solver            CaptchaSolver
	maxCaptchaRetries int
	bearerToken       string
}

type Config struct {
	Ctx               context.Context
	Device            Device
	Solver            CaptchaSolver
	MaxCaptchaRetries *int // default to 3
	Platform          Platform
	Lobby             string
}

// New ...
func New(config *Config) (*Gameforge, error) {
	if config.Device == nil {
		return nil, errors.New("device is required")
	}
	if config.Ctx == nil {
		config.Ctx = context.Background()
	}
	if !config.Platform.isValid() {
		return nil, errors.New("invalid platform")
	}
	if config.MaxCaptchaRetries == nil {
		maxCaptchaRetries := 3
		config.MaxCaptchaRetries = &maxCaptchaRetries
	}
	return &Gameforge{
		ctx:               config.Ctx,
		device:            config.Device,
		platform:          config.Platform,
		lobby:             config.Lobby,
		solver:            config.Solver,
		maxCaptchaRetries: *config.MaxCaptchaRetries,
	}, nil
}

// Login do the gameforge login, if we get a captcha, solve the captcha and retry login.
// If no "solver" have been set or "maxCaptchaRetries" is 0, then it will not try to solve the captcha
func (g *Gameforge) Login(params *GfLoginParams) (out *LoginResponse, err error) {
	solver := g.solver
	maxTry := g.maxCaptchaRetries
	ctx := g.ctx
	device := g.device
LOGIN:
	out, err = login(&gfLoginParams{GfLoginParams: params, Device: device, Ctx: ctx, platform: g.platform, lobby: g.lobby})
	if err != nil {
		var captchaErr *CaptchaRequiredError
		if errors.As(err, &captchaErr) {
			if maxTry <= 0 || solver == nil {
				return nil, err
			}
			maxTry--
			if err := solveCaptcha(ctx, device, captchaErr.ChallengeID, solver); err != nil {
				return nil, err
			}
			goto LOGIN
		}
		return nil, err
	}
	g.bearerToken = out.Token
	return out, nil
}

// GetUserAccounts ...
func (g *Gameforge) GetUserAccounts() ([]Account, error) {
	return GetUserAccounts(g.ctx, g.device, g.platform, g.lobby, g.bearerToken)
}

// GetServers ...
func (g *Gameforge) GetServers() ([]Server, error) {
	return GetServers(g.ctx, g.device, g.platform, g.lobby)
}

// StartCaptchaChallenge ...
func (g *Gameforge) StartCaptchaChallenge(challengeID string) (questionRaw, iconsRaw []byte, err error) {
	return StartCaptchaChallenge(g.ctx, g.device, challengeID)
}

// SolveChallenge ...
func (g *Gameforge) SolveChallenge(challengeID string, answer int64) error {
	return SolveChallenge(g.ctx, g.device, challengeID, answer)
}

// Register ...
func (g *Gameforge) Register(email, password, challengeID, lang string) error {
	return Register(g.device, g.ctx, g.platform, g.lobby, email, password, challengeID, lang)
}

// RedeemCode ...
func (g *Gameforge) RedeemCode(code string) error {
	return RedeemCode(g.ctx, g.device, g.platform, g.lobby, g.bearerToken, code)
}

// AddAccount adds an account to a gameforge lobby
func (g *Gameforge) AddAccount(serverName, lang string) (*AddAccountResponse, error) {
	return AddAccountByServerNameLang(g.ctx, g.device, g.platform, g.lobby, g.bearerToken, serverName, lang)
}

// ValidateAccount ...
func (g *Gameforge) ValidateAccount(code string) error {
	return ValidateAccount(g.ctx, g.device, g.platform, g.lobby, code)
}

// GetServerAccount first get the gameforge servers and user accounts,
// then find the server and account that matches the given serverName, lang and playerID.
// PlayerID is optional, when given the zero value, the first server/account that matches will be returned.
func (g *Gameforge) GetServerAccount(serverName, lang string, playerID int64) (account Account, server Server, err error) {
	return GetServerAccount(g.ctx, g.device, g.platform, g.lobby, g.bearerToken, serverName, lang, playerID)
}

// GetLoginLink ...
func (g *Gameforge) GetLoginLink(userAccount Account) (string, error) {
	return GetLoginLink(g.ctx, g.device, g.platform, g.lobby, userAccount, g.bearerToken)
}

// ExecLoginLink ...
func (g *Gameforge) ExecLoginLink(loginLink string) ([]byte, error) {
	return ExecLoginLink(g.ctx, g.device, loginLink)
}

func solveCaptcha(ctx context.Context, client HttpClient, challengeID string, captchaCallback CaptchaSolver) error {
	questionRaw, iconsRaw, err := StartCaptchaChallenge(ctx, client, challengeID)
	if err != nil {
		return errors.New("failed to start captcha challenge: " + err.Error())
	}
	answer, err := captchaCallback(ctx, questionRaw, iconsRaw)
	if err != nil {
		return errors.New("failed to get answer for captcha challenge: " + err.Error())
	}
	if err := SolveChallenge(ctx, client, challengeID, answer); err != nil {
		return errors.New("failed to solve captcha challenge: " + err.Error())
	}
	return err
}

func getGameforgeLobbyBaseURL(lobby string, platform Platform) string {
	return fmt.Sprintf("https://%s.%s.gameforge.com", Or(lobby, Lobby), platform)
}

// Register a new gameforge lobby account
func Register(device Device, ctx context.Context, platform Platform, lobby, email, password, challengeID, lang string) error {
	blackbox, err := device.GetBlackbox()
	if err != nil {
		return err
	}
	var payload struct {
		Blackbox    string `json:"blackbox"`
		Credentials struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		} `json:"credentials"`
		Language string `json:"language"`
		Kid      string `json:"kid"`
	}
	payload.Blackbox = blackboxPrefix + blackbox
	payload.Credentials.Email = email
	payload.Credentials.Password = password
	payload.Language = Or(lang, "en")
	jsonPayloadBytes, err := json.Marshal(&payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, getGameforgeLobbyBaseURL(lobby, platform)+"/api/users", strings.NewReader(string(jsonPayloadBytes)))
	if err != nil {
		return err
	}
	if challengeID != "" {
		req.Header.Set(ChallengeIDCookieName, challengeID)
	}
	req.Header.Set(contentTypeHeaderKey, applicationJson)
	req.Header.Set(acceptEncodingHeaderKey, gzipEncoding)
	req.WithContext(ctx)
	resp, err := device.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusInternalServerError {
		return fmt.Errorf("gameforme internal server error : %s", resp.Status)
	}
	if resp.StatusCode == http.StatusConflict {
		if newChallengeID := extractChallengeID(resp); newChallengeID != "" {
			return NewCaptchaRequiredError(newChallengeID)
		}
	}
	by, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var res struct {
		MigrationRequired bool   `json:"migrationRequired"`
		Error             string `json:"error"`
	}
	if err := json.Unmarshal(by, &res); err != nil {
		return errors.New(err.Error() + " : " + string(by))
	}
	if res.Error == "email_invalid" {
		return ErrEmailInvalid
	} else if res.Error == "email_used" {
		return ErrEmailUsed
	} else if res.Error == "password_invalid" {
		return ErrPasswordInvalid
	} else if res.Error != "" {
		return errors.New(res.Error)
	}
	return nil
}

// ValidateAccount validate a gameforge account eg: ________-____-____-____-____________
func ValidateAccount(ctx context.Context, client HttpClient, platform Platform, lobby, code string) error {
	if len(code) != 36 {
		return errors.New("invalid validation code")
	}
	req, err := http.NewRequest(http.MethodPut, getGameforgeLobbyBaseURL(lobby, platform)+"/api/users/validate/"+code, strings.NewReader(`{"language":"en"}`))
	if err != nil {
		return err
	}
	req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to validate account: %s", resp.Status)
	}
	return nil
}

func buildBearerHeaderValue(token string) string { return "Bearer " + token }

func setDefaultParams(params *gfLoginParams) {
	if params.Ctx == nil {
		params.Ctx = context.Background()
	}
}

// RedeemCode ...
func RedeemCode(ctx context.Context, client HttpClient, platform Platform, lobby, bearerToken, code string) error {
	var payload struct {
		Token string `json:"token"`
	}
	payload.Token = code
	jsonPayloadBytes, err := json.Marshal(&payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, getGameforgeLobbyBaseURL(lobby, platform)+"/api/token", bytes.NewReader(jsonPayloadBytes))
	if err != nil {
		return err
	}
	req.Header.Set(authorizationHeaderKey, buildBearerHeaderValue(bearerToken))
	req.Header.Set(contentTypeHeaderKey, applicationJson)
	req.Header.Set(acceptEncodingHeaderKey, gzipEncoding)
	req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// {"tokenType":"accountTrading"}
	by, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusBadRequest {
		return errors.New("invalid request, token invalid ?")
	}
	var respParsed struct {
		TokenType string `json:"tokenType"`
	}
	if err := json.Unmarshal(by, &respParsed); err != nil {
		return errors.New(err.Error() + " : " + string(by))
	}
	if respParsed.TokenType != "accountTrading" {
		return errors.New("tokenType is not accountTrading")
	}
	return nil
}

func AddAccountByServerNameLang(ctx context.Context, device Device, platform Platform, lobby, bearerToken, serverName, lang string) (*AddAccountResponse, error) {
	servers, err := GetServers(ctx, device, platform, lobby)
	if err != nil {
		return nil, err
	}
	server, found := FindServer(serverName, lang, servers)
	if !found {
		return nil, errors.New("server not found")
	}
	return AddAccount(ctx, device, platform, lobby, server.AccountGroup, bearerToken)
}

// AddAccountResponse response from creating a new account
type AddAccountResponse struct {
	ID     int `json:"id"`
	Server struct {
		Language string `json:"language"`
		Number   int    `json:"number"`
	} `json:"server"`
	AccountGroup string `json:"accountGroup"`
	Error        string `json:"error"`
	BearerToken  string `json:"bearerToken"` // Added by us; not part of ogame response
}

func (r AddAccountResponse) GetBearerToken() string { return r.BearerToken }

func AddAccount(ctx context.Context, device Device, platform Platform, lobby, accountGroup, sessionToken string) (*AddAccountResponse, error) {
	blackbox, err := device.GetBlackbox()
	if err != nil {
		return nil, err
	}
	var payload struct {
		AccountGroup string `json:"accountGroup"`
		Blackbox     string `json:"blackbox"`
		Locale       string `json:"locale"`
		Kid          string `json:"kid"`
	}
	payload.AccountGroup = accountGroup // en_181
	payload.Blackbox = blackboxPrefix + blackbox
	payload.Locale = "en_GB"
	jsonPayloadBytes, err := json.Marshal(&payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPut, getGameforgeLobbyBaseURL(lobby, platform)+"/api/users/me/accounts", strings.NewReader(string(jsonPayloadBytes)))
	if err != nil {
		return nil, err
	}
	req.Header.Set(authorizationHeaderKey, buildBearerHeaderValue(sessionToken))
	req.Header.Set(contentTypeHeaderKey, applicationJson)
	req.Header.Set(acceptEncodingHeaderKey, gzipEncoding)
	req.WithContext(ctx)
	resp, err := device.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	by, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusBadRequest { // Same status is returned when IP is temporarily blocked
		return nil, errors.New("invalid request, account already in lobby ?")
	}
	var newAccount AddAccountResponse
	if err := json.Unmarshal(by, &newAccount); err != nil {
		return nil, errors.New(err.Error() + " : " + string(by))
	}
	if newAccount.Error != "" {
		return nil, errors.New(newAccount.Error)
	}
	newAccount.BearerToken = sessionToken
	return &newAccount, nil
}

type LoginResponse struct {
	Token                     string `json:"token"`
	IsPlatformLogin           bool   `json:"isPlatformLogin"`
	IsGameAccountMigrated     bool   `json:"isGameAccountMigrated"`
	PlatformUserID            string `json:"platformUserId"`
	IsGameAccountCreated      bool   `json:"isGameAccountCreated"`
	HasUnmigratedGameAccounts bool   `json:"hasUnmigratedGameAccounts"`
}

func (r LoginResponse) GetBearerToken() string { return r.Token }

func extractChallengeID(resp *http.Response) (challengeID string) {
	gfChallengeID := resp.Header.Get(ChallengeIDCookieName)
	if gfChallengeID != "" {
		// c434aa65-a064-498f-9ca4-98054bab0db8;https://challenge.gameforge.com
		parts := strings.Split(gfChallengeID, ";")
		challengeID = parts[0]
	}
	return
}

func login(params *gfLoginParams) (out *LoginResponse, err error) {
	setDefaultParams(params)
	if params.Device == nil {
		return out, errors.New("device is nil")
	}
	client := params.Device
	ctx := params.Ctx
	gameEnvironmentID, platformGameID, err := getConfiguration(ctx, client, params.platform, params.lobby)
	if err != nil {
		return out, err
	}

	req, err := postSessionsReq(params, gameEnvironmentID, platformGameID)
	if err != nil {
		return out, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	by, err := io.ReadAll(resp.Body)
	if err != nil {
		return out, err
	}

	if resp.StatusCode == http.StatusConflict {
		if challengeID := extractChallengeID(resp); challengeID != "" {
			return out, NewCaptchaRequiredError(challengeID)
		}
	}

	if resp.StatusCode == http.StatusForbidden {
		return out, errors.New(resp.Status + " : " + string(by))
	} else if resp.StatusCode >= http.StatusInternalServerError {
		return out, errors.New("gameforge server error code : " + resp.Status)
	} else if resp.StatusCode != http.StatusCreated {
		if string(by) == `{"reason":"OTP_REQUIRED"}` {
			return out, ErrOTPRequired
		}
		if string(by) == `{"reason":"OTP_INVALID"}` {
			return out, ErrOTPInvalid
		}
		return out, ErrBadCredentials
	}

	if err := json.Unmarshal(by, &out); err != nil {
		return out, err
	}
	return out, nil
}

func getConfiguration(ctx context.Context, client HttpClient, platform Platform, lobby string) (string, string, error) {
	ogURL := getGameforgeLobbyBaseURL(lobby, platform) + "/config/configuration.js"
	req, err := http.NewRequest(http.MethodGet, ogURL, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set(acceptEncodingHeaderKey, gzipEncoding)
	req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	by, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	gameEnvironmentIDRgx := regexp.MustCompile(`"gameEnvironmentId":"([^"]+)"`)
	m := gameEnvironmentIDRgx.FindSubmatch(by)
	if len(m) != 2 {
		return "", "", errors.New("failed to get gameEnvironmentId")
	}
	gameEnvironmentID := m[1]

	platformGameIDRgx := regexp.MustCompile(`"platformGameId":"([^"]+)"`)
	m = platformGameIDRgx.FindSubmatch(by)
	if len(m) != 2 {
		return "", "", errors.New("failed to get platformGameId")
	}
	platformGameID := m[1]

	return string(gameEnvironmentID), string(platformGameID), nil
}

func postSessionsReq(params *gfLoginParams, gameEnvironmentID, platformGameID string) (*http.Request, error) {
	dev := params.Device
	ctx := params.Ctx
	username := params.Username
	password := params.Password
	otpSecret := params.OtpSecret
	challengeID := params.ChallengeID

	blackbox, err := dev.GetBlackbox()
	if err != nil {
		return nil, err
	}

	var payload = struct {
		Identity                string `json:"identity"`
		Password                string `json:"password"`
		Locale                  string `json:"locale"`
		GfLang                  string `json:"gfLang"`
		PlatformGameID          string `json:"platformGameId"`
		Blackbox                string `json:"blackbox"`
		GameEnvironmentID       string `json:"gameEnvironmentId"`
		AutoGameAccountCreation bool   `json:"autoGameAccountCreation"`
	}{
		Identity:                username,
		Password:                password,
		Locale:                  "en_GB",
		GfLang:                  "en",
		PlatformGameID:          platformGameID,
		Blackbox:                blackboxPrefix + blackbox,
		GameEnvironmentID:       gameEnvironmentID,
		AutoGameAccountCreation: false,
	}
	by, err := json.Marshal(&payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, "https://gameforge.com/api/v1/auth/thin/sessions", bytes.NewReader(by))
	if err != nil {
		return nil, err
	}

	if challengeID != "" {
		req.Header.Set(ChallengeIDCookieName, challengeID)
	}

	if otpSecret != "" {
		passcode, err := totp.GenerateCodeCustom(otpSecret, time.Now(), totp.ValidateOpts{
			Period:    30,
			Skew:      1,
			Digits:    otp.DigitsSix,
			Algorithm: otp.AlgorithmSHA1,
		})
		if err != nil {
			return nil, err
		}
		req.Header.Set(twoFactorHeaderKey, passcode)
		req.Header.Set(installationIDHeaderKey, "")
	}

	req.Header.Set(contentTypeHeaderKey, applicationJson)
	req.Header.Set(acceptEncodingHeaderKey, gzipEncoding)
	req.WithContext(ctx)
	return req, nil
}

func StartCaptchaChallenge(ctx context.Context, client HttpClient, challengeID string) (questionRaw, iconsRaw []byte, err error) {
	doReq := func(u string) ([]byte, error) {
		req, err := http.NewRequest(http.MethodGet, u, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set(acceptEncodingHeaderKey, gzipEncoding)
		req.WithContext(ctx)
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		raw, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return raw, nil
	}
	challengeURL := getChallengeURL(challengeBaseURL, challengeID)
	imgDropURL := getChallengeURL(imgDropChallengeBaseURL, challengeID) + "/" + endpointLoc
	if _, err = doReq(challengeURL); err != nil {
		return
	}
	if _, err = doReq(imgDropURL); err != nil {
		return
	}
	if questionRaw, err = doReq(imgDropURL + "/text"); err != nil {
		return
	}
	if iconsRaw, err = doReq(imgDropURL + "/drag-icons"); err != nil {
		return
	}
	return
}

func SolveChallenge(ctx context.Context, client HttpClient, challengeID string, answer int64) error {
	challengeURL := getChallengeURL(imgDropChallengeBaseURL, challengeID) + "/" + endpointLoc
	body := strings.NewReader(fmt.Sprintf(`{"answer":%d}`, answer))
	req, _ := http.NewRequest(http.MethodPost, challengeURL, body)
	req.Header.Set(contentTypeHeaderKey, applicationJson)
	req.Header.Set(acceptEncodingHeaderKey, gzipEncoding)
	req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to solve captcha (%s)", resp.Status)
	}
	return nil
}

// Server gameforge information for their servers
type Server struct {
	Language      string
	Number        int64
	AccountGroup  string
	Name          string
	PlayerCount   int64
	PlayersOnline int64
	Opened        string
	StartDate     string
	EndDate       *string
	ServerClosed  int64
	Prefered      int64
	SignupClosed  int64
	MultiLanguage int64
	AvailableOn   []string
	Settings      any
}

func (s Server) OGameSettings() OGameServerSettings {
	return convertToStruct[OGameServerSettings](s.Settings)
}

func (s Server) IkariamSettings() IkariamServerSettings {
	return convertToStruct[IkariamServerSettings](s.Settings)
}

func convertToStruct[T any](v any) (out T) {
	by, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(by, &out); err != nil {
		panic(err)
	}
	return out
}

// OGameServerSettings ...
type OGameServerSettings struct {
	AKS                      int64
	FleetSpeedWar            int64
	FleetSpeedHolding        int64
	FleetSpeedPeaceful       int64
	WreckField               int64
	ServerLabel              string
	EconomySpeed             any // can be 8 or "x8"
	PlanetFields             int64
	UniverseSize             int64 // Nb of galaxies
	ServerCategory           string
	EspionageProbeRaids      int64
	PremiumValidationGift    int64
	DebrisFieldFactorShips   int64
	ResearchDurationDivisor  float64
	DebrisFieldFactorDefence int64
}

// IkariamServerSettings ...
type IkariamServerSettings struct {
	MaxCities                  int64   `json:"maxCities"`
	FleetSpeed                 int64   `json:"fleetSpeed"`
	ServerType                 string  `json:"serverType"`
	ServerLabel                string  `json:"serverLabel"`
	EconomySpeed               int64   `json:"economySpeed"`
	ArmyCostFactor             float64 `json:"armyCostFactor"`
	ServerCategory             string  `json:"serverCategory"`
	ArmySpeedFactor            float64 `json:"armySpeedFactor"`
	ResearchCostFactor         float64 `json:"researchCostFactor"`
	CombatWithoutMorale        bool    `json:"combatWithoutMorale"`
	WineProductionFactor       float64 `json:"wineProductionFactor"`
	GoldPlunderingAllowed      bool    `json:"goldPlunderingAllowed"`
	PremiumValidationGift      int64   `json:"premiumValidationGift"`
	ArmyConstructionFactor     float64 `json:"armyConstructionFactor"`
	MarbleProductionFactor     float64 `json:"marbleProductionFactor"`
	SatisfactionWineFactor     float64 `json:"satisfactionWineFactor"`
	SulfurProductionFactor     float64 `json:"sulfurProductionFactor"`
	TransporterSpeedFactor     float64 `json:"transporterSpeedFactor"`
	CrystalProductionFactor    float64 `json:"crystalProductionFactor"`
	FleetConstructionFactor    float64 `json:"fleetConstructionFactor"`
	GoldSafeCapacityPerLevel   int64   `json:"goldSafeCapacityPerLevel"`
	ResearchProductionFactor   float64 `json:"researchProductionFactor"`
	ResourceProductionFactor   float64 `json:"resourceProductionFactor"`
	BuildingConstructionFactor float64 `json:"buildingConstructionFactor"`
	ConversionProductionFactor float64 `json:"conversionProductionFactor"`
}

func (s OGameServerSettings) ProbeRaidsEnabled() bool {
	return s.EspionageProbeRaids == 1
}

func GetServers(ctx context.Context, client HttpClient, platform Platform, lobby string) ([]Server, error) {
	var servers []Server
	req, err := http.NewRequest(http.MethodGet, getGameforgeLobbyBaseURL(lobby, platform)+"/api/servers", nil)
	if err != nil {
		return servers, err
	}
	req.Header.Set(acceptEncodingHeaderKey, gzipEncoding)
	req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return servers, err
	}
	defer resp.Body.Close()
	by, err := io.ReadAll(resp.Body)
	if err != nil {
		return servers, err
	}
	if err := json.Unmarshal(by, &servers); err != nil {
		return servers, errors.New("failed to get servers : " + err.Error() + " : " + string(by))
	}
	return servers, nil
}

type Account struct {
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
		Value any // Can be string or int
	}
	Sitting struct {
		Shared       bool
		EndTime      *string
		CooldownTime *string
	}
}

func GetServerAccount(ctx context.Context, client HttpClient, platform Platform, lobby, bearerToken, serverName, lang string, playerID int64) (account Account, server Server, err error) {
	accounts, err := GetUserAccounts(ctx, client, platform, lobby, bearerToken)
	if err != nil {
		return
	}
	servers, err := GetServers(ctx, client, platform, lobby)
	if err != nil {
		return
	}
	account, server, err = FindAccount(serverName, lang, playerID, accounts, servers)
	if err != nil {
		return
	}
	return
}

func GetUserAccounts(ctx context.Context, client HttpClient, platform Platform, lobby, bearerToken string) ([]Account, error) {
	var userAccounts []Account
	req, err := http.NewRequest(http.MethodGet, getGameforgeLobbyBaseURL(lobby, platform)+"/api/users/me/accounts", nil)
	if err != nil {
		return userAccounts, err
	}
	req.Header.Set(authorizationHeaderKey, buildBearerHeaderValue(bearerToken))
	req.Header.Set(acceptEncodingHeaderKey, gzipEncoding)
	req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return userAccounts, err
	}
	defer resp.Body.Close()
	by, err := io.ReadAll(resp.Body)
	if err != nil {
		return userAccounts, err
	}
	if err := json.Unmarshal(by, &userAccounts); err != nil {
		return userAccounts, errors.New("failed to get user accounts : " + err.Error() + " : " + string(by))
	}
	return userAccounts, nil
}

func GetLoginLink(ctx context.Context, device Device, platform Platform, lobby string, userAccount Account, bearerToken string) (string, error) {
	ogURL := getGameforgeLobbyBaseURL(lobby, platform) + "/api/users/me/loginLink"

	blackbox, err := device.GetBlackbox()
	if err != nil {
		return "", err
	}

	var payload = struct {
		Blackbox      string `json:"blackbox"`
		Id            int64  `json:"id"`
		ClickedButton string `json:"clickedButton"`
		Server        struct {
			Language string `json:"language"`
			Number   int64  `json:"number"`
		} `json:"server"`
	}{
		Blackbox:      blackboxPrefix + blackbox,
		Id:            userAccount.ID,
		ClickedButton: "account_list",
	}

	payload.Server.Language = userAccount.Server.Language
	payload.Server.Number = userAccount.Server.Number

	by, err := json.Marshal(&payload)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, ogURL, bytes.NewReader(by))
	if err != nil {
		return "", err
	}
	req.Header.Set(contentTypeHeaderKey, applicationJson)
	req.Header.Set(authorizationHeaderKey, buildBearerHeaderValue(bearerToken))
	req.Header.Set(acceptEncodingHeaderKey, gzipEncoding)
	req.WithContext(ctx)
	resp, err := device.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	by2, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == http.StatusBadRequest && string(by2) == `[]` {
		return "", ErrLoginLink
	}

	var loginLink struct{ URL string }
	if err := json.Unmarshal(by2, &loginLink); err != nil {
		return "", errors.New("failed to get login link : " + err.Error() + " : " + string(by2))
	}
	return loginLink.URL, nil
}

// ExecLoginLink ...
func ExecLoginLink(ctx context.Context, client HttpClient, loginLink string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, loginLink, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add(acceptEncodingHeaderKey, gzipEncoding)
	req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// FindServer ...
func FindServer(serverName, lang string, servers []Server) (out Server, found bool) {
	for _, s := range servers {
		if s.Name == serverName && s.Language == lang {
			return s, true
		}
	}
	return
}

// FindAccount ...
func FindAccount(serverName, lang string, playerID int64, accounts []Account, servers []Server) (Account, Server, error) {
	if lang == "ba" {
		lang = "yu"
	}
	var account Account
	server, found := FindServer(serverName, lang, servers)
	if !found {
		return Account{}, Server{}, fmt.Errorf("server %s, %s not found", serverName, lang)
	}
	for _, a := range accounts {
		if a.Server.Language == server.Language && a.Server.Number == server.Number {
			if playerID == 0 || a.ID == playerID {
				account = a
				break
			}
		}
	}
	if account.ID == 0 {
		return Account{}, Server{}, ErrAccountNotFound
	}
	return account, server, nil
}
