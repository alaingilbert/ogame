package ogame

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"golang.org/x/net/html"

	"regexp"

	"strconv"

	"time"

	"log"

	"github.com/PuerkitoBio/goquery"
	"github.com/Sirupsen/logrus"
)

const defaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) " +
	"AppleWebKit/537.36 (KHTML, like Gecko) " +
	"Chrome/51.0.2704.103 Safari/537.36"

type ogameClient struct {
	http.Client
	UserAgent string
}

func (c *ogameClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", c.UserAgent)
	return c.Client.Do(req)
}

// OGame library
type OGame struct {
	Quiet        bool
	Universe     string
	Username     string
	password     string
	ogameSession string
	serverURL    string
	client       *ogameClient
}

// New creates a new instance of OGame wrapper.
func New(universe, username, password string) *OGame {
	b := new(OGame)
	b.Quiet = false

	b.Universe = universe
	b.Username = username
	b.password = password

	jar, _ := cookiejar.New(nil)
	b.client = &ogameClient{}
	b.client.Jar = jar
	b.client.UserAgent = defaultUserAgent

	if err := b.Login(); err != nil {
		log.Println(err)
	}

	return b
}

// SetUserAgent change the user-agent used by the http client
func (b *OGame) SetUserAgent(newUserAgent string) {
	b.client.UserAgent = newUserAgent
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
		UniverseSize             int
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

func findAccountByName(universe string, accounts []account, servers []server) (account, error) {
	for _, a := range accounts {
		for _, s := range servers {
			if universe == s.Name && a.Server.Language == s.Language && a.Server.Number == s.Number {
				return a, nil
			}
		}
	}
	return account{}, fmt.Errorf("server %s not found", universe)
}

type loginLinkResp struct {
	URL string
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
	var loginLink loginLinkResp
	if err := json.Unmarshal(by, &loginLink); err != nil {
		return "", err
	}
	return loginLink.URL, nil
}

// Login to ogame server
func (b *OGame) Login() error {
	phpSessionID, err := getPhpSessionID(b.client, b.Username, b.password)
	if err != nil {
		return err
	}
	accounts, err := getUserAccounts(b.client, phpSessionID)
	if err != nil {
		return err
	}
	servers, err := getServers(b.client)
	if err != nil {
		return err
	}
	userAccount, err := findAccountByName(b.Universe, accounts, servers)
	if err != nil {
		return err
	}
	loginLink, err := getLoginLink(b.client, userAccount, phpSessionID)
	if err != nil {
		return err
	}

	r := regexp.MustCompile(`https://(.+\.ogame\.gameforge\.com)/game`)
	res := r.FindStringSubmatch(loginLink)
	if len(res) != 2 {
		return errors.New("failed to get server url")
	}
	b.serverURL = res[1]

	req, err := http.NewRequest("GET", loginLink, nil)
	if err != nil {
		return err
	}
	resp, err := b.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}
	doc.Find("meta[name=ogame-session]").Each(func(i int, s *goquery.Selection) {
		b.ogameSession, _ = s.Attr("content")
	})

	if b.ogameSession == "" {
		return errors.New("bad credentials")
	}

	return nil
}

type eventboxResp struct {
	Hostile  int
	Neutral  int
	Friendly int
}

func (b *OGame) getPageContent(page string) string {
	req, _ := http.NewRequest("GET", "https://"+b.serverURL+"/game/index.php?page="+page, nil)
	resp, _ := b.client.Do(req)
	defer resp.Body.Close()
	by, _ := ioutil.ReadAll(resp.Body)
	return string(by)
}

func (b *OGame) getURL(page string) string {
	return "https://" + b.serverURL + "/game/index.php?page=" + page
}

func (b *OGame) fetchEventbox() (eventboxResp, error) {
	req, _ := http.NewRequest("GET", "https://"+b.serverURL+"/game/index.php?page=fetchEventbox", nil)
	resp, _ := b.client.Do(req)
	defer resp.Body.Close()
	by, _ := ioutil.ReadAll(resp.Body)
	var res eventboxResp
	json.Unmarshal(by, &res)
	return res, nil
}

