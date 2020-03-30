package ogame

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// ParseInt ...
func ParseInt(val string) int64 {
	val = strings.Replace(val, ".", "", -1)
	val = strings.Replace(val, ",", "", -1)
	val = strings.TrimSpace(val)
	res, _ := strconv.ParseInt(val, 10, 64)
	return res
}

func toInt(buf []byte) (n int) {
	for _, v := range buf {
		n = n*10 + int(v-'0')
	}
	return
}

// I64Ptr returns a pointer to int64
func I64Ptr(v int64) *int64 {
	return &v
}

// IsDefenseID helper returns if an integer is a defense id
func IsDefenseID(id int64) bool {
	return ID(id).IsDefense()
}

// IsShipID helper returns if an integer is a ship id
func IsShipID(id int64) bool {
	return ID(id).IsShip()
}

// IsTechID helper returns if an integer is a tech id
func IsTechID(id int64) bool {
	return ID(id).IsTech()
}

// IsBuildingID helper returns if an integer is a building id
func IsBuildingID(id int64) bool {
	return ID(id).IsBuilding()
}

// IsResourceBuildingID helper returns if an integer is a resource defense id
func IsResourceBuildingID(id int64) bool {
	return ID(id).IsResourceBuilding()
}

// IsFacilityID helper returns if an integer is a facility id
func IsFacilityID(id int64) bool {
	return ID(id).IsFacility()
}

// ParseCoord parse a coordinate from a string
func ParseCoord(str string) (coord Coordinate, err error) {
	m := regexp.MustCompile(`^\[?(([PMD]):)?(\d{1,3}):(\d{1,3}):(\d{1,3})]?$`).FindStringSubmatch(str)
	if len(m) == 6 {
		planetTypeStr := m[2]
		galaxy, _ := strconv.ParseInt(m[3], 10, 64)
		system, _ := strconv.ParseInt(m[4], 10, 64)
		position, _ := strconv.ParseInt(m[5], 10, 64)
		planetType := PlanetType
		if planetTypeStr == "M" {
			planetType = MoonType
		} else if planetTypeStr == "D" {
			planetType = DebrisType
		}
		return Coordinate{galaxy, system, position, planetType}, nil
	}
	return coord, errors.New("unable to parse coordinate")
}

