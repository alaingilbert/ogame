package ogb

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/faunX/ogame"
)

func Start() {

}

type Ogb struct {
	Database
}

var err error

func New() *Ogb {
	b := Ogb{}

	b.ResourcesBuildings = map[ogame.CelestialID]ogame.ResourcesBuildings{}
	b.Facilities = map[ogame.CelestialID]ogame.Facilities{}
	b.ShipsInfos = map[ogame.CelestialID]ogame.ShipsInfos{}
	b.DefensesInfos = map[ogame.CelestialID]ogame.DefensesInfos{}
	b.ResourcesDetails = map[ogame.CelestialID]ogame.ResourcesDetails{}
	//b.Galaxy = map[ogame.Coordinate]ogame.PlanetInfos{}
	b.Activities = map[ogame.CelestialID]time.Time{}
	b.Constructions = map[ogame.CelestialID]BrainConstruction{}
	b.Productions = map[ogame.CelestialID]BrainProduction{}
	b.MovementsHistory = map[ogame.FleetID]ogame.Fleet{}
	b.AttackEventsHistory = map[ogame.FleetID]ogame.AttackEvent{}
	b.EventFleetsHistory = map[ogame.FleetID]ogame.Fleet{}
	b.BrainQueue = []BrainQueueType{}
	b.EventFleets = []ogame.Fleet{}
	b.SentFleets = []SentFleet{}

	b.Celestials = []ogame.Planet{}
	b.Planets = []ogame.Planet{}
	b.Moons = []ogame.Moon{}
	b.BrainQueue = []BrainQueueType{}
	b.Movements = []ogame.Fleet{}
	b.AttackEvents = []ogame.AttackEvent{}
	b.EventFleets = []ogame.Fleet{}
	// Galaxy
	b.Galaxy = map[string]struct {
		SystemInfos ogame.SystemInfos `json:"systemInfos"`
		//PlanetInfos []ogame.PlanetInfos `json:"planetInfo"`
	}{}

	//Galaxy map[ogame.Coordinate]ogame.SystemInfos
	//Galaxy map[ogame.Coordinate][]ogame.PlanetInfos

	// Messages
	b.EspionageReports = map[int64]ogame.EspionageReport{}
	b.CombatReportSummary = map[int64]ogame.CombatReportSummary{}
	b.FullCombatReports = map[int64]ogame.FullCombatReport{}
	b.EspionageReportSummary = map[int64]ogame.EspionageReportSummary{}
	b.ExpeditionMessages = map[int64]ogame.ExpeditionMessage{}
	b.Messages = map[int64]struct {
		Tabid int64
		ogame.Message
	}{}
	b.Traffic = struct {
		In  int64
		Out int64
		RPS int64
	}{}

	b.Scripts = make(map[string]Scripts)

	return &b
}

type BrainQueueType struct {
	ogame.CelestialID
	ogame.Quantifiable
}

type BrainConstruction struct {
	Quantifiable ogame.Quantifiable
	FinishAt     time.Time
}

type BrainProduction struct {
	Productions []ogame.Quantifiable
	FinishAt    time.Time
}

type Scripts struct {
	Name      string
	Script    string
	AutoStart bool
}

type Database struct {
	sync.RWMutex
	// General
	LastActiveCelestialID ogame.CelestialID                              `json:"lastActiveCelestialID"`
	Activities            map[ogame.CelestialID]time.Time                `json:"activities"`
	Celestials            []ogame.Planet                                 `json:"celestials"`
	Planets               []ogame.Planet                                 `json:"planets"`
	Moons                 []ogame.Moon                                   `json:"moons"`
	Researches            ogame.Researches                               `json:"researches"`
	ResourcesBuildings    map[ogame.CelestialID]ogame.ResourcesBuildings `json:"resourcesBuildings"`
	Facilities            map[ogame.CelestialID]ogame.Facilities         `json:"facilities"`
	ShipsInfos            map[ogame.CelestialID]ogame.ShipsInfos         `json:"shipsInfos"`
	DefensesInfos         map[ogame.CelestialID]ogame.DefensesInfos      `json:"defensesInfos"`
	ResourcesDetails      map[ogame.CelestialID]ogame.ResourcesDetails   `json:"resourcesDetails"`

	// Brain
	BrainQueue         []BrainQueueType                        `json:"brainQueue"`
	Constructions      map[ogame.CelestialID]BrainConstruction `json:"constructions"`
	ResearchInProgress BrainConstruction                       `json:"researchInProgress"`
	Productions        map[ogame.CelestialID]BrainProduction   `json:"productions"`

	// Fleets
	Slots               ogame.Slots                         `json:"slots"`
	Movements           []ogame.Fleet                       `json:"movements"`
	MovementsHistory    map[ogame.FleetID]ogame.Fleet       `json:"movementsHistory"`
	AttackEvents        []ogame.AttackEvent                 `json:"attackEvents"`
	AttackEventsHistory map[ogame.FleetID]ogame.AttackEvent `json:"attackEventsHistory"`
	EventFleets         []ogame.Fleet                       `json:"eventFleets"`
	EventFleetsHistory  map[ogame.FleetID]ogame.Fleet       `json:"eventFleetsHistory"`
	SentFleets          []SentFleet                         `json:"sentFleets"`

	// Galaxy
	Galaxy map[string]struct {
		SystemInfos ogame.SystemInfos `json:"systemInfos"`
		//PlanetInfos []ogame.PlanetInfos `json:"planetInfo"`
	} `json:"Galaxy"`

	// Messages
	EspionageReports       map[int64]ogame.EspionageReport        `json:"espionageReports"`
	CombatReportSummary    map[int64]ogame.CombatReportSummary    `json:"combatReportSummary"`
	FullCombatReports      map[int64]ogame.FullCombatReport       `json:"fullCombatReports"`
	EspionageReportSummary map[int64]ogame.EspionageReportSummary `json:"espionageReportSummary"`
	ExpeditionMessages     map[int64]ogame.ExpeditionMessage      `json:"expeditionMessages"`
	Messages               map[int64]struct {
		Tabid int64
		ogame.Message
	} `json:"messages"`

	// Scripts
	Scripts map[string]Scripts

	Traffic struct {
		In  int64
		Out int64
		RPS int64
	} `json:"traffic"`
}

