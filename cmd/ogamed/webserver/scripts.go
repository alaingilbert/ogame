//package main

package webserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	//badger "github.com/dgraph-io/badger/v3"

	"github.com/faunX/ogame"
	"github.com/faunX/ogame/cmd/ogamed/ogb"
	"github.com/faunX/ogame/cmd/ogamed/skv"
	"github.com/labstack/echo"
	"github.com/mattn/anko/env"
	"github.com/mattn/anko/vm"
)

var WebKVStore *skv.KVStore

var OnAttackCh = make(chan ogame.AttackEvent, 10)
var OnAttackFn map[string]func(ogame.AttackEvent) = map[string]func(ogame.AttackEvent){}

func HandleAttackCh() {
	for {
		newAttack := <-OnAttackCh
		log.Println("New Attack Detected in handleAttackCh()")

		for k, f := range OnAttackFn {
			f(newAttack)
			delete(OnAttackFn, k)
		}

	}
}

func registerOnAttackCh(scriptName string, fn func(ogame.AttackEvent)) {
	OnAttackFn[scriptName] = fn
}

func unregisterOnAttackCh(scriptName string) {
	delete(OnAttackFn, scriptName)
}

//func RunScript(bot *ogame.OGame, badgerDB *badger.DB, ctx context.Context, scriptName, code string) error {
func RunScript(bot *ogame.OGame, ogb *ogb.Ogb, ctx context.Context, scriptName, code string) error {
	e := env.NewEnv()

	var onAttackCh = make(chan ogame.AttackEvent)

	onAttackFn := func(a ogame.AttackEvent) {
		log.Println("New Attack")
		onAttackCh <- a
	}
	registerOnAttackCh(scriptName, onAttackFn)

	err := e.Define("println", fmt.Println)
	if err != nil {
		log.Fatalf("Define error: %v\n", err)
	}

	env.Packages["sync"] = map[string]reflect.Value{
		"NewCond": reflect.ValueOf(sync.NewCond),
	}
	env.PackageTypes["sync"] = map[string]reflect.Type{
		"Cond":      reflect.TypeOf(sync.Cond{}),
		"Mutex":     reflect.TypeOf(sync.Mutex{}),
		"Once":      reflect.TypeOf(sync.Once{}),
		"Pool":      reflect.TypeOf(sync.Pool{}),
		"RWMutex":   reflect.TypeOf(sync.RWMutex{}),
		"WaitGroup": reflect.TypeOf(sync.WaitGroup{}),
	}

	env.Packages["sort"] = map[string]reflect.Value{
		"Float64s":          reflect.ValueOf(sort.Float64s),
		"Float64sAreSorted": reflect.ValueOf(sort.Float64sAreSorted),
		"Ints":              reflect.ValueOf(sort.Ints),
		"IntsAreSorted":     reflect.ValueOf(sort.IntsAreSorted),
		"IsSorted":          reflect.ValueOf(sort.IsSorted),
		"Search":            reflect.ValueOf(sort.Search),
		"SearchFloat64s":    reflect.ValueOf(sort.SearchFloat64s),
		"SearchInts":        reflect.ValueOf(sort.SearchInts),
		"SearchStrings":     reflect.ValueOf(sort.SearchStrings),
		"Sort":              reflect.ValueOf(sort.Sort),
		"Stable":            reflect.ValueOf(sort.Stable),
		"Strings":           reflect.ValueOf(sort.Strings),
		"StringsAreSorted":  reflect.ValueOf(sort.StringsAreSorted),
	}
	env.PackageTypes["sort"] = map[string]reflect.Type{
		"Float64Slice":    reflect.TypeOf(sort.Float64Slice{}),
		"IntSlice":        reflect.TypeOf(sort.IntSlice{}),
		"StringSlice":     reflect.TypeOf(sort.StringSlice{}),
		"SortFuncsStruct": reflect.TypeOf(&SortFuncsStruct{}),
	}

	env.Packages["regexp"] = map[string]reflect.Value{
		"Match":            reflect.ValueOf(regexp.Match),
		"MatchReader":      reflect.ValueOf(regexp.MatchReader),
		"MatchString":      reflect.ValueOf(regexp.MatchString),
		"QuoteMeta":        reflect.ValueOf(regexp.QuoteMeta),
		"Compile":          reflect.ValueOf(regexp.Compile),
		"CompilePOSIX":     reflect.ValueOf(regexp.CompilePOSIX),
		"MustCompile":      reflect.ValueOf(regexp.MustCompile),
		"MustCompilePOSIX": reflect.ValueOf(regexp.MustCompilePOSIX),
	}

	env.Packages["encoding/json"] = map[string]reflect.Value{
		"Marshal":   reflect.ValueOf(json.Marshal),
		"Unmarshal": reflect.ValueOf(json.Unmarshal),
	}

	env.Packages["fmt"] = map[string]reflect.Value{
		"Errorf":   reflect.ValueOf(fmt.Errorf),
		"Fprint":   reflect.ValueOf(fmt.Fprint),
		"Fprintf":  reflect.ValueOf(fmt.Fprintf),
		"Fprintln": reflect.ValueOf(fmt.Fprintln),
		"Fscan":    reflect.ValueOf(fmt.Fscan),
		"Fscanf":   reflect.ValueOf(fmt.Fscanf),
		"Fscanln":  reflect.ValueOf(fmt.Fscanln),
		"Print":    reflect.ValueOf(fmt.Print),
		"Printf":   reflect.ValueOf(fmt.Printf),
		"Println":  reflect.ValueOf(fmt.Println),
		"Scan":     reflect.ValueOf(fmt.Scan),
		"Scanf":    reflect.ValueOf(fmt.Scanf),
		"Scanln":   reflect.ValueOf(fmt.Scanln),
		"Sprint":   reflect.ValueOf(fmt.Sprint),
		"Sprintf":  reflect.ValueOf(fmt.Sprintf),
		"Sprintln": reflect.ValueOf(fmt.Sprintln),
		"Sscan":    reflect.ValueOf(fmt.Sscan),
		"Sscanf":   reflect.ValueOf(fmt.Sscanf),
		"Sscanln":  reflect.ValueOf(fmt.Sscanln),
	}

	env.Packages["strings"] = map[string]reflect.Value{
		"Contains":       reflect.ValueOf(strings.Contains),
		"ContainsAny":    reflect.ValueOf(strings.ContainsAny),
		"ContainsRune":   reflect.ValueOf(strings.ContainsRune),
		"Count":          reflect.ValueOf(strings.Count),
		"EqualFold":      reflect.ValueOf(strings.EqualFold),
		"Fields":         reflect.ValueOf(strings.Fields),
		"FieldsFunc":     reflect.ValueOf(strings.FieldsFunc),
		"HasPrefix":      reflect.ValueOf(strings.HasPrefix),
		"HasSuffix":      reflect.ValueOf(strings.HasSuffix),
		"Index":          reflect.ValueOf(strings.Index),
		"IndexAny":       reflect.ValueOf(strings.IndexAny),
		"IndexByte":      reflect.ValueOf(strings.IndexByte),
		"IndexFunc":      reflect.ValueOf(strings.IndexFunc),
		"IndexRune":      reflect.ValueOf(strings.IndexRune),
		"Join":           reflect.ValueOf(strings.Join),
		"LastIndex":      reflect.ValueOf(strings.LastIndex),
		"LastIndexAny":   reflect.ValueOf(strings.LastIndexAny),
		"LastIndexFunc":  reflect.ValueOf(strings.LastIndexFunc),
		"Map":            reflect.ValueOf(strings.Map),
		"NewReader":      reflect.ValueOf(strings.NewReader),
		"NewReplacer":    reflect.ValueOf(strings.NewReplacer),
		"Repeat":         reflect.ValueOf(strings.Repeat),
		"Replace":        reflect.ValueOf(strings.Replace),
		"Split":          reflect.ValueOf(strings.Split),
		"SplitAfter":     reflect.ValueOf(strings.SplitAfter),
		"SplitAfterN":    reflect.ValueOf(strings.SplitAfterN),
		"SplitN":         reflect.ValueOf(strings.SplitN),
		"Title":          reflect.ValueOf(strings.Title),
		"ToLower":        reflect.ValueOf(strings.ToLower),
		"ToLowerSpecial": reflect.ValueOf(strings.ToLowerSpecial),
		"ToTitle":        reflect.ValueOf(strings.ToTitle),
		"ToTitleSpecial": reflect.ValueOf(strings.ToTitleSpecial),
		"ToUpper":        reflect.ValueOf(strings.ToUpper),
		"ToUpperSpecial": reflect.ValueOf(strings.ToUpperSpecial),
		"Trim":           reflect.ValueOf(strings.Trim),
		"TrimFunc":       reflect.ValueOf(strings.TrimFunc),
		"TrimLeft":       reflect.ValueOf(strings.TrimLeft),
		"TrimLeftFunc":   reflect.ValueOf(strings.TrimLeftFunc),
		"TrimPrefix":     reflect.ValueOf(strings.TrimPrefix),
		"TrimRight":      reflect.ValueOf(strings.TrimRight),
		"TrimRightFunc":  reflect.ValueOf(strings.TrimRightFunc),
		"TrimSpace":      reflect.ValueOf(strings.TrimSpace),
		"TrimSuffix":     reflect.ValueOf(strings.TrimSuffix),
	}

	env.Packages["math"] = map[string]reflect.Value{
		"Abs":             reflect.ValueOf(math.Abs),
		"Acos":            reflect.ValueOf(math.Acos),
		"Acosh":           reflect.ValueOf(math.Acosh),
		"Asin":            reflect.ValueOf(math.Asin),
		"Asinh":           reflect.ValueOf(math.Asinh),
		"Atan":            reflect.ValueOf(math.Atan),
		"Atan2":           reflect.ValueOf(math.Atan2),
		"Atanh":           reflect.ValueOf(math.Atanh),
		"Cbrt":            reflect.ValueOf(math.Cbrt),
		"Ceil":            reflect.ValueOf(math.Ceil),
		"Copysign":        reflect.ValueOf(math.Copysign),
		"Cos":             reflect.ValueOf(math.Cos),
		"Cosh":            reflect.ValueOf(math.Cosh),
		"Dim":             reflect.ValueOf(math.Dim),
		"Erf":             reflect.ValueOf(math.Erf),
		"Erfc":            reflect.ValueOf(math.Erfc),
		"Exp":             reflect.ValueOf(math.Exp),
		"Exp2":            reflect.ValueOf(math.Exp2),
		"Expm1":           reflect.ValueOf(math.Expm1),
		"Float32bits":     reflect.ValueOf(math.Float32bits),
		"Float32frombits": reflect.ValueOf(math.Float32frombits),
		"Float64bits":     reflect.ValueOf(math.Float64bits),
		"Float64frombits": reflect.ValueOf(math.Float64frombits),
		"Floor":           reflect.ValueOf(math.Floor),
		"Frexp":           reflect.ValueOf(math.Frexp),
		"Gamma":           reflect.ValueOf(math.Gamma),
		"Hypot":           reflect.ValueOf(math.Hypot),
		"Ilogb":           reflect.ValueOf(math.Ilogb),
		"Inf":             reflect.ValueOf(math.Inf),
		"IsInf":           reflect.ValueOf(math.IsInf),
		"IsNaN":           reflect.ValueOf(math.IsNaN),
		"J0":              reflect.ValueOf(math.J0),
		"J1":              reflect.ValueOf(math.J1),
		"Jn":              reflect.ValueOf(math.Jn),
		"Ldexp":           reflect.ValueOf(math.Ldexp),
		"Lgamma":          reflect.ValueOf(math.Lgamma),
		"Log":             reflect.ValueOf(math.Log),
		"Log10":           reflect.ValueOf(math.Log10),
		"Log1p":           reflect.ValueOf(math.Log1p),
		"Log2":            reflect.ValueOf(math.Log2),
		"Logb":            reflect.ValueOf(math.Logb),
		"Max":             reflect.ValueOf(math.Max),
		"Min":             reflect.ValueOf(math.Min),
		"Mod":             reflect.ValueOf(math.Mod),
		"Modf":            reflect.ValueOf(math.Modf),
		"NaN":             reflect.ValueOf(math.NaN),
		"Nextafter":       reflect.ValueOf(math.Nextafter),
		"Pow":             reflect.ValueOf(math.Pow),
		"Pow10":           reflect.ValueOf(math.Pow10),
		"Remainder":       reflect.ValueOf(math.Remainder),
		"Signbit":         reflect.ValueOf(math.Signbit),
		"Sin":             reflect.ValueOf(math.Sin),
		"Sincos":          reflect.ValueOf(math.Sincos),
		"Sinh":            reflect.ValueOf(math.Sinh),
		"Sqrt":            reflect.ValueOf(math.Sqrt),
		"Tan":             reflect.ValueOf(math.Tan),
		"Tanh":            reflect.ValueOf(math.Tanh),
		"Trunc":           reflect.ValueOf(math.Trunc),
		"Y0":              reflect.ValueOf(math.Y0),
		"Y1":              reflect.ValueOf(math.Y1),
		"Yn":              reflect.ValueOf(math.Yn),
	}

	env.Packages["strconv"] = map[string]reflect.Value{
		"FormatBool":  reflect.ValueOf(strconv.FormatBool),
		"FormatFloat": reflect.ValueOf(strconv.FormatFloat),
		"FormatInt":   reflect.ValueOf(strconv.FormatInt),
		"FormatUint":  reflect.ValueOf(strconv.FormatUint),
		"ParseBool":   reflect.ValueOf(strconv.ParseBool),
		"ParseFloat":  reflect.ValueOf(strconv.ParseFloat),
		"ParseInt":    reflect.ValueOf(strconv.ParseInt),
		"ParseUint":   reflect.ValueOf(strconv.ParseUint),
	}

	env.Packages["time"] = map[string]reflect.Value{
		"ANSIC":           reflect.ValueOf(time.ANSIC),
		"After":           reflect.ValueOf(time.After),
		"AfterFunc":       reflect.ValueOf(time.AfterFunc),
		"April":           reflect.ValueOf(time.April),
		"August":          reflect.ValueOf(time.August),
		"Date":            reflect.ValueOf(time.Date),
		"December":        reflect.ValueOf(time.December),
		"February":        reflect.ValueOf(time.February),
		"FixedZone":       reflect.ValueOf(time.FixedZone),
		"Friday":          reflect.ValueOf(time.Friday),
		"Hour":            reflect.ValueOf(time.Hour),
		"January":         reflect.ValueOf(time.January),
		"July":            reflect.ValueOf(time.July),
		"June":            reflect.ValueOf(time.June),
		"Kitchen":         reflect.ValueOf(time.Kitchen),
		"LoadLocation":    reflect.ValueOf(time.LoadLocation),
		"March":           reflect.ValueOf(time.March),
		"May":             reflect.ValueOf(time.May),
		"Microsecond":     reflect.ValueOf(time.Microsecond),
		"Millisecond":     reflect.ValueOf(time.Millisecond),
		"Minute":          reflect.ValueOf(time.Minute),
		"Monday":          reflect.ValueOf(time.Monday),
		"Nanosecond":      reflect.ValueOf(time.Nanosecond),
		"NewTicker":       reflect.ValueOf(time.NewTicker),
		"NewTimer":        reflect.ValueOf(time.NewTimer),
		"November":        reflect.ValueOf(time.November),
		"Now":             reflect.ValueOf(time.Now),
		"October":         reflect.ValueOf(time.October),
		"Parse":           reflect.ValueOf(time.Parse),
		"ParseDuration":   reflect.ValueOf(time.ParseDuration),
		"ParseInLocation": reflect.ValueOf(time.ParseInLocation),
		"RFC1123":         reflect.ValueOf(time.RFC1123),
		"RFC1123Z":        reflect.ValueOf(time.RFC1123Z),
		"RFC3339":         reflect.ValueOf(time.RFC3339),
		"RFC3339Nano":     reflect.ValueOf(time.RFC3339Nano),
		"RFC822":          reflect.ValueOf(time.RFC822),
		"RFC822Z":         reflect.ValueOf(time.RFC822Z),
		"RFC850":          reflect.ValueOf(time.RFC850),
		"RubyDate":        reflect.ValueOf(time.RubyDate),
		"Saturday":        reflect.ValueOf(time.Saturday),
		"Second":          reflect.ValueOf(time.Second),
		"September":       reflect.ValueOf(time.September),
		"Since":           reflect.ValueOf(time.Since),
		"Sleep":           reflect.ValueOf(time.Sleep),
		"Stamp":           reflect.ValueOf(time.Stamp),
		"StampMicro":      reflect.ValueOf(time.StampMicro),
		"StampMilli":      reflect.ValueOf(time.StampMilli),
		"StampNano":       reflect.ValueOf(time.StampNano),
		"Sunday":          reflect.ValueOf(time.Sunday),
		"Thursday":        reflect.ValueOf(time.Thursday),
		"Tick":            reflect.ValueOf(time.Tick),
		"Tuesday":         reflect.ValueOf(time.Tuesday),
		"Unix":            reflect.ValueOf(time.Unix),
		"UnixDate":        reflect.ValueOf(time.UnixDate),
		"Until":           reflect.ValueOf(time.Until),
		"Wednesday":       reflect.ValueOf(time.Wednesday),
	}
	env.PackageTypes["time"] = map[string]reflect.Type{
		"Duration": reflect.TypeOf(time.Duration(0)),
		"Ticker":   reflect.TypeOf(time.Ticker{}),
		"Time":     reflect.TypeOf(time.Time{}),
	}

	// Missions
	e.Define("ATTACK", ogame.Attack)
	e.Define("GROUPEDATTACK", ogame.GroupedAttack)
	e.Define("TRANSPORT", ogame.Transport)
	e.Define("PARK", ogame.Park)
	e.Define("PARKINTHATALLY", ogame.ParkInThatAlly)
	e.Define("SPY", ogame.Spy)
	e.Define("COLONIZE", ogame.Colonize)
	e.Define("RECYCLEDEBRISFIELD", ogame.RecycleDebrisField)
	e.Define("DESTROY", ogame.Destroy)
	e.Define("MISSILEATTACK", ogame.MissileAttack)
	e.Define("EXPEDITION", ogame.Expedition)

	// Speed
	e.Define("FIVE_PERCENT", ogame.FivePercent)
	e.Define("TEN_PERCENT", ogame.TenPercent)                  // 10
	e.Define("FIFTEEN_PERCENT", ogame.FifteenPercent)          // 15
	e.Define("TWENTY_PERCENT", ogame.TwentyPercent)            // 20
	e.Define("TWENTY_FIVE_PERCENT", ogame.TwentyFivePercent)   // 25
	e.Define("THIRTY_PERCENT", ogame.ThirtyPercent)            // 30
	e.Define("THIRTY_FIVE_PERCENT", ogame.ThirtyFivePercent)   // 35
	e.Define("FORTY_PERCENT", ogame.FourtyPercent)             // 40
	e.Define("FORTY_FIVE_PERCENT", ogame.FourtyFivePercent)    // 45
	e.Define("FIFTY_PERCENT", ogame.FiftyPercent)              // 50
	e.Define("FIFTY_FIVE_PERCENT", ogame.FiftyFivePercent)     // 55
	e.Define("SIXTY_PERCENT", ogame.SixtyPercent)              // 60
	e.Define("SIXTY_FIVE_PERCENT", ogame.SixtyFivePercent)     // 65
	e.Define("SEVENTY_PERCENT", ogame.SeventyPercent)          // 70
	e.Define("SEVENTY_FIVE_PERCENT", ogame.SeventyFivePercent) // 75
	e.Define("EIGHTY_PERCENT", ogame.EightyPercent)            // 80
	e.Define("EIGHTY_FIVE_PERCENT", ogame.EightyFivePercent)   // 85
	e.Define("NINETY_PERCENT", ogame.NinetyPercent)            // 90
	e.Define("NINETY_FIVE_PERCENT", ogame.NinetyFivePercent)   // 95
	e.Define("HUNDRED_PERCENT", ogame.HundredPercent)          // 100

	// Celestial types
	e.Define("PLANET_TYPE", ogame.PlanetType) // 1
	e.Define("DEBRIS_TYPE", ogame.DebrisType) // 2
	e.Define("MOON_TYPE", ogame.MoonType)     // 3

	// Ships
	e.Define("SMALLCARGO", ogame.SmallCargoID)
	e.Define("LARGECARGO", ogame.LargeCargoID)
	e.Define("LIGHTFIGHTER", ogame.LightFighterID)
	e.Define("HEAVYFIGHTER", ogame.HeavyFighterID)
	e.Define("CRUISER", ogame.CruiserID)
	e.Define("BATTLESHIP", ogame.BattleshipID)
	e.Define("COLONYSHIP", ogame.ColonyShipID)
	e.Define("RECYCLER", ogame.RecyclerID)
	e.Define("ESPIONAGEPROBE", ogame.EspionageProbeID)
	e.Define("BOMBER", ogame.BomberID)
	e.Define("DESTROYER", ogame.DestroyerID)
	e.Define("DEATHSTAR", ogame.DeathstarID)
	e.Define("BATTLECRUISER", ogame.BattlecruiserID)
	e.Define("CRAWLER", ogame.CrawlerID)
	e.Define("REAPER", ogame.ReaperID)
	e.Define("PATHFINDER", ogame.PathfinderID)

	// Defenses
	//, , , , , , , , ,
	e.Define("ROCKETLAUNCHER", ogame.RocketLauncherID)
	e.Define("LIGHTLASER", ogame.LightFighterID)
	e.Define("HEAVYLASER", ogame.HeavyLaserID)
	e.Define("GAUSSCANNON", ogame.GaussCannonID)
	e.Define("IONCANNON", ogame.IonCannonID)
	e.Define("PLASMATURRET", ogame.PlasmaTurretID)
	e.Define("SMALLSHIELDDOME", ogame.SmallShieldDomeID)
	e.Define("LARGESHIELDDOME", ogame.LargeShieldDomeID)
	e.Define("ANTIBALLISTICMISSILES", ogame.AntiBallisticMissilesID)
	e.Define("INTERPLANETARYMISSILES", ogame.InterplanetaryMissilesID)

	//e.Define("print", fmt.Println)
	e.Define("Login", bot.LoginWithExistingCookies)
	e.Define("Logout", bot.Logout)
	e.Define("Enable", bot.Enable)
	e.Define("IsEnabled", bot.IsEnabled)
	e.Define("IsLoggedIn", bot.IsLoggedIn)

	e.Define("Disable", bot.Disable)

	e.Define("Abandon", bot.Abandon)
	e.Define("GetFleets", bot.GetFleets)
	e.Define("GetAttacks", bot.GetAttacks)
	e.Define("GalaxyInfos", bot.GalaxyInfos)
	e.Define("GetPlanetInfo", func(coordinate interface{}) (ogame.PlanetInfos, error) {
		coordTxt := fmt.Sprintf("%v", coordinate)
		coords, err := ogame.ParseCoord(coordTxt)
		if err != nil {
			return ogame.PlanetInfos{}, err
		}
		gi, err := bot.GalaxyInfos(coords.Galaxy, coords.System)
		if err != nil {
			return ogame.PlanetInfos{}, err
		}
		if gi.Position(coords.Position) != nil {
			return *gi.Position(coords.Position), nil
		}
		return ogame.PlanetInfos{}, errors.New("No Celestial at this position")
	})

	e.Define("GetFleetsFromEventList", bot.GetFleetsFromEventList)
	e.Define("OnAttackCh", onAttackCh)

	e.Define("GetResearch", bot.GetResearch)
	e.Define("GetResourcesBuildings", bot.GetResourcesBuildings)
	e.Define("GetFacilities", bot.GetFacilities)
	e.Define("GetShips", bot.GetShips)
	e.Define("GetDefense", bot.GetDefense)

	e.Define("DBGetResourcesBuildings", func(CelestialID ogame.CelestialID) ogame.ResourcesBuildings {
		db := ogb.GetDatabase()
		return db.ResourcesBuildings[CelestialID]
	})
	e.Define("DBGetFacilities", func(CelestialID ogame.CelestialID) ogame.Facilities {
		db := ogb.GetDatabase()
		return db.Facilities[CelestialID]
	})
	e.Define("DBGetShips", func(CelestialID ogame.CelestialID) ogame.ShipsInfos {
		db := ogb.GetDatabase()
		return db.ShipsInfos[CelestialID]
	})
	e.Define("DBGetDefense", func(CelestialID ogame.CelestialID) ogame.DefensesInfos {
		db := ogb.GetDatabase()
		return db.DefensesInfos[CelestialID]
	})
	e.Define("DBGetResearches", func() ogame.Researches {
		db := ogb.GetDatabase()
		return db.Researches
	})
	e.Define("DBGetFleets", func() []ogame.Fleet {
		db := ogb.GetDatabase()
		return db.Movements
	})
	e.Define("DBGetEventFleets", func() []ogame.Fleet {
		db := ogb.GetDatabase()
		return db.EventFleets
	})

	e.Define("CancelFleet", bot.CancelFleet)
	e.Define("GetSlots", bot.GetSlots)
	e.Define("GetPlanet", bot.GetPlanet)
	e.Define("GetPlanets", bot.GetPlanets)
	e.Define("GetMoon", bot.GetMoon)
	e.Define("GetMoons", bot.GetMoons)
	e.Define("GetCelestial", bot.GetCelestial)
	e.Define("GetCelestials", bot.GetCelestials)
	e.Define("GetCachedPlanets", bot.GetCachedPlanets)
	e.Define("GetCachedMoons", bot.GetCachedMoons)
	e.Define("GetCachedCelestial", bot.GetCachedCelestial)
	e.Define("GetCachedPlayer", bot.GetCachedPlayer)
	e.Define("GetCachedCelestialByID", bot.GetCachedCelestialByID)
	e.Define("GetTechs", bot.GetTechs)
	e.Define("GetResources", bot.GetResources)
	e.Define("GetAllResources", bot.GetAllResources)
	e.Define("GetResourcesDetails", bot.GetResourcesDetails)

	e.Define("NewCoordinate", func(g, s, p, t int64) ogame.Coordinate {
		return ogame.Coordinate{
			Galaxy:   g,
			System:   s,
			Position: p,
			Type:     ogame.CelestialType(t),
		}
	})

	e.Define("NewFleetBuilder", func() *ogame.FleetBuilder {
		return ogame.NewFleetBuilder(bot)
	})
	e.Define("NewFleet", func() *ogame.FleetBuilder {
		return ogame.NewFleetBuilder(bot)
	})
	e.Define("NewShipsInfos", func() *ogame.ShipsInfos {
		return &ogame.ShipsInfos{}
	})

	e.Define("NinjaSendFleet", bot.NinjaSendFleet)
	e.Define("NjaCancelFleet", bot.NjaCancelFleet)

	e.Define("Cargo", func(shipsInfos ogame.ShipsInfos) int64 {
		return shipsInfos.Cargo(bot.GetCachedResearch(), bot.GetServer().Settings.EspionageProbeRaids == 1, bot.CharacterClass().IsCollector(), bot.IsPioneers())
	})
	e.Define("CalcCargo", bot.CalcCargo)

	e.Define("GetPrice", func(ogameID int64, nbr int64) ogame.Resources {
		o := ogame.Objs.ByID(ogame.ID(ogameID))
		return o.GetPrice(nbr)
	})

	e.Define("GetMaxExpeditionPoints", bot.GetMaxExpeditionPoints)

	e.Define("ConstructionTime", bot.ConstructionTime)

	e.Define("Sleep", func(wait int64) {
		SleepFunc(time.Duration(wait)*time.Millisecond, ctx)
	})
	e.Define("SleepMs", func(wait int64) {
		time.Sleep(time.Duration(wait) * time.Millisecond)
	})
	e.Define("SleepSec", func(wait int64) {
		//time.Sleep(time.Duration(wait) * time.Second)
		SleepFunc(time.Duration(wait)*time.Second, ctx)
	})
	e.Define("SleepMin", func(wait int64) {
		SleepFunc(time.Duration(wait)*time.Minute, ctx)
	})
	e.Define("SleepRandMs", func(min, max int64) {
		SleepFunc(time.Duration(NextRandom(min, max))*time.Millisecond, ctx)
	})
	e.Define("SleepRandSec", func(min, max int64) {
		SleepFunc(time.Duration(NextRandom(min, max))*time.Second, ctx)
	})
	e.Define("SleepRandMin", func(min, max int64) {
		SleepFunc(time.Duration(NextRandom(min, max))*time.Minute, ctx)
	})
	e.Define("SleepRandHour", func(min, max int64) {
		SleepFunc(time.Duration(NextRandom(min, max))*time.Hour, ctx)
	})
	e.Define("SleepUntil", func(timeStr string) {
		t, _ := time.Parse("15:04:05", timeStr)
		n, _ := time.Parse("15:04:05", time.Now().Format("15:04:05"))
		SleepFunc(t.Sub(n)+time.Duration(24)*time.Hour, ctx)
	})
	e.Define("SleepDur", func(dur time.Duration) {
		SleepFunc(dur, ctx)
	})
	e.Define("timeNow", func() time.Time {
		return time.Now()
	})

	e.Define("CronExec", func(time string) {
		time = ""
	})

	e.Define("Random", NextRandom)
	e.Define("FlightTime", bot.FlightTime)
	e.Define("BuyOfferOfTheDay", bot.BuyOfferOfTheDay)
	e.Define("GetCachedCelestials", bot.GetCachedCelestials)
	e.Define("Put", func(key string, value interface{}) error {
		by, _ := json.Marshal(value)
		return WebKVStore.Put(key, by)
	})
	e.Define("Get", func(key string) (interface{}, error) {
		var val []byte
		err := WebKVStore.Get(key, &val)
		var res interface{}
		json.Unmarshal(val, &res)
		return res, err
	})
	e.Define("Has", func(key string) bool {
		var val []byte
		err := WebKVStore.Get(key, &val)
		if err != nil {
			return false
		}
		return true
	})
	e.Define("Delete", WebKVStore.Delete)

	e.Define("Print", func(a ...interface{}) (n int, err error) {
		log.Println("VM", a)
		return
	})

	e.Define("SetPreferences", bot.SetPreferences)

	// e.Define("Put", func(key, val interface{}) error {
	// 	return funcPut(key, val, badgerDB)
	// })

	// e.Define("Get", func(key interface{}) (interface{}, error) {
	// 	return funcGet(key, badgerDB)
	// })

	/*
			# Global API

		    -OwnBots
		    -AllBots
		    -GetBotByID
		    -StartRemoteScript
		    -StopRemoteScript
		    +Sleep
		    +SleepMs
		    +SleepSec
		    +SleepMin
		    +SleepRandMs
		    +SleepRandSec
		    +SleepRandMin
		    +SleepRandHour
		    +SleepUntil
		    +SleepDur
		    +Random
		    +Print
		    +print
		    LogDebug
		    LogInfo
		    LogWarn
		    LogError
		    ClearOut
		    RepatriateNow
		    TriggerFleetSave
		    GetNextSleepTime
		    GetNextWakeTime
		    AddItemToQueue
		    ClearAllConstructionQueues
		    +GetPrice
		    ConstructionTime
		    GetRequirements
		    IsAvailable
		    SolarSatelliteProduction
		    IgnorePlanet
		    IgnorePlayer
		    IgnoreAlliance
		    IsPlanetIgnored
		    IsPlayerIgnored
		    IsAllianceIgnored
		    ExecIn
		    ExecAt
		    ExecAtCh
		    IntervalExec
		    CronExec
		    RemoveCron
		    RangeCronExec
		    Clock
		    Date
		    Weekday
		    GetTimestamp
		    Unix
		    OnStateChange
		    RecruitOfficer
		    HasCommander
		    HasAdmiral
		    HasEngineer
		    HasGeologist
		    HasTechnocrat
		    GetHomeWorld
		    SetHomeWorld
		    +EnableNJA
		    +DisableNJA
		    +IsNJAEnabled
		    +Login
		    +Logout
		    +IsLoggedIn
		    IsPioneers
		    IsUnderAttack
		    CharacterClass
		    IsCollector
		    IsGeneral
		    IsDiscoverer
		    SendMessage
		    SendMessageAlliance
		    +GetCelestial
		    +GetCelestials
		    +GetPlanet
		    +GetPlanets
		    +GetMoon
		    +GetMoons
		    +GetCachedPlanets
		    +GetCachedMoons
		    +GetCachedCelestial
		    +GetCachedCelestials
		    +GetCachedPlayer
		    PlayersData
		    PlayerDataByID
		    PlayerDataByName
		    FindDebrisFieldWithMinimumTravelTime
		    FindEmptyPlanetWithMinimumTravelTime
		    GetSystemsInRange
		    GetSystemsInRangeAsc
		    GetSystemsInRangeDesc
		    Shuffle
		    GetResources
		    +GetResourcesDetails
		    GetAllResources
		    +Abandon
		    CollectAllMarketplaceMessages
		    CollectMarketplaceMessage
		    GetDMCosts
		    UseDM
		    GetItems
		    GetActiveItems
		    ActivateItem
		    OfferSellMarketplace
		    OfferBuyMarketplace
		    AttackStrength
		    ShipsAttackStrength
		    ShipsAttackStrengthUsingOwnResearches
		    GetHighscore
			+GetFleetsFromEventlist
		    +GetFleets
		    +GetAttacks
		    +GalaxyInfos
		    +GetPlanetInfo
		    GetDepartureTime
		    LowestSpeed
		    Simulator
		    VersionCompare
		    Dotify
		    ShortDur
		    ShortShipsInfos
		    Bytes2Str
		    ID2Str
		    Atoi
		    Itoa
		    StartScript
		    StopScript
		    PauseScript
		    ResumeScript
		    IsPausedScript
		    IsScriptRunning
		    GetScripts
		    GetRunningScripts
		    StartBrain
		    StopBrain
		    StartScanner
		    StopScanner
		    StartHunter
		    StopHunter
		    StartSleepMode
		    StopSleepMode
		    StartPhalanxSession
		    StopPhalanxSession
		    StartFarmingBot
		    StopFarmingBot
		    PauseFarmingBot
		    ResumeFarmingBot
		    IsRunningFarmingBot
		    IsPausedFarmingBot
		    IsFarmSessionOngoing
		    FarmingBotSessionsCount
		    AbortAllFarmingSessions
		    StartExpeditionsBot
		    StopExpeditionsBot
		    IsRunningExpeditionsBot
		    StartDefenderBot
		    StopDefenderBot
		    IsRunningDefenderBot
		    SetDefenderCheckInterval
		    SetDefenderCheckOrigin
		    +DBGetResourceBuildings
		    +DBGetResourceSettings
		    +DBGetFacilities
		    +DBGetShips
		    +DBGetDefenses
		    +DBGetResearches
			+DBGetFleets
			+DBGetEventFleets
		    +Put
		    +Get
		    Has
		    +Delete
		    ParseCoord
		    NowTimeString
		    NowInTimeRange
		    DurationBetweenTimeStrings
		    MillisecondsBetweenTimeStrings
		    ParseNextDatetimeAt
		    GetNextDatetimeAt
		    +Cargo
		    +CalcCargo
		    CalcFastCargo
		    CalcFastCargoPF
		    CalcPreferredCargo
		    GetPlayerCoordinates
		    Base64
		    Base64Decode
		    GetSortedCelestials
		    GetSortedPlanets
		    GetSortedMoons
		    SendMail
		    SendTelegram
		    SendDiscord
		    NewFleet
		    NewFarmSession
		    NewACS
		    CreateUnion
		    Notify
		    PlaySound
		    GetEspionageReportMessages
		    GetEspionageReportFor
		    GetEspionageReport
		    GetCombatReportSummaryFor
		    DeleteMessage
		    DeleteAllMessagesFromTab
		    Build
		    BuildBuilding
		    BuildTechnology
		    TearDown
		    CancelBuilding
		    CancelResearch
		    +CancelFleet
		    +GetResourcesBuildings
		    +GetTechs
		    +GetFacilities
		    +GetResearch
		    +GetShips
		    +GetDefense
		    GetProduction
		    ConstructionsBeingBuilt
		    GetResourceSettings
		    SetResourceSettings
		    +GetSlots
		    Distance
		    +FlightTime
		    DestroyRockets
		    SendIPM
		    +BuyOfferOfTheDay
		    GetAuction
		    DoAuction
		    Phalanx
		    UnsafePhalanx
		    JumpGate
		    JumpGate2
		    JumpGateDestinations
		    GetFleetSlotsReserved
		    SetFleetSlotsReserved
		    NewResources
		    NewResourceSettings
		    +NewCoordinate
		    NewTemperature
		    +NewShipsInfos
		    TempFile
		    ListTempFiles
		    DeleteTempFile
		    DeleteAllTempFiles
		    JsonDecode
		    Min
		    Max
		    Ceil
		    Floor
		    Round
		    Abs
		    Pow
		    Sqrt
		    Terminate

	*/

	script := `
println("Hello World :)")
println(GetFleets())
fb = NewFleetBuilder()
fb.AddShips(LIGHTFIGHTER, 1)
_, err = fb.SendNow()
if err != nil {
	println(err)
}

Sleep(5000)
println(GetPrice(LIGHTFIGHTER, 10))

`
	// Run Function Code
	script = code

	//ctx, cancel := context.WithCancel(context.Background())
	_, err = vm.ExecuteContext(ctx, e, nil, script)
	// output: Hello World :)
	if err != nil {
		log.Println(err)
		return errors.New(err.Error())
	}
	return err
}