func name2id(name string) ID {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	name, _, _ = transform.String(t, name)
	reg, _ := regexp.Compile("[^a-zA-ZАаБбВвГгДдЕеЁёЖжЗзИиЙйКкЛлМмНнОоПпРрСсТтУуФфХхЦцЧчШшЩщЪъЫыЬьЭэЮюЯя闘残艦収型送サ小プテバイスル輸軽船ッ戦ニトタ察デヤ洋爆ラーロ機ソ重偵回骸巡撃コ大シα-ωΑ-Ω星殖重小民死輸帶太洋戰艦諜魔間能飛鬥路輕型列探履惡大彈運導衛滅者車收巡陽機回毀船]+")
	processedString := strings.ToLower(reg.ReplaceAllString(name, ""))
	nameMap := map[string]ID{
		// en
		"lightfighter":   LightFighterID,
		"heavyfighter":   HeavyFighterID,
		"cruiser":        CruiserID,
		"battleship":     BattleshipID,
		"battlecruiser":  BattlecruiserID,
		"bomber":         BomberID,
		"destroyer":      DestroyerID,
		"deathstar":      DeathstarID,
		"smallcargo":     SmallCargoID,
		"largecargo":     LargeCargoID,
		"colonyship":     ColonyShipID,
		"recycler":       RecyclerID,
		"espionageprobe": EspionageProbeID,
		"solarsatellite": SolarSatelliteID,
		"crawler":        CrawlerID,
		"reaper":         ReaperID,
		"pathfinder":     PathfinderID,

		// si
		"laheklovec":          LightFighterID,
		"tezkilovec":          HeavyFighterID,
		"krizarka":            CruiserID,
		"bojnaladja":          BattleshipID,
		"bojnakrizarka":       BattlecruiserID,
		"bombnik":             BomberID,
		"unicevalec":          DestroyerID,
		"zvezdasmrti":         DeathstarID,
		"majhnatovornaladja":  SmallCargoID,
		"velikatovornaladja":  LargeCargoID,
		"kolonizacijskaladja": ColonyShipID,
		"vohunskasonda":       EspionageProbeID,
		"soncnisatelit":       SolarSatelliteID,
		"plazilec":            CrawlerID,
		"kombajn":             ReaperID,
		"iskalecsledi":        PathfinderID,

		// ro
		"vanatorusor":      LightFighterID,
		"vanatorgreu":      HeavyFighterID,
		"crucisator":       CruiserID,
		"navaderazboi":     BattleshipID,
		"distrugator":      DestroyerID,
		"rip":              DeathstarID,
		"transportormic":   SmallCargoID,
		"transportormare":  LargeCargoID,
		"navadecolonizare": ColonyShipID,
		"reciclator":       RecyclerID,
		"probadespionaj":   EspionageProbeID,
		"satelitsolar":     SolarSatelliteID,

		// sk
		"lahkystihac":    LightFighterID,
		"tazkystihac":    HeavyFighterID,
		"bojovalod":      BattleshipID,
		"bojovykriznik":  BattlecruiserID,
		"devastator":     DestroyerID,
		"hviezdasmrti":   DeathstarID,
		"kolonizacnalod": ColonyShipID,
		"spionaznasonda": EspionageProbeID,
		"solarnysatelit": SolarSatelliteID,
		"vrtak":          CrawlerID,
		"kosa":           ReaperID,
		"prieskumnik":    PathfinderID,

		// gr
		"ελαφρυμαχητικο":         LightFighterID,
		"βαρυμαχητικο":           HeavyFighterID,
		"καταδιωκτικο":           CruiserID,
		"καταδρομικο":            BattleshipID,
		"θωρηκτοαναχαιτισης":     BattlecruiserID,
		"βομβαρδιστικο":          BomberID,
		"μικρομεταγωγικο":        SmallCargoID,
		"μεγαλομεταγωγικο":       LargeCargoID,
		"σκαφοςαποικιοποιησης":   ColonyShipID,
		"ανακυκλωτης":            RecyclerID,
		"κατασκοπευτικοστελεχος": EspionageProbeID,
		"ηλιακοισυλλεκτες":       SolarSatelliteID,

		// no
		"lettjeger":      LightFighterID,
		"tungjeger":      HeavyFighterID,
		"krysser":        CruiserID,
		"slagskip":       BattleshipID,
		"slagkrysser":    BattlecruiserID,
		"litelasteskip":  SmallCargoID,
		"stortlasteskip": LargeCargoID,
		"koloniskip":     ColonyShipID,
		"resirkulerer":   RecyclerID,
		"spionasjesonde": EspionageProbeID,
		"solarsatelitt":  SolarSatelliteID,

		// tw
		"輕型戰鬥機": LightFighterID,
		"重型戰鬥機": HeavyFighterID,
		"戰列艦":   BattleshipID,
		"戰鬥巡洋艦": BattlecruiserID,
		"導彈艦":   BomberID,
		"毀滅者":   DestroyerID,
		"死星":    DeathstarID,
		"小型運輸艦": SmallCargoID,
		"大型運輸艦": LargeCargoID,
		"殖民船":   ColonyShipID,
		"回收船":   RecyclerID,
		"間諜衛星":  EspionageProbeID,
		"太陽能衛星": SolarSatelliteID,
		"履帶車":   CrawlerID,
		"惡魔飛船":  ReaperID,
		"探路者":   PathfinderID,

		// hr
		"malilovac":         LightFighterID,
		"velikilovac":       HeavyFighterID,
		"krstarica":         CruiserID,
		"borbenibrod":       BattleshipID,
		"oklopnakrstarica":  BattlecruiserID,
		"razarac":           DestroyerID,
		"zvijezdasmrti":     DeathstarID,
		"malitransporter":   SmallCargoID,
		"velikitransporter": LargeCargoID,
		"kolonijalnibrod":   ColonyShipID,
		"recikler":          RecyclerID,
		"sondezaspijunazu":  EspionageProbeID,
		"puzavac":           CrawlerID,
		"zetelac":           ReaperID,
		"krcilac":           PathfinderID,

		// mx
		"navedelacolonia": ColonyShipID,
		"taladrador":      CrawlerID,
		"segador":         ReaperID,
		"explorador":      PathfinderID,

		// cz
		"lehkystihac":      LightFighterID,
		"tezkystihac":      HeavyFighterID,
		"kriznik":          CruiserID,
		"bitevnilod":       BattleshipID,
		"bitevnikriznik":   BattlecruiserID,
		"bombarder":        BomberID,
		"nicitel":          DestroyerID,
		"hvezdasmrti":      DeathstarID,
		"malytransporter":  SmallCargoID,
		"velkytransporter": LargeCargoID,
		"kolonizacnilod":   ColonyShipID,
		"recyklator":       RecyclerID,
		"spionaznisonda":   EspionageProbeID,
		"solarnisatelit":   SolarSatelliteID,
		"rozparovac":       ReaperID,
		"pruzkumnik":       PathfinderID,

		// it
		"caccialeggero":           LightFighterID,
		"cacciapesante":           HeavyFighterID,
		"incrociatore":            CruiserID,
		"navedabattaglia":         BattleshipID,
		"incrociatoredabattaglia": BattlecruiserID,
		"bombardiere":             BomberID,
		"corazzata":               DestroyerID,
		"mortenera":               DeathstarID,
		"cargoleggero":            SmallCargoID,
		"cargopesante":            LargeCargoID,
		"colonizzatrice":          ColonyShipID,
		"riciclatrici":            RecyclerID,
		"sondaspia":               EspionageProbeID,
		"satellitesolare":         SolarSatelliteID,

		// de
		"leichterjager":      LightFighterID,
		"schwererjager":      HeavyFighterID,
		"kreuzer":            CruiserID,
		"schlachtschiff":     BattleshipID,
		"schlachtkreuzer":    BattlecruiserID,
		"zerstorer":          DestroyerID,
		"todesstern":         DeathstarID,
		"kleinertransporter": SmallCargoID,
		"groertransporter":   LargeCargoID,
		"kolonieschiff":      ColonyShipID,
		"spionagesonde":      EspionageProbeID,
		"solarsatellit":      SolarSatelliteID,
		// "bomber":             BomberID,
		// "recycler":           RecyclerID,

		// dk
		"kravler":   CrawlerID,
		"stifinder": PathfinderID,

		// es
		"cazadorligero":      LightFighterID,
		"cazadorpesado":      HeavyFighterID,
		"crucero":            CruiserID,
		"navedebatalla":      BattleshipID,
		"acorazado":          BattlecruiserID,
		"bombardero":         BomberID,
		"destructor":         DestroyerID,
		"estrelladelamuerte": DeathstarID,
		"navepequenadecarga": SmallCargoID,
		"navegrandedecarga":  LargeCargoID,
		"colonizador":        ColonyShipID,
		"reciclador":         RecyclerID,
		"sondadeespionaje":   EspionageProbeID,
		"satelitesolar":      SolarSatelliteID,

		// fr
		"chasseurleger":          LightFighterID,
		"chasseurlourd":          HeavyFighterID,
		"croiseur":               CruiserID,
		"vaisseaudebataille":     BattleshipID,
		"traqueur":               BattlecruiserID,
		"bombardier":             BomberID,
		"destructeur":            DestroyerID,
		"etoiledelamort":         DeathstarID,
		"petittransporteur":      SmallCargoID,
		"grandtransporteur":      LargeCargoID,
		"vaisseaudecolonisation": ColonyShipID,
		"recycleur":              RecyclerID,
		"sondedespionnage":       EspionageProbeID,
		"satellitesolaire":       SolarSatelliteID,
		"foreuse":                CrawlerID,
		"faucheur":               ReaperID,
		"eclaireur":              PathfinderID,

		// br
		"cacaligeiro":       LightFighterID,
		"cacapesado":        HeavyFighterID,
		"cruzador":          CruiserID,
		"navedebatalha":     BattleshipID,
		"interceptador":     BattlecruiserID,
		"bombardeiro":       BomberID,
		"destruidor":        DestroyerID,
		"estreladamorte":    DeathstarID,
		"cargueiropequeno":  SmallCargoID,
		"cargueirogrande":   LargeCargoID,
		"navecolonizadora":  ColonyShipID,
		"sondadeespionagem": EspionageProbeID,
		//"reciclador":        RecyclerID,
		//"satelitesolar":     SolarSatelliteID,

		// jp
		"軽戦闘機":      LightFighterID,
		"重戦闘機":      HeavyFighterID,
		"巡洋艦":       CruiserID,
		"トルシッ":      BattleshipID,
		"大型戦艦":      BattlecruiserID,
		"爆撃機":       BomberID,
		"テストロイヤー":   DestroyerID,
		"テススター":     DeathstarID,
		"小型輸送機":     SmallCargoID,
		"大型輸送機":     LargeCargoID,
		"コロニーシッ":    ColonyShipID,
		"残骸回収船":     RecyclerID,
		"偵察機":       EspionageProbeID,
		"ソーラーサテライト": SolarSatelliteID,
		"ローラー":      CrawlerID,
		"ーー":        ReaperID,
		"スイター":      PathfinderID,

		// pl
		"lekkimysliwiec":      LightFighterID,
		"ciezkimysliwiec":     HeavyFighterID,
		"krazownik":           CruiserID,
		"okretwojenny":        BattleshipID,
		"pancernik":           BattlecruiserID,
		"bombowiec":           BomberID,
		"niszczyciel":         DestroyerID,
		"gwiazdasmierci":      DeathstarID,
		"maytransporter":      SmallCargoID,
		"duzytransporter":     LargeCargoID,
		"statekkolonizacyjny": ColonyShipID,
		"recykler":            RecyclerID,
		"sondaszpiegowska":    EspionageProbeID,
		"satelitasoneczny":    SolarSatelliteID,
		"pezacz":              CrawlerID,
		"rozpruwacz":          ReaperID,
		"pionier":             PathfinderID,

		// tr
		"hafifavc":           LightFighterID,
		"agravc":             HeavyFighterID,
		"kruvazoradet":       CruiserID,
		"komutagemisi":       BattleshipID,
		"firkateyn":          BattlecruiserID,
		"bombardmangemisi":   BomberID,
		"muhrip":             DestroyerID,
		"olumyildizi":        DeathstarID,
		"kucuknakliyegemisi": SmallCargoID,
		"buyuknakliyegemisi": LargeCargoID,
		"kolonigemisi":       ColonyShipID,
		"geridonusumcu":      RecyclerID,
		"casussondasi":       EspionageProbeID,
		"solaruydu":          SolarSatelliteID,
		"paletli":            CrawlerID,
		"azrail":             ReaperID,
		"rehber":             PathfinderID,

		// pt
		"interceptor":       BattlecruiserID,
		"navedecolonizacao": ColonyShipID,
		"rastejador":        CrawlerID,
		"ceifeira":          ReaperID,
		"exploradora":       PathfinderID,

		// nl
		"lichtgevechtsschip": LightFighterID,
		"zwaargevechtsschip": HeavyFighterID,
		"kruiser":            CruiserID,
		"slagschip":          BattleshipID,
		//"interceptor":          BattlecruiserID,
		"bommenwerper":     BomberID,
		"vernietiger":      DestroyerID,
		"sterdesdoods":     DeathstarID,
		"kleinvrachtschip": SmallCargoID,
		"grootvrachtschip": LargeCargoID,
		"kolonisatieschip": ColonyShipID,
		//"recycler":      RecyclerID,
		//"spionagesonde":       EspionageProbeID,
		"zonneenergiesatelliet": SolarSatelliteID,
		"kruiper":               CrawlerID,
		"hein":                  ReaperID,
		"padvinder":             PathfinderID,

		//dk
		"lillejger": LightFighterID,
		"storjger":  HeavyFighterID,
		"krydser":   CruiserID,
		"slagskib":  BattleshipID,
		//"interceptor":      BattlecruiserID,
		//"bomber":           BomberID,
		//"destroyer":        DestroyerID,
		"ddsstjerne":       DeathstarID,
		"lilletransporter": SmallCargoID,
		"stortransporter":  LargeCargoID,
		"koloniskib":       ColonyShipID,
		//"recycler":         RecyclerID,
		//"spionagesonde":    EspionageProbeID,
		//"solarsatellit":    SolarSatelliteID,

		// ru
		"легкииистребитель":  LightFighterID,
		"тяжелыиистребитель": HeavyFighterID,
		"креисер":            CruiserID,
		"линкор":             BattleshipID,
		"линеиныикреисер":    BattlecruiserID,
		"бомбардировщик":     BomberID,
		"уничтожитель":       DestroyerID,
		"звездасмерти":       DeathstarID,
		"малыитранспорт":     SmallCargoID,
		"большоитранспорт":   LargeCargoID,
		"колонизатор":        ColonyShipID,
		"переработчик":       RecyclerID,
		"шпионскиизонд":      EspionageProbeID,
		"солнечныиспутник":   SolarSatelliteID,
		"гусеничник":         CrawlerID,
		"жнец":               ReaperID,
		"первопроходец":      PathfinderID,
	}
	return nameMap[processedString]
}

// MinInt returns the minimum int64 value
func MinInt(vals ...int64) int64 {
	min := vals[0]
	for _, num := range vals {
		if num < min {
			min = num
		}
	}
	return min
}

// MaxInt returns the minimum int64 value
func MaxInt(vals ...int64) int64 {
	max := vals[0]
	for _, num := range vals {
		if num > max {
			max = num
		}
	}
	return max
}