// IsUnderAttack returns true if the user is under attack, false otherwise
func (b *OGame) IsUnderAttack() bool {
	eventbox, _ := b.fetchEventbox()
	return eventbox.Hostile > 0
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

func (b *OGame) fetchResources() (resourcesResp, error) {
	req, _ := http.NewRequest("GET", "https://"+b.serverURL+"/game/index.php?page=fetchResources", nil)
	resp, _ := b.client.Do(req)
	defer resp.Body.Close()
	by, _ := ioutil.ReadAll(resp.Body)
	var res resourcesResp
	if err := json.Unmarshal(by, &res); err != nil {
		return resourcesResp{}, err
	}
	return res, nil
}

// Resources ...
type Resources struct {
	Metal      int
	Crystal    int
	Deuterium  int
	Energy     int
	Darkmatter int
}

// GetResources gets user resources
func (b *OGame) GetResources() (Resources, error) {
	res, err := b.fetchResources()
	return Resources{
		Metal:      res.Metal.Resources.Actual,
		Crystal:    res.Crystal.Resources.Actual,
		Deuterium:  res.Deuterium.Resources.Actual,
		Energy:     res.Energy.Resources.Actual,
		Darkmatter: res.Darkmatter.Resources.Actual,
	}, err
}

// OGame constants
const (
	// Buildings
	MetalMine             = 1
	CrystalMine           = 2
	DeuteriumSynthesizer  = 3
	SolarPlant            = 4
	FusionReactor         = 12
	MetalStorage          = 22
	CrystalStorage        = 23
	DeuteriumTank         = 24
	ShieldedMetalDen      = 25
	UndergroundCrystalDen = 26
	SeabedDeuteriumDen    = 27

	// Facilities
	AllianceDepot   = 34
	RoboticsFactory = 14
	Shipyard        = 21
	ResearchLab     = 31
	MissileSilo     = 44
	NaniteFactory   = 15
	Terraformer     = 33
	SpaceDock       = 36

	// Defense
	RocketLauncher         = 401
	LightLaser             = 402
	HeavyLaser             = 403
	GaussCannon            = 404
	IonCannon              = 405
	PlasmaTurret           = 406
	SmallShieldDome        = 407
	LargeShieldDome        = 408
	AntiBallisticMissiles  = 502
	InterplanetaryMissiles = 503

	// Ships
	SmallCargo     = 202
	LargeCargo     = 203
	LightFighter   = 204
	HeavyFighter   = 205
	Cruiser        = 206
	Battleship     = 207
	ColonyShip     = 208
	Recycler       = 209
	EspionageProbe = 210
	Bomber         = 211
	SolarSatellite = 212
	Destroyer      = 213
	Deathstar      = 214
	Battlecruiser  = 215

	// Research
	EspionageTechnology          = 106
	ComputerTechnology           = 108
	WeaponsTechnology            = 109
	ShieldingTechnology          = 110
	ArmourTechnology             = 111
	EnergyTechnology             = 113
	HyperspaceTechnology         = 114
	CombustionDrive              = 115
	ImpulseDrive                 = 117
	HyperspaceDrive              = 118
	LaserTechnology              = 120
	IonTechnology                = 121
	PlasmaTechnology             = 122
	IntergalacticResearchNetwork = 123
	Astrophysics                 = 124
	GravitonTechnology           = 199

	// Missions
	Attack             = 1
	GroupedAttack      = 2
	Transport          = 3
	Park               = 4
	ParkInThatAlly     = 5
	Spy                = 6
	Colonize           = 7
	RecycleDebrisField = 8
	Destroy            = 9
	Expedition         = 15
)

// BuildBuilding ...
func (b *OGame) BuildBuilding(planetID PlanetID, buildingID int, cancel bool) error {
	planetIDStr := strconv.Itoa(int(planetID))
	url2 := "https://" + b.serverURL + "/game/index.php?page=resources&cp=" + planetIDStr
	req, _ := http.NewRequest("GET", url2, nil)
	resp, _ := b.client.Do(req)
	defer resp.Body.Close()
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	token, _ := doc.Find("form").Find("input[name=token]").Attr("value")
	modus := "1"
	if cancel {
		modus = "2"
	}
	payload := url.Values{
		"modus": {modus},
		"token": {token},
		"type":  {strconv.Itoa(buildingID)},
	}
	fmt.Println(payload)
	resp, err := b.client.PostForm(url2, payload)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

// PlanetID ...
type PlanetID int

// GetPlanetIDs returns the user planets ids
func (b *OGame) GetPlanetIDs() []PlanetID {
	req, _ := http.NewRequest("GET", "https://"+b.serverURL+"/game/index.php?page=overview", nil)
	resp, _ := b.client.Do(req)
	defer resp.Body.Close()
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	res := make([]PlanetID, 0)
	doc.Find("div.smallplanet").Each(func(i int, s *goquery.Selection) {
		el, _ := s.Attr("id")
		id, err := strconv.Atoi(strings.TrimPrefix(el, "planet-"))
		if err != nil {
			return
		}
		res = append(res, PlanetID(id))
	})
	return res
}

// PlanetInfos ...
type PlanetInfos struct {
	Img        string
	ID         PlanetID
	Name       string
	Diameter   int
	Coordinate struct {
		Galaxy   int
		System   int
		Position int
	}
	Fields struct {
		Built int
		Total int
	}
	Temperature struct {
		Min int
		Max int
	}
}

// GetPlanetInfos gets infos for planetID
func (b *OGame) GetPlanetInfos(planetID PlanetID) PlanetInfos {
	planetIDStr := strconv.Itoa(int(planetID))
	req, _ := http.NewRequest("GET", "https://"+b.serverURL+"/game/index.php?page=overview&cp="+planetIDStr, nil)
	resp, _ := b.client.Do(req)
	defer resp.Body.Close()
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	s := doc.Find("div#planet-" + planetIDStr)
	if len(s.Nodes) > 0 { // planet
		title, _ := s.Find("a").Attr("title")
		root, err := html.Parse(strings.NewReader(title))
		if err != nil {
			logrus.Error(err)
		}
		txt := goquery.NewDocumentFromNode(root).Text()
		r := regexp.MustCompile(`(\w+) \[(\d+):(\d+):(\d+)\]([\d\.]+)km \((\d+)/(\d+)\)([-\d]+).+C (?:bis|to) ([-\d]+).+C`)
		m := r.FindAllStringSubmatch(txt, -1)

		res := PlanetInfos{}
		res.ID = planetID
		res.Name = m[0][1]
		res.Coordinate.Galaxy, _ = strconv.Atoi(m[0][2])
		res.Coordinate.System, _ = strconv.Atoi(m[0][3])
		res.Coordinate.Position, _ = strconv.Atoi(m[0][4])
		res.Diameter, _ = strconv.Atoi(m[0][5])
		res.Fields.Built, _ = strconv.Atoi(m[0][6])
		res.Fields.Total, _ = strconv.Atoi(m[0][7])
		res.Temperature.Min, _ = strconv.Atoi(m[0][8])
		res.Temperature.Max, _ = strconv.Atoi(m[0][9])
		return res
	}
	return PlanetInfos{}
}

// ResourceSettings ...
type ResourceSettings struct {
	MetalMine            int
	CrystalMine          int
	DeuteriumSynthesizer int
	SolarPlant           int
	FusionReactor        int
	SolarSatellite       int
}

// GetResourceSettings gets the resources settings for specified planetID
func (b *OGame) GetResourceSettings(planetID PlanetID) ResourceSettings {
	planetIDStr := strconv.Itoa(int(planetID))
	req, _ := http.NewRequest("GET", "https://"+b.serverURL+"/game/index.php?page=resourceSettings&cp="+planetIDStr, nil)
	resp, _ := b.client.Do(req)
	defer resp.Body.Close()
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
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
	}

	res := ResourceSettings{}
	res.MetalMine = vals[0]
	res.CrystalMine = vals[1]
	res.DeuteriumSynthesizer = vals[2]
	res.SolarPlant = vals[3]
	res.FusionReactor = vals[4]
	res.SolarSatellite = vals[5]

	return res
}

// GetOgameVersion returns OGame version
func (b *OGame) GetOgameVersion() string {
	req, _ := http.NewRequest("GET", "https://"+b.serverURL+"/game/index.php?page=overview", nil)
	resp, _ := b.client.Do(req)
	defer resp.Body.Close()
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	txt := doc.Find("#siteFooter").Find("a").First().Text()
	version := strings.Trim(txt, " \n\t\r")
	return version
}

// GetServerTime returns server time
func (b *OGame) GetServerTime() time.Time {
	req, _ := http.NewRequest("GET", "https://"+b.serverURL+"/game/index.php?page=overview", nil)
	resp, _ := b.client.Do(req)
	defer resp.Body.Close()
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	txt := doc.Find("li.OGameClock").First().Text()
	serverTime, _ := time.Parse("02.01.2006 15:04:05", txt)
	return serverTime
}

// ResourcesBuildings ...
type ResourcesBuildings struct {
	MetalMine            int
	CrystalMine          int
	DeuteriumSynthesizer int
	SolarPlant           int
	FusionReactor        int
	SolarSatellite       int
	MetalStorage         int
	CrystalStorage       int
	DeuteriumTank        int
}

func parseInt(val string) int {
	val = strings.Replace(val, ".", "", -1)
	val = strings.Replace(val, ",", "", -1)
	val = strings.Trim(val, " \t\r\n")
	res, _ := strconv.Atoi(val)
	return res
}

// UserInfos ...
type UserInfos struct {
	PlayerID     int
	PlayerName   string
	Points       int
	Rank         int
	Total        int
	HonourPoints int
}

// GetUserInfos gets the user informations
func (b *OGame) GetUserInfos() UserInfos {
	req, _ := http.NewRequest("GET", "https://"+b.serverURL+"/game/index.php?page=overview", nil)
	resp, _ := b.client.Do(req)
	defer resp.Body.Close()
	by, _ := ioutil.ReadAll(resp.Body)
	html := string(by)

	playerIDRgx := regexp.MustCompile(`playerId="(\w+)"`)
	playerNameRgx := regexp.MustCompile(`playerName="([^"]+)"`)
	txtContent := regexp.MustCompile(`textContent\[7\]="([^"]+)"`)
	res := UserInfos{}
	res.PlayerID, _ = strconv.Atoi(playerIDRgx.FindAllStringSubmatch(html, -1)[0][1])
	res.PlayerName = playerNameRgx.FindAllStringSubmatch(html, -1)[0][1]
	html2 := txtContent.FindAllStringSubmatch(html, -1)[0][1]
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html2))
	infosRgx := regexp.MustCompile(`([\d\\.]+) \(Place ([\d\.]+) of ([\d\.]+)\)`)
	infos := infosRgx.FindAllStringSubmatch(doc.Text(), -1)
	res.Points = parseInt(infos[0][1])
	res.Rank = parseInt(infos[0][2])
	res.Total = parseInt(infos[0][3])
	tmpRgx := regexp.MustCompile(`textContent\[9\]="([^"]+)"`)
	res.HonourPoints = parseInt(tmpRgx.FindAllStringSubmatch(html, -1)[0][1])
	return res
}

func getNbr(doc *goquery.Document, name string) (int, error) {
	div := doc.Find("div." + name)
	level := div.Find("span.level")
	return strconv.Atoi(strings.Trim(level.Contents().Text(), " \r\t\n"))
}

// GetResourcesBuildings gets the resources buildings levels
func (b *OGame) GetResourcesBuildings(planetID PlanetID) ResourcesBuildings {
	planetIDStr := strconv.Itoa(int(planetID))
	req, _ := http.NewRequest("GET", "https://"+b.serverURL+"/game/index.php?page=resources&cp="+planetIDStr, nil)
	resp, _ := b.client.Do(req)
	defer resp.Body.Close()
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	doc.Find("span.textlabel").Remove()

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
	return res
}