func NextRandom(min, max int64) int64 {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return r1.Int63n(max) + min
}

// func funcPut(key interface{}, val interface{}, db *badger.DB) error {
// 	err := db.Update(func(txn *badger.Txn) error {
// 		keyBy, err2 := json.Marshal(key)
// 		if err2 != nil {
// 			return err2
// 		}
// 		valBy, err2 := json.Marshal(val)
// 		if err2 != nil {
// 			return err2
// 		}
// 		return txn.SetEntry(badger.NewEntry(keyBy, valBy))
// 	})
// 	return err
// }

// func funcGet(key interface{}, db *badger.DB) (interface{}, error) {
// 	var val interface{}

// 	err := db.View(func(txn *badger.Txn) error {
// 		keyBy, _ := json.Marshal(key)
// 		item, err := txn.Get(keyBy)
// 		if err != nil {
// 			return err
// 		}

// 		err = item.Value(func(itemValue []byte) error {
// 			return json.Unmarshal(itemValue, &val)
// 		})

// 		return err
// 	})
// 	return val, err
// }

// SortFuncsStruct provides functions to be used with Sort
type SortFuncsStruct struct {
	LenFunc  func() int
	LessFunc func(i, j int) bool
	SwapFunc func(i, j int)
}

func (s SortFuncsStruct) Len() int           { return s.LenFunc() }
func (s SortFuncsStruct) Less(i, j int) bool { return s.LessFunc(i, j) }
func (s SortFuncsStruct) Swap(i, j int)      { s.SwapFunc(i, j) }

