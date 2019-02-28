<img src="./logo.png" width="300" />

[![Build Status](https://travis-ci.org/alaingilbert/ogame.svg?branch=master)](https://travis-ci.org/alaingilbert/ogame) [![codecov](https://codecov.io/gh/alaingilbert/ogame/branch/master/graph/badge.svg)](https://codecov.io/gh/alaingilbert/ogame) [![codecov](https://img.shields.io/discord/546546108277719052.svg)](https://discord.gg/thnbsP)

# OGame automation toolkit

- [As a library](#ogame-library)
- [As a service (ogamed)](#ogamed-service)

---

# ogame library

### Verify attack example

```go
package main

import "fmt"
import "os"
import "github.com/alaingilbert/ogame"

func main() {
	universe := os.Getenv("UNIVERSE") // eg: Bellatrix
	username := os.Getenv("USERNAME") // eg: email@gmail.com
	password := os.Getenv("PASSWORD") // eg: *****
	language := os.Getenv("LANGUAGE") // eg: en
	bot, err := ogame.New(universe, username, password, language)
	if err != nil {
		panic(err)
	}
	attacked := bot.IsUnderAttack()
	fmt.Println(attacked) // False
}
```

### Available methods

```go
IsActive() bool
Quiet(bool)
Tx(clb func(tx *Prioritize) error) error
Begin() *Prioritize
WithPriority(priority int) *Prioritize
GetPublicIP() (string, error)
OnStateChange(clb func(locked bool, actor string))
GetState() (bool, string)
IsLocked() bool
GetSession() string
AddAccount(number int, lang string) (NewAccount, error)
GetServer() Server
SetUserAgent(newUserAgent string)
ServerURL() string
GetLanguage() string
GetPageContent(url.Values) []byte
GetAlliancePageContent(url.Values) []byte
PostPageContent(url.Values, url.Values) []byte
Login() error
Logout()
GetUsername() string
GetUniverseName() string
GetUniverseSpeed() int
GetUniverseSpeedFleet() int
IsDonutGalaxy() bool
IsDonutSystem() bool
FleetDeutSaveFactor() float64
ServerVersion() string
ServerTime() time.Time
IsUnderAttack() bool
GetUserInfos() UserInfos
SendMessage(playerID int, message string) error
GetFleets() ([]Fleet, Slots)
GetFleetsFromEventList() []Fleet
CancelFleet(FleetID) error
GetAttacks() []AttackEvent
GalaxyInfos(galaxy, system int) (SystemInfos, error)
GetResearch() Researches
GetCachedPlanets() []Planet
GetCachedMoons() []Moon
GetCachedCelestial(interface{}) Celestial
GetCachedPlayer() UserInfos
GetPlanets() []Planet
GetPlanet(interface{}) (Planet, error)
GetMoons(MoonID) []Moon
GetMoon(interface{}) (Moon, error)
GetCelestial(interface{}) (Celestial, error)
GetCelestials() ([]Celestial, error)
GetEspionageReportMessages() ([]EspionageReportSummary, error)
GetEspionageReportFor(Coordinate) (EspionageReport, error)
GetEspionageReport(msgID int) (EspionageReport, error)
GetCombatReportSummaryFor(Coordinate) (CombatReportSummary, error)
DeleteMessage(msgID int) error
Distance(origin, destination Coordinate) int
FlightTime(origin, destination Coordinate, speed Speed, ships ShipsInfos) (secs, fuel int)
RegisterChatCallback(func(ChatMsg))
RegisterHTMLInterceptor(func(method string, params, payload url.Values, pageHTML []byte))
GetSlots() Slots

// Planet or Moon functions
GetResources(CelestialID) (Resources, error)
SendFleet(celestialID CelestialID, ships []Quantifiable, speed Speed, where Coordinate, mission MissionID, resources Resources, expeditiontime int) (Fleet, error)
Build(celestialID CelestialID, id ID, nbr int) error
BuildCancelable(CelestialID, ID) error
BuildProduction(celestialID CelestialID, id ID, nbr int) error
BuildBuilding(celestialID CelestialID, buildingID ID) error
BuildDefense(celestialID CelestialID, defenseID ID, nbr int) error
BuildShips(celestialID CelestialID, shipID ID, nbr int) error
CancelBuilding(CelestialID) error
ConstructionsBeingBuilt(CelestialID) (buildingID ID, buildingCountdown int, researchID ID, researchCountdown int)
GetProduction(CelestialID) ([]Quantifiable, error)
GetFacilities(CelestialID) (Facilities, error)
GetDefense(CelestialID) (DefensesInfos, error)
GetShips(CelestialID) (ShipsInfos, error)
GetResourcesBuildings(CelestialID) (ResourcesBuildings, error)
CancelResearch(CelestialID) error
BuildTechnology(celestialID CelestialID, technologyID ID) error

// Planet specific functions
GetResourceSettings(PlanetID) (ResourceSettings, error)
SetResourceSettings(PlanetID, ResourceSettings) error
SendIPM(PlanetID, Coordinate, int, ID) (int, error)
GetResourcesProductions(PlanetID) (Resources, error)
GetResourcesProductionsLight(ResourcesBuildings, Researches, ResourceSettings, Temperature) Resources

// Moon specific functions
Phalanx(MoonID, Coordinate) ([]Fleet, error)
```

### Full documentation

[https://godoc.org/github.com/alaingilbert/ogame](https://godoc.org/github.com/alaingilbert/ogame)

---

# ogamed service

Download [ogamed binary here](https://github.com/alaingilbert/ogame/releases)  
Full documentation [here](https://github.com/alaingilbert/ogame/wiki/ogamed-full-documentation)

```
./ogamed --universe=Zibal --username=email@email.com --password=secret --language=en
```

```
$ curl 127.0.0.1:8080/bot/is-under-attack
{"Status":"ok","Code":200,"Message":"","Result":false}

$ curl 127.0.0.1:8080/bot/send-message -d 'playerID=123&message="Sup boi!"'
{"Status":"ok","Code":200,"Message":"","Result":null}

$ curl 127.0.0.1:8080/bot/user-infos
{"Status":"ok","Code":200,"Message":"","Result":{"PlayerID":106734,"PlayerName":"Commodore Nomad","Points":43825,"Rank":1130,"Total":1675,"HonourPoints":0}}
```

```
POST /bot/set-user-agent
GET  /bot/server-url
POST /bot/page-content
GET  /bot/login
GET  /bot/logout
GET  /bot/server/speed
GET  /bot/server/version
GET  /bot/server/time
GET  /bot/is-under-attack
GET  /bot/user-infos
POST /bot/send-message
GET  /bot/fleets
POST /bot/fleets/:fleetID/cancel
GET  /bot/attacks
GET  /bot/galaxy-infos/:galaxy/:system
GET  /bot/get-research
GET  /bot/planets
GET  /bot/planets/:galaxy/:system/:position
GET  /bot/planets/:planetID
GET  /bot/planets/:planetID/resource-settings
POST /bot/planets/:planetID/resource-settings
GET  /bot/planets/:planetID/resources-buildings
GET  /bot/planets/:planetID/defence
GET  /bot/planets/:planetID/ships
GET  /bot/planets/:planetID/facilities
POST /bot/planets/:planetID/build/:ogameID/:nbr
POST /bot/planets/:planetID/build/cancelable/:ogameID
POST /bot/planets/:planetID/build/production/:ogameID/:nbr
POST /bot/planets/:planetID/build/building/:ogameID
POST /bot/planets/:planetID/build/technology/:ogameID
POST /bot/planets/:planetID/build/defence/:ogameID/:nbr
POST /bot/planets/:planetID/build/ships/:ogameID/:nbr
GET  /bot/planets/:planetID/production
GET  /bot/planets/:planetID/constructions
POST /bot/planets/:planetID/cancel-building
POST /bot/planets/:planetID/cancel-research
GET  /bot/planets/:planetID/resources
POST /bot/planets/:planetID/send-fleet
```
