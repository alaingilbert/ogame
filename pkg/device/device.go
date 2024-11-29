package device

import (
	cryptoRand "crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/alaingilbert/ogame/pkg/httpclient"
	"github.com/alaingilbert/ogame/pkg/utils"
	"github.com/martinlindhe/base36"
	cookiejar "github.com/orirawlings/persistent-cookiejar"
)

type Os string

// NowFunc allow user to override it
var NowFunc = func() time.Time {
	return time.Now()
}

const (
	Android Os = "Android"
	Windows Os = "Windows"
	MacOSX  Os = "Mac OS X"
	Linux   Os = "Linux"
	Ios     Os = "iOS"
)

type Browser string

const (
	Chrome  Browser = "Chrome"
	Opera   Browser = "Opera"
	Safari  Browser = "Safari"
	Edge    Browser = "Edge"
	Firefox Browser = "Firefox"
)

type Builder struct {
	name                string
	timezone            string
	languages           string
	osName              Os
	browserName         Browser
	browserEngineName   string
	osVersion           string
	hardwareConcurrency int
	memory              int
	screenColorDepth    int
	screenWidth         int
	screenHeight        int
	userAgent           string
	navigatorVendor     string
	webglInfo           string
	offlineAudioCtx     float64
	canvas2DInfo        int
	client              *httpclient.Client
	persistor           Persistor
}

type Persistor interface {
	Load() (*JsFingerprint, error)
	Save(*JsFingerprint) error
}

type Device struct {
	client    *httpclient.Client
	persistor Persistor
}

func (d *Device) GetClient() *httpclient.Client {
	return d.client
}

func (d *Device) SetClient(client *httpclient.Client) {
	d.client = client
}

// WithClient allows to use a temporary http client for a specific call.
func (d *Device) WithClient(tmpClient *http.Client, clb func()) {
	_ = d.WithClientE(tmpClient, func() error {
		clb()
		return nil
	})
}

func (d *Device) WithClientE(tmpClient *http.Client, clb func() error) error {
	oldClient := d.client.Client
	defer func() { d.client.Client = oldClient }()
	d.client.Client = tmpClient
	return clb()
}

// WithClientFn allow to change the http client until the returned function is called.
func (d *Device) WithClientFn(tmpClient *http.Client) (restoreClientFn func()) {
	oldClient := d.client.Client
	d.client.Client = tmpClient
	return func() { d.client.Client = oldClient }
}

// NewBuilder creates a new virtual device.
// If the device already exists in ~/.ogame/devices/<name> it will be loaded from there,
// otherwise will be created when calling Build.
func NewBuilder(name string) *Builder {
	return &Builder{name: name}
}

func (d *Builder) SetPersistor(persistor Persistor) *Builder {
	d.persistor = persistor
	return d
}

func (d *Builder) SetOsName(osName Os) *Builder {
	d.osName = osName
	return d
}

func (d *Builder) SetBrowserName(browserName Browser) *Builder {
	d.browserName = browserName
	return d
}

func (d *Builder) SetOsVersion(osVersion string) *Builder {
	d.osVersion = osVersion
	return d
}

func (d *Builder) SetBrowserEngineName(browserEngineName string) *Builder {
	d.browserEngineName = browserEngineName
	return d
}

func (d *Builder) SetHardwareConcurrency(hardwareConcurrency int) *Builder {
	d.hardwareConcurrency = hardwareConcurrency
	return d
}

func (d *Builder) SetCanvas2DInfo(canvas2DInfo int) *Builder {
	d.canvas2DInfo = canvas2DInfo
	return d
}

func (d *Builder) SetMemory(memory int) *Builder {
	d.memory = memory
	return d
}

func (d *Builder) SetOfflineAudioCtx(offlineAudioCtx float64) *Builder {
	d.offlineAudioCtx = offlineAudioCtx
	return d
}

func (d *Builder) ScreenColorDepth(screenColorDepth int) *Builder {
	d.screenColorDepth = screenColorDepth
	return d
}

func (d *Builder) SetScreenWidth(screenWidth int) *Builder {
	d.screenWidth = screenWidth
	return d
}

func (d *Builder) SetScreenHeight(screenHeight int) *Builder {
	d.screenHeight = screenHeight
	return d
}

func (d *Builder) SetWebglInfo(webglInfo string) *Builder {
	d.webglInfo = webglInfo
	return d
}

func (d *Builder) SetUserAgent(userAgent string) *Builder {
	d.userAgent = userAgent
	return d
}

func (d *Builder) SetNavigatorVendor(navigatorVendor string) *Builder {
	d.navigatorVendor = navigatorVendor
	return d
}

func (d *Builder) SetTimezone(timezone string) *Builder {
	d.timezone = timezone
	return d
}

func (d *Builder) SetLanguages(languages string) *Builder {
	d.languages = languages
	return d
}

func DefaultStoragePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(home, ".ogame", "storage")
}

// FilePersistor default persistor for the javascript fingerprint
// save/load the fingerprint from ~/.ogame/storage/<device_name>/fingerprint
type FilePersistor struct {
	deviceStorageDir string
}

func (f FilePersistor) Load() (*JsFingerprint, error) {
	fingerprintFilePath := filepath.Join(f.deviceStorageDir, "fingerprint")
	diskFpBy, err := os.ReadFile(fingerprintFilePath)
	if err != nil {
		return nil, err
	}
	fingerprint, err := ParseBlackbox(string(diskFpBy))
	if err != nil {
		return nil, err
	}
	return fingerprint, nil
}

func (f FilePersistor) Save(fprt *JsFingerprint) error {
	fingerprintFilePath := filepath.Join(f.deviceStorageDir, "fingerprint")
	by, err := json.Marshal(fprt)
	if err != nil {
		return err
	}
	if err := os.WriteFile(fingerprintFilePath, by, 0644); err != nil {
		return err
	}
	return nil
}

