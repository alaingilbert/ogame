package ogame

import (
	"errors"
	"github.com/alaingilbert/ogame/pkg/utils"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// ParseCoord parse a coordinate from a string
func ParseCoord(str string) (coord Coordinate, err error) {
	m := regexp.MustCompile(`^\[?(([PMD]):)?(\d{1,3}):(\d{1,3}):(\d{1,3})]?$`).FindStringSubmatch(str)
	if len(m) == 6 {
		planetTypeStr := m[2]
		galaxy := utils.DoParseI64(m[3])
		system := utils.DoParseI64(m[4])
		position := utils.DoParseI64(m[5])
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

// DoParseCoord parse a coordinate from a string and ignore errors
func DoParseCoord(str string) (coord Coordinate) {
	coord, _ = ParseCoord(str)
	return coord
}

// MustParseCoord parse a coordinate from a string and panic if there is an error
func MustParseCoord(str string) Coordinate {
	coord, err := ParseCoord(str)
	if err != nil {
		panic(err)
	}
	return coord
}

var namesChars = "ЁАБВГДЕЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯабвгдежзийклмнопрстуфхцчшщъыьэюяёァイウオガキケコサザシスズソタダチッテデトドニノバパビフプヘマミムャヤラルレロンー偵列加反収器回型塔大太子察射導小履巡帶弾彈惡戦戰抗探撃收星機死残殖毀民洋滅漿炮爆發砲磁罩者能船艦衛諜護路車軌軽輕輸農送運道重間闘防陽際離雷電飛骸鬥魔"
var namesRgx = regexp.MustCompile("[^a-zA-Zα-ωΑ-Ω" + namesChars + "]+")

func unique(s string) string {
	//s = strings.ToLower(s)
	m := make(map[rune]struct{})
	for _, c := range s {
		m[c] = struct{}{}
	}
	arr := make([]string, 0)
	for k := range m {
		arr = append(arr, string(k))
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	return strings.Join(arr, "")
}

// DefenceName2ID ...
func DefenceName2ID(name string) ID {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	name, _, _ = transform.String(t, name)
	processedString := strings.ToLower(namesRgx.ReplaceAllString(name, ""))
	nameMap := map[string]ID{
		// en
		"rocketlauncher":         RocketLauncherID,
		"lightlaser":             LightLaserID,
		"heavylaser":             HeavyLaserID,
		"gausscannon":            GaussCannonID,
		"ioncannon":              IonCannonID,
		"plasmaturret":           PlasmaTurretID,
		"smallshielddome":        SmallShieldDomeID,
		"largeshielddome":        LargeShieldDomeID,
		"antiballisticmissiles":  AntiBallisticMissilesID,
		"interplanetarymissiles": InterplanetaryMissilesID,

		// hu
		"raketakilovo":            RocketLauncherID,
		"konnyulezer":             LightLaserID,
		"nehezlezer":              HeavyLaserID,
		"gaussagyu":               GaussCannonID,
		"ionagyu":                 IonCannonID,
		"plazmatorony":            PlasmaTurretID,
		"kispajzskupola":          SmallShieldDomeID,
		"nagypajzskupola":         LargeShieldDomeID,
		"antiballasztikusraketak": AntiBallisticMissilesID,
		"bolygokoziraketak":       InterplanetaryMissilesID,

		// si
		"raketnik":              RocketLauncherID,
		"lahkilaser":            LightLaserID,
		"tezeklaser":            HeavyLaserID,
		"gaussovtop":            GaussCannonID,
		"ionskitop":             IonCannonID,
		"plazemskitop":          PlasmaTurretID,
		"majhenscit":            SmallShieldDomeID,
		"velikscit":             LargeShieldDomeID,
		"protibalisticnerakete": AntiBallisticMissilesID,
		"medplanetarnerakete":   InterplanetaryMissilesID,

		// ro
		"lansatorderachete":     RocketLauncherID,
		"laserusor":             LightLaserID,
		"lasergreu":             HeavyLaserID,
		"tungauss":              GaussCannonID,
		"tunmagnetic":           IonCannonID,
		"turelacuplasma":        PlasmaTurretID,
		"scutplanetarmic":       SmallShieldDomeID,
		"scutplanetarmare":      LargeShieldDomeID,
		"rachetaantibalistica":  AntiBallisticMissilesID,
		"racheteinterplanetare": InterplanetaryMissilesID,

		// sk
		"raketovykomplet":       RocketLauncherID,
		"lahkylaser":            LightLaserID,
		"tazkylaser":            HeavyLaserID,
		"gaussovkanon":          GaussCannonID,
		"ionovykanon":           IonCannonID,
		"plazmovaveza":          PlasmaTurretID,
		"malyplanetarnystit":    SmallShieldDomeID,
		"velkyplanetarnystit":   LargeShieldDomeID,
		"protiraketovestrely":   AntiBallisticMissilesID,
		"medziplanetarnerakety": InterplanetaryMissilesID,

		// gr
		"εκτοξευτηςπυραυλων":      RocketLauncherID,
		"ελαφρυλειζερ":            LightLaserID,
		"βαρυλειζερ":              HeavyLaserID,
		"κανονιgauss":             GaussCannonID,
		"κανονιιοντων":            IonCannonID,
		"πυργισκοιπλασματος":      PlasmaTurretID,
		"μικροςαμυντικοςθολος":    SmallShieldDomeID,
		"μεγαλοςαμυντικοςθολος":   LargeShieldDomeID,
		"αντιβαλλιστικοιπυραυλοι": AntiBallisticMissilesID,
		"διαπλανητικοιπυραυλοι":   InterplanetaryMissilesID,

		// tw
		"飛彈發射器": RocketLauncherID,
		"輕型雷射炮": LightLaserID,
		"重型雷射炮": HeavyLaserID,
		"磁軌炮":   GaussCannonID,
		"離子加農炮": IonCannonID,
		"電漿炮塔":  PlasmaTurretID,
		"小型防護罩": SmallShieldDomeID,
		"大型防護罩": LargeShieldDomeID,
		"反彈道導彈": AntiBallisticMissilesID,
		"星際導彈":  InterplanetaryMissilesID,

		// hr
		"raketobacaci":          RocketLauncherID,
		"malilaser":             LightLaserID,
		"velikilaser":           HeavyLaserID,
		"plazmatop":             PlasmaTurretID,
		"malastitnakupola":      SmallShieldDomeID,
		"velikastitnakupola":    LargeShieldDomeID,
		"antibalistickerakete":  AntiBallisticMissilesID,
		"interplanetarnerakete": InterplanetaryMissilesID,

		// mx
		"lanzamisiles":              RocketLauncherID,
		"laserpequeno":              LightLaserID,
		"lasergrande":               HeavyLaserID,
		"canongauss":                GaussCannonID,
		"canonionico":               IonCannonID,
		"canondeplasma":             PlasmaTurretID,
		"cupulapequenadeproteccion": SmallShieldDomeID,
		"cupulagrandedeproteccion":  LargeShieldDomeID,
		"misildeintercepcion":       AntiBallisticMissilesID,
		"misilinterplanetario":      InterplanetaryMissilesID,

		// cz
		"raketomet":            RocketLauncherID,
		"lehkylaser":           LightLaserID,
		"tezkylaser":           HeavyLaserID,
		"gaussuvkanon":         GaussCannonID,
		"iontovykanon":         IonCannonID,
		"plasmovavez":          PlasmaTurretID,
		"malyplanetarnistit":   SmallShieldDomeID,
		"velkyplanetarnistit":  LargeShieldDomeID,
		"antibalistickerakety": AntiBallisticMissilesID,
		"meziplanetarnirakety": InterplanetaryMissilesID,

		// it
		"lanciamissili":         RocketLauncherID,
		"laserleggero":          LightLaserID,
		"laserpesante":          HeavyLaserID,
		"cannonegauss":          GaussCannonID,
		"cannoneionico":         IonCannonID,
		"cannonealplasma":       PlasmaTurretID,
		"cupolascudopiccola":    SmallShieldDomeID,
		"cupolascudopotenziata": LargeShieldDomeID,
		"missiliantibalistici":  AntiBallisticMissilesID,
		"missiliinterplanetari": InterplanetaryMissilesID,

		// de
		"raketenwerfer":         RocketLauncherID,
		"leichteslasergeschutz": LightLaserID,
		"schwereslasergeschutz": HeavyLaserID,
		"gaukanone":             GaussCannonID,
		"ionengeschutz":         IonCannonID,
		"plasmawerfer":          PlasmaTurretID,
		"kleineschildkuppel":    SmallShieldDomeID,
		"groeschildkuppel":      LargeShieldDomeID,
		"abfangrakete":          AntiBallisticMissilesID,
		"interplanetarrakete":   InterplanetaryMissilesID,

		// dk
		"raketkanon":         RocketLauncherID,
		"lillelaserkanon":    LightLaserID,
		"storlaserkanon":     HeavyLaserID,
		"gausskanon":         GaussCannonID,
		"ionkanon":           IonCannonID,
		"plasmakanon":        PlasmaTurretID,
		"lilleplanetskjold":  SmallShieldDomeID,
		"stortplanetskjold":  LargeShieldDomeID,
		"forsvarsraket":      AntiBallisticMissilesID,
		"interplanetarraket": InterplanetaryMissilesID,

		// es
		"misilesantibalisticos": AntiBallisticMissilesID,

		// fr
		"lanceurdemissiles":      RocketLauncherID,
		"artillerielaserlegere":  LightLaserID,
		"artillerielaserlourde":  HeavyLaserID,
		"canondegauss":           GaussCannonID,
		"artillerieaions":        IonCannonID,
		"lanceurdeplasma":        PlasmaTurretID,
		"petitbouclier":          SmallShieldDomeID,
		"grandbouclier":          LargeShieldDomeID,
		"missiledinterception":   AntiBallisticMissilesID,
		"missileinterplanetaire": InterplanetaryMissilesID,

		// br
		"lancadordemisseis":       RocketLauncherID,
		"laserligeiro":            LightLaserID,
		"laserpesado":             HeavyLaserID,
		"canhaodegauss":           GaussCannonID,
		"canhaodeions":            IonCannonID,
		"canhaodeplasma":          PlasmaTurretID,
		"pequenoescudoplanetario": SmallShieldDomeID,
		"grandeescudoplanetario":  LargeShieldDomeID,
		"missildeinterceptacao":   AntiBallisticMissilesID,
		"missilinterplanetario":   InterplanetaryMissilesID,

		// jp
		"ロケットランチャー": RocketLauncherID,
		"ライトレーサー":   LightLaserID,
		"ヘーレーサー":    HeavyLaserID,
		"ウスキャノン":    GaussCannonID,
		"イオンキャノン":   IonCannonID,
		"フラスマ砲":     PlasmaTurretID,
		"小型シールトトーム": SmallShieldDomeID,
		"大型シールトトーム": LargeShieldDomeID,
		"抗弾道ミサイル":   AntiBallisticMissilesID,
		"星間ミサイル":    InterplanetaryMissilesID,

		// pl
		"wyrzutniarakiet":         RocketLauncherID,
		"lekkiedziaolaserowe":     LightLaserID,
		"ciezkiedziaolaserowe":    HeavyLaserID,
		"dziaogaussa":             GaussCannonID,
		"dziaojonowe":             IonCannonID,
		"wyrzutniaplazmy":         PlasmaTurretID,
		"maaosonaochronna":        SmallShieldDomeID,
		"duzaosonaochronna":       LargeShieldDomeID,
		"przeciwrakieta":          AntiBallisticMissilesID,
		"rakietamiedzyplanetarna": InterplanetaryMissilesID,

		// tr
		"roketatar":               RocketLauncherID,
		"hafiflazertopu":          LightLaserID,
		"agrlazertopu":            HeavyLaserID,
		"gaustopu":                GaussCannonID,
		"iyontopu":                IonCannonID,
		"plazmaatc":               PlasmaTurretID,
		"kucukkalkankubbesi":      SmallShieldDomeID,
		"buyukkalkankubbesi":      LargeShieldDomeID,
		"yakalycroketler":         AntiBallisticMissilesID,
		"gezegenlerarasiroketler": InterplanetaryMissilesID,

		// pt
		"canhaodeioes":        IonCannonID,
		"missildeintercepcao": AntiBallisticMissilesID,

		// nl
		"raketlanceerder":              RocketLauncherID,
		"kleinelaser":                  LightLaserID,
		"grotelaser":                   HeavyLaserID,
		"kleineplanetaireschildkoepel": SmallShieldDomeID,
		"groteplanetaireschildkoepel":  LargeShieldDomeID,
		"antiballistischeraketten":     AntiBallisticMissilesID,
		"interplanetaireraketten":      InterplanetaryMissilesID,

		// ru
		"ракетнаяустановка":   RocketLauncherID,
		"легкиилазер":         LightLaserID,
		"тяжелыилазер":        HeavyLaserID,
		"пушкагаусса":         GaussCannonID,
		"ионноеорудие":        IonCannonID,
		"плазменноеорудие":    PlasmaTurretID,
		"малыищитовоикупол":   SmallShieldDomeID,
		"большоищитовоикупол": LargeShieldDomeID,
		"ракетаперехватчик":   AntiBallisticMissilesID,
		"межпланетнаяракета":  InterplanetaryMissilesID,

		// fi --
		// no --
		// ba --
	}
	return nameMap[processedString]
}

// ShipName2ID ...
func ShipName2ID(name string) ID {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	name, _, _ = transform.String(t, name)
	processedString := strings.ToLower(namesRgx.ReplaceAllString(name, ""))
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

		// fi
		"kevythavittaja":     LightFighterID,
		"raskashavittaja":    HeavyFighterID,
		"risteilija":         CruiserID,
		"taistelualus":       BattleshipID,
		"taisteluristeilija": BattlecruiserID,
		"pommittaja":         BomberID,
		"tuhoaja":            DestroyerID,
		"kuolemantahti":      DeathstarID,
		"pienirahtialus":     SmallCargoID,
		"suurirahtialus":     LargeCargoID,
		"siirtokuntaalus":    ColonyShipID,
		"kierrattaja":        RecyclerID,
		"vakoiluluotain":     EspionageProbeID,
		"aurinkosatelliitti": SolarSatelliteID,

		// hu
		"konnyuharcos": LightFighterID,
		"nehezharcos":  HeavyFighterID,
		"cirkalo":      CruiserID,
		"csatahajo":    BattleshipID,
		"csatacirkalo": BattlecruiserID,
		"bombazo":      BomberID,
		"rombolo":      DestroyerID,
		"halalcsillag": DeathstarID,
		"kisszallito":  SmallCargoID,
		"nagyszallito": LargeCargoID,
		"koloniahajo":  ColonyShipID,
		"szemetesek":   RecyclerID,
		"kemszonda":    EspionageProbeID,
		"napmuhold":    SolarSatelliteID,
		"kaszas":       ReaperID,
		"felderito":    PathfinderID,

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

		// ba
		"krstarice":      CruiserID,
		"borbenibrodovi": BattleshipID,
		"razaraci":       DestroyerID,

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
		"トルシッフ":     BattleshipID,
		"大型戦艦":      BattlecruiserID,
		"爆撃機":       BomberID,
		"テストロイヤー":   DestroyerID,
		"テススター":     DeathstarID,
		"小型輸送機":     SmallCargoID,
		"大型輸送機":     LargeCargoID,
		"コロニーシッフ":   ColonyShipID,
		"残骸回収船":     RecyclerID,
		"偵察機":       EspionageProbeID,
		"ソーラーサテライト": SolarSatelliteID,
		"ローラー":      CrawlerID,
		"ーー":        ReaperID,
		"スファインター":   PathfinderID,

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
		"kruvazor":           CruiserID,
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
		"processer":             CrawlerID,
		"ruimer":                ReaperID,
		"navigator":             PathfinderID,

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
