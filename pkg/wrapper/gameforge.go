package wrapper

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/alaingilbert/ogame/pkg/device"
	"github.com/alaingilbert/ogame/pkg/httpclient"
	"github.com/alaingilbert/ogame/pkg/ogame"
	"github.com/alaingilbert/ogame/pkg/utils"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// TokenCookieName ogame cookie name for token id
const TokenCookieName = "gf-token-production"
const ChallengeIDCookieName = "gf-challenge-id"

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

// Register a new gameforge lobby account
func Register(client httpclient.IHttpClient, ctx context.Context, lobby, email, password, challengeID, lang string) error {
	if lang == "" {
		lang = "en"
	}
	var payload struct {
		Credentials struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		} `json:"credentials"`
		Language string `json:"language"`
		Kid      string `json:"kid"`
	}
	payload.Credentials.Email = email
	payload.Credentials.Password = password
	payload.Language = lang
	jsonPayloadBytes, err := json.Marshal(&payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, "https://"+lobby+".ogame.gameforge.com/api/users", strings.NewReader(string(jsonPayloadBytes)))
	if err != nil {
		return err
	}
	if challengeID != "" {
		req.Header.Add(ChallengeIDCookieName, challengeID)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusConflict {
		gfChallengeID := resp.Header.Get(ChallengeIDCookieName) // c434aa65-a064-498f-9ca4-98054bab0db8;https://challenge.gameforge.com
		if gfChallengeID != "" {
			parts := strings.Split(gfChallengeID, ";")
			challengeID := parts[0]
			return NewCaptchaRequiredError(challengeID)
		}
	}
	by, err := utils.ReadBody(resp)
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

// ValidateAccount validate a gameforge account
func ValidateAccount(client httpclient.IHttpClient, ctx context.Context, lobby, code string) error {
	if len(code) != 36 {
		return errors.New("invalid validation code")
	}
	req, err := http.NewRequest(http.MethodPut, "https://"+lobby+".ogame.gameforge.com/api/users/validate/"+code, strings.NewReader(`{"language":"en"}`))
	if err != nil {
		return err
	}
	req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// RedeemCode ...
func RedeemCode(device *device.Device, ctx context.Context, lobby, email, password, otpSecret, token string) error {
	postSessionsRes, err := GFLogin(device, ctx, lobby, email, password, otpSecret, "")
	if err != nil {
		return err
	}
	var payload struct {
		Token string `json:"token"`
	}
	payload.Token = token
	jsonPayloadBytes, err := json.Marshal(&payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, "https://"+lobby+".ogame.gameforge.com/api/token", strings.NewReader(string(jsonPayloadBytes)))
	if err != nil {
		return err
	}
	req.Header.Add("authorization", "Bearer "+postSessionsRes.Token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.WithContext(ctx)
	resp, err := device.GetClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// {"tokenType":"accountTrading"}
	type respStruct struct {
		TokenType string `json:"tokenType"`
	}
	var respParsed respStruct
	by, err := utils.ReadBody(resp)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusBadRequest {
		return errors.New("invalid request, token invalid ?")
	}
	if err := json.Unmarshal(by, &respParsed); err != nil {
		return errors.New(err.Error() + " : " + string(by))
	}
	if respParsed.TokenType != "accountTrading" {
		return errors.New("tokenType is not accountTrading")
	}
	return nil
}

// LoginAndAddAccount adds an account to a gameforge lobby
func LoginAndAddAccount(device *device.Device, ctx context.Context, lobby, username, password, otpSecret, universe, lang string) (*AddAccountRes, error) {
	postSessionsRes, err := GFLogin(device, ctx, lobby, username, password, otpSecret, "")
	if err != nil {
		return nil, err
	}
	servers, err := GetServers(lobby, device.GetClient(), ctx)
	if err != nil {
		return nil, err
	}
	server, found := findServer(universe, lang, servers)
	if !found {
		return nil, errors.New("server not found")
	}
	return AddAccount(device.GetClient(), ctx, lobby, server.AccountGroup, postSessionsRes.Token)
}

// AddAccountRes response from creating a new account
type AddAccountRes struct {
	ID     int `json:"id"`
	Server struct {
		Language string `json:"language"`
		Number   int    `json:"number"`
	} `json:"server"`
	AccountGroup string `json:"accountGroup"`
	Error        string `json:"error"`
	BearerToken  string `json:"bearerToken"` // Added by us; not part of ogame response
}

func (r AddAccountRes) GetBearerToken() string { return r.BearerToken }

func AddAccount(client httpclient.IHttpClient, ctx context.Context, lobby, accountGroup, sessionToken string) (*AddAccountRes, error) {
	var payload struct {
		AccountGroup string `json:"accountGroup"`
		Locale       string `json:"locale"`
		Kid          string `json:"kid"`
	}
	payload.AccountGroup = accountGroup // en_181
	payload.Locale = "en_GB"
	jsonPayloadBytes, err := json.Marshal(&payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPut, "https://"+lobby+".ogame.gameforge.com/api/users/me/accounts", strings.NewReader(string(jsonPayloadBytes)))
	if err != nil {
		return nil, err
	}
	req.Header.Add("authorization", "Bearer "+sessionToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	by, err := utils.ReadBody(resp)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusBadRequest {
		return nil, errors.New("invalid request, account already in lobby ?")
	}
	var newAccount AddAccountRes
	if err := json.Unmarshal(by, &newAccount); err != nil {
		return nil, errors.New(err.Error() + " : " + string(by))
	}
	if newAccount.Error != "" {
		return nil, errors.New(newAccount.Error)
	}
	newAccount.BearerToken = sessionToken
	return &newAccount, nil
}

type GFLoginRes struct {
	Token                     string `json:"token"`
	IsPlatformLogin           bool   `json:"isPlatformLogin"`
	IsGameAccountMigrated     bool   `json:"isGameAccountMigrated"`
	PlatformUserID            string `json:"platformUserId"`
	IsGameAccountCreated      bool   `json:"isGameAccountCreated"`
	HasUnmigratedGameAccounts bool   `json:"hasUnmigratedGameAccounts"`
}

func (r GFLoginRes) GetBearerToken() string { return r.Token }

func GFLogin(dev *device.Device, ctx context.Context, lobby, username, password, otpSecret, challengeID string) (out *GFLoginRes, err error) {
	gameEnvironmentID, platformGameID, err := getConfiguration(dev.GetClient(), ctx, lobby)
	if err != nil {
		return out, err
	}

	blackbox, err := dev.GetBlackbox()
	if err != nil {
		return out, err
	}

	req, err := postSessionsReq(gameEnvironmentID, platformGameID, username, password, otpSecret, challengeID, blackbox)
	if err != nil {
		return out, err
	}

	req.WithContext(ctx)

	resp, err := dev.GetClient().Do(req)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	by, err := utils.ReadBody(resp)
	if err != nil {
		return out, err
	}

	if resp.StatusCode == http.StatusConflict {
		gfChallengeID := resp.Header.Get(ChallengeIDCookieName)
		if gfChallengeID != "" {
			parts := strings.Split(gfChallengeID, ";")
			challengeID := parts[0]
			return out, NewCaptchaRequiredError(challengeID)
		}
	}

	if resp.StatusCode == http.StatusForbidden {
		return out, errors.New(resp.Status + " : " + string(by))
	}

	if resp.StatusCode >= http.StatusInternalServerError {
		return out, errors.New("OGame server error code : " + resp.Status)
	}

	if resp.StatusCode != http.StatusCreated {
		if string(by) == `{"reason":"OTP_REQUIRED"}` {
			return out, ogame.ErrOTPRequired
		}
		if string(by) == `{"reason":"OTP_INVALID"}` {
			return out, ogame.ErrOTPInvalid
		}
		return out, ogame.ErrBadCredentials
	}

	if err := json.Unmarshal(by, &out); err != nil {
		return out, err
	}
	return out, nil
}

func getConfiguration(client httpclient.IHttpClient, ctx context.Context, lobby string) (string, string, error) {
	ogURL := "https://" + lobby + ".ogame.gameforge.com/config/configuration.js"
	req, err := http.NewRequest(http.MethodGet, ogURL, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	by, err := utils.ReadBody(resp)
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

func postSessionsReq(gameEnvironmentID, platformGameID, username, password, otpSecret, challengeID, blackbox string) (*http.Request, error) {
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
		Blackbox:                "tra:" + blackbox,
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
		req.Header.Set("gf-challenge-id", challengeID)
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
		req.Header.Add("tnt-2fa-code", passcode)
		req.Header.Add("tnt-installation-id", "")
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	return req, nil
}

func StartCaptchaChallenge(client httpclient.IHttpClient, ctx context.Context, challengeID string) (questionRaw, iconsRaw []byte, err error) {
	req, err := http.NewRequest(http.MethodGet, "https://challenge.gameforge.com/challenge/"+challengeID, nil)
	if err != nil {
		return
	}
	req.WithContext(ctx)
	challengeResp, err := client.Do(req)
	if err != nil {
		return
	}
	defer challengeResp.Body.Close()
	_, _ = ioutil.ReadAll(challengeResp.Body)

	req, err = http.NewRequest(http.MethodGet, "https://image-drop-challenge.gameforge.com/challenge/"+challengeID+"/en-GB", nil)
	if err != nil {
		return
	}
	req.WithContext(ctx)
	challengePresentedResp, err := client.Do(req)
	if err != nil {
		return
	}
	defer challengePresentedResp.Body.Close()
	_, _ = ioutil.ReadAll(challengePresentedResp.Body)

	// Question request
	req, err = http.NewRequest(http.MethodGet, "https://image-drop-challenge.gameforge.com/challenge/"+challengeID+"/en-GB/text", nil)
	if err != nil {
		return
	}
	req.WithContext(ctx)
	questionResp, err := client.Do(req)
	if err != nil {
		return
	}
	defer questionResp.Body.Close()
	questionRaw, _ = ioutil.ReadAll(questionResp.Body)

	// Icons request
	req, err = http.NewRequest(http.MethodGet, "https://image-drop-challenge.gameforge.com/challenge/"+challengeID+"/en-GB/drag-icons", nil)
	if err != nil {
		return
	}
	req.WithContext(ctx)
	iconsResp, err := client.Do(req)
	if err != nil {
		return
	}
	defer iconsResp.Body.Close()
	iconsRaw, _ = ioutil.ReadAll(iconsResp.Body)
	return
}

func SolveChallenge(client httpclient.IHttpClient, ctx context.Context, challengeID string, answer int64) error {
	challengeURL := "https://image-drop-challenge.gameforge.com/challenge/" + challengeID + "/en-GB"
	body := strings.NewReader(`{"answer":` + utils.FI64(answer) + `}`)
	req, _ := http.NewRequest(http.MethodPost, challengeURL, body)
	req.Header.Set("Content-Type", "application/json")
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

// Server ogame information for their servers
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
	Settings      struct {
		AKS                      int64
		FleetSpeed               int64
		WreckField               int64
		ServerLabel              string
		EconomySpeed             any // can be 8 or "x8"
		PlanetFields             int64
		UniverseSize             int64 // Nb of galaxies
		ServerCategory           string
		EspionageProbeRaids      int64
		PremiumValidationGift    int64
		DebrisFieldFactorShips   int64
		DebrisFieldFactorDefence int64
	}
}

func GetServers(lobby string, client httpclient.IHttpClient, ctx context.Context) ([]Server, error) {
	var servers []Server
	req, err := http.NewRequest(http.MethodGet, "https://"+lobby+".ogame.gameforge.com/api/servers", nil)
	if err != nil {
		return servers, err
	}
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return servers, err
	}
	defer resp.Body.Close()
	by, err := utils.ReadBody(resp)
	if err != nil {
		return servers, err
	}
	if err := json.Unmarshal(by, &servers); err != nil {
		return servers, errors.New("failed to get servers : " + err.Error() + " : " + string(by))
	}
	return servers, nil
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
}

// GetServerData gets the server data from xml api
func GetServerData(client httpclient.IHttpClient, ctx context.Context, serverNumber int64, serverLang string) (ServerData, error) {
	var serverData ServerData
	serverDataURL := "https://s" + utils.FI64(serverNumber) + "-" + serverLang + ".ogame.gameforge.com/api/serverData.xml"
	req, err := http.NewRequest(http.MethodGet, serverDataURL, nil)
	if err != nil {
		return serverData, err
	}
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return serverData, err
	}
	defer resp.Body.Close()
	by, err := utils.ReadBody(resp)
	if err != nil {
		return serverData, err
	}
	if err := xml.Unmarshal(by, &serverData); err != nil {
		return serverData, fmt.Errorf("failed to xml unmarshal %s : %w", serverDataURL, err)
	}
	return serverData, nil
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

func GetUserAccounts(client httpclient.IHttpClient, ctx context.Context, lobby, bearerToken string) ([]Account, error) {
	var userAccounts []Account
	req, err := http.NewRequest(http.MethodGet, "https://"+lobby+".ogame.gameforge.com/api/users/me/accounts", nil)
	if err != nil {
		return userAccounts, err
	}
	req.Header.Add("authorization", "Bearer "+bearerToken)
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return userAccounts, err
	}
	defer resp.Body.Close()
	by, err := utils.ReadBody(resp)
	if err != nil {
		return userAccounts, err
	}
	if err := json.Unmarshal(by, &userAccounts); err != nil {
		return userAccounts, errors.New("failed to get user accounts : " + err.Error() + " : " + string(by))
	}
	return userAccounts, nil
}

func GetLoginLink(dev *device.Device, ctx context.Context, lobby string, userAccount Account, bearerToken string) (string, error) {
	ogURL := fmt.Sprintf("https://%s.ogame.gameforge.com/api/users/me/loginLink", lobby)
	payload := struct {
		Server struct {
			Language string `json:"language"`
			Number   int64  `json:"number"`
		} `json:"server"`
		ID            int64  `json:"id"`
		ClickedButton string `json:"clickedButton"`
		Blackbox      string `json:"blackbox"`
	}{}
	payload.Server.Language = userAccount.Server.Language
	payload.Server.Number = userAccount.Server.Number
	payload.ID = userAccount.ID
	payload.ClickedButton = "account_list"
	payload.Blackbox, _ = dev.GetBlackbox()
	jsonPayloadBytes, err := json.Marshal(&payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, ogURL, strings.NewReader(string(jsonPayloadBytes)))
	if err != nil {
		return "", err
	}
	req.Header.Add("authorization", "Bearer "+bearerToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")

	resp, err := dev.GetClient().Do(req.WithContext(ctx))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	by, err := utils.ReadBody(resp)
	if err != nil {
		return "", err
	}
	var loginLink struct {
		URL string
	}
	if err := json.Unmarshal(by, &loginLink); err != nil {
		return "", errors.New("failed to get login link : " + err.Error() + " : " + string(by))
	}
	return loginLink.URL, nil
}