func (d *Builder) newFingerprint() (*JsFingerprint, error) {
	if d.timezone == "" {
		return nil, errors.New("timezone must be specified")
	}
	if d.osName == "" {
		return nil, errors.New("os must be specified")
	}
	if d.browserName == "" {
		return nil, errors.New("browser must be specified")
	}
	if !utils.InArr(d.timezone, timezones) {
		return nil, errors.New("timezone is not valid")
	}
	if d.languages == "" {
		d.languages = "en-US,en"
	}
	if d.offlineAudioCtx == 0 {
		d.offlineAudioCtx = utils.RandFloat(123.8, 124.9)
	}
	if d.languages != "" {
		parts := strings.Split(d.languages, ",")
		for _, part := range parts {
			if !utils.InArr(part, languages) {
				return nil, errors.New("languages is not valid")
			}
		}
	}
	if d.canvas2DInfo == 0 {
		d.canvas2DInfo = int(utils.Random(261334512, 1902830807))
	}
	if d.browserEngineName == "" {
		d.setRandomBrowserEngineName()
		if d.browserEngineName == "" {
			return nil, errors.New("browserEngineName must be specified")
		}
	}
	if d.osVersion == "" {
		d.setRandomOsVersion()
		if d.osVersion == "" {
			return nil, errors.New("osVersion must be specified")
		}
	}
	if d.memory == 0 {
		d.setRandomMemory()
		if d.memory == 0 {
			return nil, errors.New("memory must be specified")
		}
	}
	if d.hardwareConcurrency == 0 {
		d.setRandomHardwareConcurrency()
		if d.hardwareConcurrency == 0 {
			return nil, errors.New("hardwareConcurrency must be specified")
		}
	}
	if d.screenColorDepth == 0 {
		d.setRandomScreenColorDepth()
		if d.screenColorDepth == 0 {
			return nil, errors.New("screenColorDepth must be specified")
		}
	}
	if d.screenWidth == 0 || d.screenHeight == 0 {
		d.setRandomScreenSize()
		if d.screenWidth == 0 || d.screenHeight == 0 {
			return nil, errors.New("screenWidth/screenHeight must be specified")
		}
	}
	if d.navigatorVendor == "" {
		d.setRandomNavigatorVendor()
		if d.navigatorVendor == "" && d.browserName != Firefox {
			return nil, errors.New("navigatorVendor must be specified")
		}
	}
	if d.webglInfo == "" {
		d.setRandomWebglInfo()
		if d.webglInfo == "" {
			return nil, errors.New("webglInfo must be specified")
		}
	}
	if d.userAgent == "" {
		d.setRandomUserAgent()
		if d.userAgent == "" {
			return nil, errors.New("userAgent must be specified")
		}
	}

	fprt := &JsFingerprint{
		ConstantVersion:       9,
		UserAgent:             d.userAgent,
		BrowserName:           string(d.browserName),
		BrowserEngineName:     d.browserEngineName,
		NavigatorVendor:       d.navigatorVendor,
		WebglInfo:             d.webglInfo,
		XVecB64:               base64.StdEncoding.EncodeToString([]byte(GenNewXVec())),
		XGame:                 Get27RandChars(3),
		Timezone:              d.timezone,
		OsName:                string(d.osName),
		Version:               d.osVersion,
		Languages:             d.languages,
		DeviceMemory:          d.memory,
		HardwareConcurrency:   d.hardwareConcurrency,
		ScreenWidth:           d.screenWidth,
		ScreenHeight:          d.screenHeight,
		ScreenColorDepth:      d.screenColorDepth,
		OfflineAudioCtx:       d.offlineAudioCtx,
		Canvas2DInfo:          d.canvas2DInfo,
		DateIso:               NowFunc().UTC().Format(javascriptISOString),
		NavigatorDoNotTrack:   false,
		LocalStorageEnabled:   true,
		SessionStorageEnabled: true,
		VideoHash:             randFakeHash(),
		AudioCtxHash:          randFakeHash(),
		AudioHash:             randFakeHash(),
		FontsHash:             randFakeHash(),
		PluginsHash:           randFakeHash(),
		MediaDevicesHash:      randFakeHash(),
		PermissionsStatesHash: randFakeHash(),
		WebglRenderHash:       randFakeHash(),
		//Game1DateHeader:       game1DateHeader,
		//CalcDeltaMs:           elapsed,
	}
	return fprt, nil
}

