package wrapper

import (
	"encoding/hex"
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

func randFakeHash() string {
	buf := make([]byte, 32)
	rand.Read(buf)
	return hex.EncodeToString(buf)
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

var ram = []int{0, 4, 8, 16, 32}

var screend = []int{8, 15, 16, 24, 32, 48}

var systems = []struct {
	Name    string
	Version string
}{
	{Name: "Windows 10", Version: "10"},
	{Name: "Windows 8.1", Version: "8.1"},
	{Name: "Windows 8", Version: "8"},
	{Name: "Windows 7", Version: "7"},
	{Name: "Unknown", Version: "Unknown"},
}

var donottrack = []string{"unspecified", "1"}

var hardware = []int{0, 1, 4, 8}

var screenResolution = []struct {
	Width  int
	Height int
}{
	{Width: 1980, Height: 1080},
}

var browsers = []struct {
	Browser string
	Engine  string
	Vendor  string
}{{Browser: "Chrome", Engine: "Blink", Vendor: "Google Inc."},
	{Browser: "Opera", Engine: "Blink", Vendor: "Google Inc."},
	{Browser: "Edge", Engine: "Trident", Vendor: "Google Inc."},
	{Browser: "Internet Explorer", Engine: "Trident", Vendor: "Google Inc."},
	{Browser: "Safari", Engine: "WebKit", Vendor: "Apple Computer, Inc."},
	{Browser: "Firefox", Engine: "Gecko", Vendor: ""}}

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
	"Africa/Harare", "Africa/Johannesburg", "Africa/Juba", "Africa/Kampala", "Africa/Khartoum", "Africa/Kigali", "Africa/Kinshasa",
	"Africa/Lagos", "Africa/Libreville", "Africa/Lome", "Africa/Luanda", "Africa/Lubumbashi", "Africa/Lusaka", "Africa/Malabo",
	"Africa/Maputo", "Africa/Maseru", "Africa/Mbabane", "Africa/Mogadishu", "Africa/Monrovia", "Africa/Nairobi", "Africa/Ndjamena",
	"Africa/Niamey", "Africa/Nouakchott", "Africa/Ouagadougou", "Africa/Porto-Novo", "Africa/Sao_Tome", "Africa/Tripoli",
	"Africa/Tunis", "Africa/Windhoek", "America/Adak", "America/Anchorage", "America/Anguilla", "America/Antigua",
	"America/Araguaina", "America/Argentina/La_Rioja", "America/Argentina/Rio_Gallegos", "America/Argentina/Salta",
	"America/Argentina/San_Juan", "America/Argentina/San_Luis", "America/Argentina/Tucuman", "America/Argentina/Ushuaia",
	"America/Aruba", "America/Asuncion", "America/Bahia", "America/Bahia_Banderas", "America/Barbados", "America/Belem",
	"America/Belize", "America/Blanc-Sablon", "America/Boa_Vista", "America/Bogota", "America/Boise", "America/Buenos_Aires",
	"America/Cambridge_Bay", "America/Campo_Grande", "America/Cancun", "America/Caracas", "America/Catamarca", "America/Cayenne",
	"America/Cayman", "America/Chicago", "America/Chihuahua", "America/Coral_Harbour", "America/Cordoba", "America/Costa_Rica",
	"America/Creston", "America/Cuiaba", "America/Curacao", "America/Danmarkshavn", "America/Dawson", "America/Dawson_Creek",
	"America/Denver", "America/Detroit", "America/Dominica", "America/Edmonton", "America/Eirunepe", "America/El_Salvador",
	"America/Fort_Nelson", "America/Fortaleza", "America/Glace_Bay", "America/Godthab", "America/Goose_Bay", "America/Grand_Turk",
	"America/Grenada", "America/Guadeloupe", "America/Guatemala", "America/Guayaquil", "America/Guyana", "America/Halifax",
	"America/Havana", "America/Hermosillo", "America/Indiana/Knox", "America/Indiana/Marengo", "America/Indiana/Petersburg",
	"America/Indiana/Tell_City", "America/Indiana/Vevay", "America/Indiana/Vincennes", "America/Indiana/Winamac", "America/Indianapolis",
	"America/Inuvik", "America/Iqaluit", "America/Jamaica", "America/Jujuy", "America/Juneau", "America/Kentucky/Monticello",
	"America/Kralendijk", "America/La_Paz", "America/Lima", "America/Los_Angeles", "America/Louisville", "America/Lower_Princes",
	"America/Maceio", "America/Managua", "America/Manaus", "America/Marigot", "America/Martinique", "America/Matamoros", "America/Mazatlan",
	"America/Mendoza", "America/Menominee", "America/Merida", "America/Metlakatla", "America/Mexico_City", "America/Miquelon",
	"America/Moncton", "America/Monterrey", "America/Montevideo", "America/Montreal", "America/Montserrat", "America/Nassau",
	"America/New_York", "America/Nipigon", "America/Nome", "America/Noronha", "America/North_Dakota/Beulah", "America/North_Dakota/Center",
	"America/North_Dakota/New_Salem", "America/Ojinaga", "America/Panama", "America/Pangnirtung", "America/Paramaribo", "America/Phoenix",
	"America/Port-au-Prince", "America/Port_of_Spain", "America/Porto_Velho", "America/Puerto_Rico", "America/Punta_Arenas",
	"America/Rainy_River", "America/Rankin_Inlet", "America/Recife", "America/Regina", "America/Resolute", "America/Rio_Branco",
	"America/Santa_Isabel", "America/Santarem", "America/Santiago", "America/Santo_Domingo", "America/Sao_Paulo", "America/Scoresbysund",
	"America/Sitka", "America/St_Barthelemy", "America/St_Johns", "America/St_Kitts", "America/St_Lucia", "America/St_Thomas",
	"America/St_Vincent", "America/Swift_Current", "America/Tegucigalpa", "America/Thule", "America/Thunder_Bay", "America/Tijuana",
	"America/Toronto", "America/Tortola", "America/Vancouver", "America/Whitehorse", "America/Winnipeg", "America/Yakutat",
	"America/Yellowknife", "Antarctica/Casey", "Antarctica/Davis", "Antarctica/DumontDUrville", "Antarctica/Macquarie",
	"Antarctica/Mawson", "Antarctica/McMurdo", "Antarctica/Palmer", "Antarctica/Rothera", "Antarctica/Syowa", "Antarctica/Troll",
	"Antarctica/Vostok", "Arctic/Longyearbyen", "Asia/Aden", "Asia/Almaty", "Asia/Amman", "Asia/Anadyr", "Asia/Aqtau", "Asia/Aqtobe",
	"Asia/Ashgabat", "Asia/Atyrau", "Asia/Baghdad", "Asia/Bahrain", "Asia/Baku", "Asia/Bangkok", "Asia/Barnaul", "Asia/Beirut",
	"Asia/Bishkek", "Asia/Brunei", "Asia/Calcutta", "Asia/Chita", "Asia/Choibalsan", "Asia/Colombo", "Asia/Damascus", "Asia/Dhaka",
	"Asia/Dili", "Asia/Dubai", "Asia/Dushanbe", "Asia/Famagusta", "Asia/Gaza", "Asia/Hebron", "Asia/Hong_Kong", "Asia/Hovd",
	"Asia/Irkutsk", "Asia/Jakarta", "Asia/Jayapura", "Asia/Jerusalem", "Asia/Kabul", "Asia/Kamchatka", "Asia/Karachi", "Asia/Katmandu",
	"Asia/Khandyga", "Asia/Krasnoyarsk", "Asia/Kuala_Lumpur", "Asia/Kuching", "Asia/Kuwait", "Asia/Macau", "Asia/Magadan",
	"Asia/Makassar", "Asia/Manila", "Asia/Muscat", "Asia/Nicosia", "Asia/Novokuznetsk", "Asia/Novosibirsk", "Asia/Omsk",
	"Asia/Oral", "Asia/Phnom_Penh", "Asia/Pontianak", "Asia/Pyongyang", "Asia/Qatar", "Asia/Qostanay", "Asia/Qyzylorda",
	"Asia/Rangoon", "Asia/Riyadh", "Asia/Saigon", "Asia/Sakhalin", "Asia/Samarkand", "Asia/Seoul", "Asia/Shanghai", "Asia/Singapore",
	"Asia/Srednekolymsk", "Asia/Taipei", "Asia/Tashkent", "Asia/Tbilisi", "Asia/Tehran", "Asia/Thimphu", "Asia/Tokyo", "Asia/Tomsk",
	"Asia/Ulaanbaatar", "Asia/Urumqi", "Asia/Ust-Nera", "Asia/Vientiane", "Asia/Vladivostok", "Asia/Yakutsk", "Asia/Yekaterinburg",
	"Asia/Yerevan", "Atlantic/Azores", "Atlantic/Bermuda", "Atlantic/Canary", "Atlantic/Cape_Verde", "Atlantic/Faeroe",
	"Atlantic/Madeira", "Atlantic/Reykjavik", "Atlantic/South_Georgia", "Atlantic/St_Helena", "Atlantic/Stanley", "Australia/Adelaide",
	"Australia/Brisbane", "Australia/Broken_Hill", "Australia/Currie", "Australia/Darwin", "Australia/Eucla", "Australia/Hobart",
	"Australia/Lindeman", "Australia/Lord_Howe", "Australia/Melbourne", "Australia/Perth", "Australia/Sydney", "Europe/Amsterdam",
	"Europe/Andorra", "Europe/Astrakhan", "Europe/Athens", "Europe/Belgrade", "Europe/Berlin", "Europe/Bratislava", "Europe/Brussels",
	"Europe/Bucharest", "Europe/Budapest", "Europe/Busingen", "Europe/Chisinau", "Europe/Copenhagen", "Europe/Dublin",
	"Europe/Gibraltar", "Europe/Guernsey", "Europe/Helsinki", "Europe/Isle_of_Man", "Europe/Istanbul", "Europe/Jersey",
	"Europe/Kaliningrad", "Europe/Kiev", "Europe/Kirov", "Europe/Lisbon", "Europe/Ljubljana", "Europe/London", "Europe/Luxembourg",
	"Europe/Madrid", "Europe/Malta", "Europe/Mariehamn", "Europe/Minsk", "Europe/Monaco", "Europe/Moscow", "Europe/Oslo",
	"Europe/Paris", "Europe/Podgorica", "Europe/Prague", "Europe/Riga", "Europe/Rome", "Europe/Samara", "Europe/San_Marino",
	"Europe/Sarajevo", "Europe/Saratov", "Europe/Simferopol", "Europe/Skopje", "Europe/Sofia", "Europe/Stockholm", "Europe/Tallinn",
	"Europe/Tirane", "Europe/Ulyanovsk", "Europe/Uzhgorod", "Europe/Vaduz", "Europe/Vatican", "Europe/Vienna", "Europe/Vilnius",
	"Europe/Volgograd", "Europe/Warsaw", "Europe/Zagreb", "Europe/Zaporozhye", "Europe/Zurich", "Indian/Antananarivo", "Indian/Chagos",
	"Indian/Christmas", "Indian/Cocos", "Indian/Comoro", "Indian/Kerguelen", "Indian/Mahe", "Indian/Maldives", "Indian/Mauritius",
	"Indian/Mayotte", "Indian/Reunion", "Pacific/Apia", "Pacific/Auckland", "Pacific/Bougainville", "Pacific/Chatham", "Pacific/Easter",
	"Pacific/Efate", "Pacific/Enderbury", "Pacific/Fakaofo", "Pacific/Fiji", "Pacific/Funafuti", "Pacific/Galapagos", "Pacific/Gambier",
	"Pacific/Guadalcanal", "Pacific/Guam", "Pacific/Honolulu", "Pacific/Johnston", "Pacific/Kiritimati", "Pacific/Kosrae",
	"Pacific/Kwajalein", "Pacific/Majuro", "Pacific/Marquesas", "Pacific/Midway", "Pacific/Nauru", "Pacific/Niue", "Pacific/Norfolk",
	"Pacific/Noumea", "Pacific/Pago_Pago", "Pacific/Palau", "Pacific/Pitcairn", "Pacific/Ponape", "Pacific/Port_Moresby",
	"Pacific/Rarotonga", "Pacific/Saipan", "Pacific/Tahiti", "Pacific/Tarawa", "Pacific/Tongatapu", "Pacific/Truk", "Pacific/Wake",
	"Pacific/Wallis"}
