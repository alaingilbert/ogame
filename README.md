<img src="./logo.png" width="300" />

[![Build Status](https://travis-ci.org/alaingilbert/ogame.svg?branch=master)](https://travis-ci.org/alaingilbert/ogame) [![codecov](https://codecov.io/gh/alaingilbert/ogame/branch/master/graph/badge.svg)](https://codecov.io/gh/alaingilbert/ogame) [![discord](https://img.shields.io/discord/546546108277719052.svg)](https://discord.gg/4dZtZpz)

# OGame automation toolkit

- [As a library](#ogame-library)
- [As a service (ogamed)](#ogamed-service)
- [As a docker container](#docker-container)

---

Discord: https://discord.gg/4dZtZpz

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
	attacked, _ := bot.IsUnderAttack()
	fmt.Println(attacked) // False
}
```

##### How to get started

- Ensure you have go 1.12 or above `go version`
- Copy and paste the above example in a file `main.go`
- Create your own module `go mod init my_project`
- Install dependencies `go mod vendor`
- Run the code `go run main.go`

### Available methods

```go
IsV7() bool
GetExtractor() Extractor
SetOGameCredentials(username, password string)
SetProxy(proxyAddress, username, password, proxyType string, loginOnly bool) error
SetLoginWrapper(func(func() error) error)
GetClient() *OGameClient
Enable()
Disable()
IsEnabled() bool
Quiet(bool)
GetTasks() TasksOverview
Tx(clb func(tx *Prioritize) error) error
Begin() *Prioritize
BeginNamed(name string) *Prioritize
WithPriority(priority int) *Prioritize
GetPublicIP() (string, error)
OnStateChange(clb func(locked bool, actor string))
GetState() (bool, string)
IsLocked() bool
GetSession() string
AddAccount(number int, lang string) (NewAccount, error)
GetServer() Server
GetServerData() ServerData
SetUserAgent(newUserAgent string)
ServerURL() string
GetLanguage() string
GetPageContent(url.Values) ([]byte, error)
GetAlliancePageContent(url.Values) ([]byte, error)
PostPageContent(url.Values, url.Values) ([]byte, error)
LoginWithExistingCookies() (bool, error)
Login() error
Logout()
IsLoggedIn() bool
IsConnected() bool
GetUsername() string
GetUniverseName() string
GetUniverseSpeed() int64
GetUniverseSpeedFleet() int64
GetResearchSpeed() int64
GetNbSystems() int64
IsDonutGalaxy() bool
IsDonutSystem() bool
FleetDeutSaveFactor() float64
ServerVersion() string
ServerTime() time.Time
Location() *time.Location
IsUnderAttack() (bool, error)
GetUserInfos() UserInfos
SendMessage(playerID int64, message string) error
SendMessageAlliance(associationID int64, message string) error
ReconnectChat() bool
GetFleets(...Option) ([]Fleet, Slots)
GetFleetsFromEventList() []Fleet
CancelFleet(FleetID) error
GetAttacks() ([]AttackEvent, error)
GetAttacksUsing(CelestialID) ([]AttackEvent, error)
GalaxyInfos(galaxy, system int64, opts ...Option) (SystemInfos, error)
GetCachedResearch() Researches
GetResearch() Researches
GetCachedPlanets() []Planet
GetCachedMoons() []Moon
GetCachedCelestials() []Celestial
GetCachedCelestial(interface{}) Celestial
GetCachedPlayer() UserInfos
GetCachedPreferences() Preferences
IsVacationModeEnabled() bool
GetPlanets() []Planet
GetPlanet(interface{}) (Planet, error)
GetMoons() []Moon
GetMoon(interface{}) (Moon, error)
GetCelestial(interface{}) (Celestial, error)
GetCelestials() ([]Celestial, error)
Abandon(interface{}) error
GetEspionageReportMessages() ([]EspionageReportSummary, error)
GetEspionageReportFor(Coordinate) (EspionageReport, error)
GetEspionageReport(msgID int64) (EspionageReport, error)
GetCombatReportSummaryFor(Coordinate) (CombatReportSummary, error)
DeleteMessage(msgID int64) error
DeleteAllMessagesFromTab(tabID int64) error
Distance(origin, destination Coordinate) int64
FlightTime(origin, destination Coordinate, speed Speed, ships ShipsInfos) (secs, fuel int64)
RegisterWSCallback(string, func([]byte))
RemoveWSCallback(string)
RegisterChatCallback(func(ChatMsg))
RegisterAuctioneerCallback(func([]byte))
RegisterHTMLInterceptor(func(method, url string, params, payload url.Values, pageHTML []byte))
GetSlots() Slots
BuyOfferOfTheDay() error
BytesDownloaded() int64
BytesUploaded() int64
CreateUnion(fleet Fleet, unionUsers []string) (int64, error)
GetEmpire(nbr int64) (interface{}, error)
HeadersForPage(url string) (http.Header, error)
CharacterClass() CharacterClass
GetAuction() (Auction, error)
DoAuction(bid map[CelestialID]Resources) error
Highscore(category, typ, page int64) (Highscore, error)
GetAllResources() (map[CelestialID]Resources, error)
GetDMCosts(CelestialID) (DMCosts, error)
UseDM(string, CelestialID) error
GetItems(CelestialID) ([]Item, error)
ActivateItem(string, CelestialID) error

// Planet or Moon functions
GetResources(CelestialID) (Resources, error)
GetResourcesDetails(CelestialID) (ResourcesDetails, error)
SendFleet(celestialID CelestialID, ships []Quantifiable, speed Speed, where Coordinate, mission MissionID, resources Resources, expeditiontime, unionID int64) (Fleet, error)
EnsureFleet(celestialID CelestialID, ships []Quantifiable, speed Speed, where Coordinate, mission MissionID, resources Resources, expeditiontime, unionID int64) (Fleet, error)
Build(celestialID CelestialID, id ID, nbr int64) error
BuildCancelable(CelestialID, ID) error
BuildProduction(celestialID CelestialID, id ID, nbr int64) error
BuildBuilding(celestialID CelestialID, buildingID ID) error
BuildDefense(celestialID CelestialID, defenseID ID, nbr int64) error
BuildShips(celestialID CelestialID, shipID ID, nbr int64) error
CancelBuilding(CelestialID) error
TearDown(celestialID CelestialID, id ID) error
ConstructionsBeingBuilt(CelestialID) (buildingID ID, buildingCountdown int64, researchID ID, researchCountdown int64)
GetProduction(CelestialID) ([]Quantifiable, int64, error)
GetFacilities(CelestialID) (Facilities, error)
GetDefense(CelestialID) (DefensesInfos, error)
GetShips(CelestialID) (ShipsInfos, error)
GetResourcesBuildings(CelestialID) (ResourcesBuildings, error)
CancelResearch(CelestialID) error
BuildTechnology(celestialID CelestialID, technologyID ID) error

// Planet specific functions
GetResourceSettings(PlanetID) (ResourceSettings, error)
SetResourceSettings(PlanetID, ResourceSettings) error
SendIPM(PlanetID, Coordinate, int64, ID) (int64, error)
//GetResourcesProductionRatio(PlanetID) (float64, error)
GetResourcesProductions(PlanetID) (Resources, error)
GetResourcesProductionsLight(ResourcesBuildings, Researches, ResourceSettings, Temperature) Resources

// Moon specific functions
Phalanx(MoonID, Coordinate) ([]Fleet, error)
UnsafePhalanx(MoonID, Coordinate) ([]Fleet, error)
JumpGate(origin, dest MoonID, ships ShipsInfos) (bool, int64, error)
JumpGateDestinations(origin MoonID) ([]MoonID, int64, error)
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
POST /bot/delete-report/:messageID
POST /bot/delete-all-espionage-reports
POST /bot/delete-all-reports/:tabIndex
GET  /bot/attacks
GET  /bot/galaxy-infos/:galaxy/:system
GET  /bot/get-research
GET  /bot/price/:ogameID/:nbr
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
POST /bot/planets/:planetID/send-ipm
POST /bot/planets/:planetID/teardown/:ogameID
GET  /bot/moons/:moonID/phalanx/:galaxy/:system/:position
GET  /bot/get-auction
POST /bot/do-auction
```

# docker container

If you have Docker, and you are looking for a docker image just update the `.env` file specifying the universe name, credentials and language.

```properties
OGAMED_HOST=0.0.0.0
OGAMED_UNIVERSE=Bellatrix
OGAMED_USERNAME=email@gmail.com
OGAMED_PASSWORD=*****
OGAMED_LANGUAGE=en
```

`.env` file contains all the environment variables used by `ogamed` service.

Create and run the container:

```shell script
docker-compose up -d
```

Check the container logs to see if the service is correctly connected to the server:

```shell script
docker logs ogame
``` 

The container can be invoked as usual [as a service](#ogamed-service) or [as a library](#ogame-library).