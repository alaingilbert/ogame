package ogame

import (
	"strconv"
	"time"
)

type MessagesTabID int64

// CelestialID represent either a PlanetID or a MoonID
type CelestialID int64

// MoonID represent a moon id
type MoonID CelestialID

// Celestial convert a MoonID to a CelestialID
func (m MoonID) Celestial() CelestialID {
	return CelestialID(m)
}

// Speed represent a fleet speed
type Speed float64

// Float64 returns a float64 value of the speed
func (s Speed) Float64() float64 {
	return float64(s)
}

// Int64 returns an integer value of the speed
func (s Speed) Int64() int64 {
	return int64(s)
}

// Int returns an integer value of the speed
// Deprecated: backward compatibility
func (s Speed) Int() int64 {
	return int64(s)
}

func (s Speed) String() string {
	switch s {
	case FivePercent:
		return "5%"
	case TenPercent:
		return "10%"
	case FifteenPercent:
		return "15%"
	case TwentyPercent:
		return "20%"
	case TwentyFivePercent:
		return "25%"
	case ThirtyPercent:
		return "30%"
	case ThirtyFivePercent:
		return "35%"
	case FourtyPercent:
		return "40%"
	case FourtyFivePercent:
		return "45%"
	case FiftyPercent:
		return "50%"
	case FiftyFivePercent:
		return "55%"
	case SixtyPercent:
		return "60%"
	case SixtyFivePercent:
		return "65%"
	case SeventyPercent:
		return "70%"
	case SeventyFivePercent:
		return "75%"
	case EightyPercent:
		return "80%"
	case EightyFivePercent:
		return "85%"
	case NinetyPercent:
		return "90%"
	case NinetyFivePercent:
		return "95%"
	case HundredPercent:
		return "100%"
	default:
		return strconv.FormatFloat(float64(s), 'f', 1, 64)
	}
}

type ResourcesResp struct {
	Metal struct {
		Resources struct {
			ActualFormat string
			Actual       int64
			Max          int64
			Production   float64
		}
		Tooltip string
		Class   string
	}
	Crystal struct {
		Resources struct {
			ActualFormat string
			Actual       int64
			Max          int64
			Production   float64
		}
		Tooltip string
		Class   string
	}
	Deuterium struct {
		Resources struct {
			ActualFormat string
			Actual       int64
			Max          int64
			Production   float64
		}
		Tooltip string
		Class   string
	}
	Energy struct {
		Resources struct {
			ActualFormat string
			Actual       int64
		}
		Tooltip string
		Class   string
	}
	Darkmatter struct {
		Resources struct {
			ActualFormat string
			Actual       int64
		}
		String  string
		Tooltip string
	}
	HonorScore int64
}

type planetResource struct {
	Input struct {
		Metal     int64
		Crystal   int64
		Deuterium int64
	}
	Output struct {
		Metal     int64
		Crystal   int64
		Deuterium int64
	}
	IsMoon        bool
	ImageFileName string
	Name          string
	// OtherPlanet   string // can be null or apparently number (cannot unmarshal number into Go struct field planetResource.OtherPlanet of type string)
}

// PlanetResources ...
type PlanetResources map[CelestialID]planetResource

// Multiplier ...
type Multiplier struct {
	Metal     float64
	Crystal   float64
	Deuterium float64
	Honor     float64
}

// EspionageReportType type of espionage report (action or report)
type EspionageReportType int

// Action message received when an enemy is seen near your planet
const Action EspionageReportType = 0

// Report message received when you spied on someone
const Report EspionageReportType = 1

// CombatReportSummary summary of combat report
type CombatReportSummary struct {
	ID           int64
	APIKey       string
	FleetID      FleetID
	Origin       *Coordinate
	Destination  Coordinate
	AttackerName string
	DefenderName string
	Loot         int64
	Metal        int64
	Crystal      int64
	Deuterium    int64
	DebrisField  int64
	CreatedAt    time.Time
}

// EspionageReportSummary summary of espionage report
type EspionageReportSummary struct {
	ID             int64
	Type           EspionageReportType
	From           string // Fleet Command | Space Monitoring
	Target         Coordinate
	LootPercentage float64
}

// ExpeditionMessage ...
type ExpeditionMessage struct {
	ID         int64
	Coordinate Coordinate
	Content    string
	Resources  Resources
	Ships      ShipsInfos
	CreatedAt  time.Time
}

// MarketplaceMessage ...
type MarketplaceMessage struct {
	ID                  int64
	Type                int64 // 26: purchases, 27: sales
	CreatedAt           time.Time
	Token               string
	MarketTransactionID int64
}

// Preferences ...
type Preferences struct {
	SpioAnz                            int64
	SpySystemAutomaticQuantity         int64
	SpySystemTargetPlanetTypes         int64
	SpySystemTargetPlayerTypes         int64
	SpySystemIgnoreSpiedInLastXMinutes int64
	DisableChatBar                     bool // no-mobile
	DisableOutlawWarning               bool
	MobileVersion                      bool
	ShowOldDropDowns                   bool
	ActivateAutofocus                  bool
	Language                           string
	EventsShow                         int64 // Hide: 1, Above the content: 2, Below the content: 3
	SortSetting                        int64 // Order of emergence: 0, Coordinates: 1, Alphabet: 2, Size: 3, Used fields: 4
	SortOrder                          int64 // Up: 0, Down: 1
	ShowDetailOverlay                  bool
	AnimatedSliders                    bool // no-mobile
	AnimatedOverview                   bool // no-mobile
	PopupsNotices                      bool // no-mobile
	PopopsCombatreport                 bool // no-mobile
	SpioReportPictures                 bool
	MsgResultsPerPage                  int64 // 10, 25, 50
	AuctioneerNotifications            bool
	EconomyNotifications               bool
	ShowActivityMinutes                bool
	PreserveSystemOnPlanetChange       bool
	DiscoveryWarningEnabled            bool
	UrlaubsModus                       bool // Vacation mode

	// Mobile only
	Notifications struct {
		BuildList               bool
		FriendlyFleetActivities bool
		HostileFleetActivities  bool
		ForeignEspionage        bool
		AllianceBroadcasts      bool
		AllianceMessages        bool
		Auctions                bool
		Account                 bool
	}
}

type ACSValues struct {
	Galaxy        int64
	System        int64
	Position      int64
	CelestialType CelestialType
	Name          string
	ACSValues     string // Raw string is used to send fleet. eg: `1#2#3#3#name#123456`
	Union         int64
}

type Relocation struct {
	MoveLink             string
	PlanetMovePossible   bool
	SufficientDarkMatter bool
	MissionType          MissionID
	DarkMatterCost       int64
}
