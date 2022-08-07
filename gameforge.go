package ogame

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type CaptchaRequiredError struct {
	ChallengeID string
}

func NewCaptchaRequiredError(challengeID string) *CaptchaRequiredError {
	return &CaptchaRequiredError{ChallengeID: challengeID}
}

func (e CaptchaRequiredError) Error() string {
	return fmt.Sprintf("captcha required, %s", e.ChallengeID)
}

// Register a new gameforge lobby account
func Register(client *http.Client, ctx context.Context, lobby, email, password, challengeID, lang string) error {
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
		req.Header.Add(gfChallengeIDCookieName, challengeID)
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
		gfChallengeID := resp.Header.Get(gfChallengeIDCookieName) // c434aa65-a064-498f-9ca4-98054bab0db8;https://challenge.gameforge.com
		if gfChallengeID != "" {
			parts := strings.Split(gfChallengeID, ";")
			challengeID := parts[0]
			return NewCaptchaRequiredError(challengeID)
		}
	}
	by, err := readBody(resp)
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
	if res.Error != "" {
		return errors.New(res.Error)
	}
	return nil
}

// ValidateAccount validate a gameforge account
func ValidateAccount(client IHttpClient, ctx context.Context, lobby, code string) error {
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
func RedeemCode(client *http.Client, ctx context.Context, lobby, email, password, otpSecret, token string) error {
	postSessionsRes, err := GFLogin(client, ctx, lobby, email, password, otpSecret, "")
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
	req, err := http.NewRequest("POST", "https://"+lobby+".ogame.gameforge.com/api/token", strings.NewReader(string(jsonPayloadBytes)))
	if err != nil {
		return err
	}
	req.Header.Add("authorization", "Bearer "+postSessionsRes.Token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// {"tokenType":"accountTrading"}
	type respStruct struct {
		TokenType string `json:"tokenType"`
	}
	var respParsed respStruct
	by, err := readBody(resp)
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
func LoginAndAddAccount(client *http.Client, ctx context.Context, lobby, username, password, otpSecret, universe, lang string) (AddAccountRes, error) {
	var newAccount AddAccountRes
	postSessionsRes, err := GFLogin(client, ctx, lobby, username, password, otpSecret, "")
	if err != nil {
		return newAccount, err
	}
	servers, err := GetServers(lobby, client, ctx)
	if err != nil {
		return newAccount, err
	}
	server, found := findServer(universe, lang, servers)
	if !found {
		return newAccount, errors.New("server not found")
	}
	return AddAccount(client, ctx, lobby, server.AccountGroup, postSessionsRes.Token)
}

// AddAccountRes response from creating a new account
type AddAccountRes struct {
	ID     int
	Server struct {
		Language string
		Number   int
	}
	BearerToken string
	Error       string
}

func AddAccount(client IHttpClient, ctx context.Context, lobby, accountGroup, sessionToken string) (AddAccountRes, error) {
	var newAccount AddAccountRes
	var payload struct {
		AccountGroup string `json:"accountGroup"`
		Locale       string `json:"locale"`
		Kid          string `json:"kid"`
	}
	payload.AccountGroup = accountGroup // en_181
	payload.Locale = "en_GB"
	jsonPayloadBytes, err := json.Marshal(&payload)
	if err != nil {
		return newAccount, err
	}
	req, err := http.NewRequest(http.MethodPut, "https://"+lobby+".ogame.gameforge.com/api/users/me/accounts", strings.NewReader(string(jsonPayloadBytes)))
	if err != nil {
		return newAccount, err
	}
	newAccount.BearerToken = sessionToken
	req.Header.Add("authorization", "Bearer "+sessionToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return newAccount, err
	}
	defer resp.Body.Close()
	by, err := readBody(resp)
	if err != nil {
		return newAccount, err
	}
	if resp.StatusCode == http.StatusBadRequest {
		return newAccount, errors.New("invalid request, account already in lobby ?")
	}
	if err := json.Unmarshal(by, &newAccount); err != nil {
		return newAccount, errors.New(err.Error() + " : " + string(by))
	}
	if newAccount.Error != "" {
		return newAccount, errors.New(newAccount.Error)
	}
	return newAccount, nil
}

type GFLoginRes struct {
	Token                     string `json:"token"`
	IsPlatformLogin           bool   `json:"isPlatformLogin"`
	IsGameAccountMigrated     bool   `json:"isGameAccountMigrated"`
	PlatformUserID            string `json:"platformUserId"`
	IsGameAccountCreated      bool   `json:"isGameAccountCreated"`
	HasUnmigratedGameAccounts bool   `json:"hasUnmigratedGameAccounts"`
}

func GFLogin(client IHttpClient, ctx context.Context, lobby, username, password, otpSecret, challengeID string) (out *GFLoginRes, err error) {
	gameEnvironmentID, platformGameID, err := getConfiguration(client, ctx, lobby)
	if err != nil {
		return out, err
	}

	req, err := postSessionsReq(gameEnvironmentID, platformGameID, username, password, otpSecret, challengeID)
	if err != nil {
		return out, err
	}

	req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	by, err := readBody(resp)
	if err != nil {
		return out, err
	}

	if resp.StatusCode == http.StatusConflict {
		gfChallengeID := resp.Header.Get(gfChallengeIDCookieName)
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

func getConfiguration(client IHttpClient, ctx context.Context, lobby string) (string, string, error) {
	ogURL := "https://" + lobby + ".ogame.gameforge.com/config/configuration.js"
	req, err := http.NewRequest("GET", ogURL, nil)
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
	by, err := readBody(resp)
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

func postSessionsReq(gameEnvironmentID, platformGameID, username, password, otpSecret, challengeID string) (*http.Request, error) {
	payload := url.Values{
		"autoGameAccountCreation": {"false"},
		"gameEnvironmentId":       {gameEnvironmentID},
		"platformGameId":          {platformGameID},
		"gfLang":                  {"en"},
		"locale":                  {"en_GB"},
		"identity":                {username},
		"password":                {password},
	}
	req, err := http.NewRequest("POST", "https://gameforge.com/api/v1/auth/thin/sessions", strings.NewReader(payload.Encode()))
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

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	return req, nil
}

func StartCaptchaChallenge(client IHttpClient, ctx context.Context, challengeID string) (questionRaw, iconsRaw []byte, err error) {
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

func SolveChallenge(client IHttpClient, ctx context.Context, challengeID string, answer int64) error {
	challengeURL := "https://image-drop-challenge.gameforge.com/challenge/" + challengeID + "/en-GB"
	req, _ := http.NewRequest(http.MethodPost, challengeURL, strings.NewReader(`{"answer":`+strconv.FormatInt(answer, 10)+`}`))
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
		EconomySpeed             interface{} // can be 8 or "x8"
		PlanetFields             int64
		UniverseSize             int64 // Nb of galaxies
		ServerCategory           string
		EspionageProbeRaids      int64
		PremiumValidationGift    int64
		DebrisFieldFactorShips   int64
		DebrisFieldFactorDefence int64
	}
}

func GetServers(lobby string, client IHttpClient, ctx context.Context) ([]Server, error) {
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
	by, err := readBody(resp)
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
	SpeedFleet                    int64   `xml:"speedFleet"`                    // 6 // Deprecated in 8.1.0
	SpeedFleetPeaceful            int64   `xml:"speedFleetPeaceful"`            // 1
	SpeedFleetWar                 int64   `xml:"speedFleetWar"`                 // 1
	SpeedFleetHolding             int64   `xml:"speedFleetHolding"`             // 1
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
	//TopScore                      int64   `xml:"topScore"`                      // 60259362 / 1.0363090034999E+17
}

// GetServerData gets the server data from xml api
func GetServerData(client IHttpClient, ctx context.Context, serverNumber int64, serverLang string) (ServerData, error) {
	var serverData ServerData
	req, err := http.NewRequest(http.MethodGet, "https://s"+strconv.FormatInt(serverNumber, 10)+"-"+serverLang+".ogame.gameforge.com/api/serverData.xml", nil)
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
	by, err := readBody(resp)
	if err != nil {
		return serverData, err
	}
	if err := xml.Unmarshal(by, &serverData); err != nil {
		return serverData, err
	}
	return serverData, nil
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

func GetUserAccounts(client IHttpClient, ctx context.Context, lobby, token string) ([]account, error) {
	var userAccounts []account
	req, err := http.NewRequest(http.MethodGet, "https://"+lobby+".ogame.gameforge.com/api/users/me/accounts", nil)
	if err != nil {
		return userAccounts, err
	}
	req.Header.Add("authorization", "Bearer "+token)
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return userAccounts, err
	}
	defer resp.Body.Close()
	by, err := readBody(resp)
	if err != nil {
		return userAccounts, err
	}
	if err := json.Unmarshal(by, &userAccounts); err != nil {
		return userAccounts, errors.New("failed to get user accounts : " + err.Error() + " : " + string(by))
	}
	return userAccounts, nil
}

func GetLoginLink(client IHttpClient, ctx context.Context, lobby string, userAccount account, token string) (string, error) {
	ogURL := fmt.Sprintf("https://%s.ogame.gameforge.com/api/users/me/loginLink?id=%d&server[language]=%s&server[number]=%d&clickedButton=account_list",
		lobby, userAccount.ID, userAccount.Server.Language, userAccount.Server.Number)
	req, err := http.NewRequest(http.MethodGet, ogURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("authorization", "Bearer "+token)
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	by, err := readBody(resp)
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
