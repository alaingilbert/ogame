package wrapper

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alaingilbert/ogame/pkg/utils"
	"github.com/martinlindhe/base36"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

func randChar() rune {
	return rune(int64(32+rand.Float64()*94) | 0)
}

func genNewXVec() string {
	part1 := ""
	for i := 0; i < 100; i++ {
		part1 += string(randChar())
	}
	ts := time.Now().UnixMilli()
	return fmt.Sprintf("%s %d", part1, ts)
}

func rotateXVec(xvec string) string {
	nowTs := time.Now().UnixMilli()
	parts := strings.Split(xvec, " ")
	part1 := parts[0]
	prevTs := utils.DoParseI64(parts[1])
	if prevTs+1000 < nowTs {
		part1 = part1[1:] + string(randChar())
	}
	return fmt.Sprintf("%s %d", part1, nowTs)
}

func get27RandChars(n int) string {
	res := ""
	for i := 0; i < n; i++ {
		r := rand.Uint64()
		s := base36.Encode(r)[:9]
		res += s
	}
	return strings.ToLower(res)
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

func DecryptBlackbox(encrypted string) (string, error) {
	reverseRetardPseudoB64 := func(v string) []uint8 {
		extraPadding := 0
		for len(v)%4 != 0 {
			v += "A"
			extraPadding++
		}
		chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_="
		sb := make([]uint8, 0)
		for i := 0; i < len(v); {
			first := uint32(strings.Index(chars, string(v[i])))
			i++
			second := uint32(strings.Index(chars, string(v[i])))
			i++
			third := uint32(strings.Index(chars, string(v[i])))
			i++
			fourth := uint32(strings.Index(chars, string(v[i])))
			i++
			tmpp := (first << 18) | (second << 12) | (third << 6) | fourth
			sb = append(sb, uint8(tmpp>>16&255), uint8(tmpp>>8&255), uint8(tmpp&255))
		}
		sb = sb[0 : len(sb)-extraPadding]
		return sb
	}
	encrypted1 := reverseRetardPseudoB64(encrypted)
	sb := ""
	for i := len(encrypted1) - 2; i >= 0; i-- {
		sb = string(uint8(((uint32(encrypted1[i+1])+256)-uint32(encrypted1[i]))%256)) + sb
	}
	sb = string(encrypted1[0]) + sb
	out, err := url.QueryUnescape(sb)
	if err != nil {
		return "", err
	}
	return out, nil
}

func ParseBlackbox(encrypted string) (fingerprint JsFingerprint, err error) {
	decrypted, err := DecryptBlackbox(encrypted)
	if err != nil {
		return
	}
	dec := json.NewDecoder(strings.NewReader(decrypted))
	var arr []any
	if err := dec.Decode(&arr); err != nil {
		return fingerprint, err
	}
	constantVersion, ok := arr[0].(float64)
	if !ok {
		return fingerprint, errors.New("failed to parse ConstantVersion")
	}
	fingerprint.ConstantVersion = int(constantVersion)
	fingerprint.UserAgent, ok = arr[31].(string)
	if !ok {
		return fingerprint, errors.New("failed to parse UserAgent")
	}
	fingerprint.BrowserName, ok = arr[5].(string)
	if !ok {
		return fingerprint, errors.New("failed to parse BrowserName")
	}
	fingerprint.BrowserEngineName, ok = arr[3].(string)
	if !ok {
		return fingerprint, errors.New("failed to parse BrowserEngineName")
	}
	fingerprint.NavigatorVendor, ok = arr[6].(string)
	if !ok {
		return fingerprint, errors.New("failed to parse NavigatorVendor")
	}
	fingerprint.WebglInfo, ok = arr[11].(string)
	if !ok {
		return fingerprint, errors.New("failed to parse WebglInfo")
	}
	fingerprint.XVecB64, ok = arr[30].(string)
	if !ok {
		return fingerprint, errors.New("failed to parse XVecB64")
	}
	fingerprint.XGame, ok = arr[27].(string)
	if !ok {
		return fingerprint, errors.New("failed to parse XGame")
	}
	fingerprint.Timezone, ok = arr[1].(string)
	if !ok {
		return fingerprint, errors.New("failed to parse Timezone")
	}
	fingerprint.OsName, ok = arr[4].(string)
	if !ok {
		return fingerprint, errors.New("failed to parse OsName")
	}
	fingerprint.Version, ok = arr[29].(string)
	if !ok {
		return fingerprint, errors.New("failed to parse Version")
	}
	fingerprint.Languages, ok = arr[9].(string)
	if !ok {
		return fingerprint, errors.New("failed to parse Languages")
	}
	deviceMemory, ok := arr[7].(float64)
	if !ok {
		return fingerprint, errors.New("failed to parse DeviceMemory")
	}
	fingerprint.DeviceMemory = int(deviceMemory)
	hardwareConcurrency, ok := arr[8].(float64)
	if !ok {
		return fingerprint, errors.New("failed to parse DeviceMemory")
	}
	fingerprint.HardwareConcurrency = int(hardwareConcurrency)
	screenWidth, ok := arr[14].(float64)
	if !ok {
		return fingerprint, errors.New("failed to parse ScreenWidth")
	}
	fingerprint.ScreenWidth = int(screenWidth)
	screenHeight, ok := arr[15].(float64)
	if !ok {
		return fingerprint, errors.New("failed to parse ScreenHeight")
	}
	fingerprint.ScreenHeight = int(screenHeight)
	screenColorDepth, ok := arr[16].(float64)
	if !ok {
		return fingerprint, errors.New("failed to parse ScreenColorDepth")
	}
	fingerprint.ScreenColorDepth = int(screenColorDepth)
	fingerprint.OfflineAudioCtx, ok = arr[23].(float64)
	if !ok {
		return fingerprint, errors.New("failed to parse OfflineAudioCtx")
	}
	canvas2DInfo, ok := arr[25].(float64)
	if !ok {
		return fingerprint, errors.New("failed to parse Canvas2DInfo")
	}
	fingerprint.Canvas2DInfo = int(canvas2DInfo)
	fingerprint.DateIso, ok = arr[26].(string)
	if !ok {
		return fingerprint, errors.New("failed to parse DateIso")
	}
	fingerprint.Game1DateHeader, ok = arr[32].(string)
	if !ok {
		return fingerprint, errors.New("failed to parse Game1DateHeader")
	}
	calcDeltaMs, ok := arr[28].(float64)
	if !ok {
		return fingerprint, errors.New("failed to parse CalcDeltaMs")
	}
	fingerprint.CalcDeltaMs = int64(calcDeltaMs)
	fingerprint.NavigatorDoNotTrack, ok = arr[2].(bool)
	if !ok {
		return fingerprint, errors.New("failed to parse NavigatorDoNotTrack")
	}
	fingerprint.LocalStorageEnabled, ok = arr[17].(bool)
	if !ok {
		return fingerprint, errors.New("failed to parse LocalStorageEnabled")
	}
	fingerprint.SessionStorageEnabled, ok = arr[18].(bool)
	if !ok {
		return fingerprint, errors.New("failed to parse SessionStorageEnabled")
	}
	fingerprint.VideoHash, ok = arr[19].(string)
	if !ok {
		return fingerprint, errors.New("failed to parse VideoHash")
	}
	fingerprint.AudioCtxHash, ok = arr[13].(string)
	if !ok {
		return fingerprint, errors.New("failed to parse AudioCtxHash")
	}
	fingerprint.AudioHash, ok = arr[20].(string)
	if !ok {
		return fingerprint, errors.New("failed to parse AudioHash")
	}
	fingerprint.FontsHash, ok = arr[12].(string)
	if !ok {
		return fingerprint, errors.New("failed to parse FontsHash")
	}
	fingerprint.PluginsHash, ok = arr[10].(string)
	if !ok {
		return fingerprint, errors.New("failed to parse PluginsHash")
	}
	fingerprint.MediaDevicesHash, ok = arr[21].(string)
	if !ok {
		return fingerprint, errors.New("failed to parse MediaDevicesHash")
	}
	fingerprint.PermissionsStatesHash, ok = arr[22].(string)
	if !ok {
		return fingerprint, errors.New("failed to parse PermissionsStatesHash")
	}
	fingerprint.WebglRenderHash, ok = arr[24].(string)
	if !ok {
		return fingerprint, errors.New("failed to parse WebglRenderHash")
	}
	return
}