func (d *Database) GetDatabase() *Database {
	d.RLock()
	defer d.RUnlock()

	o := New()
	originalDB, _ := json.Marshal(d)
	json.Unmarshal(originalDB, &o)

	return &o.Database
}

func (d *Database) GetResources(celestialID ogame.CelestialID) ogame.Resources {
	res := ogame.Resources{}
	lastActiveTime := d.Activities[celestialID]
	TimePassedSec := time.Now().Sub(lastActiveTime).Seconds()
	// Metal
	available := d.ResourcesDetails[celestialID].Metal.Available
	capacity := d.ResourcesDetails[celestialID].Metal.StorageCapacity
	production := d.ResourcesDetails[celestialID].Metal.StorageCapacity
	if capacity > available {
		productionPerSec := float64(production/3600)*TimePassedSec + float64(available)
		res.Metal = int64(productionPerSec)
	}

	// Crystal
	available = d.ResourcesDetails[celestialID].Crystal.Available
	capacity = d.ResourcesDetails[celestialID].Crystal.StorageCapacity
	production = d.ResourcesDetails[celestialID].Crystal.StorageCapacity
	if capacity > available {
		productionPerSec := float64(production/3600)*TimePassedSec + float64(available)
		res.Crystal = int64(productionPerSec)
	}

	// Deuterium
	available = d.ResourcesDetails[celestialID].Deuterium.Available
	capacity = d.ResourcesDetails[celestialID].Deuterium.StorageCapacity
	production = d.ResourcesDetails[celestialID].Deuterium.StorageCapacity
	if capacity > available {
		productionPerSec := float64(production/3600)*TimePassedSec + float64(available)
		res.Deuterium = int64(productionPerSec)
	}
	return ogame.Resources{}
}

type Galaxy struct {
	galaxy  int64
	system  int64
	planets []struct {
		ID              int64
		Activity        int64 // no activity: 0, active: 15, inactive: [16, 59]
		Name            string
		Img             string
		Coordinate      ogame.Coordinate
		Administrator   bool
		Destroyed       bool
		Inactive        bool
		Vacation        bool
		StrongPlayer    bool
		Newbie          bool
		HonorableTarget bool
		Banned          bool
		Debris          struct {
			Metal           int64
			Crystal         int64
			RecyclersNeeded int64
		}
		Moon struct {
			ID       int64
			Diameter int64
			Activity int64
		}
		Player struct {
			ID         int64
			Name       string
			Rank       int64
			IsBandit   bool
			IsStarlord bool
		}
		Alliance struct {
			ID            int64
			Name          string
			Rank          int64
			Member        int64
			AllianceClass ogame.AllianceClass
		}
		Date time.Time
	}
	ExpeditionDebris struct {
		Metal             int64
		Crystal           int64
		PathfindersNeeded int64
	}
	Events struct {
		Darkmatter  int64
		HasAsteroid bool
	}
}

func (d *Database) AddToBrainQueue(b BrainQueueType) {
	var exists bool
	var level int64
	for _, v := range d.BrainQueue {
		if v.CelestialID == b.CelestialID && v.ID == b.ID {
			exists = true
			if level == 0 {
				level = v.Nbr + 1
			} else {
				level++
			}
		}
	}
	if !exists {
		if b.ID.IsResourceBuilding() {
			b.Nbr = d.ResourcesBuildings[b.CelestialID].ByID(b.ID) + 1
		}
		if b.ID.IsFacility() {
			b.Nbr = d.Facilities[b.CelestialID].ByID(b.ID) + 1
		}
		if b.ID.IsTech() {
			b.Nbr = d.Researches.ByID(b.ID) + 1
		}
		d.BrainQueue = append(d.BrainQueue, b)
	} else {
		b.Nbr = level
		d.BrainQueue = append(d.BrainQueue, b)
	}
}