func (d *Builder) Build() (*Device, error) {
	defaultStoragePath := DefaultStoragePath()
	deviceStorageDir := filepath.Join(defaultStoragePath, d.name)
	if err := os.MkdirAll(deviceStorageDir, 0755); err != nil {
		return nil, err
	}

	if d.persistor == nil {
		d.persistor = &FilePersistor{deviceStorageDir: deviceStorageDir}
	}

	fprt, err := d.persistor.Load()
	if err != nil {
		fprt, err = d.newFingerprint()
		if err != nil {
			return nil, err
		}
		if err := d.persistor.Save(fprt); err != nil {
			return nil, err
		}
	}
	d.userAgent = fprt.UserAgent

	if d.client == nil {
		jar, err := cookiejar.New(&cookiejar.Options{
			Filename:              filepath.Join(defaultStoragePath, d.name, "cookies"),
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

		d.client = httpclient.NewClient(d.userAgent)
		d.client.Jar = jar
	}

	return &Device{
		persistor: d.persistor,
		client:    d.client,
	}, nil
}

type JsFingerprint struct {
	ConstantVersion       int
	UserAgent             string
	BrowserName           string
	BrowserEngineName     string
	NavigatorVendor       string
	WebglInfo             string
	XVecB64               string
	XGame                 string
	Timezone              string
	OsName                string
	Version               string
	Languages             string
	DeviceMemory          int
	HardwareConcurrency   int
	ScreenWidth           int
	ScreenHeight          int
	ScreenColorDepth      int
	OfflineAudioCtx       float64
	Canvas2DInfo          int
	DateIso               string
	Game1DateHeader       string
	CalcDeltaMs           int64
	NavigatorDoNotTrack   bool
	LocalStorageEnabled   bool
	SessionStorageEnabled bool
	VideoHash             string
	AudioCtxHash          string
	AudioHash             string
	FontsHash             string
	PluginsHash           string
	MediaDevicesHash      string
	PermissionsStatesHash string
	WebglRenderHash       string
}

const javascriptISOString = "2006-01-02T15:04:05.999Z07:00"

func (d *Device) GetBlackbox() (string, error) {
	game1DateHeader, elapsed, err := getGame1Js(d.client)
	if err != nil {
		return "", err
	}

	fprt, err := d.persistor.Load()
	if err != nil {
		return "", err
	}
	xVecBy, err := base64.StdEncoding.DecodeString(fprt.XVecB64)
	xVec := string(xVecBy)
	if err != nil || xVec == "" {
		xVec = GenNewXVec()
	}
	newXVec := rotateXVec(xVec)
	fprt.XVecB64 = base64.StdEncoding.EncodeToString([]byte(newXVec))

	fprt.Game1DateHeader = game1DateHeader
	fprt.CalcDeltaMs = elapsed

	if err := d.persistor.Save(fprt); err != nil {
		return "", err
	}

	by, err := json.Marshal(fprt)
	if err != nil {
		return "", err
	}

	encrypted := EncryptBlackbox(string(by))

	return encrypted, nil
}

func (f *JsFingerprint) MarshalJSON() ([]byte, error) {
	toEnc := make([]any, 0)
	toEnc = append(toEnc, f.ConstantVersion)     // dg constant
	toEnc = append(toEnc, f.Timezone)            // dO4 Intl.DateTimeFormat()['resolvedOptions']().timeZone
	toEnc = append(toEnc, f.NavigatorDoNotTrack) // ZNJG navigator.doNotTrack || false
	toEnc = append(toEnc, f.BrowserEngineName)   // 'cOJRtSqNAQ': getBrowserRenderingEngine(browserInfo['name']),
	toEnc = append(toEnc, f.OsName)              // 'b-I2rx-E': osInfo['name'],
	toEnc = append(toEnc, f.BrowserName)         // 'YdFB': browserInfo['name'],
	toEnc = append(toEnc, f.NavigatorVendor)     // 'dttJrRyO': navigator.vendor,
	toEnc = append(toEnc, f.DeviceMemory)        // 'bdI_': navigator.deviceMemory || 0,
	toEnc = append(toEnc, f.HardwareConcurrency) // 'Y9JA': navigator.hardwareConcurrency || 0,
	toEnc = append(toEnc, f.Languages)           // 'bM07og': navigator.languages.join(','),
	toEnc = append(toEnc, f.PluginsHash)         // 'cNxRuCGPAg': produceDeterministicHash(getPluginsInfo()),
	toEnc = append(toEnc, f.WebglInfo)           // 'Z9dM': webglInfo['vendor'] + ',' + webglInfo['renderer'],
	toEnc = append(toEnc, f.FontsHash)           // 'ZtVDtyo': produceDeterministicHash(getFontsInfo()),
	toEnc = append(toEnc, f.AudioCtxHash)        // 'YdY6oxJV': produceDeterministicHash(getAudioContextInfo()),
	toEnc = append(toEnc, f.ScreenWidth)         // 'd-BEuCA': window.screen.availWidth,
	toEnc = append(toEnc, f.ScreenHeight)        // 'aM02nQV5': window.screen.availHeight,
	toEnc = append(toEnc, f.ScreenColorDepth)    // 'ZMk5rRU': window.screen.colorDepth,
	if f.ConstantVersion == 8 {
		toEnc = append(toEnc, f.LocalStorageEnabled)   // 'bL8zohR5': Boolean(localStorage),
		toEnc = append(toEnc, f.SessionStorageEnabled) // 'c8Y6qRuA': Boolean(sessionStorage),
	}
	toEnc = append(toEnc, f.VideoHash)             // 'dt9DqBc': produceDeterministicHash(getVideoPropsInfo()),
	toEnc = append(toEnc, f.AudioHash)             // 'YdY6oxI': produceDeterministicHash(getAudioPropsInfo()),
	toEnc = append(toEnc, f.MediaDevicesHash)      // 'bdI2nwA': produceDeterministicHash(promises[1]),
	toEnc = append(toEnc, f.PermissionsStatesHash) // 'cNVHtB2QA2zbSbw': produceDeterministicHash(promises[0]),
	toEnc = append(toEnc, f.OfflineAudioCtx)       // 'YdY6oxJYqA': promises[2],
	toEnc = append(toEnc, f.WebglRenderHash)       // 'd9w-pRFXpw': getWebglRenderInfoHash(),
	toEnc = append(toEnc, f.Canvas2DInfo)          // 'Y8QyqAl8whI': getCanvas2dInfo(),
	toEnc = append(toEnc, f.DateIso)               // objToEncrypt['Y9U6mw9451U'] = new Date().toISOString();
	toEnc = append(toEnc, f.XGame)                 // objToEncrypt['depTtw'] = xGame;
	toEnc = append(toEnc, f.CalcDeltaMs)           // objToEncrypt['ZA'] = new Date().getTime() - nowTimestamp;
	toEnc = append(toEnc, f.Version)               // 'b-I4nQ-C61rI': osInfo['version'],
	toEnc = append(toEnc, f.XVecB64)               // objToEncrypt['dts-siGT'] = window.btoa(newXVec);
	toEnc = append(toEnc, f.UserAgent)             // objToEncrypt['dehNvwBnzDqu'] = navigator.userAgent;
	toEnc = append(toEnc, f.Game1DateHeader)       // objToEncrypt['c9hKwCWX61TBJm_dKn0'] = new Date(httpReq.getResponseHeader('date')).toISOString();
	toEnc = append(toEnc, nil)                     // objToEncrypt['ctdIvSKVCQ'] = arg2;
	return json.Marshal(toEnc)
}

func randChar() rune {
	return rune(int64(32+rand.Float64()*94) | 0)
}

func Get27RandChars(n int) string {
	res := ""
	for i := 0; i < n; i++ {
		r := rand.Uint64()
		s := base36.Encode(r)[:9]
		res += s
	}
	return strings.ToLower(res)
}

func GenNewXVec() string {
	part1 := ""
	for i := 0; i < 100; i++ {
		part1 += string(randChar())
	}
	ts := NowFunc().UnixMilli()
	return fmt.Sprintf("%s %d", part1, ts)
}

func rotateXVec(xvec string) string {
	nowTs := NowFunc().UnixMilli()
	part1 := xvec[:100]
	prevTs := utils.DoParseI64(xvec[101:])
	if prevTs+1000 < nowTs {
		part1 = part1[1:] + string(randChar())
	}
	return fmt.Sprintf("%s %d", part1, nowTs)
}

func getGame1Js(client httpclient.IHttpClient) (dateHeader string, elapsed int64, err error) {
	before := NowFunc()
	resp, err := client.Get("https://gameforge.com/tra/game1.js")
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("invalid game1 status %s", resp.Status)
	}
	elapsed = time.Since(before).Milliseconds()
	date := resp.Header.Get("date")
	parsed, err := time.Parse(http.TimeFormat, date)
	if err != nil {
		return "", 0, err
	}
	return parsed.Format(javascriptISOString), elapsed, nil
}

func randFakeHash() string {
	buf := make([]byte, 32)
	_, _ = cryptoRand.Read(buf)
	return hex.EncodeToString(buf)
}

func (d *Builder) setRandomScreenColorDepth() {
	if d.osName == Android {
		d.screenColorDepth = 24
	} else if d.osName == Ios {
		d.screenColorDepth = 32
	} else if d.osName == Windows {
		d.screenColorDepth = 24
	} else if d.osName == MacOSX {
		d.screenColorDepth = 24
	}
}

func (d *Builder) setRandomHardwareConcurrency() {
	if d.osName == Ios {
		d.hardwareConcurrency = 4
	} else if d.osName == MacOSX {
		d.hardwareConcurrency = utils.RandChoice([]int{
			4, 8, 12,
		})
	} else if d.osName == Windows {
		d.hardwareConcurrency = utils.RandChoice([]int{
			2, 3, 6, 8, 12, 16, 24,
		})
	} else if d.osName == Android {
		// 8, 49, 234, 127, 128, 107, 113, 412
		d.hardwareConcurrency = utils.RandChoice([]int{
			4, 8, 12,
		})
	}
}

func (d *Builder) setRandomScreenSize() {
	if d.osName == Android {
		choices := [][]int{
			{393, 873},
			{873, 393},
			{412, 938},
			{384, 824},
			{360, 760},
			{384, 857},
			{412, 915},
			{396, 880},
			{412, 732},
			{400, 889},
			{915, 412},
			{412, 919},
			{800, 600},
			{412, 869},
			{360, 800},
			{360, 780},
			{384, 684},
		}
		choice := utils.RandChoice(choices)
		d.screenWidth = choice[0]
		d.screenHeight = choice[1]
	} else if d.osName == Ios {
		choices := [][]int{
			{428, 926},
			{390, 844},
			{414, 896},
			{430, 932},
			{375, 812},
		}
		choice := utils.RandChoice(choices)
		d.screenWidth = choice[0]
		d.screenHeight = choice[1]
	} else if d.osName == Windows {
		choices := [][]int{
			{2560, 1080},
			{1920, 1050},
			{1366, 738},
			{1680, 1010},
			{1280, 672},
			{1920, 1040},
			{1858, 1080},
			{3840, 1032},
			{1440, 860},
			{1920, 1160},
			{1280, 680},
			{1920, 1080},
			{1536, 834},
			{1304, 768},
			{1366, 736},
			{3440, 1440},
			{3153, 1276},
			{1504, 955},
			{1360, 720},
			{1280, 720},
			{1440, 870},
			{2560, 1392},
			{1920, 1032},
			{1366, 728},
			{1280, 920},
			{1600, 900},
			{1280, 728},
			{2560, 1400},
			{3440, 1392},
			{2195, 1235},
			{1536, 864},
			{1280, 994},
			{2560, 1080},
			{1600, 860},
			{1536, 824},
			{1366, 768},
			{2560, 1032},
			{2560, 1410},
		}
		choice := utils.RandChoice(choices)
		d.screenWidth = choice[0]
		d.screenHeight = choice[1]
	} else if d.osName == MacOSX {
		choices := [][]int{
			{1800, 1070},
			{1440, 875},
			{1680, 951},
			{1680, 967},
			{1440, 803},
			{2560, 1322},
			{1024, 768},
			{1440, 793},
			{2560, 1415},
			{1920, 991},
			{1680, 965},
		}
		choice := utils.RandChoice(choices)
		d.screenWidth = choice[0]
		d.screenHeight = choice[1]
	}
}

func (d *Builder) setRandomNavigatorVendor() {
	if d.osName == Ios {
		if d.browserName == Safari {
			d.navigatorVendor = "Apple Computer"
		}
	} else {
		if d.browserName == Chrome {
			d.navigatorVendor = "Google Inc."
		}
	}
}

func (d *Builder) setRandomWebglInfo() {
	if d.osName == Windows || d.osName == Linux {
		if d.browserName == Chrome {
			d.webglInfo = utils.RandChoice([]string{
				"Google Inc. (NVIDIA),ANGLE (NVIDIA, NVIDIA GeForce GTX 1060 Direct3D11 vs_5_0 ps_5_0, D3D11)",
				"Google Inc. (Intel),ANGLE (Intel, Intel(R) HD Graphics 530 Direct3D11 vs_5_0 ps_5_0, D3D11)",
				"Google Inc. (AMD),ANGLE (AMD, Radeon RX 580 Series Direct3D11 vs_5_0 ps_5_0, D3D11)",
				"Google Inc. (AMD),ANGLE (AMD, Radeon (TM) RX 480 Graphics Direct3D11 vs_5_0 ps_5_0, D3D11)",
				"Google Inc. (NVIDIA),ANGLE (NVIDIA, NVIDIA GeForce RTX 3090 Direct3D11 vs_5_0 ps_5_0, D3D11)",
				"Google Inc. (NVIDIA),ANGLE (NVIDIA, NVIDIA GeForce GTX 1080 Direct3D11 vs_5_0 ps_5_0, D3D11)",
				"Google Inc. (AMD),ANGLE (AMD, AMD Radeon(TM) Vega 8 Graphics Direct3D11 vs_5_0 ps_5_0, D3D11)",
				"Google Inc. (Intel),ANGLE (Intel, Intel(R) UHD Graphics Direct3D11 vs_5_0 ps_5_0, D3D11)",
				"Google Inc. (NVIDIA),ANGLE (NVIDIA, NVIDIA GeForce RTX 3060 Ti Direct3D11 vs_5_0 ps_5_0, D3D11)",
				"Google Inc. (Intel),ANGLE (Intel, Intel(R) HD Graphics Family Direct3D11 vs_5_0 ps_5_0, D3D11)",
			})
		} else if d.browserName == Edge {
			d.webglInfo = utils.RandChoice([]string{
				"Google Inc. (NVIDIA),ANGLE (NVIDIA, NVIDIA GeForce GTX 1650 Direct3D11 vs_5_0 ps_5_0, D3D11)",
				"Google Inc. (NVIDIA),ANGLE (NVIDIA, NVIDIA GeForce GTX 750 Direct3D11 vs_5_0 ps_5_0, D3D11)",
			})
		} else if d.browserName == Firefox {
			d.webglInfo = utils.RandChoice([]string{
				"Google Inc. (Intel),ANGLE (Intel, Intel(R) HD Graphics Direct3D11 vs_5_0 ps_5_0)",
			})
		} else if d.browserName == Opera {
			d.webglInfo = utils.RandChoice([]string{
				"Google Inc. (NVIDIA),ANGLE (NVIDIA, NVIDIA GeForce GTX 1660 Ti Direct3D11 vs_5_0 ps_5_0, D3D11)",
			})
		}
	} else if d.osName == MacOSX {
		if d.browserName == Safari {
			d.webglInfo = utils.RandChoice([]string{
				"Apple Inc.,Apple GPU",
			})
		} else if d.browserName == Chrome {
			d.webglInfo = utils.RandChoice([]string{
				"Google Inc. (ATI Technologies Inc.),ANGLE (ATI Technologies Inc., AMD Radeon Pro 5300 OpenGL Engine, OpenGL 4.1)",
				"Google Inc. (ATI Technologies Inc.),ANGLE (ATI Technologies Inc., AMD Radeon Pro 5300M OpenGL Engine, OpenGL 4.1)",
				"Google Inc. (Intel Inc.),ANGLE (Intel Inc., Intel(R) Iris(TM) Plus Graphics OpenGL Engine, OpenGL 4.1)",
				"Google Inc. (Apple),ANGLE (Apple, Apple M1 Pro, OpenGL 4.1)",
				"Google Inc. (Apple),ANGLE (Apple, Apple M1, OpenGL 4.1)",
			})
		} else if d.browserName == Firefox {
			d.webglInfo = utils.RandChoice([]string{
				"ATI Technologies Inc.,Radeon R9 200 Series",
			})
		}
	} else if d.osName == Ios {
		if d.browserName == Safari {
			d.webglInfo = utils.RandChoice([]string{
				"Apple Inc.,Apple GPU",
			})
		}
	}
}

func (d *Builder) setRandomUserAgent() {
	userAgent := ""
	if d.osName == MacOSX {
		if d.browserName == Firefox {
			userAgent = utils.RandChoice([]string{
				"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/109.0",
			})
		} else if d.browserName == Safari {
			userAgent = utils.RandChoice([]string{
				"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.2 Safari/605.1.15",
				"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.3 Safari/605.1.15",
				"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.2 Safari/605.1.15",
				"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.6.1 Safari/605.1.15",
				"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Safari/605.1.15",
			})
		} else if d.browserName == Chrome {
			userAgent = utils.RandChoice([]string{
				"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36",
				"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
				"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
			})
		}
	} else if d.osName == Windows {
		if d.browserName == Opera {
			userAgent = utils.RandChoice([]string{
				"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36 OPR/94.0.0.0",
				"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36 OPR/94.0.0.0 (Edition std-1)",
			})
		} else if d.browserName == Firefox {
			userAgent = utils.RandChoice([]string{
				"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/109.0",
				"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:108.0) Gecko/20100101 Firefox/108.0",
				"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:109.0) Gecko/20100101 Firefox/109.0",
			})
		} else if d.browserName == Chrome {
			userAgent = utils.RandChoice([]string{
				"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
				"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36",
				"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
				"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36",
				"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36",
				"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/106.0.0.0 Safari/537.36",
				"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36",
			})
		}
	} else if d.osName == Linux {
		if d.browserName == Chrome {
			userAgent = utils.RandChoice([]string{
				"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
			})
		} else if d.browserName == Firefox {
			userAgent = utils.RandChoice([]string{
				"Mozilla/5.0 (X11; Linux armv7l; rv:91.0) Gecko/20100101 Firefox/91.0",
				"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/109.0",
			})
		}
	} else if d.osName == Android {
		if d.browserName == Chrome {
			userAgent = utils.RandChoice([]string{
				"Mozilla/5.0 (Linux; Android 8.0; Pixel 2 Build/OPD3.170816.012) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Mobile Safari/537.36",
				"Mozilla/5.0 (Linux; Android 13; SM-G981B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Mobile Safari/537.36",
				"Mozilla/5.0 (Linux; Android 13; SM-S908E) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Mobile Safari/537.36",
				"Mozilla/5.0 (Linux; Android 11; M2102J20SG) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Mobile Safari/537.36",
				"Mozilla/5.0 (Linux; Android 12; SM-G973F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Mobile Safari/537.36",
				"Mozilla/5.0 (Linux; Android 12; 2107113SG) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Mobile Safari/537.36",
				"Mozilla/5.0 (Linux; Android 12; M2103K19PY) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Mobile Safari/537.36",
				"Mozilla/5.0 (Linux; Android 12; LM-V600) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Mobile Safari/537.36",
				"Mozilla/5.0 (Linux; Android 12; SM-N975F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Mobile Safari/537.36",
				"Mozilla/5.0 (Linux; Android 11; 2201116SG) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Mobile Safari/537.36",
				"Mozilla/5.0 (Linux; Android 12; M2102J20SG) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Mobile Safari/537.36",
				"Mozilla/5.0 (Linux; Android 12; SM-A315G) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Mobile Safari/537.36",
				"Mozilla/5.0 (Linux; Android 10; SNE-LX1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Mobile Safari/537.36",
				"Mozilla/5.0 (Linux; Android 11; SM-A205F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Mobile Safari/537.36",
				"Mozilla/5.0 (Linux; Android 13; Pixel 7 Pro) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Mobile Safari/537.36",
				"Mozilla/5.0 (Linux; Android 12; SM-G988B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Mobile Safari/537.36",
				"Mozilla/5.0 (Linux; Android 13; LE2121) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Mobile Safari/537.36",
				"Mozilla/5.0 (Linux; Android 12; SM-G970F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Mobile Safari/537.36",
				"Mozilla/5.0 (Linux; Android 10; Redmi Note 9 Pro) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Mobile Safari/537.36",
				"Mozilla/5.0 (Linux; Android 11; itel A509W) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.210 Mobile Safari/537.36",
				"Mozilla/5.0 (Linux; Android 13; SAMSUNG SM-A135F) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/19.0 Chrome/102.0.5005.125 Mobile Safari/537.36",
				"Mozilla/5.0 (Linux; Android 12; 2201117TL) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Mobile Safari/537.36",
				"Mozilla/5.0 (Android 12; Mobile; rv:109.0) Gecko/109.0 Firefox/109.0",
				"Mozilla/5.0 (Linux; Android 12; M2012K11AG) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Mobile Safari/537.36",
				"Mozilla/5.0 (Linux; Android 12; SM-G975F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Mobile Safari/537.36",
				"Mozilla/5.0 (Linux; Android 12; M2007J3SG) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Mobile Safari/537.36",
				"Mozilla/5.0 (Linux; Android 12; 2201117TY) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Mobile Safari/537.36",
				"Mozilla/5.0 (Linux; Android 11; M2103K19G) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Mobile Safari/537.36",
				"Mozilla/5.0 (Linux; Android 12; M2012K11G) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Mobile Safari/537.36",
				"Mozilla/5.0 (Linux; Android 6.0.1; Nexus 5X Build/MMB29P) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.5414.101 Mobile Safari/537.36 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
				"Mozilla/5.0 (Linux; Android 12; SM-S901B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Mobile Safari/537.36",
			})
		}
	}
	d.userAgent = userAgent
}

func (d *Builder) setRandomMemory() {
	if d.osName == Android {
		// 1, 2, 4, 8
		d.memory = 8
	}
}

func (d *Builder) setRandomOsVersion() {
	if d.osName == Android {
		d.osVersion = utils.RandChoice([]string{
			"6.0.1", "8.0", "10", "11", "12", "13",
		})
	} else if d.osName == Ios {
		d.osVersion = utils.RandChoice([]string{
			"14.4.2", "16.1.0", "16.2.0", "16.3.0", "15.6.1", "16.0.3", "15.5.0", "16.1.1", "15.3.0",
		})
	} else if d.osName == MacOSX {
		d.osVersion = utils.RandChoice([]string{
			"10_15_7", "10_15_6", "10.15",
		})
	} else if d.osName == Windows {
		d.osVersion = utils.RandChoice([]string{
			"7", "8", "10",
		})
	} else if d.osName == Linux {
		d.osVersion = utils.RandChoice([]string{
			"22.04", "20.04", "18.04",
		})
	}
}

func (d *Builder) setRandomBrowserEngineName() {
	if d.browserName == Chrome {
		d.browserEngineName = "Blink"
	} else if d.browserName == Opera {
		d.browserEngineName = "Blink"
	} else if d.browserName == Edge {
		d.browserEngineName = "Blink"
	} else if d.browserName == Safari {
		d.browserEngineName = "WebKit"
	} else if d.browserName == Firefox {
		d.browserEngineName = "Gecko"
	}
}

var languages = []string{"af", "sq", "ar-SA", "ar-IQ", "ar-EG", "ar-LY", "ar-DZ", "ar-MA", "ar-TN", "ar-OM",
	"ar-YE", "ar-SY", "ar-JO", "ar-LB", "ar-KW", "ar-AE", "ar-BH", "ar-QA", "eu", "bg",
	"be", "ca", "zh-TW", "zh-CN", "zh-HK", "zh-SG", "hr", "cs", "da", "nl", "nl-BE", "en",
	"en-US", "en-EG", "en-AU", "en-GB", "en-CA", "en-NZ", "en-IE", "en-ZA", "en-JM",
	"en-BZ", "en-TT", "et", "fo", "fa", "fi", "fr", "fr-BE", "fr-CA", "fr-CH", "fr-LU",
	"gd", "gd-IE", "de", "de-CH", "de-AT", "de-LU", "de-LI", "el", "he", "hi", "hu",
	"is", "id", "it", "it-CH", "ja", "ko", "lv", "lt", "mk", "mt", "no", "pl",
	"pt-BR", "pt", "rm", "ro", "ro-MO", "ru", "ru-MI", "sz", "sr", "sk", "sl", "sb",
	"es", "es-AR", "es-GT", "es-CR", "es-PA", "es-DO", "es-MX", "es-VE", "es-CO",
	"es-PE", "es-EC", "es-CL", "es-UY", "es-PY", "es-BO", "es-SV", "es-HN", "es-NI",
	"es-PR", "sx", "sv", "sv-FI", "th", "ts", "tn", "tr", "uk", "ur", "ve", "vi", "xh",
	"ji", "zu"}

var timezones = []string{"Africa/Abidjan", "Africa/Accra", "Africa/Addis_Ababa", "Africa/Algiers", "Africa/Asmera",
	"Africa/Bamako", "Africa/Bangui", "Africa/Banjul", "Africa/Bissau", "Africa/Blantyre", "Africa/Brazzaville",
	"Africa/Bujumbura", "Africa/Cairo", "Africa/Casablanca", "Africa/Ceuta", "Africa/Conakry", "Africa/Dakar",
	"Africa/Dar_es_Salaam", "Africa/Djibouti", "Africa/Douala", "Africa/El_Aaiun", "Africa/Freetown", "Africa/Gaborone",
	"Africa/Harare", "Africa/Johannesburg", "Africa/Juba", "Africa/Kampala", "Africa/Khartoum", "Africa/Kigali",
	"Africa/Kinshasa", "Africa/Lagos", "Africa/Libreville", "Africa/Lome", "Africa/Luanda", "Africa/Lubumbashi",
	"Africa/Lusaka", "Africa/Malabo", "Africa/Maputo", "Africa/Maseru", "Africa/Mbabane", "Africa/Mogadishu",
	"Africa/Monrovia", "Africa/Nairobi", "Africa/Ndjamena", "Africa/Niamey", "Africa/Nouakchott", "Africa/Ouagadougou",
	"Africa/Porto-Novo", "Africa/Sao_Tome", "Africa/Tripoli", "Africa/Tunis", "Africa/Windhoek", "America/Adak",
	"America/Anchorage", "America/Anguilla", "America/Antigua", "America/Araguaina", "America/Argentina/La_Rioja",
	"America/Argentina/Rio_Gallegos", "America/Argentina/Salta", "America/Argentina/San_Juan",
	"America/Argentina/San_Luis", "America/Argentina/Tucuman", "America/Argentina/Ushuaia", "America/Aruba",
	"America/Asuncion", "America/Bahia", "America/Bahia_Banderas", "America/Barbados", "America/Belem",
	"America/Belize", "America/Blanc-Sablon", "America/Boa_Vista", "America/Bogota", "America/Boise",
	"America/Buenos_Aires", "America/Cambridge_Bay", "America/Campo_Grande", "America/Cancun", "America/Caracas",
	"America/Catamarca", "America/Cayenne", "America/Cayman", "America/Chicago", "America/Chihuahua",
	"America/Coral_Harbour", "America/Cordoba", "America/Costa_Rica", "America/Creston", "America/Cuiaba",
	"America/Curacao", "America/Danmarkshavn", "America/Dawson", "America/Dawson_Creek", "America/Denver",
	"America/Detroit", "America/Dominica", "America/Edmonton", "America/Eirunepe", "America/El_Salvador",
	"America/Fort_Nelson", "America/Fortaleza", "America/Glace_Bay", "America/Godthab", "America/Goose_Bay",
	"America/Grand_Turk", "America/Grenada", "America/Guadeloupe", "America/Guatemala", "America/Guayaquil",
	"America/Guyana", "America/Halifax", "America/Havana", "America/Hermosillo", "America/Indiana/Knox",
	"America/Indiana/Marengo", "America/Indiana/Petersburg", "America/Indiana/Tell_City", "America/Indiana/Vevay",
	"America/Indiana/Vincennes", "America/Indiana/Winamac", "America/Indianapolis", "America/Inuvik", "America/Iqaluit",
	"America/Jamaica", "America/Jujuy", "America/Juneau", "America/Kentucky/Monticello", "America/Kralendijk",
	"America/La_Paz", "America/Lima", "America/Los_Angeles", "America/Louisville", "America/Lower_Princes",
	"America/Maceio", "America/Managua", "America/Manaus", "America/Marigot", "America/Martinique", "America/Matamoros",
	"America/Mazatlan", "America/Mendoza", "America/Menominee", "America/Merida", "America/Metlakatla",
	"America/Mexico_City", "America/Miquelon", "America/Moncton", "America/Monterrey", "America/Montevideo",
	"America/Montreal", "America/Montserrat", "America/Nassau", "America/New_York", "America/Nipigon", "America/Nome",
	"America/Noronha", "America/North_Dakota/Beulah", "America/North_Dakota/Center", "America/North_Dakota/New_Salem",
	"America/Ojinaga", "America/Panama", "America/Pangnirtung", "America/Paramaribo", "America/Phoenix",
	"America/Port-au-Prince", "America/Port_of_Spain", "America/Porto_Velho", "America/Puerto_Rico",
	"America/Punta_Arenas", "America/Rainy_River", "America/Rankin_Inlet", "America/Recife", "America/Regina",
	"America/Resolute", "America/Rio_Branco", "America/Santa_Isabel", "America/Santarem", "America/Santiago",
	"America/Santo_Domingo", "America/Sao_Paulo", "America/Scoresbysund", "America/Sitka", "America/St_Barthelemy",
	"America/St_Johns", "America/St_Kitts", "America/St_Lucia", "America/St_Thomas", "America/St_Vincent",
	"America/Swift_Current", "America/Tegucigalpa", "America/Thule", "America/Thunder_Bay", "America/Tijuana",
	"America/Toronto", "America/Tortola", "America/Vancouver", "America/Whitehorse", "America/Winnipeg",
	"America/Yakutat", "America/Yellowknife", "Antarctica/Casey", "Antarctica/Davis", "Antarctica/DumontDUrville",
	"Antarctica/Macquarie", "Antarctica/Mawson", "Antarctica/McMurdo", "Antarctica/Palmer", "Antarctica/Rothera",
	"Antarctica/Syowa", "Antarctica/Troll", "Antarctica/Vostok", "Arctic/Longyearbyen", "Asia/Aden", "Asia/Almaty",
	"Asia/Amman", "Asia/Anadyr", "Asia/Aqtau", "Asia/Aqtobe", "Asia/Ashgabat", "Asia/Atyrau", "Asia/Baghdad",
	"Asia/Bahrain", "Asia/Baku", "Asia/Bangkok", "Asia/Barnaul", "Asia/Beirut", "Asia/Bishkek", "Asia/Brunei",
	"Asia/Calcutta", "Asia/Chita", "Asia/Choibalsan", "Asia/Colombo", "Asia/Damascus", "Asia/Dhaka", "Asia/Dili",
	"Asia/Dubai", "Asia/Dushanbe", "Asia/Famagusta", "Asia/Gaza", "Asia/Hebron", "Asia/Hong_Kong", "Asia/Hovd",
	"Asia/Irkutsk", "Asia/Jakarta", "Asia/Jayapura", "Asia/Jerusalem", "Asia/Kabul", "Asia/Kamchatka", "Asia/Karachi",
	"Asia/Katmandu", "Asia/Khandyga", "Asia/Krasnoyarsk", "Asia/Kuala_Lumpur", "Asia/Kuching", "Asia/Kuwait",
	"Asia/Macau", "Asia/Magadan", "Asia/Makassar", "Asia/Manila", "Asia/Muscat", "Asia/Nicosia", "Asia/Novokuznetsk",
	"Asia/Novosibirsk", "Asia/Omsk", "Asia/Oral", "Asia/Phnom_Penh", "Asia/Pontianak", "Asia/Pyongyang", "Asia/Qatar",
	"Asia/Qostanay", "Asia/Qyzylorda", "Asia/Rangoon", "Asia/Riyadh", "Asia/Saigon", "Asia/Sakhalin", "Asia/Samarkand",
	"Asia/Seoul", "Asia/Shanghai", "Asia/Singapore", "Asia/Srednekolymsk", "Asia/Taipei", "Asia/Tashkent",
	"Asia/Tbilisi", "Asia/Tehran", "Asia/Thimphu", "Asia/Tokyo", "Asia/Tomsk", "Asia/Ulaanbaatar", "Asia/Urumqi",
	"Asia/Ust-Nera", "Asia/Vientiane", "Asia/Vladivostok", "Asia/Yakutsk", "Asia/Yekaterinburg", "Asia/Yerevan",
	"Atlantic/Azores", "Atlantic/Bermuda", "Atlantic/Canary", "Atlantic/Cape_Verde", "Atlantic/Faeroe",
	"Atlantic/Madeira", "Atlantic/Reykjavik", "Atlantic/South_Georgia", "Atlantic/St_Helena", "Atlantic/Stanley",
	"Australia/Adelaide", "Australia/Brisbane", "Australia/Broken_Hill", "Australia/Currie", "Australia/Darwin",
	"Australia/Eucla", "Australia/Hobart", "Australia/Lindeman", "Australia/Lord_Howe", "Australia/Melbourne",
	"Australia/Perth", "Australia/Sydney", "Europe/Amsterdam", "Europe/Andorra", "Europe/Astrakhan", "Europe/Athens",
	"Europe/Belgrade", "Europe/Berlin", "Europe/Bratislava", "Europe/Brussels", "Europe/Bucharest", "Europe/Budapest",
	"Europe/Busingen", "Europe/Chisinau", "Europe/Copenhagen", "Europe/Dublin", "Europe/Gibraltar", "Europe/Guernsey",
	"Europe/Helsinki", "Europe/Isle_of_Man", "Europe/Istanbul", "Europe/Jersey", "Europe/Kaliningrad", "Europe/Kiev",
	"Europe/Kirov", "Europe/Lisbon", "Europe/Ljubljana", "Europe/London", "Europe/Luxembourg", "Europe/Madrid",
	"Europe/Malta", "Europe/Mariehamn", "Europe/Minsk", "Europe/Monaco", "Europe/Moscow", "Europe/Oslo", "Europe/Paris",
	"Europe/Podgorica", "Europe/Prague", "Europe/Riga", "Europe/Rome", "Europe/Samara", "Europe/San_Marino",
	"Europe/Sarajevo", "Europe/Saratov", "Europe/Simferopol", "Europe/Skopje", "Europe/Sofia", "Europe/Stockholm",
	"Europe/Tallinn", "Europe/Tirane", "Europe/Ulyanovsk", "Europe/Uzhgorod", "Europe/Vaduz", "Europe/Vatican",
	"Europe/Vienna", "Europe/Vilnius", "Europe/Volgograd", "Europe/Warsaw", "Europe/Zagreb", "Europe/Zaporozhye",
	"Europe/Zurich", "Indian/Antananarivo", "Indian/Chagos", "Indian/Christmas", "Indian/Cocos", "Indian/Comoro",
	"Indian/Kerguelen", "Indian/Mahe", "Indian/Maldives", "Indian/Mauritius", "Indian/Mayotte", "Indian/Reunion",
	"Pacific/Apia", "Pacific/Auckland", "Pacific/Bougainville", "Pacific/Chatham", "Pacific/Easter", "Pacific/Efate",
	"Pacific/Enderbury", "Pacific/Fakaofo", "Pacific/Fiji", "Pacific/Funafuti", "Pacific/Galapagos", "Pacific/Gambier",
	"Pacific/Guadalcanal", "Pacific/Guam", "Pacific/Honolulu", "Pacific/Johnston", "Pacific/Kiritimati",
	"Pacific/Kosrae", "Pacific/Kwajalein", "Pacific/Majuro", "Pacific/Marquesas", "Pacific/Midway", "Pacific/Nauru",
	"Pacific/Niue", "Pacific/Norfolk", "Pacific/Noumea", "Pacific/Pago_Pago", "Pacific/Palau", "Pacific/Pitcairn",
	"Pacific/Ponape", "Pacific/Port_Moresby", "Pacific/Rarotonga", "Pacific/Saipan", "Pacific/Tahiti", "Pacific/Tarawa",
	"Pacific/Tongatapu", "Pacific/Truk", "Pacific/Wake", "Pacific/Wallis"}

func EncryptBlackbox(raw string) string {
	retardPseudoB64 := func(v []uint8) string {
		extraPadding := 0
		for len(v)%3 != 0 {
			extraPadding++
			v = append(v, 0)
		}
		chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_="
		sb := make([]uint8, 0)
		const mask = 0b11_1111
		for i := 0; i < len(v); i += 3 {
			first := uint32(v[i+0])
			second := uint32(v[i+1])
			third := uint32(v[i+2])
			packed := first<<16 | second<<8 | third<<0
			sb = append(sb, chars[(packed>>18)&mask], chars[(packed>>12)&mask], chars[(packed>>6)&mask], chars[(packed>>0)&mask])
		}
		return string(sb[:len(sb)-extraPadding])
	}
	escaped := url.QueryEscape(raw)
	sb := make([]uint8, len(escaped))
	sb[0] = escaped[0]
	for i := 1; i < len(escaped); i++ {
		sb[i] = sb[i-1] + escaped[i]
	}
	return retardPseudoB64(sb)
}

func DecryptBlackbox(encrypted string) (string, error) {
	reverseRetardPseudoB64 := func(v string) []uint8 {
		extraPadding := 0
		for len(v)%4 != 0 {
			v += "A"
			extraPadding++
		}
		chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_="
		sb := make([]uint8, 0)
		for i := 0; i < len(v); i += 4 {
			first := strings.IndexByte(chars, v[i+0])
			second := strings.IndexByte(chars, v[i+1])
			third := strings.IndexByte(chars, v[i+2])
			fourth := strings.IndexByte(chars, v[i+3])
			packed := first<<18 | second<<12 | third<<6 | fourth<<0
			sb = append(sb, uint8(packed>>16), uint8(packed>>8), uint8(packed>>0))
		}
		return sb[0 : len(sb)-extraPadding]
	}
	encrypted1 := reverseRetardPseudoB64(encrypted)
	sb := make([]uint8, len(encrypted1))
	for i := len(encrypted1) - 2; i >= 0; i-- {
		sb[i+1] = encrypted1[i+1] - encrypted1[i]
	}
	sb[0] = encrypted1[0]
	out, err := url.QueryUnescape(string(sb))
	if err != nil {
		return "", err
	}
	return out, nil
}

func ParseEncryptedBlackbox(encrypted string) (fingerprint *JsFingerprint, err error) {
	decrypted, err := DecryptBlackbox(encrypted)
	if err != nil {
		return
	}
	return ParseBlackbox(decrypted)
}

func ParseBlackbox(decrypted string) (*JsFingerprint, error) {
	dec := json.NewDecoder(strings.NewReader(decrypted))
	var arr []any
	if err := dec.Decode(&arr); err != nil {
		return nil, err
	}
	constantVersion, ok := arr[0].(float64)
	if !ok {
		return nil, errors.New("failed to parse ConstantVersion")
	}
	if constantVersion == 8 {
		return ParseBlackboxV8(decrypted)
	} else if constantVersion == 9 {
		return ParseBlackboxV9(decrypted)
	}
	return nil, errors.New("unknown blackbox version")
}

func ParseBlackboxV8(decrypted string) (*JsFingerprint, error) {
	fingerprint := &JsFingerprint{}
	dec := json.NewDecoder(strings.NewReader(decrypted))
	var arr []any
	if err := dec.Decode(&arr); err != nil {
		return nil, err
	}
	constantVersion, ok := arr[0].(float64)
	if !ok {
		return nil, errors.New("failed to parse ConstantVersion")
	}
	fingerprint.Timezone, ok = arr[1].(string)
	if !ok {
		return nil, errors.New("failed to parse Timezone")
	}
	if arr2Str, ok := arr[2].(string); ok { // arr[2] can be: '1'
		var err error
		fingerprint.NavigatorDoNotTrack, err = strconv.ParseBool(arr2Str)
		if err != nil {
			return nil, errors.New("failed to parse NavigatorDoNotTrack")
		}
	} else {
		fingerprint.NavigatorDoNotTrack, ok = arr[2].(bool)
		if !ok {
			return nil, errors.New("failed to parse NavigatorDoNotTrack")
		}
	}
	fingerprint.BrowserEngineName, ok = arr[3].(string)
	if !ok {
		return nil, errors.New("failed to parse BrowserEngineName")
	}
	fingerprint.OsName, ok = arr[4].(string)
	if !ok {
		return nil, errors.New("failed to parse OsName")
	}
	fingerprint.BrowserName, ok = arr[5].(string)
	if !ok {
		return nil, errors.New("failed to parse BrowserName")
	}
	fingerprint.NavigatorVendor, ok = arr[6].(string)
	if !ok {
		return nil, errors.New("failed to parse NavigatorVendor")
	}
	deviceMemory, ok := arr[7].(float64)
	if !ok {
		return nil, errors.New("failed to parse DeviceMemory")
	}
	hardwareConcurrency, ok := arr[8].(float64)
	if !ok {
		return nil, errors.New("failed to parse DeviceMemory")
	}
	fingerprint.HardwareConcurrency = int(hardwareConcurrency)
	fingerprint.Languages, ok = arr[9].(string)
	if !ok {
		return nil, errors.New("failed to parse Languages")
	}
	fingerprint.PluginsHash, ok = arr[10].(string)
	if !ok {
		return nil, errors.New("failed to parse PluginsHash")
	}
	fingerprint.WebglInfo, ok = arr[11].(string)
	if !ok {
		return nil, errors.New("failed to parse WebglInfo")
	}
	fingerprint.FontsHash, ok = arr[12].(string)
	if !ok {
		return nil, errors.New("failed to parse FontsHash")
	}
	fingerprint.AudioCtxHash, ok = arr[13].(string)
	if !ok {
		return nil, errors.New("failed to parse AudioCtxHash")
	}
	screenWidth, ok := arr[14].(float64)
	if !ok {
		return nil, errors.New("failed to parse ScreenWidth")
	}
	fingerprint.ScreenWidth = int(screenWidth)
	screenHeight, ok := arr[15].(float64)
	if !ok {
		return nil, errors.New("failed to parse ScreenHeight")
	}
	fingerprint.ScreenHeight = int(screenHeight)
	screenColorDepth, ok := arr[16].(float64)
	if !ok {
		return nil, errors.New("failed to parse ScreenColorDepth")
	}
	fingerprint.LocalStorageEnabled, ok = arr[17].(bool)
	if !ok {
		return nil, errors.New("failed to parse LocalStorageEnabled")
	}
	fingerprint.SessionStorageEnabled, ok = arr[18].(bool)
	if !ok {
		return nil, errors.New("failed to parse SessionStorageEnabled")
	}
	fingerprint.ScreenColorDepth = int(screenColorDepth)
	fingerprint.VideoHash, ok = arr[19].(string)
	if !ok {
		return nil, errors.New("failed to parse VideoHash")
	}
	fingerprint.AudioHash, ok = arr[20].(string)
	if !ok {
		return nil, errors.New("failed to parse AudioHash")
	}
	fingerprint.MediaDevicesHash, ok = arr[21].(string)
	if !ok {
		return nil, errors.New("failed to parse MediaDevicesHash")
	}
	fingerprint.PermissionsStatesHash, ok = arr[22].(string)
	if !ok {
		return nil, errors.New("failed to parse PermissionsStatesHash")
	}
	fingerprint.OfflineAudioCtx, ok = arr[23].(float64)
	if !ok {
		return nil, errors.New("failed to parse OfflineAudioCtx")
	}
	fingerprint.WebglRenderHash, ok = arr[24].(string)
	if !ok {
		return nil, errors.New("failed to parse WebglRenderHash")
	}
	canvas2DInfo, ok := arr[25].(float64)
	if !ok {
		return nil, errors.New("failed to parse Canvas2DInfo")
	}
	fingerprint.Canvas2DInfo = int(canvas2DInfo)
	fingerprint.DateIso, ok = arr[26].(string)
	if !ok {
		return nil, errors.New("failed to parse DateIso")
	}
	fingerprint.XGame, ok = arr[27].(string)
	if !ok {
		return nil, errors.New("failed to parse XGame")
	}
	calcDeltaMs, ok := arr[28].(float64)
	if !ok {
		return nil, errors.New("failed to parse CalcDeltaMs")
	}
	fingerprint.CalcDeltaMs = int64(calcDeltaMs)
	fingerprint.Version, ok = arr[29].(string)
	if !ok {
		return nil, errors.New("failed to parse Version")
	}
	fingerprint.DeviceMemory = int(deviceMemory)
	fingerprint.XVecB64, ok = arr[30].(string)
	if !ok {
		return nil, errors.New("failed to parse XVecB64")
	}
	fingerprint.ConstantVersion = int(constantVersion)
	fingerprint.UserAgent, ok = arr[31].(string)
	if !ok {
		return nil, errors.New("failed to parse UserAgent")
	}
	fingerprint.Game1DateHeader, ok = arr[32].(string)
	if !ok {
		return nil, errors.New("failed to parse Game1DateHeader")
	}
	return fingerprint, nil
}

func ParseBlackboxV9(decrypted string) (*JsFingerprint, error) {
	fingerprint := &JsFingerprint{}
	dec := json.NewDecoder(strings.NewReader(decrypted))
	var arr []any
	if err := dec.Decode(&arr); err != nil {
		return nil, err
	}
	constantVersion, ok := arr[0].(float64)
	if !ok {
		return nil, errors.New("failed to parse ConstantVersion")
	}
	fingerprint.Timezone, ok = arr[1].(string)
	if !ok {
		return nil, errors.New("failed to parse Timezone")
	}
	if arr2Str, ok := arr[2].(string); ok { // arr[2] can be: '1'
		var err error
		fingerprint.NavigatorDoNotTrack, err = strconv.ParseBool(arr2Str)
		if err != nil {
			return nil, errors.New("failed to parse NavigatorDoNotTrack")
		}
	} else {
		fingerprint.NavigatorDoNotTrack, ok = arr[2].(bool)
		if !ok {
			return nil, errors.New("failed to parse NavigatorDoNotTrack")
		}
	}
	fingerprint.BrowserEngineName, ok = arr[3].(string)
	if !ok {
		return nil, errors.New("failed to parse BrowserEngineName")
	}
	fingerprint.OsName, ok = arr[4].(string)
	if !ok {
		return nil, errors.New("failed to parse OsName")
	}
	fingerprint.BrowserName, ok = arr[5].(string)
	if !ok {
		return nil, errors.New("failed to parse BrowserName")
	}
	fingerprint.NavigatorVendor, ok = arr[6].(string)
	if !ok {
		return nil, errors.New("failed to parse NavigatorVendor")
	}
	deviceMemory, ok := arr[7].(float64)
	if !ok {
		return nil, errors.New("failed to parse DeviceMemory")
	}
	hardwareConcurrency, ok := arr[8].(float64)
	if !ok {
		return nil, errors.New("failed to parse DeviceMemory")
	}
	fingerprint.HardwareConcurrency = int(hardwareConcurrency)
	fingerprint.Languages, ok = arr[9].(string)
	if !ok {
		return nil, errors.New("failed to parse Languages")
	}
	fingerprint.PluginsHash, ok = arr[10].(string)
	if !ok {
		return nil, errors.New("failed to parse PluginsHash")
	}
	fingerprint.WebglInfo, ok = arr[11].(string)
	if !ok {
		return nil, errors.New("failed to parse WebglInfo")
	}
	fingerprint.FontsHash, ok = arr[12].(string)
	if !ok {
		return nil, errors.New("failed to parse FontsHash")
	}
	fingerprint.AudioCtxHash, ok = arr[13].(string)
	if !ok {
		return nil, errors.New("failed to parse AudioCtxHash")
	}
	screenWidth, ok := arr[14].(float64)
	if !ok {
		return nil, errors.New("failed to parse ScreenWidth")
	}
	fingerprint.ScreenWidth = int(screenWidth)
	screenHeight, ok := arr[15].(float64)
	if !ok {
		return nil, errors.New("failed to parse ScreenHeight")
	}
	fingerprint.ScreenHeight = int(screenHeight)
	screenColorDepth, ok := arr[16].(float64)
	if !ok {
		return nil, errors.New("failed to parse ScreenColorDepth")
	}
	fingerprint.ScreenColorDepth = int(screenColorDepth)
	fingerprint.VideoHash, ok = arr[17].(string)
	if !ok {
		return nil, errors.New("failed to parse VideoHash")
	}
	fingerprint.AudioHash, ok = arr[18].(string)
	if !ok {
		return nil, errors.New("failed to parse AudioHash")
	}
	fingerprint.MediaDevicesHash, ok = arr[19].(string)
	if !ok {
		return nil, errors.New("failed to parse MediaDevicesHash")
	}
	fingerprint.PermissionsStatesHash, ok = arr[20].(string)
	if !ok {
		return nil, errors.New("failed to parse PermissionsStatesHash")
	}
	fingerprint.OfflineAudioCtx, ok = arr[21].(float64)
	if !ok {
		return nil, errors.New("failed to parse OfflineAudioCtx")
	}
	fingerprint.WebglRenderHash, ok = arr[22].(string)
	if !ok {
		return nil, errors.New("failed to parse WebglRenderHash")
	}
	canvas2DInfo, ok := arr[23].(float64)
	if !ok {
		return nil, errors.New("failed to parse Canvas2DInfo")
	}
	fingerprint.Canvas2DInfo = int(canvas2DInfo)
	fingerprint.DateIso, ok = arr[24].(string)
	if !ok {
		return nil, errors.New("failed to parse DateIso")
	}
	fingerprint.XGame, ok = arr[25].(string)
	if !ok {
		return nil, errors.New("failed to parse XGame")
	}
	calcDeltaMs, ok := arr[26].(float64)
	if !ok {
		return nil, errors.New("failed to parse CalcDeltaMs")
	}
	fingerprint.CalcDeltaMs = int64(calcDeltaMs)
	fingerprint.Version, ok = arr[27].(string)
	if !ok {
		return nil, errors.New("failed to parse Version")
	}
	fingerprint.DeviceMemory = int(deviceMemory)
	fingerprint.XVecB64, ok = arr[28].(string)
	if !ok {
		return nil, errors.New("failed to parse XVecB64")
	}
	fingerprint.ConstantVersion = int(constantVersion)
	fingerprint.UserAgent, ok = arr[29].(string)
	if !ok {
		return nil, errors.New("failed to parse UserAgent")
	}
	fingerprint.Game1DateHeader, ok = arr[30].(string)
	if !ok {
		return nil, errors.New("failed to parse Game1DateHeader")
	}
	return fingerprint, nil
}