var RunningScripts map[string]context.CancelFunc = make(map[string]context.CancelFunc)

func RunScriptHandler(c echo.Context) error {
	bot := c.Get("bot").(*ogame.OGame)
	db := c.Get("ogb").(*ogb.Ogb)

	//badgerDB := c.Get("badgerDB").(*badger.DB)
	c.Request().ParseForm()
	name := c.FormValue("name")
	code := c.FormValue("code")

	_, ok := RunningScripts[name]
	if !ok {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			RunningScripts[name] = cancel
			//RunScript(bot, badgerDB, ctx, name, code)
			RunScript(bot, db, ctx, name, code)
			cancel()
			delete(RunningScripts, name)
		}()
	} else {
		return c.JSON(http.StatusOK, ogame.SuccessResp("Is running"))
	}

	return c.JSON(http.StatusOK, ogame.SuccessResp("success"))
}

func GetScriptsHandler(c echo.Context) error {
	c.Request().ParseForm()
	bot := c.Get("bot").(*ogame.OGame)
	database := c.Get("database").(*ogb.Ogb)
	cache, _ := json.Marshal(database)
	db := ogb.New()
	json.Unmarshal(cache, &db)

	db.Celestials = bot.GetCachedPlanets()

	obj := struct {
		Bot             *ogame.OGame
		DB              *ogb.Ogb
		ObjsStruct      ogame.ObjsStruct
		PlanetBuildings []ogame.Building
		MoonBuildings   []ogame.Building
		Buildings       []ogame.Building
		Ships           []ogame.Ship
		Defenses        []ogame.Defense
		Technologies    []ogame.Technology
		Scripts         map[string]ogb.Scripts
		ScriptName      string
		IsRunning       bool
		StopScript      context.CancelFunc
		RunningScripts  map[string]context.CancelFunc
	}{
		bot,
		db,
		ogame.Objs,
		ogame.PlanetBuildings,
		ogame.MoonBuildings,
		ogame.Buildings,
		ogame.Ships,
		ogame.Defenses,
		ogame.Technologies,
		db.Scripts,
		"",
		false,
		nil,
		RunningScripts,
	}
	scriptName := c.QueryParam("name")
	stopScript := c.QueryParam("stop")
	if scriptName != "" {
		_, ok := db.Scripts[scriptName]
		if ok {
			obj.ScriptName = scriptName
		}
		_, ok = RunningScripts[scriptName]
		if ok {
			obj.IsRunning = true
			obj.StopScript = RunningScripts[scriptName]
			if stopScript == "true" {
				unregisterOnAttackCh(scriptName)
				obj.StopScript()
			}
		}
	}
	return c.Render(http.StatusOK, "scripts", obj)
}

func PostNewScriptHandler(c echo.Context) error {
	c.Request().ParseForm()
	scriptname := c.FormValue("name")
	code := c.FormValue("code")
	//bot := c.Get("bot").(*ogame.OGame)
	database := c.Get("database").(*ogb.Ogb)

	database.Lock()
	database.Scripts[scriptname] = ogb.Scripts{
		Name:   scriptname,
		Script: code,
	}
	database.Unlock()

	cache, _ := json.Marshal(database)
	db := ogb.New()
	json.Unmarshal(cache, &db)

	obj := struct {
		DB             *ogb.Ogb
		ScriptName     string
		RunningScripts map[string]context.CancelFunc
	}{
		DB:             db,
		ScriptName:     scriptname,
		RunningScripts: RunningScripts,
	}

	return c.Render(http.StatusOK, "scripts", obj)
}

func SleepFunc(wait time.Duration, ctx context.Context) {
	select {
	case <-time.After(wait):
		return
	case <-ctx.Done():
		return
	}
}
